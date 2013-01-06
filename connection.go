package main

import "net"

type Connection struct {
   conn *net.TCPConn
   auth_header *Header
}

