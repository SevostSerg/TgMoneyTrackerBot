package messageReceiver

import (
	botTypes "TgMoneyTrackerBot/botTypes"
	botDB "TgMoneyTrackerBot/dbContext"
	"fmt"
	"log"
	"strconv"

	tgAPI "github.com/go-telegram-bot-api/telegram-bot-api"
)

var tgBot *tgAPI.BotAPI
var updatesChan tgAPI.UpdatesChannel
var activeUsers map[int64]*botTypes.User

func Start(bot *tgAPI.BotAPI, updates tgAPI.UpdatesChannel) {
	tgBot = bot
	updatesChan = updates
	activeUsers = botTypes.GetUserListFromDB()
	fmt.Print(activeUsers)
	StartReceivingMessages()
}

func StartReceivingMessages() {
	for update := range updatesChan {
		if update.Message == nil { // ignore any non-Message Updates and /start
			continue
		}

		RecognizeUserMessage(&update)
	}
}

func RecognizeUserMessage(update *tgAPI.Update) {
	CheckUser(update)
	message, err := botTypes.IdentifyMessageType(update)
	if err != nil {
		tgBot.Send(tgAPI.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}

	ModifyCash(message, update.Message.Chat.ID)
}

func CheckUser(update *tgAPI.Update) {
	if _, contains := activeUsers[update.Message.Chat.ID]; !contains {
		newUser, err := botTypes.InsertUserIntoDB(update.Message.Chat.ID)
		if err != nil {
			log.Panic(err)
		}

		activeUsers[update.Message.Chat.ID] = newUser
	}
}

func ModifyCash(message *botTypes.Message, userID int64) {
	user, contatins := activeUsers[userID]
	if !contatins {
		log.Printf("User %s doesn't exist!", strconv.FormatInt(userID, 10))
		return
	}

	value, err := strconv.Atoi(message.Message)
	if err != nil {
		tgBot.Send(tgAPI.NewMessage(userID, "Please enter an integer value!"))
		return
	}

	log.Print(message.MessageType)
	switch message.MessageType {
	case botTypes.AddIncome:
		user.MoneyIncome += value
		botDB.UpdateValInDB(userID, botDB.MoneyIncomeColumnName, strconv.Itoa(user.MoneyIncome))
	case botTypes.AddWastes:
		user.CurrentWastes += value
		botDB.UpdateValInDB(userID, botDB.CurrentWastesColumnName, strconv.Itoa(user.CurrentWastes))
	default:
		log.Panic("Unknown command!")
		return
	}

	tgBot.Send(tgAPI.NewMessage(userID, fmt.Sprintf(
		"Cash:     %s\nIncome: %s\nWastes:  %s",
		strconv.Itoa(user.MoneyIncome-user.CurrentWastes),
		strconv.Itoa(user.MoneyIncome),
		strconv.Itoa(user.CurrentWastes))))
}
