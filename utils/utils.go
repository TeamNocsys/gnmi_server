package utils

import (
	"bytes"
	"os/exec"
)

func Utils_execute_cmd(name string, arg ...string) (error, string){
	cmd := exec.Command(name, arg...)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	err := cmd.Run()
	return err, buffer.String()
}
