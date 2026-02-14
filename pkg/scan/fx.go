package scan

import "go.uber.org/fx"

var Module = fx.Module("scanner", fx.Provide(NewScanner))
