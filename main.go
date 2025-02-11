package main

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

type color struct {
	r, g, b, a byte
}

type position struct {
	x, y float32
}

type ball struct {
	position
	radius int
	xVel   float32
	yVel   float32
	color  color
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x)+x, int(ball.y)+y, ball.color, pixels)
			}
		}
	}
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle) {
	ball.x += ball.xVel
	ball.y += ball.yVel

	// handle collisions
	if (ball.y-float32(ball.radius)) < 0 || (ball.y+float32(ball.radius)) > float32(winHeight) {
		ball.yVel = -ball.yVel
	}

	if ball.x < 0 || ball.x > float32(winWidth) {
		ball.x = 300
		ball.y = 300
	}

	if ball.x-float32(ball.radius) < (leftPaddle.x)+float32(leftPaddle.width/2) && ball.y < leftPaddle.y+float32(leftPaddle.height/2) && ball.y > leftPaddle.y-float32(leftPaddle.height/2) {
		ball.xVel = -ball.xVel
	}

	if ball.x+float32(ball.radius) > (rightPaddle.x)-float32(rightPaddle.width/2) && ball.y < rightPaddle.y+float32(rightPaddle.height)/2 && ball.y > rightPaddle.y-float32(rightPaddle.height)/2 {
		ball.xVel = -ball.xVel
	}
}

type paddle struct {
	position
	width  int
	height int
	color  color
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.width/2
	startY := int(paddle.y) - paddle.height/2

	for y := 0; y < paddle.height; y++ {
		for x := 0; x < paddle.width; x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y = paddle.y - 3
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y = paddle.y + 3
	}
}

func (paddle *paddle) aiUpdate(ball *ball) {
	paddle.y = ball.y
}

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

func main() {
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

	pixels := make([]byte, winWidth*winHeight*4)

	/* for y := 0; y < winHeight; y++ {
		for x := 0; x < winWidth; x++ {
			setPixel(x, y, color{byte(x % 255), byte(y % 255), 0, 0}, pixels)
		}
	} */

	player1 := paddle{position{50, 100}, 20, 100, color{255, 255, 255, 0}}
	player2 := paddle{position{float32(winWidth) - 50, 100}, 20, 100, color{255, 255, 255, 0}}
	ball := ball{position{300, 300}, 20, 3, 3, color{255, 255, 255, 0}}

	keyState := sdl.GetKeyboardState()

	for {
		// Close window with X button
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		clear(pixels)

		player1.update(keyState)
		player2.aiUpdate(&ball)
		ball.update(&player1, &player2)

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		texture.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}

}
