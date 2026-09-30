package repository

// Penyimpanan rekap uang dan peserta warisan - tiket 05a bagian 2.
//
// Untuk apa berkas ini: tiga pekerjaan yang SELALU berjalan di dalam satu
// transaksi yang dipegang layanan (ADR-U-0015, ADR-U-0029):
//
//  1. membaca baris uang peserta satu polis - bahan rumus murni
//     `models.RekapPerMataUang`;
//  2. mengganti rekap `T_PREMIUM_LIST_SUMMARY` (hapus lalu sisip);
//  3. pl2 - menyalin peserta ke tabel warisan `M_LIFE_PREMIUM_DETAIL`, yang
//     dibaca Claim Life (`GET /api/peserta-life`), dan membaca kepala rekap
//     warisan `M_LIFE_PREMIUM_SUMMARY` (PL-09; penulisnya polis_warisan.go).
//
// ⛔ PEMBACA TIDAK MENJUMLAH. Pilihan brief giliran 10: SUM/GROUP BY di SQL,
// ATAU baca lalu hitung dengan rumus murni. Yang dipilih yang KEDUA. Cabang
// `Type` (empat `WHEN` yang saling meniadakan) dan tiga keanehan tanda
// `BALANCE` hidup di `models`, dan sudah dikunci uji literal 975/978/897/1789.
// Menulis ulang cabang itu sebagai `CASE` di SQL berarti DUA salinan rumus
// uang, dan salinan kedua tidak diuji oleh satu pun literal itu.
//
// ⛔ Uang tiba sebagai TEKS lewat `TO_CHAR` ber-NLS dan diurai `uraiDesimal`
// - konversi teks ↔ desimal hanya di lapisan ini (AC 15 spec, ADR-U-0003).
//
// ⛔ NOL `COMMIT` di teks SQL mana pun di sini (ADR-U-0029). Rule warisan
// `SaveMasterLPDet` menulis `COMMIT;` di baris 252 - itu yang membuat nomor
// dan peserta di Pega TIDAK atomik. Kami tidak menirunya.
//
// Dibaca sesudah: models/polis_summary.go (rumusnya), polis_nomor.go.

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/premiumlistlife/backend/models"
)

// SummaryPolis membaca peserta dan menyimpan rekap satu polis.
type SummaryPolis struct{ db *db.DB }

// NewSummaryPolis menyusunnya.
func NewSummaryPolis(db *db.DB) *SummaryPolis { return &SummaryPolis{db: db} }

// sqlBarisUangPolis merakit pembacaan baris uang peserta.
//
// ⛔ Kolomnya `models.KolomBacaSummary()`, bukan daftar yang diketik ulang di
// sini. Kolom yang dirumuskan di `models` tetapi lupa dibaca di sini akan
// bernilai NOL tanpa satu galat pun - `BarisUang.ambil` membaca absen sebagai
// nol, dan itu memang perilaku rule aslinya.
func sqlBarisUangPolis(detail string) string {
	kolom := models.KolomBacaSummary()
	pilih := make([]string, 0, len(kolom))
	for _, k := range kolom {
		pilih = append(pilih, fmt.Sprintf(db.FmtDesimal, "d."+k))
	}
	return fmt.Sprintf(`SELECT d.ID, d.CURRENCY, %s
	   FROM %s d WHERE d.PREMIUM_LIST_ID = :1 ORDER BY d.ID`,
		strings.Join(pilih, ", "), detail)
}

// BarisUang membaca seluruh baris uang peserta satu polis.
//
// Mengembalikan baris uang dan mata uangnya, sejajar - bentuk yang diminta
// `models.RekapPerMataUang`.
func (r *SummaryPolis) BarisUang(ctx context.Context, tx *db.Tx, polisID string) (
	[]models.BarisUang, []string, error) {

	detail, err := r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	if err != nil {
		return nil, nil, err
	}
	q := sqlBarisUangPolis(detail)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, nil, err
	}
	rows, err := tx.QueryContext(ctx, q, polisID)
	if err != nil {
		return nil, nil, fmt.Errorf("repository: membaca uang peserta: %w", err)
	}
	defer func() { _ = rows.Close() }()

	kolom := models.KolomBacaSummary()
	var baris []models.BarisUang
	var mataUang []string
	for rows.Next() {
		var id string
		var cur sql.NullString
		sel := make([]sql.NullString, len(kolom))
		tujuan := []any{&id, &cur}
		for i := range sel {
			tujuan = append(tujuan, &sel[i])
		}
		if err := rows.Scan(tujuan...); err != nil {
			return nil, nil, fmt.Errorf("repository: memindai uang peserta: %w", err)
		}
		b := models.BarisUang{}
		for i, k := range kolom {
			v, err := db.UraiDesimal("peserta "+id, k, sel[i])
			if err != nil {
				return nil, nil, err
			}
			// Kosong dibiarkan absen - `BarisUang.ambil` membacanya nol.
			if v != nil {
				b[k] = v
			}
		}
		baris = append(baris, b)
		mataUang = append(mataUang, cur.String)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("repository: membaca uang peserta: %w", err)
	}
	return baris, mataUang, nil
}

// PengenalRekap menyusun pengenal baris rekap satu mata uang.
//
// ⛔ DETERMINISTIK, dengan alasan yang sama dengan `PengenalPesertaUnggah`:
// rekap MENGGANTI isinya, jadi rekap IDR polis yang sama selalu memperoleh
// pengenal yang sama, dan submit ulang tidak menumpuk pengenal yatim.
func PengenalRekap(polisID, mataUang string) string {
	sum := md5.Sum([]byte(polisID + "\x00rekap\x00" + mataUang))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// kolomSisipRekap adalah kolom `T_PREMIUM_LIST_SUMMARY` yang diisi, urut.
//
// ID, PREMIUM_LIST_ID, CURRENCY, tiga turunan, lalu ke-33 kolom jumlah -
// seluruhnya 39, yaitu SELURUH kolom migrasi 055. Penjaga
// `TestKolomRekapSamaDenganMigrasi055` menagihnya dua arah.
func kolomSisipRekap() []string {
	k := []string{"ID", "PREMIUM_LIST_ID", "CURRENCY", "BALANCE", "PREMIUM", "COMMISSION"}
	return append(k, models.KolomJumlahSummary...)
}

// sqlHapusRekap merakit pembersihan rekap satu polis.
func sqlHapusRekap(summary string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE PREMIUM_LIST_ID = :1`, summary)
}

// sqlSisipRekap merakit penyisipan satu baris rekap.
func sqlSisipRekap(summary string) string {
	kolom := kolomSisipRekap()
	penanda := make([]string, len(kolom))
	for i := range kolom {
		penanda[i] = fmt.Sprintf(":%d", i+1)
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,
		summary, strings.Join(kolom, ", "), strings.Join(penanda, ", "))
}

// nilaiSisipRekap menyusun argumen satu baris rekap, urut `kolomSisipRekap`.
//
// ⛔ Uang dikirim sebagai TEKS desimal - bukan float. Rekap yang sudah
// dibulatkan empat angka di `models` tidak boleh dibulatkan ulang oleh driver.
func nilaiSisipRekap(polisID string, r models.RekapMataUang) []any {
	teks := func(v *apd.Decimal) any {
		if v == nil {
			return nil
		}
		return utils.FormatDecimal(v)
	}
	arg := []any{
		PengenalRekap(polisID, r.Currency), polisID, r.Currency,
		teks(r.Balance), teks(r.Premium), teks(r.Commission),
	}
	for _, k := range models.KolomJumlahSummary {
		arg = append(arg, teks(r.Jumlah[k]))
	}
	return arg
}

// ErrRekapKosong - rekap tanpa satu mata uang pun tidak disimpan.
//
// ⛔ Polis tanpa peserta tidak dapat dinomori (polis_nomor.go); bila rekapnya
// tetap kosong sampai di sini, sesuatu di antara keduanya berubah, dan
// menyimpan "tidak ada rekap" di atas rekap lama berarti MENGHAPUS rekap.
var ErrRekapKosong = errors.New(
	"repository: rekap mata uang kosong; tidak ada peserta yang direkap")

// GantiRekap menghapus rekap lama polis lalu menyisipkan yang baru.
//
// Mengembalikan cacah baris terhapus dan tersisip.
func (r *SummaryPolis) GantiRekap(ctx context.Context, tx *db.Tx, polisID string,
	rekap []models.RekapMataUang) (dihapus, disisip int, err error) {

	if len(rekap) == 0 {
		return 0, 0, ErrRekapKosong
	}
	summary, err := r.db.Qualify("T_PREMIUM_LIST_SUMMARY")
	if err != nil {
		return 0, 0, err
	}
	hapus := sqlHapusRekap(summary)
	if err := db.PeriksaSQL(hapus); err != nil {
		return 0, 0, err
	}
	hasil, err := tx.ExecContext(ctx, hapus, polisID)
	if err != nil {
		return 0, 0, fmt.Errorf("repository: menghapus rekap lama: %w", err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return 0, 0, fmt.Errorf("repository: membaca cacah rekap terhapus: %w", err)
	}
	sisip := sqlSisipRekap(summary)
	if err := db.PeriksaSQL(sisip); err != nil {
		return 0, 0, err
	}
	for _, rm := range rekap {
		if _, err := tx.ExecContext(ctx, sisip, nilaiSisipRekap(polisID, rm)...); err != nil {
			// ⛔ Mata uangnya IKUT di galat - dua rekap satu polis hanya
			// dibedakan olehnya.
			return 0, 0, fmt.Errorf("repository: menyisipkan rekap %q: %w", rm.Currency, err)
		}
	}
	return int(n), len(rekap), nil
}

// sqlKepalaSummaryWarisan merakit pembacaan kepala rekap warisan - PL-09.
//
// `COB` <- `pyWorkPage.BusinessName` (`BUSINESS_NAME`, migrasi 051) dan
// `PL_NUMBER_EDM` <- `pyWorkPage.PremiumListSummary.PL_NUMBER_EDM` - sama
// dengan sumber `p.PL_NUMBER_EDM` salinan detail.
func sqlKepalaSummaryWarisan(polis string) string {
	return fmt.Sprintf(`SELECT BUSINESS_NAME, PL_NUMBER_EDM FROM %s WHERE ID = :1`, polis)
}

// KepalaSummaryWarisan membaca kepala rekap warisan satu polis.
//
// `nomorPL` dan `idPega` datang dari pemanggil - nomor yang baru terbit di
// transaksi yang sama, dan pengenal work yang sama dengan salinan detail.
func (r *SummaryPolis) KepalaSummaryWarisan(ctx context.Context, tx *db.Tx,
	polisID, nomorPL string) (KepalaSummaryWarisan, error) {

	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return KepalaSummaryWarisan{}, err
	}
	q := sqlKepalaSummaryWarisan(polis)
	if err := db.PeriksaSQL(q); err != nil {
		return KepalaSummaryWarisan{}, err
	}
	var cob, edm sql.NullString
	if err := tx.QueryRowContext(ctx, q, polisID).Scan(&cob, &edm); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return KepalaSummaryWarisan{}, fmt.Errorf("%w: %q", ErrPolisTakDitemukan, polisID)
		}
		return KepalaSummaryWarisan{}, fmt.Errorf("repository: membaca kepala summary warisan: %w", err)
	}
	// ⛔ Apa adanya, TANPA pangkas - sama dengan `p.PL_NUMBER_EDM` salinan
	// detail (temuan /code-review GILIRAN-18): dua tabel warisan satu submit
	// tidak boleh berselisih pada kunci yang sama.
	return KepalaSummaryWarisan{
		NomorPL:  nomorPL,
		NomorEDM: edm.String,
		COB:      cob.String,
		IDPega:   polisID,
	}, nil
}

// sqlSumberWarisan merakit pembacaan baris sumber salinan warisan - pl2.
//
// ⛔ Ia membaca TABEL KAMI (`T_PREMIUM_LIST_DETAIL` + header), bukan tabel
// warisan, dan karena itu tinggal di sini - bukan di polis_warisan.go, yang
// hanya MENULIS ke tabel warisan dan dijaga penjaga 66,8 juta baris.
//
// ⚠️ Angka lewat `TO_CHAR` ber-NLS, tanggal lewat `TO_CHAR(…,'DD/MM/YYYY')`
// - bentuk yang persis diminta `To_date(…,'DD/MM/YYYY')` rule aslinya di sisi
// tulis. Kolom teks-tanggal (`STNC`, `WPC`) sudah berbentuk itu di migrasi 052.
func sqlSumberWarisan(detail, polis string) string {
	pilih := []string{"d.ID", "p.TYPE"}
	for _, k := range kolomPesertaWarisan {
		switch {
		case k.Sumber == "":
			continue
		case k.Jenis == nilaiAngka:
			pilih = append(pilih, fmt.Sprintf(db.FmtDesimal, k.Sumber))
		case k.Jenis == nilaiTanggal:
			pilih = append(pilih, "TO_CHAR("+k.Sumber+", '"+BentukTanggalOracle+"')")
		default:
			pilih = append(pilih, k.Sumber)
		}
	}
	return fmt.Sprintf(`SELECT %s
	   FROM %s d JOIN %s p ON p.ID = d.PREMIUM_LIST_ID
	  WHERE d.PREMIUM_LIST_ID = :1 ORDER BY d.ID`,
		strings.Join(pilih, ", "), detail, polis)
}

// BarisWarisan adalah satu baris sumber salinan, berkunci kolom warisan.
type BarisWarisan struct {
	// IDPeserta adalah `T_PREMIUM_LIST_DETAIL.ID` - hanya untuk pesan galat.
	IDPeserta string
	// Tipe adalah `T_PREMIUM_LIST.TYPE` - bahan `STATUS`.
	Tipe string
	// Nilai memuat teks tiap kolom warisan yang bersumber; NULL = absen.
	Nilai map[string]sql.NullString
}

// SumberWarisan membaca baris peserta polis sebagai bahan salinan warisan.
func (r *SummaryPolis) SumberWarisan(ctx context.Context, tx *db.Tx, polisID string) (
	[]BarisWarisan, error) {

	detail, err := r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	if err != nil {
		return nil, err
	}
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return nil, err
	}
	q := sqlSumberWarisan(detail, polis)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, q, polisID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca sumber salinan warisan: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var bersumber []string
	for _, k := range kolomPesertaWarisan {
		if k.Sumber != "" {
			bersumber = append(bersumber, k.Kolom)
		}
	}
	var keluar []BarisWarisan
	for rows.Next() {
		var id string
		var tipe sql.NullString
		sel := make([]sql.NullString, len(bersumber))
		tujuan := []any{&id, &tipe}
		for i := range sel {
			tujuan = append(tujuan, &sel[i])
		}
		if err := rows.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: memindai sumber salinan warisan: %w", err)
		}
		b := BarisWarisan{
			IDPeserta: id,
			Tipe:      strings.TrimSpace(tipe.String),
			Nilai:     map[string]sql.NullString{},
		}
		for i, k := range bersumber {
			b.Nilai[k] = sel[i]
		}
		keluar = append(keluar, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca sumber salinan warisan: %w", err)
	}
	return keluar, nil
}
