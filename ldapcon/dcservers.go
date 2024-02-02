package ldapcon

import "net"

type DCServers struct {
	Domain string
}

func ADdomain(Domain string) *DCServers {
	return &DCServers{Domain: Domain}
}

func (ad *DCServers) DiscoveryDCs() ([]string, error) {
	_, srvRecords, err := net.LookupSRV("ldap", "tcp", "dc._msdcs."+ad.Domain)
	if err != nil {
		return nil, err
	}

	var dcs []string
	for _, server := range srvRecords {
		dcs = append(dcs, server.Target)
	}

	return dcs, nil

}
