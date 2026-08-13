package models

import (
	"bytes"
	"fmt"
)

type Folder struct {
	Collections []Collection
}

func (f Folder) String() string {
	var out bytes.Buffer
	for _, c := range f.Collections {
		fmt.Fprintf(&out, "%s", c)
		out.WriteString("\n")
	}

	return out.String()
}
