package go_ethcli

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-homework/go-ethcli/build"
	"log"
	"math/big"
	"strconv"
	"strings"
)

type LogFill struct {
	Maker                  common.Address
	Taker                  common.Address
	FeeRecipient           common.Address
	MakerToken             common.Address
	TakerToken             common.Address
	FilledMakerTokenAmount *big.Int
	FilledTakerTokenAmount *big.Int
	PaidMakerFee           *big.Int
	PaidTakerFee           *big.Int
	Tokens                 [32]byte
	OrderHash              [32]byte
}

type LogCancel struct {
	Maker                     common.Address
	FeeRecipient              common.Address
	MakerToken                common.Address
	TakerToken                common.Address
	CancelledMakerTokenAmount *big.Int
	CancelledTakerTokenAmount *big.Int
	Tokens                    [32]byte
	OrderHash                 [32]byte
}

type LogError struct {
	ErrorID   uint8
	OrderHash [32]byte
}

func DeployExchange() {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}

	// 0x Protocol Exchange smart contract address
	contractAddress := common.HexToAddress("0x12459C951127e0c374FF9105DdA097662A027093")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6383482),
		ToBlock:   big.NewInt(6383488),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(build.ExchangeMetaData.ABI))
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: keccak256("LogFill(address,address,address,address,address,uint256,uint256,uint256,uint256,bytes32,bytes32)")
	logFillEvent := common.HexToHash("0d0b9391970d9a25552f37d436d2aae2925e2bfe1b2a923754bada030c498cb3")

	// NOTE: keccak256("LogCancel(address,address,address,address,uint256,uint256,bytes32,bytes32)")
	logCancelEvent := common.HexToHash("67d66f160bc93d925d05dae1794c90d2d6d6688b29b84ff069398a9b04587131")

	// NOTE: keccak256("LogError(uint8,bytes32)")
	logErrorEvent := common.HexToHash("36d86c59e00bd73dc19ba3adfe068e4b64ac7e92be35546adeddf1b956a87e90")

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logFillEvent.Hex():
			fmt.Printf("Log Name: LogFill\n")

			var fillEvent LogFill

			ress, err := contractAbi.Unpack("LogFill", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(ress)
			fillEvent.Maker = common.HexToAddress(vLog.Topics[1].Hex())
			fillEvent.FeeRecipient = common.HexToAddress(vLog.Topics[2].Hex())
			fillEvent.Tokens = vLog.Topics[3]

			fmt.Printf("Maker: %s\n", fillEvent.Maker.Hex())
			fmt.Printf("Taker: %s\n", fillEvent.Taker.Hex())
			fmt.Printf("Fee Recipient: %s\n", fillEvent.FeeRecipient.Hex())
			fmt.Printf("Maker Token: %s\n", fillEvent.MakerToken.Hex())
			fmt.Printf("Taker Token: %s\n", fillEvent.TakerToken.Hex())
			fmt.Printf("Filled Maker Token Amount: %s\n", fillEvent.FilledMakerTokenAmount.String())
			fmt.Printf("Filled Taker Token Amount: %s\n", fillEvent.FilledTakerTokenAmount.String())
			fmt.Printf("Paid Maker Fee: %s\n", fillEvent.PaidMakerFee.String())
			fmt.Printf("Paid Taker Fee: %s\n", fillEvent.PaidTakerFee.String())
			fmt.Printf("Tokens: %s\n", hexutil.Encode(fillEvent.Tokens[:]))
			fmt.Printf("Order Hash: %s\n", hexutil.Encode(fillEvent.OrderHash[:]))

		case logCancelEvent.Hex():
			fmt.Printf("Log Name: LogCancel\n")

			var cancelEvent LogCancel

			ress, err := contractAbi.Unpack("LogCancel", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(ress)

			cancelEvent.Maker = common.HexToAddress(vLog.Topics[1].Hex())
			cancelEvent.FeeRecipient = common.HexToAddress(vLog.Topics[2].Hex())
			cancelEvent.Tokens = vLog.Topics[3]

			fmt.Printf("Maker: %s\n", cancelEvent.Maker.Hex())
			fmt.Printf("Fee Recipient: %s\n", cancelEvent.FeeRecipient.Hex())
			fmt.Printf("Maker Token: %s\n", cancelEvent.MakerToken.Hex())
			fmt.Printf("Taker Token: %s\n", cancelEvent.TakerToken.Hex())
			fmt.Printf("Cancelled Maker Token Amount: %s\n", cancelEvent.CancelledMakerTokenAmount.String())
			fmt.Printf("Cancelled Taker Token Amount: %s\n", cancelEvent.CancelledTakerTokenAmount.String())
			fmt.Printf("Tokens: %s\n", hexutil.Encode(cancelEvent.Tokens[:]))
			fmt.Printf("Order Hash: %s\n", hexutil.Encode(cancelEvent.OrderHash[:]))

		case logErrorEvent.Hex():
			fmt.Printf("Log Name: LogError\n")

			errorID, err := strconv.ParseInt(vLog.Topics[1].Hex(), 16, 64)
			if err != nil {
				log.Fatal(err)
			}

			errorEvent := &LogError{
				ErrorID:   uint8(errorID),
				OrderHash: vLog.Topics[2],
			}

			fmt.Printf("Error ID: %d\n", errorEvent.ErrorID)
			fmt.Printf("Order Hash: %s\n", hexutil.Encode(errorEvent.OrderHash[:]))
		}

		fmt.Printf("\n\n")
	}
}

// OUTPUT
// Log Block Number: 6383482
//Log Index: 35
//Log Name: LogFill
//[0xe269E891A2Ec8585a378882fFA531141205e92E9 0xD7732e3783b0047aa251928960063f863AD022D8 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2 240000000000000000000000 6930282000000000000 0 0 [48 106 154 126 203 217 68 101 89 162 198 80 180 207 193 109 31 182 21 170 43 63 79 99 7 141 166 208 33 38 132 64]]
//Maker: 0x8dd688660ec0BaBD0B8a2f2DE3232645F73cC5eb
//Taker: 0x0000000000000000000000000000000000000000
//Fee Recipient: 0xe269E891A2Ec8585a378882fFA531141205e92E9
//Maker Token: 0x0000000000000000000000000000000000000000
//Taker Token: 0x0000000000000000000000000000000000000000
//Filled Maker Token Amount: <nil>
//Filled Taker Token Amount: <nil>
//Paid Maker Fee: <nil>
//Paid Taker Fee: <nil>
//Tokens: 0xf08499c9e419ea8c08c4b991f88632593fb36baf4124c62758acb21898711088
//Order Hash: 0x0000000000000000000000000000000000000000000000000000000000000000
//
//
//Log Block Number: 6383482
//Log Index: 38
//Log Name: LogFill
//[0xe269E891A2Ec8585a378882fFA531141205e92E9 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2 0xD7732e3783b0047aa251928960063f863AD022D8 6941718000000000000 240000000000000000000000 0 0 [172 39 14 136 206 39 182 187 120 238 91 104 235 174 246 102 167 113 149 2 10 106 184 146 40 52 240 123 201 224 213 36]]
//Maker: 0x04aa059b2e31B5898fAB5aB24761e67E8a196AB8
//Taker: 0x0000000000000000000000000000000000000000
//Fee Recipient: 0xe269E891A2Ec8585a378882fFA531141205e92E9
//Maker Token: 0x0000000000000000000000000000000000000000
//Taker Token: 0x0000000000000000000000000000000000000000
//Filled Maker Token Amount: <nil>
//Filled Taker Token Amount: <nil>
//Paid Maker Fee: <nil>
//Paid Taker Fee: <nil>
//Tokens: 0x97ef123f2b566f36ab1e6f5d462a8079fbe34fa667b4eae67194b3f9cce60f2a
//Order Hash: 0x0000000000000000000000000000000000000000000000000000000000000000
//
//
//Log Block Number: 6383488
//Log Index: 43
//Log Name: LogCancel
//[0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2 0x89d24A6b4CcB1B6fAA2625fE562bDD9a23260359 30000000000000000000 7274848425000000000000 [228 62 255 56 220 39 175 4 107 251 212 49 146 105 38 192 114 187 199 165 9 213 111 111 26 122 225 245 173 126 254 79]]
//Maker: 0x0004E79C978B95974dCa16F56B516bE0c50CC652
//Fee Recipient: 0xA258b39954ceF5cB142fd567A46cDdB31a670124
//Maker Token: 0x0000000000000000000000000000000000000000
//Taker Token: 0x0000000000000000000000000000000000000000
//Cancelled Maker Token Amount: <nil>
//Cancelled Taker Token Amount: <nil>
//Tokens: 0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391
//Order Hash: 0x0000000000000000000000000000000000000000000000000000000000000000
