package models

// Penawaran polis Life - tiket 01 PremiumList Life.
//
// Untuk apa berkas ini: tangga kerja polis dan ketiga keputusan atasnya.
// Seluruhnya MURNI - nol basis data, nol jam.
//
// Sumbernya SATU berkas, dibaca sebagai POHON 28-09-2026:
// `PremiumList Life/InputPolicyHolder.xml`. Petanya dirakit dengan menyusuri
// `pyMOId` tiap blok, lalu `pyTo` dan `pyFrom` di dalamnya.
//
// ⚠️ `pyTo` MENDAHULUI `pyFrom` di DOM. Membaca berpasangan dari atas tanpa
// menyadari itu menghasilkan graf yang seluruh panahnya terbalik - dan graf
// terbalik tetap terlihat masuk akal.
//
// Peta lengkapnya, dua belas konektor, seluruh `pyExpression` VERBATIM:
//
//	Start1        --Transition1 -------------------> Assignment2   (Input Offer Life b1340)
//	Assignment2   --Transition3  [InputDataOfferLife b1735] -> Decision1
//	Decision1     --Transition4  [Confirm b1574] ---> Decision3
//	Decision1     --Transition5  [Decline b2235] ---> End1         (Resolved-Rejected b848)
//	Decision3     --Transition10 [Premium b1658] ---> ASSIGNMENT63 (Input Premium Detail b1069)
//	Decision3     --Transition11 [Offer   b1807] ---> END52        (Resolved-Completed b947)
//	ASSIGNMENT63  --TRANSITION54 [ShowLifePremiumDetail b1505] -> Decision2
//	Decision2     --Transition7  [Confirm b2090] ---> Utility1     (InsertJsonPolisLife b765)
//	Decision2     --Transition6  [Decline b2162] ---> End1         (Resolved-Rejected)
//	Decision2     --Transition9  [Reject  b2306] ---> Assignment2  (kembali ke Input Offer Life)
//	Utility1      --Transition8 -------------------> END52         (Resolved-Completed)
//	Assignment1   --Transition2  [ShowLifePremiumSummary b1881] -> END52
//
// Dibaca sesudah: tahap.go (tangga kerja Claim Life; bentuknya sama,
// isinya berbeda).

import (
	"errors"
	"fmt"
	"strings"
)

// Status kerja polis - VERBATIM dari shape akhir flow.
//
// ⛔ DUA, dan tidak lebih. Nol shape lain di `InputPolicyHolder.xml`
// ber-`pyWorkStatus` selain kedua ini dan ketiga nama assignment di bawah.
// Status untuk polis yang sedang berjalan TIDAK ada di ekspor, dan tidak
// dikarang: kosong berarti belum selesai (ADR-U-0027).
const (
	// StatusPolisDitolak - shape `End1` b832, `pyWorkStatus` b848.
	//
	// ⚠️ Ia dicapai oleh `Decline`, BUKAN oleh `Reject`. Nama statusnya
	// menyebut "Rejected" sedangkan keputusan yang menuju ke sana bernama
	// "Decline" - kosakata sistem lama, ditiru apa adanya. Menyeragamkannya
	// berarti status yang tersimpan tidak lagi cocok dengan data warisan.
	StatusPolisDitolak = "Resolved-Rejected"
	// StatusPolisSelesai - shape `END52` b934, `pyWorkStatus` b947.
	StatusPolisSelesai = "Resolved-Completed"
)

// Nama tahap polis - VERBATIM `pyWorkStatus` tiap assignment.
//
// ⛔ RALAT 28-09-2026. Komentar ronde pertama berbunyi *"nilai inilah yang
// masuk `T_WORK_POLIS.POSITION`"*. KELIRU. Nilai ini `pyWorkStatus`, dan yang
// masuk kolom `STATUS`. Kolom `POSITION` menyimpan hal yang BERBEDA - lihat
// `PosisiOffer`/`PosisiPremium` di bawah.
//
// Yang membantahnya `Activity/ProtectAccept.xml`, yang membandingkan
// `pyWorkPage.Position` dengan `"Offer"` (b1207) dan `"Premium"` (b2288,
// b3792, b3958, b4144, b4376) - bukan dengan nama assignment mana pun.
//
// ⚠️ Akibat kekeliruannya nyata: baris warisan menyimpan `Offer` atau
// `Premium` di `POSITION`, sedangkan kode yang mencarinya dengan
// `"Input Offer Life"` akan menemukan NOL baris - kotak masuk kosong untuk
// pekerjaan yang benar-benar ada.
const (
	// TahapPolisPenawaran - `Assignment2` b1316, status b1340.
	TahapPolisPenawaran = "Input Offer Life"
	// TahapPolisDetail - `ASSIGNMENT63` b1033, status b1069.
	TahapPolisDetail = "Input Premium Detail"
	// TahapPolisSummary - `Assignment1` b1215, status b1232.
	//
	// ⛔ `[terbuka - work owner]` TAHAP INI TIDAK PUNYA SATU PUN KONEKTOR
	// MASUK. Kedua belas konektor flow disusuri satu per satu; nol di
	// antaranya ber-`pyTo` `Assignment1`. Yang keluar ada - `Transition2`
	// b1866 `[ShowLifePremiumSummary]` menuju `END52`.
	//
	// Bentuk yang sama pernah ditemukan di Claim Life (`Assignment4`, yang
	// nol `pyPosition`-nya). Artinya tahap ini dicapai lewat TICKET
	// (`Ticket1` b958/b1456) atau lewat jalur yang tidak ikut diekspor.
	// Perpindahan KE tahap ini karena itu TIDAK disediakan di sini - ia akan
	// menjadi jalur yang kami karang sendiri.
	TahapPolisSummary = "Input Premium Summary"
)

// Posisi layar polis - VERBATIM nilai `pyWorkPage.Position`.
//
// ⛔ DUA nilai, dan hanya dua di seluruh korpus:
//
//	"Offer"    `InputPolicyHolder.xml` b712 (property-set `Assignment2`),
//	           b2352 (assign konektor), `InputOfferLife_preAct.xml` b271
//	"Premium"  `countCategoryAttachment_act.xml` b272
//
// ⚠️ Ia BUKAN tahap, dan bukan pula hasil penggolong - walau kata-katanya
// sama dengan `LanjutOffer`/`LanjutPremium`. Yang ini keadaan tersimpan pada
// halaman kerja; yang itu `pyExpression` sebuah konektor. Keduanya sengaja
// dibedakan namanya di sini supaya tidak ada yang menyamakan keduanya tanpa
// memeriksa - walau nilainya kebetulan sama.
//
// Gunanya: `Activity/ProtectAccept.xml` memakainya untuk memilih pemeriksaan
// mana yang berlaku. `COB can't null` hanya saat `Offer` (b1207); `SOB can't
// null` dan seluruh pemeriksaan saldo hanya saat `Premium` (b2288, b3958).
const (
	PosisiOffer   = "Offer"
	PosisiPremium = "Premium"
)

// Keputusan penawaran - VERBATIM `pyExpression` konektor keputusan.
//
// `[keputusan work owner]` Ketiganya dibuat MANUAL oleh inputor. Kedua
// decision table (`IsLifeAccepted`, `IsFlagOnGoingPolicy`) mengekspor NOL
// baris keputusan, jadi yang ditiru AKIBAT keputusan - bukan formula yang
// memilihnya.
//
// ⛔ RALAT 29-09-2026 (GILIRAN-13): "NOL baris keputusan" KELIRU untuk
// kedua decision table. `IsLifeAccepted` memetakan `ProposalAcceptStatus`
// 1 -> Confirm, 2 -> Reject (b286/b287 -> b321/b322); `IsFlagOnGoingPolicy`
// memetakan `FlagOnGoingPolicy` "0" -> Offer, "1" -> Premium (b293/b294 ->
// b328/b329). Untuk Decision1/2 keputusan manual tetap selaras - inputor yang
// mengisi `ProposalAcceptStatus`. Untuk Decision3 TIDAK: benderanya lahir
// bersama kasus (`CreateInputLife` b618, kini kolom `FLAG_ONGOING_POLICY`,
// 057), jadi Pega merutekannya otomatis. Perilaku di sini TIDAK diubah -
// OQ-PL-16.
const (
	KeputusanConfirm = "Confirm"
	KeputusanReject  = "Reject"
	KeputusanDecline = "Decline"
)

// Hasil penggolong lanjutan - VERBATIM `pyExpression` `Decision3`.
//
// `[keputusan work owner]` `IsFlagOnGoingPolicy`: `1` = Offer, `2` = Premium.
// ⛔ RALAT 29-09-2026: nilai benderanya "0" = Offer, "1" = Premium - VERBATIM
// decision table (b293/b294) dan tombol portal (b3310, b3958); butir bn
// menetapkannya untuk kolom `FLAG_ONGOING_POLICY`.
const (
	LanjutOffer   = "Offer"
	LanjutPremium = "Premium"
)

var (
	// ErrKeputusanTidakDikenal - bukan Confirm/Reject/Decline.
	ErrKeputusanTidakDikenal = errors.New("models: keputusan penawaran tidak dikenal")
	// ErrKeputusanTidakAdaDiTahapIni - keputusannya sah, konektornya tidak ada.
	ErrKeputusanTidakAdaDiTahapIni = errors.New(
		"models: tahap ini tidak punya konektor untuk keputusan itu")
	// ErrTahapPolisTidakDikenal - bukan salah satu dari ketiga tahap.
	ErrTahapPolisTidakDikenal = errors.New("models: tahap polis tidak dikenal")
)

// AkibatKeputusan adalah apa yang terjadi sesudah sebuah keputusan.
//
// ⚠️ TIGA kemungkinan, dan ketiganya berbeda - itulah sebab tipe ini ada
// alih-alih sepasang string. `Confirm` di tahap penawaran TIDAK menutup dan
// TIDAK memindahkan: ia menyerahkan kasus ke penggolong berikutnya, dan
// hasil penggolong itulah yang memindahkan.
type AkibatKeputusan struct {
	// TahapTujuan terisi bila kasus BERPINDAH tahap.
	TahapTujuan string
	// StatusWork terisi bila kasus DITUTUP.
	StatusWork string
	// MenungguPenggolong true bila yang berikutnya adalah `Decision3`,
	// yaitu `IsFlagOnGoingPolicy` - Offer atau Premium.
	MenungguPenggolong bool
	// SimpanPolis true bila jalurnya melewati `InsertJsonPolisLife_Act` -
	// `Utility1` b765 sesudah `Transition7`, atau tombol `Submit` layar
	// summary (b26414). Tiket 05b: yang tersisa dari activity itu sesudah
	// JSON dibuang adalah penomoran, rekap, dan salinan peserta warisan.
	//
	// ⛔ BUKAN turunan `StatusWork`. `Offer` juga berakhir di `END52`
	// Resolved-Completed, tetapi lewat `Transition11` LANGSUNG - tanpa
	// `Utility1`. Menyimpulkan "tutup selesai = simpan" akan menyimpan
	// premium list untuk penawaran yang tidak pernah punya rincian.
	SimpanPolis bool
}

// Ditutup menyatakan keputusan ini mengakhiri kasus.
func (a AkibatKeputusan) Ditutup() bool { return a.StatusWork != "" }

// TransisiPenawaran menjawab apa yang terjadi pada sebuah keputusan.
//
// ⛔ TABELNYA DARI KONEKTOR, bukan dari kalimat tiket. Dan di satu titik
// keduanya BERBEDA - lihat ralat tiket 01 bertanggal 28-09-2026:
//
//	tiket  : "Reject pada tahap MANA PUN mengembalikan case ke Input Offer"
//	konektor: `Reject` hanya ada pada `Decision2` (b2306), yaitu SESUDAH
//	          Input Premium Detail. `Decision1` - penggolong sesudah tahap
//	          penawaran - hanya punya `Confirm` b1574 dan `Decline` b2235.
//
// Menolak di tahap penawaran karena itu TIDAK punya jalur di sistem lama:
// orang yang berada di layar Input Offer dan tidak ingin melanjutkan memakai
// `Decline`. Menyediakan `Reject` di sana berarti membuat jalur yang tidak
// pernah ada, dan kasus yang menempuhnya akan mendarat di tempat yang tidak
// dikenal sistem hilir.
func TransisiPenawaran(tahap, keputusan string) (AkibatKeputusan, error) {
	k := strings.TrimSpace(keputusan)
	switch k {
	case KeputusanConfirm, KeputusanReject, KeputusanDecline:
	default:
		return AkibatKeputusan{}, fmt.Errorf("%w: %q", ErrKeputusanTidakDikenal, keputusan)
	}

	switch strings.TrimSpace(tahap) {
	case TahapPolisPenawaran:
		switch k {
		case KeputusanConfirm:
			// Transition4 b1560 [Confirm b1574] -> Decision3.
			return AkibatKeputusan{MenungguPenggolong: true}, nil
		case KeputusanDecline:
			// Transition5 b2220 [Decline b2235] -> End1 b832.
			return AkibatKeputusan{StatusWork: StatusPolisDitolak}, nil
		}
	case TahapPolisDetail:
		switch k {
		case KeputusanConfirm:
			// Transition7 b2076 [Confirm b2090] -> Utility1 b765
			//   --Transition8 b2016--> END52 b934.
			//
			// ⚠️ Lewat `Utility1`, yang menjalankan `InsertJsonPolisLife`
			// (b261/b782) - tiket 05b: `SimpanPolis` menandainya, dan
			// layanan menjalankannya di transaksi yang sama dengan
			// penutupan, SEBELUM kasus ditutup.
			return AkibatKeputusan{StatusWork: StatusPolisSelesai, SimpanPolis: true}, nil
		case KeputusanDecline:
			// Transition6 b2147 [Decline b2162] -> End1 b832.
			return AkibatKeputusan{StatusWork: StatusPolisDitolak}, nil
		case KeputusanReject:
			// Transition9 b2291 [Reject b2306] -> Assignment2 b1316.
			return AkibatKeputusan{TahapTujuan: TahapPolisPenawaran}, nil
		}
	case TahapPolisSummary:
		// Nol konektor keputusan dari tahap ini - satu-satunya yang keluar
		// `Transition2` b1866, dan ia bukan keputusan melainkan penyelesaian
		// layar summary (tiket 05a).
		return AkibatKeputusan{}, fmt.Errorf(
			"%w: tahap %q tidak punya konektor keputusan sama sekali",
			ErrKeputusanTidakAdaDiTahapIni, tahap)
	default:
		return AkibatKeputusan{}, fmt.Errorf("%w: %q", ErrTahapPolisTidakDikenal, tahap)
	}
	return AkibatKeputusan{}, fmt.Errorf("%w: %q pada %q",
		ErrKeputusanTidakAdaDiTahapIni, k, tahap)
}

// LanjutanPenggolong menjawab akibat hasil `IsFlagOnGoingPolicy`.
//
// `Decision3` b879, dua konektor keluar:
//
//	Transition10 b1643 [Premium b1658] -> ASSIGNMENT63 (Input Premium Detail)
//	Transition11 b1793 [Offer   b1807] -> END52        (Resolved-Completed)
//
// ⛔ `Offer` MENUTUP kasus, tidak "berhenti di tahap penawaran" dalam arti
// menunggu. Penawarannya tetap tersimpan - yang selesai adalah pekerjaannya.
func LanjutanPenggolong(hasil string) (AkibatKeputusan, error) {
	switch strings.TrimSpace(hasil) {
	case LanjutPremium:
		return AkibatKeputusan{TahapTujuan: TahapPolisDetail}, nil
	case LanjutOffer:
		return AkibatKeputusan{StatusWork: StatusPolisSelesai}, nil
	}
	return AkibatKeputusan{}, fmt.Errorf(
		"%w: %q bukan Offer maupun Premium", ErrKeputusanTidakDikenal, hasil)
}

// TahapPolisDikenal menyatakan tahap itu salah satu dari ketiganya.
func TahapPolisDikenal(tahap string) bool {
	switch strings.TrimSpace(tahap) {
	case TahapPolisPenawaran, TahapPolisDetail, TahapPolisSummary:
		return true
	}
	return false
}

// KasusPolisTertutup menyatakan status kerja itu keadaan akhir.
//
// ⛔ Perbandingan TEPAT, bukan awalan `Resolved-`. Pega punya banyak status
// berawalan itu; yang dikenal flow ini hanya dua, dan menerima awalan berarti
// menerima status yang alurnya tidak pernah hasilkan.
func KasusPolisTertutup(statusWork string) bool {
	s := strings.TrimSpace(statusWork)
	return s == StatusPolisDitolak || s == StatusPolisSelesai
}

// PenyelesaianSummary adalah akibat tombol `Submit` di layar summary.
//
// `[terverifikasi]` `Section/ShowLifePremiumSummary.xml`: `Submit` (b27471)
// menjalankan `InsertJsonPolisLife_Act` (b26414) lalu `finishAssignment`
// (b26442). Assignment yang diselesaikan adalah `Assignment1`, dan satu-
// satunya konektor keluarnya `Transition2` [ShowLifePremiumSummary b1881] ->
// `END52` - `pyWorkStatus` Resolved-Completed b947.
//
// ⚠️ Tahap itu tetap nol konektor MASUK (lihat `TahapPolisSummary`): fungsi
// ini hanya menjawab apa yang terjadi SESUDAH kasus berada di sana.
func PenyelesaianSummary() AkibatKeputusan {
	return AkibatKeputusan{StatusWork: StatusPolisSelesai, SimpanPolis: true}
}
