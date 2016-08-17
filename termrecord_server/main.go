package main

import "fmt"
import "container/list"
import fb "github.com/jarro2783/featherbyte"
import "github.com/jarro2783/termrecorder"

type publisherChannel struct {
    done <-chan struct{}
    data chan<- []byte
}

type subscriberChannel struct {
    done chan<- struct{}
    data <-chan []byte
}

type subscribeRequest struct {
    user string
    response chan<- subscriberChannel
}

type publishRequest struct {
    user string
    listeners chan<- publisherChannel
}

func coordinator(subscribe <-chan subscribeRequest,
    publish <-chan publishRequest) {

    users := make(map[string]chan<- publisherChannel)

    for {
        select {
            case s, ok := <-subscribe:

            fmt.Printf("Coordinator got subscribe request for %s\n",
                s.user)

            if ok {
                if u := users[s.user]; u != nil {
                    done := make(chan struct{})
                    data := make(chan []byte)

                    s.response <- subscriberChannel {
                        done,
                        data,
                    }

                    u <- publisherChannel {
                        done,
                        data,
                    }
                } else {
                    close(s.response)
                }
            }

            case p, ok := <-publish:

            fmt.Printf("Coordinator got publish request\n")

            if ok {
                users[p.user] = p.listeners
            }
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
    fmt.Printf("Got request to send %s\n", user.User)
    var subscribe chan publisherChannel = make(chan publisherChannel)

    var pub = publishRequest {
        user.User,
        subscribe,
    }

    l.publish <- pub

    dataChannel := make(chan []byte)

    l.send = &dataSender{dataChannel}

    go publisher(dataChannel, subscribe)
}

func publisher(data <-chan []byte, register <-chan publisherChannel) {
    subscribers := list.New()
    remove := make([]*list.Element, 0, 5)

    for {
        select {
            case bytes, ok := <-data:

            if !ok {
                break
            }

            fmt.Printf("%s", string(bytes))

            for e := subscribers.Front(); e != nil; e = e.Next() {
                pc := e.Value.(publisherChannel)

                select {
                    case pc.data <- bytes:

                    case <-pc.done:
                    //remove the current channel
                    remove = append(remove, e)
                }
            }

            for r := range(remove) {
                subscribers.Remove(remove[r])
            }

            remove = remove[0:0]

            case subscriber, ok := <-register:

            if !ok {
                break
            }

            fmt.Printf("Adding subscriber\n")

            subscribers.PushBack(subscriber)
        }
    }
}

func subscriber(endpoint *fb.Endpoint, data subscriberChannel) {
    for {
        select {
            case d, ok := <-data.data:
            if ok {
                fmt.Printf("%s", string(d))
                endpoint.WriteBytes(d)
            } else {
                break
            }
        }
    }
}

func (l *listener) Watch(user *termrecorder.UserRequest) {
    fmt.Printf("Got request to watch %s\n", user.User)
    response := make(chan subscriberChannel)
    request := subscribeRequest {
        user.User,
        response,
    }

    l.subscribe <- request

    data, ok := <-response

    if !ok {
        fmt.Printf("Invalid watch request")
        return
    }

    fmt.Printf("starting subscriber\n")

    go subscriber(l.endpoint, data)
}

func makeListener(endpoint *fb.Endpoint, subscribe chan<- subscribeRequest,
    publish chan<- publishRequest) *listener {

    l := new(listener)
    l.subscribe = subscribe
    l.publish = publish
    l.endpoint = endpoint

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
