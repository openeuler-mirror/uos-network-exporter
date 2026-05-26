package main

import (
	"uos_network_exporter/internal/server"
	"uos_network_exporter/pkg/logger"
)

func Run(name string, version string) error {
	logger.InitDefaultLog()
