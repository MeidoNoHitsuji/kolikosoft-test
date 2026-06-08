package model

import "encoding/json"

const (
	MarketHashNameKey   = "market_hash_name"
	MinPriceKey         = "min_price"
	TradableMinPriceKey = "tradable_min_price"
)

// С одной стороны хранить невалидируемые значения плохо
// С другой в самом ТЗ было указано, что нам нужны просто все параметры, так что для оптимизации я решил не создавать
// под них отдельную структуру
type Item map[string]json.RawMessage

func (i Item) MarketHashName() string {
	var name string
	_ = json.Unmarshal(i[MarketHashNameKey], &name)
	return name
}
