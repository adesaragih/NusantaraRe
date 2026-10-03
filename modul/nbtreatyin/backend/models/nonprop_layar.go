package models

// Untuk apa berkas ini: ATURAN LAYAR jalur NonProporsional - kontainer yang
// tersembunyi bagi polis XOL baru, dan daftar yang boleh dikirim layar.
//
// Bukti XML (sel layout, `pyContainerVisibleWhen`):
//   - `Section/DetailPolicyTreatyIn`: seluruh medan uang, grid spreading, dan
//     jadwal angsuran proporsional berada di kontainer
//     `.IsNewPolicyNonProp != 1 && .IsNewPolicyListFormat != 1` (IsNewPolicyListFormat
//     tidak pernah ditulis rule mana pun di korpus NB);
//   - `Section/DetailDeptHeadTreatyIn_UW`: kontainer yang sama, `.IsNewPolicyNonProp != 1`;
//   - kedua layar: subsection `DetailPoliciesNonProportional` di kontainer
//     `.IsNewPolicyNonProp = 1 && .ClaimType != 'XOL Retro'`.
//
// Medan di kontainer tersembunyi tidak dapat diisi dan tidak pernah memicu
// refresh - jadi tidak wajib, dan rantai hitungnya tidak berjalan dari layar.

// medanKontainerProp - medan polis di kontainer proporsional kedua layar.
var medanKontainerProp = func() map[string]bool {
	m := map[string]bool{}
	for _, n := range []string{
		"GrossPremium", "GrossClaim", "PremiOgp", "RiCommOgp", "ResultOgp1", "OveriddingCommOgp", "ResultOgp2",
		"PremiOnp", "RiCommOnp", "ResultOnp1", "OveriddingCommOnp", "ResultOnp2", "Claim", "OutstandingClaim",
		"SalvageValue", "ExcessLoss", "NetPremium", "BalanceDueTo", "BalanceBeforeTax", "BalanceBeforePPH",
		"Deduction1", "Deduction2", "PPHValue", "PPNValue", "Installment",
	} {
		m[HalamanPolis+"."+n] = true
	}
	return m
}()

// PolisNonPropBaru = `.IsNewPolicyNonProp = 1` (penanda subsection NonProp).
func PolisNonPropBaru(h *Halaman) bool { return samaDenganSatu(h.Ambil(pt + "IsNewPolicyNonProp")) }

// TersembunyiNonProp - medan `jalur` berada di kontainer proporsional yang
// tersembunyi bagi polis NonProp baru.
func TersembunyiNonProp(h *Halaman, jalur string) bool {
	return PolisNonPropBaru(h) && medanKontainerProp[jalur]
}

// DaftarDariLayar - PageList yang boleh dikirim layar admin.
//
//	SpreadingRiskList  grid proporsional; bagi NonProp, section `SpreadingRiskList`
//	                   di dalam DetailPolicyTreatyInNonProportional - Add/Delete dan
//	                   TreatyType/%Share terbuka hanya bila TreatyIn.FacultativeShare
//	                   = 0 atau '' (selain itu terkunci)
//	ListInstallment    grid proporsional; bagi NonProp grid ber-`pyRowEditing`
//	                   masterDetail dengan Section `InstallmentList` ber-pyReadOnly -
//	                   hanya-baca
//	TreatyXOLList      TIDAK tampil di section mana pun - tidak pernah dari layar
func DaftarDariLayar(h *Halaman) []string {
	if !PolisNonPropBaru(h) {
		return []string{DaftarSpreading, DaftarAngsuran}
	}
	fac := h.Ambil(jMaster + "FacultativeShare")
	if d, err := AngkaTeks("FacultativeShare", fac); fac == "" || (err == nil && d.IsZero()) {
		return []string{DaftarSpreading}
	}
	return nil
}
