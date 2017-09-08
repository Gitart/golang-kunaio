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
	"io"
)

// Convert decoded JSON object to Stats struct.
func decodeLatestStats(v interface{}) (s Stats, err error) {
	m, err := jsonGetMap(v)
	if err != nil {
		return s, err
	}
	timestamp, err := jsonGetTime(m["at"])
	if err != nil {
		return s, err
	}
	m, err = jsonGetMap(m["ticker"])
	if err != nil {
		return s, err
	}
	buy, err := jsonGetFloat(m["buy"])
	if err != nil {
		return s, err
	}
	sell, err := jsonGetFloat(m["sell"])
	if err != nil {
		return s, err
	}
	low, err := jsonGetFloat(m["low"])
	if err != nil {
		return s, err
	}
	high, err := jsonGetFloat(m["high"])
	if err != nil {
		return s, err
	}
	last, err := jsonGetFloat(m["last"])
	if err != nil {
		return s, err
	}
	vol, err := jsonGetFloat(m["vol"])
	if err != nil {
		return s, err
	}
	amount, err := jsonGetFloatDef(m["amount"], 0)
	if err != nil {
		return s, err
	}
	return Stats{
		Time:   timestamp,
		Buy:    buy,
		Sell:   sell,
		Low:    low,
		High:   high,
		Last:   last,
		Vol:    vol,
		Amount: amount,
	}, nil
}

// Convert decoded JSON object to Order Book.
func decodeOrderBook(v interface{}) (OrderBook, error) {
	m, err := jsonGetMap(v)
	if err != nil {
		return OrderBook{}, err
	}
	asks, err := decodeOrders(m["asks"])
	if err != nil {
		return OrderBook{}, err
	}
	bids, err := decodeOrders(m["bids"])
	if err != nil {
		return OrderBook{}, err
	}
	return OrderBook{
		Asks: asks,
		Bids: bids,
	}, nil
}

// Convert decoded JSON object to list of orders.
func decodeOrders(v interface{}) (Orders, error) {
	orderList, err := jsonGetList(v)
	if err != nil {
		return nil, err
	}
	res := Orders{}
	for _, v := range orderList {
		order, err := decodeOrder(v)
		if err != nil {
			return nil, err
		}
		res = append(res, order)
	}
	return res, nil
}

// Convert decoded JSON object to Order struct.
func decodeOrder(v interface{}) (Order, error) {
	m, err := jsonGetMap(v)
	if err != nil {
		return Order{}, err
	}
	id, err := jsonGetInt(m["id"])
	if err != nil {
		return Order{}, err
	}
	side, err := jsonGetString(m["side"])
	if err != nil {
		return Order{}, err
	}
	ord_type, err := jsonGetString(m["ord_type"])
	if err != nil {
		return Order{}, err
	}
	price, err := jsonGetFloat(m["price"])
	if err != nil {
		return Order{}, err
	}
	avg_price, err := jsonGetFloat(m["avg_price"])
	if err != nil {
		return Order{}, err
	}
	state, err := jsonGetString(m["state"])
	if err != nil {
		return Order{}, err
	}
	market, err := jsonGetString(m["market"])
	if err != nil {
		return Order{}, err
	}
	created_at, err := jsonGetTimeFromText(m["created_at"])
	if err != nil {
		return Order{}, err
	}
	volume, err := jsonGetFloat(m["volume"])
	if err != nil {
		return Order{}, err
	}
	remaining_volume, err := jsonGetFloat(m["remaining_volume"])
	if err != nil {
		return Order{}, err
	}
	executed_volume, err := jsonGetFloat(m["executed_volume"])
	if err != nil {
		return Order{}, err
	}
	trades_count, err := jsonGetInt(m["trades_count"])
	if err != nil {
		return Order{}, err
	}
	return Order{
		ID:              id,
		Side:            side,
		OrdType:         ord_type,
		Price:           price,
		AvgPrice:        avg_price,
		State:           state,
		Market:          market,
		CreatedAt:       created_at,
		Volume:          volume,
		RemainingVolume: remaining_volume,
		ExecutedVolume:  executed_volume,
		TradesCount:     trades_count,
	}, nil
}

// Convert decoded JSON object to History list.
func decodeHistory(v interface{}) (History, error) {
	entries, err := jsonGetList(v)
	if err != nil {
		return nil, err
	}
	history := History{}
	for _, v := range entries {
		m, err := jsonGetMap(v)
		if err != nil {
			return nil, err
		}
		id, err := jsonGetInt(m["id"])
		if err != nil {
			return nil, err
		}
		price, err := jsonGetFloat(m["price"])
		if err != nil {
			return nil, err
		}
		volume, err := jsonGetFloat(m["volume"])
		if err != nil {
			return nil, err
		}
		funds, err := jsonGetFloat(m["funds"])
		if err != nil {
			return nil, err
		}
		market, err := jsonGetString(m["market"])
		if err != nil {
			return nil, err
		}
		created_at, err := jsonGetTimeFromText(m["created_at"])
		if err != nil {
			return nil, err
		}
		history = append(history, HistoryEntry{
			ID:        id,
			Price:     price,
			Volume:    volume,
			Funds:     funds,
			Market:    market,
			CreatedAt: created_at,
		})
	}
	return history, nil
}

// Convert decoded JSON to user info struct.
func decodeUserInfo(v interface{}) (*UserInfo, error) {
	m, err := jsonGetMap(v)
	if err != nil {
		return nil, err
	}
	email, err := jsonGetString(m["email"])
	if err != nil {
		return nil, err
	}
	activated, err := jsonGetBool(m["activated"])
	if err != nil {
		return nil, err
	}
	accList, err := jsonGetList(m["accounts"])
	if err != nil {
		return nil, err
	}
	accounts := []Account{}
	for _, e := range accList {
		m, err := jsonGetMap(e)
		if err != nil {
			return nil, err
		}
		currency, err := jsonGetString(m["currency"])
		if err != nil {
			return nil, err
		}
		balance, err := jsonGetFloat(m["balance"])
		if err != nil {
			return nil, err
		}
		locked, err := jsonGetFloat(m["locked"])
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, Account{
			Currency: currency,
			Balance:  balance,
			Locked:   locked,
		})
	}
	return &UserInfo{
		Email:     email,
		Activated: activated,
		Accounts:  accounts,
	}, nil
}

func decodeUserTrades(v interface{}) (Trades, error) {
	tradeList, err := jsonGetList(v)
	if err != nil {
		return nil, err
	}
	res := Trades{}
	for _, e := range tradeList {
		m, err := jsonGetMap(e)
		if err != nil {
			return nil, err
		}
		id, err := jsonGetInt(m["id"])
		if err != nil {
			return nil, err
		}
		price, err := jsonGetFloat(m["price"])
		if err != nil {
			return nil, err
		}
		volume, err := jsonGetFloat(m["volume"])
		if err != nil {
			return nil, err
		}
		funds, err := jsonGetFloat(m["funds"])
		if err != nil {
			return nil, err
		}
		market, err := jsonGetString(m["market"])
		if err != nil {
			return nil, err
		}
		created_at, err := jsonGetTimeFromText(m["created_at"])
		if err != nil {
			return nil, err
		}
		side, err := jsonGetString(m["side"])
		if err != nil {
			return nil, err
		}
		res = append(res, Trade{
			ID:        id,
			Price:     price,
			Volume:    volume,
			Funds:     funds,
			Market:    market,
			CreatedAt: created_at,
			Side:      side,
		})
	}
	return res, nil
}

func decodeError(r io.Reader) (error, bool) {
	j, err := DecodeJSON(r)
	if err != nil {
		return nil, false
	}
	m, err := jsonGetMap(j)
	if err != nil {
		return nil, false
	}
	m, err = jsonGetMap(m["error"])
	if err != nil {
		return nil, false
	}
	code, err := jsonGetInt(m["code"])
	if err != nil {
		return nil, false
	}
	msg, err := jsonGetString(m["message"])
	if err != nil {
		return nil, false
	}
	return fmt.Errorf("%d: %s", code, msg), true
}
