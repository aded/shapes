package shapes

import (
	"math"

	"github.com/remogatto/mathgl"
	gl "github.com/remogatto/opengles2"
	"github.com/remogatto/shaders"
)

type Segment struct {
	shape

	// Points of the segment
	x1, y1, x2, y2 float32
	vertices       [4]float32
}

func NewSegment(x1, y1, x2, y2 float32) *Segment {

	segment := new(Segment)

	// Set the default color

	segment.Color(DefaultColor)

	// Set the geometry

	segment.x1, segment.x2 = x1, x2
	segment.y1, segment.y2 = y1, y2

	segment.vertices = [4]float32{
		segment.x1, segment.y1,
		segment.x2, segment.y2,
	}

	// Size of the segment bounding box

	segment.width = float32(math.Abs(float64(x1 - x2)))
	segment.height = float32(math.Abs(float64(y1 - y2)))

	// Center of the segment
	segment.x = (segment.x1 + segment.x2) / 2
	segment.y = (segment.y1 + segment.y2) / 2

	// Shader sources

	vShaderSrc := (shaders.VertexShader)(
		`precision mediump float;
                 attribute vec4 pos;
                 attribute vec4 color;
                 varying vec4 vColor;
                 uniform mat4 model;
                 uniform mat4 projection;
                 uniform mat4 view;
                 void main() {
                     gl_Position = projection*model*view*pos;
                     vColor = color;
                 }`)
	fShaderSrc := (shaders.FragmentShader)(
		`precision mediump float;
                 varying vec4 vColor;
                 void main() {
                     gl_FragColor = vColor;
                 }`)

	// Link the program
	segment.program = shaders.NewProgram(vShaderSrc.Compile(), fShaderSrc.Compile())
	segment.program.Use()

	// Get variables IDs from shaders
	segment.posId = segment.program.GetAttribute("pos")
	segment.colorId = segment.program.GetAttribute("color")
	segment.projMatrixId = segment.program.GetUniform("projection")
	segment.modelMatrixId = segment.program.GetUniform("model")
	segment.viewMatrixId = segment.program.GetUniform("view")

	// Fill the model matrix with the identity.
	segment.modelMatrix = mathgl.Ident4f()

	return segment
}

// Draw actually renders the object on the surface.
func (segment *Segment) Draw() {
	// Color is the same for each vertex
	vertexColor := [8]float32{
		segment.nColor[0], segment.nColor[1], segment.nColor[2], segment.nColor[3],
		segment.nColor[0], segment.nColor[1], segment.nColor[2], segment.nColor[3],
	}

	segment.program.Use()
	gl.VertexAttribPointer(segment.posId, 2, gl.FLOAT, false, 0, &segment.vertices[0])
	gl.EnableVertexAttribArray(segment.posId)

	gl.VertexAttribPointer(segment.colorId, 4, gl.FLOAT, false, 0, &vertexColor[0])
	gl.EnableVertexAttribArray(segment.colorId)

	gl.UniformMatrix4fv(int32(segment.modelMatrixId), 1, false, (*float32)(&segment.modelMatrix[0]))
	gl.UniformMatrix4fv(int32(segment.projMatrixId), 1, false, (*float32)(&segment.projMatrix[0]))
	gl.UniformMatrix4fv(int32(segment.viewMatrixId), 1, false, (*float32)(&segment.viewMatrix[0]))

	gl.DrawArrays(gl.LINES, 0, 2)

	gl.Flush()
	gl.Finish()
}
