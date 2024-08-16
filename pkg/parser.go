package pkg

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// Read determines the RESP type and delegates to the appropriate read function.
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case STRING:
		return r.readString()
	case ERROR:
		return r.readError()
	case INTEGER:
		return r.readInteger()
	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	default:
		fmt.Printf("unknown type: %v", string(_type))
		return Value{}, fmt.Errorf("unkown type: %v", string(_type))
	}
}

// readString reads a simple string (prefixed with '+') and returns it as a Value.
func (r *Resp) readString() (Value, error) {
	v := Value{}
	v.typ = "string"

	line, _, err := r.readLine()
	if err != nil {
		return v, err
	}
	v.str = string(line)

	return v, nil
}

// readError reads an error message (prefixed with '-') and returns it as a Value.
func (r *Resp) readError() (Value, error) {
	v := Value{}
	v.typ = "error"

	line, _, err := r.readLine()
	if err != nil {
		return v, err
	}
	v.str = string(line)

	return v, nil
}

// readInteger reads an integer (prefixed with ':') and returns it as a Value.
func (r *Resp) readInteger() (Value, error) {
	v := Value{}
	v.typ = "integer"

	line, _, err := r.readLine()
	if err != nil {
		return v, err
	}

	num, err := strconv.Atoi(string(line))
	if err != nil {
		return v, err
	}
	v.num = num

	return v, nil
}

// readBulk reads a bulk string (prefixed with '$') and returns it as a Value.
func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	len, _, err := r.readLen()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	_, err = r.reader.Read(bulk)
	if err != nil {
		return v, err
	}

	v.bulk = string(bulk)

	_, _, err = r.readLine()
	if err != nil {
		return v, err
	}

	return v, nil
}

// readArray reads an array (prefixed with '*') and returns it as a Value.
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	len, _, err := r.readLen()
	if err != nil {
		return v, err
	}
	v.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.array = append(v.array, val)
	}

	return v, nil
}

// readLine reads a line of input, terminating at '\r\n', and returns the line minus the '\r\n'.
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

// readLen reads the length prefix of a bulk string or array and returns it as an integer.
func (r *Resp) readLen() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}
