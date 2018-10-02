package main

import (
	"log"

	"github.com/defectus/glutton"
)

const Glutton = "Glutton"

func main() {
	printVersionInfo()
	err := glutton.Run()
	if err != nil {
		log.Printf("error running %s: %+v", Glutton, err)
		return
	}
	log.Printf("Hope you had fun running %s", Glutton)
}

var (
	// VERSION of application
	VERSION string
	// COMMIT of application
	COMMIT string
	// BRANCH of application
	BRANCH string
	// TAG of application (closest)
	TAG string
	// BUILDTIME of application
	BUILDTIME string
	// AUTHOR of application
	AUTHOR string
)

func printVersionInfo() {
	log.Printf("%s version info:", Glutton)
	log.Printf("author:%s\tversion:%s\ttag:%s\tfinger print:%s\tbranch:%s", AUTHOR, VERSION, TAG, COMMIT, BRANCH)
	log.Printf("build time:%s", BUILDTIME)
}
