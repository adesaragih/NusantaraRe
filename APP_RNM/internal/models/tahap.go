package models

// Tahap kerja klaim - tiket 08.
//
// Untuk apa berkas ini: tangga kerja klaim, dan siapa memegang tiap anak
// tangganya. Tahap BUKAN status: memindahkan kasus ke Medical Check tidak
// memutuskan apa pun tentang barisnya (ADR-U-0011).
//
// Dibaca sesudah: statusklaim.go.
//
// ⛔ RALAT BESAR 27-09-2026. Ronde pertama menganggap kolom `PY_POSITION`
// menyimpan pengenal shape (`"Assignment2"`). **Ia menyimpan NAMA PERAN.**
// `[terverifikasi]` `Flow/Register_Flow.xml` berkas pecahan baris 582, 605,
// 628, 668, dan 731: seluruhnya menyetel `pyWorkPage.pyPosition` ke
// `"ReasLifeAdmin"`, `"ReasLifeMedicalAdvisor"`, atau `"ReasLifeSPV"` - dan
// nol baris menyetelnya ke `"Assignment<n>"`. Gerbang di korpus pun
// membandingkannya dengan nama peran (`pyPosition == 'ReasLifeSPV'` dan
// kerabatnya, 11 kemunculan di `Claim Life`).
//
// Akibat kekeliruan itu: penerjemahan kolom atas baris nyata SELALU
// mengembalikan "tidak dikenal", dan perpindahan tahap selalu gagal - cacat
// yang tidak akan terlihat sampai ada Oracle.
//
// Istilah:
//   - tahap   : anak tangga kerja - konsep layar.
//   - pemegang: PERAN yang sedang memegang kasus - yang benar-benar tersimpan.

// Tahap adalah anak tangga kerja sebuah klaim.
type Tahap int

const (
	// TahapTidakDikenal - tahap di luar keempat yang dikenal.
	TahapTidakDikenal Tahap = iota
	// TahapInputRegister - `Assignment2`, dipegang ReasLifeAdmin.
	TahapInputRegister
	// TahapOutstanding - `Assignment1`, dipegang ReasLifeAdmin.
	TahapOutstanding
	// TahapMedicalCheck - `Assignment3`, dipegang ReasLifeMedicalAdvisor.
	TahapMedicalCheck
	// TahapClaimAnalis - `Assignment4`, dipegang ReasLifeSPV.
	TahapClaimAnalis
)

// Nama peran - satu tempat, dipakai kolom `PY_POSITION` apa adanya.
const (
	PeranAdminLife   = "ReasLifeAdmin"
	PeranSPVLife     = "ReasLifeSPV"
	PeranMedicalLife = "ReasLifeMedicalAdvisor"
)

// namaTahap adalah kata yang dibaca pengguna.
var namaTahap = map[Tahap]string{
	TahapInputRegister: "Input Register",
	TahapOutstanding:   "Outstanding Claim",
	TahapMedicalCheck:  "Medical Check",
	TahapClaimAnalis:   "Claim Analis",
}

// String menulis tahap sebagai kata, bukan sebagai pengenal shape.
func (t Tahap) String() string {
	if n, ada := namaTahap[t]; ada {
		return n
	}
	return "Tidak diketahui"
}

// Diketahui membedakan tahap yang benar-benar ada dari yang tidak.
func (t Tahap) Diketahui() bool { return t != TahapTidakDikenal }

// peranTahap memetakan tahap ke peran pemegangnya.
//
// `[terverifikasi work owner]` ADR-U-0002 tabel "Peran dan pemetaan tahap":
// Register + Outstanding → `ReasLifeAdmin`; Medical Check →
// `ReasLifeMedicalAdvisor`; Claim Analis → `ReasLifeSPV`.
//
// ⛔ Ronde pertama menandai Claim Analis `[terbuka]` dengan alasan "nol
// `pyPosition` menyebut `Assignment4` di Flow". Keliru dua kali: ADR-U-0002
// sudah memutuskannya, dan ADR itu sendiri menjelaskan kenapa korpus diam -
// *"model peran DIRANCANG, bukan dimigrasikan"*. Diamnya korpus adalah celah
// yang ADR itu TUTUP, bukan bukti untuk membukanya lagi.
var peranTahap = map[Tahap]string{
	TahapInputRegister: PeranAdminLife,
	TahapOutstanding:   PeranAdminLife,
	TahapMedicalCheck:  PeranMedicalLife,
	TahapClaimAnalis:   PeranSPVLife,
}

// PeranPemegangTahap mencari peran yang memegang sebuah tahap.
//
// ⛔ Fungsi, bukan peta telanjang. Peta di Go dibagikan lewat rujukan, dan
// satu pemanggil yang menulis ke dalamnya mengubah wewenang seluruh aplikasi
// tanpa jejak - pelajaran yang sama dengan BusinessCodeContentNote.
func PeranPemegangTahap(t Tahap) (string, bool) {
	peran, ada := peranTahap[t]
	return peran, ada
}

// TahapDariPeran menerjemahkan nilai kolom `PY_POSITION` menjadi tahap.
//
// ⚠️ TIDAK TUNGGAL bagi Admin. `ReasLifeAdmin` memegang DUA tahap - Input
// Register dan Outstanding - sehingga kolom itu tidak dapat membedakan
// keduanya. Yang dikembalikan adalah Outstanding, tahap Admin yang lebih jauh
// di tangga; pembedaan Register vs Outstanding tidak tersimpan di kolom mana
// pun. `[terbuka - work owner]` bila pembedaan itu kelak diperlukan.
func TahapDariPeran(peran string) Tahap {
	switch peran {
	case PeranAdminLife:
		return TahapOutstanding
	case PeranMedicalLife:
		return TahapMedicalCheck
	case PeranSPVLife:
		return TahapClaimAnalis
	default:
		return TahapTidakDikenal
	}
}

// serahTerimaSah adalah perpindahan PERAN yang boleh terjadi.
//
// `[terverifikasi]` ADR-U-0002: `SendtoAdmin = 1` mengembalikan dari
// `ReasLifeMedicalAdvisor` **atau** `ReasLifeSPV` ke `ReasLifeAdmin`;
// `SendtoMedical = 1` mengembalikan dari `ReasLifeSPV` ke
// `ReasLifeMedicalAdvisor`. Arah majunya dari penyambung `Register_Flow`:
// Admin → Medical Advisor → SPV.
//
// ⛔ Tanpa daftar ini setiap pasangan sah - termasuk lompatan yang tidak ada
// di tangga. Tangga yang setiap anaknya dapat dilompati bukan tangga.
//
// ⚠️ Register → Outstanding TIDAK ada di sini: keduanya dipegang peran yang
// sama, sehingga `PY_POSITION` tidak berubah dan tidak ada serah terima.
var serahTerimaSah = map[string]map[string]bool{
	PeranAdminLife:   {PeranMedicalLife: true},
	PeranMedicalLife: {PeranSPVLife: true, PeranAdminLife: true},
	PeranSPVLife:     {PeranAdminLife: true, PeranMedicalLife: true},
}

// SerahTerimaSah menyatakan kasus boleh berpindah dari satu peran ke lainnya.
func SerahTerimaSah(dari, ke string) bool { return serahTerimaSah[dari][ke] }

// JalurBalikPeran menyatakan serah terima itu PENGEMBALIAN, bukan kemajuan.
//
// Diturunkan dari pasangan perannya, tidak diterima sebagai bendera bebas:
// bendera bebas membolehkan `SENDTO_ADMIN=1` ditulis pada perpindahan menuju
// Medical Check, dan membolehkan keduanya menyala sekaligus - padahal keduanya
// menunjuk tujuan yang berbeda dan tidak mungkin bersamaan.
func JalurBalikPeran(dari, ke string) (keAdmin, keMedical bool) {
	if !SerahTerimaSah(dari, ke) {
		return false, false
	}
	switch {
	case ke == PeranAdminLife:
		return true, false
	case ke == PeranMedicalLife && dari == PeranSPVLife:
		return false, true
	default:
		return false, false
	}
}
