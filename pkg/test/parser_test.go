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

	if test.SystemUnderTest == nil || test.SystemUnderTest.Command != "test-net-link await" {
		t.Fatalf("expected test.SystemUnderTest = %v but was %s", &Process{Command: "test-net-link await"}, test.SystemUnderTest)
	}

	if test.TestDriver == nil || test.TestDriver.Command != "test-net-link establish" {
		t.Fatalf("expected test.TestDriver = %v but was %s", &Process{Command: "test-net-link await"}, test.TestDriver)
	}
}
