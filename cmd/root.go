package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// The environment variable prefix of all environment variables bound to our command line flags.
	envPrefix = "BINENV"
)

// RootCmd returns the root cobra command
func RootCmd() *cobra.Command {
	var (
		bindir, linkdir, cachedir, confdir string
		global, verbose                    bool
	)

	a, err := app.New()
	if err != nil {
		fmt.Printf("got error %v\n", err)
		panic(err)
	}

	rootCmd := &cobra.Command{
		Use:   "binenv",
		Short: "Install binary distributions easily",
		Long: `binenv lets you install binary-distributed applications
(e.g. terraform, kubectl, ...) easily and switch between any version.
		
If your directory has a '.binenv.lock', proper versions will always be
selected.`,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := initializeConfig(cmd.Root())

			a.SetVerbose(verbose)
			a.SetGlobal(global)

			// if options has been changed by flag or env, we apply it
			// otherwise we keep preceding setting
			if cmd.Root().PersistentFlags().Lookup("bindir").Changed {
				a.SetBinDir(bindir)
			}
			// if options has been changed by flag or env, we apply it
			// otherwise we keep preceding setting
			if cmd.Root().PersistentFlags().Lookup("linkdir").Changed {
				a.SetLinkDir(linkdir)
			}
			// if options has been changed by flag or env, we apply it
			// otherwise we keep preceding setting
			if cmd.Root().PersistentFlags().Lookup("confdir").Changed {
				a.SetConfigDir(confdir)
			}
			// if options has been changed by flag or env, we apply it
			// otherwise we keep preceding setting
			if cmd.Root().PersistentFlags().Lookup("cachedir").Changed {
				a.SetCacheDir(cachedir)
			}

			a.DumpConfig()

			return err
		},
		Run: func(cmd *cobra.Command, args []string) {
			if !strings.HasSuffix(os.Args[0], "binenv") {
				a.Execute(os.Args)
			}
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose operation [BINENV_VERBOSE]")
	rootCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "global mode [BINENV_GLOBAL]")

	rootCmd.PersistentFlags().StringVarP(&bindir, "bindir", "B", app.GetDefaultBinDir(), "binaries directory [BINENV_BINDIR]")
	rootCmd.PersistentFlags().StringVarP(&linkdir, "linkdir", "L", app.GetDefaultLinkDir(), "link directory [BINENV_LINKDIR]")
	rootCmd.PersistentFlags().StringVarP(&cachedir, "cachedir", "K", app.GetDefaultCacheDir(), "cache directory [BINENV_CACHEDIR]")
	rootCmd.PersistentFlags().StringVarP(&confdir, "confdir", "C", app.GetDefaultConfDir(), "distributions configuration directory [BINENV_CONFDIR]")

	rootCmd.AddCommand(
		completionCmd(),
		expandCmd(a),
		installCmd(a),
		localCmd(a),
		searchCmd(a),
		uninstallCmd(a),
		updateCmd(a),
		versionCmd(),
		versionsCmd(a),
		upgradeCmd(a),
	)

	return rootCmd
}

func truthify(s string) bool {
	s = strings.ToLower(s)
	// trueness suggestions courtesy of Github co-pilot
	return s == "true" || s == "1" || s == "yes" || s == "y" || s == "on" || s == "enable" || s == "enabled" || s == "active"
}

// func envOrDefault(k, def string) string {
// 	if v := os.Getenv(k); v != "" {
// 		return v
// 	}
// 	return def
// }

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix("BINENV")

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.PersistentFlags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
