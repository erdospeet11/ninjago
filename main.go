package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 400
	ScreenHeight = 350
)

func main() {
	isMulti := false
	for _, arg := range os.Args[1:] {
		if strings.ToLower(arg) == "--multi" {
			isMulti = true
			break
		}
	}

	if isMulti {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			cmd := exec.Command("go", "run", ".", "--instance", "1")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Printf("Instance 1 error: %v", err)
			}
		}()

		time.Sleep(100 * time.Millisecond)

		go func() {
			defer wg.Done()
			cmd := exec.Command("go", "run", ".", "--instance", "2")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Printf("Instance 2 error: %v", err)
			}
		}()

		wg.Wait()
	} else {
		instanceTitle := ""
		for i, arg := range os.Args[1:] {
			if strings.ToLower(arg) == "--instance" && i+1 < len(os.Args[1:]) {
				instanceNum := os.Args[i+2]
				instanceTitle = " - Player " + instanceNum
				break
			}
		}

		ebiten.SetWindowSize(ScreenWidth*2, ScreenHeight*2)
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		ebiten.SetScreenClearedEveryFrame(false)
		ebiten.SetVsyncEnabled(false)
		ebiten.SetWindowTitle("Ninja Go Bomberman" + instanceTitle)
		if err := ebiten.RunGame(&Game{}); err != nil {
			log.Fatal(err)
		}
	}
}
