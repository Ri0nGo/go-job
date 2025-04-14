package job

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"go-job/internal/dto"
	"go-job/internal/model"
	"log/slog"
	"os"
	"runtime"
	"testing"
	"time"
)

type JobTest struct {
	CronExpression string
	Name           string
	entryID        cron.EntryID
	cron           *cron.Cron
}

func NewJobTest(cronExpr, name string) (*JobTest, error) {
	return &JobTest{
		CronExpression: cronExpr,
		Name:           name,
	}, nil
}

func (j *JobTest) Run() {
	fmt.Printf("Running job name: %s, next exec time: %s\n",
		j.Name, j.cron.Entry(j.entryID).Next.Format(time.DateTime))
}

func TestJobRun(t *testing.T) {
	jobTest, err := NewJobTest("*/2 * * * * *", "test")
	if err != nil {
		t.Fatal(err)
	}
	c := cron.New(cron.WithSeconds())
	entryID, err := c.AddJob(jobTest.CronExpression, jobTest)
	if err != nil {
		t.Fatal(err)
	}
	jobTest.entryID = entryID
	jobTest.cron = c
	c.Start()
	time.Sleep(100 * time.Second)
}

func TestJobTest2(t *testing.T) {
	c := cron.New(cron.WithSeconds())
	entryID, err := c.AddFunc("*/2 * * * * *", func() {
		fmt.Println("hello world")
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(entryID)
	c.Start()
	time.Sleep(100 * time.Second)
}

// ============= 测试执行器 ============= //

type fileExecutor struct {
	Name string
}

func (executor *fileExecutor) Execute() (string, error) {
	fmt.Println("executing file executor")
	return "", nil
}

func (executor *fileExecutor) OnResultChange(f func(result model.JobExecResult)) {
}

func (executor *fileExecutor) Run() {
	fmt.Printf("[%s]\texecutor start, name: %s\n", time.Now().Format(time.DateTime), executor.Name)
	slog.Info("test file exec", "name", executor.Name)
}

func buildExecutor(name string) *fileExecutor {
	return &fileExecutor{
		Name: name,
	}
}

func generateTestCase() []dto.ReqJob {
	return []dto.ReqJob{
		{
			Id:       1,
			Name:     "test-01",
			CronExpr: "*/2 * * * * *",
		}, {
			Id:       2,
			Name:     "test-02",
			CronExpr: "*/10 * * * * *",
		}, {
			Id:       3,
			Name:     "test-03",
			CronExpr: "*/20 * * * * *",
		}, {
			Id:       4,
			Name:     "test-04",
			CronExpr: "*/30 * * * * *",
		}, {
			Id:       5,
			Name:     "test-05",
			CronExpr: "10 * * * * *",
		}, {
			Id:       6,
			Name:     "test-06",
			CronExpr: "20 * * * * *",
		},
	}
}

func TestAddJob(t *testing.T) {
	f := initLog("test-run.log")
	defer f.Close()

	var m1, m2 runtime.MemStats
	// 获取执行前内存情况
	runtime.ReadMemStats(&m1)

	handleCron()

	time.Sleep(96 * time.Minute)

	// 执行后再读取内存情况
	runtime.ReadMemStats(&m2)
	// 输出内存变化
	fmt.Printf("Memory Allocated: %d KB\n", (m2.Alloc-m1.Alloc)/1024)
	// Memory Allocated: 121 KB
}

func handleCron() {
	testCases := generateTestCase()
	for _, testCase := range testCases {
		ctx, cancel := context.WithCancel(context.Background())
		exec := buildExecutor(testCase.Name)

		jj := NewJob(ctx, cancel, testCase, exec)
		if err := jj.BuildCrontab(); err != nil {
			fmt.Println("build cron err:", err)
		}
		// 设置状态回调事件
		exec.OnResultChange(jj.OnResultChange)
		jj.Start()
		fmt.Println("add job, name: ", testCase.Name)
		AddJob(jj)
	}
}

// ============= utils ============= //
func initLog(name string) *os.File {
	// 创建或打开日志文件
	logFile, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("无法打开日志文件", "error", err)
		os.Exit(1)
	}

	// 创建文本格式的日志处理器，输出到文件
	handler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug, // 设置日志级别
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 特别处理时间字段
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					// 格式化为 "2006-01-02 15:04:05" 格式
					a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
				}
			}
			return a
		},
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logFile
}
