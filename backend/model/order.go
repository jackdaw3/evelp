package model

import "time"

type Order struct {
	OrderId      int
	ItemId       int
	Issued       time.Time
	Duration     int
	SystemId     int
	Price        int64
	VolumeRemain int64
	VolumeTotal  int64
}
