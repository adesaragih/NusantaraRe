package models

// Untuk apa berkas ini: HALAMAN POLIS `pyWorkPage.OfferFacIn` - salinan dokumen penawaran fakultatif yang dipilih di
// pop-up Choose Polis.
//
// `[terverifikasi]` `CopyNB_Act` (langkah 3 RDB-List `GetCopyNBForClaim`: `SELECT DATA_JSON FROM JSON_POLIS WHERE
// nopolis = CARI2 AND prodke = CARI3`; langkah 7 Java "copy NB and remove the comments, and fac retro list for facin")
// menyalin dokumen itu UTUH ke `pyWorkPage.OfferFacIn`, lalu `Quotation` / `Policy.Quotation` = `OfferFacIn.QuotationData`
// (langkah 8) dan `Policy` = `OfferFacIn.PolicyData` (langkah 9); langkah 11 `IsFacRetro` / `IsRISlip` = 0, langkah 13
// Property-Remove `OfferFacIn.FacRetro` + `OfferFacIn.OldData`, langkah 14 Property-Remove
// `.PrintRISlip.FacOfferList(1).Object` setiap baris FacRetroList.
//
// ⛔ Halaman polis TIDAK disimpan (spec Bab 17 titik 1, prompt §6 butir 2 "baca dari polis bila XML hanya menyalinnya"):
// kasus menyimpan NOPOLIS + PRODKE, dan setiap muat membaca ulang JSON_POLIS lalu `TerapkanPolis`. Isi Java langkah 7
// tidak diekspor ("remove the comments") - medan sistem `px*`/`pz*` tidak dimuat.

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// AwalanPolis - jalur halaman polis.
const AwalanPolis = "OfferFacIn"

// Jalur polis yang dibaca lintas berkas.
const (
	JalurNoPolis      = AwalanPolis + ".PolicyData.PolicyNo"
	JalurMulaiPolis   = AwalanPolis + ".PolicyData.StartDateTime"
	JalurAkhirPolis   = AwalanPolis + ".PolicyData.EndDateTime"
	JalurProdke       = "ClaimData.Prodke" // CopyNB_Act Param.Prodke (FACINPRODUCTION/JSON_POLIS.PRODKE)
	DaftarCedingCo    = OQ + "CedingCoList"
	DaftarFacRetro    = AwalanPolis + ".FacRetroList"
	DaftarLokasiPolis = AwalanPolis + ".LocationList"
)

// kunciBuangAkar - subpohon akar dokumen polis yang dibuang (CopyNB_Act langkah 13).
var kunciBuangAkar = map[string]bool{"OldData": true, "FacRetro": true}

// TerapkanPolis menulis dokumen JSON polis ke halaman di bawah `OfferFacIn.` (seluruh isi lama halaman polis dibuang
// lebih dulu - Page-Remove CopyNB_Act langkah 6). Nilai skalar menjadi teks apa adanya (angka JSON ditulis ulang
// tanpa float); objek bersarang menjadi jalur bertitik; larik objek menjadi PageList (indeks 1..n), dan larik di dalam
// baris menjadi daftar anak `induk(n).anak`.
func TerapkanPolis(h *Halaman, dokumen []byte) error {
	h.HapusAwalan(AwalanPolis)
	dec := json.NewDecoder(strings.NewReader(string(dokumen)))
	dec.UseNumber()
	var akar map[string]any
	if err := dec.Decode(&akar); err != nil {
		return fmt.Errorf("models: dokumen polis JSON_POLIS tidak terurai: %w", err)
	}
	for k := range kunciBuangAkar {
		delete(akar, k)
	}
	ratakanHalaman(h, AwalanPolis, akar)
	h.Setel(AwalanPolis+".IsFacRetro", "0")        // 11
	h.Setel(AwalanPolis+".IsRISlip", "0")          // 11
	for i := range h.AmbilDaftar(DaftarFacRetro) { // 14.1.1 / 14.2.1
		h.SetelDaftar(JalurAnak(JalurAnak(DaftarFacRetro, i+1, "PrintRISlip.FacOfferList"), 1, "Object"), nil)
	}
	return nil
}

func lewatiKunci(k string) bool {
	return strings.HasPrefix(k, "px") || strings.HasPrefix(k, "pz")
}

// teksJSON - nilai skalar JSON sebagai teks halaman.
func teksJSON(v any) (string, bool) {
	switch x := v.(type) {
	case string:
		return x, true
	case json.Number:
		return x.String(), true
	case bool:
		if x {
			return "true", true
		}
		return "false", true
	case nil:
		return "", true
	}
	return "", false
}

func kunciUrut(m map[string]any) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// ratakanHalaman menulis objek `m` di jalur halaman `jalur`.
func ratakanHalaman(h *Halaman, jalur string, m map[string]any) {
	for _, k := range kunciUrut(m) {
		if lewatiKunci(k) {
			continue
		}
		p := jalur + "." + k
		switch x := m[k].(type) {
		case map[string]any:
			ratakanHalaman(h, p, x)
		case []any:
			ratakanDaftar(h, p, x)
		default:
			if s, ok := teksJSON(x); ok && s != "" {
				h.Setel(p, s)
			}
		}
	}
}

// ratakanDaftar menulis larik objek sebagai PageList `jalur`.
func ratakanDaftar(h *Halaman, jalur string, xs []any) {
	var rows []Baris
	for i, e := range xs {
		m, ok := e.(map[string]any)
		if !ok {
			continue
		}
		b := Baris{}
		ratakanBaris(h, jalur, i+1, b, "", m)
		rows = append(rows, b)
	}
	if len(rows) > 0 {
		h.SetelDaftar(jalur, rows)
	}
}

// ratakanBaris menulis objek `m` ke baris `b` (awalan properti `pre`); larik bersarang menjadi daftar anak baris n.
func ratakanBaris(h *Halaman, daftar string, n int, b Baris, pre string, m map[string]any) {
	for _, k := range kunciUrut(m) {
		if lewatiKunci(k) {
			continue
		}
		p := pre + k
		switch x := m[k].(type) {
		case map[string]any:
			ratakanBaris(h, daftar, n, b, p+".", x)
		case []any:
			ratakanDaftar(h, JalurAnak(daftar, n, p), x)
		default:
			if s, ok := teksJSON(x); ok && s != "" {
				b[p] = s
			}
		}
	}
}

// NormalkanTanggalPolis meniru CopyNB_Act langkah 11 (`@FormatDateTime(StartDateTime / EndDateTime, "yyyyMMdd")`):
// tanggal polis disimpan halaman sebagai tanggal saja ("2006-01-02").
func NormalkanTanggalPolis(h *Halaman) {
	for _, j := range []string{JalurMulaiPolis, JalurAkhirPolis} {
		if t, ok := UraiTanggal(h.Ambil(j)); ok {
			h.Setel(j, TeksTanggal(t))
		}
	}
}

// IsTreatyInDari - CopyNB_Act langkah 4-5: nomor polis "RNM-Q" -> "1", "RNM-F" -> "0", selain itu tidak diubah ("").
func IsTreatyInDari(nopolis string) string {
	switch {
	case strings.Contains(nopolis, "RNM-Q"):
		return "1"
	case strings.Contains(nopolis, "RNM-F"):
		return "0"
	}
	return ""
}
