package pkg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	tests := []struct {
		input    string
		expected Value
	}{
		{
			input:    "+OK\r\n",
			expected: Value{typ: "string", str: "OK"},
		},
		{
			input:    "-Error message\r\n",
			expected: Value{typ: "error", str: "Error message"},
		},
		{
			input:    ":123\r\n",
			expected: Value{typ: "integer", num: 123},
		},
		{
			input:    "$6\r\nfoobar\r\n",
			expected: Value{typ: "bulk", bulk: "foobar"},
		},
		{
			input: "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			expected: Value{
				typ: "array",
				array: []Value{
					{typ: "bulk", bulk: "foo"},
					{typ: "bulk", bulk: "bar"},
				},
			},
		},
	}

	for _, tt := range tests {
		r := NewResp(bytes.NewBufferString(tt.input))
		val, err := r.Read()

		assert.NoError(t, err, "unexpected error for input %s", tt.input)
		assert.Equal(t, tt.expected.typ, val.typ, "type mismatch for input %s", tt.input)

		switch val.typ {
		case "string", "error":
			assert.Equal(t, tt.expected.str, val.str, "string mismatch for input %s", tt.input)
		case "integer":
			assert.Equal(t, tt.expected.num, val.num, "integer mismatch for input %s", tt.input)
		case "bulk":
			assert.Equal(t, tt.expected.bulk, val.bulk, "bulk mismatch for input %s", tt.input)
		case "array":
			assert.Equal(t, len(tt.expected.array), len(val.array), "array length mismatch for input %s", tt.input)
			for i, v := range val.array {
				assert.Equal(t, tt.expected.array[i].typ, v.typ, "array element type mismatch at index %d for input %s", i, tt.input)
				assert.Equal(t, tt.expected.array[i].bulk, v.bulk, "array element bulk mismatch at index %d for input %s", i, tt.input)
			}
		}
	}
}

func TestReadInvalidType(t *testing.T) {
	input := "#InvalidType\r\n"
	r := NewResp(bytes.NewBufferString(input))
	_, err := r.Read()

	assert.Error(t, err, "expected error for invalid RESP type")
}
