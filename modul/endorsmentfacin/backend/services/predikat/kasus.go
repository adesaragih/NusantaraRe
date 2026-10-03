package predikat

// Sumber data per-rule (tiket E01) - berkas milik EDM, bukan salinan nbfacin.

import (
	"fmt"
	"sort"
	"strings"
)

// JalurJenisBisnisLama - kolom yang dibaca keenam predikat lini bisnis EDM
// (IsAneka, IsFire, IsGolfInsurance, IsMBU, IsMarineCargo, IsPA):
// `[terverifikasi]` `grep -l "OutData.pxResults(1).CARI2" "Endorsment Fac
// In"/When/*.xml` → 6 berkas, sama dengan hitungan atas registry
// (`TestSumberJenisBisnisLamaEnamPredikat`). Isinya jenis bisnis POLIS LAMA,
// diisi `RDBList/GetBusinessType_Sql` lewat `SetErrorBatalEndorsement_Act`
// langkah 7-8 (services.JenisBisnisUntukPredikat).
const JalurJenisBisnisLama = "OutData.pxResults(1).CARI2"

// KasusEDM - Kasus dengan dua sumber: properti kasus (jalur PERSIS seperti di
// rule) dan hasil query jenis bisnis polis lama.
type KasusEDM struct {
	Properti map[string]string
	// JenisBisnisLama - keluaran services.JenisBisnisUntukPredikat. nil = query
	// belum dijalankan.
	JenisBisnisLama *string
	// HalamanBernilai - jawaban PropertyHasValue per jalur (lihat PunyaNilai).
	HalamanBernilai map[string]bool
}

// Nilai - ⛔ urutan mengikat (E01): predikat lini yang dinilai SEBELUM query
// berjalan akan membaca kosong dan bernilai salah DIAM-DIAM; di sini itu
// kesalahan program → panic.
func (k KasusEDM) Nilai(jalur string) (string, bool) {
	if jalur == JalurJenisBisnisLama {
		if k.JenisBisnisLama == nil {
			panic("predikat: " + JalurJenisBisnisLama + " dibaca sebelum query jenis bisnis polis lama dijalankan")
		}
		return *k.JenisBisnisLama, true
	}
	v, ada := k.Properti[jalur]
	return v, ada
}

// Sumber - jalur properti yang dibaca baris kondisi predikat `nama` sendiri
// (tanpa menelusuri rujukan `evaluateWhen`), urut. Predikat panic → nil.
func Sumber(nama string) []string {
	pr, ada := registry[strings.ToUpper(nama)]
	if !ada {
		panic("predikat: predikat " + nama + " tidak ada di registry")
	}
	set := map[string]bool{}
	for _, kd := range pr.kondisi {
		if kd.jenis == banding {
			set[kd.kiri] = true
		}
	}
	hasil := make([]string, 0, len(set))
	for j := range set {
		hasil = append(hasil, j)
	}
	sort.Strings(hasil)
	return hasil
}

// PemeriksaNilai - Kasus yang dapat menjawab `@(Pega-RULES:Utilities).
// PropertyHasValue(jalur)`.
//
// ⚠️ `[belum terverifikasi]` Arti "punya nilai" untuk HALAMAN (bukan properti
// skalar) tidak dijelaskan korpus - mis. halaman ada tetapi kosong. Karena itu
// jawabannya diserahkan ke pemanggil (lapisan yang memegang clipboard/kasus),
// tidak ditebak di sini.
type PemeriksaNilai interface {
	PunyaNilai(jalur string) (bool, error)
}

// PunyaNilai - jawaban HalamanBernilai per jalur, diisi pemanggil.
// Jalur yang tidak terdaftar → galat (bukan "tidak punya nilai").
func (k KasusEDM) PunyaNilai(jalur string) (bool, error) {
	v, ada := k.HalamanBernilai[jalur]
	if !ada {
		return false, fmt.Errorf("predikat: PropertyHasValue(%s) tidak dijawab pemanggil", jalur)
	}
	return v, nil
}

// sikapKhusus - predikat EDM yang entri registry-nya panic karena bentuk
// ekspresinya tidak diport generator, diport tangan di sini dari ekspresi
// tersimpan (`pyConditionValue1`).
var sikapKhusus = map[string]func(Kasus) (bool, error){
	// `Endorsment Fac In/When/IsClaim.xml` (WORK-!ISCLAIM), logika A:
	// `@(Pega-RULES:Utilities).PropertyHasValue(pyWorkPage.ClaimData)`.
	// ⛔ K-046 - varian NB (`NB FacIn/When/IsClaim.xml`) menguji
	// `pyWorkPage.pyWorkIDPrefix = "CLM-"`; EDM menguji keberadaan halaman.
	"ISCLAIM": func(k Kasus) (bool, error) {
		pn, ok := k.(PemeriksaNilai)
		if !ok {
			panic("predikat: IsClaim butuh Kasus yang menjawab PropertyHasValue (PemeriksaNilai)")
		}
		return pn.PunyaNilai(JalurHalamanKlaim)
	},
}

// JalurHalamanKlaim - halaman yang diuji IsClaim EDM.
const JalurHalamanKlaim = "pyWorkPage.ClaimData"
