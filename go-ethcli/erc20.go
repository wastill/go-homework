package go_ethcli

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-homework/go-ethcli/build"
	token "go-homework/go-ethcli/constract/erc20"
	"log"
	"math/big"
	"strings"
)

const (
	Erc20Address = "0x85d4d57a07fe43ee1a1df3a640f425ccfbc04bd170df21900b6e17793d2ee0b7"
)

// LogTransfer ..
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

// LogApproval ..
type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}

// DeployErc20 部署合约
func DeployErc20() {
	client, err := ethclient.Dial(LocalGanache)
	if err != nil {
		log.Fatal(err)
	}
	// 根据私钥获取公钥
	privateKey, err := crypto.HexToECDSA(LocalPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 根据公钥拿到发送者地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 获取 gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 构建交易参数
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address, tx, instance, err := build.DeployErc20(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   // 0x07bEF311952B63c54Bd2f7A2edFe0721aa58a740
	fmt.Println(tx.Hash().Hex()) // 0x85d4d57a07fe43ee1a1df3a640f425ccfbc04bd170df21900b6e17793d2ee0b7

	_ = instance
}

// Erc20 监听合约事件
func Erc20() {
	client, err := ethclient.Dial(LocalGanache)
	if err != nil {
		log.Fatal(err)
	}

	// 0x Protocol (ZRX) token address
	contractAddress := common.HexToAddress(Erc20Address)
	// 构建区块查询参数
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(8),
		ToBlock:   big.NewInt(8),
		Addresses: []common.Address{
			contractAddress,
		},
	}
	// 查询区块上的日志
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	// 获取合约 abi
	contractAbi, err := abi.JSON(strings.NewReader(token.TokenMetaData.ABI))
	if err != nil {
		log.Fatal(err)
	}

	// 事件签名
	logTransferSig := []byte("Transfer(address,address,uint256)")
	LogApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Name: Transfer\n")

			var transferEvent LogTransfer

			ress, err := contractAbi.Unpack("Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(ress)

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("Tokens: %s\n", transferEvent.Tokens.String())

		case logApprovalSigHash.Hex():
			fmt.Printf("Log Name: Approval\n")

			var approvalEvent LogApproval

			ress, err := contractAbi.Unpack("Approval", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(ress)

			approvalEvent.TokenOwner = common.HexToAddress(vLog.Topics[1].Hex())
			approvalEvent.Spender = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf("Token Owner: %s\n", approvalEvent.TokenOwner.Hex())
			fmt.Printf("Spender: %s\n", approvalEvent.Spender.Hex())
			fmt.Printf("Tokens: %s\n", approvalEvent.Tokens.String())
		}

		fmt.Printf("\n\n")
	}
}
