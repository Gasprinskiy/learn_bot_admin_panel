package bot_users

import (
	"fmt"
)

const (
	FilterKeyQueryPrefix     = "query"
	FilterKeyJoinDatePrefix  = "join_date"
	FilterKeyBirthDatePrefix = "birth_date"
	FilterKeyPurchasesPrefix = "purchases"
)

func rangeKey(prefix string, fromValid, tillValid bool) string {
	return fmt.Sprintf("%s:%t:%t", prefix, fromValid, tillValid)
}

func validKey(prefix string, valid bool) string {
	return fmt.Sprintf("%s:%t", prefix, valid)
}

var (
	FilterKeyQuery = func(valid bool) string {
		return validKey(FilterKeyQueryPrefix, valid)
	}

	FilterKeyBirthDate = func(fromValid, tillValid bool) string {
		return rangeKey(FilterKeyBirthDatePrefix, fromValid, tillValid)
	}

	FilterKeyJoinDate = func(fromValid, tillValid bool) string {
		return rangeKey(FilterKeyBirthDatePrefix, fromValid, tillValid)
	}

	FilterKeyPurchases = func(valid bool) string {
		return validKey(FilterKeyPurchasesPrefix, valid)
	}
)
