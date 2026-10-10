package repository

// Untuk apa berkas ini: TULISAN KE TABEL WARISAN di titik yang sama dengan XML, ditulis ulang tanpa procedure (isi
// procedure dibaca dari ALL_SOURCE DEV 10-10-2026):
//
//	PEGA_JSON_OS_AKSEP_KLAIM (SaveOSClaim_SQL)    -> OS_AKSEPTASI_KLAIM: SELALU INSERT CASEID, NOCLAIM, DATA_JSON,
//	                                                TANGGAL (hari ini), NOPOLIS, STS_REJECT, STS_KONVERSI, STS_DLA,
//	                                                MASTERID, CLAIMOLD (hitungan baris CASEID procedure tidak dipakai)
//	PEGA_JSON_KLAIM_PNC (InsertClaimPNC)          -> JSON_KLAIM (baris IDPEGA yang ada dibiarkan, selainnya INSERT;
//	                                                DATA_JSON tidak diisi - keputusan work owner)
//	PEGA_PROGRESSCLAIM (InsertProgressClaim_SQL)  -> PROGRESSCLAIM: INSERT bila (IDPEGA, POSITION) belum ada, lalu posisi
//	                                                lain STATUS 'Done'
//	PEGA_SUBPROGRESSCLAIM                         -> SUBPROGRESSCLAIM: INSERT bila belum ada baris (IDPROGRES, IDPEGA,
//	                                                JENISPROGRES) ber-STATUS 'Auto Create%'
//	InsertLogServiceClaim                         -> MONITORING_KLAIM_LOG
//	SaveCatasrtope_Act (Obj-Save)                 -> CATASTROPHE
//	efek keluar                                   -> T_LOG_SERVICE_RNM (inti/backend/outbox)

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/claimfacin/backend/models"
)

// hariJakarta - `TO_DATE(to_char(sysdate,'dd/MM/yyyy'),'dd/MM/yyyy')`: tanggal hari aksi tanpa jam.
func hariJakarta(t time.Time) time.Time {
	l := t.In(models.Jakarta)
	return time.Date(l.Year(), l.Month(), l.Day(), 0, 0, 0, 0, time.UTC)
}

// sqlSisipOS - isi PEGA_JSON_OS_AKSEP_KLAIM baris 16 (tanpa COMMIT). STS_KONVERSI / MASTERID / CLAIMOLD = NULL
// (`InputData.CARI16 = ""`, CARI18 / CARI20 tidak diisi pemanggil Claim Fac In).
func sqlSisipOS(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, STS_DLA)
		VALUES (:1, :2, :3, :4, :5, :6, :7)`, tabel)
}

// SisipOS menyisipkan satu baris OS_AKSEPTASI_KLAIM (TGL_PROD diisi trigger TRG_TLG_PROD_OS_AKSEPTASI).
func (g *Gudang) SisipOS(ctx context.Context, tx *db.Tx, b models.BarisOS, saat time.Time) error {
	tabel, err := g.db.Qualify("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return err
	}
	q := sqlSisipOS(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, b.CaseID, teksAtauNil(b.NoClaim), teksAtauNil(b.DataJSON), hariJakarta(saat),
		teksAtauNil(b.NoPolis), teksAtauNil(b.StsReject), teksAtauNil(b.StsDLA))
	if err != nil {
		return fmt.Errorf("repository: menyisipkan OS_AKSEPTASI_KLAIM: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "baris OS_AKSEPTASI_KLAIM")
}

// sqlTandaiDLAOS = InsertDLA_OS_SQL (`UPDATE os_akseptasi_klaim SET DLA_NO, DLA_DATE = SYSDATE WHERE
// data_json.AcceptedNo = :no`; tanpa COMMIT). Nomor akseptasi unik lintas kasus (dibentuk komite) - kunci Pega
// dipertahankan; DLA_DATE = waktu aksi (bukan SYSDATE server basis data).
func sqlTandaiDLAOS(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET DLA_NO = :1, DLA_DATE = :2 WHERE JSON_VALUE(DATA_JSON, '$.AcceptedNo') = :3`, tabel)
}

// TandaiDLAOS menandai baris OS akseptasi `noAkseptasi` dengan nomor DLA (GenerateDLAFacin_Act 20.1.2). Nol baris =
// akseptasi belum ditulis komite (tahap 2) - bukan galat.
func (g *Gudang) TandaiDLAOS(ctx context.Context, tx *db.Tx, noAkseptasi, noDLA string, saat time.Time) error {
	if noAkseptasi == "" {
		return nil
	}
	tabel, err := g.db.Qualify("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return err
	}
	q := sqlTandaiDLAOS(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, noDLA, saat, noAkseptasi); err != nil {
		return fmt.Errorf("repository: menandai DLA OS_AKSEPTASI_KLAIM: %w", err)
	}
	return nil
}

// sqlBacaOS - baris OS nomor klaim ini (pop-up Outstanding Summary: GetDataOutstanding_SQL dan kawan-kawan). CASEID
// kasus Pega lama ber-pzInsKey dan kasus sistem baru ber-ID apa adanya (prompt §6 butir 6).
func sqlBacaOS(tabel string) string {
	return fmt.Sprintf(`SELECT NOCLAIM, NOPOLIS, TO_CHAR(TANGGAL, 'YYYY-MM-DD'), TO_CHAR(STS_REJECT), DATA_JSON FROM %s
		 WHERE NOCLAIM = :1 AND NOCLAIM IS NOT NULL AND STS_REJECT IN (0, 1) AND CASEID IN (:2, :3)
		 ORDER BY TANGGAL, ROWID`, tabel)
}

// BacaOS membaca baris OS_AKSEPTASI_KLAIM nomor klaim kasus `kasusID` untuk pop-up Outstanding Summary.
func (g *Gudang) BacaOS(ctx context.Context, kasusID, noKlaim string) ([]models.BarisOSTersimpan, error) {
	tabel, err := g.db.Qualify("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return nil, err
	}
	q := sqlBacaOS(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, noKlaim, models.KunciInstans(kasusID), models.KunciPegaLama(kasusID))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca OS_AKSEPTASI_KLAIM: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []models.BarisOSTersimpan
	for rows.Next() {
		var n [5]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4]); err != nil {
			return nil, fmt.Errorf("repository: memindai OS_AKSEPTASI_KLAIM: %w", err)
		}
		out = append(out, models.BarisOSTersimpan{NoClaim: n[0].String, NoPolis: n[1].String, Tanggal: n[2].String,
			StsReject: n[3].String, DataJSON: n[4].String})
	}
	return out, rows.Err()
}

// sqlAdaJSONKlaim / sqlSisipJSONKlaim - isi PEGA_JSON_KLAIM_PNC tanpa DATA_JSON.
func sqlAdaJSONKlaim(tabel string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE IDPEGA = :1`, tabel)
}

func sqlSisipJSONKlaim(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (MNK_NO_KLAIM, IDPEGA, TGL_INPUT, NOPOLIS, IDPROD) VALUES (:1, :2, :3, :4, 1)`, tabel)
}

// SalinJSONKlaim = PEGA_JSON_KLAIM_PNC (InsertJsonClaimNonMBU_act): baris IDPEGA yang sudah ada dibiarkan (UPDATE
// procedure hanya menyetel DATA_JSON, yang tidak diisi), selainnya INSERT (MNK_NO_KLAIM, IDPEGA, TGL_INPUT, NOPOLIS,
// IDPROD 1).
func (g *Gudang) SalinJSONKlaim(ctx context.Context, tx *db.Tx, idPega, noKlaim, noPolis string, saat time.Time) error {
	tabel, err := g.db.Qualify("JSON_KLAIM")
	if err != nil {
		return err
	}
	qAda, q := sqlAdaJSONKlaim(tabel), sqlSisipJSONKlaim(tabel)
	for _, x := range []string{qAda, q} {
		if err := db.PeriksaSQL(x); err != nil {
			return err
		}
	}
	var n int
	if err := tx.QueryRowContext(ctx, qAda, idPega).Scan(&n); err != nil {
		return fmt.Errorf("repository: memeriksa JSON_KLAIM: %w", err)
	}
	if n > 0 {
		return nil
	}
	hasil, err := tx.ExecContext(ctx, q, teksAtauNil(noKlaim), idPega, saat, teksAtauNil(noPolis))
	if err != nil {
		return fmt.Errorf("repository: menyisipkan JSON_KLAIM: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "baris JSON_KLAIM")
}

// sqlSisipKatastrofe - satu baris CATASTROPHE.
func sqlSisipKatastrofe(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, NOTE, KLAIMTYPE, TGL_INPUT, USER_INPUT, STSKATASTROFE, NONKATASTROFETYPE)
		VALUES (:1, :2, :3, :4, :5, :6, :7)`, tabel)
}

// SisipKatastrofe = SaveCatasrtope_Act langkah 3 (Obj-Save CATASTROPHE).
func (g *Gudang) SisipKatastrofe(ctx context.Context, tx *db.Tx, k models.BarisKatastrofe, saat time.Time) error {
	tabel, err := g.db.Qualify("CATASTROPHE")
	if err != nil {
		return err
	}
	q := sqlSisipKatastrofe(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, k.ID, teksAtauNil(potong(k.Note, 2000)), k.KlaimType, saat,
		teksAtauNil(k.UserInput), teksAtauNil(k.StsKatastrofe), teksAtauNil(k.NonKatastrofeType))
	if err != nil {
		return fmt.Errorf("repository: menyisipkan katastrofe: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "baris katastrofe")
}

// LogLayanan - satu baris MONITORING_KLAIM_LOG (InsertLogServiceClaim).
type LogLayanan struct {
	IDPega, Parameter, JenisService, NoAkseptasi, NoDLA, StsMessage, ResponMessage string
}

// sqlLogLayanan - satu baris MONITORING_KLAIM_LOG.
func sqlLogLayanan(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (TGL_INPUT, IDPEGA, PARAMETER, JN_SERVICE, NO_AKSEPTASI, NO_DLA, STS_MESSAGE,
		RESPON_MESSAGE) VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, tabel)
}

// CatatLogLayanan = InsertLogServiceClaim (CLaimFaceSheet_Act 48-49, ChooseDla_Act).
func (g *Gudang) CatatLogLayanan(ctx context.Context, tx *db.Tx, l LogLayanan, saat time.Time) error {
	tabel, err := g.db.Qualify("MONITORING_KLAIM_LOG")
	if err != nil {
		return err
	}
	q := sqlLogLayanan(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, saat, teksAtauNil(l.IDPega), teksAtauNil(potong(l.Parameter, 1000)),
		teksAtauNil(l.JenisService), teksAtauNil(l.NoAkseptasi), teksAtauNil(l.NoDLA),
		teksAtauNil(potong(l.StsMessage, 500)), teksAtauNil(potong(l.ResponMessage, 1000)))
	if err != nil {
		return fmt.Errorf("repository: mencatat log layanan klaim: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "log layanan klaim")
}

// sqlSisipProgres / sqlSelesaiProgres - isi PEGA_PROGRESSCLAIM (tanpa COMMIT).
func sqlSisipProgres(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (IDPEGA, POSITION, TANGGAL, STATUS)
		SELECT :1, :2, :3, :4 FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM %s WHERE IDPEGA = :5 AND POSITION = :6)`, tabel, tabel)
}

func sqlSelesaiProgres(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET STATUS = 'Done' WHERE IDPEGA = :1 AND POSITION <> :2`, tabel)
}

// sqlSisipSubProgres - isi PEGA_SUBPROGRESSCLAIM (tanpa COMMIT).
func sqlSisipSubProgres(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (IDPROGRES, IDPEGA, JENISPROGRES, TANGGALINPUT, POSITION1, POSITION2, NEXT_FU,
		USER_INPUT, STATUS)
		SELECT :1, :2, :3, :4, :5, :6, :7, :8, :9 FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM %s
		 WHERE IDPEGA = :10 AND IDPROGRES = :11 AND JENISPROGRES = :12 AND STATUS LIKE 'Auto Create%%')`, tabel, tabel)
}

// detik - `To_date(..., 'DD/MM/YYYY HH24:MI:SS')`: waktu Jakarta tanpa pecahan detik.
func detik(t time.Time) time.Time {
	l := t.In(models.Jakarta)
	return time.Date(l.Year(), l.Month(), l.Day(), l.Hour(), l.Minute(), l.Second(), 0, time.UTC)
}

// TulisProgres = InsertProgressClaim (PEGA_PROGRESSCLAIM + PEGA_SUBPROGRESSCLAIM). Procedure PROGRESSCLAIM pada baris
// yang sudah ada KELUAR tanpa UPDATE 'Done' (RETURN di cabang id_count > 0) - ditiru.
func (g *Gudang) TulisProgres(ctx context.Context, tx *db.Tx, p models.Progres, s models.SubProgres) error {
	tp, err := g.db.Qualify("PROGRESSCLAIM")
	if err != nil {
		return err
	}
	ts, err := g.db.Qualify("SUBPROGRESSCLAIM")
	if err != nil {
		return err
	}
	q1, q2, q3 := sqlSisipProgres(tp), sqlSelesaiProgres(tp), sqlSisipSubProgres(ts)
	for _, q := range []string{q1, q2, q3} {
		if err := db.PeriksaSQL(q); err != nil {
			return err
		}
	}
	hasil, err := tx.ExecContext(ctx, q1, p.IDPega, p.Posisi, detik(p.Tanggal), p.Status, p.IDPega, p.Posisi)
	if err != nil {
		return fmt.Errorf("repository: menulis PROGRESSCLAIM: %w", err)
	}
	if n, _ := hasil.RowsAffected(); n > 0 {
		if _, err := tx.ExecContext(ctx, q2, p.IDPega, p.Posisi); err != nil {
			return fmt.Errorf("repository: menutup posisi PROGRESSCLAIM: %w", err)
		}
	}
	if _, err := tx.ExecContext(ctx, q3, s.IDProgres, s.IDPega, s.Jenis, detik(s.Tanggal), teksAtauNil(s.Posisi1),
		teksAtauNil(s.Posisi2), detik(s.TindakLanjut), teksAtauNil(potong(s.Pengguna, 150)), s.Status, s.IDPega,
		s.IDProgres, s.Jenis); err != nil {
		return fmt.Errorf("repository: menulis SUBPROGRESSCLAIM: %w", err)
	}
	return nil
}

// Lini dan modul outbox Claim Fac In.
const (
	LiniOutbox  = models.LiniFacIn
	ModulOutbox = "claimfacin"
)

// AntreEfek menulis satu efek keluar ke outbox di transaksi aksi.
func (g *Gudang) AntreEfek(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) (string, error) {
	return outbox.NewPenyimpan(g.db).AntreEfek(ctx, tx, LiniOutbox, ModulOutbox, jenis, rujukan, muatan, saat)
}
