package core

import "sort"

type IComponent interface {
	GetName() string
}

type ISystem interface {
	Execute()
	GetPriority() int
	Disable()
	Enable()
	IsEnabled() bool
}

// RegisterSystems to the game. Cannot disamount systems currently.
func RegisterSystem(sys ...ISystem) {
	systemPriorityList = append(systemPriorityList, sys...)
	SystemSort()
	// re-assign pos
	for i, s := range systemPriorityList {
		systemMap[s] = i
	}
}

// SystemSort re-sort all registered systems' priorities from low to hign.
func SystemSort() {
	sort.Slice(systemPriorityList, func(i, j int) bool {
		return systemPriorityList[i].GetPriority() > systemPriorityList[j].GetPriority()
	})
}

// UnregisterSystem delete an system
func UnregisterSystem(sys ISystem) {
	pos := systemMap[sys]
	delete(systemMap, sys)
	slen := len(systemPriorityList)
	// removal of element
	for i := pos + 1; i < slen; i++ {
		systemPriorityList[i-1] = systemPriorityList[i]
	}
	systemPriorityList = systemPriorityList[:slen-1]
	SystemSort()
}
