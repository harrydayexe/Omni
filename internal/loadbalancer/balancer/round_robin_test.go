package balancer

import (
	"net/url"
	"reflect"
	"sync/atomic"
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

func TestRoundRobinBalance(t *testing.T) {
	type expected struct {
		isErrorExpected bool
		want            url.URL
		error           error
	}
	cases := []struct {
		name string
		bal  Balancer
		reqs []expected
	}{
		{
			name: "Balance with no URL",
			bal: &RoundRobinBalancer{
				BaseBalancer: BaseBalancer{
					servers: []url.URL{},
				},
			},
			reqs: []expected{
				{
					isErrorExpected: true,
					want:            url.URL{},
					error:           NoHealthyHostsError,
				},
			},
		},
		{
			name: "Balance with one URL",
			bal: &RoundRobinBalancer{
				BaseBalancer: BaseBalancer{
					servers: []url.URL{
						{Scheme: "http", Host: "localhost:8080"},
					},
				},
			},
			reqs: []expected{
				{
					isErrorExpected: false,
					want:            url.URL{Scheme: "http", Host: "localhost:8080"},
					error:           nil,
				},
			},
		},
		{
			name: "Balance with two URLs",
			bal: &RoundRobinBalancer{
				BaseBalancer: BaseBalancer{
					servers: []url.URL{
						{Scheme: "http", Host: "localhost:4040"},
						{Scheme: "http", Host: "localhost:8080"},
					},
				},
				current: atomic.Uint64{},
			},
			reqs: []expected{
				{
					isErrorExpected: false,
					want:            url.URL{Scheme: "http", Host: "localhost:8080"},
					error:           nil,
				},
				{
					isErrorExpected: false,
					want:            url.URL{Scheme: "http", Host: "localhost:4040"},
					error:           nil,
				},
				{
					isErrorExpected: false,
					want:            url.URL{Scheme: "http", Host: "localhost:8080"},
					error:           nil,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, req := range c.reqs {
				got, err := c.bal.Balance()
				if req.isErrorExpected && err != req.error {
					t.Errorf("Test: %s failed\n got error: %v\nwant error: %v", c.name, err, req.error)
				}
				if !req.isErrorExpected && !reflect.DeepEqual(got, req.want) {
					t.Errorf("Test: %s failed\n got: %v,\nwant: %v", c.name, got, req.want)
				}
			}
		})
	}
}
