package testFunction

import (
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

var (
	errWrongParamType = "输入参数类型错误"
	errWrongParamNum  = "输入参数数量错误"
	errWrongParamZero = "输入参数不能为0"
)

func Add(inArgs ...interface{}) ([]string, string) {
	var ret float64
	if len(inArgs) < 2 {
		return []string{""}, errWrongParamNum
	}
	for i, v := range inArgs {
		parseFloat, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return []string{""}, errWrongParamType
		}
		if i == 0 {
			ret = parseFloat
		} else {
			ret += parseFloat
		}
	}
	return []string{strconv.FormatFloat(ret, 'f', -1, 64)}, ""
}

func Subtract(inArgs ...interface{}) ([]string, string) {
	var ret float64
	if len(inArgs) < 2 {
		return []string{""}, errWrongParamNum
	}
	for i, v := range inArgs {
		parseFloat, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return []string{""}, errWrongParamType
		}
		if i == 0 {
			ret = parseFloat
		} else {
			ret -= parseFloat
		}
	}
	return []string{strconv.FormatFloat(ret, 'f', -1, 64)}, ""
}

func Multiply(inArgs ...interface{}) ([]string, string) {
	var ret float64
	if len(inArgs) < 2 {
		return []string{""}, errWrongParamNum
	}
	for i, v := range inArgs {
		parseFloat, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return []string{""}, errWrongParamType
		}
		if i == 0 {
			ret = parseFloat
		} else {
			ret *= parseFloat
		}
	}
	return []string{strconv.FormatFloat(ret, 'f', -1, 64)}, ""
}

func Divide(inArgs ...interface{}) ([]string, string) {
	var ret float64
	if len(inArgs) < 2 {
		return []string{""}, errWrongParamNum
	}
	for i, v := range inArgs {
		parseFloat, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return []string{""}, errWrongParamType
		}
		if i == 0 {
			ret = parseFloat
		} else {
			if parseFloat == 0 {
				return []string{""}, "除数" + errWrongParamZero
			}
			ret /= parseFloat
		}
	}
	return []string{strconv.FormatFloat(ret, 'f', -1, 64)}, ""
}

func Pow(inArgs ...interface{}) ([]string, string) {
	var ret float64
	if len(inArgs) != 2 {
		return []string{""}, errWrongParamNum
	}
	parseFloatX, err1 := strconv.ParseFloat(inArgs[0].(string), 64)
	parseFloatY, err2 := strconv.ParseFloat(inArgs[1].(string), 64)
	if err1 != nil || err2 != nil {
		return []string{""}, errWrongParamType
	}
	if parseFloatX == 0 {
		return []string{""}, "底数" + errWrongParamZero
	}
	ret = math.Pow(parseFloatX, parseFloatY)
	return []string{strconv.FormatFloat(ret, 'f', -1, 64)}, ""
}

func Sqrt(inArgs ...interface{}) ([]string, string) {
	var ret float64
	if len(inArgs) != 1 {
		return []string{""}, errWrongParamNum
	}
	parseFloat, err := strconv.ParseFloat(inArgs[0].(string), 64)
	if err != nil {
		return []string{""}, errWrongParamType
	}
	ret = math.Pow(parseFloat, 0.5)
	return []string{strconv.FormatFloat(ret, 'f', -1, 64)}, ""
}

func Random(inArgs ...interface{}) ([]string, string) {
	var ret int
	if len(inArgs) != 0 {
		return []string{""}, errWrongParamNum
	}
	rand.Seed(time.Now().Unix())
	ret = rand.Int()
	return []string{strconv.Itoa(ret)}, ""
}

func Swap(inArgs ...interface{}) ([]string, string) {
	if len(inArgs) != 2 {
		return []string{"", ""}, errWrongParamNum
	}
	return []string{inArgs[1].(string), inArgs[0].(string)}, ""
}

func Sort(inArgs ...interface{}) ([]string, string) {
	ret := make([]float64, 0, len(inArgs))
	if len(inArgs) == 0 {
		return []string{""}, errWrongParamNum
	}
	for _, v := range inArgs {
		parseFloat, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return []string{""}, errWrongParamType
		}
		ret = append(ret, parseFloat)
	}
	sort.Float64s(ret)
	sorted := make([]string, 0, len(ret))
	for _, v := range ret {
		sorted = append(sorted, strconv.FormatFloat(v, 'f', -1, 64))
	}
	return sorted, ""
}

func Sleep(inArgs ...interface{}) ([]string, string) {
	if len(inArgs) != 1 {
		return []string{""}, errWrongParamNum
	}
	parseFloat, err := strconv.ParseFloat(inArgs[0].(string), 64)
	if err != nil {
		return []string{""}, errWrongParamType
	}
	time.Sleep(time.Duration(parseFloat) * time.Second)
	return []string{"睡眠完成"}, ""
}

func HeartCheck(inArgs ...interface{}) ([]string, string) {
	if len(inArgs) != 0 {
		return []string{""}, errWrongParamNum
	}
	return []string{"0"}, ""
}
