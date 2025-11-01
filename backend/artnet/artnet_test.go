package artnet

import (
	"backend/config"
	"encoding/json"
	"log/slog"
	"net"
	"testing"
)

func TestArtnet_Initialize(t *testing.T) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		t.Error(err)
	}
	targets := []net.Addr{}
	targetsIpv4 := []net.Addr{}
	for _, addr := range addrs {
		ip := addr.(*net.IPNet)
		if ip.IP.IsGlobalUnicast() {
			targets = append(targets, addr)
		}
		if ip.IP.To4() != nil {
			targetsIpv4 = append(targetsIpv4, addr)
		}
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		log    *slog.Logger
		config *config.Config
		want   bool
	}{
		struct {
			name   string
			log    *slog.Logger
			config *config.Config
			want   bool
		}{
			name: "localhost",
			log:  slog.New(slog.NewJSONHandler(t.Output(), nil)),
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{"artnet"},
					Artnet: config.Artnet{
						Address:     "127.0.0.1",
						Universe:    0,
						SubUniverse: 0,
						Net:         0,
					},
				},
			},
			want: true,
		},
		struct {
			name   string
			log    *slog.Logger
			config *config.Config
			want   bool
		}{
			name: "empty address",
			log:  slog.New(slog.NewJSONHandler(t.Output(), nil)),
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{"artnet"},
					Artnet: config.Artnet{
						Address:     "",
						Universe:    0,
						SubUniverse: 0,
						Net:         0,
					},
				},
			},
			want: false,
		},
		struct {
			name   string
			log    *slog.Logger
			config *config.Config
			want   bool
		}{
			name: "outrange test",
			log:  slog.New(slog.NewJSONHandler(t.Output(), nil)),
			config: &config.Config{
				Output: config.OutputTargets{
					Target: []string{"artnet"},
					Artnet: config.Artnet{
						Address:     "254.254.254.254",
						Universe:    0,
						SubUniverse: 0,
						Net:         0,
					},
				},
			},
			want: false,
		},
	}
	for _, target := range targets {
		tests = append(tests,
			struct {
				name   string
				log    *slog.Logger
				config *config.Config
				want   bool
			}{
				name: "Can only use IPv4",
				log:  slog.New(slog.NewJSONHandler(t.Output(), nil)),
				config: &config.Config{
					Output: config.OutputTargets{
						Target: []string{"artnet"},
						Artnet: config.Artnet{
							Address:     target.(*net.IPNet).IP.String(),
							Universe:    0,
							SubUniverse: 0,
							Net:         0,
						},
					},
				},
				want: target.(*net.IPNet).IP.To4() != nil,
			})
	}
	for _, target := range targetsIpv4 {
		outrange := target.(*net.IPNet).IP
		outrange[3] += 1
		tests = append(tests,
			struct {
				name   string
				log    *slog.Logger
				config *config.Config
				want   bool
			}{
				name: "Outrange",
				log:  slog.New(slog.NewJSONHandler(t.Output(), nil)),
				config: &config.Config{
					Output: config.OutputTargets{
						Target: []string{"artnet"},
						Artnet: config.Artnet{
							Address:     outrange.String(),
							Universe:    0,
							SubUniverse: 0,
							Net:         0,
						},
					},
				},
				want: false,
			})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			encodedJson, _ := json.Marshal(tt.config)
			t.Log("data: " + string(encodedJson))
			var a Artnet = Artnet{
				TargetAddr: tt.config.Output.Artnet.Address,
			}
			got := a.Initialize(tt.log, tt.config)
			t.Logf("result: %v, want %v", got, tt.want)
			if tt.want != got {
				t.Errorf("Initialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArtnet_listen(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var a Artnet
			a.listen()
		})
	}
}

func TestArtnet_Start(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var a Artnet
			got := a.Start()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Start() = %v, want %v", got, tt.want)
			}
		})
	}
}
