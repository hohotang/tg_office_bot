package constant

const (
	RECOMMEND_NUM = 3

	PRICE_LOW  = "低"
	PRICE_MID  = "中"
	PRICE_HIGH = "高"

	CALLBACK_LOW        = "low"
	CALLBACK_MID        = "mid"
	CALLBACK_HIGH       = "high"
	CALLBACK_RAND       = "rand"
	CALLBACK_ALL        = "all"
	CALLBACK_RESTAURANT = "restaurant"
	CALLBACK_ADD        = "add"
	CALLBACK_CONFIRM    = "confirm"
	CALLBACK_DENY       = "deny"
	CALLBACK_UPDATE     = "update"
	CALL_BACK_TEST      = "update_test"

	BRANCH_TEST = "test/tg_auto"
)

const (
	ADD_STATE_NONE = iota
	ADD_STATE_NAME
	ADD_STATE_PRICE
	ADD_STATE_DESCR
	ADD_STATE_CHECK
)

// PriceLevelMap maps price levels to integers.
var PriceLevelMap = map[string]int{
	PRICE_LOW:     0,
	PRICE_MID:     1,
	PRICE_HIGH:    2,
	CALLBACK_LOW:  0,
	CALLBACK_MID:  1,
	CALLBACK_HIGH: 2,
}

var PriceStrMap = map[string]string{
	CALLBACK_LOW:  PRICE_LOW,
	CALLBACK_MID:  PRICE_MID,
	CALLBACK_HIGH: PRICE_HIGH,
	PRICE_LOW:     CALLBACK_LOW,
	PRICE_MID:     CALLBACK_MID,
	PRICE_HIGH:    CALLBACK_HIGH,
}
