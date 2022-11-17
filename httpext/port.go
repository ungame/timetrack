package httpext

import "fmt"

type Port int

func (p Port) Addr() string {
	return fmt.Sprintf(":%d", p)
}
