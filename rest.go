package main

import (
	"context"
	"fmt"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"log"
	"sync"
	"time"
)

func TradeLoop(config ApplicationConfig, command chan string, response chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := InitClient(&config.Token)
	accounts, err := GetAccounts(*client)
	if err != nil {
		log.Fatalln(err)
	}
	iisAccountID := GetAccountByType(accounts, sdk.AccountTinkoffIIS)

	for c := range command {
		log.Printf("Got command: %s\n", c)
		if c == "getPositions" {
			positions, err := GetPositionsPortfolio(*client, *iisAccountID)
			if err != nil {
				log.Println("GetPositionsPortfolio error: ", err)
				continue
			}
			response <- FormatPositions(positions)
		}
	}
}

func FormatPositions(positions []sdk.PositionBalance) string {
	var response string
	for _, pos := range positions {
		response = response + "<b>" + pos.Name + "<b>\n" +
			fmt.Sprintf("%.2f", pos.AveragePositionPrice.Value) + " " +
			string(pos.AveragePositionPrice.Currency) + "\n" +
			fmt.Sprintf("%.2f", pos.Balance) + " " + string(pos.InstrumentType) + "\n" +
			"Стоимость: " + fmt.Sprintf("%.2f", pos.Balance*pos.AveragePositionPrice.Value) + "\n\n"
	}
	return response
}

func InitClient(token *string) *sdk.RestClient {
	return sdk.NewRestClient(*token)
}

func GetAccounts(client sdk.RestClient) ([]sdk.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("Получение всех брокерских счетов")
	accounts, err := client.Accounts(ctx)
	if err != nil {
		log.Println(errorHandle(err))
		return nil, err
	}
	log.Printf("%+v\n", accounts)
	return accounts, nil
}

func GetAccountByType(accounts []sdk.Account, accountType sdk.AccountType) *string {
	for _, acc := range accounts {
		if acc.Type == accountType {
			return &acc.ID
		}
	}
	return nil
}

func GetPositionsPortfolio(client sdk.RestClient, accountId string) ([]sdk.PositionBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	positions, err := client.PositionsPortfolio(ctx, accountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return positions, nil
}

func rest(token *string) {
	client := sdk.NewRestClient(*token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Получение всех брокерских счетов")
	accounts, err := client.Accounts(ctx)
	if err != nil {
		log.Fatalln(errorHandle(err))
	}
	log.Printf("%+v\n", accounts)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	iisAccountID := GetAccountByType(accounts, sdk.AccountTinkoffIIS)
	/*
		log.Println("Получение списка операций для счета по-умолчанию за последнюю неделю по инструменту(FIGI) BBG000BJSBJ0")
		// Получение списка операций за период по конкретному инструменту(FIGI)
		// Например: ниже запрашиваются операции за последнюю неделю по инструменту NEE
		operations, err := client.Operations(ctx, iisAccount.ID, time.Now().AddDate(0, 0, -7), time.Now(), "BBG000BJSBJ0")
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%+v\n", operations)

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	*/
	log.Printf("Получение списка НЕ валютных активов портфеля для счета %s\n", *iisAccountID)
	positions, err := client.PositionsPortfolio(ctx, *iisAccountID)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Printf("%+v\n", positions)
	ListPositions(positions)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Получение списка валютных активов портфеля для счета %s\n", *iisAccountID)
	positionCurrencies, err := client.CurrenciesPortfolio(ctx, *iisAccountID)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", positionCurrencies)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Получение списка валютных и НЕ валютных активов портфеля для счета %s\n", *iisAccountID)
	// Метод является совмещеним PositionsPortfolio и CurrenciesPortfolio
	portfolio, err := client.Portfolio(ctx, *iisAccountID)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", portfolio)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Получение списка выставленных заявок(ордеров) для счета %s\n", *iisAccountID)
	orders, err := client.Orders(ctx, *iisAccountID)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", orders)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

}
