package flagd_definitions

import _ "embed"

//go:embed flags.json
var FlagSchema string

//go:embed targeting.json
var TargetingSchema string
