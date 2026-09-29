package repository

// Nama tabel dan sequence yang dipakai repository modul SEBELUM paket 2 tco4
// mengarahkannya ke tabel warisan. ⛔ Tabel-tabel ini TIDAK ADA lagi (migrasi
// 300-307 dibuang, keputusan work owner 29-09-2026); berkas ini hilang di
// paket 2.
// Nama enam tabel BARU (keputusan tco1: awalan T_).
const (
	TabelTahunTCO     = "T_TREATYYEAR"
	TabelKontrakTCO   = "T_TREATYCONTRACT"
	TabelReinsurerTCO = "T_TREATYREINSURER"
	TabelSecurityTCO  = "T_MTREATYSECURITY"
	TabelBusinessTCO  = "T_TREATYBUSINESS"
	TabelKlausulTCO   = "T_PROPORTIONALARRG"
	TabelJejakTCO     = "T_TREATYCO_JEJAK"
	// TabelLampiranTCO - lampiran tahun treaty, fitur BARU tiket 12 (tanpa warisan).
	TabelLampiranTCO = "T_TREATYYEAR_LAMPIRAN"
)

// Sequence tiap tabel baru, beserta LEBAR digit di belakang awalan '1'
// (spec §7: enam digit, kecuali klausul tujuh).
const (
	SeqTahunTCO     = "SEQ_T_TREATYYEAR"
	SeqKontrakTCO   = "SEQ_T_TREATYCONTRACT"
	SeqReinsurerTCO = "SEQ_T_TREATYREINSURER"
	SeqSecurityTCO  = "SEQ_T_MTREATYSECURITY"
	SeqBusinessTCO  = "SEQ_T_TREATYBUSINESS"
	SeqKlausulTCO   = "SEQ_T_PROPORTIONALARRG"
	SeqJejakTCO     = "SEQ_T_TREATYCO_JEJAK"
	// SeqLampiranTCO - tiket 12.
	SeqLampiranTCO = "SEQ_T_TREATYYEAR_LAMPIRAN"

	LebarIdentitasTCO        = 6
	LebarIdentitasKlausulTCO = 7
	// LebarIdentitasLampiranTCO - tiket 12, tanpa padanan warisan (keputusan kami).
	LebarIdentitasLampiranTCO = 9
)
