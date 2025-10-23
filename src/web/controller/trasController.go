package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	savedFiles, err := uploadFiles(c, "file", "./uploads")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"files":   savedFiles,
	})
}

func uploadFiles(c *gin.Context, formField string, saveDir string) ([]string, error) {
	var savedFiles []string

	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	file, err := c.FormFile(formField)
	if err == nil {
		filename := filepath.Base(file.Filename)
		savePath := filepath.Join(saveDir, filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			return nil, fmt.Errorf("failed to save file: %v", err)
		}
		savedFiles = append(savedFiles, filename)
		return savedFiles, nil
	}
	form, err := c.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %v", err)
	}

	files := form.File[formField]
	for _, f := range files {
		filename := filepath.Base(f.Filename)
		savePath := filepath.Join(saveDir, filename)
		if err := c.SaveUploadedFile(f, savePath); err != nil {
			return nil, fmt.Errorf("failed to save file %s: %v", filename, err)
		}
		savedFiles = append(savedFiles, filename)
	}

	if len(savedFiles) == 0 {
		return nil, fmt.Errorf("no files uploaded")
	}

	return savedFiles, nil
}
