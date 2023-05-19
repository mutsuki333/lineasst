/*
	util.go
	Purpose: String utiliy functions.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/02  v1.0.0 Evan Chen   Initial release
*/

package util

import (
	"bytes"
	"strconv"
	"text/template"
)

// Tpl is a text/template executer
func Tpl(tmpl string, data any) string {
	engine, err := template.New("").Parse(tmpl)
	if err != nil {
		return tmpl
	}
	buff := &bytes.Buffer{}
	if err := engine.Execute(buff, data); err != nil {
		return err.Error()
	}
	return buff.String()
}

// StripAfterChar strips all chars after the given seperator,
// if the seperator does not exist in @str, @str is returened
func StripAfterChar(str string, sep byte) string {
	for i := len(str) - 1; i > 0; i-- {
		if str[i] == sep {
			return str[i+1:]
		}
	}
	return str
}

// ToInt64X parses string to int64 and panics if error
func ToInt64X(str string) int64 {
	p, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return p
}

// ToInt32X parses string to int32 and panics if error
func ToInt32X(str string) int32 {
	p, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(p)
}

// ToBoolX parses string to bool and panics if error
func ToBoolX(str string) bool {
	p, err := strconv.ParseBool(str)
	if err != nil {
		panic(err)
	}
	return p
}

func ToFloat64X(str string) float64 {
	p, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return p
}
