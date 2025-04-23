package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
		}).Error("main: Redirect server error")
		os.Exit(1)
	}
}

func runServer() (err error) {
	webPort, err := strconv.Atoi(os.Getenv("WEB_PORT"))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrapf(err, "main: Failed to parse web port"),
		}
		return
	}

	logger.WithFields(logger.Fields{
		"port":     80,
		"web_port": webPort,
	}).Info("main: Starting HTTP server")

	server := &http.Server{
		Addr:         "0.0.0.0:80",
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,

		Handler: http.HandlerFunc(func(
			w http.ResponseWriter, req *http.Request) {

			if strings.HasPrefix(req.URL.Path,
				"/.well-known/acme-challenge/") {

				pathSplit := strings.Split(req.URL.Path, "/")
				token := pathSplit[len(pathSplit)-1]

				challenge := GetChallenge(token)
				if challenge == "" {
					w.WriteHeader(404)
					fmt.Fprint(w, "404 page not found")
					return
				}

				w.WriteHeader(200)
				fmt.Fprint(w, challenge)
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

	err = server.ListenAndServe()
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "main: Failed to bind web server"),
		}
		return
	}

	return
}
