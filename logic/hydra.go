package logic

import (
	"HydraServer/util/conf"
	"HydraServer/util/log"
	"HydraServer/util/net"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Reflexion struct {
	TaskId int    `json:"task_id"`
	Report []byte `json:"report"`
	Status bool   `json:"status"`
}

func WriteSlice2File(fname string, s []string) {
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Debug.Printf("用户手写目标文件%v\t产生的错误%v\n", fname, err)
	}
	defer f.Close()
	for _, v := range s {
		_, err := f.WriteString(strings.Join([]string{v, "\n"}, ""))
		if err != nil {
			log.Debug.Printf("用户err:%s\n", err)
		} else {
			log.Debug.Printf("用户手写:%s\n", v)
		}
	}
}
func execCommand(s []string) error {
	var bash string
	for _, v := range s {
		bash = strings.Join([]string{bash, v}, " ")
	}
	log.CMD.Printf("执行前确认命令:%v\n", bash)
	cmd := exec.Command("bash", "-c", bash)
	err := cmd.Start()
	if err != nil {
		log.Debug.Printf("执行hydra命令出错:%v\n", err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		log.Debug.Printf("执行hydra命令过程出错%v\n", err)
		return err
	}
	return nil
}
func Run(tid int, cLine []string) {
	defer func() {
		if err := recover(); err != nil {
			log.Debug.Printf("任务%v执行命令过程失败,准备发送空结构体,详情查阅日志\n", tid)
			//todo 发送空结构体
			r := &Reflexion{
				TaskId: tid,
				Report: nil,
				Status: false,
			}
			ip := strings.Join([]string{conf.GetVal("client", "host"), conf.GetVal("client", "port")}, ":")
			url := strings.Join([]string{ip, "hydra", "receive"}, "/")
			net.HttpPostJson(nil, r, url)
		}
	}()
	if err := execCommand(cLine); err != nil {
		log.Debug.Panicf("执行命令产生错误:%v\n", err)
	}
	tstring := strconv.Itoa(tid)
	reportName := strings.Join([]string{"report", tstring}, "/")
	log.Info.Printf("生成报告的路径:%v\n", reportName)
	//todo 读取为字节数组直接传回
	report, err := readReport(tid)
	if err != nil {
		log.Debug.Panicf("读取报告过程出错:%v\n", err)
	}
	r := &Reflexion{
		TaskId: tid,
		Report: report,
		Status: true,
	}
	ip := strings.Join([]string{conf.GetVal("client", "host"), conf.GetVal("client", "port")}, ":")
	url := strings.Join([]string{ip, "hydra", "receive"}, "/")
	log.Info.Printf("回传报告的url:%v\n", url)
	if json, err := net.HttpPostJson(nil, r, url); err != nil {
		log.Debug.Panicf("报告回传客户端出错:%v", err.Error())
	} else {
		log.Info.Printf("得到报告回传客户端之后的返回%v\n", json)
	}

}
func readReport(tid int) ([]byte, error) {
	tstring := strconv.Itoa(tid)
	ext := strings.Join([]string{tstring, "json"}, ".")
	path := strings.Join([]string{"report", ext}, "/")
	log.Info.Printf("任务%d即将要读取的报告路径:%s\n", tid, path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	log.Info.Printf("任务%d读取的报告内容:%s\n", tid, string(content))
	return content, nil
}
