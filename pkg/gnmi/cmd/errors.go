package cmd

import "fmt"

var (
    ErrUnknown = fmt.Errorf("unknown")
    ErrInvalidOperType = fmt.Errorf("operation type not supported")
    ErrTypeConversion = fmt.Errorf("type conversion failed")
)