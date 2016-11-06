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
