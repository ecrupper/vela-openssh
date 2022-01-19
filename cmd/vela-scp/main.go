package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/go-vela/vela-openssh/internal/openssh"
	"github.com/go-vela/vela-openssh/internal/scp"
	"github.com/go-vela/vela-openssh/pkg/binarywrapper"
)

func main() {
	app := &cli.App{
		Name:      "vela-scp",
		Usage:     "Vela plugin wrapping the scp binary.",
		Copyright: "Copyright (c) 2022 Target Brands, Inc. All rights reserved.",
		Authors: []*cli.Author{
			{
				Name:  "Vela Admins",
				Email: "vela@target.com",
			},
		},
		Action:  run,
		Version: openssh.PluginVersion.Semantic(),
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "source",
				Usage:    "source parameter for scp (see manual 'man scp')",
				EnvVars:  []string{"PARAMETER_SOURCE", "VELA_SOURCE"},
				FilePath: "/vela/parameters/vela-scp/source,/vela/secrets/vela-scp/source",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "target",
				Usage:    "target parameter for scp (see manual 'man scp')",
				EnvVars:  []string{"PARAMETER_TARGET", "VELA_TARGET"},
				FilePath: "/vela/parameters/vela-scp/target,/vela/secrets/vela-scp/target",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "identity-file.path",
				Usage:    "path to the identity file parameter for scp (see manual 'man scp')",
				EnvVars:  []string{"PARAMETER_IDENTITY_FILE_PATH", "VELA_IDENTITY_FILE_PATH", "PARAMETER_SSH_KEY_PATH", "VELA_SSH_KEY_PATH"},
				FilePath: "/vela/parameters/vela-scp/identity-file.path,/vela/secrets/vela-scp/identity-file.path",
			},
			&cli.StringFlag{
				Name:     "identity-file.contents",
				Usage:    "contents of the identity-file (not the path, the real deal)",
				EnvVars:  []string{"PARAMETER_IDENTITY_FILE_CONTENTS", "VELA_IDENTITY_FILE_CONTENTS", "PARAMETER_SSH_KEY", "VELA_SSH_KEY"},
				FilePath: "/vela/parameters/vela-scp/identity-file.contents,/vela/secrets/vela-scp/identity-file.contents",
			},
			&cli.StringSliceFlag{
				Name:     "scp.flag",
				Usage:    "any additional flags for scp can be specified here",
				EnvVars:  []string{"PARAMETER_SCP_FLAG", "VELA_SCP_FLAG"},
				FilePath: "/vela/parameters/vela-scp/scp.flag,/vela/secrets/vela-scp/scp.flag",
			},
			&cli.StringFlag{
				Name:     "sshpass.password",
				Usage:    "password for use with destination target (used with sshpass)",
				EnvVars:  []string{"PARAMETER_SSHPASS_PASSWORD", "PARAMETER_PASSWORD", "VELA_SSHPASS_PASSWORD", "VELA_PASSWORD"},
				FilePath: "/vela/parameters/vela-scp/sshpass.password,/vela/parameters/vela-scp/password,/vela/secrets/vela-scp/sshpass.password,/vela/secrets/vela-scp/password",
			},
			&cli.StringFlag{
				Name:     "sshpass.passphrase",
				Usage:    "passphrase for use with identity file (used with sshpass)",
				EnvVars:  []string{"PARAMETER_SSHPASS_PASSPHRASE", "VELA_SSHPASS_PASSPHRASE"},
				FilePath: "/vela/parameters/vela-scp/sshpass.passphrase,/vela/parameters/vela-scp/passphrase,/vela/secrets/vela-scp/sshpass.passphrase,/vela/secrets/vela-scp/passphrase",
			},
			&cli.StringSliceFlag{
				Name:     "sshpass.flag",
				Usage:    "any additional flags for sshpass can be specified here)",
				EnvVars:  []string{"PARAMETER_SSHPASS_FLAG", "VELA_SSHPASS_FLAG"},
				FilePath: "/vela/parameters/vela-scp/sshpass.flag,/vela/secrets/vela-scp/sshpass.flag",
			},
			&cli.StringFlag{
				Name:     "ci",
				Usage:    "set the CI environment (if $CI is set output tries to be friendlier)",
				EnvVars:  []string{"PARAMETER_CI", "CI"},
				FilePath: "/vela/parameters/vela-scp/ci,/vela/secrets/vela-scp/ci",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.IsSet("ci") {
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: false,
			PadLevelText:  true,
		})
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-openssh",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/scp",
		"registry": "https://hub.docker.com/r/target/vela-scp",
		"version":  openssh.PluginVersion.Semantic(),
		"commit":   openssh.PluginVersion.Metadata.GitCommit,
	}).Info("Vela SCP Plugin")

	bp := binarywrapper.Plugin{
		PluginConfig: &scp.Config{
			Source:               c.StringSlice("source"),
			Target:               c.String("target"),
			IdentityFilePath:     c.StringSlice("identity-file.path"),
			IdentityFileContents: c.String("identity-file.contents"),
			SCPFlags:             c.StringSlice("scp.flag"),
			SSHPassword:          c.String("sshpass.password"),
			SSHPassphrase:        c.String("sshpass.passphrase"),
			SSHPASSFlags:         c.StringSlice("sshpass.flag"),
		},
	}

	return bp.Exec()
}
