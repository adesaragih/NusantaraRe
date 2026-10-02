package models

// Bentuk jalur baca warisan — tiket 40 dan 41.
//
// Keduanya TURUNAN: tidak ada kolom yang menyimpannya, dan itu bukan
// kebetulan melainkan syaratnya (`INV-58`, `INV-59`).

// VersiSebelumnya adalah jawaban atas "apa nilai versi dasarnya" — tiket 40.
//
// ⛔ `Ada` adalah ruasnya sendiri, dan ia yang membedakan DUA jawaban yang
// sama sekali berbeda: *"versi ini yang pertama, tidak ada sebelumnya"* dan
// *"ada versi sebelumnya, nilainya kosong"*. Tiket 40 menyebutnya di jalur
// gagalnya: hasilnya PERNYATAAN KETIADAAN, bukan nol - dan struct tanpa `Ada`
// memaksa pemanggilnya menyimpulkannya dari `Versi == nil`, yang berarti
// setiap pemanggil mengulang kesimpulan yang sama.
type VersiSebelumnya struct {
	Ada bool `json:"ada"`
	// Versi dasar APA ADANYA saat dibaca - bukan salinan yang disimpan saat
	// versi baru dibuat. Nil ketika Ada bernilai false.
	Versi *VersiKontrak `json:"versi,omitempty"`
}

// IdentitasLama adalah identitas kontrak dalam bentuk sistem lama — tiket 41.
//
// ⛔ DITURUNKAN saat dibaca, TIDAK disimpan. `INV-58` melarang menyimpannya,
// dan `ADR-0051` memerintahkan menyediakannya: tabel datar `TREATYINDETAIL`
// ditulis pada setiap penyimpanan dan tidak satu pun aturan di dalam ekspor
// membacanya kembali - yang berarti pembacanya ada di luar ekspor, dan kita
// tidak tahu siapa.
type IdentitasLama struct {
	// NomorLamaAda false berarti kontrak ini LAHIR DI SISTEM BARU dan memang
	// tidak punya nomor lama. Itu keadaan yang dinyatakan, bukan ditambal:
	// membangkitkan nomor berpola lama untuknya berarti mengarang pengenal
	// yang tidak pernah ada di sistem mana pun.
	NomorLamaAda bool   `json:"nomorLamaAda"`
	NomorLama    string `json:"nomorLama,omitempty"`

	IDKontrak       int64  `json:"idKontrak"`
	IDCedant        int64  `json:"idCedant"`
	IDAsalBisnis    int64  `json:"idAsalBisnis"`
	SifatProporsi   string `json:"sifatProporsi"`
	TanggalMulai    string `json:"tanggalMulai"`
	TanggalBerakhir string `json:"tanggalBerakhir"`
}
