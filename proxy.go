package main

import (
   "fmt"
   "io"
   "net"
)

type SocksProxy struct {
   listener *net.TCPListener
   listen_type string
   listen_addr *net.TCPAddr
   
   // More socks related things
   auth_header *Header
}

// Create a New Socks proxy, using the specified net (tcp4, tcp6, tcp),
// address, and port.
func NewSocksProxy(nettype, addr string) (*SocksProxy, error) {
   sp := new(SocksProxy)
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
   fmt.Printf("Listening on %s port %s\n", sp.listen_addr.IP.String(), (int)(sp.listen_addr.Port))
   return sp, nil
}
func (sp *SocksProxy) handleTCPConnection(c *net.TCPConn) error {
   fmt.Printf("handleTCPConnection received connection %+v\n", c)
   c.SetReadBuffer(1024)
   c.SetWriteBuffer(1024)
   defer c.Close()

   read_buffer := make([]byte, 1024)

   for {
      // Clear the buffer
      for i := 0; i < cap(read_buffer); i+=1 {
         read_buffer[i] = 0
      }
      count, err := c.Read(read_buffer)
      if err != io.EOF && err != nil {
         panic(err)
      } else if err == io.EOF {
         fmt.Printf("handleTCPConnection (%+v) caught an EOF\n", c)
         break
      }
      fmt.Printf("Iterating buffer. . .\n\t")
      for buffer_index := 0; buffer_index < count; buffer_index+=1 {
         fmt.Printf("%d ", read_buffer[buffer_index])
         if (buffer_index + 1) % 30 == 0 {
            fmt.Printf("\n\t")
         }
      }
      socks_header, err := ParseHeader(read_buffer, count)
      if err == nil {
         // Valid socks request header
         fmt.Printf("Received socks header, version %d\n", (int)(socks_header.version))
         c.Write([]byte{0x05, 0x00})
         goto NEXT
      }

      _, err = ParseRequestHeader(read_buffer, count)
      if err == nil {
         // Valid socks request header
      } else {
         panic(err)
      }
NEXT:
      println("\nFinished iterating buffer")

   }
   fmt.Printf("handleTCPConnection for connection %+v is exiting/closing\n", c)
   return nil
}
func (sp *SocksProxy) ListenAndHandle() error {
   for {
      conn, err := sp.listener.AcceptTCP()
      if err != nil {
         panic(err)
      }
      println("Blocking on incoming connection")
      sp.handleTCPConnection(conn)
      println("Unblocked!")
   }
   return nil
}


func main() {
   proxy, err := NewSocksProxy("tcp4", "0.0.0.0:1080")
   if err != nil {
      panic(err)
   }
   proxy.ListenAndHandle()
}
