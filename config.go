package main

type Config struct {
	Label string `envconfig:"WATCHLABEL" default:"RECEPTIONIST"`
}
