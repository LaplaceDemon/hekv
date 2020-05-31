package redis

func HSET(s *Server, conn *Connection, command Command) error {
	size := command.Size()
	if size != 4 {
		// error
	}
	key := command.GetArg(1)
	field := command.GetArg(2)
	value := command.GetArg(3)

	k := make([]byte, 0, len(key)+1+len(field))
	k = append(k, key...)
	k = append(k, byte(0))
	k = append(k, field...)
	err := s.KV.Put(k, []byte(value))
	if err != nil {

	}
	rs := Replys{StatusReply{"OK"}}
	return conn.WriteAndFlush(rs)
}

func HGET(s *Server, conn *Connection, command Command) error {
	size := command.Size()
	if size%2 != 0 {
		// error
	}
	key := command.GetArg(1)
	field := command.GetArg(2)

	k := make([]byte, 0, len(key)+1+len(field))
	k = append(k, key...)
	k = append(k, byte(0))
	k = append(k, field...)

	var value []byte
	err := s.KV.Get(k, func(val []byte) error {
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

func HMSET(s *Server, conn *Connection, command Command) error {
	size := command.Size()
	if size%2 != 0 {
		// error
	}
	key := command.GetArg(1)

	i := 2
	for {
		if i < size {
			break
		}
		field := command.GetArg(i)
		i++
		value := command.GetArg(i)
		i++

		k := make([]byte, 0, len(key)+1+len(field))
		k = append(k, key...)
		k = append(k, byte(0))
		k = append(k, field...)
		err := s.KV.Put(k, []byte(value))
		if err != nil {

		}
	}
	rs := Replys{StatusReply{"OK"}}
	return conn.WriteAndFlush(rs)
}
