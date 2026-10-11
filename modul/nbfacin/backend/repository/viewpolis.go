package repository

// View Policy tab Facultative (List Offer) / Treaty / Life - tiket 71 perluasan (permintaan work owner 11-10-2026 "menu view policy ini
// buatin untuk life dan treaty juga"). Sumber `[terverifikasi]` rule kiriman work owner di `D:\migrasi\RNM\DDL\` (BUKAN korpus):
// `GetListViewPolicy_act.xml` blok 9 (OFFERFAC) / 10 (TREATYPOL) / 11 (POLISLIFE) / 12 (OFFERLIFE) dan Rule-RDB-SQL
// ASM-FW-GISFW-INT-POLICYJSON!ASM BrowseDataViewOfferF / BrowsePolisTreaty / BrowsePolisLife / BrowseOfferLife (`pyBrowseSQL`).
// Pega menyusun WHERE dengan {ASIS:ParamViewIn.CARI5} (gabungan teks); di sini SEMUA nilai lewat BIND (anti injeksi) - bentuk klausa
// sama. IDPEGA dicocokkan dua bentuk (polos / berawalan, idpegaproduksi.go); case id dari IDPEGA = SUBSTR(ekspresi dua bentuk, 1, 10)
// (Pega `substr(IDPEGA, 19, 10)`). Tabel POOLDATA milik modul lain - BACA SAJA, tanpa impor modul lain. Urutan halaman: Pega tanpa
// ORDER BY -> urutan stabil `[keputusan agent]` (lihat sqlHalamanView).
// ⚠️ DDL LIFEINPRODUCTION / JSON_OFFER_LIFE TIDAK ada di D:\migrasi\RNM\DDL (endorsementlife: "nol DDL") - kolom yang dibaca dicatat di
// kepala tiap SQL, tipe / lebar `[belum terverifikasi]` (dibaca sebagai teks).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// Tabel POOLDATA yang dibaca View Policy (milik modul lain; baca saja).
const (
	TabelTreatyInProduction = "TREATYINPRODUCTION"
	TabelLifeInProduction   = "LIFEINPRODUCTION"
	TabelJSONOfferLife      = "JSON_OFFER_LIFE"
	// TabelGeneralPolisTreaty - kepala berkas NB / EDM Treaty In sistem baru (migrasi nbtreatyin 320; PRODKE NUMBER(10) NOT NULL);
	// baca saja untuk tombol View tab Treaty.
	TabelGeneralPolisTreaty = "T_GENERAL_POLIS_TREATY"
)

// Jenis klausa View Policy (dipilih services menurut View Option / deteksi isi teks; nilai = bind akhir, pola LIKE disusun services).
const (
	// OFFERFAC (blok 9, View Option keputusan work owner "samain kaya polis facultative").
	OfferFacIDPega  = "offer-idpega"
	OfferFacNoOffer = "offer-nooffer"
	OfferFacReff    = "offer-reff"
	OfferFacInsured = "offer-insured"
	OfferFacSOB     = "offer-sob"
	// TREATYPOL (blok 10 / FilterForViewPolicy).
	TreatyIDPega   = "treaty-idpega"
	TreatyEDMT     = "treaty-edmt"
	TreatyNoPolis  = "treaty-nopolis"
	TreatyNoEndors = "treaty-noendors"
	TreatyMaster   = "treaty-master"
	TreatyInsured  = "treaty-insured"
	TreatySOB      = "treaty-sob"
	// POLISLIFE (blok 11, deteksi isi teks).
	LifeNoPolis    = "life-nopolis"
	LifeDigitAkhir = "life-digit"
	LifeIDPega     = "life-idpega"
	LifeNoEndors   = "life-noendors"
	LifeEDML       = "life-edml"
	LifeSOB        = "life-sob"
	// OFFERLIFE (blok 12, deteksi isi teks).
	OfferLifeIDPega = "offerlife-idpega"
	OfferLifeOFR    = "offerlife-ofr"
	OfferLifeCeding = "offerlife-ceding"
)

// klausaView - WHERE per jenis (bind :1..:n berurutan) dan jenis IDPEGA (nilai diperluas dua bentuk). `t` = tabel berkualifikasi
// tab itu. Teks klausa = CARI5 Pega dengan {ASIS} diganti bind.
type klausaView struct {
	teks   string
	idPega bool // nilai = case id -> (polos, berawalan)
	ulang  int  // nilai yang sama diikat n kali (OR LIKE berulang); 0 = 1
}

func klausaViewDari(jenis, t string) (klausaView, bool) {
	k, ada := map[string]klausaView{
		// 9.x (View Option): IDPEGA dua bentuk; NOOFFER =; RISLIPCEDING / INSUREDNAME / SOB LIKE.
		OfferFacIDPega:  {teks: "PC0.IDPEGA IN (:1, :2)", idPega: true},
		OfferFacNoOffer: {teks: "PC0.NOOFFER = :1"},
		OfferFacReff:    {teks: "PC0.RISLIPCEDING LIKE :1"},
		OfferFacInsured: {teks: "PC0.INSUREDNAME LIKE :1"},
		OfferFacSOB:     {teks: "PC0.SOB LIKE :1"},
		// 10.1-10.7.
		TreatyNoPolis:  {teks: "PC0.NOPOLIS = :1 and PC0.NOENDORS is null"},
		TreatyNoEndors: {teks: "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM " + t + " WHERE NOENDORS = :1)  and PC0.NOENDORS is null"},
		TreatyIDPega:   {teks: " PC0.IDPEGA IN (:1, :2) and PC0.NOENDORS is null", idPega: true},
		TreatyEDMT:     {teks: "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM " + t + " WHERE IDPEGA IN (:1, :2))  and PC0.NOENDORS is null", idPega: true},
		TreatyMaster:   {teks: "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM " + t + " WHERE NOOFFER = :1)  and PC0.NOENDORS is null"},
		TreatyInsured:  {teks: "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM " + t + " WHERE INSUREDNAME LIKE :1)  and PC0.NOENDORS is null"},
		TreatySOB:      {teks: "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM " + t + " WHERE SOB LIKE :1)  and PC0.NOENDORS is null"},
		// 11.1-11.6 (BrowsePolisLife `WHERE {ASIS}` tanpa kurung). 11.6 DIBERI KURUNG `(A OR B OR C) AND NOENDORS IS NULL` - hanya polis
		// induk (keputusan work owner 11-10-2026; Pega `A OR B OR C AND …` mengikat AND ke suku terakhir saja).
		LifeNoPolis:    {teks: "NOPOLIS = :1 and NOENDORS is null"},
		LifeDigitAkhir: {teks: "NOPOLIS like :1 and NOENDORS is null"},
		LifeIDPega:     {teks: " IDPEGA IN (:1, :2) and NOENDORS is null", idPega: true},
		LifeNoEndors:   {teks: "NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM " + t + " WHERE NOENDORS  = :1) and NOENDORS is null"},
		LifeEDML:       {teks: "NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM " + t + " WHERE IDPEGA IN (:1, :2)) and NOENDORS is null", idPega: true},
		LifeSOB:        {teks: "(SOBNAME LIKE :1 OR CEDINGCONAME LIKE :2 OR POLICYHOLDERNAME LIKE :3)  AND NOENDORS IS NULL", ulang: 3},
		// 12.1-12.3.
		OfferLifeIDPega: {teks: " A.IDPEGA IN (:1, :2)", idPega: true},
		OfferLifeOFR:    {teks: " A.ID = SUBSTR(:1, 5)"},
		OfferLifeCeding: {teks: " A.CEDINGCONAME like :1 OR  A.POLICYHOLDERNAME like :2", ulang: 2},
	}[jenis]
	return k, ada
}

// bindView - nilai bind klausa jenis itu.
func bindView(k klausaView, nilai string) []any {
	if k.idPega {
		p, a := DuaBentukIDPega(nilai)
		return []any{p, a}
	}
	n := max(k.ulang, 1)
	arg := make([]any, n)
	for i := range arg {
		arg[i] = nilai
	}
	return arg
}

// caseIDView - case id 10 karakter dari kolom IDPEGA dua bentuk (Pega substr(IDPEGA, 19, 10)).
func caseIDView(kolom string) string { return "SUBSTR(" + sqlCaseIDDariIDPega(kolom) + ", 1, 10)" }

// sqlDasarOfferFac - BrowseDataViewOfferF: DISTINCT IDPEGA / NOOFFER / BEGINDATE / ENDDATE / SOB / INSUREDNAME (QQ teks kosong, tidak dibaca).
// Pega memilih IDPEGA utuh (tanpa substr 19,10 seperti tiga tab lain) - di sini awalan dibuang bila ada, tanpa dipotong 10 karakter.
func sqlDasarOfferFac(fo, klausa string) string {
	return "SELECT DISTINCT " + sqlCaseIDDariIDPega("PC0.IDPEGA") + " AS CASEID, PC0.NOOFFER, PC0.BEGINDATE, PC0.ENDDATE, PC0.SOB, PC0.INSUREDNAME" +
		" FROM " + fo + " PC0 WHERE (" + klausa + ")"
}

// sqlDasarPolisTreaty - BrowsePolisTreaty (DISTINCT; pzInsKey tidak dibaca - case id cukup).
func sqlDasarPolisTreaty(tip, klausa string) string {
	return "SELECT DISTINCT " + caseIDView("PC0.IDPEGA") + " AS CASEID, PC0.NOPOLIS, PC0.BEGINDATE, PC0.ENDDATE, PC0.SOB, PC0.INSUREDNAME" +
		" FROM " + tip + " PC0 WHERE (" + klausa + ")"
}

// sqlDasarPolisLife - BrowsePolisLife; kolom LIFEINPRODUCTION IDPEGA, NOPOLIS, NOENDORS, CEDINGCONAME, POLICYHOLDERNAME, TYPE,
// SOBNAME, RISLIPRNM `[belum terverifikasi tipe / lebar]`. WHERE tanpa kurung (verbatim).
func sqlDasarPolisLife(lip, klausa string) string {
	return "SELECT " + caseIDView("IDPEGA") + " AS CASEID, NOPOLIS, NOENDORS, CEDINGCONAME, POLICYHOLDERNAME AS QQ, TYPE AS TYPEPOLICY," +
		" SOBNAME, RISLIPRNM FROM " + lip + " WHERE " + klausa
}

// sqlDasarOfferLife - BrowseOfferLife verbatim (subkueri D Pega memakai JSON_OFFER_LIFE tanpa skema - di sini berkualifikasi sama);
// kolom JSON_OFFER_LIFE ID, IDPEGA, TGL_INPUT, OLDID, CEDINGCONAME, POLICYHOLDERNAME, TYPECEDING, BUSINESSNAME, STATUS
// `[belum terverifikasi tipe / lebar]`.
func sqlDasarOfferLife(jol, klausa string) string {
	return "SELECT " + caseIDView("A.IDPEGA") + " AS CASEID, 'OFR-' || D.ID AS NOOFFER, A.TGL_INPUT, A.CEDINGCONAME, A.POLICYHOLDERNAME AS QQ," +
		" A.TYPECEDING AS TYPEPOLICY, A.BUSINESSNAME, A.STATUS AS PRODKE FROM " + jol + " A, (SELECT C.ID, C.IDPEGA FROM " + jol +
		" C WHERE C.TGL_INPUT = (SELECT MIN(b.TGL_INPUT) FROM " + jol + " b WHERE b.IDPEGA = C.IDPEGA)) D WHERE A.OLDID IS NULL AND" +
		" A.TGL_INPUT = (SELECT MAX(b.TGL_INPUT) FROM " + jol + " b WHERE b.IDPEGA = A.IDPEGA) AND A.IDPEGA = D.IDPEGA AND (" + klausa + ")"
}

// tabView - satu tab: tabel, SQL dasar, kolom halaman (urut sesudah bind klausa), urutan stabil `[keputusan agent]`.
type tabView struct {
	tabel, nama string
	dasar       func(t, klausa string) string
	kolom       string
	urut        string
}

var (
	tabOfferFac = tabView{tabel: TabelFacInOffer, nama: "offer Facultative", dasar: sqlDasarOfferFac,
		kolom: "CASEID, NOOFFER, INSUREDNAME, SOB, TO_CHAR(BEGINDATE, 'DD-MM-YYYY'), TO_CHAR(ENDDATE, 'DD-MM-YYYY')",
		urut:  "NOOFFER, CASEID, BEGINDATE, ENDDATE, SOB, INSUREDNAME"}
	tabPolisTreaty = tabView{tabel: TabelTreatyInProduction, nama: "polis Treaty", dasar: sqlDasarPolisTreaty,
		kolom: "CASEID, NOPOLIS, INSUREDNAME, SOB, TO_CHAR(BEGINDATE, 'DD-MM-YYYY'), TO_CHAR(ENDDATE, 'DD-MM-YYYY')",
		urut:  "NOPOLIS, CASEID, BEGINDATE, ENDDATE, SOB, INSUREDNAME"}
	tabPolisLife = tabView{tabel: TabelLifeInProduction, nama: "polis Life", dasar: sqlDasarPolisLife,
		kolom: "CASEID, NOPOLIS, NOENDORS, TYPEPOLICY, SOBNAME, CEDINGCONAME, QQ, RISLIPRNM",
		urut:  "NOPOLIS, NOENDORS NULLS FIRST, CASEID, TYPEPOLICY, SOBNAME, CEDINGCONAME, QQ, RISLIPRNM"}
	tabOfferLife = tabView{tabel: TabelJSONOfferLife, nama: "offer Life", dasar: sqlDasarOfferLife,
		kolom: "CASEID, NOOFFER, TYPEPOLICY, BUSINESSNAME, CEDINGCONAME, QQ, PRODKE",
		urut:  "TGL_INPUT DESC, CASEID, NOOFFER"}
)

func sqlJumlahView(dasar string) string { return "SELECT COUNT(*) FROM (" + dasar + ")" }

// sqlHalamanView - kolom tab dari dasar, urut stabil, OFFSET / FETCH sesudah n bind klausa.
func sqlHalamanView(tab tabView, dasar string, n int) string {
	return "SELECT " + tab.kolom + " FROM (" + dasar + ") ORDER BY " + tab.urut + " OFFSET :" + strconv.Itoa(n+1) +
		" ROWS FETCH NEXT :" + strconv.Itoa(n+2) + " ROWS ONLY"
}

// PembacaViewPolis - satu halaman tiap tab (jenis klausa + nilai bind akhir) beserta jumlah seluruhnya.
type PembacaViewPolis interface {
	CariOfferFac(ctx context.Context, jenis, nilai string, offset, batas int) ([]models.BarisOfferFac, int, error)
	CariPolisTreaty(ctx context.Context, jenis, nilai string, offset, batas int) ([]models.BarisPolisTreaty, int, error)
	CariPolisLife(ctx context.Context, jenis, nilai string, offset, batas int) ([]models.BarisPolisLife, int, error)
	CariOfferLife(ctx context.Context, jenis, nilai string, offset, batas int) ([]models.BarisOfferLife, int, error)
	// BerkasPolisTreaty - berkas NB / EDM Treaty In terbaru untuk nomor polis; false = tanpa berkas di sistem baru.
	BerkasPolisTreaty(ctx context.Context, noPolis string) (models.BerkasPolisTreaty, bool, error)
}

// ViewPolisOracle - PembacaViewPolis atas Oracle.
type ViewPolisOracle struct{ db *db.DB }

// NewViewPolisOracle merakit pembaca tab View Policy tambahan.
func NewViewPolisOracle(d *db.DB) *ViewPolisOracle { return &ViewPolisOracle{db: d} }

// cari - jumlah + halaman satu tab; isi(n) menerima sel teks berurutan kolom tab.
func (r *ViewPolisOracle) cari(ctx context.Context, tab tabView, jenis, nilai string, offset, batas, nKolom int,
	isi func(v []string)) (int, error) {
	t, err := r.db.Qualify(tab.tabel)
	if err != nil {
		return 0, err
	}
	k, ada := klausaViewDari(jenis, t)
	if !ada {
		return 0, fmt.Errorf("repository: jenis cari %s %q tidak dikenal", tab.nama, jenis)
	}
	dasar := tab.dasar(t, k.teks)
	arg := bindView(k, nilai)
	q := sqlJumlahView(dasar)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var total int
	if err := r.db.QueryRowContext(ctx, q, arg...).Scan(&total); err != nil {
		return 0, fmt.Errorf("repository: jumlah %s: %w", tab.nama, err)
	}
	q = sqlHalamanView(tab, dasar, len(arg))
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	rows, err := r.db.QueryContext(ctx, q, append(arg, offset, batas)...)
	if err != nil {
		return 0, fmt.Errorf("repository: cari %s: %w", tab.nama, err)
	}
	defer rows.Close()
	for rows.Next() {
		v := make([]sql.NullString, nKolom)
		tuju := make([]any, nKolom)
		for i := range v {
			tuju[i] = &v[i]
		}
		if err := rows.Scan(tuju...); err != nil {
			return 0, fmt.Errorf("repository: cari %s: %w", tab.nama, err)
		}
		s := make([]string, nKolom)
		for i := range v {
			s[i] = v[i].String
		}
		isi(s)
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("repository: cari %s: %w", tab.nama, err)
	}
	return total, nil
}

// CariOfferFac - lihat PembacaViewPolis.
func (r *ViewPolisOracle) CariOfferFac(ctx context.Context, jenis, nilai string, offset, batas int) ([]models.BarisOfferFac, int, error) {
	hasil := []models.BarisOfferFac{}
	total, err := r.cari(ctx, tabOfferFac, jenis, nilai, offset, batas, 6, func(v []string) {
		hasil = append(hasil, models.BarisOfferFac{CaseID: v[0], NoOffer: v[1], InsuredName: v[2], SourceOfBusiness: v[3],
			StartPeriod: v[4], EndPeriod: v[5]})
	})
	return hasil, total, err
}

// CariPolisTreaty - lihat PembacaViewPolis.
func (r *ViewPolisOracle) CariPolisTreaty(ctx context.Context, jenis, nilai string, offset, batas int) ([]models.BarisPolisTreaty, int, error) {
	hasil := []models.BarisPolisTreaty{}
	total, err := r.cari(ctx, tabPolisTreaty, jenis, nilai, offset, batas, 6, func(v []string) {
		hasil = append(hasil, models.BarisPolisTreaty{CaseID: v[0], PolicyNo: v[1], InsuredName: v[2], SourceOfBusiness: v[3],
			StartPeriod: v[4], EndPeriod: v[5]})
	})
	return hasil, total, err
}

// CariPolisLife - lihat PembacaViewPolis.
func (r *ViewPolisOracle) CariPolisLife(ctx context.Context, jenis, nilai string, offset, batas int) ([]models.BarisPolisLife, int, error) {
	hasil := []models.BarisPolisLife{}
	total, err := r.cari(ctx, tabPolisLife, jenis, nilai, offset, batas, 8, func(v []string) {
		hasil = append(hasil, models.BarisPolisLife{CaseID: v[0], PolicyNo: v[1], EndorsNo: v[2], TypePolicy: v[3],
			SourceOfBusiness: v[4], CedingCoName: v[5], QQ: v[6], RISlipCeding: v[7]})
	})
	return hasil, total, err
}

// CariOfferLife - lihat PembacaViewPolis.
func (r *ViewPolisOracle) CariOfferLife(ctx context.Context, jenis, nilai string, offset, batas int) ([]models.BarisOfferLife, int, error) {
	hasil := []models.BarisOfferLife{}
	total, err := r.cari(ctx, tabOfferLife, jenis, nilai, offset, batas, 7, func(v []string) {
		hasil = append(hasil, models.BarisOfferLife{CaseID: v[0], NoOffer: v[1], TypePolicy: v[2], BusinessName: v[3],
			CedingCoName: v[4], QQ: v[5], LastStatus: v[6]})
	})
	return hasil, total, err
}

// sqlBerkasPolisTreaty - berkas NB / EDM Treaty In terbaru (PRODKE terbesar) untuk satu nomor polis; salinan sqlBerkasPolis Claim Non
// Prop (OQ-CNP-13) - JOIN T_WORK_POLIS = hanya berkas yang punya case.
func sqlBerkasPolisTreaty(gen, kerja string) string {
	return "SELECT g.ID, g.PRODKE FROM " + gen + " g JOIN " + kerja + " w ON w.ID = g.ID WHERE g.NOPOLIS = :1 ORDER BY g.PRODKE DESC" +
		" FETCH FIRST 1 ROWS ONLY"
}

// BerkasPolisTreaty - lihat PembacaViewPolis.
func (r *ViewPolisOracle) BerkasPolisTreaty(ctx context.Context, noPolis string) (models.BerkasPolisTreaty, bool, error) {
	gen, err := r.db.Qualify(TabelGeneralPolisTreaty)
	if err != nil {
		return models.BerkasPolisTreaty{}, false, err
	}
	kerja, err := r.db.Qualify(TabelWorkPolis)
	if err != nil {
		return models.BerkasPolisTreaty{}, false, err
	}
	q := sqlBerkasPolisTreaty(gen, kerja)
	if err := db.PeriksaSQL(q); err != nil {
		return models.BerkasPolisTreaty{}, false, err
	}
	var id string
	var prodke int64
	err = r.db.QueryRowContext(ctx, q, noPolis).Scan(&id, &prodke)
	if errors.Is(err, sql.ErrNoRows) {
		return models.BerkasPolisTreaty{}, false, nil
	}
	if err != nil {
		return models.BerkasPolisTreaty{}, false, fmt.Errorf("repository: berkas polis Treaty: %w", err)
	}
	return models.BerkasPolisTreaty{Modul: models.ModulBerkasPolisTreaty(prodke), Kasus: id}, true, nil
}
