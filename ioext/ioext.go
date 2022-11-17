package ioext

import (
	"io"
	"log"
)

func Close(closer io.Closer) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			log.Printf("Error on close %T: %s\n", closer, err.Error())
		}
	}
}
