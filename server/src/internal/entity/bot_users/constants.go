package bot_users

import (
	"fmt"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusNotExists SubscriptionStatus = "not_exists"
)

const (
	FilterKeyQueryPrefix              = "query"
	FilterKeyJoinDatePrefix           = "join_date"
	FilterKeyBirthDatePrefix          = "birth_date"
	FilterKeySubscriptionStatusPrefix = "purchases"
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
		return rangeKey(FilterKeyJoinDatePrefix, fromValid, tillValid)
	}

	FilterKeySubscriptionStatus = func(valid bool, status SubscriptionStatus) string {
		return fmt.Sprintf("%s:%t:%s", FilterKeySubscriptionStatusPrefix, valid, status)
	}
)

const BilFileNameTemplate = "bill_%d"

const (
	PaymentTypeIDP2P = 1
)
