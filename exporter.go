package main

import (
	"uos_network_exporter/internal/server"
	"uos_network_exporter/pkg/logger"
)

func Run(name string, version string) error {
	logger.InitDefaultLog()
	s := server.NewServer(name, version)

	s.PrintVersion()
	err := s.SetUp()
	if err != nil {
		return err
	}
	go func() {
		err := s.Run()
		if err != nil {
			s.Error = err
		}
		s.Exit()
	}()
