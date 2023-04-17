package combo

import (
	"fmt"
	"genytest/utilities"
)

type Combo struct {
	Items   map[string]string `yaml:"items"`
	Hash    string            `yaml:"hash"`
	Useless bool
}

func (c Combo) String() (output string) {
	for k, v := range c.Items {
		output += fmt.Sprintf("%s:%s,", k, v)
	}
	return
}

func (c *Combo) GenHash() {
	c.Hash = utilities.Hash(c.String())
}

func (c Combo) Clone() *Combo {
	clone := Combo{
		Items: map[string]string{},
	}
	for key, value := range c.Items {
		clone.Items[key] = value
	}

	clone.Hash = c.Hash
	return &clone
}
