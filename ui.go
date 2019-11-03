package cronic

import (
	"os/exec"

	"github.com/rivo/tview"
)

func NewUI(crontabs *crontabCollection) *tview.Application {
	app := tview.NewApplication()

	files := tview.NewList()
	files.ShowSecondaryText(false)
	files.SetBorder(true)
	files.SetTitle(" Files ")
	files.SetHighlightFullLine(true)
	files.SetDoneFunc(func() {
		app.Stop()
	})
	environment := tview.NewList()
	environment.ShowSecondaryText(false)
	environment.SetBorder(true)
	environment.SetTitle(" Environment ")
	environment.SetSelectedFocusOnly(true)
	environment.SetHighlightFullLine(true)

	commands := tview.NewList()
	commands.ShowSecondaryText(false)
	commands.SetBorder(true)
	commands.SetTitle(" Commands ")
	commands.SetSelectedFocusOnly(true)
	commands.SetHighlightFullLine(true)
	commands.SetDoneFunc(func() {
		app.SetFocus(files)
	})

	files.SetChangedFunc(func(i int, file, t string, s rune) {
		environment.Clear()
		commands.Clear()

		crontab := crontabs.named(file)
		if err := crontab.Error(); err != nil {
			environment.AddItem(err.Error(), "", 0, nil)
			return
		}

		for _, e := range crontab.Environment() {
			environment.AddItem(e, "", 0, nil)
		}

		for _, c := range crontab.Commands() {
			command := c.command
			user := c.user

			commands.AddItem(command, "", 0, func() {
				file, err := crontab.WriteCommand(command)
				if err != nil {
					return
				}

				cmd := exec.Command("sudo", "-u", user, file)

				if err := cmd.Start(); err != nil {
					panic(err)
				}

				app.Stop()
			})
		}
	})

	for _, path := range crontabs.Paths() {
		files.AddItem(path, "", 0, func() {
			app.SetFocus(commands)
		})
	}

	flex := tview.NewFlex().
		AddItem(files, 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(environment, 0, 1, false).
			AddItem(commands, 0, 2, true), 0, 3, false)

	app.SetRoot(flex, true)

	return app
}
