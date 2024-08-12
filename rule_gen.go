// c:\code\rule_gen\rule_gen.go

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Condition struct {
	Fact     string  `json:"fact"`
	Operator string  `json:"operator"`
	Value    float64 `json:"value"`
}

type Action struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	Value  string `json:"value"`
}

type Rule struct {
	Name       string `json:"name"`
	Conditions struct {
		All []Condition `json:"all"`
	} `json:"conditions"`
	Actions []Action `json:"actions"`
}

type Rules struct {
	Rules []Rule `json:"rules"`
}

func main() {
	numSensors := 1000 // Change this to match the DIU simulator configuration
	sensorChannels := []string{"temperature", "pressure", "humidity"}
	rules := Rules{}

	for i := 0; i < numSensors; i++ {
		channel := sensorChannels[i%len(sensorChannels)]
		sensorName := fmt.Sprintf("%s:sensor_%03d", channel, i)

		switch channel {
		case "temperature":
			rules.Rules = append(rules.Rules, Rule{
				Name: fmt.Sprintf("High%sAlert_%s", capitalize(channel), sensorName),
				Conditions: struct {
					All []Condition `json:"all"`
				}{
					All: []Condition{
						{Fact: sensorName, Operator: "GT", Value: 32},
					},
				},
				Actions: []Action{
					{Type: "updateStore", Target: fmt.Sprintf("alerts:sensor_%03d", i), Value: fmt.Sprintf("%s exceeds 32 degrees", capitalize(channel))},
				},
			})
		case "pressure":
			rules.Rules = append(rules.Rules, Rule{
				Name: fmt.Sprintf("Low%sAlert_%s", capitalize(channel), sensorName),
				Conditions: struct {
					All []Condition `json:"all"`
				}{
					All: []Condition{
						{Fact: sensorName, Operator: "LT", Value: 0.9},
					},
				},
				Actions: []Action{
					{Type: "updateStore", Target: fmt.Sprintf("alerts:sensor_%03d", i), Value: fmt.Sprintf("%s below 0.9 atm", capitalize(channel))},
				},
			})
		case "humidity":
			rules.Rules = append(rules.Rules, Rule{
				Name: fmt.Sprintf("High%sAlert_%s", capitalize(channel), sensorName),
				Conditions: struct {
					All []Condition `json:"all"`
				}{
					All: []Condition{
						{Fact: sensorName, Operator: "GT", Value: 85},
					},
				},
				Actions: []Action{
					{Type: "updateStore", Target: fmt.Sprintf("alerts:sensor_%03d", i), Value: fmt.Sprintf("%s exceeds 85%%", capitalize(channel))},
				},
			})
		}
	}

	// Write rules to JSON file
	file, err := os.Create("rules.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(rules); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-'a'+'A') + s[1:]
}
