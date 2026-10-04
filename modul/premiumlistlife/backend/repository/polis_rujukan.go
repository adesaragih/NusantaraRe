package repository

// Pencarian master untuk layar Input Offer - tiket 01 bagian 3.
//
// Untuk apa berkas ini: kedua popup pilihan layar `InputOfferLife.xml` -
// `Ceding_Harness` (`Section/Ceding_Section.xml`) dan `PolicyHolder_Harness`
// (`Section/PolicyHolder_Section.xml`). Keduanya dibaca saja.
//
// `[terverifikasi]` grid kedua section dan report definition-nya:
//
//	Ceding_Section       pyRDName BrowseCedingCoLife_RD (ASM-FW-GISFW-Int-AGENT)
//	                     param CedingCoLeader = SearchPolicyHolder.CARI1 (teks cari)
//	                     tombol Choose -> setCeding_act (clientid=.ID, clientname=.ClientName)
//	BrowseCedingCoLife_RD  logika "B AND A AND C":
//	                     B .ID Contains "L0" · A .ClientName Contains Param.CedingCoLeader
//	                     C .StatusActive = 1
//	PolicyHolder_Section pyRDName BrowseClientNusaRe_RD (ASM-FW-GISFW-Int-CLIENT)
//	                     param PolicyHolderName = CARI1, Business = "LIFE INSURANCE"
//	                     tombol Choose -> setPolicyHolder_act (.ID, .Name)
//	BrowseClientNusaRe_RD  "A AND B AND C AND D":
//	                     A .Name Contains param (pyCaseInsensitive true) · B .Name != "-"
//	                     C .Name IS NOT NULL · D .BU_Note = Param.Business
//	SearchPolicyHolder_act  satu langkah: CARI1 = @toUpperCase(CARI1)
//
// ⛔ Tabel fisik:
//   - Ceding -> `AGENT` (`ID`, `CLIENTNAME`, `STATUSACTIVE`): kelas Int-AGENT
//     yang SAMA dipetakan ke `POOLDATA.AGENT` oleh Claim Life (PARITAS §5.2,
//     `[data DBA — katalog DEV]`) dan kolomnya dipakai Treaty Contract Out
//     (`MasterReinsurerAgentTCO`, OQ-TCO-12 ditutup work owner).
//   - Policy Holder -> `CLIENT`: kelas Int-CLIENT tercatat di discovery sebagai
//     `POOLDATA.CLIENT` (NB FacIn `GetInsuredID.xml`). ⚠️ `[belum terverifikasi]`
//     NAMA KOLOM `NAME` dan `BU_NOTE` - diturunkan dari nama properti RD
//     (`.Name`, `.BU_Note`) dengan pola yang terbukti di AGENT (`.ClientName`
//     -> `CLIENTNAME`). Dikonfirmasi work owner / DBA sebelum dipakai di DEV.
//
// ⛔ Nama tabel dari KONSTANTA, tidak pernah dari permintaan; teks cari lewat
// penanda dengan ESCAPE (pola `polaLikeTCO`).
//
// ⛔ Nilai `CLIENTNAME`/`NAME` dapat memuat nama; dibaca saat jalan, nol baris
// disalin ke fixture, tiket, maupun log.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

const (
	// MasterCedingPolis - master ceding, kelas Int-AGENT.
	MasterCedingPolis = "AGENT"
	// MasterPemegangPolis - master pemegang polis, kelas Int-CLIENT.
	MasterPemegangPolis = "CLIENT"
	// NilaiCedingAktif - BrowseCedingCoLife_RD filter C `pyFilterValue 1`.
	NilaiCedingAktif = "1"
	// PenandaIDCeding - BrowseCedingCoLife_RD filter B `pyFilterValue "L0"`.
	PenandaIDCeding = "L0"
	// BisnisPemegangPolisLife - PolicyHolder_Section param `Business`.
	BisnisPemegangPolisLife = "LIFE INSURANCE"
	// NamaPemegangTanpaIsi - BrowseClientNusaRe_RD filter B `!= "-"`.
	NamaPemegangTanpaIsi = "-"
)

// batasCariRujukanPolis membatasi hasil popup.
//
// ⚠️ RD aslinya ber-`pyMaxRecords` 10000/100000; grid popup menampilkan
// halaman. Batas 100 pola `batasCariReinsurerTCO` - pemakai mempersempit
// dengan mengetik, bukan dengan menggulir sepuluh ribu baris.
const batasCariRujukanPolis = 100

// BarisRujukan adalah satu baris pilihan popup.
type BarisRujukan struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
	// Keterangan - kolom `Business` grid Policy Holder (`.BU_Note`); kosong
	// untuk ceding.
	Keterangan string `json:"keterangan,omitempty"`
}

// Rujukan membaca master popup layar Input Offer.
type Rujukan struct{ db *db.DB }

// NewRujukan menyusunnya.
func NewRujukan(db *db.DB) *Rujukan { return &Rujukan{db: db} }

// polaLikeRujukan membungkus teks ketikan menjadi pola `Contains` yang aman.
//
// ⛔ Dibesarkan lebih dahulu - SearchPolicyHolder_act `@toUpperCase(CARI1)`.
func polaLikeRujukan(teks string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + r.Replace(strings.ToUpper(strings.TrimSpace(teks))) + "%"
}

func sqlCariCeding(tabel string) string {
	return fmt.Sprintf(`SELECT ID, CLIENTNAME FROM %s
	  WHERE STATUSACTIVE = :1 AND ID LIKE :2 ESCAPE '\' AND UPPER(CLIENTNAME) LIKE :3 ESCAPE '\'
	  ORDER BY CLIENTNAME, ID
	  FETCH FIRST %d ROWS ONLY`, tabel, batasCariRujukanPolis)
}

func sqlCariPemegangPolis(tabel string) string {
	return fmt.Sprintf(`SELECT ID, NAME, BU_NOTE FROM %s
	  WHERE NAME IS NOT NULL AND NAME <> :1 AND BU_NOTE = :2 AND UPPER(NAME) LIKE :3 ESCAPE '\'
	  ORDER BY NAME, ID
	  FETCH FIRST %d ROWS ONLY`, tabel, batasCariRujukanPolis)
}

// CariCeding membaca ceding aktif yang namanya memuat teks.
func (r *Rujukan) CariCeding(ctx context.Context, teks string) ([]BarisRujukan, error) {
	return r.cari(ctx, MasterCedingPolis, sqlCariCeding,
		NilaiCedingAktif, "%"+PenandaIDCeding+"%", polaLikeRujukan(teks))
}

// CariPemegangPolis membaca pemegang polis Life yang namanya memuat teks.
func (r *Rujukan) CariPemegangPolis(ctx context.Context, teks string) ([]BarisRujukan, error) {
	return r.cari(ctx, MasterPemegangPolis, sqlCariPemegangPolis,
		NamaPemegangTanpaIsi, BisnisPemegangPolisLife, polaLikeRujukan(teks))
}

func (r *Rujukan) cari(ctx context.Context, master string, rakit func(string) string, arg ...any) (
	[]BarisRujukan, error) {

	tabel, err := r.db.Qualify(master)
	if err != nil {
		return nil, err
	}
	q := rakit(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca master %s: %w", master, err)
	}
	defer func() { _ = rows.Close() }()
	kolom, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	hasil := []BarisRujukan{}
	for rows.Next() {
		n := make([]sql.NullString, len(kolom))
		tujuan := make([]any, len(n))
		for i := range n {
			tujuan[i] = &n[i]
		}
		if err := rows.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: memindai master %s: %w", master, err)
		}
		b := BarisRujukan{ID: n[0].String, Nama: n[1].String}
		if len(n) > 2 {
			b.Keterangan = n[2].String
		}
		hasil = append(hasil, b)
	}
	return hasil, rows.Err()
}

// ——— Layar Input Premium Detail (tiket 03 bagian 2) ———
//
// `[terverifikasi]` sumber ketiganya:
//
//	Product Name     ChooseProdName -> BrowseProductForNB_Life (Int-PRODUCT_LIFE),
//	                 param Ceding = pyWorkPage.CedingCo -> filter `.CEDINGID = Param.Ceding`
//	                 (`.POLICYHODER = Param.PolicyHolder` - param tidak diisi section)
//	Marketing Officer BrowseMarketingOfficer_RD (Int-marketingofficer), filter
//	                 `.MOStatus = 1`; pilihan: .ID -> MOID, .ClientID -> MarketingCode,
//	                 .ClientName -> MarketingName
//	R/I SLIP RNM No. GetPLandNopolis_sql: `NOPOLIS LIKE '%RNML-Q%'` UNION
//	                 `NOPOLIS LIKE '%RNML-F%'` dari JSON_POLIS
//
// ⛔ Tabel fisik:
//   - `PRODUCT_LIFE` `[terverifikasi]` nama tabelnya (RDBList GetProductDtlPL,
//     kelas yang sama). ⚠️ `[belum terverifikasi]` kolom CEDINGID, SOBID,
//     SOBNAME, POLICYHODER, POLICYHODERNAME, INWARDNAME - nama properti RD
//     (seluruhnya huruf besar, pola nama kolom), satu-satunya join RD ke tipe
//     produk. GetProductDtlPL hanya membuktikan ID, PRODUCTNAME, CEDING.
//   - `MARKETINGOFFICER` - Claim Life PARITAS §5.2 `[data DBA - katalog DEV]`,
//     `CLIENTNAME` terbukti. ⚠️ `[belum terverifikasi]` `CLIENTID`, `MOSTATUS`.
//   - R/I SLIP dari `T_PREMIUM_LIST.NO_POLIS` baris utama: `JSON_POLIS` dibuang
//     (spec §12). Polis warisan baru muncul setelah migrasi data (tiket 09).

const (
	// MasterProdukLife - master produk inward Life.
	MasterProdukLife = "PRODUCT_LIFE"
	// MasterMarketingOfficer - master marketing officer.
	MasterMarketingOfficer = "MARKETINGOFFICER"
	// NilaiMarketingAktif - BrowseMarketingOfficer_RD `pyFilterValue 1`.
	NilaiMarketingAktif = "1"
)

// PolaNomorRISlip - kedua pola `GetPLandNopolis_sql`, VERBATIM.
var PolaNomorRISlip = []string{"%RNML-Q%", "%RNML-F%"}

// BarisProduk adalah satu baris popup Choose Product Name.
type BarisProduk struct {
	ID               string `json:"id"`
	InwardName       string `json:"inwardName"`
	CedingID         string `json:"cedingId"`
	Ceding           string `json:"ceding"`
	SobID            string `json:"sobId"`
	SobName          string `json:"sobName"`
	PolicyHolder     string `json:"policyHolder"`
	PolicyHolderName string `json:"policyHolderName"`
}

// BarisMarketing adalah satu pilihan Marketing Officer.
type BarisMarketing struct {
	ID   string `json:"id"`
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

func sqlCariProduk(tabel string) string {
	return fmt.Sprintf(`SELECT ID, INWARDNAME, CEDINGID, CEDING, SOBID, SOBNAME, POLICYHODER, POLICYHODERNAME
	   FROM %s WHERE CEDINGID = :1 AND UPPER(NVL(INWARDNAME, ' ')) LIKE :2 ESCAPE '\'
	  ORDER BY INWARDNAME, ID
	  FETCH FIRST %d ROWS ONLY`, tabel, batasCariRujukanPolis)
}

func sqlCariMarketing(tabel string) string {
	return fmt.Sprintf(`SELECT ID, CLIENTID, CLIENTNAME FROM %s
	  WHERE MOSTATUS = :1 AND UPPER(CLIENTNAME) LIKE :2 ESCAPE '\'
	  ORDER BY CLIENTNAME, ID
	  FETCH FIRST %d ROWS ONLY`, tabel, batasCariRujukanPolis)
}

// sqlCariRISlip - baris utama saja (`ID = ID_PEGA`), nomor unik.
func sqlCariRISlip(tabel string) string {
	return fmt.Sprintf(`SELECT DISTINCT NO_POLIS FROM %s
	  WHERE ID = ID_PEGA AND (NO_POLIS LIKE :1 OR NO_POLIS LIKE :2)
	    AND UPPER(NO_POLIS) LIKE :3 ESCAPE '\'
	  ORDER BY NO_POLIS
	  FETCH FIRST %d ROWS ONLY`, tabel, batasCariRujukanPolis)
}

// CariProduk membaca produk milik ceding kasus.
func (r *Rujukan) CariProduk(ctx context.Context, cedingID, teks string) ([]BarisProduk, error) {
	tabel, err := r.db.Qualify(MasterProdukLife)
	if err != nil {
		return nil, err
	}
	q := sqlCariProduk(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, cedingID, polaLikeRujukan(teks))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca master %s: %w", MasterProdukLife, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []BarisProduk{}
	for rows.Next() {
		var n [8]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7]); err != nil {
			return nil, fmt.Errorf("repository: memindai master %s: %w", MasterProdukLife, err)
		}
		hasil = append(hasil, BarisProduk{ID: n[0].String, InwardName: n[1].String, CedingID: n[2].String,
			Ceding: n[3].String, SobID: n[4].String, SobName: n[5].String,
			PolicyHolder: n[6].String, PolicyHolderName: n[7].String})
	}
	return hasil, rows.Err()
}

// CariMarketing membaca marketing officer aktif yang namanya memuat teks.
func (r *Rujukan) CariMarketing(ctx context.Context, teks string) ([]BarisMarketing, error) {
	tabel, err := r.db.Qualify(MasterMarketingOfficer)
	if err != nil {
		return nil, err
	}
	q := sqlCariMarketing(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, NilaiMarketingAktif, polaLikeRujukan(teks))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca master %s: %w", MasterMarketingOfficer, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []BarisMarketing{}
	for rows.Next() {
		var n [3]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2]); err != nil {
			return nil, fmt.Errorf("repository: memindai master %s: %w", MasterMarketingOfficer, err)
		}
		hasil = append(hasil, BarisMarketing{ID: n[0].String, Kode: n[1].String, Nama: n[2].String})
	}
	return hasil, rows.Err()
}

// CariRISlip membaca nomor polis RNML-Q / RNML-F yang memuat teks.
func (r *Rujukan) CariRISlip(ctx context.Context, teks string) ([]BarisRujukan, error) {
	tabel, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return nil, err
	}
	q := sqlCariRISlip(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, PolaNomorRISlip[0], PolaNomorRISlip[1], polaLikeRujukan(teks))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca nomor polis R/I SLIP: %w", err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []BarisRujukan{}
	for rows.Next() {
		var no sql.NullString
		if err := rows.Scan(&no); err != nil {
			return nil, err
		}
		hasil = append(hasil, BarisRujukan{ID: no.String, Nama: no.String})
	}
	return hasil, rows.Err()
}
