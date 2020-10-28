package cronjob

var execList map[string]JobsExec
var jobList []JobCore

// 需要将定义的struct 添加到字典中；
// 字典 key 可以配置到 自动任务 调用目标 中；
func InitJob() {
	execList = map[string]JobsExec{
		"ExamplesOne": ExamplesOne{},
		// ...
	}
	jobList = []JobCore{
		{
			InvokeTarget:   "ExamplesOne",
			Name:           "exec test",
			JobId:          1,
			EntryId:        0,
			CronExpression: "*/5 * * * * *",
			Args:           "参数",
			JobType:           2,
		},
		{
			// 只支持http
			InvokeTarget:   "http://postman-echo.com/get",
			Name:           "http test",
			JobId:          2,
			EntryId:        0,
			CronExpression: "*/5 * * * * *",
			Args:           "参数",
			JobType:           1,
		},
	}
}
