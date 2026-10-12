package repository

// Untuk apa berkas ini: master kontrak treaty jalur NB NonProporsional / XOL dari
// TABEL relasional `T_TREATY_*` modul treatyin - cadangan `MasterXOLDariJSON`.
//
// ⭐ Laporan pemakai 11-10-2026 (kasus NB-22454, master 1002318): Choose Business
// atas master Non-Prop memunculkan "master treaty (M_TREATY_IN / M_TREATY_IN_EDM)
// kontrak ini tidak ditemukan". Penulisan `JSONDATA` `M_TREATY_IN` ditutup pemilik
// proses 07-10-2026 (ID terakhir 1002304); sejak itu Master Treaty In menyimpan
// kontrak ke `T_TREATY_*`. Jalur Proportional tidak terkena karena hanya membaca
// view `TREATYINDETAILJOINEDM`. Permintaan pemakai: pembacaan master Non-Prop
// disamakan - tanpa galat untuk master baru.
//
//   - BACA-SAJA: SELECT saja, nol INSERT/UPDATE/DELETE.
//   - Bentuk keluaran SAMA dengan `MasterXOLDariJSON` (kunci relatif `TreatyIn`,
//     `models.JalurAnak` untuk daftar bersarang) dan medannya tetap DAFTAR
//     TERTUTUP `models.SkalarMasterXOL` / `models.DaftarMasterXOL` - anggota yang
//     tidak punya kolom (mis. `NetPremiAfterPPN`) tidak diisi, sama dengan dokumen
//     yang tidak memuatnya.
//   - Dipanggil HANYA bila dokumen JSON tidak ada (`services.penyimpanOracle`), jadi
//     master lama tetap dibaca persis seperti sebelumnya.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// sumberDaftarMaster - satu PageList master di tabel pendaratan treatyin.
type sumberDaftarMaster struct {
	tabel string
	// jenis - nilai kolom `JENIS` (tabel gabungan beberapa larik); "" = tanpa saringan.
	jenis string
	// kolom - anggota skema -> kolom tabel. Anggota skema tanpa kolom tidak diisi.
	kolom map[string]string
	// anak - daftar bersarang, berkunci `IDINDUK` = `ID` baris induk.
	anak map[string]sumberDaftarMaster
}

var (
	kolomMataNilai = map[string]string{"Currency": "CURRENCY", "Value": "VALUE"}
	kolomDeduksi   = map[string]string{"Currency": "CURRENCY", "Deduction": "DEDUCTION"}
	kolomRingkas   = map[string]string{
		"Note": "NOTE", "Limit": "LIMITVAL", "Limit2": "LIMITVAL2", "MDP": "MDP", "MDP2": "MDP2",
		"Deductible": "DEDUCTIBLE", "Deductible2": "DEDUCTIBLE2", "NetPremi": "NETPREMI", "NetPremi2": "NETPREMI2",
	}
)

// totalMaster - larik total per mata uang di `T_TREATY_TOTAL` (kolom `JENIS` = nama larik).
func totalMaster(jenis string) sumberDaftarMaster {
	return sumberDaftarMaster{tabel: "T_TREATY_TOTAL", jenis: jenis, kolom: kolomMataNilai}
}

// SumberDaftarMasterXOL - setiap PageList `models.DaftarMasterXOL` dan tabelnya.
// Dikunci uji `TestSumberMasterTabelMeliputiDaftarMaster`.
var SumberDaftarMasterXOL = map[string]sumberDaftarMaster{
	"Share": {
		tabel: "T_TREATY_SHARE",
		kolom: map[string]string{
			"LayerType": "LAYERTYPE", "Layer": "LAYER", "LayerPartType": "LAYERPARTTYPE", "LayerPart": "LAYERPART",
			"SpreadingTypeXOL": "SPREADINGTYPEXOL", "SpreadingTypeIDXOL": "SPREADINGTYPEIDXOL",
		},
		anak: map[string]sumberDaftarMaster{
			"GrossPremiumList":   {tabel: "T_TREATY_SHARE_AMOUNT", jenis: "GrossPremiumList", kolom: kolomMataNilai},
			"NetPremiumList":     {tabel: "T_TREATY_SHARE_AMOUNT", jenis: "NetPremiumList", kolom: kolomMataNilai},
			"DeductionTotalList": {tabel: "T_TREATY_SHARE_AMOUNT", jenis: "DeductionTotalList", kolom: kolomMataNilai},
			"RnmLimitList":       {tabel: "T_TREATY_SHARE_AMOUNT", jenis: "RnmLimitList", kolom: kolomMataNilai},
			"DeductionList":      {tabel: "T_TREATY_SHARE_DEDUCTION", kolom: kolomDeduksi},
			"SpreadingListXOL": {tabel: "T_TREATY_SHARE_SPREADING", kolom: map[string]string{
				"ReinsTypeID": "REINSTYPEID", "ReinsTypeName": "REINSTYPENAME", "Pct": "PCT",
			}},
		},
	},
	"Installment": {
		tabel: "T_TREATY_INSTALLMENT",
		kolom: map[string]string{"Currency": "CURRENCY", "AmountTotal": "AMOUNTTOTAL", "PctTotal": "PCTTOTAL"},
		anak: map[string]sumberDaftarMaster{
			"InstallmentList": {tabel: "T_TREATY_INSTALLMENT_ITEM", kolom: map[string]string{
				"Installment": "INSTALLMENT", "PaymentDate": "PAYMENTDATE", "InstallmentPct": "INSTALLMENTPCT",
				"Amount": "AMOUNT", "Currency": "CURRENCY",
			}},
		},
	},
	"FacultativeShareList": {
		tabel: "T_TREATY_FAC_SHARE",
		kolom: map[string]string{},
		anak: map[string]sumberDaftarMaster{
			"GrossPremiumList":   {tabel: "T_TREATY_FAC_SHARE_AMOUNT", jenis: "GrossPremiumList", kolom: kolomMataNilai},
			"DeductionTotalList": {tabel: "T_TREATY_FAC_SHARE_AMOUNT", jenis: "DeductionTotalList", kolom: kolomMataNilai},
			"DeductionList":      {tabel: "T_TREATY_FAC_SHARE_DEDUCTION", kolom: kolomDeduksi},
		},
	},
	"Limits": {
		tabel: "T_TREATY_LIMITS",
		kolom: map[string]string{
			"ReinstatementValue": "REINSTATEMENTVALUE", "ReinstatementNote": "REINSTATEMENTNOTE",
			"Limit": "LIMIT", "Limit2": "LIMIT2",
		},
		anak: map[string]sumberDaftarMaster{
			"MDPList":            {tabel: "T_TREATY_LIMIT_MEASURE", jenis: "MDPList", kolom: kolomMataNilai},
			"Reinstatement_List": {tabel: "T_TREATY_LIMIT_REINSTATEMENT", kolom: map[string]string{"ReinstatementValue": "REINSTATEMENTVALUE"}},
		},
	},
	"LimitShareSummaryList":    {tabel: "T_TREATY_SHARE_SUMMARY", jenis: "LimitShareSummaryList", kolom: kolomRingkas},
	"LimitFacShareSummaryList": {tabel: "T_TREATY_SHARE_SUMMARY", jenis: "LimitFacShareSummaryList", kolom: kolomRingkas},
	"LimitSummaryList": {tabel: "T_TREATY_LIMIT_SUMMARY", kolom: map[string]string{
		"Note": "NOTE", "Limit": "LIMITVAL", "Limit2": "LIMITVAL2", "Deductible": "DEDUCTIBLE",
		"Deductible2": "DEDUCTIBLE2", "MDP": "MDP", "MDP2": "MDP2",
	}},
	"TotalLimitIOONP":          totalMaster("TotalLimitIOONP"),
	"TotalLimitDeductblNP":     totalMaster("TotalLimitDeductblNP"),
	"TotalLimitMDPNP":          totalMaster("TotalLimitMDPNP"),
	"TotalShareRnmNP":          totalMaster("TotalShareRnmNP"),
	"TotalShareGrossNP":        totalMaster("TotalShareGrossNP"),
	"TotalShareDeductionNP":    totalMaster("TotalShareDeductionNP"),
	"TotalShareNetNP":          totalMaster("TotalShareNetNP"),
	"TotalSpreadedNetPremi":    totalMaster("TotalSpreadedNetPremi"),
	"TotalSpreadedNetPremiRI":  totalMaster("TotalSpreadedNetPremiRI"),
	"TotalFacShareRnmNP":       totalMaster("TotalFacShareRnmNP"),
	"TotalFacShareGrossNP":     totalMaster("TotalFacShareGrossNP"),
	"TotalFacShareDeductionNP": totalMaster("TotalFacShareDeductionNP"),
	"TotalFacShareNetNP":       totalMaster("TotalFacShareNetNP"),
}

// KolomSkalarMasterXOL - `models.SkalarMasterXOL` -> kolom `T_TREATY_REVISION`.
var KolomSkalarMasterXOL = map[string]string{
	"RNMShare": "RNMSHARE", "RnmShareDeducted": "RNMSHAREDEDUCTED", "FacultativeShare": "FACULTATIVESHARE",
	"EDMState": "EDMSTATE", "ProportionType": "PROPORTIONTYPE",
}

const tabelRevisiMaster = "T_TREATY_REVISION"

// barisMasterTabel - satu baris tabel pendaratan: pengenal, induk, dan nilai anggota.
type barisMasterTabel struct {
	id, induk string
	nilai     models.Baris
}

// MasterXOLDariTabel membaca master satu kontrak (`noOffer` = `MASTERID`) dari
// tabel `T_TREATY_*`. Nol baris revisi DAN nol baris Share/Limits =
// `ErrMasterXOLTidakAda`.
func (g *Gudang) MasterXOLDariTabel(ctx context.Context, noOffer string) (models.MasterXOL, error) {
	m := models.MasterXOL{Nilai: map[string]string{}, Daftar: map[string][]models.Baris{}}
	adaRevisi, err := g.skalarMasterTabel(ctx, noOffer, m.Nilai)
	if err != nil {
		return models.MasterXOL{}, err
	}
	for nama, skema := range models.DaftarMasterXOL {
		sumber, ada := SumberDaftarMasterXOL[nama]
		if !ada {
			continue
		}
		induk, err := g.barisMasterTabel(ctx, sumber, noOffer, nama, skema.Anggota, false)
		if err != nil {
			return models.MasterXOL{}, err
		}
		if len(induk) == 0 {
			continue
		}
		baris := make([]models.Baris, 0, len(induk))
		urutan := map[string]int{}
		for i, b := range induk {
			baris = append(baris, b.nilai)
			urutan[b.id] = i + 1
		}
		m.Daftar[nama] = baris
		for anak, anggota := range skema.Anak {
			sa, ada := sumber.anak[anak]
			if !ada {
				continue
			}
			isi, err := g.barisMasterTabel(ctx, sa, noOffer, nama+"."+anak, anggota, true)
			if err != nil {
				return models.MasterXOL{}, err
			}
			per := map[int][]models.Baris{}
			for _, b := range isi {
				if n, ok := urutan[b.induk]; ok {
					per[n] = append(per[n], b.nilai)
				}
			}
			for n, sub := range per {
				m.Daftar[models.JalurAnak(nama, n, anak)] = sub
			}
		}
	}
	if !adaRevisi && len(m.Daftar["Share"]) == 0 && len(m.Daftar["Limits"]) == 0 {
		return models.MasterXOL{}, fmt.Errorf("%w: %s", ErrMasterXOLTidakAda, noOffer)
	}
	return m, nil
}

// skalarMasterTabel mengisi `models.SkalarMasterXOL` dari `T_TREATY_REVISION`;
// kolom NULL tidak diisi (sama dengan dokumen yang tidak memuat kuncinya).
func (g *Gudang) skalarMasterTabel(ctx context.Context, noOffer string, nilai map[string]string) (bool, error) {
	nama, err := g.nama(tabelRevisiMaster)
	if err != nil {
		return false, err
	}
	kolom := make([]string, 0, len(models.SkalarMasterXOL))
	for _, k := range models.SkalarMasterXOL {
		kolom = append(kolom, KolomSkalarMasterXOL[k])
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE MASTERID = :1`, strings.Join(kolom, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	sel := make([]sql.NullString, len(kolom))
	tuju := make([]any, len(kolom))
	for i := range sel {
		tuju[i] = &sel[i]
	}
	err = g.db.QueryRowContext(ctx, q, noOffer).Scan(tuju...)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("repository: membaca master treaty (tabel revisi): %w", err)
	}
	for i, k := range models.SkalarMasterXOL {
		if sel[i].Valid {
			nilai[k] = sel[i].String
		}
	}
	return true, nil
}

// barisMasterTabel membaca baris satu sumber, urut `URUTAN` (anak: per induk).
func (g *Gudang) barisMasterTabel(ctx context.Context, s sumberDaftarMaster, noOffer, jalur string, anggota []string, anak bool) ([]barisMasterTabel, error) {
	nama, err := g.nama(s.tabel)
	if err != nil {
		return nil, err
	}
	var dipakai []string
	pilih := []string{"TO_CHAR(ID)"}
	if anak {
		pilih = append(pilih, "TO_CHAR(IDINDUK)")
	} else {
		pilih = append(pilih, "''")
	}
	for _, a := range anggota {
		if k, ada := s.kolom[a]; ada {
			dipakai = append(dipakai, a)
			pilih = append(pilih, k)
		}
	}
	args := []any{noOffer}
	saring := "MASTERID = :1"
	if s.jenis != "" {
		args = append(args, s.jenis)
		saring += " AND JENIS = :2"
	}
	urut := "URUTAN"
	if anak {
		urut = "IDINDUK, URUTAN"
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE %s ORDER BY %s`, strings.Join(pilih, ", "), nama, saring, urut)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca master treaty (%s): %w", s.tabel, err)
	}
	defer rows.Close()
	var out []barisMasterTabel
	for rows.Next() {
		sel := make([]sql.NullString, len(pilih))
		tuju := make([]any, len(pilih))
		for i := range sel {
			tuju[i] = &sel[i]
		}
		if err := rows.Scan(tuju...); err != nil {
			return nil, fmt.Errorf("repository: membaca master treaty (%s): %w", s.tabel, err)
		}
		b := barisMasterTabel{id: sel[0].String, induk: sel[1].String, nilai: models.Baris{}}
		for i, a := range dipakai {
			v := sel[2+i]
			if !v.Valid {
				continue
			}
			t := v.String
			if models.TanggalMasterXOL[jalur+"."+a] {
				t = models.TanggalMasterPega(t)
			}
			b.nilai[a] = t
		}
		out = append(out, b)
	}
	return out, rows.Err()
}
