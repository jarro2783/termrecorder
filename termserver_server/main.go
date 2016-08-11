package main

import fb "github.com/jarro2783/featherbyte"
import "github.com/jarro2783/termserver"
import "fmt"

type connectionHandler struct {
}

func (handler *connectionHandler) Connection(*fb.Endpoint) {
    fmt.Printf("New connection")
}

type listener struct {
}

func (l *listener) Bytes(data []byte) {
    fmt.Printf("%s", string(data))
}

func (l *listener) Send(user *termserver.UserRequest) {
    fmt.Printf("Got user request for %s\n", user.User)
}

func (l *listener) Watch(user *termserver.UserRequest) {
    fmt.Printf("Watch user")
}

func main() {
    termserver.Listen("localhost", 34234, new(connectionHandler),
        new(listener))
}
