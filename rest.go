package main

import (
	"context"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"log"
	"time"
)

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

	iisAccountID := GetAccount(accounts, sdk.AccountTinkoffIIS)
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

func GetAccount(accounts []sdk.Account, accountType sdk.AccountType) *string {
	for _, acc := range accounts {
		if acc.Type == accountType {
			return &acc.ID
		}
	}
	return nil
}
