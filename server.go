package termrecorder

import fb "github.com/jarro2783/featherbyte"
import "fmt"

func Listen(address string, port int,
    connections fb.ConnectionHandler) {
    fb.Listen("tcp", fmt.Sprintf("%s:%d", address, port),
        connections)
}
