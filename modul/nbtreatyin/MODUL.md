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
  ⭐ Pengecualian K8 (03-10-2026): master jalur NonProp/XOL dibaca BACA-SAJA dari `JSONDATA`
  `M_TREATY_IN` / `M_TREATY_IN_EDM` di satu fungsi `repository.MasterXOLDariJSON`
  (`services.PembacaMasterTreaty`, kelak kontrak modul `treatyin`); treaty keluar tidak dibaca.
- **Penyimpanan** = `T_WORK_POLIS` (tabel kasus lintas-lini, migrasi premiumlistlife 050/057/059 —
  modul ini MEMBUTUHKAN premiumlistlife aktif lebih dulu) + `T_GENERAL_POLIS` dan anak `T_POLIS_*`
  (migrasi 320-327, tepat delapan tabel diagram grilling) menurut katalog `backend/models/katalog.go`;
  satu transaksi per tindakan.
- **Tangga** tiga posisi; berkas menunggu POSISI (workbasket), keanggotaan diperiksa menurut nama
  antrean (`inti.Pelaku.Peran`). Sec Head menyetujui → selalu Dept Head (AC 8, keputusan work owner —
  bertentangan dengan `CekLimitTreatyAcc_Act`, dicatat). Dept Head menerbitkan nomor polis
  (`inti/backend/penomor`, `KODE_PRODUKSI` NONLIFE).
- **Riwayat** `HISTORYAKSEPTASIPEGA` di transaksi submit; `OPERATORID` = identitas login, `USERNAME`
  = nama tampilan (`M_LOGIN_GO.NAME`).
- **Catatan usulan** (`PolicyTreatyIn.SuggestList`) = tabel lama `POOLDATA.HISTORYAKSEPTASIPRODUCTION`
  (keputusan work owner K4 03-10-2026; pengganti `SaveViewSuggest -> InsertViewSuggest_SQL`), ditulis
  di transaksi submit SETIAP jenjang dan dibaca balik untuk layar. `[penyimpangan sadar]` K4: syarat
  `BusinessFac == "F"` tidak ditiru; `[penyimpangan sadar — disetujui WO 04-10-2026]` (F2): tiga jenjang,
  TGL_INP 24 jam, NOURUT dari repository (= `.pxListSubscript`, terbukti uji) - rincian di
  `backend/models/usulan.go`. Tabel warisan: tidak dibuat, tidak diubah strukturnya.
- **Peran pengganti nama orang** = konstanta kode `backend/models/peran_tempat.go`
  (`PemetaanPeranTempat`, keputusan work owner K16 03-10-2026 - BUKAN tabel), KOSONG sampai IAM
  menjawab; kedua belas tempat tiket 05 terdaftar di `DaftarTempat`; tanpa baris = tempat tertunda.
  Peran pengguna dari `inti.Pelaku.Peran`.

## Migrasi

Rentang `320-359`: 320-327 - TEPAT delapan tabel diagram grilling (`Diagram-Skema-Tabel-NusantaraRe.xlsx`
sheet *NB Treaty In Prop* / *NonProp*), bangkitan `docs/alat/skema.py` dari katalog. Tidak ada tabel lain
(bab 0 butir 11 PROMPT putaran 2; K4, K16, K17): catatan usulan ke tabel warisan di bawah, pemetaan
peran-tempat konstanta kode, medan tak dikenal pemuat dokumen lama ke berkas laporan CSV.
Perbandingan kolom lawan diagram: `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`. Slot menu `968`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`.

⛔ **Penahan migrasi — tulis di deskripsi PR (K18, `docs/PERMINTAAN-TIM-INTI.md` C10).** `POOLDATA.T_GENERAL_POLIS`
sudah ada: tabel FacIn 7 kolom dari migrasi `182_t_general_polis` (`nbfacin`, di luar repo ini); diagram grilling
memberi nama itu ke Treaty In. Sampai WO dan pemilik `nbfacin` memutuskan C10, `-migrate` dengan modul ini aktif
**berhenti di pra-terbang inti** (`praTerbangBentuk`) sebelum satu pernyataan pun dikirim — seluruh langkah yang
belum tercatat, semua modul, ikut tertahan — dengan galat *"tabel T_GENERAL_POLIS sudah ada … tetapi BENTUKNYA
BERBEDA"*; tidak ada yang berubah. Migrasi 320 sengaja tanpa blok PL/SQL (koreksi WO 04-10-2026); bentuk yang
dibandingkan pra-terbang dikunci `backend/migrasi_test.go`.

| Migrasi | Tabel yang dibuat | Sheet diagram |
| --- | --- | --- |
| 320 | `T_GENERAL_POLIS` | Prop + NonProp F9–F33 |
| 321 | `T_POLIS_QUOTATION` | Prop + NonProp J35–J38 |
| 322 | `T_POLIS_CEDING` | Prop + NonProp R40–R50 |
| 323 | `T_POLIS_INSTALMENT` | Prop J52–J56 · NonProp J52–J55 |
| 324 | `T_POLIS_INSTALMENT_DETAIL` | NonProp R57–R61 (NonProp saja) |
| 325 | `T_POLIS_SPREADING` | Prop J66–J69 · NonProp J71–J74 |
| 326 | `T_POLIS_XOL` | NonProp J76–J80 (NonProp saja) |
| 327 | `T_POLIS_XOL_LAYER` | NonProp R82–R87 (NonProp saja) |

## Tabel warisan yang dibaca dan ditulis

Tidak dibuat dan tidak diubah strukturnya oleh modul ini; ditulis hanya bila diagram grilling menyebut penulisnya.

| Tabel | Akses | Dasar |
| --- | --- | --- |
| `T_WORK_POLIS`, `SEQ_WORK_POLIS` | tulis + baca | akar diagram (B5); milik premiumlistlife |
| `HISTORYAKSEPTASIPEGA` | tulis + baca | riwayat akseptasi (`InsertHistoryAkseptasiPega_Sql`, diagram Prop F98–F99) |
| `HISTORYAKSEPTASIPRODUCTION` | tulis + baca | catatan usulan (`SaveViewSuggest -> InsertViewSuggest_SQL`, diagram Prop J74–J76; K4) |
| `GENERATE_SEQUENCE_NUMBER` | tulis lewat `inti/backend/penomor` | deret nomor polis (padanan `PROC_GENERATE_SEQUENCE_NUMBER`); diagram F103/F118 menyebutnya "dibaca saja" — penulisan lewat penomor bersama disetujui WO 04-10-2026 (`docs/PERMINTAAN-TIM-INTI.md` F8; RALAT catatan di `docs/rancangan-tabel-datar-treaty-in.md` §4bis.4) |
| `TANGGAL_CLOSING`, `KODE_PRODUKSI`, `CURRENCY`, `BUSINESS`, `REINSURANCETYPE`, `TREATYGROUP`, `MARKETINGOFFICER`, `CLIENT`, `AGENT`, `M_LOGIN_GO`, view `TREATYINDETAILJOINEDM`, `TREATYINDETAIL`, `TREATYINPRODUCTION` | baca | RD/RDB terjangkau; `TREATYINPRODUCTION` ditulis modul EDM |
| `M_TREATY_IN`, `M_TREATY_IN_EDM` | baca (JSON, satu fungsi `repository.MasterXOLDariJSON`) | K8 |
| `JSON_POLIS` | baca (pemuat dokumen lama saja) | tiket 22 |

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

Uji lawan Oracle (`-tags=db`, seam repository; belum masuk `make test-db` — PERMINTAAN A4):

```powershell
go test -tags=db ./modul/nbtreatyin/...
```

| Env | Isi (nama saja — nilai tidak pernah ditulis di repo) |
| --- | --- |
| `ORACLE_DSN` | sambungan ke instance yang memuat skema uji |
| `ORACLE_SCHEMA` | skema uji **kosong** dari DBA — **bukan** `POOLDATA` dan tidak memuat kata itu (PERMINTAAN C9) |
| `ORACLE_SKEMA_UJI` | `true` — pengakuan sadar: uji memasang migrasi semua modul lalu **membuang** tabelnya di skema itu |
| `IS_PEGA_PROD` | tidak `true` (ditolak) |

Tanpa `ORACLE_DSN` seluruh uji db **melewati**; dengan DSN tetapi skema tidak diakui atau memuat `POOLDATA`, uji
**gagal** (bukan melewati). Objek warisan yang tidak ada di skema uji dibuat tiruan sementara lalu dibuang (U1,
`backend/repository/tiruan_db_test.go`) — kecuali view `TREATYINDETAILJOINEDM` untuk spec AC 89, yang menuntut view
sungguhan.

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per
butir. Judul yang tidak ada berarti modul ini tidak menyatakan apa pun untuk jenis itu.

### Tabel warisan: dibaca, tidak dibuat

Tabel yang `docs/STRUKTUR-TABEL-NB-TREATY-IN.md` gambarkan tetapi SENGAJA tidak dibuat migrasi mana pun
(`TestKolomDDLCocokDenganStruktur`, `TestTabelBukanMilikKitaTidakDibuat`). Mencabut satu baris =
kepemilikan tabel berpindah — keputusan work owner.

| Tabel | Alasan |
| --- | --- |
| `HISTORYAKSEPTASIPRODUCTION` | tabel warisan POOLDATA (15 kolom, sudah datar); catatan `SuggestList` ditulis dan dibaca di sini menurut diagram grilling (`SaveViewSuggest -> InsertViewSuggest_SQL`) dan keputusan work owner K4 03-10-2026 — tidak dibuat, tidak diubah strukturnya |
