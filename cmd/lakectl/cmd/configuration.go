package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// configuration is the user-visible configuration structure in Golang form.  When editing
// make sure *all* fields have a `mapstructure:"..."` tag, to simplify future refactoring.
type configuration struct {
	Credentials struct {
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
	}
	Server struct {
		EndpointURL string `mapstructure:"endpoint_url"`
	}
	Metastore struct {
		Type string `mapstructure:"type"`
		Hive struct {
			URI            string `mapstructure:"uri"`
			DBLocationtURI string `mapstructure:"db_location_uri"`
		} `mapstructure:"hive"`
		Glue struct {
			// TODO(ariels): Refactor credentials to share with server side.
			Profile         string `mapstructure:"profile"`
			CredentialsFile string `mapstructure:"credentials_file"`
			DBLocationtURI  string `mapstructure:"db_location_uri"`
			Credentials     *struct {
				AccessKeyID     string `mapstructure:"access_key_id"`
				AccessSecretKey string `mapstructure:"access_secret_key"`
				SessionToken    string `mapstructure:"session_token"`
			} `mapstructure:"credentials"`

			Region    string `mapstructure:"region"`
			CatalogID string `mapstructure:"catalog_id"`
		}
	}
}

func (c *configuration) GetMetastoreAwsConfig() *aws.Config {
	cfg := &aws.Config{
		Region: aws.String(c.Metastore.Glue.Region),
	}
	if c.Metastore.Glue.Profile != "" || c.Metastore.Glue.CredentialsFile != "" {
		cfg.Credentials = credentials.NewSharedCredentials(
			c.Metastore.Glue.CredentialsFile,
			c.Metastore.Glue.Profile,
		)
	}
	if c.Metastore.Glue.Credentials != nil {
		cfg.Credentials = credentials.NewStaticCredentials(
			c.Metastore.Glue.Credentials.AccessKeyID,
			c.Metastore.Glue.Credentials.AccessSecretKey,
			c.Metastore.Glue.Credentials.SessionToken,
		)
	}
	return cfg
}

func (c *configuration) GetMetastoreHiveURI() string {
	return c.Metastore.Hive.URI
}

func (c *configuration) GetMetastoreGlueCatalogID() string {
	return c.Metastore.Glue.CatalogID
}
func (c *configuration) GetMetastoreType() string {
	return c.Metastore.Type
}

func (c *configuration) GetHiveDBLocationURI() string {
	return c.Metastore.Hive.DBLocationtURI
}

func (c *configuration) GetGlueDBLocationURI() string {
	return c.Metastore.Glue.DBLocationtURI
}
