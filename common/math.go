package ricommon

import (
    "math"
)

var (
    Epsilon = math.Nextafter(1, 2) - 1
)

func IsEqualFloat32(value1, value2 float32) bool {
    return value1 < value2 && float64(value2 - value1) < Epsilon ||
           value1 > value2 && float64(value1 - value2) < Epsilon
}

func IsEqualFloat64(value1, value2 float64) bool {
    return value1 < value2 && float64(value2 - value1) < Epsilon ||
           value1 > value2 && float64(value1 - value2) < Epsilon
}
