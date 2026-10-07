package repository

// Tab Correspondence (tiket 56) - tabel BARU migrasi 207 (tidak ada di rancangan loader / DDL; sampel CONTOH semua kosong):
// T_CORRESPONDENCELIST (induk T_GENERAL_POLIS, urut SEQ_NO = urutan Add) -> T_RECIPIENTLIST (penerima, urut SEQ_NO). Pesan
// = HTML yang sudah dibersihkan layanan (VARCHAR2(4000)). Status Pega `.IsEditable` "1" draf / "0" terkirim. Akun pengirim dan saran
// penerima dibaca dari M_LOGIN_GO (inti login, baca saja).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelKorespondensi, TabelPenerimaKorespondensi - tabel migrasi 207.
	TabelKorespondensi         = "T_CORRESPONDENCELIST"
	TabelPenerimaKorespondensi = "T_RECIPIENTLIST"

	sequenceKorespondensi         = "SEQ_T_CORRESPONDENCELIST"
	sequencePenerimaKorespondensi = "SEQ_T_RECIPIENTLIST"

	// statusDraf / statusTerkirim - `.IsEditable`.
	statusDraf     = "1"
	statusTerkirim = "0"
)

var (
	// ErrKorespondensiTidakAda - nomor bukan milik kasus. 404.
	ErrKorespondensiTidakAda = errors.New("repository: korespondensi tidak ada di kasus ini")
	// ErrKorespondensiTerkirim - korespondensi sudah terkirim (terkunci). 409.
	ErrKorespondensiTerkirim = errors.New("repository: korespondensi sudah terkirim dan tidak dapat diubah")
)

// PenyimpanKorespondensi - korespondensi satu kasus.
type PenyimpanKorespondensi interface {
	// DaftarKorespondensi - urut SEQ_NO; ErrKasusTidakAda bila case FAC tidak ada.
	DaftarKorespondensi(ctx context.Context, id string) ([]models.Korespondensi, error)
	// TambahKorespondensi - draf baru di akhir (SEQ_NO terbesar + 1); mengembalikan nomornya.
	TambahKorespondensi(ctx context.Context, tx *db.Tx, id string, k models.Korespondensi, pembuat string) (string, error)
	// UbahKorespondensi - subject / message / penerima draf (penerima diganti utuh).
	UbahKorespondensi(ctx context.Context, tx *db.Tx, id, no string, isi models.IsiKorespondensi) error
	// HapusKorespondensi - draf beserta penerimanya.
	HapusKorespondensi(ctx context.Context, tx *db.Tx, id, no string) error
	// TandaiTerkirim - draf -> terkirim (`.IsEditable` 0 untuk korespondensi dan penerimanya), SENT_AT = waktu.
	TandaiTerkirim(ctx context.Context, tx *db.Tx, id, no, waktu string) error
}

// PembacaAkunLogin - akun login M_LOGIN_GO untuk korespondensi.
type PembacaAkunLogin interface {
	// AkunLogin - NAME / EMAIL akun; ada=false bila akun tidak ada.
	AkunLogin(ctx context.Context, akunID string) (nama, email string, ada bool, err error)
	// SaranPenerima - akun aktif ber-EMAIL, NAME / EMAIL mengandung kata (tidak peka huruf); id LOGIN_ID, label NAME,
	// ekstra EMAIL; urut NAME; batas BatasSaranAkumulasi.
	SaranPenerima(ctx context.Context, kata string) ([]models.SaranAkumulasi, error)
}

// KorespondensiOracle - PenyimpanKorespondensi dan PembacaAkunLogin atas Oracle.
type KorespondensiOracle struct{ db *db.DB }

// NewKorespondensiOracle merakit penyimpan korespondensi.
func NewKorespondensiOracle(d *db.DB) *KorespondensiOracle { return &KorespondensiOracle{db: d} }

type tabelKorespondensi struct{ work, general, kor, pen, login string }

func (r *KorespondensiOracle) tabel() (tabelKorespondensi, error) {
	var t tabelKorespondensi
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelGeneralPolis, &t.general}, {TabelKorespondensi, &t.kor},
		{TabelPenerimaKorespondensi, &t.pen}, {TabelLoginGo, &t.login}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

// ---- SQL ---------------------------------------------------------------------------------------------------------

func sqlKasusAdaKorespondensi(work string) string {
	return "SELECT COUNT(*) FROM " + work + " WHERE ID = :1 AND LINI = :2"
}

func sqlDaftarKorespondensi(t tabelKorespondensi) string {
	return "SELECT TO_CHAR(k.ID), k.TANGGAL, k.SENDER, k.SENDER_MAIL, k.SUBJECT, k.MESSAGE, k.IS_EDITABLE, k.SENT_AT FROM " + t.kor +
		" k WHERE k.PARENT_ID = :1 ORDER BY k.SEQ_NO, k.ID"
}

func sqlDaftarPenerima(t tabelKorespondensi) string {
	return "SELECT TO_CHAR(p.PARENT_ID), p.RECIPIENT, p.RECIPIENT_MAIL FROM " + t.pen + " p WHERE p.PARENT_ID IN (SELECT k.ID FROM " +
		t.kor + " k WHERE k.PARENT_ID = :1) ORDER BY p.PARENT_ID, p.SEQ_NO, p.ID"
}

func sqlSisipKorespondensi(t tabelKorespondensi) string {
	return "INSERT INTO " + t.kor + " (ID, PARENT_ID, SEQ_NO, ROW_UID, TANGGAL, SENDER, SENDER_MAIL, SUBJECT, MESSAGE, IS_EDITABLE," +
		" CREATE_OP, TGL_CREATE) SELECT :1, :2, NVL(MAX(SEQ_NO), 0) + 1, :3, :4, :5, :6, :7, :8, '" + statusDraf + "', :9, SYSDATE FROM " +
		t.kor + " WHERE PARENT_ID = :10"
}

func sqlUbahKorespondensi(kor string) string {
	return "UPDATE " + kor + " SET SUBJECT = :1, MESSAGE = :2 WHERE ID = :3 AND PARENT_ID = :4 AND IS_EDITABLE = '" + statusDraf + "'"
}

func sqlStatusKorespondensi(kor string) string {
	return "SELECT IS_EDITABLE FROM " + kor + " WHERE ID = :1 AND PARENT_ID = :2"
}

func sqlHapusPenerima(pen string) string { return "DELETE FROM " + pen + " WHERE PARENT_ID = :1" }

func sqlSisipPenerima(pen string) string {
	return "INSERT INTO " + pen + " (ID, PARENT_ID, SEQ_NO, ROW_UID, RECIPIENT, RECIPIENT_MAIL, IS_EDITABLE) VALUES (:1, :2, :3, :4, :5, :6, '" +
		statusDraf + "')"
}

func sqlHapusKorespondensi(kor string) string {
	return "DELETE FROM " + kor + " WHERE ID = :1 AND PARENT_ID = :2 AND IS_EDITABLE = '" + statusDraf + "'"
}

func sqlTandaiTerkirim(kor string) string {
	return "UPDATE " + kor + " SET IS_EDITABLE = '" + statusTerkirim + "', SENT_AT = :1 WHERE ID = :2 AND PARENT_ID = :3 AND IS_EDITABLE = '" +
		statusDraf + "'"
}

func sqlTandaiPenerimaTerkirim(pen string) string {
	return "UPDATE " + pen + " SET IS_EDITABLE = '" + statusTerkirim + "' WHERE PARENT_ID = :1"
}

func sqlAkunLogin(login string) string {
	return "SELECT NAME, EMAIL FROM " + login + " WHERE LOGIN_ID = :1"
}

// sqlSaranPenerima - akun aktif (IS_ACTIVE '1', inti 902) ber-EMAIL; kata opsional.
func sqlSaranPenerima(login string, denganKata bool) string {
	syarat, batas := "", ":1"
	if denganKata {
		syarat, batas = ` AND (UPPER(NAME) LIKE :1 ESCAPE '\' OR UPPER(EMAIL) LIKE :2 ESCAPE '\')`, ":3"
	}
	return "SELECT LOGIN_ID, NAME, EMAIL FROM " + login + " WHERE IS_ACTIVE = '1' AND EMAIL IS NOT NULL" + syarat +
		" ORDER BY UPPER(NAME), LOGIN_ID FETCH FIRST " + batas + " ROWS ONLY"
}

// ---- baca ---------------------------------------------------------------------------------------------------------

// DaftarKorespondensi - lihat PenyimpanKorespondensi.
func (r *KorespondensiOracle) DaftarKorespondensi(ctx context.Context, id string) ([]models.Korespondensi, error) {
	t, err := r.tabel()
	if err != nil {
		return nil, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, sqlKasusAdaKorespondensi(t.work), id, LiniFacIn).Scan(&n); err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelWorkPolis, err)
	}
	if n == 0 {
		return nil, ErrKasusTidakAda
	}
	hasil := []models.Korespondensi{}
	indeks := map[string]int{}
	baris, err := r.db.QueryContext(ctx, sqlDaftarKorespondensi(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelKorespondensi, err)
	}
	defer baris.Close()
	for baris.Next() {
		var no, tgl, pengirim, surel, subjek, pesan, status, kirim sql.NullString
		if err := baris.Scan(&no, &tgl, &pengirim, &surel, &subjek, &pesan, &status, &kirim); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelKorespondensi, err)
		}
		indeks[no.String] = len(hasil)
		hasil = append(hasil, models.Korespondensi{No: no.String, Tanggal: tgl.String, Sender: pengirim.String, SenderMail: surel.String,
			Subject: subjek.String, Message: pesan.String, Penerima: []models.PenerimaKorespondensi{},
			Terkirim: status.String == statusTerkirim, DikirimPada: kirim.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelKorespondensi, err)
	}
	pen, err := r.db.QueryContext(ctx, sqlDaftarPenerima(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelPenerimaKorespondensi, err)
	}
	defer pen.Close()
	for pen.Next() {
		var induk, nama, surel sql.NullString
		if err := pen.Scan(&induk, &nama, &surel); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelPenerimaKorespondensi, err)
		}
		if i, ada := indeks[induk.String]; ada {
			hasil[i].Penerima = append(hasil[i].Penerima, models.PenerimaKorespondensi{Recipient: nama.String, RecipientMail: surel.String})
		}
	}
	if err := pen.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelPenerimaKorespondensi, err)
	}
	return hasil, nil
}

// statusSesudahGagal - UPDATE / DELETE draf tidak mengenai baris: tidak ada -> ErrKorespondensiTidakAda, ada (berarti
// terkirim) -> ErrKorespondensiTerkirim.
func (r *KorespondensiOracle) statusSesudahGagal(ctx context.Context, tx *db.Tx, kor, id, no string) error {
	var status sql.NullString
	err := tx.QueryRowContext(ctx, sqlStatusKorespondensi(kor), no, id).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrKorespondensiTidakAda
	}
	if err != nil {
		return fmt.Errorf("repository: baca %s: %w", TabelKorespondensi, err)
	}
	return ErrKorespondensiTerkirim
}

// ---- tulis --------------------------------------------------------------------------------------------------------

// sentuhKasus - case FAC dikunci (UPDATE TGL_UPDATE) dan T_GENERAL_POLIS dipastikan ada.
func (r *KorespondensiOracle) sentuhKasus(ctx context.Context, tx *db.Tx, t tabelKorespondensi, id string) error {
	hasil, err := jalankan(ctx, tx, sqlSentuhCase(t.work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrKasusTidakAda
	}
	_, err = jalankan(ctx, tx, sqlPastikanGeneral(t.general), "memastikan "+TabelGeneralPolis, id, id)
	return err
}

// TambahKorespondensi - lihat PenyimpanKorespondensi. Kunci baris case (sentuh) menyerialkan SEQ_NO Add serentak.
func (r *KorespondensiOracle) TambahKorespondensi(ctx context.Context, tx *db.Tx, id string, k models.Korespondensi, pembuat string) (string, error) {
	t, err := r.tabel()
	if err != nil {
		return "", err
	}
	if err := r.sentuhKasus(ctx, tx, t, id); err != nil {
		return "", err
	}
	no, err := r.db.NomorBerikut(ctx, tx, sequenceKorespondensi)
	if err != nil {
		return "", err
	}
	uid, err := uidAcak()
	if err != nil {
		return "", err
	}
	_, err = jalankan(ctx, tx, sqlSisipKorespondensi(t), "menyisipkan "+TabelKorespondensi, no, id, uid, k.Tanggal,
		db.KosongJadiNil(k.Sender), db.KosongJadiNil(k.SenderMail), db.KosongJadiNil(k.Subject), db.KosongJadiNil(k.Message),
		db.KosongJadiNil(pembuat), id)
	return no, err
}

// UbahKorespondensi - lihat PenyimpanKorespondensi.
func (r *KorespondensiOracle) UbahKorespondensi(ctx context.Context, tx *db.Tx, id, no string, isi models.IsiKorespondensi) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	if err := r.sentuhKasus(ctx, tx, t, id); err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlUbahKorespondensi(t.kor), "mengubah "+TabelKorespondensi, db.KosongJadiNil(isi.Subject),
		db.KosongJadiNil(isi.Message), no, id)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return r.statusSesudahGagal(ctx, tx, t.kor, id, no)
	}
	if _, err := jalankan(ctx, tx, sqlHapusPenerima(t.pen), "menghapus "+TabelPenerimaKorespondensi, no); err != nil {
		return err
	}
	for i, p := range isi.Penerima {
		idP, err := r.db.NomorBerikut(ctx, tx, sequencePenerimaKorespondensi)
		if err != nil {
			return err
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		if _, err := jalankan(ctx, tx, sqlSisipPenerima(t.pen), "menyisipkan "+TabelPenerimaKorespondensi, idP, no, i+1, uid,
			db.KosongJadiNil(p.Recipient), db.KosongJadiNil(p.RecipientMail)); err != nil {
			return err
		}
	}
	return nil
}

// HapusKorespondensi - lihat PenyimpanKorespondensi.
func (r *KorespondensiOracle) HapusKorespondensi(ctx context.Context, tx *db.Tx, id, no string) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	if err := r.sentuhKasus(ctx, tx, t, id); err != nil {
		return err
	}
	var status sql.NullString
	err = tx.QueryRowContext(ctx, sqlStatusKorespondensi(t.kor), no, id).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrKorespondensiTidakAda
	}
	if err != nil {
		return fmt.Errorf("repository: baca %s: %w", TabelKorespondensi, err)
	}
	if status.String != statusDraf {
		return ErrKorespondensiTerkirim
	}
	if _, err := jalankan(ctx, tx, sqlHapusPenerima(t.pen), "menghapus "+TabelPenerimaKorespondensi, no); err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlHapusKorespondensi(t.kor), "menghapus "+TabelKorespondensi, no, id)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return r.statusSesudahGagal(ctx, tx, t.kor, id, no)
	}
	return nil
}

// TandaiTerkirim - lihat PenyimpanKorespondensi.
func (r *KorespondensiOracle) TandaiTerkirim(ctx context.Context, tx *db.Tx, id, no, waktu string) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlTandaiTerkirim(t.kor), "menandai terkirim "+TabelKorespondensi, waktu, no, id)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return r.statusSesudahGagal(ctx, tx, t.kor, id, no)
	}
	_, err = jalankan(ctx, tx, sqlTandaiPenerimaTerkirim(t.pen), "menandai terkirim "+TabelPenerimaKorespondensi, no)
	return err
}

// AkunLogin - lihat PembacaAkunLogin.
func (r *KorespondensiOracle) AkunLogin(ctx context.Context, akunID string) (string, string, bool, error) {
	t, err := r.db.Qualify(TabelLoginGo)
	if err != nil {
		return "", "", false, err
	}
	var nama, surel sql.NullString
	err = r.db.QueryRowContext(ctx, sqlAkunLogin(t), akunID).Scan(&nama, &surel)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", false, nil
	}
	if err != nil {
		return "", "", false, fmt.Errorf("repository: baca %s: %w", TabelLoginGo, err)
	}
	return nama.String, surel.String, true, nil
}

// SaranPenerima - lihat PembacaAkunLogin.
func (r *KorespondensiOracle) SaranPenerima(ctx context.Context, kata string) ([]models.SaranAkumulasi, error) {
	t, err := r.db.Qualify(TabelLoginGo)
	if err != nil {
		return nil, err
	}
	q, arg := sqlSaranPenerima(t, false), []any{BatasSaranAkumulasi}
	if pola := PolaCari(kata); pola != "" {
		q, arg = sqlSaranPenerima(t, true), []any{pola, pola, BatasSaranAkumulasi}
	}
	baris, err := r.db.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelLoginGo, err)
	}
	defer baris.Close()
	hasil := []models.SaranAkumulasi{}
	for baris.Next() {
		var id, nama, surel sql.NullString
		if err := baris.Scan(&id, &nama, &surel); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelLoginGo, err)
		}
		hasil = append(hasil, models.SaranAkumulasi{ID: id.String, Label: nama.String, Ekstra: surel.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelLoginGo, err)
	}
	return hasil, nil
}
