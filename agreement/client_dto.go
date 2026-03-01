package agreement

type AddClientRequest struct {
	ScannerClientID uint   `json:"scanner_client_id" binding:"required"`
	Role            string `json:"role" binding:"omitempty,max=100"`
}
