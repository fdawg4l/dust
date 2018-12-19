package lib

type Datum struct {
	T struct {
		Humidity float32 `json:"humidity_P"`
		Temp     float32 `json:"temp_F"`
	}
	A []struct {
		SP_1_0  float32
		SP_2_5  float32
		SP_10_0 float32
		AE_1_0  float32
		AE_2_5  float32
		AE_10_0 float32
	}
}
