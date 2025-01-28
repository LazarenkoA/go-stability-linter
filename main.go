package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"stabilityIndicators/app_module"
)

func init() {
	//stdPackages, _ = packages.Load(nil, "std")
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "go-stability-linter [directory]",
		Short: "Go Dependency Stability Linter",
		Long:  "Go Dependency Stability Linter is a tool for analyzing the stability of code dependencies in Go projects. Inspired by Robert C. Martin's \"Clean Architecture,\" it calculates the instability of components based on their Fan-in and Fan-out metrics, helping developers identify and manage tightly coupled code.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Error: The directory path is not specified.\nUsage: go-stability-linter [directory]")
				os.Exit(1)
			}

			l, err := app_module.NewLint(args[0])
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}

			if err = l.Check(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
