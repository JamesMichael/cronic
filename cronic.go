package cronic

import (
	"fmt"
	"github.com/rivo/tview"
	"log"
)

type application struct {
	app      *tview.Application
	crontabs *crontabCollection
}

func New() *application {
	paths, err := crontabPaths()
	if err != nil {
		log.Fatal(err.Error())
	}

	crontabs, err := loadCrontabs(paths)
	if err != nil {
		log.Fatal(err.Error())
	}

	ui := NewUI(crontabs)

	return &application{
		ui,
		crontabs,
	}
}

func (a application) Run() {
	err := a.app.Run()
	if err != nil {
		fmt.Printf("Error running application: %s\n", err)
	}
}
