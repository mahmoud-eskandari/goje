package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mahmoud-eskandari/goje/internal/config"
	"github.com/mahmoud-eskandari/goje/internal/generator"
	"github.com/mahmoud-eskandari/goje/internal/repo"
	"github.com/mahmoud-eskandari/goje/internal/version"
)

func main() {
	Execute()
}

func Execute() {
	conf := flag.String("c", "config.goje.yaml", "Set config path")

	if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == "init" {
		initConfig(*conf)
		fmt.Println("Config file created successfully!")
		fmt.Printf("Edit [%s] to set database connection and other options", *conf)
		return
	}

	if len(os.Args) > 1 && strings.Contains(os.Args[1], "ersion") {
		fmt.Printf("goje Ver(%s) developed by Mahmoud Eskandari.", version.Ver)
		return
	}

	err := config.ReadConfig(*conf)
	if err != nil {
		log.Fatalf("Error when reading config: %+v \n **** Run: `goje init` to create a default config file ****", err)
	}

	err = repo.Connect()
	if err != nil {
		log.Fatalf("Error on connecting to database: %+v", err)
	}

	err = generator.Generate()
	if err != nil {
		log.Fatalf("Error on generate files: %+v", err)
	}

}
