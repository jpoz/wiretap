package disk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jpoz/wiretap"
)

const TimeFormat = "2006-01-02T15_04_05Z07_00"

type Storage struct {
	Location string
}

func (s Storage) Save(session wiretap.Session) {
	r := session.Request
	u := session.Request.URL
	t := session.Started.Format(TimeFormat)

	fullPath := filepath.Join(s.Location, u.Host, u.Path, r.Method, t)
	reqPath := filepath.Join(fullPath, "request.txt")
	respPath := filepath.Join(fullPath, "response.txt")

	os.MkdirAll(fullPath, 0777)

	writeFile(reqPath, session.RequestBody)
	writeFile(respPath, session.ResponseBody)
}

func writeFile(path string, bytes []byte) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("err %+v\n", err)
		return
	}

	f.Write(bytes)
	f.Close()
}
