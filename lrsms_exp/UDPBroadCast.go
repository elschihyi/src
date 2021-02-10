package main

import (
  "net"
  //"fmt"
)

func main() {
  pc, err := net.ListenPacket("udp4", ":12345")
  if err != nil {
    panic(err)
  }
  defer pc.Close()

  addr,err := net.ResolveUDPAddr("udp4", "172.16.1.5:5683")
  if err != nil {
    panic(err)
  }

  _,err = pc.WriteTo([]byte("data to transmit"), addr)
  if err != nil {
    panic(err)
  }
}
