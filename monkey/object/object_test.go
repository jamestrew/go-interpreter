package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	foobar1 := &String{Value: "foo bar baz"}
	foobar2 := &String{Value: "foo bar baz"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with the same content has different hash keys")
	}

	if foobar1.HashKey() != foobar2.HashKey() {
		t.Errorf("strings with the same content has different hash keys")
	}

	if hello1.HashKey() == foobar1.HashKey() {
		t.Errorf("strings with different content has same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	int1a := &Integer{Value: 69}
	int1b := &Integer{Value: 69}
	int2a := &Integer{Value: 420}
	int2b := &Integer{Value: 420}

	if int1a.HashKey() != int1b.HashKey() {
		t.Errorf("ints with the same content has different hash keys")
	}

	if int2a.HashKey() != int2b.HashKey() {
		t.Errorf("ints with the same content has different hash keys")
	}

	if int1a.HashKey() == int2a.HashKey() {
		t.Errorf("ints with different content has same hash keys")
	}
}

func TestBoolHashKey(t *testing.T) {
	bool1a := &Boolean{Value: true}
	bool1b := &Boolean{Value: true}
	bool2a := &Boolean{Value: false}
	bool2b := &Boolean{Value: false}

	if bool1a.HashKey() != bool1b.HashKey() {
		t.Errorf("bools with the same content has different hash keys")
	}

	if bool2a.HashKey() != bool2b.HashKey() {
		t.Errorf("bools with the same content has different hash keys")
	}

	if bool1a.HashKey() == bool2a.HashKey() {
		t.Errorf("bools with different content has same hash keys")
	}
}
