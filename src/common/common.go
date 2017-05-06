package common

import (
	"log"
)

type Config struct {
	Serve struct {
		Host string
		Port int
	}
	Auth struct {
		Username string
		Password string
	}
	Mods []Mod
}

type Mod struct {
	Prefix  string
	Comment string
	Token   string
}

func HandleError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
