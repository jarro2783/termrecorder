package termrecorder

import "os"

type Uploader interface {
    Upload(user string, gameid string, filename string, source *os.File)
}
