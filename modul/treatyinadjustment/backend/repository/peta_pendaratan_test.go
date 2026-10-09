package repository

// ⛔ TIAP LARIK YANG KERANGKA BACA HARUS PUNYA BARISNYA DI PETA — kalau tidak,
// gridnya berbunyi `No items` selamanya dan nol yang menunjukkan sebabnya.
//
// Laporan pemilik proses 8 Oktober 2026: *"spreading di adjustment juga
// kosong"*. Panel ▸ baris spreading (`SpreadingTPDtl`: Spread · Currency ·
// Share (%) · Amount) kosong, padahal kerangkanya memuat panel itu dan
// tabelnya sudah terpasang. Yang hilang satu baris di `petaPendaratanPenyesuaian`.

import "testing"

// urutanPeta - indeks entri pertama yang menulis tabel ini.
func urutanPeta(tabel string) int {
	for i, pd := range petaPendaratanPenyesuaian {
		if pd.Tabel == tabel {
			return i
		}
	}
	return -1
}

func TestPecahanSpreadingAdaDiPeta(t *testing.T) {
	induk := urutanPeta("T_TREATY_LIMIT_SPREADING")
	anak := urutanPeta("T_TREATY_LIMIT_SPRD_BREAKDOWN")
	if induk < 0 {
		t.Fatal("T_TREATY_LIMIT_SPREADING tidak ada di peta")
	}
	if anak < 0 {
		t.Fatal("T_TREATY_LIMIT_SPRD_BREAKDOWN tidak ada di peta — panel ▸ baris spreading akan kosong")
	}

	// ⚠️ `bacaPohon` merangkai anak ke induk yang SUDAH terbaca. Entri anak
	// yang mendahului induknya membuat setiap barisnya yatim — dan yatim
	// dihitung diam-diam, tidak digagalkan, jadi gejalanya sama persis
	// dengan entri yang hilang: grid kosong tanpa sepatah galat.
	if anak < induk {
		t.Errorf("pecahan (indeks %d) mendahului induknya (%d) — barisnya akan yatim", anak, induk)
	}

	pd := petaPendaratanPenyesuaian[anak]
	if pd.Induk != "T_TREATY_LIMIT_SPREADING" || pd.KunciAnak != "BreakDownSprdList" {
		t.Errorf("induk/kunci pecahan = %q/%q", pd.Induk, pd.KunciAnak)
	}
	// Kunci = nama properti yang dibaca kerangka `SpreadingTPDtl`.
	mau := map[string]bool{"Amount": true, "Currency": true, "ReinsID": true, "ReinsName": true, "SharePct": true}
	if len(pd.Kunci) != len(mau) {
		t.Fatalf("kunci pecahan %v", pd.Kunci)
	}
	for _, k := range pd.Kunci {
		if !mau[k] {
			t.Errorf("kunci tak dikenal %q", k)
		}
	}
	if len(pd.Kolom) != len(pd.Kunci) {
		t.Errorf("kunci %d ≠ kolom %d", len(pd.Kunci), len(pd.Kolom))
	}
}

// ⭐ Dan induknya benar-benar diakui punya anak — `punyaAnak` dan `larikAnak`
// dihitung dari peta, jadi keduanya ikut hidup begitu entrinya ada. Tanpa ini
// simpulnya tidak pernah disimpan ke `perTabel`, dan anaknya yatim.
func TestIndukSpreadingDiakuiPunyaAnak(t *testing.T) {
	if !punyaAnak["T_TREATY_LIMIT_SPREADING"] {
		t.Fatal("T_TREATY_LIMIT_SPREADING tidak diakui punya anak")
	}
	var ada bool
	for _, n := range larikAnak["T_TREATY_LIMIT_SPREADING"] {
		if n == "BreakDownSprdList" {
			ada = true
		}
	}
	if !ada {
		t.Errorf("larikAnak = %v, mau memuat BreakDownSprdList", larikAnak["T_TREATY_LIMIT_SPREADING"])
	}
}
