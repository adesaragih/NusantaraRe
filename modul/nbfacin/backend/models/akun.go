package models

// Akun - satu baris POOLDATA.T_M_ACCOUNT untuk popup ChooseAccount (tiket 27). Seluruh
// kolom teks apa adanya; NULL dibaca sebagai "".
type Akun struct {
	ID              string
	GroupBusinessID string
	GroupBusiness   string
	InsuredID       string
	InsuredName     string
}
