package empire

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// DefaultTag is used when json decoding an image. If there is no tag part
// present, this will be used as the tag.
const DefaultTag = "latest"

// Image represents a container image, which is tied to a repository.
type Image struct {
	ID   string
	Repo string
}

func (i Image) String() string {
	return encodeImage(i)
}

// Scan implements the sql.Scanner interface.
func (i *Image) Scan(src interface{}) error {
	if src, ok := src.([]byte); ok {
		*i = decodeImage(string(src))
	}

	return nil
}

// Value implements the driver.Value interface.
func (i Image) Value() (driver.Value, error) {
	return driver.Value(i.String()), nil
}

func (i *Image) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*i = decodeImage(s)
	return nil
}

func encodeImage(i Image) string {
	return fmt.Sprintf("%s:%s", i.Repo, i.ID)
}

func decodeImage(s string) (image Image) {
	p := strings.Split(s, ":")

	if len(p) == 0 {
		return
	}

	image.Repo = p[0]

	if len(p) == 1 {
		image.ID = DefaultTag
		return
	}

	image.ID = p[1]
	return
}
