package ricommon

import (
    "io"
    "time"

    "github.com/rwcarlsen/goexif/exif"
    "github.com/dsoprea/go-logging"
)

// Other
var (
    exifLog = log.NewLogger("ri.common.exif")
)

// ImageExif Describes information retrieved from EXIF
type ImageExif struct {
    Timestamp time.Time
    Latitude float64
    Longitude float64
}

func NewImageExifWithReader(r io.Reader) (ie *ImageExif, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    ie = new(ImageExif)

    x, err := exif.Decode(r)
    log.PanicIf(err)

    exifTimestamp, err := x.DateTime()
    if err == nil {
        ie.Timestamp = exifTimestamp
    }

    exifLat, exifLong, err := x.LatLong()
    if err == nil {
        ie.Latitude, ie.Longitude = exifLat, exifLong
    }

    return ie, nil
}
