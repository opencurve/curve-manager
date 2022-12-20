/*
prometheus http api resonse data struct

{
  "status": "success" | "error",
  "data": <data>,

  // Only set if status is "error". The data field may still hold additional data.
  "errorType": "<string>",
  "error": "<string>"
}
*/

package metrics

const (
	VECTOR_METRIC_PATH = "/api/v1/query"
	VECTOR_METRIC_QUERY_KEY = "query"

	// etcd
	ETCD_CLUSTER_VERSION_NAME  = "etcd_cluster_version"
	ETCD_SERVER_IS_LEADER_NAME = "etcd_server_is_leader"

	ETCD_CLUSTER_VERSION  = "cluster_version"

)

// data struct: https://prometheus.io/docs/prometheus/latest/querying/api/#expression-query-result-formats
type QueryResponseOfVector struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type EtcdStatus struct {
	Version  string   `json:"version"`
	Leader   []string `json:"leader"`
	Follower []string `json:"follower"`
}
