# 22: Pemuat dokumen lama

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Seluruh dokumen polis dipindahkan — setiap polis, setiap generasinya. Tidak ada penyaringan.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 5.

---


**Status:** selesai — pemuat dibangun; F3 diputuskan dan dibangun, F6 gugur, F7 urutan resmi ditulis (putaran 3, 04-10-2026, bab *Putaran 3* di bawah); tertahan hanya pihak luar: AC 55 uji `db` (K11), AC 59 atas data nyata (uji-kering F7 langkah 2), pemuatan di produksi oleh WO/DBA (F7), `SEQ_WORK_POLIS` (C5), empat angka DBA (C7) *(semula putaran 2: "… AC 59 keputusan WO F3, penanda `SUMBER='PEGA'` (F6), pemuatan di produksi (F7) …"; putaran 2, konsolidasi P10 04-10-2026 — rincian `docs/HASIL-IMPLEMENTASI.md` bab 9; semula: sebagian — dibangun, putaran 2 03-10-2026 cabang `modul/nbtreatyin/p7-pemuat`; awalnya ready-for-agent)*
~~**Blocked by:** **19** · **20** · ⛔ `[work owner]` **dokumen lama dipindahkan seluruhnya atau sebagian** — belum diputuskan~~ ⛔ **penahan gugur 23-09-2026**
**Menutup:** NB AC **55–59** *(5 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-3

## Hasil & nilai pengguna

Dokumen polis lama dimuat ke dalam tabel baru **lewat antarmuka yang sama** dengan jalur biasa,
sehingga data lama melewati pemeriksaan yang sama dan tidak ada pintu belakang.

⭐⭐ **Nilai lama tidak pernah dihitung ulang** — disalin apa adanya, lengkap dengan galat
presisinya, supaya laporan lama masih dapat direkonsiliasi.

## Yang dibangun

Pemuat massal yang: membaca dokumen lama · memecahnya dengan pemecah yang sama · menulis lewat
antarmuka penyimpanan yang sama · dan **melaporkan** dokumen yang gagal diurai beserta sebabnya.

⛔ **Dokumen yang gagal tidak dilewati diam-diam.**

## Batas — yang TIDAK termasuk

⛔ Penanda migrasi khas endorsemen — tiket **28**.
⛔ Keputusan lingkup pemindahan — `[work owner]`, lihat di bawah.

## Cara mengujinya

Lewat seam `repository` yang sama. ⭐ Uji utama: dokumen bergalat presisi dimuat, lalu nilainya
**dibaca dari kolomnya** — pembulatan ke presisi mata uang berarti gagal.

## Acceptance criteria

- [ ] 🟡 **AC 55** — dokumen lama dimuat **tanpa pembulatan ke presisi mata uang** *(putaran 2: pemecah membawa `592629512.880000276` utuh — `TestPecahDokumenProporsionalDatar`; nilai kolom `592629512.88000028` diperiksa `TestPemuatLamaMenulisLewatAntarmukaSama` bertag `db`, **belum dijalankan**, K11)*
- [x] **AC 56** — pemuat menulis lewat antarmuka penyimpanan yang **sama** *(putaran 2: `services/pemuat.go` memakai `SisipKasus`, `SimpanHalaman`, `SetelNomorPolis`, `TutupKasus`; jalur tulis terpisah dijaga `repository/lama_test.go` `TestPemuatTanpaJalurTulisTerpisah`)*
- [x] **AC 57** — medan tak dikenal **tersimpan di penampung**, bukan dibuang *(⛔ RALAT K17: penampung = **berkas CSV** `POLIS_ID,JALUR,NILAI`, bukan tabel — `TestLaporanMedanTakDikenalBerkasCSV`)* *(⛔ RALAT F3 04-10-2026: berkas = **arsip audit** `POLIS_ID,JALUR,NILAI,KEPUTUSAN` setiap medan yang tidak masuk kolom, termasuk yang dibuang menurut keputusan tertulis — `TestLaporanArsipMedanTanpaKolomBerkasCSV`)*
- [x] **AC 58** — dokumen yang gagal diurai **dilaporkan beserta sebabnya** *(putaran 2: berkas CSV galat — `TestLaporanGalatBerkasCSVDanTanggalAmbiguDihitung`, `TestDokumenBergalatTidakDimuatDanSebabnyaDisebut`, `TestDokumenDiLuarLingkupAtauRusak`)*
- [ ] 🟡 **AC 59** — penampung medan tak dikenal **wajib kosong** sebelum pekerjaan dinyatakan selesai *(putaran 2: mekanisme dibangun — jumlah dicetak, kode keluar 1 bila > 0, `TestRingkasanSelesaiHanyaBilaNolGalatDanNolTakDikenal`; jumlah atas data nyata belum diketahui sampai pemuat dijalankan work owner)* *(⛔ RALAT F3 04-10-2026: = **nol medan BELUM DIPUTUSKAN**; setiap jalur panduan bentuk dokumen diputuskan — `TestPanduanBentukDokumenNolMedanBelumDiputuskan`; kode keluar ≠ 0 hanya bila ada medan belum diputuskan atau galat — `TestRingkasanSelesaiHanyaBilaNolGalatDanNolBelumDiputuskan`. Tetap 🟡: panduan basi (diagram R45), data nyata baru terlihat pada uji-kering F7 langkah 2 — penahan kini F7/K11, bukan F3)*

## ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Belum diputuskan apakah dokumen lama dipindahkan seluruhnya, sebagian, atau tetap
dibaca lewat jalur lama.** Lingkup pemuat bergantung langsung padanya — memuat seluruh riwayat dan
memuat dua tahun terakhir adalah pekerjaan yang berbeda besarnya.

⚠️ **AC 59 tidak dapat dipenuhi sebelum tiket 19 lepas** — penampung medan tak dikenal tidak akan
kosong selama daftar kolom lengkap belum ada.

## ⭐ Penerapan KEPUTUSAN-RONDE-12 — 2026-10-03

- **Butir 5** — seluruh dokumen `POOLDATA.JSON_POLIS`, setiap polis dan setiap generasinya, dipindah;
  tanpa penyaring. Bunyi lama (*"seluruhnya atau sebagian — belum diputuskan"*) digantikan.
- ⛔ **Pemuat belum dibangun pada implementasi 2026-10-03.** Penampung `T_POLIS_MEDAN_LAIN` (migrasi 329)
  sudah ada; penulisnya lahir bersama pemuat. Empat angka DBA butir 5 (cacah baris, tahun terawal,
  ukuran, PRODKE tertinggi) masih diperlukan untuk merencanakan pemuatan penuh.

## ⭐ Putaran 2 — pemuat DIBANGUN (03-10-2026)

Acuan: PROMPT-NB-TREATY-IN-PUTARAN-2.md bab 0 butir 11–12, bab 2 **K15**, **K17**, bab 5;
PESAN-KOREKSI-PUTARAN-2.md A.3. Perintah dan flag: `MODUL.md` bab *Pemuat dokumen lama*.

### ⛔ RALAT — penampung medan tak dikenal

> Bunyi lama (bab *Penerapan KEPUTUSAN-RONDE-12*): *"⛔ **Pemuat belum dibangun pada implementasi
> 2026-10-03.** Penampung `T_POLIS_MEDAN_LAIN` (migrasi 329) sudah ada; penulisnya lahir bersama pemuat."*
>
> Bunyi baru: `[keputusan work owner]` **K17** — `T_POLIS_MEDAN_LAIN` **tidak ada di diagram grilling**
> dan dihapus (paket penyimpanan). Penampung medan tak dikenal = **berkas laporan CSV per jalankan**
> `nbtreatyin-medan-tak-dikenal-<stempel>.csv` berkolom `POLIS_ID`, `JALUR`, `NILAI` di folder
> `-keluaran` pilihan operator. Medannya tetap tersimpan beserta nilainya, jumlahnya dicetak, dan wajib
> **0** sebelum pekerjaan dinyatakan selesai. Laporan galat juga berkas, bukan tabel. Nol tabel baru.

### ⛔ RALAT — seam uji

> Bunyi lama (*Cara mengujinya*): *"Lewat seam `repository` yang sama."*
>
> Bunyi baru: pemecah dokumen dan penulis laporan adalah **fungsi murni** di `models` (seam 3
> `spec.md` §6.2) supaya uji-kering pemuat berjalan tanpa menyentuh tabel baru; penulisan diuji lewat
> seam `repository` bertag `db` (`repository/lama_db_test.go`, ditulis, **belum dijalankan** — K11).
> `[penyimpangan sadar]` atas spec-penyimpanan ID-4 (*"pemecah … rincian di dalam `repository`"*).

### Yang dibangun — dari XML ke kode

| Bukti XML / dokumen | Kode |
| --- | --- |
| `Activity\SaveJsonPolisTreatyIn_Act.xml` langkah 6: halaman `pyWorkPage.PolicyTreatyIn` (kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`), `InputData.CARI3 = @ASM.GetPageJSONString()` — DATA_JSON = halaman `PolicyTreatyIn` | `models.PecahDokumenLama` menelusuri akar sebagai `PolicyTreatyIn.*`; dokumen berkelas lain (JSON_POLIS dipakai bersama Fac In) hanya dihitung |
| `RDBList\SavePolisTreatyIn_SQL.xml`: `PEGA_JSON_POLIS_TREATYIN({pzInsKey}, {PolicyTreatyIn.PolicyNo}, NULL, '0', {CARI21}, {OperatorID.pyUserIdentifier}, {CARI3})` | IDPEGA → pyID = ID kasus; NOPOLIS = `SetelNomorPolis` (PolicyNo dokumen wajib sama); PRODKE `'0'` = generasi NB; IDPEGA, NOENDORS, TGL_INPUT, USERNAME apa adanya (`SetelKolomDatarLama`, ID-21) |
| `Flow\InputRealizationTreatyIn`: `Decision8` *Nopolis not empty* → `Utility1` (Save json policy) → `Utility2` → `End3` | kasus lama ditutup `Resolved-Completed` (`TutupKasus`), tidak muncul di antrean portal |
| `Start1 → Assignment2`: `.FlagOnGoingPolicy = 1` | `SisipKasus` yang sama dengan jalur biasa |
| spec-penyimpanan ID-19, AC 21–22; P29 sifat 2 | `models.BacaTanggalLama`: `YYYYMMDD`; cap waktu `… GMT` → jam dinding Asia/Jakarta (rule `GeneratePolicyNoTreaty_Act` langkah 5.3 membaca hari dalam Asia/Jakarta; `20170930T170000.000 GMT` = 1-10-2017 00.00 WIB) |
| K15, P32 | `05/06/2017` → `ErrTanggalAmbigu` (tidak ditebak, dihitung); bentuk lain di luar dua format → `ErrFormatTanggal` (cara menentukan susunan per baris belum diputuskan, P32 butir 1) |
| katalog `models/katalog.go` | `petaKatalogDokumen` menurunkan pola jalur dokumen → kolom dari `SemuaTabel` (ikut perubahan katalog otomatis); daftar bersarang lewat `IndukDaftarBersarang` (dijaga uji) |
| spec AC 69, §5.8 `[keputusan work owner]` | EndDate kosong = StartDate |
| KEPUTUSAN-RONDE-12 butir 5; edmtreatyin tiket 10 | seluruh generasi NB tanpa penyaring; PRODKE > 0 hanya dihitung (milik pemuat EDM) |

Medan yang sengaja **tidak** disimpan — dihitung per alasan di ringkasan, tidak hilang diam-diam:
`pxObjClass` (P29 sifat 4, rancangan §4.1) · `pxListSubscript` (= NOURUT, ID-11) · `Show` `ViewState`
`pxResults` `FillPaymentInstallmentEDMT` (keadaan layar, rancangan §4.1/4q5.7) · `Total*` tingkat polis
(turunan, katalog/RALAT AC 38) · `Layer*` tingkat polis (ID-22) · `BreakDownSpreadList` (KR-12 butir 3/3b)
· `isApprovedtoDeptHead` (spec AC 64) · `OldData` (rancangan 4ter.1) · `TreatyDifference`
`TreatyXOLDifferenceList` (rancangan 4ter.2, nol baris di NB). Selain itu — termasuk `pyExpanded` dan
`SuggestList` sesudah `T_POLIS_SUGGEST` dihapus (K4) — masuk berkas CSV medan tak dikenal.

### Butir terbuka

1. ⛔ `[work owner]` **`SuggestList` dokumen lama** — sesudah `T_POLIS_SUGGEST` dihapus (K4) tidak punya
   kolom; ia masuk berkas CSV. Apakah pemuat menyalinnya ke `HISTORYAKSEPTASIPRODUCTION`
   (`InsertViewSuggest_SQL`)? Belum diputuskan — tidak dikarang.
2. ⛔ `[work owner]` **Penanda `SUMBER='PEGA'`** (KEPUTUSAN-RONDE-12 butir 5 poin 2) — tidak ada kolom di
   diagram grilling; baris hasil pemuat dikenali dari IDPEGA berbentuk `<kelas> <pyID>` dan status
   `Resolved-Completed`. ⛔ "tidak ada kolom di diagram grilling".
3. ⛔ `[tim inti / premiumlistlife]` **`SEQ_WORK_POLIS`** — ID kasus lama = pyID `NB-<n>`; sequence wajib
   dimajukan melewati nomor terbesar yang dicetak pemuat sebelum kasus baru dibuat (pola OQ-PL-15).
4. ⛔ `[work owner]` **`-jalankan` di `IS_PEGA_PROD=true` ditolak** (pola `pindahflat`) — pemindahan di
   produksi butuh keputusan tersendiri.
5. ⛔ `[data DBA]` empat angka butir 5 (cacah baris, tahun terawal, ukuran, PRODKE tertinggi) — tetap.
6. Uji repository bertag `db` belum dijalankan (K11 kosong).

## ⭐ Putaran 2 — P9 (04-10-2026): AC 59 / K17 tidak dapat mencapai 0 tanpa keputusan WO

Dasar: tinjauan spec P9 (temuan 4) — bukan perbaikan kode. K17 mewajibkan jumlah medan tak dikenal di berkas CSV
**0** sebelum pekerjaan dinyatakan selesai, tetapi sebagian medan dokumen lama **secara struktur** tidak punya
tempat simpan di delapan tabel diagram (bab 0 butir 11–12). Pemuat tidak mengarang tempatnya; ia menulisnya ke
CSV (`POLIS_ID`, `JALUR`, `NILAI`) dan kode keluarnya ≠ 0. Butir keputusan `[work owner]` yang dibutuhkan:

1. **Medan dokumen lama tanpa kolom.** Dari panduan bentuk dokumen (`docs/dataguide-json-polis.json`, 378 jalur;
   penambal — `DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md` menyebut 394 batas bawah), dipetakan lewat katalog dan alasan
   tertulis pemuat (`models.PecahDokumenLama`), **22 pola** tetap tak dikenal:
   `EDMNo`, `EDMType`, `ProdKe`; `QuotationData.` `BranchCode`, `BranchName`, `BusinessType2`, `CedingCo`,
   `CedingCoName`, `MarketingCode`, `OperatorID`, `SobLsg`, `SobName`, `StatusBusiness`, `StatusSyariah`,
   `TeamGroup`, `TypeFacultative`; `ListInstallment().pyExpanded`, `TreatyXOLList().pyExpanded`; dan empat medan
   `SuggestList()` (butir 2). Ditambah medan yang dibuang paket penyimpanan dari kolom (`IsOJKNopolis`,
   `InstallmentList().PPN/PPh`; `docs/PERBANDINGAN-KOLOM-DIAGRAM.md` bab 11 butir 1). Pilihan per medan:
   **(a)** tambah kolom lewat RALAT diagram grilling, **(b)** nyatakan dibuang dengan alasan tertulis (masuk
   `Diabaikan` ringkasan, bukan CSV), atau **(c)** biarkan di CSV dan longgarkan syarat "wajib 0" K17.
2. **`SuggestList` dokumen lama** (`Date`, `IsApproved`, `OperatorName`, `Suggest`) — butir terbuka 1 di atas:
   salin ke `POOLDATA.HISTORYAKSEPTASIPRODUCTION` lewat pemetaan `InsertViewSuggest_SQL` yang sama dengan jalur
   biasa (`models.UsulanBelumTersimpan`), atau nyatakan dibuang. Belum diputuskan — tidak dikarang.

Sampai kedua butir dijawab, AC 59 spec penyimpanan tetap 🟡 dan tiket ini tidak dapat dinyatakan selesai atas data
nyata. Dicatat juga di PERMINTAAN-TIM-INTI bagian F3.

## ⭐ Putaran 3 — F3, F6, F7 diputuskan dan diterapkan (04-10-2026)

Acuan: `PROMPT-NB-TREATY-IN-PUTARAN-3.md` bab 0 dan bab 2 (**F3**, **F6**, **F7**), keputusan work owner
04-10-2026. Setiap medan dibaca ulang di korpus `NB Treaty In (Done)` (`docs/alat/pemakai.py`, `graf.py`
— terjangkau dari titik masuk nyata), bukan dari klaim putaran 1–2.

### ⛔ RALAT — butir terbuka putaran 2

1. Bunyi lama (butir terbuka 1), dikutip: *"⛔ `[work owner]` **`SuggestList` dokumen lama** — sesudah
   `T_POLIS_SUGGEST` dihapus (K4) tidak punya kolom; ia masuk berkas CSV. Apakah pemuat menyalinnya ke
   `HISTORYAKSEPTASIPRODUCTION` (`InsertViewSuggest_SQL`)? Belum diputuskan — tidak dikarang."* Bunyi baru:
   `[keputusan work owner]` **F3** — **disalin** ke `POOLDATA.HISTORYAKSEPTASIPRODUCTION`, pemetaan
   `Activity\SaveViewSuggest.xml` langkah 2 "UNTUK TREATY" (2.1 `.IsSave==""`; 2.1.2 CARI1..CARI10 →
   `RDBList\InsertViewSuggest_SQL.xml`), penulis `CatatUsulan` jalur biasa, **penjaga dobel menurut `IDPEGA`**
   (`repository.SalinUsulanLama`: IDPEGA yang sudah punya baris tidak disalin lagi — pemuat diulang tidak
   menggandakan baris). `AKSES_LOGIN` / `PIC` dari isi baris (`OperatorID` / `OperatorName`) bila ada; kosong
   ditulis apa adanya (NULL) dan dihitung di ringkasan — tidak dikarang. `CARI7`/`BUSINESS_CODE`
   (`Quotation.BusinessFac/BusinessCode`) dari `QuotationData` dokumen = salinan halaman `Quotation`
   (`Activity\GeneratePolicyNoTreaty_Act.xml` langkah 10 Page-Copy). `TGL_INP` dari `.Date` lewat
   `BacaTanggalLama` (jam dinding Asia/Jakarta, 24 jam — F2 butir 2); tanggal ambigu = galat dokumen (K15).
2. Bunyi lama (butir terbuka 2), dikutip: *"⛔ `[work owner]` **Penanda `SUMBER='PEGA'`** … ⛔ \"tidak ada kolom
   di diagram grilling\"."* Bunyi baru: `[keputusan work owner]` **F6** — penanda **gugur**; baris hasil pemuat
   dikenali dari `IDPEGA` (`<kelas> <pyID>`, jalur biasa menulis `NB-<n>`) dan status `Resolved-Completed`.
   RALAT `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md` butir 5 poin 2.
3. Bunyi lama (butir terbuka 4), dikutip: *"⛔ `[work owner]` **`-jalankan` di `IS_PEGA_PROD=true` ditolak** (pola
   `pindahflat`) — pemindahan di produksi butuh keputusan tersendiri."* Bunyi baru: `[keputusan work owner]`
   **F7** — tidak dijalankan agen; urutan resmi di `MODUL.md` bab *Migrasi*: migrasi 320–327 oleh WO →
   uji-kering pemuat di skema uji (K11) → F3 tuntas (nol medan `BELUM DIPUTUSKAN`) → pemuatan produksi oleh
   WO/DBA. Penolakan `IS_PEGA_PROD=true` tetap (ADR-U-0005, sama dengan `-migrate`). spec AC 68 tetap 🟡.
4. Bunyi lama (bab *Yang dibangun*), dikutip: *"Selain itu — termasuk `pyExpanded` dan `SuggestList` sesudah
   `T_POLIS_SUGGEST` dihapus (K4) — masuk berkas CSV medan tak dikenal."* Bunyi baru: `pyExpanded` **dibuang**
   berbukti (tabel di bawah); `SuggestList` **disalin** (butir 1). Berkas CSV menjadi **arsip audit pemuatan**
   `nbtreatyin-arsip-medan-<stempel>.csv` (`POLIS_ID`, `JALUR`, `NILAI`, `KEPUTUSAN`) — setiap medan yang tidak
   masuk kolom, `dibuang: <kunci alasan>` atau `BELUM DIPUTUSKAN`; hanya yang terakhir menahan selesai (kode
   keluar ≠ 0). Bab *P9* di atas (pilihan (a)/(b)/(c)) terjawab: (a) dan (b) per medan, (c) tidak dipakai.

### Keputusan per medan (F3)

Kunci alasan dan bukti tercatat di `backend/models/medan_abaikan_lama.json` bagian `pola` (berkas data:
penjaga claimlife `polaIndeksPosisi`); uji `TestPenggolongMedanDiabaikan`,
`TestPanduanBentukDokumenNolMedanBelumDiputuskan` (seluruh jalur daun panduan bentuk dokumen terputuskan).

| Medan (relatif `PolicyTreatyIn`) | (a)/(b) | Tempat / keputusan | Bukti XML |
| --- | :---: | --- | --- |
| `EDMType` | **a** | **kolom** `T_GENERAL_POLIS.EDM_TYPE` (RALAT rancangan §4sexies; rancangan §4.1 *penentu bentuk* `EDM_TYPE`) | `Activity\InputPolicyTreatyInPre_Act.xml` langkah 10 prasyarat `.PolicyTreatyIn.EDMType=="3"` → lewati `Call TreatyRealizationCheckXOLList` (terjangkau; `models.PerluCekDaftarXOL`) |
| `SuggestList().Date/IsApproved/OperatorName/Suggest` (+ `OperatorID`, `IsSave`) | — | **disalin** `HISTORYAKSEPTASIPRODUCTION` (butir 1) | `SaveViewSuggest` langkah 2.1–2.1.4; `InsertViewSuggest_SQL` |
| `EDMNo` | b | dibuang `f3_salinan_generasi` | 0 rujukan korpus NB; rancangan §4.1 `EDM_NO` = turunan NOPOLIS + PRODKE (P55) |
| `ProdKe` | b | dibuang `f3_salinan_generasi` | 0 rujukan korpus NB; generasi di kolom `PRODKE` (`SavePolisTreatyIn_SQL` `'0'`) |
| `IsOJKNopolis` | b | dibuang `f3_tanpa_pembaca` | satu-satunya rujukan `GeneratePolicyNoTreaty_Act` langkah 23 `.IsOJKNopolis = "1"` — **berlabel `//`**; nol pembaca |
| `BrokerageFee` | b | dibuang | ditulis `SetPPNPPH` langkah **4** (bukan 3 seperti tertulis di PERBANDINGAN 1d); nol pembaca |
| `pyMessageLabel` | b | dibuang | ditulis `ProtectDate` langkah 1 (dipanggil `Section\DetailPolicyTreatyIn`); nol pembaca |
| `QuotationData.BranchCode`, `.BranchName` | b | dibuang | ditulis `CheckDataMkt` langkah 4 (`pyWorkPage.Quotation.*`, masuk dokumen lewat Page-Copy `GeneratePolicyNoTreaty_Act` 10); nol pembaca |
| `QuotationData.TeamGroup`, `.MarketingCode` | b | dibuang | ditulis `CheckDataMkt` langkah 4; pembaca `.MarketingCode` hanya langkah 6 berlabel `//`; `When\IsTBonding` tak terjangkau |
| `QuotationData.BusinessType` | b | dibuang | ditulis `InputPolicyTreatyInDetail_preACT` 14.8 (DT `BusinessType_DeT`); pembaca hanya rule tak terjangkau (`When\Is*`, `ProtectCoverage_Act`, `ProtectFIREMBUPA_Act`) |
| `QuotationData.CedingCo`, `.CedingCoName` | b | dibuang (tingkat polis tetap `T_GENERAL_POLIS.CEDING_CO/_NAME`, R47) | ditulis preACT 3, `SearchHierarkiSourceBizAgent_PostDT` 2.1/2.2; pembaca hanya treaty keluar (K8 butir 4) / halaman `OfferFacIn` tak terjangkau |
| `QuotationData.SobName`, `.SobLeader0`, `.SobLeader1` | b | dibuang | ditulis preACT 3, `SearchHierarkiSourceBizAgent_PostDT` 1.2/4/5; pembaca hanya `CheckDataMkt` 7 (`//`) dan `InputPolicyTreatyOutDetail_preACT` 4 (treaty keluar) |
| `QuotationData.StatusBusiness` | b | dibuang | 9 rule perujuk, seluruhnya tak terjangkau (`Protection_Act`, `CheckSpreadingProtect_ACT`, `When\IsEDM`, …) |
| `QuotationData.OperatorID`, `.BusinessType2`, `.SobLsg`, `.StatusSyariah`, `.TypeFacultative` | b | dibuang | 0 penulis dan 0 pembaca di korpus NB (diisi di luar modul) |
| `ListInstallment().InstallmentList().PPN`, `.PPh` | b | dibuang | ditulis preACT 18.3.4.2.1; pembaca `.PPN/.PPh` (18.3.4.1) adalah baris `ListInstallment` — kolom `T_POLIS_INSTALMENT.PPN/PPH` |
| `ListInstallment().pyExpanded`, `TreatyXOLList().pyExpanded` | b | dibuang `f3_keadaan_baris` | ditulis `InputPolicyTreatyInDetail_NonProp` 12.1 (`.pyExpanded = true`); nol pembaca properti |

⚠️ Panduan bentuk dokumen basi (diagram R45: satu dokumen memuat 95 jalur di luar panduan). Medan di luarnya
baru terlihat pada uji-kering atas data nyata (F7 langkah 2); pemuat mencatatnya `BELUM DIPUTUSKAN` di arsip,
kode keluar ≠ 0, dan ia diputuskan dengan cara yang sama (pola + bukti) sebelum pemuatan produksi.

### Yang dibangun — dari XML ke kode

| Bukti XML / keputusan | Kode | Uji |
| --- | --- | --- |
| F3 (b), 25 pola | `models/medan_abaikan_lama.json` bagian `pola` (alasan + bukti wajib, `muatPenggolongAbaikan` menolak pola tanpa bukti); `kunciDibuang` | `TestPenggolongMedanDiabaikan`, `TestBerkasPenggolongCacatDitolak` |
| F3 (a) `EDMType` | `models/katalog.go` `kKode(pt+"EDMType", "EDM_TYPE", 16)`; migrasi 320 + `STRUKTUR` dibangkitkan `docs/alat/skema.py` | `TestPetaKatalogDokumenDigerakkanKatalog`, `TestTabelDanKolomMengikutiDiagramGrilling` (tetap tepat delapan `CREATE TABLE`) |
| `SaveViewSuggest` 2.1–2.1.2 | `models/usulanlama.go` `UsulanDokumenLama` (pemetaan bersama `petaUsulan` dengan `UsulanBelumTersimpan`); `PecahDokumenLama` membaca `.Date` | `TestSuggestListLamaDisalinMenurutSaveViewSuggest`, `TestSuggestListLamaTanggalAmbiguMenggagalkanDokumen` |
| penjaga dobel `IDPEGA` | `repository/lama.go` `SalinUsulanLama` (baca `COUNT(*) … WHERE IDPEGA = :1`, tulis lewat `CatatUsulan`); `services/pemuat.go` `muat` dalam transaksi dokumen | `TestSQLPemuatLamaBerskemaTanpaCommit`; db `TestPemuatLamaMenyalinSuggestListSekaliMenurutIDPega` (ditulis, **belum dijalankan** — K11) |
| arsip audit, AC 59 RALAT | `models/laporanlama.go` (`KepalaArsipMedan`, `KeputusanBelumDiputuskan`, `Selesai`), `alat/pemuatlama/main.go` | `TestLaporanArsipMedanTanpaKolomBerkasCSV`, `TestRingkasanSelesaiHanyaBilaNolGalatDanNolBelumDiputuskan`, `TestRingkasanSalinanUsulanLama` |
