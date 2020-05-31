package redis

func GET(s *Server, conn *Connection, command Command) error {
	arg := command.GetArg(1)
	var value []byte
	err := s.KV.Get([]byte(arg), func(val []byte) error {
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
}

func SET(s *Server, conn *Connection, command Command) error {
	key := command.GetArg(1)
	val := command.GetArg(2)
	err := s.KV.Put([]byte(key), []byte(val))
	if err != nil {

	}
	rs := Replys{StatusReply{"OK"}}
	return conn.WriteAndFlush(rs)
}

func DEL(s *Server, conn *Connection, command Command) error {
	key := command.GetArg(1)
	err := s.KV.Del([]byte(key))
	if err != nil {

	}
	rs := Replys{IntegerReply{1}}
	return conn.WriteAndFlush(rs)
}
