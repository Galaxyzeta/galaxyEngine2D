package core

import (
	"sort"

	"galaxyzeta.io/engine/base"
)

// RegisterSystems to the game. Not thread safe.
func RegisterSystem(sys ...base.ISystem) {
	doRegisterSystem(&systemPriorityList, sys...)
}

// RegisterGfxSystem to the game. Will be called every render FPS in render loop. Not thread safe.
func RegisterGfxSystem(sys ...base.ISystem) {
	doRegisterGfxSystem(&gfxSystemPriorityList, sys...)
}

func doRegisterSystem(systemList *[]base.ISystem, sys ...base.ISystem) {
	*systemList = append(*systemList, sys...)
	doSystemSort(*systemList)
	// re-assign pos
	for i, s := range *systemList {
		system2Priority[s] = i
		name2System[s.GetName()] = s
	}
}

func doRegisterGfxSystem(gfxSystemList *[]base.ISystem, sys ...base.ISystem) {
	*gfxSystemList = append(*gfxSystemList, sys...)
	doSystemSort(gfxSystemPriorityList)
	// re-assign pos
	for i, s := range *gfxSystemList {
		system2Priority[s] = i
		name2System[s.GetName()] = s
	}
}

func doSystemSort(systemList []base.ISystem) {
	sort.Slice(systemList, func(i, j int) bool {
		return systemList[i].GetSystemBase().GetPriority() > systemList[j].GetSystemBase().GetPriority()
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

	// TODO need improvement
	doSystemSort(systemPriorityList)
	doSystemSort(gfxSystemPriorityList)
}

// SubscribeSystem registers an object into given system.
// Will panic if the system was not found.
func SubscribeSystem(iobj base.IGameObject2D, sysname string) {
	sys := name2System[sysname]
	// sys.Register(iobj)	delayed execution
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
