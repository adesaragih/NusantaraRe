package models

// Untuk apa berkas ini: MEDAN TURUNAN (dihitung ulang setiap muat dan sesudah setiap aksi, tidak disimpan) dan pengurai
// kunci panel masterDetail.

import (
	"regexp"
	"strconv"
	"strings"
)

// HitungTurunan menulis medan turunan halaman: indeks item dari objek induknya (`TurunkanIndeksItem`), nomor baris
// (`.pxListSubscript`) setiap daftar di pohon objek, dan label baris item Marine Cargo.
func HitungTurunan(h *Halaman) {
	TurunkanIndeksItem(h)
	for jalur, rows := range h.Daftar {
		if !strings.HasPrefix(jalur, DaftarObjek) && jalur != DaftarCalon {
			continue
		}
		for n, b := range rows {
			b[PropNomorBaris] = strconv.Itoa(n + 1)
		}
	}
	if IsMarineCargo(h) {
		for o := range h.AmbilDaftar(DaftarObjek) {
			for _, it := range h.AmbilDaftar(DaftarItem(o + 1)) {
				it[PropKlikMarine] = "Click Here to View Coverage and Input Estimation"
			}
		}
	}
	LabelAdjMarine(h)
	LengkapiMataUangAdj(h)
}

// LengkapiMataUangAdj - `.CurencyAdjustment` setiap adjustment (CountTotalEstimasi_Act 12-14, 17.1: mata uang
// estimasi item ber-PrintFaceClaim 1, unik). Pilihan "Choose Currency" bukan kolom; disusun ulang dari estimasi item
// yang tersimpan (`[inferensi]`, PARITAS).
func LengkapiMataUangAdj(h *Halaman) {
	for o := range h.AmbilDaftar(DaftarObjek) {
		for i := range h.AmbilDaftar(DaftarItem(o + 1)) {
			var mu []Baris
			lihat := map[string]bool{}
			for _, e := range h.AmbilDaftar(DaftarDiItem(o+1, i+1, AnakEstimasi)) {
				if e["PrintFaceClaim"] == "1" && !lihat[e["CurrencyID"]] {
					lihat[e["CurrencyID"]] = true
					mu = append(mu, Baris{"CurrencyID": e["CurrencyID"], "Currency": e["Currency"], "KursValue": e["KursValue"]})
				}
			}
			for a := range h.AmbilDaftar(DaftarAdj(o+1, i+1)) {
				h.SetelDaftar(DaftarDiAdj(o+1, i+1, a+1, DaftarMataUangAdj), SalinDaftar(mu))
			}
		}
	}
}

// BuangTurunan membuang medan turunan baris sebelum disimpan (katalog tidak memuatnya - tidak perlu, tetapi tiruan
// menyalin baris apa adanya).
func BuangTurunan(h *Halaman) {
	for jalur, rows := range h.Daftar {
		if !strings.HasPrefix(jalur, DaftarObjek) {
			continue
		}
		for _, b := range rows {
			delete(b, PropNomorBaris)
			delete(b, PropKlikMarine)
			delete(b, PropKlikAdjMarine)
		}
	}
}

var indeksPanel = regexp.MustCompile(`\((\d+)\)`)

// UraiPanel memecah kunci panel `prefiks:jalur(n)...` (`KunciPanel`) menjadi prefiks dan indeks berurutan dari luar ke
// dalam: "estitem:ClaimData.ObjectList(1).ObjectItemList(2)" -> ("estitem", [1 2]).
func UraiPanel(kunci string) (string, []int, bool) {
	prefiks, jalur, ok := strings.Cut(kunci, ":")
	if !ok || prefiks == "" || jalur == "" {
		return "", nil, false
	}
	var idx []int
	for _, m := range indeksPanel.FindAllStringSubmatch(jalur, -1) {
		n, err := strconv.Atoi(m[1])
		if err != nil || n < 1 {
			return "", nil, false
		}
		idx = append(idx, n)
	}
	return prefiks, idx, len(idx) > 0
}
