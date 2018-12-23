package lib

type Datum struct {
	Host string
	T    struct {
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

func (d *Datum) Map() map[string]interface{} {
	m := make(map[string]float32)

	m["humidity"] = d.T.Humidity
	m["temperature"] = d.T.Temp

	samples := float32(len(d.A))

	for idx := 0; idx < len(d.A); idx++ {
		m["SP_1_0"] += d.A[idx].SP_1_0
		m["SP_2_5"] += d.A[idx].SP_2_5
		m["SP_10_0"] += d.A[idx].SP_10_0
		m["AE_1_0"] += d.A[idx].AE_1_0
		m["AE_2_5"] += d.A[idx].AE_2_5
		m["AE_10_0"] += d.A[idx].AE_10_0
	}

	m["SP_1_0"] = m["SP_1_0"] / samples
	m["SP_2_5"] = m["SP_2_5"] / samples
	m["SP_10_0"] = m["SP_10_0"] / samples

	m["AE_1_0"] = m["AE_1_0"] / samples
	m["AE_2_5"] = m["AE_2_5"] / samples
	m["AE_10_0"] = m["AE_10_0"] / samples

	p := make(map[string]interface{})

	for k, v := range m {
		p[k] = v
	}
	return p
}
