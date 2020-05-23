package redis

import (
	"errors"
	"fmt"
	"hekv/kv"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	Config Config
	kv     kv.KV
}

func CreateServer() *Server {
	//tmpPath, err := ioutil.TempDir("hekv_tmp")
	tempDir := os.TempDir()
	hekvTempPath, err := ioutil.TempDir(tempDir, "hekv")
	if err != nil {
		return nil
	}
	kv, err := kv.OpenPebbleKV(hekvTempPath)
	if err != nil {
		return nil
	}
	server := &Server{
		Config: Config{
			Port: 6380,
		},
		kv: kv,
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

func (s *Server) handleCommand(conn *connection, command Command) error {
	cmd := command.GetArg(0)
	if strings.EqualFold(cmd, "GET") {
		arg := command.GetArg(1)
		var value []byte
		err := s.kv.Get([]byte(arg), func(val []byte) error {
			value = val
			return nil
		})

		if err != nil {
		}

		rs := Replys{}
		if value != nil {
			rs = append(rs, BulkReply{
				val: string(value),
			})
		}

		return conn.WriteAndFlush(rs)
	} else if strings.EqualFold(cmd, "SET") {
		key := command.GetArg(1)
		val := command.GetArg(2)
		err := s.kv.Put([]byte(key), []byte(val))
		if err != nil {

		}
		rs := Replys{StatusReply{"OK"}}
		return conn.WriteAndFlush(rs)
	} else if strings.EqualFold(cmd, "DEL") {
		key := command.GetArg(1)
		err := s.kv.Del([]byte(key))
		if err != nil {

		}
		rs := Replys{IntegerReply{1}}
		return conn.WriteAndFlush(rs)
	}

	return errors.New("Unknow command")
}
