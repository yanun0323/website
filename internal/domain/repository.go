package domain

type IRepository interface {
	GetMarkdownList(url string) []string
	GetMarkdown(url string) ([]byte, error)
	GetTemplate(url string) ([]byte, error)
	GetHttpBodyBuf(url string) ([]byte, error)
}
