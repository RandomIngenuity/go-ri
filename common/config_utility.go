package ricommon

import (
    "io/ioutil"
    "os"

    "golang.org/x/net/context"
    "gopkg.in/yaml.v2"
    "github.com/dsoprea/go-logging"
)

// Misc
var (
    cuLogger = log.NewLogger("ri.common.config_utility")
)

func GetStaticConfig(ctx context.Context, relFilepath string, output interface{}) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
            cuLogger.Errorf(ctx, nil, "Could not get static config: [%s]", err)
        }
    }()

    f, err := os.Open(relFilepath)
    if err != nil {
        log.Panic(err)
    }

    defer f.Close()

    rawData, err := ioutil.ReadAll(f)
    if err != nil {
        log.Panic(err)
    }

    if err := yaml.Unmarshal([]byte(rawData), output); err != nil {
        log.Panic(err)
    }

    return nil
}
