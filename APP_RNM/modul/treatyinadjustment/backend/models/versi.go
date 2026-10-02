// Package models memegang bentuk data modul Treaty In Adjustment
// (`treatyinadjustment`).
//
// ⛔ Modul ini TIDAK punya §10 sendiri. Model datanya satu dengan Treaty In
// (`treaty-in/SPEC-MODEL-DATA.md`); yang dipisahkan 25-09-2026 hanya PAPAN
// TIKETNYA. Bentuk di bawah karena itu menggambarkan baris `VERSI_KONTRAK`
// yang sama, dibaca dari sudut penyesuaian.
package models

// Kontrak adalah kepala kontrak - lapisan BEKU yang seluruh versinya bagi
// tanpa menyalinnya (ADR-0040).
type Kontrak struct {
	ID                  int64  `json:"id"`
	NomorKontrakWarisan string `json:"nomorKontrakWarisan"`
	SifatProporsi       string `json:"sifatProporsi"`
	TanggalMulai        string `json:"tanggalMulai"`
	TanggalBerakhir     string `json:"tanggalBerakhir"`
}

// Versi adalah satu baris rantai versi sebuah kontrak, dibaca dari sudut
// penyesuaian: nomor urutnya, keadaannya, dan dasarnya.
//
// NomorUrutVersi bertipe penunjuk sebab tiket 05 membuatnya BOLEH KOSONG -
// kosong di sana berarti "baris warisan yang belum dinomori ulang" (tiket 10),
// bukan nol.
//
// IDVersiDasar kosong pada versi PERTAMA sebuah kontrak, dan itu keadaan yang
// benar, bukan data yang hilang (bahan to-spec `B-3`).
type Versi struct {
	ID                     int64  `json:"id"`
	IDKontrak              int64  `json:"idKontrak"`
	NomorUrutVersi         *int64 `json:"nomorUrutVersi"`
	KeadaanSiklusHidup     string `json:"keadaanSiklusHidup"`
	JenisAddendum          string `json:"jenisAddendum"`
	SifatMaterialAddendum  string `json:"sifatMaterialAddendum"`
	TanggalBerlakuAddendum string `json:"tanggalBerlakuAddendum"`
	IDVersiDasar           *int64 `json:"idVersiDasar"`
	NamaKontrak            string `json:"namaKontrak"`
}
