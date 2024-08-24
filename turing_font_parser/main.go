package turingfontparser

type TuringFontParser struct {
}

func NewTuringFont() {
	var p = NewReader("./assets/font/comic_sans_ms.ttf")
	p.parse()
}
