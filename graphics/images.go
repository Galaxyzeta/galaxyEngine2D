package graphics

import (
	"fmt"
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

func ReadAllPngsUnderDirectory(dirPath string) (fileNames []string, images []image.Image, err error) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, nil, err
	}
	for _, entry := range dirEntries {
		fileNames = append(fileNames, entry.Name())
		img, err := ReadPng(fmt.Sprintf("%s/%s", dirPath, entry.Name()))
		if err != nil {
			return nil, nil, err
		}
		images = append(images, img)
	}
	return fileNames, images, nil
}
