package nimblemossbilliard

import (
	"image/png"
	"os"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type NimbleMoss struct {
	ctx                     *closedGL.ClosedGLContext
	circles                 []Circle
	lines                   [][]glm.Vec2
	currImgIdx, currLineIdx int
	mesh                    *closedGL.CircleMesh
	lineMesh                *closedGL.LineMesh
	rectMesh                *closedGL.RectangleMesh
	paths                   []string
	speedSlider, gravSlider closedGL.Slider
	speed, grav             float32
	currLine                []glm.Vec2
	pMesh                   *closedGL.PixelMesh
	cols                    int
	lineColours             [][]glm.Vec4
}

func NewNimbleMoss(ctx *closedGL.ClosedGLContext) NimbleMoss {
	var middleY = ctx.Window.Wh - 50
	var middleX = ctx.Window.Ww / 2
	var moreThanAbs = []glm.Vec2{{0, 250}, {middleX, middleY}, {middleX, middleY}, {800, 250}}
	var Abs = []glm.Vec2{{0, 150}, {middleX, middleY}, {middleX, middleY}, {800, 150}}
	var lessThanAbs = []glm.Vec2{{200, 0}, {middleX, middleY}, {middleX, middleY}, {600, 0}}
	var straight = []glm.Vec2{{0, middleY - 50}, {800, middleY - 50}}
	var lines = [][]glm.Vec2{moreThanAbs, Abs, lessThanAbs, straight}

	var mesh = ctx.CreateCircleMesh()
	var lineMesh = ctx.CreateLineMesh()
	var rectMesh = ctx.CreateRectMesh()
	var pMesh = ctx.CreatePixelMesh()
	pMesh.SetPixelSize(1)

	var first = []glm.Vec4{{1, 0, 0, 1}, {0, 0, 1, 1}}
	var second = []glm.Vec4{{0, 0, 1, 1}, {0, 1, 0, 1}}
	var colours = [][]glm.Vec4{first, second}

	var moose = NimbleMoss{
		ctx:         ctx,
		circles:     []Circle{},
		lines:       lines,
		mesh:        &mesh,
		lineMesh:    &lineMesh,
		rectMesh:    &rectMesh,
		pMesh:       &pMesh,
		speedSlider: closedGL.NewSlider(&rectMesh, ctx, glm.Vec4{20, 550, 75, 20}),
		gravSlider:  closedGL.NewSlider(&rectMesh, ctx, glm.Vec4{620, 550, 75, 20}),
		currImgIdx:  0,
		currLineIdx: 0,
		speed:       260 * 0.5,
		grav:        0.31,
		lineColours: colours,
	}
	moose.loadPaths("./nimble_moss_billiard/forms/")
	moose.loadLines()
	moose.loadColourBalls(1)
	//moose.loadImage(0.1)
	return moose
}

func (this *NimbleMoss) loadColourBalls(amount int) []Circle {
	var vel = glm.Vec2{0, 1}
	var circles = []Circle{}
	var red = glm.Vec4{1, 0, 0, 1}
	var blue = glm.Vec4{0, 0, 1, 1}
	for i := 0; i < amount; i++ {
		var c = newCircle(glm.Vec2{300 + float32(i)*0.02, 230}, vel.Normalized(), 10, closedGL.LerpVec4(red, blue, float32(i)/float32(amount)))
		c.drawInto(this.mesh)
		this.circles = append(this.circles, c)
	}
	return circles
}

func (this *NimbleMoss) ProcessNTimes(startPos glm.Vec2, amountBounces int, delta float32) {
	var c = &this.circles[0]
	c.pos = startPos
	var cols = 0
	for cols < amountBounces {
		c.process(delta, this.speed, this.grav)
		//this.pMesh.AddPixel(c.getCentre(), glm.Vec4{1, 1, 1, 1})
		var col, refl = this.didCircleCollide(c)
		if col {
			cols++
			c.vel = refl
			c.pos[1] -= 1
		}
	}
	var centre = c.getCentre()
	var left = this.currLine[0]
	var middle = this.currLine[1]
	var right = this.currLine[3]
	var cType = this.lineColours[0]
	if centre[0] > middle[0] {
		left = middle
		middle = right
		cType = this.lineColours[1]
	}
	var percentage = closedGL.CalcPercentage(left[0], middle[0], centre[0])
	this.pMesh.AddPixel(startPos, closedGL.LerpVec4(cType[0], cType[1], percentage))
	this.pMesh.Copy()
}

func (this *NimbleMoss) ProcessUntil(amountBounces int, delta float32) {
	var c = &this.circles[0]
	if this.cols < amountBounces {
		this.pMesh.AddPixel(c.pos, glm.Vec4{1, 1, 1, 1})
		c.process(delta, this.speed, this.grav)
		var col, refl = this.didCircleCollide(c)
		if col {
			this.cols++
			c.vel = refl
			c.pos[1] -= 1
		}
	}
	this.pMesh.Copy()
}

func (this *NimbleMoss) loadPaths(dirPath string) {
	var paths = []string{}
	var dir, err = os.ReadDir(dirPath)
	if err != nil {
		panic(err.Error())
	}
	for _, x := range dir {
		paths = append(paths, dirPath+x.Name())
	}
	this.paths = paths
}

func (this *NimbleMoss) loadImage(scaleFactor float32) {
	var vel = glm.Vec2{0, 1}
	this.circles = []Circle{}
	this.mesh.Clear()
	var f, _ = os.Open(this.paths[this.currImgIdx])
	var img, _ = png.Decode(f)
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			var r, g, b, a = img.At(x, y).RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			a = a >> 8
			var cVec = glm.Vec4{float32(r) / 255.0, float32(g) / 255.0, float32(b) / 255.0, float32(a) / 255.0}
			if a == 0 {
				continue
			}
			this.circles = append(this.circles, newCircle(glm.Vec2{300 + scaleFactor*float32(x), 20 + scaleFactor*float32(y)}, vel.Normalized(), 10, cVec))
		}
	}
}

func (this *NimbleMoss) loadLines() {
	var line = this.lines[this.currLineIdx]
	this.lineMesh.Clear()

	for i := 0; i < len(line); i += 2 {
		this.lineMesh.AddLine(line[i], line[i+1], this.lineColours[i/2][0], this.lineColours[i/2][1])
	}
	this.lineMesh.CopyToGPU()
	this.currLine = line
}

func (this *NimbleMoss) Process(delta float32) {

	this.speedSlider.Process()
	this.gravSlider.Process()

	var newIdx = this.currImgIdx
	if this.ctx.IsKeyPressed(glfw.KeyD) {
		newIdx++
	}
	if this.ctx.IsKeyPressed(glfw.KeyA) {
		newIdx--
	}
	newIdx = int(closedGL.Clamp(0, float32(len(this.paths)-1), float32(newIdx)))
	if newIdx != this.currImgIdx {
		this.currImgIdx = newIdx
		this.loadImage(0.1)
	}

	var newLineIdX = this.currLineIdx
	if this.ctx.IsKeyPressed(glfw.KeyE) {
		newLineIdX++
	}
	if this.ctx.IsKeyPressed(glfw.KeyQ) {
		newLineIdX--
	}
	newLineIdX = int(closedGL.Clamp(0, float32(len(this.lines))-1, float32(newLineIdX)))
	if newLineIdX != this.currLineIdx {
		this.currLineIdx = newLineIdX
		this.loadLines()
		this.loadImage(0.1)
	}
	this.speed = 260 * this.speedSlider.GetPercentage()
	this.grav = 0.62 * this.gravSlider.GetPercentage()
	for i := 0; i < len(this.circles); i++ {
		var c = &this.circles[i]
		c.process(delta, this.speed, this.grav)
		var col, refl = this.didCircleCollide(c)
		if col {
			c.vel = refl
			c.pos[1] -= 1
		}
	}
}

func (this *NimbleMoss) didCircleCollide(circle *Circle) (bool, glm.Vec2) {
	var col = false
	var reflVec glm.Vec2
	for j := 0; j < len(this.currLine); j += 2 {
		var p1 = this.currLine[j]
		var p2 = this.currLine[j+1]
		var closestPoint = closedGL.ClosestPointOnLine(p1, p2, circle.getCentre())
		var dist = closestPoint.Sub(&circle.pos)
		if dist.Len() < circle.radius {
			col = true
			reflVec = closedGL.GetReflectionVec(p1, p2, circle.vel)
			break
		}
	}
	return col, reflVec
}

func (this *NimbleMoss) Draw(delta float32) {
	this.mesh.Clear()
	for i := 0; i < len(this.circles); i++ {
		this.circles[i].drawInto(this.mesh)
	}

	this.mesh.Copy()
	this.rectMesh.Draw()
	this.lineMesh.Draw()
	this.pMesh.Draw()
	this.mesh.Draw()
}
