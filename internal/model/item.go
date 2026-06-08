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

// Т.к. на данный момент нам нужно только 1 поле у предмета, то при доставании предмета из кеша ограничусь им
// ну и именем.. Просто чтобы было)
type ItemInfo struct {
	MarketHashName string `json:"market_hash_name"`
	// Данное поле было заполнено в большинстве случаев и представляет оптимальную цену за товар.
	// Так что я отталкиваюсь от данной цены.
	SuggestedPrice *float64 `json:"suggested_price"`
}
