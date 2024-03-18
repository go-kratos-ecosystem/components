package math

import "testing"

func TestMax(t *testing.T) {
	m, err := Max(1, 2, 3, 4, 5)
	if err != nil {
		t.Fatal(err)
	}
	if m != 5 {
		t.Fatalf("expected 5, got %d", m)
	}

	m2, err := Max(1.1, 2.2, 3.3, 4.4, 5.5)
	if err != nil {
		t.Fatal(err)
	}
	if m2 != 5.5 {
		t.Fatalf("expected 5.5, got %f", m2)
	}

	_, err = Max[int]()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestMin(t *testing.T) {
	m, err := Min(1, 2, 3, 4, 5)
	if err != nil {
		t.Fatal(err)
	}
	if m != 1 {
		t.Fatalf("expected 1, got %d", m)
	}

	m2, err := Min(1.2, 1.1, 2.2, 3.3, 4.4, 5.5)
	if err != nil {
		t.Fatal(err)
	}
	if m2 != 1.1 {
		t.Fatalf("expected 1.1, got %f", m2)
	}

	_, err = Min[int]()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestSum(t *testing.T) {
	s, err := Sum(1, 2, 3, 4, 5)
	if err != nil {
		t.Fatal(err)
	}
	if s != 15 {
		t.Fatalf("expected 15, got %d", s)
	}

	s2, err := Sum(1.1, 2.2, 3.3, 4.4, 5.5)
	if err != nil {
		t.Fatal(err)
	}
	if s2 != 16.5 {
		t.Fatalf("expected 16.5, got %f", s2)
	}

	_, err = Sum[int]()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
