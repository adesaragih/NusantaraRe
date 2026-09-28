# PROMPT — GILIRAN 7 *(sesi tunggal di `OUTPUT_HASIL_RNM`, cabang `main` @ `8f69682`)*: **paket 0 catatan ralat bulan → PremiumList 04 + av + 05a + 05b → 06 + 08 + 09 → Komite 01 → 09 → Claim Life §3.1** — laporan pertama paling cepat sesudah **05b**

> GILIRAN-3 *(A/B/C)*, 4, 5, 6 **tetap berlaku**, termasuk **aturan berhenti** GILIRAN-6 *(laporan hanya pada batas kelompok;
> berhenti dini hanya dengan "konteks menipis: kira-kira N% tersisa")*. Mekanisme giliran = lanjutan 8 §1.

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 28-09-2026

| Klaim laporan | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 4 commit `04609c7` → `8f69682`; 51 berkas | 4 commit; 49 berkas, +6.386 baris; pohon bersih | ✅ |
| Go 463 · 0 · 38; JS 322; gofmt, vet, vet `-tags db`, tsc, build | dijalankan ulang: **463 PASS · 0 FAIL · 38 SKIP** · vitest **322** · **59** modul · semuanya bersih | ✅ |
| bg: Beranda, 17 kelompok, PaletMenu, pola grid | `pages/Beranda.tsx`, `components/PaletMenu.tsx`, `pages/premiumlist/`, Shell "belum dimigrasi" | ✅ |
| bentuk nomor PL b3126; satu penghitung b2717; tahun dua angka b3021/b3075; periode tidak dikirim (memo) | `Activity/SubmitPremiumList_Act.xml` *(bukan `GetPLNumber_Act`)*: b13 memo `buang ParamSeq.CARI3`; b2717 `ParamSeq.CARI2 = ParamSeq.HASIL3+"QR/QP/TP/TR"`; b3021 `@substring(ParamSeq.HASIL1,0,2)`; b3075 `@substring(ParamSeq.HASIL1,5,7)`; b3126 `ParamSeq.HASIL3+InputData.CARI20+pyWorkPage.BusinessCode+"."+ParamSeq.HASIL1+"."+ParamSeq.HASIL2` | ✅ *(nama berkas dikoreksi di tiket 03 bila tertulis `GetPLNumber_Act`)* |
| tiga AC tiket 03 diralat | tiket 03 bab `[terbuka]` baris 231–245, `/code-review`, Verifikasi | ✅ |
| `HitungPeriodeNomor` `AddDate(0,1,0)` → jepit akhir bulan seperti `ADD_MONTHS` | `penomor.go:94–118`; 7 kasus; **mengubah keluaran penomoran Claim Life** | ✅ kode; ⛔ **belum dicatat** di tiket 02 Claim Life maupun LAPORAN Claim Life → paket 0 |
| OQ-PL-01 kunci `CLASS` penghitung | **Terjawab dari data DEV** *(agregat, 28-09-2026)*: tabel `POOLDATA.GENERATE_SEQUENCE_NUMBER` memuat `CLASS = ASM-FW-GISFW-Work-LIFE`, `JENIS = RNML-QR/QP/TP/TR` *(dua baris: tahun 2025 dan 2026; `NO_SEQ` tertinggi 39)*; kelas `RNM-FW-LIFEFW-Work-LIFE` **tidak ada** di tabel. Pembanding: Claim Life `RNML-K` *(seq 31)*, Komite `RNML-A` *(seq 3)*, NB `RNM-QR/QP/TP` | ✅ **ditutup**: kelas konkret = `ASM-FW-GISFW-Work-LIFE`; `JENIS` = `HASIL3` *(awalan dari `KODE_PRODUKSI`)* + `"QR/QP/TP/TR"` = `RNML-QR/QP/TP/TR` VERBATIM |
| OQ-PL-02/03/04 | tetap terbuka | — |

## 1. PAKET 0 — dua catatan yang tertinggal *(dokumen saja, satu commit)*

1. Tiket 02 Claim Life *(penomoran)* + `LAPORAN-GILIRAN-F0.md`: bab bertanggal *"Ralat 28-09-2026 — pergeseran bulan
   `HitungPeriodeNomor` menjepit ke akhir bulan seperti `ADD_MONTHS`"* dengan 7 kasus uji dan dampaknya pada nomor klaim
   tanggal 29–31.
2. Tiket 03 PremiumList: OQ-PL-01 **ditutup** dengan blok bertanggal *(bukti §0)*; konstanta kelas penghitung di kode =
   `ASM-FW-GISFW-Work-LIFE`, dikunci uji; sumber "b2717" disebut sebagai `SubmitPremiumList_Act.xml`.

Commit `docs: ralat bulan penomoran di tiket 02 Claim Life; OQ-PL-01 ditutup dari data DEV`.

## 2. URUTAN — semua di `main`, satu commit per tiket, laporan hanya pada batas kelompok

| # | Kelompok | Isi | Rujukan |
| ---: | --- | --- | --- |
| 1 | **PL 04 + av + 05a + 05b** | 04 unggah CSV *(`UploadCSV_LifePremium` → `UploadCSVLifePremium_Act`, `ValidasiUploadPL_act`, `InsertLifePremiumDetail_act`, `GenerateDetailPeserta`, harness `ViewCSVResult_*`; fixture `UJI-*`)* + **pl4** `repository.PolisRingkas`; **segera** av Claim Life *(`PanelDataPolis` sepuluh medan dari `T_PREMIUM_LIST`, ambang ba, penanda dicabut)*; 05a summary uang *(`ShowLifePremiumSummary`, `SavePremiumList_Act`, `SubmitPremiumList_Act` langkah 15 yang tiket 03 tandai `[terbuka]` — dibaca ulang dan diputuskan di sini, pl2 satu transaksi)*; 05b JSON polis | brief 3-PREMIUMLIST §2; 3-CLAIM-LIFE §4 |
| — | **laporan pertama boleh di sini** | | |
| 2 | **PL 06 + 08 + 09** | outbox Arasapas/email dengan stub; kontrak hilir dua sisi; verifikasi skema 050–056 lawan STRUKTUR | brief 3-PREMIUMLIST §2 |
| 3 | **Komite 01 → 09** | Inbox Komite pola grid; `ShowTransfer` layar keputusan; `KomitePostAdjustment` utuh; km1–km5 | brief 3-KOMITE §1–§2 |
| 4 | **Claim Life §3.1** | 41 baris sensus diputuskan; celah dibangun | brief 3-CLAIM-LIFE §3.1 |

## 3. ATURAN · LAPORAN · TELEMETRI

Sama dengan GILIRAN-4 §2–§3 dan GILIRAN-6 aturan berhenti. Laporan: baris pertama menyebut selesai atau konteks menipis;
tabel **tiket → commit → tombol/kelas XML → rute/komponen**; ralat tiket bertanggal; OQ dibuka/ditutup; angka uji tiap commit;
bab **TELEMETRI EKSEKUSI** per tiket.

---

*Disusun 28 September 2026 sesudah verifikasi `8f69682` (uji, vet, gofmt, tsc, build dijalankan ulang), pembacaan ulang
`SubmitPremiumList_Act.xml` b13/b2717/b3021/b3075/b3126, `penomor.go`, tiket 03 PremiumList, dan katalog DEV
`GENERATE_SEQUENCE_NUMBER` (agregat per kelas dan jenis).*
