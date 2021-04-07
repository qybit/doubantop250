package main

import (
	"bufio"
	"fmt"
	"os"
	"qybit.com/doubantop250/spider"
	"strconv"
)

func main() {
	fmt.Print("请输入要爬取的页数，最大10页：")
	cin := bufio.NewScanner(os.Stdin)
	cin.Scan()
	page, err := strconv.Atoi(cin.Text())
	if err != nil {
		fmt.Println("输入数据不合法，请按照要求输入！")
		os.Exit(1)
	}
	fmt.Print("是否需要持久化？（请输入 yes/y 或者 no/n）")
	cin.Scan()
	ok := cin.Text()
	isOk := false
	if ok == "yes" || ok == "y" {
		isOk = true
	}
	spider.Run(page, isOk)
}
