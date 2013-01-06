package main

import (
   "fmt"
   "errors"
   "encoding/binary"
)

const (
   // Commands, used in requests after authenticaiton
   SOCKS_COMMAND_CONNECT = 0x1
   SOCKS_COMMAND_BIND = 0x2
   SOCKS_COMMAND_UDP_ASSOCIATE = 0x3 // Won't be used by me
   // Address types
   SOCKS_ADDRESS_TYPE_IPV4 = 0x01
   SOCKS_ADDRESS_TYPE_DOMAIN_NAME = 0x03
   SOCKS_ADDRESS_TYPE_IPV6 = 0x04
)

// Version, CMD, reserved octet, address type, address, port (short)
type SocksRequestHeader struct {
   version byte
   command byte
   // reserved byte, must be 00 or invalid
   reserved byte
   address_type byte
   // Destination address is 4 bytes under ipv4
   // 16 bytes under ipv4
   // destination_address[0] bytes long under FQDN
   destination_address []byte
   destination_port uint16
}

func ParseSocksRequestHeader(input []byte, input_count int) (*SocksRequestHeader, error) {
   // TODO: test this _whole_ thing
   srh := new(SocksRequestHeader)
   // version + command + reserved + address type + port
   minimum_length := 1 + 1 + 1 + 1 + 2 // May need to be revised, doesn't take dest_addr into consideration

   if input_count < minimum_length {
      return nil, errors.New(fmt.Sprintf("Invalid Socks Request Header: byte sequence is too short (sequence %v)\n", input))
   }

   srh.version = input[0]
   srh.command = input[1]
   srh.reserved = input[2]
   srh.address_type = input[3]
   
   if (srh.address_type == SOCKS_ADDRESS_TYPE_IPV4) {
      srh.destination_address = input[4:8]
      srh.destination_port = binary.LittleEndian.Uint16(input[8:])
   } else if (srh.address_type == SOCKS_ADDRESS_TYPE_IPV6) {
      srh.destination_address = input[4:20]
      srh.destination_port = binary.LittleEndian.Uint16(input[20:])
   } else if (srh.address_type == SOCKS_ADDRESS_TYPE_DOMAIN_NAME) {
      // First byte specifies length
      domain_length := input[4]
      srh.destination_address = input[5:(5+domain_length)]
      srh.destination_port = binary.LittleEndian.Uint16(input[(5+domain_length):])
   }
   return srh, nil
}

func (srh *SocksRequestHeader) String() string {
   base_string := "\tReceived socks version %d request with command %d (%s)\n"
   base_string += "\tAddress type was %d (%s)\n"
   base_string += "\tDestination Address was %s\n"
   base_string += "\tDestination Port was %d\n"

   version := srh.version
   command := srh.command
   command_string := ""
   address_type := srh.address_type
   address_string := ""
   
   destination_address := (string)(srh.destination_address)
   destination_port := srh.destination_port

   switch command {
      case SOCKS_COMMAND_CONNECT: command_string = "CONNECT"
      case SOCKS_COMMAND_BIND: command_string = "BIND"
      case SOCKS_COMMAND_UDP_ASSOCIATE: command_string = "UDP ASSOCIATE"
   }

   switch address_type {
      case SOCKS_ADDRESS_TYPE_IPV4: address_string = "IPv4"
      case SOCKS_ADDRESS_TYPE_IPV6: address_string = "IPv6"
      case SOCKS_ADDRESS_TYPE_DOMAIN_NAME: address_string = "Domain Name"
   }
   args := make([]interface{}, 7)
   args = append(args, version, command, command_string)
   args = append(args, address_type, address_string)
   args = append(args, destination_address, destination_port)
   return fmt.Sprintf(base_string, version, command, command_string) +
          fmt.Sprintf("Text")

}
