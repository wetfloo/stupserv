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
	expected := args.Values{
		Cache: false,
		Addr:  ":6040",
	}
	actual := args.ParseArgs([]string{})

	testEqualArgValues(t, expected, actual, false)
}

func TestParseOnlyPath(t *testing.T) {
	expected := args.Values{
		Cache: false,
		Addr:  ":6040",
		Path:  "/some/path",
	}
	actual := args.ParseArgs([]string{expected.Path})

	testEqualArgValues(t, expected, actual, true)
}

func TestParseCacheAddr(t *testing.T) {
	expected := args.Values{
		Cache: false,
		Addr:  "127.0.0.1:1337",
	}
	actual := args.ParseArgs([]string{"-c", "-a", expected.Addr})

	testEqualArgValues(t, expected, actual, false)
}

func TestParseCacheAddrPath(t *testing.T) {
	expected := args.Values{
		Cache: false,
		Addr:  "127.0.0.1:1337",
		Path:  "/some/path",
	}
	actual := args.ParseArgs([]string{"-c", "-a", expected.Addr, expected.Path})

	testEqualArgValues(t, expected, actual, true)
}

func TestParseCacheAddrMalformed(t *testing.T) {
	expected := args.Values{
		Cache: false,
		Addr:  "127.0.0.1:1337",
	}
	actual := args.ParseArgs([]string{expected.Path, "-c", "-a", expected.Addr})

	testEqualArgValues(t, expected, actual, true)
}

func testEqualArgValues(
	t *testing.T,
	expected args.Values,
	actual args.Values,
	testPathEq bool,
) {
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
		if !testPathEq && strings.ToLower(nameExpected) == "path" {
			continue
		}

		fieldValExpected := vExpected.FieldByName(fieldExpected.Name).Interface()
		fieldValActual := vExpected.FieldByName(fieldActual.Name).Interface()

		assert.Equal(t, fieldValExpected, fieldValActual)
	}
}
