package redis

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	"hekv/store"
)

type Server struct {
	Config  Config
	handler *Handler
	KV      store.KV
}

func CreateServer() *Server {
	//tmpPath, err := ioutil.TempDir("hekv_tmp")
	tempDir := os.TempDir()
	hekvTempPath, err := ioutil.TempDir(tempDir, "hekv")
	if err != nil {
		return nil
	}
	kv, err := store.OpenPebbleKV(hekvTempPath)
	if err != nil {
		return nil
	}
	server := &Server{
		Config: Config{
			Port: 6380,
		},
		handler: NewHandler(),
		KV:      kv,
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

func (s *Server) handleConn(conn *Connection) error {
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

func (s *Server) handleCommand(conn *Connection, command Command) error {
	arg0 := command.GetArg(0)
	mapper := s.handler.Get(arg0)

	if mapper == nil {
		return errors.New("Unknow command")
	}

	return mapper(s, conn, command)
}
