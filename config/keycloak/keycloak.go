package config

import "github.com/Nerzal/gocloak/v12"

type keycloak struct {
	gocloak      *gocloak.GoCloak // keycloak client
	clientId     string           // clientId specified in Keycloak
	clientSecret string           // client secret specified in Keycloak
	realm        string           // realm specified in Keycloak
}

func newKeycloak() *keycloak {
	return &keycloak{
		gocloak:      gocloak.NewClient("http://localhost:8080"),
		clientId:     "rehab-go",
		clientSecret: "h5rQCqL7dgofp4OdCBZOFUVBxWJXRNBC",
		realm:        "Rehab",
	}
}
