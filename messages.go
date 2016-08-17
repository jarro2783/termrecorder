package termrecorder

import fb "github.com/jarro2783/featherbyte"

const (
    SendUser  = iota + fb.UserMessageStart
    WatchUser = iota + fb.UserMessageStart
)

type UserRequest struct {
    User string
}
