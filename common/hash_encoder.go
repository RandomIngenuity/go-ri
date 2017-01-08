package ricommon

import (
    "fmt"

    "encoding/gob"
    "crypto/sha1"

    "github.com/dsoprea/go-logging"
)

func EncodeToSha1DigestString(parts []interface{}) (digest string) {
    h := sha1.New()
    g := gob.NewEncoder(h)

    for _, x := range parts {
        if err := g.Encode(x); err != nil {
            log.Panic(fmt.Errorf("error encoding [%v]: %s", x, err))
        }
    }

    digest = fmt.Sprintf("%x", h.Sum(nil))
    return digest
}
