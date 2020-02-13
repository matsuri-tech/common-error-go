package merrors

import (
	"errors"
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
