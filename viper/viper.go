package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func GetYAMLFromString() {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var yamlExample = []byte(`
Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
jacket: leather
trousers: denim
age: 35
eyes : brown
beard: true
`)

	err := viper.ReadConfig(bytes.NewBuffer(yamlExample))
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	fmt.Printf("name: %v\n", viper.Get("name"))         // this would be "steve"
	fmt.Printf("hobbies: %v\n", viper.Get("hobbies"))   // this would be []string{"skateboarding", "snowboarding", "go"}
	fmt.Printf("jacket: %v\n", viper.Get("jacket"))     // this would be "leather"
	fmt.Printf("trousers: %v\n", viper.Get("trousers")) // this would be "denim"
	fmt.Printf("age: %v\n", viper.Get("age"))           // this would be "35"
	fmt.Printf("eyes: %v\n", viper.Get("eyes"))         // this would be "brown"
	fmt.Printf("beard: %v\n", viper.Get("beard"))       // this would be "true"
	fmt.Printf("Hacker: %v\n", viper.Get("Hacker"))     // this would be "true"

}

func GetYAMLFromFile() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Config file not found")
		} else {
			// Config file was found but another error was produced
			fmt.Printf("Config file found but error: %s", err)
		}
	}
	fmt.Printf("name: %v\n", viper.Get("name"))         // this would be "steve"
	fmt.Printf("hobbies: %v\n", viper.Get("hobbies"))   // this would be []string{"skateboarding", "snowboarding", "go"}
	fmt.Printf("jacket: %v\n", viper.Get("jacket"))     // this would be "leather"
	fmt.Printf("trousers: %v\n", viper.Get("trousers")) // this would be "denim"
	fmt.Printf("age: %v\n", viper.Get("age"))           // this would be "35"
	fmt.Printf("eyes: %v\n", viper.Get("eyes"))         // this would be "brown"
	fmt.Printf("beard: %v\n", viper.Get("beard"))       // this would be "true"
	fmt.Printf("Hacker: %v\n", viper.Get("Hacker"))     // this would be "true"
}

func WriteYAMLToFile() {
	viper.SetConfigName("config2")
	viper.AddConfigPath(".")
	viper.Set("name", "steve")
	viper.Set("hobbies", []string{"skateboarding", "snowboarding", "go"})
	viper.Set("jacket", "leather")
	viper.Set("trousers", "denim")
	viper.Set("age", 35)
	viper.Set("eyes", "brown")
	viper.Set("beard", true)
	viper.Set("Hacker", true)

	viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
	// viper.SafeWriteConfig()
	// viper.WriteConfigAs("/path/to/my/.config")
	// viper.SafeWriteConfigAs("/path/to/my/.config") // will error since it has already been written
	// viper.SafeWriteConfigAs("/path/to/my/.other_config")
}

func EnvExample() {
	viper.SetEnvPrefix("spf") // will be uppercased automatically
	viper.BindEnv("id")

	os.Setenv("SPF_ID", "13") // typically done outside of the app

	id := viper.Get("id") // 13
	fmt.Println(id)
}

func WatchConfigFile() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		fmt.Println("Config file op:", e.Op)
	})
	viper.WatchConfig()
}

func main() {
	GetYAMLFromString()
	GetYAMLFromFile()
	WriteYAMLToFile()
	EnvExample()
	WatchConfigFile()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
}
