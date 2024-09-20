package main

import (
	"image/png"
	_ "image/png"
	"os"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

func main() {
	//StartClosedGL()
	/* StartTuwuing() */
	NimbleMoss()
}

type Circle struct {
	pos                  glm.Vec2
	vel                  glm.Vec2
	radius, speed, gravY float32
	circleMesh           *closedGL.CircleMesh
	colour               glm.Vec4
}

func (this *Circle) process(delta float32) {
	this.move(delta)
}

func (this *Circle) move(delta float32) {
	this.vel[1] += this.gravY * delta
	this.pos.AddScaledVec(delta*this.speed, &this.vel)

}

func (this *Circle) draw() {
	this.circleMesh.AddCircle(this.pos, this.colour, this.colour, this.radius, 0)
}

func newCircle(pos, vel glm.Vec2, radius, speed float32, mesh *closedGL.CircleMesh, colour glm.Vec4) Circle {
	return Circle{
		pos:        pos,
		vel:        vel,
		radius:     radius,
		speed:      speed,
		circleMesh: mesh,
		colour:     colour,
		gravY:      0.31,
	}
}

// kp, warum das in ss funktioniert
func getReflectionVec(p1, p2, vec glm.Vec2) glm.Vec2 {
	var line = p2.Sub(&p1)
	var perp = line.Perp()
	var dot = vec.Dot(&perp) / (perp.Dot(&perp))
	var projNV = perp.Mul(dot)
	projNV = projNV.Mul(-2)
	var retVec = vec.Add(&projNV)
	retVec.Normalize()

	return retVec
}

func loadColourBall(mesh *closedGL.CircleMesh) []Circle {
	var vel = glm.Vec2{0, 1}
	var circles = []Circle{}
	var red = glm.Vec4{1, 0, 0, 1}
	var blue = glm.Vec4{0, 0, 1, 1}
	for i := 0; i < 5000; i++ {
		circles = append(circles, newCircle(glm.Vec2{300 + float32(i)*0.02, 230}, vel.Normalized(), 10, 130, mesh, closedGL.LerpVec4(red, blue, float32(i)/5000)))
	}
	return circles
}

func loadImage(mesh *closedGL.CircleMesh) []Circle {
	var vel = glm.Vec2{0, 1}
	var circles = []Circle{}
	var f, _ = os.Open("./nimble_moss_billiard/apple2.png")
	var img, _ = png.Decode(f)
	println(img.Bounds().Max.X)
	println(img.At(0, 0).RGBA())
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
			circles = append(circles, newCircle(glm.Vec2{300 + 0.1*float32(x), 20 + 0.1*float32(y)}, vel.Normalized(), 10, 230, mesh, cVec))
		}
	}
	return circles
}
func NimbleMoss() {
	var openGL = closedGL.InitClosedGL(800, 600, "nimble moss")
	openGL.LimitFPS(true)
	var lineMesh = openGL.CreateLineMesh()
	var circleMesh = openGL.CreateCircleMesh()

	var middleY = openGL.Window.Wh - 50
	var middleX = openGL.Window.Ww / 2
	var lines = []glm.Vec2{{0, 250}, {middleX, middleY}, {middleX, middleY}, {800, 250}}
	for i := 0; i < len(lines); i += 2 {
		lineMesh.AddLine(lines[i], lines[i+1], glm.Vec4{1, 1, 1, 1}, glm.Vec4{1, 1, 1, 1})
	}

	//	var circles = loadColourBall(&circleMesh)
	var circles = loadImage(&circleMesh)

	lineMesh.Copy()
	circleMesh.Copy()

	var rectMesh = openGL.CreateRectMesh()
	var speedSlider = closedGL.NewSlider(&rectMesh, &openGL, glm.Vec4{20, 550, 75, 20})
	var gravSlider = closedGL.NewSlider(&rectMesh, &openGL, glm.Vec4{620, 550, 75, 20})
	openGL.LimitFPS(false)
	for !openGL.WindowShouldClose() {
		speedSlider.Process()
		gravSlider.Process()
		var delta = openGL.GetDelta()
		circleMesh.Clear()

		openGL.ClearBG(glm.Vec4{0, 0})
		openGL.BeginDrawing()
		lineMesh.Draw()
		circleMesh.Draw()
		openGL.DrawFPS(0, 0, 1)
		for i := 0; i < len(circles); i++ {
			circles[i].process(delta)
			circles[i].draw()
			for j := 0; j < len(lines); j += 2 {
				var p1 = lines[j]
				var p2 = lines[j+1]
				var line = closedGL.CalculateLine(p1, p2)
				var closestPoint = line.ClosestPoint(circles[i].pos)
				var dist = closestPoint.Sub(&circles[i].pos)
				circles[i].speed = 260 * speedSlider.GetPercentage()
				circles[i].gravY = 0.62 * gravSlider.GetPercentage()
				if dist.Len() < circles[i].radius {
					circles[i].vel = getReflectionVec(p1, p2, circles[i].vel)
					circles[i].pos[1] -= 1
				}
			}
		}
		circleMesh.Copy()
		rectMesh.Draw()
		circleMesh.Draw()
		openGL.EndDrawing()

	}
	openGL.Free()
}

/* func StartTuwuing() {
	var openGL = closedGL.InitClosedGL(1400, 800, "demo")

	openGL.LimitFPS(true)
	var complete = tuwuing_complete.NewTuwuingComplete(&openGL)

	for !openGL.WindowShouldClose() {
		complete.Process()

		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.BeginDrawing()
		complete.Draw()

		openGL.DrawFPS(500, 0, 1)

		openGL.EndDrawing()
	}
	openGL.Free()
} */

/*
func StartClosedGL() {
	var openGL = closedGL.InitClosedGL(800, 600, "demo")
	openGL.LimitFPS(false)
	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)

	var val = true

	//	var chunks = []ynnebcraft.Chunk{}
	//	var mesher = ynnebcraft.NewGreedyMesher()

	/* 	for i := 0; i < 20; i++ {
		for j := 0; j < 10; j++ {
			chunks = append(chunks, ynnebcraft.NewChunk(glm.Vec3{float32(i) * 32, 0, float32(j) * 32}, glm.Vec3{32, 32, 32}, &openGL, &mesher))
		}
	}
	openGL.Logger.Enabled = true

	var pxTest = openGL.CreatePixelMesh()
	var lines = openGL.CreateLineMesh()

	for i := 0; i < 800; i++ {
		for j := 0; j < 600; j++ {
			pxTest.AddPixel(glm.Vec2{float32(i), float32(j)}, glm.Vec4{float32(j) / 600, float32(i) / 800, 1, 1})
		}
	}
	pxTest.Copy()
	lines.AddPath([]glm.Vec2{{0, 0}, {200, 200}}, []glm.Vec4{{1, 1, 1, 1}, glm.Vec4{1, 1, 1, 1}})
	lines.AddQuadraticBezier(glm.Vec2{0, 0}, glm.Vec2{200, 200}, glm.Vec2{0, 200}, glm.Vec4{0, 0.54, 0.57, 1})
	lines.Copy()

	var cam = openGL.NewCam2D()
	openGL.LimitFPS(false)
	for !openGL.WindowShouldClose() {
		cam.Process(float32(openGL.FPSCounter.Delta))
		pxTest.View = cam.ViewMat
		if openGL.IsKeyPressed(glfw.Key(glfw.KeyF)) {
			val = !val
			closedGL.SetWireFrameMode(val)
		}
		openGL.BeginDrawing()
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.DrawFPS(500, 0, 1)

		//openGL.DrawSprite(glm.Vec4{0, 0, 20, 20}, "./assets/sprites/fence_small.png", glm.Vec4{0, 0, 1, 1}, glm.Vec2{32, 1024}, 1)

		openGL.Text.DrawText(500, 50, "x:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[0]), 'f', 2, 64), 1)
		openGL.Text.DrawText(500, 75, "y:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[1]), 'f', 2, 64), 1)
		openGL.Text.DrawText(500, 100, "z:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[2]), 'f', 2, 64), 1)

		openGL.Text.DrawText(600, 50, "x:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[0]), 'f', 2, 64), 1)
		openGL.Text.DrawText(600, 75, "y:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[1]), 'f', 2, 64), 1)
		openGL.Text.DrawText(600, 100, "z:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[2]), 'f', 2, 64), 1)

		/* for i := 0; i < len(chunks); i++ {
			chunks[i].Draw()
		}
		//	lines.Draw()
		pxTest.Draw()
		//openGL.DrawRect(glm.Vec4{0, 0, 100, 100}, glm.Vec4{0, 1, 0, 1}, 1)
		openGL.EndDrawing()
	}
	openGL.Free()
}
*/

func assert(val bool) {
	if !val {
		panic("assert failed")
	}
}
