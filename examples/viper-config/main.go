package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

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

func main() {
	// set env var
	os.Setenv("BEARD", "false")

	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	// bind target env
	viper.BindEnv("BEARD")

	viper.ReadConfig(bytes.NewBuffer(yamlExample))

	fmt.Printf("name : %v keys: %+v \n", viper.Get("name"), viper.AllKeys()) // this would be "steve"
	fmt.Printf("keys: %+v \n", viper.AllKeys())

	// will be false because OS Env Vars is set
	fmt.Println(viper.GetBool("beard"))
}
