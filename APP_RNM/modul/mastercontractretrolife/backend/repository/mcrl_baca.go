package repository

// Pembaca kelima tabel warisan - padanan grid Pega (paket 1).
//
//	tahun     `BrowseTreatyYear_Life_RD`         urut `.ID ASC` b765
//	kontrak   `BrowseTreatyContract_Life_RD`     `.IDTREATYYEAR =` b628, urut `.ID ASC` b940
//	reinsurer `BrowseDetailTreatyReisurerLife_RD` `.TREATYYEARID =` b613, `.TREATYCONTRACTID =` b631, urut `.ID DESC` b896
//	security  `BrowseSecurityReinsurer_Life_RD`  tiga saringan b579/b592/b609; RD tanpa urut -> `ID ASC` (stabil)
//	business  `BrowseTreatyBusiness_Life_RD`     b615/b627, urut `.TGLUPDATE ASC` b939 (+ `ID` pemecah seri)
//
// ⛔ Saringan `REINSTYPEID` RD kontrak (b642) TIDAK dipasang: grid tidak
// mengirim parameternya (`InputRetroLimitReinsurers.xml` b9252 hanya
// `IDTREATYYEAR`).
//
// ⛔ Anak dibaca dengan KEDUA kunci induk (tahun + kontrak), persis RD: baris
// yang salinan tahunnya menyimpang dari kontraknya tidak ikut tampil - sama
// seperti di Pega.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

// ErrTidakAda - baris dengan ID itu tidak ada di tabelnya.
var ErrTidakAda = errors.New("repository: row not found")

// barisTeks adalah satu baris hasil, dikunci nama kolom.
type barisTeks map[string]sql.NullString

func pindai(baca interface{ Scan(...any) error }, kolom []string) (barisTeks, error) {
	nilai := make([]sql.NullString, len(kolom))
	arah := make([]any, len(kolom))
	for i := range nilai {
		arah[i] = &nilai[i]
	}
	if err := baca.Scan(arah...); err != nil {
		return nil, err
	}
	b := barisTeks{}
	for i, k := range kolom {
		b[k] = nilai[i]
	}
	return b, nil
}

func (b barisTeks) s(k string) string { return b[k].String }

func (b barisTeks) tanggal(k string) (time.Time, error) {
	v := b[k]
	if !v.Valid || v.String == "" {
		return time.Time{}, nil
	}
	t, err := utils.ParseTanggal(v.String)
	if err != nil {
		return time.Time{}, fmt.Errorf("repository: row %s column %s has value %q: %w", b.s("ID"), k, v.String, err)
	}
	return t, nil
}

func (b barisTeks) desimal(k string) (*apd.Decimal, error) { return db.UraiDesimal(b.s("ID"), k, b[k]) }

// urai mengumpulkan galat urai pertama - supaya pembangun model tetap ringkas.
type urai struct{ err error }

func (u *urai) t(b barisTeks, k string) time.Time {
	v, err := b.tanggal(k)
	if u.err == nil {
		u.err = err
	}
	return v
}

func (u *urai) d(b barisTeks, k string) *apd.Decimal {
	v, err := b.desimal(k)
	if u.err == nil {
		u.err = err
	}
	return v
}

func keTahun(b barisTeks) (models.TahunTreaty, error) {
	var u urai
	t := models.TahunTreaty{ID: b.s("ID"), TreatyYear: b.s("TREATYYEAR"), UnderwritingYear: b.s("UNDERWRITINGYEAR"),
		UserID: b.s("USERID"), StartDate: u.t(b, "STARTDATE"), EndDate: u.t(b, "ENDDATE"), TglUpdate: u.t(b, "TGLUPDATE")}
	return t, u.err
}

func keKontrak(b barisTeks) (models.Kontrak, error) {
	var u urai
	k := models.Kontrak{ID: b.s("ID"), IDTreatyYear: b.s("IDTREATYYEAR"), ReinsTypeID: b.s("REINSTYPEID"),
		ReinsTypeName: b.s("REINSTYPENAME"), UserID: b.s("USERID"), TglUpdate: u.t(b, "TGLUPDATE"),
		TreatyStartDate: u.t(b, "TREATYSTARTDATE"), TreatyEndDate: u.t(b, "TREATYENDDATE"),
		IDR: u.d(b, "IDR"), USD: u.d(b, "USD"), BIDR: u.d(b, "B_IDR"), BUSD: u.d(b, "B_USD"),
		IDRSelisih: u.d(b, "IDR_SELISIH"), USDSelisih: u.d(b, "USD_SELISIH")}
	return k, u.err
}

func keReinsurer(b barisTeks) (models.Reinsurer, error) {
	var u urai
	r := models.Reinsurer{ID: b.s("ID"), TreatyYearID: b.s("TREATYYEARID"), TreatyContractID: b.s("TREATYCONTRACTID"),
		ReinsTypeID: b.s("REINSTYPEID"), ReinsTypeName: b.s("REINSTYPENAME"), ReinsurerID: b.s("REINSURERID"),
		ReinsurerName: b.s("REINSURERNAME"), PctShare: u.d(b, "PCTSHARE"), Komisi: u.d(b, "COMMISION"),
		OvrComm: u.d(b, "OVR_COMM"), UserID: b.s("USERID"), TglUpdate: u.t(b, "TGLUPDATE")}
	return r, u.err
}

func keSecurity(b barisTeks) (models.SecurityReinsurer, error) {
	var u urai
	s := models.SecurityReinsurer{ID: b.s("ID"), TreatyYearID: b.s("TREATYYEARID"),
		TreatyContractID: b.s("TREATYCONTRACTID"), TreatyReinsurerID: b.s("TREATYREINSURERID"),
		ReinsurerID: b.s("REINSURERID"), ReinsurerName: b.s("REINSURERNAME"), PctShare: u.d(b, "PCTSHARE"),
		UserID: b.s("USERID"), TglUpdate: u.t(b, "TGLUPDATE")}
	return s, u.err
}

func keBusiness(b barisTeks) (models.Business, error) {
	var u urai
	x := models.Business{ID: b.s("ID"), TreatyYearID: b.s("TREATYYEARID"), TreatyYear: b.s("TREATYYEAR"),
		TreatyContractID: b.s("TREATYCONTRACTID"), ReinsTypeID: b.s("REINSTYPEID"), ReinsTypeName: b.s("REINSTYPENAME"),
		BizCode: b.s("BIZCODE"), BizName: b.s("BIZNAME"), RIRateID: b.s("RIRATEID"), RIRate: b.s("RIRATE"),
		UserID: b.s("USERID"), TglUpdate: u.t(b, "TGLUPDATE")}
	return x, u.err
}

// daftar menjalankan kueri banyak baris dan membangun modelnya.
func daftar[T any](ctx context.Context, q kuerier, teks string, kolom []string, bangun func(barisTeks) (T, error),
	nama string, args ...any) ([]T, error) {

	rows, err := q.QueryContext(ctx, teks, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: reading %s list: %w", nama, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []T
	for rows.Next() {
		b, err := pindai(rows, kolom)
		if err != nil {
			return nil, fmt.Errorf("repository: scanning %s row: %w", nama, err)
		}
		m, err := bangun(b)
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: reading %s list: %w", nama, err)
	}
	return hasil, nil
}

// satu menjalankan kueri satu baris; nol baris = ErrTidakAda.
func satu[T any](ctx context.Context, q kuerier, teks string, kolom []string, bangun func(barisTeks) (T, error),
	nama, id string) (T, error) {

	var nol T
	b, err := pindai(q.QueryRowContext(ctx, teks, id), kolom)
	if errors.Is(err, sql.ErrNoRows) {
		return nol, fmt.Errorf("%w: %s %s", ErrTidakAda, nama, id)
	}
	if err != nil {
		return nol, fmt.Errorf("repository: reading %s %s: %w", nama, id, err)
	}
	return bangun(b)
}

func sqlPilih(kolom []string, saring, urut string) func(string) string {
	return func(tabel string) string {
		q := fmt.Sprintf(`SELECT %s FROM %s`, daftarPilih(kolom), tabel)
		if saring != "" {
			q += " WHERE " + saring
		}
		if urut != "" {
			q += " ORDER BY " + urut
		}
		return q
	}
}

// DaftarTahun - seluruh tahun treaty, `ID ASC`.
func (g *Gudang) DaftarTahun(ctx context.Context) ([]models.TahunTreaty, error) {
	q, err := g.siapkan(TabelTahun, sqlPilih(KolomTahun, "", "ID ASC"))
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, q, KolomTahun, keTahun, "treaty year")
}

// AmbilTahun membaca satu tahun treaty (di dalam tx bila ada).
func (g *Gudang) AmbilTahun(ctx context.Context, tx *db.Tx, id string) (models.TahunTreaty, error) {
	q, err := g.siapkan(TabelTahun, sqlPilih(KolomTahun, "ID = :1", ""))
	if err != nil {
		return models.TahunTreaty{}, err
	}
	return satu(ctx, g.kueri(tx), q, KolomTahun, keTahun, "treaty year", id)
}

// DaftarKontrak - kontrak satu tahun treaty, `ID ASC`.
func (g *Gudang) DaftarKontrak(ctx context.Context, tx *db.Tx, tahunID string) ([]models.Kontrak, error) {
	q, err := g.siapkan(TabelKontrak, sqlPilih(KolomKontrak, "IDTREATYYEAR = :1", "ID ASC"))
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.kueri(tx), q, KolomKontrak, keKontrak, "treaty contract", tahunID)
}

// AmbilKontrak membaca satu kontrak.
func (g *Gudang) AmbilKontrak(ctx context.Context, tx *db.Tx, id string) (models.Kontrak, error) {
	q, err := g.siapkan(TabelKontrak, sqlPilih(KolomKontrak, "ID = :1", ""))
	if err != nil {
		return models.Kontrak{}, err
	}
	return satu(ctx, g.kueri(tx), q, KolomKontrak, keKontrak, "treaty contract", id)
}

// DaftarReinsurer - reinsurer satu kontrak, `ID DESC` (RD b896).
func (g *Gudang) DaftarReinsurer(ctx context.Context, tx *db.Tx, tahunID, kontrakID string) ([]models.Reinsurer, error) {
	q, err := g.siapkan(TabelReinsurer, sqlPilih(KolomReinsurer, "TREATYYEARID = :1 AND TREATYCONTRACTID = :2", "ID DESC"))
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.kueri(tx), q, KolomReinsurer, keReinsurer, "reinsurer", tahunID, kontrakID)
}

// AmbilReinsurer membaca satu reinsurer.
func (g *Gudang) AmbilReinsurer(ctx context.Context, tx *db.Tx, id string) (models.Reinsurer, error) {
	q, err := g.siapkan(TabelReinsurer, sqlPilih(KolomReinsurer, "ID = :1", ""))
	if err != nil {
		return models.Reinsurer{}, err
	}
	return satu(ctx, g.kueri(tx), q, KolomReinsurer, keReinsurer, "reinsurer", id)
}

// DaftarSecurity - security di bawah satu reinsurer.
func (g *Gudang) DaftarSecurity(ctx context.Context, tx *db.Tx, tahunID, kontrakID, reinsurerID string) (
	[]models.SecurityReinsurer, error) {

	q, err := g.siapkan(TabelSecurity, sqlPilih(KolomSecurity,
		"TREATYYEARID = :1 AND TREATYCONTRACTID = :2 AND TREATYREINSURERID = :3", "ID ASC"))
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.kueri(tx), q, KolomSecurity, keSecurity, "security reinsurer", tahunID, kontrakID, reinsurerID)
}

// AmbilSecurity membaca satu security reinsurer.
func (g *Gudang) AmbilSecurity(ctx context.Context, tx *db.Tx, id string) (models.SecurityReinsurer, error) {
	q, err := g.siapkan(TabelSecurity, sqlPilih(KolomSecurity, "ID = :1", ""))
	if err != nil {
		return models.SecurityReinsurer{}, err
	}
	return satu(ctx, g.kueri(tx), q, KolomSecurity, keSecurity, "security reinsurer", id)
}

// DaftarBusiness - business satu kontrak, `TGLUPDATE ASC` (RD b939).
func (g *Gudang) DaftarBusiness(ctx context.Context, tx *db.Tx, tahunID, kontrakID string) ([]models.Business, error) {
	q, err := g.siapkan(TabelBusiness, sqlPilih(KolomBusiness, "TREATYYEARID = :1 AND TREATYCONTRACTID = :2",
		"TGLUPDATE ASC, ID ASC"))
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.kueri(tx), q, KolomBusiness, keBusiness, "business", tahunID, kontrakID)
}

// AmbilBusiness membaca satu business.
func (g *Gudang) AmbilBusiness(ctx context.Context, tx *db.Tx, id string) (models.Business, error) {
	q, err := g.siapkan(TabelBusiness, sqlPilih(KolomBusiness, "ID = :1", ""))
	if err != nil {
		return models.Business{}, err
	}
	return satu(ctx, g.kueri(tx), q, KolomBusiness, keBusiness, "business", id)
}
