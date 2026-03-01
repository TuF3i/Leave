package oauth

import (
	"leave/core/pkg/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type OAuthCore struct {
	conf      *config.Config
	oauthConf *oauth2.Config
}

func GetOAuth2(conf *config.Config) *OAuthCore {
	o := &OAuthCore{
		conf: conf,
		oauthConf: &oauth2.Config{
			ClientID:     conf.OAuth.ClientID,
			ClientSecret: conf.OAuth.ClientSecret,
			Endpoint:     github.Endpoint,
			RedirectURL:  conf.OAuth.RedirectUrl,
			Scopes:       []string{"user:email"},
		},
	}

	return o
}
