//go:build db

package repository_test

// Bukti Oracle untuk §0 dan §1 ronde 5 Oktober 2026 — syarat tampil tab, dan
// `AccountingModeNonProp` sebagai properti KEDUA.
//
// ⛔ BACA SAJA. Nol transaksi, nol tulisan, nol `Commit`, nol `DROP`.
//
// ⚠️ Yang dijaga di sini bukan "kodenya jalan" melainkan ANGKANYA. Setiap
// angka di bawah diukur dari POOLDATA pada 5 Oktober 2026; uji ini merah
// pada hari angkanya berubah, dan hari itu keputusannya perlu ditinjau,
// bukan angkanya diperbarui diam-diam.

import (
	"database/sql"
	"encoding/json"
	"testing"

	"nusantarare/inti/backend/config"
)

// Kontrak yang briefing minta dipakai sebagai bukti.
const kontrakBukti = "1001856"

// Kelima kontrak yang `IsMultipleRetro`-nya `"true"` — terukur, bukan dipilih.
var kontrakRetroBerganda = []string{"1000493", "1000755", "1001043", "1001405", "1001853"}

// ⭐ Keempat kunci baru benar-benar TERBACA dari dokumen.
func TestKunciSyaratTabTerbaca(t *testing.T) {
	g, ctx := gudangBaca(t)

	for _, id := range kontrakRetroBerganda {
		k, err := g.BacaKontrakWarisan(ctx, id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if k.RetroBerganda != "true" {
			t.Errorf("%s: RetroBerganda %q, mau \"true\"", id, k.RetroBerganda)
		}
		if !k.AdaDiJSON["IsMultipleRetro"] {
			t.Errorf("%s: AdaDiJSON[IsMultipleRetro] false padahal nilainya terbaca", id)
		}
		// ⛔ Kelimanya NON-proporsional, dan itu penting: tab `Retro`
		// bersyarat hanya di cabang PROPORSIONAL. Jadi kelima kontrak ini
		// justru TIDAK memakai syarat itu di layar.
		t.Logf("%s sifat=%q retro=%q", id, k.SifatProporsiAsli, k.RetroBerganda)
	}
}

// ⛔ Kontrak bukti yang briefing tunjuk adalah DRAFT, dan itu menjelaskan
// tangkapan layarnya.
//
// ⚠️ Terukur: `JSONDATA`-nya 1.764 bita, sementara kontrak nyata 33.000 -
// 56.000 bita, dan NOL dari delapan kunci kepala ada di dalamnya. Tangkapan
// layar yang kosong di sana kosong karena KONTRAKNYA kosong — bukan karena
// layarnya salah membaca.
func TestKontrakBuktiAdalahDraftKosong(t *testing.T) {
	g, ctx := gudangBaca(t)
	k, err := g.BacaKontrakWarisan(ctx, kontrakBukti)
	if err != nil {
		t.Fatalf("%s: %v", kontrakBukti, err)
	}
	for _, kunci := range []string{
		"Bordeaux", "AccountingMode", "AccountingModeNonProp",
		"IsMultipleRetro", "EDMState", "EDMMaterialType",
	} {
		if k.AdaDiJSON[kunci] {
			t.Errorf("%s: kunci %q ADA — kontrak ini terukur nol kunci pada 5 Oktober 2026; "+
				"bila ia kini terisi, bukti §0 perlu kontrak lain", kontrakBukti, kunci)
		}
	}
	if k.SifatProporsiAsli != "NonProportional" {
		t.Errorf("%s: sifat %q, terukur NonProportional", kontrakBukti, k.SifatProporsiAsli)
	}
}

// ⭐ `AccountingModeNonProp` ADA, dan ia BUKAN salinan `AccountingMode`.
//
// ⛔ Inilah yang membuktikan layar non-prop sebelumnya membaca properti yang
// salah: kedua properti itu ada di dokumen yang sama dengan nilai yang
// berbeda, jadi menampilkan yang satu di tempat yang lain menampilkan nilai
// yang sungguh keliru — bukan sekadar nama medan yang berbeda.
func TestDuaPropertiCaraPembukuanBerbeda(t *testing.T) {
	g, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	h, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = h.Close() }()

	rows, err := h.QueryContext(ctx,
		`SELECT t.ID FROM `+cfg.OracleSchema+`.TREATY_IN t
		  WHERE t.PROPORTIONTYPE = 'NonProportional'
		  ORDER BY t.ID FETCH FIRST 40 ROWS ONLY`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()

	var berbeda, diperiksa int
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
		k, err := g.BacaKontrakWarisan(ctx, id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if !k.AdaDiJSON["AccountingModeNonProp"] {
			continue
		}
		diperiksa++
		if k.CaraPembukuanNonProp != k.CaraPembukuan {
			berbeda++
		}
		// Domainnya terukur: `loss` 1.830 · `risk` 21 pada 1.854 dokumen.
		if k.CaraPembukuanNonProp != "loss" && k.CaraPembukuanNonProp != "risk" {
			t.Errorf("%s: AccountingModeNonProp %q di luar {loss,risk} yang terukur",
				id, k.CaraPembukuanNonProp)
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol kontrak non-prop punya AccountingModeNonProp — terukur 772 dari 775")
	}
	if berbeda == 0 {
		t.Errorf("%d kontrak diperiksa dan NOL yang kedua propertinya berbeda; "+
			"bila benar, medan kedua di layar tidak lagi membeli apa pun", diperiksa)
	}
	t.Logf("%d kontrak non-prop diperiksa, %d berbeda nilai antara kedua properti",
		diperiksa, berbeda)
}

// ⛔ `EDMMaterialType` NOL di seluruh 1.854 dokumen — jadi syarat tab
// `Value Difference` tidak pernah terpenuhi lewat JSONDATA.
//
// ⚠️ Itu TIDAK berarti tabnya mati. Pega menilai syaratnya atas clipboard,
// yang dapat diisi Activity saat jalan; yang terukur hanya bahwa nilainya
// tidak TERSIMPAN. Perbedaan itu adalah pertanyaan terbuka, bukan kesimpulan.
func TestEdmJenisMaterialNolDiSeluruhDokumen(t *testing.T) {
	cfg, _ := config.Load()
	h, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = h.Close() }()

	_, ctx := gudangBaca(t)
	rows, err := h.QueryContext(ctx,
		`SELECT JSONDATA FROM `+cfg.OracleSchema+`.M_TREATY_IN`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()

	var total, punyaMaterial, punyaState int
	for rows.Next() {
		var dok sql.NullString
		if err := rows.Scan(&dok); err != nil {
			t.Fatal(err)
		}
		total++
		if !dok.Valid || dok.String == "" {
			continue
		}
		var m map[string]json.RawMessage
		if err := json.Unmarshal([]byte(dok.String), &m); err != nil {
			continue
		}
		if _, ada := m["EDMMaterialType"]; ada {
			punyaMaterial++
		}
		if _, ada := m["EDMState"]; ada {
			punyaState++
		}
	}
	// ⛔ LANTAI, bukan angka PERSIS — ralat 6 Oktober 2026.
	//
	// Angka persis membuat gerbang ini merah setiap kali bisnis berjalan
	// normal: Pega menambah baris ke tabel warisan ini, dan penambahan itu
	// bukan regresi. Terukur hari itu — `TREATY_IN` dan `M_TREATY_IN` sudah
	// 1.855 dari 1.854 (kontrak `1002059`), `M_TREATY_IN_DETAIL` 27.618 dari
	// 27.617.
	//
	// ⚠️ YANG HILANG karenanya, dan itu dinyatakan: penambahan baris oleh
	// MODUL INI tidak lagi tertangkap cacahnya. Yang menangkapnya
	// `TestWarisanHanyaDibaca`, yang menyapu naskah SQL dan menolak
	// `INSERT`/`UPDATE`/`DELETE` terhadap tabel ini — penjaga yang
	// sesungguhnya, dan yang 6 Oktober 2026 diperluas supaya mencakupnya.
	if total < 1854 {
		t.Errorf("%d dokumen, terukur >= 1.854 pada 6 Oktober 2026 — TURUN berarti baris warisan hilang", total)
	}
	if punyaMaterial != 0 {
		t.Errorf("EDMMaterialType ada di %d dokumen; terukur NOL pada 5 Oktober 2026 — "+
			"bila kini terisi, syarat tab Value Difference perlu diukur ulang", punyaMaterial)
	}
	if punyaState != 2 {
		t.Errorf("EDMState ada di %d dokumen, terukur 2", punyaState)
	}
	t.Logf("%d dokumen: EDMMaterialType %d, EDMState %d", total, punyaMaterial, punyaState)
}
