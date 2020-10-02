package aqi

const (
	Good                        = Concern("Good")
	Moderate                    = Concern("Moderate")
	UnhealthyForSensitiveGroups = Concern("Unhealthy For Sensitive Groups")
	Unhealthy                   = Concern("Unhealthy")
	VeryUnhealthy               = Concern("Very Unhealthy")
	Hazardous                   = Concern("Hazardous")
)

type Concern string

func (c Concern) String() string {
	return string(c)
}

type Index struct {
	Concern             Concern
	ILow, IHigh         float32
	PM25CLow, PM25CHigh float32
}

func (i Index) InRange(pm25_ugm3 float32) bool {
	return (i.PM25CLow <= pm25_ugm3) && (i.PM25CHigh >= pm25_ugm3)
}

func (i Index) Coefficient() float32 {
	return ((i.IHigh - i.ILow) / (i.PM25CHigh - i.PM25CLow))
}

var Table = map[Concern]Index{
	Good:                        Index{Good, float32(0), float32(50.0), float32(0), float32(12.0)},
	Moderate:                    Index{Moderate, float32(51), float32(100), float32(12.1), float32(35.4)},
	UnhealthyForSensitiveGroups: Index{UnhealthyForSensitiveGroups, float32(101), float32(150), float32(35.5), float32(55.4)},
	Unhealthy:                   Index{Unhealthy, float32(151), float32(200), float32(55.5), float32(150.4)},
	VeryUnhealthy:               Index{VeryUnhealthy, float32(201), float32(300), float32(150.5), float32(250.4)},
	Hazardous:                   Index{Hazardous, float32(301), float32(500), float32(250.5), float32(500.4)},
}

// I returns the AQI index and the concern string
func I(pm25_ugm3 float32) (float32, string) {
	var category Index

	for idx, i := range Table {
		if i.InRange(pm25_ugm3) {
			category = Table[idx]
		}
	}

	return (category.Coefficient()*(pm25_ugm3-category.PM25CLow) + category.ILow), category.Concern.String()
}
