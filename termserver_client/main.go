package main

import "flag"
import "fmt"
import "os"

import "github.com/jarro2783/termserver"
//import fb "github.com/jarro2783/featherbyte"

func main() {
    user := flag.String("user", "", "The name of the user to record")
    host := flag.String("host", "", "The host to send the session to")
    port := flag.Int("port", 45123, "The port to connect to on the host")

    flag.Parse()

    if *user == "" {
        cmdError("Username not specified")
    }

    if *host == "" {
        cmdError("Host not specified")
    }

    fmt.Printf("Session for %s connecting to %s:%d\n", *user, *host, *port)

    _, err := termserver.CreateWriter(*host, *port, *user)

    if err != nil {
        fmt.Printf("Error connecting to host: %s\n", err.Error())
        os.Exit(1)
    }
}

func cmdError(s string) {
    fmt.Printf("%s\n", s)
    os.Exit(1)
}
