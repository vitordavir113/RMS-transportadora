package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const templatesDir = "web/templates/"

var funcMap = template.FuncMap{
	// formatLitros formata 8000 -> "8.000" e 8000.5 -> "8.000,5" (padrão BR).
	"formatLitros": formatLitros,
	"formatDate":   func(t time.Time) string { return t.Format("02/01/2006 15:04") },
	"formatDatePtr": func(t *time.Time) string {
		if t == nil {
			return "-"
		}
		return t.Format("02/01/2006 15:04")
	},
}

func formatLitros(f float64) string {
	inteiro := int64(f)
	decimal := f - float64(inteiro)

	s := strconv.FormatInt(inteiro, 10)
	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}

	var groups []string
	for len(s) > 3 {
		groups = append([]string{s[len(s)-3:]}, groups...)
		s = s[:len(s)-3]
	}
	groups = append([]string{s}, groups...)
	result := strings.Join(groups, ".")
	if neg {
		result = "-" + result
	}

	if decimal > 0.009 {
		result = fmt.Sprintf("%s,%d", result, int64(decimal*10))
	}
	return result
}

// RenderPage renderiza uma página completa (layout + sidebar + navbar + conteúdo).
func RenderPage(c *gin.Context, page string, data gin.H) {
	files := []string{
		templatesDir + "layout.html",
		templatesDir + "partials/sidebar.html",
		templatesDir + "partials/navbar.html",
		templatesDir + page,
	}

	tmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		c.String(http.StatusInternalServerError, "erro ao carregar template: %v", err)
		return
	}

	if data == nil {
		data = gin.H{}
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(c.Writer, "layout", data); err != nil {
		c.String(http.StatusInternalServerError, "erro ao renderizar template: %v", err)
	}

}

// RenderFragment renderiza apenas um fragmento (usado em respostas HTMX),
// sem o layout completo.
func RenderFragment(c *gin.Context, file string, blockName string, data interface{}) {
	tmpl, err := template.New(blockName).Funcs(funcMap).ParseFiles(templatesDir + file)
	if err != nil {
		c.String(http.StatusInternalServerError, "erro ao carregar fragmento: %v", err)
		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(c.Writer, blockName, data); err != nil {
		c.String(http.StatusInternalServerError, "erro ao renderizar fragmento: %v", err)
	}
}

func RenderAuthPage(c *gin.Context, page string, data gin.H) {
	files := []string{
		templatesDir + "auth/layout.html",
		templatesDir + page,
	}

	tmpl, err := template.New("auth_layout.html").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		c.String(http.StatusInternalServerError, "erro ao carregar template: %v", err)
		return
	}

	if data == nil {
		data = gin.H{}
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(c.Writer, "auth_layout", data); err != nil {
		c.String(http.StatusInternalServerError, "erro ao renderizar template: %v", err)
	}
}
