package termrecorder

import "os"

type Uploader interface {
    Upload(user string, filename string, source *os.File)
}
