package botTypes

import (
	botDB "TgMoneyTrackerBot/dbContext"
	"fmt"
	"log"
	"strconv"
)

type User struct {
	ChatID             int64
	MoneySpendingLimit int
	MoneyIncome        int
	CurrentWastes      int
}

func GetUserListFromDB() map[int64]*User {
	userListChatIDKey := make(map[int64]*User)
	rows, err := botDB.GetBotDB().Query(fmt.Sprintf("select * from %s", botDB.TableName))
	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ChatID, &user.MoneySpendingLimit, &user.CurrentWastes, &user.MoneyIncome)
		if err != nil {
			log.Panic(err)
		}

		userListChatIDKey[user.ChatID] = &user
	}

	return userListChatIDKey
}

func InsertUserIntoDB(userID int64) (*User, error) {
	a := fmt.Sprintf(botDB.InsertUserIntoDBCommand, botDB.TableName, strconv.FormatInt(userID, 10), "0", "0", "0")
	log.Print(a)
	_, err := botDB.GetBotDB().Exec(a)
	return &User{
		ChatID:             userID,
		MoneySpendingLimit: 0,
		MoneyIncome:        0,
		CurrentWastes:      0,
	}, err
}
