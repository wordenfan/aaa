package worden_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	baseURL := "https://hq.sinajs.cn/list=sh600036,sz002241"
	referer := "https://finance.sina.com.cn"

	// 创建请求
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	// 设置Referer头部
	req.Header.Set("Referer", referer)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	// 解析数据
	content := string(body)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			fields := strings.Split(line, ",")
			if len(fields) >= 31 { // 确保数据足够解析
				stockInfo := map[string]string{
					"股票名称":     fields[0],
					"今日开盘价":   fields[1],
					"昨日收盘价":   fields[2],
					"当前价格":     fields[3],
					"今日最高价":   fields[4],
					"今日最低价":   fields[5],
					"竞买价":       fields[6],
					"竞卖价":       fields[7],
					"成交量(股)":   fields[8],
					"成交金额(元)": fields[9],
					"买一量":       fields[10],
					"买一价":       fields[11],
					"买二量":       fields[12],
					"买二价":       fields[13],
					"日期":         fields[30],
					"时间":         fields[31],
				}
				// 打印股票信息
				fmt.Println("股票名称:", stockInfo["股票名称"])
				fmt.Println("今日开盘价:", stockInfo["今日开盘价"])
				// ... 其他信息打印省略，可以根据需要添加
				fmt.Println("日期:", stockInfo["日期"])
				fmt.Println("时间:", stockInfo["时间"])
				fmt.Println("-------------------------")
			}
		}
	}
}
