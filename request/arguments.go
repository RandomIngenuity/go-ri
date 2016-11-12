package rirequest

import (
    "net/http"

    "fmt"

    "github.com/gorilla/mux"
    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-multiparse"
)

// Other
var (
    argumentLogger = log.NewLogger("ri.request.arguments")
)

func GetStringRequestPathParameter(r *http.Request, name string) string {
    vars := mux.Vars(r)
    value := vars[name]

    if value == "" {
        log.Panic(fmt.Errorf("[%s] path parameter empty in request", name))
    }

    return value
}

func GetRequestPathParameter(r *http.Request, name, kindName string) interface{} {
    valueRaw := GetStringRequestPathParameter(r, name)
    return parse.Parse(valueRaw, kindName)
}
