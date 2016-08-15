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
    listeners chan<- (chan<- []byte)
}

func coordinator(subscribe <-chan subscribeRequest,
    publish <-chan publishRequest) {

    select {
        case _, ok := <-subscribe:

        if ok {
        }

        case _, ok := <-publish:

        if ok {
        }
    }
}

func (handler *connectionHandler) Connection(*fb.Endpoint) {
    fmt.Printf("New connection")
}

type listener struct {
    subscribe chan<- subscribeRequest
    publish chan<- publishRequest
}

func (l *listener) Bytes(data []byte) {
    fmt.Printf("%s", string(data))
}

func (l *listener) Send(user *termrecorder.UserRequest) {
    fmt.Printf("Got user request for %s\n", user.User)

    var response chan chan<- []byte = make(chan chan<- []byte)

    var pub = publishRequest {
        user.User,
        response,
    }

    l.publish <- pub

    dataChannel := make(chan []byte)

    go publisher(dataChannel, response)
}

func publisher(data <-chan []byte, register <-chan chan<- []byte) {
    subscribers := make([](chan<- []byte), 0, 10)
    for {
        select {
            case bytes, ok := <-data:

            if !ok {
                break
            }

            for s := range(subscribers) {
                subscribers[s] <- bytes
            }

            case subscriber, ok := <-register:

            if !ok {
                break
            }

            subscribers = append(subscribers, subscriber)
        }
    }
}

func (l *listener) Watch(user *termrecorder.UserRequest) {
    fmt.Printf("Watch user")
}

func makeListener(subscribe chan<- subscribeRequest,
    publish chan<- publishRequest) *listener {

    l := new(listener)
    l.subscribe = subscribe
    l.publish = publish

    return l
}

type connectionHandler struct {
    *listener
}

func makeHandler(listen *listener) *connectionHandler {
    h := new(connectionHandler)

    h.listener = listen

    return h
}

func main() {
    subscribe := make(chan subscribeRequest)
    publish := make(chan publishRequest)

    go coordinator(subscribe, publish)

    termrecorder.Listen("localhost", 34234,
        makeHandler(makeListener(subscribe, publish)))
}
