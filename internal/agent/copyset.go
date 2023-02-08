package agent

import (
	"fmt"
	"strconv"
	"strings"

	set "github.com/deckarep/golang-set/v2"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
)

const (
	// copyset status check result
	COPYSET_CHECK_HEALTHY               = "healthy"
	COPYSET_CHECK_PARSE_ERROR           = "parse_error"
	COPYSET_CHECK_NO_LEADER             = "no leader"
	COPYSET_CHECK_PEERS_NO_SUFFICIENT   = "peers_no_sufficient"
	COPYSET_CHECK_LOG_INDEX_TOO_BIG     = "log_index_gap_too_big"
	COPYSET_CHECK_INSTALLING_SNAPSHOT   = "installing_snapshot"
	COPYSET_CHECK_MINORITY_PEER_OFFLINE = "minority_peer_offline"
	COPYSET_CHECK_MAJORITY_PEER_OFFLINE = "majority_peer_offline"
	COPYSET_CHECK_INCONSISTENT          = "Three copies inconsistent"
	COPYSET_CHECK_OTHER_ERROR           = "other_error"
	COPYSET_TOTAL                       = "total"
)

type Copyset struct {
	// recored copysets in all status; key is status, value is groupIds
	copyset map[string]set.Set[string]
	// record chunkservers which send rpc failed
	serviceExceptionChunkservers set.Set[string]
	// record chunkservers which load copyset error
	copysetLoacExceptionChunkServers set.Set[string]
	// record copysets which belongs to different chunkservers, chunkserverEndpoint <-> copysets
	chunkServerCopysets map[string]set.Set[string]
}

func NewCopyset() *Copyset {
	var cs Copyset
	cs.copyset = make(map[string]set.Set[string])
	cs.serviceExceptionChunkservers = set.NewSet[string]()
	cs.copysetLoacExceptionChunkServers = set.NewSet[string]()
	cs.chunkServerCopysets = make(map[string]set.Set[string])
	return &cs
}

func getGroupId(poolId, copysetId uint32) uint64 {
	return uint64(poolId)<<32 | uint64(copysetId)
}

func getPoolIdFormGroupId(groupId uint64) uint32 {
	return uint32(groupId >> 32)
}

func getCopysetIdFromGroupId(groupId uint64) uint32 {
	return uint32(groupId & (uint64(1)<<32 - 1))
}

func (cs *Copyset) recordCopyset(s, groupId string) {
	if _, ok := cs.copyset[s]; !ok {
		cs.copyset[s] = set.NewSet[string]()
	}
	cs.copyset[s].Add(groupId)
}

func (cs *Copyset) checkCopysetOnline(addr, groupId string) bool {
	copysets := cs.chunkServerCopysets[addr]
	if copysets != nil {
		return copysets.Contains(groupId)
	}
	return false
}

func (cs *Copyset) checkPeerOnlineStatus(groupId string, peers []string) string {
	offlineNum := 0
	for _, peer := range peers {
		addr := peer[0:strings.LastIndex(peer, ":")]
		online := cs.checkCopysetOnline(addr, groupId)
		if !online {
			offlineNum += 1
		}
	}
	if offlineNum > 0 {
		if offlineNum < len(peers)/2+1 {
			return COPYSET_CHECK_MINORITY_PEER_OFFLINE
		} else {
			return COPYSET_CHECK_MAJORITY_PEER_OFFLINE
		}
	}
	return COPYSET_CHECK_HEALTHY
}

func (cs *Copyset) updatePeerOfflineCopysets(csAddr string) error {
	// get copysets on offline chunkserver
	item := strings.Split(csAddr, ":")
	port, err := strconv.Atoi(item[1])
	if err != nil {
		return err
	}
	copysets, err := bsrpc.GMdsClient.GetCopySetsInChunkServer(item[0], uint32(port))
	if err != nil {
		return fmt.Errorf("GetCopySetsInChunkServer failed, %s", err)
	}
	// get copyset's all members
	var logicalPoolId uint32
	var copysetIds []uint32
	for _, copyset := range copysets {
		copysetIds = append(copysetIds, copyset.CopysetId)
	}
	if len(copysets) > 0 {
		logicalPoolId = copysets[0].LogicalPoolId
	}
	memberInfo, err := bsrpc.GMdsClient.GetChunkServerListInCopySets(logicalPoolId, copysetIds)
	if err != nil {
		return fmt.Errorf("GetChunkServerListInCopySets failed, %s", err)
	}
	// check all members
	for _, info := range memberInfo {
		var peers []string
		for _, loc := range info.CsLocs {
			endpoint := fmt.Sprintf("%s:%d:0", loc.HostIp, loc.Port)
			peers = append(peers, endpoint)
		}
		groupId := strconv.FormatUint(getGroupId(logicalPoolId, info.CopysetId), 10)
		ret := cs.checkPeerOnlineStatus(groupId, peers)
		if err != nil {
			return err
		}
		switch ret {
		case COPYSET_CHECK_MINORITY_PEER_OFFLINE:
			cs.recordCopyset(COPYSET_CHECK_MINORITY_PEER_OFFLINE, groupId)
		case COPYSET_CHECK_MAJORITY_PEER_OFFLINE:
			cs.recordCopyset(COPYSET_CHECK_MAJORITY_PEER_OFFLINE, groupId)
		default:
			cs.recordCopyset(COPYSET_TOTAL, groupId)
		}
	}
	return nil
}

func (cs *Copyset) updateChunkServerCopysets(csAddr string, status []map[string]string) {
	copysetGroupIds := set.NewSet[string]()
	if status != nil {
		for _, s := range status {
			copysetGroupIds.Add(s[common.RAFT_STATUS_KEY_GROUPID])
		}
	}
	cs.chunkServerCopysets[csAddr] = copysetGroupIds
}

func (cs *Copyset) ifChunkServerInCopysets(csAddr string, groupIds *set.Set[string]) (map[string]bool, error) {
	var logicalPoolId uint32
	var copysetIds []uint32
	result := make(map[string]bool)
	for gid := range (*groupIds).Iter() {
		ngid, err := strconv.ParseUint(gid, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse string groupId to uint64 error: %s", gid)
		}
		logicalPoolId = getPoolIdFormGroupId(ngid)
		copysetIds = append(copysetIds, getCopysetIdFromGroupId(ngid))
	}
	memberInfo, err := bsrpc.GMdsClient.GetChunkServerListInCopySets(logicalPoolId, copysetIds)
	if err != nil {
		return nil, fmt.Errorf("GetChunkServerListInCopySets failed, %s", err)
	}
	for _, info := range memberInfo {
		csId := info.CopysetId
		gId := getGroupId(logicalPoolId, csId)
		for _, csloc := range info.CsLocs {
			addr := fmt.Sprintf("%s:%d", csloc.HostIp, csloc.Port)
			if csAddr == addr {
				result[strconv.FormatUint(gId, 10)] = true
				break
			}
		}
	}
	return result, nil
}

func (cs *Copyset) checkCopysetsNoLeader(csAddr string, peersMap *map[string][]string) bool {
	healthy := true
	if len(*peersMap) == 0 {
		return true
	}
	groupIds := set.NewSet[string]()
	for k, _ := range *peersMap {
		groupIds.Add(k)
	}
	result, err := cs.ifChunkServerInCopysets(csAddr, &groupIds)
	if err != nil {
		return false
	}

	for k, v := range result {
		if v {
			healthy = false
			ret := cs.checkPeerOnlineStatus(k, (*peersMap)[k])
			if ret == COPYSET_CHECK_MAJORITY_PEER_OFFLINE {
				cs.recordCopyset(COPYSET_CHECK_MAJORITY_PEER_OFFLINE, k)
				continue
			}
			cs.recordCopyset(COPYSET_CHECK_NO_LEADER, k)
		}
	}
	return healthy
}

func (cs *Copyset) checkHealthOnLeader(raftStatus *map[string]string) string {
	// 1. check peers number
	peers, ok := (*raftStatus)[common.RAFT_STATUS_KEY_PEERS]
	if !ok {
		return COPYSET_CHECK_PARSE_ERROR
	}
	peerArray := strings.Split(peers, " ")
	if len(peerArray) < common.RAFT_REPLICAS_NUMBER {
		return COPYSET_CHECK_PEERS_NO_SUFFICIENT
	}

	// 2. check offline peers
	groupId := (*raftStatus)[common.RAFT_STATUS_KEY_GROUPID]
	ret := cs.checkPeerOnlineStatus(groupId, peerArray)
	if ret != COPYSET_CHECK_HEALTHY {
		return ret
	}

	// 3. check log index gap between replicas
	strLastLogId := (*raftStatus)[common.RAFT_STATUS_KEY_LAST_LOG_ID]
	lastLogId, err := strconv.ParseUint(strLastLogId[strings.Index(strLastLogId, "=")+1:strings.Index(strLastLogId, ",")], 10, 64)
	if err != nil {
		return COPYSET_CHECK_PARSE_ERROR
	}
	var gap, nextIndex, flying uint64
	gap = 0
	nextIndex = 0
	flying = 0
	for i, size := 0, len(peerArray)-1; i < size; i++ {
		key := fmt.Sprintf("%s%d", common.RAFT_STATUS_KEY_REPLICATOR, i)
		repInfos := strings.Split((*raftStatus)[key], " ")
		for _, info := range repInfos {
			if !strings.Contains(info, "=") {
				if strings.Contains(info, common.RAFT_STATUS_KEY_SNAPSHOT) {
					return COPYSET_CHECK_INSTALLING_SNAPSHOT
				} else {
					continue
				}
			}
			pos := strings.Index(info, "=")
			if info[0:pos] == common.RAFT_STATUS_KEY_NEXT_INDEX {
				nextIndex, err = strconv.ParseUint(info[pos+1:], 10, 64)
				if err != nil {
					return COPYSET_CHECK_PARSE_ERROR
				}
			} else if info[0:pos] == common.RAFT_STATUS_KEY_FLYING_APPEND_ENTRIES_SIZE {
				flying, err = strconv.ParseUint(info[pos+1:], 10, 64)
				if err != nil {
					return COPYSET_CHECK_PARSE_ERROR
				}
			}
			gap = common.Max(gap, lastLogId-(nextIndex-1-flying))
		}
	}
	if gap > common.RAFT_MARGIN {
		return COPYSET_CHECK_LOG_INDEX_TOO_BIG
	}
	return COPYSET_CHECK_HEALTHY
}

func (cs *Copyset) checkCopysetsOnChunkServer(csAddr string, status []map[string]string) bool {
	healthy := true
	// query copysets' status failed on chunkserver and think it offline
	if status == nil {
		cs.updatePeerOfflineCopysets(csAddr)
		cs.serviceExceptionChunkservers.Add(csAddr)
		return false
	}

	if len(status) == 0 {
		cs.copysetLoacExceptionChunkServers.Add(csAddr)
		return false
	}

	noLeaderCopysetPeers := make(map[string][]string)
	for _, statMap := range status {
		groupId := statMap[common.RAFT_STATUS_KEY_GROUPID]
		state := statMap[common.RAFT_STATUS_KEY_STATE]
		cs.recordCopyset(COPYSET_TOTAL, groupId)
		if state == common.RAFT_STATUS_STATE_LEADER {
			ret := cs.checkHealthOnLeader(&statMap)
			switch ret {
			case COPYSET_CHECK_PEERS_NO_SUFFICIENT:
				cs.recordCopyset(COPYSET_CHECK_PEERS_NO_SUFFICIENT, groupId)
				healthy = false
				break
			case COPYSET_CHECK_LOG_INDEX_TOO_BIG:
				cs.recordCopyset(COPYSET_CHECK_LOG_INDEX_TOO_BIG, groupId)
				healthy = false
				break
			case COPYSET_CHECK_INSTALLING_SNAPSHOT:
				cs.recordCopyset(COPYSET_CHECK_INSTALLING_SNAPSHOT, groupId)
				healthy = false
				break
			case COPYSET_CHECK_MINORITY_PEER_OFFLINE:
				cs.recordCopyset(COPYSET_CHECK_MINORITY_PEER_OFFLINE, groupId)
				healthy = false
				break
			case COPYSET_CHECK_MAJORITY_PEER_OFFLINE:
				cs.recordCopyset(COPYSET_CHECK_MAJORITY_PEER_OFFLINE, groupId)
				healthy = false
				break
			case COPYSET_CHECK_PARSE_ERROR:
				cs.recordCopyset(COPYSET_CHECK_PARSE_ERROR, groupId)
				healthy = false
				break
			default:
				break
			}
		} else if state == common.RAFT_STATUS_STATE_FOLLOWER {
			v, ok := statMap[common.RAFT_STATUS_KEY_LEADER]
			if !ok || v == common.RAFT_EMPTY_ADDR {
				noLeaderCopysetPeers[groupId] = strings.Split(statMap[common.RAFT_STATUS_KEY_PEERS], " ")
				continue
			}
		} else if state == common.RAFT_STATUS_STATE_TRANSFERRING || state == common.RAFT_STATUS_STATE_CANDIDATE {
			cs.recordCopyset(COPYSET_CHECK_NO_LEADER, groupId)
			healthy = false
		} else {
			// other state: ERROR,UNINITIALIZED,SHUTTING和SHUTDOWN
			key := fmt.Sprintf("state %s", state)
			cs.recordCopyset(key, groupId)
			healthy = false
		}
	}

	health := cs.checkCopysetsNoLeader(csAddr, &noLeaderCopysetPeers)
	if !health {
		healthy = false
	}
	return healthy
}

func (cs *Copyset) checkCopysetsWithMds() (bool, error) {
	// get copysets in cluster
	csInfos, err := bsrpc.GMdsClient.GetCopySetsInCluster()
	if err != nil {
		return false, fmt.Errorf("GetCopySetsInCluster failed, %s", err)
	}

	// check copyset number
	if len(csInfos) != cs.copyset[COPYSET_TOTAL].Cardinality() {
		return false, fmt.Errorf("Copyset numbers in chunkservers not consistent with mds,"+
			"please check! copysets on chunkserver: %d; copysets in mds: %d",
			cs.copyset[COPYSET_TOTAL].Cardinality(), len(csInfos))
	}

	// check copyset groupId difference
	groupIdsInMds := set.NewSet[string]()
	for _, info := range csInfos {
		groupIdsInMds.Add(strconv.FormatUint(getGroupId(info.LogicalPoolId, info.CopysetId), 10))
	}
	notInChunkServer := groupIdsInMds.Difference(cs.copyset[COPYSET_TOTAL])
	if notInChunkServer.Cardinality() != 0 {
		return false, fmt.Errorf("some copysets in mds but not in chunkserver: %v", notInChunkServer)
	}
	notInMds := cs.copyset[COPYSET_TOTAL].Difference(groupIdsInMds)
	if notInMds.Cardinality() != 0 {
		return false, fmt.Errorf("some copysets in chunkserver but not in mds: %v", notInMds)
	}

	// check copyset data consistency scanned result
	count := 0
	for _, info := range csInfos {
		if info.LastScanSec == 0 || (info.LastScanSec != 0 && info.LastScanConsistent) {
			continue
		}
		groupId := strconv.FormatUint(getGroupId(info.LogicalPoolId, info.CopysetId), 10)
		cs.recordCopyset(COPYSET_CHECK_INCONSISTENT, groupId)
		count++
	}
	if count > 0 {
		return false, fmt.Errorf("There are %d inconsistent copyset", count)
	}
	return true, nil
}

func (cs *Copyset) checkCopysetsInCluster() (bool, error) {
	healthy := true
	// 2.1 get chunkservers in cluster
	chunkservers, err := bsrpc.GMdsClient.GetChunkServerInCluster()
	if err != nil {
		return false, fmt.Errorf("GetChunkServerInCluster failed, %s", err)
	}

	// 2.2 get copyset raft status
	var csAddrs []string
	for _, cs := range chunkservers {
		csAddrs = append(csAddrs, fmt.Sprintf("%s:%d", cs.HostIp, cs.Port))
	}
	status, err := bsmetric.GetCopysetRaftStatus(&csAddrs)
	if err != nil {
		return false, err
	}
	for k, v := range status {
		cs.updateChunkServerCopysets(k, v)
	}

	// 2.3 check copyset status on chunkserver
	for k, v := range status {
		if !cs.checkCopysetsOnChunkServer(k, v) {
			healthy = false
		}
	}

	// 2.4 check copysets queried from chunkserver with mds
	health, err := cs.checkCopysetsWithMds()
	if err != nil {
		return false, err
	}
	if !health {
		healthy = false
	}

	// TODO: check operator on mds if needed
	return healthy, nil
}

func (cs *Copyset) getCopysetTotalNum() uint32 {
	if total, ok := cs.copyset[COPYSET_TOTAL]; ok {
		return uint32(total.Cardinality())
	}
	return 0
}

func (cs *Copyset) getCopysetUnhealthyNum() uint32 {
	var number uint32
	for k, v := range cs.copyset {
		if k != COPYSET_TOTAL {
			number += uint32(v.Cardinality())
		}
	}
	return number
}