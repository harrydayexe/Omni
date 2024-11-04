package balancer

import (
	"net/url"
	"reflect"
	"testing"
)

func TestRoundRobinAdd(t *testing.T) {
	// Test code here
	cases := []struct {
		name string
		bal  Balancer
		args []url.URL
		want Balancer
	}{
		{
			name: "Add one URL to empty balancer",
			bal:  &RoundRobinBalancer{},
			args: []url.URL{{Scheme: "http", Host: "localhost:8080"}},
			want: &RoundRobinBalancer{
				BaseBalancer: BaseBalancer{
					servers: []url.URL{{Scheme: "http", Host: "localhost:8080"}},
				},
			},
		},
		{
			name: "Add two URLs to empty balancer",
			bal:  &RoundRobinBalancer{},
			args: []url.URL{
				{Scheme: "http", Host: "localhost:8080"},
				{Scheme: "http", Host: "localhost:4040"},
			},
			want: &RoundRobinBalancer{
				BaseBalancer: BaseBalancer{
					servers: []url.URL{
						{Scheme: "http", Host: "localhost:8080"},
						{Scheme: "http", Host: "localhost:4040"},
					},
				},
			},
		},
		{
			name: "Add URL to filled balancer",
			bal: &RoundRobinBalancer{
				BaseBalancer: BaseBalancer{
					servers: []url.URL{
						{Scheme: "http", Host: "localhost:8080"},
					},
				},
			},
			args: []url.URL{
				{Scheme: "http", Host: "localhost:4040"},
			},
			want: &RoundRobinBalancer{
				BaseBalancer: BaseBalancer{
					servers: []url.URL{
						{Scheme: "http", Host: "localhost:8080"},
						{Scheme: "http", Host: "localhost:4040"},
					},
				},
			},
		},
		{
			name: "Add existing URL to filled balancer",
			bal: &RoundRobinBalancer{
				BaseBalancer: BaseBalancer{
					servers: []url.URL{
						{Scheme: "http", Host: "localhost:8080"},
					},
				},
			},
			args: []url.URL{
				{Scheme: "http", Host: "localhost:8080"},
			},
			want: &RoundRobinBalancer{
				BaseBalancer: BaseBalancer{
					servers: []url.URL{
						{Scheme: "http", Host: "localhost:8080"},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, arg := range c.args {
				c.bal.Add(arg)
			}
			if !reflect.DeepEqual(c.bal, c.want) {
				t.Errorf("Adding args: %v,\n got: %v,\nwant: %v", c.args, c.bal, c.want)
			}
		})
	}
}
