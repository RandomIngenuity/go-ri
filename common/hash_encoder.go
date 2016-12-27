package ricommon

import (
    "fmt"

    "encoding/binary"
    "crypto/sha1"

    "github.com/dsoprea/go-logging"
)

func EncodeToSha1DigestString(parts []interface{}) (digest string) {
    h := sha1.New()

    for _, x := range parts {
        if err := binary.Write(h, binary.LittleEndian, x); err != nil {
            log.Panic(err)
        }
    }

    digest = fmt.Sprintf("%x", h.Sum(nil))
    return digest
}
