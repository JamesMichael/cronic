package main

import (
	"flag"

	"github.com/jamesmichael/cronic"
)

func main() {
	path := flag.String("path", cronic.Getenv("CRONIC_PATH", "/etc/cron.d"), "Path to crontabs")
	flag.Parse()

	options := cronic.NewOptions()
	options.SetPath(*path)

	app := cronic.New(options)
	app.Run()
}
