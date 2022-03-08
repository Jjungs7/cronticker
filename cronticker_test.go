package cronticker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const CRITERION_TIMESTAMP_20220308_123030_KST = 1646710230

var CRITERION_TIME = time.Unix(CRITERION_TIMESTAMP_20220308_123030_KST, 0)

func TestCreateNewCronTicker(t *testing.T) {
	var testCases = []struct {
		cronSpec string
	}{
		{"0 * * * *"},
		{"* * * * * *"},
		{"* * * * * ?"},
		{"0,30 * * * * ?"},
		{"0,30 */10 * * * ?"},
	}

	for _, tc := range testCases {
		_, err := NewCronWithOptionalSecondsTicker(tc.cronSpec)
		assert.NoError(t, err)
	}
}

func TestSpecIsNotCorrect(t *testing.T) {
	var testCases = []struct {
		cronSpec string
	}{
		{"* * * *"},
		{"* * * * * * *"},
	}

	for _, tc := range testCases {
		_, err := NewCronWithOptionalSecondsTicker(tc.cronSpec)
		assert.Error(t, err)
	}
}

func TestReturnsCorrectTick(t *testing.T) {
	var testCases = []struct {
		cronSpec   string
		assertTick time.Time
	}{
		{"* * * * *", CRITERION_TIME.Add(30 * time.Second)},
		{"*/30 * * * * *", CRITERION_TIME.Add(30 * time.Second)},
		{"30 * * * * *", CRITERION_TIME.Add(1 * time.Minute)},
		{"*/1,31 * * * * *", CRITERION_TIME.Add(1 * time.Second)},
		{"*/31,51 * * * * *", CRITERION_TIME.Add(1 * time.Second)},
		{"* * 7 * *", CRITERION_TIME.Add(24*30*time.Hour - 12*time.Hour - 30*time.Minute - 30*time.Second)},
	}

	for _, tc := range testCases {
		ticker, _ := NewCronWithOptionalSecondsTicker(tc.cronSpec)
		ticker.currentTick = time.Unix(CRITERION_TIMESTAMP_20220308_123030_KST, 0)

		ticker.Stop()
		ticker.runTimer()

		assert.Equal(t, tc.assertTick, ticker.nextTick)
	}
}
