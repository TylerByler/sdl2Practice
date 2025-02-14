package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// UTILITY
const winWidth, winHeight int = 800, 600

type gameState uint8

const (
	start gameState = iota
	play
	pause
)

var state = start

type color struct {
	r, g, b, a byte
}

var white color = color{255, 255, 255, 255}

type position struct {
	x, y float32
}

type ball struct {
	position
	radius float32
	xVel   float32
	yVel   float32
	color  color
}

type paddle struct {
	position
	width  float32
	height float32
	speed  float32
	score  int
	color  color
}

var scoreNums = [][]byte{
	{1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1},

	{1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1},

	{1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1},

	{1, 1, 1,
		0, 0, 1,
		0, 1, 1,
		0, 0, 1,
		1, 1, 1},

	{1, 0, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		0, 0, 1},

	{1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1},

	{1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1},

	{1, 1, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1},

	{1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1},

	{1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1},
}

var titleLetters = []byte{
	1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 0, 0, 1, 1, 1,
	1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0,
	1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0,
	0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0,
	1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0,
}

// BALL FUNCITONS
func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle, elapsedTime float32) {
	ball.x += ball.xVel * elapsedTime
	ball.y += ball.yVel * elapsedTime

	if (ball.y - float32(ball.radius)) < 0 {
		ball.yVel = -ball.yVel
		ball.y = ball.radius
	}

	if (ball.y + float32(ball.radius)) > float32(winHeight) {
		ball.yVel = -ball.yVel
		ball.y = float32(winHeight) - ball.radius
	}

	if ball.x-ball.radius < 0 {
		ball.position = getScreenCenter()
		leftPaddle.y = float32(winHeight) / 2
		rightPaddle.y = float32(winHeight) / 2
		rightPaddle.score++
	}

	if ball.x+ball.radius > float32(winWidth) {
		ball.position = getScreenCenter()
		leftPaddle.y = float32(winHeight) / 2
		rightPaddle.y = float32(winHeight) / 2
		leftPaddle.score++
	}

	if rightPaddle.score == 9 || leftPaddle.score == 9 {
		resetGame(leftPaddle, rightPaddle, ball)
	}

	if ball.x-ball.radius < leftPaddle.x+leftPaddle.width/2 &&
		ball.y < leftPaddle.y+leftPaddle.height/2 &&
		ball.y > leftPaddle.y-leftPaddle.height/2 {
		ball.xVel = -ball.xVel
		ball.x = leftPaddle.x + leftPaddle.width/2 + ball.radius
	}

	if ball.x+ball.radius > rightPaddle.x-rightPaddle.width/2 &&
		ball.y < rightPaddle.y+rightPaddle.height/2 &&
		ball.y > rightPaddle.y-rightPaddle.height/2 {
		ball.xVel = -ball.xVel
		ball.x = rightPaddle.x - rightPaddle.width/2 - ball.radius
	}

	ball.xVel *= 1.0002
	ball.yVel *= 1.0002
}

// PADDLE FUNCTIONS
func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.width/2)
	startY := int(paddle.y - paddle.height/2)

	for y := 0; y < int(paddle.height); y++ {
		for x := 0; x < int(paddle.width); x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}

	numX := lerp(paddle.x, getScreenCenter().x, .2)
	drawScore(position{numX, 60}, paddle.color, 15, paddle.score, pixels)
}

func (paddle *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= paddle.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += paddle.speed * elapsedTime
	}
}

func (paddle *paddle) aiUpdate(ball *ball) {
	paddle.y = ball.y
}

// UTILITY FUNCTIONS
func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels) && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func getScreenCenter() position {
	return position{float32(winWidth / 2), float32(winHeight / 2)}
}

func drawScore(pos position, color color, size int, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i, v := range scoreNums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

func drawTitle(size int, pixels []byte) {
	startX := (winWidth - (19 * size)) / 2
	startY := (winHeight - (5 * size)) / 2

	for i, v := range titleLetters {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, white, pixels)
				}
			}
		}
		startX += size
		if (i+1)%19 == 0 {
			startY += size
			startX -= size * 19
		}
	}

}

func lerp(a float32, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

func resetGame(leftPaddle *paddle, rightPaddle *paddle, ball *ball) {
	leftPaddle.y = float32(winHeight) / 2
	rightPaddle.y = float32(winHeight) / 2
	leftPaddle.score = 0
	rightPaddle.score = 0
	ball.position = getScreenCenter()
	state = start
}

// -----------GAME START-------------
func main() {
	// ACTIVATE ALL SDL THINGS
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Tyler's Window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, 1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer texture.Destroy()

	pixels := make([]byte, winWidth*winHeight*4) // CREATE PIXEL BUFFER

	// CREATE ENTITIES
	player1 := paddle{position{50, float32(winHeight) / 2}, 20, 100, 300, 0, white}
	player2 := paddle{position{float32(winWidth - 50), float32(winHeight) / 2}, 20, 100, 300, 0, white}
	ball := ball{position{300, 300}, 20, 350, 350, white}

	keyState := sdl.GetKeyboardState()
	var frameStart time.Time
	var elapsedTime float32

	// --------- GAME LOOP -----------
	for {
		frameStart = time.Now()

		// POLL FOR EXIT BUTTON EVENT
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		// ENTER ENTITY PIXELS INTO PIXEL BUFFER
		clear(pixels) // CLEARS ENTIRE PIXEL BUFFER EVERY FRAME
		if state != start {
			player1.draw(pixels)
			player2.draw(pixels)
			ball.draw(pixels)
		}

		if state == play {
			// UPDATE ENTITIES
			player1.update(keyState, elapsedTime)
			player2.aiUpdate(&ball)
			ball.update(&player1, &player2, elapsedTime)
		} else if state == start {
			drawTitle(15, pixels)
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				player1.score = 0
				player2.score = 0
				state = play
			}
		}

		// SET TEXTURE TO PIXEL BUFFER
		texture.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		// UPDATE IS FRAMERATE INDEPENDENT AND CAPS FRAMERATE AT 200
		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < .005 {
			sdl.Delay(5 - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
		fmt.Println("(xVel,yVel):  (", ball.xVel, ",", ball.yVel, ")")
	}
}
