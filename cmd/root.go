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
		panic(err)
	}

	rootCmd := &cobra.Command{
		Use:   "binenv",
		Short: "Install binary distributions easily",
		Long: `binenv lets you install binary-distributed applications
(e.g. terraform, kubectl, ...) easily and switch between any version.
		
If your directory has a '.binenv.lock', proper versions will always be
selected.

This is version ` + Version + ` built on ` + BuildDate + `.`,
		// this is required since in shim mode, we have to accept any number of
		// arguments
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := initializeConfig(cmd.Root())
			if err != nil {
				return err
			}

			a.SetVerbose(verbose)

			// Set defaults or explicitely set directories
			a.SetBinDir(bindir)
			a.SetLinkDir(linkdir)
			a.SetConfigDir(confdir)
			a.SetCacheDir(cachedir)

			// Apply dir changes for global mode
			a.SetGlobal(global)

			// If some directories have been set explicitely, overwrite them
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

			err = a.Init()
			if err != nil {
				os.Exit(0)
				return err
			}

			// short circuit ShellCompNoDescRequestCmd handling
			// for binaries completion completion handling
			// (bit not for binenv)
			if (cmd.CalledAs() == cobra.ShellCompNoDescRequestCmd || cmd.CalledAs() == cobra.ShellCompRequestCmd) && !isItMe() {
				// we do not want the internal (cobra-generated)
				// __completeNoDesc to be called since this would prevent
				// shimmed binary to return its own completions
				a.Execute(os.Args)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if !isItMe() {
				a.Execute(os.Args)
			}
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	// If we can not guess those directories, we have to bailout
	dbin, err := app.GetDefaultBinDir()
	if err != nil {
		err = fmt.Errorf("unable to guess binaries directory: %v", err)
		panic(err)
	}
	dlink, err := app.GetDefaultLinkDir()
	if err != nil {
		err = fmt.Errorf("unable to guess link directory: %v", err)
		panic(err)
	}
	dcache, err := app.GetDefaultCacheDir()
	if err != nil {
		err = fmt.Errorf("unable to guess cache directory: %v", err)
		panic(err)
	}
	dconf, err := app.GetDefaultConfDir()
	if err != nil {
		err = fmt.Errorf("unable to guess conf directory: %v", err)
		panic(err)
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose operation [BINENV_VERBOSE]")
	rootCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "global mode [BINENV_GLOBAL]")

	rootCmd.PersistentFlags().StringVarP(&bindir, "bindir", "B", dbin, "binaries directory [BINENV_BINDIR]")
	rootCmd.PersistentFlags().StringVarP(&linkdir, "linkdir", "L", dlink, "link directory [BINENV_LINKDIR]")
	rootCmd.PersistentFlags().StringVarP(&cachedir, "cachedir", "K", dcache, "cache directory [BINENV_CACHEDIR]")
	rootCmd.PersistentFlags().StringVarP(&confdir, "confdir", "C", dconf, "distributions configuration directory [BINENV_CONFDIR]")

	// disable flag parsing if we're called as a shim
	if !isItMe() {
		debugCompletion("binenv called in shim mode for %q\n", strings.Join(os.Args, " "))
		rootCmd.DisableFlagParsing = true
		rootCmd.Args = cobra.ArbitraryArgs
		rootCmd.SilenceUsage = true
		rootCmd.CompletionOptions.DisableDefaultCmd = true

		// no need to add commands
		return rootCmd
	}

	debugCompletion("binenv called in binenv mode for %q\n", strings.Join(os.Args, " "))

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

func isItMe() bool {
	// get filename part from args[0], remove path
	filename := os.Args[0]
	if strings.Contains(filename, "/") {
		parts := strings.Split(filename, "/")
		filename = parts[len(parts)-1]
	}
	// if the filename is "binenv" or starts with "__debug_bin" we are in binenv mode
	// otherwise we are in shim mode
	return filename == "binenv" || // this is us
		strings.HasPrefix(filename, "__debug_bin") // for debugging in vscode
}

func isCompletionDebug() bool {
	return os.Getenv("BASH_COMP_DEBUG_FILE") != ""
}

func debugCompletion(msg string, args ...interface{}) {
	if !isCompletionDebug() {
		return
	}

	f, err := os.OpenFile(os.Getenv("BASH_COMP_DEBUG_FILE"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open bash completion debug file %s: %v", os.Getenv("BASH_COMP_DEBUG_FILE"), err)
		return
	}

	defer f.Close()

	fmt.Fprintf(f, msg, args...)
	// fmt.Fprintf(os.Stderr, "[DBG] called in completion mode for %s", os.Args[0])
}
