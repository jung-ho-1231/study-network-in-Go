package chapter08

import (
	"net/http"
	"testing"
	"time"
)

func TestHeadTime(t *testing.T) {
	resp, err := http.Head("https://time.gov/")

	if err != nil {
		t.Fatal(err)
	}

	_ = resp.Body.Close() // 예외 상황 처리 없이 항상 보디를 닫습니다.

	now := time.Now().Round(time.Second)
	date := resp.Header.Get("Date")

	if date == "" {
		t.Fatal("no date header received from time.gov")
	}

	dt, err := time.Parse(time.RFC1123, date)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("time.gov : %s (skew %s)", dt, now.Sub(dt))
}
