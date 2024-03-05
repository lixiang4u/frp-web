package utils

import (
	"bufio"
	"os"
)

func WaitInput(tips ...string) string {
	if len(tips) == 0 {
		tips = append(tips, "按回车结束输入")
	}
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	return str
}

func WaitInputExit(tips ...string) {
	WaitInput(tips...)
	os.Exit(1)
}
