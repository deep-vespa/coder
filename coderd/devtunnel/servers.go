package devtunnel

import (
	"runtime"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"

	"github.com/coder/coder/cryptorand"
)

type Region struct {
	ID           int
	LocationName string
	Nodes        []Node
}

type Node struct {
	ID                int    `json:"id"`
	RegionID          int    `json:"region_id"`
	HostnameHTTP      string `json:"hostname_https"`
	Secure            bool   `json:"secure"` // Whether to use HTTPS.
	HostnameWireguard string `json:"hostname_wireguard"`
	WireguardPort     uint16 `json:"wireguard_port"`

	AvgLatency time.Duration `json:"avg_latency"`
}

var Regions = []Region{
	// {
	// 	ID:           0,
	// 	LocationName: "US East Pittsburgh",
	// 	Nodes: []Node{
	// 		{
	// 			ID:                1,
	// 			RegionID:          0,
	// 			HostnameHTTP:      "pit-1.try.coder.app",
	// 			Secure:            true,
	// 			HostnameWireguard: "pit-1.try.coder.app",
	// 			WireguardPort:     55551,
	// 		},
	// 	},
	// },
	{
		ID:           0,
		LocationName: "Local",
		Nodes: []Node{
			{
				ID:                1,
				RegionID:          0,
				HostnameHTTP:      "local.try.coder.app:8080",
				HostnameWireguard: "local.try.coder.app",
				WireguardPort:     55555,
			},
		},
	},
}

func FindClosestNode() (Node, error) {
	nodes := []Node{}

	for _, region := range Regions {
		// Pick a random node from each region.
		i, err := cryptorand.Intn(len(region.Nodes))
		if err != nil {
			return Node{}, err
		}
		nodes = append(nodes, region.Nodes[i])
	}

	var (
		nodesMu sync.Mutex
		eg      = errgroup.Group{}
	)
	for i, node := range nodes {
		i, node := i, node
		eg.Go(func() error {
			pinger, err := ping.NewPinger(node.HostnameHTTP)
			if err != nil {
				return err
			}

			if runtime.GOOS == "windows" {
				pinger.SetPrivileged(true)
			}

			pinger.Count = 5
			err = pinger.Run()
			if err != nil {
				return err
			}

			nodesMu.Lock()
			nodes[i].AvgLatency = pinger.Statistics().AvgRtt
			nodesMu.Unlock()
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return Node{}, err
	}

	slices.SortFunc(nodes, func(i, j Node) bool {
		return i.AvgLatency < j.AvgLatency
	})
	return nodes[0], nil
}
