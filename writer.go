package termserver

import "fmt"
import fb "github.com/jarro2783/featherbyte"

type Writer struct {
    endpoint *fb.Endpoint
    reader fb.ReadData
}

type writeListener struct {
}

func (*writeListener) Data(messageType byte, data []byte) {
}

func (*writeListener) Message(messageType byte, data []byte) {
}

func CreateWriter(host string, port int, user string) (*Writer, error) {
    var writer *Writer
    var err error

    writer = new(Writer)
    writer.reader = new(writeListener)
    writer.endpoint, err = fb.Connect("tcp", fmt.Sprintf("%s:%d", host, port),
        writer.reader)

    return writer, err
}
