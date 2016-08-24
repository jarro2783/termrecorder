package main

import "bytes"
import "encoding/binary"
import "os"

var clear = "\033[2J"

type frame struct {
    time int
    micro int
    data []byte
}

type framebuffer struct {
    file *os.File
    frames []frame
    data []byte
}

func newFramebuffer(user string, filename string) *framebuffer {
    f := new(framebuffer)
    f.file, _ = os.Create(user + "-" + filename)
    f.frames = make([]frame, 0, 1024)

    return f
}

func (fb *framebuffer) write() {
    for i := range(fb.frames) {
        f := fb.frames[i]
        buffer := new(bytes.Buffer)
        binary.Write(buffer, binary.LittleEndian, uint32(f.time))
        binary.Write(buffer, binary.LittleEndian, uint32(f.micro))
        binary.Write(buffer, binary.LittleEndian, uint32(len(f.data)))
        fb.file.Write(buffer.Bytes())
        fb.file.Write(f.data)
    }
}

func (fb *framebuffer) flush() {
    fb.write()
}

func (fb *framebuffer) close() {
    fb.flush()
    fb.file.Close()
}

func (fb *framebuffer) addFrame(f frame) {
    fb.frames = append(fb.frames, f)

    //look for clear screen in the current frame
    needlePos := 0
    nextClear := 0
    clearStart := 0
    for ; nextClear < len(f.data); nextClear++ {
        if needlePos != len(clear) {
            if f.data[nextClear] == clear[needlePos] {
                needlePos++
            } else {
                needlePos = 0
                clearStart = nextClear
            }
        } else {
            fb.data = f.data[clearStart:]
            break
        }
    }

    //no clear was found
    if nextClear == len(f.data) {
        fb.data = append(fb.data, f.data...)
    }

    if len(fb.frames) == cap(fb.frames) {
        //write out the frames and start again
        fb.write()
        fb.frames = fb.frames[0:0]
    }
}
