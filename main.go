package main

import (
	"fmt"
	"log"
	"sync"

	"ldapcon/ldapcon"
)

func main() {

	baseDN := "dc=example,dc=com"
	filter := "(cn=John Doe)"
	attributes := []string{"dn", "cn", "mail"}
	users := []string{"user1", "user2", "user3", "userN"}

	credential := ldapcon.NewCredential()
	fmt.Println(credential.Username)

	var domainname string
	fmt.Print("DC Domain name: ")
	fmt.Scanln(&domainname)
	domain := ldapcon.ADdomain(domainname)
	dcServers, err := domain.DiscoveryDCs()
	if err != nil {
		panic("DC ler tespit edilemedi")
	}
	fmt.Println(dcServers)

	manager := ldapcon.GetLDAPConnManager()
	var ldapConnections []*ldapcon.LDAPConnection

	for _, server := range dcServers {
		templdapconnection, err := manager.GetInstance(server, credential, 636)
		if templdapconnection != nil && err == nil {
			ldapConnections = append(ldapConnections, templdapconnection)
		}
	}

	usersPerServer := len(users) / len(ldapConnections)

	var wg sync.WaitGroup
	for i, conn := range ldapConnections {
		startIndex := i * usersPerServer
		endIndex := startIndex + usersPerServer

		if i == len(ldapConnections)-1 {
			endIndex = len(users)
		}

		serverUsers := users[startIndex:endIndex]

		wg.Add(1)
		go func(conn *ldapcon.LDAPConnection, users []string, baseDN string, filter string, attributes []string, wg *sync.WaitGroup) {
			defer wg.Done()
			processUsers(conn, users, baseDN, filter, attributes, wg)
		}(conn, serverUsers, baseDN, filter, attributes, &wg)
	}

	wg.Wait()

}

func processUsers(conn *ldapcon.LDAPConnection, users []string, basedn string, filter string, attributes []string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Kullanıcılar için basit bir işlem
	for _, user := range users {
		fmt.Printf("Processing %s on %v\n", user, conn)

		result, err := conn.Search(basedn, filter, attributes) // Doğru şekilde Search çağrısı
		if err != nil {
			log.Printf("Arama hatası %s: %v", user, err)
			continue
		}
		// Eğer arama sonucu varsa, sonuçları işle
		for _, entry := range result.Entries {
			fmt.Printf("Kullanıcı bulundu: DN: %s, CN: %s\n", entry.DN, entry.GetAttributeValue("cn"))
		}

	}
}
