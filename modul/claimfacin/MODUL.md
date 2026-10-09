# Modul `claimfacin` — Claim Fac In

Klaim fakultatif inward - kasus `ASM-FW-GCNMFW-Work-PNC` (`Flow/Register_Flow.xml`: Input Register -> Input Estimasi
-> Choose Surveyor / Input Adjustment -> Resolved-Completed). Dimigrasi 10-10-2026 (prompt "IMPLEMENTASI CLAIM FAC IN,
TAHAP 1 DARI 2"); spec `docs/spec.md`, paritas tombol XML `docs/PARITAS.md`, OQ `docs/OQ.md`, pindai `docs/PINDAI.md`,
tabel `docs/STRUKTUR-TABEL-CLAIM-FACIN.md`. Tahap 2 (keputusan komite KMT-) = modul `komiteclaimfacin`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `claimfacin` |
| Folder korpus | `Claim Fac In` |
| GROUPMENU | `KLAIM` |
| Pemilik | `@PEMILIK-CLAIMFACIN` |
| Status | dimigrasi |
| Rentang migrasi | `560-599` |
| Slot menu | `978-979` |
| Prefix rute API | `/api/claim-fac-in` |
| Kontrak disediakan | — (kontrak untuk `komiteclaimfacin` = tahap 2) |
| Kontrak dipakai | — (nol kontrak `inti/backend/kontrak`; pola Claim Prop / Non Prop disalin, tidak diimpor) |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `PINDAI.md`, `PARITAS.md`, `OQ.md`, `spec.md` (+ RALAT), `issues/` (+ RALAT), `STRUKTUR-TABEL-CLAIM-FACIN.md`, `RELASI-TABEL-CLAIM-FACIN.md`, grilling |
| `backend/` | `modul.go` (`Pendaftaran()`), `models/` (port activity, katalog, tata layar), `repository/` (Oracle), `services/`, `handlers/`, `tiruan/` (uji), `migrations/` |
| `frontend/` | `menu.ts`, `rute.tsx`, renderer tata, halaman awal |

## Migrasi

Rentang `560-599` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `978-979` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026), di folder `backend/migrations/` modul ini sendiri — bentuk SQL-nya di
`APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6. Nomor selalu tiga digit. `-migrate` dijalankan work owner.

| Nomor | Isi |
| --- | --- |
| `560` | 19 kolom khas FACIN di `T_GENERAL_CLAIM` (ALTER ADD nullable; baris kolom di STRUKTUR Claim Life) |
| `561` | tabel baru `T_CLAIM_OBJECT` (`ClaimData.ObjectList`) |
| `562` | tabel baru `T_CLAIM_OBJECT_ITEM` (`.ObjectItemList`) |
| `563`–`566` | `OBJECT_ITEM_ID` (FK CASCADE ke item) + kolom khas FAC nullable di tabel Claim Prop `T_CLAIM_ESTIMATION`, `T_CLAIM_SPREADING` (+ `JENIS`), `T_CLAIM_BREAK_QS`, `T_CLAIM_ADJUSTMENT` (baris kolom di STRUKTUR Claim Prop) |
| `567` | `T_CLAIM_FAC_RETRO.ADJUSTMENT_ID` (FK CASCADE ke adjustment) |
| `978` | menu `claimfacin` `DIMIGRASI = '1'` |

Bergantung pada migrasi Claim Prop yang berjalan lebih dulu: tabel `T_CLAIM_*` (`521`–`531`), `SEQ_T_CLAIM` (`533`),
`T_VIEW_SUGGEST.CLAIM_ID` (`532`), `T_KATEGORI_DOC_KLAIM` (`535`/`536`, baris FAC). Nol `MODIFY` / `DROP` kolom yang
sudah ada (OQ-CFI-01).

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per butir.

### Kaskade ON DELETE CASCADE

Kaskade HANYA pada berkas migrasi modul ini yang berawalan di bawah; berkas lain modul ini tanpa
`ON DELETE CASCADE` (`TestKaskadeHanyaPadaRelasiTerdaftar`).

| Awalan berkas | Relasi |
| --- | --- |
| `561_` | T_CLAIM_OBJECT.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.ObjectList` |
| `562_` | T_CLAIM_OBJECT_ITEM.OBJECT_ID -> T_CLAIM_OBJECT dan .CLAIM_ID -> T_GENERAL_CLAIM - `.ObjectList(n).ObjectItemList` |
| `563_` | T_CLAIM_ESTIMATION.OBJECT_ITEM_ID -> T_CLAIM_OBJECT_ITEM - `.ObjectItemList(i).EstimationList` |
| `564_` | T_CLAIM_SPREADING.OBJECT_ITEM_ID -> T_CLAIM_OBJECT_ITEM - `.ObjectItemList(i).SpreadingList / SpreadingClaim` |
| `565_` | T_CLAIM_BREAK_QS.OBJECT_ITEM_ID -> T_CLAIM_OBJECT_ITEM - `.ObjectItemList(i).SpreadingAdjustment` |
| `566_` | T_CLAIM_ADJUSTMENT.OBJECT_ITEM_ID -> T_CLAIM_OBJECT_ITEM - `.ObjectItemList(i).Adjustment` |
| `567_` | T_CLAIM_FAC_RETRO.ADJUSTMENT_ID -> T_CLAIM_ADJUSTMENT - `.Adjustment(a).FacRetroList` |

## Keputusan work owner 09-10-2026 yang mengikat modul ini

- **OQ-CFI-01 tabel Prop + kolom item**: tabel baru `T_CLAIM_OBJECT` / `T_CLAIM_OBJECT_ITEM`; tabel `T_CLAIM_*` Claim Prop
  dipakai ulang dengan `OBJECT_ITEM_ID` tambahan; baris FAC mengisi `CLAIM_ID` dan `OBJECT_ITEM_ID`; nol MODIFY.
- **OQ-CFI-02 awalan ID sama dengan Pega**: klaim `CLM-`, komite `KMT-`, nomor `SEQ_WORK_CLAIM`; lini disaring lewat
  `T_WORK_CLAIM.LINI = 'FACIN'`, tidak pernah lewat awalan; `TAHAP` = nama FlowAction (`InputRegister`,
  `InputEstimasi`, `InputSurveyor`) supaya tidak bertabrakan dengan label TAHAP Claim Life.
- **OQ-CFI-03** kelainan XML diperbaiki (`[penyimpangan sadar]` di `docs/PARITAS.md`).
- **OQ-CFI-04 komite pola Claim Prop**: tanpa menu komite; kasus `KMT-` TT2 lahir di `T_WORK_CLAIM` + `T_GENERAL_KOMITE` +
  `T_KOMITE_KOMITELIST` (hanya `backend/models/komite.go` + `backend/repository/komite.go`, dikecualikan penjaga batas
  Claim Life `komite_statik_test.go` atas izin work owner 09-10-2026); `POSITION` = `OPERATOR_ID` tingkat 1. TT3 / TT4
  (tanpa adjustment) = OQ-CFI-27.
- **Izin menyunting STRUKTUR Claim Life / Claim Prop** untuk kolom yang ditambahkan ke tabel bersama (dokumen saja).
- `OS_AKSEPTASI_KLAIM.DATA_JSON` diisi (format `GetPageJSONString`), `JSON_KLAIM` tanpa `DATA_JSON`; nol stored
  procedure, nol `COMMIT` di teks SQL; uang / persen `NUMBER(38,10)`; efek luar hanya produksi; nilai DB tampil apa
  adanya.
- Kata sandi Save Spreading tidak tertulis di mana pun. Tombolnya tidak dibangun: container `LS42` section `Estimasi`
  bergerbang `ClaimData.ExGratia = 1`, sedangkan penulis satu-satunya (`InsertObjects_dt` 11) menulis `0` (OQ-CFI-17,
  `docs/PARITAS.md`).
