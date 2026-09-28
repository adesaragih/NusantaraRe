# PROMPT — GILIRAN 8 *(sesi tunggal di `OUTPUT_HASIL_RNM`, cabang `main` @ `5bf3539`)*: **paket 0 rapikan Claim Life (bi) + ambang ba dibuka (bh) → PL 05a (OQ-PL-08 terjawab) → 05b → 06 + 08 + 09 → Komite 01 → 09 → Claim Life §3.1**

> GILIRAN-3 *(A/B/C)*, 4, 5, 6, 7 **tetap berlaku**, termasuk aturan berhenti GILIRAN-6 *(laporan pertama giliran ini paling
> cepat sesudah **05b**; berhenti dini hanya dengan "konteks menipis: kira-kira N% tersisa")*. Mekanisme giliran = lanjutan 8 §1.

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 28-09-2026

| Klaim laporan | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 4 commit `2cbdf2d` → `5bf3539`; berhenti "konteks menipis 15%" | 4 commit, 31 berkas, +3.571; pohon bersih; alasan berhenti sah | ✅ |
| Go 501 · 0 · 38; JS 339 | dijalankan ulang: **501 PASS · 0 FAIL · 38 SKIP** · vitest **339** · **60** modul · vet, gofmt, tsc, build bersih | ✅ |
| av + pl4 | `repository/polis_ringkas.go` *(`PolisRingkas`)*; `PanelDataPolis` tanpa penanda "menunggu" | ✅ |
| `utils.ParseDecimal` menolak `NaN`/`Infinity` | `pkg/utils/decimal.go:46–54` | ✅ lintas modul, benar |
| ralat tiket 04 *(32 kolom uang; `M_TEMPUPLOADLIFE` tujuh kolom)*; OQ-PL-05/06/07 | tiket 04 baris 196–226 | ✅ |
| OQ-PL-08 tiga parameter summary tidak tertentukan | **Terjawab dari korpus** — §1 | ✅ ditutup |
| ba tidak dapat ditutup *(view `PRODUCTINWARD_LIFE`, penjaga `TestMasterViewTidakDisentuh`)* | keputusan **bh** — §1 | ✅ dibuka |

## 1. TIGA KEPUTUSAN GILIRAN INI

### bh — ambang `ba` dibaca dari view `PRODUCTINWARD_LIFE` `[DIPUTUSKAN; veto work owner]`

`Claim Life/RDBList/GetProductName.xml`: `SELECT * FROM POOLDATA.PRODUCTINWARD_LIFE WHERE ID = {pyWorkPage.PolicyDataLife.ProductNameID}`
— Pega **sendiri** membaca view itu; AC 38 melarang **mengurai `JSONDATA`** di aplikasi, bukan membaca kolom bertipe dari view
milik basis data. Maka: satu pembaca `repository` **read-only** atas **dua kolom** `MAXEXPIREDCLAIM`, `MAXDATARECEIVE` berkunci
`ID` = `PRODUCT_NAME_ID` *(av)*; `TestMasterViewTidakDisentuh` **dipersempit** untuk mengizinkan tepat pembaca itu *(nama fungsi
dan dua kolom dikunci)*; nol `JSONDATA`, nol tabel `m_product_life` di kode. Ambang `ba` disambungkan ke Detail; blok bertanggal di
tiket 06 Claim Life dan OQ-001 *(pertanyaan DDL produksi tetap terbuka; ini tidak menunggunya)*.

### OQ-PL-08 — ditutup dari `DataTransform/AppendCurrencySummary_DT.xml` *(dipanggil `SavePremiumList_Act`)*

| Parameter | Rumus di DT | Baris |
| --- | --- | --- |
| `.PREMIUM` | akumulasi `Param.Premium` = Σ `.GROSS_PREMIUM` + Σ `.GROSS_PREMIUM_REFUND` + Σ `.GROSS_PREMIUM_RETRO` + Σ `.GROSS_PREMIUM_REFUND_RETRO` *(per cabang jenis baris)* | b2298, b2359, b2419, b2478 |
| `.COMMISSION` | Σ `.COMM` saja; `PROF_COMM` dan `OVR_COMM` masuk parameter lain *(`Param.Prof_Comm` b2626, `Param.Ovr_Comm` b2568)* | b2509 |
| `.BALANCE` | Σ per cabang: **normal** `GROSS_PREMIUM − DEDUCTION − (RI_ADMIN_FEE + BROKERAGE_FEE + TAX + PROF_COMM + CLAIM)` b3526/b3541/b3645; **refund** `GROSS_PREMIUM_REFUND + CLAIM_AMOUNT − (DEDUCTION_REFUND + BROKERAGE_FEE_REFUND + RI_ADMIN_FEE_REFUND + TAX + PROF_COMM + CLAIM…)` b3695/b3709/b3712; **retro** `GROSS_PREMIUM_RETRO − DISCOUNT_PREMIUM_RETRO − RI_ADMIN_FEE_RETRO + BROKERAGE_FEE_RETRO` b3761/b3775/b3778; **refund retro** `GROSS_PREMIUM_REFUND_RETRO − DISCOUNT_PREMIUM_REFUND_RETRO − RI_ADMIN_FEE_REFUND_RETRO + BROKERAGE_FEE…` b3827/b3841/b3844 | ⚠️ tanda `+ BROKERAGE_FEE_RETRO` pada cabang retro berbeda dari cabang normal — **VERBATIM**, dicatat sebagai keanehan warisan, bukan "diperbaiki" |

Executor membaca DT itu **utuh** *(syarat tiap cabang, parameter lain b2538–b2829: brokerage, tax, claim, refund-refund)* sebelum
05a; rumus ditiru per baris lalu dijumlah per `CURRENCY`; contoh terhitung dari literal; tiket 05a diralat bertanggal; OQ-PL-08 ditutup.

### bi — rapikan Claim Life di frontend `[permintaan work owner 28-09-2026]`

Pindahkan `pages/InboxClaimLife.tsx`, `OutstandingClaimLife.tsx`, `RegisterKlaim.tsx`, `KlaimLife.tsx` *(+ uji)* ke
`pages/claimlife/`; pecah `assets/labels.ts`: konstanta milik Claim Life *(`TAHAP`, `TAHAP_ID`, `TOMBOL`, `LAYAR`, `PERAN`,
`PERAN_ID`, `PRODUK`, `REGISTER`, `OUTSTANDING`, `TOMBOL_OS`, `DETAIL`, `DOKUMEN`, `TOMBOL_MEDIS`, `TOMBOL_AKSEPTASI`,
`TOMBOL_KOMITE`, `DIAGNOSA`)* → `assets/labels.claimlife.ts`; yang bersama *(`MENU`, `MODUL`, `MENU_MODUL`,
`KETERANGAN_BELUM_DIMIGRASI`, `BERANDA`, `MODUL_LAIN_TERLARANG`)* tetap di `labels.ts`. **Nol perubahan perilaku**: uji yang ada
tetap hijau tanpa disunting kecuali jalur impor; komponen `PanelDataPolis`, `PanelDokumenPeserta`, `PanelTotalPeserta`,
`PanelPindahTahap`, `GridDiagnosa`, `CariDiagnosa` ikut ke `components/claimlife/`. Commit
`refactor: halaman dan label Claim Life ke folder modulnya (bi)`.

## 2. URUTAN — semua di `main`, satu commit per paket

| # | Paket | Rujukan |
| ---: | --- | --- |
| 0 | **bi** rapikan Claim Life · **bh** ambang ba disambungkan | §1 |
| 1 | PL **05a** summary uang *(32 kolom; tiga parameter dari §1; pl2 `M_LIFE_PREMIUM_SUMMARY`/`M_LIFE_PREMIUM_DETAIL` satu transaksi; `SubmitPremiumList_Act` langkah 15 diputuskan)* | brief 3-PREMIUMLIST §2; tiket 05a |
| 2 | PL **05b** JSON polis | tiket 05b |
| — | **laporan pertama boleh di sini** | |
| 3 | PL **06 + 08 + 09** | brief 3-PREMIUMLIST §2 |
| 4 | Komite **01 → 09** | brief 3-KOMITE |
| 5 | Claim Life **§3.1** | brief 3-CLAIM-LIFE §3.1 |

## 3. ATURAN · LAPORAN · TELEMETRI

Sama dengan GILIRAN-4 §2–§3 dan GILIRAN-6 aturan berhenti. Laporan: baris pertama alasan berhenti; tabel **paket → commit →
tombol/rule XML → rute/komponen**; ralat tiket; OQ dibuka/ditutup; angka uji tiap commit; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 28 September 2026 sesudah verifikasi `5bf3539` (uji, vet, gofmt, tsc, build dijalankan ulang), pembacaan
`AppendCurrencySummary_DT.xml` (akumulasi b2298–b2829, Balance b3526–b3844), `GetProductName.xml`, `decimal.go`, dan inventaris
`frontend/src/pages` serta `assets/labels.ts`.*
