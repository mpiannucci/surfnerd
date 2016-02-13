package surfnerd

type WindModelType int64

const (
	GFS WindModelType = iota
	NAM
)

const (
	gfsURL  = "http://nomads.ncep.noaa.gov:9090/dods/nam/nam20160213/nam_na_12z"
	nameURL = "http://nomads.ncep.noaa.gov:9090/dods/nam/nam20160213/nam_na_12z"
)
