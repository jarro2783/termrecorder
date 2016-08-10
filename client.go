package termserver

import "fmt"
import fb "github.com/jarro2783/featherbyte"
import "encoding/json"

type Writer struct {
    endpoint *fb.Endpoint
    reader   fb.ReadData
}

type writeListener struct {
    listener Listener
}

func (listen *writeListener) Data(messageType byte, data []byte) {
    listen.listener.Bytes(data)
}

func (listen *writeListener) Message(messageType byte, data []byte) {
    switch messageType {
    case SendUser:
        var user UserRequest
        json.Unmarshal(data, user)
        listen.listener.Send(&user)

    case WatchUser:
        var user UserRequest
        json.Unmarshal(data, user)
        listen.listener.Watch(&user)
    }
}

func Connect(host string, port int, listener Listener) (*Writer, error) {
    var writer *Writer
    var err error

    writer = new(Writer)
    writer.reader = new(writeListener)
    writer.endpoint, err = fb.Connect("tcp", fmt.Sprintf("%s:%d", host, port),
        writer.reader)

    return writer, err
}

func (writer *Writer) Watch(user string) {
    ustruct := UserRequest{user}
    juser, _ := json.Marshal(ustruct)
    writer.endpoint.WriteMessage(SendUser, juser)
}

func (writer *Writer) Write(data []byte) {
    writer.endpoint.WriteBytes(data)
}
