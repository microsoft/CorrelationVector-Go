// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package correlationvector contains library functions to manipulate CorrelationVectors.
package correlationvector

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSpinSortValidation(t *testing.T) {
	vector := NewCorrelationVector()
	spinParameters := SpinParameters{FineInterval, ShortPeriodicity, TwoEntropy}

	lastSpinValue := uint64(0)
	wrappedCounter := 0
	for i := 0; i < 100; i++ {
		spin, _ := SpinWithParameters(vector.Value(), &spinParameters)

		// The cV after a Spin will look like <cvBase>.0.<spinValue>.0, so the spinValue is at index = 2.
		spinValue, _ := strconv.ParseUint(strings.Split(spin.Value(), ".")[2], 10, 64)

		// Count the number of times the counter wraps.
		if spinValue <= lastSpinValue {
			wrappedCounter++
		}

		lastSpinValue = spinValue

		time.Sleep(10 * time.Millisecond)
	}

	if wrappedCounter > 1 {
		t.Errorf("Expecting the extension to wrap at most 1 time, actually wrapped %d times", wrappedCounter)
	}
}

func TestSpinOverMaxCVLength(t *testing.T) {

	var baseVector = "tul4NUsfs9Cl7mOf.2147483647.2147483647.2147483647.214748364.23"
	cv, _ := Spin(baseVector)
	if cv.Value() != (baseVector + CVTerminator) {
		t.Errorf("Termination should be applied for CV that goes beyond max length after spin operation")
	}
}

func TestSpinOverMaxCVLengthV2(t *testing.T) {

	var baseVector = "KZY+dsX2jEaZesgCPjJ2Ng.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.214"
	cv, _ := Spin(baseVector)
	if cv.Value() != (baseVector + CVTerminator) {
		t.Errorf("Termination should be applied for CV that goes beyond max length after spin operation")
	}
}
