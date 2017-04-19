package main

import (
	"log"
	"time"

	"github.com/jessevdk/go-flags"
)

type CommandOptions struct {
	Name string  `short:"n" long:"name" description:"name of the sensor" required:"true"`
	Freq int     `short:"f" long:"freq" description:"update frequency in cycles/sec" required:"true"`
	Max  float64 `long:"max" description:"maximum value for generated readings" required:"true"`
	Min  float64 `long:"min" description:"minimum value for generated readings" required:"true"`
	Step float64 `short:"s" long:"step" description:"maximum allowable change per measurement" required:"true"`
}

func main() {
	start := time.Now()

	options := CommandOptions{}

	parser := initParser(&options)

	_, err := parser.Parse()

	if err != nil {
		log.Print(err)
	}

	run(&options)

	log.Print(time.Since(start))
}

func init() {
	var banner = `
  |------------------------------------------------------------------------------|
  |                                                                              |
  |              SensorStreaming Golang - Create and Publish sensor data         |
  |                                                                              |
  |------------------------------------------------------------------------------|
  `

	log.Print(banner)
}

func initParser(options *CommandOptions) (parser *flags.Parser) {
	//default behaviour is HelpFlag | PrintErrors | PassDoubleDash - we need to override the stderr output
	return flags.NewParser(options, flags.HelpFlag)
}

func run(options *CommandOptions) {
	log.Print(options)
}
