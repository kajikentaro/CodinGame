package main

import (
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

func dist(a Site, b Unit) int {
	return (a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y)
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

var waitingTrain int = 0

func calcTrainingSite(siteList []Site, gold int) []Site {
	VARRACKS_TYPE := 3
	var isTrainable = make([]bool, VARRACKS_TYPE)
	var trainableObj = make([]Site, VARRACKS_TYPE)
	for _, val := range siteList {
		if val.owner == 0 && val.structureType == 2 {
			if isTrainable[val.param2] {
				continue
			}
			isTrainable[val.param2] = true
			trainableObj[val.param2] = val
		}
	}

	setNextWaiting := func() {
		for i := 0; i < VARRACKS_TYPE; i++ {
			waitingTrain = (waitingTrain + 1) % 3
			if isTrainable[waitingTrain] {
				break
			}
		}
	}

	TRAIN_COST := []int{80, 100, 140}
	log(waitingTrain)
	log(TRAIN_COST[waitingTrain], gold, TRAIN_COST[waitingTrain] > gold)
	log([]Site{trainableObj[waitingTrain]})
	if isTrainable[waitingTrain] == false {
		// 訓練予定のサイトが存在しないとき
		setNextWaiting()
		return []Site{}
	} else if TRAIN_COST[waitingTrain] > gold {
		// お金が足りずに待機するとき
		return []Site{}
	} else {
		// 予定通り訓練できるとき
		defer setNextWaiting()
		return []Site{trainableObj[waitingTrain]}
	}
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

func calcBuildingSite(siteList []Site) (Site, string, error) {
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
			if decideStructure(idx) == "MINE" && val.gold == 0 {
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

		//unitList, queen := getUnitList(numUnits)
		_, queen := getUnitList(numUnits)

		updateSiteList(numSites, siteList, queen, newSiteList)

		targetSite, targetType, err := calcBuildingSite(siteList)

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
