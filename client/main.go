package main

import "net"
import "fmt"
import "bufio"
import "os"


const(
  IP = "127.0.0.1"
  PORT = "9055"
)

var (
  chatConn net.Conn
  ip, port string
)

func main() {

  if len(os.Args) == 1 {
    ip = IP
    port = PORT
  }
  if len(os.Args) == 2 {
    port = PORT
    ip = os.Args[1]
  }
  if len(os.Args) == 3 {
    ip = os.Args[1]
    port = os.Args[2]
  }

  //addr := fmt.Sprintf("%s:%s", ip, port)
  //conn, _ := net.Dial("tcp", addr)
  conn, _ := connectToChatServer(ip, port)

  go func() {
    buf := bufio.NewReader(conn)
    for {
        msg,_ := buf.ReadString('\n')
        if len(msg)>0 {
            fmt.Printf(">>%s\n",string(msg))
        }
    }
  }()
  for {
    input := bufio.NewReader(os.Stdin)
    msg, _ := input.ReadString('\n')
    fmt.Fprintf(conn, msg + "")
  }
}


func connectToChatServer(ip, port string) (net.Conn, error) {
  addr := fmt.Sprintf("%s:%s", ip, port)

  conn, err := net.Dial("tcp", addr)
  if err != nil {
    panic(fmt.Sprintf("%s: %v","Chat connection failure: ",err))
  }
  chatConn = conn
  return chatConn, nil
}
