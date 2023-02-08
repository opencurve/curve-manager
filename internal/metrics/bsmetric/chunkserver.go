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

func GetCopysetRaftStatus(endpoints *[]string) (map[string][]map[string]string, error) {
	size := len(*endpoints)
	results := make(chan comm.MetricResult, size)
	comm.GetRaftStatusMetric(*endpoints, &results)

	count := 0
	// key: chunkserver's addr, value: copysets' raft status
	ret := map[string][]map[string]string{}
	for res := range results {
		if res.Err == nil {
			v, e := comm.ParseRaftStatusMetric(res.Result.(string))
			if e != nil {
				return nil, e
			} else {
				ret[res.Key.(string)] = v
			}
		} else {
			ret[res.Key.(string)] = nil
		}
		count += 1
		if count >= size {
			break
		}
	}
	return ret, nil
}
