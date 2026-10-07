package repository

// Untuk apa berkas ini: KASUS - baris `T_WORK_POLIS` (tabel kasus lintas-lini, premiumlistlife 050/059) dan baris
// generasi `T_GENERAL_POLIS_TREATY` (milik nbtreatyin, 320) yang berbagi kunci utama dengannya. Padanan work object
// `ASM-FW-GISFW-Work-EndorsementTreaty`: `CreateEDMT` (CreateWorkPage + AddWork), perpindahan assignment
// `Flow/InputAddendumTreatyIn`, dan daftar portal `Section/SFAPortal_Endorsement_Treaty` (RD `InboxEDM_RD2`).
// Pola dan sebagian SQL: salinan `modul/nbtreatyin/backend/repository/kasus.go` (06-10-2026).
//
// ⛔ Setiap ubah tahap MENYERTAKAN tahap lama di WHERE: dua permintaan serentak tidak dapat sama-sama memindahkan
// kasus yang sama.
//
// ⛔ Generasi endorsemen: `PRODKE >= 1` (daftar NB menyaring `PRODKE = 0` - kedua modul tidak saling melihat
// kasusnya). NOPOLIS generasi endorsemen KOSONG selama berjalan dan diisi nomor polis induk saat selesai
// (`SetelNomorPolisSelesai`): indeks unik berfungsi `UQ_GP_TREATY_NOPOLIS` (320) hanya memuat baris bernomor,
// sehingga endorsemen yang ditolak tidak menahan nomor generasi (Pega: ProdKe = PRODKE json_polis terbesar + 1 -
// json_polis hanya ditulis Utility1 sesudah disetujui).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
)

// IDKasusBerikut menerbitkan pengenal `EDMT-<n>` dari SEQ_WORK_POLIS.
func (g *Gudang) IDKasusBerikut(ctx context.Context, tx *db.Tx) (string, error) {
	urut, err := g.db.NomorBerikut(ctx, tx, sequenceKerjaPolis)
	if err != nil {
		return "", err
	}
	return models.RakitIDKasus(urut), nil
}

func sqlSisipKerja(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, LINI, POSITION, STATUS_WORK, FLAG_ONGOING_POLICY, COVER_KEY, CREATE_OP, CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE)
		VALUES (:1, :2, :3, :4, :5, NULL, :6, :7, SYSDATE, SYSDATE)`, t)
}

func sqlSisipGenerasi(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, PRODKE, NOENDORS, OLD_POLIS_ID, EDM_TYPE, TGL_INPUT, USERNAME, POSITION_NOTE)
		VALUES (:1, :2, :3, :4, :5, SYSDATE, :6, :7)`, t)
}

// GenerasiBaru - kunci generasi endorsemen saat lahir (`CreateEDMT` langkah 14-15 / `SetEDMTNoPolis`).
type GenerasiBaru struct {
	// ProdKe - nomor generasi (>= 1).
	ProdKe int
	// EDMNo - `PolicyTreatyIn.EDMNo` -> NOENDORS.
	EDMNo string
	// OldPolisID - ID generasi tepat sebelumnya -> OLD_POLIS_ID (UQ: satu penerus; ID-10, AC 3-4).
	OldPolisID string
	// EDMType - `PolicyTreatyIn.EDMType` (param edmtype `CreateEDMT` langkah 14) -> EDM_TYPE.
	EDMType string
}

// SisipKasus melahirkan kasus di posisi admin - connector `Start1 -> Assignment2`: `.Position = "4";
// .PositionNote = "ReasTreatyInAdmin"; .FlagOnGoingPolicy = "1"` - bersama baris generasinya.
//
// Penerus kedua atas generasi yang sama ditolak BASIS DATA (UQ_GP_TREATY_OLD) -> `ErrGenerasiSudahDiendorse`.
func (g *Gudang) SisipKasus(ctx context.Context, tx *db.Tx, id, pembuat, namaPembuat string, gb GenerasiBaru) error {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return err
	}
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, "menyisipkan kasus", sqlSisipKerja(kerja),
		id, models.LiniKasus, models.PositionAdmin, models.AssignmentAdmin, "1",
		db.KosongJadiNil(pembuat), db.KosongJadiNil(namaPembuat))
	if err != nil {
		return err
	}
	if err := db.PastikanSatuBaris(hasil, "penyisipan kasus"); err != nil {
		return err
	}
	hasil, err = jalankan(ctx, tx, "menyisipkan generasi endorsemen", sqlSisipGenerasi(gen),
		id, gb.ProdKe, db.KosongJadiNil(gb.EDMNo), db.KosongJadiNil(gb.OldPolisID), db.KosongJadiNil(gb.EDMType),
		db.KosongJadiNil(pembuat), models.PosisiAdmin)
	if err != nil {
		if strings.Contains(err.Error(), "ORA-00001") {
			return ErrGenerasiSudahDiendorse
		}
		return err
	}
	return db.PastikanSatuBaris(hasil, "penyisipan generasi endorsemen")
}

func sqlKeadaan(kerja, gen string) string {
	return fmt.Sprintf(`SELECT w.ID, w.POSITION, w.STATUS_WORK, g.POSITION_NOTE, g.NOPOLIS, g.PRODKE, g.OLD_POLIS_ID,
	        CASE WHEN %s THEN 0 ELSE 1 END, w.CREATE_OP,
	        TO_CHAR(w.TGL_CREATE, '%s')
	   FROM %s w JOIN %s g ON g.ID = w.ID
	  WHERE w.ID = :1 AND g.PRODKE >= 1`, syaratTerbuka(gen), fmtTanggal, kerja, gen)
}

// Keadaan membaca keadaan kerja satu kasus endorsemen. Di dalam transaksi bila `tx` terisi.
func (g *Gudang) Keadaan(ctx context.Context, tx *db.Tx, id string) (models.Kasus, error) {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return models.Kasus{}, err
	}
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return models.Kasus{}, err
	}
	q := sqlKeadaan(kerja, gen)
	if err := db.PeriksaSQL(q); err != nil {
		return models.Kasus{}, err
	}
	var k models.Kasus
	var pos, status, catatan, nopol, lama, op, lahir sql.NullString
	var prodke sql.NullInt64
	var tutup int
	err = g.pembaca(tx).QueryRowContext(ctx, q, id).Scan(&k.ID, &pos, &status, &catatan, &nopol, &prodke, &lama, &tutup, &op, &lahir)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Kasus{}, ErrKasusTidakAda
	}
	if err != nil {
		return models.Kasus{}, fmt.Errorf("repository: membaca keadaan kasus: %w", err)
	}
	k.Position, k.StatusWork, k.PositionNote = teks(pos), teks(status), teks(catatan)
	k.NoPolis, k.OldPolisID, k.CreateOp, k.TglCreate = teks(nopol), teks(lama), teks(op), teks(lahir)
	k.ProdKe = int(prodke.Int64)
	k.GenerasiTertutup = tutup == 1
	return k, nil
}

// KunciKasus mengunci baris kasus (`FOR UPDATE`) di transaksi pemanggil dan memastikan tahapnya masih
// `statusHarap` - tahap yang DIBACA sebelum transaksi.
func (g *Gudang) KunciKasus(ctx context.Context, tx *db.Tx, id, statusHarap string) error {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`SELECT STATUS_WORK FROM %s WHERE ID = :1 FOR UPDATE`, kerja)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	var status sql.NullString
	err = tx.QueryRowContext(ctx, q, id).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrKasusTidakAda
	}
	if err != nil {
		return fmt.Errorf("repository: mengunci kasus: %w", err)
	}
	if status.String != statusHarap {
		return ErrTahapBerubah
	}
	return nil
}

func sqlPindahKerja(t string) string {
	return fmt.Sprintf(`UPDATE %s SET POSITION = NVL(:1, POSITION), STATUS_WORK = :2, TGL_UPDATE = SYSDATE
		  WHERE ID = :3 AND STATUS_WORK = :4`, t)
}

func sqlPindahGenerasi(t string) string {
	return fmt.Sprintf(`UPDATE %s g SET POSITION_NOTE = :1 WHERE g.ID = :2 AND %s`, t, syaratTerbuka(t))
}

// PindahPosisi memindahkan kasus ke posisi lain - tugas properti connector (`pyWorkPage.Position` bila ditulis
// connector, `PositionNote`) dan assignment tujuannya. `position` kosong = Position tidak diubah connector.
func (g *Gudang) PindahPosisi(ctx context.Context, tx *db.Tx, id, statusLama, posisiBaru, position string) error {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return err
	}
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, "memindahkan kasus", sqlPindahKerja(kerja),
		db.KosongJadiNil(position), models.AssignmentPosisi(posisiBaru), id, statusLama)
	if err != nil {
		return err
	}
	if n, _ := hasil.RowsAffected(); n != 1 {
		return ErrTahapBerubah
	}
	hasil, err = jalankan(ctx, tx, "memindahkan posisi generasi", sqlPindahGenerasi(gen), posisiBaru, id)
	if err != nil {
		return err
	}
	if n, _ := hasil.RowsAffected(); n != 1 {
		return ErrGenerasiTertutup
	}
	return nil
}

func sqlTutupKerja(t string) string {
	return fmt.Sprintf(`UPDATE %s SET STATUS_WORK = :1, POSITION = NULL, TGL_UPDATE = SYSDATE
		  WHERE ID = :2 AND STATUS_WORK = :3`, t)
}

// TutupKasus menyelesaikan kasus (End3) dengan status akhirnya. POSITION dikosongkan (pola NB / premiumlistlife).
func (g *Gudang) TutupKasus(ctx context.Context, tx *db.Tx, id, statusLama, statusAkhir string) error {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, "menutup kasus", sqlTutupKerja(kerja), statusAkhir, id, statusLama)
	if err != nil {
		return err
	}
	if n, _ := hasil.RowsAffected(); n != 1 {
		return ErrTahapBerubah
	}
	return nil
}

// kolomCariPortal - kolom kotak saring portal (perintah work owner 07-10-2026 "pencarian ... buat bisa mencari nomor
// nb/edm, insured name dll"): nomor kasus EDMT-n, Offer No, nomor polis (generasi dan polis NB lama), EDM No, insured
// (generasi dan quotation), group business, SOB, ceding, marketing, treaty group, class of business, nama pembuat -
// seragam dengan portal NB (sesi NB TREATY 07-10-2026) ditambah dua kolom khas EDM (OLD_POLICY_NO, NOENDORS). XML
// filter B hanya `q.OLD_POLICY_NO`.
var kolomCariPortal = []string{"w.ID", "g.NO_OFFER", "g.NOPOLIS", "q.OLD_POLICY_NO", "g.NOENDORS", "g.INSURED_NAME",
	"q.INSURED_NAME", "q.BUSINESS_NAME", "g.SOB_NAME", "g.CEDING_CO_NAME", "q.MARKETING_NAME", "g.TREATY_GROUP_NAME",
	"g.BIZ_NAME", "w.CREATE_OP_NAME"}

// polaLike - % _ dan garis miring terbalik kotak saring dicari harfiah (`ESCAPE` garis miring terbalik).
var polaLike = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// sqlDaftarKasus merakit daftar portal beserta argumen ikatnya - RD `InboxEDM_RD2` (kelas
// ASM-FW-GISFW-Work-EndorsementTreaty, pyMaxRecords 500, urut `.pxUpdateDateTime DESC`, logika `B AND C AND A
// AND D AND E`):
//
//	A  `.pxCreateOperator = Param.UserNameID`                                  -> s.Pembuat
//	B  `.OfferFacIn.QuotationData.OldPolicyNo Contains Param.FilterTermForEndorsement` -> s.Cari, DIPERLUAS
//	   (WO 07-10-2026): setiap kata `kolomCariPortal`
//	C  `.pyStatusWork != "Resolved-Completed"`                                  -> !s.Selesai (switch Resolved =
//	   kebalikannya, aturan portal NB keputusan WO 07-10-2026)
//	D  `.Quotation.TeamGroup = Param.TeamGroup` - Param dari `InputParam.CARI33`, pengisinya TIDAK ada di korpus
//	   -> tidak dibangun (butir terbuka)
//	E  `.PolicyTreatyIn.PolicyNo = Param.PolicyNo` - tidak dikirim grid portal -> tidak dibangun
//
// ⚠️ Grid portal memanggil `Assign-WorkBasket.InboxEDM_RD2` (tidak ada di korpus); yang dibangun versi kelas
// EndorsementTreaty yang ADA di korpus.
// Kotak masuk Beranda memakai `s.Posisi` / `s.Antrean` / `s.PembuatPosisi` (pola NB).
func sqlDaftarKasus(kerja, gen, quot string, s models.SaringanKasus) (string, []any) {
	var b strings.Builder
	args := []any{models.LiniKasus, models.StatusDitolak, models.StatusSelesai}
	pen := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf(":%d", len(args))
	}
	fmt.Fprintf(&b, `SELECT w.ID, g.NO_OFFER, NVL(g.NOPOLIS, q.OLD_POLICY_NO), g.NOENDORS, g.SOB_NAME, g.CEDING_CO_NAME,
	        g.EDM_TYPE, q.PROPORTIONAL_TYPE, q.MARKETING_NAME, g.NB_STATUS, w.STATUS_WORK, g.POSITION_NOTE,
	        TO_CHAR(w.TGL_CREATE, '%s'), TO_CHAR(g.START_DATE, '%s'), TO_CHAR(w.TGL_UPDATE, '%s'), TO_CHAR(g.TGL_PROD, '%s')
	   FROM %s w
	   JOIN %s g ON g.ID = w.ID
	   LEFT JOIN %s q ON q.POLIS_ID = g.ID
	  WHERE g.PRODKE >= 1 AND w.LINI = :1
	    AND w.ID LIKE '%s%%'`, fmtTanggal, fmtTanggal, fmtTanggal, fmtTanggal, kerja, gen, quot, models.AwalanKasus)
	if !s.Selesai {
		// filter C RD InboxEDM_RD2 (`!= "Resolved-Completed"`). ⚠️ End3 Pega menutup kasus tanpa status tertulis;
		// sistem baru menutup dengan dua status (ketetapan NB P24) - keduanya "sudah selesai" bagi filter C.
		b.WriteString(`
	    AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK NOT IN (:2, :3))`)
	} else {
		// switch Resolved (aturan portal NB, keputusan WO 07-10-2026)
		b.WriteString(`
	    AND w.STATUS_WORK IN (:2, :3)`)
	}
	// kotak saring (perintah WO 07-10-2026): setiap kata cocok dengan salah satu kolom portal, tanpa beda huruf
	for _, k := range models.KataCari(s.Cari) {
		pola := "%" + polaLike.Replace(k) + "%"
		atau := make([]string, len(kolomCariPortal))
		for i, c := range kolomCariPortal {
			atau[i] = "UPPER(" + c + ") LIKE " + pen(pola) + ` ESCAPE '\'`
		}
		fmt.Fprintf(&b, `
	    AND (%s)`, strings.Join(atau, " OR "))
	}
	if s.Posisi != "" {
		fmt.Fprintf(&b, `
	    AND g.POSITION_NOTE = %s`, pen(s.Posisi))
	}
	pembuat := ""
	if s.Pembuat != "" {
		pembuat = "w.CREATE_OP = " + pen(s.Pembuat)
		if s.PembuatPosisi != "" {
			pembuat = "(" + pembuat + " AND g.POSITION_NOTE = " + pen(s.PembuatPosisi) + ")"
		}
	}
	antrean := ""
	if len(s.Antrean) > 0 {
		var ps []string
		for _, a := range s.Antrean {
			ps = append(ps, pen(a))
		}
		antrean = "g.POSITION_NOTE IN (" + strings.Join(ps, ", ") + ")"
	}
	switch {
	case pembuat != "" && antrean != "":
		fmt.Fprintf(&b, `
	    AND (%s OR %s)`, pembuat, antrean)
	case pembuat != "":
		fmt.Fprintf(&b, `
	    AND %s`, pembuat)
	case antrean != "":
		fmt.Fprintf(&b, `
	    AND %s`, antrean)
	}
	fmt.Fprintf(&b, `
	  ORDER BY w.TGL_UPDATE DESC
	  FETCH FIRST %d ROWS ONLY`, models.BatasDaftarPortal)
	return b.String(), args
}

// DaftarKasus membaca kasus endorsemen untuk portal / kotak masuk (`sqlDaftarKasus`).
func (g *Gudang) DaftarKasus(ctx context.Context, s models.SaringanKasus) ([]models.RingkasanKasus, error) {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return nil, err
	}
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return nil, err
	}
	quot, err := g.nama(models.TabelQuotation.Nama)
	if err != nil {
		return nil, err
	}
	q, args := sqlDaftarKasus(kerja, gen, quot, s)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca daftar kasus: %w", err)
	}
	defer rows.Close()
	out := []models.RingkasanKasus{}
	for rows.Next() {
		var r models.RingkasanKasus
		var off, np, edm, sob, ced, jenis, prop, mkt, nb, st, pn, tg, mulai, ubah, prod sql.NullString
		if err := rows.Scan(&r.ID, &off, &np, &edm, &sob, &ced, &jenis, &prop, &mkt, &nb, &st, &pn, &tg, &mulai, &ubah, &prod); err != nil {
			return nil, fmt.Errorf("repository: membaca baris daftar kasus: %w", err)
		}
		r.NoOffer, r.NoPolis, r.EDMNo, r.SOBName, r.CedingCoName = teks(off), teks(np), teks(edm), teks(sob), teks(ced)
		r.EDMType, r.ProportionalType, r.MarketingName, r.NBStatus = teks(jenis), teks(prop), teks(mkt), teks(nb)
		r.StatusWork, r.PositionNote, r.TglCreate, r.StartDate, r.TglUpdate = teks(st), teks(pn), teks(tg), teks(mulai), teks(ubah)
		r.TglProd = teks(prod)
		out = append(out, r)
	}
	return out, rows.Err()
}

// sqlHitungKotakMasuk - cacah berkas endorsemen terbuka per posisi yang MENUNGGU akun (kotak masuk Beranda, pola
// NB): Admin = buatan akun yang masih di Admin, atasan = seluruh antrean workbasket yang ia pegang.
func sqlHitungKotakMasuk(kerja, gen, akun string, admin bool, atasan []string) (string, []any) {
	args := []any{models.LiniKasus, models.StatusDitolak, models.StatusSelesai}
	pen := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf(":%d", len(args))
	}
	var syarat []string
	if admin {
		syarat = append(syarat, fmt.Sprintf("(g.POSITION_NOTE = %s AND w.CREATE_OP = %s)", pen(models.PosisiAdmin), pen(akun)))
	}
	if len(atasan) > 0 {
		var ps []string
		for _, a := range atasan {
			ps = append(ps, pen(a))
		}
		syarat = append(syarat, "g.POSITION_NOTE IN ("+strings.Join(ps, ", ")+")")
	}
	q := fmt.Sprintf(`SELECT g.POSITION_NOTE, COUNT(*)
	   FROM %s w
	   JOIN %s g ON g.ID = w.ID
	  WHERE g.PRODKE >= 1 AND w.LINI = :1 AND w.ID LIKE '%s%%'
	    AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK NOT IN (:2, :3))
	    AND (%s)
	  GROUP BY g.POSITION_NOTE`, kerja, gen, models.AwalanKasus, strings.Join(syarat, " OR "))
	return q, args
}

// HitungKotakMasuk - cacah kotak masuk Beranda per posisi; tanpa syarat = peta kosong tanpa query.
func (g *Gudang) HitungKotakMasuk(ctx context.Context, akun string, admin bool, atasan []string) (map[string]int, error) {
	out := map[string]int{}
	if !admin && len(atasan) == 0 {
		return out, nil
	}
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return nil, err
	}
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return nil, err
	}
	q, args := sqlHitungKotakMasuk(kerja, gen, akun, admin, atasan)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: menghitung kotak masuk: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var pos sql.NullString
		var n int
		if err := rows.Scan(&pos, &n); err != nil {
			return nil, fmt.Errorf("repository: membaca cacah kotak masuk: %w", err)
		}
		out[pos.String] = n
	}
	return out, rows.Err()
}
