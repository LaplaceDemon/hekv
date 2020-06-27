package redis

import (
	"bufio"
	"fmt"
	"net"
)

type Connection struct {
	c      net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewConnection(conn net.Conn) *Connection {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	return &Connection{
		c:      conn,
		reader: reader,
		writer: writer,
	}
}

func (c *Connection) ReadLine() (string, error) {
	line, _, err := c.reader.ReadLine()
	if err != nil {
		return "", err
	}
	str := string(line)
	return str, nil
}

func (c *Connection) WriteAndFlush(rs Replys) error {
	if len(rs) == 0 {
		if _, err := c.writer.WriteString("$-1\r\n"); err != nil {
			return err
		}
	} else {
		for _, r := range rs {
			switch r.(type) {
			case BulkReply:
				val := r.value()
				valLen := len(val)
				if _, err := c.writer.WriteString(fmt.Sprintf("$%d\r\n", valLen)); err != nil {
					return err
				}
				if _, err := c.writer.WriteString(fmt.Sprintf("%s\r\n", val)); err != nil {
					return err
				}
			case StatusReply:
				val := r.value()
				if _, err := c.writer.WriteString(fmt.Sprintf("+%s\r\n", val)); err != nil {
					return err
				}
			case IntegerReply:
				val := r.value()
				if _, err := c.writer.WriteString(fmt.Sprintf(":%s\r\n", val)); err != nil {
					return err
				}
			default:

			}

		}
	}
	return c.writer.Flush()
}
