package repository

// Pemindah data warisan Treaty Contract Out ke tabel T_* - bagian ORACLE
// (tiket 01, keputusan tco2: kode + uji bertag db; TIDAK dijalankan di DEV
// oleh sesi modul - `-migrate-data-treaty-contract-out` milik work owner
// sesudah penyatuan ke main).
//
// Alurnya, dalam SATU transaksi:
//  1. tabel T_* wajib KOSONG - pemindahan tidak pernah diulang di atas isi;
//  2. enam tabel warisan dibaca sebagai TEKS (TO_CHAR berformat, RTRIM
//     untuk CHAR), kolom mati PROPORTIONALLIST/OBJECT dicacah;
//  3. konversi MURNI (tco_migrasidata.go); temuan yang memblokir
//     MEMBATALKAN seluruhnya - tidak ada baris yang ditulis separuh;
//  4. seluruh baris ditulis, induk sebelum anak; MTREATYSECURITY menerima
//     ID surrogate dari sequence;
//  5. rekonsiliasi: tabel T_* dibaca KEMBALI dan dibandingkan tepat dengan
//     nilai yang dikonversi - selisih membatalkan transaksi;
//  6. commit; lalu sequence diselaraskan ke ekor identitas terbesar + 1.
//
// ⛔ Tabel warisan HANYA DIBACA (spec b106: modul ini penulis tunggal tabel
// BARU; tabel warisan tidak disentuh - tco1). Penjaga TestTCOWarisanHanyaDibaca.
//
// ⛔ Nol pembacaan tabel dokumen M_* (AC 64).
//
// ⚠️ Penyelarasan sequence memakai DDL (DROP + CREATE START WITH). DDL Oracle
// menutup transaksinya sendiri - itu sifat Oracle, bukan alasan menulis
// COMMIT (migrasi.go) - karena itu ia dijalankan SESUDAH data ter-commit dan
// kegagalannya dilaporkan menyebut sequence-nya; datanya tetap utuh.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/pkg/utils"
)

var (
	// ErrMigrasiTCOTidakBersih - ada temuan yang memblokir; nol baris ditulis.
	ErrMigrasiTCOTidakBersih = errors.New(
		"repository: migrasi data Treaty Contract Out dibatalkan: ada temuan yang memblokir; lihat laporan")
	// ErrTabelBaruTCOSudahBerisi - pemindahan tidak diulang di atas isi.
	ErrTabelBaruTCOSudahBerisi = errors.New(
		"repository: tabel T_TREATY* sudah berisi; migrasi data Treaty Contract Out tidak diulang")
	// ErrRekonsiliasiTCOBerselisih - nilai yang dibaca kembali berbeda.
	ErrRekonsiliasiTCOBerselisih = errors.New(
		"repository: rekonsiliasi migrasi data Treaty Contract Out menemukan selisih; transaksi dibatalkan")
)

// MigrasiTCO memindahkan data warisan enam tabel master arrangement.
type MigrasiTCO struct{ db *DB }

// NewMigrasiTCO menyusunnya.
func NewMigrasiTCO(db *DB) *MigrasiTCO { return &MigrasiTCO{db: db} }

// sqlSamakanNLSTCO memaksa titik sebagai pemisah desimal untuk sesi ini.
//
// Nilai desimal di-bind sebagai TEKS ke kolom NUMBER (ADR-U-0034), dan
// konversi implisit itu mengikuti NLS sesi. Pola yang sama dengan skemauji.
const sqlSamakanNLSTCO = `ALTER SESSION SET NLS_NUMERIC_CHARACTERS = '.,'`

// ekspresiBacaWarisanTCO membungkus satu kolom warisan supaya tiba sebagai
// TEKS yang bentuknya pasti, menurut tipe deklarasinya `[data DBA]`.
func ekspresiBacaWarisanTCO(tabel, kolom string) string {
	switch TipeWarisanTCO(tabel, kolom) {
	case WarisanAngka:
		return fmt.Sprintf(fmtDesimal, kolom)
	case WarisanTanggal:
		return fmt.Sprintf(fmtTanggalOracle, kolom)
	case WarisanChar:
		return fmt.Sprintf("RTRIM(%s)", kolom)
	}
	return kolom
}

// ekspresiBacaBaruTCO membungkus satu kolom tabel BARU untuk rekonsiliasi,
// menurut tipe tujuannya.
func ekspresiBacaBaruTCO(warisan, kolom string) string {
	switch TipeTujuanTCO(warisan, kolom) {
	case TujuanDesimal:
		return fmt.Sprintf(fmtDesimal, kolom)
	case TujuanTanggal:
		return fmt.Sprintf(fmtTanggalOracle, kolom)
	}
	return kolom
}

// sqlBacaWarisanTCO merakit pembacaan satu tabel warisan.
//
// Urutannya dijaga: ID bila ada, ROWID bila tidak (MTREATYSECURITY).
func sqlBacaWarisanTCO(tabelBerskema, tabel string) string {
	kolom := KolomWarisanTCO(tabel)
	ekspresi := make([]string, 0, len(kolom))
	for _, k := range kolom {
		ekspresi = append(ekspresi, ekspresiBacaWarisanTCO(tabel, k))
	}
	urut := "ID"
	if tabel == warisanSecurityTCO {
		urut = "ROWID"
	}
	return fmt.Sprintf("SELECT %s FROM %s ORDER BY %s", strings.Join(ekspresi, ", "), tabelBerskema, urut)
}

// sqlBacaBaruTCO merakit pembacaan kembali satu tabel baru (ID dahulu).
func sqlBacaBaruTCO(tabelBerskema, warisan string) string {
	kolom := KolomWarisanTCO(warisan)
	ekspresi := []string{"ID"}
	for _, k := range kolom {
		if k == "ID" {
			continue
		}
		ekspresi = append(ekspresi, ekspresiBacaBaruTCO(warisan, k))
	}
	return fmt.Sprintf("SELECT %s FROM %s ORDER BY ID", strings.Join(ekspresi, ", "), tabelBerskema)
}

// sqlSisipBaruTCO merakit INSERT satu tabel baru, kolom BERNAMA seluruhnya.
//
// ⛔ Tidak ada INSERT posisional (AC 19). Kolom ID selalu pertama; untuk
// MTREATYSECURITY ia surrogate yang lahir di sini.
func sqlSisipBaruTCO(tabelBerskema string, kolom []string) string {
	nama := []string{"ID"}
	for _, k := range kolom {
		if k != "ID" {
			nama = append(nama, k)
		}
	}
	penampung := make([]string, len(nama))
	for i := range nama {
		penampung[i] = fmt.Sprintf(":%d", i+1)
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tabelBerskema,
		strings.Join(nama, ", "), strings.Join(penampung, ", "))
}

// sqlCacahIsiTCO mencacah baris sebuah tabel.
func sqlCacahIsiTCO(tabelBerskema string) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s", tabelBerskema)
}

// sqlCacahKolomMatiTCO mencacah baris PROPORTIONALARRG yang kolom matinya
// terisi (AC 70).
func sqlCacahKolomMatiTCO(tabelBerskema, kolom string) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s IS NOT NULL", tabelBerskema, kolom)
}

// nilaiBindTCO mengubah nilai siap tulis menjadi argumen bind.
//
// Desimal menyeberang sebagai TEKS (ADR-U-0034), tanggal sebagai time.Time,
// kosong sebagai NULL.
func nilaiBindTCO(v any) any {
	switch x := v.(type) {
	case nil:
		return nil
	case *apd.Decimal:
		if x == nil {
			return nil
		}
		return utils.FormatDecimal(x)
	default:
		return x
	}
}

// kanonikDesimalTCO menyamakan bentuk teks desimal untuk perbandingan TEPAT:
// nilai yang sama dengan nol di ekor berbeda ("12.50" dari konversi,
// "12.5" dari TO_CHAR TM9) adalah nilai yang SAMA.
func kanonikDesimalTCO(teks string) (string, error) {
	t := strings.TrimSpace(teks)
	if t == "" {
		return "", nil
	}
	d, err := utils.ParseDecimal(t)
	if err != nil {
		return "", err
	}
	d.Reduce(d)
	return utils.FormatDecimal(d), nil
}

// SelisihRekonsiliasiTCO membandingkan baris yang dikonversi dengan baris
// yang dibaca KEMBALI dari tabel baru. MURNI.
//
// `kembali` berkunci nama tabel BARU; tiap barisnya memuat ID dan seluruh
// kolom warisan sebagai teks. Tabel ber-ID dibandingkan per ID dan per
// kolom; MTREATYSECURITY - yang ID-nya surrogate - dibandingkan sebagai
// HIMPUNAN GANDA kunci kanonik.
func SelisihRekonsiliasiTCO(siap []BarisSiapTCO, kembali map[string][]BarisWarisanTCO) []string {
	var selisih []string
	perTabel := map[string][]BarisSiapTCO{}
	for _, s := range siap {
		perTabel[s.Warisan] = append(perTabel[s.Warisan], s)
	}
	for _, warisan := range urutanWarisanTCO {
		tabelBaru := tabelBaruDariWarisanTCO[warisan]
		harap := perTabel[warisan]
		dapat := kembali[tabelBaru]
		if len(harap) != len(dapat) {
			selisih = append(selisih, fmt.Sprintf("%s: ditulis %d baris, terbaca kembali %d",
				tabelBaru, len(harap), len(dapat)))
		}
		kolom := KolomWarisanTCO(warisan)
		if warisan == warisanSecurityTCO {
			selisih = append(selisih, selisihHimpunanGandaTCO(tabelBaru, kolom, harap, dapat)...)
			continue
		}
		dapatPerID := map[string]BarisWarisanTCO{}
		for _, b := range dapat {
			dapatPerID[b.Nilai["ID"]] = b
		}
		for _, h := range harap {
			b, ada := dapatPerID[h.ID]
			if !ada {
				selisih = append(selisih, fmt.Sprintf("%s/%s: tidak terbaca kembali", tabelBaru, h.ID))
				continue
			}
			delete(dapatPerID, h.ID)
			for _, k := range kolom {
				if k == "ID" {
					continue
				}
				mau, dapatTeks, err := bandingkanKolomTCO(warisan, k, h.Kanonik[k], b.Nilai[k])
				if err != nil {
					selisih = append(selisih, fmt.Sprintf("%s/%s.%s: %v", tabelBaru, h.ID, k, err))
					continue
				}
				if mau != dapatTeks {
					selisih = append(selisih, fmt.Sprintf("%s/%s.%s: ditulis %q, terbaca %q",
						tabelBaru, h.ID, k, mau, dapatTeks))
				}
			}
		}
		for id := range dapatPerID {
			selisih = append(selisih, fmt.Sprintf("%s/%s: ada di tabel baru tetapi tidak dikonversi", tabelBaru, id))
		}
	}
	return selisih
}

// bandingkanKolomTCO menyamakan bentuk kedua sisi menurut tipe tujuan.
func bandingkanKolomTCO(warisan, kolom, kanonik, terbaca string) (string, string, error) {
	if TipeTujuanTCO(warisan, kolom) == TujuanDesimal {
		a, err := kanonikDesimalTCO(kanonik)
		if err != nil {
			return "", "", fmt.Errorf("nilai konversi %q: %w", kanonik, err)
		}
		b, err := kanonikDesimalTCO(terbaca)
		if err != nil {
			return "", "", fmt.Errorf("nilai terbaca %q: %w", terbaca, err)
		}
		return a, b, nil
	}
	return kanonik, terbaca, nil
}

// selisihHimpunanGandaTCO membandingkan MTREATYSECURITY sebagai himpunan ganda.
func selisihHimpunanGandaTCO(tabelBaru string, kolom []string, harap []BarisSiapTCO,
	dapat []BarisWarisanTCO) []string {

	var selisih []string
	cacah := map[string]int{}
	for _, h := range harap {
		kanonik := map[string]string{}
		for _, k := range kolom {
			v, _, err := bandingkanKolomTCO(warisanSecurityTCO, k, h.Kanonik[k], "")
			if err != nil {
				selisih = append(selisih, fmt.Sprintf("%s.%s: %v", tabelBaru, k, err))
			}
			kanonik[k] = v
		}
		cacah[KunciKanonikTCO(kolom, kanonik)]++
	}
	for _, b := range dapat {
		kanonik := map[string]string{}
		for _, k := range kolom {
			_, v, err := bandingkanKolomTCO(warisanSecurityTCO, k, "", b.Nilai[k])
			if err != nil {
				selisih = append(selisih, fmt.Sprintf("%s/%s.%s: %v", tabelBaru, b.Nilai["ID"], k, err))
			}
			kanonik[k] = v
		}
		kunci := KunciKanonikTCO(kolom, kanonik)
		if cacah[kunci] == 0 {
			selisih = append(selisih, fmt.Sprintf("%s/%s: terbaca tetapi tidak dikonversi: %s",
				tabelBaru, b.Nilai["ID"], kunci))
			continue
		}
		cacah[kunci]--
	}
	for kunci, n := range cacah {
		if n > 0 {
			selisih = append(selisih, fmt.Sprintf("%s: %d baris tidak terbaca kembali: %s", tabelBaru, n, kunci))
		}
	}
	return selisih
}

// bacaTeksTCO menjalankan satu SELECT dan mengembalikan tiap baris sebagai
// peta kolom -> teks (NULL menjadi teks kosong).
func bacaTeksTCO(ctx context.Context, tx *Tx, q string, kolom []string) ([]map[string]string, error) {
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := tx.tx.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var hasil []map[string]string
	for rows.Next() {
		nilai := make([]sql.NullString, len(kolom))
		tujuan := make([]any, len(kolom))
		for i := range nilai {
			tujuan[i] = &nilai[i]
		}
		if err := rows.Scan(tujuan...); err != nil {
			return nil, err
		}
		baris := make(map[string]string, len(kolom))
		for i, k := range kolom {
			baris[k] = nilai[i].String
		}
		hasil = append(hasil, baris)
	}
	return hasil, rows.Err()
}

// bacaWarisan membaca seluruh baris satu tabel warisan.
func (m *MigrasiTCO) bacaWarisan(ctx context.Context, tx *Tx, tabel string) ([]BarisWarisanTCO, error) {
	nama, err := m.db.Qualify(tabel)
	if err != nil {
		return nil, err
	}
	kolom := KolomWarisanTCO(tabel)
	baris, err := bacaTeksTCO(ctx, tx, sqlBacaWarisanTCO(nama, tabel), kolom)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca warisan %s (tipe kolom `[data DBA]`; selisih tipe gagal di sini): %w",
			tabel, err)
	}
	hasil := make([]BarisWarisanTCO, 0, len(baris))
	for _, b := range baris {
		hasil = append(hasil, BarisWarisanTCO{Tabel: tabel, Nilai: b})
	}
	return hasil, nil
}

// bacaKembali membaca seluruh baris satu tabel baru untuk rekonsiliasi.
func (m *MigrasiTCO) bacaKembali(ctx context.Context, tx *Tx, warisan string) ([]BarisWarisanTCO, error) {
	tabelBaru := tabelBaruDariWarisanTCO[warisan]
	nama, err := m.db.Qualify(tabelBaru)
	if err != nil {
		return nil, err
	}
	kolom := []string{"ID"}
	for _, k := range KolomWarisanTCO(warisan) {
		if k != "ID" {
			kolom = append(kolom, k)
		}
	}
	baris, err := bacaTeksTCO(ctx, tx, sqlBacaBaruTCO(nama, warisan), kolom)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca kembali %s: %w", tabelBaru, err)
	}
	hasil := make([]BarisWarisanTCO, 0, len(baris))
	for _, b := range baris {
		hasil = append(hasil, BarisWarisanTCO{Tabel: tabelBaru, Nilai: b})
	}
	return hasil, nil
}

// cacah menjalankan satu SELECT COUNT(*).
func (m *MigrasiTCO) cacah(ctx context.Context, tx *Tx, q string) (int, error) {
	if err := PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int
	if err := tx.tx.QueryRowContext(ctx, q).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

// pastikanTabelBaruKosong menolak pemindahan di atas isi.
func (m *MigrasiTCO) pastikanTabelBaruKosong(ctx context.Context, tx *Tx) error {
	for _, warisan := range urutanWarisanTCO {
		tabelBaru := tabelBaruDariWarisanTCO[warisan]
		nama, err := m.db.Qualify(tabelBaru)
		if err != nil {
			return err
		}
		n, err := m.cacah(ctx, tx, sqlCacahIsiTCO(nama))
		if err != nil {
			return fmt.Errorf("repository: memeriksa isi %s: %w", tabelBaru, err)
		}
		if n > 0 {
			return fmt.Errorf("%w: %s memuat %d baris", ErrTabelBaruTCOSudahBerisi, tabelBaru, n)
		}
	}
	return nil
}

// tulis menulis satu baris siap ke tabel barunya.
func (m *MigrasiTCO) tulis(ctx context.Context, tx *Tx, s *BarisSiapTCO) error {
	nama, err := m.db.Qualify(s.Tabel)
	if err != nil {
		return err
	}
	if s.ID == "" {
		if s.Warisan != warisanSecurityTCO {
			return fmt.Errorf("repository: %s tanpa identitas tidak dapat ditulis", s.Tabel)
		}
		id, err := m.db.IdentitasBerikutTCO(ctx, tx, SeqSecurityTCO)
		if err != nil {
			return err
		}
		s.ID = id
	}
	q := sqlSisipBaruTCO(nama, s.Kolom)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	arg := []any{s.ID}
	for _, k := range s.Kolom {
		if k == "ID" {
			continue
		}
		arg = append(arg, nilaiBindTCO(s.Nilai[k]))
	}
	hasil, err := tx.tx.ExecContext(ctx, q, arg...)
	if err != nil {
		return fmt.Errorf("repository: menulis %s/%s: %w", s.Tabel, s.ID, err)
	}
	return pastikanSatuBaris(hasil, "penulisan "+s.Tabel)
}

// Pindahkan menjalankan seluruh alur pemindahan.
//
// Laporan SELALU dikembalikan, juga saat gagal - supaya temuan dan selisih
// terbaca orang yang menjalankannya.
func (m *MigrasiTCO) Pindahkan(ctx context.Context) (LaporanMigrasiTCO, error) {
	lap := laporanKosongTCO()
	if m == nil || m.db == nil || m.db.sql == nil {
		return lap, ErrTanpaOracle
	}
	tx, err := m.db.Mulai(ctx)
	if err != nil {
		return lap, err
	}
	defer func() { _ = tx.Rollback() }()

	if err := PeriksaSQL(sqlSamakanNLSTCO); err != nil {
		return lap, err
	}
	if _, err := tx.tx.ExecContext(ctx, sqlSamakanNLSTCO); err != nil {
		return lap, fmt.Errorf("repository: menyamakan NLS sesi: %w", err)
	}
	if err := m.pastikanTabelBaruKosong(ctx, tx); err != nil {
		return lap, err
	}

	masuk := map[string][]BarisWarisanTCO{}
	for _, tabel := range urutanWarisanTCO {
		baris, err := m.bacaWarisan(ctx, tx, tabel)
		if err != nil {
			return lap, err
		}
		masuk[tabel] = baris
	}

	// Kolom mati (AC 70): dicacah dan dilaporkan, tidak dibawa.
	klausulLama, err := m.db.Qualify(warisanKlausulTCO)
	if err != nil {
		return lap, err
	}
	for _, kolom := range kolomMatiKlausulTCO {
		n, err := m.cacah(ctx, tx, sqlCacahKolomMatiTCO(klausulLama, kolom))
		if err != nil {
			return lap, fmt.Errorf("repository: mencacah kolom mati %s: %w", kolom, err)
		}
		lap.KolomMatiBerisi[kolom] = n
		if n > 0 {
			lap.Temuan = append(lap.Temuan, Temuan{Jenis: TemuanKolomMatiBerisi,
				Sumber: warisanKlausulTCO, Medan: kolom, Nilai: fmt.Sprint(n),
				Catatan: "tidak dibawa ke skema baru (AC 70); keputusan work owner bila isinya hidup"})
		}
	}

	siap, lapKonversi := KonversiWarisanTCO(masuk)
	lap.Masuk = lapKonversi.Masuk
	lap.SequenceBerikut = lapKonversi.SequenceBerikut
	lap.Temuan = append(lap.Temuan, lapKonversi.Temuan...)
	if lap.Memblokir() {
		return lap, ErrMigrasiTCOTidakBersih
	}

	for i := range siap {
		if err := m.tulis(ctx, tx, &siap[i]); err != nil {
			return lap, err
		}
		lap.Ditulis[siap[i].Tabel]++
	}

	kembali := map[string][]BarisWarisanTCO{}
	for _, warisan := range urutanWarisanTCO {
		baris, err := m.bacaKembali(ctx, tx, warisan)
		if err != nil {
			return lap, err
		}
		tabelBaru := tabelBaruDariWarisanTCO[warisan]
		kembali[tabelBaru] = baris
		lap.DibacaKembali[tabelBaru] = len(baris)
	}
	lap.Selisih = SelisihRekonsiliasiTCO(siap, kembali)
	if len(lap.Selisih) > 0 {
		return lap, ErrRekonsiliasiTCOBerselisih
	}

	if err := tx.Commit(); err != nil {
		return lap, fmt.Errorf("repository: commit migrasi data Treaty Contract Out: %w", err)
	}

	for _, warisan := range urutanWarisanTCO {
		seq, ada := sequenceDariWarisanTCO[warisan]
		if !ada {
			continue
		}
		if err := m.selaraskanSequence(ctx, seq.Nama, lap.SequenceBerikut[seq.Nama]); err != nil {
			return lap, err
		}
	}
	return lap, nil
}

// selaraskanSequence memasang nilai NEXTVAL berikutnya (AC 69).
//
// DROP + CREATE START WITH, bukan ALTER INCREMENT lalu NEXTVAL: bentuknya
// deterministik dan berlaku di setiap versi Oracle. Nama sequence dibatasi
// daftar modul; nilainya bilangan, bukan teks dari luar.
func (m *MigrasiTCO) selaraskanSequence(ctx context.Context, sequence string, berikut int64) error {
	if _, dikenal := sequenceDikenalTCO[sequence]; !dikenal {
		return fmt.Errorf("%w: %q", ErrSequenceTakDikenal, sequence)
	}
	if berikut < 1 {
		berikut = 1
	}
	nama, err := m.db.Qualify(sequence)
	if err != nil {
		return err
	}
	for _, q := range []string{
		fmt.Sprintf("DROP SEQUENCE %s", nama),
		fmt.Sprintf("CREATE SEQUENCE %s START WITH %d INCREMENT BY 1 NOCACHE", nama, berikut),
	} {
		if err := PeriksaSQL(q); err != nil {
			return err
		}
		if _, err := m.db.sql.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("repository: menyelaraskan %s ke %d (data sudah ter-commit): %w",
				sequence, berikut, err)
		}
	}
	return nil
}
