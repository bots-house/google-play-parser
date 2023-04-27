package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/bots-house/google-play-parser/internal/shared"
)

func Parse(data []byte) (*shared.ParsedObject, error) {
	scriptPattern, err := regexp.Compile(`>AF_initDataCallback[\s\S]*?<\/script`)
	if err != nil {
		return nil, fmt.Errorf("script pattern compile: %w", err)
	}

	keyPattern, err := regexp.Compile(`'(?P<key>ds:.*?)'`)
	if err != nil {
		return nil, fmt.Errorf("key pattern compile: %w", err)
	}

	valuePattern, err := regexp.Compile(`data:(?P<value>[\s\S]*?), sideChannel: {}}\);<\/`)
	if err != nil {
		return nil, fmt.Errorf("value pattern compile: %w", err)
	}

	if err := matches(data, scriptPattern, keyPattern, valuePattern); err != nil {
		return nil, fmt.Errorf("matches: %w", err)
	}

	scriptData := scriptPattern.FindAll(data, -1)

	parsed := parseScriptData(scriptData, keyPattern, valuePattern)

	resultData := make(map[string][]any)
	if err := json.Unmarshal(parsed, &resultData); err != nil {
		return nil, fmt.Errorf("unmarshal app data: %w", err)
	}

	serviceData, err := parseServiceRequests(data)
	if err != nil {
		return nil, fmt.Errorf("parse service requests: %w", err)
	}

	return &shared.ParsedObject{
		Data:        resultData,
		ServiceData: serviceData,
	}, nil
}

func parseScriptData(scriptData [][]byte, keyPattern, valuePattern *regexp.Regexp) []byte {
	parsed := []byte(`{`)

	for idx, data := range scriptData {
		keyData := keyPattern.Find(data)
		valueData := valuePattern.Find(data)

		parsed = append(parsed, keyPattern.ReplaceAll(keyData, []byte(`"$key":`))...)
		parsed = append(parsed, valuePattern.ReplaceAll(valueData, []byte(`$value`))...)

		if idx < len(scriptData)-1 {
			parsed = append(parsed, []byte(`,`)...)
		}
	}

	return append(parsed, []byte(`}`)...)
}

func parseServiceRequests(data []byte) (map[string]shared.Service, error) {
	scriptPattern, err := regexp.Compile(`; var AF_dataServiceRequests[\s\S]*?; var AF_initDataChunkQueue`)
	if err != nil {
		return nil, fmt.Errorf("script pattern compile: %w", err)
	}

	valuePattern, err := regexp.Compile(`{'ds:[\s\S]*}}`)
	if err != nil {
		return nil, fmt.Errorf("value pattern compile: %w", err)
	}

	parsedData := make([]byte, 0, 1024)

	scriptData := scriptPattern.FindAll(data, -1)
	for _, data := range scriptData {
		parsedData = append(parsedData, valuePattern.Find(data)...)
	}

	serviceData, err := parseServiceData(parsedData)
	if err != nil {
		return nil, fmt.Errorf("parse service data: %w", err)
	}

	return serviceData, nil
}

func parseServiceData(data []byte) (map[string]shared.Service, error) {
	pattern, err := regexp.Compile(
		`('(?P<dsKey>ds:\d+)')\s?:\s?({((?P<idKey>id):('(?P<idValue>\w*)'))+[,\s]?((?P<extKey>ext):\s?(?P<extValue>[\w\d.]*)\s?)?,\s?((?P<requestKey>request):(?P<requestValue>[\[\]\w,."\\:\d]*)?)})`,
	)
	if err != nil {
		return nil, fmt.Errorf("value pattern compile: %w", err)
	}

	matches := pattern.FindAll(data, -1)

	result := strings.Builder{}
	result.WriteString("{")

	for idx := range matches {
		value := pattern.ReplaceAll(matches[idx], []byte(`"$dsKey": {"$idKey": "$idValue", "$requestKey": $requestValue, "$extKey": $extValue}`))

		result.Write(value)

		if idx < len(matches)-1 {
			result.WriteString(",")
		}
	}

	result.WriteString("}")

	out := strings.ReplaceAll(result.String(), `, "":`, "")

	serviceMap := make(map[string]shared.Service)
	if err := json.Unmarshal([]byte(out), &serviceMap); err != nil {
		return nil, fmt.Errorf("unmarshal service data: %w", err)
	}

	return serviceMap, nil
}
