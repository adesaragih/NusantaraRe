package models

// Batas umur dan sum insured produk - `SavePremiumList_Act` langkah 6-8.2,
// dijalankan Validate CSV sebagai PENOLAKAN (keputusan work owner 05-10-2026).
//
// `[terverifikasi]` langkah HIDUP SavePremiumList_Act dan padanannya:
//
//	6-7   ProductNameID -> GetRateProductLife (`PRODUCTINWARD_LIFE`)  -> BatasProduk
//	8.1   Protect Age: `.ENTRY_AGE < MINAGE || .ENTRY_AGE > MAXAGE`
//	      -> `<NAME_OF_INSURED> Age exceeds the limit, at list <idx>`
//	8.2   Protect Sum Insured: `SUM_INSURED < MINSUMINSURED || > MAXSUMINSURED`,
//	      DILEWATI bila `(Type TP || TR) && @contains(RISLIPRNM,"RNML-FL")`
//	      -> `<NAME_OF_INSURED> Sum Insured exceeds the limit, at list <idx>`
//	9     Page-Set-Messages bila ada galat
//	15    Obj-Save (WithErrors=true) - data TETAP tersimpan walau ada galat
//
// ⛔ BEDA DARI PEGA (keputusan work owner 05-10-2026): batas dibaca dari tabel
// flat Master Product Name Life `M_PRODUCTNAME_LIFE` (bukan `PRODUCTINWARD_LIFE`),
// nilai peserta dari BARIS CSV, dan yang melewati batas DITOLAK di Validate CSV
// - Calculate CSV tidak dapat diklik sampai CSV atau master dibetulkan. Pega
// hanya memberi pesan dan tetap menyimpan (langkah 15).
//
// 8.3, 8.5-8.7, 8.9, 11, 13 ter-remark `//`.

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// BatasProduk - batas `M_PRODUCTNAME_LIFE` satu produk; nil = tidak terisi.
type BatasProduk struct {
	MinAge, MaxAge               *apd.Decimal
	MinSumInsured, MaxSumInsured *apd.Decimal
}

// PeriksaBatasProduk menjalankan langkah 8.1 dan 8.2 atas BARIS CSV dan
// mengembalikan penolakannya (kolom ENTRY_AGE / SUM_INSURED, pesan VERBATIM,
// `at list` = nomor baris CSV).
//
// ⚠️ Batas yang kosong, atau nilai peserta yang kosong / tidak terurai, TIDAK
// diperiksa: nilai rusak sudah ditolak validasi bentuk, dan baris di `lewati`
// (sudah ditolak validasi) tidak dicek lagi.
func PeriksaBatasProduk(typePolis, riSlip string, b BatasProduk, baris []BarisUnggah, lewati map[int]bool) []Penolakan {
	tolak := []Penolakan{}
	luar := func(v, min, maks *apd.Decimal) bool {
		if v == nil {
			return false
		}
		return (min != nil && v.Cmp(min) < 0) || (maks != nil && v.Cmp(maks) > 0)
	}
	rentang := func(min, maks *apd.Decimal) string {
		return fmt.Sprintf("%s - %s", AngkaTampil(min), AngkaTampil(maks))
	}
	lewatiSI := TypeRetro(typePolis) && strings.Contains(riSlip, "RNML-FL")
	for _, br := range baris {
		if lewati[br.Nomor] {
			continue
		}
		nama := strings.TrimSpace(br.Nilai["NAME_OF_INSURED"])
		idx := fmt.Sprintf("%d", br.Nomor)
		var umur *apd.Decimal
		if v := strings.TrimSpace(br.Nilai["ENTRY_AGE"]); polaBulat.MatchString(v) {
			umur, _ = utils.ParseDecimal(v)
		}
		if luar(umur, b.MinAge, b.MaxAge) {
			tolak = append(tolak, Penolakan{Baris: br.Nomor, Kolom: "ENTRY_AGE",
				Pesan: nama + " Age exceeds the limit, at list " + idx,
				Sebab: fmt.Sprintf("ENTRY_AGE %s outside product limit %s", AngkaTampil(umur), rentang(b.MinAge, b.MaxAge))})
		}
		if lewatiSI {
			continue
		}
		si, err := UangCSV(br.Nilai["SUM_INSURED"])
		if err != nil {
			si = nil
		}
		if luar(si, b.MinSumInsured, b.MaxSumInsured) {
			tolak = append(tolak, Penolakan{Baris: br.Nomor, Kolom: "SUM_INSURED",
				Pesan: nama + " Sum Insured exceeds the limit, at list " + idx,
				Sebab: fmt.Sprintf("SUM_INSURED %s outside product limit %s", AngkaTampil(si), rentang(b.MinSumInsured, b.MaxSumInsured))})
		}
	}
	return tolak
}

// AngkaTampil - desimal untuk kalimat sebab penolakan, bentuk Indonesia: TITIK
// pemisah ribuan, KOMA desimal (`-27300184` → `-27.300.184`, `1000000000.5` →
// `1.000.000.000,5`) - sama dengan grid layar (permintaan work owner
// 05-10-2026). Hanya tampilan; nil → teks kosong.
func AngkaTampil(d *apd.Decimal) string {
	t := utils.FormatDecimal(d)
	if t == "" {
		return ""
	}
	tanda := ""
	if strings.HasPrefix(t, "-") {
		tanda, t = "-", t[1:]
	}
	bulat, pecahan, _ := strings.Cut(t, ".")
	var b strings.Builder
	for i, r := range bulat {
		if i > 0 && (len(bulat)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(r)
	}
	if pecahan != "" {
		b.WriteString("," + pecahan)
	}
	return tanda + b.String()
}
