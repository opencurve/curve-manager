package bsmetric

import comm "github.com/opencurve/curve-manager/internal/metrics/common"

func GetChunkServerVersion(endpoints *[]string) (*map[string]int, error) {
	size := len(*endpoints)
	results := make(chan comm.MetricResult, size)
	comm.GetBvarMetric(*endpoints, comm.CURVEBS_VERSION, &results)

	count := 0
	ret := make(map[string]int)
	for res := range results {
		if res.Err == nil {
			v, e := comm.ParseBvarMetric(res.Result.(string))
			if e != nil {
				return nil, e
			} else {
				ret[(*v)[comm.CURVEBS_VERSION]] += 1
			}
		} else {
			return nil, res.Err
		}
		count += 1
		if count >= size {
			break
		}
	}
	return &ret, nil
}
