package validator

import (
	"bytes"
	"fmt"
)

const (
	maxImageSize = 1 << 20
)

func getAllowedImageTypes() []string {
	return []string{
		".jpg",
		".jpeg",
		".png",
	}
}

func ValidateType(imageType string) error {
	for _, t := range getAllowedImageTypes() {
		if imageType == t {
			return nil
		}
	}

	return fmt.Errorf("provided image type '%s' is invalid", imageType)
}

func IsSizeExceededMaxSize(chunk []byte, buffer bytes.Buffer) error {
	size := len(chunk)

	imageSize := buffer.Len()
	imageSize += size
	fmt.Println("image size: ", imageSize)
	if imageSize > maxImageSize {
		return fmt.Errorf("image is too large, max image size is %d", maxImageSize)
	}

	return nil
}
