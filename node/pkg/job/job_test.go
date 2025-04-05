package job

import (
	"fmt"
	"github.com/robfig/cron/v3"
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
