package system

type Map map[string]interface{}

func ResultMap(err Error, data interface{}) Map {
	if data != nil {
		return Map{"status": err.Map(), "data": data}
	} else {
		return Map{"status": err.Map()}
	}
}

func SuccessResultMap(data interface{}) Map {
	if data != nil {
		return Map{"status": SuccessStatus.Map(), "data": data}
	} else {
		return Map{"status": SuccessStatus.Map()}
	}
}
