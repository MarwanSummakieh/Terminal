package main

import (
  "bufio"
  "io"
  "os"
  "os/exec"
  "time"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/app"
  "fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/widget"
  "github.com/creack/pty"
)

const MaxBufferSize = 16

func main() {
  a := app.New()
  w := a.NewWindow("germ")

  ui := widget.NewTextGrid()

  os.Setenv("TERM", "dumb")
  c := exec.Command("/bin/bash")
  p, err := pty.Start(c)

  if err != nil {
    fyne.LogError("Failed to open pty", err)
    os.Exit(1)
  }

  defer c.Process.Kill()

  onTypedKey := func(e *fyne.KeyEvent) {
    if e.Name == fyne.KeyEnter || e.Name == fyne.KeyReturn {
      _, _ = p.Write([]byte{'
'})
    }
  }

  onTypedRune := func(r rune) {
    _, _ = p.WriteString(string(r))
  }

  w.Canvas().SetOnTypedKey(onTypedKey)
  w.Canvas().SetOnTypedRune(onTypedRune)

  buffer := [][]rune{}
  reader := bufio.NewReader(p)

  go func() {
    line := []rune{}
    buffer = append(buffer, line)
    for {
      r, _, err := reader.ReadRune()
      if err != nil {
        if err == io.EOF {
          return
        }
        os.Exit(0)
      }
      line = append(line, r)
      buffer[len(buffer)-1] = line
      if r == '
' {
        if len(buffer) > MaxBufferSize {
          buffer = buffer[1:]
        }
        line = []rune{}
        buffer = append(buffer, line)
      }
    }
  }()

  go func() {
    for {
      time.Sleep(100 * time.Millisecond)
      ui.SetText("")
      var lines string
      for _, line := range buffer {
        lines += string(line)
      }
      ui.SetText(string(lines))
    }
  }()

  w.SetContent(
    fyne.NewContainerWithLayout(
      layout.NewGridWrapLayout(fyne.NewSize(900, 325)),
      ui,
    ),
  )
  w.ShowAndRun()
}