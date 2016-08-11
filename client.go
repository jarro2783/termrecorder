package termserver

import "fmt"
import fb "github.com/jarro2783/featherbyte"
import "encoding/json"

type Writer struct {
    endpoint *fb.Endpoint
    reader   fb.ReadData
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
