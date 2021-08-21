package core

import (
	"sort"

	"galaxyzeta.io/engine/base"
)

// RegisterSystems to the game. Cannot disamount systems currently.
func RegisterSystem(sys ...base.ISystem) {
	systemPriorityList = append(systemPriorityList, sys...)
	SystemSort()
	// re-assign pos
	for i, s := range systemPriorityList {
		system2Priority[s] = i
		name2System[s.GetName()] = s
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
	pos := system2Priority[sys]
	delete(system2Priority, sys)
	delete(name2System, sys.GetName())
	slen := len(systemPriorityList)
	// removal of element
	for i := pos + 1; i < slen; i++ {
		systemPriorityList[i-1] = systemPriorityList[i]
	}
	systemPriorityList = systemPriorityList[:slen-1]
	SystemSort()
}

// SubscribeSystem registers an object into given system.
// Will panic if the system was not found.
func SubscribeSystem(iobj base.IGameObject2D, sysname string) {
	sys := name2System[sysname]
	sys.Register(iobj)
	iobj.Obj().AppendSubscribedSystem(sys)
}

// UnsubscribeSystem unregisters an object from given system.
// Will panic if the system was not found.
func UnsubscribeSystem(iobj base.IGameObject2D, sysname string) {
	sys := name2System[sysname]
	sys.Unregister(iobj)
	iobj.Obj().RemoveSubscribedSystem(sys)

}

func GetSystem(sysname string) base.ISystem {
	mu := GetRWMutex(Mutex_System)
	mu.RLock()
	defer mu.RUnlock()
	return name2System[sysname]
}
