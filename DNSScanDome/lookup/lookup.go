package lookup

import (
	"errors"
	"github.com/miekg/dns"
)

//结果结构体
type Result struct {
	IPAddress string	//IP地址
	HostName string		//域名
}

//查询域名
func Lookup(fqdn , serverAddr string)[]Result{
	var results []Result	//实例一个结果切片
	var cfqdn = fqdn	//创建一个fqdn副本，确保原始fqdn不会丢失
	for {	//查询域名的A记录和CNAME记录，并将结果保存至results结果切片中
		cname , err := lookupCNAME(cfqdn,serverAddr)
		if err == nil && len(cname) > 0 {
			cfqdn = cname[0]
			continue 				//未查询到CNAME，跳转处理下一个CNAME
		}
		ips , err := lookupA(cfqdn,serverAddr)
		if err != nil {
			break				//未查询到A记录
		}
		for _,ip := range ips{	//将查询到结果，以result结构保存在results切片中
			results = append(results,Result{IPAddress: ip,HostName: fqdn})
		}
		break
	}
	return results	//返回结果切片
}

//查询A记录
func lookupA(fqdn,serverAddr string)([]string,error){
	var m dns.Msg		//实例dns。Msg类型
	var ips []string	//创建一个结果切片
	m.SetQuestion(dns.Fqdn(fqdn),dns.TypeA)		//在应答函数中将fqdn加入，并指明查询A记录
	in , err := dns.Exchange(&m,serverAddr)		//将参数m和域名服务器serverAddr传入并查询
	if err != nil {
		return ips, err
	}
	if len(in.Answer) <1 {			//若查询应答记录<1的话，说明没有查询到A记录。则返回当前结果
		return ips, errors.New("没有查询到A记录")
	}
	for _,answer := range in.Answer {	//若查询到应答记录，则将结果保存在ips结果切片中
		if a ,ok :=answer.(*dns.A); ok{
			ips = append(ips,a.A.String())
		}
	}
	return ips, nil	//当查询结束后，返回结果
}

//查询CNAME记录
func lookupCNAME(fqdn,serverAddr string)([]string,error){
	var m dns.Msg
	var ips []string
	m.SetQuestion(dns.Fqdn(fqdn),dns.TypeCNAME) //在应答函数中将fqdn加入，并指明查询CNAME记录
	in , err := dns.Exchange(&m,serverAddr)
	if err != nil {
		return ips, err
	}
	if len(in.Answer) <1 {		//若查询应答记录<1的话，说明没有查询到CNAME记录。则返回当前结果
		return ips, errors.New("没有查询到CNAME记录")
	}
	for _,answer := range in.Answer {		//若查询到应答记录，则将结果保存在ips结果切片中
		if a ,ok :=answer.(*dns.CNAME); ok{
			ips = append(ips,a.Target)
		}
	}
	return ips, nil //当查询结束后，返回结果
}

