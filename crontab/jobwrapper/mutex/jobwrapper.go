package mutex

import (
	"fmt"

	"github.com/robfig/cron/v3"

	"github.com/go-kratos-ecosystem/components/v2/features"
)

type Job interface {
	cron.Job
	features.Named

	IsMutexJob()
}

type MutexJob struct{} //nolint:revive

func (m *MutexJob) IsMutexJob() {}

type ExpirableJob features.Expirable

func SkipIfStillMutexRunning(opts ...Option) cron.JobWrapper {
	o := newOptions(opts...)

	return func(job cron.Job) cron.Job {
		return cron.FuncJob(func() {
			j, ok := job.(Job)
			if !ok {
				o.logger.Info(fmt.Sprintf("crontab/jobwrapper/mutex: the job %v is not a mutex job, continue to run", job))
				job.Run()
				return
			}

			name := o.prefix + j.Name()
			expiration := o.expiration
			if ej, ok := job.(ExpirableJob); ok {
				expiration = ej.Expiration()
			}

			if err := o.locker.Lock(name, expiration); err != nil {
				o.logger.Info(fmt.Sprintf("crontab/jobwrapper/mutex: skip job [%s], because still mutex lock", j.Name()))
				return
			}
			defer o.locker.Unlock(name) //nolint:errcheck

			job.Run()
		})
	}
}
