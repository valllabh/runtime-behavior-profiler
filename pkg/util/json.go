package util

import "encoding/json"

func JsonPrint(i interface{}) {
	j, _ := json.Marshal(i)
	println("\n" + string(j))
}
