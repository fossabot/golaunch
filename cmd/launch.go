package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	termbox "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"

	"github.com/micnncim/golaunch/util"
	"github.com/micnncim/picker"
)

func Launch(c *cli.Context) error {
	dataDir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		return err
	}
	dataDir += DataDir

	appDetails, err := util.ReadAppDataFiles(dataDir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	appDetailRows, err := appDetails.Render()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := termbox.Init(); err != nil {
		return errors.Wrap(err, "failed to initialize termbox")
	}

	x, y := 0, 1
	prompt := "QUERY>"
	rows := picker.MakeRows(appDetailRows, termbox.ColorDefault, termbox.ColorDefault)
	marker := termbox.ColorCyan
	s := picker.NewScreen(rows, x, y, "", prompt, termbox.ColorDefault, termbox.ColorDefault, termbox.ColorBlack, termbox.ColorWhite, marker)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer termbox.Close()

	for {
		s.Clear()
		s.DrawAll()

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				text := picker.GetRow(s.CursorY, s.BaseFg, s.BaseBg).Text
				termbox.Close()
				app := fmt.Sprintf("%s/%s.app", AppDir, strings.Split(text, " ")[0])
				if err := open.Run(app); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				os.Exit(0)
			case termbox.KeyEsc, termbox.KeyCtrlD, termbox.KeyCtrlC:
				termbox.Close()
				os.Exit(0)
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				s.DecreaseInput()
			case termbox.KeySpace:
				s.IncreaseInput(' ')
			case termbox.KeyArrowUp, termbox.KeyCtrlK, termbox.KeyCtrlP:
				s.MoveCursorUpper()
			case termbox.KeyArrowDown, termbox.KeyCtrlJ, termbox.KeyCtrlN:
				s.MoveCursorLower()
			default:
				if ev.Ch != 0 {
					s.IncreaseInput(ev.Ch)
				}
			}
		}

		if err := termbox.Flush(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	return nil
}
