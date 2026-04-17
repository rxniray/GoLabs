package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Server struct {
	Host       string   `json:"host"`
	Port       int      `json:"port"`
	Debug      bool     `json:"debug"`
	AllowedIPs []string `json:"allowed_ips"`
}

func ToJSON(v any) (string, error) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v) 
 
	if val.Kind() != reflect.Struct {
		return "", errors.New("ToJSON підтримує лише структури")
	}

	var parts []string

	
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)   
		fieldType := typ.Field(i)  

		tag := fieldType.Tag.Get("json")
		if tag == "" {
			tag = strings.ToLower(fieldType.Name) 
		}

		var valStr string

		switch fieldVal.Kind() {
		case reflect.String:
			valStr = fmt.Sprintf(`"%s"`, fieldVal.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			valStr = fmt.Sprintf(`%d`, fieldVal.Int())
		case reflect.Bool:
			valStr = fmt.Sprintf(`%t`, fieldVal.Bool())
		case reflect.Slice:
			var sliceParts []string
			for j := 0; j < fieldVal.Len(); j++ {
				sliceParts = append(sliceParts, fmt.Sprintf("\n\t\t\"%s\"", fieldVal.Index(j).String()))
			}
			valStr = fmt.Sprintf("[%s\n\t]", strings.Join(sliceParts, ","))
		default:
			return "", fmt.Errorf("непідтримуваний тип поля: %s", fieldVal.Kind())
		}

		parts = append(parts, fmt.Sprintf("\t\"%s\": %s", tag, valStr))
	}

	result := "{\n" + strings.Join(parts, ",\n") + "\n}"
	return result, nil
}

func main() {
	srv := Server{
		Host:       "localhost",
		Port:       8080,
		Debug:      true,
		AllowedIPs: []string{"192.168.1.1", "10.0.0.1"},
	}

	jsonStr, err := ToJSON(srv)
	if err != nil {
		fmt.Println("Сталася помилка:", err)
		return
	}

	fmt.Println(jsonStr)
}