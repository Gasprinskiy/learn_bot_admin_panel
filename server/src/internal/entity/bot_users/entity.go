package bot_users

import (
	"learn_bot_admin_panel/tools/sql_null"
	"time"
)

type BotUserCommonData struct {
	UID          int64     `db:"u_id" json:"u_id"`
	TgID         int64     `db:"tg_id" json:"tg_id"`
	TgUserName   string    `db:"tg_user_name" json:"tg_user_name" excel_head:"Юзернейм" excel_cell:"string"`
	FirstName    string    `db:"first_name" json:"first_name" excel_head:"Имя" excel_cell:"string"`
	LastName     string    `db:"last_name" json:"last_name" excel_head:"Фамилия" excel_cell:"string"`
	BirthDate    time.Time `db:"birth_date" json:"birth_date" excel_head:"Дата рождения" excel_cell:"date"`
	PhoneNumber  string    `db:"phone_number" json:"phone_number" excel_head:"Номер телефона" excel_cell:"string"`
	JoinDate     time.Time `db:"join_date" json:"join_date" excel_head:"Дата вступления" excel_cell:"date"`
	RegisterDate time.Time `db:"register_date" json:"register_date"`
}

type BotUserProfile struct {
	UID                int64              `db:"u_id" json:"u_id"`
	TgID               int64              `db:"tg_id" json:"tg_id"`
	TgUserName         string             `db:"tg_user_name" json:"tg_user_name" excel_head:"Юзернейм" excel_cell:"string"`
	FirstName          string             `db:"first_name" json:"first_name" excel_head:"Имя" excel_cell:"string"`
	LastName           string             `db:"last_name" json:"last_name" excel_head:"Фамилия" excel_cell:"string"`
	BirthDate          time.Time          `db:"birth_date" json:"birth_date" excel_head:"Дата рождения" excel_cell:"date"`
	PhoneNumber        string             `db:"phone_number" json:"phone_number" excel_head:"Номер телефона" excel_cell:"string"`
	JoinDate           time.Time          `db:"join_date" json:"join_date" excel_head:"Дата вступления" excel_cell:"date"`
	RegisterDate       time.Time          `db:"register_date" json:"register_date"`
	CommonTotalCount   int                `db:"total_count" json:"-"`
	SubscrID           sql_null.NullInt64 `db:"sub_id" json:"subscription_id"`
	SubscrPurchaseDate sql_null.NullTime  `db:"p_time" json:"-"`
	SubscrTerm         sql_null.NullInt64 `db:"term_in_month" json:"subscription_term"`
	SubscrKickTime     sql_null.NullTime  `db:"kick_time" json:"kick_time"`
	//
	SubscrStatus SubscriptionStatus `json:"subscription_status" excel_head:"Подписка" excel_cell:"string"`
}

func (b *BotUserProfile) SetSubscriptionStatus(value SubscriptionStatus) {
	b.SubscrStatus = value
}

type BotUnregistredUserProfile struct {
	UID              int64               `db:"u_id" json:"u_id"`
	TgID             int64               `db:"tg_id" json:"tg_id"`
	TgUserName       string              `db:"tg_user_name" json:"tg_user_name" excel_head:"Юзернейм" excel_cell:"string"`
	FirstName        sql_null.NullString `db:"first_name" json:"first_name" excel_head:"Имя" excel_cell:"sql_null_string"`
	LastName         sql_null.NullString `db:"last_name" json:"last_name" excel_head:"Фамилия" excel_cell:"sql_null_string"`
	BirthDate        sql_null.NullTime   `db:"birth_date" json:"birth_date" excel_head:"Дата рождения" excel_cell:"sql_null_time"`
	JoinDate         time.Time           `db:"join_date" json:"join_date" excel_head:"Дата вступления" excel_cell:"date"`
	CommonTotalCount int                 `db:"total_count" json:"-"`
}

type BotSubscriptionType struct {
	SubID       int     `db:"sub_id" json:"sub_id"`
	TermInMonth int     `db:"term_in_month" json:"term_in_month"`
	Price       float64 `db:"price" json:"price"`
}

type Purchase struct {
	SubID         int                  `db:"sub_id"`
	UserID        int                  `db:"u_id"`
	PurchaseTime  time.Time            `db:"p_time"`
	Discount      sql_null.NullFloat64 `db:"discount"`
	ManagerID     sql_null.NullInt64   `db:"manager_id"`
	PaymentTypeID int                  `db:"payment_type_id"`
}

type BotUserPurchase struct {
	PID              int                  `db:"p_id" json:"p_id"`
	SubID            int                  `db:"sub_id" json:"sub_id"`
	PurchaseTime     time.Time            `db:"p_time" json:"p_time"`
	KickTime         sql_null.NullTime    `db:"kick_time" json:"kick_time"`
	KickReason       sql_null.NullInt64   `db:"kick_reason" json:"kick_reason"`
	PaymentTypeID    int                  `db:"payment_type_id" json:"payment_type_id"`
	Discount         sql_null.NullFloat64 `db:"discount" json:"discount"`
	ManagerID        sql_null.NullInt64   `db:"manager_id" json:"manager_id"`
	Term             int                  `db:"term_in_month" json:"subscription_term"`
	Price            float64              `db:"price" json:"price"`
	ReceiptFileName  sql_null.NullString  `db:"receipt_file_name" json:"receipt_file_name"`
	ManagerFirstName sql_null.NullString  `db:"manager_first_name" json:"manager_first_name"`
	ManagerLastName  sql_null.NullString  `db:"manager_last_name" json:"manager_last_name"`
	//
	Status SubscriptionStatus `json:"subscription_status"`
}

func (b *BotUserPurchase) SetSubscriptionStatus(value SubscriptionStatus) {
	b.Status = value
}

func NewPurchase(
	subID, userID int,
	pTime time.Time,
	discount sql_null.NullFloat64,
	managerID sql_null.NullInt64,
	paymentTypeID int,
) Purchase {
	return Purchase{
		SubID:         subID,
		UserID:        userID,
		PurchaseTime:  pTime,
		Discount:      discount,
		ManagerID:     managerID,
		PaymentTypeID: paymentTypeID,
	}
}

type BotUserDetailData struct {
	BotUserCommonData
	PurchaseData []BotUserPurchase `json:"purchase_data"`
}

func NewBotUserDetailData(
	commonData BotUserCommonData,
	purchaseData []BotUserPurchase,
) BotUserDetailData {
	return BotUserDetailData{
		BotUserCommonData: commonData,
		PurchaseData:      purchaseData,
	}
}
