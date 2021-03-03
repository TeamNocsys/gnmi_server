package cmd

import (
    "errors"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/utils"
)

type OperType int32

const (
    ADD     OperType = 0
    UPDATE  OperType = 1
    DEL     OperType = 2
)

type Adapter struct {
    client command.Client
}

func (adpt *Adapter) exec(cmdstr string) error {
    logrus.Trace(cmdstr + "|EXEC")
    if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
        return errors.New(r)
    }
    return nil
}
