package main

import (
	"fmt"
	"github.com/LazarenkoA/go-stability-linter/app_module"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	//stdPackages, _ = packages.Load(nil, "std")
}

func main() {
	var verbose bool

	var rootCmd = &cobra.Command{
		Use:   "go-stability-linter [directory]",
		Short: "Go Dependency Stability Linter",
		Long:  "Go Dependency Stability Linter is a tool for analyzing the stability of code dependencies in Go projects. Inspired by Robert C. Martin's \"Clean Architecture,\" it calculates the instability of components based on their Fan-in and Fan-out metrics, helping developers identify and manage tightly coupled code.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			l, err := app_module.NewLint(args[0])
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}

			if !verbose {
				if err = l.Check(); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					os.Exit(1)
				}
			} else {
				print(l.GetPackageInfoTree())
			}
		},
	}

	rootCmd.Flags().BoolVarP(&verbose, "tree", "t", false, "visual tree output")

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
