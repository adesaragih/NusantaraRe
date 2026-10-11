# Modul `claimnonprop` — Claim Non Prop

Klaim treaty inward non proporsional / XoL - kasus `ASM-FW-GCNMFW-Work-ClaimTreatyNonProp` (`Flow/Flow_TreatyIn.xml`:
Outstanding Claim -> Input Acceptation -> Resolved-Completed). Dimigrasi 09-10-2026 (prompt "IMPLEMENTASI CLAIM NON
PROP, TAHAP 1 DARI 2"); spec `docs/spec.md`, paritas tombol XML `docs/PARITAS.md`, OQ `docs/OQ.md`, tabel
`docs/STRUKTUR-TABEL-CLAIM-NON-PROP.md`. Tahap 2 (keputusan komite KMTNP-) = modul `komiteclaimnonprop`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `claimnonprop` |
| Folder korpus | `Claim Non Prop` |
| GROUPMENU | `KLAIM` |
| Pemilik | `@PEMILIK-CLAIMNONPROP` |
| Status | dimigrasi |
| Rentang migrasi | `600-639` |
| Slot menu | `982-983` |
| Prefix rute API | `/api/claim-non-prop` |
| Kontrak disediakan | `kontrak.KlaimTreatyNonPropKomite` (`inti/backend/kontrak/klaimtreatynonprop.go`, dipakai `komiteclaimnonprop`) — perintah work owner 09-10-2026 |
| Kontrak dipakai | — (nol kontrak `inti/backend/kontrak`; pola Claim Prop disalin, tidak diimpor) |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `PINDAI.md` (pindai korpus + DEV 09-10-2026), `OQ.md`, `spec.md`, `PARITAS.md`, `STRUKTUR-TABEL-CLAIM-NON-PROP.md` |
| `backend/` | `modul.go` (`Pendaftaran()`), `models/` (port activity, mesin XoL, katalog, tata layar), `repository/` (Oracle), `services/`, `handlers/`, `tiruan/` (uji), `konfigurasi/kasir.json` |
| `frontend/` | `menu.ts`, `rute.tsx`, renderer tata (`components/`), halaman awal (`pages/ClaimNonProp.tsx`) |

## Migrasi

Rentang `600-639` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `982-983` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026), di folder `backend/migrations/` modul ini sendiri — bentuk SQL-nya di
`APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6. Nomor selalu tiga digit. `-migrate` dijalankan work owner.

| Nomor | Isi |
| --- | --- |
| `600` | 25 kolom khas NONPROP di `T_GENERAL_CLAIM` (ALTER ADD nullable; baris kolom di STRUKTUR Claim Life) |
| `601`–`607` | kolom tambahan nullable di tabel Claim Prop `T_CLAIM_INTEREST`, `T_CLAIM_CLAIM_AMOUNT`, `T_CLAIM_SPREADING`, `T_CLAIM_BREAK_QS`, `T_CLAIM_ADJUSTMENT`, `T_CLAIM_ADJ_SPREADING`, `T_CLAIM_ADJ_QUOTA_SHARE` (baris kolom di STRUKTUR Claim Prop) |
| `608`–`610` | tabel baru `T_CLAIM_NP_LOSS_ALLOC`, `T_CLAIM_NP_XOL_ALLOC` (ber-`JENIS`), `T_CLAIM_NP_CLAIM_ACCEPT` |
| `611` | roster `EMAILKOMITE` NONPROP DEGREE 1–4 -> workbasket `ReasClaimDeptHead` / `ReasClaimTechDivHead` / `ReasClaimOpsDir` / `ReasClaimTechDir` (UPDATE di tempat, pola claimprop `537`) |
| `982` | menu `claimnonprop` `DIMIGRASI = '1'` |

Bergantung pada migrasi Claim Prop yang berjalan lebih dulu: tabel `T_CLAIM_*` (`521`–`531`), `SEQ_T_CLAIM` (`533`),
`T_VIEW_SUGGEST.CLAIM_ID` (`532`), workbasket roster (`537`).

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per butir.

### Kaskade ON DELETE CASCADE

Kaskade HANYA pada berkas migrasi modul ini yang berawalan di bawah; berkas lain modul ini tanpa
`ON DELETE CASCADE` (`TestKaskadeHanyaPadaRelasiTerdaftar`).

| Awalan berkas | Relasi |
| --- | --- |
| `608_` | T_CLAIM_NP_LOSS_ALLOC.CLAIM_ID -> T_GENERAL_CLAIM dan .ADJUSTMENT_ID -> T_CLAIM_ADJUSTMENT - `.ClaimData.CNPSpreadLoss` / `.AdjustmentList(n).CNPSpreadLoss` (tepat satu induk) |
| `609_` | T_CLAIM_NP_XOL_ALLOC.CLAIM_ID -> T_GENERAL_CLAIM dan .ADJUSTMENT_ID -> T_CLAIM_ADJUSTMENT - `.ClaimData.SpreadingRisk` / `.AdjustmentList(n).SpreadingRisk / LossAllocation / AlokasiXOLPaid` (tepat satu induk) |
| `610_` | T_CLAIM_NP_CLAIM_ACCEPT.ADJUSTMENT_ID -> T_CLAIM_ADJUSTMENT - `.AdjustmentList(n).ListClaimAcceptation` |

## Keputusan work owner 09-10-2026 yang mengikat modul ini

- **OQ-CNP-01 tangga komite ikut XML**: RNM Share <= 30 dan ValueAdjustment akseptasi terakhir <= 30.000.000, atau
  subjectivity -> hanya tingkat roster DEGREE 1; selain itu semua baris roster NONPROP aktif. Konstanta di
  `backend/models/komite.go`.
- **OQ-CNP-04** kata sandi Edit XOL Allocation ke konfigurasi: tidak ada jalur konfigurasi sah tanpa sandi di repo
  (modul dilarang membaca env) -> tombol tampil nonaktif + OQ.
- **OQ-CNP-05** ketujuh kelainan XML diperbaiki (`[penyimpangan sadar]` di `docs/PARITAS.md`).
- **OQ-CNP-06** pola Claim Prop + tabel XoL (`docs/STRUKTUR-TABEL-CLAIM-NON-PROP.md`).
- **Komite ikut pola Claim Prop**: tanpa menu komite; kasus `KMTNP-` lahir di `T_WORK_CLAIM` + `T_GENERAL_KOMITE` +
  `T_KOMITE_KOMITELIST` (hanya `backend/models/komite.go` + `backend/repository/komite.go`, dikecualikan penjaga batas
  Claim Life `komite_statik_test.go` atas izin work owner 09-10-2026); `POSITION` = workbasket tingkat 1.
- **Izin menyunting STRUKTUR Claim Life / Claim Prop** untuk kolom yang ditambahkan ke tabel bersama (dokumen saja).
- `OS_AKSEPTASI_KLAIM.DATA_JSON` diisi (format `GetPageJSONString`), `JSON_KLAIM` tanpa `DATA_JSON`; nol stored
  procedure, nol `COMMIT` di teks SQL; uang / persen `NUMBER(38,10)`; nilai DB tampil apa adanya.
- `konfigurasi/kasir.json` - kode tetap muatan Kasir (`HitServiceToKasir_Act` 9.3-9.4), bukan literal kode.
