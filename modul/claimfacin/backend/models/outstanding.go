package models

// Untuk apa berkas ini: POP-UP "Outstanding Summary" baris objek (`GetAllData_Act` + harness Outstanding /
// Outstanding_SC): baris OS_AKSEPTASI_KLAIM nomor klaim ini (STS 0 estimasi, STS 1 akseptasi) dengan TSI RNM item, kurs
// standar per mata uang, nilai konversi, dan total.
//
// RDB-List GetDataOutstanding_SQL / GetAcceptation_SQL / GetDataOutstandingMBU_Sql / GetAllDataTravel_Act /
// GetAllDataPA_SQL membaca DATA_JSON dengan dot-notation Oracle dalam beberapa bentuk (halaman datar,
// `osAkseptasi`, `AcceptationList`, `CurrencyList[*].AcceptationList[*]`). Di sini repository hanya membaca baris
// (NOCLAIM, NOPOLIS, TANGGAL, STS_REJECT, DATA_JSON; CASEID dua bentuk - prompt §6 butir 6) dan penguraian bentuknya
// dikerjakan Go (`RingkasOutstanding`) - satu tempat untuk tujuh varian SQL.

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// BarisOSTersimpan - satu baris OS_AKSEPTASI_KLAIM kasus yang dibaca pop-up.
type BarisOSTersimpan struct {
	NoClaim, NoPolis, Tanggal, StsReject, DataJSON string
}

// BarisRingkasan - satu baris grid Outstanding_SC.
type BarisRingkasan struct {
	NoClaim    string `json:"noClaim"`
	NoPolis    string `json:"noPolis"`
	AcceptedNo string `json:"acceptedNo"`
	Nama       string `json:"nama"`
	Coverage   string `json:"coverage"`
	Currency   string `json:"currency"`
	CurrencyID string `json:"-"`
	TSIRNM     string `json:"tsiRnm"`
	Nilai      string `json:"nilai"`
	Kurs       string `json:"kurs"`
	Konversi   string `json:"konversi"`
}

// Ringkasan - isi pop-up Outstanding.
type Ringkasan struct {
	// LabelNama - judul kolom nama per lini ("Object Item Name" / "Trading Name" / "Brand Name" / "Object Name").
	LabelNama string           `json:"labelNama"`
	Baris     []BarisRingkasan `json:"baris"`
	Total     string           `json:"total"`
	// MataUang - `TempOutstanding...Currency := IDR` (langkah 19).
	MataUang string `json:"mataUang"`
}

// medanNama - nama kolom DATA_JSON yang menjadi kolom "nama" grid per lini dan status.
func medanNama(h *Halaman, sts string) (string, string) {
	switch {
	case IsMarineCargo(h):
		return "TradingName", "Trading Name"
	case IsMBU(h):
		return "BrandName", "Brand Name"
	case IsTravel(h) || IsPA(h):
		return "ObjectName", "Object Name"
	}
	if sts == "1" { // GetDataOutstanding_SQL cabang STS 1: ObjectItemName
		return "ObjectItemName", "Object Item Name"
	}
	return "ObjectName", "Object Item Name"
}

// RingkasOutstanding = GetAllData_Act. Perbaikan (PARITAS `[penyimpangan sadar]`):
//   - langkah 18.2 mencocokkan mata uang baris Travel / PA pada `.CARI15` yang tidak dipilih SQL keduanya (mata uangnya
//     di `.CARI8`) - kurs dan konversi Travel / PA selalu kosong; di sini memakai mata uang baris;
//   - kolom Currency grid Marine Cargo menampilkan `.CARI10` yang tidak dipilih SQL (selalu kosong); di sini mata uang
//     baris.
func RingkasOutstanding(k *Konteks, h *Halaman, rows []BarisOSTersimpan) (Ringkasan, error) {
	_, label := medanNama(h, "0")
	out := Ringkasan{LabelNama: label, MataUang: "IDR"}
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].Tanggal < rows[j].Tanggal })
	for _, r := range rows {
		if r.StsReject != "0" && r.StsReject != "1" {
			continue
		}
		var akar map[string]any
		dec := json.NewDecoder(strings.NewReader(r.DataJSON))
		dec.UseNumber()
		if err := dec.Decode(&akar); err != nil {
			continue // DATA_JSON tak terurai: dot-notation Oracle pun mengembalikan NULL
		}
		nama, _ := medanNama(h, r.StsReject)
		tambah := func(m map[string]any, induk map[string]any, noAksep string) {
			cov, nilai := teksDari(m, "CoverageName"), teksDari(m, "Value")
			if cov == "" || nilai == "" {
				return
			}
			b := BarisRingkasan{NoClaim: r.NoClaim, NoPolis: r.NoPolis, AcceptedNo: noAksep, Nama: teksDari(m, nama),
				Coverage: cov, Currency: teksDari(m, "Currency"), CurrencyID: teksDari(m, "CurrencyID"), Nilai: nilai}
			if b.Currency == "" && induk != nil {
				b.Currency, b.CurrencyID = teksDari(induk, "Currency"), teksDari(induk, "CurrencyID")
			}
			out.Baris = append(out.Baris, b)
		}
		for _, m := range objekDari(akar["osAkseptasi"]) {
			tambah(m, nil, teksDari(m, "AcceptedNo"))
		}
		tambah(akar, nil, teksDari(akar, "AcceptedNo"))
		if r.StsReject == "1" {
			for _, c := range objekDari(akar["CurrencyList"]) {
				for _, a := range objekDari(c["AcceptationList"]) {
					tambah(a, c, teksDari(akar, "AcceptedNo"))
				}
			}
			if IsMarineCargo(h) {
				for _, a := range objekDari(akar["AcceptationList"]) {
					tambah(a, nil, teksDari(akar, "AcceptedNo"))
				}
			}
		}
	}
	var kk Kalkulator
	for i := range out.Baris { // 11-14: nilai akseptasi negatif
		b := &out.Baris[i]
		if b.AcceptedNo != "" {
			b.Nilai = Teks(kk.Neg(kk.Teks("Value", b.Nilai)))
		}
	}
	if err := kk.Galat(); err != nil {
		return out, err
	}
	isiTSI(h, out.Baris) // 16
	kurs := map[string]string{}
	total := apd.New(0, 0)
	for i := range out.Baris { // 17-18
		b := &out.Baris[i]
		if b.CurrencyID == "" {
			continue
		}
		v, ada := kurs[b.CurrencyID]
		if !ada {
			var err error
			if v, err = k.Acuan.KursStandar(k.Ctxt(), b.CurrencyID); err != nil {
				return out, err
			}
			kurs[b.CurrencyID] = v
		}
		b.Kurs = v
		konv := kk.Kali(kk.Teks("Value", b.Nilai), kk.Teks("KURS", v))
		b.Konversi = Teks(konv)
		total = kk.Tambah(total, konv)
	}
	if err := kk.Galat(); err != nil {
		return out, err
	}
	out.Total = Teks(total) // 19
	return out, nil
}

// isiTSI = GetAllData_Act 16: TSI RNM baris = TSINusare item halaman yang cocok - estimasi (tanpa nomor akseptasi):
// nama item, mata uang, dan nilai estimasi sama; akseptasi: nama item sama dan nomor akseptasi adjustment item sama.
// Kecocokan terakhir menang (loop menimpa).
func isiTSI(h *Halaman, rows []BarisRingkasan) {
	for o := range h.AmbilDaftar(DaftarObjek) {
		for i, it := range h.AmbilDaftar(DaftarItem(o + 1)) {
			for _, e := range h.AmbilDaftar(DaftarDiItem(o+1, i+1, AnakEstimasi)) {
				for r := range rows {
					b := &rows[r]
					if b.AcceptedNo == "" && it["ObjectItemName"] == b.Nama && e["Currency"] == b.Currency &&
						samaAngka(e["EstimationValue"], b.Nilai) {
						b.TSIRNM = it["TSINusare"]
					}
				}
			}
			for _, a := range h.AmbilDaftar(DaftarDiItem(o+1, i+1, AnakAdj)) {
				for r := range rows {
					b := &rows[r]
					if b.AcceptedNo != "" && it["ObjectItemName"] == b.Nama && a["AcceptedNo"] == b.AcceptedNo {
						b.TSIRNM = it["TSINusare"]
					}
				}
			}
		}
	}
}

// objekDari - objek JSON tunggal atau larik objek.
func objekDari(v any) []map[string]any {
	switch x := v.(type) {
	case map[string]any:
		return []map[string]any{x}
	case []any:
		var out []map[string]any
		for _, e := range x {
			if m, ok := e.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out
	}
	return nil
}

// teksDari - nilai skalar kunci `k` sebagai teks ("" bila tidak ada / bukan skalar).
func teksDari(m map[string]any, k string) string {
	if m == nil {
		return ""
	}
	s, ok := teksJSON(m[k])
	if !ok {
		return ""
	}
	return s
}
