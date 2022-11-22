# go-ldap-client

Simple ldap client to authenticate, retrieve basic information and groups for a user.

## Usage

The only external dependency is [gopkg.in/ldap.v2](http://gopkg.in/ldap.v2).

```
package main

import (
	"log"

	ldap "github.com/colynn/go-ldap-client"
)

func main() {
	client := &ldap.Client{
		Base:         "dc=example,dc=com",
		Host:         "ldap.example.com",
		Port:         389,
		UseSSL:       false,
		BindDN:       "uid=readonlysuer,ou=People,dc=example,dc=com",
		BindPassword: "readonlypassword",
		UserFilter:   "(uid=%s)",
		GroupFilter: "(memberUid=%s)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}
	// It is the responsibility of the caller to close the connection
	defer client.Close()

	ok, user, err := client.Authenticate("username", "password")
	if err != nil {
		log.Fatalf("Error authenticating user %s: %+v", "username", err)
	}
	if !ok {
		log.Fatalf("Authenticating failed for user %s", "username")
	}
	log.Printf("User: %+v", user)
	
	groups, err := client.GetGroupsOfUser("username")
	if err != nil {
		log.Fatalf("Error getting groups for user %s: %+v", "username", err)
	}
	log.Printf("Groups: %+v", groups) 
}
```

## SSL(LDAPS)
If you use SSL, you will need to pass the server name for certificate verification or skip domain name verification e.g.`client.ServerName = "ldap.example.com"`.

## Why?
Because [go-ldap-client](https://github.com/jtblin/go-ldap-client) been a long time didn't maintenance from 2017 to now.
So re-create it, make it better for everyone to use and maintain.

## Later
we plan to create `go-ldap-client` base on [gopkg.in/ldap.v3](http://gopkg.in/ldap.v3).
