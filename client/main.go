package main

import (
	_ "github.com/TheTitanrain/w32"
	"main/frontend"
	"os"
)

func main() {
	frontend.Parse(os.Args)
}
