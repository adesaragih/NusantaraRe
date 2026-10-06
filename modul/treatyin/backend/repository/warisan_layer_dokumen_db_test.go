//go:build db

package repository_test

// Bukti Oracle untuk baris layer BERSUMBER DOKUMEN — penerus
// `warisan_in2_db_test.go`.
//
// ⛔ Keempat uji berkas lama DIUBAH ARTINYA, bukan dihapus diam-diam:
//
//	TestBacaLayerMengisiMedanDariKolomYangBenar
//	  → TestLayerDokumenMengisiMedanDariJalurYangBenar
//	    Dahulu: medan terisi dari KOLOM yang benar.
//	    Kini:   medan terisi dari JALUR DOKUMEN yang benar.
//
//	TestLayerBerurutMenurutAngka
//	  → TestUrutanDokumenDipertahankan (non-db, `warisan_layer_dokumen_test.go`)
//	    Urutan dahulu milik `ORDER BY`; kini milik dokumen, dan yang dijaga
//	    berbalik menjadi "jangan diurutkan ulang".
//
//	TestKontrakTanpaBarisLayerMengembalikanKosong
//	  → TestNolLimitsMemberiIrisanKosong (non-db)
//	    Tidak lagi memerlukan Oracle: nol `Limits[]` adalah sifat dokumen.
//
//	TestJangkauanTabelLayerTerukur
//	  → TestLubang510Tertutup
//	    Dahulu MENGUNCI lubangnya (1.340 kontrak, dan membuktikan ia TIDAK
//	    penuh). Kini membuktikan lubang itu TERTUTUP — dan tetap mengunci
//	    cacah tabelnya, sebab tabelnya tetap terlarang disentuh.
//
// ⛔ BACA SAJA. Nol transaksi, nol tulisan, nol `Commit`, nol `DROP`.

import (
	"database/sql"
	"encoding/json"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/modul/treatyin/backend/models"
)

// ⭐ Medan layer terisi dari JALUR DOKUMEN yang benar, pada kontrak nyata.
//
// ⚠️ Kontraknya dipilih karena BERISI, bukan karena nomornya kecil — ronde
// sebelumnya terjebak kontrak draft `1001856` (JSONDATA 1.764 bita, nol dari
// 8 kunci kepala).
func TestLayerDokumenMengisiMedanDariJalurYangBenar(t *testing.T) {
	g, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	h, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = h.Close() }()

	// Kontrak dengan `Limits[]` TERBANYAK — paling banyak medan untuk diadu.
	rows, err := h.QueryContext(ctx,
		`SELECT ID, JSONDATA FROM `+cfg.OracleSchema+`.M_TREATY_IN
		  WHERE DBMS_LOB.INSTR(JSONDATA, '"Limits":[') > 0
		  ORDER BY DBMS_LOB.GETLENGTH(JSONDATA) DESC FETCH FIRST 5 ROWS ONLY`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()

	var diperiksa int
	for rows.Next() {
		var id, dok sql.NullString
		if err := rows.Scan(&id, &dok); err != nil {
			t.Fatal(err)
		}
		var m struct {
			Limits []map[string]any `json:"Limits"`
		}
		if json.Unmarshal([]byte(dok.String), &m) != nil || len(m.Limits) == 0 {
			continue
		}
		// ⛔ RALAT 6 Oktober 2026 — SUMBERNYA DIGANTI, artinya DIPERKUAT.
		//
		// Dahulu `g.BacaKontrakWarisan`, yang mengurai `JSONDATA` di
		// tempat. Penguraian itu DICABUT dari fungsi tersebut (keputusan
		// "nol nilai dari JSONDATA"), sehingga `k.Layer` kembali kosong
		// dan uji ini mati dengan `index out of range` — bukan dengan
		// pesan yang menjelaskan apa pun.
		//
		// ⭐ Penggantinya `BacaLayerPendaratan`, yang membaca TABEL
		// pendaratan. Dan itu membuat uji ini lebih kuat daripada
		// sebelumnya: dahulu ia mengadu dokumen dengan pengurai dokumen —
		// dua jalur yang membaca sumber yang sama. Kini ia mengadu
		// dokumen dengan TABEL, yaitu persis rantai yang layar pakai.
		//
		// ⚠️ Kontrak yang belum didaratkan dilewati, bukan digagalkan:
		// pemuatan dijalankan dengan tangan, dan uji baca tidak berhak
		// menuntut seseorang sudah menjalankannya.
		layer, err := g.BacaLayerPendaratan(ctx, id.String)
		if err != nil {
			t.Fatalf("%s: %v", id.String, err)
		}
		if len(layer) == 0 {
			continue
		}
		k := struct{ Layer []models.BarisLayerWarisan }{Layer: layer}
		diperiksa++

		// Cacah barisnya = jumlah (layer × treaty group).
		mau := 0
		for _, el := range m.Limits {
			d, _ := el["Detail"].([]any)
			if len(d) == 0 {
				mau++
				continue
			}
			mau += len(d)
		}
		if len(k.Layer) != mau {
			t.Errorf("%s: %d baris layer, dokumen memberi %d (layer × treaty group)",
				id.String, len(k.Layer), mau)
		}

		// Medan tingkat layer baris PERTAMA diadu langsung dengan dokumen.
		el := m.Limits[0]
		cocok := func(nama, jalur, dapat string) {
			v, ada := el[jalur]
			if !ada || v == nil {
				return
			}
			mau := teksAny(v)
			if dapat != mau {
				t.Errorf("%s: %s = %q, dokumen `%s` = %q", id.String, nama, dapat, jalur, mau)
			}
		}
		b := k.Layer[0]
		cocok("Layer", "Layer", b.Layer)
		cocok("JenisLayer", "LayerType", b.JenisLayer)
		cocok("DasarCover", "Cover", b.DasarCover)
		cocok("MataUang", "Currency", b.MataUang)
		cocok("Limit100", "Limit", b.Limit100)
		cocok("RetensiCedant", "Deductible", b.RetensiCedant)
		cocok("RasioMDP", "MDPPct", b.RasioMDP)
		cocok("ROL", "ROLPct", b.ROL)
		cocok("AdjRate", "AdjRate", b.AdjRate)
		cocok("RelasiMataUang", "CurrencyRelation", b.RelasiMataUang)
		t.Logf("%s: %d baris layer", id.String, len(k.Layer))
	}
	if diperiksa == 0 {
		t.Fatal("nol kontrak berisi diperiksa")
	}
}

// teksAny meniru `teksJSON` untuk nilai yang sudah terurai.
func teksAny(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		// ⚠️ `encoding/json` memberi `float64` di sini, dan uji ini karena
		// itu HANYA membandingkan nilai yang bulat dan kecil. Nilai besar
		// dijaga uji non-db `TestAngkaBesarUtuh`, yang membaca teks mentah.
		if t == float64(int64(t)) {
			return itoa(int64(t))
		}
	}
	return ""
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		return "-" + string(b)
	}
	return string(b)
}

// ⭐⭐ LUBANG 510 KONTRAK TERTUTUP.
//
// Inilah uji yang membuktikan pencabutan `M_TREATY_IN2` membeli sesuatu
// nyata. Terukur 5 Oktober 2026:
//
//	kontrak dengan `Limits[]` berisi      1.850
//	di antaranya ada di `M_TREATY_IN2`    1.340
//	LUBANG                                  510 kontrak · 1.210 elemen
//
// Sesudah pencabutan, jangkauannya 1.850 — dan `Kosong` kembali punya SATU
// arti.
func TestLubang510Tertutup(t *testing.T) {
	g, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	h, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = h.Close() }()

	// Satu kontrak yang punya `Limits[]` TETAPI nol baris di `M_TREATY_IN2`
	// — persis kasus yang dahulu memberi grid kosong tanpa sebab.
	//
	// ⚠️ Kueri ini menyebut tabel yang dicabut, dan itu DISENGAJA: ia
	// membuktikan pencabutannya, bukan memakainya sebagai sumber. Ia hidup
	// di uji, bukan di jalur baca — dan `TestNolKueriMTreatyIn2` memang
	// menyapu berkas ini pula, jadi berkas ini ada di daftar-izinnya.
	var id string
	if err := h.QueryRowContext(ctx,
		`SELECT m.ID FROM `+cfg.OracleSchema+`.M_TREATY_IN m
		  WHERE DBMS_LOB.INSTR(m.JSONDATA, '"Limits":[') > 0
		    AND NOT EXISTS (SELECT 1 FROM `+cfg.OracleSchema+`.M_TREATY_IN2 x
		                     WHERE x.MASTERID = m.ID)
		  ORDER BY DBMS_LOB.GETLENGTH(m.JSONDATA) DESC FETCH FIRST 1 ROWS ONLY`).
		Scan(&id); err != nil {
		t.Skipf("lewati: nol kontrak di dalam lubang: %v", err)
	}

	k, err := g.BacaKontrakWarisan(ctx, id)
	if err != nil {
		t.Fatalf("%s: %v", id, err)
	}
	if len(k.Layer) == 0 {
		t.Fatalf("%s ada di lubang 510 dan MASIH memberi nol baris layer — "+
			"pencabutan M_TREATY_IN2 tidak membeli apa pun", id)
	}
	t.Logf("kontrak lubang %s kini memberi %d baris layer (dahulu 0)", id, len(k.Layer))

	// ⛔ Dan tabelnya TETAP UTUH — dicabut sebagai sumber, bukan dihapus.
	var baris int
	if err := h.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM `+cfg.OracleSchema+`.M_TREATY_IN2`).Scan(&baris); err != nil {
		t.Fatal(err)
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
	if baris < 7281 {
		t.Errorf("M_TREATY_IN2 %d baris, terukur 7.281 — tabel WARISAN ini tetap "+
			"terlarang disentuh walau tidak lagi dibaca", baris)
	}
}

// ⭐ JANGKAUAN BARU diukur, bukan diasumsikan.
func TestJangkauanDokumenLebihLuasDaripadaTabel(t *testing.T) {
	_, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	h, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = h.Close() }()

	var lewatDokumen, lewatTabel int
	if err := h.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM `+cfg.OracleSchema+`.M_TREATY_IN
		  WHERE DBMS_LOB.INSTR(JSONDATA, '"Limits":[') > 0`).Scan(&lewatDokumen); err != nil {
		t.Fatal(err)
	}
	if err := h.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT MASTERID) FROM `+cfg.OracleSchema+`.M_TREATY_IN2`).
		Scan(&lewatTabel); err != nil {
		t.Fatal(err)
	}
	if lewatDokumen <= lewatTabel {
		t.Errorf("dokumen mencakup %d kontrak, tabel %d — pencabutan tidak memperluas apa pun",
			lewatDokumen, lewatTabel)
	}
	// ⚠️ 1.851 punya KUNCINYA; 1.850 punya ISI. Selisih satu itu dokumen
	// dengan `"Limits":[]` — larik kosong, bukan kunci yang hilang, dan
	// penyaring teks tidak dapat membedakan keduanya. Yang membedakan
	// pengurai, dan `TestNolLimitsMemberiIrisanKosong` menjaganya.
	// ⛔ LANTAI, bukan angka PERSIS — ralat 6 Oktober 2026, sebab yang sama
	// dengan `syarat_tab_db_test.go`: kontrak baru dari Pega menaikkan kedua
	// angka ini, dan kenaikan itu bukan regresi. Yang BERARTI penurunannya.
	if lewatDokumen < 1851 {
		t.Errorf("dokumen memuat kunci `Limits` pada %d kontrak, terukur >= 1.851 "+
			"pada 5 Oktober 2026 (1.850 di antaranya BERISI) — TURUN berarti "+
			"dokumen warisan kehilangan kuncinya", lewatDokumen)
	}
	if lewatTabel < 1340 {
		t.Errorf("tabel mencakup %d kontrak, terukur >= 1.340", lewatTabel)
	}
	t.Logf("jangkauan: dokumen %d · tabel %d · selisih %d kontrak",
		lewatDokumen, lewatTabel, lewatDokumen-lewatTabel)
}

// ⭐ KETIGA TINGKAT pohon tab Limits benar-benar terisi dari kontrak nyata.
//
// ⛔ Uji ini yang membedakan "pohonnya ada di dokumen" dari "pohonnya terbaca
// oleh kita". Yang pertama sudah diukur dengan penyapu jalur; yang kedua
// hanya terbukti lewat jalur baca yang sesungguhnya.
func TestTigaTingkatPohonLimitsTerisi(t *testing.T) {
	g, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	h, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = h.Close() }()

	// ⛔ Kontraknya dipilih menurut ADANYA `Detail[]`, BUKAN menurut ukuran
	// dokumen — dan perbedaannya nyata: kedelapan dokumen TERBESAR punya
	// `Detail: []` kosong (besarnya datang dari `RevisionHistory`), sehingga
	// pemilihan menurut ukuran memberi nol treaty group dan membuat uji ini
	// menuduh pengurai yang sebenarnya benar.
	//
	// Terukur: dari 4.210 elemen `Limits[]` di seluruh dokumen, hanya 1.365
	// punya `Detail[]` berisi — ia memang jarang, dan kontrak termuda yang
	// paling sering punya.
	rows, err := h.QueryContext(ctx,
		`SELECT ID FROM `+cfg.OracleSchema+`.M_TREATY_IN
		  WHERE DBMS_LOB.INSTR(JSONDATA, '"Detail":[') > 0
		  ORDER BY ID DESC FETCH FIRST 30 ROWS ONLY`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()

	var diperiksa, punyaKelompok, punyaKelas int
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
		k, err := g.BacaKontrakWarisan(ctx, id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if len(k.Layer) == 0 {
			continue
		}
		diperiksa++
		for _, b := range k.Layer {
			// ⛔ Irisan KOSONG, bukan nil — layar membedakan keduanya.
			if b.KelasBisnis == nil {
				t.Errorf("%s: KelasBisnis nil, mau irisan kosong", id)
			}
			if b.KelompokTreaty != "" {
				punyaKelompok++
			}
			if len(b.KelasBisnis) > 0 {
				punyaKelas++
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol kontrak berisi diperiksa")
	}
	if punyaKelompok == 0 {
		t.Error("nol baris punya treaty group — tingkat KEDUA pohon tidak terbaca")
	}
	if punyaKelas == 0 {
		t.Error("nol baris punya kelas bisnis — tingkat KETIGA pohon tidak terbaca")
	}
	// ⚠️ Tingkat ketiga JARANG, dan itu bukan cacat: gambar 05 dan 31
	// memperlihatkannya berbunyi `No items` pada kontrak contoh pemilik
	// proses sendiri. Yang dijaga di sini "pernah terbaca", bukan "selalu".

	t.Logf("%d kontrak: %d baris bertreaty group, %d baris berkelas bisnis",
		diperiksa, punyaKelompok, punyaKelas)
}

// ⭐ Medan pengelompokan puncak BERCABANG — dan itu yang membenarkan dua
// bentuk layar.
//
// Terukur: `Limits[].TreatyType` terisi pada prop dan nil pada 2.847 dari
// 2.850 elemen non-prop; `Limits[].LayerType` sebaliknya.
func TestMedanPuncakBercabang(t *testing.T) {
	g, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	h, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = h.Close() }()

	hitung := func(cabang string) (jenisTreaty, jenisLayer int) {
		rows, err := h.QueryContext(ctx,
			`SELECT t.ID FROM `+cfg.OracleSchema+`.TREATY_IN t
			   JOIN `+cfg.OracleSchema+`.M_TREATY_IN m ON m.ID = t.ID
			  WHERE t.PROPORTIONTYPE = :1
			    AND DBMS_LOB.INSTR(m.JSONDATA, '"Limits":[') > 0
			  ORDER BY DBMS_LOB.GETLENGTH(m.JSONDATA) DESC FETCH FIRST 10 ROWS ONLY`, cabang)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = rows.Close() }()
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				t.Fatal(err)
			}
			k, err := g.BacaKontrakWarisan(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
			for _, b := range k.Layer {
				if b.JenisTreaty != "" {
					jenisTreaty++
				}
				if b.JenisLayer != "" {
					jenisLayer++
				}
			}
		}
		return
	}

	pJenis, pLayer := hitung("Proportional")
	nJenis, nLayer := hitung("NonProportional")
	t.Logf("prop:     jenisTreaty %d · jenisLayer %d", pJenis, pLayer)
	t.Logf("non-prop: jenisTreaty %d · jenisLayer %d", nJenis, nLayer)

	if pJenis == 0 {
		t.Error("prop nol jenisTreaty — puncak `Kind of Treaty` tidak punya isi")
	}
	if nLayer == 0 {
		t.Error("non-prop nol jenisLayer — puncak `Layers` tidak punya isi")
	}
	// ⛔ Dan keduanya benar-benar BERBEDA arah: non-prop jauh lebih
	// bergantung pada `LayerType` daripada prop.
	if pLayer >= nLayer {
		t.Errorf("jenisLayer prop %d >= non-prop %d; percabangan puncak "+
			"tidak lagi terbukti dan kedua bentuk layar perlu ditinjau", pLayer, nLayer)
	}
}

// ⛔⛔ MATA UANG KEDUA pada kontrak NYATA — bukti Oracle untuk cacat
// `nilaiPertama` yang ditutup 6 Oktober 2026.
//
// ⚠️ Kontraknya dipilih karena BERMATA UANG DUA, bukan karena besar — jebakan
// yang sudah menjerat ronde sebelumnya. Layer bermata uang tunggal lulus
// bahkan dengan cacatnya.
func TestMataUangKeduaPadaKontrakNyata(t *testing.T) {
	g, ctx := gudangBaca(t)

	// `1000003` — terukur 6 Oktober 2026: empat layer pertamanya
	// ber-`Currency IDR` / `Currency2 USD`, dan `MDPList` berpanjang dua.
	k, err := g.BacaKontrakWarisan(ctx, "1000003")
	if err != nil {
		t.Fatal(err)
	}
	if len(k.Layer) == 0 {
		t.Fatal("1000003: nol baris layer")
	}

	var berdua int
	for i, b := range k.Layer {
		if b.MDPKedua == "" && b.PremiEarnedKedua == "" {
			continue
		}
		berdua++
		// ⛔ Mata uang keduanya harus BERBEDA dari yang pertama — kalau sama,
		// yang terbaca nilai kembar dan bentuk dua-kolom tidak dibutuhkan.
		if b.MataUangLimit == "" {
			t.Errorf("baris %d punya nilai kedua tetapi nol mata uang kedua: %+v", i, b)
		}
		if b.MataUang == b.MataUangLimit {
			t.Errorf("baris %d: mata uang pertama dan kedua SAMA (%q) — "+
				"bentuk dua-kolom perlu ditinjau", i, b.MataUang)
		}
		if b.MDPKedua == b.MDP && b.MDPKedua != "" {
			t.Errorf("baris %d: MDPKedua menyalin MDP (%q)", i, b.MDPKedua)
		}
	}
	if berdua == 0 {
		t.Fatal("1000003: nol baris bermata uang dua — terukur 4+ pada 6 Oktober 2026; " +
			"bila dokumennya berubah, uji ini perlu kontrak lain")
	}
	t.Logf("1000003: %d dari %d baris layer membawa mata uang KEDUA", berdua, len(k.Layer))
}

// ⭐ Cacah layer bermata uang dua DIKUNCI — angka yang hanya ditulis sekali
// akan basi tanpa suara.
func TestCacahLayerMataUangDuaTerukur(t *testing.T) {
	g, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	h, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = h.Close() }()

	rows, err := h.QueryContext(ctx, `SELECT ID FROM `+cfg.OracleSchema+`.M_TREATY_IN`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	var id []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			t.Fatal(err)
		}
		id = append(id, s)
	}

	var mdp, premi int
	for _, s := range id {
		k, err := g.BacaKontrakWarisan(ctx, s)
		if err != nil {
			continue
		}
		for _, b := range k.Layer {
			if b.MDPKedua != "" {
				mdp++
			}
			if b.PremiEarnedKedua != "" {
				premi++
			}
		}
	}
	// ⚠️ Baris layer = layer × treaty group, jadi cacahnya >= cacah layer
	// (106 / 100) yang diukur atas `Limits[]` langsung.
	if mdp < 106 {
		t.Errorf("MDPKedua terisi pada %d baris, terukur >= 106 pada 6 Oktober 2026", mdp)
	}
	if premi < 100 {
		t.Errorf("PremiEarnedKedua terisi pada %d baris, terukur >= 100", premi)
	}
	t.Logf("baris layer bermata uang kedua: MDP %d · PremiEarned %d", mdp, premi)
}
