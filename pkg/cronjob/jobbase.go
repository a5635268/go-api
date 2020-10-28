package cronjob

import (
	"fmt"
	"go-api/tools"
	"sync"
	"time"

	"go-api/common/global"
	"go-api/pkg"
)

var timeFormat = "2006-01-02 15:04:05"
var retryCount = 3
var lock sync.Mutex

type JobCore struct {
	InvokeTarget   string
	Name           string
	JobId          uint
	EntryId        int
	CronExpression string
	Args           string
	JobType           int  // 1.http,2.exec
}

// 任务类型 http
type HttpJob struct {
	JobCore
}

type ExecJob struct {
	JobCore
}

func (e *ExecJob) Run() {
	startTime := time.Now()
	var obj = execList[e.InvokeTarget]
	if obj == nil {
		global.JobLogger.Warning(" ExecJob Run job nil", e)
		return
	}
	CallExec(obj.(JobsExec), e.Args)
	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)
	global.JobLogger.Info(time.Now().Format(timeFormat), " [INFO] JobCore ", e, "exec success , spend :", latencyTime)
}

//http 任务接口
func (h *HttpJob) Run() {

	startTime := time.Now()
	var count = 0
	/* 循环 */
LOOP:
	if count < retryCount {
		/* 跳过迭代 */
		str, err := pkg.Get(h.InvokeTarget)
		if err != nil {
			// 如果失败暂停一段时间重试
			fmt.Println(time.Now().Format(timeFormat), " [ERROR] mission failed! ", err)
			fmt.Printf(time.Now().Format(timeFormat)+" [INFO] Retry after the task fails %d seconds! %s \n", time.Duration(count)*time.Second, str)
			time.Sleep(time.Duration(count) * time.Second)
			goto LOOP
		}
		count = count + 1
		global.JobLogger.Info("response: ", str)
	}
	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)
	//TODO: 待完善部分

	global.JobLogger.Info(time.Now().Format(timeFormat), " [INFO] JobCore ", h, "exec success , spend :", latencyTime)
}

// 初始化
func Setup() {

	fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore Starting...")

	global.GADMCron = NewWithSeconds()

	if len(jobList) == 0 {
		fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore total:0")
	}

	var err error
	for i := 0; i < len(jobList); i++ {
		//
		if jobList[i].JobType == 1 {
			j := &HttpJob{}
			j.InvokeTarget = jobList[i].InvokeTarget
			j.CronExpression = jobList[i].CronExpression
			j.JobId = jobList[i].JobId
			j.Name = jobList[i].Name
			j.EntryId, err = AddJob(j)
			if err == nil {
				fmt.Println(time.Now().Format(timeFormat), " [INFO] " + jobList[i].Name + " 注册成功！ id:" + tools.IntToString(j.EntryId))
			}
		} else if jobList[i].JobType == 2 {
			j := &ExecJob{}
			j.InvokeTarget = jobList[i].InvokeTarget
			j.CronExpression = jobList[i].CronExpression
			j.JobId = jobList[i].JobId
			j.Name = jobList[i].Name
			j.Args = jobList[i].Args
			j.EntryId, err = AddJob(j)
			if err == nil {
				fmt.Println(time.Now().Format(timeFormat), " [INFO] " + jobList[i].Name + " 注册成功！ id:" + tools.IntToString(j.EntryId))
			}
		}
	}

	// 其中任务
	global.GADMCron.Start()
	fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore start success.")
	// 关闭任务
	defer global.GADMCron.Stop()
	// 空跑一个select，阻塞
	select {}
}

// 添加任务 AddJob(invokeTarget string, jobId int, jobName string, cronExpression string)
func AddJob(job Job) (int, error) {
	if job == nil {
		fmt.Println("unknown")
		return 0, nil
	}
	return job.addJob()
}

func (h *HttpJob) addJob() (int, error) {
	id, err := global.GADMCron.AddJob(h.CronExpression, h)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore AddJob error", err)
		return 0, err
	}
	EntryId := int(id)
	return EntryId, nil
}

func (h *ExecJob) addJob() (int, error) {
	id, err := global.GADMCron.AddJob(h.CronExpression, h)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore AddJob error", err)
		return 0, err
	}
	EntryId := int(id)
	return EntryId, nil
}
