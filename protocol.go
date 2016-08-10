package termserver

import "fmt"
import fb "github.com/jarro2783/featherbyte"
import "encoding/json"

type Writer struct {
    endpoint *fb.Endpoint
    reader   fb.ReadData
}

type username struct {
    User string
}

type writeListener struct {
}

func (*writeListener) Data(messageType byte, data []byte) {
}

func (*writeListener) Message(messageType byte, data []byte) {
    switch messageType {
    case RequestUser:
    }
}

func Connect(host string, port int, user string) (*Writer, error) {
    var writer *Writer
    var err error

    writer = new(Writer)
    writer.reader = new(writeListener)
    writer.endpoint, err = fb.Connect("tcp", fmt.Sprintf("%s:%d", host, port),
        writer.reader)

    if err != nil {
        ustruct := username{user}
        juser, _ := json.Marshal(ustruct)
        writer.endpoint.WriteMessage(RequestUser, juser)
    }

    return writer, err
}

func (writer *Writer) Write(data []byte) {
    writer.endpoint.WriteBytes(data)
}

type connection struct {
}

func (c *connection) Connection(ep *fb.Endpoint) {
}

func Listen(address string, port int) {
    connections := new(connection)
    fb.Listen("tcp", fmt.Sprintf("%s:%d", address, port),
        connections, new(writeListener))
}
