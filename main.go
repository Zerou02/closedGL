package main

import (
	_ "image/png"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
	turingfontparser "github.com/Zerou02/closedGL/turing_font_parser"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	StartTTF()
	//StartClosedGL()
	//StartTuwuing()
}

func StartTTF() {
	isOnLineTest()
	var opengl = closedGL.InitClosedGL(800, 800, "comic")

	opengl.LimitFPS(false)
	var p = turingfontparser.NewTuringFont("./assets/font/jetbrains_mono_medium.ttf", &opengl)
	var glyfs = []turingfontparser.Glyf{}
	/* var test = p.ParseGlyf(uint32('a'), 1)
	var factor = test.CalcScaleFactor(12) */
	/* 	factor = 1 */
	for _, x := range "S" {
		var glyf = p.ParseGlyf(uint32(x), 1)
		/* 		glyf.Scale(factor) */
		glyfs = append(glyfs, glyf)
	}

	var offset float32 = 0
	for i := 0; i < len(glyfs); i++ {
		/* 	for j := 0; j < len(glyfs[i].SimpleGlyfs); j++ {
			glyfs[i].SimpleGlyfs[j].AddOffset(glm.Vec2{offset, 100})
		} */
		offset += glyfs[i].AdvanceWidth
	}

	var lines = opengl.CreateLineMesh()
	var pixelMesh = opengl.CreatePixelMesh()
	var polys = []turingfontparser.Polygon2{}
	var tri = opengl.CreateTriMesh()
	for _, x := range glyfs {
		polys = append(polys, turingfontparser.NewPolygon2(x.SimpleGlyfs[0], &opengl, &tri, &lines, &pixelMesh))
	}

	pixelMesh.SetPixelSize(10)
	pixelMesh.Copy()

	lines.Copy()

	tri.Copy()

	tri.Copy()
	for !opengl.WindowShouldClose() {
		opengl.BeginDrawing()

		if opengl.IsKeyPressed(glfw.KeyL) {
			closedGL.PrintlnVec2(opengl.GetMousePos())
		}
		if opengl.IsKeyPressed(glfw.KeyI) {
			closedGL.PrintlnVec2(opengl.GetMousePos())
		}

		opengl.SetWireFrameMode(opengl.IsKeyDown(glfw.KeyF))
		opengl.ClearBG(glm.Vec4{0, 0, 0, 0})
		tri.Draw()
		/* 	opengl.SetWireFrameMode(true)
		tri.Draw() */
		opengl.SetWireFrameMode(opengl.IsKeyDown(glfw.KeyF))
		pixelMesh.Draw()
		lines.Draw()

		opengl.DrawFPS(600, 0, 1)
		opengl.EndDrawing()
	}
	opengl.Free()
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

func testLines() {
	var p1 = glm.Vec2{0, 0}
	var p2 = glm.Vec2{100, 100}
	var p3 = glm.Vec2{100, 0}
	var p4 = glm.Vec2{0, 100}
	//normal-normal
	var l1 = closedGL.CalculateLine(p1, p2)
	var l2 = closedGL.CalculateLine(p3, p4)
	var poin, _ = l1.GetIntersection(l2)
	var point2, _ = l2.GetIntersection(l1)
	assert(point2.Equal(&poin))
	println("expect 50,50")
	closedGL.PrintlnVec2(poin)

	//vert-horiz
	var l3 = closedGL.CalculateLine(p2, p3)
	var l4 = closedGL.CalculateLine(p1, p3)
	var cp2, _ = l3.GetIntersection(l4)
	var cp22, _ = l4.GetIntersection(l3)
	assert(cp2.Equal(&cp22))
	println("expect 0,100")
	closedGL.PrintlnVec2(cp2)

	//vert-normal
	var l5 = closedGL.CalculateLine(p1, p4)
	var l6 = closedGL.CalculateLine(p4, p3)
	var cp3, _ = l5.GetIntersection(l6)
	var cp32, _ = l6.GetIntersection(l5)

	assert(cp3.Equal(&cp32))
	println("expect 0,100")
	closedGL.PrintlnVec2(cp3)

	//vert-vert
	var l7 = closedGL.CalculateLine(p2, p3)
	var l8 = closedGL.CalculateLine(p1, p4)
	var cp4, succ = l7.GetIntersection(l8)
	var cp42, _ = l8.GetIntersection(l7)
	assert(cp4.Equal(&cp42))
	println("expect 0,0,false")
	print(succ, ",")
	closedGL.PrintlnVec2(cp4)

	//horiz,horiz
	var l9 = closedGL.CalculateLine(p2, p4)
	var l10 = closedGL.CalculateLine(p1, p3)
	var cp5, succ2 = l9.GetIntersection(l10)
	var cp52, _ = l10.GetIntersection(l9)
	assert(cp5.Equal(&cp52))
	println("expect 0,0,false")
	print(succ2, ",")
	closedGL.PrintlnVec2(cp5)

	//horiz,normal
	var l11 = closedGL.CalculateLine(p1, p2)
	var l12 = closedGL.CalculateLine(p1, p3)
	var cp6, succ3 = l11.GetIntersection(l12)
	var cp62, _ = l12.GetIntersection(l11)
	assert(cp62.Equal(&cp6))
	println("expect 0,0,true")
	print(succ3, ",")
	closedGL.PrintlnVec2(cp6)
}

func isOnLineTest() {
	var p1 = glm.Vec2{0, 0}     //ul
	var p2 = glm.Vec2{100, 100} //or
	var p3 = glm.Vec2{100, 0}   //ur
	var p4 = glm.Vec2{0, 100}   //ol

	var vertical = closedGL.CalculateLine(p1, p4)
	var horizontal = closedGL.CalculateLine(p2, p4)
	var normal = closedGL.CalculateLine(p3, p4)
	assert(vertical.IsOnLine(glm.Vec2{0, 50}))
	assert(!vertical.IsOnLine(glm.Vec2{0, 150}))
	assert(horizontal.IsOnLine(glm.Vec2{50, 100}))
	assert(!horizontal.IsOnLine(glm.Vec2{120, 100}))
	assert(normal.IsOnLine(glm.Vec2{50, 50}))
	assert(!normal.IsOnLine(glm.Vec2{60, 60}))

}
