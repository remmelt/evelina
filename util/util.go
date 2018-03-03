package util

import "log"

func L(delivery string, msg1 interface{}, msg ...interface{}) {
	log.Println(delivery, msg1, msg)
}

func Pr(msg ...interface{}) {
	log.Println(msg)
}
