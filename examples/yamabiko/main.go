package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	msg := Yamabiko(os.Stdin)
	fmt.Print(msg)
}

// TODO: io.Reader を引数にする
func Yamabiko(stdin io.Reader) string {
	in, err := io.ReadAll(stdin)
	if err != nil {
		log.Print(err)
	}
	return fmt.Sprintf("%v%v", string(in), string(in))
}
