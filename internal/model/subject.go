package model

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"noeru/egarim/internal/util"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type Subject struct {
	gorm.Model
	Userid string
	Points int
	Secret string
}

var Size int = 10

func FirstOrCreateProfile(user discordgo.User) Subject {
	var subject Subject
	util.DB.FirstOrCreate(&subject, Subject{Userid: user.ID})

	// if Secret not set, set it
	if len(subject.Secret) == 0 {
		util.DB.Model(&subject).Update("Secret", GenerateSecret(Size))
	}

	return subject
}

/*
Generate a new Secret for the Member

TODO: prefix suffix?
*/
func GenerateSecret(size int) string {
	rb := make([]byte, size)

	_, err := rand.Read(rb)

	if err != nil {
		log.Fatalf("Unable to generate Secret!")
	}

	rs := base64.URLEncoding.EncodeToString(rb)

	return rs
}

func FindSubjectByKey(key string) (*Subject, *error) {
	var target Subject
	tx := util.DB.First(&target, Subject{Secret: key})

	if tx.Error != nil {
		return nil, &gorm.ErrRecordNotFound
	}

	return &target, nil
}