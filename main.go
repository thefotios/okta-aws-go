package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/davecgh/go-spew/spew"
	"github.com/tcnksm/go-input"
	"gopkg.in/urfave/cli.v1" // imports as package "cli"
)

var ui = &input.UI{
	Writer: os.Stdout,
	Reader: os.Stdin,
}

type user_creds struct {
	username, password string
}

func getUserCreds(username string) (user_creds, error) {
	if username == "" {
		u, err := ui.Ask("Username", &input.Options{
			Required:  true,
			HideOrder: true,
		})
		if err != nil {
			return user_creds{}, cli.NewExitError("Username not provided", 86)
		}
		username = u
	}

	password, err := ui.Ask(fmt.Sprintf("Password for %s", username), &input.Options{
		Required:  true,
		HideOrder: true,
		Mask:      true,
	})

	if err != nil {
		return user_creds{}, cli.NewExitError("Password not provided", 86)
	}

	return user_creds{
		username: username,
		password: password,
	}, nil
}

type OktaAwsConfigFile struct {
	path string
}

// TODO: Can we use the SDK for this?
// https://aws.amazon.com/blogs/security/a-new-and-standardized-way-to-manage-credentials-in-the-aws-sdks/
type AWSConfigFile struct {
	path string
}

type SIDCacheFile struct {
	path string
}

type paths struct {
	okta_aws_login_config_file OktaAwsConfigFile
	aws_config_file            AWSConfigFile
	sid_cache_file             SIDCacheFile
}

func getPaths() (paths, error) {
	// file_root: Path in which all file interaction will be relative to.
	// Defaults to the users home dir.
	usr, err := user.Current()
	if err != nil {
		return paths{}, err
	}
	fileRoot := usr.HomeDir

	return paths{
		// okta_aws_login_config_file: The file were the config parameters for the
		// okta_aws_login tool is stored
		okta_aws_login_config_file: OktaAwsConfigFile{path: fileRoot + "/.okta_aws_login_config"},
		// aws_config_file: The file where this script will store the temp
		// credentials under the saml profile.
		aws_config_file: AWSConfigFile{path: fileRoot + "/.aws/credentials"},
		// sid_cache_file: The file where the Okta sid is stored.
		// only used if cache_sid is True.
		sid_cache_file: SIDCacheFile{path: fileRoot + "/.okta_sid"},
	}, nil
}

func main() {
	var username string
	var profile string

	app := cli.NewApp()
	app.Name = "okta-aws-go"
	app.Description = "Gets a STS token to use for AWS CLI based on a SAML assertion from Okta"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "username, u",
			Usage:       "The username to use when logging into Okta. If not provided you will be prompted to enter a username.",
			EnvVar:      "OKTA_USERNAME",
			Destination: &username,
		},
		cli.StringFlag{
			Name:        "profile, p",
			Usage:       "The name of the profile to use when storing the credentials in the AWS credentials file. If not provided then the name of the role assumed will be used as the profile name.",
			Destination: &profile,
		},
	}

	app.Action = func(c *cli.Context) error {
		creds, err := getUserCreds(username)
		if err != nil {
			return err
		}

		filePaths, err := getPaths()
		if err != nil {
			return err
		}

		spew.Dump(creds, filePaths)

		return nil
	}

	app.Run(os.Args)
}
