package model

type Driver struct {
	DriverId string `bson:"driver_id"`
	Name     string `bson:"name"`
	Auto     string `bson:"auto"`
}
