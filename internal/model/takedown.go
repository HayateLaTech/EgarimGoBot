package model

import (
	"noeru/egarim/internal/util"

	"gorm.io/gorm"
)

type Takedown struct {
	gorm.Model
	PreyId string
	SuccId string
	NotifMsgId string
}

func RegNewKill(prey Subject, succ Subject) Takedown {
	var takedown Takedown = Takedown{PreyId: prey.Userid, SuccId: succ.Userid}

	util.DB.Create(&takedown)

	return takedown
}