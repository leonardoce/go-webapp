/*
Copyright © 2024 Leonardo Cecchi
*/
package cmd

import (
	"github.com/leonardoce/go-webapp/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := server.New(cmd.Context())
		if err != nil {
			return err
		}

		if err := s.Start(cmd.Context()); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP(
		"listen",
		"l",
		"127.0.0.1:8000",
		"Where should we listen for HTTP requests?")
	viper.BindPFlag("listen", serveCmd.Flags().Lookup("listen"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
