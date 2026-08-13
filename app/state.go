package main

type AppState struct {
	Running  bool
	Settings SettingsState
}

type SettingsState struct {
	DarkMode   bool
	ShowHidden bool
}
