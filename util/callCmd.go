package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

func CallCmd(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	var outerr bytes.Buffer
	cmd.Stderr = &outerr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	if outerr.Len() != 0 {
		return "", fmt.Errorf(outerr.String())
	}

	return out.String(), nil
}
