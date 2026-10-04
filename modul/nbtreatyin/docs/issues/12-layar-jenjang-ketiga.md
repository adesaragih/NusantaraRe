# 12: Layar jenjang ketiga — sembilan medan wajib SEKALIGUS terkunci

**Status:** selesai *(putaran 2 2026-10-03, paket P2 cabang `modul/nbtreatyin/p2-layar` — AC 52 ✅; semula: sebagian, implementasi 2026-10-03 cabang `modul/nbtreatyin/implementasi`; awalnya ready-for-agent)*
**Blocked by:** 11

> ⭐⭐ **LEPAS DARI BLOKIR.** `[penyimpangan sadar]` 2026-09-22 — dua baris ini semula berbunyi:
> > *"**Status:** blocked · **Blocked by:** 11 · ⛔ **P18** *(isi 268 langkah penetapan nilai)*"*
>
> **P18 ditarik oleh tim migrasi.** Isi langkah ada di dalam ekspor, di tag
> `PropertiesName`/`PropertiesValue` **tanpa awalan `py`**. Lihat `VERIFIKASI-P18.md`.
**Menutup:** AC 47 · 48 · 52 · 78 · 80 *(5 AC)* — US 11 · 30

## Hasil & nilai pengguna

Hari ini pemegang jenjang ketiga **mengisi tujuh medan sendiri**, lalu memutuskan — ⭐ ia **bukan**
sekadar melihat. ⛔ Tetapi layarnya memuat **sembilan medan yang wajib diisi sekaligus terkunci**:
medan itu wajib, **tidak dapat diisi pengguna**, dan yang seharusnya mengisinya adalah
**perhitungan yang isinya belum terkirim**.

> ⛔ **RALAT** 2026-10-03 (putaran 2, paket P2) — bunyi lama: *"Hari ini pemegang jenjang ketiga
> **mengisi tujuh medan sendiri**"* → bunyi baru: pemegang jenjang ketiga mengisi **Approval** dan
> **Suggest** (selalu), ditambah **Production Date** hanya bila medan itu tampil baginya; selebihnya
> layarnya hanya-baca. Bukti: bab Hasil implementasi putaran 2 di bawah.

Sesudah tiket ini, pemegang jenjang ketiga dapat mengisi medan keputusannya dan menyimpan berkas.

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat berikut dikutip utuh lalu ditarik:
> > *"⛔ **Selama P18 kosong, layar itu tidak dapat disimpan** — dan itu bukan cacat pembangunan,
> > melainkan akibat bahan yang belum lengkap."*
>
> ⭐ `[terverifikasi]` **Layar itu dapat disimpan.** Dari **34** medan wajib-sekaligus-terkunci di
> seluruh `Section`, **28** punya langkah pengisi yang isinya terbaca. Tiga sisanya — `.ExcessLoss`,
> `.OutstandingClaim`, `.SalvageValue` — **diketik underwriter** di `DetailPolicyTreatyIn` /
> `GeneralPolicyTreatyIn` *(`bacasaja=false`)*, dan hanya **ditampilkan terkunci** di sini.

## ⭐ Blocker — **NIHIL**

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — bab Blocker semula berbunyi:
> > *"⛔⛔ **Sembilan medan wajib-sekaligus-terkunci** … | ⛔ penyimpanan layar jenjang ketiga |
> > kesembilan medan wajib **tidak dapat diisi** | ⛔ penyelesaian tangga | berkas dapat naik ke
> > jenjang ketiga, ⛔ **tidak dapat keluar** dari sana | … ⭐ **Dua jalan keluar, keduanya bukan
> > keputusan teknis:** ekspor ulang **dengan isi langkah disertakan** *(P18)*, atau **pencabutan
> > kewajiban** atas kesembilan medan oleh Product & Underwriting."*

⭐ **Tidak ada jalan keluar yang perlu ditempuh.** Kesembilan medan itu memang wajib sekaligus
terkunci, dan itu tetap benar — tetapi **yang mengisinya terbaca di ekspor**. `[terverifikasi]`

| Medan | Diisi oleh |
| --- | --- |
| `Deduction2` | `InputPolicyTreatyEDMDetail_NP` · `InputPolicyTreatyInDetail_NonProp` |
| `OveriddingCommOgp` | `CountOGPONP_Act` · `CountOverridingCommOgp_Act` |
| `OveriddingCommOnp` | `CountOGPONP_Act` · `CountOverridingCommOnp_Act` |
| `ResultOnp1` | `CountOGPONP_Act` · `CountResult1Onp_Act` |
| `RiCommOgp` | `CalculatePremi_Act` · `CountOGPONP_Act` |
| `RiCommOnp` | `CountOGPONP_Act` · `CountResult1Onp_Act` |
| ⭐ `ExcessLoss` | **diketik underwriter** di layar `DetailPolicyTreatyIn` / `GeneralPolicyTreatyIn` |
| ⭐ `OutstandingClaim` | idem |
| ⭐ `SalvageValue` | idem |

## Area codebase

- Antarmuka: layar jenjang ketiga
- Lapisan handler: validasi medan wajib per tingkat

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Medan wajib per tingkat | layar jenjang ketiga **tidak** mewajibkan enam medan yang wajib di layar admin, ⭐ tetapi **mewajibkan satu** yang tidak wajib di sana |
| Medan terkunci | **36** dari 38 ada di dua layar jenjang ketiga |
| Medan yang dapat diisi | **tujuh** |
| Pembersihan pesan yang **ditahan** | **5** tempat yang menghapus **sesudah** validasi memasang pesan; seluruhnya di rantai yang sama |

> ⛔ **RALAT** 2026-10-03 (paket P2) — baris *"Medan yang dapat diisi | **tujuh**"* → **dua selalu
> (`.IsApproved`, `.Suggest`) + satu bersyarat (`.ProductionDate`)** — rule: `Section/ListSuggest`;
> `DueTo`, `FlagPPH`, `NoOfferSlip` di `Section/DetailDeptHeadTreatyIn_UW` ber-`pyDisabled=always`.

## ADR terkait

- **ADR-0007** — jejak audit setiap transisi

## Acceptance criteria

- [x] **AC 47** — layar jenjang ketiga **mewajibkan** satu medan yang tidak wajib di layar admin
- [x] **AC 48** — berkas **tidak dapat disimpan** bila medan wajib pada tingkat itu kosong
- [x] ✅ **AC 52** — pemegang jenjang ketiga **dapat mengisi Approval dan Suggest** (+ Production
      Date bersyarat); ⛔ layarnya **bukan** sepenuhnya hanya-baca, ⛔ medan lain tidak dapat diisi
      > ⛔ **RALAT** 2026-10-03 — bunyi lama: *"- [ ] 🟡 **AC 52** — pemegang jenjang ketiga **dapat
      > mengisi tujuh medan**; ⛔ layarnya **bukan** sepenuhnya hanya-baca"*. Uji:
      > `handlers/alur_test.go` `TestMedanTerkunciAtasanDanTurunanAdmin`, `models/layar_test.go`
      > `TestAtasanHanyaMengisiPutusanCatatanDanTanggalProduksiBersyarat`, `frontend/medan.test.ts`.
- [x] **AC 78** — **5** tempat yang menghapus pesan **sesudah** validasi memasangnya **ditahan**
      dan tidak dibangun
- [x] **AC 80** — ⭐ layar **dapat disimpan**; enam medan diisi langkah, tiga diketik underwriter
      di layar sebelumnya *(AC 80 dicabut lalu diganti 2026-09-22 — lihat `spec.md`)*

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| ⭐ ~~**P18**~~ | ~~isi 268 langkah penetapan nilai~~ | ⭐ **DITARIK 2026-09-22** — tidak menahan |
| **13** | **5** tempat penghapus pesan galat sesudah validasi | ⛔ menahan kelimanya |

## Perintah verifikasi

⭐ **Dapat dijalankan sepenuhnya.**

1. Buka layar jenjang ketiga — ⭐ **tujuh** medan dapat diisi, **36** terkunci.
   > ⛔ **RALAT** 2026-10-03 — bunyi lama: *"⭐ **tujuh** medan dapat diisi"* → bunyi baru: ⭐
   > **Approval** dan **Suggest** dapat diisi; **Production Date** hanya bila tampil (tiket 05);
   > `DueTo`, `FlagPPH`, `No Offer Slip` tampil terkunci.
2. Coba simpan dengan kesembilan medan kosong — ⭐ **ditolak**, dan pesannya menyebut bahwa
   nilainya berasal dari perhitungan yang belum tersedia.

## Catatan

⭐ **Tiket ini semula ditulis lengkap walau berstatus `blocked`** — supaya begitu P18 dijawab,
pekerjaannya sudah terumus. ⭐⭐ **P18 tidak pernah perlu dijawab**, dan pekerjaannya kini dapat
diambil. Kalimat lama ⛔ *"Jangan ditandai `ready-for-agent` sebelum itu"* — **ditarik**.

## Hasil implementasi 2026-10-03

Layar atasan: seluruh medan terkunci kecuali DueTo, FlagPPH, No Offer Slip, Approval, Suggest,
ProductionDate (enam; spec menyebut tujuh — AC 52 sebagian). Tombol: Sec Head selalu Submit; Dept Head
IsApproved 1 → Generate nomor polis → ShowPolicyNoTreaty (OK = Submit) — posisi menggantikan
`<ID-operator-1>` (tiket 05).

> ⛔ **RALAT** 2026-10-03 — kalimat *"seluruh medan terkunci kecuali DueTo, FlagPPH, No Offer Slip,
> Approval, Suggest, ProductionDate (enam; spec menyebut tujuh — AC 52 sebagian)"* ditarik: DueTo,
> FlagPPH, No Offer Slip ber-`pyDisabled=always` (lihat di bawah).

## Hasil implementasi putaran 2 — 2026-10-03 (paket P2, cabang `modul/nbtreatyin/p2-layar`)

- **AC 52 ✅.** `models.medanAtasan` = `IsApproved`, `Suggest`; `ProductionDate` diterima hanya bila
  tampil (`IsApproved == 1` + tempat `LISTSUGGEST_PRODUCTIONDATE`, juga di layar admin) —
  `GabungMasukanLayar(h, masuk, posisi, tempatTanggalProduksi)`. Layar: `FlagPPH` kotak terkunci,
  `DueTo` / `No Offer Slip` tampil terkunci.
- Wadah layar atasan: bagian uang/spreading/angsuran `.IsNewPolicyNonProp != 1` (juga syarat
  wajibnya — medan di wadah tersembunyi tidak wajib); `No Polis` / `Production Date` di wadah
  `pyWorkPage.FlagViewPolicy = 1` (tidak diisi rule terjangkau mana pun ⇒ tidak tampil);
  `Ceding Company`, `Marketing Officer`, `Type Tax` `NOTBLANK`; total berlabel `.Total*` (2 desimal).
- Format penyajian per sel (K14) — tiket 07.

> Bukti (dibaca ulang 2026-10-03): `FlowAction/DeptHeadTreatyIn_UW.xml` `pySectionReference =
> GeneralDeptHeadTreatyIn_UW` (nol sel sendiri; menyertakan `DetailDeptHeadTreatyIn_UW` atas
> `.PolicyTreatyIn`). `Section/DetailDeptHeadTreatyIn_UW.xml`: sel berkontrol ber-`pyReadOnly=false`
> hanya `.DueTo`, `.FlagPPH`, `.QuotationData.NoOfferSlip` — ketiganya `pyModes` baris 1
> `pyDisabled=true`/`pyDisabledNew=always`, **tidak dapat diisi** — dan lima `pxButton`.
> `Section/ListSuggest.xml`: `.IsApproved` (wajib), `.Suggest` (wajib), `.ProductionDate` (tampil dan
> wajib hanya `.IsApproved == 1 && OperatorID.pyUserIdentifier == '<ID-operator-3>' ||
> '<ID-operator-4>'` → tempat berperan tiket 05). Hitungan putaran 1 ("enam") memasukkan tiga sel
> `pyDisabled`; dasar "tujuh" tidak ditemukan. Subsection NonProp (wadah `.IsNewPolicyNonProp = 1`)
> milik paket P5.

## ⛔ RALAT putaran 2 (P9, 04-10-2026) — tempat `ProductionDate`

Bunyi lama (Hasil implementasi P2, AC 52), dikutip: *"`ProductionDate` diterima hanya bila tampil (`IsApproved == 1`
+ tempat `LISTSUGGEST_PRODUCTIONDATE`, juga di layar admin) — `GabungMasukanLayar(h, masuk, posisi,
tempatTanggalProduksi)`"*. Bunyi baru: syarat tampil `ListSuggest.ProductionDate` (`pyVisible`) dan syarat wajibnya
(`pyRequiredWhen`) masing-masing DUA tempat berperan tiket 05 (`LISTSUGGEST_PRODUCTIONDATE_TAMPIL_OPERATOR_3/4`,
`…_WAJIB_OPERATOR_3/4`); diterima bila `models.TanggalProduksiTampil(h, tempat)` —
`GabungMasukanLayar(h, masuk, posisi, tempat)`; wajib lewat satu sumber `models.MedanWajibBerlaku` (layar, Save,
submit). AC 52 tetap ✅.

## ⛔ RALAT putaran 3 (paket R6, 04-10-2026) — audit silang P3 W2 dan 7.4

### W2 — atasan menyunting grid spreading NonProp

Bunyi lama (Hasil implementasi putaran 2, AC 52), dikutip: *"`models.medanAtasan` = `IsApproved`,
`Suggest`; `ProductionDate` diterima hanya bila tampil"* dan baris tabel RALAT P2 *"**dua selalu
(`.IsApproved`, `.Suggest`) + satu bersyarat (`.ProductionDate`)**"*; serta komentar kode lama
`models.GabungMasukanLayar` *"Atasan  hanya `medanAtasan`; seluruh daftar terkunci"* dan
`services/aksiposisi.go` *"Yang terbuka hanya radio `ListSuggest .IsApproved`"*.

Bunyi baru: selain ketiga medan itu, pemegang layar `DetailDeptHeadTreatyIn_UW` dapat menyunting grid
`SpreadingRiskList` subsection NonProp — **syarat XML sama dengan layar admin**: polis NonProp baru
(wadah S88 `.IsNewPolicyNonProp = 1 && .ClaimType != 'XOL Retro'`), tombol Add/Delete bila
`pyWorkPage.TreatyIn.FacultativeShare = 0 || = ''`, sel `.TreatyType` / `.SharePercentage` /
`.ClaimPercentage` terkunci bila `FacultativeShare > 0`, sel %Share memicu refresh `CountSpreading_Act`.
Grid Proporsional layar atasan (S96 `pyReadOnly=true` tanpa syarat) **tetap hanya-baca**.

Bukti XML: `Section/DetailDeptHeadTreatyIn_UW.xml` S88 SUB_SECTION `DetailPoliciesNonProportional`
(`pyReadOnly=false`, `pyEditOptions=Auto`) → `DetailPoliciesNonProportional` S1 →
`DetailPolicyTreatyInNonProportional` S73 SUB_SECTION `SpreadingRiskList` (`pyReadOnly=false`,
`Auto`) → `Section/SpreadingRiskList.xml` S2 (Add/Delete `pyVisible` FacultativeShare 0/'';
`pyReadOnlyCondition pyWorkPage.TreatyIn.FacultativeShare >0`; CountSpreading_Act).
Spec AC 52 di-RALAT sama (bunyi lama dikutip di `spec.md`); AC 52 tetap ✅.

Dibangun: `models.SpreadingDariLayar(h, posisi)` (menggantikan `DaftarDariLayar`; kedua posisi),
`models.GabungMasukanLayar` menerima daftar itu di cabang atasan, `services/aksiposisi.go` aksi
`CountSpreading` terbuka di kedua layar hanya bila grid terbuka, `frontend` `DetailNonProp`
`sunting={boleh}`. Kolom `.PremiumSpreaded` / `.ClaimSpreaded` tetap dihitung server (W5, tiket 11).
Uji: `handlers/masukanlayar_test.go` `TestAtasanMenyuntingSpreadingNonProp`,
`TestAtasanSpreadingTerkunciBilaFakultatifAtauProporsional`; `models/layar_daftar_test.go`
`TestGabungMasukanDaftarMenurutGridXML`; `frontend/components/DetailNonProp.test.tsx`.

### 7.4 — `FetchTreatyGroupOldID` langkah 3 menahan Submit Dept Head

Sebelumnya tidak dibangun (status.json hanya menyebut langkah 6). Terjangkau: tombol Submit Dept Head
(IsApproved 1) → `runActivity GeneratePolicyNoTreaty_Act` → langkah 11 (`TreatyGroupOldID==""`)
`Call FetchTreatyGroupOldID` → langkah 2 `Param.Errmsg = "Cannot fetch Treaty Group ID, Contact IT"`,
langkah 3 `Page-Set-Messages pyWorkPage` bila `pyWorkPage.PolicyTreatyIn.TreatyGroupID==""`. Nomor
polis tetap dibentuk dan disimpan (langkah 30 `Obj-Save WithErrors=true`); halaman berpesan tidak dapat
di-submit (OK `ShowPolicyNoTreaty` = finishAssignment).

Dibangun: `models.GrupTreatyTakTerbaca`, `models.PesanGrupTreatyTakTerbaca` (VERBATIM);
`services.validasiKirim` menahan submit Dept Head (tombol nomor polis) dengan pesan itu (422), berkas
tetap di Dept Head. ⚠️ Di Oracle, TreatyGroupID kosong biasanya juga berarti `OJKBusinessID` kosong →
`TerbitkanNomor` sudah berhenti lebih dulu (`ErrOJKKosong`); pesan ini menahan kasus yang OJK-nya
terisi (mis. dokumen lama). Uji: `handlers/masukanlayar_test.go`
`TestGrupTreatyKosongMenahanSubmitDeptHead` (tiruan `OldIDGrup` dapat diatur).
