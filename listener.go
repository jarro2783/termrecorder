package termrecorder

import "encoding/json"

type Listener interface {
    Watch(user *UserRequest)
    Send(user *UserRequest)

    Bytes(data []byte)

    Exiting()
}

type writeListener struct {
    listener Listener
}

func NewListener(l Listener) *writeListener {
    w := new(writeListener)
    w.listener = l

    return w
}

func (listen *writeListener) Data(messageType byte, data []byte) {
    listen.listener.Bytes(data)
}

func (listen *writeListener) Message(messageType byte, data []byte) {
    switch messageType {
    case SendUser:
        var user UserRequest
        json.Unmarshal(data, &user)
        listen.listener.Send(&user)

    case WatchUser:
        var user UserRequest
        json.Unmarshal(data, &user)
        listen.listener.Watch(&user)
    }
}

func (listen *writeListener) Exiting() {
    listen.listener.Exiting()
}
