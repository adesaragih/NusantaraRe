package models

// Isi pemilih "Choose Ceding" dan "Choose Source of Business".
//
// ⛔ SUMBERNYA `TREATY_IN` SENDIRI, dan itu keputusan pemilik proses
// 4 Oktober 2026 — bukan jalan pintas. Nol tabel master yang terisi ada di
// POOLDATA untuk keduanya, dan ketiganya sudah diperiksa:
// `T_CEDINGCOLIST` 2 baris, `M_ORGANIZATION` 1 baris, `CATEGORY_ATTACH_REAS`
// isinya FACIN/FACOUT. Yang benar-benar ada adalah nilai yang sudah menempel
// di 1.853 dari 1.854 kontrak: 131 pasangan pengenal+nama cedant (94 nama
// berbeda) dan 96 pasangan asal bisnis (91 nama berbeda).
//
// ⚠️ SATU NAMA DAPAT BER-ID LEBIH DARI SATU, dan itu TIDAK disembunyikan.
// Terukur 4 Oktober 2026: **36 nama cedant** dipakai lebih dari satu
// pengenal. `REASURANSI MAIPARK INDONESIA` punya 3,
// `ASURANSI ADIRA DINAMIKA` punya 2 (`G0000006` dan ID kerja Pega
// `ASM-SFAGIS-WORK-ORG ORG-34`). Pemilih yang menyatukannya menjadi satu
// baris harus MEMILIH pengenal mana yang ditulis — dan pilihan itu diam,
// salahnya baru terbaca saat rekonsiliasi. Karena itu tiap pasangan
// pengenal+nama berdiri sebagai barisnya sendiri, dan yang memilih melihat
// bahwa ada dua.
// OpsiLimits - isi ketiga jenis dropdown tab Limits proporsional.
type OpsiLimits struct {
	JenisTreaty    []PilihanWarisan `json:"jenisTreaty"`
	KelompokTreaty []PilihanWarisan `json:"kelompokTreaty"`
	MataUang       []PilihanWarisan `json:"mataUang"`
}

type PilihanWarisan struct {
	// Pengenal yang DITULIS ke kontrak bila baris ini dipilih.
	ID string `json:"id"`
	// Nama sebagaimana tersimpan — tidak dirapikan, tidak disatukan.
	Nama string `json:"nama"`
	// Benar bila nama ini dipakai LEBIH DARI SATU pengenal. Diisi services,
	// bukan repository, dan dipakai layar untuk menandainya.
	Kembar bool `json:"kembar"`
}
