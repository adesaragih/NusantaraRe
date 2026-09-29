package outbox

// Penjaga penyambungan jalur menyerah - TANPA Oracle.
//
// Refactor bentuk B (30-09-2026): penjaga ini dipindah apa adanya dari
// `internal/services/jejak_statik_test.go`, mengikuti subjeknya - `antrean.go`
// kini tinggal di paket ini. Pemotong fungsi dan pembuang komentarnya disalin
// (bukan dibagi) supaya paket bersama tidak bergantung pada berkas uji modul.

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

type fungsiGo struct {
	nama  string
	tubuh string
}

// polaAwalFungsi mencocokkan baris pembuka sebuah fungsi atau method.
var polaAwalFungsi = regexp.MustCompile(`(?m)^func (?:\([^)]*\) )?([A-Za-z_][A-Za-z0-9_]*)`)

// pecahPerFungsi memotong sebuah berkas Go menjadi tubuh fungsi-fungsinya.
func pecahPerFungsi(isi string) []fungsiGo {
	letak := polaAwalFungsi.FindAllStringSubmatchIndex(isi, -1)
	out := make([]fungsiGo, 0, len(letak))
	for i, l := range letak {
		akhir := len(isi)
		if i+1 < len(letak) {
			akhir = letak[i+1][0]
		}
		out = append(out, fungsiGo{nama: isi[l[2]:l[3]], tubuh: isi[l[0]:akhir]})
	}
	return out
}

// buangKomentar menghapus baris komentar sebelum pencocokan.
func buangKomentar(isi string) string {
	var b strings.Builder
	for _, baris := range strings.Split(isi, "\n") {
		if strings.HasPrefix(strings.TrimSpace(baris), "//") {
			continue
		}
		b.WriteString(baris)
		b.WriteString("\n")
	}
	return b.String()
}

// AC 20 tiket 12: setiap jalur yang MENYERAH meninggalkan jejak audit.
//
// ⛔ Penjaga penyambungan, bukan penjaga perilaku. Uji perilakunya memanggil
// `rekamMenyerah` langsung - dan itu terbukti TIDAK cukup: mencabut
// pemanggilan `rekamMenyerah` dari `Antre` membiarkan uji itu HIJAU, sebab ia
// tidak pernah melewati `Antre`. Yang hilang bukan logikanya melainkan
// sambungannya, dan sambungan dijaga di sini.
//
// Aturannya: fungsi mana pun di `antrean.go` yang menulis status
// `gagal-permanen` WAJIB merekam jejaknya di fungsi yang sama.
func TestSetiapJalurMenyerahMerekamJejak(t *testing.T) {
	isi, err := os.ReadFile("antrean.go")
	if err != nil {
		t.Fatal(err)
	}
	diperiksa := 0
	for _, fn := range pecahPerFungsi(string(isi)) {
		tubuh := buangKomentar(fn.tubuh)
		// Yang menulis kegagalan permanen - bukan yang sekadar menyebutnya.
		if !strings.Contains(tubuh, "StatusEfekGagalPermanen") {
			continue
		}
		diperiksa++
		if !strings.Contains(tubuh, "rekamMenyerah(ctx, tx,") {
			t.Errorf("%s menulis status gagal-permanen tetapi tidak merekam "+
				"jejaknya di fungsi yang sama; AC 20 tiket 12 menuntut "+
				"kegagalan masuk jalur audit, bukan hanya antrean", fn.nama)
		}
	}
	// DUA jalur menyerah: gagal permanen sejak awal, dan jatah habis.
	const mauJalur = 2
	if diperiksa != mauJalur {
		t.Errorf("%d jalur menyerah ditemukan, mau %d - bila jalurnya bertambah, "+
			"penjaga ini harus ikut tahu", diperiksa, mauJalur)
	}
}
