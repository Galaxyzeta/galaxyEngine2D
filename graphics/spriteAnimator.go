package graphics

import "sync"

type Animator struct {
	mu           sync.RWMutex
	stateMap     map[string]*SpriteInstance
	currentState string
}

type StateSpritePair struct {
	State string
	Spr   *SpriteInstance
}

func NewAnimator(cfgs ...StateSpritePair) (anmt *Animator) {
	anmt = &Animator{
		mu:           sync.RWMutex{},
		stateMap:     map[string]*SpriteInstance{},
		currentState: cfgs[0].State,
	}
	for _, state := range cfgs {
		anmt.stateMap[state.State] = state.Spr
	}
	return anmt
}

func (a *Animator) Spr() (ret *SpriteInstance) {
	a.mu.RLock()
	ret = a.stateMap[a.currentState]
	a.mu.RUnlock()
	return
}

func (a *Animator) AlterState(toState string) {
	a.mu.Lock()
	a.currentState = toState
	a.mu.Unlock()
}

func (a *Animator) RegisterState(spr *SpriteInstance, s string) {
	a.mu.Lock()
	a.stateMap[s] = spr
	a.mu.Unlock()
}

func (a *Animator) RegisterStates(sps ...StateSpritePair) {
	a.mu.Lock()
	for _, pair := range sps {
		a.stateMap[pair.State] = pair.Spr
	}
	a.mu.Unlock()
}
