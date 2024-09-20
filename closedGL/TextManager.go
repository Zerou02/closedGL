package closedGL

import (
	"github.com/EngoEngine/glm"
)

type TextManager struct {
	reader Reader
	ctx    *ClosedGLContext
}

func newTextManager(path string, ctx *ClosedGLContext) TextManager {
	return TextManager{
		reader: NewReader(path, ctx),
		ctx:    ctx,
	}
}

func (this *TextManager) readGlyf(unicode uint32) Glyf {
	return this.reader.ReadGlyf(unicode)
}

func (this *TextManager) drawText(x, y float32, size float32, text string, triMesh *TriangleMesh) {
	var base = this.readGlyf('a')
	var scaleFactor = base.CalcScaleFactor(size)
	var glyfs = []Glyf{}
	for _, x := range text {
		glyfs = append(glyfs, this.readGlyf(uint32(x)))
	}

	this.setAdvanceWidth(&glyfs)
	for i := 0; i < len(glyfs); i++ {
		glyfs[i].Scale(scaleFactor)
		glyfs[i].AddOffset(glm.Vec2{x, -y + this.ctx.Window.Wh})
		this.drawGlyf(&glyfs[i], this.ctx, triMesh)
	}
}

func (this *TextManager) drawGlyf(glyf *Glyf, ctx *ClosedGLContext, triMesh *TriangleMesh) {
	for _, x := range glyf.SimpleGlyfs {
		drawSimpleGlyf(x, ctx, triMesh)
	}
}

func (this *TextManager) setAdvanceWidth(glyfs *[]Glyf) {
	var offset float32 = 0
	for i := 0; i < len(*glyfs); i++ {
		for j := 0; j < len((*glyfs)[i].SimpleGlyfs); j++ {
			(*glyfs)[i].SimpleGlyfs[j].AddOffset(glm.Vec2{offset, 100})
		}
		offset += (*glyfs)[i].AdvanceWidth
	}
}
