package repository

// Tab Scoring Risk kasus FIRE (tiket 54) - halaman `.OfferFacIn.ScoringRisk` ke lima tabel RANCANGAN (migrasi 205):
// T_SCORINGRISK (induk T_GENERAL_POLIS, satu per kasus) -> T_SCORINGRESULT (satu) dan T_DATASCORINGRISKLIST (kolom Top /
// Dominant, SEQ_NO 1 / 2) -> T_SCORING_FACTOR (24 sub-faktor, FACTOR_GROUP / FACTOR_NAME = nama halaman Pega, aturan
// loader V-30) -> T_SCORING_OPTION (ChechBoxN / ScoreN: OPTION_NO = N, CHECK_BOX "true" / "false", SCORE = poin bila
// dicentang selain itu 0). Boolean = teks "true" / "false" (butir 68.1). Simpan = ganti utuh di satu transaksi.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelScoringRisk ... TabelOpsiScoring - tabel rancangan migrasi 205.
	TabelScoringRisk   = "T_SCORINGRISK"
	TabelHasilScoring  = "T_SCORINGRESULT"
	TabelDataScoring   = "T_DATASCORINGRISKLIST"
	TabelFaktorScoring = "T_SCORING_FACTOR"
	TabelOpsiScoring   = "T_SCORING_OPTION"

	sequenceScoringRisk   = "SEQ_T_SCORINGRISK"
	sequenceHasilScoring  = "SEQ_T_SCORINGRESULT"
	sequenceDataScoring   = "SEQ_T_DATASCORINGRISKLIST"
	sequenceFaktorScoring = "SEQ_T_SCORING_FACTOR"
	sequenceOpsiScoring   = "SEQ_T_SCORING_OPTION"
)

// PenyimpanScoring - baca / ganti Scoring Risk satu kasus.
type PenyimpanScoring interface {
	// BacaScoring - POSITION / FLAG_ONGOING_POLICY kasus dan masukan kolom tersimpan; ErrKasusTidakAda bila case FAC tidak
	// ada.
	BacaScoring(ctx context.Context, id string) (models.KasusScoring, error)
	// TulisScoring - ganti utuh (di transaksi pemanggil); ErrKasusTidakAda bila case FAC tidak ada.
	TulisScoring(ctx context.Context, tx *db.Tx, id string, s models.Scoring) error
}

// ScoringOracle - PenyimpanScoring atas Oracle.
type ScoringOracle struct{ db *db.DB }

// NewScoringOracle merakit penyimpan Scoring Risk.
func NewScoringOracle(d *db.DB) *ScoringOracle { return &ScoringOracle{db: d} }

type tabelScoring struct{ work, general, sr, hasil, data, faktor, opsi, seqOpsi string }

func (r *ScoringOracle) tabel() (tabelScoring, error) {
	var t tabelScoring
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelGeneralPolis, &t.general}, {TabelScoringRisk, &t.sr}, {TabelHasilScoring, &t.hasil},
		{TabelDataScoring, &t.data}, {TabelFaktorScoring, &t.faktor}, {TabelOpsiScoring, &t.opsi}, {sequenceOpsiScoring, &t.seqOpsi}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

// ---- SQL ---------------------------------------------------------------------------------------------------------

// sqlKepalaScoringCatatan - sqlKepalaScoring dengan posisi efektif (tiket 77).
func sqlKepalaScoringCatatan(work string) string {
	return "SELECT " + sqlPosisiEfektif("w") + ", w.FLAG_ONGOING_POLICY FROM " + work + " w WHERE w.ID = :1 AND w.LINI = :2"
}

func sqlKepalaScoring(work string) string {
	return "SELECT POSITION, FLAG_ONGOING_POLICY FROM " + work + " WHERE ID = :1 AND LINI = :2"
}

// subkueriData - ID T_DATASCORINGRISKLIST kasus :1.
func subkueriData(t tabelScoring) string {
	return "SELECT d.ID FROM " + t.data + " d WHERE d.PARENT_ID IN (SELECT s.ID FROM " + t.sr + " s WHERE s.PARENT_ID = :1)"
}

func sqlBacaDataScoring(t tabelScoring) string {
	return "SELECT TO_CHAR(d.ID), d.SEQ_NO, d.INDEX_LOCATION FROM " + t.data + " d WHERE d.PARENT_ID IN (SELECT s.ID FROM " +
		t.sr + " s WHERE s.PARENT_ID = :1) ORDER BY d.SEQ_NO, d.ID"
}

func sqlBacaFaktorScoring(t tabelScoring) string {
	return "SELECT TO_CHAR(f.ID), TO_CHAR(f.PARENT_ID), f.FACTOR_GROUP, f.FACTOR_NAME, " + angkaKeluar("f.PERCENT_LOSS_RATIO") +
		" FROM " + t.faktor + " f WHERE f.PARENT_ID IN (" + subkueriData(t) + ") ORDER BY f.PARENT_ID, f.SEQ_NO, f.ID"
}

// sqlBacaOpsiScoring - opsi tercentang saja (CHECK_BOX 'true', teks Pega).
func sqlBacaOpsiScoring(t tabelScoring) string {
	return "SELECT TO_CHAR(o.PARENT_ID), o.OPTION_NO FROM " + t.opsi + " o WHERE o.CHECK_BOX = 'true' AND o.PARENT_ID IN (SELECT f.ID FROM " +
		t.faktor + " f WHERE f.PARENT_ID IN (" + subkueriData(t) + "))"
}

// sqlHapusScoring - urutan anak -> induk (FK); semua ber-bind :1 = ID kasus.
func sqlHapusScoring(t tabelScoring) []string {
	faktor := "SELECT f.ID FROM " + t.faktor + " f WHERE f.PARENT_ID IN (" + subkueriData(t) + ")"
	sr := "SELECT s.ID FROM " + t.sr + " s WHERE s.PARENT_ID = :1"
	return []string{
		"DELETE FROM " + t.opsi + " WHERE PARENT_ID IN (" + faktor + ")",
		"DELETE FROM " + t.faktor + " WHERE PARENT_ID IN (" + subkueriData(t) + ")",
		"DELETE FROM " + t.data + " WHERE PARENT_ID IN (" + sr + ")",
		"DELETE FROM " + t.hasil + " WHERE PARENT_ID IN (" + sr + ")",
		"DELETE FROM " + t.sr + " WHERE PARENT_ID = :1",
	}
}

func sqlSisipScoringRisk(sr string) string {
	return "INSERT INTO " + sr + " (ID, PARENT_ID, FINAL_SCORE, NOTE_FINAL_SCORE, STATUS) VALUES (:1, :2, " + angkaMasuk(":3") +
		", :4, :5)"
}

// kolomHasilScoring - urutan kolom T_SCORINGRESULT = urutan argHasilScoring.
var kolomHasilScoring = []string{"GOOD", "ABOVE_AVERAGE", "AVERAGE", "BELOW_AVERAGE", "POOR", "ACCEPTABLE_NOT_APPROVAL",
	"ACCEPTABLE_WITH_APPROVAL", "NOT_ACCEPTABLE", "OCCUPATION", "OBJECT_CONDITIONS", "LOSS_RATIO", "OTHERS",
	"OPERATIONAL_DIV_HEAD", "OPERATIONAL_DIRECTOR", "TECHNICAL_DIRECTOR", "PRESIDENT_DIRECTOR"}

func sqlSisipHasilScoring(hasil string) string {
	q := "INSERT INTO " + hasil + " (ID, PARENT_ID"
	nilai := ":1, :2"
	for i, k := range kolomHasilScoring {
		q += ", " + k
		nilai += ", :" + strconv.Itoa(i+3)
	}
	return q + ") VALUES (" + nilai + ")"
}

func argHasilScoring(h models.HasilScoring) []any {
	arg := []any{}
	for _, b := range []bool{h.Good, h.AboveAverage, h.Average, h.BelowAverage, h.Poor, h.AcceptableNotApproval,
		h.AcceptableWithApproval, h.NotAcceptable, h.Occupation, h.ObjectConditions, h.LossRatio, h.Others,
		h.OperationalDivHead, h.OperationalDirector, h.TechnicalDirector, h.PresidentDirector} {
		arg = append(arg, teksBool(b))
	}
	return arg
}

func sqlSisipDataScoring(data string) string {
	return "INSERT INTO " + data + " (ID, PARENT_ID, SEQ_NO, ROW_UID, OCCUPATION_CODE, OCCUPATION_NOTE, SCORE, TOP_RISK_LOCATION," +
		" INDEX_LOCATION) VALUES (:1, :2, :3, :4, :5, :6, " + angkaMasuk(":7") + ", :8, :9)"
}

func sqlSisipFaktorScoring(faktor string) string {
	return "INSERT INTO " + faktor + " (ID, PARENT_ID, SEQ_NO, ROW_UID, CHECKED, FACTOR_GROUP, FACTOR_NAME, SCORE_PER_FACTOR," +
		" PERCENT_LOSS_RATIO) VALUES (:1, :2, :3, :4, :5, :6, :7, " + angkaMasuk(":8") + ", " + angkaMasuk(":9") + ")"
}

// sqlSisipOpsiScoring - ID langsung dari sequence berskema (baris daun, ID-nya tidak dipakai lagi; pola rw.go) - satu
// perjalanan ke basis data per opsi, bukan dua.
func sqlSisipOpsiScoring(opsi, seq string) string {
	return "INSERT INTO " + opsi + " (ID, PARENT_ID, SEQ_NO, ROW_UID, CHECK_BOX, OPTION_NO, SCORE) VALUES (" + seq +
		".NEXTVAL, :1, :2, :3, :4, :5, " + angkaMasuk(":6") + ")"
}

// ---- baca ---------------------------------------------------------------------------------------------------------

// faktorMenurutHalaman - (FACTOR_GROUP, FACTOR_NAME) -> kunci kontrak.
var faktorMenurutHalaman = func() map[[2]string]string {
	m := map[[2]string]string{}
	for _, f := range models.DaftarFaktorScoring {
		m[[2]string{f.Grup, f.Nama}] = f.Kunci
	}
	return m
}()

// BacaScoring - lihat PenyimpanScoring. Kolom ber-SEQ_NO selain 1 / 2 diabaikan; opsi tercentang ganda -> nomor terkecil;
// faktor tak dikenal diabaikan.
func (r *ScoringOracle) BacaScoring(ctx context.Context, id string) (models.KasusScoring, error) {
	var k models.KasusScoring
	t, err := r.tabel()
	if err != nil {
		return k, err
	}
	var posisi, flag sql.NullString
	q := sqlKepalaScoring(t.work)
	if adaCatatanPosisi(ctx, r.db, t.work) {
		q = sqlKepalaScoringCatatan(t.work)
	}
	err = r.db.QueryRowContext(ctx, q, id, LiniFacIn).Scan(&posisi, &flag)
	if errors.Is(err, sql.ErrNoRows) {
		return k, ErrKasusTidakAda
	}
	if err != nil {
		return k, fmt.Errorf("repository: baca %s: %w", TabelWorkPolis, err)
	}
	k.Position, k.FlagOnGoingPolicy = posisi.String, flag.String

	kolomData := map[string]int{} // ID baris data -> indeks kolom
	if err := r.baca(ctx, TabelDataScoring, sqlBacaDataScoring(t), id, func(b *sql.Rows) error {
		var idData, lokasi sql.NullString
		var seq sql.NullInt64
		if err := b.Scan(&idData, &seq, &lokasi); err != nil {
			return err
		}
		n := int(seq.Int64) - 1
		if n < 0 || n > 1 || k.Kolom[n].Ada {
			return nil
		}
		k.Kolom[n] = models.KolomScoring{Ada: true, IndexLocation: lokasi.String, Pilihan: map[string]int{}}
		kolomData[idData.String] = n
		return nil
	}); err != nil {
		return k, err
	}
	type faktorBaris struct {
		kolom int
		kunci string
	}
	faktor := map[string]faktorBaris{}
	if err := r.baca(ctx, TabelFaktorScoring, sqlBacaFaktorScoring(t), id, func(b *sql.Rows) error {
		var idF, idData, grup, nama, persen sql.NullString
		if err := b.Scan(&idF, &idData, &grup, &nama, &persen); err != nil {
			return err
		}
		n, ada := kolomData[idData.String]
		kunci := faktorMenurutHalaman[[2]string{grup.String, nama.String}]
		if !ada || kunci == "" {
			return nil
		}
		faktor[idF.String] = faktorBaris{n, kunci}
		if kunci == models.KunciLossRatio {
			d, err := bacaDesimal(TabelFaktorScoring+" "+idF.String, "PERCENT_LOSS_RATIO", &persen)
			if err != nil {
				return err
			}
			k.Kolom[n].PersenLossRatio = d
		}
		return nil
	}); err != nil {
		return k, err
	}
	err = r.baca(ctx, TabelOpsiScoring, sqlBacaOpsiScoring(t), id, func(b *sql.Rows) error {
		var idF, no sql.NullString
		if err := b.Scan(&idF, &no); err != nil {
			return err
		}
		f, ada := faktor[idF.String]
		n, errNo := strconv.Atoi(no.String)
		if !ada || errNo != nil || n < 1 {
			return nil
		}
		if lama := k.Kolom[f.kolom].Pilihan[f.kunci]; lama == 0 || n < lama {
			k.Kolom[f.kolom].Pilihan[f.kunci] = n
		}
		return nil
	})
	return k, err
}

// baca - satu kueri ber-bind ID kasus, setiap baris ke fn.
func (r *ScoringOracle) baca(ctx context.Context, tabel, q, id string, fn func(*sql.Rows) error) error {
	baris, err := r.db.QueryContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("repository: baca %s: %w", tabel, err)
	}
	defer baris.Close()
	for baris.Next() {
		if err := fn(baris); err != nil {
			return fmt.Errorf("repository: %s: %w", tabel, err)
		}
	}
	if err := baris.Err(); err != nil {
		return fmt.Errorf("repository: %s: %w", tabel, err)
	}
	return nil
}

// ---- tulis --------------------------------------------------------------------------------------------------------

// TulisScoring - lihat PenyimpanScoring: sentuh case (FAC), pastikan T_GENERAL_POLIS, hapus lima tabel, sisip ulang.
func (r *ScoringOracle) TulisScoring(ctx context.Context, tx *db.Tx, id string, s models.Scoring) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlSentuhCase(t.work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrKasusTidakAda
	}
	if _, err := jalankan(ctx, tx, sqlPastikanGeneral(t.general), "memastikan "+TabelGeneralPolis, id, id); err != nil {
		return err
	}
	for _, q := range sqlHapusScoring(t) {
		if _, err := jalankan(ctx, tx, q, "menghapus Scoring Risk", id); err != nil {
			return err
		}
	}
	idSR, err := r.db.NomorBerikut(ctx, tx, sequenceScoringRisk)
	if err != nil {
		return err
	}
	if _, err := jalankan(ctx, tx, sqlSisipScoringRisk(t.sr), "menyisipkan "+TabelScoringRisk, idSR, id,
		ikatDesimal(s.FinalScore), db.KosongJadiNil(s.NoteFinalScore), db.KosongJadiNil(s.Status)); err != nil {
		return err
	}
	idHasil, err := r.db.NomorBerikut(ctx, tx, sequenceHasilScoring)
	if err != nil {
		return err
	}
	if _, err := jalankan(ctx, tx, sqlSisipHasilScoring(t.hasil), "menyisipkan "+TabelHasilScoring,
		append([]any{idHasil, idSR}, argHasilScoring(s.Hasil)...)...); err != nil {
		return err
	}
	for n, k := range s.Kolom {
		if !k.Ada {
			continue
		}
		if err := r.tulisKolom(ctx, tx, t, idSR, n+1, k); err != nil {
			return err
		}
	}
	return nil
}

// tulisKolom - satu baris T_DATASCORINGRISKLIST, 24 faktor, dan seluruh opsinya.
func (r *ScoringOracle) tulisKolom(ctx context.Context, tx *db.Tx, t tabelScoring, idSR string, seq int, k models.KolomScoring) error {
	idData, err := r.db.NomorBerikut(ctx, tx, sequenceDataScoring)
	if err != nil {
		return err
	}
	uid, err := uidAcak()
	if err != nil {
		return err
	}
	if _, err := jalankan(ctx, tx, sqlSisipDataScoring(t.data), "menyisipkan "+TabelDataScoring, idData, idSR, seq, uid,
		db.KosongJadiNil(k.OccupationCode), db.KosongJadiNil(k.OccupationNote), ikatDesimal(k.Total),
		db.KosongJadiNil(k.TopRiskLocation), db.KosongJadiNil(k.IndexLocation)); err != nil {
		return err
	}
	for i, f := range models.DaftarFaktorScoring {
		idF, err := r.db.NomorBerikut(ctx, tx, sequenceFaktorScoring)
		if err != nil {
			return err
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		pilih := k.Pilihan[f.Kunci]
		var persen *apd.Decimal
		if f.Kunci == models.KunciLossRatio {
			persen = k.PersenLossRatio
		}
		if _, err := jalankan(ctx, tx, sqlSisipFaktorScoring(t.faktor), "menyisipkan "+TabelFaktorScoring, idF, idData, i+1, uid,
			teksBool(pilih > 0), db.KosongJadiNil(f.Grup), f.Nama, ikatDesimal(k.Skor[f.Kunci]), ikatDesimal(persen)); err != nil {
			return err
		}
		for o := 1; o <= f.JumlahOpsi; o++ {
			uid, err := uidAcak()
			if err != nil {
				return err
			}
			skor := apd.New(0, 0)
			if o == pilih && k.Skor[f.Kunci] != nil {
				skor = k.Skor[f.Kunci]
			}
			if _, err := jalankan(ctx, tx, sqlSisipOpsiScoring(t.opsi, t.seqOpsi), "menyisipkan "+TabelOpsiScoring, idF, o, uid,
				teksBool(o == pilih), strconv.Itoa(o), ikatDesimal(skor)); err != nil {
				return err
			}
		}
	}
	return nil
}
