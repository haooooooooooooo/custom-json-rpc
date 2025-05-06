package testFunction

import (
	"errors"
	"math"
	"math/rand"
	"sort"
	"time"
)

var (
	errWrongParamType = "Invalid parameter type"
	errWrongParamNum  = "Invalid number of parameters"
	errWrongParamZero = "Parameter cannot be zero"
)

func Add(inArgs ...any) ([]any, error) {
	var res float64
	if len(inArgs) < 2 {
		return []any{}, errors.New(errWrongParamNum)
	}
	for i, v := range inArgs {
		val, ok := v.(float64)
		if !ok {
			return []any{}, errors.New(errWrongParamType)
		}
		if i == 0 {
			res = val
		} else {
			res += val
		}
	}
	return []any{res}, nil
}

func Subtract(inArgs ...any) ([]any, error) {
	var res float64
	if len(inArgs) < 2 {
		return []any{}, errors.New(errWrongParamNum)
	}
	for i, v := range inArgs {
		val, ok := v.(float64)
		if !ok {
			return []any{}, errors.New(errWrongParamType)
		}
		if i == 0 {
			res = val
		} else {
			res -= val
		}
	}
	return []any{res}, nil
}

func Multiply(inArgs ...any) ([]any, error) {
	var res float64 = 1
	if len(inArgs) < 2 {
		return []any{}, errors.New(errWrongParamNum)
	}
	for _, v := range inArgs {
		val, ok := v.(float64)
		if !ok {
			return []any{}, errors.New(errWrongParamType)
		}
		res *= val
	}
	return []any{res}, nil
}

func Divide(inArgs ...any) ([]any, error) {
	var res float64
	if len(inArgs) < 2 {
		return []any{}, errors.New(errWrongParamNum)
	}
	for i, v := range inArgs {
		val, ok := v.(float64)
		if !ok {
			return []any{}, errors.New(errWrongParamType)
		}
		if i == 0 {
			res = val
		} else {
			if val == 0 {
				return []any{}, errors.New("divisor: " + errWrongParamZero)
			}
			res /= val
		}
	}
	return []any{res}, nil
}

func Pow(inArgs ...any) ([]any, error) {
	if len(inArgs) != 2 {
		return []any{}, errors.New(errWrongParamNum)
	}
	base, ok1 := inArgs[0].(float64)
	exp, ok2 := inArgs[1].(float64)
	if !ok1 || !ok2 {
		return []any{}, errors.New(errWrongParamType)
	}
	if base == 0 {
		return []any{}, errors.New("base: " + errWrongParamZero)
	}
	res := math.Pow(base, exp)
	return []any{res}, nil
}

func Sqrt(inArgs ...any) ([]any, error) {
	if len(inArgs) != 1 {
		return []any{}, errors.New(errWrongParamNum)
	}
	val, ok := inArgs[0].(float64)
	if !ok {
		return []any{}, errors.New(errWrongParamType)
	}
	res := math.Sqrt(val)
	return []any{res}, nil
}

func Random(inArgs ...any) ([]any, error) {
	if len(inArgs) != 0 {
		return []any{}, errors.New(errWrongParamNum)
	}
	rand.Seed(time.Now().Unix())
	res := rand.Float64()
	return []any{res}, nil
}

func Swap(inArgs ...any) ([]any, error) {
	if len(inArgs) != 2 {
		return []any{}, errors.New(errWrongParamNum)
	}
	return []any{inArgs[1], inArgs[0]}, nil
}

func Sort(inArgs ...any) ([]any, error) {
	if len(inArgs) == 0 {
		return []any{}, errors.New(errWrongParamNum)
	}
	ret := make([]float64, 0, len(inArgs))
	for _, v := range inArgs {
		val, ok := v.(float64)
		if !ok {
			return []any{}, errors.New(errWrongParamType)
		}
		ret = append(ret, val)
	}
	sort.Float64s(ret)
	res := make([]any, len(ret))
	for i, v := range ret {
		res[i] = v
	}
	return res, nil
}

func Sleep(inArgs ...any) ([]any, error) {
	if len(inArgs) != 1 {
		return []any{}, errors.New(errWrongParamNum)
	}
	val, ok := inArgs[0].(float64)
	if !ok {
		return []any{}, errors.New(errWrongParamType)
	}
	time.Sleep(time.Duration(val) * time.Second)
	return []any{"sleep done"}, nil
}

func HeartCheck(inArgs ...any) ([]any, error) {
	if len(inArgs) != 0 {
		return []any{}, errors.New(errWrongParamNum)
	}
	return []any{0}, nil
}
