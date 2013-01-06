package main

import "net"

type SocksConnection struct {
   conn *net.TCPConn
   auth_header *Header
}

