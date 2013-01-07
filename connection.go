package main

import (
   "io"
   "fmt"
   "os"
   "net"
   "log"
   "time"
)

func printBytesFormatted(buf []byte, length, line_length uint, hex bool) {
   as_decimal := ""
   as_hex := ""
   var idx uint = 0
   print("\t")
   for idx = 0; idx < length; idx++ {
      if (idx + 1) % 30 == 0 {
         as_hex += "\n\t"
         as_decimal += "\n\t"
      }
      as_hex += fmt.Sprintf("%3X", buf[idx])
      as_decimal += fmt.Sprintf("%3d", buf[idx])
   }
   
   if hex {
      print(as_hex)
   } else {
      print(as_decimal)
   }
   if (idx + 1) % 30 != 0 {
      println()
   }
}
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
   stdout.Printf("Connection.Handle received connection %+v\n", conn, conn.tcp_conn)
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

      // Set a 5 second timeout for the read
      now := time.Now()
      conn.tcp_conn.SetReadDeadline(now.Add(5 * time.Second))
      count, err := conn.tcp_conn.Read(read_buffer)
      if err != nil {
         if err == io.EOF {
            debug.Printf("Connection.Handle caught an EOF on TCPConn %+v\n", conn, conn.tcp_conn)
            break
         } else if err.(net.Error).Timeout() {
            debug.Printf("Connection.Handle timed out, connection closed")
            return nil
         } else {
            panic(err)
         }
      }

      debug.Println("Iterating buffer. . .")
      printBytesFormatted(read_buffer, (uint)(count), 30, false) // 30 lines dec

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
         println(req_header.String())
      } else {
         panic(err)
      }
NEXT:
      debug.Println("Finished iterating buffer")

   }
   debug.Printf("Connection.Handle (%+v) for connection %+v is exiting/closing\n", conn, conn.tcp_conn)
   return nil
}

