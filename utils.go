package main

import "github.com/golang/glog"

//custom log level
const (
	Info = iota
	Warning
	Debug
	Error
)

//CheckErr checks given error
func checkErr(messge string, err error, level int) {
	if err != nil {
		switch level {
		case Info, Warning, Debug:
			glog.Infoln(messge, err)
		case Error:
			glog.Fatalln(messge, err)
		}
	}
}

//convert type map[string]string to array
func convertMap2Array(m map[string]string) (s []string) {
	for _, v := range m {
		s = append(s, v)
	}

	return s
}

//convert type array to map[string]string
func convertArray2Map(s []string) (m map[string]string) {
	m = make(map[string]string)
	for _, v := range s {
		m[v] = v
	}
	return m
}
