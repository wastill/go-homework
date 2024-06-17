package go_ethcli

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"log"
	"math/big"
)

func NewHdWallet() {
	// 助记词
	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	// 根据助记词生成钱包
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	//  "m/44'/60'/0'/0/0" 是什么意思？
	// 这个字符串代表了一个在分层确定性钱包（HD Wallets）中的钱包路径。具体来说，这个钱包路径通常用于以太坊钱包，并遵循了BIP-44标准（Bitcoin Improvement Proposal 44）。
	//
	//在这个路径中，每个数字代表了特定的目的。具体解释如下：
	//
	//m: 代表一条主链（master chain），作为路径的起始点。
	//44': 代表BIP-44规范中定义的Coin Type（币种类型）。44代表以太坊，后面的'表示这是强制的硬分叉（hardened derivation）。
	//60': 代表以太坊的币种代码。在这里60代表以太币（Ether）。
	//0': 代表第一个以太坊钱包账户。
	//0: 代表这个账户中的第一个外部（外部的意思是用于接收资金的）地址。
	//综合起来，这个路径指定了以太坊中第一个账户的第一个外部地址（通常用于接收资金）的HD Wallets路径。在以太坊开发和加密货币探索中，这种路径的概念非常常见和重要。
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 0xC49926C4124cEe1cbA0Ea94Ea31a6c12318df947
	// 生成另外一个账户
	path = hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/1")
	account, err = wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 0x8230645aC28A4EdD1b0B53E7Cd8019744E9dD559
}

func SignTransaction() {

	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, true)
	if err != nil {
		log.Fatal(err)
	}

	nonce := uint64(0)
	value := big.NewInt(1000000000000000000)
	toAddress := common.HexToAddress("0x0")
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(21000000000)
	var data []byte

	//tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data) 该方法已过时
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})
	// 签名交易
	signedTx, err := wallet.SignTx(account, tx, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 把交易打印出来
	spew.Dump(signedTx)
}
