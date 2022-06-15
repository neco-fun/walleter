package model

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"neco-wallet-center/internal/comm"
)

// Wallet 钱包
type Wallet struct {
	gorm.Model       `swagger-ignore:"true"`
	AccountId        string           `json:"account_id" gorm:"unique; not null"` // 玩家账户ID
	PublicAddress    string           `json:"address" gorm:"unique;not null"`     // 玩家的钱包地址
	ERC20TokenData   []ERC20TokenData `json:"erc_20_token_data" gorm:"foreignKey:AccountId;references:AccountId"`
	ERC1155TokenData ERC1155TokenData `json:"erc_1155_token_data" gorm:"foreignKey:AccountId;references:AccountId"`
	CheckSign        string           `json:"check_sign" gorm:"type:varchar(128);not null;comment:'安全签名'"`
}

type ERC20TokenData struct {
	gorm.Model    `swagger-ignore:"true"`
	AccountId     string  `json:"account_id"`     //往家账户的ID
	TokenType     string  `json:"token_type"`     //代币类型 NFISH, BUSD
	TokenBalance  float64 `json:"token_balance"`  // 玩家当前代币的余额
	TokenIncome   float64 `json:"token_income"`   // 玩家通过玩游戏的总收入
	TokenSpend    float64 `json:"token_spend"`    // 玩家通过玩游戏的总花费
	TokenDeposit  float64 `json:"token_deposit"`  // 玩家通过质押的总额度
	TokenWithdraw float64 `json:"token_withdraw"` // 玩家提取代币的总金额
	TokenFee      float64 `json:"token_fee"`      // 玩家使用当前代币付的总手续费
}

type ERC1155TokenData struct {
	gorm.Model `swagger-ignore:"true"`
	AccountId  string `json:"account_id"` //往家账户的ID
	Ids        string `json:"ids"`        //玩家拥有的NFT所有的id
	Values     string `json:"values"`     // 玩家拥有的NFT的数量
}

// ERC20WalletLog 钱包流水日志
type ERC20WalletLog struct {
	gorm.Model     `swagger-ignore:"true"`
	AccountId      string  `json:"account_id"` //往家账户的ID
	PublicAddress  string  `json:"address" gorm:"type:varchar(256);uniqueIndex;not null;comment:'钱包地址'"`
	BusinessModule string  `json:"business_module" gorm:"type:varchar(64);not null;comment:'业务模块'"`
	ActionType     string  `json:"action_type" gorm:"type:varchar(64);not null;comment:'操作类型'"`
	TokenType      string  `json:"token_type;comment:'变更的代币数据'"`
	Value          float64 `json:"value"` // 代币金额
	Fee            float64 `json:"fee"`   // 手续费
	Status         string  `json:"status" gorm:"type:varchar(64);not null;comment:处理状态"`
	OriginalWallet Wallet  `json:"original_wallet" gorm:"type:json;not null;comment:'变更前的钱包数据'"`
	DisposeWallet  Wallet  `json:"dispose_wallet" gorm:"type:json;not null;comment:'变更后的钱包数据'"`
}

// ERC1155WalletLog 钱包流水日志
type ERC1155WalletLog struct {
	gorm.Model     `swagger-ignore:"true"`
	AccountId      string  `json:"account_id"` //往家账户的ID
	PublicAddress  string  `json:"address" gorm:"type:varchar(256);uniqueIndex;not null;comment:'钱包地址'"`
	BusinessModule string  `json:"business_module" gorm:"type:varchar(64);not null;comment:'业务模块'"`
	ActionType     string  `json:"action_type" gorm:"type:varchar(64);not null;comment:'操作类型'"`
	Ids            string  `json:"ids;comment:'变更的NFT IDs'"`
	Values         float64 `json:"values"` // 变更的NFT数量
	Fee            float64 `json:"fee"`    // 手续费
	Status         string  `json:"status" gorm:"type:varchar(64);not null;comment:处理状态"`
	OriginalWallet Wallet  `json:"original_wallet" gorm:"type:json;not null;comment:'变更前的钱包数据'"`
	DisposeWallet  Wallet  `json:"dispose_wallet" gorm:"type:json;not null;comment:'变更后的钱包数据'"`
}

type walletDA0 struct{}

var WalletDAO = &walletDA0{}

func (dao *walletDA0) getWalletTableName(gameClient comm.GameClient) string {
	return fmt.Sprintf("t_wallet_%d", gameClient)
}

func (dao walletDA0) InitWallet(ctx context.Context, gameClient comm.GameClient, accountId string, publicAddress string) bool {
	err := getDb(ctx).Transaction(
		func(tx1 *gorm.DB) error {
			//首先创建data
			var erc20Data = ERC20TokenData{
				AccountId:     accountId,
				TokenType:     comm.NFISH.String(),
				TokenBalance:  0,
				TokenIncome:   0,
				TokenSpend:    0,
				TokenDeposit:  0,
				TokenWithdraw: 0,
				TokenFee:      0,
			}
			if err := tx1.Create(&erc20Data).Error; err != nil {
				return err
			}

			var erc1155Data = ERC1155TokenData{
				AccountId: accountId,
				Ids:       "",
				Values:    "",
			}
			if err := tx1.Create(&erc1155Data).Error; err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return false
	}
	return true
}

type erc20TokenDataDAO struct{}

var ERC20TokenDataDAO = &erc20TokenDataDAO{}

func (dao *erc20TokenDataDAO) getERC20TokenDataTableName(gameClient comm.GameClient) string {
	return fmt.Sprintf("t_erc20_token_data_%d", gameClient)
}

type erc1155TokenDataDAO struct{}

var ERC1155TokenDataDAO = &erc1155TokenDataDAO{}

func (dao *erc1155TokenDataDAO) getERC1155TokenDataTableName(gameClient comm.GameClient) string {
	return fmt.Sprintf("t_erc1155_token_data_%d", gameClient)
}

type erc20WalletLogDAO struct{}

var ERC20WalletLogDAO = &erc20WalletLogDAO{}

func (dao *erc20WalletLogDAO) getERC20WalletLogTableName(gameClient comm.GameClient) string {
	return fmt.Sprintf("t_erc20_wallet_log_%d", gameClient)
}

type erc1155WalletLogDAO struct{}

var ERC1155WalletLogDAO = &erc1155WalletLogDAO{}

func (dao *erc1155WalletLogDAO) getERC1155WalletLogTableName(gameClient comm.GameClient) string {
	return fmt.Sprintf("t_erc1155_wallet_log_%d", gameClient)
}
