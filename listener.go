package termserver

type Listener interface {
    Watch(user *UserRequest)
}
