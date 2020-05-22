package redis

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Server struct {
	Config Config
}

func CreateServer() *Server {
	server := &Server{
		Config: Config{
			Port: 6380,
		},
	}

	return server
}

func (s *Server) Run() error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Config.Port))
	if err != nil {
		return err
	}

	for {
		tcpConn, err := listen.Accept()
		if err != nil {
			return err
		}

		go func(tcpConn net.Conn) {
			defer tcpConn.Close()
			redisConnection := NewConnection(tcpConn)
			if err := s.handleConn(redisConnection); err != nil {
				return
			}
		}(tcpConn)
	}
}

func (s *Server) handleConn(conn *connection) error {
	for {
		//conn.Read
		argsCountStr, err := conn.ReadLine()
		if err != nil {
			return err
		}

		cmd := Command{}
		if strings.HasPrefix(argsCountStr, "*") {
			argsCount, err := strconv.Atoi(argsCountStr[1:])
			if err != nil {
				return err
			}
			for i := 0; i < argsCount; i++ {
				argLengthStr, err := conn.ReadLine()
				if err != nil {
					return err
				}

				if strings.HasPrefix(argLengthStr, "$") {
					var argValueStr string
					argLength, err := strconv.Atoi(argLengthStr[1:])
					if err != nil {
						return err
					}

					argValueStr, err = conn.ReadLine()
					if err != nil {
						return err
					}
					if argLength != len(argValueStr) {
						return errors.New("error protocol. argLength != len(line)")
					}
					cmd.AddArgs(argValueStr)
				}
			}
		}

		s.handleCommand(conn, cmd)
	}

	return nil
}

func (s *Server) handleCommand(conn *connection, command Command) {
	if command.GetArg(0) == "GET" {
		v := Values{}
		v.values = append(v.values, "helloworld")
		conn.WriteAndFlush(v)
	}
}
