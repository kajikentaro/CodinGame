package main

import (
	"fmt"
	"os"
	"sort"
)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/
type Site struct {
	siteId                                                                      SiteId
	x, y, radius, gold, maxMineSize, structureType, owner, param1, param2, dist int
}

type SiteId int

type BuildOrder struct {
	siteId        SiteId
	structureType StructureType
}

type StructureType string

const (
	BARRACKS_KNIGHT = StructureType("BARRACKS-KNIGHT")
	BARRACKS_ARCHER = StructureType("BARRACKS-ARCHER")
	BARRACKS_GIANT  = StructureType("BARRACKS-GIANT")
	TOWER           = StructureType("TOWER")
	MINE            = StructureType("MINE")
)

type Unit struct {
	x, y, owner, unitType, health int
}

func distUnit(a Unit, b Unit) int {
	return (a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y)
}

func distSite(a Site, b Site) int {
	return (a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y)
}

func dist(a Site, b Unit) int {
	return (a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y)
}

func updateSiteList(numSites int, idToSite []Site, queen Unit, tmpList []Site) []Site {
	for si := 0; si < numSites; si++ {
		for ti := 0; ti < numSites; ti++ {
			if tmpList[ti].siteId == idToSite[si].siteId {
				// NOTE: バグ？goldが0になったり-1になったりする
				if idToSite[si].gold != 0 {
					idToSite[si].gold = tmpList[ti].gold
				}
				idToSite[si].maxMineSize = tmpList[ti].maxMineSize
				idToSite[si].structureType = tmpList[ti].structureType
				idToSite[si].owner = tmpList[ti].owner
				idToSite[si].param1 = tmpList[ti].param1
				idToSite[si].param2 = tmpList[ti].param2
			}
		}
	}
	return idToSite
}

func inputSiteList(numSites int) []Site {
	tmpList := make([]Site, numSites)
	for i := 0; i < numSites; i++ {
		// gold: 発掘できる残りの金
		// maxMineSize: used in future leagues
		// structureType: -1 = No structure, 2 = Barracks
		// owner: -1 = No structure, 0 = Friendly, 1 = Enemy
		var siteId SiteId
		fmt.Scan(&siteId)
		site := &tmpList[siteId]
		site.siteId = siteId
		fmt.Scan(&site.gold, &site.maxMineSize, &site.structureType, &site.owner, &site.param1, &site.param2)
	}
	return tmpList
}

func getUnitList(numUnits int) ([]Unit, Unit) {
	var queen Unit

	unitList := make([]Unit, numUnits)
	for i := 0; i < numUnits; i++ {
		// unitType: -1 = QUEEN, 0 = KNIGHT, 1 = ARCHER
		unit := &unitList[i]
		fmt.Scan(&unit.x, &unit.y, &unit.owner, &unit.unitType, &unit.health)
		if unit.owner == 0 && unit.unitType == -1 {
			queen = *unit
		}
	}
	return unitList, queen
}

func pow(a, b int) int {
	if b == 0 {
		return 1
	}
	res := 1
	if b%2 == 1 {
		res *= a
		b--
	}
	k := pow(a, b/2)
	res *= k
	res *= k
	return res
}

func travelingSalesman(nearSiteList []*Site) []SiteId {
	n := len(nearSiteList)
	s := 0

	type P struct {
		first, second int
	}

	path := make([][]P, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			p := P{distSite(*nearSiteList[i], *nearSiteList[j]), j}
			path[i] = append(path[i], p)
		}
	}

	inf := int(1e18)

	dp := make([][]int, pow(2, n))
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = inf
		}
	}

	dp_log := make([][][]int, pow(2, n))
	for i := range dp_log {
		dp_log[i] = make([][]int, n)
		for j := range dp_log[i] {
			dp_log[i][j] = []int{}
		}
	}
	for _, i := range path[s] {
		w := i.first
		next := i.second
		dp[1<<next][next] = w
		dp_log[1<<next][next] = []int{s, next}
	}
	for i := 0; i < pow(2, n); i++ {
		for j := 0; j < n; j++ {
			if ((1 << j) & i) == 0 {
				continue
			}
			for _, tmp := range path[j] {
				w := tmp.first
				next := tmp.second
				next_status := i | (1 << next)
				if next_status == i {
					continue
				}
				if dp[next_status][next] > w+dp[i][j] {
					dp[next_status][next] = w + dp[i][j]
					dp_log[next_status][next] = make([]int, len(dp_log[i][j])+1)
					copy(dp_log[next_status][next], append(dp_log[i][j], next))
				}
			}
		}
	}

	siteIdList := []SiteId{}
	for _, val := range dp_log[pow(2, int(n))-1][s] {
		siteIdList = append(siteIdList, nearSiteList[val].siteId)
	}
	return siteIdList
}

// 末尾が新しい
var trainingHistory []int

func calcTrainingSite(idToSite []Site, gold int) []Site {
	VARRACKS_TYPE := 3
	var isSiteTrainable = make([]bool, VARRACKS_TYPE)
	var siteObj = make([]Site, VARRACKS_TYPE)
	for _, val := range idToSite {
		if val.owner == 0 && val.structureType == 2 {
			if isSiteTrainable[val.param2] {
				continue
			}
			isSiteTrainable[val.param2] = true
			siteObj[val.param2] = val
		}
	}

	for _, val := range trainingHistory {
		// val == ARCHER
		if val == 1 {
			// アーチャーの2回目はつくらない
			isSiteTrainable[1] = false
		}
	}

	// 末尾のほうが優先度が高い
	var potentialCandidate []int

	// 履歴があるものをpotentialClientに追加する
	for i := len(trainingHistory) - 1; i >= 0; i-- {
		varracksType := trainingHistory[i]
		if isSiteTrainable[varracksType] {
			potentialCandidate = append(potentialCandidate, varracksType)
			isSiteTrainable[varracksType] = false
		}
	}
	// 履歴がないものをpotentialClientに追加する
	for varracksType := 0; varracksType < VARRACKS_TYPE; varracksType++ {
		if isSiteTrainable[varracksType] {
			potentialCandidate = append(potentialCandidate, varracksType)
		}
	}

	TRAIN_COST := []int{80, 100, 140}
	var response []Site
	costSum := 0
	for i := len(potentialCandidate) - 1; i >= 0; i-- {
		val := potentialCandidate[i]
		costSum += TRAIN_COST[val]
		if costSum <= gold {
			trainingHistory = append(trainingHistory, val)
			response = append(response, siteObj[val])
		}
	}
	return response
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func log(data ...any) {
	fmt.Fprintf(os.Stderr, "%+v\n", data...)
}

func log2(data ...any) {
	fmt.Fprintln(os.Stderr, data...)
}

// KNIGHTの攻撃を受けている最中かどうか
func isUnderAttack(queen Unit, unitList []Unit) bool {
	for _, val := range unitList {
		if val.owner == 1 && val.unitType == 0 {
			dist := distUnit(queen, val)
			if dist < 30*30+10 {
				return true
			}
		}
	}
	return false
}

func decideBuildType(buildOrderList []BuildOrder, idToSite []Site, unitList []Unit) StructureType {
	return "TOWER"

	/*

		// 作れるサイトが余るときは末尾が採用される
		decideStructure := func(idx int) string {
			buildingPriority := []string{"BARRACKS-ARCHER", "MINE", "BARRACKS-KNIGHT", "TOWER", "MINE", "TOWER"}
			return buildingPriority[min(len(buildingPriority)-1, idx)]
		}


		// まだ建築されていないサイトが有るときは優先して作る
		var targetSite Site
		structureType := ""
		for idx, val := range idToSite {
			if idx >= len(idToSite)/3 {
				break
			}
			if val.owner == -1 {
				if decideStructure(idx) == "MINE" && (val.gold == 0 || isUnderAttack) {
					continue
				}
				targetSite = val
				structureType = decideStructure(idx)
				break
			}
		}
		if structureType != "" {
			return targetSite, structureType, nil
		}

		// 鉱山がレベルアップ可能なときは優先してレベルアップ
		for idx, val := range idToSite {
			// owner == 自分 && structureType == 鉱山 && 採掘可能goldが0でない
			if val.owner == 0 && (val.structureType == 0) {
				if val.param1 < val.maxMineSize {
					targetSite = val
					structureType = decideStructure(idx)
					break
				} else {
					continue
				}
			}
		}
		if structureType != "" {
			return targetSite, structureType, nil
		}

		// タワーがレベルアップ可能なときは優先してレベルアップ
		for idx, val := range idToSite {
			// owner == 自分 && structureType == タワー
			if val.owner == 0 && (val.structureType == 1) {
				targetSite = val
				structureType = decideStructure(idx)
				break
			}
		}
		if structureType != "" {
			return targetSite, structureType, nil
		}

		return Site{}, "", errors.New("可能な行動が存在しません")
	*/
}

func calcOptimalRoute(idToSite []Site, nearSiteList []*Site, queen Unit, buildOrderList []BuildOrder) ([]BuildOrder, []*Site) {
	// 初回
	if len(nearSiteList) == 0 {
		type Marker struct {
			idx  int
			dist int
		}
		distIdxMarker := make([]Marker, len(idToSite))

		for i := 0; i < len(idToSite); i++ {
			distIdxMarker[i].dist = dist(idToSite[i], queen)
			distIdxMarker[i].idx = i
		}
		sort.Slice(distIdxMarker, func(i, j int) bool {
			return distIdxMarker[i].dist < distIdxMarker[j].dist
		})

		nearSiteList = make([]*Site, len(idToSite)/3)
		for idx, v := range distIdxMarker {
			if idx >= len(idToSite)/3 {
				break
			}
			nearSiteList[idx] = &idToSite[v.idx]
		}

		optimalRouteIdx := travelingSalesman(nearSiteList)
		for _, r := range optimalRouteIdx {
			buildOrder := BuildOrder{}
			buildOrder.siteId = idToSite[r].siteId
			buildOrderList = append(buildOrderList, buildOrder)
		}
		return buildOrderList, nearSiteList
	}

	// siteIdがbuildOrderに含まれているかどうか判定する関数
	isScheduled := func(buildOrderList []BuildOrder, siteId SiteId) bool {
		for _, v := range buildOrderList {
			if v.siteId == siteId {
				return true
			}
		}
		return false
	}

	// 1/3のidToSiteで、自分の所有でない&buildOrderに含まれていないものを buildOrderの末尾に追加する
	// TODO 常に末尾は効率が悪い
	for _, v := range nearSiteList {
		if v.owner == -1 && !isScheduled(buildOrderList, v.siteId) {
			newBuildOrder := BuildOrder{}
			newBuildOrder.siteId = v.siteId
			buildOrderList = append(buildOrderList, newBuildOrder)
		}
	}

	// buildOrderListの中の更地ではないものを削除する(先頭のみ)
	removeLength := 0
	for _, v := range buildOrderList {
		if idToSite[v.siteId].owner == -1 {
			break
		} else {
			removeLength++
			log(v)
		}
	}
	buildOrderList = buildOrderList[removeLength:]

	return buildOrderList, nearSiteList
}

func main() {
	var numSites int
	fmt.Scan(&numSites)

	idToSite := make([]Site, numSites)
	var nearSiteList []*Site

	// サイトを回る順番
	buildOrderList := []BuildOrder{}

	for i := 0; i < numSites; i++ {
		var site Site
		fmt.Scan(&site.siteId, &site.x, &site.y, &site.radius)
		idToSite[site.siteId] = site
		// NOTE: バグ？-1になったり0になったりすることがある
		idToSite[site.siteId].gold = -1
	}

	for {
		// touchedSite: -1 if none
		var gold int
		var touchedSite SiteId
		fmt.Scan(&gold, &touchedSite)

		newSiteList := inputSiteList(numSites)

		var numUnits int
		fmt.Scan(&numUnits)

		unitList, queen := getUnitList(numUnits)

		idToSite = updateSiteList(numSites, idToSite, queen, newSiteList)

		buildOrderList, nearSiteList = calcOptimalRoute(idToSite, nearSiteList, queen, buildOrderList)

		if len(buildOrderList) == 0 {
			fmt.Println("MOVE", nearSiteList[0].x, nearSiteList[0].y)
		} else if touchedSite == buildOrderList[0].siteId {
			fmt.Println("BUILD", buildOrderList[0].siteId, decideBuildType(buildOrderList, idToSite, unitList))
			buildOrderList = buildOrderList[1:]
		} else {
			fmt.Println("BUILD", buildOrderList[0].siteId, "MINE")
		}

		trainingSite := calcTrainingSite(idToSite, gold)

		outputStr := "TRAIN"
		for _, t := range trainingSite {
			outputStr = fmt.Sprintf("%s %d", outputStr, t.siteId)
		}

		fmt.Println(outputStr)
	}
}
