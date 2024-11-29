package controller

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func IndexAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		fp := filepath.Join("views", "adminIndex.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		err = tmpl.Execute(c.Writer, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
}
