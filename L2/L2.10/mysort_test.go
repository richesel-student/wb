package main

import (
	"reflect"
	"testing"
)

func TestSortString(t *testing.T) {
	lines := []string{"c", "a", "b"}
	opts := options{}

	sortLines(lines, opts)

	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("expected %v, got %v", expected, lines)
	}
}

func TestSortNumeric(t *testing.T) {
	lines := []string{"10", "2", "1"}
	opts := options{numeric: true}

	sortLines(lines, opts)

	expected := []string{"1", "2", "10"}
	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("expected %v, got %v", expected, lines)
	}
}

func TestSortReverse(t *testing.T) {
	lines := []string{"a", "b", "c"}
	opts := options{reverse: true}

	sortLines(lines, opts)

	expected := []string{"c", "b", "a"}
	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("expected %v, got %v", expected, lines)
	}
}

func TestSortByColumn(t *testing.T) {
	lines := []string{
		"b\t2",
		"a\t3",
		"c\t1",
	}
	opts := options{column: 2, numeric: true}

	sortLines(lines, opts)

	expected := []string{
		"c\t1",
		"b\t2",
		"a\t3",
	}
	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("expected %v, got %v", expected, lines)
	}
}

func TestUnique(t *testing.T) {
	lines := []string{"a", "a", "b", "b", "c"}
	result := unique(lines)

	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}