package testlib

import (
	"fmt"
	"image/color"

	"github.com/aded/shapes"
	"github.com/remogatto/imagetest"
	"github.com/remogatto/mandala/test/src/testlib"
	gl "github.com/remogatto/opengles2"
)

const (
	distanceThreshold = 0.002
	texFilename       = "gopher.png"
	texDistThreshold  = 0.004
)

func distanceError(distance float64, filename string) string {
	return fmt.Sprintf("Image differs by distance %f, result saved in %s", distance, filename)
}

func (t *TestSuite) TestShape() {
	box := shapes.NewBox(t.renderState.boxProgram, 10, 20)

	// Color

	c := box.Color()
	nc := box.NColor()
	t.Equal(color.RGBA{0, 0, 255, 255}, c)
	t.Equal([4]float32{0, 0, 1, 1}, nc)

	box.SetColor(color.RGBA{0xaa, 0xaa, 0xaa, 0xff})
	c = box.Color()
	nc = box.NColor()
	t.Equal(color.RGBA{170, 170, 170, 255}, c)
	t.Equal([4]float32{0.6666667, 0.6666667, 0.6666667, 1}, nc)

	// GetSize
	size := box.Bounds().Size()
	w, h := size.X, size.Y
	t.True(w == 10)
	t.True(h == 20)

	// Center

	x, y := box.Center()
	t.Equal(float32(0), x)
	t.Equal(float32(0), y)

	// Center after translation

	winW, winH := t.renderState.window.GetSize()
	world := newWorld(winW, winH)
	box.AttachToWorld(world)
	box.MoveTo(10, 20)

	x, y = box.Center()
	t.Equal(float32(10), x)
	t.Equal(float32(20), y)

	// Angle after rotation

	box.Rotate(10)
	angle := box.Angle()
	t.Equal(float32(10), angle)

	// String representation

	t.Equal("(5,10)-(15,30)", box.String())
}

func (t *TestSuite) TestBox() {
	filename := "expected_box.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Create a box
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.MoveTo(float32(w/2), 0)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestRotatedBox() {
	filename := "expected_box_rotated_20.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Create a 100x100 pixel² box
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		box.AttachToWorld(world)
		// Place the box at the center of the screen
		box.MoveTo(float32(w/2), 0)
		// Rotate the box 20 degrees
		box.Rotate(20.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestTranslatedBox() {
	filename := "expected_box_translated_10_10.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Place a box on the center of the window
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		box.AttachToWorld(world)
		box.MoveTo(111, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestColoredBox() {
	filename := "expected_box_yellow.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		// Color is yellow
		box.SetColor(color.RGBA{255, 255, 0, 255})
		box.AttachToWorld(world)
		box.MoveTo(float32(w/2), 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestScaledBox() {
	filename := "expected_box_scaled.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		// Color is yellow
		box.SetColor(color.RGBA{0, 0, 255, 255})
		box.AttachToWorld(world)
		box.MoveTo(float32(w/2), 0)
		box.Scale(1.5, 1.5)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestSegment() {
	filename := "expected_line.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		segment := shapes.NewSegment(t.renderState.segmentProgram, 81.5, -40, 238.5, 44)

		// Color is yellow
		segment.SetColor(color.RGBA{255, 0, 0, 255})
		segment.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		segment.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < 0.0009, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestSegmentCenter() {
	segment := shapes.NewSegment(t.renderState.segmentProgram, 10, 15, 20, 20)

	x, y := segment.Center()
	t.Equal(float32(15), x)
	t.Equal(float32(17.5), y)

	size := segment.Bounds().Size()
	w, h := size.X, size.Y
	t.Equal(10, w)
	t.Equal(5, h)
}

func (t *TestSuite) TestTexturedBox() {
	filename := "expected_box_textured.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)

		// Create a box
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.MoveTo(float32(w/2), 0)

		// Add an image as a texture
		gopherTexture := world.addImageAsTexture(texFilename)

		texCoords := []float32{
			0, 0,
			1, 0,
			0, 1,
			1, 1,
		}

		box.SetTexture(gopherTexture, texCoords)

		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < texDistThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestTexturedRotatedBox() {
	filename := "expected_box_textured_rotated_20.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Create a box
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.MoveTo(float32(w/2), 0)

		// Add an image as a texture
		gopherTexture := world.addImageAsTexture(texFilename)
		texCoords := []float32{
			0, 0,
			1, 0,
			0, 1,
			1, 1,
		}
		box.SetTexture(gopherTexture, texCoords)

		box.Rotate(20.0)

		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < texDistThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestPartialTextureRotatedBox() {
	filename := "expected_box_partial_texture_rotated_20.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)

		// Create a box
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.MoveTo(float32(w/2), 0)

		// Add an image as a texture
		gopherTexture := world.addImageAsTexture(texFilename)

		texCoords := []float32{
			0, 0,
			0.5, 0,
			0, 0.5,
			0.5, 0.5,
		}
		box.SetTexture(gopherTexture, texCoords)

		box.Rotate(20.0)

		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < texDistThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestGroup() {
	filename := "expected_group.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)

		// Create first group, 2 small boxes
		group1 := shapes.NewGroup()
		b1 := shapes.NewBox(t.renderState.boxProgram, 20, 20)
		b1.MoveTo(30, 40)
		b2 := shapes.NewBox(t.renderState.boxProgram, 50, 50)
		b2.MoveTo(45, -25)
		b2.Rotate(20.0)
		group1.Append(b1)
		group1.Append(b2)

		// Create the main group
		group := shapes.NewGroup()
		group.Append(group1)
		group.Append(shapes.NewBox(t.renderState.boxProgram, 100, 100))

		// Get the second element of the group
		b3 := group.GetAt(1)
		b3.MoveTo(float32(w/2), 0)

		group.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		group.Draw()

		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestGroupTranslated() {
	filename := "expected_group_translated_20_15.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)

		// Create first group, 2 small boxes
		group1 := shapes.NewGroup()
		b1 := shapes.NewBox(t.renderState.boxProgram, 20, 20)
		b1.MoveTo(30, 40)
		b2 := shapes.NewBox(t.renderState.boxProgram, 50, 50)
		b2.MoveTo(45, -25)
		b2.Rotate(10.0)
		group1.Append(b1)
		group1.Append(b2)

		// Create the main group
		group := shapes.NewGroup()
		group.Append(group1)
		b3 := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		b3.MoveTo(float32(w/2), 0)
		group.Append(b3)

		// Translate by (20, 15)
		cx, cy := group.Center()
		group.MoveTo(cx+20, cy+15)

		group.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		group.Draw()

		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestGroupRotated() {
	filename := "expected_group_rotated_25.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)

		// Create first group, 2 small boxes
		group1 := shapes.NewGroup()
		b1 := shapes.NewBox(t.renderState.boxProgram, 20, 20)
		b1.MoveTo(30, 40)
		b2 := shapes.NewBox(t.renderState.boxProgram, 50, 50)
		b2.MoveTo(45, -25)
		b2.Rotate(10.0)
		group1.Append(b1)
		group1.Append(b2)

		// Create the main group
		group := shapes.NewGroup()
		group.Append(group1)
		b3 := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		b3.MoveTo(float32(w/2), 0)
		group.Append(b3)

		// Rotate by 25
		group.Rotate(25.0)

		group.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		group.Draw()

		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	} else {
		saveExpAct(t.outputPath, "success_"+filename, exp, act)
	}
}

// func getBufferDataFromImage(img image.Image) ([]byte, int, int) {
// 	bounds := img.Bounds()
// 	imgWidth, imgHeight := bounds.Size().X, bounds.Size().Y
// 	buffer := make([]byte, imgWidth*imgHeight*4)
// 	index := 0
// 	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
// 		for x := bounds.Min.X; x < bounds.Max.X; x++ {
// 			r, g, b, a := img.At(x, y).RGBA()
// 			buffer[index] = byte(r)
// 			buffer[index+1] = byte(g)
// 			buffer[index+2] = byte(b)
// 			buffer[index+3] = byte(a)
// 			index += 4
// 		}
// 	}

// 	return buffer, imgWidth, imgHeight
// }
