package route

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/testy/lightf/pkg/slice"
	"github.com/testy/lightf/pkg/token"
	"github.com/testy/lightf/server/auth"
	"github.com/testy/lightf/server/storage"
)

func Upload(cfg *auth.Config, gen *token.Generator, store *storage.Storager) gin.HandlerFunc {
	return func(c *gin.Context) {
		upload(c, cfg, gen, store)
	}
}

// /upload is an POST endpoint.
//
// this endpoint accepts a file to be uploaded
// and puts it into the storage, if the data
// from the request is valid.
//
// the request needs:
// - a header field "Token" with the auth token
// - a valid content-length
// - a valid file content-type
func upload(c *gin.Context, cfg *auth.Config, gen *token.Generator, store *storage.Storager) {
	token := c.Request.Header.Get("Token")
	length := c.Request.ContentLength

	// ================================
	// check token
	// ================================
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no token provided",
		})
		return
	}

	user := cfg.GetUser(token)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}

	// ================================
	// check length
	// ================================
	if length <= 0 {
		c.Status(http.StatusLengthRequired)
		return
	}
	if length > user.Size {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"message": fmt.Sprintf("expected content length <= %d byte(s)", user.Size),
		})
		return
	}

	// ================================
	// check file in body
	// ================================
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "could not decode multipart",
			"error":   err.Error(),
		})
		return
	}

	files := form.File[user.FileKey]
	if s := len(files); s != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": fmt.Sprintf("file amount should not be %d", s),
		})
		return
	}

	f := files[0]
	f0, err := f.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not open file from body",
			"error":   err.Error(),
		})
		return
	}

	// ================================
	// check content-type
	// ================================
	cType := getContentType(f)
	if !slice.ContainsString(user.Types, cType) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": fmt.Sprintf("content type %v not allowed", cType),
		})
		return
	}

	name, err := gen.Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not generate file name",
			"error":   err.Error(),
		})
		return
	}

	addr, err := store.Store(name.Str, user.Expire, f0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not store file",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"address": addr,
		"key":     name.Str,
	})
}

// workaround for file Content-Type headers
// which contain multiple values such as "...; charset=utf-8"
func getContentType(f *multipart.FileHeader) string {
	c := f.Header.Get("Content-Type")
	if strings.Contains(c, ";") {
		c = strings.Split(c, ";")[0]
	} else if strings.Contains(c, ",") {
		c = strings.Split(c, ",")[0]
	}
	return c
}
