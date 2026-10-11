package models

// Untuk apa berkas ini: PERLUASAN TANGGA - `ApprovalKomite_Act` S3-S8 (keputusan work owner 10-10-2026 KCF-02 dan
// pita SPV B KCF-01). Kasus komite TT2 lahir dengan tangga calon `SetListKomite_act` (non-retro: roster ber-LIMIT_BOTTOM
// <= 1 = Claim Supervisor); di tingkat 1 tangga diperluas dengan jenjang di atasnya menurut total adjustment IDR.
//
//	S3     Fac Retro -> keluar (tangga retro dibentuk sisi klaim)
//	S4     KomiteCount = 1: roster FACIN aktif DEGREE > 1 ber-LIMIT_BOTTOM < total, urut DEGREE
//	S5     pemutus SPV B dan 30.000.000 < total <= 57.750.000: HANYA jenjang DEGREE > 1 pertama
//	S6.2   buang calon ber-OPERATOR_ID = KomiteList(1).KomiteID
//	S7     KomiteCount = 1 dan tangga satu baris: calon ditambahkan (IDKomite = JABATAN, approval 0)
//	S8     KomiteLoop = cacah tangga
//
// Dihitung saat tingkat 1 membuka kasus (ditampilkan, tanpa menulis saat GET) dan disimpan saat tingkat 1 Submit
// (KCF-02). PERBAIKAN prompt §5 butir 2: S7 `@LengthOfPageList(...) =1` (tanda `=` tunggal) dibaca pembandingan.

import (
	"sort"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// Batas pita SPV B (ApprovalKomite_Act S5, konstanta bernama KCF-01).
var (
	// BatasBawahPitaSPVB - `Local.TotalAdj > 30000000.00` (eksklusif).
	BatasBawahPitaSPVB = apd.New(30000000, 0)
	// BatasAtasPitaSPVB - `Local.TotalAdj <= 57750000.00` (inklusif).
	BatasAtasPitaSPVB = apd.New(57750000, 0)
)

// AnggotaRoster - satu baris roster EMAILKOMITE STS_KLAIM FACIN aktif (Obj-Browse GetKomite).
type AnggotaRoster struct {
	ID, OperatorID, Email, Jabatan, Degree string
	// LimitBottom - teks angka (NULL = kosong: saringan `LIMIT_BOTTOM < :total` SQL tidak meloloskannya).
	LimitBottom string
}

func derajat(s string) (int, bool) {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	return n, err == nil
}

// PerluasTangga = ApprovalKomite_Act S3-S8: anggota tangga BARU (urut lanjutan) untuk kasus `k`; nil = tanpa perluasan.
// `pemutusSPVB` = pemutus tingkat 1 anggota ReasClaimSPVB (`AnggotaSPVB`).
func PerluasTangga(k Kasus, pr PraProses, roster []AnggotaRoster, pemutusSPVB bool) ([]Anggota, error) {
	if k.TransferType != TransferAdjustment || pr.Retro || k.Count != 1 { // S1 SetValueKomite, S3, S4 / S7
		return nil, nil
	}
	if len(k.Tangga) != 1 { // S7 `@LengthOfPageList(pyWorkPage.KomiteList) = 1`
		return nil, nil
	}
	total := pr.TotalAdj
	if total == nil {
		total = apd.New(0, 0)
	}
	type calon struct {
		AnggotaRoster
		deg int
	}
	var atas []calon // DEGREE > 1 (S4 / S5.2)
	for _, r := range roster {
		d, ok := derajat(r.Degree)
		if !ok || d <= 1 {
			continue
		}
		atas = append(atas, calon{r, d})
	}
	sort.SliceStable(atas, func(i, j int) bool { return atas[i].deg < atas[j].deg })
	var pilih []calon
	if pemutusSPVB && lebih(total, BatasBawahPitaSPVB) && !lebih(total, BatasAtasPitaSPVB) { // S5
		if len(atas) > 0 {
			pilih = atas[:1] // S5.2 MaxRecords 1 (KCF-01: DEGREE > 1 pertama)
		}
	} else { // S4
		for _, c := range atas {
			if strings.TrimSpace(c.LimitBottom) == "" {
				continue
			}
			lb, err := desimal("LIMIT_BOTTOM", c.LimitBottom)
			if err != nil {
				return nil, err
			}
			if lb.Cmp(total) < 0 { // IsLessThan LIMIT_BOTTOM < Local.TotalAdj
				pilih = append(pilih, c)
			}
		}
	}
	satu := k.Tangga[0].OperatorID
	urut := k.Tangga[0].Urut
	var out []Anggota
	for _, c := range pilih {
		if c.OperatorID == satu { // S6.2
			continue
		}
		urut++
		out = append(out, Anggota{Urut: urut, OperatorID: c.OperatorID, Jabatan: c.Jabatan, Email: c.Email,
			Keputusan: KeputusanMenunggu}) // S7.1
	}
	return out, nil
}
