package command

import (
	"encoding/json/v2"
	"fmt"
)

// Parser takes the raw binaries from a sender
// and create a Command from it if possible.
// It needs a mapping of existing command and a dto for
// each existing command to know how to parse.
type Parser struct {
	commandMapping map[string]any
}

func NewParser(commandMapping map[string]any) *Parser {
	return &Parser{commandMapping}
}

func (p *Parser)Parse(b []byte, senderId string) (Command, error) {
	cmd, err := p.getCommand(b)
	if err != nil {
		return ErrorCommand(senderId), err
	}

	args, ok := p.commandMapping[cmd]
	if !ok {
		return	ErrorCommand(senderId),
				fmt.Errorf("%s: %s", cmd, "command not found")
	}

	if err := p.getArgs(b, &args); err != nil {
		return ErrorCommand(senderId), err
	}

	return NewCommand(senderId, cmd, args), nil
}

func (p *Parser)getCommand(b []byte) (string, error) {
	var command struct {
		Cmd string `json:"cmd"`
	}
	if err := json.Unmarshal(b, &command); err != nil {
		return "", err
	}
	return command.Cmd, nil
}

func (p * Parser)getArgs(b []byte, dest any) error {
	return json.Unmarshal(b, &dest)
}