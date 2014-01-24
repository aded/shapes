package shapes

import (
	"fmt"
	"image/color"

	"github.com/remogatto/mathgl"
	"github.com/remogatto/shaders"
)

var (
	// The default color for shapes is blue.
	DefaultColor = color.RGBA{0, 0, 0xff, 0xff}
)

type World interface {
	// Projection returns the projection matrix used to render
	// the objects in the World.
	Projection() mathgl.Mat4f

	// View returns the view matrix used to render the World from
	// the point-of-view of a camera.
	View() mathgl.Mat4f
}

type shape struct {
	x, y          float32
	width, height float32

	angle float32
	color color.Color

	// normalized RGBA color
	nColor [4]float32

	// Matrices
	projMatrix  mathgl.Mat4f
	modelMatrix mathgl.Mat4f
	viewMatrix  mathgl.Mat4f

	// GLSL program
	program shaders.Program

	// GLSL variables IDs
	colorId       uint32
	posId         uint32
	projMatrixId  uint32
	modelMatrixId uint32
	viewMatrixId  uint32
}

func (shape *shape) GetSize() (float32, float32) {
	return shape.width, shape.height
}

// Center returns the coordinates of the transformed center of the
// shape.
func (shape *shape) Center() (float32, float32) {
	return shape.x, shape.y
}

// Angle returns the current angle of the shape in degrees.
func (shape *shape) Angle() float32 {
	return shape.angle
}

// AttachToWorld fills projection and view matrices.
func (shape *shape) AttachToWorld(world World) {
	shape.projMatrix = world.Projection()
	shape.viewMatrix = world.View()
}

// Rotate the box around its center, by the given angle in degrees.
func (shape *shape) Rotate(angle float32) {
	shape.modelMatrix = mathgl.Translate3D(shape.x, shape.y, 0).Mul4(mathgl.HomogRotate3DZ(angle))
	shape.angle = angle
}

// Place the box at the given position
func (shape *shape) Position(x, y float32) {
	shape.modelMatrix = mathgl.Translate3D(x, y, 0)
	shape.x, shape.y = x, y
}

// Set the color of the shape.
func (shape *shape) Color(c color.Color) {

	shape.color = c

	// Convert to RGBA
	rgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	r, g, b, a := rgba.R, rgba.G, rgba.B, rgba.A

	// Normalize the color components
	shape.nColor = [4]float32{
		float32(r) / 255,
		float32(g) / 255,
		float32(b) / 255,
		float32(a) / 255,
	}
}

// Get the color of the shape.
func (shape *shape) GetColor() color.Color {
	return shape.color
}

// Get the color is a normalized float32 array.
func (shape *shape) GetNColor() [4]float32 {
	return shape.nColor
}

// String return a string representation of the shape in the form
// "(cx,cy),(w,h)".
func (shape *shape) String() string {
	return fmt.Sprintf("(%f,%f)-(%f,%f)", shape.x, shape.y, shape.width, shape.height)
}
