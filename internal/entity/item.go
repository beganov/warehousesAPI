package entity

type Item struct {
	Name        string `json:"name"`
	Size        string `json:"size"`
	UniqueId    string `json:"unique_id"`
	Quantity    int    `json:"quantity"`
	WarehouseId int    `json:"warehouse_id"`
}
