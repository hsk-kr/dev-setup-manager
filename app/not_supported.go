package app

import "fmt"

func NotSupported(command string) {
	fmt.Printf("%s command is not supported yet\n", command)
}
