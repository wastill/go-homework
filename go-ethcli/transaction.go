package go_ethcli

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	//"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func QueryBlock() {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	// 获取最新区块
	number, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(number.Number.String())

	// 获取指定区块
	blockNumber := big.NewInt(number.Number.Int64())
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(block.Number().Uint64())
	log.Println(block.Time())
	log.Println(block.Difficulty())
	log.Println(block.Hash().Hex())
	log.Println(block.ParentHash().Hex())
	//log.Println(len(block.Transactions()))

	// 返回区块的交易数量
	count, err := client.TransactionCount(ctx, block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("transaction count ", count)
}

func QueryTransaction() {
	// 1. 根据区块编号查询交易信息
	// 2. 根据区块 hash 查询交易信息
	// 3. 根据交易 hash 查询交易信息
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}
	// 根据区块编号获取交易信息
	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
		fmt.Println(tx.Value().String())    // 10000000000000000
		fmt.Println(tx.Gas())               // 105000
		fmt.Println(tx.GasPrice().Uint64()) // 102000000000
		fmt.Println(tx.Nonce())             // 110644
		fmt.Println(tx.Data())              // []
		fmt.Println(tx.To().Hex())          // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e
		// 获取发送者信息
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		log.Println(chainID)
		// 通过 tx 获取发送者信息

		//msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
		list := tx.AccessList()
		fmt.Println(list)
		signer := types.NewEIP155Signer(chainID)
		from, err := types.Sender(signer, tx)
		if err != nil {
			//log.Fatal(err)
		}
		fmt.Println("=================") // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
		fmt.Println(from.Hex())          // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
		////
		//if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err == nil {
		//	fmt.Println(msg.From().Hex()) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
		//}
		// 交易收据
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(receipt.Status) // 1
		break
	}
	// 根据区块 hash 获取交易信息
	blockHash := common.HexToHash("0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9")
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}

	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(tx.Hash().Hex()) // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
		if idx > 3 {
			break
		}
	}
	// 根据交易 hash 获取交易信息
	txHash := common.HexToHash("0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx.Hash().Hex()) // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
	fmt.Println(isPending)       // false
}
func TransferToken() {
	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 发送者地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 获取发送者地址的 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 获取建议的 gas 价格
	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 获取接收者地址
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	tokenAddress := common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")

	// 调用方法的签名
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.NewKeccakState() // sh3 已经过时了
	hash.Write(transferFnSignature)
	// 获取方法的 ID
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	// 接收者地址
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	// 转账金额 代币
	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10) // 1000 tokens
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	// 构建交易数据 方法 id + 接收者地址 + 转账金额
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// 估算 gas 限制
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gasLimit) // 23256

	// 构建交易
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)
	// 签名交易
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 发送交易 交易发送成功后会返回交易哈希
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0xa56316b637a94c4cc0331c73ef26389d6c097506d581073f927275e7a6ece0bc
}

func PublishBlock() {
	client, err := ethclient.Dial("wss://ropsten.infura.io/ws")
	if err != nil {
		log.Fatal(err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f

			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(block.Hash().Hex())        // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			fmt.Println(block.Number().Uint64())   // 3477413
			fmt.Println(block.Time())              // 1529525947
			fmt.Println(block.Nonce())             // 130524141876765836
			fmt.Println(len(block.Transactions())) // 7
		}
	}
}
