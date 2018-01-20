package ricommon

import (
    "bytes"

    "encoding/json"

    "github.com/dsoprea/go-logging"
)

func ParseJsonFromString(raw string) (p map[string]interface{}, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = DistillError(state)
        }
    }()

    b := bytes.NewBufferString(raw)
    d := json.NewDecoder(b)

    p = make(map[string]interface{})
    err = d.Decode(&p)
    log.PanicIf(err)

    return p, err
}

func ParseJsonFromStringTo(raw string, out interface{}) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = DistillError(state)
        }
    }()

    b := bytes.NewBufferString(raw)
    d := json.NewDecoder(b)

    err = d.Decode(out)
    log.PanicIf(err)

    return err
}
