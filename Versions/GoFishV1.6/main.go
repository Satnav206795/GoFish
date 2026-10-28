package main

func main() {
	settings := load_settings("settings.json")
	uci_loop(settings)
}
