package main

import "github.com/suvaidkhan/envsecrets/internal/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		println("error occured")
	}
}
