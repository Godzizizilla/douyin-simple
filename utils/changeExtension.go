package utils

import (
	"path/filepath"
	"strings"
)

func ChangeExtension(filePath, newExt string) string {
	ext := filepath.Ext(filePath)
	return strings.TrimSuffix(filePath, ext) + newExt
}
