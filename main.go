package main

import "github.com/suvaidkhan/envsecret/internal/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		println("error occured")
	}
}
