package main

import (
	"flag"
	"fmt"
	"imgprocessor/process"
	"os"

	"github.com/golang/glog"
)

var from string
var to string
var label string
var limit int
var help bool

func init() {
	flag.BoolVar(&help, "help", false, "the preprocessing command usage")
	flag.StringVar(&from, "from", "./3", "the path of images need to process")
	flag.StringVar(&to, "to", "./4", "the path of  processed images")
	flag.StringVar(&label, "label", "./labels.txt", "the path of the label")
	flag.IntVar(&limit, "limit", 100, "the min number of images")
}

func main() {

	flag.Parse()
	handleHelp()
	start()
}

func handleHelp() {
	if help {
		fmt.Println("--from			the path of images need to process.\n\t\t\tdefault value is ./3")
		fmt.Println("--to			the path of processed images.\n\t\t\tdefault value is ./4")
		fmt.Println("--label			the path of the label.\n\t\t\tdefault value is ./labels.txt")
		fmt.Println("--limit			the min number of images to filter.\n\t\t\tdefault value is 100")
		os.Exit(0)
	}
}

func start() {
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(error); ok {
				glog.Errorln(err)
				os.Exit(-1)
			} else {
				glog.Info(err)
			}
		}
	}()

	p := process.New(from, to, limit)
	d := p.ParseLabel(label)
	p.Process(d)
	process.Detail(d)
	process.Summary(d)

}
