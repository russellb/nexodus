package nexodus

import (
	"errors"

	"github.com/nexodus-io/nexodus/internal/api/public"
)

const (
	persistentKeepalive = "20"
)

var (
	securityGroupErr = errors.New("nftables setup error")
)

func (nx *Nexodus) DeployWireguardConfig(updatedPeers map[string]public.ModelsDevice) error {
	cfg := &wgConfig{
		Interface: nx.wgConfig.Interface,
		Peers:     nx.wgConfig.Peers,
	}

	if nx.TunnelIP != nx.getIPv4Iface(nx.tunnelIface).String() {
		if err := nx.setupInterface(); err != nil {
			return err
		}
	}

	// add routes and tunnels for the new peers only according to the cache diff
	for _, updatedPeer := range updatedPeers {
		if updatedPeer.Id == "" {
			continue
		}
		// add routes for each peer candidate (unless the key matches the local nodes key)
		peer, ok := cfg.Peers[updatedPeer.PublicKey]
		if !ok || peer.PublicKey == nx.wireguardPubKey {
			continue
		}
		if err := nx.handlePeerRoute(peer); err != nil {
			return err
		}
		if err := nx.handlePeerTunnel(peer); err != nil {
			return err
		}
	}

	nx.logger.Debug("Peer setup complete")
	return nil
}
