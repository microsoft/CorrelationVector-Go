// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package correlationvector contains library functions to manipulate CorrelationVectors.
package correlationvector

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	// MaxVectorLength is the max length of a V1 correlation vector
	MaxVectorLength int = 63

	// MaxVectorLengthV2 is the max length of a V2 correlation vector
	MaxVectorLengthV2 int = 127

	// BaseLength is the max length of a V1 correlation vector base
	BaseLength int = 16

	// BaseLengthV2 is the max length of a V2 correlation vector base
	BaseLengthV2 int = 22
)

// ValidateCorrelationVectorDuringCreation indicates whether or not to validate the
// correlation vector on creation.
var ValidateCorrelationVectorDuringCreation = false

// CorrelationVector represents a lightweight vector for identifying and measuring causality.
type CorrelationVector struct {
	baseVector string
	extension  int32
	version    Version
}

// Version represents a version of the correlation vector protocol.
type Version int

const (
	// V1Version represents the V1 correlation vector version
	V1Version Version = 1

	// V2Version represents the V2 correlation vector version
	V2Version Version = 2
)

// NewCorrelationVector initializes a new instance of the CorrelationVector struct.
// This should only be called when no correlation vector was found in the message header.
func NewCorrelationVector() *CorrelationVector {
	cv, _ := NewCorrelationVectorWithVersion(V1Version)
	return cv
}

// NewCorrelationVectorWithVersion initializes a new instance of the
// CorrelationVector struct of the given protocol version. This should
// only be called when no correlation vector was found in the message header.
func NewCorrelationVectorWithVersion(version Version) (*CorrelationVector, error) {
	base, err := getUniqueValue(version)
	if err != nil {
		return nil, err
	}
	return newCorrelationVector(base, 0, version), nil
}

// Extend creates a new correlation vector by extending an existing value.
// this should be done at the entry point of an operation.
func Extend(correlationVector string) (*CorrelationVector, error) {
	version, err := inferVersion(correlationVector)

	if ValidateCorrelationVectorDuringCreation {
		if err = validate(correlationVector, version); err != nil {
			return nil, err
		}
	}

	return newCorrelationVector(correlationVector, 0, version), err
}

// Parse creates a new correlation vector by parsing its string representation.
func Parse(correlationVector string) (*CorrelationVector, error) {
	version, err := inferVersion(correlationVector)

	p := strings.LastIndex(correlationVector, ".")
	if p > 0 {
		if extension, exterr := strconv.Atoi(correlationVector[p+1:]); exterr == nil && extension >= 0 {
			return newCorrelationVector(correlationVector[:p], extension, version), err
		}
		return nil, errors.New("correlationvector: invalid extension")
	}

	return nil, errors.New("correlationvector: invalid correlation vector string")
}

// Increment increments the current extension by one. Do this before passing
// the value to an outbound message header.
func (cv *CorrelationVector) Increment() string {
	var snapshot int32
	var next int32
	for {
		snapshot = cv.extension
		if snapshot == math.MaxInt32 {
			return cv.Value()
		}
		next = snapshot + 1
		size := len(cv.baseVector) + 1 + int(math.Log10(float64(next))) + 1
		if (cv.version == V1Version && size > int(MaxVectorLength)) || (cv.version == V2Version && size > int(MaxVectorLengthV2)) {
			return cv.Value()
		}
		if atomic.CompareAndSwapInt32(&cv.extension, snapshot, next) {
			return cv.baseVector + "." + strconv.Itoa(int(next))
		}
	}
}

// Value gets the value of the correlation vector as a string.
func (cv *CorrelationVector) Value() string {
	return cv.baseVector + "." + strconv.Itoa(int(cv.extension))
}

// Version gets the version of the correlation vector protocol.
func (cv *CorrelationVector) Version() Version {
	return cv.version
}

func newCorrelationVector(baseVector string, extension int, version Version) *CorrelationVector {
	cv := CorrelationVector{baseVector, int32(extension), version}
	return &cv
}

func getUniqueValue(version Version) (string, error) {
	switch version {
	case V1Version:
		bytes := make([]byte, 12)
		rand.Read(bytes)
		return base64.StdEncoding.EncodeToString(bytes), nil
	case V2Version:
		bytes := make([]byte, 16)
		rand.Read(bytes)
		return base64.StdEncoding.EncodeToString(bytes)[:BaseLengthV2], nil
	}
	return "", errors.New("correlationvector: invalid Version")
}

func inferVersion(correlationVector string) (Version, error) {
	index := strings.Index(correlationVector, ".")

	switch index {
	case BaseLength:
		return V1Version, nil
	case BaseLengthV2:
		return V2Version, nil
	}

	// Default to V1
	return V1Version, errors.New("correlationvector: invalid correlation vector string")
}

func validate(correlationVector string, version Version) error {
	var maxVectorLength int
	var baseLength int

	switch version {
	case V1Version:
		maxVectorLength = MaxVectorLength
		baseLength = BaseLength
	case V2Version:
		maxVectorLength = MaxVectorLengthV2
		baseLength = BaseLengthV2
	default:
		return errors.New("correlationvector: invalid Version")
	}

	if correlationVector == "" || len(correlationVector) > maxVectorLength {
		return fmt.Errorf("correlationvector: the V%d correlation vector cannot be empty or bigger than %d characters", int(version), maxVectorLength)
	}

	parts := strings.Split(correlationVector, ".")

	if len(parts) < 2 || len(parts[0]) != baseLength {
		return fmt.Errorf("correlationvector: invalid correlation vector %s. invalid base value %s", correlationVector, parts[0])
	}

	for i := 1; i < len(parts); i++ {
		if result, err := strconv.Atoi(parts[i]); err != nil || result < 0 {
			return fmt.Errorf("correlationvector: invalid correlation vector %s. invalid extension value %s", correlationVector, parts[i])
		}
	}

	return nil
}
