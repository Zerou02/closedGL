package turingfontparser

import "github.com/Zerou02/closedGL/closedGL"

type TuringFontParser struct {
	reader Reader
	ctx    *closedGL.ClosedGLContext
}

func NewTuringFont(fontPath string, ctx *closedGL.ClosedGLContext) TuringFontParser {
	var p = NewReader(fontPath, ctx)
	p.init()
	return TuringFontParser{
		reader: p,
		ctx:    ctx,
	}
}

func (this *TuringFontParser) ParseGlyf(unicodeVal uint32, scale float32) Glyf {
	return this.reader.readGlyf(unicodeVal, scale)
}
