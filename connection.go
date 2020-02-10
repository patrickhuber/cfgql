package cfgql

// Connection defines a connection to a cloud foundry foundation
type Connection struct {
	ID       *string `json:"id"`
	Name     *string `json:"name"`
	API      *string `json:"api"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}
