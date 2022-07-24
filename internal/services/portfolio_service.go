package services

import (
	"dividend_portfolio_tracker_cli/internal"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const portfolio_data_file_path = ".dividend_portfolio_tracker"

func PurchaseStock(ticker string, shares int, price float64, date time.Time) {
	newStock := internal.StockData{
		Shares:       shares,
		Price:        price,
		PurchaseDate: date.Unix(),
		Operation:    internal.PurchaseOperation,
	}

	insertStockOperation(ticker, newStock)
}

func SellStock(ticker string, shares int, price float64, date time.Time) {
	newStock := internal.StockData{
		Shares:       shares,
		Price:        price,
		PurchaseDate: date.Unix(),
		Operation:    internal.SellOperation,
	}

	insertStockOperation(ticker, newStock)
}

func GetPortfolioData() internal.PortfolioData {
	path, err := getPortfolioDataFilePath()
	if err != nil {
		return internal.PortfolioData{
			Stocks: make(map[string][]internal.StockData),
		}
	}

	portfolio_data_content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(path); err != nil {
			fmt.Printf("Error occured while creating data.json file: %v", err)
		}

		return internal.PortfolioData{
			Stocks: make(map[string][]internal.StockData),
		}
	} else {
		var portfolio_data internal.PortfolioData
		if err := json.Unmarshal(portfolio_data_content, &portfolio_data); err != nil {
			fmt.Printf("Unable to unmarshal user data file: %v", err)
		}

		return portfolio_data
	}
}

func insertStockOperation(ticker string, newStock internal.StockData) {
	portfolio_data := GetPortfolioData()

	if stocks, ok := portfolio_data.Stocks[ticker]; ok {
		portfolio_data.Stocks[ticker] = append(stocks, newStock)
	} else {
		portfolio_data.Stocks[ticker] = []internal.StockData{newStock}
	}

	err := savePortfolioData(portfolio_data)
	if err != nil {
		fmt.Printf("Error occured while saving portfolio data: %v", err)
	}
}

func savePortfolioData(portfolioData internal.PortfolioData) error {
	path, err := getPortfolioDataFilePath()
	if err != nil {
		return err
	}

	portfolio_data_json, err := json.Marshal(portfolioData)
	if err != nil {
		return err
	}

	return os.WriteFile(path, portfolio_data_json, os.ModeAppend)
}

func getPortfolioDataFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + string(os.PathSeparator) + portfolio_data_file_path, nil
}