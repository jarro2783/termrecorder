package termrecorder

import fb "github.com/jarro2783/featherbyte"

const (
    SendUser  = iota + fb.UserMessageStart
    WatchUser = iota + fb.UserMessageStart
)

type WatchRequest struct {
    User string
}

type PublishRequest struct {
    User string
    Gameid string
}
