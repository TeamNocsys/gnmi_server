package utils

import (
    "bytes"
    "context"
    "os/exec"
    "time"
)

func Utils_execute_cmd(name string, arg ...string) (error, string) {
    // TODO: 后期根据不同命令，传入不同超时时间
    ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, name, arg...)
    var buffer bytes.Buffer
    cmd.Stdout = &buffer
    err := cmd.Run()
    return err, buffer.String()
}
