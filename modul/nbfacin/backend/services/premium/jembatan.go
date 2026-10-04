// Berkas ini menyambungkan registry predikat (`services/rules`, Seam 1) dengan
// resolver satuan (`resolver.go`, Seam 2) - tiket NB-10. Lini bisnis sebuah kasus
// tidak diterima mentah dari pemanggil, tetapi diturunkan dari gerbang predikat
// yang sungguh dipakai sistem lama untuk memilih rumus premi.

package premium

import (
	"fmt"
	"sort"
	"strings"

	"nusantarare/modul/nbfacin/backend/services/rules"
)

// gerbangLini - gerbang pemilih rumus premi, URUT seperti langkahnya.
//
// [terverifikasi] Asal, `NB FacIn\Activity\`:
//   - CountGrossPremi_Act.xml (ASM-FW-GISFW-WORK / COUNTGROSSPREMI_ACT) langkah
//     5.1 `IsFire`, 5.2 `IsAneka`, 5.3 `isGolfInsurance`, 5.4 memanggil
//     CountGPWMarinePAMbu_Act tanpa syarat;
//   - CountGPWMarinePAMbu_Act.xml langkah 1.1 `IsMarineCargo`, 1.2 `IsMBU`,
//     1.3 `IsPA`.
//
// Gerbang-gerbang itu BERDIRI SENDIRI (setiap langkah WhenFalse 3, lewati):
// satu kasus dapat membuka lebih dari satu. Karena itu hasilnya daftar, bukan
// satu lini.
//
// [terverifikasi] BONDING tidak punya gerbang sendiri: IsAneka
// (NB FacIn\When\IsAneka.xml, ASM-FW-GISFW-DATA!ISANEKA) baris AA merujuk
// IsBondingAndCustomBonds, sehingga kasus Bonding dihitung lewat jalur ANEKA
// (keputusan agent A11, dikonfirmasi work owner 01-10-2026).
// LiniBonding tetap di peta resolver (K-018) tetapi tidak dihasilkan di sini.
// Layering juga tidak: tidak dipakai di sistem baru (butir 30).
//
// ⚠️ Langkah 5 dan 1 sendiri bergerbang
// `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness!=3` - 3 = Endorsement
// (keputusan work owner 01-10-2026, butir 29): premi endorsement tidak lewat sini.
// Gerbang itu milik alur pemanggil dan TIDAK diport di sini.
var gerbangLini = []struct {
	predikat string
	lini     LiniBisnis
}{
	{"IsFire", LiniFire},
	{"IsAneka", LiniAneka},
	{"isGolfInsurance", LiniGolf},
	{"IsMarineCargo", LiniMarineCargo},
	{"IsMBU", LiniMBU},
	{"IsPA", LiniPA},
}

// jalurJenisBisnis - jalur `BusinessType` yang dibaca predikat-predikat lini.
// [terverifikasi] dari registry: dua absolut, dua relatif. Keempatnya
// DIPERTAHANKAN apa adanya (NB-10): apakah selalu sinkron belum terjawab, jadi
// tidak dinormalisasi menjadi satu medan.
var jalurJenisBisnis = []string{
	"pyWorkPage.Quotation.BusinessType",
	"pyWorkPage.OfferFacIn.QuotationData.BusinessType",
	".Quotation.BusinessType",
	".OfferFacIn.QuotationData.BusinessType",
}

// LiniTerpilih - satu lini yang gerbangnya terbuka, beserta satuannya dari
// resolver.
type LiniTerpilih struct {
	Lini   LiniBisnis
	Satuan Satuan
}

// HasilLini - lini yang gerbangnya terbuka, urut langkah, dan peringatan yang
// wajib dibaca pemanggil (jalur tak sinkron, spasi dipangkas).
type HasilLini struct {
	Lini       []LiniTerpilih
	Peringatan []string
}

// LiniDariPredikat menurunkan lini bisnis kasus `k` dari gerbang predikat.
//
// Galat - keraguan data dari registry (rules.ErrTafsirBerbeda). Panic - tidak
// satu gerbang pun terbuka: lini tidak dikenali, gagal keras lewat resolver.
func LiniDariPredikat(k rules.Kasus) (HasilLini, error) {
	kp := &kasusTerpangkas{asli: k}
	var hasil HasilLini
	for _, g := range gerbangLini {
		buka, err := rules.Eval(g.predikat, kp)
		if err != nil {
			return HasilLini{}, err
		}
		if buka {
			// Tidak ada jalur yang melewati resolver.
			hasil.Lini = append(hasil.Lini, LiniTerpilih{Lini: g.lini, Satuan: SatuanRate(g.lini)})
		}
	}
	// peringatanJalur ikut membaca lewat kp (dan dapat menambah catatan
	// pangkas), jadi dihitung LEBIH DULU: urutan evaluasi operan append tidak
	// dijamin spesifikasi Go.
	jalur := peringatanJalur(kp)
	hasil.Peringatan = append(append([]string{}, kp.dipangkas...), jalur...)
	if len(hasil.Lini) == 0 {
		// Gagal keras, bukan lini bawaan. Panic ditulis langsung: memanggil
		// SatuanRate dengan nilai BusinessType justru LOLOS bila nilai itu
		// kebetulan nama lini resolver (mis. "BONDING") - temuan review.
		nilai, _ := kp.Nilai(jalurJenisBisnis[0])
		panic(fmt.Sprintf("premium: lini bisnis tidak dikenali - tidak satu gerbang pemilih rumus terbuka (%s = %q)",
			jalurJenisBisnis[0], nilai))
	}
	return hasil, nil
}

// peringatanJalur - satu peringatan bila jalur-jalur BusinessType yang terisi
// bernilai berbeda.
func peringatanJalur(k rules.Kasus) []string {
	nilai := map[string][]string{}
	for _, j := range jalurJenisBisnis {
		if v, ada := k.Nilai(j); ada && v != "" {
			nilai[v] = append(nilai[v], j)
		}
	}
	if len(nilai) <= 1 {
		return nil
	}
	var bagian []string
	for v, jalur := range nilai {
		bagian = append(bagian, fmt.Sprintf("%q di %s", v, strings.Join(jalur, ", ")))
	}
	sort.Strings(bagian)
	return []string{"jalur BusinessType berbeda: " + strings.Join(bagian, "; ")}
}

// kasusTerpangkas memangkas spasi di ujung nilai jalur BusinessType - batas
// input, NB-10. Kandidat perbaikan: sumber datanya menyimpan spasi.
type kasusTerpangkas struct {
	asli      rules.Kasus
	dipangkas []string
}

func (k *kasusTerpangkas) Nilai(jalur string) (string, bool) {
	v, ada := k.asli.Nilai(jalur)
	for _, j := range jalurJenisBisnis {
		if jalur == j {
			if p := strings.TrimSpace(v); p != v {
				k.catat(fmt.Sprintf("spasi dipangkas pada %s (%q → %q) - kandidat perbaikan", jalur, v, p))
				v = p
			}
		}
	}
	return v, ada
}

func (k *kasusTerpangkas) catat(s string) {
	for _, ada := range k.dipangkas {
		if ada == s {
			return
		}
	}
	k.dipangkas = append(k.dipangkas, s)
}
