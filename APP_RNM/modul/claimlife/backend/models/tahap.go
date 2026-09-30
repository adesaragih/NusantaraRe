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
// ⛔ CADANGAN, bukan sumber - sejak butir **at** (27-09-2026). Sumber tahap
// adalah kolom `TAHAP`; fungsi ini hanya dipakai untuk baris LAMA yang
// kolomnya masih kosong.
//
// ⚠️ TIDAK TUNGGAL bagi Admin, dan di situlah batasnya: `ReasLifeAdmin`
// memegang DUA tahap - Input Register dan Outstanding Claim - sehingga
// `PY_POSITION` sendirian tidak dapat membedakan keduanya. Yang
// dikembalikan Outstanding, tahap Admin yang lebih jauh di tangga; baris
// lama yang sebenarnya berada di Input Register karena itu akan tampak
// Outstanding sampai kolom `TAHAP`-nya terisi. Itu diterima dan dicatat,
// bukan ditebak lebih jauh.
//
// `[terbuka]` yang dahulu berdiri di sini DITUTUP butir at: kolomnya ada.
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

// serahTerimaSah adalah perpindahan TAHAP yang boleh terjadi.
//
// ⛔ BUTIR at (27-09-2026): dahulu peta PERAN. Peta peran tidak dapat
// menyatakan Admin → Admin, sehingga satu perpindahan yang XML tunjukkan
// dengan terang - `Send Back to Register` - tidak punya tempat di dalamnya.
// `[terverifikasi]` `Section/InputOSClaimLife.xml` 21404 `<pyLabel>Send
// Back to Register</pyLabel>`, 21433 `<pyLocalAction>SendtoAdmin`.
//
// `[terverifikasi]` ADR-U-0002 tetap berlaku untuk jalur baliknya:
// `SendtoAdmin = 1` mengembalikan dari Medical Check **atau** Claim Analis
// ke tangan Admin; `SendtoMedical = 1` mengembalikan dari Claim Analis ke
// Medical Check. Arah majunya dari penyambung `Register_Flow`.
//
// ⛔ Tanpa daftar ini setiap pasangan sah - termasuk lompatan yang tidak ada
// di tangga. Tangga yang setiap anaknya dapat dilompati bukan tangga.
//
// ⚠️ Input Register ⇄ Outstanding Claim keduanya dipegang `ReasLifeAdmin`,
// sehingga `PY_POSITION` TIDAK berubah pada perpindahan itu - yang berubah
// hanya `TAHAP`. Itulah sebabnya peta ini tidak dapat lagi berupa peta peran.
var serahTerimaSah = map[Tahap]map[Tahap]bool{
	TahapInputRegister: {TahapOutstanding: true},
	TahapOutstanding:   {TahapInputRegister: true, TahapMedicalCheck: true},
	TahapMedicalCheck:  {TahapClaimAnalis: true, TahapOutstanding: true},
	TahapClaimAnalis:   {TahapOutstanding: true, TahapMedicalCheck: true},
}

// SerahTerimaSah menyatakan kasus boleh berpindah dari satu tahap ke lainnya.
func SerahTerimaSah(dari, ke Tahap) bool { return serahTerimaSah[dari][ke] }

// JalurBalikTahap menyatakan serah terima itu PENGEMBALIAN, bukan kemajuan.
//
// Diturunkan dari pasangan tahapnya, tidak diterima sebagai bendera bebas:
// bendera bebas membolehkan `SENDTO_ADMIN=1` ditulis pada perpindahan menuju
// Medical Check, dan membolehkan keduanya menyala sekaligus - padahal keduanya
// menunjuk tujuan yang berbeda dan tidak mungkin bersamaan.
//
// ⚠️ Kemajuan Input Register → Outstanding Claim BUKAN jalur balik, dan
// kembalinya Outstanding → Input Register JUGA bukan `SENDTO_ADMIN`:
// keduanya di tangan Admin yang sama, jadi tidak ada yang "dikembalikan"
// kepada siapa pun. Benderanya menandai pengembalian ANTARPERAN.
func JalurBalikTahap(dari, ke Tahap) (keAdmin, keMedical bool) {
	if !SerahTerimaSah(dari, ke) {
		return false, false
	}
	peranDari, _ := PeranPemegangTahap(dari)
	peranKe, _ := PeranPemegangTahap(ke)
	if peranDari == peranKe {
		return false, false
	}
	switch {
	case peranKe == PeranAdminLife:
		return true, false
	case ke == TahapMedicalCheck && dari == TahapClaimAnalis:
		return false, true
	default:
		return false, false
	}
}

// TahapBerlaku adalah tahap BERLAKU sebuah kasus - SATU sumber aturan
// cadangannya (GILIRAN-11 paket 4; sebelumnya tersalin di lima layanan).
//
// Kolom `TAHAP` menang; `PY_POSITION` hanya CADANGAN untuk baris lama yang
// kolomnya masih kosong (butir at) - dengan batas yang TahapDariPeran
// nyatakan. Hasil tak dikenal dikembalikan apa adanya: pemanggil yang memutus
// apakah itu galat.
func TahapBerlaku(kolomTahap, peranPemegang string) Tahap {
	if t := TahapDariNama(kolomTahap); t.Diketahui() {
		return t
	}
	return TahapDariPeran(peranPemegang)
}

// TahapDariNama menerjemahkan isi kolom `TAHAP` menjadi tahap - butir at.
//
// ⛔ Ia kebalikan `String()`, dan keduanya memakai peta yang SAMA
// (`namaTahap`). Dua peta terpisah berarti ada saat ketika sebuah tahap dapat
// ditulis tetapi tidak dapat dibaca kembali.
//
// ⚠️ Nama yang tidak dikenal menjadi `TahapTidakDikenal`, bukan tebakan.
// Kolom yang berisi nilai asing adalah kolom yang seseorang tulis di luar
// aplikasi ini, dan menebak artinya menyembunyikan itu.
func TahapDariNama(nama string) Tahap {
	for t, n := range namaTahap {
		if n == nama {
			return t
		}
	}
	return TahapTidakDikenal
}
