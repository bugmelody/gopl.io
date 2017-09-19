package snippet

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readStdin() {
	fmt.Println("请输入验证码:")
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	if err != nil {
		fmt.Println("读取标准输入出错")
		return
	}
	fmt.Println("你输入的验证码是: ", input)
}
