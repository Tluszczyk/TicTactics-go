package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestParseRequestBody(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	testCases := []struct {
		name         string
		requestBody  string
		expectedData TestStruct
		expectedErr  reflect.Type
	}{
		{
			name:        "valid request body",
			requestBody: `{"name": "John", "age": 30}`,
			expectedData: TestStruct{
				Name: "John",
				Age:  30,
			},
			expectedErr: nil,
		},
		{
			name:        "invalid age type",
			requestBody: `{"name": "John", "age": "30"}`,
			expectedData: TestStruct{
				Name: "John",
			},
			expectedErr: reflect.TypeOf(&json.UnmarshalTypeError{}),
		},
		{
			name:         "invalid field names",
			requestBody:  `{"first": "John", "second": 30}`,
			expectedData: TestStruct{},
			expectedErr:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var data TestStruct
			err := ParseRequestBody(tc.requestBody, &data)

			if reflect.TypeOf(err) != tc.expectedErr {
				fmt.Println(err)
				fmt.Println(tc.expectedErr)

				t.Errorf("expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if data != tc.expectedData {
				t.Errorf("expected data: %v, but got: %v", tc.expectedData, data)
			}
		})
	}
}
