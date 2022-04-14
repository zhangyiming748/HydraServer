package controller

import (
	"HydraServer/logic"
	"HydraServer/util/log"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

type HydraArgs struct {
	TaskId       int    `json:"task_id"`
	TaskName     string `json:"task_name"`
	Address      string `json:"address"`
	Port         string `json:"port"`
	Protocol     string `json:"protocol"`
	Username     string `json:"username"`
	UsernameFile string `json:"username_file"`
	UserNameType int    `json:"user_name_type"` // 1 默认 2 手写 3 上传
	Password     string `json:"password"`
	PasswordFile string `json:"password_file"`
	PasswordType int    `json:"passwd_type"` // 1 默认 2 手写 3 上传
	UserId       int    `json:"user_id"`
	Path         string `json:"path"`
	Form         string `json:"form"`
	Sid          string `json:"sid"`
	RequestHost  string `json:"request_host"`
}
type Transport struct {
	TaskId       int      `json:"task_id"`
	TaskString   string   `json:"task_string"`
	Username     string   `json:"username"`
	UsernameType int      `json:"username_type"`
	Password     string   `json:"password"`
	PasswordType int      `json:"password_type"`
	CmdLine      []string `json:"cmd_line"`
}

const PREFIX = "/Users/zen/Github/HydraServer" // $PWD
func Hydra(ctx *gin.Context) {
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	log.Info.Printf("ctx.Request.body: %v", string(data))
	var t Transport
	if err := json.Unmarshal(data, &t); err != nil {
		log.Debug.Println("解析为结构体失败")
		ctx.JSON(400, Transport{})
	}
	if t.UsernameType == 2 {
		path := strings.Join([]string{"username", t.TaskString}, "/")
		log.Info.Printf("接收手写的用户名字典保存到%v\n", path)
		var items []string
		if strings.Contains(t.Username, "|") {
			items = strings.Split(t.Username, "|")
		} else {
			items = append(items, t.Username)
		}
		logic.WriteSlice2File(path, items)
	}
	if t.PasswordType == 2 {
		path := strings.Join([]string{"password", t.TaskString}, "/")
		log.Info.Printf("接收手写的密码字典保存到%v\n", path)
		var items []string
		if strings.Contains(t.Password, "|") {
			items = strings.Split(t.Password, "|")
		} else {
			items = append(items, t.Password)
		}
		logic.WriteSlice2File(path, items)
	}
	log.Info.Printf("任务id = %v\n", t.TaskId)
	log.Info.Printf("任务id的字符串 = %v\n", t.TaskString)
	log.Info.Printf("手写用户名字典 = %v\n", t.Username)
	log.Info.Printf("用户名字典类型 = %v\n", t.UsernameType)
	log.Info.Printf("手写密码字典 = %v\n", t.Password)
	log.Info.Printf("密码字典类型 = %v\n", t.PasswordType)
	log.Info.Printf("拼写好的命令行 = %v\n", t.CmdLine)
	logic.Run(t.TaskId, t.CmdLine)
	ctx.JSON(200, data)
}

func SaveFile(ctx *gin.Context) {
	t_type := ctx.Request.FormValue("type")
	tid := ctx.Request.FormValue("task_id")

	ctx.JSON(200, map[string]string{"task_id": tid, "type": t_type})
}
func SaveUsername(ctx *gin.Context) {
	tid := ctx.Request.FormValue("task_id")
	_, header, err := ctx.Request.FormFile("username_file")
	if err != nil {
		log.Debug.Println("上传文件失败")
	}

	dst := strings.Join([]string{PREFIX, "username", tid}, "/")
	ctx.SaveUploadedFile(header, dst)
	ctx.JSON(200, map[string]string{"tid": tid, "save": dst})
}
func SavePassword(ctx *gin.Context) {
	tid := ctx.Request.FormValue("task_id")
	_, header, err := ctx.Request.FormFile("password_file")
	if err != nil {
		log.Debug.Println("上传文件失败")
	}

	dst := strings.Join([]string{PREFIX, "password", tid}, "/")
	ctx.SaveUploadedFile(header, dst)
	ctx.JSON(200, map[string]string{"tid": tid, "save": dst})
}
