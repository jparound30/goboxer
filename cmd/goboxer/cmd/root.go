package cmd

import (
	"fmt"
	"github.com/jparound30/goboxer"
	"io/ioutil"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var clientId string
var clientSecret string
var accessToken string
var refreshToken string

var apiConn *goboxer.ApiConn

func init() {
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goboxer",
	Short: "Utility CLI TOOL made with goboxer",
	Long: `Utility CLI TOOL made with goboxer
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var mainState Main

func init() {
	mainState = Main{}

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goboxer.yaml)")

	rootCmd.PersistentFlags().StringVar(&clientId, "clientId", "", "ClientID")
	rootCmd.PersistentFlags().StringVar(&clientSecret, "clientSecret", "", "ClientSecret")
	rootCmd.MarkPersistentFlagRequired("clientId")
	rootCmd.MarkPersistentFlagRequired("clientSecret")
	rootCmd.PersistentFlags().StringVar(&accessToken, "accessToken", "", "AccessToken")
	rootCmd.PersistentFlags().StringVar(&refreshToken, "refreshToken", "", "RefreshToken")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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

		// Search config in home directory with name ".goboxer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".goboxer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

var (
	StateFilename = "apiconnstate.json"
)

func createGobxerApiConn() error {
	apiConn = goboxer.NewApiConnWithRefreshToken(clientId, clientSecret, accessToken, refreshToken)

	_, err := os.Stat(StateFilename)
	if err == nil {
		bytes, err := ioutil.ReadFile(StateFilename)
		err = apiConn.RestoreApiConn(bytes)
		if err != nil {
			return err
		}
	}
	if accessToken != "" {
		apiConn.AccessToken = accessToken
	}
	if refreshToken != "" {
		apiConn.RefreshToken = refreshToken
	}

	apiConn.SetApiConnRefreshNotifier(&mainState)
	goboxer.Log = &mainState
	return nil
}

type Main struct {
}

func (*Main) RequestDumpf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) ResponseDumpf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) Debugf(format string, args ...interface{}) {
	fmt.Printf("[goboxer] "+format, args...)
}

func (*Main) Infof(format string, args ...interface{}) {
	fmt.Printf("[goboxer] "+format, args...)
}

func (*Main) Warnf(format string, args ...interface{}) {
	fmt.Printf("[goboxer] "+format, args...)
}

func (*Main) Errorf(format string, args ...interface{}) {
	fmt.Printf("[goboxer] "+format, args...)
}

func (*Main) Fatalf(format string, args ...interface{}) {
	fmt.Printf("[goboxer] "+format, args...)
}
func (*Main) EnabledLoggingResponseBody() bool {
	return true
}
func (*Main) EnabledLoggingRequestBody() bool {
	return false
}

func (*Main) Success(apiConn *goboxer.ApiConn) {
	fmt.Printf("access_token: %s\n", apiConn.AccessToken)
	fmt.Printf("refresh_token: %s\n", apiConn.RefreshToken)
	bytes, err := apiConn.SaveState()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	err = ioutil.WriteFile(StateFilename, bytes, 0666)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

func (*Main) Fail(apiConn *goboxer.ApiConn, err error) {
	fmt.Printf("%v\n", err)
}

var UTF8_BOM = []byte{0xEF, 0xBB, 0xBF}
