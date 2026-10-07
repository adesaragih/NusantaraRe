package models

// Untuk apa berkas ini: ATURAN LAYAR jalur NonProporsional - kontainer yang
// tersembunyi bagi polis XOL baru, dan grid yang boleh dikirim layar.
//
// Bukti XML (sel layout, `pyContainerVisibleWhen`):
//   - `Section/DetailPolicyTreatyIn`: seluruh medan uang, grid spreading, dan
//     jadwal angsuran proporsional berada di kontainer S19
//     `.IsNewPolicyNonProp != 1 && .IsNewPolicyListFormat != 1` (IsNewPolicyListFormat
//     tidak pernah ditulis rule mana pun di korpus NB);
//   - `Section/DetailDeptHeadTreatyIn_UW`: kontainer S90 `.IsNewPolicyNonProp != 1`;
//   - kedua layar: subsection `DetailPoliciesNonProportional` di kontainer
//     `.IsNewPolicyNonProp = 1 && .ClaimType != 'XOL Retro'` (admin S17, atasan
//     S88; sel SUB_SECTION `pyEditOptions=Auto`, tidak `pyReadOnly`).
//
// Medan di kontainer tersembunyi tidak dapat diisi dan tidak pernah memicu
// refresh - jadi tidak wajib (`wadahUangAdmin` / `bukanNonPropBaru`, layar.go),
// tidak diterima dari layar (`medanAdmin`, W4 audit silang P3), dan rantai
// hitungnya tidak berjalan dari layar (`services.turunkan`, `validasiKirim`).

// PolisNonPropBaru = `.IsNewPolicyNonProp = 1` (penanda subsection NonProp).
func PolisNonPropBaru(h *Halaman) bool { return samaDenganSatu(h.Ambil(pt + "IsNewPolicyNonProp")) }

// tampilSubsectionNonProp = wadah `.IsNewPolicyNonProp = 1 && .ClaimType != 'XOL
// Retro'` (admin S17, atasan S88) yang menyertakan `DetailPoliciesNonProportional`.
// Varian EDM subsection itu (wadah S2 `TreatyMasterInEDM && IsEDMInputOnNB !=
// true`) tak terjangkau: InputPolicyTreatyInDetail_NonProp langkah 10 mengisi
// `IsEDMInputOnNB = true` tepat ketika `TreatyMasterInEDM` (frontend `nonprop.ts`).
func tampilSubsectionNonProp(h *Halaman) bool { return PolisNonPropBaru(h) && bukanXOLRetro(h) }

// fakultatifNolAtauKosong = syarat tombol Add/Delete `Section/SpreadingRiskList`:
// `pyWorkPage.TreatyIn.FacultativeShare = 0 || pyWorkPage.TreatyIn.FacultativeShare
// = ”`. Sel TreatyType/%Share/%Share Claim terkunci `FacultativeShare > 0`.
// ⚠️ Nilai negatif / bukan angka (sel terbuka, tombol tidak) diperlakukan
// tertutup - master tidak pernah memuatnya; syarat tombol yang dipakai.
func fakultatifNolAtauKosong(h *Halaman) bool {
	fac := h.Ambil(jMaster + "FacultativeShare")
	d, err := AngkaTeks("FacultativeShare", fac)
	return fac == "" || (err == nil && d.IsZero())
}

// SpreadingDariLayar - grid `.SpreadingRiskList` dapat disunting di layar posisi
// `posisi` (baris Add/Delete; sel `.TreatyType`, `.SharePercentage`,
// `.ClaimPercentage`). Hanya bila benar, daftar kiriman layar diterima
// (`GabungMasukanLayar`) dan refresh `CountSpreading_Act` sel %Share dapat
// terpicu (`services.aksiTerbuka`).
//
//	subsection NonProp  layar ADMIN saja: `tampilSubsectionNonProp` DAN
//	                   `fakultatifNolAtauKosong`. `[keputusan work owner 07-10-2026]`
//	                   "HANYA ADMIN YANG BISA EDIT" - menyimpang sadar dari XML/W2:
//	                   `DetailDeptHeadTreatyIn_UW` S88 menyertakan subsection yang sama
//	                   ber-`pyEditOptions=Auto`, tetapi Sec Head / Dept Head hanya-baca
//	admin Proporsional `DetailPolicyTreatyIn` S30 (wadah S19 `wadahUangAdmin`):
//	                   sel terbuka, Add/Delete
//	atasan Proporsional `DetailDeptHeadTreatyIn_UW` S96: sel `pyReadOnly=true`
//	                   tanpa syarat, tanpa Add/Delete - TIDAK
//
// `InputParam.CARI12 != 'treaty' && != 'TreatyPolicy'` (syarat tombol Add/Delete)
// tidak dibaca: parameter harness itu tidak diisi alur NB (layar dibuka lewat
// `openWorkByHandle` / flow action), jadi selalu benar.
func SpreadingDariLayar(h *Halaman, posisi string) bool {
	if tampilSubsectionNonProp(h) {
		return posisi == PosisiAdmin && fakultatifNolAtauKosong(h)
	}
	return posisi == PosisiAdmin && wadahUangAdmin(h)
}

// IsianNonProp - isian layar yang ikut ditulis ulang preACT 16 saat pajak NonProp dihitung ulang
// (`services.pajakNonProp`, `[keputusan work owner 06-10-2026]`): langkah 13 StartDate/EndDate dan
// langkah 23 %Share / %Share Claim spreading.
type IsianNonProp struct {
	mulai, akhir string
	spreading    []Baris
}

// SimpanIsianNonProp mencatat isian layar sebelum preACT 16 dijalankan ulang.
func SimpanIsianNonProp(h *Halaman) IsianNonProp {
	return IsianNonProp{h.Ambil(pt + "StartDate"), h.Ambil(pt + "EndDate"), salinBaris(h.AmbilDaftar(DaftarSpreading))}
}

// Kembalikan memasang lagi isian yang dicatat. %Share kembali ke baris ber-Treaty Type sama di urutan
// yang sama (baris spreading milik server); spreading lalu dihitung ulang (`CountSpreading_Act`).
func (s IsianNonProp) Kembalikan(h *Halaman) error {
	h.Setel(pt+"StartDate", s.mulai)
	h.Setel(pt+"EndDate", s.akhir)
	baru := h.AmbilDaftar(DaftarSpreading)
	for i, r := range baru {
		if i >= len(s.spreading) || s.spreading[i]["TreatyType"] != r["TreatyType"] {
			continue
		}
		for _, m := range selSpreading {
			r[m] = s.spreading[i][m]
		}
	}
	h.SetelDaftar(DaftarSpreading, baru)
	return CountSpreading(h, 0)
}
