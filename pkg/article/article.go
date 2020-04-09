package article

import (
	"encoding/json"
	"fmt"
)

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"Desc"`
	Content string `json:"Content"`
}

type Articles []Article

func (a *Article) MarshalBinary() ([]byte, error) {
	if a == nil {
		return []byte(``), fmt.Errorf("blank entry at MarshalBinary")
	}
	return json.Marshal(a)
}

func (a *Article) UnmarshalBinary(d []byte) error {
	if a == nil {
		return fmt.Errorf("nil article at UnmarshalBinary")
	}
	return json.Unmarshal(d, a)
}
