package main

import (
	"testing"
)

func TestConvertTimestampToDate(t *testing.T) {
	var want string = "2021 Oct 09"
	if got := ConvertTimestampToDate("1633717104"); got != want {
		t.Errorf("ConvertTimestampToTime() = %q, want %q", got, want)
	}
}

func TestConvertTextToArticleBySpecialDivider(t *testing.T) {
	testStr := "TestName$TestTitle$Lorem$1633717104"
	var want Article = Article{
		Author: "TestName",
		Name:   "TestTitle",
		Body:   "Lorem",
		Time:   "2021 Oct 09",
	}
	if got := ConvertTextToArticleBySpecialDivider(testStr, "$"); got != want {
		t.Errorf("ConvertTextToArticleBySpecialDivider() = %q, want %q", got, want)
	}
}
