package redis

import (
	"strings"
)

type Handler struct {
	mapper map[string]func(s *Server, conn *Connection, command Command) error
}

func NewHandler() *Handler {
	handleMapper := make(map[string]func(s *Server, conn *Connection, command Command) error)
	return &Handler{
		mapper: handleMapper,
	}
}

func (h *Handler) Init() {
	// string
	h.mapper["GET"] = GET
	h.mapper["SET"] = SET
	h.mapper["DEL"] = DEL

	// hash
	h.mapper["HSET"] = HSET
	h.mapper["HGET"] = HGET
	h.mapper["HMSET"] = HMSET

}

func (h *Handler) Func(name string) func(s *Server, conn *Connection, command Command) error {
	cmdName := strings.ToUpper(name)
	return h.mapper[cmdName]
}
