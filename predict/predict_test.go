package predict

import (
	"fmt"
	"log"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
}

func TestPredict(t *testing.T) {
	predictor, err := NewPredictor("test_fixtures/model/saved_model_145.pb")
	if err != nil {
		log.Printf("NewPredictor: %v", err)
		t.Fail()
	}

	model := predictor.Model
	defer model.Session.Close()
	cc, err := predictor.Predict("test_fixtures/images/16_20180430T0500Z.jpg")
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fail()
	}
	assertEqual(t, 0, cc) // This zero does not come from the labels, but from what a python implementation we know predict works.
	fmt.Printf("CC: %d\n", cc)

	cc, err = predictor.Predict("test_fixtures/images/176_20190516T1700Z.jpg")
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fail()
	}
	assertEqual(t, 7, cc) // This 7 does not come from the labels, but from what a python implementation we know predict works .
	fmt.Printf("CC: %d\n", cc)

	cc, err = predictor.Predict("test_fixtures/images/114_20180922T0900Z.jpg")
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fail()
	}
	assertEqual(t, 5, cc) // This 5 does not come from the labels, but from what a python implementation we know predict works .
	fmt.Printf("CC: %d\n", cc)

}
