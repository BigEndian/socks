package main
import (
   "fmt"
   "errors"
)

const (
   // Methods, used when a client is authenticating
   SOCKS_METHOD_NO_AUTH_REQUIRED = 0x00
   SOCKS_METHOD_GSSAPI = 0x01
   SOCKS_METHOD_USERNAME_PASSWORD = 0x02
   SOCKS_METHOD_IANA_ASSIGNED = 0x03
   SOCKS_METHOD_IANA_ASSIGNED_START = 0x03
   SOCKS_METHOD_IANA_ASSIGNED_END = 0x7F
   SOCKS_METHOD_RESERVED_PRIVATE_METHODS = 0x80
   SOCKS_METHOD_RESERVED_PRIVATE_METHODS_START = 0x80
   SOCKS_METHOD_RESERVED_PRIVATE_METHODS_END = 0xFE
   SOCKS_METHOD_NOT_ACCEPTABLE = 0xFF // Used as a reply when no given method sates the server
)
type Header struct {
   version byte
   method_count byte
   methods []byte
}

func ParseHeader(input []byte, input_count int) (*Header, error) {
   sh := new(Header)
   if len(input) < 3 {
      // Can't be valid
      return nil, errors.New(fmt.Sprintf("Invalid Socks Header Byte Sequence: %v", input))
   }
   sh.version = input[0]
   sh.method_count = input[1]
   expected_count := input_count - 2
   if expected_count != (int)(sh.method_count) {
      // Too few or too many methods given
      return nil, errors.New(fmt.Sprintf("Invalid Socks Header Byte Sequence: invalid method count (sequence %v)\n", input))
   }
   sh.methods = make([]byte, sh.method_count)
   sh.methods = append(sh.methods, input[2:]...)
   return sh, nil
}
