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
	siteId, x, y, radius, ignore1, ignore2, structureType, owner, param1, param2, dist int
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
				siteList[si].ignore2 = tmpList[ti].ignore2
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
		// ignore1: used in future leagues
		// ignore2: used in future leagues
		// structureType: -1 = No structure, 2 = Barracks
		// owner: -1 = No structure, 0 = Friendly, 1 = Enemy
		var siteId int
		fmt.Scan(&siteId)
		site := &tmpList[siteId]
		site.siteId = siteId
		fmt.Scan(&site.ignore1, &site.ignore2, &site.structureType, &site.owner, &site.param1, &site.param2)
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
	if isTrainable[waitingTrain] == false {
		// 訓練予定のサイトが存在しないとき
		setNextWaiting()
		return []Site{}
	} else if TRAIN_COST[waitingTrain] > gold {
		// お金が足りずに待機するとき
		return []Site{}
	} else {
		// 予定通り訓練できるとき
		setNextWaiting()
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

func main() {
	var numSites int
	fmt.Scan(&numSites)

	siteList := make([]Site, numSites)

	for i := 0; i < numSites; i++ {
		var site Site
		fmt.Scan(&site.siteId, &site.x, &site.y, &site.radius)
		siteList[site.siteId] = site
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

		var targetSite Site
		targetIdx := -1
		for idx, val := range siteList {
			if idx >= len(siteList)/3 {
				break
			}
			if val.owner != 0 {
				targetSite = val
				targetIdx = idx
				break
			}
		}

		if targetIdx == -1 {
			for idx, val := range siteList {
				if val.owner == 0 && (val.structureType == 0) {
					targetSite = val
					targetIdx = idx
					break
				}
			}
		}

		buildings := []string{"BARRACKS-KNIGHT", "MINE", "TOWER"}
		fmt.Fprintln(os.Stderr, targetIdx)
		if targetIdx != -1 {
			fmt.Println("BUILD", targetSite.siteId, buildings[min(2, targetIdx)])
		} else {
			fmt.Println("MOVE", siteList[0].x, siteList[0].y)
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		trainingSite := calcTrainingSite(siteList, gold)

		outputStr := "TRAIN"
		for _, t := range trainingSite {
			outputStr = fmt.Sprintf("%s %d", outputStr, t.siteId)
		}

		fmt.Println(outputStr)
	}
}
