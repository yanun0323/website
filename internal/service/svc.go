package service

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"website/internal/domain"
	"website/util"

	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	_GITHUB_TEMPLATE_DIR = "./asset/template/"
	_GITHUB_MD_LIST_URL  = "https://raw.githubusercontent.com/yanun0323/memo/main/article"
	_GITHUB_MD_URL       = "https://raw.githubusercontent.com/yanun0323/memo/main/markdown/"
)

type Service struct {
	template *template.Template
	mdReader goldmark.Markdown
	repo     domain.IRepository
	l        *log.Logger

	list []string
	btn  string
}

func NewService(repo domain.IRepository) Service {
	files := []string{
		util.Url(_GITHUB_TEMPLATE_DIR, "main.html"),
		util.Url(_GITHUB_TEMPLATE_DIR, "style.html"),
		util.Url(_GITHUB_TEMPLATE_DIR, "sidebar.html"),
		util.Url(_GITHUB_TEMPLATE_DIR, "markdown.html"),
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("xcode-dark"),
			),
		),
		goldmark.WithParser(goldmark.DefaultParser()),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		))

	ls := repo.GetMarkdownList(_GITHUB_MD_LIST_URL)

	return Service{
		template: template.Must(template.ParseFiles(files...)),
		mdReader: md,
		repo:     repo,
		l:        log.Default(),
		list:     ls,
		btn:      getButtonString(ls),
	}
}

func (svc *Service) SetHomePage(e *echo.Echo, m ...echo.MiddlewareFunc) {
	e.GET("/", func(c echo.Context) error {
		url := url(_GITHUB_MD_URL, "home")
		html, err := svc.getMarkdownHtml(url)
		if err != nil {
			return c.HTML(http.StatusNotFound, err.Error())
		}
		return c.HTML(http.StatusOK, html)
	}, m...)
}

func (svc *Service) SetAllArticlePage(e *echo.Echo, m ...echo.MiddlewareFunc) {
	for _, name := range svc.list {
		svc.setArticlePage(e, name, m...)
	}
}

func (svc *Service) setArticlePage(e *echo.Echo, name string, m ...echo.MiddlewareFunc) {
	e.GET("/"+name, func(c echo.Context) error {
		url := url(_GITHUB_MD_URL, name)
		html, err := svc.getMarkdownHtml(url)
		if err != nil {
			return c.HTML(http.StatusNotFound, err.Error())
		}
		return c.HTML(http.StatusOK, html)
	}, m...)
}

func (svc *Service) getMarkdownHtml(url string) (string, error) {
	md, err := svc.repo.GetMarkdown(url)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := svc.mdReader.Convert(md, buf); err != nil {
		return "", err
	}

	html := new(bytes.Buffer)

	if err := svc.template.Execute(html,
		struct {
			Body   template.HTML
			Button template.HTML
		}{
			Body:   template.HTML(buf.String()),
			Button: template.HTML(svc.btn),
		}); err != nil {
		return "", fmt.Errorf("execute buf, %w", err)
	}
	replace := `img src="` + _GITHUB_MD_URL
	return strings.ReplaceAll(html.String(), `img src="./`, replace), nil
}

func url(host, name string) string {
	return host + name + ".md"
}

func getButtonString(list []string) string {
	result := make([]byte, 0, 30)
	for _, str := range list {
		result = append(result, []byte(getButtonStr(str))...)
	}
	return string(result)
}

func getButtonStr(name string) string {
	return `<a id="0" href="http://localhost:8080/` +
		name + `"><button class="sidebar-button sidebar-font-small">` +
		"・" + name + ` </button></a>`
}
