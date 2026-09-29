package models

// Status simpan - tiket 09 Treaty Contract Out.
//
// `[data DBA]` keluaran prosedur penulis `StsSimpan`: 1 = sukses / 0 = gagal
// (AC 39). Prosedur tidak dipanggil sistem baru; aturan yang sama berlaku di
// batas API simpan utuh: klien membaca SELAIN "1" - kosong, NULL, "0", " 1",
// apa pun - sebagai KEGAGALAN. Tidak ada tebakan "mungkin sukses".

// StatusSimpanSuksesTeksTCO - satu-satunya nilai sukses.
const StatusSimpanSuksesTeksTCO = "1"

// StatusSimpanSuksesTCO: true HANYA untuk "1" persis; nil (NULL) gagal.
func StatusSimpanSuksesTCO(status *string) bool {
	return status != nil && *status == StatusSimpanSuksesTeksTCO
}
