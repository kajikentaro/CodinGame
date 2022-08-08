package main

import (
	"container/heap"
	"errors"
	"fmt"
	"os"
	"sort"
)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/
type Site struct {
	siteId, x, y, radius, gold, maxMineSize, structureType, owner, param1, param2, dist int
}

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

type HeapItem struct {
	idx, distSum int
}

type PairHeap []*HeapItem

func (h PairHeap) Len() int           { return len(h) }
func (h PairHeap) Less(i, j int) bool { return h[i].distSum < h[j].distSum }
func (h PairHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *PairHeap) Push(x any) {
	*h = append(*h, x.(*HeapItem))
}
func (h *PairHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func dike(siteList []Site) {
	n := len(siteList)
	path := make([][]int, n)
	for i := 0; i < n; i++ {
		path[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			path[i][j] = distSite(siteList[i], siteList[j])
		}
	}

	pq := make(PairHeap, n)
	res := make([]int, n)
	for i := range res {
		res[i] = -1
	}
	heap.Init(&pq)
	pq.Push(HeapItem{0, 0})

	for {
		if pq.Len() == 0 {
			break
		}
		now := (pq.Pop()).(*HeapItem)
		if res[now.idx] != -1 {
			continue
		}
		res[now.idx] = 
		for i, p := range path[now.idx] {
			if res[i] != -1 {
				continue
			}
			pq.Push(HeapItem{i, now.distSum + p})
		}
	}


}

func updateSiteList(numSites int, siteList []Site, queen Unit, tmpList []Site) {
	for si := 0; si < numSites; si++ {
		for ti := 0; ti < numSites; ti++ {
			if tmpList[ti].siteId == siteList[si].siteId {
				// NOTE: バグ？goldが0になったり-1になったりする
				if siteList[si].gold != 0 {
					siteList[si].gold = tmpList[ti].gold
				}
				siteList[si].maxMineSize = tmpList[ti].maxMineSize
				siteList[si].structureType = tmpList[ti].structureType
				siteList[si].owner = tmpList[ti].owner
				siteList[si].param1 = tmpList[ti].param1
				siteList[si].param2 = tmpList[ti].param2

			}
		}
	}

	// 初回のみ
	if siteList[0].dist == 0 && siteList[1].dist == 0 {
		for i := 0; i < len(siteList); i++ {
			siteList[i].dist = dist(siteList[i], queen)
		}
		sort.Slice(siteList, func(i, j int) bool {
			return siteList[i].dist < siteList[j].dist
		})
	}
}

func inputSiteList(numSites int) []Site {
	tmpList := make([]Site, numSites)
	for i := 0; i < numSites; i++ {
		// gold: 発掘できる残りの金
		// maxMineSize: used in future leagues
		// structureType: -1 = No structure, 2 = Barracks
		// owner: -1 = No structure, 0 = Friendly, 1 = Enemy
		var siteId int
		fmt.Scan(&siteId)
		site := &tmpList[siteId]
		site.siteId = siteId
		fmt.Scan(&site.gold, &site.maxMineSize, &site.structureType, &site.owner, &site.param1, &site.param2)
	}
	return tmpList
}

func inputUnitList(numUnits int) ([]Unit, Unit) {
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

// 末尾が新しい
var trainingHistory []int

func calcTrainingSite(siteList []Site, gold int) []Site {
	VARRACKS_TYPE := 3
	var isSiteTrainable = make([]bool, VARRACKS_TYPE)
	var siteObj = make([]Site, VARRACKS_TYPE)
	for _, val := range siteList {
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
	fmt.Fprintf(os.Stderr, "%+v\n", data)
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

func calcBuildingSite(siteList []Site, isUnderAttack bool) (Site, string, error) {
	// 作れるサイトが余るときは末尾が採用される
	decideStructure := func(idx int) string {
		buildingPriority := []string{"BARRACKS-ARCHER", "MINE", "BARRACKS-KNIGHT", "TOWER", "MINE", "TOWER"}
		return buildingPriority[min(len(buildingPriority)-1, idx)]
	}

	// まだ建築されていないサイトが有るときは優先して作る
	var targetSite Site
	structureType := ""
	for idx, val := range siteList {
		if idx >= len(siteList)/3 {
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
	for idx, val := range siteList {
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
	for idx, val := range siteList {
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
}

// func calcBuildingStrategy(siteList){

func main() {
	var numSites int
	fmt.Scan(&numSites)

	siteList := make([]Site, numSites)

	for i := 0; i < numSites; i++ {
		var site Site
		fmt.Scan(&site.siteId, &site.x, &site.y, &site.radius)
		siteList[site.siteId] = site
		// NOTE: バグ？-1になったり0になったりすることがある
		siteList[site.siteId].gold = -1
	}

	for {
		// touchedSite: -1 if none
		var gold, touchedSite int
		fmt.Scan(&gold, &touchedSite)

		newSiteList := inputSiteList(numSites)

		var numUnits int
		fmt.Scan(&numUnits)

		unitList, queen := inputUnitList(numUnits)

		updateSiteList(numSites, siteList, queen, newSiteList)

		targetSite, targetType, err := calcBuildingSite(siteList, isUnderAttack(queen, unitList))

		if err == nil {
			fmt.Println("BUILD", targetSite.siteId, targetType)
		} else {
			fmt.Println("MOVE", siteList[0].x, siteList[0].y)
		}

		trainingSite := calcTrainingSite(siteList, gold)
		log(trainingSite)

		outputStr := "TRAIN"
		for _, t := range trainingSite {
			outputStr = fmt.Sprintf("%s %d", outputStr, t.siteId)
		}

		fmt.Println(outputStr)
	}
}
