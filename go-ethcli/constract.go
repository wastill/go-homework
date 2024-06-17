package go_ethcli

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-homework/go-ethcli/constract/store"
	"log"
	"math/big"
	"strings"
)

const (
	LocalGanache    = "http://127.0.0.1:8545"
	LocalGanacheWS  = "ws://localhost:8545"
	LocalPrivateKey = "90a47f1aa3dcb54c5f157822f76742d188afdd86b418d08244feeebcfbfc0609"
	StoreAddress    = "0xB06D12947AF05CE47e59A8879B546f70d5e3157F"
)

func DeployStore() {
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

	// 部署合约 调用合约生成的 go 方法
	input := "1.0"
	address, tx, instance, err := store.DeployStore(auth, client, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   // 0xB06D12947AF05CE47e59A8879B546f70d5e3157F
	fmt.Println(tx.Hash().Hex()) // 0x1f733fe3838ecf95fe761711e6e18d342ea5ac9f9f23152fd4363fca8a82ede6

	_ = instance
}

func LoadStore() {
	client, err := ethclient.Dial(LocalGanache)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress(StoreAddress)
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("contract is loaded")

	fmt.Println(instance.Version(nil))
}

func WriteStore() {
	// 连接到本地的 ganache
	client, err := ethclient.Dial(LocalGanache)
	if err != nil {
		log.Fatal(err)
	}
	// 根据私钥获取公钥
	privateKey, err := crypto.HexToECDSA(LocalPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 根据公钥获取发送者地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

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

	// 连接到合约
	address := common.HexToAddress(StoreAddress)
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))
	// 调用合约的方法
	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s \n", tx.Hash().Hex()) // tx sent: 0x4a453932132beaeaf63acbb5a5f30556d24c1f66e34394e041562c90acb7edc7
	//
	result, err := instance.Items(nil, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result[:])) // bar
}

func StoreCode() {
	client, err := ethclient.Dial(LocalGanache)
	if err != nil {
		log.Fatal(err)
	}
	// 获取合约地址
	contractAddress := common.HexToAddress(StoreAddress)
	// 获取合约的 bytecode
	bytecode, err := client.CodeAt(context.Background(), contractAddress, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hex.EncodeToString(bytecode)) // 60806...10029
}
func StoreEvent() {
	client, err := ethclient.Dial(LocalGanacheWS)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(StoreAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log
			// {0xB06D12947AF05CE47e59A8879B546f70d5e3157F [0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4] [102 111 111 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 98 97 114 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] 7 0x72521165e1cf8f7da9d347c02d96bd0f2e22b09e86b2ecabcf87b381f1a39a3c 0 0x95719c9763ae2287f545ab4ee78b28a66f3c369fe2e62155eddc737f945bd4c5 0 false}
		}
	}
}

func PublishStoreEvent() {
	client, err := ethclient.Dial(LocalGanacheWS)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(StoreAddress)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(7),
		ToBlock:   big.NewInt(7),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(store.StoreABI)))
	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {
		fmt.Println(vLog.BlockHash.Hex()) // 0x3404b8c050aa0aacd0223e91b5c32fee6400f357764771d0684fa7b3f448f1a8
		fmt.Println(vLog.BlockNumber)     // 2394201
		fmt.Println(vLog.TxHash.Hex())    // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6

		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		results, err := contractAbi.Unpack("ItemSet", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(results)                // [102 111 111
		fmt.Println(string(event.Key[:]))   // foo
		fmt.Println(string(event.Value[:])) // bar

		var topics [4]string
		for i := range vLog.Topics {
			topics[i] = vLog.Topics[i].Hex()
		}

		fmt.Println(topics[0]) // 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4
	}

	eventSignature := []byte("ItemSet(bytes32,bytes32)")
	hash := crypto.Keccak256Hash(eventSignature)
	fmt.Println(hash.Hex()) // 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4
}
