package config

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type CronJob struct {
	JobName string
	JobFunc func() error
	JobSpec string
	EntryID cron.EntryID
}

func NewCronJob(jobName string, jobFunc func() error, jobSpec string) *CronJob {
	return &CronJob{JobName: jobName, JobFunc: jobFunc, JobSpec: jobSpec}
}

func (cr *CronJob) Run(cronJob *cron.Cron) error {
	var err error
	cr.EntryID, err = cronJob.AddFunc(cr.JobSpec, func() {
		err := cr.JobFunc()
		if err != nil {
			log.Fatalf("Error when do job %s: %v", cr.JobName, err)
		}
	})
	return err
}

type CronJobs struct {
	CronJobPackage *cron.Cron
	Jobs           map[string]*CronJob
}

func (crs CronJobs) Run() {
	for _, cr := range crs.Jobs {
		err := cr.Run(crs.CronJobPackage)
		if err != nil {
			log.Fatalf("Error when add job %s: %v", cr.JobName, err)
		}
	}
	crs.CronJobPackage.Start()
}

func (crs CronJobs) Stop() {
	crs.CronJobPackage.Stop()
}

func (crs *CronJobs) Remove(key string) {
	job, found := crs.Jobs[key]
	if !found {
		log.Printf("Job %s not found\n", key)
		return
	}
	if job.EntryID != 0 {
		crs.CronJobPackage.Remove(job.EntryID)
	}
	delete(crs.Jobs, key)
}

func (crs *CronJobs) Add(key string, job *CronJob) {
	crs.Jobs[key] = job
}

func NewCronJobs(timezone string) (*CronJobs, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("error loading location %s: %v", timezone, err)
	}

	return &CronJobs{
		CronJobPackage: cron.New(
			cron.WithLocation(location),
			cron.WithChain(
				cron.SkipIfStillRunning(cron.DefaultLogger),
				cron.Recover(cron.DefaultLogger),
			),
		),
		Jobs: make(map[string]*CronJob),
	}, nil
}
