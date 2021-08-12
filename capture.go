package main

import (
	"fmt"
	"image"
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/vova616/screenshot"
)

func main() {
	systray.Run(onMenuReady, func() {})
}

func onMenuReady() {
	systray.SetTitle("Screen Capture")
	mRecord := systray.AddMenuItem("Capture", "Starts the screen capturer");
	systray.AddSeparator();
	mHelp := systray.AddMenuItem("Help", "Open the help menu");
	mFullscreen := systray.AddMenuItem("Enter Fullscreen", "Enter fullscreen");
	mQuit := systray.AddMenuItem("Quit", "Quit now");

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
			case <-mRecord.ClickedCh:
				fmt.Println("Starting capture");
				showRecorder();
			case <-mHelp.ClickedCh:
				fmt.Println("Show help menu");
			case <-mFullscreen.ClickedCh:
				fmt.Println("Use Fullscreen");
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-sigc:
				systray.Quit()
		}
	}
}

func showRecorder() {
	screenWidth, screenHeight := robotgo.GetScreenSize()
	wnd, cv, err := sdlcanvas.CreateWindow(screenWidth, screenHeight, "Screen Capture")
	if err != nil {
		panic(err)
	}

	defer wnd.Destroy()

	cursor, err := cv.LoadImage("pointer.png")
	if err != nil {
		panic(err)
	}

	wnd.MainLoop(func() {
		img, err := screenshot.CaptureRect(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: screenWidth, Y: screenHeight}})
		if err != nil {
			panic(err)
		}

		cv.DrawImage(img)
		x, y := robotgo.GetMousePos()
		cv.DrawImage(cursor, float64(x), float64(y))
	})
}