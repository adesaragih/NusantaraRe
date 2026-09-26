package models

// Data acuan BusinessCode -> ContentNote - tiket 06.
//
// Untuk apa berkas ini: satu tabel data, bukan rangkaian percabangan. Jenis
// klaim ditentukan kode produk, dan daftar kodenya tumbuh dari sisi master
// produk - bukan dari sisi klaim. Tabel dapat ditambah tanpa menyentuh satu
// baris pun logika; `if` bercabang tidak.
//
// Dibaca sesudah: klaimlife.go.
//
// Sumber: `CONTEXT.md` bab "Daftar BusinessCode -> produk -> jenis klaim"
// `[terverifikasi work owner 2026-09-14]`, dua puluh satu baris.
//
// `[terverifikasi]` XML menguatkan sebelas baris pertamanya: satu precondition
// di `SaveOutStandingLife_Act` (baris 3261 berkas pecahan) menguji
// `BusinessCode` "L1" sampai "L11" berderet, dan jalur yang dijaganya adalah
// jalur `ContentNote = "DEATH"` (baris 3446, 3679, 3823; kebalikannya
// `!= "DEATH"` di 4054 dan 4199). Kesebelasnya memang `DEATH` di tabel ini -
// nol selisih, jadi nol ralat diperlukan. `L12`-`L21` tidak diuji modul klaim
// sama sekali; keduapuluh satu baris tetap ditulis supaya kode yang belum
// pernah lewat modul ini tidak jatuh ke galat "tidak dikenal" hanya karena
// modul klaim belum pernah melihatnya.
//
// ⛔ Kode tetap TEKS (ADR-U-0022): "L1" bukan 1. Mengubahnya menjadi bilangan
// menghapus awalan huruf dan memecahkan penggolong.

// businessCodeContentNote memetakan kode produk ke jenis klaim.
var businessCodeContentNote = map[string]string{
	"L1":  "DEATH",  // INDIVIDUAL TERM LIFE
	"L2":  "DEATH",  // GROUP TERM LIFE
	"L3":  "DEATH",  // INDIVIDUAL WHOLE LIFE
	"L4":  "DEATH",  // GROUP PA
	"L5":  "DEATH",  // GROUP LEVEL TERM LIFE
	"L6":  "DEATH",  // GROUP DECREASING TERM LIFE
	"L7":  "DEATH",  // INDIVIDU ENDOWMENT LIFE
	"L8":  "DEATH",  // INDIVIDU INCREASING TERM LIFE
	"L9":  "DEATH",  // GROUP ENDOWMENT LIFE
	"L10": "DEATH",  // GROUP INCREASING TERM LIFE
	"L11": "DEATH",  // INDIVIDU PA
	"L12": "HEALTH", // INDIVIDU EXPENSE HEALTH
	"L13": "HEALTH", // INDIVIDU DISABILITY HEALTH
	"L14": "HEALTH", // GROUP EXPENSE HEALTH
	"L15": "HEALTH", // GROUP DISABILITY HEALTH
	"L16": "CI",     // INDIVIDU CRITICAL ILLNESS
	"L17": "TPD",    // INDIVIDU TPD
	"L18": "HEALTH", // INDIVIDU HOSPITAL CASH PLAN
	"L19": "CI",     // GROUP CRITICAL ILLNESS
	"L20": "TPD",    // GROUP TPD
	"L21": "TI",     // GROUP TERMINAL ILLNESS
}

// ContentNoteUntuk mencari jenis klaim sebuah kode produk.
//
// Kedua nilai kembaliannya bermakna: teks kosong dengan `false` berarti kode
// itu TIDAK ADA di daftar - berbeda dari kode yang ada tetapi jenisnya kosong,
// yang tidak boleh terjadi dan akan terlihat sebagai teks kosong dengan `true`.
func ContentNoteUntuk(kodeBisnis string) (string, bool) {
	note, ada := businessCodeContentNote[kodeBisnis]
	return note, ada
}

// BusinessCodeContentNote mengeluarkan SALINAN seluruh daftar.
//
// Salinan, bukan petanya sendiri: peta di Go dibagikan lewat rujukan, dan satu
// pemanggil yang menulis ke dalamnya akan mengubah data acuan bagi seluruh
// aplikasi tanpa jejak.
func BusinessCodeContentNote() map[string]string {
	out := make(map[string]string, len(businessCodeContentNote))
	for k, v := range businessCodeContentNote {
		out[k] = v
	}
	return out
}
