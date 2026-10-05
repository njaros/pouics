package command

type Command struct {
	from string
	command string
	args any
}

func NewCommand(from string, command string, args any) Command {
	return Command{from, command, args}
}

func ErrorCommand(from string) Command {
	return NewCommand(from, "", nil)
}