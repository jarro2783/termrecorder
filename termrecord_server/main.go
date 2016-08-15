package main

import fb "github.com/jarro2783/featherbyte"
import "github.com/jarro2783/termrecorder"
import "fmt"

type subscribeRequest struct {
    user string
    response chan<- <-chan []byte
}

type publishRequest struct {
    user string
    response chan<- (chan<- []byte)
}

type subscribeChannel <-chan subscribeRequest
type publishChannel <-chan publishRequest

func coordinator(subscribe subscribeChannel,
    publish publishChannel) {
}

func (handler *connectionHandler) Connection(*fb.Endpoint) {
    fmt.Printf("New connection")
}

type listener struct {
}

func (l *listener) Bytes(data []byte) {
    fmt.Printf("%s", string(data))
}

func (l *listener) Send(user *termrecorder.UserRequest) {
    fmt.Printf("Got user request for %s\n", user.User)
}

func (l *listener) Watch(user *termrecorder.UserRequest) {
    fmt.Printf("Watch user")
}

type connectionHandler struct {
    *listener
    subscribe subscribeChannel
    publish publishChannel
}

func makeHandler(listen *listener, subscribe subscribeChannel,
    publish publishChannel) *connectionHandler {
    h := new(connectionHandler)

    h.listener = listen
    h.subscribe = subscribe
    h.publish = publish

    return h
}

func main() {
    subscribe := make(subscribeChannel)
    publish := make(publishChannel)

    go coordinator(subscribe, publish)

    termrecorder.Listen("localhost", 34234,
        makeHandler(new(listener), subscribe, publish))
}
