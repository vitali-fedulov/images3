package images

import (
	"reflect"
	"testing"
)

func TestHyperPoints10(t *testing.T) {
	s := HyperPoints10
	if len(s) != 10 {
		t.Errorf(
			"Mismatching set length. Expected 10, got %v",
			len(s))
	}
	for _, p1 := range s {
		for _, p2 := range s {
			if p1 == p2 {
				continue
			}
			d := distance(p1, p2)
			if d < 2 {
				t.Errorf(
					`Too small distance between points, got
					distance %v, between points %v and %v`,
					d, p1, p2)
			}
		}
	}
}

func TestSelectPoints(t *testing.T) {
	got := CustomPoints(5)
	want := map[Point]bool{
		{2, 2}: true, {2, 8}: true, {5, 5}: true,
		{8, 2}: true, {8, 8}: true}
	for k := range got {
		if _, ok := want[k]; !ok {
			t.Errorf(
				"Missing point %v in expected %v", k, want)
		}
	}
	got = CustomPoints(12)
	if len(got) != 12 {
		t.Errorf(
			"Mismatching set length. Expected 12, got %v",
			len(got))
	}
	if _, ok := got[Point{1, 1}]; !ok {
		t.Errorf(
			"Missing point{1, 1} in the result %v", got)
	}
	for p1 := range got {
		for p2 := range got {
			if p1 == p2 {
				continue
			}
			d := distance(p1, p2)
			if d < 2 {
				t.Errorf(
					`Too small distance between points, got
					distance %v, between points %v and %v`,
					d, p1, p2)
			}
		}
	}
}

func TestDistance(t *testing.T) {
	got := distance(
		Point{5, 7}, Point{2, 8})
	want := 3.1622776601683795
	if got != want {
		t.Errorf("Want %v, got %v", want, got)
	}
}

func TestMinKey(t *testing.T) {
	got := minKey(
		map[Point]float64{
			{1, 1}: 1.9, {2, 2}: 0.3, {3, 3}: 0.01,
			{7, 7}: 12.0, {9, 9}: 3.0})
	want := Point{3, 3}
	if got != want {
		t.Errorf("Want %v, got %v", want, got)
	}
}

func TestMaxKey(t *testing.T) {
	got := maxKey(
		map[Point]float64{
			{1, 1}: 1.9, {2, 2}: 0.3, {3, 3}: 0.01,
			{7, 7}: 12.0, {9, 9}: 3.0})
	want := Point{7, 7}
	if got != want {
		t.Errorf("Want %v, got %v", want, got)
	}
}

func TestExlude(t *testing.T) {
	got := exclude(
		Point{2, 2}, map[Point]bool{
			{1, 1}: true, {2, 2}: true, {3, 3}: true,
			{7, 7}: true, {9, 9}: true})
	want := map[Point]bool{
		{1, 1}: true, {3, 3}: true,
		{7, 7}: true, {9, 9}: true}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Want %v, got %v", want, got)
	}
}

func TestNearest(t *testing.T) {
	got := nearest(
		map[Point]bool{
			{1, 1}: true, {2, 2}: true, {3, 3}: true,
			{7, 7}: true, {9, 9}: true},
		Point{6, 6})
	want := Point{7, 7}
	if got != want {
		t.Errorf("Want %v, got %v", want, got)
	}
	got = nearest(
		map[Point]bool{
			{1, 1}: true, {2, 2}: true, {3, 3}: true,
			{7, 7}: true, {9, 9}: true},
		Point{3, 3})
	want = Point{2, 2}
	if got != want {
		t.Errorf("Want %v, got %v", want, got)
	}
}
