package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/fajarardiyanto/steganography/config"
	"github.com/fajarardiyanto/steganography/models"
	"github.com/fajarardiyanto/steganography/pkg/bimg"
	"github.com/fajarardiyanto/steganography/pkg/qpdf"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
)

type RemoveMetadata struct{}

func New() *RemoveMetadata {
	return &RemoveMetadata{}
}

func (r *RemoveMetadata) RemoveMetadataService(c *gin.Context) {
	file, _, err := c.Request.FormFile("files")
	defer file.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	buffer := bytes.NewBuffer(nil)
	if _, err = io.Copy(buffer, file); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	imgs := bimg.NewImage(buffer.Bytes())
	size, _ := imgs.Size()
	imgt := imgs.Type()
	id := uuid.NewString()

	// Original
	fileName := fmt.Sprintf("/original/%s.%s", id, imgt)
	originalUrl, err := config.UploadFile(context.Background(), fileName, buffer.Bytes())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	imgType := bimg.PNG
	switch imgt {
	case bimg.ImageTypeName(bimg.JPEG):
		imgType = bimg.JPEG
	case bimg.ImageTypeName(bimg.PDF):
		fileName = fmt.Sprintf("/compressed/%s.%s", id, imgt)
		tempFile := fmt.Sprintf("tempfile-%s.%s", id, imgt)
		if err = r.PDFRemoveMetadata(buffer.Bytes(), tempFile, fileName); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, "OK PDF")
		return
	}

	newImage, err := imgs.
		Process(bimg.Options{
			RemoveAllMetaData: true,
			StripMetadata:     true,
			Width:             size.Width,
			Height:            size.Height,
			Embed:             true,
			Quality:           50,
			Compression:       9,
			Interpolator:      bimg.Bilinear,
			Trim:              true,
			Type:              bimg.JPEG,
			NoProfile:         true,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	dat, err := bimg.NewImage(newImage).Process(bimg.Options{
		Type:              imgType,
		RemoveAllMetaData: true,
		Trim:              true,
		StripMetadata:     true,
		NoProfile:         true,
		Embed:             true,
		Compression:       9,
		Quality:           50,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	fileName = fmt.Sprintf("/compressed/%s.%s", id, imgt)
	compresseUrl, err := config.UploadFile(context.Background(), fileName, dat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Error: false,
		Data: models.GCSResponse{
			OriginalURL: originalUrl,
			CompressURL: compresseUrl,
		},
	})
}

func (r *RemoveMetadata) PDFRemoveMetadata(file []byte, output, fileName string) error {
	f, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	if _, err = f.Write(file); err != nil {
		return err
	}

	q := qpdf.Init()
	defer q.GC()

	if err = q.Open(f.Name()); err != nil {
		return errors.New("err open file: " + err.Error())
	}

	err = q.SetOutput(output)
	if err != nil {
		return err
	}

	q.Linearize()
	q.Write()

	if q.HasError() {
		log.Println(q.LastError())
		return nil
	}

	buff, err := os.ReadFile(output)
	if err != nil {
		return err
	}

	if err = config.UploadToMinio(context.Background(), fileName, buff); err != nil {
		return err
	}

	defer os.Remove(output)

	return nil
}
