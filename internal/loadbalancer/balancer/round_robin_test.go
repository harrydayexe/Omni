package balancer

import (
	"net/url"
	"reflect"
	"testing"
)

func TestRoundRobinFactory(t *testing.T) {
	cases := []struct {
		name string
		args []url.URL
		want Balancer
	}{
		{
			name: "Create Round Robin Balancer with one URL",
			args: []url.URL{{Scheme: "http", Host: "localhost:8080"}},
			want: &RoundRobinBalancer{
				BaseBalancer: BaseBalancer{
					servers: []url.URL{{Scheme: "http", Host: "localhost:8080"}},
				},
			},
		},
		{
			name: "Create Round Robin Balancer with two URLs",
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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bal, err := BuildBalancer("round-robin", c.args)
			if err != nil {
				t.Errorf("Test: %s failed\nunexpected error: %v", c.name, err)
			}
			if !reflect.DeepEqual(bal, c.want) {
				t.Errorf("Test: %s failed\n got: %v,\nwant: %v", c.name, bal, c.want)
			}
		})
	}
}

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
				t.Errorf("Test: %s failed\ngot: %v,\nwant: %v", c.name, c.bal, c.want)
			}
		})
	}
}

func TestRoundRobinRemove(t *testing.T) {
	// Test code here
	cases := []struct {
		name string
		bal  Balancer
		args []url.URL
		want Balancer
	}{
		{
			name: "Remove one URL from empty balancer",
			bal:  &RoundRobinBalancer{},
			args: []url.URL{{Scheme: "http", Host: "localhost:8080"}},
			want: &RoundRobinBalancer{},
		},
		{
			name: "Remove URL from filled balancer",
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
					servers: []url.URL{},
				},
			},
		},
		{
			name: "Remove non-existent URL from filled balancer",
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
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, arg := range c.args {
				c.bal.Remove(arg)
			}
			if !reflect.DeepEqual(c.bal, c.want) {
				t.Errorf("Test: %s failed\n got: %v,\nwant: %v", c.name, c.bal, c.want)
			}
		})
	}
}
