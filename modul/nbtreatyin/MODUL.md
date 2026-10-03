# Modul `nbtreatyin` — NB Treaty In

Realisasi penutupan treaty inward: dari kontrak treaty yang disetujui menjadi polis treaty, melalui
tangga **Admin → Sec Head → Dept Head** (`docs/spec.md`). Satu folder, satu modul, satu pemilik: kode
backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur tim satu folder per modul,
keputusan work owner 30-09-2026). Padanan Pega: Harness `SFAPortalOpportunities` dan
`Flow/InputRealizationTreatyIn` — seluruh rule terjangkau beserta nasibnya (dibangun/tidak + alasan)
di `docs/INVENTARIS-XML.md`; hasil per acceptance criteria di `docs/HASIL-IMPLEMENTASI.md`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `nbtreatyin` |
| Folder korpus | `NB Treaty In` |
| GROUPMENU | `TREATY` |
| Pemilik | `@PEMILIK-NBTREATYIN` |
| Status | dimigrasi |
| Rentang migrasi | `320-359` |
| Slot menu | `968-969` |
| Prefix rute API | `/api/nb-treaty-in` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | spec, tiket (`issues/`), grilling, `INVENTARIS-XML.md` (bangkitan `docs/alat/`), `STRUKTUR-TABEL-NB-TREATY-IN.md`, `HASIL-IMPLEMENTASI.md`, `PERMINTAAN-TIM-INTI.md` |
| `backend/` | `models/` (fungsi murni: halaman kerja, katalog, rantai uang, tangga, layar, pemecah dokumen lama) `repository/` (seluruh SQL) `services/` (aturan, transaksi) `handlers/` (HTTP) `migrations/` `alat/pemuatlama/` (perintah pemuat dokumen lama - dijalankan manusia, bukan bagian aplikasi) `modul.go` |
| `frontend/` | `pages/` `components/` `api.ts` `labels.ts` `medan.ts` `menu.ts` `rute.tsx` `nbtreatyin.css` dan `*.test.ts` |

## Aturan (ringkas)

- **Sumber data kontrak** = view `POOLDATA.TREATYINDETAILJOINEDM` (pilih bisnis) dan tabel
  `TREATYINDETAIL` (grid popup) — bukan JSONDATA (P29). Nol penulisan JSON.
- **Penyimpanan** = `T_WORK_POLIS` (tabel kasus lintas-lini, migrasi premiumlistlife 050/057/059 —
  modul ini MEMBUTUHKAN premiumlistlife aktif lebih dulu) + `T_GENERAL_POLIS` dan anak `T_POLIS_*`
  (migrasi 320-328) menurut katalog `backend/models/katalog.go`; satu transaksi per tindakan.
- **Tangga** tiga posisi; berkas menunggu POSISI (workbasket), keanggotaan diperiksa menurut nama
  antrean (`inti.Pelaku.Peran`). Sec Head menyetujui → selalu Dept Head (AC 8, keputusan work owner —
  bertentangan dengan `CekLimitTreatyAcc_Act`, dicatat). Dept Head menerbitkan nomor polis
  (`inti/backend/penomor`, `KODE_PRODUKSI` NONLIFE).
- **Riwayat** `HISTORYAKSEPTASIPEGA` di transaksi submit; `OPERATORID` = identitas login, `USERNAME`
  = nama tampilan (`M_LOGIN_GO.NAME`).
- **Peran pengganti nama orang** = `M_NBTRIN_PERAN_TEMPAT` (migrasi 330), DIISI KEMUDIAN oleh IAM;
  tanpa baris = tempat tertunda (tiket 05).

## Migrasi

Rentang `320-359`: 320-328 tabel polis (bangkitan `docs/alat/skema.py` dari katalog), 329
`T_POLIS_MEDAN_LAIN` (penampung medan tak dikenal untuk pemuat dokumen lama), 330
`M_NBTRIN_PERAN_TEMPAT`. Slot menu `968`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`.

## Pemuat dokumen lama (tiket 22)

Memindah seluruh dokumen polis generasi NB (`PRODKE 0`) dari `POOLDATA.JSON_POLIS` ke 8 tabel diagram
grilling - setiap polis, tanpa penyaring (KEPUTUSAN-RONDE-12 butir 5). Dijalankan **manusia** dari akar
repo, sesudah lingkungan dimuat seperti `cmd/api` (`ORACLE_DSN`, `ORACLE_SCHEMA`, `IS_PEGA_PROD`); tidak
pernah berjalan saat aplikasi menyala. Migrasi 320-327 wajib sudah dijalankan work owner sebelum `-jalankan`.

```powershell
go run ./modul/nbtreatyin/backend/alat/pemuatlama -keluaran D:\laporan-pemuat            # uji-kering
go run ./modul/nbtreatyin/backend/alat/pemuatlama -keluaran D:\laporan-pemuat -jalankan  # tulis
```

| Flag | Arti |
| --- | --- |
| `-keluaran <folder>` | wajib; folder berkas laporan, dibuat bila belum ada |
| `-jalankan` | tulis ke tabel baru, satu transaksi per dokumen, aman diulang (kasus ber-IDPEGA sama dilewati). Tanpa flag ini: uji-kering - hanya `JSON_POLIS` yang dibaca, nol pernyataan ke tabel baru. Ditolak bila `IS_PEGA_PROD=true` |

Berkas per jalankan (`<stempel>` = `YYYYMMDD-HHMMSS`):

| Berkas | Kolom | Aturan |
| --- | --- | --- |
| `nbtreatyin-medan-tak-dikenal-<stempel>.csv` | `POLIS_ID`, `JALUR`, `NILAI` | medan dokumen tanpa kolom katalog dan tanpa keputusan tertulis - disimpan di sini, bukan tabel (K17). **Wajib 0 baris data** sebelum pekerjaan dinyatakan selesai (spec-penyimpanan AC 57, 59) |
| `nbtreatyin-galat-<stempel>.csv` | `IDPEGA`, `NOPOLIS`, `JALUR`, `NILAI`, `SEBAB` | dokumen yang tidak dimuat beserta sebabnya (AC 58); tanggal ambigu tidak ditebak (K15) |

Ringkasan dicetak ke layar (cacah per jenis galat, per alasan medan diabaikan, per pola medan tak dikenal,
nomor kasus `NB-` terbesar yang dimuat - `SEQ_WORK_POLIS` wajib dimajukan melewatinya). Kode keluar 0
hanya bila nol dokumen gagal dan nol medan tak dikenal. Penulisan lewat antarmuka yang sama dengan
aplikasi (`SisipKasus`, `SimpanHalaman`, `SetelNomorPolis`, `TutupKasus`, ditambah kolom datar
json_polis); `JSON_POLIS` hanya dibaca. Generasi endorsemen (`PRODKE > 0`) milik pemuat EDM
(`modul/edmtreatyin`, tiket 10) dan hanya dihitung.

## Menjalankan uji modul ini saja

Dari akar repo:

```powershell
go test ./modul/nbtreatyin/...
npx vitest run modul/nbtreatyin
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per
butir. Judul yang tidak ada berarti modul ini tidak menyatakan apa pun untuk jenis itu.
