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
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	_GITHUB_TEMPLATE_DIR = ".internal/asset/template/"
)

var (
	_GITHUB_MD_LIST_URL = ""
	_GITHUB_MD_URL      = ""
	_SITE_URL           = ""
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
	l := log.Default()
	_GITHUB_MD_LIST_URL = viper.GetString("resource.list")
	_GITHUB_MD_URL = viper.GetString("resource.markdown")
	_SITE_URL = viper.GetString("server.site.url")
	l.Print(_GITHUB_MD_LIST_URL)
	l.Print(_GITHUB_MD_URL)
	l.Print(_SITE_URL)

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
		l:        l,
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

	result = append(result, []byte(getHomeButtonStr())...)
	for _, str := range list {
		result = append(result, []byte(getButtonStr(str))...)
	}
	return string(result)
}

func getHomeButtonStr() string {
	return `<a href="` + _SITE_URL + `">
        <button class="sidebar-button sidebar-font-medium"> Home </button>
    </a>`
}

func getButtonStr(name string) string {
	return `<a id="0" href="` +
		name + `"><button class="sidebar-button sidebar-font-small">` +
		"・" + name + ` </button></a>`
}
