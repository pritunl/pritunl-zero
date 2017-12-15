package keybase

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var saltpackReg = regexp.MustCompile("[^a-zA-Z0-9]+")

func FilterSig(input string) string {
	input = strings.TrimSpace(input)
	input = strings.Replace(input,
		"BEGIN KEYBASE SALTPACK SIGNED MESSAGE", "", 1)
	input = strings.Replace(input,
		"END KEYBASE SALTPACK SIGNED MESSAGE", "", 1)
	input = strings.Replace(input, " ", "", -1)
	input = strings.Replace(input, ".", "", -1)

	input = saltpackReg.ReplaceAllString(input, "")

	return "BEGIN KEYBASE SALTPACK SIGNED MESSAGE. " + input +
		". END KEYBASE SALTPACK SIGNED MESSAGE."
}

func VerifySig(message, signature, username string) (valid bool, err error) {
	signature = FilterSig(signature)

	usr, err := user.Lookup("nobody")
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrapf(err, "keybase: Failed to lookup user"),
		}
		return
	}

	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrapf(err, "keybase: Failed to parse uid"),
		}
		return
	}

	gid, err := strconv.Atoi(usr.Gid)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrapf(err, "keybase: Failed to parse gid"),
		}
		return
	}

	cmd := exec.Command(
		"pritunl-keybase",
		"verify",
		"--signed-by",
		username,
		"--message",
		signature,
	)

	env := []string{}
	env = append(env, "TERM="+os.Getenv("TERM"))
	env = append(env, "PATH="+os.Getenv("PATH"))
	env = append(env, "LANG="+os.Getenv("LANG"))
	env = append(env, "SHELL="+os.Getenv("SHELL"))
	env = append(env, "HOME=/tmp")
	env = append(env, "KEYBASE_STANDALONE=true")
	env = append(env, "KEYBASE_CONFIG_FILE=/dev/null")
	env = append(env, "KEYBASE_UPDATER_CONFIG_FILE=/dev/null")
	cmd.Env = env

	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(uid),
		Gid: uint32(gid),
	}

	outputByt, err := cmd.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			eOutput := strings.ToLower(string(e.Stderr))

			if strings.Contains(eOutput, "local assertions for user") ||
				strings.Contains(eOutput, "bad signature") ||
				strings.Contains(eOutput, "invalid encoding length") ||
				strings.Contains(eOutput, "unexpected eof") ||
				strings.Contains(eOutput, "unknown stream format") {

				err = nil
			} else {
				err = &errortypes.UnknownError{
					errors.Wrapf(err, "keybase: Failed to exec keybase %s",
						eOutput),
				}
			}
		} else {
			err = &errortypes.UnknownError{
				errors.Wrap(err, "keybase: Failed to exec keybase"),
			}
		}

		return
	}
	output := string(outputByt)
	output = strings.TrimSpace(output)

	valid = output == message

	return
}
