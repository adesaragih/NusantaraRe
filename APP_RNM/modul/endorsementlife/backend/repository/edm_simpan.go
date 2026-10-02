package repository

// Simpan kasus - `Save` b37202 → `Activity/SetPremi_EDM.xml` (9 langkah, nol `//`).
//
// ⛔ JURNAL BALIK IDEMPOTEN (RALAT R06). Pega mengalikan 32 kolom `× -1` setiap
// kali langkah 2.1 berjalan, tanpa penjaga baris yang sudah minus. Di sini
// pembalikan hanya mengenai baris yang MASIH `Old`, dan statusnya berubah di
// pernyataan yang sama - pembalikan kedua mustahil.
//
// ⛔ Rekap mata uang dihitung ulang dari peserta kasus (`AppendCurrencySummary_DT`
// versi Endorsement, R08): tanpa pembulatan `@divide(…,1,20)` (kolom kita
// `NUMBER(38,8)` sudah ≤ 8 desimal) dan tanpa empat kolom SUM.
//
// Dibaca sesudah: edm_buat.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
)

// batasDaftarID - satu `IN (…)` Oracle menampung paling banyak 1000 nilai.
const batasDaftarID = 500

// batasKecuali - pengecualian `DELETE ALL` dalam satu pernyataan (beberapa
// `NOT IN` ber-AND); jauh di bawah batas 65.535 penampung Oracle.
const batasKecuali = 30000

// ErrDaftarKecualiTerlaluPanjang - `DELETE ALL` lalu terlalu banyak pengecualian.
var ErrDaftarKecualiTerlaluPanjang = errors.New("repository: daftar pengecualian melebihi batas satu pernyataan")

// setJurnalBalik - `SET EDM_STATUS = :1, X = -X …` untuk 32 kolom `SetPremi_EDM` 2.1.
func setJurnalBalik() string {
	bagian := []string{"d.EDM_STATUS = :1"}
	for _, k := range models.KolomJurnalBalik {
		bagian = append(bagian, fmt.Sprintf("d.%s = -d.%s", k, k))
	}
	return strings.Join(bagian, ", ")
}

// daftarPenampung - `:awal, :awal+1, …` sebanyak n.
func daftarPenampung(awal, n int) string {
	p := make([]string, n)
	for i := range p {
		p[i] = fmt.Sprintf(":%d", awal+i)
	}
	return strings.Join(p, ", ")
}

// sqlTandai - peserta `Old` kasus menjadi `Delete`/`Batal` dan dibalik tandanya.
// `pilih` > 0: hanya ID di daftar; `kecuali` > 0: semua selain daftar (per 500
// dalam `NOT IN` ber-AND); keduanya nol: seluruh peserta `Old`. Penampung :1
// status baru, :2 kasus, :3 'Old', :4… ID.
func sqlTandai(peserta string, pilih, kecuali int) string {
	q := fmt.Sprintf(`UPDATE %s d SET %s WHERE d.PREMIUM_LIST_ID = :2 AND d.EDM_STATUS = :3`, peserta, setJurnalBalik())
	switch {
	case pilih > 0:
		q += " AND d.ID IN (" + daftarPenampung(4, pilih) + ")"
	case kecuali > 0:
		for awal := 0; awal < kecuali; awal += batasDaftarID {
			q += " AND d.ID NOT IN (" + daftarPenampung(4+awal, min(batasDaftarID, kecuali-awal)) + ")"
		}
	}
	return q
}

// Tandai menjalankan `SetPremi_EDM` 2.1-2.3: peserta `Old` terpilih (atau
// seluruhnya) menjadi `statusBaru` dan 32 kolom uangnya dibalik. Mengembalikan
// cacah baris yang berubah.
func (g *Gudang) Tandai(ctx context.Context, tx *db.Tx, kasusID, statusBaru string, p models.PilihanHapus) (int, error) {
	if statusBaru != models.StatusDelete && statusBaru != models.StatusBatal {
		return 0, fmt.Errorf("repository: status tandai %q tidak sah", statusBaru)
	}
	n, err := g.nama(tabelPeserta)
	if err != nil {
		return 0, err
	}
	jalankan := func(q string, args ...any) (int, error) {
		if err := db.PeriksaSQL(q); err != nil {
			return 0, err
		}
		h, err := tx.ExecContext(ctx, q, args...)
		if err != nil {
			return 0, fmt.Errorf("repository: menandai peserta kasus %q: %w", kasusID, err)
		}
		c, err := h.RowsAffected()
		return int(c), err
	}
	dasar := []any{statusBaru, kasusID, models.StatusOld}
	switch {
	case p.Semua:
		if len(p.Kecuali) > batasKecuali {
			return 0, ErrDaftarKecualiTerlaluPanjang
		}
		args := dasar
		for _, id := range p.Kecuali {
			args = append(args, id)
		}
		return jalankan(sqlTandai(n[0], 0, len(p.Kecuali)), args...)
	default:
		total := 0
		for awal := 0; awal < len(p.Pilih); awal += batasDaftarID {
			akhir := min(awal+batasDaftarID, len(p.Pilih))
			args := append([]any{}, dasar...)
			for _, id := range p.Pilih[awal:akhir] {
				args = append(args, id)
			}
			c, err := jalankan(sqlTandai(n[0], akhir-awal, 0), args...)
			if err != nil {
				return 0, err
			}
			total += c
		}
		return total, nil
	}
}

// --- rekap mata uang ------------------------------------------------------------

// kolomJumlahRekap - kolom `T_PREMIUM_LIST_SUMMARY` yang = jumlah kolom peserta
// bernama sama (`AppendCurrencySummary_DT` langkah 2.34.1.5-2.34.1.34, versi
// Endorsement). `COMMISSION` ← `COMM` (2.34.1.5).
var kolomJumlahRekap = [][2]string{
	{"COMMISSION", "COMM"}, {"BROKERAGE_FEE", "BROKERAGE_FEE"}, {"OVR_COMM", "OVR_COMM"}, {"TAX", "TAX"},
	{"PROF_COMM", "PROF_COMM"}, {"CLAIM", "CLAIM"}, {"NET_PREMIUM_REFUND", "NET_PREMIUM_REFUND"},
	{"GROSS_PREMIUM_REFUND", "GROSS_PREMIUM_REFUND"}, {"COMM_REFUND", "COMM_REFUND"},
	{"BROKERAGE_FEE_REFUND", "BROKERAGE_FEE_REFUND"}, {"OVR_COMM_REFUND", "OVR_COMM_REFUND"},
	{"TAX_REFUND", "TAX_REFUND"}, {"SHARE_RETRO", "SHARE_RETRO"}, {"GROSS_PREMIUM_RETRO", "GROSS_PREMIUM_RETRO"},
	{"DISCOUNT_PREMIUM_RETRO", "DISCOUNT_PREMIUM_RETRO"}, {"OVR_COMM_RETRO", "OVR_COMM_RETRO"},
	{"BROKERAGE_FEE_RETRO", "BROKERAGE_FEE_RETRO"}, {"NET_PREMIUM_RETRO", "NET_PREMIUM_RETRO"},
	{"GROSS_PREMIUM_REFUND_RETRO", "GROSS_PREMIUM_REFUND_RETRO"},
	{"DISCOUNT_PREMIUM_REFUND_RETRO", "DISCOUNT_PREMIUM_REFUND_RETRO"},
	{"OVR_COMM_REFUND_RETRO", "OVR_COMM_REFUND_RETRO"}, {"BROKERAGE_FEE_REFUND_RETRO", "BROKERAGE_FEE_REFUND_RETRO"},
	{"NET_PREMIUM_REFUND_RETRO", "NET_PREMIUM_REFUND_RETRO"}, {"CLAIM_AMOUNT", "CLAIM_AMOUNT"},
	{"RI_ADMIN_FEE_RETRO", "RI_ADMIN_FEE_RETRO"}, {"RI_ADMIN_FEE_REFUND_RETRO", "RI_ADMIN_FEE_REFUND_RETRO"},
	{"RI_ADMIN_FEE_REFUND", "RI_ADMIN_FEE_REFUND"}, {"DEDUCTION_REFUND", "DEDUCTION_REFUND"},
	{"RI_ADMIN_FEE", "RI_ADMIN_FEE"}, {"DEDUCTION", "DEDUCTION"},
}

// nv - `NVL(d.X, 0)`: properti desimal Pega kosong dijumlah sebagai nol.
func nv(k string) string { return "NVL(d." + k + ", 0)" }

// rumusPremiBalance - `PREMIUM` dan `BALANCE` per `.Type`
// (`AppendCurrencySummary_DT` 2.34.1.1-2.34.1.4 dan 2.34.1.35-2.34.1.38). Tipe
// lain: kedua cabang tidak berjalan, nilai awal "0" (2.2) bertahan.
func rumusPremiBalance(tipe string) (premi, balance string) {
	switch tipe {
	case models.TypeQR:
		return "SUM(" + nv("GROSS_PREMIUM") + ")",
			"SUM(" + nv("GROSS_PREMIUM") + " - " + nv("DEDUCTION") + " - (" + nv("RI_ADMIN_FEE") + " + " + nv("BROKERAGE_FEE") +
				" + " + nv("TAX") + " + " + nv("PROF_COMM") + " + " + nv("CLAIM") + "))"
	case models.TypeQP:
		return "SUM(" + nv("GROSS_PREMIUM_REFUND") + ")",
			"SUM(" + nv("GROSS_PREMIUM_REFUND") + " + " + nv("CLAIM_AMOUNT") + " - (" + nv("DEDUCTION_REFUND") + " + " +
				nv("BROKERAGE_FEE_REFUND") + " + " + nv("RI_ADMIN_FEE_REFUND") + " + " + nv("TAX") + " + " + nv("PROF_COMM") +
				" + " + nv("CLAIM") + "))"
	case models.TypeTP:
		return "SUM(" + nv("GROSS_PREMIUM_RETRO") + ")",
			"SUM(" + nv("GROSS_PREMIUM_RETRO") + " - " + nv("DISCOUNT_PREMIUM_RETRO") + " - " + nv("RI_ADMIN_FEE_RETRO") +
				" + " + nv("BROKERAGE_FEE_RETRO") + ")"
	case models.TypeTR:
		return "SUM(" + nv("GROSS_PREMIUM_REFUND_RETRO") + ")",
			"SUM(" + nv("GROSS_PREMIUM_REFUND_RETRO") + " - " + nv("DISCOUNT_PREMIUM_REFUND_RETRO") + " - " +
				nv("RI_ADMIN_FEE_REFUND_RETRO") + " + " + nv("BROKERAGE_FEE_REFUND_RETRO") + ")"
	}
	return "0", "0"
}

// sqlHapusRekap / sqlSisipRekap - rekap kasus dihapus lalu dihitung ulang
// (satu baris per mata uang). Kunci `IDX_PLSUM_PL` (PREMIUM_LIST_ID).
func sqlHapusRekap(rekap string) string {
	return fmt.Sprintf(`DELETE FROM %s r WHERE r.PREMIUM_LIST_ID = :1`, rekap)
}

func sqlSisipRekap(rekap, peserta, tipe string) string {
	premi, balance := rumusPremiBalance(tipe)
	kolom := []string{"ID", "PREMIUM_LIST_ID", "CURRENCY", "PREMIUM", "BALANCE"}
	nilai := []string{"RAWTOHEX(STANDARD_HASH(:1 || '/C/' || NVL(d.CURRENCY, '-'), 'MD5'))", ":2", "d.CURRENCY", premi, balance}
	for _, k := range kolomJumlahRekap {
		kolom = append(kolom, k[0])
		nilai = append(nilai, "SUM("+nv(k[1])+")")
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) SELECT %s FROM %s d WHERE d.PREMIUM_LIST_ID = :3 GROUP BY d.CURRENCY`,
		rekap, strings.Join(kolom, ", "), strings.Join(nilai, ", "), peserta)
}

// HitungRekap menghitung ulang rekap mata uang kasus; mengembalikan cacah baris.
func (g *Gudang) HitungRekap(ctx context.Context, tx *db.Tx, kasusID, tipe string) (int, error) {
	n, err := g.nama(tabelRekap, tabelPeserta)
	if err != nil {
		return 0, err
	}
	hapus, sisip := sqlHapusRekap(n[0]), sqlSisipRekap(n[0], n[1], tipe)
	for _, q := range []string{hapus, sisip} {
		if err := db.PeriksaSQL(q); err != nil {
			return 0, err
		}
	}
	if _, err := tx.ExecContext(ctx, hapus, kasusID); err != nil {
		return 0, fmt.Errorf("repository: menghapus rekap kasus %q: %w", kasusID, err)
	}
	h, err := tx.ExecContext(ctx, sisip, kasusID, kasusID, kasusID)
	if err != nil {
		return 0, fmt.Errorf("repository: menghitung rekap kasus %q: %w", kasusID, err)
	}
	c, err := h.RowsAffected()
	return int(c), err
}

// sqlRekapKasus - rekap mata uang kasus (grid `InputEDMLife` b23064 …).
func sqlRekapKasus(rekap string) string {
	kolom := []string{"r.CURRENCY"}
	for _, k := range models.KolomRekapKasus[1:] {
		kolom = append(kolom, "TO_CHAR(r."+k+", 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')")
	}
	return fmt.Sprintf(`SELECT %s FROM %s r WHERE r.PREMIUM_LIST_ID = :1 ORDER BY r.CURRENCY`, strings.Join(kolom, ", "), rekap)
}

// RekapKasus membaca rekap mata uang kasus.
func (g *Gudang) RekapKasus(ctx context.Context, tx *db.Tx, kasusID string) ([]map[string]string, error) {
	n, err := g.nama(tabelRekap)
	if err != nil {
		return nil, err
	}
	q := sqlRekapKasus(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.pakai(tx).QueryContext(ctx, q, kasusID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca rekap kasus %q: %w", kasusID, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []map[string]string{}
	for rows.Next() {
		v := make([]sql.NullString, len(models.KolomRekapKasus))
		ptr := make([]any, len(v))
		for i := range v {
			ptr[i] = &v[i]
		}
		if err := rows.Scan(ptr...); err != nil {
			return nil, fmt.Errorf("repository: memindai rekap kasus: %w", err)
		}
		baris := map[string]string{"CURRENCY": strings.TrimSpace(v[0].String)}
		for i, k := range models.KolomRekapKasus[1:] {
			s, err := rapikanAngka(k, v[i+1])
			if err != nil {
				return nil, err
			}
			baris[k] = s
		}
		hasil = append(hasil, baris)
	}
	return hasil, rows.Err()
}
