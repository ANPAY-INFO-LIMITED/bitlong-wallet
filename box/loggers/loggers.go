package loggers

import (
	"log"
	"os"
)

var (
	cronLogger   *log.Logger
	lntLogger    *log.Logger
	tokenLogger  *log.Logger
	bdInfoLogger *log.Logger
	litLogger    *log.Logger
	boxLogger    *log.Logger
	frpLogger    *log.Logger
	chanLogger   *log.Logger
)

func Cron() *log.Logger {
	return cronLogger
}

func SetCron(f *os.File) {
	cronLogger = log.New(f, "[CRON]", log.Ldate|log.Ltime|log.Lshortfile)
}

func Lnt() *log.Logger {
	return lntLogger
}

func SetLnt(f *os.File) {
	lntLogger = log.New(f, "[LNT]", log.Ldate|log.Ltime|log.Lshortfile)
}

func Token() *log.Logger {
	return tokenLogger
}

func SetToken(f *os.File) {
	tokenLogger = log.New(f, "[Tokn]", log.Ldate|log.Ltime|log.Lshortfile)
}

func BdInfo() *log.Logger {
	return bdInfoLogger
}

func SetBdInfo(f *os.File) {
	bdInfoLogger = log.New(f, "[BDIF]", log.Ldate|log.Ltime|log.Lshortfile)
}

func Lit() *log.Logger {
	return litLogger
}

func SetLit(f *os.File) {
	litLogger = log.New(f, "[LIT]", log.Ldate|log.Ltime|log.Lshortfile)
}

func Box() *log.Logger {
	return boxLogger
}

func SetBox(f *os.File) {
	boxLogger = log.New(f, "[BOX]", log.Ldate|log.Ltime|log.Lshortfile)
}

func Frp() *log.Logger {
	return frpLogger
}

func SetFrp(f *os.File) {
	frpLogger = log.New(f, "[FRP]", log.Ldate|log.Ltime|log.Lshortfile)
}

func Chan() *log.Logger {
	return chanLogger
}

func SetChan(f *os.File) {
	chanLogger = log.New(f, "[CHAN]", log.Ldate|log.Ltime|log.Lshortfile)
}
