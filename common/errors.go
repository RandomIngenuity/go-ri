package ricommon

import (
    "errors"
    "github.com/dsoprea/go-logging"
)

var (
    ErrNotFound = errors.New("query not found")
    ErrAlreadyExists = errors.New("already exists")
    ErrArgumentError = errors.New("argument error")
    ErrMaintenanceRequiresDebug = errors.New("this type of maintenance requires debug mode")
)

func IsErrorMessage(err error, message string) bool {
    return err != nil && err.Error() == message
}

// DistillError makes sure the error is an error (some system errors, for
// example reflect, will emit a string) and that it's wrapped with an error-
// stack.
//
// Note that this will culminate in the first stack-frame being from here if an
// error is not already wrapped. However, it's the price we pay for the
// simplicity of having this one function for manage this.
func DistillError(state interface{}) (err error) {
    if err, ok := state.(error); ok == true {
        return log.Wrap(err)
    }

    // If Go generates an error that's just a string.
    // (e.g. reflection might do this, reflect.MakeSlice of non-slice type)
    return log.Errorf("Non-error error: [%v]", state)
}
