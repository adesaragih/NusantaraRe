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

// TanpaKosong - salinan dokumen TANPA nilai kosong: teks kosong, `nil`, dan
// larik kosong dibuang, menurun ke halaman dan baris tertanam.
//
// ⭐ Untuk LAPORAN tombol Save saja (`kunciTakTersimpan`): properti kosong
// tidak membawa apa pun yang dapat hilang, jadi menyebutnya "TIDAK tersimpan"
// hanya derau — 8 Oktober 2026 pembaca pohon menyisipkan sembilan larik
// kosong, dan laporan Save memuat kesembilannya. Dokumen yang DITULIS tidak
// tersentuh.
func TanpaKosong(doc map[string]any) map[string]any {
	out := make(map[string]any, len(doc))
	for k, v := range doc {
		if b, ada := tanpaKosongNilai(v); ada {
			out[k] = b
		}
	}
	return out
}

func tanpaKosongNilai(v any) (any, bool) {
	switch x := v.(type) {
	case nil:
		return nil, false
	case string:
		return x, x != ""
	case map[string]any:
		return TanpaKosong(x), true
	case []any:
		out := make([]any, 0, len(x))
		for _, e := range x {
			if b, ada := tanpaKosongNilai(e); ada {
				out = append(out, b)
			}
		}
		return out, len(out) > 0
	case []map[string]any:
		out := make([]any, 0, len(x))
		for _, e := range x {
			out = append(out, TanpaKosong(e))
		}
		return out, len(out) > 0
	}
	return v, true
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
	elemen := elemenPerTabel(doc)
	out := map[string][]string{}
	for _, p := range PetaPendaratan {
		dikenal := map[string]bool{}
		for _, k := range p.Kunci {
			// ⭐ Jalur BERTITIK dikenal lewat segmen PERTAMANYA.
			// `ValueDifference.RNMShare` memakai kunci akar
			// `ValueDifference`; tanpa ini halaman tertanam itu
			// dilaporkan asing pada setiap dokumen yang punya.
			if i := strings.IndexByte(k, '.'); i > 0 {
				dikenal[k[:i]] = true
				continue
			}
			dikenal[k] = true
		}
		// Kunci yang menampung tabel ANAK bukan kolom yang hilang — ia
		// larik yang mendarat ke tabelnya sendiri.
		for k := range kunciAnak[p.Tabel] {
			dikenal[k] = true
		}
		if p.Akar {
			// Tabel skalar akar melihat SELURUH dokumen, termasuk setiap
			// larik tingkat pertama. Larik itu milik tabel lain.
			for k := range larikTingkatPertama {
				dikenal[k] = true
			}
			// ⭐ Sejak `446` dokumen yang sama mendarat ke DUA tabel akar
			// (`T_TREATY_REVISION`, `T_TREATY_HAZARD_LIMIT`). Kunci milik
			// tabel akar lain bukan kunci yang hilang.
			for k := range kunciSemuaAkar {
				dikenal[k] = true
			}
			// Sisa yang sungguh asing dilaporkan SEKALI, di tabel akar
			// pertama — bukan diulang di setiap tabel akar.
			if p.Tabel != tabelAkarPertama {
				continue
			}
		}
		asing := map[string]bool{}
		for _, el := range elemen[p.Tabel] {
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

// kunciAnak memetakan nama tabel INDUK -> himpunan kunci larik anaknya.
//
// ⛔ DIBANGKITKAN dari peta, bukan ditulis tangan. Bentuk sebelumnya
// mengecualikan SATU kunci anak (`InstallmentList`) lewat tetapan
// `LarikAnakAngsuran`, sehingga ke-13 kunci anak lain — `Detail`,
// `TreatyGroupList`, `MDPList`, `DeductionList`, … — dilaporkan sebagai
// "kunci tanpa kolom" pada hampir setiap dokumen.
//
// ⚠️ Akibatnya bukan sekadar berisik: daftar itu gunanya MEMPERLIHATKAN
// properti baru dari Pega, dan daftar yang selalu memuat dua puluh nama
// palsu membuat nama ke-21 yang sungguhan tidak terlihat siapa pun.
var kunciAnak = func() map[string]map[string]bool {
	m := map[string]map[string]bool{}
	for _, p := range PetaPendaratan {
		if p.Induk == "" {
			continue
		}
		if m[p.Induk] == nil {
			m[p.Induk] = map[string]bool{}
		}
		if len(p.LarikGabung) > 0 {
			for _, n := range p.LarikGabung {
				m[p.Induk][n] = true
			}
			continue
		}
		m[p.Induk][p.KunciAnak] = true
	}
	return m
}()

// kunciSemuaAkar — gabungan kunci SELURUH tabel akar; tabelAkarPertama —
// tempat sisa kunci akar yang asing dilaporkan.
var kunciSemuaAkar, tabelAkarPertama = func() (map[string]bool, string) {
	m, pertama := map[string]bool{}, ""
	for _, p := range PetaPendaratan {
		if !p.Akar {
			continue
		}
		if pertama == "" {
			pertama = p.Tabel
		}
		for _, k := range p.Kunci {
			if i := strings.IndexByte(k, '.'); i > 0 {
				k = k[:i]
			}
			m[k] = true
		}
	}
	return m, pertama
}()

// larikTingkatPertama adalah setiap kunci larik di PUNCAK dokumen yang
// sudah punya tabelnya sendiri — yang harus tidak terlihat oleh tabel
// skalar akar.
var larikTingkatPertama = func() map[string]bool {
	m := map[string]bool{}
	for _, p := range PetaPendaratan {
		if p.Induk != "" || p.Akar {
			continue
		}
		if len(p.LarikGabung) > 0 {
			for _, n := range p.LarikGabung {
				m[n] = true
			}
			continue
		}
		m[p.Larik] = true
	}
	return m
}()

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
//
// ---------------------------------------------------------------------
// ⛔ RALAT 6 Oktober 2026 — fungsi ini BUTA terhadap migrasi 438 dan 439.
// ---------------------------------------------------------------------
// Bentuk sebelumnya hanya mengenal DUA dari empat bentuk yang `MuatKontrak`
// tangani: larik akar (`Larik`) dan larik anak (`Induk` + `KunciAnak`). Dua
// bentuk yang lahir bersama migrasi 438/439 — `LarikGabung` (beberapa larik
// ke satu tabel) dan `Akar` (dokumen itu sendiri sebagai satu baris) —
// jatuh ke cabang `p.Induk == ""` dan dihitung lewat `BarisLarik(doc, "")`,
// yang selalu mengembalikan NOL.
//
// ⚠️ Akibatnya BUKAN galat, dan itu yang membuatnya berbahaya:
// rekonsiliasi akan melaporkan "dokumen memberi 0" untuk sembilan tabel
// yang pemuatnya justru mengisi, lalu SELISIHNYA dibaca sebagai kerusakan
// pemuat. Yang rusak justru pembandingnya.
//
// ⛔ Keempat cabang di bawah kini CERMINAN `MuatKontrak` satu lawan satu,
// dengan urutan `switch` yang sama. Keduanya harus berubah bersama; yang
// menjaganya `TestCacahLarikMengenalKeempatBentuk`.
func CacahLarik(doc map[string]any) map[string]int {
	out := map[string]int{}
	for tabel, baris := range elemenPerTabel(doc) {
		out[tabel] = len(baris)
	}
	return out
}

// elemenPerTabel menapaki dokumen sekali dan mengembalikan elemen yang akan
// mendarat di tiap tabel.
//
// ⛔ SATU penapak untuk `CacahLarik` DAN `KunciTakTerpetakan`. Keduanya dulu
// menapaki pohonnya sendiri-sendiri, dan keduanya buta terhadap bentuk yang
// berbeda — yang pertama terhadap `Akar` dan `LarikGabung`, yang kedua
// terhadap SELURUH tabel anak kecuali butir angsuran. Dua penapak untuk satu
// pohon akan berselisih lagi; yang ini tidak dapat.
//
// ⚠️ Urutan `switch` di bawah CERMINAN `MuatKontrak` satu lawan satu.
// Keduanya harus berubah bersama; yang menjaganya
// `TestCacahLarikMengenalKeempatBentuk`.
func elemenPerTabel(doc map[string]any) map[string][]map[string]any {
	// SEJAJAR urutan peta — anak membacanya untuk menemukan lariknya di
	// dalam induk, jadi induk harus sudah terisi lebih dulu.
	elemen := map[string][]map[string]any{}
	for _, p := range PetaPendaratan {
		var baris []map[string]any
		switch {
		case p.Akar:
			// Dokumennya SENDIRI satu-satunya elemen — `URUTAN` 0.
			baris = []map[string]any{doc}
		case p.Induk == "" && len(p.LarikGabung) > 0:
			for _, nama := range p.LarikGabung {
				baris = append(baris, BarisLarik(doc, nama)...)
			}
		case p.Induk == "":
			baris = BarisLarik(doc, p.Larik)
		case len(p.LarikGabung) > 0:
			for _, el := range elemen[p.Induk] {
				for _, nama := range p.LarikGabung {
					baris = append(baris, larikObjek(el, nama)...)
				}
			}
		default:
			for _, el := range elemen[p.Induk] {
				baris = append(baris, larikObjek(el, p.KunciAnak)...)
			}
		}
		elemen[p.Tabel] = baris
	}
	return elemen
}

// larikObjek mengambil larik bernama `kunci` dari satu elemen, hanya
// anggotanya yang berupa objek.
//
// ⚠️ Anggota yang BUKAN objek dilewati diam-diam, sama seperti di
// `MuatKontrak`. Nol anggota semacam itu di korpus hari ini; yang penting
// kedua tempat memperlakukannya SAMA, sebab beda perlakuan di sini
// muncul sebagai selisih rekonsiliasi yang tak dapat dijelaskan.
func larikObjek(el map[string]any, kunci string) []map[string]any {
	larik, ok := el[kunci].([]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(larik))
	for _, e := range larik {
		if o, ok := e.(map[string]any); ok {
			out = append(out, o)
		}
	}
	return out
}
