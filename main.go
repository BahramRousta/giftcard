package main

import "giftcard/cmd"

func main() {
	cmd.Execute()
	//var config config2.Config
	//
	//// Set the name of the config file (without extension)
	//viper.SetConfigName("config")
	//// Set the path to look for the config file
	//viper.AddConfigPath(".")
	//// Read the config file
	//if err := viper.ReadInConfig(); err != nil {
	//	fmt.Printf("Error reading config file: %s\n", err)
	//	return
	//}
	//
	//// Unmarshal the config into the Config struct
	//if err := viper.Unmarshal(&config); err != nil {
	//	fmt.Printf("Error unmarshaling config: %s\n", err)
	//	return
	//}
	//
	//// Use the config...
	//fmt.Printf("Service name: %s\n", config.Service.Name)
	//fmt.Printf("Redis host: %s\n", config.Redis.Host)
	//fmt.Printf("Postgres username: %s\n", config.DataBase.Username)
	//fmt.Printf("Tracer hostPort: %s\n", config.Jaeger.HostPort)
	//fmt.Printf("Tracer LogSpans: %v\n", config.Jaeger.LogSpans)
	//fmt.Printf("Logstash endpoint: %s\n", config.Logstash.Endpoint)
}
