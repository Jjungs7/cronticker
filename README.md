# CronTicker

Ticker, but with cron specification

## Download

```bash
go get github.com/Jjungs7/CronTicker
```

## Example

```golang
ticker, _ := NewCronWithOptionalSecondsTicker("*/5 * * * *")
ticker.Start()

for {
    <-ticker.C

    // do job
}
```
