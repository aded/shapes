package shapes

import (
	gl "github.com/remogatto/opengles2"

	"github.com/remogatto/mathgl"
	"github.com/remogatto/shaders"
)

// A Box
type Box struct {
	shape

	// 4x4 color matrix (four color component for each vertex)
	vertexColor [16]float32

	// Vertices of the box
	vertices [8]float32
}

// NewBox creates a new box of given sizes.
func NewBox(width, height float32) *Box {

	box := new(Box)

	// The box is built around its center at (0, 0)
	box.vertices = [8]float32{
		-width / 2, -height / 2,
		width / 2, -height / 2,
		-width / 2, height / 2,
		width / 2, height / 2,
	}

	// Set the default color
	box.Color(DefaultColor)

	// Shader sources

	vShaderSrc := (shaders.VertexShader)(
		`precision mediump float;
                 attribute vec4 pos;
                 attribute vec4 color;
                 attribute vec2 texIn;
                 varying vec2 texOut;
                 varying vec4 vColor;
                 uniform mat4 model;
                 uniform mat4 projection;
                 uniform mat4 view;
                 void main() {
                     gl_Position = projection*model*view*pos;
                     vColor = color;
                     texOut = texIn;
                 }`)
	fShaderSrc := (shaders.FragmentShader)(
		`precision mediump float;
                 varying vec4 vColor;
	             varying vec2 texOut;
                 uniform sampler2D texture;
                 uniform float texRatio;
                 void main() {
                     vec2 flippedTexCoords = vec2(texOut.x, 1.0 - texOut.y);
                     vec4 texColor = texture2D(texture, flippedTexCoords) * texRatio;
                     vec4 vertColor = vColor * (1.0 - texRatio);
                     gl_FragColor = texColor + vertColor;
                 }`)

	// Link the program
	program := shaders.NewProgram(vShaderSrc.Compile(), fShaderSrc.Compile())
	box.program = program
	box.program.Use()

	// Get variables IDs from shaders
	box.posId = program.GetAttribute("pos")
	box.colorId = program.GetAttribute("color")
	box.projMatrixId = program.GetUniform("projection")
	box.modelMatrixId = program.GetUniform("model")
	box.viewMatrixId = program.GetUniform("view")
	box.texInId = program.GetAttribute("texIn")
	box.textureId = program.GetUniform("texture")
	box.texRatioId = program.GetUniform("texRatio")

	// Fill the model matrix with the identity.
	box.modelMatrix = mathgl.Ident4f()

	// Size of the box
	box.width = width
	box.height = height

	return box
}

// Draw actually renders the shape on the surface.
func (box *Box) Draw() {

	// Color is the same for each vertex
	vertexColor := [16]float32{
		box.nColor[0], box.nColor[1], box.nColor[2], box.nColor[3],
		box.nColor[0], box.nColor[1], box.nColor[2], box.nColor[3],
		box.nColor[0], box.nColor[1], box.nColor[2], box.nColor[3],
		box.nColor[0], box.nColor[1], box.nColor[2], box.nColor[3],
	}

	box.program.Use()

	gl.VertexAttribPointer(box.posId, 2, gl.FLOAT, false, 0, &box.vertices[0])
	gl.EnableVertexAttribArray(box.posId)

	gl.VertexAttribPointer(box.colorId, 4, gl.FLOAT, false, 0, &vertexColor[0])
	gl.EnableVertexAttribArray(box.colorId)

	gl.UniformMatrix4fv(int32(box.modelMatrixId), 1, false, (*float32)(&box.modelMatrix[0]))
	gl.UniformMatrix4fv(int32(box.projMatrixId), 1, false, (*float32)(&box.projMatrix[0]))
	gl.UniformMatrix4fv(int32(box.viewMatrixId), 1, false, (*float32)(&box.viewMatrix[0]))

	gl.Uniform1f(int32(box.texRatioId), 0.0)

	// Texture
	if len(box.texCoords) > 0 {
		gl.Uniform1f(int32(box.texRatioId), 1.0)
		gl.VertexAttribPointer(box.texInId, 2, gl.FLOAT, false, 0, &box.texCoords[0])
		gl.EnableVertexAttribArray(box.texInId)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, box.texBuffer)
		gl.Uniform1i(int32(box.textureId), 0)
	}

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	gl.Flush()
	gl.Finish()
}
