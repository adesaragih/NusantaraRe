package repository

// Rincian Product Name - tombol "View" di samping Product Name layar Input
// Premium Detail (permintaan work owner 05-10-2026: "tombol atau apapun yang
// bisa untuk melihat apa isi dari product name ini").
//
// Untuk apa berkas ini: membaca SATU produk dari tabel flat Master Product Name
// Life - `M_PRODUCTNAME_LIFE` + lima tabel anaknya - untuk ditampilkan saja.
//
// `[keputusan work owner 05-10-2026]` sumbernya tabel flat Master Product Name
// Life (yang ditulis menu itu), BUKAN view lama `PRODUCT_LIFE` /
// `PRODUCTINWARD_LIFE`. Nama tabel dan kolomnya = `modul/masterproductnamelife/
// backend/repository/mpnl_flat.go` (DDL migrasi 140-147); modul itu tidak
// diimpor - batas modul.
//
// ⛔ BACA-SAJA, kolom DISEBUT satu per satu (nol `SELECT *`), berkunci ID.
// Angka dibaca sebagai TEKS desimal (`db.FmtDesimal`, ADR-0003), tanggal
// `YYYY-MM-DD`.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// Tabel flat Master Product Name Life yang dibaca.
const (
	TabelProdukMPNL        = "M_PRODUCTNAME_LIFE"
	TabelProdukMPNLLien    = "M_PRODUCTNAME_LIFE_LIEN"
	TabelProdukMPNLDokumen = "M_PRODUCTNAME_LIFE_DOCCLAIM"
	TabelProdukMPNLPlan    = "M_PRODUCTNAME_LIFE_PLAN"
	TabelProdukMPNLFinUW   = "M_PRODUCTNAME_LIFE_FINUW"
	TabelProdukMPNLUWLimit = "M_PRODUCTNAME_LIFE_UWLIMIT"
)

// ErrRincianProdukTidakAda - ID produk tidak ada di tabel flat Master Product
// Name Life (mis. produk lama yang belum dipindah).
var ErrRincianProdukTidakAda = errors.New("repository: produk tidak ditemukan di M_PRODUCTNAME_LIFE")

// Jenis nilai rincian - menentukan cara layar menampilkannya.
const (
	JenisRincianTeks    = "teks"
	JenisRincianAngka   = "angka"
	JenisRincianTanggal = "tanggal"
)

// kolomRincian - satu kolom yang dibaca: nama fisik, label layar Master
// Product Name Life, jenis.
type kolomRincian struct{ Nama, Label, Jenis string }

func teksR(n, l string) kolomRincian    { return kolomRincian{n, l, JenisRincianTeks} }
func angkaR(n, l string) kolomRincian   { return kolomRincian{n, l, JenisRincianAngka} }
func tanggalR(n, l string) kolomRincian { return kolomRincian{n, l, JenisRincianTanggal} }

// bagianRincian - satu kelompok medan induk.
type bagianRincian struct {
	Judul string
	Kolom []kolomRincian
}

// BagianRincianProduk - medan induk `M_PRODUCTNAME_LIFE`, urutan layar Master
// Product Name Life (sisi umum lalu sisi inward). Kolom ID dan jejak operator
// tidak ditampilkan.
var BagianRincianProduk = []bagianRincian{
	{Judul: "TREATY NAME", Kolom: []kolomRincian{
		teksR("PRODUCTNAME", "Product Name"), teksR("ID", "Product Code"), teksR("CEDING", "Ceding"),
		teksR("SOBNAME", "SOB"), angkaR("RICOMM", "Deduction (%)"), teksR("RIRISK", "R/I Risk Name"),
		teksR("INWARDNAME", "Treaty Name"), teksR("TREATYNUMBER", "Treaty Number"), teksR("CAUSE", "Cause Of Loss"),
	}},
	{Judul: "INWARD", Kolom: []kolomRincian{
		teksR("POLICYHOLDERNAME", "Policy Holder"), teksR("INSURED", "Insured"),
		angkaR("ADDENDUMNO", "Addendum No."), teksR("ADDENDUMWORD", "Addendum"),
		angkaR("AMANDEMENTNO", "Amandement No."), teksR("AMANDEMENTSCHD", "Amandement"),
		angkaR("MAXEXPIREDCLAIM", "Max Notification Claim Expired"), tanggalR("BEGIN_DATE", "Begin Date"),
		tanggalR("STNC", "STNC"), tanggalR("MATURE", "Expired Date"),
		angkaR("CEDINGRETENTIONNUM", "Ceding Retention (%)"), angkaR("CEDINGLIMIT", "Ceding's Limit"),
		angkaR("BROKERAGE", "Brokerage Fee (%)"), angkaR("MINAGE", "Minimum Age (Years)"),
		angkaR("MAXAGE", "Maximum Age (Years)"), angkaR("EXPIRYAGE", "Expiry Age (Years)"),
		angkaR("EXTRAPREMI", "Extra Premium"), angkaR("MINSUMINSURED", "Min Sum Insured"),
		angkaR("MAXSUMINSURED", "Max Sum Insured"), angkaR("MAXSUMREASURED", "Max Sum Reasured"),
		angkaR("RNMSHARE", "Nusantara Re Share (%)"), angkaR("RNMLIMITNUM", "Nusantara Re's Limit"),
		angkaR("PREMIUMFACTOR", "Premium Factor (%)"), teksR("PAYMENT", "Premium Payment Method"),
		teksR("SUBJECTTO", "Subject To"), angkaR("ANNUITYINTEREST", "Annuity Interest (%)"),
		angkaR("PREMIUMREFUNDFACTOR", "Premium Refund Factor (%)"),
		angkaR("MAXDATARECEIVE", "Max Production Data Receive"), teksR("BIRTHDAY", "Birthday"),
		teksR("CURRENCY", "Currency"), angkaR("EXTRAMORTALITY", "Extra Mortality (%)"),
		angkaR("MAXCONTRACT", "Max Contract (year)"), angkaR("PROPORTIONALTABLE", "Proportional Table"),
	}},
}

// gridRincian - satu tabel anak (berkunci PRODUCTID, urut URUT).
type gridRincian struct {
	Judul, Tabel string
	Kolom        []kolomRincian
	// Kunci - kolom tersembunyi per baris (PLAN LIST: RIRATEID, untuk View Rate).
	Kunci string
}

// GridRincianProduk - kelima grid layar Master Product Name Life.
var GridRincianProduk = []gridRincian{
	{Judul: "PLAN LIST", Tabel: TabelProdukMPNLPlan, Kunci: "RIRATEID", Kolom: []kolomRincian{
		teksR("PLAN", "Plan Name"), teksR("NAME", "Bussines"), teksR("BENEFIT", "Benefit"), teksR("RIRATE", "R/I Rate"),
	}},
	{Judul: "UNDERWRITING LIMIT", Tabel: TabelProdukMPNLUWLimit, Kolom: []kolomRincian{
		angkaR("MININSURED", "Min Insured"), angkaR("MAXINSURED", "Max Insured"), angkaR("MINAGE", "Min Age"),
		angkaR("MAXAGE", "Max Age"), teksR("MEDICAL", "Medical"), teksR("DESCRIPTION", "Description"),
	}},
	{Judul: "FINANCIAL UNDERWRITING", Tabel: TabelProdukMPNLFinUW, Kolom: []kolomRincian{
		angkaR("MININSURED", "Min Insured"), angkaR("MAXINSURED", "Max Insured"),
		teksR("EMPLOYEE", "Employee"), teksR("NON_EMPLOYEE", "Non-Employee"),
	}},
	{Judul: "LIEN CLAUSE (Potongan Manfaat Klaim)", Tabel: TabelProdukMPNLLien, Kolom: []kolomRincian{
		teksR("USIA", "Usia saat Klaim"), teksR("MANFAAT", "% Manfaat yang dibayarkan"),
	}},
	{Judul: "DOCUMENT CLAIM", Tabel: TabelProdukMPNLDokumen, Kolom: []kolomRincian{
		teksR("DOCUMENT", "Document List"),
	}},
}

// MedanRincian - satu medan induk yang dibaca.
type MedanRincian struct {
	Label string `json:"label"`
	Nilai string `json:"nilai"`
	Jenis string `json:"jenis"`
}

// BagianRincian - satu kelompok medan induk.
type BagianRincian struct {
	Judul string         `json:"judul"`
	Medan []MedanRincian `json:"medan"`
}

// KolomGridRincian - kepala satu kolom grid.
type KolomGridRincian struct {
	Label string `json:"label"`
	Jenis string `json:"jenis"`
}

// GridRincian - satu tabel anak beserta barisnya.
type GridRincian struct {
	Judul string             `json:"judul"`
	Kolom []KolomGridRincian `json:"kolom"`
	Baris [][]string         `json:"baris"`
	// Kunci - nilai kolom tersembunyi per baris (PLAN LIST: RIRATEID); kosong
	// untuk grid tanpa kunci.
	Kunci []string `json:"kunci,omitempty"`
}

// RincianProduk - isi satu produk untuk ditampilkan.
type RincianProduk struct {
	Bagian []BagianRincian `json:"bagian"`
	Grid   []GridRincian   `json:"grid"`
	// RIRiskID - kunci View R/I Risk (`RIRISK_LIFE.IDUSEDBY`); tidak ditampilkan.
	RIRiskID string `json:"riRiskId"`
}

// ekspresiRincian - kolom → ekspresi SELECT yang mengembalikan teks.
func ekspresiRincian(k kolomRincian) string {
	switch k.Jenis {
	case JenisRincianAngka:
		return fmt.Sprintf(db.FmtDesimal, k.Nama)
	case JenisRincianTanggal:
		return fmt.Sprintf("TO_CHAR(%s, 'YYYY-MM-DD')", k.Nama)
	default:
		return k.Nama
	}
}

// angkaUtuh - `TO_CHAR(.., 'TM9')` membuang nol di depan titik (`.5`); layar
// menuntut `0.5`.
func angkaUtuh(s string) string {
	switch {
	case strings.HasPrefix(s, "."):
		return "0" + s
	case strings.HasPrefix(s, "-."):
		return "-0" + s[1:]
	}
	return s
}

// PembayaranProduk - kode PAYMENT → teks layar Master Product Name Life
// (`GenerateUpload_Act` b1141). Kode lain tampil apa adanya.
var PembayaranProduk = map[string]string{"1": "Annual", "2": "Semi Annual", "3": "Quarterly", "4": "Monthly"}

func nilaiRincian(k kolomRincian, v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	s := strings.TrimSpace(v.String)
	switch {
	case k.Jenis == JenisRincianAngka:
		return angkaUtuh(s)
	case k.Nama == "PAYMENT" && PembayaranProduk[s] != "":
		return PembayaranProduk[s]
	}
	return s
}

func sqlRincianInduk(tabel string) string {
	var ekspresi []string
	for _, b := range BagianRincianProduk {
		for _, k := range b.Kolom {
			ekspresi = append(ekspresi, ekspresiRincian(k))
		}
	}
	// Kunci tersembunyi View R/I Risk - kolom TERAKHIR.
	ekspresi = append(ekspresi, "RIRISKID")
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, strings.Join(ekspresi, ", "), tabel)
}

func sqlRincianGrid(g gridRincian, tabel string) string {
	ekspresi := make([]string, 0, len(g.Kolom)+1)
	for _, k := range g.Kolom {
		ekspresi = append(ekspresi, ekspresiRincian(k))
	}
	if g.Kunci != "" {
		ekspresi = append(ekspresi, g.Kunci)
	}
	return fmt.Sprintf(`SELECT %s FROM %s WHERE PRODUCTID = :1 ORDER BY URUT`, strings.Join(ekspresi, ", "), tabel)
}

// RincianProduk membaca satu produk Master Product Name Life beserta kelima
// gridnya. `ErrRincianProdukTidakAda` bila ID-nya tidak ada.
func (r *Rujukan) RincianProduk(ctx context.Context, id string) (RincianProduk, error) {
	tabel, err := r.db.Qualify(TabelProdukMPNL)
	if err != nil {
		return RincianProduk{}, err
	}
	q := sqlRincianInduk(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return RincianProduk{}, err
	}
	var semua []kolomRincian
	for _, b := range BagianRincianProduk {
		semua = append(semua, b.Kolom...)
	}
	nilai := make([]sql.NullString, len(semua)+1)
	tujuan := make([]any, len(semua)+1)
	for i := range nilai {
		tujuan[i] = &nilai[i]
	}
	err = r.db.QueryRowContext(ctx, q, id).Scan(tujuan...)
	if errors.Is(err, sql.ErrNoRows) {
		return RincianProduk{}, fmt.Errorf("%w: %s", ErrRincianProdukTidakAda, id)
	}
	if err != nil {
		return RincianProduk{}, fmt.Errorf("repository: membaca %s: %w", TabelProdukMPNL, err)
	}
	hasil := RincianProduk{Bagian: []BagianRincian{}, Grid: []GridRincian{},
		RIRiskID: strings.TrimSpace(nilai[len(semua)].String)}
	i := 0
	for _, b := range BagianRincianProduk {
		bag := BagianRincian{Judul: b.Judul, Medan: []MedanRincian{}}
		for _, k := range b.Kolom {
			bag.Medan = append(bag.Medan, MedanRincian{Label: k.Label, Nilai: nilaiRincian(k, nilai[i]), Jenis: k.Jenis})
			i++
		}
		hasil.Bagian = append(hasil.Bagian, bag)
	}
	for _, g := range GridRincianProduk {
		grid, err := r.gridRincian(ctx, g, id)
		if err != nil {
			return RincianProduk{}, err
		}
		hasil.Grid = append(hasil.Grid, grid)
	}
	return hasil, nil
}

func (r *Rujukan) gridRincian(ctx context.Context, g gridRincian, id string) (GridRincian, error) {
	tabel, err := r.db.Qualify(g.Tabel)
	if err != nil {
		return GridRincian{}, err
	}
	q := sqlRincianGrid(g, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return GridRincian{}, err
	}
	rows, err := r.db.QueryContext(ctx, q, id)
	if err != nil {
		return GridRincian{}, fmt.Errorf("repository: membaca %s: %w", g.Tabel, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := GridRincian{Judul: g.Judul, Kolom: []KolomGridRincian{}, Baris: [][]string{}}
	for _, k := range g.Kolom {
		hasil.Kolom = append(hasil.Kolom, KolomGridRincian{Label: k.Label, Jenis: k.Jenis})
	}
	cacah := len(g.Kolom)
	if g.Kunci != "" {
		cacah++
	}
	for rows.Next() {
		nilai := make([]sql.NullString, cacah)
		tujuan := make([]any, cacah)
		for i := range nilai {
			tujuan[i] = &nilai[i]
		}
		if err := rows.Scan(tujuan...); err != nil {
			return GridRincian{}, fmt.Errorf("repository: memindai %s: %w", g.Tabel, err)
		}
		baris := make([]string, len(g.Kolom))
		for i, k := range g.Kolom {
			baris[i] = nilaiRincian(k, nilai[i])
		}
		hasil.Baris = append(hasil.Baris, baris)
		if g.Kunci != "" {
			hasil.Kunci = append(hasil.Kunci, strings.TrimSpace(nilai[len(g.Kolom)].String))
		}
	}
	return hasil, rows.Err()
}

// ——— View Rate: isi satu R/I Rate baris PLAN LIST (permintaan work owner 05-10-2026) ———
//
// Sama dengan dialog `View Rate` Master Product Name Life (`BrowseRateLife_RD`):
// `M_RATE_LIFE` (dulu view `RATE_LIFE`) berkunci `IDUSEDBY` = RIRATEID baris plan, enam kolom grid
// `ViewRate`, urut `ID DESC, RATE ASC`, paling banyak BatasRateProduk baris.
// BACA-SAJA, nol `SELECT *`, nol `JSONDATA`.

// ViewRateProduk - objek rate yang dibaca. RALAT 07-10-2026 (keputusan work owner, `modul/riratelife/MODUL.md`
// RALAT R7): kini TABEL flat `M_RATE_LIFE` berkolom sama dengan view `RATE_LIFE` lama (dibuang migrasi inti 930), teks
// apa adanya. Nama konstanta tetap.
const ViewRateProduk = "M_RATE_LIFE"

// BatasRateProduk - `pyMaxRecords` `BrowseRateLife_RD`; lebih dari ini dinyatakan terpotong.
const BatasRateProduk = 500

// BarisRateProduk - satu baris rate.
type BarisRateProduk struct {
	ID       string `json:"id"`
	UsedBy   string `json:"usedBy"`
	Gender   string `json:"gender"`
	Contract string `json:"contract"`
	Age      string `json:"age"`
	Rate     string `json:"rate"`
}

// RateProduk - isi satu R/I Rate.
type RateProduk struct {
	Baris     []BarisRateProduk `json:"baris"`
	Terpotong bool              `json:"terpotong"`
}

func sqlRateProduk(tabel string) string {
	return fmt.Sprintf(`SELECT ID, USEDBY, GENDER, CONTRACT, AGE, RATE FROM %s
		WHERE IDUSEDBY = :1
		ORDER BY ID DESC, RATE ASC FETCH FIRST %d ROWS ONLY`, tabel, BatasRateProduk+1)
}

// RateProduk membaca isi satu R/I Rate (`IDUSEDBY` = RIRATEID).
func (r *Rujukan) RateProduk(ctx context.Context, riRateID string) (RateProduk, error) {
	tabel, err := r.db.Qualify(ViewRateProduk)
	if err != nil {
		return RateProduk{}, err
	}
	q := sqlRateProduk(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return RateProduk{}, err
	}
	rows, err := r.db.QueryContext(ctx, q, riRateID)
	if err != nil {
		return RateProduk{}, fmt.Errorf("repository: membaca %s: %w", ViewRateProduk, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := RateProduk{Baris: []BarisRateProduk{}}
	for rows.Next() {
		var n [6]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5]); err != nil {
			return RateProduk{}, fmt.Errorf("repository: memindai %s: %w", ViewRateProduk, err)
		}
		hasil.Baris = append(hasil.Baris, BarisRateProduk{ID: n[0].String, UsedBy: n[1].String,
			Gender: n[2].String, Contract: n[3].String, Age: n[4].String, Rate: n[5].String})
	}
	if err := rows.Err(); err != nil {
		return RateProduk{}, fmt.Errorf("repository: membaca %s: %w", ViewRateProduk, err)
	}
	if len(hasil.Baris) > BatasRateProduk {
		hasil.Baris = hasil.Baris[:BatasRateProduk]
		hasil.Terpotong = true
	}
	return hasil, nil
}

// ——— View R/I Risk: isi R/I Risk Name produk (permintaan work owner 05-10-2026) ———
//
// Tabel `RIRISK_LIFE` berkunci `IDUSEDBY` = `M_PRODUCTNAME_LIFE.RIRISKID` - dulu
// VIEW (semua VARCHAR2, ALL_TAB_COLUMNS POOLDATA, work owner 05-10-2026), sejak
// migrasi inti 938-940 TABEL (tabel Pega M_RIRISK_LIFE berganti nama, keputusan
// work owner 08-10-2026 K1) berkolom sama: ID, IDUSEDBY, USEDBY, AGE, YEAR,
// MONTH, CONTRACT teks; RISK NUMBER (dulu teks "921,9", kini dipindai driver
// sebagai teks angka bertitik - tampilan saja). BACA-SAJA, kolom disebut satu
// per satu, paling banyak BatasRiskProduk baris (lebih dari itu dinyatakan
// terpotong).

// ViewRiskProduk - tabel R/I Risk yang dibaca (nama tetap; dulu view).
const ViewRiskProduk = "RIRISK_LIFE"

// BatasRiskProduk - baris paling banyak satu R/I Risk.
const BatasRiskProduk = 500

// BarisRiskProduk - satu baris R/I Risk.
type BarisRiskProduk struct {
	ID       string `json:"id"`
	UsedBy   string `json:"usedBy"`
	Age      string `json:"age"`
	Year     string `json:"year"`
	Month    string `json:"month"`
	Risk     string `json:"risk"`
	Contract string `json:"contract"`
}

// RiskProduk - isi satu R/I Risk.
type RiskProduk struct {
	Baris     []BarisRiskProduk `json:"baris"`
	Terpotong bool              `json:"terpotong"`
}

func sqlRiskProduk(tabel string) string {
	return fmt.Sprintf(`SELECT ID, USEDBY, AGE, YEAR, MONTH, RISK, CONTRACT FROM %s
		WHERE IDUSEDBY = :1
		ORDER BY ID ASC FETCH FIRST %d ROWS ONLY`, tabel, BatasRiskProduk+1)
}

// RiskProduk membaca isi satu R/I Risk (`IDUSEDBY` = RIRISKID produk).
func (r *Rujukan) RiskProduk(ctx context.Context, riRiskID string) (RiskProduk, error) {
	tabel, err := r.db.Qualify(ViewRiskProduk)
	if err != nil {
		return RiskProduk{}, err
	}
	q := sqlRiskProduk(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return RiskProduk{}, err
	}
	rows, err := r.db.QueryContext(ctx, q, riRiskID)
	if err != nil {
		return RiskProduk{}, fmt.Errorf("repository: membaca %s: %w", ViewRiskProduk, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := RiskProduk{Baris: []BarisRiskProduk{}}
	for rows.Next() {
		var n [7]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6]); err != nil {
			return RiskProduk{}, fmt.Errorf("repository: memindai %s: %w", ViewRiskProduk, err)
		}
		hasil.Baris = append(hasil.Baris, BarisRiskProduk{ID: n[0].String, UsedBy: n[1].String, Age: n[2].String,
			Year: n[3].String, Month: n[4].String, Risk: n[5].String, Contract: n[6].String})
	}
	if err := rows.Err(); err != nil {
		return RiskProduk{}, fmt.Errorf("repository: membaca %s: %w", ViewRiskProduk, err)
	}
	if len(hasil.Baris) > BatasRiskProduk {
		hasil.Baris = hasil.Baris[:BatasRiskProduk]
		hasil.Terpotong = true
	}
	return hasil, nil
}
