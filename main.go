package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"website/model"

	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	_HOST          = "http://localhost:8080"
	_HTML_FILEPATH = "./asset/html/portfolio"
	_MD_DIR        = "./asset/markdown/"
)

func main() {

	entries, err := os.ReadDir(_MD_DIR)
	if err != nil {
		fmt.Println(err)
	}
	btn := make([]byte, 0, len(entries))
	btn = append(btn, []byte(GetHomeButtonStr())...)
	btn = append(btn, []byte(GetRootButtonStr())...)

	for _, entry := range entries {
		if entry.IsDir() || entry.Name()[0] == '.' {
			continue
		}
		btn = append(btn, []byte(GetButtonStr(StrTakeOffMD(entry.Name())))...)
	}

	e := echo.New()
	btnStr := string(btn)
	SetHomePage(e, "/", btnStr)
	for _, entry := range entries {
		if entry.IsDir() || entry.Name()[0] == '.' {
			continue
		}
		name := StrTakeOffMD(entry.Name())
		SetPage(e, name, _MD_DIR, name, ".md", btnStr)
	}
	SetRawPage(e, "/home", _HTML_FILEPATH, _HOST)
	SetRawPage(e, "/illustration", _HTML_FILEPATH, _HOST)
	SetRawPage(e, "/sketch", _HTML_FILEPATH, _HOST)
	SetRawPage(e, "/about", _HTML_FILEPATH, _HOST)
	SetRawPage(e, "/animation-2d", _HTML_FILEPATH, _HOST)
	SetRawPage(e, "/animation-3d", _HTML_FILEPATH, _HOST)
	SetRawPage(e, "/animation-spine", _HTML_FILEPATH, _HOST)
	SetRawPage(e, "/animation-fx", _HTML_FILEPATH, _HOST)

	e.Logger.Fatal(e.Start(":8080"))
}
func SetHomePage(e *echo.Echo, path, btn string) {
	e.GET(path, func(c echo.Context) error {
		files := []string{
			"./asset/html/main.html",
			"./asset/html/markdown.html",
			"./asset/html/sidebar.html",
			"./asset/html/style.html",
		}

		tmpl := new(bytes.Buffer)
		t := template.Must(template.ParseFiles(files...))

		if err := t.Execute(tmpl,
			struct {
				Body   template.HTML
				Button template.HTML
			}{
				Body:   template.HTML(""),
				Button: template.HTML(btn),
			}); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("execute buf, %w", err).Error())
		}
		return c.HTML(http.StatusOK, tmpl.String())
	})
}

func StrTakeOffMD(str string) string {
	return strings.ReplaceAll(str, ".md", "")
}

func SetPage(e *echo.Echo, path, dir, filename, suffix, btn string) {
	e.GET(path, func(c echo.Context) error {
		buf := new(bytes.Buffer)

		file, err := os.ReadFile(dir + filename + suffix)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("read file, %w", err).Error())
		}

		reader := NewMarkdownReader()
		if err := reader.Convert(file, buf); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("convert markdown, %w", err).Error())
		}

		files := []string{
			"./asset/html/main.html",
			"./asset/html/markdown.html",
			"./asset/html/sidebar.html",
			"./asset/html/style.html",
		}

		tmpl := new(bytes.Buffer)
		t := template.Must(template.ParseFiles(files...))

		if err := t.Execute(tmpl,
			struct {
				Body   template.HTML
				Button template.HTML
			}{
				Body:   template.HTML(buf.String()),
				Button: template.HTML(btn),
			}); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("execute buf, %w", err).Error())
		}
		return c.HTML(http.StatusOK, tmpl.String())
	})
}

func GetRootButtonStr() string {
	return `<a href="http://localhost:8080/"> <button class="sidebar-button"> Root </button> </a>`
}

func GetHomeButtonStr() string {
	return `<a href="http://localhost:8080/home"> <button class="sidebar-button"> Home </button> </a>`
}

func GetButtonStr(filename string) string {
	return `<a href="http://localhost:8080/` +
		filename + `"> <button class="sidebar-button"> ` +
		filename + ` </button> </a>`
}

func NewMarkdownReader() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParser(goldmark.DefaultParser()),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
}

func SetRawPage(e *echo.Echo, path, filepath, host string) {
	e.GET(path, func(c echo.Context) error {
		buf := new(bytes.Buffer)
		t := template.Must(template.ParseFiles(filepath + path + ".html"))

		custom := model.NewCustomTemplate(host)
		if err := t.Execute(buf, custom); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("execute buf, %w", err).Error())
		}

		return c.HTML(http.StatusOK, buf.String())
	})
}
