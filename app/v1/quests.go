package v1

	import (
		"database/sql"
		"github.com/JustinOsu/api/common"
	)

type singleQuest struct {
	ID			int		`json:"id,omitempty"`
	Name		string	`json:"name"`
	Description	string	`json:"description"`
	Conditions	string	`json:"conditions"`
	Amount		int		`json:"amount"`
	Reward		string	`json:"reward"`
	Time		int		`json:"time"`
	TimeLeft	int		`json:"timeleft"`
	Icon		string	`json:"icon"`
	Done		int		`json:"done"`
}

type singleQuest2 struct {
	ID			int		`json:"id,omitempty"`
	Name		string	`json:"name"`
	Description	string	`json:"description"`
	Conditions	string	`json:"conditions"`
	Amount		int		`json:"amount"`
	Reward		string	`json:"reward"`
	Time		int		`json:"time"`
	TimeLeft	int		`json:"timeleft"`
	Icon		string	`json:"icon"`
	Done		int		`json:"done"`
}

type multiQuestsData struct {
	common.ResponseBase
	Quests []singleQuest `json:"quests"`
}

func QuestsGet(md common.MethodData) common.CodeMessager {
	var (
		r		multiQuestsData
		rows 	*sql.Rows
		err		error
	)
	if md.Query("uid") != "" {
		rows, err = md.DB.Query("SELECT quests.id, quests.name, quests.description, quests.conditions, quests.amount, quests.icon, quests.time, quests.timeleft, quests.reward, user_quests.done FROM quests, user_quests WHERE user_quests.user_id = ? AND user_quests.quest_id = quests.id", md.Query("uid"))
	} else {
		rows, err = md.DB.Query("SELECT quests.id, quests.name, quests.description, quests.conditions, quests.amount, quests.icon, quests.time, quests.timeleft, quests.reward, discord_roles.done FROM quests, discord_roles")
	}
	if err != nil {
		md.Err(err)
		return Err500
	}
	defer rows.Close()
	for rows.Next() {
		nc := singleQuest{}
		err = rows.Scan(&nc.ID, &nc.Name, &nc.Description, &nc.Conditions, &nc.Amount, &nc.Icon, &nc.Time, &nc.TimeLeft, &nc.Reward, &nc.Done)
		if err != nil {
			md.Err(err)
		}
		r.Quests = append(r.Quests, nc)
	}
	if err := rows.Err(); err != nil {
		md.Err(err)
	}
	r.ResponseBase.Code = 200
	return r
}

type UserQuestsData struct {
	User		int `json:"user_id"`
	Quest		int `json:"quest_id"`
	Amount		int `json:"amount"`
	Time		int `json:"time"`
}

type UserQuests struct {
	common.ResponseBase
	Quests []UserQuestsData `json:"quests"`
}

func UserQuestsGet(md common.MethodData) common.CodeMessager {
	ui := md.Query("uid")

	if ui == "0" {
		return ErrMissingField("uid")
	}

	var (
		r		UserQuests
		rows	*sql.Rows
		err 	error
	)

	rows, err = md.DB.Query("SELECT user_id, quest_id, amount, time FROM user_quests WHERE user_id = ?", ui)

	if err != nil {
		md.Err(err)
		return Err500
	}

	defer rows.Close()
	for rows.Next() {
		nc := UserQuestsData{}
		err = rows.Scan(&nc.User, &nc.Quest, &nc.Amount, &nc.Time)
		if err != nil {
			md.Err(err)
		}
		r.Quests = append(r.Quests, nc)
	}
	if err := rows.Err(); err != nil {
		md.Err(err)
	}
	r.ResponseBase.Code = 200
	return r
}