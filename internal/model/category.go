package model

import "fmt"

type Category struct {
	Name    string
	Buttons []SideButton
}

func NewCategory(name string) *Category {
	return &Category{
		Name:    name,
		Buttons: make([]SideButton, 0, 10),
	}
}

func (c *Category) AppendButton(button SideButton) {
	c.Buttons = append(c.Buttons, button)
}

func (c *Category) Byte2() []byte {
	// <a href="javascript:;" onclick="document.getElementById('hide123').style.display=(document.getElementById('hide123').style.display=='none'?'':'none');">點我展開／隱藏(方法四)</a>

	// <span id="hide123" style="display:none">哎呀又又被發現了啦(つд⊂) </span>

	result := make([]byte, 0, 50)

	result = append(result, `<div class="sidebar-category">
	`...)
	result = append(result, []byte(c.Name)...)

	for _, button := range c.Buttons {
		result = append(result, button.Byte()...)
	}

	result = append(result, `
	</div>`...)

	return result
}

func (c *Category) Byte() []byte {
	result := make([]byte, 0, 50)
	rowCategory := fmt.Sprintf(`
		<a href="javascript:;" 
			onclick="document.getElementById('%s').style.display=(document.getElementById('%s').style.display=='none'?'block':'none');">
			<button class="sidebar-category">
				%s
			</button>
		</a>`, c.Name, c.Name, c.Name)
	rowButtonBlockPrefix := fmt.Sprintf(`
	<div id="%s">`, c.Name)

	result = append(result, rowCategory...)
	result = append(result, rowButtonBlockPrefix...)

	for _, button := range c.Buttons {
		result = append(result, button.Byte()...)
	}

	result = append(result, `
	</div>`...)

	return result
}
