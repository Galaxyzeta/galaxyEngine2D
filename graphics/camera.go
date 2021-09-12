package graphics

import (
	"sync"

	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
)

type Camera struct {
	pos        linalg.Vector2f64
	resolution linalg.Vector2f64
	rwmutex    *sync.RWMutex
}

func NewCamera(pos linalg.Vector2f64, resolution linalg.Vector2f64) *Camera {
	return &Camera{
		pos:        pos,
		resolution: resolution,
		rwmutex:    &sync.RWMutex{},
	}
}

func (c *Camera) GetPos() (pos linalg.Vector2f64) {
	c.rwmutex.RLock()
	pos = c.pos
	c.rwmutex.RUnlock()
	return
}

func (c *Camera) GetResolution() (res linalg.Vector2f64) {
	c.rwmutex.RLock()
	res = c.resolution
	c.rwmutex.RUnlock()
	return
}

func (c *Camera) SetPos(x float64, y float64) {
	c.rwmutex.Lock()
	c.pos.X = x
	c.pos.Y = y
	c.rwmutex.Unlock()
}

func (c *Camera) Translate(x float64, y float64) {
	c.rwmutex.Lock()
	c.pos.X += x
	c.pos.Y += y
	c.rwmutex.Unlock()
}

func (c *Camera) SetPosX(x float64) {
	c.rwmutex.Lock()
	c.pos.X = x
	c.rwmutex.Unlock()
}

func (c *Camera) SetPosY(y float64) {
	c.rwmutex.Lock()
	c.pos.Y = y
	c.rwmutex.Unlock()
}

func (c *Camera) GetPolygon() physics.Polygon {
	return physics.NewRectangle(c.pos.X, c.pos.Y, c.resolution.X, c.resolution.Y).ToPolygon()
}
