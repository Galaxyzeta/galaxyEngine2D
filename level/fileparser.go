package level

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
	FrameMetas       FrameMetas       `xml:"frame-metas"`
	SpriteMetas      SpriteMetas      `xml:"sprite-metas"`
	ObjectMetas      ObjectMetas      `xml:"object-metas"`
	ApplicationMetas ApplicationMetas `xml:"application-metas"`
}

type LevelDetails struct {
	ObjectDetails []ObjectDetail `xml:"object"`
}

type ObjectDetail struct {
	Name string `xml:"name,attr"`
	X    int64  `xml:"x,attr"`
	Y    int64  `xml:"y,attr"`
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

type Reolution struct {
	W float64 `xml:"w,attr"`
	H float64 `xml:"h,attr"`
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
	Resolution  Reolution `xml:"resolution"`
	FPS         FPS       `xml:"fps"`
	Parallelism int       `xml:"parallelism"`
	Title       string    `xml:"title"`
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
