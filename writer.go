package termserver

import "fmt"
import fb "github.com/jarro2783/featherbyte"
import "encoding/json"

type Writer struct {
    endpoint *fb.Endpoint
    reader   fb.ReadData
}

type UserRequest struct {
    User string
}

type writeListener struct {
    listener Listener
}

func (*writeListener) Data(messageType byte, data []byte) {
}

func (listen *writeListener) Message(messageType byte, data []byte) {
    switch messageType {
    case SendUser:

    case WatchUser:
        var user UserRequest
        json.Unmarshal(data, user)
        listen.listener.Watch(&user)
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
        ustruct := UserRequest{user}
        juser, _ := json.Marshal(ustruct)
        writer.endpoint.WriteMessage(SendUser, juser)
    }

    return writer, err
}

func (writer *Writer) Write(data []byte) {
    writer.endpoint.WriteBytes(data)
}
