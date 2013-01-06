package main

const (
   SOCKS_VERSION = 0x05
   // Response Reply field options
   SOCKS_REPLY_SUCCESS = 0x00
   SOCKS_REPLY_GENERAL_FAILURE = 0x01
   SOCKS_REPLY_CONNECTION_DISALLOWED = 0x02
   SOCKS_REPLY_NETWORK_UNREACHABLE = 0x03
   SOCKS_REPLY_HOST_UNREACHABLE = 0x04
   SOCKS_REPLY_CONNECTION_REFUSED = 0x05
   SOCKS_REPLY_TTL_EXPIRED = 0x06
   SOCKS_REPLY_COMMAND_UNSUPPORTED = 0x07
   SOCKS_REPLY_ADDRESS_TYPE_UNSUPPORTED = 0x08
   SOCKS_REPLY_UNASSIGNED = 0x09
   SOCKS_REPLY_UNASSIGNED_START = 0x09
   SOCKS_REPLY_UNASSIGNED_END = 0xFF
)

type SocksResponse struct {
   version byte
   reply byte
   // Reserved here, *must* be 0x00
   reserved byte
}
