package ldapcon

import (
	"crypto/tls"
	"fmt"
	"sync"

	"gopkg.in/ldap.v2"
)

type LDAPConnection struct {
	conn *ldap.Conn
}

type LDAPConnManager struct {
	connections map[string]*LDAPConnection
	mu          sync.Mutex
}

var manager *LDAPConnManager
var once sync.Once

func GetLDAPConnManager() *LDAPConnManager {
	once.Do(func() {
		manager = &LDAPConnManager{
			connections: make(map[string]*LDAPConnection),
		}
	})
	return manager
}

func (m *LDAPConnManager) GetInstance(server string, credential *Credential, port int) (*LDAPConnection, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Eğer bağlantı zaten varsa, mevcut bağlantıyı döndür.
	if conn, exists := m.connections[server]; exists {
		return conn, nil
	}

	// Yeni bir LDAPS bağlantısı oluştur.
	conn, err := ldap.DialTLS("tcp", server, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}

	// Bind işlemi
	if err := conn.Bind(credential.Username, credential.Password); err != nil {
		conn.Close()
		return nil, err
	}

	// Yeni bağlantıyı sakla ve döndür.
	ldapConn := &LDAPConnection{conn: conn}
	m.connections[server] = ldapConn
	return ldapConn, nil
}

func (lc *LDAPConnection) Search(baseDN, filter string, attributes []string) (*ldap.SearchResult, error) {
	searchRequest := ldap.NewSearchRequest(
		baseDN, // Arama yapılacak base DN
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,     // LDAP filtresi
		attributes, // Geri döndürülecek öznitelikler
		nil,
	)

	// Search işlemi
	result, err := lc.conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search failed: %v", err)
	}
	return result, nil
}
