package surfnerd

const (
	baseMultigridUrl = "http://nomads.ncep.noaa.gov:9090/dods/wave/mww3/%[1]s/%[2]s%[1]s_%[3]s.ascii?time[0:60],dirpwsfc.dirpwsfc[0:60][%[4]d][%[5]d],htsgwsfc.htsgwsfc[0:60][%[4]d][%[5]d],perpwsfc.perpwsfc[0:60][%[4]d][%[5]d],swdir_1.swdir_1[0:60][%[4]d][%[5]d],swdir_2.swdir_2[0:60][%[4]d][%[5]d],swell_1.swell_1[0:60][%[4]d][%[5]d],swell_2.swell_2[0:60][%[4]d][%[5]d],swper_1.swper_1[0:60][%[4]d][%[5]d],swper_2.swper_2[0:60][%[4]d][%[5]d],ugrdsfc.ugrdsfc[0:60][%[4]d][%[5]d],vgrdsfc.vgrdsfc[0:60][%[4]d][%[5]d],wdirsfc.wdirsfc[0:60][%[4]d][%[5]d],windsfc.windsfc[0:60][%[4]d][%[5]d],wvdirsfc.wvdirsfc[0:60][%[4]d][%[5]d],wvhgtsfc.wvhgtsfc[0:60][%[4]d][%[5]d],wvpersfc.wvpersfc[0:60][%[4]d][%[5]d]"
)

type WaveModel struct {
	Name               string
	Description        string
	BottomLeftLocation Location
	TopRightLocation   Location
	LocationResolution float64
	TimeResolution     float64
}

func (w WaveModel) ContainsLocation(loc Location) bool {
	if loc.Latitude > w.BottomLeftLocation.Latitude && loc.Latitude < w.TopRightLocation.Latitude {
		if loc.Longitude > w.BottomLeftLocation.Longitude && loc.Longitude < w.TopRightLocation.Longitude {
			return true
		}
	}
	return false
}

func (w WaveModel) LocationIndices(loc Location) (int, int) {
	if !w.ContainsLocation(loc) {
		return -1, -1
	}

	// Find the offsets from the minimum lat and long
	latOffset := loc.Latitude - w.BottomLeftLocation.Latitude
	lonOffset := loc.Longitude - w.BottomLeftLocation.Longitude

	// Get the indexes and return them
	latIndex := int(latOffset / w.LocationResolution)
	lonIndex := int(lonOffset / w.LocationResolution)
	return latIndex, lonIndex
}

func NewEastCoastWaveModel() *WaveModel {
	return &WaveModel{
		Name:               "multi_1.at_10m",
		Description:        "Multi-grid wave model: US East Coast 10 arc-min grid",
		BottomLeftLocation: Location{0.00, 260.00},
		TopRightLocation:   Location{55.00011, 310.00011},
		LocationResolution: 0.167,
		TimeResolution:     0.125,
	}
}

func NewWestCoastWaveModel() *WaveModel {
	return &WaveModel{
		Name:               "multi_1.wc_10m",
		Description:        "Multi-grid wave model: US West Coast 10 arc-min grid",
		BottomLeftLocation: Location{25.00, 210.00},
		TopRightLocation:   Location{50.00005, 250.00008},
		LocationResolution: 0.167,
		TimeResolution:     0.125,
	}
}

func NewPacificIslandsWaveModel() *WaveModel {
	return &WaveModel{
		Name:               "multi_1.ep_10m",
		Description:        "Multi-grid wave model: Pacific Islands (including Hawaii) 10 arc-min grid",
		BottomLeftLocation: Location{-20.00, 130.00},
		TopRightLocation:   Location{30.0001, 215.00017},
		LocationResolution: 0.167,
		TimeResolution:     0.125,
	}
}

func GetAllAvailableWaveModels() []*WaveModel {
	eastCoastModel := NewEastCoastWaveModel()
	westCoastModel := NewWestCoastWaveModel()
	pacificIslandsModel := NewPacificIslandsWaveModel()
	return []*WaveModel{
		eastCoastModel,
		westCoastModel,
		pacificIslandsModel,
	}
}
