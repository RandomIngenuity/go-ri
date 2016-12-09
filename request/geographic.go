package rirequest

import (
    "bitbucket.com/dsoprea/go-ri/common"
)

// Mime-type mappings
var (
    FormatMimetypeMapping = map[string]string{
        CtKml: ricommon.FormatKml,
        CtGeojson: ricommon.FormatGeoJson,
    }
)
