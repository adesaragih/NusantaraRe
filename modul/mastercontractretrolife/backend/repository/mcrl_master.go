package repository

// Master rujukan layar - DIBACA SAJA (paket 1).
//
//	REINS TYPE      `BrowseReinsuranceTypeLimit_RD` b578 `.Flag = 1`, urut `.ID DESC` b765; dropdown nilai `.ID`, tampil `.Note`
//	REINSURER NAME  `BrowseCedingCoLife_RD` b556 `B AND A AND C`: `.ID Contains "L0"` b565,
//	                `.ClientName Contains Param.CedingCoLeader` b584 (kata yang diketik), `.StatusActive = 1` b601;
//	                urut `.ClientName ASC` b745; isi `REINSURERID ← .ID`
//	BUSINESS NAME   `BrowseBusinessLife_RD` b651 `.OLDID StartsWith "L"`, urut `.ID ASC` b1057; tampil `.Note`,
//	                isi `BIZCODE ← .ID`, kolom tampil `.OLDID`
//	R/I RATE        `BrowseRateLifeSummary` (tabel `M_RATE_LIFE_SUMMARY`) `A AND B`: `.ID = param.id` b794,
//	                `.USEDBY Contains param.idusedby` b809 - autocomplete b4534/b4542 mengirim keduanya
//	                KOSONG, jadi kedua saringan gugur; kata yang diketik dicari di `.USEDBY` (medan cari
//	                b4494, pola autocomplete modul ini); urut `.ID ASC` b692; isi `RIRATEID ← .ID` b4527
//	Rate List       `BrowseRateLife_RD` (`M_RATE_LIFE` (dulu view `RATE_LIFE`)) `.IDUSEDBY = Param.idusedby` b860/b868
//	                (`ViewRate.xml` b1055 `ParamID.RIRATEID`); urut `.ID DESC` b747, `.RATE ASC` b784;
//	                `pyMaxRecords` 500 b729
//
// ⛔ K1 (keputusan work owner 01-10-2026, OQ-MCRL-13 + OQ-MCRL-05): kedua objek rate dibaca SAJA, kolom
// yang dibaca RD XML saja - autocomplete `ID`, `USEDBY`; Rate List enam kolom yang ditampilkan grid
// `ViewRate`. Nol `SELECT *`, nol `JSONDATA`. Kolom `M_RATE_LIFE` teks (migrasi inti 929/930, RALAT R7
// riratelife; dulu view VARCHAR2(4000)): dibaca teks apa adanya. Rate List berkunci `IDUSEDBY` (indeks
// `IX_M_RATE_LIFE_IDUSEDBY`) - `IDUSEDBY` kosong ditolak layanan, nol pembacaan tanpa kunci.
//
// ⛔ Kolom fisik AGENT (`ID`, `CLIENTNAME`, `STATUSACTIVE`) dan BUSINESS
// (`ID`, `NOTE`) sama dengan yang dibaca Treaty Contract Out; `OLDID` dari
// nama properti RD. Pencarian autocomplete tidak peka huruf besar-kecil.
// ⛔ Nol INSERT/UPDATE/DELETE atas master (penjaga modul).

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/modul/mastercontractretrolife/backend/models"
)

// BatasPilihan membatasi jumlah pilihan autocomplete per ketikan.
const BatasPilihan = 100

// ErrMasterTidakTerbaca - objek master tidak dapat dibaca (mis. tidak ada di
// skema); pesannya menyebut objeknya, tidak diam.
var ErrMasterTidakTerbaca = errors.New("repository: reference master cannot be read")

// GalatMaster membungkus kegagalan membaca satu master: `Error()` membawa
// sebab aslinya (untuk log server), `PesanLayar()` hanya objeknya - teks mentah
// Oracle tidak pernah sampai ke layar (K8, tiket 10).
type GalatMaster struct {
	Objek string
	Sebab error
}

// MasterTidakTerbaca menyusun GalatMaster - dipakai repository dan tiruan uji.
func MasterTidakTerbaca(objek string, sebab error) error {
	return GalatMaster{Objek: objek, Sebab: sebab}
}

func (g GalatMaster) Error() string {
	if g.Sebab == nil {
		return fmt.Sprintf("%v: %s", ErrMasterTidakTerbaca, g.Objek)
	}
	return fmt.Sprintf("%v: %s: %v", ErrMasterTidakTerbaca, g.Objek, g.Sebab)
}

// Is membuat errors.Is(err, ErrMasterTidakTerbaca) benar.
func (GalatMaster) Is(target error) bool { return target == ErrMasterTidakTerbaca }

// Unwrap membuka sebab aslinya.
func (g GalatMaster) Unwrap() error { return g.Sebab }

// PesanLayar - kalimat untuk layar: objeknya saja.
func (g GalatMaster) PesanLayar() string {
	return fmt.Sprintf("reference master cannot be read: %s", g.Objek)
}

// Teks SQL master - fungsi murni atas nama tabel berkualifikasi (diuji tanpa Oracle).

func sqlJenisReasuransi(t string) string {
	return fmt.Sprintf(`SELECT ID, NOTE FROM %s WHERE FLAG = :1 ORDER BY ID DESC`, t)
}

func sqlCariReinsurer(t string) string {
	return fmt.Sprintf(`SELECT ID, CLIENTNAME FROM %s
	 WHERE ID LIKE :1 AND UPPER(CLIENTNAME) LIKE :2 ESCAPE '\' AND STATUSACTIVE = :3
	 ORDER BY CLIENTNAME ASC, ID ASC FETCH FIRST %d ROWS ONLY`, t, BatasPilihan)
}

// sqlAmbilReinsurer - satu reinsurer menurut ID, dengan kolom saringan (life, aktif) yang layanan
// periksa untuk pilihan BARU (`models.MasterReinsurer.Life`/`Aktif`).
func sqlAmbilReinsurer(t string) string {
	return fmt.Sprintf(`SELECT ID, CLIENTNAME, STATUSACTIVE FROM %s WHERE ID = :1`, t)
}

func sqlCariBusiness(t string) string {
	return fmt.Sprintf(`SELECT ID, NOTE, OLDID FROM %s
	 WHERE OLDID LIKE :1 AND UPPER(NOTE) LIKE :2 ESCAPE '\'
	 ORDER BY ID ASC FETCH FIRST %d ROWS ONLY`, t, BatasPilihan)
}

func sqlAmbilBusiness(t string) string {
	return fmt.Sprintf(`SELECT ID, NOTE, OLDID FROM %s WHERE ID = :1`, t)
}

// BatasRate - `BrowseRateLife_RD` b729 `pyMaxRecords` 500. Satu baris lebih dibaca untuk mengetahui
// daftar terpotong (Pega memotong diam-diam).
const BatasRate = 500

func sqlCariRingkasanRate(t string) string {
	return fmt.Sprintf(`SELECT ID, USEDBY FROM %s
	 WHERE UPPER(USEDBY) LIKE :1 ESCAPE '\'
	 ORDER BY ID ASC FETCH FIRST %d ROWS ONLY`, t, BatasPilihan)
}

func sqlAmbilRingkasanRate(t string) string {
	return fmt.Sprintf(`SELECT ID, USEDBY FROM %s WHERE ID = :1`, t)
}

func sqlDaftarRate(t string) string {
	return fmt.Sprintf(`SELECT ID, USEDBY, GENDER, CONTRACT, AGE, RATE FROM %s
	 WHERE IDUSEDBY = :1
	 ORDER BY ID DESC, RATE ASC FETCH FIRST %d ROWS ONLY`, t, BatasRate+1)
}

// Argumen saringan tetap - nilai VERBATIM RD (flag dan status di models).
const (
	idReinsurerLife = "%" + models.PenandaReinsurerLife + "%" // `BrowseCedingCoLife_RD` b565 Contains "L0"
	awalanBizLife   = models.AwalanBusinessLife + "%"         // `BrowseBusinessLife_RD` b651 StartsWith "L"
)

// PolaCari merakit pola LIKE ber-escape untuk "Contains" RD.
func PolaCari(kata string) string {
	k := strings.ToUpper(strings.TrimSpace(kata))
	k = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(k)
	return "%" + k + "%"
}

func (g *Gudang) bacaMaster(ctx context.Context, objek string, susun func(string) string, kolom []string,
	args ...any) ([]barisTeks, error) {

	q, err := g.siapkan(objek, susun)
	if err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, MasterTidakTerbaca(objek, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []barisTeks
	for rows.Next() {
		b, err := pindai(rows, kolom)
		if err != nil {
			return nil, MasterTidakTerbaca(objek, err)
		}
		hasil = append(hasil, b)
	}
	if err := rows.Err(); err != nil {
		return nil, MasterTidakTerbaca(objek, err)
	}
	return hasil, nil
}

// JenisReasuransiLife - lima jenis life (`FLAG = '1'`, "1 for life").
func (g *Gudang) JenisReasuransiLife(ctx context.Context) ([]models.JenisReasuransi, error) {
	bb, err := g.bacaMaster(ctx, MasterJenisReasuransi, sqlJenisReasuransi, []string{"ID", "NOTE"},
		models.FlagJenisReasuransiLife)
	if err != nil {
		return nil, err
	}
	hasil := make([]models.JenisReasuransi, 0, len(bb))
	for _, b := range bb {
		hasil = append(hasil, models.JenisReasuransi{ID: b.s("ID"), Note: b.s("NOTE")})
	}
	return hasil, nil
}

// CariMasterReinsurer - pilihan `REINSURER NAME` / `SECURITY REINSURER NAME`.
func (g *Gudang) CariMasterReinsurer(ctx context.Context, kata string) ([]models.MasterReinsurer, error) {
	bb, err := g.bacaMaster(ctx, MasterReinsurer, sqlCariReinsurer, []string{"ID", "CLIENTNAME"},
		idReinsurerLife, PolaCari(kata), models.StatusMasterReinsurerAktif)
	if err != nil {
		return nil, err
	}
	hasil := make([]models.MasterReinsurer, 0, len(bb))
	for _, b := range bb {
		hasil = append(hasil, models.MasterReinsurer{ID: b.s("ID"), ClientName: b.s("CLIENTNAME")})
	}
	return hasil, nil
}

// AmbilMasterReinsurer - satu reinsurer master menurut ID, TANPA menyaring: baris lama yang
// reinsurernya kini nonaktif tetap dapat disunting; saringan life/aktif untuk pilihan BARU
// diperiksa layanan dari kolom yang dibaca di sini.
func (g *Gudang) AmbilMasterReinsurer(ctx context.Context, id string) (models.MasterReinsurer, bool, error) {
	bb, err := g.bacaMaster(ctx, MasterReinsurer, sqlAmbilReinsurer, []string{"ID", "CLIENTNAME", "STATUSACTIVE"}, id)
	if err != nil || len(bb) == 0 {
		return models.MasterReinsurer{}, false, err
	}
	return models.MasterReinsurer{ID: bb[0].s("ID"), ClientName: bb[0].s("CLIENTNAME"),
		StatusActive: bb[0].s("STATUSACTIVE")}, true, nil
}

// CariMasterBusiness - pilihan `BUSINESS NAME`.
func (g *Gudang) CariMasterBusiness(ctx context.Context, kata string) ([]models.MasterBusiness, error) {
	bb, err := g.bacaMaster(ctx, MasterBusiness, sqlCariBusiness, []string{"ID", "NOTE", "OLDID"},
		awalanBizLife, PolaCari(kata))
	if err != nil {
		return nil, err
	}
	hasil := make([]models.MasterBusiness, 0, len(bb))
	for _, b := range bb {
		hasil = append(hasil, models.MasterBusiness{ID: b.s("ID"), Note: b.s("NOTE"), OldID: b.s("OLDID")})
	}
	return hasil, nil
}

// AmbilMasterBusiness - satu business master menurut ID.
func (g *Gudang) AmbilMasterBusiness(ctx context.Context, id string) (models.MasterBusiness, bool, error) {
	bb, err := g.bacaMaster(ctx, MasterBusiness, sqlAmbilBusiness, []string{"ID", "NOTE", "OLDID"}, id)
	if err != nil || len(bb) == 0 {
		return models.MasterBusiness{}, false, err
	}
	return models.MasterBusiness{ID: bb[0].s("ID"), Note: bb[0].s("NOTE"), OldID: bb[0].s("OLDID")}, true, nil
}

// CariRingkasanRate - pilihan autocomplete `R/I RATE`.
func (g *Gudang) CariRingkasanRate(ctx context.Context, kata string) ([]models.RingkasanRate, error) {
	bb, err := g.bacaMaster(ctx, MasterRingkasanRate, sqlCariRingkasanRate, []string{"ID", "USEDBY"}, PolaCari(kata))
	if err != nil {
		return nil, err
	}
	hasil := make([]models.RingkasanRate, 0, len(bb))
	for _, b := range bb {
		hasil = append(hasil, models.RingkasanRate{ID: b.s("ID"), UsedBy: b.s("USEDBY")})
	}
	return hasil, nil
}

// AmbilRingkasanRate - satu tabel rate menurut ID (pilihan BARU diperiksa layanan).
func (g *Gudang) AmbilRingkasanRate(ctx context.Context, id string) (models.RingkasanRate, bool, error) {
	bb, err := g.bacaMaster(ctx, MasterRingkasanRate, sqlAmbilRingkasanRate, []string{"ID", "USEDBY"}, id)
	if err != nil || len(bb) == 0 {
		return models.RingkasanRate{}, false, err
	}
	return models.RingkasanRate{ID: bb[0].s("ID"), UsedBy: bb[0].s("USEDBY")}, true, nil
}

// DaftarRate - section `Rate List` satu `IDUSEDBY` (= `RIRATEID`), paling banyak BatasRate baris;
// `terpotong` benar bila view memuat lebih.
func (g *Gudang) DaftarRate(ctx context.Context, idUsedBy string) ([]models.BarisRate, bool, error) {
	kolom := []string{"ID", "USEDBY", "GENDER", "CONTRACT", "AGE", "RATE"}
	bb, err := g.bacaMaster(ctx, MasterRate, sqlDaftarRate, kolom, idUsedBy)
	if err != nil {
		return nil, false, err
	}
	terpotong := len(bb) > BatasRate
	if terpotong {
		bb = bb[:BatasRate]
	}
	hasil := make([]models.BarisRate, 0, len(bb))
	for _, b := range bb {
		hasil = append(hasil, models.BarisRate{ID: b.s("ID"), UsedBy: b.s("USEDBY"), Gender: b.s("GENDER"),
			Contract: b.s("CONTRACT"), Age: b.s("AGE"), Rate: b.s("RATE")})
	}
	return hasil, terpotong, nil
}
