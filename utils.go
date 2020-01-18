package main

func AssertInt64(obj interface{}) int64 {
	res, ok := obj.(float64)
	if ok {
		return int64(res)
	}
	return int64(0)
}

func AssertFloat64(obj interface{}) float64 {
	res, ok := obj.(float64)
	if ok {
		return res
	}
	return float64(0)
}

func AssertBool(obj interface{}) bool {
	res, ok := obj.(bool)
	if ok {
		return res
	}
	return false
}

func AssertString(obj interface{}) string {
	res, ok := obj.(string)
	if ok {
		return res
	}
	return ""
}
