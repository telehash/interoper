package test

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	test, err := Load("testdata/example.md")
	if err != nil {
		t.Fatal(err)
	}

	if test.Name != "testdata/example.md" {
		t.Fatalf("expected test.Name = %q but was %q", "testdata/example.md", test.Name)
	}

	if test.Timeout() != time.Minute {
		t.Fatalf("expected test.Timeout() = %s but was %s", time.Minute, test.Timeout())
	}

	if test.Containers["worker"] == nil || test.Containers["worker"].Command != "th-test example worker" {
		t.Fatalf("expected test.Containers[\"worker\"] = %v but was %s", &Process{Command: "th-test example worker"}, test.Containers["worker"])
	}

	if test.Containers["driver"] == nil || test.Containers["driver"].Command != "th-test example driver" {
		t.Fatalf("expected test.Containers[\"driver\"] = %v but was %s", &Process{Command: "th-test example driver"}, test.Containers["driver"])
	}
}
