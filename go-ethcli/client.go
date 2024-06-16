package go_ethcli

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	token "go-homework/go-ethcli/constract"
	"golang.org/x/crypto/sha3"
	"log"
	"math"
	"math/big"
	"os"
)

func NewClient() {
	Client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	id, err := Client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain id: %v", err)
	}
	log.Println(id)
	log.Println("Connected to the Ethereum client")
	address := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")

	fmt.Println(address.Hex()) // 0x71C7656EC7ab88b098defB751B7401B5f6d8976F
	//fmt.Println(address.Hash().Hex()) // 0x00000000000000000000000071c7656ec7ab88b098defb751b7401b5f6d8976f
	hasher := sha256.New()
	hasher.Write(address.Bytes())
	hash := common.BytesToHash(hasher.Sum(nil))

	// 打印哈希值
	println("Address Hash:", hash.Hex())
	fmt.Println(address.Bytes())
}

func Account() {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	account := common.HexToAddress("0xD3Cdf736d3F53f55070271CF4796b5121237dfBA")
	fmt.Println("Account:", account.Hex())
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance) // 25893180161173005034
	blockNumber := big.NewInt(0)
	balance1, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance1) // 25729324269165216042

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	fmt.Println(ethValue) // 25.729324269165216041

	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println(pendingBalance) // 25729324269165216042

}

func GolemErc20() {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}

	// Golem (GNT) Address
	tokenAddress := common.HexToAddress("0xa74476443119A942dE498590Fe1f2454d7D4aC0d")
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x0536806df512d6cdde913cf95c9886f65b1d3462")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name: %s\n", name)         // "name: Golem Network"
	fmt.Printf("symbol: %s\n", symbol)     // "symbol: GNT"
	fmt.Printf("decimals: %v\n", decimals) // "decimals: 18"

	fmt.Printf("wei: %s\n", bal) // "wei: 74605500647408739782407023"

	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

	fmt.Printf("balance: %f", value) // "balance: 74605500.647409"
}

func NewWallet() {
	//client, err := ethclient.Dial("http://
	// 生成私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	// 把私钥转化成字节
	privateKeyBytes := crypto.FromECDSA(privateKey)
	// 转成 16 进制，去除 0x
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:])
	// 生成公钥
	publicKey := privateKey.Public()
	// 转成公钥结构体
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 转成字节，去除通用前缀
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:])
	// 生成地址 公共地址其实就是公钥的Keccak-256哈希，然后我们取最后40个字符（20个字节）并用“0x”作为前缀
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))

}

func GenKeyStore() {

	store := keystore.NewKeyStore("./keystore/key1", keystore.StandardScryptN, keystore.StandardScryptP)
	pd := "123456"
	account, err := store.NewAccount(pd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account.Address.Hex())
}

func ImportKeyStore() {
	file := "./keystore/key1/UTC--2024-06-16T15-13-53.262838000Z--46de161b97dd849674e8b6031a38a128effc2b43"
	pd := "123456"
	readFile, err := os.ReadFile(file) // 读取文件 go 1.16 之后从 ioutil 换到了 os 包下
	if err != nil {
		log.Fatalf("Failed to read keystore file: %v", err)
	}
	// 直接读取
	k, err := keystore.DecryptKey(readFile, pd)
	if err != nil {
		log.Fatalf("Failed to decrypt key: %v", err)
	}
	fmt.Println(k.Address.Hex())

	// 生成到新文件
	ks := keystore.NewKeyStore("./keystore/key2", keystore.StandardScryptN, keystore.StandardScryptP)
	importECDSA, err := ks.ImportECDSA(k.PrivateKey, pd)
	if err != nil {
		log.Fatalf("Failed to import key: %v", err)
	}
	fmt.Println(importECDSA.Address.Hex())
}
