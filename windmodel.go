package surfnerd

type WindModelType int64

const (
    GFS = iota WindModelType,
    gfsURL = "http://nomads.ncep.noaa.gov:9090/dods/nam/nam20160213/nam_na_12z"

    NAM = iota WindModelType,
    nameURL = "http://nomads.ncep.noaa.gov:9090/dods/nam/nam20160213/nam_na_12z"
)
