package redis

type Command struct {
	cmds []string
}

func (c *Command) AddArgs(cmd string) {
	c.cmds = append(c.cmds, cmd)
}

func (c *Command) GetArg(i int) string {
	return c.cmds[i]
}

type Values struct {
	values []string
}
