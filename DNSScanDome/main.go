package main

import (
	"DnsScan/export"
	"DnsScan/lookup"
	"DnsScan/scan"
	"flag"
	"fmt"
	"os"
	"time"
)




func main() {
	var (
		flDomain = flag.String("t","","要爆破的域名")
		flWordlist = flag.String("w","","字典")
		flWorkCount = flag.Int("n",1000,"线程")
		flServerAddr = flag.String("s","","添加指定域名服务器，默认使用公共域名服务器114.114.114.114:53，8.8.8.8:53,202.101.172.35:53等")
		flServerAddrs = []string{
			"202.101.172.35:53",
			"114.114.114.114:53",
			"8.8.8.8:53"}
	)
	flag.Parse()			//格式化flag参数
	if *flDomain == "" || *flWordlist == "" {
		fmt.Println("域名或字典为空,命令格式不正确")
		flag.PrintDefaults()		//展示flag默认参数
		os.Exit(0)			//如果命令参数不正确即退出
	}
	if *flServerAddr != ""{		//判断用户是否要添加指定域名服务器
		flServerAddrs = append(flServerAddrs,fmt.Sprintf("%s:53",*flServerAddr))
	}

	var results []lookup.Result

	start := time.Now()				//计算程序运行时间，开始时间
	fmt.Println("程序开始运行......")
	for _,server := range flServerAddrs{
		results = append(results,scan.Scan(*flDomain,*flWordlist,server,*flWorkCount)...)
	}

	//整理结果，并生成CSV文件
	export.Export(results)

	end := time.Now()			//计算程序运行时间，结束时间
	delte := end.Sub(start)		//计算程序运行时间
	fmt.Printf("程序执行时长: %s\n", delte)	//输出运行时间
}
