package services

// Spreading Type dipilih dari RD, Reins Type · Pct dari `PROPORTIONALARRG` —
// keputusan pemilik proses 8 Oktober 2026, angka tangkapan layar Pega-nya.

import (
	"reflect"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
)

// Layar Pega pemakai: PROPERTY, RNM Share IDR 15.000.000.000, Spreading Type
// 10252 → QS (R/I) 60 · QS (OR) 40, Total Spreading Pct 100,00%, Value
// Spreading OR 6.000.000.000, R/I 9.000.000.000.
func TestSpreadingTypeDariRDSamaLayarPega(t *testing.T) {
	src := &sumberPropUji{
		induk: []models.SusunanSpreading{{ReinsTypeID: "10252", ReinsTypeName: "2025 QS 82M FAC", TreatyYearID: "1000679", TreatyYear: "2025"}},
		anak: map[string][]models.SusunanSpreading{"10252": {
			{ReinsTypeID: "10004", ReinsTypeName: "QS (R/I)", ParentReinsTypeID: "10252", Pct: "60"},
			{ReinsTypeID: "10028", ReinsTypeName: "QS (OR)", ParentReinsTypeID: "10252", Pct: "40"},
		}},
	}
	limits := pohonUji(t, `[{"TreatyType":"QUOTA SHARE","Detail":[{"TreatyGroup":"PROPERTY","TreatyGroupID":"10007","RNMShare":"25",
		"SpreadingTypeID":"10252","RNMShareList":[{"Currency":"IDR","Value":"15000000000"}]}]}]`)
	h := shareProp(t, MasukanShareProp{Aksi: AksiSharePropSpreading, Limits: limits, Commencement: "20250101"}, src)
	d := detailUji(h, 0, 0)

	// ⛔ GRUP KOSONG — penyimpangan yang diputuskan pemilik proses 8 Oktober
	// 2026 (*"gimana pun caranya asal itu ada isinya"*): menyaring per Treaty
	// Group mengosongkan dropdown untuk kontrak yang susunannya ADA di
	// `PROPORTIONALARRG`, hanya di grup lain. Dibuang di dropdown DAN di
	// pencarian ini sekaligus — lihat `fetchQS` (`hitung_share_np.go`).
	if !reflect.DeepEqual(src.panggil, [][]string{{"2025", "", "10001", "10252", "1000679"}}) {
		t.Errorf("RD anak %v", src.panggil)
	}
	// ⛔ RD INDUK pun dipanggil TANPA grup — penyimpangan ini hanya utuh
	// bila KEDUANYA dibuang. Menyaring induk per grup sementara anaknya
	// tidak membuat Spreading Type tidak pernah ditemukan, dan gridnya
	// kosong tanpa satu pun galat.
	for _, g := range src.grupInduk {
		if g != "" {
			t.Errorf("RD induk diminta dengan grup %q, mau kosong", g)
		}
	}
	if len(src.grupInduk) == 0 {
		t.Error("RD induk tidak dipanggil sama sekali")
	}
	var baris []string
	for _, r := range larikSimpul(d, "SpreadingList") {
		baris = append(baris, teksSimpul(r, "ReinsTypeName")+" "+teksSimpul(r, "Pct"))
	}
	if !reflect.DeepEqual(baris, []string{"QS (R/I) 60", "QS (OR) 40"}) {
		t.Errorf("Reins Type · Pct %v", baris)
	}
	if teksSimpul(d, "SpreadingTotalPct") != "100" || teksSimpul(d, "SpreadingType") != "2025 QS 82M FAC" {
		t.Errorf("total %v nama %v", d["SpreadingTotalPct"], d["SpreadingType"])
	}
	if got := nilaiUji(larikSimpul(d, "RNMSpreadedList")); !reflect.DeepEqual(got, []string{"IDR 6000000000"}) {
		t.Errorf("Value Spreading OR %v", got)
	}
	if got := nilaiUji(larikSimpul(d, "RNMSpreadedListRI")); !reflect.DeepEqual(got, []string{"IDR 9000000000"}) {
		t.Errorf("Value Spreading R/I %v", got)
	}
}

// Spreading Type dikosongkan → grid dan nilai sebaran ikut kosong.
func TestSpreadingTypeDikosongkan(t *testing.T) {
	limits := pohonUji(t, `[{"TreatyType":"QUOTA SHARE","Detail":[{"TreatyGroupID":"10007","SpreadingTypeID":"",
		"SpreadingList":[{"ReinsTypeName":"QS (OR)","Pct":"40"}],"RNMShareList":[{"Currency":"IDR","Value":"100"}]}]}]`)
	h := shareProp(t, MasukanShareProp{Aksi: AksiSharePropSpreading, Limits: limits}, &sumberPropUji{})
	d := detailUji(h, 0, 0)
	if len(larikSimpul(d, "SpreadingList")) != 0 || teksSimpul(d, "SpreadingTotalPct") != "0" {
		t.Errorf("sesudah dikosongkan %v", d)
	}
}
