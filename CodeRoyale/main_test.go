package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcCrossProduct(t *testing.T) {
	actual := calcCrossProduct(Point{1, 2}, PointF{2, 2}, Point{1, 1})
	expected := float64(-1)
	assert.Equal(t, expected, actual)

	actual2 := calcCrossProduct(Point{3, -2}, PointF{5, 1}, Point{1, 2})
	expected2 := float64(14)
	assert.Equal(t, expected2, actual2)
}

func TestCalcContact(t *testing.T) {
	actual1, actual2 := calcContact(PointF{3, 1}, PointF{0, 0}, float64(math.Sqrt(2)))
	expected1 := PointF{float64(1) / 5, float64(7) / 5}
	expected2 := PointF{float64(1), -1}

	ok := false
	if actual1.equal(expected1) && actual2.equal(expected2) {
		ok = true
	}
	if actual1.equal(expected2) && actual2.equal(expected1) {
		ok = true
	}
	if !ok {
		t.Errorf(`接点の計算結果が間違っています. 
		 expected: (%f %f), (%f, %f), actual: (%f %f), (%f %f)`,
			expected1.x, expected1.y, expected2.x, expected2.y,
			actual1.x, actual1.y, actual2.x, actual2.y)
	}
}

func TestCalcContact2(t *testing.T) {
	actual1, actual2 := calcContact(PointF{5, 4}, PointF{2, 3}, float64(math.Sqrt(2)))
	expected1 := PointF{float64(1)/5 + 2, float64(7)/5 + 3}
	expected2 := PointF{float64(1) + 2, -1 + 3}

	ok := false
	if actual1.equal(expected1) && actual2.equal(expected2) {
		ok = true
	}
	if actual1.equal(expected2) && actual2.equal(expected1) {
		ok = true
	}
	if !ok {
		t.Errorf(`接点の計算結果が間違っています. 
		 expected: (%f %f), (%f, %f), actual: (%f %f), (%f %f)`,
			expected1.x, expected1.y, expected2.x, expected2.y,
			actual1.x, actual1.y, actual2.x, actual2.y)
	}
}

func TestDist(t *testing.T) {
	actual := dist(Point{1, 2}, Point{4, 6})
	expected := 25
	assert.Equal(t, actual, expected)
}

func TestPow(t *testing.T) {
	actual := pow(2, 10)
	expected := 1024
	assert.Equal(t, actual, expected)
}
