package closedGL

/*
import (
	"strconv"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.3-core/gl"
)

type TileSheet struct {
	tex     *Texture
	dim     glm.Vec2
	tileDim glm.Vec2
	name    string
}
type TileMap struct {
	SheetIdxMap []TileSheet
	entries     []int //mapIdx,idxInMap
	dim         glm.Vec2
	tileDim     glm.Vec2
	vao, vbo    Buffer
	shader      *Shader
	projection  *glm.Mat4
	vertexLen   int
	tex         Texture
	scale       float32
}

func NewTileMap(shader *Shader, projection *glm.Mat4) TileMap {

	var retMap = TileMap{SheetIdxMap: []TileSheet{
		{
			tex:     LoadImage("./assets/tilemaps/dungeon_a1.png", gl.RGBA),
			dim:     glm.Vec2{512, 512},
			name:    "City_TileB",
			tileDim: glm.Vec2{16, 16},
		},
		{
			tex:     LoadImage("./assets/tilemaps/dungeon_a2.png", gl.RGBA),
			dim:     glm.Vec2{512, 512},
			name:    "cityC",
			tileDim: glm.Vec2{16, 16},
		},
		{
			tex:     LoadImage("./assets/tilemaps/dungeon_a1.png", gl.RGBA),
			dim:     glm.Vec2{512, 512},
			name:    "cityC",
			tileDim: glm.Vec2{16, 16},
		},
	},
		dim:        glm.Vec2{512, 512},
		tileDim:    glm.Vec2{16, 16},
		entries:    make([]int, 32*32*2),
		shader:     shader,
		projection: projection,
		scale:      1}
	for i := 0; i < len(retMap.entries); i++ {
		retMap.entries[i] = -1
	}
	for i := 0; i < 32*32; i++ {
		retMap.AddEntry(i, i%3, 32*5)
	}
	retMap.generateVertices()
	return retMap
}

func (this *TileMap) AddEntry(at, tileSheet, tileIdx int) {
	this.entries[at*2] = tileSheet
	this.entries[at*2+1] = tileIdx
}

func (this *TileMap) generateVertices() {
	//posX,posY,texX,texY,texID
	var vertices = []float32{}

	//createVertices
	for i := 0; i < 2*int(this.dim[0]/this.tileDim[0])*int(this.dim[1]/this.tileDim[1]); i += 2 {
		var sheetIdx = this.entries[i]
		var sheet = this.SheetIdxMap[sheetIdx]
		var tileIdx = this.entries[i+1]
		var tileSizeX, tileSizeY = this.tileDim[0], this.tileDim[1]
		var tileX, tileY = IdxToGridPos(tileIdx, int(sheet.dim[0]/sheet.tileDim[0]), int(sheet.dim[1]/sheet.tileDim[1]))
		var currX, currY = IdxToGridPos(i/2, int(this.dim[0]/this.tileDim[0]), int(this.dim[1]/this.tileDim[1]))

		var upperLeftVertices = []float32{
			this.scale * float32(currX) * this.tileDim[0],
			this.scale * float32(currY) * this.tileDim[1],
			(float32(tileX) * tileSizeX) / sheet.dim[0],
			(float32(tileY) * tileSizeY) / sheet.dim[1],

			float32(sheetIdx),
		}
		var upperRightVertices = []float32{
			this.scale * float32(currX+1) * this.tileDim[0],
			this.scale * float32(currY) * this.tileDim[1],
			(float32(tileX+1) * tileSizeX) / sheet.dim[0],
			(float32(tileY) * tileSizeY) / sheet.dim[1],
			float32(sheetIdx),
		}
		var lowerRightVertices = []float32{
			this.scale * float32(currX+1) * this.tileDim[0],
			this.scale * float32(currY+1) * this.tileDim[1],
			(float32(tileX+1) * tileSizeX) / sheet.dim[0],
			(float32(tileY+1) * tileSizeY) / sheet.dim[1],
			float32(sheetIdx),
		}
		var lowerLeftVertices = []float32{
			this.scale * float32(currX) * this.tileDim[0],
			this.scale * float32(currY+1) * this.tileDim[1],
			(float32(tileX) * tileSizeX) / sheet.dim[0],
			(float32(tileY+1) * tileSizeY) / sheet.dim[1],
			float32(sheetIdx),
		}
		vertices = append(vertices, upperRightVertices...)
		vertices = append(vertices, lowerRightVertices...)
		vertices = append(vertices, upperLeftVertices...)
		vertices = append(vertices, lowerRightVertices...)
		vertices = append(vertices, lowerLeftVertices...)
		vertices = append(vertices, upperLeftVertices...)
	}
	this.vertexLen = len(vertices)
	PrintlnFloat(vertices[0+4*5])
	PrintlnFloat(vertices[1+4*5])
	PrintlnFloat(vertices[2+4*5])
	PrintlnFloat(vertices[3+4*5])
	PrintlnFloat(vertices[4+4*5])

	GenerateBuffers(&this.vao, &this.vbo, nil, vertices, 0, nil, []int{2, 2, 1})
}
func (this *TileMap) Draw() {
	this.shader.use()
	for i, x := range this.SheetIdxMap {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, *x.tex)
		this.shader.setUniform1i("tex"+strconv.FormatInt(int64(i), 10), int32(i))
	}

	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)

	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.VERTEX_ARRAY, this.vbo)
	this.shader.setUniformMatrix4("projection", this.projection)
	var model = glm.Ident4()
	this.shader.setUniformMatrix4("model", &model)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(this.vertexLen))
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
}
*/
