package ethereum

import (
	"fmt"
	"strings"
	"sync"
)

type Data struct {
	Ethereum Ethereum `json:"ethereum"`
}

type Ethereum struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Time           string  `json:"time"`
	GasPrice       float64 `json:"gasPrice"`
	GasValue       float64 `json:"gasValue"`
	Average        float64 `json:"average"`
	MaxGasPrice    float64 `json:"maxGasPrice"`
	MedianGasPrice float64 `json:"medianGasPrice"`
}

type AverageData struct {
	Sum   float64
	Count float64
}

type Distribution struct {
	Max float64 `json:"max"`
	Min float64 `json:"min"`
}

type Answer struct {
	AveragePricePerDay    map[string]float64      `json:"average_price_per_day"`
	SpentInMonth          map[string]float64      `json:"spent_in_month"`
	FrequencyDistribution map[string]Distribution `json:"frequency_distribution"`
	TotalСosts            float64                 `json:"total_costs"`
}

// Считает частотное распределение цены по часам
func (ethereum *Ethereum) FrequencyDistribution(wg *sync.WaitGroup, answer *Answer) {
	defer wg.Done()

	transaction := ethereum.Transactions
	freqDistr := make(map[string]Distribution)
	for i, v := range transaction {
		time := strings.Split(transaction[i].Time, " ")[1]
		hour := strings.Split(time, ":")[0]
		if _, ok := freqDistr[hour]; !ok {
			freqDistr[hour] = Distribution{
				Max: 0,
				Min: v.MaxGasPrice,
			}
		}

		if freqDistr[hour].Max < v.GasPrice {
			freqDistr[hour] = Distribution{
				Max: v.GasPrice,
				Min: freqDistr[hour].Min,
			}
		}
		if freqDistr[hour].Min > v.GasPrice {
			freqDistr[hour] = Distribution{
				Max: freqDistr[hour].Max,
				Min: v.GasPrice,
			}
		}
	}

	answer.FrequencyDistribution = freqDistr
}

// Считает среднюю цену gas за день
func (ethereum *Ethereum) AveragePricePerDay(wg *sync.WaitGroup, answer *Answer) {
	defer wg.Done()

	transaction := ethereum.Transactions
	average := make(map[string]AverageData)
	for i := range transaction {
		time := strings.Split(transaction[i].Time, " ")[0]
		month := strings.Split(time, "-")[1]
		day := strings.Split(time, "-")[2]
		date := fmt.Sprintf("%s.%s", month, day)

		average[date] = AverageData{
			Sum:   average[date].Sum + transaction[i].GasPrice,
			Count: average[date].Count + 1,
		}
	}

	an := make(map[string]float64)
	for i, v := range average {
		an[i] = v.Sum / v.Count
	}

	answer.AveragePricePerDay = an
}

// Считает расходы gas помесячно
func (ethereum *Ethereum) SpentInMonth(wg *sync.WaitGroup, answer *Answer) {
	defer wg.Done()

	transaction := ethereum.Transactions
	spent := make(map[string]float64)
	for i := range transaction {
		month := strings.Split(transaction[i].Time, "-")[1]
		spent[month] += transaction[i].GasValue
	}

	answer.SpentInMonth = spent
}

// Считает затраты на gas за весь период
func (ethereum *Ethereum) TotalСosts(wg *sync.WaitGroup, answer *Answer) {
	defer wg.Done()

	transaction := ethereum.Transactions
	var result float64
	for _, tr := range transaction {
		result += tr.GasPrice * tr.GasValue
	}

	answer.TotalСosts = result
}
