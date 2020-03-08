package main

/*
 * Module Dependencies
 */

import (
	"fmt"
	"strings"

	"github.com/mozzzzy/arguments"
	"github.com/mozzzzy/arguments/argumentOption"
	"github.com/mozzzzy/config/json/config"
	"github.com/mozzzzy/config/json/configOption"
	"github.com/mozzzzy/logger"
)

/*
 * Types
 */

/*
 * Constants and Package Scope Variables
 */

/*
 * Functions
 */

func contain(slice []string, target string) bool {
	for _, elem := range slice {
		if elem == target {
			return true
		}
	}
	return false
}

func initLogger(c config.Config) error {
	// Get categories
	var categories []string
	keys := c.GetAllKeys()
	for _, key := range keys {
		elems := strings.Split(key, ".")
		category := elems[0]
		if !contain(categories, category) {
			categories = append(categories, category)
		}
	}

	for _, category := range categories {
		categoryConfig, err := c.GetObject(category)
		if err != nil {
			return err
		}
		path, err := categoryConfig.GetString("path")
		if err != nil {
			return err
		}
		levelStr, err := categoryConfig.GetString("level")
		if err != nil {
			return err
		}
		level, err := logger.GetLogLevelByStr(levelStr)
		if err != nil {
			return err
		}
		backup, err := categoryConfig.GetInt("backup")
		if err != nil {
			return err
		}
		maxSize, err := categoryConfig.GetInt64("max_size")
		if err != nil {
			return err
		}

		if err := logger.AddCategory(
			category,
			path,
			level,
			maxSize,
			backup,
		); err != nil {
			return err
		}
	}
	return nil
}

func configArgOptions() (arguments.Args, error) {
	var args arguments.Args
	err := args.AddOptions([]argumentOption.Option{
		{
			LongKey:     "config",
			ShortKey:    "c",
			Description: "Specify config file.",
			ValueType:   "string",
		},
		{
			LongKey:     "help",
			ShortKey:    "h",
			Description: "Show help message and exit.",
		},
	})
	return args, err
}

func configConfigOptions() (config.Config, error) {
	var conf config.Config
	err := conf.AddOptions([]configOption.Option{
		// logger section
		{
			Key:         "logger",
			Description: "Section about logger.",
			ValueType:   "object",
			Required:    true,
		},
		// logger.diagnostic section
		{
			Key:         "logger.diagnostic",
			Description: "Section about diagnostic logger.",
			ValueType:   "object",
			Required:    true,
		},
		{
			Key:         "logger.diagnostic.path",
			Description: "File path of diagnostic log file.",
			ValueType:   "string",
			Required:    true,
		},
		{
			Key:          "logger.diagnostic.level",
			Description:  "Log level of diagnostic log.",
			ValueType:    "string",
			DefaultValue: "INFO",
		},
		{
			Key:          "logger.diagnostic.backup",
			Description:  "Number of lotated old diagnostic log files.",
			ValueType:    "int",
			DefaultValue: 5,
		},
		{
			Key:          "logger.diagnostic.max_size",
			Description:  "Max size of diagnostic log file.",
			ValueType:    "int64",
			DefaultValue: int64(1073741824), // 1gb
		},
	})
	return conf, err
}

func parseArgs() (arguments.Args, error) {
	// Configure argument options
	args, configArgOptsErr := configArgOptions()
	if configArgOptsErr != nil {
		fmt.Println(configArgOptsErr)
		return args, configArgOptsErr
	}
	// Parse argument options
	optParseErr := args.Parse()
	return args, optParseErr
}

func parseConfig(path string) (config.Config, error) {
	// Configure configuration options
	config, configConfOptsErr := configConfigOptions()
	if configConfOptsErr != nil {
		return config, configConfOptsErr
	}
	// Parse config file
	configParseErr := config.Parse(path)
	return config, configParseErr
}

func main() {
	// Configure configuration options
	config, err := configConfigOptions()
	if err != nil {
		fmt.Printf("Falied to setup configuration rules.\n")
		fmt.Println(err)
		return
	}

	// Parse argument options
	args, err := parseArgs()
	if err != nil {
		fmt.Printf("Falied to parse argument options.\n")
		fmt.Println(err)
		fmt.Println(args)
		fmt.Println(config)
		return
	}

	// If --help, -h option is specified, show usage and return.
	if args.IsSet("help") {
		fmt.Println(args)
		fmt.Println(config)
		return
	}

	// Get config file path
	configFilePath, err := args.GetString("config")
	if err != nil {
		fmt.Printf("Falied to get config file path.\n")
		fmt.Println(err)
		return
	}

	// Parse config file
	if err := config.Parse(configFilePath); err != nil {
		fmt.Printf("Failed to parse config file %v.\n", configFilePath)
		fmt.Println(err)
		fmt.Println(config)
		return
	}

	// Get logger section config
	loggerConfig, err := config.GetObject("logger")
	if err != nil {
		fmt.Println(err)
	}

	// Initialize logger
	if err := initLogger(loggerConfig); err != nil {
		fmt.Println(err)
		return
	}

	log, err := logger.New("diagnostic")
	log.Info("Start.")

	/*
		if err := exec(); err != nil {
			log.Error(err.Error())
			return
		}
	*/

	log.Info("Finish.")
}
