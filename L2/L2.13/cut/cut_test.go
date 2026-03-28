package main

import "testing"

func TestPF1(t *testing.T) {
	fields := ParseFields("2")

	if !fields[2] {
		t.Fatal("expected field 2")
	}
}

func TestPFRng(t *testing.T) {
	fields := ParseFields("2-4")

	if !fields[2] || !fields[3] || !fields[4] {
		t.Fatal("range parsing failed")
	}
}

func TestPFMix(t *testing.T) {
	fields := ParseFields("1,3-5")

	if !fields[1] || !fields[3] || !fields[4] || !fields[5] {
		t.Fatal("mixed parsing failed")
	}
}

func TestPLBase(t *testing.T) {

	fields := map[int]bool{
		1: true,
		3: true,
	}

	res, ok := ProcessLine("a:b:c:d", ":", fields, false)

	if !ok {
		t.Fatal("expected line")
	}

	if res != "a:c" {
		t.Fatalf("expected a:c got %s", res)
	}
}

func TestPLComma(t *testing.T) {

	fields := map[int]bool{
		2: true,
		4: true,
	}

	res, ok := ProcessLine("a,b,c,d", ",", fields, false)

	if !ok || res != "b,d" {
		t.Fatalf("expected b,d got %s", res)
	}
}

func TestPLSep(t *testing.T) {

	fields := map[int]bool{1: true}

	_, ok := ProcessLine("hello", ":", fields, true)

	if ok {
		t.Fatal("should skip")
	}
}

func TestPLNoSep(t *testing.T) {

	fields := map[int]bool{1: true}

	res, ok := ProcessLine("hello", ":", fields, false)

	if !ok || res != "hello" {
		t.Fatalf("expected hello got %s", res)
	}
}

func TestPLEmpty(t *testing.T) {

	fields := map[int]bool{3: true}

	res, ok := ProcessLine("1:Alice::London", ":", fields, false)

	if !ok || res != "" {
		t.Fatalf("expected empty got %s", res)
	}
}

func TestPLRange(t *testing.T) {

	fields := map[int]bool{5: true}

	_, ok := ProcessLine("a:b:c", ":", fields, false)

	if ok {
		t.Fatal("expected no output")
	}
}

func TestPLLead(t *testing.T) {

	fields := map[int]bool{
		1: true,
		2: true,
	}

	res, ok := ProcessLine(":a:b", ":", fields, false)

	if !ok || res != ":a" {
		t.Fatalf("expected :a got %s", res)
	}
}

func TestPLTrail(t *testing.T) {

	fields := map[int]bool{
		1: true,
		3: true,
	}

	res, ok := ProcessLine("a:b:", ":", fields, false)

	if !ok || res != "a:" {
		t.Fatalf("expected a: got %s", res)
	}
}