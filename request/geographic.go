package rirequest

import (
    "github.com/randomingenuity/go-ri/common"
)

// Mime-type mappings
var (
    FormatMimetypeMapping = map[string]string{
        CtKml: ricommon.FormatKml,
        CtGeojson: ricommon.FormatGeoJson,
    }
)
