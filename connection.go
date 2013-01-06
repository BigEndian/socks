package main

import (
   "io"
   "fmt"
   "os"
   "net"
   "log"
)

type Connection struct {
   tcp_conn *net.TCPConn
   auth_header *Header
   req_header *RequestHeader
   debug  *log.Logger
   stdout *log.Logger
   stderr *log.Logger
}

func NewConnection(conn *net.TCPConn) *Connection {
   debug :=  log.New(os.Stderr, fmt.Sprintf("[DEBUG]  [Connection %p] ", conn), log.LstdFlags)
   stdout := log.New(os.Stderr, fmt.Sprintf("[OUTPUT] [Connection %p] ", conn), log.LstdFlags)
   stderr := log.New(os.Stderr, fmt.Sprintf("[ERROR]  [Connection %p] ", conn), log.LstdFlags)
   return &Connection{conn, nil, nil, debug, stdout, stderr}
}

func (conn *Connection) Handle(count *ConnectionCount) error {
   stdout := conn.stdout
   //stderr := conn.stderr
   debug  := conn.debug
   stdout.Printf("Connection.Handle (%+v) received connection %+v\n", conn, conn.tcp_conn)
   conn.tcp_conn.SetReadBuffer(1024)
   conn.tcp_conn.SetWriteBuffer(1024)
   read_buffer := make([]byte, 1024)
   defer func() {
      count.Lock()
      count.Count -= 1
      count.Unlock()
   }()

   for {
      // Clear the buffer
      for i := 0; i < cap(read_buffer); i+=1 {
         read_buffer[i] = 0
      }
      count, err := conn.tcp_conn.Read(read_buffer)
      if err != io.EOF && err != nil {
         panic(err)
      } else if err == io.EOF {
         debug.Printf("Connection.Handle (%+v) caught an EOF on TCPConn %+v\n", conn, conn.tcp_conn)
         break
      }
      debug.Printf("Iterating buffer. . .\n\t")
      var buffer_index int
      for buffer_index = 0; buffer_index < count; buffer_index+=1 {
         fmt.Printf("%d ", read_buffer[buffer_index])
         if (buffer_index + 1) % 30 == 0 {
            fmt.Printf("\n\t")
         }
      }
      if (buffer_index + 1) % 30 != 0 {
         println()
      }

      var socks_header *Header
      var req_header *RequestHeader
      socks_header, err = ParseHeader(read_buffer, count)
      if err == nil {
         // Valid socks request header
         debug.Printf("Received socks header, version %d\n", (int)(socks_header.version))
         conn.tcp_conn.Write([]byte{0x05, 0x00})
         goto NEXT
      }

      req_header, err = ParseRequestHeader(read_buffer, count)
      if err == nil {
         // Valid socks request header
         debug.Printf("Received socks request header, version %d\n", (int)(req_header.version))
         debug.Printf("%s\n", req_header.String())
      } else {
         panic(err)
      }
NEXT:
      debug.Println("Finished iterating buffer")

   }
   debug.Printf("Connection.Handle (%+v) for connection %+v is exiting/closing\n", conn, conn.tcp_conn)
   return nil
}

