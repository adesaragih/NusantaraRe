package repository

// Pembacaan modul Endorsement Life: kotak masuk, kepala kasus, peserta kasus,
// dan versi polis lama di kedua sumber (RALAT bab 4).
//
// Dibaca sesudah: edm_tabel.go.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/endorsementlife/backend/models"
)

// penanya - `*db.DB` atau `*db.Tx`, keduanya menjawab kueri yang sama.
type penanya interface {
	QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row
	ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error)
}

// pakai memilih transaksi bila ada, selain itu koneksi.
func (g *Gudang) pakai(tx *db.Tx) penanya {
	if tx.Terisi() {
		return tx
	}
	return g.db
}

// polaPengenalKasus - `LIKE` atas PK, kasus modul ini saja (`EDMLF-%`).
const polaPengenalKasus = models.AwalanKasus + "%"

// ekspresiBaca membaca satu kolom sebagai teks: tanggal `YYYY-MM-DD`, angka
// lewat `TM9` bertitik desimal (dirapikan `rapikanAngka`).
func ekspresiBaca(alias string, k models.KolomPeserta) string {
	kol := alias + "." + k.Nama
	switch k.Jenis {
	case models.KolomTanggal:
		return "TO_CHAR(" + kol + ", 'YYYY-MM-DD')"
	case models.KolomAngka:
		return "TO_CHAR(" + kol + ", 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')"
	}
	return kol
}

// rapikanAngka mengubah keluaran `TM9` (mis. `.5`, `-1E+03` tidak muncul pada
// NUMBER(38,8)) menjadi teks desimal baku. Gagal urai = galat, bukan kosong.
func rapikanAngka(kolom string, v sql.NullString) (string, error) {
	if !v.Valid || strings.TrimSpace(v.String) == "" {
		return "", nil
	}
	d, err := utils.ParseDecimal(strings.TrimSpace(v.String))
	if err != nil {
		return "", fmt.Errorf("repository: kolom %s bernilai %q yang tidak terurai: %w", kolom, v.String, err)
	}
	return utils.FormatDecimal(d), nil
}

// pindaiNilai memindai deretan kolom teks ke peta berkunci nama kolom.
func pindaiNilai(kolom []models.KolomPeserta, mentah []sql.NullString) (map[string]string, error) {
	n := make(map[string]string, len(kolom))
	for i, k := range kolom {
		if k.Jenis == models.KolomAngka {
			v, err := rapikanAngka(k.Nama, mentah[i])
			if err != nil {
				return nil, err
			}
			n[k.Nama] = v
			continue
		}
		n[k.Nama] = strings.TrimSpace(mentah[i].String)
	}
	return n, nil
}

// --- kotak masuk ----------------------------------------------------------

// sqlInbox - `ReportDefinition/InboxEDMLife.xml`: bukan `Resolved-Completed`
// (b659) dan bukan `Resolved-Rejected` (b673) = `STATUSS` kosong; urut
// `pxCreateDateTime` `DESC` (b1040) + pemutus seri `ID`.
//
// ⚠️ Kolom `Endorsement No` (`InboxEndorsementLife.xml` b8556) menampilkan
// `A.PremiumListSummary.PL_NUMBER` (b10516) = nomor polis yang di-endorse
// (`MappingEDMLife` 9 b2184) - RALAT R25.
func sqlInbox(polis string) string {
	return fmt.Sprintf(`SELECT p.ID, p.OLD_POLICY_NO, p.TYPE, p.EDM_TYPE, p.SOB_NAME, p.CEDING_CO_NAME,
	        p.POLICY_HOLDER_NAME, p.MARKETING_NAME, TO_CHAR(p.TGL_INPUT, 'YYYY-MM-DD HH24:MI:SS'),
	        p.CREATE_OP_NAME, p.STATUSS
	   FROM %s p
	  WHERE p.ID LIKE :1 AND p.STATUSS IS NULL
	  ORDER BY p.TGL_INPUT DESC, p.ID
	 OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY`, polis)
}

func sqlCacahInbox(polis string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s p WHERE p.ID LIKE :1 AND p.STATUSS IS NULL`, polis)
}

// Inbox membaca satu halaman kasus terbuka beserta cacah seluruhnya.
func (g *Gudang) Inbox(ctx context.Context, halaman, ukuran int) ([]models.BarisInbox, int, error) {
	n, err := g.nama(tabelPolis)
	if err != nil {
		return nil, 0, err
	}
	var total int
	qc := sqlCacahInbox(n[0])
	if err := db.PeriksaSQL(qc); err != nil {
		return nil, 0, err
	}
	if err := g.db.QueryRowContext(ctx, qc, polaPengenalKasus).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: mencacah kotak masuk endorsement: %w", err)
	}
	q := sqlInbox(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return nil, 0, err
	}
	rows, err := g.db.QueryContext(ctx, q, polaPengenalKasus, (halaman-1)*ukuran, ukuran)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: membaca kotak masuk endorsement: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []models.BarisInbox
	for rows.Next() {
		var v [11]sql.NullString
		ptr := make([]any, len(v))
		for i := range v {
			ptr[i] = &v[i]
		}
		if err := rows.Scan(ptr...); err != nil {
			return nil, 0, fmt.Errorf("repository: memindai kotak masuk endorsement: %w", err)
		}
		hasil = append(hasil, models.BarisInbox{
			CaseID: v[0].String, EndorsementNo: v[1].String, PolicyNo: v[1].String,
			Tipe: v[2].String, EdmType: v[3].String, Sob: v[4].String, Ceding: v[5].String,
			PolicyHolder: v[6].String, MarketingName: v[7].String, CreateDate: v[8].String,
			CreateOperator: v[9].String, Status: v[10].String,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository: membaca kotak masuk endorsement: %w", err)
	}
	return hasil, total, nil
}

// --- kepala kasus -----------------------------------------------------------

// kolomKasus - kolom tetap kepala kasus, sebelum `KolomKepalaSalin`.
const kolomKasus = `p.ID, p.OLD_POLICY_NO, p.PL_NUMBER_EDM, NVL(p.PROD_KE, 1), p.EDM_TYPE, p.EDM_NOTE,
	        TO_CHAR(p.EDM_DATE, 'YYYY-MM-DD'), p.STATUSS, TO_CHAR(p.TGL_INPUT, 'YYYY-MM-DD HH24:MI:SS'),
	        p.CREATE_OP_NAME`

// sqlKasus membaca kepala satu kasus; `kunci` menambah `FOR UPDATE` (kasus
// dikunci selama transaksi yang mengubahnya - dua tab tidak dapat menyimpan
// atau memutuskan kasus yang sama bersamaan).
func sqlKasus(polis string, kunci bool) string {
	var b strings.Builder
	b.WriteString("SELECT " + kolomKasus)
	for _, k := range models.KolomKepalaSalin {
		b.WriteString(", ")
		if k.Tanggal {
			b.WriteString("TO_CHAR(p." + k.Kolom + ", 'YYYY-MM-DD')")
		} else {
			b.WriteString("p." + k.Kolom)
		}
	}
	fmt.Fprintf(&b, " FROM %s p WHERE p.ID = :1 AND p.ID LIKE :2", polis)
	if kunci {
		b.WriteString(" FOR UPDATE")
	}
	return b.String()
}

// AmbilKasus membaca kepala kasus; `ErrTidakAda` bila tidak ada.
func (g *Gudang) AmbilKasus(ctx context.Context, tx *db.Tx, id string, kunci bool) (models.Kasus, error) {
	n, err := g.nama(tabelPolis)
	if err != nil {
		return models.Kasus{}, err
	}
	q := sqlKasus(n[0], kunci)
	if err := db.PeriksaSQL(q); err != nil {
		return models.Kasus{}, err
	}
	v := make([]sql.NullString, 10+len(models.KolomKepalaSalin))
	var prodKe int
	ptr := make([]any, len(v))
	for i := range v {
		ptr[i] = &v[i]
	}
	ptr[3] = &prodKe
	err = g.pakai(tx).QueryRowContext(ctx, q, id, polaPengenalKasus).Scan(ptr...)
	if err == sql.ErrNoRows {
		return models.Kasus{}, ErrTidakAda
	}
	if err != nil {
		return models.Kasus{}, fmt.Errorf("repository: membaca kasus endorsement %q: %w", id, err)
	}
	k := models.Kasus{
		ID: v[0].String, NomorPolis: v[1].String, PLNumber: v[1].String, PLNumberEDM: v[2].String,
		ProdKe: prodKe, EdmType: v[4].String, EdmNote: v[5].String, EdmDate: v[6].String,
		Status: v[7].String, TglInput: v[8].String, Pembuat: v[9].String,
		Kepala: make(map[string]string, len(models.KolomKepalaSalin)),
	}
	for i, kol := range models.KolomKepalaSalin {
		k.Kepala[kol.Kolom] = strings.TrimSpace(v[10+i].String)
	}
	return k, nil
}

// --- peserta kasus ----------------------------------------------------------

func sqlPeserta(peserta string, kolom []models.KolomPeserta) string {
	var b strings.Builder
	b.WriteString("SELECT d.ID, d.PARENT_ID, d.EDM_STATUS")
	for _, k := range kolom {
		b.WriteString(", " + ekspresiBaca("d", k))
	}
	fmt.Fprintf(&b, ` FROM %s d WHERE d.PREMIUM_LIST_ID = :1
	  ORDER BY d.CERTIFICATE_NO, d.NAME_OF_INSURED, d.ID
	 OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY`, peserta)
	return b.String()
}

func sqlCacahPeserta(peserta string) string {
	return fmt.Sprintf(`SELECT NVL(d.EDM_STATUS, ' '), COUNT(*) FROM %s d
	  WHERE d.PREMIUM_LIST_ID = :1 GROUP BY NVL(d.EDM_STATUS, ' ')`, peserta)
}

// DaftarPeserta membaca satu halaman peserta kasus (grid b11899 / b17500).
func (g *Gudang) DaftarPeserta(ctx context.Context, kasusID string, halaman, ukuran int) ([]models.Peserta, error) {
	n, err := g.nama(tabelPeserta)
	if err != nil {
		return nil, err
	}
	q := sqlPeserta(n[0], models.KolomPesertaGrid)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, kasusID, (halaman-1)*ukuran, ukuran)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca peserta kasus %q: %w", kasusID, err)
	}
	defer func() { _ = rows.Close() }()
	return pindaiPeserta(rows, models.KolomPesertaGrid)
}

// pindaiPeserta memindai `ID, PARENT_ID, EDM_STATUS, <kolom>`.
func pindaiPeserta(rows *sql.Rows, kolom []models.KolomPeserta) ([]models.Peserta, error) {
	var hasil []models.Peserta
	for rows.Next() {
		v := make([]sql.NullString, 3+len(kolom))
		ptr := make([]any, len(v))
		for i := range v {
			ptr[i] = &v[i]
		}
		if err := rows.Scan(ptr...); err != nil {
			return nil, fmt.Errorf("repository: memindai peserta: %w", err)
		}
		nilai, err := pindaiNilai(kolom, v[3:])
		if err != nil {
			return nil, err
		}
		s := strings.TrimSpace(v[2].String)
		hasil = append(hasil, models.Peserta{
			ID: v[0].String, ParentID: v[1].String, EdmStatus: s,
			Terkunci: s == models.StatusDelete || s == models.StatusBatal, Nilai: nilai,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca peserta: %w", err)
	}
	return hasil, nil
}

// CacahPeserta menghitung peserta kasus per `EDM_STATUS` (kosong = " ").
func (g *Gudang) CacahPeserta(ctx context.Context, tx *db.Tx, kasusID string) (map[string]int, error) {
	n, err := g.nama(tabelPeserta)
	if err != nil {
		return nil, err
	}
	q := sqlCacahPeserta(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.pakai(tx).QueryContext(ctx, q, kasusID)
	if err != nil {
		return nil, fmt.Errorf("repository: mencacah peserta kasus %q: %w", kasusID, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := map[string]int{}
	for rows.Next() {
		var s string
		var c int
		if err := rows.Scan(&s, &c); err != nil {
			return nil, fmt.Errorf("repository: memindai cacah peserta: %w", err)
		}
		hasil[strings.TrimSpace(s)] = c
	}
	return hasil, rows.Err()
}

// AdaRekap - `IsJsonPolis` (`SetPremi_EDM` 8 b5600): kasus sudah disimpan bila
// rekap mata uangnya ada.
func (g *Gudang) AdaRekap(ctx context.Context, tx *db.Tx, kasusID string) (bool, error) {
	n, err := g.nama(tabelRekap)
	if err != nil {
		return false, err
	}
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s r WHERE r.PREMIUM_LIST_ID = :1`, n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var c int
	if err := g.pakai(tx).QueryRowContext(ctx, q, kasusID).Scan(&c); err != nil {
		return false, fmt.Errorf("repository: membaca rekap kasus %q: %w", kasusID, err)
	}
	return c > 0, nil
}

// --- versi polis lama -------------------------------------------------------

// sqlVersiEDM - versi endorsement RESMI sistem baru: `NO_POLIS` diisi saat
// Confirm (`IDX_PL_NOPOLIS_PRODKE`, migrasi 480). `PROD_KE DESC` - satu urutan
// (penyimpangan sadar 1, RALAT R19). `sebelum` > 0 membatasi ke versi lebih tua.
func sqlVersiEDM(polis string, sebelum bool) string {
	syarat := ""
	if sebelum {
		syarat = " AND NVL(p.PROD_KE, 1) < :3"
	}
	return fmt.Sprintf(`SELECT p.ID, NVL(p.PROD_KE, 1), p.EDM_TYPE FROM %s p
	  WHERE p.NO_POLIS = :1 AND p.STATUSS = :2%s
	  ORDER BY NVL(p.PROD_KE, 1) DESC, p.ID FETCH FIRST 1 ROWS ONLY`, polis, syarat)
}

// sqlVersiNB - versi new business sistem baru: kepala tanpa `EDM_TYPE` yang
// pesertanya ber-`PL_NUMBER` = nomor polis (`IDX_PLD_PL_NUMBER`, migrasi 480;
// nomor polis = `PL_NUMBER`, RALAT R01).
func sqlVersiNB(polis, peserta string, sebelum bool) string {
	syarat := ""
	if sebelum {
		syarat = " AND NVL(p.PROD_KE, 1) < :2"
	}
	// ⛔ Penampung muncul URUT (`:1` lalu `:2`): godror mengikat menurut urutan
	// kemunculan, bukan menurut angkanya.
	return fmt.Sprintf(`SELECT p.ID, NVL(p.PROD_KE, 1), p.EDM_TYPE FROM %s p
	  WHERE p.ID IN (SELECT d.PREMIUM_LIST_ID FROM %s d WHERE d.PL_NUMBER = :1)
	    AND p.EDM_TYPE IS NULL%s
	  ORDER BY NVL(p.PROD_KE, 1) DESC, p.ID FETCH FIRST 1 ROWS ONLY`, polis, peserta, syarat)
}

// sqlVersiWarisan - versi sistem lama: `GetProdkeNopolis` b84 (`IDPEGA … ORDER BY
// PRODKE DESC`) dan `GetEdmTypeLife` b84 (`A.DATA_JSON.EdmType`). `PRODKE`
// kosong dibaca 1 (E1, OQ-EDM-008). `JSON_VALUE` dipakai sebagai ganti notasi
// titik Pega, yang menuntut constraint `IS JSON` pada kolomnya.
func sqlVersiWarisan(jsonPolis string, sebelum bool) string {
	syarat := ""
	if sebelum {
		syarat = " AND NVL(j.PRODKE, 1) < :2"
	}
	return fmt.Sprintf(`SELECT j.IDPEGA, NVL(j.PRODKE, 1), JSON_VALUE(j.DATA_JSON, '$.EdmType' NULL ON ERROR)
	   FROM %s j WHERE j.NOPOLIS = :1%s
	  ORDER BY NVL(j.PRODKE, 1) DESC, j.TGL_INPUT DESC FETCH FIRST 1 ROWS ONLY`, jsonPolis, syarat)
}

// VersiBerjalan membaca versi berjalan sebuah polis - `PRODKE` terbesar di
// ketiga jalur. `sebelumProdKe` > 0: versi terbesar di bawah nilai itu (polis
// lama sebuah kasus yang sudah diresmikan). `ada` false = polis tidak dikenal.
func (g *Gudang) VersiBerjalan(ctx context.Context, tx *db.Tx, nomorPolis string, sebelumProdKe int) (models.Versi, bool, error) {
	n, err := g.nama(tabelPolis, tabelPeserta, tabelPolisWarisan)
	if err != nil {
		return models.Versi{}, false, err
	}
	sebelum := sebelumProdKe > 0
	type jalur struct {
		jenis models.JenisSumber
		q     string
		args  []any
	}
	argsEDM := []any{nomorPolis, models.StatusKasusSelesai}
	argsLain := []any{nomorPolis}
	if sebelum {
		argsEDM = append(argsEDM, sebelumProdKe)
		argsLain = append(argsLain, sebelumProdKe)
	}
	daftar := []jalur{
		{models.SumberAplikasi, sqlVersiEDM(n[0], sebelum), argsEDM},
		{models.SumberAplikasi, sqlVersiNB(n[0], n[1], sebelum), argsLain},
		{models.SumberWarisan, sqlVersiWarisan(n[2], sebelum), argsLain},
	}
	var terbaik models.Versi
	for _, j := range daftar {
		if err := db.PeriksaSQL(j.q); err != nil {
			return models.Versi{}, false, err
		}
		var id, edm sql.NullString
		var prodKe int
		err := g.pakai(tx).QueryRowContext(ctx, j.q, j.args...).Scan(&id, &prodKe, &edm)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return models.Versi{}, false, fmt.Errorf("repository: membaca versi %s polis %q: %w", j.jenis, nomorPolis, err)
		}
		terbaik = models.LebihBaru(terbaik, models.Versi{
			Jenis: j.jenis, ID: id.String, ProdKe: prodKe, EdmType: strings.TrimSpace(edm.String),
		})
	}
	return terbaik, terbaik.ID != "", nil
}
