package main

import fb "github.com/jarro2783/featherbyte"
import "github.com/jarro2783/termrecorder"

type publisherChannel struct {
    done <-chan struct{}
    data chan<- []byte
}

type subscribeRequest struct {
    user string
    response chan<- <-chan []byte
}

type publishRequest struct {
    user string
    listeners chan<- publisherChannel
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

func (handler *connectionHandler) Connection(endpoint *fb.Endpoint) {
    l := makeListener(endpoint, handler.subscribe, handler.publish)

    go endpoint.StartReader(termrecorder.NewListener(l))
}

type dataSender struct {
    send chan []byte
}

type listener struct {
    endpoint *fb.Endpoint
    subscribe chan<- subscribeRequest
    publish chan<- publishRequest

    send *dataSender
}

func (l *listener) Bytes(data []byte) {
    if l.send != nil {
        l.send.send <- data
    }
}

func (l *listener) Send(user *termrecorder.UserRequest) {
    var response chan publisherChannel = make(chan publisherChannel)

    var pub = publishRequest {
        user.User,
        response,
    }

    l.publish <- pub

    dataChannel := make(chan []byte)

    l.send = &dataSender{dataChannel}

    go publisher(dataChannel, response)
}

func publisher(data <-chan []byte, register <-chan publisherChannel) {
    subscribers := make([]publisherChannel, 0, 10)
    for {
        select {
            case bytes, ok := <-data:

            if !ok {
                break
            }

            for s := range(subscribers) {
                pc := subscribers[s]

                select {
                    case pc.data <- bytes:

                    case <-pc.done:
                    //remove the current channel
                }
            }

            case subscriber, ok := <-register:

            if !ok {
                break
            }

            subscribers = append(subscribers, subscriber)
        }
    }
}

func subscriber(endpoint *fb.Endpoint, data <-chan []byte) {
    for {
        select {
            case d, ok := <-data:
            if ok {
                endpoint.WriteBytes(d)
            } else {
                break
            }
        }
    }
}

func (l *listener) Watch(user *termrecorder.UserRequest) {
    response := make(chan (<-chan []byte))
    request := subscribeRequest {
        user.User,
        response,
    }

    l.subscribe <- request

    data := <-response

    go subscriber(l.endpoint, data)
}

func makeListener(endpoint *fb.Endpoint, subscribe chan<- subscribeRequest,
    publish chan<- publishRequest) *listener {

    l := new(listener)
    l.subscribe = subscribe
    l.publish = publish

    return l
}

type connectionHandler struct {
    subscribe chan subscribeRequest
    publish chan publishRequest
}

func makeHandler(subscribe chan subscribeRequest,
    publish chan publishRequest) *connectionHandler {
    h := new(connectionHandler)

    h.subscribe = subscribe
    h.publish = publish

    return h
}

func main() {
    subscribe := make(chan subscribeRequest)
    publish := make(chan publishRequest)

    go coordinator(subscribe, publish)

    termrecorder.Listen("localhost", 34234,
        makeHandler(subscribe, publish))
}
