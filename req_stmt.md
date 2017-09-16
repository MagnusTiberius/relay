# The problem

This homework problem is written to be challenging. There might be terms or
systems you are unfamiliar with in this problem but that's okay. What
we're interested in is *how* you tackle and solve problems, especially problems
you haven't seen before. This problem might require some research and certainly
some questions and clarification from us. Please don't get stuck! If you get
stuck, just let me know and I'll be more than happy to provide help of any 
kind. Not everyone we've hired has completely solved this problem!

Problem: write a generic TCP relay.

Suppose you have two programs that can talk to a relay computer via TCP, but 
can not reach each other (because of NAT or firewall issues). 

You are asked to:
 1. Implement a relay server program. We will run this program on the relay
    computer that the other programs *can* talk to. This program will listen 
    for TCP connections on behalf of other programs, alert the other programs
    when an incoming connection has been made, and forward all 
    incoming/outgoing traffic from the incoming connection to the other program 
    and vice versa, regardless of application protocol (so, this is a level
    lower than HTTP or SSL or something). The relay server program will be 
    contacted by two types of applications - programs needing a relay and 
    other programs wanting to talk to the program needing a relay.
 2. Describe in clear terms to another programmer how they would enable a
    program they already have, behind some firewall, to be able to use your 
    relay server. Remember that every program in this scenario can contact the 
    relay server, but the relay server can't *initiate* communication with any 
    other program (due to firewalls, etc).
 3. To help yourself debug and for us to evaluate your work, write a small
    application that uses the relay server (maybe an echo server). Once your
    server has established a successfully relayed port, it should output 
    what its new public address is, and we should be able to contact it, 
    through the relay, with existing programs like telnet or netcat or 
    something. Note that this requires your relay server interface to notify 
    relayed clients of their public address.
 4. Your relay server should work generally for any application level protocol
    (not just the echo server), for multiple concurrent clients, where each 
    client may have its own multiple concurrent clients.

Use whatever language and tools you feel most comfortable with, though we
ask that you try to limit yourself to standard libraries where possible.

In order to test this, make sure that your relay program can be started on the
command line:

```
./relay port
```

where port is the port the relay listens on, and how your echoserver will 
contact the relay. You can (and probably should) accept other arguments, but 
choose reasonable defaults.

Your echo server program should be started on the command line:

```
./echoserver relayhost relayport
```

where relayhost is the host the relay is running on, and port is the port you
provided to the relay program. Upon receiving a (host, port) pair from the
relay, the echo server program should notify us what its relayed hostname/port
pair is.

Example session:

```
$ ./relay 8080 &
$ ./echoserver localhost 8080 &
established relay address: localhost:8081
$ telnet localhost 8081
Hello, world
Hello, world
```

Thanks, and good luck!
