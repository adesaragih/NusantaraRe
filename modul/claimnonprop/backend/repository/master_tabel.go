package repository

// Untuk apa berkas ini: master treaty Non Prop dari TABEL relasional `T_TREATY_*` modul Treaty In - cadangan
// `MasterTreaty` bila dokumen JSON (M_TREATY_IN / M_TREATY_IN_EDM) tidak ada.
//
// ⭐ Laporan pemakai 11-10-2026 (master Non Prop 1002320): Choose Master berhasil (popup membaca TREATYINDETAIL) tetapi
// halaman `TreatyInMaster` kosong. Penulisan `JSONDATA` `M_TREATY_IN` ditutup pemilik proses 07-10-2026 (ID terakhir
// 1002304); sejak itu Master Treaty In menyimpan kontrak ke `T_TREATY_*`. Pola NB Treaty In
// (`modul/nbtreatyin/backend/repository/masterxol_tabel.go`) dan Claim Prop (membaca tabel flat Treaty In).
//
//   - BACA-SAJA: SELECT saja, nol INSERT / UPDATE / DELETE.
//   - Bentuk keluaran SAMA dengan `bacaMaster` (JSON): medan yang sama, urutan baris = `URUTAN`, anak per `IDINDUK`.
//   - Dipanggil HANYA bila dokumen JSON tidak ada, jadi master lama tetap dibaca persis seperti sebelumnya.
//   - Angka tersimpan berformat tampilan Indonesia di sebagian baris (DEV 11-10-2026: Limit / Deductible master
//     1001855, 1002307, 1002309 "7.100.000.000,00"; RNMSHARE "2,5") diluruskan `angkaMaster` - dokumen JSON selalu
//     memuat angka polos, dan mesin XoL (`models.AngkaTeks`) menolak format tampilan.
//
// Sumber per medan (nama kolom = ejaan dokumen, `PetaPendaratan` Treaty In):
//
//	header            T_TREATY_REVISION `MASTERID = IDMaster` (baris terakhir) - kolom `medanMaster`;
//	                  AccountingMode kosong -> AccountingModeNonProp (`modePembukuanMaster`)
//	Limits            T_TREATY_LIMITS
//	  .TreatyGroups   T_TREATY_LIMIT_GROUP (anak Limits, KunciAnak TreatyGroupList)
//	  .MDPList        T_TREATY_LIMIT_MEASURE `JENIS = 'MDPList'`
//	LimitSummaryList  T_TREATY_LIMIT_SUMMARY (LIMITVAL dst. tidak dibaca rule Claim Non Prop)
//	CurrencyList      T_TREATY_KURS -> TREATYEXCHANGEYEARLY (`Conversion` = TOIDR, `CurrencyList` dokumen adalah salinan
//	                  tabel ini - pendaratan_akar.go Treaty In); kontrak tanpa baris kurs sendiri: kurs tahun treaty
//	                  yang periodenya mencakup Commencement
//	Share             T_TREATY_SHARE; jenis spreading kosong + rincian satu baris -> jenis baris itu
//	                  (`lengkapiJenisSpreading`)
//	  .SpreadingListXOL  T_TREATY_SHARE_SPREADING (anak Share)

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"nusantarare/modul/claimnonprop/backend/models"
)

// masterDariTabel - `MasterTreaty` dari tabel `T_TREATY_*`. Nol baris T_TREATY_REVISION = tidak ada.
func (a *Acuan) masterDariTabel(ctx context.Context, id string) (models.MasterTreaty, bool, error) {
	rev, err := a.q("T_TREATY_REVISION")
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	kol := make([]string, len(medanMaster))
	for i, m := range medanMaster {
		kol[i] = strings.ToUpper(m)
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT %s FROM %s WHERE MASTERID = :1 ORDER BY ID DESC
		FETCH FIRST 1 ROWS ONLY`, strings.Join(kol, ", "), rev), len(kol), id)
	if err != nil || len(rows) == 0 {
		return models.MasterTreaty{}, false, err
	}
	r := rows[0]
	m := models.MasterTreaty{ID: id, ProportionType: r[0], Ceding: r[1], CedingID: r[2], LeadingReinsSource: r[3],
		LeadingReinsSourceID: r[4], Bordeaux: r[5], BordereauxNote: r[6], AccountingMode: modePembukuanMaster(r[7], r[8]),
		AccountingModeNonProp: r[8], TeritorialScope: r[9], Commencement: tglPega(r[10]), Termination: tglPega(r[11]),
		TreatyYear: r[12], RNMShare: angkaMaster(r[13]), EDMState: r[14], StatusAkseptasi: r[15]}

	t, err := a.q("T_TREATY_LIMITS")
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	lim, err := a.banyak(ctx, fmt.Sprintf(`SELECT TO_CHAR(ID), LAYER, LAYERTYPE, LAYERPART, LAYERPARTTYPE, CURRENCY,
		CURRENCY2, LIMIT, LIMIT2, DEDUCTIBLE, DEDUCTIBLE2, REINSTATEMENTPCT, ISCOMBINEMDP, NORIPCALCULATION
		FROM %s WHERE MASTERID = :1 ORDER BY URUTAN, ID`, t), 14, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	idxLimit := map[string]int{}
	for i, l := range lim {
		idxLimit[l[0]] = i
		m.Limits = append(m.Limits, models.LimitXOL{Layer: l[1], LayerType: l[2], LayerPart: l[3], LayerPartType: l[4],
			Currency: l[5], Currency2: l[6], Limit: angkaMaster(l[7]), Limit2: angkaMaster(l[8]),
			Deductible: angkaMaster(l[9]), Deductible2: angkaMaster(l[10]), ReinstatementPct: angkaMaster(l[11]),
			IsCombineMDP: l[12], NoRIPCalculation: l[13]})
	}
	if t, err = a.q("T_TREATY_LIMIT_GROUP"); err != nil {
		return models.MasterTreaty{}, false, err
	}
	grup, err := a.banyak(ctx, fmt.Sprintf(`SELECT TO_CHAR(IDINDUK), TREATYGROUP FROM %s
		WHERE MASTERID = :1 AND TREATYGROUP IS NOT NULL ORDER BY IDINDUK, URUTAN, ID`, t), 2, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, g := range grup {
		if i, ok := idxLimit[g[0]]; ok {
			m.Limits[i].TreatyGroups = append(m.Limits[i].TreatyGroups, g[1])
		}
	}
	if t, err = a.q("T_TREATY_LIMIT_MEASURE"); err != nil {
		return models.MasterTreaty{}, false, err
	}
	mdp, err := a.banyak(ctx, fmt.Sprintf(`SELECT TO_CHAR(IDINDUK), CURRENCY, VALUE FROM %s
		WHERE MASTERID = :1 AND JENIS = 'MDPList' AND (CURRENCY IS NOT NULL OR VALUE IS NOT NULL)
		ORDER BY IDINDUK, URUTAN, ID`, t), 3, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, x := range mdp {
		if i, ok := idxLimit[x[0]]; ok {
			m.Limits[i].MDPList = append(m.Limits[i].MDPList, models.NilaiMataUang{Currency: x[1], Value: angkaMaster(x[2])})
		}
	}

	if t, err = a.q("T_TREATY_LIMIT_SUMMARY"); err != nil {
		return models.MasterTreaty{}, false, err
	}
	ring, err := a.banyak(ctx, fmt.Sprintf(`SELECT LAYER, LAYERTYPE, MDP, MDP2, CURRENCY, CURRENCY2 FROM %s
		WHERE MASTERID = :1 ORDER BY URUTAN, ID`, t), 6, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, s := range ring {
		m.LimitSummaryList = append(m.LimitSummaryList, models.RingkasanLimit{Layer: s[0], LayerType: s[1],
			MDP: angkaMaster(s[2]), MDP2: angkaMaster(s[3]), Currency: s[4], Currency2: s[5]})
	}

	if m.CurrencyList, err = a.kursMasterTabel(ctx, id, m.TreatyYear, r[10]); err != nil {
		return models.MasterTreaty{}, false, err
	}

	if t, err = a.q("T_TREATY_SHARE"); err != nil {
		return models.MasterTreaty{}, false, err
	}
	share, err := a.banyak(ctx, fmt.Sprintf(`SELECT TO_CHAR(ID), SPREADINGTYPEIDXOL, SPREADINGTYPEXOL,
		SPREADINGTOTALPCTXOL FROM %s WHERE MASTERID = :1 ORDER BY URUTAN, ID`, t), 4, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	idxShare := map[string]int{}
	for i, s := range share {
		idxShare[s[0]] = i
		m.Share = append(m.Share, models.ShareXOL{SpreadingTypeIDXOL: s[1], SpreadingTypeXOL: s[2],
			SpreadingTotalPctXOL: angkaMaster(s[3])})
	}
	if t, err = a.q("T_TREATY_SHARE_SPREADING"); err != nil {
		return models.MasterTreaty{}, false, err
	}
	spr, err := a.banyak(ctx, fmt.Sprintf(`SELECT TO_CHAR(IDINDUK), REINSTYPEID, REINSTYPENAME, PCT FROM %s
		WHERE MASTERID = :1 AND (REINSTYPEID IS NOT NULL OR REINSTYPENAME IS NOT NULL)
		ORDER BY IDINDUK, URUTAN, ID`, t), 4, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, s := range spr {
		if i, ok := idxShare[s[0]]; ok {
			m.Share[i].SpreadingListXOL = append(m.Share[i].SpreadingListXOL,
				models.SpreadingMaster{ReinsTypeID: s[1], ReinsTypeName: s[2], Pct: angkaMaster(s[3])})
		}
	}
	for i := range m.Share {
		lengkapiJenisSpreading(&m.Share[i])
	}
	return m, true, nil
}

// lengkapiJenisSpreading - Share tanpa jenis spreading yang rinciannya SATU baris memakai jenis baris itu.
//
// ⭐ 11-10-2026 (laporan pemakai, kasus CLMNP-000017 / master 1002320): kolom Treaty Type Spreading List klaim kosong.
// Master Pega ber-ORS menyimpan `SpreadingTypeIDXOL = 10007`, `SpreadingTypeXOL = ORS` DAN rincian ORS (59 dokumen
// DEV); Master Treaty In kini menyimpan ORS yang dipilih lewat Reins Type hanya di rincian (`T_TREATY_SHARE_SPREADING`)
// dengan jenis Share kosong (1002311-1002320). Rincian lebih dari satu baris tanpa jenis bernama (mis. "2024 QS 155M
// TRT") dibiarkan - nama jenis spreading tidak dapat disimpulkan dari rinciannya.
func lengkapiJenisSpreading(sh *models.ShareXOL) {
	if strings.TrimSpace(sh.SpreadingTypeIDXOL) != "" || strings.TrimSpace(sh.SpreadingTypeXOL) != "" ||
		len(sh.SpreadingListXOL) != 1 {
		return
	}
	sh.SpreadingTypeIDXOL, sh.SpreadingTypeXOL = sh.SpreadingListXOL[0].ReinsTypeID, sh.SpreadingListXOL[0].ReinsTypeName
}

// kursMasterTabel - `CurrencyList` master dari tabel: baris kurs milik kontrak (T_TREATY_KURS, urut URUTAN; kunci
// TREATYEXCHANGEYEARLY = ID + TREATYYEAR karena ID tidak unik antartahun). Kontrak tanpa baris milik sendiri (sebelum
// tabel hubung ada): kurs tahun treaty yang periodenya mencakup Commencement ("yyyyMMdd").
func (a *Acuan) kursMasterTabel(ctx context.Context, id, tahun, mulai string) ([]models.KursMaster, error) {
	hub, err := a.q("T_TREATY_KURS")
	if err != nil {
		return nil, err
	}
	kurs, err := a.q("TREATYEXCHANGEYEARLY")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT k.CURRENCY, k.TOIDR FROM %s h JOIN %s k ON k.ID = h.IDKURS
		WHERE h.MASTERID = :1 AND k.TREATYYEAR = :2 ORDER BY h.URUTAN`, hub, kurs), 2, id, strings.TrimSpace(tahun))
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		rows, err = a.banyak(ctx, fmt.Sprintf(`SELECT CURRENCY, TOIDR FROM %s WHERE TREATYYEAR = :1
			AND SUBSTR(STARTDATE, 1, 8) <= :2 AND SUBSTR(ENDDATE, 1, 8) >= :3 ORDER BY CURRENCY, ID`, kurs), 2,
			strings.TrimSpace(tahun), strings.TrimSpace(mulai), strings.TrimSpace(mulai))
		if err != nil {
			return nil, err
		}
	}
	var out []models.KursMaster
	for _, k := range rows {
		out = append(out, models.KursMaster{Currency: k[0], Conversion: angkaMaster(k[1])})
	}
	return out, nil
}

// ribuanTitik - bilangan bulat bertitik ribuan minimal dua kelompok ("7.100.000"); satu kelompok ("1.500") tetap
// dibaca titik desimal.
var ribuanTitik = regexp.MustCompile(`^-?\d{1,3}(\.\d{3}){2,}$`)

// angkaMaster meluruskan angka berformat tampilan Indonesia ("7.100.000.000,00", "2,5") menjadi angka polos
// ("7100000000.00", "2.5"). Teks lain dikembalikan apa adanya.
func angkaMaster(s string) string {
	t := strings.TrimSpace(s)
	switch {
	case strings.Contains(t, ","):
		return strings.ReplaceAll(strings.ReplaceAll(t, ".", ""), ",", ".")
	case ribuanTitik.MatchString(t):
		return strings.ReplaceAll(t, ".", "")
	}
	return s
}

// modePembukuanMaster - `TreatyInMaster.AccountingMode` master dari tabel.
//
// ⭐ 11-10-2026 (laporan pemakai, kasus CLMNP-000017 / master 1002320): Accounting Mode Claim Treaty kosong. Layar
// Master Treaty In cabang Non Prop hanya punya `AccountingModeNonProp` (loss / risk; `AccountingMode` accounting /
// underwriting milik cabang Proportional - pyCondition TreatyInNONProportional), jadi ACCOUNTINGMODE master Non Prop
// baru selalu kosong (14 dari 14 di DEV). Yang tampil = Accounting Mode yang diisi di master. Master yang punya
// AccountingMode (dokumen lama) tidak berubah. Medan ini hanya tampilan; logika loss / risk klaim tetap membaca
// AccountingModeNonProp.
func modePembukuanMaster(mode, modeNonProp string) string {
	if strings.TrimSpace(mode) != "" {
		return mode
	}
	return modeNonProp
}
