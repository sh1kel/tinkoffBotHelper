package main

import (
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"log"
)

func ListPositions(positins []sdk.PositionBalance) {
	for i, pos := range positins {
		log.Printf("%d. %s:\n\t%.2f %s\n\t%.0f %s\n\t%.2f\n", i,
			pos.Name,
			pos.AveragePositionPrice.Value,
			pos.AveragePositionPrice.Currency,
			pos.Balance,
			pos.InstrumentType,
			pos.AveragePositionPrice.Value*pos.Balance,
		)
	}
}
