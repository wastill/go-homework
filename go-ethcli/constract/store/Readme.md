# 编译.sol 文件

## 1. 安装solc

```shell
pip3 install solc-select==0.2.0

solc-select install 0.8.0

solc-select use 0.8.0
```

安装 abigen

```shell
go get -u github.com/ethereum/go-ethereum

git clone https://github.com/ethereum/go-ethereum.git

make

make devtools
```

## 2. 编译.sol 文件

```shell
solc --abi --bin --optimize --overwrite -o ./build/ ./contracts/Store.sol
```

## 3. 生成 go 代码

```shell
abigen --bin=./build/Store.bin --abi=./build/Store.abi --pkg=store --out=./store/store.go
```

