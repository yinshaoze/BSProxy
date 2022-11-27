package service

import (
	"fmt"
	"log"
	"net"

	"github.com/yinshaoze/BSProxy/config"
	"github.com/yinshaoze/BSProxy/service/minecraft"
	"github.com/yinshaoze/BSProxy/service/tls"
	"github.com/yinshaoze/BSProxy/service/transfer"
)

func newConnReceiver(s *config.ConfigProxyService,
	conn *net.TCPConn,
	options *transfer.Options,
) {
	log.Println("Service", s.Name, ": A new connection request sent by", conn.RemoteAddr().String(), "is received.")
	defer log.Println("Service", s.Name, ": A connection with", conn.RemoteAddr().String(), "is closed.")
	var err error // in order to avoid scoop problems
	var remote net.Conn = nil

	if options.IsTLSHandleNeeded {
		remote, err = tls.NewConnHandler(s, conn, options.Out)
		if err != nil {
			conn.Close()
			return
		}
	} else if options.IsMinecraftHandleNeeded {
		remote, err = minecraft.NewConnHandler(s, conn, options)
		if err != nil {
			conn.Close()
			return
		}
	}

	if remote == nil {
		remote, err = options.Out.Dial("tcp", fmt.Sprintf("%v:%v", s.TargetAddress, s.TargetPort))
		if err != nil {
			log.Printf("Service %s: Failed to dial to target server: %v", s.Name, err.Error())
			conn.Close()
			return
		}
	}
	options.AddCount(1)
	defer options.AddCount(-1)
	transfer.SimpleTransfer(conn, remote, options.FlowType)
}
