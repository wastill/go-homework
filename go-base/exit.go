package go_base

import (
	"fmt"
	"os"
)

func Exit() {

	defer fmt.Println("!")

	os.Exit(3) // 执行该方法后，程序立刻结束，不会执行 defer
}
