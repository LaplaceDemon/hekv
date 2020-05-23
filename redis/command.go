package redis

import "strconv"

type Command struct {
	cmds []string
}

func (c *Command) AddArgs(cmd string) {
	c.cmds = append(c.cmds, cmd)
}

func (c *Command) GetArg(i int) string {
	return c.cmds[i]
}

type Replys []Reply

type Reply interface {
	value() string
}

type BulkReply struct {
	val string
}

func (r BulkReply) value() string {
	return r.val
}

type StatusReply struct {
	val string
}

func (r StatusReply) value() string {
	return r.val
}

type IntegerReply struct {
	val int
}

func (r IntegerReply) value() string {
	return strconv.Itoa(r.val)
}
