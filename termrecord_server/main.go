package main

import "container/list"
import fb "github.com/jarro2783/featherbyte"
import "flag"
import "fmt"
import "github.com/jarro2783/termrecorder"
import "golang.org/x/net/context"
import "time"

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
                    data := make(chan []byte, 100)

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
            } else {
                fmt.Printf("Error reading from subscribe")
            }

            case p, ok := <-publish:

            fmt.Printf("Coordinator got publish request\n")

            if ok {
                users[p.user] = p.listeners
            } else {
                fmt.Printf("Error reading from publish")
            }
        }
    }
}

func (handler *connectionHandler) Connection(endpoint *fb.Endpoint) {
    l := makeListener(endpoint, handler.subscribe, handler.publish,
        handler.uploaders)

    go endpoint.StartReader(termrecorder.NewListener(l))
}

type listener struct {
    endpoint *fb.Endpoint
    subscribe chan<- subscribeRequest
    publish chan<- publishRequest

    send chan []byte

    context.CancelFunc

    uploaders []termrecorder.Uploader
}

func (l *listener) Bytes(data []byte) {
    if l.send != nil {
        l.send <- data
    }
}

func (l *listener) Send(user *termrecorder.PublishRequest) {
    fmt.Printf("Got request to send %s\n", user.User)
    var subscribe chan publisherChannel = make(chan publisherChannel, 5)

    var pub = publishRequest {
        user.User,
        subscribe,
    }

    l.publish <- pub

    dataChannel := make(chan []byte, 20)

    l.send = dataChannel

    ctx, cancel := context.WithCancel(context.Background())

    l.CancelFunc = cancel

    go publisher(ctx,
        user.User,
        user.Gameid,
        dataChannel,
        subscribe,
        l.uploaders)
}

func publisher(ctx context.Context, user string,
        gameid string, data <-chan []byte,
        register <-chan publisherChannel,
        uploaders []termrecorder.Uploader) {

    subscribers := list.New()
    remove := make([]*list.Element, 0, 5)

    now := time.Now()
    thetime := now.UTC().Format("2006-01-02.15-04-05")

    filename := thetime + ".ttyrec"

    framebuffer := newFramebuffer(user, filename)

    fmt.Printf("Starting session for %s at %s\n", user, thetime)

    Loop:
    for {
        select {
            case <-ctx.Done():
            break Loop

            case bytes, ok := <-data:

            if !ok {
                break Loop
            }

            now = time.Now()

            //store data for this session
            framebuffer.addFrame(
                frame{int(now.Unix()), now.Nanosecond() / 1000, bytes})

            for e := subscribers.Front(); e != nil; e = e.Next() {
                pc := e.Value.(publisherChannel)

                select {
                    case pc.data <- bytes:

                    case <-pc.done:
                    //remove the current channel
                    fmt.Printf("Channel done\n")
                    remove = append(remove, e)
                }
            }

            for _, r := range remove {
                fmt.Printf("Removing a subscriber\n")
                subscribers.Remove(r)
            }

            remove = remove[0:0]

            case subscriber, ok := <-register:

            if !ok {
                break Loop
            }

            subscribers.PushBack(subscriber)

            //send them everything since the last clear screen
            subscriber.data <- framebuffer.data
            //subscriber.data <- sessionData
        }
    }

    framebuffer.flush()
    file := framebuffer.file
    file.Sync()

    for _, u := range uploaders {
        file.Seek(0, 0)
        u.Upload(user, gameid, filename, file)
    }

    framebuffer.close()
    fmt.Printf("Terminating publisher for %s\n", user)
}

func subscriber(ctx context.Context,
    endpoint *fb.Endpoint,
    data subscriberChannel) {
    Loop:
    for {
        select {
            case <- ctx.Done():
            fmt.Printf("subscriber done\n")
            close(data.done)
            break Loop

            case d, ok := <-data.data:
            if ok {
                //fmt.Printf("%s", string(d))
                err := endpoint.WriteBytes(d)

                if err != nil {
                    fmt.Printf("error writing to listener\n")
                    close(data.done)
                    break Loop
                }
            } else {
                fmt.Printf("No more data for subscriber\n")
                close(data.done)
                break Loop
            }
        }
    }

    fmt.Printf("Terminating subscriber\n")
}

func (l *listener) Watch(user *termrecorder.WatchRequest) {
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
        fmt.Printf("Closing sender\n")
        close(l.send)
    }

    if l.CancelFunc != nil {
        fmt.Printf("Cancel goroutines\n")
        l.CancelFunc()
    }
}

func makeListener(
    endpoint *fb.Endpoint,
    subscribe chan<- subscribeRequest,
    publish chan<- publishRequest,
    uploaders []termrecorder.Uploader) *listener {

    l := new(listener)
    l.subscribe = subscribe
    l.publish = publish
    l.endpoint = endpoint
    l.uploaders = uploaders

    return l
}

type connectionHandler struct {
    subscribe chan subscribeRequest
    publish chan publishRequest
    uploaders []termrecorder.Uploader
}

func makeHandler(subscribe chan subscribeRequest,
    publish chan publishRequest,
    bucket string,
    region string,
    root string,
    subpath string) *connectionHandler {

    h := new(connectionHandler)

    h.subscribe = subscribe
    h.publish = publish
    h.uploaders = make([]termrecorder.Uploader, 0, 2)

    if bucket != "" {
        h.uploaders = append(h.uploaders, termrecorder.MakeAwsUploader(
            region, bucket, root, subpath))
    }

    return h
}

func main() {
    port := flag.Int("port", 34234, "The port to listen on")
    bucket := flag.String("bucket", "", "The S3 bucket to upload to")
    region := flag.String("region", "us-east-1", "The region to use for AWS")
    subpath := flag.String("subpath", "",
        "A subpath to put the recording in under the user's name.")
    root := flag.String("root", "",
        "The root directory in the bucket to put the recordings.")

    flag.Parse()

    subscribe := make(chan subscribeRequest, 10)
    publish := make(chan publishRequest, 10)

    go coordinator(subscribe, publish)

    termrecorder.Listen("", *port,
        makeHandler(subscribe, publish, *bucket, *region, *root, *subpath))
}
