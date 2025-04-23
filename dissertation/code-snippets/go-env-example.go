func main() {
	ctx := context.Background()
	cfg, err := env.ParseAs[config.DatabaseConfig]()
	if err != nil {
		panic(err)
	}

	// Start-Up Code...
}
