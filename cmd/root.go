package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tphakala/birdnet-go/cmd/authors"

	//"github.com/tphakala/birdnet-go/cmd/directory"
	//"github.com/tphakala/birdnet-go/cmd/file"
	"github.com/tphakala/birdnet-go/cmd/license"
	"github.com/tphakala/birdnet-go/cmd/realtime"
	"github.com/tphakala/birdnet-go/internal/conf"
)

// RootCommand creates and returns the root command
func RootCommand(settings *conf.Settings) *cobra.Command {
	//ctx := config.GetGlobalContext()

	rootCmd := &cobra.Command{
		Use:   "birdnet",
		Short: "BirdNET-Go CLI",
	}

	// Set up the global flags for the root command.
	setupFlags(rootCmd, settings)

	// Add sub-commands to the root command.
	//fileCmd := file.Command(ctx)
	//directoryCmd := directory.Command(ctx)
	realtimeCmd := realtime.Command(settings)
	authorsCmd := authors.Command()
	licenseCmd := license.Command()

	subcommands := []*cobra.Command{
		//fileCmd,
		//directoryCmd,
		realtimeCmd,
		authorsCmd,
		licenseCmd,
	}

	rootCmd.AddCommand(subcommands...)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Skip setup for authors and license commands
		if cmd.Name() != authorsCmd.Name() && cmd.Name() != licenseCmd.Name() {
			if err := initialize(); err != nil {
				return fmt.Errorf("error initializing: %v", err)
			}
		}

		return nil
	}

	return rootCmd
}

// initialize is called before any subcommands are run, but after the context is ready
// This function is responsible for setting up configurations, ensuring the environment is ready, etc.
func initialize() error {
	return nil
}

// defineGlobalFlags defines flags that are global to the command line interface
func setupFlags(rootCmd *cobra.Command, settings *conf.Settings) error {
	rootCmd.PersistentFlags().BoolVarP(&settings.Debug, "debug", "d", viper.GetBool("debug"), "Enable debug output")
	rootCmd.PersistentFlags().StringVar(&settings.BirdNET.Locale, "locale", viper.GetString("node.locale"), "Set the locale for labels. Accepts full name or 2-letter code.")
	rootCmd.PersistentFlags().IntVarP(&settings.BirdNET.Threads, "threads", "j", viper.GetInt("node.threads"), "Number of CPU threads to use for analysis (default 0 which is all CPUs)")
	rootCmd.PersistentFlags().Float64VarP(&settings.BirdNET.Sensitivity, "sensitivity", "s", viper.GetFloat64("birdnet.sensitivity"), "Sigmoid sensitivity value between 0.0 and 1.5")
	rootCmd.PersistentFlags().Float64VarP(&settings.BirdNET.Threshold, "threshold", "t", viper.GetFloat64("birdnet.threshold"), "Confidency threshold for detections, value between 0.1 to 1.0")
	rootCmd.PersistentFlags().Float64Var(&settings.BirdNET.Overlap, "overlap", viper.GetFloat64("birdnet.overlap"), "Overlap value between 0.0 and 2.9")
	rootCmd.PersistentFlags().Float64Var(&settings.BirdNET.Latitude, "latitude", viper.GetFloat64("birdnet.latitude"), "Latitude for species prediction")
	rootCmd.PersistentFlags().Float64Var(&settings.BirdNET.Longitude, "longitude", viper.GetFloat64("birdnet.longitude"), "Longitude for species prediction")

	// Bind flags to the viper settings
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		return fmt.Errorf("error binding flags: %v", err)
	}

	return nil
}
