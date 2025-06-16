package test

import (
	"log"
	"srv/userSrv/global"
	"testing"
)

func TestGetPort(t *testing.T) {
	log.Println(global.GetPort())
}

func TestUUID(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Println(global.UUID())
	}
}
