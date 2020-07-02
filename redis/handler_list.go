package redis

func LPUSH(s *Server, conn *Connection, command Command) error {
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
