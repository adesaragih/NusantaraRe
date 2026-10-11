package repository

// Untuk apa berkas ini: POLIS DARI TABEL DATAR NB FAC IN - perintah work owner 10-10-2026 ("pertama cek json_polis,
// jika json nya kosong ambil dari flat nya"); susunan dokumennya di `models/polis_datar.go`. BACA SAJA: tabel datar
// milik nbfacin (`modul/nbfacin/docs/STRUKTUR-TABEL-NB-FACIN.md`; peta rancangan `models/polis_datar_peta.go`) dan
// master warisan BUSINESS / T_M_ACCOUNT / MARKETINGOFFICER / CURRENCY / FACINPRODUCTION / JSON_POLIS.
//
// Peta itu RANCANGAN loader (kolom yang belum dimigrasi termasuk): setiap kueri hanya memilih tabel / kolom yang ADA di
// katalog (`ALL_TAB_COLUMNS` pemilik skema, dibaca sekali per proses), jadi skema yang berbeda antarlingkungan tidak
// menjadi ORA-00904 / ORA-00942.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"nusantarare/modul/claimfacin/backend/models"
)

// batasIndukDatar - paling banyak bind PARENT_ID per kueri anak (sisa dipecah).
const batasIndukDatar = 500

// katalogDatar - kolom tabel datar polis yang ada di basis data (tabel -> kolom -> DATA_TYPE) dan indeks induknya
// (tabel -> nama indeks berawalan PARENT_ID / PARENT_TABLE,PARENT_ID), dibaca sekali.
type katalogDatar struct {
	mu     sync.Mutex
	kolom  map[string]map[string]string
	indeks map[string]string
}

// tembolokDatar - dokumen polis datar hasil susun per (nopolis, prodke, kasus) selama `umurTembolokDatar`: satu aksi
// memuat kasus dua kali dan layar memanggil beberapa rute sekaligus, sedangkan polis hasil produksi (satu PRODKE) tidak
// berubah (pengukuran DEV 10-10-2026: satu susun ~0,7 detik sebelum hint indeks). Keluhan work owner 11-10-2026 "SAAT
// BUKA CLM NYA LAMBAT!": DEV kini ~1 juta baris tabel datar dan `T_PROPERTYITEMLIST` hanya berindeks PK - setiap susun
// menyapu seluruh tabel (0,1 s hangat, beberapa detik saat cache Oracle dingin; susun pertama terukur 5,9 s). Umur
// diperpanjang dari 5 menit ke 8 jam supaya penyapuan itu terjadi sekali per polis per proses, bukan tiap 5 menit.
type tembolokDatar struct {
	mu  sync.Mutex
	isi map[string]entriDatar
}

type entriDatar struct {
	dok  []byte
	saat time.Time
}

// umurTembolokDatar / batasTembolokDatar - umur satu entri dan jumlah entri paling banyak.
const (
	umurTembolokDatar  = 8 * time.Hour
	batasTembolokDatar = 256
)

func (c *tembolokDatar) ambil(k string, kini time.Time) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ada := c.isi[k]
	if !ada || kini.Sub(e.saat) > umurTembolokDatar {
		return nil, false
	}
	return e.dok, true
}

func (c *tembolokDatar) simpan(k string, dok []byte, kini time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isi == nil {
		c.isi = map[string]entriDatar{}
	}
	for kk, e := range c.isi { // entri kedaluwarsa dibuang
		if kini.Sub(e.saat) > umurTembolokDatar {
			delete(c.isi, kk)
		}
	}
	for len(c.isi) >= batasTembolokDatar { // penuh: entri tertua dibuang
		tua, saat := "", kini
		for kk, e := range c.isi {
			if !e.saat.After(saat) {
				tua, saat = kk, e.saat
			}
		}
		delete(c.isi, tua)
	}
	c.isi[k] = entriDatar{dok: dok, saat: kini}
}

// indeksIndukDatar - kunci `kolom` katalog untuk nama indeks induk tabel (bukan nama kolom Oracle).
const indeksIndukDatar = "#INDEKS"

// sqlIndeksDatar - dua kolom terdepan indeks tabel datar (SYS.ALL_IND_COLUMNS, pemilik `:1`, `n` tabel bind :2..).
// Pengukuran DEV 10-10-2026: tabel datar bersegmen ratusan MB berstatistik basi - tanpa hint indeks Oracle menyapu
// seluruh segmen (100-190 ms per kueri), dengan hint indeks induk 1 ms.
func sqlIndeksDatar(n int) string {
	return fmt.Sprintf(`SELECT TABLE_NAME, INDEX_NAME, COLUMN_NAME, COLUMN_POSITION FROM SYS.ALL_IND_COLUMNS
		WHERE TABLE_OWNER = NVL(:1, USER) AND COLUMN_POSITION <= 2 AND TABLE_NAME IN (%s)
		ORDER BY TABLE_NAME, INDEX_NAME, COLUMN_POSITION`, daftarBind(2, n))
}

// pilihIndeksInduk - per tabel: indeks berkolom pertama PARENT_ID, atau PARENT_TABLE lalu PARENT_ID.
func pilihIndeksInduk(rows [][]string) map[string]string {
	kol := map[[2]string][]string{}
	var urut [][2]string
	for _, r := range rows {
		k := [2]string{r[0], r[1]}
		if kol[k] == nil {
			urut = append(urut, k)
		}
		kol[k] = append(kol[k], r[2])
	}
	out := map[string]string{}
	for _, k := range urut {
		c := kol[k]
		cocok := c[0] == "PARENT_ID" || (len(c) > 1 && c[0] == "PARENT_TABLE" && c[1] == "PARENT_ID")
		if _, ada := out[k[0]]; cocok && !ada {
			out[k[0]] = k[1]
		}
	}
	return out
}

// sqlKatalogDatar - SYS.ALL_TAB_COLUMNS pemilik `:1` (kosong = skema sesi, USER) atas `n` tabel (bind :2..:n+1).
func sqlKatalogDatar(n int) string {
	return fmt.Sprintf(`SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = NVL(:1, USER)
		AND TABLE_NAME IN (%s)`, daftarBind(2, n))
}

func daftarBind(awal, n int) string {
	b := make([]string, n)
	for i := range b {
		b[i] = fmt.Sprintf(":%d", awal+i)
	}
	return strings.Join(b, ", ")
}

// kolomAda - katalog tabel datar polis (dibaca sekali per proses; galat tidak disimpan).
func (a *Acuan) kolomAda(ctx context.Context) (map[string]map[string]string, error) {
	a.katalog.mu.Lock()
	defer a.katalog.mu.Unlock()
	if a.katalog.kolom != nil {
		return a.katalog.kolom, nil
	}
	akar, err := a.q(models.TabelAkarPolisDatar)
	if err != nil {
		return nil, err
	}
	var pemilik any // nama tanpa skema: skema sesi
	if p, _, ada := strings.Cut(akar, "."); ada {
		pemilik = strings.ToUpper(strings.Trim(p, `"`))
	}
	nama := models.NamaTabelPolisDatar()
	args := make([]any, 0, len(nama)+1)
	args = append(args, pemilik)
	for _, t := range nama {
		args = append(args, t)
	}
	rows, err := a.banyak(ctx, sqlKatalogDatar(len(nama)), 3, args...)
	if err != nil {
		return nil, err
	}
	ind, err := a.banyak(ctx, sqlIndeksDatar(len(nama)), 4, args...)
	if err != nil {
		return nil, err
	}
	k := map[string]map[string]string{}
	for _, r := range rows {
		if k[r[0]] == nil {
			k[r[0]] = map[string]string{}
		}
		k[r[0]][r[1]] = r[2]
	}
	for t, nm := range pilihIndeksInduk(ind) {
		if k[t] != nil {
			k[t][indeksIndukDatar] = nm
		}
	}
	a.katalog.kolom = k
	return k, nil
}

// kolomDatar - kolom rancangan `spec` yang ada di katalog `ada`.
func kolomDatar(spec models.TabelPolisDatar, ada map[string]string) []models.KolomPolisDatar {
	var out []models.KolomPolisDatar
	for _, k := range spec.Kolom {
		if _, ok := ada[k.Nama]; ok {
			out = append(out, k)
		}
	}
	return out
}

// pilihKolomDatar - daftar SELECT kolom (NUMBER lewat TO_CHAR 'TM9').
func pilihKolomDatar(ks []models.KolomPolisDatar, ada map[string]string) string {
	b := make([]string, len(ks))
	for i, k := range ks {
		b[i] = k.Nama
		if ada[k.Nama] == "NUMBER" {
			b[i] = dec(k.Nama)
		}
	}
	return strings.Join(b, ", ")
}

// sqlAkarPolisDatar - baris T_GENERAL_POLIS kasus.
func sqlAkarPolisDatar(t string, ks []models.KolomPolisDatar, ada map[string]string) string {
	pilih := "ID"
	if len(ks) > 0 {
		pilih += ", " + pilihKolomDatar(ks, ada)
	}
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, pilih, t)
}

// sqlAnakPolisDatar - baris tabel anak ber-PARENT_ID di `n` bind; tabel ber-SRC_PATH (katalog) juga PARENT_TABLE =
// bind n+1; urut SEQ_NO bila ada; hint indeks induk bila katalog menemukannya (`sqlIndeksDatar`).
func sqlAnakPolisDatar(t string, ks []models.KolomPolisDatar, ada map[string]string, n int) string {
	pilih, saring, urut := "ID, PARENT_ID", "", "ID"
	if nm := ada[indeksIndukDatar]; nm != "" {
		pilih = "/*+ INDEX(x " + nm + ") */ " + pilih
	}
	if srcDatar(ada) {
		pilih += ", SRC_PATH"
		saring = fmt.Sprintf(" AND PARENT_TABLE = :%d", n+1)
	}
	if _, ok := ada["SEQ_NO"]; ok {
		urut = "SEQ_NO, ID"
	}
	if len(ks) > 0 {
		pilih += ", " + pilihKolomDatar(ks, ada)
	}
	return fmt.Sprintf(`SELECT %s FROM %s x WHERE PARENT_ID IN (%s)%s ORDER BY %s`, pilih, t, daftarBind(1, n), saring, urut)
}

func srcDatar(ada map[string]string) bool {
	_, src := ada["SRC_PATH"]
	_, induk := ada["PARENT_TABLE"]
	return src && induk
}

// sqlKasusPolisDatar - ID kasus sistem baru = ruas terakhir JSON_POLIS.IDPEGA ("ASM-FW-GISFW-WORK NB-n" -> "NB-n";
// DEV 10-10-2026: `SUBSTR(IDPEGA, 19)` = `T_WORK_POLIS.ID`), baris terbaru ber-DATA_JSON lebih dulu (baris yang sama
// dengan `DokumenPolis`).
func sqlKasusPolisDatar(jp string) string {
	return fmt.Sprintf(`SELECT REGEXP_SUBSTR(IDPEGA, '[^ ]+$') FROM %s WHERE NOPOLIS = :1 AND PRODKE = :2
		ORDER BY CASE WHEN DATA_JSON IS NULL THEN 1 ELSE 0 END, TGL_INPUT DESC NULLS LAST FETCH FIRST 1 ROWS ONLY`, jp)
}

// sqlKasusProduksiDatar - cadangan tanpa baris JSON_POLIS: IDPEGA baris FACINPRODUCTION polis itu yang terbaru.
func sqlKasusProduksiDatar(fp string) string {
	return fmt.Sprintf(`SELECT REGEXP_SUBSTR(IDPEGA, '[^ ]+$') FROM %s WHERE NOPOLIS = :1 AND IDPEGA IS NOT NULL
		ORDER BY SPREAD_DATE DESC NULLS LAST FETCH FIRST 1 ROWS ONLY`, fp)
}

// sqlOpportunityDatar - Business Status, Insured ID / Name (T_M_ACCOUNT tanpa PK: MAX, pola nbfacin kasus.go), Class
// Of Business, Group Business kasus.
func sqlOpportunityDatar(opp, akun string) string {
	return fmt.Sprintf(`SELECT o.BUSINESS_STATUS, o.INSURED_ID, o.CLASS_OF_BUSINESS, o.GROUP_BUSINESS_ID,
		       (SELECT MAX(a.INSUREDNAME) FROM %s a WHERE a.ID = o.ACCOUNT_ID)
		  FROM %s o WHERE o.ID = :1`, akun, opp)
}

// sqlBisnisProduksiDatar - master BUSINESS menurut BUSINESSCODE baris FACINPRODUCTION polis ber-kasus `:2`.
func sqlBisnisProduksiDatar(bis, fp string) string {
	return fmt.Sprintf(`SELECT b.ID, b.OLDID, b.GROUPPANEL, b.NOTE FROM %s b
		WHERE b.ID = (SELECT MAX(f.BUSINESSCODE) FROM %s f WHERE f.NOPOLIS = :1 AND REGEXP_SUBSTR(f.IDPEGA, '[^ ]+$') = :2)`,
		bis, fp)
}

// sqlBisnisNamaDatar - cadangan: master BUSINESS ber-NOTE = Class Of Business opportunity, di BUSINESSGROUPID
// Group Business-nya bila terisi (pilihan Class Of Business NB: RD filter `.BusinessGroupID = Param.Group`).
func sqlBisnisNamaDatar(bis string) string {
	return fmt.Sprintf(`SELECT ID, OLDID, GROUPPANEL, NOTE FROM %s WHERE NOTE = :1 AND (:2 IS NULL OR BUSINESSGROUPID = :3)
		ORDER BY ID FETCH FIRST 1 ROWS ONLY`, bis)
}

// sqlMarketingDatar - MARKETINGOFFICER menurut MOID (pasangan medan Quotation: `estimasi_aksi.go` MarketingOfficer).
func sqlMarketingDatar(mo string) string {
	return fmt.Sprintf(`SELECT CLIENTID, CLIENTNAME, TEAMGROUP, BRANCHDETAILID, BRANCHDETAILNAME FROM %s WHERE ID = :1
		FETCH FIRST 1 ROWS ONLY`, mo)
}

// sqlIDMataUangDatar - CURRENCY.ID menurut kode (T_PROPERTYITEMLIST.CURRENCY = CURRENCY.CURRENCY, nbfacin item.go).
func sqlIDMataUangDatar(cur string) string {
	return fmt.Sprintf(`SELECT MAX(ID) FROM %s WHERE CURRENCY = :1`, cur)
}

// DokumenPolisDatar - lihat `models.Acuan`.
func (a *Acuan) DokumenPolisDatar(ctx context.Context, nopolis, prodke, kasus string) ([]byte, bool, error) {
	kunci := nopolis + "|" + prodke + "|" + kasus
	if dok, ada := a.tembolok.ambil(kunci, time.Now()); ada {
		return dok, true, nil
	}
	if kasus == "" {
		var err error
		if kasus, err = a.kasusPolisDatar(ctx, nopolis, prodke); err != nil || kasus == "" {
			return nil, false, err
		}
	}
	kat, err := a.kolomAda(ctx)
	if err != nil {
		return nil, false, err
	}
	akar, ada, err := a.akarPolisDatar(ctx, kat, kasus)
	if err != nil || !ada {
		return nil, false, err
	}
	doc, err := models.SusunPolisDatar(akar, func(s models.SisiPolisDatar, induk []string) ([]models.BarisPolisDatar,
		error) {
		return a.anakPolisDatar(ctx, kat, s, induk)
	})
	if err != nil {
		return nil, false, err
	}
	t, err := a.tambahanPolisDatar(ctx, nopolis, kasus, models.MOIDDatar(doc))
	if err != nil {
		return nil, false, err
	}
	ids := map[string]string{}
	if kode := models.KodeMataUangDatar(doc); len(kode) > 0 {
		cur, err := a.q("CURRENCY")
		if err != nil {
			return nil, false, err
		}
		for _, k := range kode {
			if ids[k], _, err = a.satu(ctx, sqlIDMataUangDatar(cur), k); err != nil {
				return nil, false, err
			}
		}
	}
	models.LengkapiPolisDatar(doc, nopolis, t, ids)
	b, err := json.Marshal(doc)
	if err != nil {
		return nil, false, fmt.Errorf("repository: dokumen polis datar: %w", err)
	}
	a.tembolok.simpan(kunci, b, time.Now())
	return b, true, nil
}

func (a *Acuan) kasusPolisDatar(ctx context.Context, nopolis, prodke string) (string, error) {
	jp, err := a.q("JSON_POLIS")
	if err != nil {
		return "", err
	}
	if v, ada, err := a.satu(ctx, sqlKasusPolisDatar(jp), nopolis, prodke); err != nil || (ada && v != "") {
		return v, err
	}
	fp, err := a.q("FACINPRODUCTION")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, sqlKasusProduksiDatar(fp), nopolis)
	return v, err
}

func (a *Acuan) akarPolisDatar(ctx context.Context, kat map[string]map[string]string, kasus string) (
	models.BarisPolisDatar, bool, error) {
	ada, ok := kat[models.TabelAkarPolisDatar]
	if !ok {
		return models.BarisPolisDatar{}, false, nil
	}
	spec, _ := models.TabelPolisDatarDari(models.TabelAkarPolisDatar)
	ks := kolomDatar(spec, ada)
	t, err := a.q(models.TabelAkarPolisDatar)
	if err != nil {
		return models.BarisPolisDatar{}, false, err
	}
	rows, err := a.banyak(ctx, sqlAkarPolisDatar(t, ks, ada), 1+len(ks), kasus)
	if err != nil || len(rows) == 0 {
		return models.BarisPolisDatar{}, false, err
	}
	return models.BarisPolisDatar{ID: rows[0][0], Nilai: nilaiDatar(ks, rows[0][1:])}, true, nil
}

// anakPolisDatar - baris tabel `s.Tabel` ber-induk di `induk` (dipecah per batasIndukDatar bind); tabel yang tidak ada
// di katalog = tanpa baris; tabel khusus tampilan polis NB tanpa indeks induk tidak dibaca (`models.TabelTampilanNB`).
func (a *Acuan) anakPolisDatar(ctx context.Context, kat map[string]map[string]string, s models.SisiPolisDatar,
	induk []string) ([]models.BarisPolisDatar, error) {
	ada, ok := kat[s.Tabel]
	if _, punyaInduk := ada["PARENT_ID"]; !ok || !punyaInduk {
		return nil, nil
	}
	if models.TabelTampilanNB(s.Tabel) && ada[indeksIndukDatar] == "" {
		return nil, nil
	}
	spec, _ := models.TabelPolisDatarDari(s.Tabel)
	ks := kolomDatar(spec, ada)
	t, err := a.q(s.Tabel)
	if err != nil {
		return nil, err
	}
	src := srcDatar(ada)
	lebih := 2
	if src {
		lebih = 3
	}
	var out []models.BarisPolisDatar
	for len(induk) > 0 {
		n := min(len(induk), batasIndukDatar)
		args := make([]any, 0, n+1)
		for _, id := range induk[:n] {
			args = append(args, id)
		}
		if src {
			args = append(args, s.Induk)
		}
		rows, err := a.banyak(ctx, sqlAnakPolisDatar(t, ks, ada, n), lebih+len(ks), args...)
		if err != nil {
			return nil, err
		}
		for _, r := range rows {
			b := models.BarisPolisDatar{ID: r[0], Induk: r[1], Nilai: nilaiDatar(ks, r[lebih:])}
			if src {
				b.Src = r[2]
			}
			out = append(out, b)
		}
		induk = induk[n:]
	}
	return out, nil
}

// nilaiDatar - nilai baris per medan; sentinel mata uang `UNKNOWN` (DEFAULT CURRENCY / CURRENCY_CODE) dibaca kosong.
func nilaiDatar(ks []models.KolomPolisDatar, r []string) map[string]string {
	m := make(map[string]string, len(ks))
	for i, k := range ks {
		v := r[i]
		if v == "" || ((k.Nama == "CURRENCY" || k.Nama == "CURRENCY_CODE") && v == models.KodeMataUangTakDikenal) {
			continue
		}
		m[k.Medan] = v
	}
	return m
}

func (a *Acuan) tambahanPolisDatar(ctx context.Context, nopolis, kasus, moid string) (models.TambahanPolisDatar, error) {
	var t models.TambahanPolisDatar
	opp, err := a.q("T_NB_OPPORTUNITY")
	if err != nil {
		return t, err
	}
	akun, err := a.q("T_M_ACCOUNT")
	if err != nil {
		return t, err
	}
	rows, err := a.banyak(ctx, sqlOpportunityDatar(opp, akun), 5, kasus)
	if err != nil {
		return t, err
	}
	kelas, grup := "", ""
	if len(rows) > 0 {
		t.StatusBusiness, t.InsuredID, kelas, grup, t.InsuredName = rows[0][0], rows[0][1], rows[0][2], rows[0][3], rows[0][4]
	}
	bis, err := a.q("BUSINESS")
	if err != nil {
		return t, err
	}
	fp, err := a.q("FACINPRODUCTION")
	if err != nil {
		return t, err
	}
	b, err := a.banyak(ctx, sqlBisnisProduksiDatar(bis, fp), 4, nopolis, kasus)
	if err != nil {
		return t, err
	}
	if len(b) == 0 && kelas != "" {
		var g any
		if grup != "" {
			g = grup
		}
		if b, err = a.banyak(ctx, sqlBisnisNamaDatar(bis), 4, kelas, g, g); err != nil {
			return t, err
		}
	}
	if len(b) > 0 {
		t.BusinessCode, t.BusinessOldId, t.GroupPanel, t.BusinessName = b[0][0], b[0][1], b[0][2], b[0][3]
	}
	if moid == "" {
		return t, nil
	}
	mo, err := a.q("MARKETINGOFFICER")
	if err != nil {
		return t, err
	}
	m, err := a.banyak(ctx, sqlMarketingDatar(mo), 5, moid)
	if err != nil || len(m) == 0 {
		return t, err
	}
	t.MarketingCode, t.MarketingName, t.TeamGroup, t.BranchCode, t.BranchName = m[0][0], m[0][1], m[0][2], m[0][3], m[0][4]
	return t, nil
}
