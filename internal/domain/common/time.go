package common

import "time"

var Fujian = time.FixedZone("Asia/Shanghai", 8*60*60)

func Now() time.Time                        { return time.Now().In(Fujian) }
func Parse(value string) (time.Time, error) { return time.Parse(time.RFC3339, value) }
