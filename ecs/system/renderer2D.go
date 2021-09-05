package system

import (
	"sort"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/graphics"
	cc "galaxyzeta.io/engine/infra/concurrency"
	"galaxyzeta.io/engine/infra/logger"
)

var NameRenderer2DSystem = "sys_Renderer2D"

type Renderer2DSystem struct {
	*base.SystemBase
	spriteRenderers       []*component.SpriteRenderer       // dynamically re-arranged according to elements' Z coordinate.
	staticSpriteRenderers []*component.SpriteRenderer       // will not be sorted, has a static Z coordinate. Register to this slice to optimize your game performace.
	indexer               map[*component.SpriteRenderer]int // indexer is meant to be updated after every iteration of rendering. It is useful when we try to delete an element
	logger                *logger.Logger
}

func NewRenderer2DSystem(priority int) *Renderer2DSystem {
	return &Renderer2DSystem{
		SystemBase:      base.NewSystemBase(priority),
		indexer:         map[*component.SpriteRenderer]int{},
		spriteRenderers: []*component.SpriteRenderer{},
		logger:          logger.New(NameRenderer2DSystem),
	}
}

func (ren *Renderer2DSystem) execute(_ *cc.Executor) {
	cam := graphics.GetCurrentCamera()
	// sort spriteRenderers first
	sort.SliceStable(ren.spriteRenderers, func(i, j int) bool {
		return ren.spriteRenderers[i].Z < ren.spriteRenderers[j].Z
	})
	ptr1, ptr2 := 0, 0
	idx := 0
	for ptr1 < len(ren.spriteRenderers) && ptr2 < len(ren.staticSpriteRenderers) {
		if ren.spriteRenderers[ptr1].Z >= ren.staticSpriteRenderers[ptr2].Z {
			ren.doRenderExecute(&ptr1, &idx, cam, ren.spriteRenderers)
		} else {
			ren.doRenderExecute(&ptr2, &idx, cam, ren.staticSpriteRenderers)
		}
	}
	for ptr1 < len(ren.spriteRenderers) {
		ren.doRenderExecute(&ptr1, &idx, cam, ren.spriteRenderers)
	}
	for ptr2 < len(ren.staticSpriteRenderers) {
		ren.doRenderExecute(&ptr2, &idx, cam, ren.staticSpriteRenderers)
	}
}

func (ren *Renderer2DSystem) doRenderExecute(ptr *int, idx *int, cam *graphics.Camera, targetSlice []*component.SpriteRenderer) {
	sr := ren.spriteRenderers[*ptr]
	sr.Render(cam)
	ren.indexer[sr] = *idx
	sr.Spr().DoFrameStep()
	*ptr++
	*idx++
}

// ===== IMPLEMENTATION =====

func (s *Renderer2DSystem) Execute(executor *cc.Executor) {
	s.execute(executor)
}

func (s *Renderer2DSystem) GetSystemBase() *base.SystemBase {
	return s.SystemBase
}

func (s *Renderer2DSystem) GetName() string {
	return NameRenderer2DSystem
}

func (s *Renderer2DSystem) Register(iobj base.IGameObject2D) {
	sr := getSpriteRenderer(iobj)
	if sr.IsStatic() {
		pos := s.binarySearchStatic(sr.Z)
		copy(s.staticSpriteRenderers[pos:], s.staticSpriteRenderers[pos+1:])
		s.staticSpriteRenderers[pos] = sr
		s.indexer[sr] = pos
	} else {
		s.spriteRenderers = append(s.spriteRenderers, sr)
		s.indexer[sr] = len(s.spriteRenderers) - 1
	}
}

// TODO need test
func (s *Renderer2DSystem) binarySearchStatic(z int64) int {
	left, right := 0, len(s.staticSpriteRenderers)-1
	for left <= right {
		mid := left + (right-left)<<1
		if s.staticSpriteRenderers[mid].Z < z {
			left = mid + 1
		} else if s.staticSpriteRenderers[mid].Z == z {
			return mid
		} else {
			right = mid - 1
		}
	}
	if left < 0 {
		return 0
	}
	return left
}

func (s *Renderer2DSystem) Unregister(iobj base.IGameObject2D) {
	sr := getSpriteRenderer(iobj)
	index, ok := s.indexer[sr]
	if !ok {
		panic("should not happen")
	}
	if sr.IsStatic() {
		// order is important, need to rearrange
		for i, j := index+1, index; i < len(s.staticSpriteRenderers); i, j = i+1, j+1 {
			s.staticSpriteRenderers[j] = s.staticSpriteRenderers[i]
			s.indexer[s.staticSpriteRenderers[j]] = j
		}
		s.staticSpriteRenderers = s.staticSpriteRenderers[:len(s.staticSpriteRenderers)-1]
	} else {
		// order here is not important here
		// because all elements will be sorted again before next rendering process.
		s.spriteRenderers[index] = s.spriteRenderers[len(s.spriteRenderers)-1]
		s.indexer[s.spriteRenderers[index]] = index
		s.spriteRenderers[len(s.spriteRenderers)-1] = nil
		s.spriteRenderers = s.spriteRenderers[:len(s.spriteRenderers)-1]
		delete(s.indexer, sr)
	}
}

func getSpriteRenderer(iobj base.IGameObject2D) *component.SpriteRenderer {
	return iobj.Obj().GetComponent(component.NameSpriteRenderer).(*component.SpriteRenderer)
}
