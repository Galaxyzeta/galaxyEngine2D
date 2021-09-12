package parser

import (
	"encoding/xml"
	"io"
	"os"
)

type LevelConfig struct {
	LevelMetas   LevelMetas   `xml:"level-metas"`
	LevelDetails LevelDetails `xml:"level-details"`
}

type LevelMetas struct {
	Static           string           `xml:"static"`
	CameraCount      int              `xml:"camera-count"`
	FrameMetas       FrameMetas       `xml:"frame-metas"`
	SpriteMetas      SpriteMetas      `xml:"sprite-metas"`
	ObjectMetas      ObjectMetas      `xml:"object-metas"`
	ApplicationMetas ApplicationMetas `xml:"application-metas"`
}

type LevelDetails struct {
	Scene []Scene `xml:"scene"`
}

type ObjectDetail struct {
	Name string `xml:"name,attr"`
	X    int64  `xml:"x,attr"`
	Y    int64  `xml:"y,attr"`
}

type CameraWrapper struct {
	Cameras []CameraDetail
}

type CameraDetail struct {
	RXYAttr
	Index int `xml:"index,attr"`
}

type Scene struct {
	SceneName     string              `xml:"name,attr"`
	SceneMetas    SceneMetas          `xml:"scene-metas"`
	ObjectDetails ObjectDetailWrapper `xml:"objects"`
}

type ObjectDetailWrapper struct {
	Objects []ObjectDetail `xml:"object"`
}

type SceneMetas struct {
	RoomSize RWHAttr       `xml:"room-size"`
	Cameras  CameraWrapper `xml:"cameras"`
}

type FrameMetas struct {
	Dirs []FrameDir `xml:"dir"`
}

type FrameDir struct {
	Name   string `xml:"name,attr"`
	Prefix string `xml:"prefix,attr"`
}

type SpriteMetas struct {
	Sprites []Sprite `xml:"sprite"`
}

type ObjectMetas struct {
	Objects []Object `xml:"object"`
}

type RWHAttr struct {
	W float64 `xml:"w,attr"`
	H float64 `xml:"h,attr"`
}

type RXYAttr struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
}

type Object struct {
	Name string `xml:"name,attr"`
}

type Sprite struct {
	Name   string  `xml:"name,attr"`
	Frames []Frame `xml:"frame"`
}

type Frame struct {
	Name string `xml:"name,attr"`
}

type ApplicationMetas struct {
	Resolution  RWHAttr `xml:"resolution"`
	FPS         FPS     `xml:"fps"`
	Parallelism int     `xml:"parallelism"`
	Title       string  `xml:"title"`
}

type FPS struct {
	Physics int `xml:"physics,attr"`
	Render  int `xml:"render,attr"`
}

func ParseGameLevelFile(filePath string) (ret *LevelConfig) {
	fp, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	err = xml.Unmarshal(data, &ret)
	if err != nil {
		panic(err)
	}
	return
}
