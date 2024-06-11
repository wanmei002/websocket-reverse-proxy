package unixs

import "net"

func UnixListener() {
	net.Listen("unix")
}
