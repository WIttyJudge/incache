package incache

import "time"

type expirationsQueue map[string]time.Time
