# 22: Pemuat dokumen lama

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Seluruh dokumen polis dipindahkan — setiap polis, setiap generasinya. Tidak ada penyaringan.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 5.

---


**Status:** ⭐ **sebagian — dibangun** *(putaran 2, 03-10-2026, cabang `modul/nbtreatyin/p7-pemuat`; semula: ~~belum~~ *(implementasi 2026-10-03)* · ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*)*. Tersisa: uji repository bertag `db` belum dijalankan (K11 kosong) dan jumlah medan tak dikenal atas data nyata belum diketahui (AC 59) — lihat bab terakhir.
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
- [x] **AC 57** — medan tak dikenal **tersimpan di penampung**, bukan dibuang *(⛔ RALAT K17: penampung = **berkas CSV** `POLIS_ID,JALUR,NILAI`, bukan tabel — `TestLaporanMedanTakDikenalBerkasCSV`)*
- [x] **AC 58** — dokumen yang gagal diurai **dilaporkan beserta sebabnya** *(putaran 2: berkas CSV galat — `TestLaporanGalatBerkasCSVDanTanggalAmbiguDihitung`, `TestDokumenBergalatTidakDimuatDanSebabnyaDisebut`, `TestDokumenDiLuarLingkupAtauRusak`)*
- [ ] 🟡 **AC 59** — penampung medan tak dikenal **wajib kosong** sebelum pekerjaan dinyatakan selesai *(putaran 2: mekanisme dibangun — jumlah dicetak, kode keluar 1 bila > 0, `TestRingkasanSelesaiHanyaBilaNolGalatDanNolTakDikenal`; jumlah atas data nyata belum diketahui sampai pemuat dijalankan work owner)*

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
