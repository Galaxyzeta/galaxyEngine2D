package graphics

var spriteMap map[string]*Sprite
var shaderMap map[string]*Shader

func init() {
	shaderMap = make(map[string]*Shader, 0)
}
