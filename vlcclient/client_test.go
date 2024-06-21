package vlcclient

import "testing"

func TestClient_constructUrl(t *testing.T) {
	type fields struct {
		Addr         string
		HttpUser     string
		HttpPassword string
	}
	type args struct {
		endpoint    string
		queryParams map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "With port, without query string",
			fields: fields{
				Addr:         "localhost:9090",
				HttpUser:     "",
				HttpPassword: "tvcontroller123",
			},
			args: args{
				endpoint:    "/requests/status.json",
				queryParams: nil,
			},
			want: "localhost:9090/requests/status.json",
		},
		{
			name: "With query string",
			fields: fields{
				Addr:         "localhost:9090",
				HttpUser:     "",
				HttpPassword: "tvcontroller123",
			},
			args: args{
				endpoint: "/requests/status.json",
				queryParams: map[string]string{
					"key1": "val1",
					"key2": "val2",
				},
			},
			want: "localhost:9090/requests/status.json?key1=val1&key2=val2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				Addr:         tt.fields.Addr,
				HttpUser:     tt.fields.HttpUser,
				HttpPassword: tt.fields.HttpPassword,
			}
			if got := c.constructUrl(tt.args.endpoint, tt.args.queryParams); got != tt.want {
				t.Errorf("constructUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
