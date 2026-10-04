package repository

// Isian penawaran polis - tiket 01 bagian 3 (form Input Offer).
//
// Untuk apa berkas ini: membaca dan menulis kolom header `T_PREMIUM_LIST` yang
// layar `Section/InputOfferLife.xml` isi, beserta riwayatnya `T_VIEW_SUGGEST`
// (`Activity/AddHistorySuggest.xml`). Penggantinya `INSERTJSONOFFERLIFE` -
// yang TIDAK dipanggil (spec §12): penawaran kini kolom bernama.
//
// ⛔ Setiap query menyebut skemanya lewat `Qualify` (ADR-U-0033), nol
// `COMMIT` (ADR-U-0029); tulisan hanya di transaksi PEMANGGIL.
//
// Dibaca sesudah: polis_kasus.go (yang melahirkan header kosongnya).

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/premiumlistlife/backend/models"
)

// ErrHeaderPolisTidakAda - baris `T_PREMIUM_LIST` polis itu tidak ada.
var ErrHeaderPolisTidakAda = errors.New("repository: header polis tidak ada")

// Penawaran membaca dan menulis isian penawaran.
type Penawaran struct{ db *db.DB }

// NewPenawaran menyusunnya.
func NewPenawaran(db *db.DB) *Penawaran { return &Penawaran{db: db} }

// IsianTersimpan adalah isian penawaran yang sudah ada di header.
type IsianTersimpan struct {
	NoOffer                string
	CedingCo               string
	CedingCoName           string
	PolicyHolder           string
	PolicyHolderName       string
	TypeCeding             string
	TypeCedingName         string
	BusinessCode           string
	BusinessName           string
	DateReceived           *time.Time
	Description            string
	BatasUsiaPeserta       *int
	PeriodePertanggungan   string
	SumInsured             *apd.Decimal
	TanggalPenawaran       *time.Time
	TanggalRespon          *time.Time
	TanggalKonfirmasi      *time.Time
	TBC                    *int
	TanggalTBC             *time.Time
	StatusUpdate           string
	KeteranganMarketing    string
	QQName                 string
	JenisUsaha             string
	KetentuanUnderwriting  string
	TanggalKonfirmasiBalik *time.Time
	TanggalRealisasi       *time.Time
	TanggalBind            *time.Time
	StatusFinal            string
	JenisAsuransi          string
	// Status - `STATUS_PENAWARAN` baris utama: status TERAKHIR disimpan.
	Status string
}

// Jenis kolom penawaran - menentukan bungkus SELECT dan pengurai.
const (
	kolomPenawaranTeks = iota
	kolomPenawaranTanggal
	kolomPenawaranAngka
)

// kolomBacaPenawaran - urutan SELECT `sqlBacaPenawaran` DAN urutan pengurai
// `Baca`. SATU daftar untuk keduanya: dua daftar akan berselisih, dan
// selisihnya memasukkan nilai satu kolom ke medan lain tanpa galat apa pun.
var kolomBacaPenawaran = []struct {
	nama  string
	jenis int
}{
	{"NO_OFFER", kolomPenawaranTeks}, {"CEDING_CO", kolomPenawaranTeks},
	{"CEDING_CO_NAME", kolomPenawaranTeks}, {"POLICY_HOLDER", kolomPenawaranTeks},
	{"POLICY_HOLDER_NAME", kolomPenawaranTeks}, {"TYPE_CEDING", kolomPenawaranTeks},
	{"TYPE_CEDING_NAME", kolomPenawaranTeks}, {"BUSINESS_CODE", kolomPenawaranTeks},
	{"BUSINESS_NAME", kolomPenawaranTeks}, {"DATE_RECEIVED", kolomPenawaranTanggal},
	{"DESCRIPTION", kolomPenawaranTeks}, {"BATAS_USIA_PESERTA", kolomPenawaranAngka},
	{"PERIODE_PERTANGGUNGAN", kolomPenawaranTeks}, {"SUM_INSURED", kolomPenawaranAngka},
	{"TANGGAL_PENAWARAN", kolomPenawaranTanggal}, {"TANGGAL_RESPON", kolomPenawaranTanggal},
	{"TANGGAL_KONFIRMASI", kolomPenawaranTanggal}, {"TBC", kolomPenawaranAngka},
	{"TANGGAL_TBC", kolomPenawaranTanggal}, {"STATUS_UPDATE", kolomPenawaranTeks},
	{"KETERANGAN_MARKETING", kolomPenawaranTeks},
	{"QQ_NAME", kolomPenawaranTeks}, {"JENIS_USAHA", kolomPenawaranTeks},
	{"KETENTUAN_UNDERWRITING", kolomPenawaranTeks}, {"TANGGAL_KONFIRMASI_BALIK", kolomPenawaranTanggal},
	{"TANGGAL_REALISASI", kolomPenawaranTanggal}, {"TANGGAL_BIND", kolomPenawaranTanggal},
	{"STATUS_FINAL", kolomPenawaranTeks},
	{"JENIS_ASURANSI", kolomPenawaranTeks},
	{"STATUS_PENAWARAN", kolomPenawaranTeks},
}

// sqlBacaPenawaran merakit pembacaan header.
//
// ⚠️ Tanggal keluar sebagai TEKS berpola tetap (`db.FmtTanggalOracle`) dan
// angka sebagai TEKS (`db.FmtDesimal`) - nol `float64` di jalur uang.
func sqlBacaPenawaran(polis string) string {
	bagian := make([]string, len(kolomBacaPenawaran))
	for i, k := range kolomBacaPenawaran {
		medan := "p." + k.nama
		switch k.jenis {
		case kolomPenawaranTanggal:
			bagian[i] = fmt.Sprintf(db.FmtTanggalOracle, medan)
		case kolomPenawaranAngka:
			bagian[i] = fmt.Sprintf(db.FmtDesimal, medan)
		default:
			bagian[i] = medan
		}
	}
	return fmt.Sprintf(`SELECT %s FROM %s p WHERE p.ID = :1`, strings.Join(bagian, ", "), polis)
}

// Baca membaca isian penawaran satu polis.
func (r *Penawaran) Baca(ctx context.Context, id string) (IsianTersimpan, error) {
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return IsianTersimpan{}, err
	}
	q := sqlBacaPenawaran(polis)
	if err := db.PeriksaSQL(q); err != nil {
		return IsianTersimpan{}, err
	}
	n := make([]sql.NullString, len(kolomBacaPenawaran))
	tujuan := make([]any, len(n))
	for i := range n {
		tujuan[i] = &n[i]
	}
	err = r.db.QueryRowContext(ctx, q, id).Scan(tujuan...)
	if errors.Is(err, sql.ErrNoRows) {
		return IsianTersimpan{}, ErrHeaderPolisTidakAda
	}
	if err != nil {
		return IsianTersimpan{}, fmt.Errorf("repository: membaca penawaran polis: %w", err)
	}
	return uraiPenawaran(n)
}

// uraiPenawaran mengurai satu baris `kolomBacaPenawaran` menjadi isian.
func uraiPenawaran(n []sql.NullString) (IsianTersimpan, error) {
	nilai := map[string]sql.NullString{}
	for i, k := range kolomBacaPenawaran {
		nilai[k.nama] = n[i]
	}
	teks := func(k string) string { return nilai[k].String }
	h := IsianTersimpan{
		NoOffer: teks("NO_OFFER"), CedingCo: teks("CEDING_CO"), CedingCoName: teks("CEDING_CO_NAME"),
		PolicyHolder: teks("POLICY_HOLDER"), PolicyHolderName: teks("POLICY_HOLDER_NAME"),
		TypeCeding: teks("TYPE_CEDING"), TypeCedingName: teks("TYPE_CEDING_NAME"),
		BusinessCode: teks("BUSINESS_CODE"), BusinessName: teks("BUSINESS_NAME"),
		Description: teks("DESCRIPTION"), PeriodePertanggungan: teks("PERIODE_PERTANGGUNGAN"),
		StatusUpdate: teks("STATUS_UPDATE"), KeteranganMarketing: teks("KETERANGAN_MARKETING"),
		QQName: teks("QQ_NAME"), JenisUsaha: teks("JENIS_USAHA"),
		KetentuanUnderwriting: teks("KETENTUAN_UNDERWRITING"), StatusFinal: teks("STATUS_FINAL"),
		JenisAsuransi: teks("JENIS_ASURANSI"),
		Status:        teks("STATUS_PENAWARAN"),
	}
	var err error
	for _, t := range []struct {
		kolom  string
		tujuan **time.Time
	}{
		{"DATE_RECEIVED", &h.DateReceived}, {"TANGGAL_PENAWARAN", &h.TanggalPenawaran},
		{"TANGGAL_RESPON", &h.TanggalRespon}, {"TANGGAL_KONFIRMASI", &h.TanggalKonfirmasi},
		{"TANGGAL_TBC", &h.TanggalTBC},
		{"TANGGAL_KONFIRMASI_BALIK", &h.TanggalKonfirmasiBalik},
		{"TANGGAL_REALISASI", &h.TanggalRealisasi}, {"TANGGAL_BIND", &h.TanggalBind},
	} {
		if *t.tujuan, err = uraiTanggalOracle(nilai[t.kolom], t.kolom); err != nil {
			return IsianTersimpan{}, err
		}
	}
	for _, b := range []struct {
		kolom  string
		tujuan **int
	}{{"BATAS_USIA_PESERTA", &h.BatasUsiaPeserta}, {"TBC", &h.TBC}} {
		if *b.tujuan, err = uraiBulat(nilai[b.kolom], b.kolom); err != nil {
			return IsianTersimpan{}, err
		}
	}
	if s := nilai["SUM_INSURED"]; s.Valid && strings.TrimSpace(s.String) != "" {
		if h.SumInsured, err = utils.ParseDecimal(strings.TrimSpace(s.String)); err != nil {
			return IsianTersimpan{}, fmt.Errorf("repository: kolom SUM_INSURED bernilai %q: %w", s.String, err)
		}
	}
	return h, nil
}

// sqlSimpanPenawaran menimpa kolom isian penawaran.
//
// ⛔ HANYA kolom yang layar itu isi. Kolom milik tahap lain (TYPE,
// PRODUCT_NAME, SOB, ...) tidak disentuh: penawaran yang disimpan ulang tidak
// boleh menghapus isian Premium List Detail. `NO_OFFER` pun tidak: ia nomor
// yang lahir di tempat lain, bukan ketikan.
func sqlSimpanPenawaran(polis string) string {
	return fmt.Sprintf(`UPDATE %s
	    SET CEDING_CO = :1, CEDING_CO_NAME = :2, POLICY_HOLDER = :3, POLICY_HOLDER_NAME = :4,
	        TYPE_CEDING = :5, TYPE_CEDING_NAME = :6, BUSINESS_CODE = :7, BUSINESS_NAME = :8,
	        DATE_RECEIVED = :9, DESCRIPTION = :10,
	        BATAS_USIA_PESERTA = :11, PERIODE_PERTANGGUNGAN = :12, SUM_INSURED = :13,
	        TANGGAL_PENAWARAN = :14, TANGGAL_RESPON = :15, TANGGAL_KONFIRMASI = :16,
	        TBC = :17, TANGGAL_TBC = :18, STATUS_UPDATE = :19, KETERANGAN_MARKETING = :20,
	        QQ_NAME = :21, JENIS_USAHA = :22, KETENTUAN_UNDERWRITING = :23,
	        TANGGAL_KONFIRMASI_BALIK = :24, TANGGAL_REALISASI = :25, TANGGAL_BIND = :26,
	        STATUS_FINAL = :27, JENIS_ASURANSI = :28, STATUS_PENAWARAN = :29
	  WHERE ID = :30`, polis)
}

func argSimpanPenawaran(id string, p models.PenawaranTersimpan) []any {
	tanggal := func(t *time.Time) any {
		if t == nil {
			return nil
		}
		return *t
	}
	bulat := func(v *int) any {
		if v == nil {
			return nil
		}
		return *v
	}
	// ⛔ Uang dikirim sebagai TEKS desimal - pola `nilaiSisipRekap`.
	var sumInsured any
	if p.SumInsured != nil {
		sumInsured = utils.FormatDecimal(p.SumInsured)
	}
	return []any{
		db.KosongJadiNil(p.CedingCo), db.KosongJadiNil(p.CedingCoName),
		db.KosongJadiNil(p.PolicyHolder), db.KosongJadiNil(p.PolicyHolderName),
		db.KosongJadiNil(p.TypeCeding), db.KosongJadiNil(p.TypeCedingName),
		db.KosongJadiNil(p.BusinessCode), db.KosongJadiNil(p.BusinessName),
		tanggal(p.DateReceived), db.KosongJadiNil(p.Description),
		bulat(p.BatasUsiaPeserta), db.KosongJadiNil(p.PeriodePertanggungan), sumInsured,
		tanggal(p.TanggalPenawaran), tanggal(p.TanggalRespon), tanggal(p.TanggalKonfirmasi),
		bulat(p.TBC), tanggal(p.TanggalTBC), db.KosongJadiNil(p.StatusUpdate),
		db.KosongJadiNil(p.KeteranganMarketing),
		db.KosongJadiNil(p.QQName), db.KosongJadiNil(p.JenisUsaha), db.KosongJadiNil(p.KetentuanUnderwriting),
		tanggal(p.TanggalKonfirmasiBalik), tanggal(p.TanggalRealisasi), tanggal(p.TanggalBind),
		db.KosongJadiNil(p.StatusFinal), db.KosongJadiNil(p.JenisAsuransi),
		db.KosongJadiNil(p.Status), id,
	}
}

// uraiBulat membaca teks `db.FmtDesimal` kolom NUMBER(5); NULL tetap nil.
func uraiBulat(v sql.NullString, kolom string) (*int, error) {
	if !v.Valid || strings.TrimSpace(v.String) == "" {
		return nil, nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(v.String))
	if err != nil {
		return nil, fmt.Errorf("repository: kolom %s bernilai %q: %w", kolom, v.String, err)
	}
	return &n, nil
}

// Simpan menulis isian penawaran ke header - di transaksi pemanggil.
//
// ⚠️ UPDATE ini juga MENGUNCI baris header sampai commit, dan `SisipSuggest`
// bergantung padanya: dua simpanan serentak atas polis yang sama berurutan di
// sini, sehingga nomor riwayat berikutnya tidak dibaca ganda.
func (r *Penawaran) Simpan(ctx context.Context, tx *db.Tx, id string, p models.PenawaranTersimpan) error {
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return err
	}
	q := sqlSimpanPenawaran(polis)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, argSimpanPenawaran(id, p)...)
	if err != nil {
		return fmt.Errorf("repository: menyimpan penawaran polis: %w", err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrHeaderPolisTidakAda
	}
	return db.PastikanSatuBaris(hasil, "penawaran polis")
}

// PengenalSuggest menyusun `T_VIEW_SUGGEST.ID` - 32 heksa, deterministik.
//
// ⛔ Pola `PengenalPesertaUnggah`: nol sequence untuk keluarga tabel ini.
// Riwayat hanya bertambah, dan `NO` unik per polis, jadi pasangan (polis, NO)
// memberi pengenal yang unik dan dapat dibandingkan.
func PengenalSuggest(polisID string, no int) string {
	sum := md5.Sum([]byte(polisID + "\x00suggest\x00" + strconv.Itoa(no)))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// sqlNomorSuggestBerikut - `pxListSubscript` baris yang baru ditambahkan.
func sqlNomorSuggestBerikut(suggest string) string {
	return fmt.Sprintf(`SELECT NVL(MAX(NO), 0) + 1 FROM %s WHERE PREMIUM_LIST_ID = :1`, suggest)
}

func sqlSisipSuggest(suggest string) string {
	return fmt.Sprintf(`INSERT INTO %s
	    (ID, PREMIUM_LIST_ID, NO, DATE_SUGGEST, PIC_SUGGEST, IS_CEDING_CONFIRM, COMMENT_SUGGEST, INITIAL_SUGGEST)
	  VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, suggest)
}

// SisipSuggest menambah satu baris riwayat dan mengembalikan nomornya.
func (r *Penawaran) SisipSuggest(ctx context.Context, tx *db.Tx, polisID string, b models.BarisSuggest) (int, error) {
	suggest, err := r.db.Qualify("T_VIEW_SUGGEST")
	if err != nil {
		return 0, err
	}
	qNo := sqlNomorSuggestBerikut(suggest)
	if err := db.PeriksaSQL(qNo); err != nil {
		return 0, err
	}
	var noTeks sql.NullString
	if err := tx.QueryRowContext(ctx, qNo, polisID).Scan(&noTeks); err != nil {
		return 0, fmt.Errorf("repository: membaca nomor riwayat penawaran: %w", err)
	}
	no, err := strconv.Atoi(strings.TrimSpace(noTeks.String))
	if err != nil {
		return 0, fmt.Errorf("repository: nomor riwayat penawaran %q: %w", noTeks.String, err)
	}
	q := sqlSisipSuggest(suggest)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	hasil, err := tx.ExecContext(ctx, q, PengenalSuggest(polisID, no), polisID, no, b.DateSuggest,
		db.KosongJadiNil(b.PICSuggest), db.KosongJadiNil(b.IsCedingConfirm),
		db.KosongJadiNil(b.CommentSuggest), db.KosongJadiNil(b.InitialSuggest))
	if err != nil {
		return 0, fmt.Errorf("repository: menyisipkan riwayat penawaran: %w", err)
	}
	return no, db.PastikanSatuBaris(hasil, "riwayat penawaran")
}

// sqlRiwayatSuggest - terbaru dahulu, padanan `Obj-Sort .No Descending`
// (AddHistorySuggest langkah 3).
func sqlRiwayatSuggest(suggest string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(s.NO), `+db.FmtTanggalOracle+`, s.PIC_SUGGEST,
	        s.IS_CEDING_CONFIRM, s.COMMENT_SUGGEST, s.INITIAL_SUGGEST
	   FROM %s s WHERE s.PREMIUM_LIST_ID = :1
	  ORDER BY s.NO DESC, s.ID`, "s.DATE_SUGGEST", suggest)
}

// Riwayat membaca seluruh riwayat penawaran satu polis.
func (r *Penawaran) Riwayat(ctx context.Context, polisID string) ([]models.BarisSuggest, error) {
	suggest, err := r.db.Qualify("T_VIEW_SUGGEST")
	if err != nil {
		return nil, err
	}
	q := sqlRiwayatSuggest(suggest)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, polisID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca riwayat penawaran: %w", err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []models.BarisSuggest{}
	for rows.Next() {
		var n [6]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5]); err != nil {
			return nil, fmt.Errorf("repository: memindai riwayat penawaran: %w", err)
		}
		b := models.BarisSuggest{PICSuggest: n[2].String, IsCedingConfirm: n[3].String,
			CommentSuggest: n[4].String, InitialSuggest: n[5].String}
		if n[0].Valid {
			if b.No, err = strconv.Atoi(strings.TrimSpace(n[0].String)); err != nil {
				return nil, fmt.Errorf("repository: kolom NO riwayat %q: %w", n[0].String, err)
			}
		}
		t, err := uraiTanggalOracle(n[1], "DATE_SUGGEST")
		if err != nil {
			return nil, err
		}
		if t != nil {
			b.DateSuggest = *t
		}
		hasil = append(hasil, b)
	}
	return hasil, rows.Err()
}

// uraiTanggalOracle membaca teks `db.FmtTanggalOracle`; NULL tetap nil.
//
// ⛔ Waktu dibaca sebagai waktu SETEMPAT tanpa zona - bentuk yang Oracle
// simpan di kolom DATE. Mengarang zona di sini menggeser tanggal di dekat
// tengah malam.
func uraiTanggalOracle(v sql.NullString, kolom string) (*time.Time, error) {
	if !v.Valid || strings.TrimSpace(v.String) == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(v.String))
	if err != nil {
		return nil, fmt.Errorf("repository: kolom %s bernilai %q: %w", kolom, v.String, err)
	}
	return &t, nil
}
