package repository

// Spreading peserta ke retro - Calculate CSV Type QR (keputusan work owner 06-10-2026).
// Aturannya murni di models/polis_spreading.go.
//
// Untuk apa berkas ini:
//   - MEMBACA bahan Contract Retro Life (TREATYYEAR_LIFE, TREATYCONTRACT_LIFE,
//     TREATYREINSURER_LIFE, TREATYBUSINESS_LIFE - milik modul Master Contract Retro
//     Life, dibaca saja, modulnya TIDAK diimpor) dan BEGIN_DATE produk;
//   - MENULIS T_PREMIUM_LIST_SPREADING (053) dan T_PREMIUM_LIST_SPREADING_RETRO (054)
//     di transaksi yang sama dengan sisip peserta;
//   - MEMBACA hasilnya untuk popup spreading satu peserta (07-10-2026) dan rekap retro tab Summary.
//
// ⛔ Kolom disebut satu per satu, angka lewat db.FmtDesimal (teks desimal), tanggal
// TO_CHAR 'YYYY-MM-DD'; nol COMMIT (db.PeriksaSQL).

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
)

// Tabel Contract Retro Life yang dibaca (nama = modul mastercontractretrolife, repository/mcrl_tabel.go).
const (
	TabelTahunRetro     = "TREATYYEAR_LIFE"
	TabelKontrakRetro   = "TREATYCONTRACT_LIFE"
	TabelReinsurerRetro = "TREATYREINSURER_LIFE"
	TabelBusinessRetro  = "TREATYBUSINESS_LIFE"
	TabelSpreading      = "T_PREMIUM_LIST_SPREADING"
	TabelSpreadingRetro = "T_PREMIUM_LIST_SPREADING_RETRO"
)

const fmtTanggalISO = `TO_CHAR(%s, 'YYYY-MM-DD')`

func des(k string) string { return fmt.Sprintf(db.FmtDesimal, k) }
func tgl(k string) string { return fmt.Sprintf(fmtTanggalISO, k) }

// ——— baca bahan ———

func sqlBeginProduk(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, tgl("BEGIN_DATE"), t)
}

func sqlTahunRetro(t string) string {
	return fmt.Sprintf(`SELECT ID, %s, %s FROM %s ORDER BY ID`, tgl("STARTDATE"), tgl("ENDDATE"), t)
}

func sqlKontrakRetro(t string) string {
	return fmt.Sprintf(`SELECT ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME, %s, %s, %s, %s, %s, %s, %s, %s
	   FROM %s WHERE IDTREATYYEAR = :1 ORDER BY ID`,
		tgl("TREATYSTARTDATE"), tgl("TREATYENDDATE"), des("IDR"), des("USD"), des("B_IDR"), des("B_USD"),
		des("IDR_SELISIH"), des("USD_SELISIH"), t)
}

func sqlReinsurerRetro(t string) string {
	return fmt.Sprintf(`SELECT TREATYCONTRACTID, REINSURERNAME, %s, %s, %s FROM %s WHERE TREATYYEARID = :1 ORDER BY ID`,
		des("PCTSHARE"), des("COMMISION"), des("OVR_COMM"), t)
}

func sqlBusinessRetro(t string) string {
	return fmt.Sprintf(`SELECT TREATYCONTRACTID, BIZNAME, RIRATEID, RIRATE FROM %s WHERE TREATYYEARID = :1 ORDER BY ID`, t)
}

// BeginDateProduk - BEGIN_DATE produk (YYYY-MM-DD) dari M_PRODUCTNAME_LIFE; kosong bila tidak terisi.
func (r *Rujukan) BeginDateProduk(ctx context.Context, produkID string) (string, error) {
	tabel, err := r.db.Qualify(TabelProdukMPNL)
	if err != nil {
		return "", err
	}
	q := sqlBeginProduk(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var v sql.NullString
	err = r.db.QueryRowContext(ctx, q, produkID).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%w: %s", ErrRincianProdukTidakAda, produkID)
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca BEGIN_DATE %s: %w", TabelProdukMPNL, err)
	}
	return strings.TrimSpace(v.String), nil
}

// bacaBaris - satu kueri, setiap baris dipindai ke n teks (NULL → "").
func (r *Rujukan) bacaBaris(ctx context.Context, tabelLogis, q string, n int, args ...any) ([][]string, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", tabelLogis, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil [][]string
	for rows.Next() {
		sel := make([]sql.NullString, n)
		tujuan := make([]any, n)
		for i := range sel {
			tujuan[i] = &sel[i]
		}
		if err := rows.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: memindai %s: %w", tabelLogis, err)
		}
		b := make([]string, n)
		for i, s := range sel {
			b[i] = strings.TrimSpace(s.String)
		}
		hasil = append(hasil, b)
	}
	return hasil, rows.Err()
}

// DaftarTahunRetro - seluruh TREATYYEAR_LIFE (pemilihan S1 di models.PilihTahunRetro).
func (r *Rujukan) DaftarTahunRetro(ctx context.Context) ([]models.TahunRetro, error) {
	t, err := r.db.Qualify(TabelTahunRetro)
	if err != nil {
		return nil, err
	}
	rows, err := r.bacaBaris(ctx, TabelTahunRetro, sqlTahunRetro(t), 3)
	if err != nil {
		return nil, err
	}
	hasil := make([]models.TahunRetro, 0, len(rows))
	for _, b := range rows {
		hasil = append(hasil, models.TahunRetro{ID: b[0], Mulai: b[1], Akhir: b[2]})
	}
	return hasil, nil
}

// BahanSpreadingTahun - kontrak, reinsurer, dan business satu treaty year.
func (r *Rujukan) BahanSpreadingTahun(ctx context.Context, tahunID string) (
	[]models.KontrakRetro, []models.ReinsurerRetro, []models.BusinessRetro, error) {

	tk, err := r.db.Qualify(TabelKontrakRetro)
	if err != nil {
		return nil, nil, nil, err
	}
	tr, err := r.db.Qualify(TabelReinsurerRetro)
	if err != nil {
		return nil, nil, nil, err
	}
	tb, err := r.db.Qualify(TabelBusinessRetro)
	if err != nil {
		return nil, nil, nil, err
	}
	rk, err := r.bacaBaris(ctx, TabelKontrakRetro, sqlKontrakRetro(tk), 12, tahunID)
	if err != nil {
		return nil, nil, nil, err
	}
	var kontrak []models.KontrakRetro
	for _, b := range rk {
		kontrak = append(kontrak, models.KontrakRetro{ID: b[0], IDTahun: b[1], ReinsTypeID: b[2], ReinsTypeName: b[3],
			Mulai: b[4], Akhir: b[5], IDR: angkaUtuh(b[6]), USD: angkaUtuh(b[7]), BIDR: angkaUtuh(b[8]),
			BUSD: angkaUtuh(b[9]), IDRSelisih: angkaUtuh(b[10]), USDSelisih: angkaUtuh(b[11])})
	}
	rr, err := r.bacaBaris(ctx, TabelReinsurerRetro, sqlReinsurerRetro(tr), 5, tahunID)
	if err != nil {
		return nil, nil, nil, err
	}
	var reins []models.ReinsurerRetro
	for _, b := range rr {
		reins = append(reins, models.ReinsurerRetro{KontrakID: b[0], Nama: b[1], PctShare: angkaUtuh(b[2]),
			Komisi: angkaUtuh(b[3]), OvrComm: angkaUtuh(b[4])})
	}
	rb, err := r.bacaBaris(ctx, TabelBusinessRetro, sqlBusinessRetro(tb), 4, tahunID)
	if err != nil {
		return nil, nil, nil, err
	}
	var biz []models.BusinessRetro
	for _, b := range rb {
		biz = append(biz, models.BusinessRetro{KontrakID: b[0], BizName: b[1], RIRateID: b[2], RIRate: b[3]})
	}
	return kontrak, reins, biz, nil
}

// ——— tulis ———

// PengenalSpreading - ID baris T_PREMIUM_LIST_SPREADING: pola PengenalPesertaUnggah
// (sidik polis + nomor baris) ditambah urutan lapisan. 32 karakter.
func PengenalSpreading(polisID string, nomorBaris, lapisan int) string {
	sum := md5.Sum([]byte(polisID + "\x00spreading"))
	return strings.ToUpper(hex.EncodeToString(sum[:]))[:20] + fmt.Sprintf("%08d%04d", nomorBaris, lapisan)
}

// PengenalSpreadingRetro - ID baris T_PREMIUM_LIST_SPREADING_RETRO. 32 karakter.
func PengenalSpreadingRetro(polisID string, nomorBaris, lapisan, reinsurer int) string {
	sum := md5.Sum([]byte(polisID + "\x00retro"))
	return strings.ToUpper(hex.EncodeToString(sum[:]))[:16] + fmt.Sprintf("%08d%04d%04d", nomorBaris, lapisan, reinsurer)
}

func sqlHapusSpreadingRetro(retro, spreading, detail string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE SPREADING_ID IN (SELECT s.ID FROM %s s JOIN %s d ON d.ID = s.DETAIL_ID
	   WHERE d.PREMIUM_LIST_ID = :1)`, retro, spreading, detail)
}

func sqlHapusSpreading(spreading, detail string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE DETAIL_ID IN (SELECT d.ID FROM %s d WHERE d.PREMIUM_LIST_ID = :1)`, spreading, detail)
}

func sqlSisipSpreading(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, DETAIL_ID, TREATY_TYPE_ID, TREATY_TYPE_NAME, TREATY_YEAR_LIFE, RETROCADED_SHARE,
	   IDR, USD, IDR_SELISIH, USD_SELISIH, B_IDR, B_USD, TGL_UPDATE, USER_ID)
	   VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, SYSDATE, :13)`, t)
}

func sqlSisipSpreadingRetro(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, SPREADING_ID, REINSURER_NAME, PERCENT_SHARE, AMOUNT, COMMISION, OVR_COMM, RATE,
	   PREMIUM_SPREADED_GROSS, PREMIUM_SPREADED_NET, TREATY_TYPE_ID, TREATY_TYPE_NAME, TREATY_START_DATE,
	   TREATY_END_DATE, TGL_UPDATE, USER_ID)
	   VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, TO_DATE(:13, 'YYYY-MM-DD'),
	   TO_DATE(:14, 'YYYY-MM-DD'), SYSDATE, :15)`, t)
}

// Spreading menulis dan membaca hasil spreading satu polis.
type Spreading struct{ db *db.DB }

// NewSpreading menyusunnya.
func NewSpreading(db *db.DB) *Spreading { return &Spreading{db: db} }

func (r *Spreading) tabel() (spreading, retro, detail string, err error) {
	if spreading, err = r.db.Qualify(TabelSpreading); err != nil {
		return
	}
	if retro, err = r.db.Qualify(TabelSpreadingRetro); err != nil {
		return
	}
	detail, err = r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	return
}

// HapusSpreadingPolis menghapus spreading (dan retro-nya) seluruh peserta satu polis,
// EKSPLISIT - tidak bergantung ON DELETE CASCADE (053/054) - sebelum peserta dihapus.
func (r *Spreading) HapusSpreadingPolis(ctx context.Context, tx *db.Tx, polisID string) error {
	spreading, retro, detail, err := r.tabel()
	if err != nil {
		return err
	}
	for _, q := range []string{sqlHapusSpreadingRetro(retro, spreading, detail), sqlHapusSpreading(spreading, detail)} {
		if err := db.PeriksaSQL(q); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, q, polisID); err != nil {
			return fmt.Errorf("repository: menghapus spreading polis: %w", err)
		}
	}
	return nil
}

// nilaiOpsional - teks kosong → NULL.
func nilaiOpsional(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

// NilaiSisipSpreading - argumen sisip satu lapisan (urut sqlSisipSpreading).
func NilaiSisipSpreading(id, detailID, userID string, l models.LapisanSpreading) []any {
	k := l.Kontrak
	return []any{id, detailID, k.ReinsTypeID, k.ReinsTypeName, l.TahunID, models.TeksSimpan(l.Bagian),
		nilaiOpsional(k.IDR), nilaiOpsional(k.USD), nilaiOpsional(k.IDRSelisih), nilaiOpsional(k.USDSelisih),
		nilaiOpsional(k.BIDR), nilaiOpsional(k.BUSD), userID}
}

// NilaiSisipSpreadingRetro - argumen sisip satu reinsurer (urut sqlSisipSpreadingRetro).
func NilaiSisipSpreadingRetro(id, spreadingID, userID string, l models.LapisanSpreading, x models.ReinsurerSpreading) []any {
	return []any{id, spreadingID, x.Nama, models.TeksSimpan(x.PctShare), models.TeksSimpan(x.Amount),
		models.TeksSimpan(x.Komisi), models.TeksSimpan(x.Ovr), models.TeksSimpan(x.Rate), models.TeksSimpan(x.Gross),
		models.TeksSimpan(x.Net), l.Kontrak.ReinsTypeID, l.Kontrak.ReinsTypeName, nilaiOpsional(l.Kontrak.Mulai),
		nilaiOpsional(l.Kontrak.Akhir), userID}
}

// SisipSpreading menulis spreading setiap baris (\`b.Spreading\`) dengan DETAIL_ID dari
// \`idPeserta\` (nomor baris → ID T_PREMIUM_LIST_DETAIL). Mengembalikan cacah lapisan dan reinsurer.
func (r *Spreading) SisipSpreading(ctx context.Context, tx *db.Tx, polisID, userID string, idPeserta map[int]string,
	baris []models.BarisUnggah) (int, int, error) {

	spreading, retro, _, err := r.tabel()
	if err != nil {
		return 0, 0, err
	}
	qs, qr := sqlSisipSpreading(spreading), sqlSisipSpreadingRetro(retro)
	for _, q := range []string{qs, qr} {
		if err := db.PeriksaSQL(q); err != nil {
			return 0, 0, err
		}
	}
	var nLap, nRe int
	for _, b := range baris {
		if len(b.Spreading) == 0 {
			continue
		}
		detailID, ada := idPeserta[b.Nomor]
		if !ada {
			return 0, 0, fmt.Errorf("repository: ID peserta baris %d tidak ada untuk spreading", b.Nomor)
		}
		for i, l := range b.Spreading {
			sid := PengenalSpreading(polisID, b.Nomor, i+1)
			if _, err := tx.ExecContext(ctx, qs, NilaiSisipSpreading(sid, detailID, userID, l)...); err != nil {
				return 0, 0, fmt.Errorf("repository: menyisipkan spreading baris %d: %w", b.Nomor, err)
			}
			nLap++
			for j, x := range l.Reinsurer {
				rid := PengenalSpreadingRetro(polisID, b.Nomor, i+1, j+1)
				if _, err := tx.ExecContext(ctx, qr, NilaiSisipSpreadingRetro(rid, sid, userID, l, x)...); err != nil {
					return 0, 0, fmt.Errorf("repository: menyisipkan spreading retro baris %d: %w", b.Nomor, err)
				}
				nRe++
			}
		}
	}
	return nLap, nRe, nil
}

// ——— baca hasil (popup spreading satu peserta, rekap Summary) ———

// SpreadingReinsurer - satu reinsurer satu lapisan, teks desimal dari Oracle.
type SpreadingReinsurer struct {
	Nama     string `json:"nama"`
	PctShare string `json:"pctShare"`
	Amount   string `json:"amount"`
	Rate     string `json:"rate"`
	Komisi   string `json:"komisi"`
	OvrComm  string `json:"ovrComm"`
	Gross    string `json:"gross"`
	Net      string `json:"net"`
}

// SpreadingLapisan - satu lapisan satu peserta.
type SpreadingLapisan struct {
	ReinsType  string               `json:"reinsType"`
	TreatyYear string               `json:"treatyYear"`
	Bagian     string               `json:"bagian"`
	Reinsurer  []SpreadingReinsurer `json:"reinsurer"`
}

// SpreadingSatuPeserta - popup spreading satu baris peserta (keputusan work owner
// 07-10-2026). Tanpa spreading (bukan Type QR, atau belum dihitung) = Lapisan kosong.
type SpreadingSatuPeserta struct {
	ID              string `json:"id"`
	Nama            string `json:"nama"`
	Sertifikat      string `json:"sertifikat"`
	Currency        string `json:"currency"`
	RetrocededShare string `json:"retrocededShare"`
	// TreatyYear - TREATYYEAR_LIFE lapisan pertama (satu treaty year per polis):
	// label Transaction Year bila ada, selain itu ID-nya; periode ISO YYYY-MM-DD.
	TreatyYear  string             `json:"treatyYear"`
	TreatyMulai string             `json:"treatyMulai"`
	TreatyAkhir string             `json:"treatyAkhir"`
	Lapisan     []SpreadingLapisan `json:"lapisan"`
}

// ErrPesertaSpreadingTidakAda - DETAIL_ID bukan peserta polis itu (404).
var ErrPesertaSpreadingTidakAda = errors.New("repository: peserta tidak ada di polis ini")

// RekapRetro - total retro per reinsurer per mata uang (nama dari data, tanpa daftar tetap).
type RekapRetro struct {
	Reinsurer string `json:"reinsurer"`
	Currency  string `json:"currency"`
	Gross     string `json:"gross"`
	Net       string `json:"net"`
}

// sqlIdentitasPesertaSpreading - kepala popup; ⛔ peserta HARUS milik polis :1.
func sqlIdentitasPesertaSpreading(detail string) string {
	return fmt.Sprintf(`SELECT d.ID, d.NAME_OF_INSURED, d.CERTIFICATE_NO, d.CURRENCY, %s FROM %s d
	   WHERE d.PREMIUM_LIST_ID = :1 AND d.ID = :2`, des("d.RETROCEDED_SHARE"), detail)
}

// sqlRincianSpreadingPeserta - lapisan → reinsurer satu peserta. Urut ID: pengenal
// lapisan/reinsurer berakhiran nomor urut berlebar tetap (PengenalSpreading*),
// jadi urut ID = urut lapisan lalu urut reinsurer saat dihitung.
func sqlRincianSpreadingPeserta(spreading, retro, tahun string) string {
	return fmt.Sprintf(`SELECT s.ID, s.TREATY_TYPE_NAME, s.TREATY_YEAR_LIFE, t.TREATYYEAR, %s, %s, %s,
	   r.REINSURER_NAME, %s, %s, %s, %s, %s, %s, %s
	   FROM %s s JOIN %s r ON r.SPREADING_ID = s.ID LEFT JOIN %s t ON TO_CHAR(t.ID) = s.TREATY_YEAR_LIFE
	   WHERE s.DETAIL_ID = :1
	   ORDER BY s.ID, r.ID`,
		tgl("t.STARTDATE"), tgl("t.ENDDATE"), des("s.RETROCADED_SHARE"), des("r.PERCENT_SHARE"), des("r.AMOUNT"),
		des("r.RATE"), des("r.COMMISION"), des("r.OVR_COMM"), des("r.PREMIUM_SPREADED_GROSS"),
		des("r.PREMIUM_SPREADED_NET"), spreading, retro, tahun)
}

func sqlRekapRetro(detail, spreading, retro string) string {
	return fmt.Sprintf(`SELECT r.REINSURER_NAME, d.CURRENCY, %s, %s
	   FROM %s r JOIN %s s ON s.ID = r.SPREADING_ID JOIN %s d ON d.ID = s.DETAIL_ID
	   WHERE d.PREMIUM_LIST_ID = :1
	   GROUP BY r.REINSURER_NAME, d.CURRENCY
	   ORDER BY d.CURRENCY, r.REINSURER_NAME`,
		des("SUM(r.PREMIUM_SPREADED_GROSS)"), des("SUM(r.PREMIUM_SPREADED_NET)"), retro, spreading, detail)
}

func (r *Spreading) baca(ctx context.Context, q string, n int, args ...any) ([][]string, error) {
	return (&Rujukan{db: r.db}).bacaBaris(ctx, "spreading", q, n, args...)
}

// SusunSpreadingPeserta mengelompokkan baris sqlRincianSpreadingPeserta (urut) ke
// lapisan → reinsurer dan mengisi treaty year dari baris pertama. Nol baris → Lapisan kosong.
func SusunSpreadingPeserta(p *SpreadingSatuPeserta, rows [][]string) {
	p.Lapisan = []SpreadingLapisan{}
	lapisanTerakhir := ""
	for i, b := range rows {
		if i == 0 {
			p.TreatyYear, p.TreatyMulai, p.TreatyAkhir = b[2], b[4], b[5]
			if b[3] != "" {
				p.TreatyYear = b[3]
			}
		}
		if b[0] != lapisanTerakhir || len(p.Lapisan) == 0 {
			lapisanTerakhir = b[0]
			p.Lapisan = append(p.Lapisan, SpreadingLapisan{ReinsType: b[1], TreatyYear: b[2], Bagian: angkaUtuh(b[6]),
				Reinsurer: []SpreadingReinsurer{}})
		}
		l := &p.Lapisan[len(p.Lapisan)-1]
		l.Reinsurer = append(l.Reinsurer, SpreadingReinsurer{Nama: b[7], PctShare: angkaUtuh(b[8]), Amount: angkaUtuh(b[9]),
			Rate: angkaUtuh(b[10]), Komisi: angkaUtuh(b[11]), OvrComm: angkaUtuh(b[12]), Gross: angkaUtuh(b[13]),
			Net: angkaUtuh(b[14])})
	}
}

// Peserta membaca spreading retro SATU peserta polis (popup baris peserta, 07-10-2026).
// ErrPesertaSpreadingTidakAda bila detailID bukan peserta polisID.
func (r *Spreading) Peserta(ctx context.Context, polisID, detailID string) (SpreadingSatuPeserta, error) {
	spreading, retro, detail, err := r.tabel()
	if err != nil {
		return SpreadingSatuPeserta{}, err
	}
	tahun, err := r.db.Qualify(TabelTahunRetro)
	if err != nil {
		return SpreadingSatuPeserta{}, err
	}
	rp, err := r.baca(ctx, sqlIdentitasPesertaSpreading(detail), 5, polisID, detailID)
	if err != nil {
		return SpreadingSatuPeserta{}, err
	}
	if len(rp) == 0 {
		return SpreadingSatuPeserta{}, ErrPesertaSpreadingTidakAda
	}
	p := SpreadingSatuPeserta{ID: rp[0][0], Nama: rp[0][1], Sertifikat: rp[0][2], Currency: rp[0][3],
		RetrocededShare: angkaUtuh(rp[0][4])}
	rr, err := r.baca(ctx, sqlRincianSpreadingPeserta(spreading, retro, tahun), 15, detailID)
	if err != nil {
		return SpreadingSatuPeserta{}, err
	}
	SusunSpreadingPeserta(&p, rr)
	return p, nil
}

// Rekap membaca total retro per reinsurer per mata uang satu polis.
func (r *Spreading) Rekap(ctx context.Context, polisID string) ([]RekapRetro, error) {
	spreading, retro, detail, err := r.tabel()
	if err != nil {
		return nil, err
	}
	rows, err := r.baca(ctx, sqlRekapRetro(detail, spreading, retro), 4, polisID)
	if err != nil {
		return nil, err
	}
	hasil := []RekapRetro{}
	for _, b := range rows {
		hasil = append(hasil, RekapRetro{Reinsurer: b[0], Currency: b[1], Gross: angkaUtuh(b[2]), Net: angkaUtuh(b[3])})
	}
	return hasil, nil
}
