package v1

import (
	"database/sql"
	"github.com/JustinOsu/api/common"
	"time"
)

type singleShop struct {
	ID			int		`json:"id,omitempty"`
	Name		string	`json:"name"`
	Description	string	`json:"description"`
	Icon		string	`json:"icon"`
	Price		int		`json:"price"`
}

type multiShopData struct {
	common.ResponseBase
	Shop []singleShop `json:shop`
}

func BuyItemGet(md common.MethodData) common.CodeMessager {
	var (
		Exists	bool
	)

	err := md.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM shop WHERE `id` = ?)", md.Query("id")).Scan(&Exists)
	if err != nil && err != sql.ErrNoRows {
		md.Err(err)
		return Err500
	}
	if Exists {
		_, err := md.DB.Exec("UPDATE users u, users_stats s, shop p SET s.coins = s.coins - p.price AND u.privieges = 7 AND u.donor_expire = ? + 2628000 WHERE u.id = ? AND u.priviles != 2", time.Now().Unix(), md.ID())
	}
}

func ShopGet(md common.MethodData) common.CodeMessager {
	var (
		r		multiShopData
		rows	*sql.Rows
		err		error
	)
	if md.Query("id") != "" {
		rows, err = md.DB.Query("SELECT id, name, description, icon, price FROM shop WHERE id = ? LIMIT 1", md.Query("id"))		
	} else {
		rows, err = md.DB.Query("SELECT id, name, description, icon, price FROM shop")
	}
	if err != nil {
		md.Err(err)
		return Err500
	}
	defer rows.Close()
	for rows.Next() {
		nc := singleShop{}
		err = rows.Scan(&nc.ID, &nc.Name, &nc.Description, &nc.Icon, &nc.Price)
		if err != nil {
			md.Err(err)
		}
		r.Shop = append(r.Shop, nc)
	}
	if err := rows.Err(); err != nil {
		md.Err(err)
	}
	r.ResponseBase.Code = 200
	return r
}