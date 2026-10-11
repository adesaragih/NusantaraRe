package models

// Arsip muatan keluar — tiket 42 dan 74.
//
// ⛔ DUA bentuk, dan pemisahannya adalah inti tiket 42. `ArsipMuatanKeluar`
// dipakai MENULIS dan membawa muatannya; `BuktiArsip` dipakai MEMBACA dan
// **tidak punya ruas muatan sama sekali**.
//
// Pemisahan itu bukan kerapian: `INV-61` menyatakan arsip JSON **tidak punya
// jalur baca** dan bukan sumber kanonik, dan `ADR-0034` menyebut sebabnya -
// arsip yang punya jalur baca akan dipakai sebagai sumber oleh orang yang
// sedang buru-buru, dan sejak saat itu ada DUA kebenaran. Larangan yang hanya
// ditulis di komentar dilanggar dalam enam bulan; larangan yang diwujudkan
// sebagai struct tanpa ruasnya tidak dapat dilanggar tanpa menyuntingnya.

// ArsipMuatanKeluar adalah satu pengiriman ke hilir, disimpan apa adanya.
//
// ⚠️ Dipakai HANYA pada jalur tulis. Nol fungsi di repository yang
// mengembalikannya.
type ArsipMuatanKeluar struct {
	ID        int64  `json:"id,omitempty"`
	IDKontrak int64  `json:"idKontrak"`
	Tujuan    string `json:"tujuan"`
	// RFC 3339; jam pengiriman, bukan jam penyimpanan.
	DikirimPada string `json:"dikirimPada"`
	// Pengiriman yang GAGAL di hilir tetap diarsipkan - arsip yang hanya
	// memuat yang berhasil tidak menyelesaikan satu pun perselisihan.
	Berhasil bool `json:"berhasil"`
	// Muatan APA ADANYA. Nol jalur baca mengembalikannya.
	Muatan string `json:"muatan"`
}

// BuktiArsip menjawab "apakah arsipnya ada", dan tidak lebih.
//
// ⛔ NOL ruas muatan, dan itu seluruh maksudnya. Tiket 42 menyebut
// kemampuannya "dapat menunjukkan bahwa arsip itu ada" - bukan "dapat
// membacanya". Menambahkan ruas muatan ke sini membalikkan `ADR-0034` tanpa
// membukanya.
type BuktiArsip struct {
	IDKontrak int64 `json:"idKontrak"`
	// Cacah pengiriman terarsip. Nol berarti kontrak ini belum pernah
	// dikirim ke hilir - jawaban, bukan galat.
	Cacah int64 `json:"cacah"`
	// RFC 3339; kosong ketika Cacah nol.
	TerakhirDikirim string `json:"terakhirDikirim,omitempty"`
	// Tujuan yang pernah menerima, berurut. Namanya saja - bukan isinya.
	Tujuan []string `json:"tujuan"`
}
