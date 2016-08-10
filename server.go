package termserver

import fb "github.com/jarro2783/featherbyte"
import "fmt"

type connection struct {
}

func (c *connection) Connection(ep *fb.Endpoint) {
}

func Listen(address string, port int) {
    connections := new(connection)
    fb.Listen("tcp", fmt.Sprintf("%s:%d", address, port),
        connections, new(writeListener))
}
