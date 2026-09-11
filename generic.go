package env_config

import (
	"fmt"
	"strconv"
)

type FloatType interface {
	float32 | float64
}

type IntType interface {
	int | int8 | int16 | int32 | int64
}

type UintType interface {
	uint | uint8 | uint16 | uint32 | uint64
}

func StringArrayToFloatArray[F FloatType](strings []string) ([]F, error) {
	floats := make([]F, len(strings))
	bitSize := 32
	if _, ok := any(F(0)).(float64); ok {
		bitSize = 64
	}
	for i, s := range strings {
		f, err := strconv.ParseFloat(s, bitSize)
		if err != nil {
			return nil, fmt.Errorf("parsing element %q: %w", s, err)
		}
		floats[i] = F(f)
	}
	return floats, nil
}

func StringArrayToIntArray[I IntType](strings []string) ([]I, error) {
	ints := make([]I, len(strings))
	bitSize := intBitSize(any(I(0)))
	for i, s := range strings {
		n, err := strconv.ParseInt(s, 10, bitSize)
		if err != nil {
			return nil, fmt.Errorf("parsing element %q: %w", s, err)
		}
		ints[i] = I(n)
	}
	return ints, nil
}

func StringArrayToUintArray[U UintType](strings []string) ([]U, error) {
	uints := make([]U, len(strings))
	bitSize := intBitSize(any(U(0)))
	for i, s := range strings {
		n, err := strconv.ParseUint(s, 10, bitSize)
		if err != nil {
			return nil, fmt.Errorf("parsing element %q: %w", s, err)
		}
		uints[i] = U(n)
	}
	return uints, nil
}

// intBitSize returns the strconv bit size of the concrete element type, so an
// out-of-range element fails instead of being silently truncated.
func intBitSize(sample any) int {
	switch sample.(type) {
	case int8, uint8:
		return 8
	case int16, uint16:
		return 16
	case int32, uint32:
		return 32
	default:
		return 64
	}
}
