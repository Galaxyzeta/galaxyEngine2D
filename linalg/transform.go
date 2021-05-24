package linalg

func Cam2OpenGLScreen(camSpacePoint Point2f, camResolution Vector2f64, windowResolution Vector2f64) Point2f {
	ratioX := windowResolution.X / camResolution.X
	ratioY := windowResolution.Y / camResolution.Y
	return Point2f{
		camSpacePoint.X*ratioX - 0.5,
		camSpacePoint.Y*ratioY - 0.5,
	}
}

func World2Cam(worldSpacePoint Point2f, camLTPos Point2f) Point2f {
	return Point2f{
		worldSpacePoint.X - camLTPos.X,
		worldSpacePoint.Y - camLTPos.Y,
	}
}

func World2Screen(worldSpacePoint Point2f, camLTPos Point2f, camResolution Vector2f64, windowResolution Vector2f64) Point2f {
	return Cam2OpenGLScreen(World2Cam(worldSpacePoint, camLTPos), camResolution, windowResolution)
}
