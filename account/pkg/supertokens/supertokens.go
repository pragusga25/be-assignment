package supertokens

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type Config struct {
	ConnectionURI string
	APIKey        string
	AppName       string
	APIDomain     string
	WebsiteDomain string
}

func Init(config Config) error {
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: config.ConnectionURI,
			APIKey:        config.APIKey,
		},
		AppInfo: supertokens.AppInfo{
			AppName:       config.AppName,
			APIDomain:     config.APIDomain,
			WebsiteDomain: config.WebsiteDomain,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			session.Init(nil),
		},
	})

	return err
}
