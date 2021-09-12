package core

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/infra"
	"galaxyzeta.io/engine/parser"
	"galaxyzeta.io/engine/physics"
)

type injectContext struct {
	iobj                       base.IGameObject2D
	toInject                   reflect.Value
	delayedInjectionFieldIndex []int
	cachedTf                   *component.Transform2D
}

func Inject(iobj base.IGameObject2D) {
	injCtx := injectContext{
		iobj:                       iobj,
		toInject:                   reflect.ValueOf(iobj),
		delayedInjectionFieldIndex: []int{},
		cachedTf:                   &component.Transform2D{},
	}
	// 1st inject attempt
	injectAll(&injCtx)
	// inject delayed:
	for _, i := range injCtx.delayedInjectionFieldIndex {
		doInject(&injCtx, i, false)
	}
}

func injectAll(injCtx *injectContext) {
	for i := 0; i < injCtx.toInject.NumField(); i++ {
		doInject(injCtx, i, true)
	}
}

func doInject(injCtx *injectContext, i int, tolerateMissingDep bool) {
	ifv := injCtx.toInject
	ift := ifv.Type()
	fdt := ift.Field(i)
	fdv := ifv.Field(i)
	iobj := injCtx.iobj

	if fdt.Anonymous && ift.Kind() == reflect.Struct {
		if fdt.Type.Kind() == reflect.Ptr {
			injCtx.toInject = reflect.ValueOf(fdv.Interface()).Elem()
		} else {
			injCtx.toInject = reflect.ValueOf(fdv.Interface())
		}
		injectAll(injCtx)
		injCtx.toInject = ifv
	}

	tags := fdt.Tag
	tag := tags.Get("gxen")
	attrs := strings.Split(tag, "|")
	switch attrs[0] {
	case "tf":
		injectValue(&fdv, component.NewTransform2D(), iobj)
	case "rb":
		injectValue(&fdv, component.NewRigidBody2D(), iobj)
	case "sr":
		if injCtx.cachedTf == nil {
			if !tolerateMissingDep {
				panic("required Transform2D for SpriteRenderer2D is not found")
			}
			injCtx.delayedInjectionFieldIndex = append(injCtx.delayedInjectionFieldIndex, i)
			return
		}
		animator, isStatic, opts := mustParseSr(attrs)
		injectValue(&fdv, component.NewSpriteRendererWithOptions(animator, injCtx.cachedTf, isStatic, opts), iobj)
	}
}

func injectValue(val *reflect.Value, any interface{}, iobj base.IGameObject2D) {
	replica := reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr()))
	replica.Set(reflect.ValueOf(any))
	iobj.Obj().RegisterComponentIfAbsent(any.(base.IComponent))
}

func mustParseSr(attrs []string) (animator *graphics.Animator, isStatic bool, options graphics.RenderOptions) {
	// "sr|clips=clip,status,clip,status,...|static=true|pivot=tl"
	clipPairs := make([]graphics.StateClipPair, 0)
	params := mustResolveParams(attrs)
	clipPairDefs, ok := params["clips"]
	if ok {
		terms := strings.Split(clipPairDefs, ",")
		if len(terms)&1 == 1 {
			panic("clip definition length must be even")
		}
		for i, j := 0, 1; j < len(terms); i, j = i+2, j+2 {
			clipPairs = append(clipPairs, graphics.StateClipPair{
				State: terms[j],
				Clip:  graphics.NewSpriteInstance(terms[i]),
			})
		}
	}
	isStatic = false
	isStaticAttr, ok := params["static"]
	if ok && (isStaticAttr == "true" || isStaticAttr == "1") {
		isStatic = true
	}
	// handle pivot
	pivot, ok := params["pivot"]
	if ok {
		pivotEnum := physics.PivotOption_Disable
		switch infra.TrimSpace(pivot) {
		case "tl":
			pivotEnum = physics.PivotOption_TopLeft
		case "tc":
			pivotEnum = physics.PivotOption_TopCenter
		case "tr":
			pivotEnum = physics.PivotOption_TopRight
		case "cl":
			pivotEnum = physics.PivotOption_CenterLeft
		case "c":
			pivotEnum = physics.PivotOption_Center
		case "cr":
			pivotEnum = physics.PivotOption_CenterRight
		case "bl":
			pivotEnum = physics.PivotOption_BottomLeft
		case "bc":
			pivotEnum = physics.PivotOption_BottomCenter
		case "br":
			pivotEnum = physics.PivotOption_BottomRight
		}
		options.Pivot = &physics.Pivot{
			Option: pivotEnum,
		}
	}
	// handle scale
	scaleAttr, ok := params["scale"]
	if ok {
		s := parser.MustParseNumericStringTuple(scaleAttr)
		options.Scale = &s
	}
	animator = graphics.NewAnimator(clipPairs...)
	return
}

func mustResolveParams(attrs []string) (ret map[string]string) {
	ret = make(map[string]string)
	for _, attr := range attrs {
		eq := strings.Split(attr, "=")
		if len(eq) != 2 {
			panic(fmt.Sprintf("cannot parse %s, must be like a = b pattern", attr))
		}
		ret[infra.TrimSpace(eq[0])] = infra.TrimSpace(eq[1])
	}
	return
}
