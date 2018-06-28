// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package correlationvector contains library functions to manipulate CorrelationVectors.
package correlationvector

import (
	"strings"
	"testing"
)

func TestCorrelationVectorIncrementIsUniqueAcrossThreads(t *testing.T) {
	root := NewCorrelationVector()
	vector, _ := Extend(root.Value())

	all := make(chan string, 1000)
	for i := 0; i < 1000; i++ {
		go func() {
			all <- vector.Increment()
		}()
	}

	unique := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		actual := <-all
		if _, ok := unique[actual]; ok {
			t.Errorf("Non unique CV found: %s", actual)
		}
		unique[actual] = true
	}
}

func TestCreateAndIncrementCorrelationVectorDefault(t *testing.T) {
	vector := NewCorrelationVector()
	splitVector := strings.Split(vector.Value(), ".")

	if len(splitVector) != 2 {
		t.Errorf("New vector should have 2 components, got %d", len(splitVector))
		return
	}
	if len(splitVector[0]) != 16 {
		t.Errorf("New vector base should have length 16, got %d", len(splitVector[0]))
	}
	if splitVector[1] != "0" {
		t.Errorf("New vector extension should be 0, got %s", splitVector[1])
	}

	incrementedVector := vector.Increment()
	splitVector = strings.Split(incrementedVector, ".")

	if len(splitVector) != 2 {
		t.Errorf("Incremented vector should have 2 components, got %d", len(splitVector))
		return
	}
	if splitVector[1] != "1" {
		t.Errorf("Incremented vector extension should be 1, got %s", splitVector[1])
	}
}

func TestCreateAndIncrementCorrelationVectorV1(t *testing.T) {
	vector, _ := NewCorrelationVectorWithVersion(V1Version)
	splitVector := strings.Split(vector.Value(), ".")

	if len(splitVector) != 2 {
		t.Errorf("New vector should have 2 components, got %d", len(splitVector))
		return
	}
	if len(splitVector[0]) != 16 {
		t.Errorf("New vector base should have length 16, got %d", len(splitVector[0]))
	}
	if splitVector[1] != "0" {
		t.Errorf("New vector extension should be 0, got %s", splitVector[1])
	}

	incrementedVector := vector.Increment()
	splitVector = strings.Split(incrementedVector, ".")

	if len(splitVector) != 2 {
		t.Errorf("Incremented vector should have 2 components, got %d", len(splitVector))
		return
	}
	if splitVector[1] != "1" {
		t.Errorf("Incremented vector extension should be 1, got %s", splitVector[1])
	}
}

func TestCreateAndIncrementCorrelationVectorV2(t *testing.T) {
	vector, _ := NewCorrelationVectorWithVersion(V2Version)
	splitVector := strings.Split(vector.Value(), ".")

	if len(splitVector) != 2 {
		t.Errorf("New vector should have 2 components, got %d", len(splitVector))
		return
	}
	if len(splitVector[0]) != 22 {
		t.Errorf("New vector base should have length 22, got %d", len(splitVector[0]))
	}
	if splitVector[1] != "0" {
		t.Errorf("New vector extension should be 0, got %s", splitVector[1])
	}

	incrementedVector := vector.Increment()
	splitVector = strings.Split(incrementedVector, ".")

	if len(splitVector) != 2 {
		t.Errorf("Incremented vector should have 2 components, got %d", len(splitVector))
		return
	}
	if splitVector[1] != "1" {
		t.Errorf("Incremented vector extension should be 1, got %s", splitVector[1])
	}
}

func TestCreateCorrelationVectorFromString(t *testing.T) {
	vector, _ := Extend("tul4NUsfs9Cl7mOf.1")
	splitVector := strings.Split(vector.Value(), ".")

	if len(splitVector) != 3 {
		t.Errorf("Extended vector should have 3 components, got %d", len(splitVector))
		return
	}
	if splitVector[2] != "0" {
		t.Errorf("Extended vector extension should be 0, got %s", splitVector[2])
	}

	incrementedVector := vector.Increment()
	splitVector = strings.Split(incrementedVector, ".")

	if len(splitVector) != 3 {
		t.Errorf("Incremented vector should have 3 components, got %d", len(splitVector))
		return
	}
	if splitVector[2] != "1" {
		t.Errorf("Incremented vector extension should be 1, got %s", splitVector[2])
	}
	if vector.Value() != "tul4NUsfs9Cl7mOf.1.1" {
		t.Errorf("Incremented vector value should be tul4NUsfs9Cl7mOf.1.1, got %s", vector.Value())
	}
}

func TestCreateCorrelationVectorFromStringV2(t *testing.T) {
	vector, _ := Extend("KZY+dsX2jEaZesgCPjJ2Ng.1")
	splitVector := strings.Split(vector.Value(), ".")

	if len(splitVector) != 3 {
		t.Errorf("Extended vector should have 3 components, got %d", len(splitVector))
		return
	}
	if splitVector[2] != "0" {
		t.Errorf("Extended vector extension should be 0, got %s", splitVector[2])
	}

	incrementedVector := vector.Increment()
	splitVector = strings.Split(incrementedVector, ".")

	if len(splitVector) != 3 {
		t.Errorf("Incremented vector should have 3 components, got %d", len(splitVector))
		return
	}
	if splitVector[2] != "1" {
		t.Errorf("Incremented vector extension should be 1, got %s", splitVector[2])
	}
	if vector.Value() != "KZY+dsX2jEaZesgCPjJ2Ng.1.1" {
		t.Errorf("Incremented vector value should be KZY+dsX2jEaZesgCPjJ2Ng.1.1, got %s", vector.Value())
	}
}

func TestExtendEmptyCorrelationVector(t *testing.T) {
	vector, err := Extend("")
	if vector.Value() != ".0" {
		t.Errorf("Extending empty correlation vector string should result in value .0, got %s", vector.Value())
	}
	if err == nil {
		t.Errorf("Extending empty correlation vector string should return error")
	}

	ValidateCorrelationVectorDuringCreation = true
	vector, err = Extend("")
	if vector != nil {
		t.Errorf("Extending empty correlation vector string with validation should return nil")
	}
	if err == nil {
		t.Errorf("Extending empty correlation vector string with validation should return error")
	}
	ValidateCorrelationVectorDuringCreation = false
}

func TestExtendInsufficientCharsCorrelationVector(t *testing.T) {
	vector, err := Extend("tul4NUsfs9Cl7mO.1")
	if vector.Value() != "tul4NUsfs9Cl7mO.1.0" {
		t.Errorf("Extending insufficient characters correlation vector string should result in value tul4NUsfs9Cl7mO.1.0, got %s", vector.Value())
	}
	if err == nil {
		t.Errorf("Extending insufficient characters correlation vector string should return error")
	}

	ValidateCorrelationVectorDuringCreation = true
	vector, err = Extend("tul4NUsfs9Cl7mO.1")
	if vector != nil {
		t.Errorf("Extending insufficient characters correlation vector string with validation should return nil")
	}
	if err == nil {
		t.Errorf("Extending insufficient characters correlation vector string with validation should return error")
	}
	ValidateCorrelationVectorDuringCreation = false
}

func TestExtendTooManyCharsCorrelationVector(t *testing.T) {
	vector, err := Extend("tul4NUsfs9Cl7mOfN/dupsl.1")
	if vector.Value() != "tul4NUsfs9Cl7mOfN/dupsl.1.0" {
		t.Errorf("Extending too many characters correlation vector string should result in value tul4NUsfs9Cl7mOfN/dupsl.1.0, got %s", vector.Value())
	}
	if err == nil {
		t.Errorf("Extending too many characters correlation vector string should return error")
	}

	ValidateCorrelationVectorDuringCreation = true
	vector, err = Extend("tul4NUsfs9Cl7mOfN/dupsl.1")
	if vector != nil {
		t.Errorf("Extending too many characters correlation vector string with validation should return nil")
	}
	if err == nil {
		t.Errorf("Extending too many characters correlation vector string with validation should return error")
	}
	ValidateCorrelationVectorDuringCreation = false
}

func TestExtendTooBigCorrelationVector(t *testing.T) {
	ValidateCorrelationVectorDuringCreation = true
	// Bigger than 63 chars
	vector, err := Extend("tul4NUsfs9Cl7mOf.2147483647.2147483647.2147483647.2147483647.2147483647")
	if vector != nil {
		t.Errorf("Extending too big correlation vector string with validation should return nil")
	}
	if err == nil {
		t.Errorf("Extending too big correlation vector string with validation should return error")
	}
	ValidateCorrelationVectorDuringCreation = false
}

func TestExtendTooBigCorrelationVectorV2(t *testing.T) {
	ValidateCorrelationVectorDuringCreation = true
	// Bigger than 127 chars
	vector, err := Extend("KZY+dsX2jEaZesgCPjJ2Ng.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647")
	if vector != nil {
		t.Errorf("Extending too big correlation vector string with validation should return nil")
	}
	if err == nil {
		t.Errorf("Extending too big correlation vector string with validation should return error")
	}
	ValidateCorrelationVectorDuringCreation = false
}

func TestExtendTooBigExtensionCorrelationVector(t *testing.T) {
	ValidateCorrelationVectorDuringCreation = true
	// Bigger int32
	vector, err := Extend("tul4NUsfs9Cl7mOf.11111111111111111111111111111")
	if vector != nil {
		t.Errorf("Extending too big extension correlation vector string with validation should return nil")
	}
	if err == nil {
		t.Errorf("Extending too big extension correlation vector string with validation should return error")
	}
	ValidateCorrelationVectorDuringCreation = false
}

func TestIncrementPastMaxWithTerminator(t *testing.T) {
	vector, _ := Extend("tul4NUsfs9Cl7mOf.2147483647.2147483647.2147483647.21474836479")
	vector.Increment()
	if vector.Value() != "tul4NUsfs9Cl7mOf.2147483647.2147483647.2147483647.21474836479.1" {
		t.Errorf("Incrementing correlation vector should return tul4NUsfs9Cl7mOf.2147483647.2147483647.2147483647.21474836479.1, got %s", vector.Value())
	}

	for i := 0; i < 20; i++ {
		vector.Increment()
	}

	// We hit 63 chars so we silently stopped counting and add the terminator
	if vector.Value() != "tul4NUsfs9Cl7mOf.2147483647.2147483647.2147483647.21474836479.9!" {
		t.Errorf("Incrementing past max correlation vector should return tul4NUsfs9Cl7mOf.2147483647.2147483647.2147483647.21474836479.9!, got %s", vector.Value())
	}
}

func TestIncrementPastMaxWithTerminatorV2(t *testing.T) {
	vector, _ := Extend("KZY+dsX2jEaZesgCPjJ2Ng.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.214")
	vector.Increment()
	if vector.Value() != "KZY+dsX2jEaZesgCPjJ2Ng.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.214.1" {
		t.Errorf("Incrementing correlation vector should return KZY+dsX2jEaZesgCPjJ2Ng.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.214.1, got %s", vector.Value())
	}

	for i := 0; i < 20; i++ {
		vector.Increment()
	}

	// We hit 127 chars so we silently stopped counting and add the terminator
	if vector.Value() != "KZY+dsX2jEaZesgCPjJ2Ng.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.214.9!" {
		t.Errorf("Incrementing past max correlation vector should return KZY+dsX2jEaZesgCPjJ2Ng.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.214.9!, got %s", vector.Value())
	}
}

func TestExtendOverMaxCVLength(t *testing.T) {
	var baseVector = "tul4NUsfs9Cl7mOf.2147483647.2147483647.2147483647.214748364.23"
	cv, _ := Extend(baseVector)
	if cv.Value() != (baseVector + CVTerminator) {
		t.Errorf("Extending cv with max length should be appended with !, got %s", cv.Value())
	}
}

func TestExtendOverMaxCVLengthV2(t *testing.T) {
	var baseVector = "KZY+dsX2jEaZesgCPjJ2Ng.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2141"
	cv, _ := Extend(baseVector)
	if cv.Value() != (baseVector + CVTerminator) {
		t.Errorf("Extending cv with max length should be appended with !, got %s", cv.Value())
	}
}

func TestImmutableCVWithTerminator(t *testing.T) {
	var cvStr = "tul4NUsfs9Cl7mOf.2147483647.2147483647.2147483647.21474836479.0!"

	cv1, _ := Parse(cvStr)
	if cvStr != cv1.Increment() {
		t.Errorf("Terminated CV should remain unchanged after increment operation")
	}
	cv2, _ := Extend(cvStr)
	if cvStr != cv2.Value() {
		t.Errorf("Terminated CV should remain unchanged after extend operation")
	}
	cv3, _ := Spin(cvStr)
	if cvStr != cv3.Value() {
		t.Errorf("Terminated CV should remain unchanged after spin operation")
	}
}

func TestImmutableCVWithTerminatorV2(t *testing.T) {
	var cvStr = "KZY+dsX2jEaZesgCPjJ2Ng.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.2147483647.214.0!"

	cv1, _ := Parse(cvStr)
	if cvStr != cv1.Increment() {
		t.Errorf("Terminated CV should remain unchanged after increment operation")
	}
	cv2, _ := Extend(cvStr)
	if cvStr != cv2.Value() {
		t.Errorf("Terminated CV should remain unchanged after extend operation")
	}
	cv3, _ := Spin(cvStr)
	if cvStr != cv3.Value() {
		t.Errorf("Terminated CV should remain unchanged after spin operation")
	}
}
