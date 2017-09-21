package handshake

import (
	"net"

	"github.com/bifurcation/mint"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var cookieValid bool
var mockCallback = func(net.Addr, *Cookie) bool {
	return cookieValid
}

var _ = Describe("Cookie Handler", func() {
	var ch *cookieHandler
	var conn *mint.Conn

	BeforeEach(func() {
		var err error
		ch, err = newCookieHandler(mockCallback)
		Expect(err).ToNot(HaveOccurred())
		addr := &net.UDPAddr{IP: net.IPv4(42, 43, 44, 45), Port: 46}
		conn = mint.NewConn(&fakeConn{remoteAddr: addr}, &mint.Config{}, false)
	})

	It("generates and validates a token", func() {
		cookie, err := ch.Generate(conn)
		Expect(err).ToNot(HaveOccurred())
		Expect(ch.Validate(conn, cookie)).To(BeFalse())
		cookieValid = true
		Expect(ch.Validate(conn, cookie)).To(BeTrue())
	})

	It("correctly handles a token that it can't decode", func() {
		cookie := &mint.CookieExtension{Cookie: []byte("unparseable cookie")}
		Expect(ch.Validate(conn, cookie)).To(BeFalse())
	})
})
