package main

import "flag"
import "fmt"
import "os"

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

func main() {
    user := flag.String("user", "", "The name of the user to record")
    host := flag.String("host", "", "The host to send the session to")
    port := flag.Int("port", 34234, "The port to connect to on the host")
    watch := flag.Bool("watch", false, "Watch the requested user")
    gameid := flag.String("gameid", "",
        "An optional game id if sending a session")
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

    var watcher termrecorder.Listener

    if (*watch) {
        watcher = new(WatchListener)
    }

    if (*send) {
        watcher = new(SendListener)
    }

    //fmt.Printf("Session for %s connecting to %s:%d\n", *user, *host, *port)

    writer, err := termrecorder.Connect(*host, *port, watcher)

    if err != nil {
        fmt.Printf("Error connecting to host: %s\n", err.Error())
        os.Exit(1)
    }

    if (*send) {
        writer.Send(*user, *gameid)

        var data []byte = make([]byte, 1024)
        for true {
            n, err := os.Stdin.Read(data)

            if err != nil {
                break
            }

            if n != 0 {
                writer.Write(data[0: n])
            }
        }
    } else {
        writer.Watch(*user)
        var data[1]byte
        os.Stdin.Read(data[:])
    }
}

func cmdError(s string) {
    fmt.Printf("%s\n", s)
    os.Exit(1)
}
