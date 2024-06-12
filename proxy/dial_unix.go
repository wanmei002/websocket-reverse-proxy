package proxy

import (
	"github.com/wanmei002/websocket-reverse-proxy/unixs"
	"net"
)

func DialUnix() (net.Conn, error) {
	return net.Dial("unix", unixs.SocketFile)

}
