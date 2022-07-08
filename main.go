package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"website/model"

	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	_HOST          = "http://localhost:8080"
	_HTML_FILEPATH = "./asset/html/portfolio"
	_MD_DIR        = "./asset/markdown/"
	_GITHUB_MD_URL = "https://raw.githubusercontent.com/yanun0323/memo/main/markdown/"
)

func main() {

	entries, err := os.ReadDir(_MD_DIR)
	if err != nil {
		fmt.Println(err)
	}
	btn := make([]byte, 0, len(entries))

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
		// SetPage(e, name, _MD_DIR, name, ".md", btnStr)
		SetHttpPage(e, name, _GITHUB_MD_URL, name, ".md", btnStr)
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

func SetHttpPage(e *echo.Echo, path, githubUrl, filename, suffix, btn string) {
	e.GET(path, func(c echo.Context) error {
		buf := new(bytes.Buffer)
		url := githubUrl + filename + suffix
		response, err := http.Get(url)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("get http file, %w", err).Error())
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("http status, %d", response.StatusCode).Error())
		}

		file, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("read http body, %w", err).Error())
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

func GetButtonStr(filename string) string {
	return `<a id="0" href="http://localhost:8080/` +
		filename + `"> <button class="sidebar-button sidebar-font-small"> ` +
		"・" + filename + ` </button> </a>`
}

func NewMarkdownReader() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("xcode-dark"),
			),
		),
		goldmark.WithParser(goldmark.DefaultParser()),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
}

func SetRawPage(e *echo.Echo, path, filepath, host string) {
	e.GET(path, func(c echo.Context) error {
		buf := new(bytes.Buffer)
		t := template.Must(template.ParseFiles(filepath+path+".html", filepath+"/back-to-home.html"))

		custom := model.NewCustomTemplate(host)
		if err := t.Execute(buf, custom); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("execute buf, %w", err).Error())
		}

		return c.HTML(http.StatusOK, buf.String())
	})
}
