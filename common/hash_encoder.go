package ricommon

import (
    "fmt"
    "io"

    "encoding/gob"
    "encoding/binary"
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

func EncodeStringsToSha1DigestString(parts []string) (digest string) {
    h := sha1.New()

    for i, s := range parts {
        if c, err := io.WriteString(h, s); err != nil {
            log.Panic(fmt.Errorf("error writing [%s]: %s", s, err))
        } else if c != len(s) {
            log.Panic(fmt.Errorf("incorrect number of bytes written"))
        }

        // Push a NUL byte to separate between parts but not at the end. This 
        // way, we can also use this function to hash single strings.
        if i < len(parts) - 1 {
            if c, err := h.Write([]byte { 0 }); err != nil {
                log.Panic(err)
            } else if c != 1 {
                log.Panic(fmt.Errorf("incorrect number of bytes written"))
            }
        }
    }

    digest = fmt.Sprintf("%x", h.Sum(nil))
    return digest
}

func EncodeInt64ToSha1DigestString(n int64) (digest string) {
    h := sha1.New()
    binary.Write(h, binary.LittleEndian, n)

    digest = fmt.Sprintf("%x", h.Sum(nil))
    return digest
}
