package log

import (
	"fmt"

	"github.com/jpoz/wiretap"
)

const TimeFormat = "2006-01-02T15_04_05Z07_00"

type Storage struct {
	Log Logger
}

type Logger interface {
	Println(v ...interface{})
}

func (s Storage) Save(session wiretap.Session) {
	u := session.Request.URL

	str := fmt.Sprintf("%s", u)

	s.Log.Println(str)
}
