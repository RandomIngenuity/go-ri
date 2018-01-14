package ridata

import (
    "crypto/rand"
    "crypto/sha256"

    "fmt"
    "strings"

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
func GetLargeRandomHexString() string {
    return RandomHexString(LargeUuidLength)
}

// GetMediumRandomHexString returns a random-string that can usually be used as
// a reusable identifier.
func GetMediumRandomHexString() string {
    return RandomHexString(MediumUuidLength)
}
