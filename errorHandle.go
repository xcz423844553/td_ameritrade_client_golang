package main

import "log"

func handleFatalErr(tag string, err error) {
	if err != nil {
		log.Fatalln("[" + tag + "]" + err.Error())
	}
}
