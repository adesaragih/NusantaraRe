package repository

// Untuk apa berkas ini: DAFTAR SOURCE OF BUSINESS pemilih XOL Retro - RD
// `BrowseAgentHierarkiList_RD` (kelas `ASM-FW-GISFW-Int-AGENT`, tabel AGENT)
// SEBAGAIMANA EFEKTIF di kasus NB. HANYA MEMBACA.
//
// RD: kolom `.ID .ClientName .Leader0 .ChildCount .ClientID`; filter
// `A: .Leader0 = Param.Leader` (pyUseNullIfEmpty=true) AND `B: .StatusActive IS
// NULL`; urut `.ClientName` ASC; pyMaxRecords 10000; paging mati.
//
// ⭐ `Param.Leader` TIDAK PERNAH TERISI di NB: satu-satunya pengisinya
// `AgentSourceBizTreatyIn_Act` langkah 1, bersyarat
// `pyWorkPage.OfferTreatyIn.QuotationData.btnQuotation=="SOB"` - properti yang
// ditulis NOL rule korpus (btnSOB_DT menulis `pyWorkPage.Quotation.btnQuotation`).
// Langkah 3 (`pxShowReport`, pyStepsPreCondition=false = selalu jalan) lalu
// memanggil RD dengan Leader kosong; pyUseNullIfEmpty menjadikan filter A
// `LEADER0 IS NULL`. Hasilnya daftar akar - untuk akar pohon DAN untuk setiap
// simpul yang diperluas (pyDeferLoadActivity yang sama tanpa parameter).
//
// `[data DBA - belum dikonfirmasi]` nama kolom CLIENTNAME, LEADER0, CHILDCOUNT,
// CLIENTID diambil dari nama properti RD (pola `StsPKPAgen`); ID dan
// STATUSACTIVE sudah dicek di katalog `ALL_TAB_COLUMNS` (prompt putaran 2 bab 1).
// Setiap kolom dibaca lewat TO_CHAR supaya angka tidak lewat float.
// ⚠️ ORDER BY ditambah ID sesudah CLIENTNAME supaya urutan nama kembar tetap -
// penyimpangan kecil yang tidak mengubah isi daftar.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// batasDaftarAgen = pyMaxRecords RD `BrowseAgentHierarkiList_RD`.
const batasDaftarAgen = 10000

// kolomAgen - kelima kolom RD, berurutan seperti `models.BarisAgen`.
const kolomAgen = `TO_CHAR(ID), TO_CHAR(CLIENTNAME), TO_CHAR(LEADER0), TO_CHAR(CHILDCOUNT), TO_CHAR(CLIENTID)`

func pindaiAgen(sc interface{ Scan(...any) error }) (models.BarisAgen, error) {
	var v [5]sql.NullString
	if err := sc.Scan(&v[0], &v[1], &v[2], &v[3], &v[4]); err != nil {
		return models.BarisAgen{}, err
	}
	return models.BarisAgen{ID: v[0].String, ClientName: v[1].String, Leader0: v[2].String,
		ChildCount: v[3].String, ClientID: v[4].String}, nil
}

// DaftarAgenHierarki = RD `BrowseAgentHierarkiList_RD` dengan `Param.Leader`
// kosong (keadaan NB) - grid TreeGrid `Section/SourceHierarki`, dan pencarian
// yang DIJALANKAN ULANG saat Save/Submit untuk mencocokkan pilihan Source Of
// Business yang dipegang layar (F4, `services.terimaSumberBisnis`).
func (g *Gudang) DaftarAgenHierarki(ctx context.Context) ([]models.BarisAgen, error) {
	t, err := g.nama(tabelAgen)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE LEADER0 IS NULL AND STATUSACTIVE IS NULL
	  ORDER BY CLIENTNAME, ID FETCH FIRST %d ROWS ONLY`, kolomAgen, t, batasDaftarAgen)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca daftar sumber bisnis: %w", err)
	}
	defer rows.Close()
	var out []models.BarisAgen
	for rows.Next() {
		b, err := pindaiAgen(rows)
		if err != nil {
			return nil, fmt.Errorf("repository: membaca daftar sumber bisnis: %w", err)
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// AgenHierarki membaca SATU baris daftar yang sama (filter RD + `ID`), supaya
// baris yang diklik dibaca ulang server - bukan dipercaya dari layar. `ada`
// false = ID itu tidak ada di daftar.
func (g *Gudang) AgenHierarki(ctx context.Context, id string) (models.BarisAgen, bool, error) {
	t, err := g.nama(tabelAgen)
	if err != nil {
		return models.BarisAgen{}, false, err
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE LEADER0 IS NULL AND STATUSACTIVE IS NULL AND TO_CHAR(ID) = :1
	  FETCH FIRST 1 ROWS ONLY`, kolomAgen, t)
	if err := db.PeriksaSQL(q); err != nil {
		return models.BarisAgen{}, false, err
	}
	b, err := pindaiAgen(g.db.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return models.BarisAgen{}, false, nil
	}
	if err != nil {
		return models.BarisAgen{}, false, fmt.Errorf("repository: membaca sumber bisnis: %w", err)
	}
	return b, true, nil
}
