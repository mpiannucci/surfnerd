package wavewatch

const (
	BASE_MULTIGRID_URL = "http://nomads.ncep.noaa.gov:9090/dods/wave/mww3/%s/%s%s_%s.ascii?dirpwsfc.dirpwsfc[0:60][%d][%d],htsgwsfc.htsgwsfc[0:60][%d][%d],perpwsfc.perpwsfc[0:60][%d][%d],swdir_1.swdir_1[0:60][%d][%d],swdir_2.swdir_2[0:60][%d][%d],swell_1.swell_1[0:60][%d][%d],swell_2.swell_2[0:60][%d][%d],swper_1.swper_1[0:60][%d][%d],swper_2.swper_2[0:60][%d][%d],ugrdsfc.ugrdsfc[0:60][%d][%d],vgrdsfc.vgrdsfc[0:60][%d][%d],wdirsfc.wdirsfc[0:60][%d][%d],windsfc.windsfc[0:60][%d][%d],wvdirsfc.wvdirsfc[0:60][%d][%d],wvhgtsfc.wvhgtsfc[0:60][%d][%d],wvpersfc.wvpersfc[0:60][%d][%d]"
)

func fetchRawWaveWatchData(loc *Location) {
	eastCoastModel := EastCoastModel{}
	if !eastCoastModel.ContainsLocation(loc) {
		return
	}

	// TODO: Find the latitude and longitde indexes

	// TODO: Format the data url

	// TODO: Fetch the raw data

	// TODO: Call to parse the raw data into containers
}
