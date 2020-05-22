package redis

import (
	"bufio"
	"fmt"
	"net"
)

type connection struct {
	c      net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewConnection(conn net.Conn) *connection {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	return &connection{
		c:      conn,
		reader: reader,
		writer: writer,
	}
}

func (c *connection) ReadLine() (string, error) {
	line, _, err := c.reader.ReadLine()
	if err != nil {
		return "", err
	}
	fmt.Printf("%s\n", string(line))
	str := string(line)
	return str, nil
}

func (c *connection) WriteAndFlush(v Values) error {
	for _, val := range v.values {
		valLen := len(val)
		if _, err := c.writer.WriteString(fmt.Sprintf("$%d\r\n", valLen)); err != nil {
			return err
		}
		if _, err := c.writer.WriteString(fmt.Sprintf("%s\r\n", val)); err != nil {
			return err
		}
	}
	return c.writer.Flush()
}
