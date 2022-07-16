package model

import "fmt"

type SideButton struct {
	Medium bool
	Name   string
	Url    string
}

func NewSideButton(name, url string, medium bool) SideButton {
	return SideButton{
		Medium: medium,
		Name:   name,
		Url:    url,
	}
}

func (b SideButton) String() string {
	class := "sidebar-button sidebar-button-small"
	if b.Medium {
		class = "sidebar-button sidebar-button-medium"
	}

	content := fmt.Sprintf(`
		<button class="%s">
			%s 
		</button>`, class, b.Name)

	return fmt.Sprintf(`
		<a href="%s"> %s
		</a>`, b.Url, content)
}

func (b SideButton) Byte() []byte {
	return []byte(b.String())
}
