package shell

import (
	"io/ioutil"
	"os/exec"
)

func Pipe(name string, args ...string) (error, string) {
	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err, ""
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err, ""
	}

	if err := cmd.Start(); err != nil {
		return err, ""
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return err, ""
	}

	if len(bytesErr) != 0 {
		return err, ""
	}

	byteOut, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err, ""
	}

	if err := cmd.Wait(); err != nil {
		return err, ""
	}

	return nil, string(byteOut)
}
