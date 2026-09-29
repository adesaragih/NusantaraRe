package repository

// Kontrak BACA hilir atas master arrangement - tiket 01, keputusan tco3.
//
// Untuk apa berkas ini: `Claim Prop`, `Komite Claim Prop`, dan `Claim Fac In`
// membaca `PROPORTIONALARRG`, `TREATYREINSURER`, `TREATYBUSINESS`, dan
// `TREATYCONTRACT` warisan hari ini (spec b224-b232, OQ-042). Ketiganya belum
// dimigrasi. Kontraknya ditulis DI SINI - satu pembaca per kueri hilir, dengan
// kolom VERBATIM kueri itu - supaya saat mereka dibangun, mereka membaca tabel
// warisan yang SAMA (tco4) lewat pembaca yang sudah dikunci uji.
//
// Sumber tiap pembaca `[terverifikasi]` (path korpus, kolom dan WHERE
// disalin apa adanya):
//
//	KlausulUntukHilir        Claim Prop/RDBList/GetLimitPLATreatyin.xml,
//	                         Claim Fac In/RDBList/GetLimitPLADLA_Sql.xml
//	KlausulIndukUntukHilir   Claim Fac In/RDBList/GetQuotaShare.xml
//	                         (PARENTREINSTYPEID, bukan REINSTYPEID)
//	ReinsurerUntukHilir      {Claim Prop, Komite Claim Prop, Claim Fac In}/
//	                         RDBList/GetListRetro_Sql.xml
//	BusinessUntukHilir       Claim Fac In/RDBList/GetTreatyGroup_Sql.xml
//	GrupTreatyAktifUntukHilir Claim Prop/RDBList/GetTreatyGroupID.xml
//	LimitTreatyUntukHilir    Claim Fac In/RDBList/GetDataTreatyLimit_Sql.xml
//
// ⛔ BACA SAJA. Nol INSERT/UPDATE/DELETE/MERGE di berkas ini - modul ini
// penulis tunggal, dan pembaca hilir tidak pernah menulis (AC 1, 2).
//
// ⛔ Nilai TREATYDESCID '10001' yang hilir tulis sebagai literal TIDAK ditanam
// di sini: ia parameter, milik pemanggil (pola penyimpangan sadar 7).
//
// Uang dan persen keluar sebagai TEKS desimal bertitik, nol float. tco4:
// `PCT`/`RP`/`USD` PROPORTIONALARRG VARCHAR2 - dibaca apa adanya seperti
// kueri hilir Pega (`a.Pct`), dinormalkan di Go (`UraiDesimalWarisanTCO`);
// `RICOMM`/`PCTSHARE` NUMBER - TO_CHAR ber-NLS.

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// KontrakHilirTCO adalah pembaca read-only untuk konteks hilir.
type KontrakHilirTCO struct{ db *DB }

// NewKontrakHilirTCO menyusunnya.
func NewKontrakHilirTCO(db *DB) *KontrakHilirTCO { return &KontrakHilirTCO{db: db} }

// KlausulHilir adalah satu baris PROPORTIONALARRG seperti dibaca hilir.
type KlausulHilir struct {
	TreatyYear    string
	TreatyGroupID string
	TreatyDescID  string
	ReinsTypeID   string
	ReinsTypeName string
	Pct           string
	Rp            string
	Usd           string
}

// ReinsurerHilir adalah satu baris TREATYREINSURER seperti dibaca hilir.
type ReinsurerHilir struct {
	ID          string
	ReinsurerID string
	ClientID    string
	Name        string
	Ricomm      string
	PctShare    string
}

// BusinessHilir adalah satu baris TREATYBUSINESS seperti dibaca hilir.
type BusinessHilir struct {
	ID            string
	TreatyYear    string
	TreatyGroupID string
	ReinsTypeID   string
}

// Kolom VERBATIM tiap kueri hilir - dikunci uji terhadap korpus dan DDL.
var (
	KolomKlausulHilir   = []string{"TREATYYEAR", "TREATYGROUPID", "TREATYDESCID", "REINSTYPEID", "REINSTYPENAME", "PCT", "RP", "USD"}
	KolomReinsurerHilir = []string{"ID", "REINSURERID", "CLIENTID", "NAME", "RICOMM", "PCTSHARE"}
	KolomBusinessHilir  = []string{"ID", "TREATYYEAR", "TREATYGROUPID", "REINSTYPEID"}
)

const daftarPilihKlausulHilir = `TREATYYEAR, TREATYGROUPID, TREATYDESCID, REINSTYPEID, REINSTYPENAME,
	       PCT, RP, USD`

// sqlKlausulHilir - GetLimitPLATreatyin / GetLimitPLADLA_Sql.
func sqlKlausulHilir(tabel string) string {
	return fmt.Sprintf(`SELECT %s
	  FROM %s
	 WHERE TREATYDESCID = :1 AND TREATYYEAR = :2 AND TREATYGROUPID = :3 AND REINSTYPEID = :4
	 ORDER BY ID`, daftarPilihKlausulHilir, tabel)
}

// sqlKlausulIndukHilir - GetQuotaShare: disaring PARENTREINSTYPEID.
func sqlKlausulIndukHilir(tabel string) string {
	return fmt.Sprintf(`SELECT %s
	  FROM %s
	 WHERE TREATYYEAR = :1 AND TREATYGROUPID = :2 AND TREATYDESCID = :3 AND PARENTREINSTYPEID = :4
	 ORDER BY ID`, daftarPilihKlausulHilir, tabel)
}

// sqlReinsurerHilir - GetListRetro_Sql (tiga konteks, teks identik).
func sqlReinsurerHilir(tabel string) string {
	return fmt.Sprintf(`SELECT ID, REINSURERID, CLIENTID, NAME,
	       TO_CHAR(RICOMM, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''),
	       TO_CHAR(PCTSHARE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')
	  FROM %s
	 WHERE REINSTYPEID = :1 AND TREATYYEAR = :2 AND TREATYGROUPID = :3
	 ORDER BY ID`, tabel)
}

// sqlBusinessHilir - GetTreatyGroup_Sql.
func sqlBusinessHilir(tabel string) string {
	return fmt.Sprintf(`SELECT ID, TREATYYEAR, TREATYGROUPID, REINSTYPEID
	  FROM %s
	 WHERE BIZCODE = :1 AND TREATYYEAR = :2 AND REINSTYPEID = :3
	 ORDER BY ID`, tabel)
}

// sqlGrupTreatyAktifHilir - GetTreatyGroupID: `isactive='1'` VERBATIM.
func sqlGrupTreatyAktifHilir(tabel string) string {
	return fmt.Sprintf(`SELECT DISTINCT TREATYGROUPID
	  FROM %s
	 WHERE BIZCODE = :1 AND TREATYYEAR = :2 AND ISACTIVE = '1'
	 ORDER BY TREATYGROUPID`, tabel)
}

// sqlLimitTreatyHilir - GetDataTreatyLimit_Sql: business JOIN klausul,
// dibatasi kontrak yang berlaku pada tanggal.
func sqlLimitTreatyHilir(business, klausul, kontrak string) string {
	return fmt.Sprintf(`SELECT p.TREATYYEAR, p.TREATYGROUPID, p.TREATYDESCID, p.REINSTYPEID, p.REINSTYPENAME,
	       p.PCT, p.RP, p.USD
	  FROM %s b
	  JOIN %s p
	    ON p.TREATYYEARID = b.TREATYYEARID
	   AND p.REINSTYPEID = b.REINSTYPEID
	 WHERE b.BIZCODE = :1
	   AND b.REINSTYPEID = :2
	   AND p.TREATYDESCID = :3
	   AND EXISTS (
	         SELECT 1
	           FROM %s tc
	          WHERE tc.IDTREATYYEAR = b.TREATYYEARID
	            AND tc.REINSTYPEID = b.REINSTYPEID
	            AND :4 BETWEEN tc.TREATYSTARTDATE AND tc.TREATYENDDATE
	       )
	 ORDER BY p.ID`, business, klausul, kontrak)
}

func (r *KontrakHilirTCO) bacaKlausul(ctx context.Context, q string, arg ...any) ([]KlausulHilir, error) {
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.sql.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca klausul untuk hilir: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []KlausulHilir
	for rows.Next() {
		var n [8]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7]); err != nil {
			return nil, err
		}
		var angka [3]string
		for i, kolom := range []string{"PCT", "RP", "USD"} {
			d, err := desimalWarisanTCO(n[5+i], kolom)
			if err != nil {
				return nil, err
			}
			if v := TulisDesimalWarisanTCO(d); v != nil {
				angka[i] = v.(string)
			}
		}
		out = append(out, KlausulHilir{
			TreatyYear: n[0].String, TreatyGroupID: n[1].String, TreatyDescID: n[2].String,
			ReinsTypeID: n[3].String, ReinsTypeName: n[4].String,
			Pct: angka[0], Rp: angka[1], Usd: angka[2],
		})
	}
	return out, rows.Err()
}

// KlausulUntukHilir membaca klausul satu jenis pada satu kombinasi.
func (r *KontrakHilirTCO) KlausulUntukHilir(ctx context.Context,
	treatyDescID, treatyYear, treatyGroupID, reinsTypeID string) ([]KlausulHilir, error) {

	tabel, err := r.db.Qualify(TabelKlausulTCO)
	if err != nil {
		return nil, err
	}
	return r.bacaKlausul(ctx, sqlKlausulHilir(tabel), treatyDescID, treatyYear, treatyGroupID, reinsTypeID)
}

// KlausulIndukUntukHilir membaca klausul satu jenis menurut induknya.
func (r *KontrakHilirTCO) KlausulIndukUntukHilir(ctx context.Context,
	treatyYear, treatyGroupID, treatyDescID, parentReinsTypeID string) ([]KlausulHilir, error) {

	tabel, err := r.db.Qualify(TabelKlausulTCO)
	if err != nil {
		return nil, err
	}
	return r.bacaKlausul(ctx, sqlKlausulIndukHilir(tabel), treatyYear, treatyGroupID, treatyDescID, parentReinsTypeID)
}

// LimitTreatyUntukHilir membaca klausul lewat business, dibatasi kontrak
// yang berlaku pada `tanggal`.
func (r *KontrakHilirTCO) LimitTreatyUntukHilir(ctx context.Context,
	bizCode, reinsTypeID, treatyDescID string, tanggal time.Time) ([]KlausulHilir, error) {

	business, err := r.db.Qualify(TabelBusinessTCO)
	if err != nil {
		return nil, err
	}
	klausul, err := r.db.Qualify(TabelKlausulTCO)
	if err != nil {
		return nil, err
	}
	kontrak, err := r.db.Qualify(TabelKontrakTCO)
	if err != nil {
		return nil, err
	}
	return r.bacaKlausul(ctx, sqlLimitTreatyHilir(business, klausul, kontrak),
		bizCode, reinsTypeID, treatyDescID, tanggal)
}

// ReinsurerUntukHilir membaca reinsurer satu kombinasi.
func (r *KontrakHilirTCO) ReinsurerUntukHilir(ctx context.Context,
	reinsTypeID, treatyYear, treatyGroupID string) ([]ReinsurerHilir, error) {

	tabel, err := r.db.Qualify(TabelReinsurerTCO)
	if err != nil {
		return nil, err
	}
	q := sqlReinsurerHilir(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.sql.QueryContext(ctx, q, reinsTypeID, treatyYear, treatyGroupID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca reinsurer untuk hilir: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []ReinsurerHilir
	for rows.Next() {
		var n [6]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5]); err != nil {
			return nil, err
		}
		out = append(out, ReinsurerHilir{ID: n[0].String, ReinsurerID: n[1].String,
			ClientID: n[2].String, Name: n[3].String, Ricomm: n[4].String, PctShare: n[5].String})
	}
	return out, rows.Err()
}

// BusinessUntukHilir membaca baris business satu kode bisnis.
func (r *KontrakHilirTCO) BusinessUntukHilir(ctx context.Context,
	bizCode, treatyYear, reinsTypeID string) ([]BusinessHilir, error) {

	tabel, err := r.db.Qualify(TabelBusinessTCO)
	if err != nil {
		return nil, err
	}
	q := sqlBusinessHilir(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.sql.QueryContext(ctx, q, bizCode, treatyYear, reinsTypeID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca business untuk hilir: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []BusinessHilir
	for rows.Next() {
		var n [4]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3]); err != nil {
			return nil, err
		}
		out = append(out, BusinessHilir{ID: n[0].String, TreatyYear: n[1].String,
			TreatyGroupID: n[2].String, ReinsTypeID: n[3].String})
	}
	return out, rows.Err()
}

// GrupTreatyAktifUntukHilir membaca grup treaty aktif satu kode bisnis.
func (r *KontrakHilirTCO) GrupTreatyAktifUntukHilir(ctx context.Context,
	bizCode, treatyYear string) ([]string, error) {

	tabel, err := r.db.Qualify(TabelBusinessTCO)
	if err != nil {
		return nil, err
	}
	q := sqlGrupTreatyAktifHilir(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.sql.QueryContext(ctx, q, bizCode, treatyYear)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca grup treaty untuk hilir: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []string
	for rows.Next() {
		var n sql.NullString
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		out = append(out, n.String)
	}
	return out, rows.Err()
}
