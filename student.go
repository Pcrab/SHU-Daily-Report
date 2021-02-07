package main

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"strings"
)

func MakeStudents() {

	files, err := ioutil.ReadDir(configPath)
	if err != nil {
		log.Fatal(err)
	}
	configs := []string{}
	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, configExtension) {
			configs = append(configs, file.Name()[:len(file.Name())-5])
		}
	}

	for _, config := range configs {
		students = append(students, MakeStudent(config, configPath, configExtension))
	}
	return
}

func MakeStudent(config, configPath, configExtension string) (student Student) {
	v := viper.New()
	v.SetConfigName(config)
	v.SetConfigType(configExtension)
	v.AddConfigPath(configPath)
	err := v.ReadInConfig() // Find and read the config file
	if err != nil {         // Handle errors reading the config file
		log.Printf("Fatal error config file: %s \n", err)
	}

	student.Username = v.GetString("username")
	student.Password = v.GetString("password")
	student.Cookie = v.GetString("cookie")
	student.Config = config

	return
}
