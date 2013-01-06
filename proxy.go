package main

import (
//   "fmt"
//   "io"
   "net"
   "os"
   "log"
   "sync"
)


var stdout = log.New(os.Stdout, "[Out]   ", log.LstdFlags)
var debug  = log.New(os.Stderr, "[Debug] ", log.LstdFlags)
var stderr = log.New(os.Stderr, "[Error] ", log.LstdFlags)
type ConnectionCount struct {
   sync.Mutex
   Count int
}

const (
   MAX_CONNECTIONS = 10
)
type Proxy struct {
   listener *net.TCPListener
   listen_type string
   listen_addr *net.TCPAddr
   connections []*Connection
   
   // More socks related things
   auth_header *Header
}

// Create a New Socks proxy, using the specified net (tcp4, tcp6, tcp),
// address, and port.
func NewProxy(nettype, addr string) (*Proxy, error) {
   sp := new(Proxy)
   sp.listen_type = nettype

   resolved_addr, err := net.ResolveTCPAddr(sp.listen_type, addr)
   if err != nil {
      return nil, err
   }
   sp.listen_addr = resolved_addr

   listener, err := net.ListenTCP(sp.listen_type, sp.listen_addr)
   if err != nil {
      return nil, err
   }
   sp.listener = listener
   sp.auth_header = nil
   sp.connections = make([]*Connection, MAX_CONNECTIONS)
   stdout.Printf("Listening on %s port %s\n", sp.listen_addr.IP.String(), (int)(sp.listen_addr.Port))
   return sp, nil
}
func (sp *Proxy) handleTCPConnection(c *net.TCPConn) error {
   return nil
}
func (sp *Proxy) Listen() error {
   connection_count := new(ConnectionCount)
   connection_count.Lock()
   connection_count.Count = 0
   connection_count.Unlock()

   for {
      tcp_conn, err := sp.listener.AcceptTCP()
      if err != nil {
         panic(err)
      }
      connection_count.Lock()
      connection_count.Count++
      connection_count.Unlock()

      socks_conn := NewConnection(tcp_conn)
      go socks_conn.Handle(connection_count)
      /*if err != nil {
         stderr.Panicf("NewConnection panicked upon receiving %+v (err %s)\n", socks_conn, err)
      }
      debug.Println("Blocking on incoming connection")
      sp.handleTCPConnection(tcp_conn)
      debug.Println("Unblocked!")*/
      if connection_count.Count == MAX_CONNECTIONS {
         connection_count.Lock()
         debug.Printf("%d connections are currently in use, where %d is the max\n", connection_count.Count, MAX_CONNECTIONS)
         debug.Println("Waiting for connection to open up")
         connection_count.Count--
         connection_count.Unlock()
      }
   }
   return nil
}


func main() {
   proxy, err := NewProxy("tcp4", "0.0.0.0:1080")
   if err != nil {
      panic(err)
   }
   proxy.Listen()
}
