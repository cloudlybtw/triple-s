package main

import (
	"flag"
	"fmt"
	"os"
	"triple-s/launchServer"
)

var (
	ver  bool
	port int
	dir  string
)

func init() {
	flag.BoolVar(&ver, "ver", true, "for logging process")
	flag.IntVar(&port, "port", 8080, "Port number")
	flag.StringVar(&dir, "dir", "data", "Path to the directory")
}

func Usage() {
	fmt.Println(`$ ./triple-s --help  
Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`)
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	triples, err := launchServer.CoreServer(ver, port, dir)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	if err := triples.Launch(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
