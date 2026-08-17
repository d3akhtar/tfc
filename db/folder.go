package db

import (
	"bytes"
	"fmt"
)

type Folder struct {
	Name        string       `json:"name"`
	Collections []Collection `json:"collections"`
}

func (f Folder) String() string {
	var out bytes.Buffer
	for _, c := range f.Collections {
		fmt.Fprintf(&out, "%s", c)
		out.WriteString("\n")
	}

	return out.String()
}
