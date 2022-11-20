package dbContext

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	config "TgMoneyTrackerBot/configuration"

	_ "github.com/glebarez/go-sqlite"
)

const (
	sqlDriver                    = "sqlite"
	TableName                    = "TgMoneyBotDB"
	ChatIDColumnName             = "ChatID"
	MoneySpendingLimitColumnName = "MoneySpendingLimit"
	CurrentWastesColumnName      = "CurrentWastes"
	MoneyIncomeColumnName        = "MoneyIncome"
	UpdateDBCommand              = "UPDATE \"%s\" SET %s = %s WHERE chatID = %s"
	CreateUsersTableCommand      = "CREATE TABLE \"%s\" (\"%s\" INTEGER, \"%s\" INTEGER, \"%s\" INTEGER, \"%s\" INTEGER)"
	InsertUserIntoDBCommand      = "INSERT INTO %s VALUES (%s, %s, %s, %s)"
)

var botDB *sql.DB

func StartDB() *sql.DB {
	dbPath := path.Join(config.GetInfo().DBFolderPath, config.GetInfo().DBFileName)
	botDB = CheckDB(dbPath)
	return botDB
}

func CheckDB(path string) *sql.DB {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		db, err := sql.Open(sqlDriver, path)
		if err != nil {
			log.Fatal(err)
		}

		return db
	}

	dbFolder := config.GetInfo().DBFolderPath
	if _, err := os.Stat(dbFolder); os.IsNotExist(err) {
		err := os.Mkdir(dbFolder, 0755)
		if err != nil {
			log.Panic(err)
		}
	}

	os.Create(path)
	db, err := sql.Open(sqlDriver, path)
	if err != nil {
		log.Fatal(err)
	}

	err = CreateTable(db)
	if err != nil {
		log.Panic(err)
	}

	return db
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(
		CreateUsersTableCommand,
		TableName,
		ChatIDColumnName,
		MoneySpendingLimitColumnName,
		CurrentWastesColumnName,
		MoneyIncomeColumnName))
	if err != nil {
		return err
	}

	return nil
}

func UpdateValInDB(userID int64, parameter string, newValue string) {
	result, err := botDB.Exec(fmt.Sprintf(UpdateDBCommand, TableName, parameter, newValue, strconv.FormatInt(userID, 10)))
	if err != nil {
		log.Print("DB error:" + err.Error())
	}

	log.Print(result)
}

func GetBotDB() *sql.DB {
	return botDB
}
