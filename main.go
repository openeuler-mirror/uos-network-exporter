package main

var (
	Name    = "network_exporter"
	Version = "1.0.0"
)

func mian() {
	err := Run(Name, Version)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
