package bastion

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/ssh"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Bastion struct {
	Authority bson.ObjectId
	Container string
	authr     *authority.Authority
	state     bool
	kill      bool
	path      string
}

func (b *Bastion) wait() {
	defer func() {
		os.Remove(b.path)
		b.state = false
		delete(state, b.Authority)
	}()

	output, err := utils.ExecOutput("", "docker", "wait", b.Container)
	if b.state && err != nil {
		err = &errortypes.RequestError{
			errors.Wrapf(err, "utils: Failed to exec docker"),
		}
		return
	}

	output = strings.TrimSpace(output)
	if output != "0" && output != "137" {
		logrus.WithFields(logrus.Fields{
			"exit_code": output,
		}).Error("bastion: Bastion process error")
	}
}

func (b *Bastion) Start(db *database.Database,
	authr *authority.Authority) (err error) {

	logrus.WithFields(logrus.Fields{
		"authority_id": b.Authority.Hex(),
	}).Info("bastion: Starting bastion server")

	if b.state || b.path != "" {
		err = &errortypes.UnknownError{
			errors.Wrapf(err, "bastion: Bastion server already running"),
		}
		return
	}

	b.authr = authr
	b.path = utils.GetTempPath()

	if authr.ProxyPublicKey == "" || authr.ProxyPrivateKey == "" {
		err = authr.GenerateRsaProxyPrivateKey()
		if err != nil {
			return
		}

		err = authr.CommitFields(
			db,
			set.NewSet("proxy_private_key", "proxy_public_key"),
		)
		if err != nil {
			return
		}
	}

	b.state = true

	err = utils.ExistsMkdir(b.path, 0755)
	if err != nil {
		b.state = false
		os.Remove(b.path)
		return
	}

	if authr.HostCertificates {
		cert, e := ssh.NewBastionHostCertificate(db,
			authr.ProxyHostname, authr.ProxyPublicKey, authr)
		if e != nil {
			b.state = false
			os.Remove(b.path)
			err = e
			return
		}

		hostCertPath := filepath.Join(b.path, "ssh_host_rsa_key-cert.pub")

		if len(cert.Certificates) == 0 {
			b.state = false
			os.Remove(b.path)
			err = &errortypes.UnknownError{
				errors.Wrapf(err, "bastion: Missing host certificate"),
			}
			return
		}

		err = utils.CreateWrite(hostCertPath, cert.Certificates[0], 0644)
		if err != nil {
			b.state = false
			os.Remove(b.path)
			return
		}
	}

	output, err := utils.ExecOutput("",
		"docker",
		"run",
		"--rm",
		"-d",
		"-u", "bastion",
		"--name", DockerGetName(authr.Id),
		"-p", fmt.Sprintf("%d:9722", authr.ProxyPort),
		"-v", fmt.Sprintf("%s:/ssh_mount", b.path),
		"-e", fmt.Sprintf(
			"BASTION_TRUSTED=%s", authr.PublicKey),
		"-e", fmt.Sprintf(
			"BASTION_HOST_KEY=%s", authr.ProxyPrivateKey),
		"-e", fmt.Sprintf(
			"BASTION_HOST_PUB_KEY=%s", authr.ProxyPublicKey),
		settings.System.BastionDockerImage,
	)
	if err != nil {
		b.state = false
		os.Remove(b.path)

		return
	}

	b.Container = strings.TrimSpace(output)

	go b.wait()

	return
}

func (b *Bastion) Stop() (err error) {
	if b.kill {
		return
	}
	b.kill = true

	logrus.WithFields(logrus.Fields{
		"authority_id": b.Authority.Hex(),
	}).Info("bastion: Stopping bastion server")

	_, err = utils.ExecOutputLogged(nil,
		"docker", "stop", "-t", "3", b.Container)
	if err != nil {
		return
	}

	go func() {
		time.Sleep(15 * time.Second)

		if b.state {
			utils.ExecOutputLogged(nil, "docker", "kill", b.Container)
		}
	}()

	return
}

func (b *Bastion) State() bool {
	return b.state
}

func (b *Bastion) Diff(authr *authority.Authority) bool {
	if b.authr.ProxyPublicKey != authr.ProxyPublicKey ||
		b.authr.ProxyPrivateKey != authr.ProxyPrivateKey ||
		b.authr.HostCertificates != authr.HostCertificates ||
		b.authr.ProxyPort != authr.ProxyPort ||
		b.authr.PublicKey != authr.PublicKey {

		return true
	}
	return false
}
