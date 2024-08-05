package config

import "flag"

type Flags struct {
	Addr     *string
	BaseAddr *string
}

func ParseFlags() Flags {
	var flags Flags

	flags.Addr = flag.String("a", "localhost:8080", "address for HTTP-server (addr:port); default localhost:8080")
	flags.BaseAddr = flag.String("b", "", "sets base address for all resulting short urls; if not set uses -a flag address")

	flag.Parse()

	// checking is -b flag was set
	// if it's not, then flags.Addr used as base
	if *flags.BaseAddr == "" {
		baseAddr := "http://" + *flags.Addr + "/"
		flags.BaseAddr = &baseAddr
	}

	return flags
}

type Config struct {
	Addr     string
	BaseAddr string
}

func New() Config {
	flags := ParseFlags()

	return Config{
		Addr:     *flags.Addr,
		BaseAddr: *flags.BaseAddr,
	}
}
