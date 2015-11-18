package wavewatch

const (
	baseMultigridUrl = "http://nomads.ncep.noaa.gov:9090/dods/wave/mww3/%[1]s/%[2]s%[1]s_%[3]s.ascii?dirpwsfc.dirpwsfc[0:60][%[4]d][%[5]d],htsgwsfc.htsgwsfc[0:60][%[4]d][%[5]d],perpwsfc.perpwsfc[0:60][%[4]d][%[5]d],swdir_1.swdir_1[0:60][%[4]d][%[5]d],swdir_2.swdir_2[0:60][%[4]d][%[5]d],swell_1.swell_1[0:60][%[4]d][%[5]d],swell_2.swell_2[0:60][%[4]d][%[5]d],swper_1.swper_1[0:60][%[4]d][%[5]d],swper_2.swper_2[0:60][%[4]d][%[5]d],ugrdsfc.ugrdsfc[0:60][%[4]d][%[5]d],vgrdsfc.vgrdsfc[0:60][%[4]d][%[5]d],wdirsfc.wdirsfc[0:60][%[4]d][%[5]d],windsfc.windsfc[0:60][%[4]d][%[5]d],wvdirsfc.wvdirsfc[0:60][%[4]d][%[5]d],wvhgtsfc.wvhgtsfc[0:60][%[4]d][%[5]d],wvpersfc.wvpersfc[0:60][%[4]d][%[5]d]"
)

type WaveModel interface {
	Name() string
	Description() string
	ContainsLocation(loc *Location) bool
	LocationResolution() float64
	TimeResolution() float64
	LocationIndices(loc *Location) (int, int)
}
