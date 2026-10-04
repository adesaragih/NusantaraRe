package repository

// Rincian satu peserta (`PL_DetailAction` → `PL_Detail_Sec`, retro `RetroLife`
// → `RetroDetailLife`) dan isi popup polis lama (`ViewOldPolicy_EDM*`).
//
// Dibaca sesudah: edm_baca.go.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
)

// --- rincian peserta ----------------------------------------------------------

func sqlRincianPeserta(peserta string) string {
	var b strings.Builder
	b.WriteString("SELECT d.ID, d.PARENT_ID, d.EDM_STATUS")
	for _, k := range models.KolomPesertaRinci {
		b.WriteString(", " + ekspresiBaca("d", k))
	}
	fmt.Fprintf(&b, " FROM %s d WHERE d.ID = :1 AND d.PREMIUM_LIST_ID = :2", peserta)
	return b.String()
}

// sqlSpreadingPeserta - grid `Treaty Type` / `Retroceded Share`
// (`PL_Detail_Sec.xml` b11001/b11147) dan `RetroDetailLife` b1205 sekaligus.
// Kunci `IDX_PLS_DETAIL` (DETAIL_ID).
func sqlSpreadingPeserta(spreading, retro string) string {
	return fmt.Sprintf(`SELECT s.ID, s.TREATY_TYPE_NAME,
	        TO_CHAR(s.RETROCADED_SHARE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''),
	        r.ID, r.REINSURER_NAME,
	        TO_CHAR(r.PERCENT_SHARE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''),
	        TO_CHAR(r.AMOUNT, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''),
	        TO_CHAR(r.RATE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''),
	        TO_CHAR(r.PREMIUM_SPREADED_GROSS, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''),
	        TO_CHAR(r.COMMISION, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''),
	        TO_CHAR(r.OVR_COMM, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''),
	        TO_CHAR(r.PREMIUM_SPREADED_NET, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')
	   FROM %s s LEFT JOIN %s r ON r.SPREADING_ID = s.ID
	  WHERE s.DETAIL_ID = :1
	  ORDER BY s.TREATY_TYPE_NAME, s.ID, r.REINSURER_NAME, r.ID`, spreading, retro)
}

// RincianPeserta membaca satu peserta kasus beserta spreading dan retronya.
func (g *Gudang) RincianPeserta(ctx context.Context, kasusID, pesertaID string) (models.RincianPeserta, error) {
	n, err := g.nama(tabelPeserta, tabelSpreading, tabelSpreadingRetro)
	if err != nil {
		return models.RincianPeserta{}, err
	}
	q := sqlRincianPeserta(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return models.RincianPeserta{}, err
	}
	rows, err := g.db.QueryContext(ctx, q, pesertaID, kasusID)
	if err != nil {
		return models.RincianPeserta{}, fmt.Errorf("repository: membaca rincian peserta %q: %w", pesertaID, err)
	}
	daftar, err := pindaiPeserta(rows, models.KolomPesertaRinci)
	_ = rows.Close()
	if err != nil {
		return models.RincianPeserta{}, err
	}
	if len(daftar) == 0 {
		return models.RincianPeserta{}, ErrTidakAda
	}
	hasil := models.RincianPeserta{Peserta: daftar[0], Spreading: []models.Spreading{}}
	qs := sqlSpreadingPeserta(n[1], n[2])
	if err := db.PeriksaSQL(qs); err != nil {
		return models.RincianPeserta{}, err
	}
	srows, err := g.db.QueryContext(ctx, qs, pesertaID)
	if err != nil {
		return models.RincianPeserta{}, fmt.Errorf("repository: membaca spreading peserta %q: %w", pesertaID, err)
	}
	defer func() { _ = srows.Close() }()
	indeks := map[string]int{}
	for srows.Next() {
		var v [12]sql.NullString
		ptr := make([]any, len(v))
		for i := range v {
			ptr[i] = &v[i]
		}
		if err := srows.Scan(ptr...); err != nil {
			return models.RincianPeserta{}, fmt.Errorf("repository: memindai spreading: %w", err)
		}
		i, ada := indeks[v[0].String]
		if !ada {
			share, err := rapikanAngka("RETROCADED_SHARE", v[2])
			if err != nil {
				return models.RincianPeserta{}, err
			}
			hasil.Spreading = append(hasil.Spreading, models.Spreading{
				ID: v[0].String, TreatyTypeName: v[1].String, RetrocadedShare: share, Retro: []models.SpreadingRetro{},
			})
			i = len(hasil.Spreading) - 1
			indeks[v[0].String] = i
		}
		if !v[3].Valid {
			continue
		}
		var angka [7]string
		for j, nama := range []string{"PERCENT_SHARE", "AMOUNT", "RATE", "PREMIUM_SPREADED_GROSS", "COMMISION", "OVR_COMM", "PREMIUM_SPREADED_NET"} {
			if angka[j], err = rapikanAngka(nama, v[5+j]); err != nil {
				return models.RincianPeserta{}, err
			}
		}
		hasil.Spreading[i].Retro = append(hasil.Spreading[i].Retro, models.SpreadingRetro{
			ReinsurerName: v[4].String, PercentShare: angka[0], Amount: angka[1], Rate: angka[2],
			PremiumSpreadedGross: angka[3], Commision: angka[4], OvrComm: angka[5], PremiumSpreadedNet: angka[6],
		})
	}
	return hasil, srows.Err()
}

// --- popup polis lama -----------------------------------------------------------

// sqlPesertaWarisanHalaman - peserta versi sistem lama, berhalaman.
//
// ⛔ E4: kunci `PL_NUMBER = :1` → `M_LIFE_PREMIUM_DETAIL_INDEX4`, `IDPEGA`
// penyaring sesudahnya, berbatas `FETCH NEXT`. `STNC`/`WPC` DATE di sana.
func sqlPesertaWarisanHalaman(warisan string) string {
	var b strings.Builder
	b.WriteString("SELECT m.ID, NULL, NULL")
	for _, k := range models.KolomPesertaRinci {
		switch {
		case k.Nama == "STNC" || k.Nama == "WPC":
			b.WriteString(", TO_CHAR(m." + k.Nama + ", 'DD/MM/YYYY')")
		case k.Nama == "RISK":
			b.WriteString(", TO_CHAR(m.RISK)")
		default:
			b.WriteString(", " + ekspresiBaca("m", k))
		}
	}
	fmt.Fprintf(&b, ` FROM %s m WHERE m.PL_NUMBER = :1 AND m.IDPEGA = :2
	  ORDER BY m.CERTIFICATE_NO, m.NAME_OF_INSURED, m.ID
	 OFFSET :3 ROWS FETCH NEXT :4 ROWS ONLY`, warisan)
	return b.String()
}

// sqlCacahPesertaWarisan - ⛔ E4: `M_LIFE_PREMIUM_DETAIL_INDEX4` (`PL_NUMBER`).
func sqlCacahPesertaWarisan(warisan string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s m WHERE m.PL_NUMBER = :1 AND m.IDPEGA = :2`, warisan)
}

// sqlRekapPolis - `BrowsePremiumList_RD` (b800 `.PL_NUMBER = Param.PL_NUMBER`,
// urut b1435 `PL_NUMBER DESC`): rekap seluruh versi bernomor PL itu di tabel
// warisan rekap - ditulis PremiumList (NB) dan modul ini (E2) untuk polis
// sistem baru, dan Pega untuk polis sistem lama.
//
// ⚠️ `pyCaseInsensitive=true` di RD tidak ditiru dengan `UPPER(kolom)`: fungsi
// atas kolom mematikan index; nomor PL dirakit sistem, bukan diketik bebas.
func sqlRekapPolis(rekapWarisan string) string {
	kolom := make([]string, len(models.KolomRekapPolisLama))
	for i, k := range models.KolomRekapPolisLama {
		if k == "COB" || k == "PL_NUMBER" || k == "PL_NUMBER_EDM" || k == "CURRENCY" {
			kolom[i] = "r." + k
			continue
		}
		kolom[i] = "TO_CHAR(r." + k + ", 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')"
	}
	return fmt.Sprintf(`SELECT %s FROM %s r WHERE r.PL_NUMBER = :1
	  ORDER BY r.PL_NUMBER DESC, r.PL_NUMBER_EDM DESC NULLS LAST, r.CURRENCY FETCH FIRST 500 ROWS ONLY`,
		strings.Join(kolom, ", "), rekapWarisan)
}

// PesertaVersi membaca satu halaman peserta sebuah versi beserta cacahnya.
func (g *Gudang) PesertaVersi(ctx context.Context, v models.Versi, nomorPolis string, halaman, ukuran int) ([]models.Peserta, int, error) {
	n, err := g.nama(tabelPeserta, tabelPesertaWarisanEDM)
	if err != nil {
		return nil, 0, err
	}
	var q, qc string
	var args, argsC []any
	switch v.Jenis {
	case models.SumberAplikasi:
		q, args = sqlPeserta(n[0], models.KolomPesertaRinci), []any{v.ID, (halaman - 1) * ukuran, ukuran}
		qc, argsC = fmt.Sprintf(`SELECT COUNT(*) FROM %s d WHERE d.PREMIUM_LIST_ID = :1`, n[0]), []any{v.ID}
	case models.SumberWarisan:
		q, args = sqlPesertaWarisanHalaman(n[1]), []any{nomorPolis, v.ID, (halaman - 1) * ukuran, ukuran}
		qc, argsC = sqlCacahPesertaWarisan(n[1]), []any{nomorPolis, v.ID}
	default:
		return nil, 0, fmt.Errorf("repository: sumber versi %q tidak dikenal", v.Jenis)
	}
	for _, s := range []string{q, qc} {
		if err := db.PeriksaSQL(s); err != nil {
			return nil, 0, err
		}
	}
	var total int
	if err := g.db.QueryRowContext(ctx, qc, argsC...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: mencacah peserta versi %q: %w", v.ID, err)
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: membaca peserta versi %q: %w", v.ID, err)
	}
	defer func() { _ = rows.Close() }()
	p, err := pindaiPeserta(rows, models.KolomPesertaRinci)
	return p, total, err
}

// RekapPolis membaca rekap seluruh versi bernomor PL itu (`BrowsePremiumList_RD`).
func (g *Gudang) RekapPolis(ctx context.Context, nomorPolis string) ([]map[string]string, error) {
	n, err := g.nama(tabelRekapWarisan)
	if err != nil {
		return nil, err
	}
	q := sqlRekapPolis(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, nomorPolis)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca rekap polis %q: %w", nomorPolis, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []map[string]string{}
	for rows.Next() {
		v := make([]sql.NullString, len(models.KolomRekapPolisLama))
		ptr := make([]any, len(v))
		for i := range v {
			ptr[i] = &v[i]
		}
		if err := rows.Scan(ptr...); err != nil {
			return nil, fmt.Errorf("repository: memindai rekap polis: %w", err)
		}
		baris := map[string]string{}
		for i, k := range models.KolomRekapPolisLama {
			if i < 4 {
				baris[k] = strings.TrimSpace(v[i].String)
				continue
			}
			s, err := rapikanAngka(k, v[i])
			if err != nil {
				return nil, err
			}
			baris[k] = s
		}
		hasil = append(hasil, baris)
	}
	return hasil, rows.Err()
}
