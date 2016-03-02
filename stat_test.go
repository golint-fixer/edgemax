package edgemax

import (
	"encoding/json"
	"net"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestSystemStats(t *testing.T) {
	var tests = []struct {
		desc    string
		b       []byte
		errType reflect.Type
		s       *SystemStats
	}{
		{
			desc:    "invalid JSON",
			b:       []byte(`foo`),
			errType: reflect.TypeOf(&json.SyntaxError{}),
		},
		{
			desc:    "invalid CPU integer",
			b:       []byte(`{"cpu":"foo"}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid uptime integer",
			b:       []byte(`{"cpu":"0","uptime":"foo"}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid memory integer",
			b:       []byte(`{"cpu":"0","uptime":"1","mem":"foo"}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc: "OK",
			b:    []byte(`{"cpu":"10","uptime":"20","mem":"30"}`),
			s: &SystemStats{
				CPU:    10,
				Uptime: 20 * time.Second,
				Memory: 30,
			},
		},
	}

	for i, tt := range tests {
		t.Logf("[%02d] test %q", i, tt.desc)

		s := new(SystemStats)
		err := s.UnmarshalJSON(tt.b)

		if want, got := tt.errType, reflect.TypeOf(err); !reflect.DeepEqual(want, got) {
			t.Fatalf("unexpected error type:\n- want: %v\n-  got: %v", want, got)
		}
		if err != nil {
			continue
		}

		if want, got := tt.s, s; !reflect.DeepEqual(want, got) {
			t.Fatalf("unexpected SystemStats:\n- want: %v\n-  got: %v", want, got)
		}
	}
}

func TestInterfaces(t *testing.T) {
	var tests = []struct {
		desc    string
		b       []byte
		errType reflect.Type
		ifis    Interfaces
	}{
		{
			desc:    "invalid JSON",
			b:       []byte(`foo`),
			errType: reflect.TypeOf(&json.SyntaxError{}),
		},
		{
			desc:    "invalid speed",
			b:       []byte(`{"eth0":{"speed":"foo"}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid MTU",
			b:       []byte(`{"eth0":{"mtu":"foo"}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid receive packets",
			b:       []byte(`{"eth0":{"stats":{"rx_packets":"foo"}}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid transmit packets",
			b:       []byte(`{"eth0":{"stats":{"tx_packets":"foo"}}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid receive bytes",
			b:       []byte(`{"eth0":{"stats":{"rx_bytes":"foo"}}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid transmit bytes",
			b:       []byte(`{"eth0":{"stats":{"tx_bytes":"foo"}}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid receive errors",
			b:       []byte(`{"eth0":{"stats":{"rx_errors":"foo"}}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid transmit errors",
			b:       []byte(`{"eth0":{"stats":{"tx_errors":"foo"}}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid receive dropped",
			b:       []byte(`{"eth0":{"stats":{"rx_dropped":"foo"}}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid transmit dropped",
			b:       []byte(`{"eth0":{"stats":{"tx_dropped":"foo"}}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid multicast",
			b:       []byte(`{"eth0":{"stats":{"multicast":"foo"}}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid receive bps",
			b:       []byte(`{"eth0":{"stats":{"rx_bps":"foo"}}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid transmit bps",
			b:       []byte(`{"eth0":{"stats":{"tx_bps":"foo"}}}`),
			errType: reflect.TypeOf(&strconv.NumError{}),
		},
		{
			desc:    "invalid MAC",
			b:       []byte(`{"eth0":{"mac":"foo"}}`),
			errType: reflect.TypeOf(&net.AddrError{}),
		},
		{
			desc:    "invalid CIDR IP",
			b:       []byte(`{"eth0":{"addresses":["foo"]}}`),
			errType: reflect.TypeOf(&net.ParseError{}),
		},
		{
			desc: "OK one interface",
			b:    []byte(`{"eth0":{"up":"true","autoneg":"true","duplex":"full","speed":"10","mac":"de:ad:be:ef:de:ad","mtu":"1500","addresses":["192.168.1.1/24"],"stats":{"rx_packets":"1","tx_packets":"2","rx_bytes":"3","tx_bytes":"4","rx_errors":"5","tx_errors":"6","rx_dropped":"7","tx_dropped":"8","multicast":"9","rx_bps":"10","tx_bps":"11"}}}`),
			ifis: Interfaces{{
				Name:            "eth0",
				Up:              true,
				Autonegotiation: true,
				Duplex:          "full",
				Speed:           10,
				MAC:             net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
				MTU:             1500,
				Addresses:       []net.IP{net.IPv4(192, 168, 1, 1)},
				Stats: InterfaceStats{
					ReceivePackets:  1,
					TransmitPackets: 2,
					ReceiveBytes:    3,
					TransmitBytes:   4,
					ReceiveErrors:   5,
					TransmitErrors:  6,
					ReceiveDropped:  7,
					TransmitDropped: 8,
					Multicast:       9,
					ReceiveBPS:      10,
					TransmitBPS:     11,
				},
			}},
		},
		{
			desc: "OK two interfaces",
			b:    []byte(`{"eth1":{"mac":"ab:ad:1d:ea:ab:ad","addresses":["192.168.1.2/24"]},"eth0":{"mac":"de:ad:be:ef:de:ad","addresses":["192.168.1.1/24"]}}`),
			ifis: Interfaces{
				{
					Name:      "eth0",
					MAC:       net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
					Addresses: []net.IP{net.IPv4(192, 168, 1, 1)},
				},
				{
					Name:      "eth1",
					MAC:       net.HardwareAddr{0xab, 0xad, 0x1d, 0xea, 0xab, 0xad},
					Addresses: []net.IP{net.IPv4(192, 168, 1, 2)},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Logf("[%02d] test %q", i, tt.desc)

		var ifis Interfaces
		err := ifis.UnmarshalJSON(tt.b)

		if want, got := tt.errType, reflect.TypeOf(err); !reflect.DeepEqual(want, got) {
			t.Fatalf("unexpected error type:\n- want: %v\n-  got: %v", want, got)
		}
		if err != nil {
			continue
		}

		if want, got := tt.ifis, ifis; !reflect.DeepEqual(want, got) {
			t.Fatalf("unexpected Interfaces:\n- want: %+v\n-  got: %+v", want, got)
		}
	}
}

func Test_ipLess(t *testing.T) {
	var tests = []struct {
		a    net.IP
		b    net.IP
		less bool
	}{
		{
			less: false,
		},
		{
			a:    net.ParseIP("10.0.0.1"),
			b:    net.ParseIP("10.0.0.2"),
			less: true,
		},
		{
			a:    net.ParseIP("10.0.0.2"),
			b:    net.ParseIP("10.0.0.1"),
			less: false,
		},
		{
			a:    net.ParseIP("10.0.0.1"),
			b:    net.ParseIP("10.1.0.1"),
			less: true,
		},
		{
			a:    net.ParseIP("10.1.0.1"),
			b:    net.ParseIP("10.0.0.1"),
			less: false,
		},
		{
			a:    net.ParseIP("2001:db8::1"),
			b:    net.ParseIP("10.0.0.1"),
			less: false,
		},
		{
			a:    net.ParseIP("10.0.0.1"),
			b:    net.ParseIP("2001:db8::1"),
			less: true,
		},
		{
			a:    net.ParseIP("2001:db8::1"),
			b:    net.ParseIP("2001:db8::2"),
			less: true,
		},
		{
			a:    net.ParseIP("2001:db8::2"),
			b:    net.ParseIP("2001:db8::1"),
			less: false,
		},
		{
			a:    net.ParseIP("2001:db9::2"),
			b:    net.ParseIP("2001:db8::1"),
			less: false,
		},
	}

	for i, tt := range tests {
		t.Logf("[%02d] test %q < %q?", i, tt.a, tt.b)

		less := ipLess(tt.a, tt.b)
		if want, got := tt.less, less; !reflect.DeepEqual(want, got) {
			t.Fatalf("unexpected ipLess(%q, %q) result:\n- want: %v\n-  got: %v",
				tt.a, tt.b, want, got)
		}
	}
}
