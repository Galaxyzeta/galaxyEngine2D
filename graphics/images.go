package graphics

import (
	"image"
	"image/png"
	"os"
	"strings"
)

type ImageFormat uint8

const (
	UNKNOWN ImageFormat = iota
	PNG
)

func GetImgFormat(fileName string) ImageFormat {
	split := strings.SplitN(fileName, ".", 2)
	if len(split) != 2 {
		return UNKNOWN
	}
	suffix := strings.ToLower(split[1])
	switch suffix {
	case "png":
		return PNG
	default:
		return UNKNOWN
	}
}

func ReadPng(fileName string) (image.Image, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(fp)
	return img, err
}
