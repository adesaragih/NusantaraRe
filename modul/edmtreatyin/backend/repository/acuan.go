package repository

// Asal: salinan sebagian modul/nbtreatyin/backend/repository/acuan.go (06-10-2026) - pembacaan view kontrak
// TREATYINDETAILJOINEDM (NB Choose Business) DIBUANG: EDM membaca master dari JSONDATA M_TREATY_IN /
// M_TREATY_IN_EDM / M_TREATY_OUT (`SetValueEDM_Act`, pengecualian K8 - masterxol.go).
//
// Untuk apa berkas ini: TABEL ACUAN warisan yang dibaca layar dan aktivitas EDM (satu fungsi = satu rule; nama
// rulenya di atas tiap fungsi). Seluruhnya HANYA MEMBACA.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
)

// ------------------------------------------------------------------ satu nilai

// satuTeks menjalankan kueri satu kolom satu baris; nol baris = "" (perilaku
// `pxResults(1).X` atas hasil kosong di Pega).
func (g *Gudang) satuTeks(ctx context.Context, apa, q string, args ...any) (string, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var v sql.NullString
	err := g.db.QueryRowContext(ctx, q, args...).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: %s: %w", apa, err)
	}
	return v.String, nil
}

// IDMataUangDariNama = RDB `GetCurrencyIDByName` (`SetTreatyCurrencyID`).
func (g *Gudang) IDMataUangDariNama(ctx context.Context, nama string) (string, error) {
	t, err := g.nama(tabelMataUang)
	if err != nil {
		return "", err
	}
	return g.satuTeks(ctx, "membaca ID mata uang",
		fmt.Sprintf(`SELECT TO_CHAR(ID) FROM %s WHERE CURRENCY = :1 FETCH FIRST 1 ROWS ONLY`, t), nama)
}

// NamaMataUang = RDB `GetCurrency` (`SetCurrency_act`).
func (g *Gudang) NamaMataUang(ctx context.Context, id string) (string, error) {
	t, err := g.nama(tabelMataUang)
	if err != nil {
		return "", err
	}
	return g.satuTeks(ctx, "membaca nama mata uang",
		fmt.Sprintf(`SELECT CURRENCY FROM %s WHERE TO_CHAR(ID) = :1 FETCH FIRST 1 ROWS ONLY`, t), id)
}

// StsPKPAgen = RD `BrowseClientName_RD` (`SetPPNPPH` langkah 1-3): filter
// D `.ID = Param.ID` (SourceOfBusiness) dan E `.StatusActive IS NULL`.
// `[data DBA - belum dikonfirmasi]` nama kolom STS_PKP / STATUSACTIVE tabel
// AGENT diambil dari nama properti RD kelas `ASM-FW-GISFW-Int-AGENT`.
func (g *Gudang) StsPKPAgen(ctx context.Context, sobID string) (string, error) {
	if sobID == "" {
		return "", nil
	}
	t, err := g.nama(tabelAgen)
	if err != nil {
		return "", err
	}
	return g.satuTeks(ctx, "membaca status PKP agen",
		fmt.Sprintf(`SELECT TO_CHAR(STS_PKP) FROM %s WHERE TO_CHAR(ID) = :1 AND STATUSACTIVE IS NULL FETCH FIRST 1 ROWS ONLY`, t), sobID)
}

// BisnisDariKunci = RDB `GetOldIDBusiness_SQL`:
// `where (ID = CARI2 OR NOTE = CARI2) AND GROUPPANEL IS NOT NULL`.
// ⚠️ RDB tanpa ORDER BY (baris pertama urutan basis data); di sini ORDER BY ID
// supaya hasilnya tetap - penyimpangan kecil, dicatat.
func (g *Gudang) BisnisDariKunci(ctx context.Context, kunci string) (models.BarisBisnis, error) {
	t, err := g.nama(tabelBisnis)
	if err != nil {
		return models.BarisBisnis{}, err
	}
	q := fmt.Sprintf(`SELECT TO_CHAR(OLDID), TO_CHAR(GROUPPANEL), TO_CHAR(ID) FROM %s
	  WHERE (TO_CHAR(ID) = :1 OR NOTE = :2) AND GROUPPANEL IS NOT NULL ORDER BY ID FETCH FIRST 1 ROWS ONLY`, t)
	if err := db.PeriksaSQL(q); err != nil {
		return models.BarisBisnis{}, err
	}
	var a, b, c sql.NullString
	err = g.db.QueryRowContext(ctx, q, kunci, kunci).Scan(&a, &b, &c)
	if errors.Is(err, sql.ErrNoRows) {
		return models.BarisBisnis{}, nil
	}
	if err != nil {
		return models.BarisBisnis{}, fmt.Errorf("repository: membaca bisnis: %w", err)
	}
	return models.BarisBisnis{OldID: a.String, GroupPanel: b.String, ID: c.String}, nil
}

// MO = `CheckDataMkt` langkah 3 (`Obj-Browse` marketing officer, `.ID =
// QuotationData.MOID`).
func (g *Gudang) MO(ctx context.Context, id string) (models.BarisMO, error) {
	t, err := g.nama(tabelMO)
	if err != nil {
		return models.BarisMO{}, err
	}
	q := fmt.Sprintf(`SELECT TO_CHAR(ID), TO_CHAR(CLIENTID), CLIENTNAME, TO_CHAR(TEAMGROUP), TO_CHAR(BRANCHDETAILID), BRANCHDETAILNAME
	  FROM %s WHERE TO_CHAR(ID) = :1 FETCH FIRST 1 ROWS ONLY`, t)
	if err := db.PeriksaSQL(q); err != nil {
		return models.BarisMO{}, err
	}
	var v [6]sql.NullString
	err = g.db.QueryRowContext(ctx, q, id).Scan(&v[0], &v[1], &v[2], &v[3], &v[4], &v[5])
	if errors.Is(err, sql.ErrNoRows) {
		return models.BarisMO{}, nil
	}
	if err != nil {
		return models.BarisMO{}, fmt.Errorf("repository: membaca marketing officer: %w", err)
	}
	return models.BarisMO{ID: v[0].String, ClientID: v[1].String, ClientName: v[2].String,
		TeamGroup: v[3].String, BranchDetailID: v[4].String, BranchDetailName: v[5].String}, nil
}

// ------------------------------------------------------------------ pilihan layar

func (g *Gudang) daftarPilihan(ctx context.Context, apa, q string, args ...any) ([]models.Pilihan, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: %s: %w", apa, err)
	}
	defer rows.Close()
	var out []models.Pilihan
	for rows.Next() {
		var n, l sql.NullString
		if err := rows.Scan(&n, &l); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", apa, err)
		}
		out = append(out, models.Pilihan{Nilai: n.String, Label: l.String})
	}
	return out, rows.Err()
}

// DaftarMataUang = RD `BrowseCurrencyTreatyIn_RD`: filter B `.Currency !=
// "ITL"` (AC 54); filter A ber-parameter kosong diabaikan.
func (g *Gudang) DaftarMataUang(ctx context.Context) ([]models.Pilihan, error) {
	t, err := g.nama(tabelMataUang)
	if err != nil {
		return nil, err
	}
	return g.daftarPilihan(ctx, "membaca daftar mata uang",
		fmt.Sprintf(`SELECT TO_CHAR(ID), CURRENCY FROM %s WHERE CURRENCY <> :1 ORDER BY CURRENCY`, t), "ITL")
}

// DaftarMO = RD `BrowseMarketingOfficer_RD`: filter A `.MOStatus = 1`.
func (g *Gudang) DaftarMO(ctx context.Context) ([]models.Pilihan, error) {
	t, err := g.nama(tabelMO)
	if err != nil {
		return nil, err
	}
	return g.daftarPilihan(ctx, "membaca daftar marketing officer",
		fmt.Sprintf(`SELECT TO_CHAR(ID), CLIENTNAME FROM %s WHERE MOSTATUS = :1 ORDER BY UPPER(CLIENTNAME), ID`, t), "1")
}

// DaftarJenisReas = RD `BrowseReinsuranceType_RD` (dropdown TreatyType layar
// atasan); keempat filternya ber-parameter, tidak diisi layar -> diabaikan.
func (g *Gudang) DaftarJenisReas(ctx context.Context) ([]models.Pilihan, error) {
	t, err := g.nama(tabelJenisReas)
	if err != nil {
		return nil, err
	}
	return g.daftarPilihan(ctx, "membaca daftar jenis reasuransi",
		fmt.Sprintf(`SELECT TO_CHAR(ID), NOTE FROM %s ORDER BY ID`, t))
}

// ------------------------------------------------------------------ duplikat
