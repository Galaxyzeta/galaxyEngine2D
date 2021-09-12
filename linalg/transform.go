package linalg

// Note: about space convertion:
// The game world uses coordinate system which the origin point locates at top-left, having y-axis pointing down, and
// x-axis pointing right.
// First we need to clip an object into camera space, which has the same definition of coordinates as described above.
// Second we convert camera space to window space, an isometric scale is performed to do this.
// Finally we convert window space into OpenGL rendering space, which has a completely different coordinate system, with
// its origin point located at center of the window, and y-axis pointing upwards and x-axis pointing rightwards.

func World2Cam(worldSpacePoint Vector2f64, camLTPos Vector2f64) Vector2f64 {
	return Vector2f64{
		worldSpacePoint.X - camLTPos.X,
		worldSpacePoint.Y - camLTPos.Y,
	}
}

func Cam2OpenGL(camSpacePoint Vector2f64, camResolution Vector2f64, windowResolution Vector2f64) Vector2f64 {
	ratioX := windowResolution.X / camResolution.X
	ratioY := windowResolution.Y / camResolution.Y
	return ScreenNormalizeToOpenGL(Vector2f64{
		camSpacePoint.X * ratioX,
		camSpacePoint.Y * ratioY,
	}, windowResolution)
}

func ScreenNormalizeToOpenGL(p Vector2f64, windowResolution Vector2f64) Vector2f64 {
	return Vector2f64{
		p.X*2/windowResolution.X - 1,
		1 - p.Y*2/windowResolution.Y,
	}
}

func World2OpenGL(worldSpacePoint Vector2f64, camLTPos Vector2f64, camResolution Vector2f64, windowResolution Vector2f64) Vector2f64 {
	return Cam2OpenGL(World2Cam(worldSpacePoint, camLTPos), camResolution, windowResolution)
}

func WorldVertice2OpenGL(arr *[]float64, offset int, stride int, camLTPos Vector2f64, camResolution Vector2f64, windowResolution Vector2f64) {
	var pos int = offset
	for pos < len(*arr) {
		spos := World2OpenGL(NewVector2f64((*arr)[pos], (*arr)[pos+1]), camLTPos, camResolution, windowResolution)
		(*arr)[pos] = spos.X
		(*arr)[pos+1] = spos.Y
		pos += stride
	}
}
