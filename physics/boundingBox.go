package physics

import "galaxyzeta.io/engine/linalg"

type BoundingBox [4]linalg.Vector2f64

const BB_TopLeft = 0
const BB_TopRight = 1
const BB_BotLeft = 2
const BB_BotRight = 3
