package merrors

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestErrorTypeEqual(t *testing.T) {
	tests := []struct {
		x    error
		y    error
		want bool
	}{
		// 非CommonErrorの場合、reflect.DeepEqualで判断
		{
			x:    errors.New("same_message"),
			y:    errors.New("same_message"),
			want: true,
		},
		{
			x:    errors.New("different_message_1"),
			y:    errors.New("different_message_2"),
			want: false,
		},
		// CommonError同士の場合、ErrorTypeが同じであればMessageが異なってもtrue
		{
			x:    ErrorBadRequest("different_message_1", "same_type"),
			y:    ErrorBadRequest("different_message_2", "same_type"),
			want: true,
		},
		// Messageが同じでもErrorTypeが違うならfalse
		{
			x:    ErrorBadRequest("same_message", "different_type_1"),
			y:    ErrorBadRequest("same_message", "different_type_2"),
			want: false,
		},
		// CommonErrorとそれ以外のエラーの比較は必ずfalse
		{
			x:    errors.New("same_message"),
			y:    ErrorBadRequest("same_message", "type"),
			want: false,
		},
		{
			x:    ErrorBadRequest("same_message", "type"),
			y:    errors.New("same_message"),
			want: false,
		},
	}

	for _, tt := range tests {
		if got := ErrorTypeEqual(tt.x, tt.y); got != tt.want {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	}
}

func TestErrorResponseFromJson(t *testing.T) {
	tests := []struct {
		in   []byte
		want ErrorResponse
	}{
		{
			in: []byte(`{"error":"case_1", "errorType": "snake_case"}`),
			want: ErrorResponse{
				Error:     "case_1",
				ErrorType: "snake_case",
			},
		},
		{
			in: []byte(`{"error":"case2", "errorType": "camelCase"}`),
			want: ErrorResponse{
				Error:     "case2",
				ErrorType: "camelCase",
			},
		},
		{
			in: []byte(`{"error":"case_3", "error_type": "snake_case"}`),
			want: ErrorResponse{
				Error:     "case_3",
				ErrorType: "snake_case",
			},
		},
		{
			in: []byte(`{"error":"case4", "error_type": "camelCase"}`),
			want: ErrorResponse{
				Error:     "case4",
				ErrorType: "camelCase",
			},
		},
	}
	for _, tt := range tests {
		e := &ErrorResponse{}
		if err := json.Unmarshal(tt.in, e); err != nil {
			t.Errorf("err: %v", err)
			continue
		}

		if !reflect.DeepEqual(*e, tt.want) {
			t.Errorf("expected: %v, result: %v", tt.want, e)
		}
	}
}
