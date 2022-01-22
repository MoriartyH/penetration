package export

import (
	"DnsScan/lookup"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	//"text/tabwriter"
)

//创建result目录
func ismkdir() bool {
	if err := os.Mkdir("result",os.ModeDir); err !=nil{
		if strings.Contains(err.Error(),"file already exists"){
			return true
		}
		fmt.Println("创建目录失败！", err)
		return false
	}
	return true
}

func Export(results []lookup.Result){
	fmt.Println("包括重复值的结果：",len(results))
	results = removeDuplicateElement(results)	//删除结果的重复值
	fmt.Println("去重后的结果：",len(results))

	filename := ""
	if ismkdir() {
		filename = fmt.Sprintf("result\\result%s.csv",time.Now().Format("200601021504"))
	}else {
		filename = fmt.Sprintf("result%s.csv",time.Now().Format("200601021504"))
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err!=nil{
		log.Println("文件打开失败！")
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，防止中文乱码
	WriterCsv:=csv.NewWriter(f)
	var s = []string{"域名：","IP："}
	err1:=WriterCsv.Write(s)
	if err1!=nil{
		log.Println("WriterCsv写入文件失败")
	}
	WriterCsv.Flush() //刷新，不刷新是无法写入的

	for _,r := range results{	//遍历results切片，并将结果输出到页面中
		s := []string{r.HostName,r.IPAddress}
		err1:=WriterCsv.Write(s)
		if err1!=nil{
			log.Println("WriterCsv写入文件失败")
			continue
		}
		WriterCsv.Flush() //刷新，不刷新是无法写入的
	}
	fmt.Printf("数据写入文件成功：%s\n",filename)
}


//删除重复值
func removeDuplicateElement(oldresult []lookup.Result)[]lookup.Result{
	newresult := make([]lookup.Result,0,len(oldresult))	//创建一个跟旧的切片一样大小的切片
	temp := map[lookup.Result]struct{}{}		//创建map字典
	for _,item := range oldresult{
		if _,ok := temp[item];!ok{	//如果字典中找不到元素，ok=false，!ok为true，就往切片中append元素。
			temp[item] = struct{}{}
			newresult = append(newresult,item)
		}
	}
	return newresult
}
