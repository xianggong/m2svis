package database

// ActiveInsts contains a count of active instructions in a cycle
type ActiveInsts struct {
	Cycle []int `json:"cycle"`
	Count []int `json:"count"`
}
