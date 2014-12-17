package disk

import (
	"fmt"

	"github.com/jpoz/wiretap"
)

type Storage struct {
}

func (s Storage) Save(session wiretap.Session) {
	fmt.Printf("WRITE")
}
