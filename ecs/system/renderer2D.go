package system

import (
	"sort"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/graphics"
	cc "galaxyzeta.io/engine/infra/concurrency"
	"galaxyzeta.io/engine/infra/logger"
)

var NameRenderer2DSystem = "sys_Renderer2D"

type Renderer2DSystem struct {
	*base.SystemBase
	renderers       []graphics.IRenderable       // dynamically re-arranged according to elements' Z coordinate.
	staticRenderers []graphics.IRenderable       // will not be sorted, has a static Z coordinate. Register to this slice to optimize your game performace.
	indexer         map[graphics.IRenderable]int // indexer is meant to be updated after every iteration of rendering. It is useful when we try to delete an element
	logger          *logger.Logger
}

func NewRenderer2DSystem(priority int) *Renderer2DSystem {
	return &Renderer2DSystem{
		SystemBase:      base.NewSystemBase(priority),
		indexer:         map[graphics.IRenderable]int{},
		renderers:       []graphics.IRenderable{},
		staticRenderers: []graphics.IRenderable{},
		logger:          logger.New(NameRenderer2DSystem),
	}
}

func (ren *Renderer2DSystem) execute(_ *cc.Executor) {
	cam := graphics.GetCurrentCamera()
	// sort spriteRenderers first
	sort.SliceStable(ren.renderers, func(i, j int) bool {
		return ren.renderers[i].Z() < ren.renderers[j].Z()
	})
	ptr1, ptr2 := 0, 0
	idx := 0
	for ptr1 < len(ren.renderers) && ptr2 < len(ren.staticRenderers) {
		if ren.renderers[ptr1].Z() >= ren.staticRenderers[ptr2].Z() {
			ren.doRenderExecute(&ptr1, &idx, cam, ren.renderers)
		} else {
			ren.doRenderExecute(&ptr2, &idx, cam, ren.staticRenderers)
		}
	}
	for ptr1 < len(ren.renderers) {
		ren.doRenderExecute(&ptr1, &idx, cam, ren.renderers)
	}
	for ptr2 < len(ren.staticRenderers) {
		ren.doRenderExecute(&ptr2, &idx, cam, ren.staticRenderers)
	}
}

func (ren *Renderer2DSystem) doRenderExecute(ptr *int, idx *int, cam *graphics.Camera, targetSlice []graphics.IRenderable) {
	sr := ren.renderers[*ptr]
	sr.Render(cam)
	ren.indexer[sr] = *idx
	sr.PostRender()
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
	comps := iobj.Obj().GetAllComponents()
	for _, comp := range comps {
		ren, ok := comp.(graphics.IRenderable)
		if !ok {
			continue
		}
		if ren.IsStatic() {
			pos := s.binarySearchStatic(ren.Z())
			copy(s.staticRenderers[pos:], s.staticRenderers[pos+1:])
			s.staticRenderers[pos] = ren
			s.indexer[ren] = pos
		} else {
			s.renderers = append(s.renderers, ren)
			s.indexer[ren] = len(s.renderers) - 1
		}
	}

}

// TODO need test
func (s *Renderer2DSystem) binarySearchStatic(z int64) int {
	left, right := 0, len(s.staticRenderers)-1
	for left <= right {
		mid := left + (right-left)<<1
		if s.staticRenderers[mid].Z() < z {
			left = mid + 1
		} else if s.staticRenderers[mid].Z() == z {
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
	rens := getIRenderers(iobj)
	for _, ren := range rens {
		index, ok := s.indexer[ren]
		if !ok {
			panic("should not happen")
		}
		if ren.IsStatic() {
			// order is important, need to rearrange
			for i, j := index+1, index; i < len(s.staticRenderers); i, j = i+1, j+1 {
				s.staticRenderers[j] = s.staticRenderers[i]
				s.indexer[s.staticRenderers[j]] = j
			}
			s.staticRenderers = s.staticRenderers[:len(s.staticRenderers)-1]
		} else {
			// order here is not important here
			// because all elements will be sorted again before next rendering process.
			s.renderers[index] = s.renderers[len(s.renderers)-1]
			s.indexer[s.renderers[index]] = index
			s.renderers[len(s.renderers)-1] = nil
			s.renderers = s.renderers[:len(s.renderers)-1]
			delete(s.indexer, ren)
		}
	}

}

func (s *Renderer2DSystem) Activate(iobj base.IGameObject2D) {
	s.Register(iobj)
}

func (s *Renderer2DSystem) Deactivate(iobj base.IGameObject2D) {
	s.Unregister(iobj)
}

func getIRenderers(iobj base.IGameObject2D) (rens []graphics.IRenderable) {
	rens = make([]graphics.IRenderable, 0)
	comps := iobj.Obj().GetAllComponents()
	for _, comp := range comps {
		ren, ok := comp.(graphics.IRenderable)
		if !ok {
			continue
		}
		rens = append(rens, ren)
	}
	return
}
