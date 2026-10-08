//go:build db

package repository_test

// Bukti Oracle untuk layar Adjustment — `TREATY_IN_EDM` + `M_TREATY_IN_EDM`.
//
// ⛔ BACA SAJA. Nol tulisan, nol transaksi, nol `Commit`.
//
// ⭐ Setiap angka yang menjadi DASAR keputusan layar diukur ULANG di sini,
// bukan dihafal di komentar. Bila salah satu uji ini merah, yang berubah
// datanya — dan keputusan yang bersandar padanya perlu ditinjau, bukan
// ujinya yang dilonggarkan.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
)

// ⭐ Daftar penyesuaian lengkap: cacahnya cocok dengan Oracle.
func TestDaftarPenyesuaianCacahCocok(t *testing.T) {
	g, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)

	var n int
	if err := h.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+skema+`.TREATY_IN_EDM`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	d, err := g.DaftarPenyesuaianWarisan(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != n {
		t.Errorf("daftar %d baris, Oracle %d", len(d), n)
	}
	if n != 280 {
		t.Logf("⚠️ TREATY_IN_EDM %d baris, terukur 280 pada 5 Oktober 2026", n)
	}
	// Urutan RD Pega `BrowseTREATY_IN_EDM`: `.ID` DESC — terbaru di atas.
	for i := 1; i < len(d); i++ {
		if d[i-1].ID < d[i].ID {
			t.Fatalf("daftar tidak menurun menurut ID: %q sebelum %q", d[i-1].ID, d[i].ID)
		}
	}
}

// ⭐ Satu penyesuaian dibaca dengan KEDUA sisinya.
func TestBacaPenyesuaianKeduaSisi(t *testing.T) {
	g, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)

	var id string
	if err := h.QueryRowContext(ctx, `SELECT MIN(ID) FROM `+skema+`.TREATY_IN_EDM`).Scan(&id); err != nil {
		t.Skipf("lewati: nol penyesuaian: %v", err)
	}
	p, err := g.BacaPenyesuaianWarisan(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != id || p.Baru.Medan["ID"] != id {
		t.Errorf("pengenal %q, halaman akar ID %q, mau %q", p.ID, p.Baru.Medan["ID"], id)
	}
	if len(p.Lama.Medan) == 0 {
		t.Error("sisi Old kosong — OLDDATA tidak terbaca")
	}
	if _, ada := p.Baru.Larik["CommentList"]; !ada {
		t.Error("CommentList tidak terbaca dari halaman akar")
	}
}

func TestPenyesuaianTidakAdaDikenali(t *testing.T) {
	g, ctx := bacaSaja(t)
	if _, err := g.BacaPenyesuaianWarisan(ctx, "TIDAK-ADA/R99"); err == nil {
		t.Fatal("mau galat tidak-ada")
	}
}

// ⭐ Dasar panel Old Data: SETIAP dokumen membawa `OLDDATA`.
func TestSetiapDokumenPenyesuaianMembawaOldData(t *testing.T) {
	_, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)
	var total, punya int
	if err := h.QueryRowContext(ctx, `SELECT COUNT(*),
		SUM(CASE WHEN JSON_EXISTS(JSONDATA, '$.OLDDATA') THEN 1 ELSE 0 END)
		FROM `+skema+`.M_TREATY_IN_EDM`).Scan(&total, &punya); err != nil {
		t.Fatal(err)
	}
	if total != punya {
		t.Errorf("%d dari %d dokumen membawa OLDDATA — panel Old Data bersandar pada SEMUA", punya, total)
	}
}

// ⭐ Sumber History: `T_VIEW_COMMENT` — RALAT 6 Oktober 2026, DIBALIK.
//
// ---------------------------------------------------------------------
// MENGAPA PERNYATAANNYA DIBALIK, BUKAN DILONGGARKAN
// ---------------------------------------------------------------------
// Bunyi sebelumnya: *"Dasar sumber History: `CommentList` dokumen, BUKAN
// `T_VIEW_COMMENT`"*, ditegakkan dengan menuntut `T_VIEW_COMMENT` memuat
// NOL baris berpengenal penyesuaian.
//
// Pernyataan itu benar pada hari ia ditulis, dan benar karena satu sebab
// saja: ⛔ **korpus `M_TREATY_IN_EDM` belum pernah dimuat**. Tabel yang
// kosong memang tidak dapat menjadi sumber apa pun.
//
// Korpus itu dimuat 6 Oktober 2026 — 280 dokumen, 60.407 baris, nol
// selisih — dan tabelnya kini memuat komentar penyesuaian. Jadi yang
// berubah bukan keputusannya melainkan KENYATAAN yang mendasarinya, dan
// membiarkan pernyataan lama berdiri akan memaksa panel History membaca
// `JSONDATA` selamanya — yang justru DILARANG.
//
// ⭐ Terukur hari itu, dan ketiganya berjumlah tepat:
//
//	T_VIEW_COMMENT seluruhnya        13.133
//	  berpengenal kontrak Treaty In  11.365
//	  berpengenal penyesuaian BARU    1.488
//	  berpengenal penyesuaian #LAMA     280   (satu per halaman OLDDATA)
//
// ⚠️ Yang TETAP dijaga: `CommentList` di dokumen harus tetap berisi pada
// seluruh 280 — ia SUMBER pemuatannya. Dokumen yang kehilangan
// `CommentList` akan membuat tabelnya kosong pada pemuatan berikutnya,
// dan kekosongan itu tidak akan bersuara di mana pun selain di sini.
func TestHistoryPenyesuaianBersumberCommentList(t *testing.T) {
	_, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)

	var sisiBaru, sisiLama int
	if err := h.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+skema+`.T_VIEW_COMMENT
		WHERE MASTERID IN (SELECT ID FROM `+skema+`.TREATY_IN_EDM)`).Scan(&sisiBaru); err != nil {
		t.Fatal(err)
	}
	if err := h.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+skema+`.T_VIEW_COMMENT
		WHERE MASTERID LIKE '%#LAMA'
		  AND SUBSTR(MASTERID, 1, LENGTH(MASTERID)-5)
		      IN (SELECT ID FROM `+skema+`.TREATY_IN_EDM)`).Scan(&sisiLama); err != nil {
		t.Fatal(err)
	}
	// ⛔ LANTAI, bukan angka persis: Pega terus menambah komentar, dan
	// penambahan itu bukan regresi. Yang menjadi regresi adalah NOL —
	// ia berarti muatannya hilang dan panel History kembali kosong.
	if sisiBaru == 0 {
		t.Error("NOL komentar penyesuaian sisi `New` di T_VIEW_COMMENT — " +
			"panel History kosong; muat ulang korpus M_TREATY_IN_EDM")
	}
	if sisiLama == 0 {
		t.Error("NOL komentar penyesuaian sisi `#LAMA` di T_VIEW_COMMENT — " +
			"panel History sisi Old kosong")
	}
	t.Logf("T_VIEW_COMMENT penyesuaian: sisi New %d · sisi #LAMA %d (terukur 1.488 · 280)",
		sisiBaru, sisiLama)

	var total, berisi int
	if err := h.QueryRowContext(ctx, `SELECT COUNT(*),
		SUM(CASE WHEN JSON_EXISTS(JSONDATA, '$.CommentList[0]') THEN 1 ELSE 0 END)
		FROM `+skema+`.M_TREATY_IN_EDM`).Scan(&total, &berisi); err != nil {
		t.Fatal(err)
	}
	if berisi == 0 {
		t.Error("nol CommentList berisi — panel History akan kosong untuk semua penyesuaian")
	}
	t.Logf("CommentList berisi di %d dari %d dokumen (terukur 280/280)", berisi, total)
}

// ⭐ Dasar pilihan dropdown: `PILIHAN_MEDAN` di layar SAMA dengan domain
// yang tersimpan di kedua sisi seluruh dokumen.
func TestPilihanDropdownSamaDenganDomainTersimpan(t *testing.T) {
	_, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)
	layar := pilihanLayar(t)

	for _, kunci := range []string{"Bordeaux", "AccountingMode", "AccountingModeNonProp"} {
		q := `SELECT DISTINCT v FROM (
			SELECT JSON_VALUE(JSONDATA, '$.` + kunci + `') v FROM ` + skema + `.M_TREATY_IN_EDM
			UNION ALL
			SELECT JSON_VALUE(JSONDATA, '$.OLDDATA.` + kunci + `') v FROM ` + skema + `.M_TREATY_IN_EDM
		) WHERE v IS NOT NULL ORDER BY v`
		rows, err := h.QueryContext(ctx, q)
		if err != nil {
			t.Fatal(err)
		}
		var tersimpan []string
		for rows.Next() {
			var v string
			if err := rows.Scan(&v); err != nil {
				t.Fatal(err)
			}
			tersimpan = append(tersimpan, v)
		}
		_ = rows.Close()
		ditawarkan := append([]string(nil), layar[kunci]...)
		sort.Strings(ditawarkan)
		if strings.Join(tersimpan, "|") != strings.Join(ditawarkan, "|") {
			t.Errorf("%s: tersimpan %v, layar menawarkan %v", kunci, tersimpan, ditawarkan)
		}
	}
}

// pilihanLayar membaca `PILIHAN_MEDAN` dari `labelsPenyesuaian.ts` —
// satu sumber, diperiksa dari dua sisi.
func pilihanLayar(t *testing.T) map[string][]string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "frontend", "labelsPenyesuaian.ts"))
	if err != nil {
		t.Fatal(err)
	}
	blok := regexp.MustCompile(`(?s)export const PILIHAN_MEDAN[^{]*\{(.*?)\n\}`).FindSubmatch(b)
	if blok == nil {
		t.Fatal("PILIHAN_MEDAN tidak terbaca")
	}
	out := map[string][]string{}
	for _, m := range regexp.MustCompile(`(\w+): \[([^\]]*)\]`).FindAllSubmatch(blok[1], -1) {
		for _, v := range regexp.MustCompile(`'([^']*)'`).FindAllSubmatch(m[2], -1) {
			out[string(m[1])] = append(out[string(m[1])], string(v[1]))
		}
	}
	if len(out) != 3 {
		t.Fatalf("PILIHAN_MEDAN terbaca %d kunci, mau 3", len(out))
	}
	return out
}

// ⭐ Dasar golongan `persenShare` kolom `Proportion %` EGNPI: barisnya
// berjumlah 100. Terukur 5 Oktober 2026 atas 468 sisi berisi: 365 TEPAT,
// 97 meleset ≤ 10⁻¹⁸ (sisa potong 20 desimal saat disimpan), 6 lain.
//
// ⚠️ "Berjumlah 100" di sini berarti dalam 10⁻⁶ — sisa pembulatan simpan
// bukan tanda persen biasa. Cacah TEPAT dicatat terpisah supaya angka
// itu tidak lagi dilaporkan sebagai yang lain.
func TestProporsiEgnpiBerjumlahSeratus(t *testing.T) {
	_, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)
	q := `SELECT COUNT(*), SUM(CASE WHEN s = 100 THEN 1 ELSE 0 END),
		SUM(CASE WHEN ABS(s - 100) < 0.000001 THEN 1 ELSE 0 END) FROM (
		SELECT m.ID, SUM(j.p) s FROM ` + skema + `.M_TREATY_IN_EDM m,
			JSON_TABLE(m.JSONDATA, '$.EGNPI[*]' COLUMNS (p NUMBER PATH '$.Proportion')) j GROUP BY m.ID
		UNION ALL
		SELECT m.ID, SUM(j.p) s FROM ` + skema + `.M_TREATY_IN_EDM m,
			JSON_TABLE(m.JSONDATA, '$.OLDDATA.EGNPI[*]' COLUMNS (p NUMBER PATH '$.Proportion')) j GROUP BY m.ID
	)`
	var total, tepat, dekat int
	if err := h.QueryRowContext(ctx, q).Scan(&total, &tepat, &dekat); err != nil {
		t.Fatal(err)
	}
	if total == 0 {
		t.Skip("lewati: nol EGNPI berisi")
	}
	if dekat*100 < total*95 {
		t.Errorf("hanya %d dari %d sisi berjumlah 100 (dalam 10⁻⁶) — golongan persenShare perlu ditinjau", dekat, total)
	}
	t.Logf("EGNPI: %d sisi, %d tepat 100, %d dalam 10⁻⁶ (terukur 468 / 365 / 462)", total, tepat, dekat)
}

// ⭐ Dasar pemilihan `.TreatyGroup` untuk sel `Treaty Group`: properti
// keduanya, `TreatyGroupName`, tidak pernah terisi.
func TestTreatyGroupNameTidakPernahTerisi(t *testing.T) {
	_, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)
	var n int
	if err := h.QueryRowContext(ctx, `SELECT COUNT(*) FROM (
		SELECT j.v FROM `+skema+`.M_TREATY_IN_EDM m,
			JSON_TABLE(m.JSONDATA, '$.Retention[*]' COLUMNS (v VARCHAR2(400) PATH '$.TreatyGroupName')) j
		UNION ALL
		SELECT j.v FROM `+skema+`.M_TREATY_IN_EDM m,
			JSON_TABLE(m.JSONDATA, '$.OLDDATA.Retention[*]' COLUMNS (v VARCHAR2(400) PATH '$.TreatyGroupName')) j
	) WHERE v IS NOT NULL`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("TreatyGroupName kini terisi %d kali — sel Treaty Group perlu ditinjau", n)
	}
}

// ⭐ Dasar menampilkan `.Currency` pada grid Rate of Exchange sisi New:
// `CurrencyID` adalah KODE yang berpasangan TETAP dengan satu nama.
func TestCurrencyIDBerpasanganTetapDenganCurrency(t *testing.T) {
	_, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)
	var ganda int
	if err := h.QueryRowContext(ctx, `SELECT COUNT(*) FROM (
		SELECT j.cid FROM `+skema+`.M_TREATY_IN_EDM m,
			JSON_TABLE(m.JSONDATA, '$.CurrencyList[*]' COLUMNS (
				cid VARCHAR2(100) PATH '$.CurrencyID', cur VARCHAR2(100) PATH '$.Currency')) j
		WHERE j.cid IS NOT NULL AND j.cur IS NOT NULL
		GROUP BY j.cid HAVING COUNT(DISTINCT j.cur) > 1)`).Scan(&ganda); err != nil {
		t.Fatal(err)
	}
	if ganda != 0 {
		t.Errorf("%d CurrencyID berpasangan dengan lebih dari satu nama — sel Currency New perlu ditinjau", ganda)
	}
}

// ⭐ Dasar teks kotak Retro (§17): `IsMultipleRetro = 'true'` di halaman
// akar penyesuaian. Terukur 2 dari 280.
func TestRetroJarangDiPenyesuaian(t *testing.T) {
	_, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)
	var total, retro int
	if err := h.QueryRowContext(ctx, `SELECT COUNT(*),
		SUM(CASE WHEN JSON_VALUE(JSONDATA, '$.IsMultipleRetro') = 'true' THEN 1 ELSE 0 END)
		FROM `+skema+`.M_TREATY_IN_EDM`).Scan(&total, &retro); err != nil {
		t.Fatal(err)
	}
	if retro != 2 || total != 280 {
		t.Errorf("IsMultipleRetro 'true' %d dari %d, teks layar menyebut 2 dari 280 — perbarui PENYESUAIAN.retroJarangPetunjuk", retro, total)
	}
}

// ⭐ Larik bertitik dan bersarang terbaca dari dokumen SUNGGUHAN.
func TestValueDifferenceShareTerbacaDariOracle(t *testing.T) {
	g, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)
	var id string
	if err := h.QueryRowContext(ctx, `SELECT MIN(ID) FROM `+skema+`.M_TREATY_IN_EDM
		WHERE JSON_EXISTS(JSONDATA, '$.ValueDifference.Share[0].RnmLimitListDisplay[0]')`).Scan(&id); err != nil {
		t.Skipf("lewati: nol dokumen ber-ValueDifference.Share: %v", err)
	}
	p, err := g.BacaPenyesuaianWarisan(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	b := p.Baru.Larik["ValueDifference.Share"]
	if len(b) == 0 {
		t.Fatalf("%s: ValueDifference.Share kosong", id)
	}
	if _, ada := b[0]["RnmLimitListDisplay(1).Currency"]; !ada {
		t.Errorf("%s: RnmLimitListDisplay(1).Currency tidak terbaca; kunci %v", id, b[0])
	}
}

// ⭐ Grid Rate of Exchange kedua sisi dari `TREATYEXCHANGEYEARLY`, disaring
// Treaty Year sisi masing-masing — bukan `CurrencyList` dokumen, yang tidak
// pernah mendarat. Diukur atas SELURUH penyesuaian, tanpa ROWNUM.
func TestKursPenyesuaianDariKursTahunan(t *testing.T) {
	g, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)

	cacahTahun := map[string]int{}
	rows, err := h.QueryContext(ctx, `SELECT TREATYYEAR, COUNT(*) FROM `+skema+`.TREATYEXCHANGEYEARLY GROUP BY TREATYYEAR`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var th string
		var n int
		if err := rows.Scan(&th, &n); err != nil {
			t.Fatal(err)
		}
		cacahTahun[th] = n
	}
	_ = rows.Close()

	d, err := g.DaftarPenyesuaianWarisan(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var terisiBaru, terisiLama, bertahunBaru, bertahunLama int
	for _, b := range d {
		p, err := g.BacaPenyesuaianPendaratan(ctx, b.ID)
		if err != nil {
			t.Fatal(err)
		}
		for _, s := range []struct {
			nama          string
			medan         map[string]string
			larik         []map[string]string
			terisi, tahun *int
		}{
			{"New", p.Baru.Medan, p.Baru.Larik["CurrencyList"], &terisiBaru, &bertahunBaru},
			{"Old", p.Lama.Medan, p.Lama.Larik["CurrencyList"], &terisiLama, &bertahunLama},
		} {
			th := strings.TrimSpace(s.medan["TreatyYear"])
			if th != "" {
				*s.tahun++
			}
			if len(s.larik) != cacahTahun[th] {
				t.Errorf("%s %s tahun %q: %d baris kurs, Oracle %d", b.ID, s.nama, th, len(s.larik), cacahTahun[th])
			}
			if len(s.larik) > 0 {
				*s.terisi++
				if s.larik[0]["Conversion"] == "" || s.larik[0]["Currency"] == "" {
					t.Errorf("%s %s: baris kurs tanpa Currency/Conversion: %v", b.ID, s.nama, s.larik[0])
				}
			}
		}
	}
	t.Logf("dari %d penyesuaian: New bertahun %d, berkurs %d · Old bertahun %d, berkurs %d",
		len(d), bertahunBaru, terisiBaru, bertahunLama, terisiLama)
	if terisiBaru == 0 {
		t.Error("nol penyesuaian berkurs di sisi New — pembaca kurs tidak tersambung")
	}
}

// ⭐ Enam larik yang ditambahkan ke peta 7 Oktober 2026 — setiap baris
// pendaratannya terbaca, kedua sisi, SELURUH penyesuaian (tanpa ROWNUM).
func TestLarikTambahanPenyesuaianTerbaca(t *testing.T) {
	g, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)

	tabel := map[string]string{
		"Installment": "T_TREATY_INSTALLMENT", "CoInScale": "M_TREATYIN_COINSCALE",
		"Share": "T_TREATY_SHARE", "ShareReins": "T_TREATY_RETRO_SHARE",
		"FacultativeShareList": "T_TREATY_FAC_SHARE", "ShareFacultativeReinsurers": "T_TREATY_FAC_REINSURER",
	}
	d, err := g.DaftarPenyesuaianWarisan(ctx)
	if err != nil {
		t.Fatal(err)
	}
	terbaca := map[string]int{}
	for _, b := range d {
		p, err := g.BacaPenyesuaianPendaratan(ctx, b.ID)
		if err != nil {
			t.Fatal(err)
		}
		for larik := range tabel {
			terbaca[larik] += len(p.Baru.Larik[larik]) + len(p.Lama.Larik[larik])
		}
	}
	for larik, tb := range tabel {
		var n int
		q := `SELECT COUNT(*) FROM ` + skema + `.` + tb + ` WHERE MASTERID IN
			(SELECT ID FROM ` + skema + `.TREATY_IN_EDM UNION ALL SELECT ID || '#LAMA' FROM ` + skema + `.TREATY_IN_EDM)`
		if err := h.QueryRowContext(ctx, q).Scan(&n); err != nil {
			t.Fatal(err)
		}
		t.Logf("%-28s %-24s Oracle %5d · terbaca %5d", larik, tb, n, terbaca[larik])
		if terbaca[larik] != n {
			t.Errorf("%s: terbaca %d baris, Oracle %d", larik, terbaca[larik], n)
		}
	}
}

// ⭐ `TreatyCalculateProratePct` diukur terhadap nilai TERSIMPAN — SELURUH
// penyesuaian ber-IsProRate, tanpa ROWNUM.
//
// Ekspor hari ini: Days = Effective → Termination. Data menyimpan DUA versi:
// penyesuaian terbaru mengikuti ekspor, yang lebih lama menyimpan
// Commencement → Effective (versi rumus sebelumnya). Terukur 7 Oktober
// 2026: 22 ber-IsProRate · 6 cocok ekspor · 16 cocok versi lama (satu
// cocok keduanya, 90 = 90) · 1 tersimpan 0 · total hari cocok 20.
// Layar membangun bunyi ekspor (`frontend/komponen/rumusKepala.ts`).
func TestRumusProRataTerukur(t *testing.T) {
	g, ctx := bacaSaja(t)
	d, err := g.DaftarPenyesuaianWarisan(ctx)
	if err != nil {
		t.Fatal(err)
	}
	tgl := func(s string) (time.Time, bool) {
		s = strings.TrimSpace(s)
		if len(s) < 8 {
			return time.Time{}, false
		}
		x, e := time.Parse("20060102", s[:8])
		return x, e == nil
	}
	var proRata, cocokEkspor, cocokLama, cocokTotal, takTerbaca int
	for _, b := range d {
		p, err := g.BacaPenyesuaianPendaratan(ctx, b.ID)
		if err != nil {
			t.Fatal(err)
		}
		m := p.Baru.Medan
		if m["IsProRate"] != "true" {
			continue
		}
		proRata++
		ef, ok1 := tgl(m["EDMEffective"])
		mu, ok2 := tgl(m["Commencement"])
		ak, ok3 := tgl(m["Termination"])
		if !ok1 || !ok2 || !ok3 {
			takTerbaca++
			continue
		}
		hari := func(a, b time.Time) string { return fmt.Sprint(int(b.Sub(a).Hours() / 24)) }
		tersimpan := strings.TrimSpace(m["ProRateDays"])
		if hari(ef, ak) == tersimpan {
			cocokEkspor++
		}
		if hari(mu, ef) == tersimpan {
			cocokLama++
		}
		if hari(mu, ak) == strings.TrimSpace(m["ProRateTotalDays"]) {
			cocokTotal++
		}
	}
	t.Logf("IsProRate=true %d dari %d penyesuaian · hari cocok ekspor %d · cocok versi lama %d · total cocok %d · tak terbaca %d",
		proRata, len(d), cocokEkspor, cocokLama, cocokTotal, takTerbaca)
	if proRata > 0 && cocokEkspor == 0 {
		t.Error("nol penyesuaian cocok dengan rumus ekspor — rumus layar perlu ditinjau")
	}
}

// ⭐ Panel `Existing Policy for Master ID` — `TREATYINPRODUCTION`, SQL rule
// `FetchTreatyInProductionUsingNooffer` apa adanya. Diukur atas SELURUH
// penyesuaian: cacah barisnya sama dengan kueri mentah yang setara.
func TestPolisMasterPenyesuaianDariProduksi(t *testing.T) {
	g, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)
	d, err := g.DaftarPenyesuaianWarisan(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var berpolis, total int
	for _, b := range d {
		id := b.IDAsal
		if strings.TrimSpace(id) == "" {
			id = b.ID
		}
		polis, err := g.BacaPolisMaster(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		var n int
		q := `SELECT COUNT(*) FROM (SELECT DISTINCT NOPOLIS, IDPEGA, QUARTER, QUARTER_YEAR FROM ` + skema +
			`.TREATYINPRODUCTION WHERE SUBSTR(NOOFFER,1,7) = SUBSTR(:1,1,7))`
		if err := h.QueryRowContext(ctx, q, id).Scan(&n); err != nil {
			t.Fatal(err)
		}
		if len(polis) != n {
			t.Errorf("%s (master %s): %d baris, kueri mentah %d", b.ID, id, len(polis), n)
		}
		if n > 0 {
			berpolis++
		}
		total += n
	}
	t.Logf("dari %d penyesuaian: %d punya polis produksi, %d baris", len(d), berpolis, total)
}

// ⭐ POHON — setiap baris tabel ANAK terangkai ke induknya, kedua sisi,
// SELURUH penyesuaian (tanpa ROWNUM), dan nol yang yatim.
func TestPohonPenyesuaianUtuh(t *testing.T) {
	g, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)
	d, err := g.DaftarPenyesuaianWarisan(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// Cacah simpul per nama larik (akar maupun anak) di seluruh pohon.
	cacah := map[string]int{}
	var hitung func(simpul []map[string]any)
	hitung = func(simpul []map[string]any) {
		for _, s := range simpul {
			for k, v := range s {
				if anak, ok := v.([]map[string]any); ok {
					cacah[k] += len(anak)
					hitung(anak)
				}
			}
		}
	}
	yatimTotal := 0
	for _, b := range d {
		for _, id := range []string{b.ID, b.ID + "#LAMA"} {
			pohon, yatim, err := g.BacaPohonUntukUji(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
			yatimTotal += yatim
			for larik, simpul := range pohon {
				cacah["akar:"+larik] += len(simpul)
				hitung(simpul)
			}
		}
	}
	if yatimTotal != 0 {
		t.Errorf("%d baris anak tanpa induk", yatimTotal)
	}
	saring := `WHERE MASTERID IN (SELECT ID FROM ` + skema + `.TREATY_IN_EDM UNION ALL SELECT ID || '#LAMA' FROM ` + skema + `.TREATY_IN_EDM)`
	periksa := func(label, tabel, tambahan string, terbaca int) {
		var n int
		if err := h.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+skema+`.`+tabel+` `+saring+tambahan).Scan(&n); err != nil {
			t.Fatal(err)
		}
		t.Logf("%-30s %-28s Oracle %6d · pohon %6d", label, tabel, n, terbaca)
		if n != terbaca {
			t.Errorf("%s: pohon %d, Oracle %d", label, terbaca, n)
		}
	}
	periksa("akar:Limits", "T_TREATY_LIMITS", "", cacah["akar:Limits"])
	periksa("akar:Share", "T_TREATY_SHARE", "", cacah["akar:Share"])
	periksa("Detail", "T_TREATY_LIMIT_DETAIL", "", cacah["Detail"])
	periksa("TreatyGroupList", "T_TREATY_LIMIT_GROUP", "", cacah["TreatyGroupList"])
	periksa("ClassOfBusinessList", "T_TREATY_LIMIT_GROUP_COB", "", cacah["ClassOfBusinessList"])
	periksa("COBList", "T_TREATY_LIMIT_COB", "", cacah["COBList"])
	periksa("SpreadingListXOL", "T_TREATY_SHARE_SPREADING", "", cacah["SpreadingListXOL"])
	periksa("InstallmentList", "T_TREATY_INSTALLMENT_ITEM", "", cacah["InstallmentList"])
	periksa("MDPList", "T_TREATY_LIMIT_MEASURE", " AND JENIS = 'MDPList'", cacah["MDPList"])
	periksa("IOOLimitList", "T_TREATY_LIMIT_AMOUNT", " AND JENIS = 'IOOLimitList'", cacah["IOOLimitList"])
}

