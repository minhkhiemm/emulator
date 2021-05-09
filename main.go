package main

import (
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kr/pty"
)

func main() {
	a := app.New()
	w := a.NewWindow("goerm")

	ui := widget.NewTextGrid()
	ui.SetText("I'm on a goerm, a fantastic terminal emulator")

	c := exec.Command("/bin/bash")
	p, err := pty.Start(c)
	if err != nil {
		fyne.LogError("Failed to open pty", err)
		os.Exit(1)
	}

	defer c.Process.Kill()

	onTypeKey := func(e *fyne.KeyEvent) {
		if e.Name == fyne.KeyEnter || e.Name == fyne.KeyReturn {
			_, _ = p.Write([]byte{'\r'})
		}
	}

	onTypeRune := func(r rune) {
		_, _ = p.WriteString(string(r))
	}

	w.Canvas().SetOnTypedKey(onTypeKey)
	w.Canvas().SetOnTypedRune(onTypeRune)

	if _, err := p.Write([]byte("ls\r")); err != nil {
		fyne.LogError("Failed to read pty", err)
		os.Exit(1)
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			b := make([]byte, 1024)
			_, err = p.Read(b)
			if err != nil {
				fyne.LogError("Failed to read pty", err)
				os.Exit(1)
			}

			ui.SetText(string(b))
		}
	}()

	w.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewGridWrapLayout(fyne.NewSize(420, 200)),
			ui,
		),
	)

	w.ShowAndRun()
}
