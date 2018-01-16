package ridata

import (
    "fmt"
    "strings"
    "time"

    "crypto/sha256"
    "math/rand"

    "github.com/dsoprea/go-logging"
)

// Other
const (
    LargeUuidLength = 40
    MediumUuidLength = 20
)

func RandomHexString(length int) string {
    list := make([]string, 0)
    n := 0

    for n < length {
        buffer := make([]byte, 10)
        _, err := rand.Read(buffer)
        log.PanicIf(err)

        s := fmt.Sprintf("%x", sha256.Sum256(buffer))
        list = append(list, s)
        n += len(s)
    }

    return strings.Join(list, "")[:length]
}

// GetLargeRandomHexString returns a very unique, long string that is good as a
// one-time password or sharing key.
func LargeRandomHexString() string {
    return RandomHexString(LargeUuidLength)
}

// GetMediumRandomHexString returns a random-string that can usually be used as
// a reusable identifier.
func MediumRandomHexString() string {
    return RandomHexString(MediumUuidLength)
}

func RandomUint64() uint64 {
    return rand.Uint64()
}

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}
