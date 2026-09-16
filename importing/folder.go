package importing

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/d3akhtar/tfc/domain"
)

func ImportFolder(path string) (*domain.Folder, error) {
	if !strings.HasSuffix(path, ".tfcf") {
		return nil, fmt.Errorf("Path should end with .tfcf for imported folders. Path: %s", path)
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error while reading file while importing folder %v", err)
	}

	res := domain.Folder{}
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshalling imported folder")
	}

	return &res, nil
}
