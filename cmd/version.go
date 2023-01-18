package cmd

import "runtime/debug"

var version string

func init() {
	if cliPlugin == "true" {
		return
	}

	if version == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			version = info.Main.Version
		}
	}

	rootCmd.Version = version
}
