package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func StoreMultipartImage(r *http.Request, maxSize int64, fieldName string) (string, error) {
	if err := r.ParseMultipartForm(maxSize); err != nil {
		return "", err
	}

	file, handler, fileErr := r.FormFile(fieldName)
	if fileErr != nil {
		return "", fileErr
	}
	defer file.Close()

	localFile, err := os.Create("uploads/" + handler.Filename)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	if _, err := io.Copy(localFile, file); err != nil {
		return "", err
	}
	fmt.Println(path.Dir(path.Dir(handler.Filename)))
	return handler.Filename, nil
}
