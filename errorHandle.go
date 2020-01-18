package td_ameritrade_client_golang

import "log"

func handleFatalErr(tag string, err error) {
	if err != nil {
		log.Fatalln("[" + tag + "]" + err.Error())
	}
}
