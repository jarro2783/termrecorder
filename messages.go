package termrecorder

import fb "github.com/jarro2783/featherbyte"

const (
    SendUser  = fb.UserMessageStart
    WatchUser = iota
)

type UserRequest struct {
    User string
}
