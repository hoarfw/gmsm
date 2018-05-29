package ca

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "gmca",
	Short: "gmca generate cert for sm2withsm3",
	Long: `gmca generate cert for sm2withsm3
				  Sm2WithSM3 Req 
				  Sm2WithSM3 Cert
				  Sm2WithSM3 CA`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println(viper.GetString("cakey"));
	},
}

func RootCmd() {
	rootCmd.AddCommand(reqCmd())
	rootCmd.AddCommand(certCmd())
	rootCmd.AddCommand(selfcaCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ca.yaml)")
	//rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	rootCmd.PersistentFlags().StringP("cakey", "k", "ca.key", "ca key file")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	//rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	viper.BindPFlag("cakey", rootCmd.PersistentFlags().Lookup("cakey"))
	//viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	//viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	//viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	//viper.SetDefault("license", "apache")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ca")

		// For environment variables.
		//viper.SetEnvPrefix(cmdRoot)
		//viper.AutomaticEnv()
		//replacer := strings.NewReplacer(".", "_")
		//viper.SetEnvKeyReplacer(replacer)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
	//viper.WatchConfig()
	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("Config file changed:", e.Name)
	//})
}
