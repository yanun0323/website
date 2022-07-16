package service

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"website/internal/domain"
	"website/internal/model"
	"website/util"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	_GITHUB_TEMPLATE_DIR = "./internal/resource/template/"
	_CATEGORY_DEFAULT    = "未分類"
)

var (
	_GITHUB_TEMPLATE_URL = ""
	_GITHUB_MD_LIST_URL  = ""
	_GITHUB_MD_URL       = ""
	_SITE_URL            = ""
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
	_GITHUB_TEMPLATE_URL = viper.GetString("resource.template")
	_GITHUB_MD_LIST_URL = viper.GetString("resource.list")
	_GITHUB_MD_URL = viper.GetString("resource.markdown")
	_SITE_URL = viper.GetString("server.site.url")

	urls := []string{
		util.Url(_GITHUB_TEMPLATE_URL, "main.html"),
	}

	files := []string{
		util.Url(_GITHUB_TEMPLATE_DIR, "main.html"),
	}

	if skip := viper.GetBool("skip.download.template"); !skip {
		downloadTemplate(repo, urls, files)
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
			html.WithUnsafe(),
		))

	ls := repo.GetMarkdownList(_GITHUB_MD_LIST_URL)

	return Service{
		template: template.Must(template.ParseFiles(files[0])),
		mdReader: md,
		repo:     repo,
		l:        l,
		list:     ls,
		btn:      getButtonString(ls),
	}
}

func (svc *Service) SetHomePage(e *echo.Echo, m ...echo.MiddlewareFunc) {
	svc.setMdPage(e, "/", getUrl(_GITHUB_MD_URL, "home", ".md"), m...)
}

func (svc *Service) SetAllArticlePage(e *echo.Echo, m ...echo.MiddlewareFunc) {
	for _, name := range svc.list {
		svc.setMdPage(e, "/"+name, getUrl(_GITHUB_MD_URL, name, ".md"), m...)
	}
}

func (svc *Service) setPage(e *echo.Echo, path, url string, m ...echo.MiddlewareFunc) {
	e.GET(path, func(c echo.Context) error {
		buf, err := svc.repo.GetHttpBodyBuf(url)
		if err != nil {
			return c.HTML(http.StatusNotFound, err.Error())
		}
		return c.HTML(http.StatusOK, string(buf))
	}, m...)
}

func (svc *Service) setMdPage(e *echo.Echo, path, url string, m ...echo.MiddlewareFunc) {
	e.GET(path, func(c echo.Context) error {
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

	css, err := svc.repo.GetHttpBodyBuf(getUrl(_GITHUB_TEMPLATE_URL, "style", ".css"))

	buf := new(bytes.Buffer)
	if err := svc.mdReader.Convert(md, buf); err != nil {
		return "", err
	}

	html := new(bytes.Buffer)

	if err := svc.template.Execute(html,
		struct {
			Body   template.HTML
			Button template.HTML
			CSS    template.HTML
		}{
			Body:   template.HTML(buf.String()),
			Button: template.HTML(svc.btn),
			CSS:    template.HTML("<style>" + string(css) + "</style>"),
		}); err != nil {
		return "", fmt.Errorf("execute buf, %w", err)
	}
	replace := `img src="` + _GITHUB_MD_URL
	return strings.ReplaceAll(html.String(), `img src="./`, replace), nil
}

func mdUrl(host, name string) string {
	return host + name + ".md"
}

func getUrl(host, name, suffix string) string {
	return host + name + suffix
}

func downloadTemplate(repo domain.IRepository, urls, files []string) {
	l := log.Default()
	for i, u := range urls {
		body, err := repo.GetTemplate(u)
		if err != nil {
			l.Println(fmt.Errorf("get template, %w", err))
			continue
		}
		if err := saveToLocal(files[i], body, true); err != nil {
			l.Println(fmt.Errorf("download to local, %w", err))
			continue
		}
	}
}

func saveToLocal(name string, data []byte, skipZero bool) error {
	if skipZero && len(data) == 0 {
		return nil
	}

	if len(data) == 0 {
		return errors.New("data is empty")
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("open file %s, %w", name, err)
	}
	defer f.Close()

	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("truncate file %s, %w", name, err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("seek file %s, %w", name, err)
	}

	if _, err := fmt.Fprintf(f, "%s", data); err != nil {
		return fmt.Errorf("fprintf file %s, %w", name, err)
	}
	return nil
}

func getButtonString(list []string) string {
	result := make([]byte, 0, 30)
	hash := make(map[string]*model.Category, 0)

	result = append(result, model.NewSideButton("Home", _SITE_URL, true).Byte()...)
	hash[_CATEGORY_DEFAULT] = model.NewCategory(_CATEGORY_DEFAULT)
	for _, name := range list {
		raw, category, single := separator(name)
		if single {
			category = _CATEGORY_DEFAULT
		}
		_, exist := hash[category]
		if !exist {
			hash[category] = model.NewCategory(category)
		}
		hash[category].AppendButton(model.NewSideButton(raw, name, false))
	}
	for _, category := range hash {
		result = append(result, category.Byte()...)
	}

	return string(result)
}

func separator(str string) (name string, category string, single bool) {
	collection := strings.Split(str, "-")

	if len(collection) == 1 {
		return str, "", true
	}

	return collection[1], collection[0], false
}
