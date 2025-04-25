package args_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wetfloo/stupserv/args"
)

func TestParseNoPath(t *testing.T) {
	expected := args.ArgValues{
		Cache: false,
		Addr:  ":6040",
	}
	actual := args.ParseArgs()

	vExpected := reflect.ValueOf(expected)
	vActual := reflect.ValueOf(actual)
	typeExpected := vExpected.Type()
	typeActual := vActual.Type()
	require.Equal(t, typeExpected, typeActual)
	require.Equal(t, vExpected.NumField(), vActual.NumField())
	require.Equal(t, vExpected.NumMethod(), vActual.NumMethod())

	visibleFieldsExpected := reflect.VisibleFields(typeExpected)
	visibleFieldsActual := reflect.VisibleFields(typeActual)
	require.EqualValues(t, visibleFieldsExpected, visibleFieldsActual)

	for i, fieldExpected := range visibleFieldsExpected {
		fieldActual := visibleFieldsActual[i]
		nameExpected := typeExpected.Field(i).Name
		nameActual := typeActual.Field(i).Name

		assert.Equal(t, nameExpected, nameActual)
		if strings.ToLower(nameExpected) == "path" {
			continue
		}

		fieldValExpected := vExpected.FieldByName(fieldExpected.Name).Interface()
		fieldValActual := vExpected.FieldByName(fieldActual.Name).Interface()

		assert.Equal(t, fieldValExpected, fieldValActual)
	}
}
