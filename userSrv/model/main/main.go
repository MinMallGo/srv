package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"hash"
	"srv/userSrv/global"
	"srv/userSrv/model"
)

func main() {
	migrate()
	pwdEncrypt("gen password")
}

func migrate() {
	db := global.DB

	err := db.AutoMigrate(&model.User{})
	if err != nil {
		panic(fmt.Sprintf("autoMigrate failed: %v", err))
	}
}

type passwordEncrypt struct {
	saltLen int
	iter    int
	keyLen  int
	h       func() hash.Hash
}

// pbkdf2
func pwdEncrypt(pwd string) string {
	opts := &passwordEncrypt{
		saltLen: 10,
		iter:    100,
		keyLen:  32,
		h:       sha512.New,
	}

	// $algorithm$salt$password
	salt, encodePsw := gen(pwd, opts)
	fmt.Printf("%#v\n", verify(pwd, salt, encodePsw, opts))
	encodePsw = fmt.Sprintf("$pbkdf2-sha256$%s$%s", salt, encodePsw)
	fmt.Println(len(encodePsw))
	return encodePsw
}

func gen(str string, opts *passwordEncrypt) (string, string) {
	salt := genSalt(opts.saltLen)
	encodePsw := pbkdf2.Key([]byte(str), []byte(salt), opts.iter, opts.keyLen, opts.h)
	return salt, hex.EncodeToString(encodePsw)
}

func verify(password string, salt, encodePwd string, opts *passwordEncrypt) bool {
	tmp := pbkdf2.Key([]byte(password), []byte(salt), opts.iter, opts.keyLen, opts.h)
	return hex.EncodeToString(tmp) == encodePwd
}

func genSalt(strlen int) string {
	b := make([]byte, strlen)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
