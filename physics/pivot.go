package physics

import "galaxyzeta.io/engine/linalg"

type PivotOption int8

const (
	PivotOption_Disable PivotOption = iota
	PivotOption_TopLeft
	PivotOption_TopCenter
	PivotOption_TopRight
	PivotOption_CenterLeft
	PivotOption_Center
	PivotOption_CenterRight
	PivotOption_BottomLeft
	PivotOption_BottomCenter
	PivotOption_BottomRight
)

type Pivot struct {
	Option PivotOption
	Point  *linalg.Vector2f64
}

func (o PivotOption) GetPivotPoint(bb BoundingBox) linalg.Vector2f64 {
	rect := bb.ToRectangle()
	switch o {
	case PivotOption_Disable, PivotOption_TopLeft:
		return bb[BB_TopLeft]
	case PivotOption_TopRight:
		return bb[BB_TopRight]
	case PivotOption_BottomLeft:
		return bb[BB_BotLeft]
	case PivotOption_BottomRight:
		return bb[BB_BotRight]
	case PivotOption_TopCenter:
		return linalg.NewVector2f64(rect.Left+rect.Width/2, rect.Top)
	case PivotOption_CenterLeft:
		return linalg.NewVector2f64(rect.Left, rect.Top+rect.Height/2)
	case PivotOption_Center:
		return linalg.NewVector2f64(rect.Left+rect.Width/2, rect.Top+rect.Height/2)
	case PivotOption_CenterRight:
		return linalg.NewVector2f64(rect.Left+rect.Width, rect.Top+rect.Height/2)
	case PivotOption_BottomCenter:
		return linalg.NewVector2f64(rect.Left+rect.Width/2, rect.Top+rect.Height)
	}
	panic("unknown pivot type")
}
