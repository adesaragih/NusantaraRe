package lintasmodul

// Peta pendaratan modul Treaty In dan salinan baca modul Adjustment WAJIB
// menyebut hal yang sama.
//
// ---------------------------------------------------------------------
// ⛔ MENGAPA UJI INI ADA, dan mengapa ia di `uji/`
// ---------------------------------------------------------------------
// Sejak 6 Oktober 2026 layar Adjustment membaca TABEL PENDARATAN milik modul
// Treaty In — keputusan pemilik proses melarang keras menarik nilai dari
// `JSONDATA`. Tetapi penjaga arsitektur melarang modul mengimpor modul lain,
// jadi Adjustment memegang SALINAN peta itu
// (`repository.petaPendaratanPenyesuaian`).
//
// ⭐ Salinan yang tidak diadu akan basi, dan basinya DIAM: kolom yang berganti
// nama di Treaty In membuat kolom Adjustment kosong di layar, bukan gagal.
// Uji ini yang membuat hari itu gagal di sini — dan ia hidup di `uji/`, satu-
// satunya lapisan yang boleh mengenal kedua modul.

import (
	"reflect"
	"sort"
	"testing"

	ti "nusantarare/modul/treatyin/backend/repository"
	adj "nusantarare/modul/treatyinadjustment/backend/repository"
)

// Akhiran sisi `Old` ditulis pemuat Treaty In dan dibaca Adjustment. Dua
// tetapan yang berbeda terbaca sebagai "sisi Old kosong", bukan sebagai galat.
func TestAkhiranSisiLamaSama(t *testing.T) {
	if ti.AkhiranSisiLama != adj.AkhiranSisiLama {
		t.Fatalf("akhiran sisi lama berselisih: treatyin %q, adjustment %q",
			ti.AkhiranSisiLama, adj.AkhiranSisiLama)
	}
	if ti.AkhiranSisiLama == "" {
		t.Fatal("akhiran sisi lama kosong; kedua sisi akan mendarat di MASTERID yang sama")
	}
}

// Tiap entri salinan Adjustment harus cocok PERSIS dengan entri bernama sama
// di peta Treaty In — tabel, kunci JSON, dan nama kolomnya.
func TestSalinanPetaAdjustmentCocokDenganTreatyIn(t *testing.T) {
	asli := map[string]ti.Pendaratan{}
	for _, p := range ti.PetaPendaratan {
		asli[p.Tabel] = p
	}
	salinan := adj.PetaPendaratanPenyesuaianUntukUji()
	if len(salinan) == 0 {
		t.Fatal("salinan peta Adjustment kosong; pembacanya yang rusak")
	}
	for _, c := range salinan {
		p, ada := asli[c.Tabel]
		if !ada {
			t.Errorf("Adjustment membaca %s yang TIDAK ada di peta Treaty In — "+
				"tabelnya tidak akan pernah terisi", c.Tabel)
			continue
		}
		if !reflect.DeepEqual(c.Kunci, p.Kunci) {
			t.Errorf("%s: kunci JSON berselisih\n  adjustment %v\n  treatyin   %v",
				c.Tabel, c.Kunci, p.Kunci)
		}
		if !reflect.DeepEqual(c.Kolom, p.Kolom) {
			t.Errorf("%s: nama kolom berselisih\n  adjustment %v\n  treatyin   %v",
				c.Tabel, c.Kolom, p.Kolom)
		}
		if c.Akar != p.Akar {
			t.Errorf("%s: penanda Akar berselisih (%v lawan %v)", c.Tabel, c.Akar, p.Akar)
		}
		if c.Larik != p.Larik {
			t.Errorf("%s: nama larik berselisih (%q lawan %q)", c.Tabel, c.Larik, p.Larik)
		}
		gab, pGab := append([]string{}, c.Gabung...), append([]string{}, p.LarikGabung...)
		sort.Strings(gab)
		sort.Strings(pGab)
		if !reflect.DeepEqual(gab, pGab) {
			t.Errorf("%s: daftar larik gabungan berselisih\n  adjustment %v\n  treatyin   %v",
				c.Tabel, gab, pGab)
		}
	}
}

// ⛔ Dan uji di atas MENGGIGIT: ia harus sungguh membandingkan sesuatu.
func TestSalinanPetaAdjustmentSungguhDibandingkan(t *testing.T) {
	salinan := adj.PetaPendaratanPenyesuaianUntukUji()
	kolom := 0
	for _, c := range salinan {
		kolom += len(c.Kolom)
	}
	if len(salinan) < 5 || kolom < 20 {
		t.Fatalf("salinan terlalu kecil untuk menjaga apa pun: %d tabel, %d kolom",
			len(salinan), kolom)
	}
}

// ⭐ 8 Oktober 2026 — halaman `ActualValue` (cabang Adjust Premium) mendarat di
// `ID + AkhiranSisiAktual`; penulis (Treaty In) dan pembaca (Adjustment)
// wajib sepakat, atau sisi Actual terbaca kosong tanpa galat.
func TestAkhiranSisiAktualSama(t *testing.T) {
	if ti.AkhiranSisiAktual != adj.AkhiranSisiAktual {
		t.Fatalf("akhiran sisi Actual berbeda: treatyin %q, treatyinadjustment %q",
			ti.AkhiranSisiAktual, adj.AkhiranSisiAktual)
	}
	if ti.AkhiranSisiAktual == "" || ti.AkhiranSisiAktual == ti.AkhiranSisiLama {
		t.Fatalf("akhiran sisi Actual %q wajib terisi dan berbeda dari sisi Old %q", ti.AkhiranSisiAktual, ti.AkhiranSisiLama)
	}
}

// ⭐ 8 Oktober 2026 — halaman `ValueBeforeProrate` mendarat di
// `ID + AkhiranSisiSebelumProrata`, pola yang sama dengan sisi Actual.
//
// ⛔ Ketiga akhiran wajib BERBEDA satu sama lain: dua sisi berakhiran sama
// akan saling menimpa di tabel yang sama, dan yang hilang tidak akan
// bergalat — ia hanya terbaca kosong.
func TestAkhiranSisiSebelumProrataSama(t *testing.T) {
	if ti.AkhiranSisiSebelumProrata != adj.AkhiranSisiSebelumProrata {
		t.Fatalf("akhiran sisi ValueBeforeProrate berbeda: treatyin %q, treatyinadjustment %q",
			ti.AkhiranSisiSebelumProrata, adj.AkhiranSisiSebelumProrata)
	}
	akhiran := []string{ti.AkhiranSisiLama, ti.AkhiranSisiAktual, ti.AkhiranSisiSebelumProrata}
	for i, a := range akhiran {
		if a == "" {
			t.Fatalf("akhiran ke-%d kosong", i)
		}
		for j, b := range akhiran {
			if i != j && a == b {
				t.Fatalf("akhiran ke-%d dan ke-%d sama (%q) — kedua sisi akan saling menimpa", i, j, a)
			}
		}
	}
}
