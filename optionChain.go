package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OptionChain struct {
	Symbol            string      `json:"symbol"`
	Status            string      `json:"status"`
	Underlying        Underlying  `json:"underlying"`
	Strategy          string      `json:"strategy"` //enum[SINGLE, ANALYTICAL, COVERED, VERTICAL, CALENDAR, STRANGLE, STRADDLE, BUTTERFLY, CONDOR, DIAGONAL, COLLAR, ROLL]
	Interval          float64     `json:"interval"`
	IsDelayed         bool        `json:"isDelayed"`
	IsIndex           bool        `json:"isIndex"`
	DaysToExpiration  float64     `json:"daysToExpiration"`
	InterestRate      float64     `json:"interestRate"`
	UnderlyingPrice   float64     `json:"underlyingPrice"`
	Volatility        float64     `json:"volatility"`
	NumberOfContracts int64       `json:"numberOfContracts"`
	CallExpDateMap    interface{} `json:"callExpDateMap"`
	PutExpDateMap     interface{} `json:"putExpDateMap"`
	CallMap           []Option
	PutMap            []Option
}

type StrikePriceMap struct {
}

type Option struct {
	PutCall                string               `json:"putCall"` //enum[PUT, CALL]
	Symbol                 string               `json:"symbol"`
	Description            string               `json:"description"`
	ExchangeName           string               `json:"exchangeName"`
	Bid                    float64              `json:"bid"`
	Ask                    float64              `json:"ask"`
	Last                   float64              `json:"last"`
	Mark                   float64              `json:"mark"`
	BidSize                int64                `json:"bidSize"`
	AskSize                int64                `json:"askSize"`
	BidAskSize             string               `json:"bidAskSize"`
	LastSize               int64                `json:"lastSize"`
	HighPrice              float64              `json:"highPrice"`
	LowPrice               float64              `json:"lowPrice"`
	OpenPrice              float64              `json:"openPrice"`
	ClosePrice             float64              `json:"closePrice"`
	TotalVolume            int64                `json:"totalVolume"`
	QuoteTimeInLong        int64                `json:"quoteTimeInLong"`
	TradeTimeInLong        int64                `json:"tradeTimeInLong"`
	NetChange              float64              `json:"netChange"`
	Volatility             float64              `json:"volatility"`
	Delta                  float64              `json:"delta"`
	Gamma                  float64              `json:"gamma"`
	Theta                  float64              `json:"theta"`
	Vega                   float64              `json:"vega"`
	Rho                    float64              `json:"rho"`
	TimeValue              float64              `json:"timeValue"`
	OpenInterest           float64              `json:"openInterest"`
	TheoreticalOptionValue float64              `json:"theoreticalOptionValue"`
	TheoreticalVolatility  float64              `json:"theoreticalVolatility"`
	OptionDeliverablesList []OptionDeliverables `json:"optionDeliverablesList"`
	StrikePrice            float64              `json:"strikePrice"`
	ExpirationDate         int64                `json:"expirationDate"`
	DaysToExpiration       int64                `json:"daysToExpiration"`
	ExpirationType         string               `json:"expirationType"`
	LastTradingDay         int64                `json:"lastTradingDay"`
	Multiplier             float64              `json:"multiplier"`
	SettlementType         string               `json:"settlementType"`
	DeliverableNote        string               `json:"deliverableNote"`
	IsIndexOption          bool                 `json:"isIndexOption"`
	PercentChange          float64              `json:"percentChange"`
	MarkChange             float64              `json:"markChange"`
	MarkPercentChange      float64              `json:"markPercentChange"`
	InTheMoney             bool                 `json:"inTheMoney"`
	Mini                   bool                 `json:"mini"`
	NonStandard            bool                 `json:"nonStandard"`
}

type OptionDeliverables struct {
	Symbol           string `json:"symbol"`
	AssetType        string `json:"assetType"`
	DeliverableUnits string `json:"deliverableUnits"`
	CurrencyType     string `json:"currencyType"`
}

type Underlying struct {
	Ask               float64 `json:"ask"`
	AskSize           int64   `json:"askSize"`
	Bid               float64 `json:"bid"`
	BidSize           int64   `json:"bidSize"`
	Change            float64 `json:"change"`
	Close             float64 `json:"close"`
	Delayed           bool    `json:"delayed"`
	Description       string  `json:"description"`
	ExchangeName      string  `json:"exchangeName"` //enum[IND, ASE, NYS, NAS, NAP, PAC, OPR, BATS]
	FiftyTwoWeekHigh  float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow   float64 `json:"fiftyTwoWeekLow"`
	HighPrice         float64 `json:"highPrice"`
	Last              float64 `json:"last"`
	LowPrice          float64 `json:"lowPrice"`
	Mark              float64 `json:"mark"`
	MarkChange        float64 `json:"markChange"`
	MarkPercentChange float64 `json:"markPercentChange"`
	OpenPrice         float64 `json:"openPrice"`
	PercentChange     float64 `json:"percentChange"`
	QuoteTime         int64   `json:"quoteTime"`
	Symbol            string  `json:"symbol"`
	TotalVolume       int64   `json:"totalVolume"`
	TradeTime         int64   `json:"tradeTime"`
}

type ExpirationDate struct {
	Date string `json:"date"`
}

func GetOptionChain(client *http.Client, consumerKey string, symbol string) OptionChain {
	request, err := http.NewRequest("GET", urlGetOptionChain, nil)
	handleFatalErr("GetOptionChain/request", err)
	q := request.URL.Query()
	q.Add("apikey", consumerKey)
	q.Add("symbol", symbol)
	q.Add("includeQuotes", "true")
	request.URL.RawQuery = q.Encode()
	resp, err := client.Do(request)
	handleFatalErr("GetOptionChain/resp", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handleFatalErr("GetOptionChain/body", err)
	var oc OptionChain
	json.Unmarshal(body, &oc)
	oc.CallMap = parseExpDateMap(oc.CallExpDateMap)
	oc.PutMap = parseExpDateMap(oc.PutExpDateMap)
	return oc
}

func parseExpDateMap(obj interface{}) []Option {
	var opts []Option
	expDateMap := obj.(map[string]interface{})
	for _, strikePriceMap := range expDateMap {
		switch strikePriceMap.(type) {
		case interface{}:
			optionMap := strikePriceMap.(map[string]interface{})
			for _, optionEntry := range optionMap {
				option := optionEntry.([]interface{})[0].(map[string]interface{})
				var opt Option
				opt.PutCall = AssertString(option["putCall"])
				opt.Symbol = AssertString(option["symbol"])
				opt.Description = AssertString(option["description"])
				opt.ExchangeName = AssertString(option["exchangeName"])
				opt.Bid = AssertFloat64(option["bid"])
				opt.Ask = AssertFloat64(option["ask"])
				opt.Last = AssertFloat64(option["last"])
				opt.Mark = AssertFloat64(option["mark"])
				opt.BidSize = AssertInt64(option["bidSize"])
				opt.AskSize = AssertInt64(option["askSize"])
				opt.BidAskSize = AssertString(option["bidAskSize"])
				opt.LastSize = AssertInt64(option["lastSize"])
				opt.HighPrice = AssertFloat64(option["highPrice"])
				opt.LowPrice = AssertFloat64(option["lowPrice"])
				opt.OpenPrice = AssertFloat64(option["openPrice"])
				opt.ClosePrice = AssertFloat64(option["closePrice"])
				opt.TotalVolume = AssertInt64(option["totalVolume"])
				opt.QuoteTimeInLong = AssertInt64(option["quoteTimeInLong"])
				opt.TradeTimeInLong = AssertInt64(option["tradeTimeInLong"])
				opt.NetChange = AssertFloat64(option["netChange"])
				opt.Volatility = AssertFloat64(option["volatility"])
				opt.Delta = AssertFloat64(option["delta"])
				opt.Gamma = AssertFloat64(option["gamma"])
				opt.Theta = AssertFloat64(option["theta"])
				opt.Vega = AssertFloat64(option["vega"])
				opt.Rho = AssertFloat64(option["rho"])
				opt.TimeValue = AssertFloat64(option["timeValue"])
				opt.OpenInterest = AssertFloat64(option["openInterest"])
				opt.TheoreticalOptionValue = AssertFloat64(option["theoreticalOptionValue"])
				opt.TheoreticalVolatility = AssertFloat64(option["theoreticalVolatility"])
				opt.StrikePrice = AssertFloat64(option["strikePrice"])
				opt.ExpirationDate = AssertInt64(option["expirationDate"])
				opt.DaysToExpiration = AssertInt64(option["daysToExpiration"])
				opt.ExpirationType = AssertString(option["expirationType"])
				opt.LastTradingDay = AssertInt64(option["lastTradingDay"])
				opt.Multiplier = AssertFloat64(option["multiplier"])
				opt.SettlementType = AssertString(option["settlementType"])
				opt.DeliverableNote = AssertString(option["deliverableNote"])
				opt.IsIndexOption = AssertBool(option["isIndexOption"])
				opt.PercentChange = AssertFloat64(option["percentChange"])
				opt.MarkChange = AssertFloat64(option["markChange"])
				opt.MarkPercentChange = AssertFloat64(option["markPercentChange"])
				opt.InTheMoney = AssertBool(option["inTheMoney"])
				opt.Mini = AssertBool(option["mini"])
				opt.NonStandard = AssertBool(option["nonStandard"])
				opts = append(opts, opt)
			}
		default:
			fmt.Println("Expecting a JSON object, got something wrong")
		}
	}
	return opts
}
