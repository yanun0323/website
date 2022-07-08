package domain

type IRepository interface {
	GetMarkdownList(url string) []string
	GetMarkdown(url string) ([]byte, error)
}
