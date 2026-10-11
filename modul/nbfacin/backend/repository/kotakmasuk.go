package repository

// Kotak masuk Beranda NB Fac In (permintaan work owner 07-10-2026 "ini kenapa tidak ada ya, padahal ini login teamleader";
// pola nbtreatyin keputusan 06-10-2026). Case yang MENUNGGU akun: T_WORK_POLIS ber-LINI FAC yang POSITION-nya salah satu
// workbasket yang ia pegang; STATUS_WORK "Resolved-…" tidak dihitung. Nama workbasket dari M_WORKBASKET (master login inti,
// baca saja). Parameter terikat, ditulis menurut urutan kemunculan.
//
// Saringan team group (keputusan work owner 08-10-2026 "yg muncul sesuai workbasketnya dan team groupnya"; PENYIMPANGAN
// SADAR - di Pega kotak masuk hanya per workbasket, `Flow\InputInwardFacultativeOffer.xml` 21 assignment ToWorkBasket):
// tahap marketing (POSITION "Offer" / "ReasFacInMarketing", menggantikan aturan CREATE_OP) dan Team Leader hanya case yang
// team-nya (T_QUOTATIONDATA.MOID -> MARKETINGOFFICER.TEAMGROUP, dibaca saat itu; Pega menyimpan QuotationData.TeamGroup) salah
// satu team akun. Workbasket UW / lainnya tanpa saringan team. Portal opportunity TETAP CREATE_OP (portal.go).

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelWorkbasket - master workbasket login inti (baca saja).
	TabelWorkbasket = "M_WORKBASKET"
	// polaSelesai - STATUS_WORK case selesai (mis. "Resolved-Decline", tiket 57) - tidak menunggu siapa pun.
	polaSelesai = "Resolved-%"
	// batasKotakMasuk - baris terbanyak daftar kotak masuk (pola BatasDaftarPortal nbtreatyin, pyMaxRecords 500).
	batasKotakMasuk = 500
)

// AntreanTim - POSITION yang dihitung hanya bila team case salah satu Team (Team kosong = antrean kosong).
type AntreanTim struct {
	Posisi, Team []string
}

// SaringanKotakMasuk - siapa yang ditunggu: antrean bersaring team (marketing, Team Leader) dan/atau posisi lain.
type SaringanKotakMasuk struct {
	Tim []AntreanTim
	// Lain - workbasket lain yang dipegang (POSITION = workbasket), tanpa saringan team.
	Lain []string
}

func (s SaringanKotakMasuk) kosong() bool {
	for _, t := range s.Tim {
		if len(t.Posisi) > 0 && len(t.Team) > 0 {
			return false
		}
	}
	return len(s.Lain) == 0
}

// PembacaKotakMasuk - cacah, daftar, dan nama workbasket kotak masuk Beranda.
type PembacaKotakMasuk interface {
	// HitungKotakMasuk - cacah per POSITION (posisi tanpa case tidak ada di peta).
	HitungKotakMasuk(ctx context.Context, s SaringanKotakMasuk) (map[string]int, error)
	// DaftarMenunggu - baris inbox urut TGL_UPDATE terbaru dulu, paling banyak 500.
	DaftarMenunggu(ctx context.Context, s SaringanKotakMasuk) ([]models.BarisKotakMasuk, error)
	// NamaWorkbasket - M_WORKBASKET.NAME per WORKBASKET_ID; tak terdaftar = tidak di peta.
	NamaWorkbasket(ctx context.Context, id []string) (map[string]string, error)
	// TeamKotakMasuk - team group akun: marketing = POOLDATA.MARKETINGOFFICER.TEAMGROUP ber-AKSES_LOGIN = login (kandidat A),
	// Team Leader = POOLDATA.M_LIMIT_PROPERTYY.TEAM_GROUP ber-LOGIN = login (kandidat B, bisa lebih dari satu). Login
	// dibandingkan UPPER(TRIM); bahwa kolom itu berisi nilai yang sama dengan M_LOGIN_GO.LOGIN_ID `belum terverifikasi`
	// (pemeriksaan DBA). Nilai di-TRIM, unik.
	TeamKotakMasuk(ctx context.Context, login string) (marketing, teamLeader []string, err error)
}

// teamCase - team group case: T_QUOTATIONDATA.MOID -> MARKETINGOFFICER.TEAMGROUP (tiket 57 / 61), di-TRIM.
func teamCase(quo, mo string) string {
	return "TRIM((SELECT MAX(m2.TEAMGROUP) FROM " + quo + " q2 JOIN " + mo + " m2 ON m2.ID = q2.MOID WHERE q2.PARENT_ID = w.ID))"
}

// saringKotakMasuk - WHERE kotak masuk; argumen menurut urutan kemunculan sesudah `awal` argumen sebelumnya.
func saringKotakMasuk(s SaringanKotakMasuk, awal int, quo, mo string) (string, []any) {
	var args []any
	pen := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf(":%d", awal+len(args))
	}
	daftar := func(xs []string) string {
		ps := make([]string, len(xs))
		for i, x := range xs {
			ps[i] = pen(x)
		}
		return strings.Join(ps, ", ")
	}
	q := " WHERE w.LINI = " + pen(LiniFacIn) + " AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK NOT LIKE " + pen(polaSelesai) + ")"
	var syarat []string
	for _, t := range s.Tim {
		if len(t.Posisi) > 0 && len(t.Team) > 0 {
			syarat = append(syarat, "(w.POSITION IN ("+daftar(t.Posisi)+") AND "+teamCase(quo, mo)+" IN ("+daftar(t.Team)+"))")
		}
	}
	if len(s.Lain) > 0 {
		syarat = append(syarat, "w.POSITION IN ("+daftar(s.Lain)+")")
	}
	return q + " AND (" + strings.Join(syarat, " OR ") + ")", args
}

func sqlHitungKotakMasuk(t tabelPortal, s SaringanKotakMasuk) (string, []any) {
	where, args := saringKotakMasuk(s, 0, t.quo, t.mo)
	return "SELECT w.POSITION, COUNT(*) FROM " + t.work + " w" + where + " GROUP BY w.POSITION", args
}

// sqlTeamMarketing / sqlTeamLeader - TeamKotakMasuk (bind :1 = login).
func sqlTeamMarketing(mo string) string {
	return "SELECT DISTINCT TRIM(TEAMGROUP) FROM " + mo + " WHERE UPPER(TRIM(AKSES_LOGIN)) = UPPER(TRIM(:1)) AND TRIM(TEAMGROUP) IS NOT NULL"
}

func sqlTeamLeader(limit string) string {
	return "SELECT DISTINCT TRIM(TO_CHAR(TEAM_GROUP)) FROM " + limit + " WHERE UPPER(TRIM(LOGIN)) = UPPER(TRIM(:1)) AND TEAM_GROUP IS NOT NULL"
}

// sqlDaftarMenunggu - kolom inbox "New Business"; aksi terakhir = APPROVAL / PIC baris HISTORYAKSEPTASIPRODUCTION TERAKHIR
// (KEEP DENSE_RANK LAST atas TGL_INP, NOURUT) ber-IDPEGA = :1 / :2 || ID (awalan pzInsKey). Bind :1, :2 lalu saringan.
func sqlDaftarMenunggu(t tabelPortal, produksi string, s SaringanKotakMasuk) (string, []any) {
	where, args := saringKotakMasuk(s, 2, t.quo, t.mo)
	q := fmt.Sprintf(`SELECT w.ID,
  (SELECT MAX(a.INSUREDNAME) FROM %s a WHERE a.ID = o.ACCOUNT_ID),
  o.CLASS_OF_BUSINESS, q.SOB_NAME, q.CEDING_CO_NAME, g.START_DATE_TIME, %s,
  (SELECT MAX(m.CLIENTNAME) FROM %s m WHERE m.ID = q.MOID),
  w.FLAG_ONGOING_POLICY, TO_CHAR(w.TGL_CREATE, 'DD-MM-YYYY'), g.OFFERING_DATE, g.NB_STATUS,
  (SELECT MAX(h.APPROVAL) KEEP (DENSE_RANK LAST ORDER BY h.TGL_INP, h.NOURUT) FROM %s h WHERE h.IDPEGA IN (w.ID, :1 || w.ID)),
  (SELECT MAX(h.PIC) KEEP (DENSE_RANK LAST ORDER BY h.TGL_INP, h.NOURUT) FROM %s h WHERE h.IDPEGA IN (w.ID, :2 || w.ID)),
  TO_CHAR(w.TGL_UPDATE, 'YYYY-MM-DD HH24:MI:SS'), g.SUBMIT_TO_UW
FROM %s w
LEFT JOIN %s o ON o.ID = w.ID
LEFT JOIN %s q ON q.PARENT_ID = w.ID
LEFT JOIN %s g ON g.ID = w.ID%s
ORDER BY w.TGL_UPDATE DESC NULLS LAST, w.ID DESC FETCH FIRST %d ROWS ONLY`, t.akun, angkaKeluar("g.PERCENT_SHARE"), t.mo, produksi, produksi,
		t.work, t.opp, t.quo, t.general, where, batasKotakMasuk)
	return q, append([]any{AwalanIDPega, AwalanIDPega}, args...)
}

func sqlNamaWorkbasket(wb string, n int) string {
	b := make([]string, n)
	for i := range b {
		b[i] = fmt.Sprintf(":%d", i+1)
	}
	return "SELECT WORKBASKET_ID, NAME FROM " + wb + " WHERE WORKBASKET_ID IN (" + strings.Join(b, ", ") + ")"
}

// HitungKotakMasuk - lihat PembacaKotakMasuk. Saringan kosong = peta kosong tanpa kueri.
func (r *PortalOracle) HitungKotakMasuk(ctx context.Context, s SaringanKotakMasuk) (map[string]int, error) {
	out := map[string]int{}
	if s.kosong() {
		return out, nil
	}
	t, err := tabelBarisCase(r.db)
	if err != nil {
		return nil, err
	}
	q, args := sqlHitungKotakMasuk(t, s)
	baris, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: cacah kotak masuk: %w", err)
	}
	defer baris.Close()
	for baris.Next() {
		var pos sql.NullString
		var n int
		if err := baris.Scan(&pos, &n); err != nil {
			return nil, fmt.Errorf("repository: cacah kotak masuk: %w", err)
		}
		out[pos.String] += n
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: cacah kotak masuk: %w", err)
	}
	return out, nil
}

// DaftarMenunggu - lihat PembacaKotakMasuk. Saringan kosong = kosong tanpa kueri.
func (r *PortalOracle) DaftarMenunggu(ctx context.Context, s SaringanKotakMasuk) ([]models.BarisKotakMasuk, error) {
	hasil := []models.BarisKotakMasuk{}
	if s.kosong() {
		return hasil, nil
	}
	t, err := tabelBarisCase(r.db)
	if err != nil {
		return nil, err
	}
	produksi, err := r.db.Qualify(TabelRiwayatProduksi)
	if err != nil {
		return nil, err
	}
	q, args := sqlDaftarMenunggu(t, produksi, s)
	baris, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca kotak masuk: %w", err)
	}
	defer baris.Close()
	for baris.Next() {
		var v [16]sql.NullString
		tujuan := make([]any, len(v))
		for i := range v {
			tujuan[i] = &v[i]
		}
		if err := baris.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: kotak masuk: %w", err)
		}
		hasil = append(hasil, models.BarisKotakMasuk{CaseID: v[0].String, InsuredName: v[1].String, BusinessName: v[2].String,
			SourceOfBusiness: v[3].String, CedingCompany: v[4].String, StartDateTime: v[5].String, ShareRNM: v[6].String,
			Marketing: v[7].String, FlagOnGoingPolicy: v[8].String, CreatedDate: v[9].String, OfferingDate: v[10].String, NBStatus: v[11].String,
			ApprovalTerakhir: v[12].String, PICAkhir: v[13].String, TglUpdate: v[14].String, SubmitToUW: v[15].String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: kotak masuk: %w", err)
	}
	return hasil, nil
}

// NamaWorkbasket - lihat PembacaKotakMasuk.
func (r *PortalOracle) NamaWorkbasket(ctx context.Context, id []string) (map[string]string, error) {
	out := map[string]string{}
	if len(id) == 0 {
		return out, nil
	}
	wb, err := r.db.Qualify(TabelWorkbasket)
	if err != nil {
		return nil, err
	}
	args := make([]any, len(id))
	for i, x := range id {
		args[i] = x
	}
	baris, err := r.db.QueryContext(ctx, sqlNamaWorkbasket(wb, len(id)), args...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelWorkbasket, err)
	}
	defer baris.Close()
	for baris.Next() {
		var k, n sql.NullString
		if err := baris.Scan(&k, &n); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelWorkbasket, err)
		}
		out[k.String] = n.String
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelWorkbasket, err)
	}
	return out, nil
}

// TeamKotakMasuk - lihat PembacaKotakMasuk.
func (r *PortalOracle) TeamKotakMasuk(ctx context.Context, login string) ([]string, []string, error) {
	mo, err := r.db.Qualify(TabelMarketingOfficer)
	if err != nil {
		return nil, nil, err
	}
	limit, err := r.db.Qualify(TabelLimitProperty)
	if err != nil {
		return nil, nil, err
	}
	marketing, err := r.daftarTeks(ctx, sqlTeamMarketing(mo), login)
	if err != nil {
		return nil, nil, fmt.Errorf("repository: team marketing akun: %w", err)
	}
	tl, err := r.daftarTeks(ctx, sqlTeamLeader(limit), login)
	if err != nil {
		return nil, nil, fmt.Errorf("repository: team Team Leader akun: %w", err)
	}
	return marketing, tl, nil
}

// daftarTeks - satu kolom teks per baris (NULL / kosong dilewati).
func (r *PortalOracle) daftarTeks(ctx context.Context, q string, args ...any) ([]string, error) {
	baris, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer baris.Close()
	var out []string
	for baris.Next() {
		var v sql.NullString
		if err := baris.Scan(&v); err != nil {
			return nil, err
		}
		if s := strings.TrimSpace(v.String); s != "" {
			out = append(out, s)
		}
	}
	return out, baris.Err()
}
