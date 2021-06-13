package core

type IComponent interface {
	GetName() string
}

type ISystem interface {
	Execute()
	GetPriority() int
}
