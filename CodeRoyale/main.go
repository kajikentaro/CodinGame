package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

// 23.11
// 23.40

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/
type PointF struct {
	x, y float64
}

func (p *PointF) equal(q PointF) bool {
	if math.Abs(p.x-q.x) < 1e-5 && math.Abs(p.y-q.y) < 1e-5 {
		return true
	}
	return false
}

type Point struct {
	x, y int
}

func (p *Point) float() PointF {
	return PointF{float64(p.x), float64(p.y)}
}

type Site struct {
	siteId                                                                SiteId
	p                                                                     Point
	radius, gold, maxMineSize, structureType, owner, param1, param2, dist int
}

type SiteId int

type BuildOrder struct {
	siteId SiteId
}

type StructureType string

// すべての敷地のうち何割を構築目標とするか
func BUILD_RATIO(allSiteSize int) int {
	return allSiteSize / 3
}

const (
	BARRACKS_KNIGHT = StructureType("BARRACKS-KNIGHT")
	BARRACKS_ARCHER = StructureType("BARRACKS-ARCHER")
	BARRACKS_GIANT  = StructureType("BARRACKS-GIANT")
	TOWER           = StructureType("TOWER")
	MINE            = StructureType("MINE")
)

type Unit struct {
	p                       Point
	owner, unitType, health int
}

func dist(a, b Point) int {
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
		fmt.Scan(&unit.p.x, &unit.p.y, &unit.owner, &unit.unitType, &unit.health)
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
			p := P{dist(nearSiteList[i].p, nearSiteList[j].p), j}
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
	for idx, val := range dp_log[pow(2, int(n))-1][s] {
		if idx == n-1 {
			break
		}
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
			dist := dist(queen.p, val.p)
			if dist < 30*30+10 {
				return true
			}
		}
	}
	return false
}

// ベクトルoa, obの外積を計算する
func calcCrossProduct(a Point, b PointF, o Point) float64 {
	af := a.float()
	of := o.float()
	return ((af.x-of.x)*(b.y-of.y) - (af.y-of.y)*(b.x-of.x))
}

func calcContact(start, circle PointF, r float64) (PointF, PointF) {
	start.x -= circle.x
	start.y -= circle.y

	r2 := math.Pow(r, 2)
	r4 := math.Pow(r, 4)
	a := start.x
	b := start.y
	a2 := math.Pow(a, 2)
	b2 := math.Pow(b, 2)
	b4 := math.Pow(b, 4)

	root := math.Sqrt(-b2*r4 + a2*b2*r2 + b4*r2)
	x1 := (a*r2 + root) / (a2 + b2)
	y1 := (r2 - x1*a) / b
	x2 := (a*r2 - root) / (a2 + b2)
	y2 := (r2 - x2*a) / b
	return PointF{x1 + circle.x, y1 + circle.y}, PointF{x2 + circle.x, y2 + circle.y}
}

func isWin(unitList []Unit) bool {
	var enemyQueen Unit
	var friendlyQueen Unit
	for _, val := range unitList {
		if val.owner == 1 && val.unitType == -1 {
			enemyQueen = val
		}
		if val.owner == 0 && val.unitType == -1 {
			friendlyQueen = val
		}
	}
	return enemyQueen.health < friendlyQueen.health
}

func isNearEnemy(site Site, unitList []Unit) bool {
	for _, val := range unitList {
		if val.owner == 1 && val.unitType == 0 && dist(site.p, val.p) < 400000 {
			return true
		}
	}
	return false
}

var continousCnt int
var preTouchedId SiteId = -1

func decideBuildType(buildOrderList []BuildOrder, idToSite []Site, unitList []Unit, touchedSite SiteId) (StructureType, bool) {
	var nowBuildingType StructureType

	friendly := map[StructureType]int{}
	enemy := map[StructureType]int{}
	if preTouchedId == touchedSite {
		continousCnt++
	} else {
		continousCnt = 1
	}
	defer func() {
		preTouchedId = touchedSite
	}()

	for _, site := range idToSite {
		var structureType StructureType

		if site.structureType == 0 {
			structureType = "MINE"
		}
		if site.structureType == 1 {
			structureType = "TOWER"
		}
		if site.structureType == 2 {
			if site.param2 == 0 {
				structureType = "BARRACKS-KNIGHT"
			}
			if site.param2 == 1 {
				structureType = "BARRACKS-ARCHER"
			}
			if site.param2 == 2 {
				structureType = "BARRACKS-GIANT"
			}
		}
		// 今建設中(レベルアップ中)のものは数に入れない
		if site.siteId == buildOrderList[0].siteId {
			nowBuildingType = structureType
			continue
		}
		if site.owner == 0 {
			friendly[structureType]++
		}
		if site.owner == 1 {
			enemy[structureType]++
		}
	}

	targetSite := idToSite[buildOrderList[0].siteId]

	/** 前回作ったものを強化するここから **/
	if nowBuildingType == "MINE" {
		if targetSite.param1 >= targetSite.maxMineSize-1 {
			// 最大まで強化したとき
			return "MINE", true
		}
		return "MINE", false
	}

	if nowBuildingType == "TOWER" {
		if continousCnt == 4 {
			return "TOWER", true
		}
		return "TOWER", false
	}
	/** 強化ここまで **/

	if friendly["TOWER"] == 0 && friendly["BARRACKS-KNIGHT"] == 0 && enemy["BARRACKS-KNIGHT"] > 0 && isNearEnemy(targetSite, unitList) {
		return "TOWER", false
	}

	if friendly["MINE"] <= 2 && targetSite.gold >= 50 && !isWin(unitList) && !isNearEnemy(targetSite, unitList) {
		// 作ったのにすぐ壊されたとき
		if continousCnt > 1 && targetSite.owner != 0 {
			// なにもしない
		} else if targetSite.maxMineSize == 1 {
			return "MINE", true
		} else {
			return "MINE", false
		}
	}

	if friendly["BARRACKS-KNIGHT"] == 0 {
		return "BARRACKS-KNIGHT", true
	}

	if enemy["TOWER"] >= 3 && friendly["BARRACKS-GIANT"] == 0 && continousCnt == 0 {
		return "BARRACKS-GIANT", true
	}

	if friendly["TOWER"] <= 5 {
		return "TOWER", false
	}
	return "TOWER", false
}

func calcOptimalCoordinate(buildOrderList []BuildOrder, idToSite []Site, queen Unit) PointF {
	if len(buildOrderList) <= 1 {
		target := idToSite[buildOrderList[0].siteId]
		return target.p.float()
	}
	nextCircle := idToSite[buildOrderList[0].siteId]
	ans1, ans2 := calcContact(queen.p.float(), nextCircle.p.float(), float64(nextCircle.radius))
	cc1 := calcCrossProduct(queen.p, ans1, nextCircle.p)
	cc2 := calcCrossProduct(queen.p, ans2, nextCircle.p)
	cc0 := calcCrossProduct(queen.p, idToSite[buildOrderList[1].siteId].p.float(), nextCircle.p)
	if cc1*cc0 > 0 {
		return ans1
	}
	if cc2*cc0 > 0 {
		return ans2
	}
	if cc0 == 0 {
		return ans1
	}
	log("an error occered in calcOptimalCoordinate", cc0, cc1, cc2)
	os.Exit(1)
	return PointF{}
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
			distIdxMarker[i].dist = dist(idToSite[i].p, queen.p)
			idToSite[i].dist = distIdxMarker[i].dist
			distIdxMarker[i].idx = i
		}
		sort.Slice(distIdxMarker, func(i, j int) bool {
			return distIdxMarker[i].dist < distIdxMarker[j].dist
		})

		nearSiteList = make([]*Site, BUILD_RATIO(len(idToSite)))
		for idx, v := range distIdxMarker {
			if idx >= BUILD_RATIO(len(idToSite)) {
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

	// 更地&buildOrderに含まれていないものを buildOrderの末尾に追加する
	// TODO 常に末尾は効率が悪い
	for _, v := range nearSiteList {
		if v.dist >= 150000 {
			continue
		}
		if v.owner == -1 && !isScheduled(buildOrderList, v.siteId) {
			newBuildOrder := BuildOrder{}
			newBuildOrder.siteId = v.siteId
			buildOrderList = append(buildOrderList, newBuildOrder)
		}
	}

	// buildOrderListの中の敵のサイトを削除する(先頭のみ)
	removeLength := 0
	for _, v := range buildOrderList {
		if idToSite[v.siteId].owner == 1 {
			removeLength++
		} else {
			break
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
		fmt.Scan(&site.siteId, &site.p.x, &site.p.y, &site.radius)
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
			// サイトがすべて構築済みの場合
			fmt.Println("BUILD", nearSiteList[0].siteId, "TOWER")
		} else if touchedSite == buildOrderList[0].siteId {
			// サイトを建設する場合
			buildType, isPopBuildOrder := decideBuildType(buildOrderList, idToSite, unitList, touchedSite)
			fmt.Println("BUILD", buildOrderList[0].siteId, buildType)
			if isPopBuildOrder {
				buildOrderList = buildOrderList[1:]
			}
		} else {
			// サイトに向かって移動する場合
			p := calcOptimalCoordinate(buildOrderList, idToSite, queen)
			// サイトに向かって移動する場合
			fmt.Println("MOVE", math.Round(p.x), math.Round(p.y))
		}

		trainingSite := calcTrainingSite(idToSite, gold)

		outputStr := "TRAIN"
		for _, t := range trainingSite {
			outputStr = fmt.Sprintf("%s %d", outputStr, t.siteId)
		}

		fmt.Println(outputStr)
	}
}
