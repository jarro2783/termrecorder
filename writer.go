package termserver

import "fmt"
import "net"

type Writer struct {
    server net.Conn
}

func CreateWriter(host string, port int, user string) (*Writer, error) {
    connection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

    var writer *Writer

    if connection != nil {
        writer = new(Writer)
        writer.server = connection
    }

    return writer, err
}
