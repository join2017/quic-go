package handshake

import (
	"net"

	"github.com/bifurcation/mint"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

type cookieHandler struct {
	callback func(net.Addr, *Cookie) bool

	cookieGenerator *CookieGenerator
}

var _ mint.CookieHandler = &cookieHandler{}

func newCookieHandler(callback func(net.Addr, *Cookie) bool) (*cookieHandler, error) {
	cookieGenerator, err := NewCookieGenerator()
	if err != nil {
		return nil, err
	}
	return &cookieHandler{
		callback:        callback,
		cookieGenerator: cookieGenerator,
	}, nil
}

func (h *cookieHandler) Generate(conn *mint.Conn) (*mint.CookieExtension, error) {
	data, err := h.cookieGenerator.NewToken(conn.RemoteAddr())
	if err != nil {
		return nil, err
	}
	return &mint.CookieExtension{Cookie: data}, nil
}

func (h *cookieHandler) Validate(conn *mint.Conn, cookie *mint.CookieExtension) bool {
	data, err := h.cookieGenerator.DecodeToken(cookie.Cookie)
	if err != nil {
		utils.Debugf("Couldn't decode cookie: %s", err.Error())
		return false
	}
	return h.callback(conn.RemoteAddr(), data)
}
