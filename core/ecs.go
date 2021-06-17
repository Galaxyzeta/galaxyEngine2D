package core

import (
	"sort"

	"galaxyzeta.io/engine/base"
)

// JoinSystem
func JoinSystem(iobj2d, IGameObject2D, sys base.ISystem) {

}

// RegisterSystems to the game. Cannot disamount systems currently.
func RegisterSystem(sys ...base.ISystem) {
	systemPriorityList = append(systemPriorityList, sys...)
	SystemSort()
	// re-assign pos
	for i, s := range systemPriorityList {
		systemPriorityMap[s] = i
	}
}

// SystemSort re-sort all registered systems' priorities from low to hign.
func SystemSort() {
	sort.Slice(systemPriorityList, func(i, j int) bool {
		return systemPriorityList[i].GetSystemBase().GetPriority() > systemPriorityList[j].GetSystemBase().GetPriority()
	})
}

// UnregisterSystem delete an system
func UnregisterSystem(sys base.ISystem) {
	pos := systemPriorityMap[sys]
	delete(systemPriorityMap, sys)
	slen := len(systemPriorityList)
	// removal of element
	for i := pos + 1; i < slen; i++ {
		systemPriorityList[i-1] = systemPriorityList[i]
	}
	systemPriorityList = systemPriorityList[:slen-1]
	SystemSort()
}

func GetSystem() {

}
