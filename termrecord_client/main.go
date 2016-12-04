package main

import "flag"
import "fmt"
import "container/list"
import "os"
import "time"

import "github.com/jarro2783/termrecorder"

type SendListener struct {
}

type WatchListener struct {
}

//These three should not get any data
func (*SendListener) Bytes([]byte) {
}

func (*SendListener) Send(*termrecorder.PublishRequest) {
}

func (*SendListener) Watch(*termrecorder.WatchRequest) {
}

func (*SendListener) Exiting() {
}

//Only bytes should get data here
func (*WatchListener) Bytes(data []byte) {
    os.Stdout.Write(data)
    os.Stdout.Sync()
}

func (*WatchListener) Send(*termrecorder.PublishRequest) {
}

func (*WatchListener) Watch(*termrecorder.WatchRequest) {
}

func (*WatchListener) Exiting() {
}

func sender(user, host string, port int, input <-chan []byte) {
    running := true

    watcher := new(SendListener)

    for running {
        writer, err := termrecorder.Connect(host, port, watcher)
        if err != nil {

            //sleep for 5 seconds before trying to connect again
            time.Sleep(5)
            continue
        }

        writer.Send(user, "")

        select {
            case data:= <-input:
            writer.Write(data)
        }
    }
}

func main() {
    user := flag.String("user", "", "The name of the user to record")
    host := flag.String("host", "", "The host to send the session to")
    port := flag.Int("port", 34234, "The port to connect to on the host")
    watch := flag.Bool("watch", false, "Watch the requested user")
    send := flag.Bool("send", false,
        "Send a session for the requested user")

    flag.Parse()

    if *watch && *send {
        cmdError("Only one of watch or send may be specified")
    }

    if !(*watch || *send) {
        cmdError("One of watch or send must be specified")
    }

    if *user == "" {
        cmdError("Username not specified")
    }

    if *host == "" {
        cmdError("Host not specified")
    }

    if (*send) {
        sendChannel := make(chan []byte, 100)
        dataBuffer := list.New()

        go sender(*user, *host, *port, sendChannel)

        var data []byte = make([]byte, 1024)
        for true {
            n, err := os.Stdin.Read(data)

            if err != nil {
                break
            }

            if n != 0 {
                dataBuffer.PushBack([]byte(data[0:n]))

                for dataBuffer.Len() > 0 {
                    element := dataBuffer.Front()
                    value := element.Value.([]byte)
                    select {
                        case sendChannel <- value:
                            dataBuffer.Remove(element)
                        default:
                            break
                    }
                }
            }
        }
    } else {
        watcher := new(WatchListener)

        writer, err := termrecorder.Connect(*host, *port, watcher)

        if err != nil {
            fmt.Printf("Error connecting to host: %s\n", err.Error())
            os.Exit(1)
        }

        writer.Watch(*user)
        var data[1]byte
        os.Stdin.Read(data[:])
    }
}

func cmdError(s string) {
    fmt.Printf("%s\n", s)
    os.Exit(1)
}
