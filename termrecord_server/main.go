package main

import "fmt"
import "container/list"
import fb "github.com/jarro2783/featherbyte"
import "github.com/jarro2783/termrecorder"
import "golang.org/x/net/context"

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

type listener struct {
    endpoint *fb.Endpoint
    subscribe chan<- subscribeRequest
    publish chan<- publishRequest

    send chan []byte

    context.CancelFunc
}

func (l *listener) Bytes(data []byte) {
    if l.send != nil {
        l.send <- data
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

    l.send = dataChannel

    ctx, cancel := context.WithCancel(context.Background())

    l.CancelFunc = cancel

    go publisher(ctx, user.User, dataChannel, subscribe)
}

func publisher(ctx context.Context, user string, data <-chan []byte,
    register <-chan publisherChannel) {
    subscribers := list.New()
    remove := make([]*list.Element, 0, 5)

    sessionData := make([]byte, 0, 1024*1024)

    var newStart int = 0
    var needle = "\033[2J"

    Loop:
    for {
        select {
            case <-ctx.Done():
            break Loop

            case bytes, ok := <-data:

            if !ok {
                break Loop
            }

            nextClear := len(sessionData)

            //store data for this session
            sessionData = append(sessionData, bytes...)

            needlePos := 0
            for ; nextClear < len(sessionData); nextClear++ {
                if needlePos != len(needle) {
                    if sessionData[nextClear] == needle[needlePos] {
                        needlePos++
                    } else {
                        needlePos = 0
                    }
                } else {
                    newStart = nextClear
                    fmt.Printf("New clear at %d\n", newStart)
                    needlePos = 0
                }
            }

            //fmt.Printf("%s", string(bytes))
            fmt.Printf("Got %d bytes\n", len(bytes))

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
                break Loop
            }

            fmt.Printf("Adding subscriber starting at byte %d\n", newStart)

            subscribers.PushBack(subscriber)

            //send them everything since the last clear screen
            subscriber.data <- sessionData[newStart:]
        }
    }

    fmt.Printf("Terminating publisher for %s\n", user)
}

func subscriber(ctx context.Context,
    endpoint *fb.Endpoint,
    data subscriberChannel) {
    Loop:
    for {
        select {
            case <- ctx.Done():
            break Loop

            case d, ok := <-data.data:
            if ok {
                //fmt.Printf("%s", string(d))
                err := endpoint.WriteBytes(d)

                if err != nil {
                    close(data.done)
                }
            } else {
                break Loop
            }
        }
    }

    fmt.Printf("Terminating subscriber\n")
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

    ctx, cancel := context.WithCancel(context.Background())

    l.CancelFunc = cancel

    go subscriber(ctx, l.endpoint, data)
}

func (l *listener) Exiting() {
    if l.send != nil {
        close(l.send)
    }

    if l.CancelFunc != nil {
        l.CancelFunc()
    }
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
