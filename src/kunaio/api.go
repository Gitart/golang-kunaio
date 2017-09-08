// Copyright 2017 Aleksey Morarash <tuxofil@gmail.com>
//
// Licensed under the BSD 2 Clause License (the "License");
// you may not use the file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://opensource.org/licenses/BSD-2-Clause
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kunaio

import (
	"fmt"
	"time"
)

const (
	// Kuna.io base URL
	gBaseURL = "https://kuna.io"
	// Supported market type names
	BTCUAH = "btcuah"
	ETHUAH = "ethuah"
)

var (
	// List of supported market types
	supportedMarkets = []string{
		BTCUAH,
		ETHUAH,
	}
)

type Stats struct {
	// Server time
	Time time.Time
	// Current ICO buy price
	Buy float64
	// Current ICO sell price
	Sell float64
	// Lowest deal price for last 24 hours
	Low float64
	// Highest deal price for last 24 hours
	High float64
	// Last deal price
	Last float64
	// Trade volume in base currency for last 24 hours
	Vol float64
	// Total trade price for last 24 hours
	Amount float64
}

type OrderBook struct {
	Asks Orders
	Bids Orders
}

type Order struct {
	// Order ID
	ID int
	// For asks always "sell", for bids always "buy"
	Side string
	// Order type: "limit" or "market"
	OrdType string
	// Price for one ICO
	Price float64
	// Average order price
	AvgPrice float64
	// Order state: wait
	State string
	// Market identifier
	Market string
	// Order creation time
	CreatedAt time.Time
	// Order volume, in ICO
	Volume float64
	// ICO amount left
	RemainingVolume float64
	// ICO amount sold (for asks) or bought (for bids)
	ExecutedVolume float64
	// Deals count for this order
	TradesCount int
}

type Orders []Order

type HistoryEntry struct {
	// Order ID
	ID int
	// Price for one ICO
	Price float64
	// Volume of ICO
	Volume float64
	// Volume of UAH
	Funds float64
	// Market identifier
	Market string
	// Deal time
	CreatedAt time.Time
}

type History []HistoryEntry

type UserInfo struct {
	// Kuna.io user email
	Email string
	// If user account activated or not
	Activated bool
	// List of user assets
	Accounts []Account
}

type Account struct {
	// Currency type
	Currency string
	// Asset balance
	Balance float64
	// Locked funds
	Locked float64
}

type Trade struct {
	// Order ID
	ID int
	// Price per one ICO
	Price float64
	// ICO amount
	Volume float64
	// UAH amount
	Funds float64
	// Market type
	Market string
	// Deal date and time
	CreatedAt time.Time
	// "bid" or "ask"
	Side string
}

// Return list of supported markets.
func SupportedMarkets() []string {
	return supportedMarkets
}

// Return server time.
func GetServerTime() (time.Time, error) {
	j, err := doGet(fmt.Sprintf(
		"%s/api/v2/timestamp", gBaseURL))
	if err != nil {
		return time.Time{}, err
	}
	return jsonGetTime(j)
}

// Return latest market stats.
func GetLatestStats(market string) (s Stats, err error) {
	j, err := gClient.Get(fmt.Sprintf(
		"%s/api/v2/tickers/%s", gBaseURL, market))
	if err != nil {
		return s, err
	}
	return decodeLatestStats(j)
}

// Return order book (lists of current asks and bids).
func GetOrderBook(market string) (OrderBook, error) {
	j, err := doGet(fmt.Sprintf("%s/api/v2/order_book?market=%s",
		gBaseURL, market))
	if err != nil {
		return OrderBook{}, err
	}
	return decodeOrderBook(j)
}

// Return trade history.
func GetTradeHistory(market string) (History, error) {
	j, err := doGet(fmt.Sprintf("%s/api/v2/trades?market=%s",
		gBaseURL, market))
	if err != nil {
		return nil, err
	}
	return decodeHistory(j)
}

// Return user info and his assets.
func GetUserInfo(access_key, secret_key string) (*UserInfo, error) {
	url := privURL("GET", "/api/v2/members/me",
		access_key, secret_key, nil)
	j, err := doGet(url)
	if err != nil {
		return nil, err
	}
	return decodeUserInfo(j)
}

// Return list of active user orders.
func GetUserOrders(access_key, secret_key, market string) ([]Order, error) {
	url := privURL("GET", "/api/v2/orders",
		access_key, secret_key, Args{{"market", market}})
	j, err := doGet(url)
	if err != nil {
		return nil, err
	}
	return decodeOrders(j)
}

// Return list of user deals.
func GetUserTrades(access_key, secret_key, market string) (Trades, error) {
	url := privURL("GET", "/api/v2/trades/my",
		access_key, secret_key, Args{{"market", market}})
	j, err := doGet(url)
	if err != nil {
		return nil, err
	}
	return decodeUserTrades(j)
}

// Create new order.
func NewOrder(access_key, secret_key, market, side string, volume, price float64) (Order, error) {
	url := privURL("POST", "/api/v2/orders",
		access_key, secret_key,
		Args{
			{"market", market},
			{"price", fmt.Sprintf("%f", price)},
			{"side", side},
			{"volume", fmt.Sprintf("%f", volume)},
		})
	j, err := doPost(url)
	if err != nil {
		return Order{}, err
	}
	return decodeOrder(j)
}

// Cancel user order, identified by order ID.
func CancelOrder(access_key, secret_key string, id int) (Order, error) {
	url := privURL("POST", "/api/v2/order/delete",
		access_key, secret_key,
		Args{{"id", fmt.Sprintf("%d", id)}})
	j, err := doPost(url)
	if err != nil {
		return Order{}, err
	}
	return decodeOrder(j)
}
