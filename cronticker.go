/*
Cron 문법으로 Ticker를 실행할 수 있습니다
CronTicker의 문법은 다음과 같습니다
공식 문법과는 다르게, 초 단위를 Optional로 제공합니다

┌───────────── 초 (0 - 59) (Optional)
│ ┌───────────── 분 (0 - 59)
│ │ ┌───────────── 시 (0 - 23)
│ │ │ ┌───────────── 일 (1 - 31)
│ │ │ │ ┌───────────── 월 (1 - 12)
│ │ │ │ │ ┌───────────── 요일 (0 - 6) (일요일부터 토요일까지;
│ │ │ │ │ │                                   특정 시스템에서는 7도 일요일)
│ │ │ │ │ │
│ │ │ │ │ │
│ │ │ │ │ │
* * * * * *

위 문법 외에도 Descriptor 문법을 제공합니다
참고: https://pkg.go.dev/github.com/robfig/cron#hdr-Predefined_schedules

	Entry                  | Description                                | Equivalent To
	-----                  | -----------                                | -------------
	@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 1 1 *
	@monthly               | Run once a month, midnight, first of month | 0 0 1 * *
	@weekly                | Run once a week, midnight between Sat/Sun  | 0 0 * * 0
	@daily (or @midnight)  | Run once a day, midnight                   | 0 0 * * *
	@hourly                | Run once an hour, beginning of hour        | 0 * * * *

*/
package cronticker

import (
	"time"

	"github.com/robfig/cron/v3"
)

type CronTicker struct {
	C           chan time.Time
	k           chan bool
	currentTick time.Time
	nextTick    time.Time
	cron.Schedule
}

func NewCronWithOptionalSecondsTicker(spec string) (*CronTicker, error) {
	sch, err := cron.NewParser(cron.SecondOptional |
		cron.Minute |
		cron.Hour |
		cron.Dom |
		cron.Month |
		cron.Dow |
		cron.Descriptor).Parse(spec)
	if err != nil {
		return nil, err
	}

	return &CronTicker{
		C:        make(chan time.Time, 1),
		k:        make(chan bool, 1),
		Schedule: sch,
	}, nil
}

func (c *CronTicker) Start() {
	c.currentTick = time.Now()
	go c.runTimer()
}

func (c *CronTicker) Stop() {
	c.k <- true
}

func (c *CronTicker) runTimer() {
	c.nextTick = c.Schedule.Next(c.currentTick)
	timer := time.NewTimer(time.Until(c.nextTick))
	defer timer.Stop()

	for {
		select {
		case <-c.k:
			return
		case c.currentTick = <-timer.C:
			c.C <- c.currentTick
			c.nextTick = c.Schedule.Next(c.currentTick)
			timer.Reset(time.Until(c.nextTick))
		}
	}
}
