package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pritunl/pritunl-zero/redirect/crypto"
	"github.com/pritunl/tools/errors"
	"github.com/pritunl/tools/errortypes"
	"github.com/pritunl/tools/logger"
)

func main() {
	logger.Init()
	logger.AddHandler(func(record *logger.Record) {
		fmt.Print(record.String())
	})

	err := runServer()
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err,
		}).Error("redirect: Redirect server error")
		os.Exit(1)
	}
}

func runServer() (err error) {
	webPort, err := strconv.Atoi(os.Getenv("WEB_PORT"))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrapf(err, "redirect: Failed to parse web port"),
		}
		return
	}

	key := crypto.AsymNaclHmacKey{
		PublicKey: os.Getenv("PUBLIC_KEY"),
		Key:       os.Getenv("KEY"),
		Secret:    os.Getenv("SECRET"),
	}

	box := &crypto.AsymNaclHmac{}

	err = box.Import(key)
	if err != nil {
		return
	}

	logger.WithFields(logger.Fields{
		"port":     80,
		"web_port": webPort,
	}).Info("redirect: Starting HTTP redirect server")

	go sandboxTest()

	file := os.NewFile(uintptr(3), "systemd-socket")
	listener, err := net.FileListener(file)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "redirect: Failed to get socket listener"),
		}
		return
	}

	server := &http.Server{
		Addr:         ":80",
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,

		Handler: http.HandlerFunc(func(
			w http.ResponseWriter, req *http.Request) {

			if req.Method == "GET" && strings.HasPrefix(req.URL.Path,
				"/.well-known/acme-challenge/") {

				pathSplit := strings.Split(req.URL.Path, "/")
				token := pathSplit[len(pathSplit)-1]

				chal := GetChallenge(token)
				if chal == nil {
					w.WriteHeader(404)
					fmt.Fprint(w, "404 page not found")
					return
				}

				w.WriteHeader(200)
				fmt.Fprint(w, chal.Response)
				return
			} else if req.Method == "POST" && req.URL.Path == "/token" {
				bodyBytes := make([]byte, 8096)
				n, err := io.LimitReader(req.Body, 8096).Read(bodyBytes)
				if err != nil && err != io.EOF {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "Internal server error")
					return
				}
				bodyBytes = bodyBytes[:n]

				chal := &Challenge{}
				err = box.UnsealJson(string(bodyBytes), chal)
				if err != nil && err != io.EOF {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprint(w, "Failed to authorize")
					return
				}

				AddChallenge(chal)

				w.WriteHeader(200)
				fmt.Fprint(w, "success")
				return
			}

			req.URL.Scheme = "https"
			req.URL.Host = StripPort(req.Host)
			if webPort != 443 {
				req.URL.Host += fmt.Sprintf(":%d", webPort)
			}

			http.Redirect(w, req, req.URL.String(),
				http.StatusMovedPermanently)
		}),
	}

	err = server.Serve(listener)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "redirect: Failed to bind web server"),
		}
		return
	}

	return
}

func sandboxTest() {
	time.Sleep(3 * time.Second)

	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 3 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest(
		"GET",
		"https://127.0.0.1",
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "acme: Sandbox test request failed"),
		}
		return
	}

	resp, err := client.Do(req)
	if err == nil {
		logger.WithFields(logger.Fields{
			"status_code": resp.StatusCode,
		}).Error("redirect: Sandbox escape test failed")
	} else {
		logger.Info("redirect: Sandbox escape test successful")
	}
}
