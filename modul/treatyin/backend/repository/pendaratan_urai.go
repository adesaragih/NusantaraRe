package repository

// Pengurai dokumen warisan untuk pendaratan - NOL Oracle.
//
// ⛔ Seluruh berkas ini fungsi murni. Yang diputuskan di sini - apa yang
// dianggap larik, bagaimana sebuah nilai menjadi teks, kunci mana yang
// belum punya kolom - adalah bagian yang paling mudah salah dan paling
// mahal diuji lewat basis data. Dipisah, ia teruji tanpa koneksi.

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// UraiDokumen mengurai `M_TREATY_IN.JSONDATA`.
//
// ⛔ `UseNumber()`. Tanpa itu `encoding/json` mengubah setiap angka menjadi
// `float64`, dan `float64` MEMBULATKAN: `AmountTotal` terpanjang di korpus
// 61 aksara, jauh di luar 15-17 digit berarti yang `float64` sanggup bawa.
// Pembulatan itu tidak akan pernah terlihat sebagai galat - hanya sebagai
// angka yang berbeda tipis dari sistem lama.
func UraiDokumen(teks string) (map[string]any, error) {
	d := json.NewDecoder(strings.NewReader(teks))
	d.UseNumber()
	var doc map[string]any
	if err := d.Decode(&doc); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJSONWarisanRusak, err)
	}
	return doc, nil
}

// NilaiTeks mengubah satu nilai JSON menjadi teks APA ADANYA.
//
// Mengembalikan `ada=false` untuk `nil`, yang mendarat sebagai NULL - beda
// dengan string kosong, yang mendarat sebagai string kosong.
//
// ⚠️ `json.Number` dikembalikan lewat `String()`, yaitu LITERAL aslinya -
// bukan hasil hitung ulang. Itulah seluruh guna `UseNumber()`.
func NilaiTeks(v any) (string, bool) {
	switch t := v.(type) {
	case nil:
		return "", false
	case string:
		return t, true
	case json.Number:
		return t.String(), true
	case bool:
		if t {
			return "true", true
		}
		return "false", true
	default:
		// Larik maupun objek tidak pernah mendarat sebagai kolom; pemanggil
		// yang menyaringnya. Sampai di sini berarti bentuk yang belum
		// pernah terlihat, dan ia dibawa apa adanya supaya terlihat.
		return fmt.Sprint(t), true
	}
}

// BarisLarik mengambil satu larik puncak sebagai daftar objek.
//
// Larik yang kuncinya tidak ada, bukan larik, atau kosong semuanya
// menghasilkan nol baris - dan ketiganya memang berarti hal yang sama bagi
// pemuat: tidak ada yang mendarat. Yang BERBEDA di antara ketiganya dicatat
// rekonsiliasi, bukan di sini.
func BarisLarik(doc map[string]any, nama string) []map[string]any {
	larik, ok := doc[nama].([]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(larik))
	for _, el := range larik {
		if objek, ok := el.(map[string]any); ok {
			out = append(out, objek)
		}
	}
	return out
}

// KunciTakTerpetakan mendaftar kunci JSON yang TIDAK punya kolom.
//
// ⛔ Ini penjaga terhadap diam. Pega menambah properti tanpa memberi tahu
// siapa pun, dan kunci baru yang tidak punya kolom hilang tanpa galat -
// pemuatnya tetap hijau, cacah barisnya tetap cocok, dan satu medan berhenti
// ada. Dipanggil rekonsiliasi; hasilnya dilaporkan, bukan ditelan.
//
// Kunci `px`/`py` TIDAK dikecualikan: `pxObjClass` dan
// `pyTemplateRichTextEditor` sudah punya kolom, jadi pengecualian akan
// menyembunyikan saudaranya yang berikutnya.
func KunciTakTerpetakan(doc map[string]any) map[string][]string {
	out := map[string][]string{}
	for i, p := range PetaPendaratan {
		var baris []map[string]any
		if i == IndeksButirAngsuran {
			baris = butirAngsuran(doc)
		} else {
			baris = BarisLarik(doc, p.Larik)
		}
		dikenal := map[string]bool{}
		for _, k := range p.Kunci {
			dikenal[k] = true
		}
		if i == IndeksAngsuran {
			dikenal[LarikAnakAngsuran] = true // anak, bukan kolom
		}
		asing := map[string]bool{}
		for _, el := range baris {
			for k := range el {
				if !dikenal[k] {
					asing[k] = true
				}
			}
		}
		if len(asing) == 0 {
			continue
		}
		daftar := make([]string, 0, len(asing))
		for k := range asing {
			daftar = append(daftar, k)
		}
		sort.Strings(daftar)
		out[p.Tabel] = daftar
	}
	return out
}

// butirAngsuran meratakan `Installment[].InstallmentList` menjadi satu
// daftar, dalam urutan baca.
func butirAngsuran(doc map[string]any) []map[string]any {
	var out []map[string]any
	for _, induk := range BarisLarik(doc, PetaPendaratan[IndeksAngsuran].Larik) {
		larik, ok := induk[LarikAnakAngsuran].([]any)
		if !ok {
			continue
		}
		for _, el := range larik {
			if objek, ok := el.(map[string]any); ok {
				out = append(out, objek)
			}
		}
	}
	return out
}

// CacahLarik menghitung elemen tiap larik di SATU dokumen - angka yang
// rekonsiliasi adu dengan cacah baris di tabelnya.
func CacahLarik(doc map[string]any) map[string]int {
	out := map[string]int{}
	for i, p := range PetaPendaratan {
		if i == IndeksButirAngsuran {
			out[p.Tabel] = len(butirAngsuran(doc))
			continue
		}
		out[p.Tabel] = len(BarisLarik(doc, p.Larik))
	}
	return out
}
