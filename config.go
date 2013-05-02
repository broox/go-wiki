package main

import (
    "io/ioutil"
    "encoding/json"
)

// A struct to represent our configuration
type Config struct {
    Port string
    Database struct {
        User string
        Password string `json: "pass"`
        Host string
        Name string
    }
}

func (conf *Config) FromJson(path string) (err error) {
    b, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }
    err = json.Unmarshal(b, &conf)
    if err != nil {
        return err
    }
    return
}