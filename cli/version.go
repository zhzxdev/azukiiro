package cli

import (
	"fmt"
	"runtime/debug"
)

func GetVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		settings := make(map[string]string)
		for _, setting := range info.Settings {
			settings[setting.Key] = setting.Value
		}
		hash, ok := settings["vcs.revision"]
		if !ok {
			return "azukiiro/unknown"
		}
		hash = hash[:7]
		return fmt.Sprintf("azukiiro/%s %s@%s", hash, info.Main.Version, settings["vcs.time"])
	}
	return "azukiiro/unknown"
}
