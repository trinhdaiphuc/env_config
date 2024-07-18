package env_config

import "strconv"

type FloatType interface {
	float32 | float64
}

type IntType interface {
	int | int8 | int16 | int32 | int64
}

type UintType interface {
	uint | uint8 | uint16 | uint32 | uint64
}

func StringArrayToFloatArray[F FloatType](strings []string) []F {
	floats := make([]F, len(strings))
	for i, s := range strings {
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			floats[i] = F(f)
		}
	}
	return floats
}

func StringArrayToIntArray[I IntType](strings []string) []I {
	ints := make([]I, len(strings))
	for i, s := range strings {
		n, err := strconv.Atoi(s)
		if err == nil {
			ints[i] = I(n)
		}
	}
	return ints
}

func StringArrayToUintArray[U UintType](strings []string) []U {
	uints := make([]U, len(strings))
	for i, s := range strings {
		n, err := strconv.ParseUint(s, 10, 64)
		if err == nil {
			uints[i] = U(n)
		}
	}
	return uints
}
