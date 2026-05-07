package serializer

import (
	"testing"
)

func TestToYAML(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    string
		wantErr bool
	}{
		{
			name: "Valid Server Struct",
			input: Server{ 
				Host:       "localhost",
				Port:       8080,
				Debug:      true,
				AllowedIPs: []string{"192.168.1.1", "10.0.0.1"},
			},
			want:    "host: localhost\nport: 8080\ndebug: true\nallowed_ips:\n    - 192.168.1.1\n    - 10.0.0.1\n",
			wantErr: false,
		},
		{
			name:    "Simple Map",
			input:   map[string]int{"age": 25},
			want:    "age: 25\n",
			wantErr: false,
		},
		{
			name:    "Nil input",
			input:   nil,
			want:    "null\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToYAML(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToYAML() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkToYAML(b *testing.B) {
	data := Server{
		Host:       "localhost",
		Port:       8080,
		Debug:      true,
		AllowedIPs: []string{"192.168.1.1", "10.0.0.1"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ToYAML(data)
	}
}

func BenchmarkToJSON(b *testing.B) {
	data := Server{
		Host:       "localhost",
		Port:       8080,
		Debug:      true,
		AllowedIPs: []string{"192.168.1.1", "10.0.0.1"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ToJSON(data)
	}
}