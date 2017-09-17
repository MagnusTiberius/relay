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
  PORT = "9055"
  BASE_PORT = 9090
)

var (
  conn_pool map[string] *net.Conn
  port string
  ip string
  relayConn net.Conn
)

type Event struct {
  Name string
  Client net.Conn
  Msg []byte
}

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
  connectToRelay(ip, port)

  setupListener(ip, port)
}

func setupListener(ip string , sport string) {
  nport,_ := strconv.Atoi(sport)
  n := nport + 1
  addr := fmt.Sprintf("%s:%d", ip, n)
  server, err := net.Listen("tcp", addr)
  for err != nil {
    n = n + 1
    addr = fmt.Sprintf("%s:%d", ip, n)
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

func connectToRelay(ip, port string) (net.Conn, error) {
  addr := fmt.Sprintf("%s:%s", ip, port)

  conn, err := net.Dial("tcp", addr)
  if err != nil {
    panic(fmt.Sprintf("%s: %v","Relay connection failure: ",err))
  }
  relayConn = conn

  go func() {
    buf := bufio.NewReader(conn)
    for {
        msg,_ := buf.ReadString('\n')
        if len(msg)>0 {
            fmt.Printf("Msg:%s\n",string(msg))
            handleEvent( Event{Name:"RELAY_MSG", Client:conn, Msg:[]byte(msg)})
        }
    }
  }()

  return relayConn, nil
}

func handleListen(addr string) {
  server, err := net.Listen("tcp", addr)

  if server == nil {
      panic(fmt.Sprintf("%s: %v","Listen failure: ",err))
  }
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
            client.Write([]byte("Welcome to echoserver utopia\n"))
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
          msg := fmt.Sprintf("ECHO>>%v:%s",e.Client.RemoteAddr(), string(e.Msg))
          conn.Write([]byte(msg))
        }
      }
    case "RELAY_MSG" :
      for _,c := range conn_pool {
        client := *c
        msg := fmt.Sprintf("RELAY>>%v:%s\n",e.Client.RemoteAddr(), string(e.Msg))
        client.Write([]byte(msg))
      }
    default:
  }
  relayEvent(e)
}

func relayEvent(e Event ) {
  //msg := fmt.Sprintf("CLIENT>>%v:%s",e.Client.RemoteAddr(), string(e.Msg))
  //relayConn.Write([]byte(msg))
  relayConn.Write(e.Msg)
}

func handleConn(c net.Conn) {
    b := bufio.NewReader(c)
    for {
        msg, err := b.ReadBytes('\n')
        if err != nil {
            break
        }
        fmt.Printf("%v:%s",c.RemoteAddr(),string(msg))
        c.Write(msg)
        handleEvent( Event{Name:"CLIENT_MSG", Client:c, Msg:msg})
    }
}
