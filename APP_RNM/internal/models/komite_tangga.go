package models

// Mesin tangga Komite - tiket 02 Komite Claim Life. MURNI.
//
// `[terverifikasi]` `Komite Claim Life/`:
//
//	Flow/KomiteLife_Flow.xml         Start2 -[Always]-> Assignment1
//	                                 Assignment1 -[ViewTransferDtl]-> Decision1
//	                                 Decision1 -[IsKomiteLoop]-> Assignment1
//	                                 Decision1 -[NoLoop]-> End1 (pyWorkStatus kosong)
//	When/IsKomiteLoop.xml            .AcceptStatus = "1" DAN .KomiteCount <= .KomiteLoop
//	Activity/KomitePostAdjustment    langkah 2  Local.Komite = pyWorkPage.KomiteCount
//	                                 langkah 3  KomiteList(Local.Komite).KomiteAproval = .AcceptStatus,
//	                                            .KomiteComment = .Comment, .DateApprove = @CurrentDateTime()
//	                                 langkah 4  "Approve Last Komite" bila AcceptStatus = 1
//	                                            DAN KomiteCount == KomiteLoop (b5695)
//	                                 langkah 5  "Reject" bila AcceptStatus==2 DAN
//	                                            KomiteCount == KomiteLoop (b8119)
//	                                 langkah 13 KomiteCount = KomiteCount + 1 (b9028-9029), tanpa syarat
//
// ⛔ URUTANNYA MENENTUKAN. Langkah 4/5 membaca `KomiteCount` SEBELUM langkah
// 13 menaikkannya; `IsKomiteLoop` membacanya SESUDAH. Tingkat akhir karena itu
// `count == loop` (sebelum naik), dan tangga berlanjut bila `count+1 <= loop`.

import (
	"errors"
	"fmt"
	"strings"
)

// Nilai `.AcceptStatus` - enum TERTUTUP.
//
// ⚠️ `[keputusan work owner]` (AC 35 spec Komite). Dropdown `.AcceptStatus`
// di `Section/ShowTransfer.xml` b32607 ber-`pyListSource associated`: daftar
// pilihannya milik aturan properti `AcceptStatus` yang TIDAK ikut diekspor.
// Korpus hanya membuktikan `1` (`IsKomiteLoop`, b5695) dan `2` (b8119).
const (
	KeputusanKomiteSetuju = "1"
	KeputusanKomiteTolak  = "2"
)

// ErrKeputusanKomiteTidakDikenal - nilai di luar {1, 2}.
//
// ⛔ Ditolak TERANG. Di Pega nilai asing hanya membuat `IsKomiteLoop` palsu -
// tangga berhenti diam-diam, tanpa akseptasi dan tanpa penolakan.
var ErrKeputusanKomiteTidakDikenal = errors.New(
	"models: keputusan komite hanya 1 (Setuju) atau 2 (Tolak)")

// KataKeputusanKomite menerjemahkan kode keputusan menjadi kata.
func KataKeputusanKomite(kode string) string {
	switch strings.TrimSpace(kode) {
	case KeputusanKomiteSetuju:
		return "Setuju"
	case KeputusanKomiteTolak:
		return "Tolak"
	}
	return ""
}

// TanggaBerlanjut adalah `When/IsKomiteLoop.xml`, VERBATIM.
func TanggaBerlanjut(acceptStatus string, komiteCount, komiteLoop int) bool {
	return strings.TrimSpace(acceptStatus) == KeputusanKomiteSetuju && komiteCount <= komiteLoop
}

// KasusDiTangga menjawab apakah kasus masih berdiri di `Assignment1`.
//
// Kasus baru (`AcceptStatus` belum pernah diisi) tiba lewat `Start2
// -[Always]->`; sesudah keputusan pertama ia di sana hanya bila
// `IsKomiteLoop` benar. `End1` tidak punya `pyWorkStatus` di ekspor, jadi
// "sudah berhenti" DITURUNKAN dari syarat yang sama - bukan dari status kerja
// yang karangan.
func KasusDiTangga(acceptStatus string, komiteCount, komiteLoop int) bool {
	if strings.TrimSpace(acceptStatus) == "" {
		return komiteCount >= 1 && komiteCount <= komiteLoop
	}
	return TanggaBerlanjut(acceptStatus, komiteCount, komiteLoop)
}

// AkibatKeputusanKomite adalah apa yang terjadi sesudah satu anggota memutuskan.
type AkibatKeputusanKomite struct {
	// TingkatDiputus = `Local.Komite` = `KomiteCount` sebelum naik.
	TingkatDiputus int
	// CountBaru = `KomiteCount + 1` (langkah 13).
	CountBaru int
	// Berlanjut = `IsKomiteLoop` atas CountBaru.
	Berlanjut bool
	// AkseptasiAkhir = langkah 4 menyala (Setuju di tingkat akhir).
	AkseptasiAkhir bool
	// TolakAkhir = langkah 5 menyala (Tolak di tingkat akhir).
	TolakAkhir bool
	Keputusan  string
}

// TerapkanKeputusanKomite menghitung akibat satu keputusan.
func TerapkanKeputusanKomite(keputusan string, komiteCount, komiteLoop int) (
	AkibatKeputusanKomite, error) {

	k := strings.TrimSpace(keputusan)
	if k != KeputusanKomiteSetuju && k != KeputusanKomiteTolak {
		return AkibatKeputusanKomite{}, fmt.Errorf("%w: %q", ErrKeputusanKomiteTidakDikenal, keputusan)
	}
	if komiteCount < 1 || komiteLoop < 1 || komiteCount > komiteLoop {
		return AkibatKeputusanKomite{}, fmt.Errorf(
			"models: tangga komite tidak sah (count %d, loop %d)", komiteCount, komiteLoop)
	}
	akhir := komiteCount == komiteLoop
	baru := komiteCount + 1
	return AkibatKeputusanKomite{
		TingkatDiputus: komiteCount,
		CountBaru:      baru,
		Berlanjut:      TanggaBerlanjut(k, baru, komiteLoop),
		AkseptasiAkhir: k == KeputusanKomiteSetuju && akhir,
		TolakAkhir:     k == KeputusanKomiteTolak && akhir,
		Keputusan:      k,
	}, nil
}
