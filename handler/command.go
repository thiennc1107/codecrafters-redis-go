package handler

import (
	"errors"
	"fmt"
)

type Command struct {
	Name      string
	Args      []string
	ArgLength int
}

type CommandHandler interface {
	Execute(ctx Context) (string, error)
}

func GetCommandHandler(commandName string) (CommandHandler, error) {
	switch commandName {
	case EchoCommand:
		return &EchoCommandHandler{}, nil
	case PingCommand:
		return &PingCommandHandler{}, nil
	default:
		return nil, fmt.Errorf("command not found")
	}
}

const EchoCommand = "ECHO"
const PingCommand = "ping"

type EchoCommandHandler struct {
}

const ErrEchoArgumentCount = "echo command require exactly 1 argument"

func (e *EchoCommandHandler) Execute(ctx Context) (string, error) {
	cmd := ctx.Command
	if cmd.ArgLength != 1 {
		return "", errors.New(ErrEchoArgumentCount)
	}
	return cmd.Args[0], nil
}

type PingCommandHandler struct {
}

func (p *PingCommandHandler) Execute(ctx Context) (string, error) {
	return "PONG", nil
}
