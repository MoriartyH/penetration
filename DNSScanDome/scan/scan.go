package scan

import (
	"DnsScan/lookup"
	"bufio"
	"fmt"
	"os"
)

//空结构体，用于跟踪线程
type empty struct {}

//扫描DNS
func Scan(flDomain,flWordlist,flServerAddr string,flWorkCount int)[]lookup.Result{
	var results []lookup.Result			//结果切片
	fqdns := make(chan string,flWorkCount)		//fqdn域名通道
	gather := make(chan []lookup.Result)		//接收结果的通道
	tracker := make(chan empty)			//判断是否关闭的通道

	fh , err := os.Open(flWordlist)	//打开指定文件
	if err != nil {
		panic(err)			//若打开错误，则panic返回该错误
	}
	defer fh.Close()		//在最后时关闭该文件
	scanner := bufio.NewScanner(fh)		//用bufio创建一个新的文件类型

	for i := 0; i < flWorkCount; i++ {//启动线程，将参数都传入线程中
		go worker(tracker, fqdns,gather, flServerAddr)
	}

	for scanner.Scan(){			//从文件类型遍历每行子域名构造成完整的域名，并传入fqdn域名通道中
		fqdns <- fmt.Sprintf("%s.%s",scanner.Text(),flDomain)
	}

	go func() {			//创建一个匿名函数可使用go并发执行
		for r := range gather{		//遍历gather通道中的结果，并添加到results切片中
			results = append(results,r...)
		}
		var e empty		//创建一个空的struct
		tracker <- e	//向判断通道中传入空的struct，为了防止出现竞态条件
	}()

	close(fqdns)	//所有结果已传入至results中，关闭fqdn域名通道
	for i:=0 ; i<flWorkCount;i++{
		<-tracker	//从判断通道中接收一次，以允许线程发出已完成工作并退出
	}
	close(gather)	//所有线程都退出后，即可关闭接收结果的通道
	<-tracker		//再接收一次判断通道，以使上面的匿名函数的go并发完全完成

	return results
}

//线程
func worker(tracker chan empty,fqdns chan string, gather chan []lookup.Result,serverAddr string){
	for fqdn := range fqdns {
		results := lookup.Lookup(fqdn, serverAddr)//将遍历出来的fqdn传入lookup函数查询
		if len(results) > 0 {               //判断如果results有结果的话，传入接收结果通道中
			gather <- results
		}
	}
	var e empty
	tracker <- e	//向判断通道中传入空的struct，为了防止出现竞态条件
}