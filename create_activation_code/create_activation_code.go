package main

import (
	"flag"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	activation "github.com/jinmukeji/jiujiantang-services/subscription/activation-code"

	"github.com/jinmukeji/go-pkg/v2/crypto/rand"

	"github.com/jinzhu/gorm"
)

const (
	// LetterDigits 字母与数字(激活码中不包含 0、1、2、O、I、Z)
	LetterDigits = "3456789ABCDEFGHJKLMNPQRSTUVWXY"
	// ExpiredDuration 过期时长 1年
	ExpiredDuration = 1
)

var (
	username         string
	password         string
	address          string
	database         string
	charset          string
	parsetime        bool
	locale           string
	subscriptionType int    // 订阅类型
	contractYear     int    // 年份
	maxUserLimits    int    // 最大人数
	key              string // 加密的key
	number           int    // 生成的数目
)

func init() {
	flag.StringVar(&username, "username", "", "missing username")
	flag.StringVar(&password, "password", "", "missing password")
	flag.StringVar(&address, "address", "", "missing address")
	flag.StringVar(&database, "database", "", "missing database")
	flag.StringVar(&charset, "charset", "utf8mb4", "missing charset")
	flag.BoolVar(&parsetime, "parsetime", true, "missing parsetime")
	flag.StringVar(&locale, "locale", "utc", "missing locale")
	flag.IntVar(&subscriptionType, "subscription_type", 2, "missing subscription type")
	flag.IntVar(&contractYear, "contract_year", 1, "missing contract year")
	flag.IntVar(&maxUserLimits, "max_user_limits", 200, "missing max user limits")
	flag.StringVar(&key, "key", "", "missing key")
	flag.IntVar(&number, "number", 0, "missing number")
}

// SubscriptionActivationCode 激活码
type SubscriptionActivationCode struct {
	Code             string     `gorm:"primary_key"`              // 激活码
	MaxUserLimits    int32      `gorm:"column:max_user_limits"`   // 组织下最大用户数量
	ContractYear     int32      `gorm:"column:contract_year"`     // 年限
	SubscriptionType int32      `gorm:"column:subscription_type"` // 0 定制化 1 试用版 2 黄金姆 3 白金姆 4 钻石姆 5 礼品版
	Checksum         string     `gorm:"column:checksum"`          // 校验位
	ExpiredAt        time.Time  `gorm:"column:expired_at"`        // 到期时间
	CreatedAt        time.Time  // 创建时间
	UpdatedAt        time.Time  // 更新时间
	DeletedAt        *time.Time // 删除时间
}

// TableName 返回 SubscriptionActivationCode 所在的表名
func (s SubscriptionActivationCode) TableName() string {
	return "subscription_activation_code"
}

// createSubscriptionActivationCode 创建激活码
func createSubscriptionActivationCode(db *gorm.DB, subscriptionActivationCode *SubscriptionActivationCode) error {
	return db.Create(subscriptionActivationCode).Error
}

// checkSubscriptionActivationCodeExsit 检查激活码是否存在
func checkSubscriptionActivationCodeExsit(db *gorm.DB, code string) (bool, error) {
	var count int
	err := db.Model(&SubscriptionActivationCode{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

func main() {
	flag.Parse()
	db := InitDB()
	defer db.Close()
	if number == 0 {
		panic("number cannot be zero")
	}
	for i := 0; i < number; i++ {
		createActivationCode(db)
	}
}

// 创建单个激活码
func createActivationCode(db *gorm.DB) {
	code, err := rand.RandomStringWithMask(LetterDigits, 16)
	if err != nil {
		return
	}
	// 存在返回上个步骤
	if exsit, _ := checkSubscriptionActivationCodeExsit(db, code); exsit {
		createActivationCode(db)
		return
	}
	helper := activation.NewActivationCodeCipherHelper()
	checksum := helper.Encrypt(code, key, int32(contractYear), int32(maxUserLimits))
	now := time.Now()
	subscriptionActivationCode := &SubscriptionActivationCode{
		Code:             code,
		MaxUserLimits:    int32(maxUserLimits),
		ContractYear:     int32(contractYear),
		SubscriptionType: int32(subscriptionType),
		Checksum:         checksum,
		ExpiredAt:        now.AddDate(ExpiredDuration, 0, 0),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	createSubscriptionActivationCode(db, subscriptionActivationCode)
	fmt.Printf("code=%s,contractYear=%d,maxUserLimits=%d,number=%d", code, contractYear, maxUserLimits, number)
}

//InitDB 初始化数据库
func InitDB() *gorm.DB {
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		username,
		password,
		address,
		database,
		charset,
		parsetime,
		locale)
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}
	db.SingularTable(true)
	db.DB().SetMaxOpenConns(1)
	db.LogMode(true)
	return db
}
