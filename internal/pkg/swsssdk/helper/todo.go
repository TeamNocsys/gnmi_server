package helper

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "gnmi_server/cmd/command"
)

type Todo struct {
    Key string
    Client command.Client
    Data *sonicpb.SonicTodo_Todo_TodoList
}
