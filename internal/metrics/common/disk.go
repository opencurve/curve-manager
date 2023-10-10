/*
*  Copyright (c) 2023 NetEase Inc.
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
 */

/*
* Project: Curve-Manager
* Created Date: 2023-02-15
* Author: wanghai (SeanHai)
 */

package common

import (
	"fmt"
	"strconv"
	"strings"

	comm "github.com/opencurve/curve-manager/internal/common"
)

const (
	// disk info
	NODE_DISK_INFO                  = "node_disk_info"
	NODE_DISK_READ_COMPLETED_TOTAL  = "node_disk_reads_completed_total"
	NODE_DISK_WRITE_COMPLETED_TOTAL = "node_disk_writes_completed_total"
	NODE_DISK_READ_BYTES_TOTAL      = "node_disk_read_bytes_total"
	NODE_DISK_WRITTEN_BYTES_TOTAL   = "node_disk_written_bytes_total"

	// filesystem
	NODE_FILESYSTEM_SIZE_TOTAL = "node_filesystem_size_bytes"
	NODE_FILESYSTEM_SIZE_AVAIL = "node_filesystem_avail_bytes"

	// disk
	NODE_DISK_ROTATION_RATE      = "node_disk_ata_rotation_rate_rpm"
	NODE_DISK_WRITE_CACHE_ENABLE = "node_disk_ata_write_cache_enabled"
	SSD_TYPE                     = "SSD"
	HDD_TYPE                     = "HDD"

	// disk modle
	MODLE = "model"
)

// @return map[string][]Performance, key: device, value: performance at different timestamp
func GetDiskPerformance(instance string, start, end, interval uint64) (interface{}, error) {
	performance := make(map[string][]Performance)
	retMap := make(map[string]map[float64]*Performance)

	// writeIOPS, writeBPS, readIOPS, readBPS
	requestSize := 4
	results := make(chan MetricResult, requestSize)
	writeIOPSName := fmt.Sprintf("%s&start=%d&end=%d&step=%ds",
		GetNodeDiskPerformanceName(NODE_DISK_WRITE_COMPLETED_TOTAL, instance, interval), start, end, interval)
	writeBPSName := fmt.Sprintf("%s&start=%d&end=%d&step=%ds",
		GetNodeDiskPerformanceName(NODE_DISK_WRITTEN_BYTES_TOTAL, instance, interval), start, end, interval)
	readIOPSName := fmt.Sprintf("%s&start=%d&end=%d&step=%ds",
		GetNodeDiskPerformanceName(NODE_DISK_READ_COMPLETED_TOTAL, instance, interval), start, end, interval)
	readBPSName := fmt.Sprintf("%s&start=%d&end=%d&step=%ds",
		GetNodeDiskPerformanceName(NODE_DISK_READ_BYTES_TOTAL, instance, interval), start, end, interval)

	go QueryRangeMetric(writeIOPSName,start,end,interval, &results)
	go QueryRangeMetric(writeBPSName,start,end,interval, &results)
	go QueryRangeMetric(readIOPSName,start,end,interval, &results)
	go QueryRangeMetric(readBPSName,start,end,interval, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := ParseMatrixMetric(res.Result.(*QueryResponseOfMatrix), DEVICE)
		for device, v := range ret {
			for _, data := range v {
				if _, ok := retMap[device]; ok {
					if p, ok := retMap[device][data.Timestamp]; ok {
						switch res.Key.(string) {
						case writeIOPSName:
							p.WriteIOPS = data.Value
						case writeBPSName:
							p.WriteBPS = data.Value
						case readIOPSName:
							p.ReadIOPS = data.Value
						case readBPSName:
							p.ReadBPS = data.Value
						}
					} else {
						switch res.Key.(string) {
						case writeIOPSName:
							retMap[device][data.Timestamp] = &Performance{
								Timestamp: data.Timestamp,
								WriteIOPS: data.Value,
							}
						case writeBPSName:
							retMap[device][data.Timestamp] = &Performance{
								Timestamp: data.Timestamp,
								WriteBPS:  data.Value,
							}
						case readIOPSName:
							retMap[device][data.Timestamp] = &Performance{
								Timestamp: data.Timestamp,
								ReadIOPS:  data.Value,
							}
						case readBPSName:
							retMap[device][data.Timestamp] = &Performance{
								Timestamp: data.Timestamp,
								ReadBPS:   data.Value,
							}
						}
					}
				} else {
					retMap[device] = make(map[float64]*Performance)
				}
			}
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	for dev, v := range retMap {
		var tmparr []Performance
		for _, p := range v {
			tmparr = append(tmparr, *p)
		}
		performance[dev] = tmparr
	}
	return performance, nil
}

// @return map[instance][]map[key]value
func ListDiskInfo(instance string) (interface{}, error) {
	disks := make(map[string][]map[string]string)
	requestSize := 1
	results := make(chan MetricResult, requestSize)
	metricName := NODE_DISK_INFO
	if instance != "" {
		metricName = fmt.Sprintf("%s{instance=%q}", metricName, instance)
	}
	QueryInstantMetric(metricName, &results)
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := res.Result.(*QueryResponseOfVector)
		for _, item := range ret.Data.Result {
			if _, ok := item.Metric[MODLE]; ok {
				disks[item.Metric[INSTANCE]] = append(disks[item.Metric[INSTANCE]], item.Metric)
			}
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	return disks, nil
}

// /dev/sda -> sda
func getShortDeviceName(name string) string {
	prefix := "/dev/"
	if !strings.HasPrefix(name, prefix) {
		return name
	}
	return strings.ReplaceAll(name, prefix, "")
}

// @return ap[instance]map[device]FileSystemInfo
func GetDiskFileSystemInfo(instance string) (interface{}, error) {
	fileSystemInfos := make(map[string]map[string]FileSystemInfo)
	// totalSpace, freeSpace
	requestSize := 2
	results := make(chan MetricResult, requestSize)
	totalName := NODE_FILESYSTEM_SIZE_TOTAL
	availName := NODE_FILESYSTEM_SIZE_AVAIL
	if instance != "" {
		totalName = fmt.Sprintf("%s{instance=%q}", totalName, instance)
		availName = fmt.Sprintf("%s{instance=%q}", availName, instance)
	}
	go QueryInstantMetric(totalName, &results)
	go QueryInstantMetric(availName, &results)
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := res.Result.(*QueryResponseOfVector)
		for _, item := range ret.Data.Result {
			inst := item.Metric[INSTANCE]
			// ListDisk return the short device name, and this info will combine with disk info at uplayer
			dev := getShortDeviceName(item.Metric[DEVICE])
			if _, ok := fileSystemInfos[inst]; !ok {
				fileSystemInfos[inst] = make(map[string]FileSystemInfo)
			}
			space, _ := strconv.ParseUint(item.Value[1].(string), 10, 64)
			if info, ok := fileSystemInfos[inst][dev]; ok {
				switch res.Key.(string) {
				case totalName:
					info.SpaceTotal = space / comm.GiB
				case availName:
					info.SpaceAvail = space / comm.GiB
				}
				fileSystemInfos[inst][dev] = info
			} else {
				info := FileSystemInfo{}
				info.Device = dev
				info.FsType = item.Metric[FSTYPE]
				info.MountPoint = item.Metric[MOUNTPOINT]
				switch res.Key.(string) {
				case totalName:
					info.SpaceTotal = space / comm.GiB
				case availName:
					info.SpaceAvail = space / comm.GiB
				}
				fileSystemInfos[inst][dev] = info
			}
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	return fileSystemInfos, nil
}

func getDiskTypeByRotationRate(rate string) string {
	switch rate {
	case "0":
		return SSD_TYPE
	default:
		return HDD_TYPE
	}
}

// @return map[instance]map[device]type
func GetDiskType(instance string) (interface{}, error) {
	diskTypes := make(map[string]map[string]string)
	requestSize := 1
	results := make(chan MetricResult, requestSize)
	metricName := NODE_DISK_ROTATION_RATE
	if instance != "" {
		metricName = fmt.Sprintf("%s{instance=%q}", metricName, instance)
	}
	go QueryInstantMetric(metricName, &results)
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := res.Result.(*QueryResponseOfVector)
		for _, item := range ret.Data.Result {
			inst := item.Metric[INSTANCE]
			dev := item.Metric[DEVICE]
			if _, ok := diskTypes[inst]; !ok {
				diskTypes[inst] = make(map[string]string)
			}
			diskTypes[inst][dev] = getDiskTypeByRotationRate(item.Value[1].(string))
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	return diskTypes, nil
}

// @return map[instance]map[device]writeCacheEnable
func GetDiskWriteCacheEnableFlag(instance string) (interface{}, error) {
	diskWriteCaches := make(map[string]map[string]string)
	requestSize := 1
	results := make(chan MetricResult, requestSize)
	metricName := NODE_DISK_WRITE_CACHE_ENABLE
	if instance != "" {
		metricName = fmt.Sprintf("%s{instance=%q}", metricName, instance)
	}
	go QueryInstantMetric(metricName, &results)
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := res.Result.(*QueryResponseOfVector)
		for _, item := range ret.Data.Result {
			inst := item.Metric[INSTANCE]
			dev := item.Metric[DEVICE]
			if _, ok := diskWriteCaches[inst]; !ok {
				diskWriteCaches[inst] = make(map[string]string)
			}
			diskWriteCaches[inst][dev] = item.Value[1].(string)
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	return diskWriteCaches, nil
}
