package main

import (
    "net"
    "bufio"
    "fmt"
    "os"
    "strconv"
)

const(
  IP = "127.0.0.1"
  PORT = 8080
)

var (
  conn_pool map[string] *net.Conn
  port int
)

type Event struct {
  Name string
  Client net.Conn
  Msg []byte
}

func main() {
    if len(os.Args) == 1 {
      port = PORT
    }
    port,_ = strconv.Atoi(os.Args[1])
    //fmt.Println("relay port:"+port)
    setupListener(IP, port)
}

func setupListener(ip string , port int) {
  addr := fmt.Sprintf("%s:%d", ip, port)
  server, err := net.Listen("tcp", addr)
  n := port + 1
  for err != nil {
    addr := fmt.Sprintf("%s:%d", ip, n)
    server, err = net.Listen("tcp", addr)
  }
  if server == nil {
      panic(fmt.Sprintf("%s: %v","Listen failure: ",err))
  }
  fmt.Println("Relay connection: "+addr)
  conn_pool = map[string] *net.Conn{}
  conns := handleConns(server)
  for {
      go handleConn(<-conns)
  }

}

func handleConns(l net.Listener) chan net.Conn {
    ch := make(chan net.Conn)
    i := 0
    go func() {
        for {
            client, err := l.Accept()
            if client == nil {
                panic(fmt.Sprintf("%s: %v","Listener Accept() failure: ",err))
                continue
            }
            i++
            fmt.Printf("%d: %v accepted %v\n", i, client.LocalAddr(), client.RemoteAddr())
            conn_pool[fmt.Sprintf("%v",client.RemoteAddr())] = &client
            client.Write([]byte("Welcome to relay utopia\n"))
            handleEvent( Event{Name:"CONNECT_EVENT", Client:client})
            ch <- client
        }
    }()
    return ch
}

func handleEvent(e Event ) {
  //fmt.Println(e.Name)
  switch e.Name {
    case "CONNECT_EVENT":
      for _,c := range conn_pool {
        conn := *c
        addr := fmt.Sprintf("%v",e.Client.RemoteAddr())
        if fmt.Sprintf("%v",conn.RemoteAddr()) != addr {
          msg_connect := []byte(fmt.Sprintf("%s has connected\n",addr))
          conn.Write(msg_connect)
        }
      }
    case "CLIENT_MSG" :
      for _,c := range conn_pool {
        addr := fmt.Sprintf("%v",e.Client.RemoteAddr())
        conn := *c
        if fmt.Sprintf("%v",conn.RemoteAddr()) != addr {
          msg := fmt.Sprintf("SEND TO:%v: MESSAGE:%s",e.Client.RemoteAddr(), string(e.Msg))
          conn.Write([]byte(msg))
        }
      }
    default:
  }
}

func handleConn(c net.Conn) {
    b := bufio.NewReader(c)
    for {
        msg, err := b.ReadBytes('\n')
        if err != nil {
            break
        }
        fmt.Printf("RECEIVE>> Addr:%v Data:%s",c.RemoteAddr(),string(msg))
        handleEvent( Event{Name:"CLIENT_MSG", Client:c, Msg:msg})
    }
}
