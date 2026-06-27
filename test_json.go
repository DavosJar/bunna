package main

import (
	"encoding/json"
	"fmt"
)

type Req struct {
	LoteID string `json:"loteID"`
}

func main() {
	var r Req
	err := json.Unmarshal([]byte(`{"loteID": null}`), &r)
	fmt.Printf("err: %v, val: '%s'\n", err, r.LoteID)
}
