package graphics

import "sync"

type Animator struct {
	mu           sync.RWMutex
	state2clip   map[string]*SpriteInstance // maps from state name to an animation clip.
	currentState string
}

type StateClipPair struct {
	State string
	Clip  *SpriteInstance
}

func NewAnimator(cfgs ...StateClipPair) (anmt *Animator) {
	anmt = &Animator{
		mu:           sync.RWMutex{},
		state2clip:   map[string]*SpriteInstance{},
		currentState: cfgs[0].State,
	}
	for _, state := range cfgs {
		anmt.state2clip[state.State] = state.Clip
	}
	return anmt
}

func (a *Animator) Spr() (ret *SpriteInstance) {
	a.mu.RLock()
	ret = a.state2clip[a.currentState]
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
	a.state2clip[s] = spr
	a.mu.Unlock()
}

func (a *Animator) RegisterStates(sps ...StateClipPair) {
	a.mu.Lock()
	for _, pair := range sps {
		a.state2clip[pair.State] = pair.Clip
	}
	a.mu.Unlock()
}
