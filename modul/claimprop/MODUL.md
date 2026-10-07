# Modul `claimprop` — Claim Prop

Klaim treaty inward proporsional - kasus `ASM-FW-GCNMFW-Work-ClaimTreaty` (`Flow/Flow_TreatyIn.xml`: Outstanding
Claim -> Input Acceptation -> Resolved-Completed). Dimigrasi 07-10-2026 (prompt
`_brief/PROMPT-IMPLEMENTASI-MODUL-CLAIM-PROP.md`); paritas tombol XML di `docs/PARITAS.md`, OQ di `docs/OQ.md`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `claimprop` |
| Folder korpus | `Claim Prop` |
| GROUPMENU | `KLAIM` |
| Pemilik | `@PEMILIK-CLAIMPROP` |
| Status | dimigrasi |
| Rentang migrasi | `520-559` |
| Slot menu | `980-981` |
| Prefix rute API | `/api/claim-prop` |
| Kontrak disediakan | — |
| Kontrak dipakai | — (nol kontrak `inti/backend/kontrak`; pola modul lain disalin, tidak diimpor) |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | spec, tiket (`issues/`), grilling, catatan — dulu `.scratch/claim-prop/` (dipindah dengan `git mv`, isi byte-identik) |
| `backend/` | `modul.go` (`Pendaftaran()`), `models/` (port activity, katalog, tata layar), `repository/` (Oracle), `services/`, `handlers/`, `tiruan/` (uji), `alat/pemuatlama/` (pemuat data lama, uji-kering) |
| `frontend/` | `menu.ts`, `rute.tsx`, renderer tata (`components/TataView.tsx`), halaman awal (`pages/ClaimProp.tsx`) |

## Migrasi

Rentang `520-559` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `980-981` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026) saat modul mendapat layar pertamanya, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per butir (pola
MODUL.md Claim Life). Sebelum 07-10-2026 kedua judul di bawah berada di bab Migrasi sehingga TIDAK terbaca penjaga
(mutasi 532 tanpa CASCADE tetap hijau); dipindah ke bab ini.

### Kaskade ON DELETE CASCADE

Kaskade HANYA pada berkas migrasi modul ini yang berawalan di bawah; berkas lain modul ini tanpa
`ON DELETE CASCADE` (`TestKaskadeHanyaPadaRelasiTerdaftar`).

| Awalan berkas | Relasi |
| --- | --- |
| `521_` | T_CLAIM_ESTIMATION.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.EstimationList` hidup di halaman klaim |
| `522_` | T_CLAIM_INTEREST.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.InterestList` |
| `523_` | T_CLAIM_CLAIM_AMOUNT.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.ListClaimAmount` |
| `524_` | T_CLAIM_LOSS_ALLOCATION.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.SpreadingRisk` |
| `525_` | T_CLAIM_SPREADING.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.SpreadingClaim` |
| `526_` | T_CLAIM_BREAK_QS.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.SpreadingBreakQS` |
| `527_` | T_CLAIM_FAC_RETRO.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.FacRetroList` |
| `528_` | T_CLAIM_ADJUSTMENT.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.AdjustmentList` (baris berkomite tidak pernah dihapus aplikasi) |
| `529_` | T_CLAIM_ADJ_SPREADING.ADJUSTMENT_ID -> T_CLAIM_ADJUSTMENT - `.AdjustmentList(n).SpreadingAdjustment` |
| `530_` | T_CLAIM_ADJ_QUOTA_SHARE.ADJUSTMENT_ID -> T_CLAIM_ADJUSTMENT - `.AdjustmentList(n).SpreadingQuotaShare` |
| `531_` | T_CLAIM_ADJ_LOSS_ALLOCATION.ADJUSTMENT_ID -> T_CLAIM_ADJUSTMENT - `.AdjustmentList(n).LossAllocation` |
| `532_` | T_VIEW_SUGGEST.CLAIM_ID -> T_GENERAL_CLAIM - riwayat `.ClaimData.SuggestList` (induk kedua, CHECK tepat satu) |

### Nama terlarang di migrasi

Nama yang tidak boleh muncul di migrasi modul MANA PUN (`TestNamaYangDibuangTidakAda`).

| Nama | Sebab |
| --- | --- |
| `FLAG_ON_GOING_COMMITTEE` | `AddKomiteTreatyChild_ACT` langkah 3 `.FlagOnGoingCommitte` dibuang (keputusan work owner 19-09-2026) - tidak ada kolom, tidak ada properti tersimpan |

## Keputusan work owner 07-10-2026 sesudah laporan pertama

- **Migrasi `532` (`T_VIEW_SUGGEST.CLAIM_ID` + CHECK satu induk) dipasang** — "1 tabel aja gabung life dan non life": Claim History satu tabel
  dengan riwayat penawaran PremiumList Life. Baris kolomnya di STRUKTUR PremiumList Life; uji modul itu
  (`strukturtipe_polis_test.go`) menghitung 244 kolom dan melewati kolom buatan modul lain (`kolomModulLain`).
- **Baris kolom 520 di STRUKTUR Claim Life** (`## T_GENERAL_CLAIM`, dokumen saja) = keputusan "Tabel bersama".
- **Penyerahan ke komite tetap dicabut** (penulis dan pembaca tangga, kolom keputusan anggota): penjaga batas Claim
  Life `komite_statik_test.go` menolaknya. Tombol tampil nonaktif: OQ-CP-16, menunggu jawaban.
- **Tambah / hapus baris Spreading Claim nonaktif**, termasuk ikon grid bawaan (dicabut). Akibat yang diketahui:
  kasus BARU tidak dapat melahirkan baris spreading, sehingga Save to issue RNM ditolak `ProteksiData_act` langkah 5
  ("please Fill SpreadingList"); kasus hasil pemuat data lama membawa barisnya.
- **`STS_REJECT = 1` ditunda** sampai modul Komite Claim Prop ("itu nanti kan dari komite") — pemuat mencatatnya
  "ditunda", bukan gagal.

## Pemuat data lama

`backend/alat/pemuatlama` (prompt §6 butir 11; AC 9–12, 123, 132). Sumber baca-saja: `OS_AKSEPTASI_KLAIM`
(CASEID `ASM-FW-GCNMFW-WORK CLMP-%`, baris berlaku per kasus menurut AC 123) dan `JSON_KLAIM` (halaman `.ClaimData`
bila ada). Tujuan: `T_WORK_CLAIM` + `T_GENERAL_CLAIM` (SUMBER `PEGA`, ID = pyID Pega `CLMP-n`) dan tabel `T_CLAIM_*`
lewat `SimpanHalaman` — jalur yang sama dengan aplikasi.

```
go run ./modul/claimprop/backend/alat/pemuatlama -keluaran <folder>             # uji-kering (bawaan, baca saja)
go run ./modul/claimprop/backend/alat/pemuatlama -keluaran <folder> -jalankan   # tulis - HANYA work owner
```

- `-jalankan` ditolak bila `IS_PEGA_PROD=true`; satu transaksi per kasus; kasus yang ID-nya sudah ada dilewati.
- Berkas keluaran: `claimprop-arsip-medan-*.csv` (medan yang tidak masuk kolom + sebab) dan `claimprop-galat-*.csv`.
  Keduanya memuat data kasus — simpan di luar repositori.
- Uji-kering DEV 07-10-2026: 2.451 kasus, 6.596 baris OS (4.145 riwayat), 6 kasus berhalaman JSON, **2.032 siap**,
  **419 ditunda** (baris berlaku `STS_REJECT = 1`, menunggu modul Komite Claim Prop), **0 gagal**; keluar 0.
- Urutan resmi: work owner menjalankan `-migrate` (520–533, 980) → uji-kering ulang → keputusan OQ-CP-18 →
  `-jalankan` oleh work owner / DBA. Tidak pernah dijalankan agen.

## Uji SQL di DEV (baca saja)

`go test -tags ujidev -run TestSQLDiDEV -v ./modul/claimprop/backend/repository/` — setiap SELECT dijalankan lewat
metode aslinya dengan masukan `UJI-*`, setiap INSERT / UPDATE / DELETE hanya diurai `DBMS_SQL.PARSE`. 07-10-2026:
53 ok, 0 gagal, 37 "objek belum ada" (34 menunggu migrasi 520–533; 3 skema luar OQ-CP-11).
