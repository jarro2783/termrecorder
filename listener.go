package termserver

type Listener interface {
    Watch(user *UserRequest)
    Send(user *UserRequest)

    Bytes(data []byte)
}
