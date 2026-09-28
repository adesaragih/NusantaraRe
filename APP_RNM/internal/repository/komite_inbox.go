package repository

// Inbox Komite dan kasus komite - tiket 01 Komite Claim Life.
//
// Untuk apa berkas ini: membaca kasus komite yang MENUNGGU pelaku, dan membaca
// satu kasus beserta tangganya. Nol tulisan - melahirkan kasus sudah milik
// `BuatKasusKomite` (kasuskomite.go, Claim Life A2), dan tidak disalin.
//
// `[terverifikasi]` `Komite Claim Life/Activity/KomiteRouter.xml` (router
// `Assignment1`, `pyRouteTo Custom` - `Flow/KomiteLife_Flow.xml` b619):
//
//	langkah 1     perulangan `.KomiteList`                    (b227)
//	langkah 1.1   `param.AssignTo = .KomiteID`                (b294-295)
//	              precondition `.KomiteAproval==0`            (b382, b396)
//	              transisi sesudah langkah: kode `6` (true/false)
//
// ⚠️ `[dugaan kuat]` kode transisi `6` = keluar dari perulangan, sehingga yang
// ditugasi adalah anggota PERTAMA ber-`KomiteAproval` 0 - bukan yang terakhir.
// Ekspor tidak menyebut arti kodenya. Yang menguatkan: `KomiteCount` dimulai
// 1 dan naik per putaran (`CreateKMTLife_Act` b1398, `KomiteLoop`), yaitu
// tingkat berjalan = anggota pertama yang belum memutuskan. Bila kode `6`
// ternyata "lanjut", worklist jatuh ke anggota TERAKHIR - dan itu harus
// dijawab pemilik ekspor Pega.
//
// ⛔ `KomiteID` = `KOMITE_OPERATORID` (roster.go: `.KomiteID = .OPERATOR_ID`,
// `CreateKMTLife_Act` b866/b972). Pelaku dicocokkan dengan kolom itu - km2:
// peran komite = roster, bukan `X-Peran` Claim Life.
//
// Dibaca sesudah: kasuskomite.go, roster.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/internal/models"
)

// ApprovalKomiteMenunggu adalah `KomiteAproval` anggota yang belum memutuskan.
//
// `[terverifikasi]` `CreateKMTLife_Act.xml` b912/b993 `= 0`; `KomiteRouter`
// b382 `.KomiteAproval==0`. Sama dengan `approvalAwal` (kasuskomite.go).
const ApprovalKomiteMenunggu = approvalAwal

// BarisInboxKomite adalah satu kasus di Inbox Komite.
type BarisInboxKomite struct {
	KasusID    string
	KlaimID    string
	NomorKlaim string
	// TingkatBerjalan adalah `KOMITE_URUT` anggota pertama yang menunggu.
	TingkatBerjalan int
	KomiteCount     int
	KomiteLoop      int
	// NilaiKlaim TEKS - uang tidak melewati float (ADR-U-0003).
	NilaiKlaim string
	MataUang   string
	StsReject  string
	StatusWork string
	TglUpdate  sql.NullTime
	// AcceptStatus = `T_GENERAL_KOMITE.ACCEPT_STATUS` - keputusan terakhir.
	AcceptStatus string
}

// kolomBarisKomite dipakai bersama inbox dan pembacaan satu kasus.
const kolomBarisKomite = `g.ID, w.COVER_KEY, c.CLAIM_NO, l.KOMITE_URUT,
	       g.KOMITE_COUNT, g.KOMITE_LOOP,
	       TO_CHAR(a.CLAIM_AMOUNT, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''),
	       a.CURRENCY, a.STS_REJECT, w.STATUS_WORK, w.TGL_UPDATE`

// sqlDariKomite adalah FROM bersama - kasus, work, klaim induk, baris adjustment.
//
// ⚠️ `LEFT JOIN` ke klaim induk dan baris adjustment: kasus yang induknya
// hilang tetap harus TERLIHAT di inbox, bukan lenyap diam-diam.
func sqlDariKomite(gen, work, list, klaim, adj string) string {
	return fmt.Sprintf(` FROM %s g
	  JOIN %s w ON w.ID = g.ID
	  JOIN %s l ON l.DATA_KOMITE_ID = g.ID
	  LEFT JOIN %s c ON c.ID = w.COVER_KEY
	  LEFT JOIN %s a ON a.ID = g.ADJUSTMENT_ID`, gen, work, list, klaim, adj)
}

// sqlSaringInboxKomite menyaring ke anggota BERJALAN yang adalah pelaku.
//
// ⛔ "Berjalan" = `KOMITE_URUT` TERKECIL yang masih ber-approval menunggu
// (lihat kepala berkas). Anggota tingkat 3 tidak melihat kasus yang tingkat
// 1-nya belum memutuskan - AC 10 spec, ADR-0014.
//
// ⛔ Kasus tertutup (`STATUS_WORK` Resolved-Completed) tidak berdiri di
// worklist siapa pun, walau barisnya masih menyimpan approval menunggu.
//
// ⛔ TIKET 02 - dan kasus yang tangganya BERHENTI pun tidak. Tolak di tingkat
// tengah meninggalkan anggota berikutnya ber-approval menunggu; tanpa syarat
// `IsKomiteLoop` (`models.KasusDiTangga`) kasus itu jatuh ke inbox mereka.
// `End1` tidak punya `pyWorkStatus`, jadi syaratnya dibaca dari kepala kasus.
const sqlSaringInboxKomite = `
	 WHERE l.KOMITE_OPERATORID = :akun
	   AND l.KOMITE_APPROVAL = :menunggu
	   AND l.KOMITE_URUT = (SELECT MIN(l2.KOMITE_URUT) FROM %s l2
	                         WHERE l2.DATA_KOMITE_ID = g.ID
	                           AND l2.KOMITE_APPROVAL = :menunggu)
	   AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK <> :tutup)
	   AND (g.ACCEPT_STATUS IS NULL
	        OR (g.ACCEPT_STATUS = :setuju AND g.KOMITE_COUNT <= g.KOMITE_LOOP))`

// InboxKomite membaca kasus komite.
type InboxKomite struct{ db *DB }

// NewInboxKomite menyusunnya.
func NewInboxKomite(db *DB) *InboxKomite { return &InboxKomite{db: db} }

// tabelKomite meng-qualify kelima tabel sekaligus.
func (r *InboxKomite) tabelKomite() (gen, work, list, klaim, adj string, err error) {
	for _, p := range []struct {
		nama string
		ke   *string
	}{
		{"T_GENERAL_KOMITE", &gen}, {"T_WORK_CLAIM", &work},
		{"T_KOMITE_KOMITELIST", &list}, {"T_GENERAL_CLAIM", &klaim},
		{"T_CLAIMLF_ADJUSTMENT", &adj},
	} {
		if *p.ke, err = r.db.Qualify(p.nama); err != nil {
			return
		}
	}
	return
}

// sqlInboxKomite merakit pembacaan satu halaman inbox.
func sqlInboxKomite(gen, work, list, klaim, adj string) string {
	return `SELECT ` + kolomBarisKomite + sqlDariKomite(gen, work, list, klaim, adj) +
		fmt.Sprintf(sqlSaringInboxKomite, list) + `
	 ORDER BY w.TGL_UPDATE DESC, g.ID DESC
	 OFFSET :offset ROWS FETCH NEXT :ukuran ROWS ONLY`
}

// sqlCacahInboxKomite merakit pencacahan totalnya - penyaring yang SAMA.
func sqlCacahInboxKomite(gen, work, list, klaim, adj string) string {
	return `SELECT COUNT(*)` + sqlDariKomite(gen, work, list, klaim, adj) +
		fmt.Sprintf(sqlSaringInboxKomite, list)
}

// pindaiBarisKomite memindai satu baris `kolomBarisKomite`.
func pindaiBarisKomite(rows interface{ Scan(...any) error }) (BarisInboxKomite, error) {
	var b BarisInboxKomite
	var klaimID, nomor, nilai, mu, sts, status sql.NullString
	var urut, count, loop sql.NullInt64
	if err := rows.Scan(&b.KasusID, &klaimID, &nomor, &urut, &count, &loop,
		&nilai, &mu, &sts, &status, &b.TglUpdate); err != nil {
		return b, err
	}
	b.KlaimID, b.NomorKlaim = klaimID.String, nomor.String
	b.TingkatBerjalan, b.KomiteCount, b.KomiteLoop = int(urut.Int64), int(count.Int64), int(loop.Int64)
	b.NilaiKlaim, b.MataUang = strings.TrimSpace(nilai.String), strings.TrimSpace(mu.String)
	b.StsReject, b.StatusWork = strings.TrimSpace(sts.String), status.String
	return b, nil
}

// Ambil membaca satu halaman inbox pelaku beserta cacah totalnya.
func (r *InboxKomite) Ambil(ctx context.Context, akunID, statusTutup string,
	offset, ukuran int) ([]BarisInboxKomite, int, error) {

	if strings.TrimSpace(akunID) == "" {
		return nil, 0, errors.New("repository: inbox komite tanpa akun pelaku")
	}
	gen, work, list, klaim, adj, err := r.tabelKomite()
	if err != nil {
		return nil, 0, err
	}
	q := sqlInboxKomite(gen, work, list, klaim, adj)
	if err := PeriksaSQL(q); err != nil {
		return nil, 0, err
	}
	// ⛔ Urutan argumen = urutan MUNCULNYA penanda di teks - driver Oracle di
	// jalur ini mengikat penanda bernama secara berurutan (pola `AmbilInbox`
	// Claim Life). `:menunggu` muncul DUA kali, jadi ia dikirim dua kali.
	saring := []any{akunID, ApprovalKomiteMenunggu, ApprovalKomiteMenunggu, statusTutup,
		models.KeputusanKomiteSetuju}
	rows, err := r.db.sql.QueryContext(ctx, q, append(saring, offset, ukuran)...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: membaca inbox komite: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var keluar []BarisInboxKomite
	for rows.Next() {
		b, err := pindaiBarisKomite(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("repository: memindai inbox komite: %w", err)
		}
		keluar = append(keluar, b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository: membaca inbox komite: %w", err)
	}

	qc := sqlCacahInboxKomite(gen, work, list, klaim, adj)
	if err := PeriksaSQL(qc); err != nil {
		return nil, 0, err
	}
	var total int
	if err := r.db.sql.QueryRowContext(ctx, qc, saring...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: mencacah inbox komite: %w", err)
	}
	return keluar, total, nil
}

// AnggotaKasus adalah satu anak tangga sebuah kasus.
//
// ⚠️ `OperatorID` DATA ORANG (pengenal akun). `Email` sengaja TIDAK dibaca:
// layar tidak memerlukannya, dan alamat orang yang tidak dipakai adalah
// alamat orang yang suatu hari tercetak di tempat yang salah.
type AnggotaKasus struct {
	Urut       int
	OperatorID string
	Jabatan    string
	Approval   string
	Komentar   string
	TglAprove  sql.NullTime
}

// KasusKomite adalah satu kasus beserta tangganya.
type KasusKomite struct {
	Baris  BarisInboxKomite
	Tangga []AnggotaKasus
	AdjID  string
	// PesertaID = `T_CLAIMLF_ADJUSTMENT.PREMIUM_LIST_DETAIL_ID` - peserta
	// pemilik baris yang diputuskan (tiket 04a menulis statusnya).
	PesertaID string
	// IsKPR = `T_GENERAL_CLAIM.IS_KPR` klaim induk - gerbang Kasir langkah 12
	// (`TempOpenPage.ClaimData.IsKPR=="KPR"`).
	IsKPR   string
	Ditemui bool
}

// ErrKasusKomiteTakDitemukan - tidak ada kasus komite dengan id itu.
var ErrKasusKomiteTakDitemukan = errors.New("repository: kasus komite tidak ditemukan")

// sqlKasusKomite membaca kepala satu kasus. Tingkat berjalannya dihitung
// seperti inbox; kasus yang tidak lagi menunggu siapa pun memberi urut 0.
func sqlKasusKomite(gen, work, list, klaim, adj string) string {
	return fmt.Sprintf(`SELECT g.ID, w.COVER_KEY, c.CLAIM_NO,
	       (SELECT MIN(l2.KOMITE_URUT) FROM %s l2
	         WHERE l2.DATA_KOMITE_ID = g.ID AND l2.KOMITE_APPROVAL = :1),
	       g.KOMITE_COUNT, g.KOMITE_LOOP,
	       TO_CHAR(a.CLAIM_AMOUNT, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''),
	       a.CURRENCY, a.STS_REJECT, w.STATUS_WORK, w.TGL_UPDATE, g.ADJUSTMENT_ID,
	       g.ACCEPT_STATUS, a.PREMIUM_LIST_DETAIL_ID, c.IS_KPR
	  FROM %s g
	  JOIN %s w ON w.ID = g.ID
	  LEFT JOIN %s c ON c.ID = w.COVER_KEY
	  LEFT JOIN %s a ON a.ID = g.ADJUSTMENT_ID
	 WHERE g.ID = :2`, list, gen, work, klaim, adj)
}

// sqlTanggaKasus membaca tangga satu kasus, urut jenjang.
func sqlTanggaKasus(list string) string {
	return fmt.Sprintf(`SELECT KOMITE_URUT, KOMITE_OPERATORID, KOMITE_JABATAN,
	       KOMITE_APPROVAL, KOMITE_COMMENT, DATE_APPROVE
	  FROM %s WHERE DATA_KOMITE_ID = :1 ORDER BY KOMITE_URUT, ID`, list)
}

// Kasus membaca satu kasus komite beserta tangganya.
func (r *InboxKomite) Kasus(ctx context.Context, kasusID string) (KasusKomite, error) {
	gen, work, list, klaim, adj, err := r.tabelKomite()
	if err != nil {
		return KasusKomite{}, err
	}
	q := sqlKasusKomite(gen, work, list, klaim, adj)
	if err := PeriksaSQL(q); err != nil {
		return KasusKomite{}, err
	}
	var k KasusKomite
	var klaimID, nomor, nilai, mu, sts, status, adjID, accept, peserta, kpr sql.NullString
	var urut, count, loop sql.NullInt64
	err = r.db.sql.QueryRowContext(ctx, q, ApprovalKomiteMenunggu, kasusID).Scan(
		&k.Baris.KasusID, &klaimID, &nomor, &urut, &count, &loop, &nilai, &mu, &sts,
		&status, &k.Baris.TglUpdate, &adjID, &accept, &peserta, &kpr)
	if errors.Is(err, sql.ErrNoRows) {
		return KasusKomite{}, fmt.Errorf("%w: %q", ErrKasusKomiteTakDitemukan, kasusID)
	}
	if err != nil {
		return KasusKomite{}, fmt.Errorf("repository: membaca kasus komite: %w", err)
	}
	k.Ditemui = true
	k.Baris.KlaimID, k.Baris.NomorKlaim = klaimID.String, nomor.String
	k.Baris.TingkatBerjalan = int(urut.Int64)
	k.Baris.KomiteCount, k.Baris.KomiteLoop = int(count.Int64), int(loop.Int64)
	k.Baris.NilaiKlaim, k.Baris.MataUang = strings.TrimSpace(nilai.String), strings.TrimSpace(mu.String)
	k.Baris.StsReject, k.Baris.StatusWork = strings.TrimSpace(sts.String), status.String
	k.AdjID = adjID.String
	k.Baris.AcceptStatus = strings.TrimSpace(accept.String)
	k.PesertaID = peserta.String
	k.IsKPR = strings.TrimSpace(kpr.String)

	qt := sqlTanggaKasus(list)
	if err := PeriksaSQL(qt); err != nil {
		return KasusKomite{}, err
	}
	rows, err := r.db.sql.QueryContext(ctx, qt, kasusID)
	if err != nil {
		return KasusKomite{}, fmt.Errorf("repository: membaca tangga komite: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var a AnggotaKasus
		var op, jab, appr, kom sql.NullString
		var u sql.NullInt64
		if err := rows.Scan(&u, &op, &jab, &appr, &kom, &a.TglAprove); err != nil {
			return KasusKomite{}, fmt.Errorf("repository: memindai tangga komite: %w", err)
		}
		a.Urut, a.OperatorID, a.Jabatan = int(u.Int64), op.String, jab.String
		a.Approval, a.Komentar = strings.TrimSpace(appr.String), kom.String
		k.Tangga = append(k.Tangga, a)
	}
	if err := rows.Err(); err != nil {
		return KasusKomite{}, fmt.Errorf("repository: membaca tangga komite: %w", err)
	}
	return k, nil
}
