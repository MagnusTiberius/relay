# relay
Relay Demo

test:

terminal 1
relay$ go run relay/main.go 8056

terminal 2
relay$ go run rechoserver/main.go 127.0.0.1 8056
relay connection 127.0.0.1 8057

terminal 3
relay$ telnet 127.0.0.1 8057

terminal 4
relay$ telnet 127.0.0.1 8057

terminal 5
relay$ go run rechoserver/main.go 127.0.0.1 8056
relay connection 127.0.0.1 8058

terminal 6
relay$ telnet 127.0.0.1 8058

terminal 7
relay$ telnet 127.0.0.1 8058
