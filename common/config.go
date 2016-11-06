package ricommon

import (
    "os"
    "strconv"
    "time"
    "fmt"
    "strings"

    "github.com/dsoprea/go-logging"
)

func GetConfigValueString(name string) string {
    value := os.Getenv(name)
    if value == "" {
        log.Panic(fmt.Errorf("string configuration value not found: [%s]", name))
    }
    
    return value
}

func GetConfigValueStringWithDefault(name string, defaultValue string) string {
    value := os.Getenv(name)
    if value == "" {
        value = defaultValue
    }

    return value
}

func GetConfigValueStringSlice(name string, separator string, forceToLower bool) []string {
    values := GetConfigValueStringSliceWithDefault(name, separator, forceToLower, []string {})

    if len(values) == 0 {
        log.Panic(fmt.Errorf("string-slice configuration value not found: [%s]", name))
    }
    
    if forceToLower == true {
        for i, v := range values {
            values[i] = strings.ToLower(v)
        }
    }

    return values
}

func GetConfigValueStringSliceWithDefault(name string, separator string, forceToLower bool, defaultValues []string) []string {
    value := os.Getenv(name)
    var values []string
    
    if value == "" {
        values = defaultValues
    } else {
        values = strings.Split(value, separator)
    }

    return values
}

func GetConfigValueInt32(name string) int32 {
    valueRaw := GetConfigValueString(name)

    value, err := strconv.ParseInt(valueRaw, 10, 32)
    if err != nil {
        log.Panic(fmt.Errorf("32-bit integer configuration value not found: [%s]", name))
    }
    
    return int32(value)
}

// TODO(dustin): Adopt these "default" functions where we should already be using them.

func GetConfigValueInt32WithDefault(name string, defaultValue int32) int32 {
    valueRaw := GetConfigValueStringWithDefault(name, "")

    if valueRaw == "" {
        return defaultValue
    }

    value, err := strconv.ParseInt(valueRaw, 10, 32)
    if err != nil {
        log.Panic(err)
    }

    return int32(value)
}

func GetConfigValueDuration(name string) time.Duration {
    return time.Duration(GetConfigValueInt32(name)) * time.Second
}

func GetConfigValueDurationWithDefault(name string, defaultValue time.Duration) time.Duration {
    valueRaw := GetConfigValueStringWithDefault(name, "")
    if valueRaw == "" {
        return defaultValue
    }
    
    value, err := strconv.ParseInt(valueRaw, 10, 32)
    if err != nil {
        log.Panic(err)
    }
    
    return time.Duration(value) * time.Second
}

func GetConfigValueInt64(name string) int64 {
    valueRaw := GetConfigValueString(name)

    value, err := strconv.ParseInt(valueRaw, 10, 64)
    if err != nil {
        log.Panic(err)
    }

    return value
}

func GetConfigValueInt64WithDefault(name string, defaultValue int64) int64 {
    valueRaw := GetConfigValueStringWithDefault(name, "")

    if valueRaw == "" {
        return defaultValue
    }

    value, err := strconv.ParseInt(valueRaw, 10, 64)
    if err != nil {
        log.Panic(err)
    }

    return value
}

func GetConfigValueBool(name string) bool {
    valueRaw := GetConfigValueString(name)

    value, err := strconv.ParseBool(valueRaw)
    if err != nil {
        log.Panic(err)
    }
    
    return value
}

func GetConfigValueBoolWithDefault(name string, defaultValue bool) bool {
    valueRaw := GetConfigValueStringWithDefault(name, "")

    if valueRaw == "" {
        return defaultValue
    }

    value, err := strconv.ParseBool(valueRaw)
    if err != nil {
        log.Panic(err)
    }

    return value
}
