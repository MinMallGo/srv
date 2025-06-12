package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
	"hash"
	"srv/userSrv/global"
	"srv/userSrv/model"
	"strconv"
	"time"
)

var (
	DefaultSaltLen   = 10
	DefaultIter      = 100
	DefaultKeyLen    = 32
	DefaultCryptFunc = sha512.New
)

func main() {
	//migrate()
	//pwdEncrypt("gen password")
	genUser()
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
		saltLen: DefaultSaltLen,
		iter:    DefaultIter,
		keyLen:  DefaultKeyLen,
		h:       DefaultCryptFunc,
	}

	// $algorithm$salt$password
	salt, encodePsw := gen(pwd, opts)
	fmt.Printf("%#v\n", verify(pwd, salt, encodePsw, opts))
	encodePsw = fmt.Sprintf("$pbkdf2-sha256$%s$%s", salt, encodePsw)
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

func genUser() {
	var user []model.User
	birthday := time.Date(2000, time.July, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 10; i++ {
		user = append(user, model.User{
			Mobile:   "1762324000" + strconv.Itoa(i),
			Password: pwdEncrypt("123456"),
			NickName: "batch_add_" + strconv.Itoa(i),
			Birthday: &birthday,
			Gender:   "female",
			Role:     1,
		})
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		for _, user := range user {
			err := tx.Create(&user).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}
