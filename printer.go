package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"text/tabwriter"
)

func TableString(iface interface{}) string {
	b := &bytes.Buffer{}
	w := tabwriter.NewWriter(b, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)

	onStruct := func(f reflect.StructField, v reflect.Value) error {
		fmt.Fprintf(w, "##### %s #####\n", f.Name)
		return nil
	}
	onStructField := func(f reflect.StructField, v reflect.Value) error {
		val := v.Interface()
		if v, ok := f.Tag.Lookup("print"); ok && v == "-" {
			val = "*** Hidden value ***"
		}
		fmt.Fprintf(w, "%s\t\x1b[0m%v\t\x1b[1;34m%s\x1b[0m \x1b[1;92m`%s`\x1b[0m\n", f.Name, val, v.Type().String(), f.Tag)
		return nil
	}

	fmt.Fprint(w, "\n-----------------------------------\n")
	walk(iface, onStruct, onStructField)
	fmt.Fprint(w, "-----------------------------------\n")
	w.Flush()

	return b.String()
}

func JSON(iface interface{}) ([]byte, error) {
	m := make(map[string]interface{})

	onStruct := func(f reflect.StructField, v reflect.Value) error {
		return nil
	}
	onStructField := func(f reflect.StructField, v reflect.Value) error {
		val := v.Interface()
		if tagValue, ok := f.Tag.Lookup("print"); ok && tagValue == "-" {
			val = "*** Hidden value ***"
		}
		if name, ok := f.Tag.Lookup("config"); ok {
			m[name] = val
		}
		return nil
	}

	walk(iface, onStruct, onStructField)

	return json.Marshal(m)
}

func walk(iface interface{}, onStruct, onStructField func(f reflect.StructField, v reflect.Value) error) {
	value := reflect.ValueOf(iface)

	if value.Kind() != reflect.Struct {
		value = value.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if !field.CanInterface() {
			continue
		}
		typeField := value.Type().Field(i)
		if field.Kind() == reflect.Struct {
			onStruct(typeField, field)
			iface := field.Interface()
			walk(iface, onStruct, onStructField)
			continue
		}

		onStructField(typeField, field)
	}
}
