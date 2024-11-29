package liner

import "testing"

func TestTakeALook(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 bool
	}{
		{
			name: "normal",
			args: args{
				line: "github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0",
			},
			want:  "github.com/grpc-ecosystem/grpc-gateway/v2",
			want1: "v2.22.0",
			want2: true,
		},
		{
			name: "with comment",
			args: args{
				line: "google.golang.org/grpc v1.67.1 // indirect",
			},
			want:  "google.golang.org/grpc",
			want1: "v1.67.1",
			want2: true,
		},
		{
			name: "wrong",
			args: args{
				line: "google.asadas1.1",
			},
			want:  "",
			want1: "",
			want2: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := TakeALook(tt.args.line)
			if got != tt.want {
				t.Errorf("TakeALook() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TakeALook() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("TakeALook() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
