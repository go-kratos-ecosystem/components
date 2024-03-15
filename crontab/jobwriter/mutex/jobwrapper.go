package mutex

import (
	"github.com/robfig/cron/v3"

	"github.com/go-kratos-ecosystem/components/v2/feature"
)

type Job interface {
	feature.Named

	IsMutex()
}

type ExpirableJob feature.Expirable

func SkipIfStillMutexLock(opts ...Option) cron.JobWrapper {
	o := newOptions(opts...)

	return func(job cron.Job) cron.Job {
		return cron.FuncJob(func() {
			j, ok := job.(Job)
			if !ok {
				job.Run()
				return
			}

			name := o.prefix + j.Name()
			expiration := o.expiration
			if ej, ok := job.(ExpirableJob); ok {
				expiration = ej.Expiration()
			}

			if err := o.locker.Lock(name, expiration); err != nil {
				o.logger.Info("crontab/jobwriter/locker: skip job %s, because still running", j.Name())
				return
			}
			defer func() {
				_ = o.locker.Unlock(name)
			}()

			job.Run()
		})
	}
}
