package util

import log "github.com/sirupsen/logrus"

//func init() {
//	log.SetFormatter(&log.JSONFormatter{})
//}

func Info(delivery string, msg1 interface{}, msg ...interface{}) {
	log.WithFields(log.Fields{"delivery": delivery}).Info(msg1, msg)
}

func Pr(msg ...interface{}) {
	log.Println(msg)
}
