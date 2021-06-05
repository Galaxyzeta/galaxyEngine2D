package file

import (
	"io"
	"os"
	"strings"
)

func OpenAndRead(file string) (string, error) {
	fp, err := os.Open(file)
	if err != nil {
		return "", err
	}
	buffer := make([]byte, 256)
	sb := strings.Builder{}
	for {
		_, err = fp.Read(buffer)
		sb.Write(buffer)
		if err != nil {
			if err == io.EOF {
				return sb.String(), nil
			}
			return "", err
		}
	}
}
