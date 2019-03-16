package main

import "github.com/nattaponra/logz"

func main() {
	log := logz.NewLogz("abc/log.log", "%Y-%m-%d")
	log.Write([]byte("dfd\r\n"))
}
