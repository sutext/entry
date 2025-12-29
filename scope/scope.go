package scope

import (
	"database/sql/driver"
	"slices"
	"strings"
)

type Scopes []string

const (
	Email         = "email"
	Phone         = "phone"
	Groups        = "groups"
	OpenID        = "openid"
	Address       = "address"
	Profile       = "profile"
	OfflineAccess = "offline_access"
)
const (
	crossClientPrefix = "audience:server:client_id:"
)

func All() Scopes {
	return Scopes{Email, Phone, Groups, OpenID, Address, Profile, OfflineAccess}
}
func Parse(scopes string) Scopes {
	return Scopes(strings.Fields(scopes))
}
func (s Scopes) Contains(scope string) bool {
	return slices.Contains(s, scope)
}
func (s Scopes) String() string {
	return strings.Join(s, " ")
}
func (s Scopes) Value() driver.Value {
	return s.String()
}
func (s *Scopes) Scan(src any) error {
	*s = Parse(src.(string))
	return nil
}
func (s Scopes) Validate() (hasOpenID bool, unrecognized, peerIDs []string) {
	hasOpenIDScope := false
	for _, sc := range s {
		switch sc {
		case OpenID:
			hasOpenIDScope = true
		case Email, Profile, Groups, Address, Phone, OfflineAccess:
		default:
			peerID, ok := ParseClientID(sc)
			if !ok {
				unrecognized = append(unrecognized, sc)
				continue
			}
			peerIDs = append(peerIDs, peerID)
		}
	}
	return hasOpenIDScope, unrecognized, peerIDs
}

// /gorm types
func (s Scopes) GormDataType() string {
	return "string"
}
func ParseClientID(scope string) (clientID string, ok bool) {
	if ok = strings.HasPrefix(scope, crossClientPrefix); ok {
		return scope[len(crossClientPrefix):], true
	}
	return clientID, false
}
