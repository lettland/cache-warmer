package structs

import (
	"reflect"
	"testing"
)

func TestCustomFlag_String(t *testing.T) {
	tests := []struct {
		name string
		cf   *CustomFlag
		want string
	}{
		{
			name: "Empty flag values",
			cf:   &CustomFlag{value: []string{}},
			want: "",
		},
		{
			name: "Single flag value",
			cf:   &CustomFlag{value: []string{"flag1"}},
			want: "flag1",
		},
		{
			name: "Multiple flag values",
			cf:   &CustomFlag{value: []string{"flag1", "flag2"}},
			want: "flag1,flag2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cf.String(); got != tt.want {
				t.Errorf("CustomFlag.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomFlag_Set(t *testing.T) {
	tests := []struct {
		name    string
		cf      *CustomFlag
		arg     string
		wantErr error
		want    []string
	}{
		{
			name:    "Set empty string",
			cf:      &CustomFlag{},
			arg:     "",
			wantErr: nil,
			want:    []string{},
		},
		{
			name:    "Set single value",
			cf:      &CustomFlag{},
			arg:     "flag1",
			wantErr: nil,
			want:    []string{"flag1"},
		},
		{
			name:    "Set multiple values",
			cf:      &CustomFlag{},
			arg:     "flag1,flag2",
			wantErr: nil,
			want:    []string{"flag1", "flag2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cf.Set(tt.arg); (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CustomFlag.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.cf.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomFlag.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomFlag_Get(t *testing.T) {
	tests := []struct {
		name string
		cf   *CustomFlag
		want []string
	}{
		{
			name: "Get when no values set",
			cf:   new(CustomFlag),
			want: []string{},
		},
		{
			name: "Get single value",
			cf:   &CustomFlag{value: []string{"flag1"}},
			want: []string{"flag1"},
		},
		{
			name: "Get multiple values",
			cf:   &CustomFlag{value: []string{"flag1", "flag2"}},
			want: []string{"flag1", "flag2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cf.Get()
			if got == nil {
				got = []string{}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomFlag.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomFlag_IsChanged(t *testing.T) {
	tests := []struct {
		name string
		cf   *CustomFlag
		want bool
	}{
		{
			name: "IsChanged when not changed",
			cf:   &CustomFlag{value: []string{"flag1", "flag2"}, changed: false},
			want: false,
		},
		{
			name: "IsChanged when changed",
			cf:   &CustomFlag{value: []string{"flag1", "flag2"}, changed: true},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cf.IsChanged(); got != tt.want {
				t.Errorf("CustomFlag.IsChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCustomFlag(t *testing.T) {
	tests := []struct {
		name string
		want *CustomFlag
	}{
		{
			name: "New custom flag",
			want: &CustomFlag{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCustomFlag(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCustomFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}
