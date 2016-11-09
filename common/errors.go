package ricommon

import (
    "errors"
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
