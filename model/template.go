package model

import "html/template"

type CustomTemplate struct {
	Home           template.HTML
	Illustration   template.HTML
	Sketch         template.HTML
	About          template.HTML
	Animation2D    template.HTML
	Animation3D    template.HTML
	AnimationFX    template.HTML
	AnimationSpine template.HTML
}

func NewCustomTemplate(host string) CustomTemplate {
	return CustomTemplate{
		Home:           template.HTML(`<a name="#" id="0" class="MenuButton" href="` + host + `/home">`),
		Illustration:   template.HTML(`<a name="#" id="0" class="MenuButton" href="` + host + `/illustration">`),
		Sketch:         template.HTML(`<a name="#" id="0" class="MenuButton" href="` + host + `/sketch">`),
		About:          template.HTML(`<a name="#" id="0" class="MenuButton" href="` + host + `/about">`),
		Animation2D:    template.HTML(`<a name="$" id="0" class="MenuBranch" href="` + host + `/animation-2d">`),
		Animation3D:    template.HTML(`<a name="$" id="0" class="MenuBranch" href="` + host + `/animation-3d">`),
		AnimationFX:    template.HTML(`<a name="$" id="0" class="MenuBranch" href="` + host + `/animation-fx">`),
		AnimationSpine: template.HTML(`<a name="$" id="0" class="MenuBranch" href="` + host + `/animation-spine">`),
	}
}
