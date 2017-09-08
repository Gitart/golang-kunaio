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

package main

import (
	"fmt"
	"kunaio"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	USAGE = "Usage:\n" +
		"\t%s -h|--help                show this memo;\n" +
		"\t%s time                     show current server time;\n" +
		"\t%s [options] stats          show latest trade statistics;\n" +
		"\t%s [options] [--uah] sell LIMIT\n" +
		"\t                            show order book - asks;\n" +
		"\t%s [options] [--uah] buy LIMIT\n" +
		"\t                            show order book - bids;\n" +
		"\t%s [options] history        show trade history;\n" +
		"\t%s [options] userinfo       show user info and assets;\n" +
		"\t%s [options] userorders     show current orders for user;\n" +
		"\t%s [options] usertrades     show history of user trades;\n" +
		"\t%s [options] [--uah] addorder SIDE VOLUME PRICE\n" +
		"\t                            create new order. SIDE - buy or sell;\n" +
		"\t                            VOLUME - in ICO; PRICE - price for 1 ICO;\n" +
		"\t%s [options] delorder ORDER_ID\n" +
		"\t                            delete existing order;\n" +
		"\t%s [options] delall         delete all existing orders;\n" +
		"Options are:\n" +
		"\t--unix                      print date/time as Unix timestamp.\n" +
		"\t--market MARKET             set market type. Default is btcuah.\n" +
		"\t--akey ACCESS_KEY           set API access key.\n" +
		"\t--skey SECRET_KEY           set API secret key.\n" +
		"Environment variables:\n" +
		"\tKUNAIO_MARKET               market type. Default is btcuah.\n" +
		"\tKUNAIO_ACCESS_KEY           API access key.\n" +
		"\tKUNAIO_SECRET_KEY           API secret key.\n"
)

// Configurations
var (
	gMarket     string = kunaio.BTCUAH
	gTimeLayout        = "2006-01-02T15:04:05-0700"
	gAKey       string
	gSKey       string
	gUAH        bool
	gUnix       bool
)

// Entry point.
func main() {
	// read environment variables
	if s, ok := os.LookupEnv("KUNAIO_MARKET"); ok {
		gMarket = s
	}
	if s, ok := os.LookupEnv("KUNAIO_ACCESS_KEY"); ok {
		gAKey = s
	}
	if s, ok := os.LookupEnv("KUNAIO_SECRET_KEY"); ok {
		gSKey = s
	}
	// parse command line args
	args := os.Args[1:]
	for 0 < len(args) && strings.HasPrefix(args[0], "-") {
		arg := args[0]
		args = args[1:]
		switch arg {
		case "-h":
			usage()
			os.Exit(0)
		case "--help":
			usage()
			os.Exit(0)
		case "--market":
			gMarket = args[0]
			valid := false
			for _, v := range kunaio.SupportedMarkets() {
				if v == gMarket {
					valid = true
					break
				}
			}
			if !valid {
				fatalf("invalid market: %#v. Valid are: %v\n",
					gMarket, kunaio.SupportedMarkets())
			}
			args = args[1:]
		case "--akey":
			gAKey = args[0]
			args = args[1:]
		case "--skey":
			gSKey = args[0]
			args = args[1:]
		case "--uah":
			gUAH = true
		case "--unix":
			gUnix = true
		case "--":
			break
		default:
			fatalf("unknown option: %v", arg)
		}
	}
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}
	cmd := args[0]
	args = args[1:]
	switch cmd {
	case "time":
		t, err := kunaio.GetServerTime()
		if err != nil {
			fatalf("get server time: %s", err)
		}
		fmt.Printf("%s\n", tts(t))
	case "stats":
		stats, err := kunaio.GetLatestStats(gMarket)
		if err != nil {
			fatalf("get stats: %s", err)
		}
		fmt.Printf("%#v\n", stats)
	case "sell":
		var limit float64
		if len(args) == 1 {
			f, err := strconv.ParseFloat(args[0], 64)
			if err != nil || f < 0 {
				fatalf("invalid limit (%s): %s", args[0], err)
			}
			limit = f
		}
		obook, err := kunaio.GetOrderBook(gMarket)
		if err != nil {
			fatalf("get order book: %s", err)
		}
		fmt.Printf("SELL:%10s %15s %15s %15s %15s %15s\n",
			"PRICE", "VOLUME", "FUNDS",
			"AVG_PRICE", "SUM_VOLUME", "SUM_FUNDS")
		var (
			sumVolume float64
			sumFunds  float64
		)
		needBreak := false
		for _, e := range obook.Asks {
			if 0 < limit {
				if !gUAH && limit <= sumVolume+e.RemainingVolume {
					e.RemainingVolume = limit - sumVolume
					needBreak = true
				} else if gUAH && limit <= (sumFunds+e.RemainingVolume*e.Price) {
					e.RemainingVolume = (limit - sumFunds) / e.Price
					needBreak = true
				}
			}
			sumVolume += e.RemainingVolume
			sumFunds += e.RemainingVolume * e.Price
			fmt.Printf("%15.7f %15.7f %15.7f %15.7f %15.7f %15.7f\n",
				e.Price, e.RemainingVolume,
				e.RemainingVolume*e.Price,
				sumFunds/sumVolume,
				sumVolume, sumFunds)
			if needBreak {
				break
			}
		}
	case "buy":
		var limit float64
		if len(args) == 1 {
			f, err := strconv.ParseFloat(args[0], 64)
			if err != nil || f < 0 {
				fatalf("invalid limit (%s): %s", args[0], err)
			}
			limit = f
		}
		obook, err := kunaio.GetOrderBook(gMarket)
		if err != nil {
			fatalf("get order book: %s\n", err)
		}
		fmt.Printf("BUY:%11s %15s %15s %15s %15s %15s\n",
			"PRICE", "VOLUME", "FUNDS",
			"AVG_PRICE", "SUM_VOLUME", "SUM_FUNDS")
		var (
			sumVolume float64
			sumFunds  float64
		)
		needBreak := false
		for _, e := range obook.Bids {
			if 0 < limit {
				if !gUAH && limit <= sumVolume+e.RemainingVolume {
					e.RemainingVolume = limit - sumVolume
					needBreak = true
				} else if gUAH && limit <= (sumFunds+e.RemainingVolume*e.Price) {
					e.RemainingVolume = (limit - sumFunds) / e.Price
					needBreak = true
				}
			}
			sumVolume += e.RemainingVolume
			sumFunds += e.RemainingVolume * e.Price
			fmt.Printf("%15.7f %15.7f %15.7f %15.7f %15.7f %15.7f\n",
				e.Price, e.RemainingVolume,
				e.RemainingVolume*e.Price,
				sumFunds/sumVolume,
				sumVolume, sumFunds)
			if needBreak {
				break
			}
		}
	case "history":
		hist, err := kunaio.GetTradeHistory(gMarket)
		if err != nil {
			fatalf("get trade history: %s", err)
		}
		fmt.Printf("%24s %15s %15s %15s\n",
			"WHEN", "PRICE", "VOLUME", "FUNDS")
		for _, e := range hist {
			fmt.Printf("%15s %15.7f %15.7f %15.7f\n",
				tts(e.CreatedAt),
				e.Price, e.Volume, e.Funds)
		}
		fmt.Printf("TOTAL%51.7f %15.7f\n",
			hist.SumVolume(), hist.SumFunds())
		fmt.Printf("MIN%37.7f %15.7f %15.7f\n",
			hist.MinPrice(), hist.MinVolume(), hist.MinFunds())
		fmt.Printf("AVG%37.7f %15.7f %15.7f\n",
			hist.AvgPrice(), hist.AvgVolume(), hist.AvgFunds())
		fmt.Printf("MAX%37.7f %15.7f %15.7f\n",
			hist.MaxPrice(), hist.MaxVolume(), hist.MaxFunds())
	case "userinfo":
		checkReqs()
		info, err := kunaio.GetUserInfo(gAKey, gSKey)
		if err != nil {
			fatalf("get user info: %s", err)
		}
		fmt.Printf("email:\t%s\nactive:\t%v\n"+
			"Accounts: %6s %15s %15s\n",
			info.Email, info.Activated,
			"CURRENCY", "BALANCE", "LOCKED")
		for _, a := range info.Accounts {
			fmt.Printf("%18s %15.7f %15.7f\n",
				a.Currency, a.Balance, a.Locked)
		}
	case "userorders":
		checkReqs()
		orders, err := kunaio.GetUserOrders(gAKey, gSKey, gMarket)
		if err != nil {
			fatalf("get user orders: %s", err)
		}
		fmt.Printf("%24s %8s %4s %8s %15s %15s %15s %15s %15s %15s %15s %15s %6s\n",
			"WHEN", "MARKET", "SIDE", "ID", "PRICE",
			"AVG_PRICE", "VOLUME", "FUNDS", "REMAINING", "REM_FUNDS",
			"EXECUTED", "EXEC_FUNDS", "TRADES")
		for _, e := range orders {
			fmt.Printf("%6s %8s %4s %8d %15.7f %15.7f %15.7f %15.7f "+
				"%15.7f %15.7f %15.7f %15.7f %6d\n",
				tts(e.CreatedAt),
				e.Market, e.Side, e.ID, e.Price,
				e.AvgPrice,
				e.Volume, e.Volume*e.Price,
				e.RemainingVolume, e.RemainingVolume*e.Price,
				e.ExecutedVolume, e.ExecutedVolume*e.Price,
				e.TradesCount)
		}
	case "usertrades":
		checkReqs()
		trades, err := kunaio.GetUserTrades(gAKey, gSKey, gMarket)
		if err != nil {
			fatalf("get user trades: %s", err)
		}
		fmt.Printf("%24s %4s %15s %15s %15s %15s %15s %15s\n",
			"WHEN", "SIDE", "PRICE", "VOLUME", "FUNDS",
			"AVG_PRICE", "SUM_VOLUME", "SUM_FUNDS")
		var (
			sumVolume float64
			sumFunds  float64
		)
		for _, e := range trades {
			sumVolume += e.Volume
			sumFunds += e.Funds
			fmt.Printf("%24s %4s %15.7f %15.7f %15.7f %15.7f %15.7f %15.7f\n",
				tts(e.CreatedAt), e.Side,
				e.Price, e.Volume, e.Funds,
				sumFunds/sumVolume,
				sumVolume, sumFunds)
		}
		fmt.Printf("AVERAGE%38.7f %15.7f %15.7f\n",
			trades.AvgPrice(), trades.AvgVolume(), trades.AvgFunds())
	case "addorder":
		checkReqs()
		if len(args) != 3 {
			fatalf("bad args count: %d (expected 3)", len(args))
		}
		side := strings.Trim(strings.ToLower(args[0]), " \t\n\r")
		if side != "sell" && side != "buy" {
			fatalf("invalid SIDE arg (%s). Valid values are: sell, buy", side)
		}
		volume, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fatalf("invalid VOLUME arg (%s): %s", args[1], err)
		}
		price, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			fatalf("invalid PRICE arg (%s): %s", args[2], err)
		}
		if gUAH {
			volume /= price
		}
		order, err := kunaio.NewOrder(gAKey, gSKey, gMarket,
			side, volume, price)
		if err != nil {
			fatalf("new order: %s", err)
		}
		fmt.Printf("Order created:\n%s", formatOrder(order))
	case "delorder":
		checkReqs()
		if len(args) != 1 {
			fatalf("bad args count: %d (expected 1)", len(args))
		}
		i, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fatalf("invalid order ID %s: %s", args[0], err)
		}
		order, err := kunaio.CancelOrder(gAKey, gSKey, int(i))
		if err != nil {
			fatalf("cancel order: %s", err)
		}
		fmt.Printf("Order deleted:\n%s", formatOrder(order))
	case "delall":
		checkReqs()
		orders, err := kunaio.GetUserOrders(gAKey, gSKey, gMarket)
		if err != nil {
			fatalf("get user orders: %s", err)
		}
		for _, order := range orders {
			order, err := kunaio.CancelOrder(gAKey, gSKey, order.ID)
			if err != nil {
				fatalf("cancel order: %s", err)
			}
			fmt.Printf("Order deleted:\n%s", formatOrder(order))
		}
	default:
		usage()
		os.Exit(1)
	}
}

func formatOrder(order kunaio.Order) string {
	return fmt.Sprintf(""+
		"  ID             : %d\n"+
		"  Side           : %s\n"+
		"  OrdType        : %s\n"+
		"  Price          : %16.8f\n"+
		"  AvgPrice       : %16.8f\n"+
		"  State          : %s\n"+
		"  Market         : %s\n"+
		"  CreatedAt      : %s\n"+
		"  Volume         : %16.8f\n"+
		"  RemainingVolume: %16.8f\n"+
		"  ExecutedVolume : %16.8f\n"+
		"  TradesCount    : %d\n",
		order.ID, order.Side, order.OrdType, order.Price,
		order.AvgPrice, order.State, order.Market,
		tts(order.CreatedAt),
		order.Volume, order.RemainingVolume,
		order.ExecutedVolume, order.TradesCount)
}

// Check access requisites.
func checkReqs() {
	if gAKey == "" {
		fatalf("--akey option is required")
	}
	if gSKey == "" {
		fatalf("--skey option is required")
	}
}

// Show usage info.
func usage() {
	s := os.Args[0]
	fmt.Printf(USAGE, s, s, s, s, s, s, s, s, s, s, s, s)
}

// Print error report and terminate with exit code 1.
func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}

// Format date/time to string.
func tts(t time.Time) string {
	if gUnix {
		return fmt.Sprintf("%d", t.Unix())
	}
	return t.Format(gTimeLayout)
}
