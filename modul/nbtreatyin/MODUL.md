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

- **Sumber data kontrak** = view `POOLDATA.TREATYINDETAILJOINEDM` — pilih bisnis dan grid popup
  (RD `BrowseTreatyJoinEDM`, putaran 3 R5) — bukan JSONDATA (P29). Nol penulisan JSON.
  ⭐ Pengecualian K8 (03-10-2026): master jalur NonProp/XOL dibaca BACA-SAJA dari `JSONDATA`
  `M_TREATY_IN` / `M_TREATY_IN_EDM` di satu fungsi `repository.MasterXOLDariJSON`
  (`services.PembacaMasterTreaty`, kelak kontrak modul `treatyin`); treaty keluar tidak dibaca.
- **Penyimpanan** = `T_WORK_POLIS` (tabel kasus lintas-lini, migrasi premiumlistlife 050/057/059 —
  modul ini MEMBUTUHKAN premiumlistlife aktif lebih dulu) + `T_GENERAL_POLIS_TREATY` dan anak `T_POLIS_*`
  (migrasi 320-327, tepat delapan tabel diagram grilling) menurut katalog `backend/models/katalog.go`;
  satu transaksi per tindakan.
- **Tangga** tiga posisi; berkas menunggu POSISI (workbasket), keanggotaan diperiksa menurut nama
  antrean (`inti.Pelaku.Peran`). Sec Head menyetujui → selalu Dept Head (AC 8, keputusan work owner —
  bertentangan dengan `CekLimitTreatyAcc_Act`, dicatat). Dept Head menerbitkan nomor polis
  (`inti/backend/penomor`, `KODE_PRODUKSI` NONLIFE).
- **Portal** = berkas buatan akun ini (`w.CREATE_OP`) + LINI non-life, switch *In Progress* / *Resolved*
  (`?status=selesai`). RALAT WO 07-10-2026 *"TAMBAHKAN KAN UNTUK PEMBUAT. MENU ITU HANYA UNTUK SI PEMBUAT, NB DAN EDM
  TREATY"*: Resolved pun hanya buatan akun ini (dulu 06-10-2026 semua berkas selesai).
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
peran-tempat konstanta kode, medan dokumen lama tanpa kolom diputuskan per medan (F3) dan diarsipkan di
berkas CSV pemuat.
Perbandingan kolom lawan diagram: `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`. Slot menu `968`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`.

✅ **K18 / C10 selesai** — perintah work owner 05-10-2026: tabel induk NB Treaty In diganti nama `T_GENERAL_POLIS_TREATY` (migrasi `320_t_general_polis_treaty`); `T_GENERAL_POLIS` tetap milik `nbfacin` dan tidak dipakai modul ini. Nama constraint `GP_TREATY` (≤ 30 byte) supaya tidak bentrok dengan
constraint tabel lama di skema yang sudah menjalankan 320 lama. Skema seperti itu (DEV) dipindah **skrip transisi**
di `SCRIPT-TABEL-KOLOM-BARU.xlsx` sheet NB TREATY (salin baris NB Treaty In, pindah FK anak, hapus baris lama);
kolom treaty yang terlanjur menempel di `T_GENERAL_POLIS` dibiarkan (keputusan WO 05-10-2026). Bentuk 320 dan
larangan menyebut `T_GENERAL_POLIS` dikunci `backend/migrasi_test.go`.

| Migrasi | Tabel yang dibuat | Sheet diagram |
| --- | --- | --- |
| 320 | `T_GENERAL_POLIS_TREATY` | Prop + NonProp F9–F33 |
| 321 | `T_POLIS_QUOTATION` | Prop + NonProp J35–J38 |
| 322 | `T_POLIS_CEDING` | Prop + NonProp R40–R50 |
| 323 | `T_POLIS_INSTALMENT` | Prop J52–J56 · NonProp J52–J55 |
| 324 | `T_POLIS_INSTALMENT_DETAIL` | NonProp R57–R61 (NonProp saja) |
| 325 | `T_POLIS_SPREADING` | Prop J66–J69 · NonProp J71–J74 |
| 326 | `T_POLIS_XOL` | NonProp J76–J80 (NonProp saja) |
| 327 | `T_POLIS_XOL_LAYER` | NonProp R82–R87 (NonProp saja) |
| 328 | `T_POLIS_SURVEY` | — di luar diagram: survei historis, Prop saja (keputusan WO 06-10-2026, membatalkan K7) |

`T_GENERAL_POLIS_TREATY.EDM_TYPE` (`PolicyTreatyIn.EDMType`) ditambahkan putaran 3 (F3, RALAT rancangan §4sexies):
medan dokumen lama yang dibaca syarat `InputPolicyTreatyInPre_Act` langkah 10.

### Urutan resmi migrasi dan pemuatan dokumen lama (F7, keputusan work owner 04-10-2026)

Agen **tidak** menjalankan migrasi maupun pemuatan apa pun. Urutannya:

1. **Migrasi 320–327 oleh work owner** (`-migrate`; menolak `IS_PEGA_PROD=true`, ADR-U-0005). Skema yang
   sudah menjalankan 320 lama (`T_GENERAL_POLIS`): sesudah `-migrate`, jalankan skrip transisi sheet NB TREATY
   `SCRIPT-TABEL-KOLOM-BARU.xlsx`, lalu mulai ulang backend.
2. **Uji-kering pemuat di skema uji** (K11) — tanpa `-jalankan`; hanya `JSON_POLIS` yang dibaca, arsip medan
   dan laporan galat ditulis (bab *Pemuat dokumen lama*).
3. **F3 tuntas**: arsip uji-kering nol baris `KEPUTUSAN = BELUM DIPUTUSKAN` dan nol dokumen gagal (kode
   keluar 0). Medan baru yang muncul diputuskan per medan dengan bukti XML di
   `backend/models/medan_abaikan_lama.json` (bagian `pola`) atau diberi kolom lewat RALAT (bab 0 butir 12).
4. **Pemuatan produksi oleh WO/DBA** (`-jalankan`, satu transaksi per dokumen, aman diulang). Sesudahnya
   `SEQ_WORK_POLIS` dimajukan melewati nomor kasus terbesar yang dicetak pemuat (C5).

spec AC 68 (berkas lama dapat dibuka) tetap 🟡 sampai langkah 4.

## Tabel warisan yang dibaca dan ditulis

Tidak dibuat dan tidak diubah strukturnya oleh modul ini; ditulis hanya bila diagram grilling menyebut penulisnya.

| Tabel | Akses | Dasar |
| --- | --- | --- |
| `T_WORK_POLIS`, `SEQ_WORK_POLIS` | tulis + baca | akar diagram (B5); milik premiumlistlife |
| `HISTORYAKSEPTASIPEGA` | tulis + baca | riwayat akseptasi (`InsertHistoryAkseptasiPega_Sql`, diagram Prop F98–F99) |
| `HISTORYAKSEPTASIPRODUCTION` | tulis + baca | catatan usulan (`SaveViewSuggest -> InsertViewSuggest_SQL`, diagram Prop J74–J76; K4) |
| `GENERATE_SEQUENCE_NUMBER` | tulis lewat `inti/backend/penomor` | deret nomor polis (padanan `PROC_GENERATE_SEQUENCE_NUMBER`); diagram F103/F118 menyebutnya "dibaca saja" — penulisan lewat penomor bersama disetujui WO 04-10-2026 (`docs/PERMINTAAN-TIM-INTI.md` F8; RALAT catatan di `docs/rancangan-tabel-datar-treaty-in.md` §4bis.4) |
| `TANGGAL_CLOSING`, `KODE_PRODUKSI`, `CURRENCY`, `BUSINESS`, `REINSURANCETYPE`, `TREATYGROUP`, `MARKETINGOFFICER`, `CLIENT`, `AGENT`, `M_LOGIN_GO`, view `TREATYINDETAILJOINEDM` | baca | RD/RDB terjangkau |
| `M_TREATY_IN`, `M_TREATY_IN_EDM` | baca (JSON, satu fungsi `repository.MasterXOLDariJSON`) | K8 |
| `M_LOGIN_GO_WORKBASKET`, `M_WORKBASKET` (+ `M_LOGIN_GO.IS_ACTIVE` / `DIVISION_CODE`) | baca | nama di NBStatus dari pemegang aktif workbasket tujuan, divisi IT tidak dihitung; nama kartu kotak masuk Beranda (`GET /api/nb-treaty-in/kotak-masuk`) - keputusan work owner 06-10-2026 |
| `JSON_POLIS` | tulis + baca | tulis: Utility1 `SaveJsonPolisTreatyIn_Act` saat realisasi selesai, **tanpa `DATA_JSON`** (`[keputusan work owner 06-10-2026]` "JSON-nya tidak disimpan, tapi tetap insert kolom lainnya" - RALAT sebagian AC 16; `repository/produksi.go`, pemetaan `models/produksi.go`); baca: pemuat dokumen lama (tiket 22) |
| `TREATYINPRODUCTION` | tulis + baca | tulis: `InsetTreatyInProd_Act` + `InsertTreatyInProd_SQL` (58 kolom) saat realisasi selesai, dilewati bila IDPEGA sudah punya baris; trigger `TRG_TREATYINPRODUCTION_INSERT` menyalin ke `_BACKUP`; baca: cek duplikat `TreatyRealizationCheckDuplicate` |
| `ACHIEVEMENT` | tulis | `SetAchivementValue` + `SaveAchievementSQL` saat realisasi selesai; prosedur `InsertUpdateAchievment` ditulis ulang sebagai satu INSERT (nol prosedur, AC 48), `TGL_PROD` = SYSDATE |
| `DOCUMENT_POLIS`, `CATEGORY_ATTACH_REAS`, `T_STORAGE_IMAGE`, `T_FOLDER_IMAGE` | tulis + baca (`DOCUMENT_POLIS`, `T_STORAGE_IMAGE`); baca (`CATEGORY_ATTACH_REAS`, `T_FOLDER_IMAGE`) | panel lampiran **Attachment File** di bawah layar kasus (grid `AttachmentGridReas` korpus NB FacIn, keputusan WO 08-10-2026) lewat `inti/backend/dokumenpolis` + `inti/backend/penyimpanan` - SQL di inti, bukan di modul ini. Lampiran baru berkunci `KunciInstans` (ID `T_WORK_POLIS` polos), lampiran Pega lama dibaca lewat pzInsKey `ASM-FW-GISFW-WORK <pyID>`. Upload / Delete selama kasus belum Resolve (WO 08-10-2026 *"SEMUA BISA ASAL BELUM RESOLVE"*); rute `/api/nb-treaty-in/kasus/{id}/lampiran` |

## Alat simpanproduksi (06-10-2026)

Menulis baris Utility1 `SaveJsonPolisTreatyIn_Act` (json_polis tanpa `DATA_JSON`, `ACHIEVEMENT`,
`TREATYINPRODUCTION`) untuk SATU kasus yang sudah Resolved-Completed sebelum penulisan produksi ada (permintaan
work owner). Dijalankan **manusia** dari akar repo sesudah `. .\muat-env.ps1`:

```powershell
go run ./modul/nbtreatyin/backend/alat/simpanproduksi -kasus NB-22445             # uji-kering: cetak baris, nol tulis
go run ./modul/nbtreatyin/backend/alat/simpanproduksi -kasus NB-22445 -jalankan   # tulis, satu transaksi
```

Pemetaan sama dengan submit (`models.SusunSimpananPolis`); ProductionDate tersimpan dipakai apa adanya; halaman
kasus tidak diubah. Aman diulang: tabel yang sudah punya baris untuk IDPEGA itu dilewati (ACHIEVEMENT juga).
`-pengguna` mengganti USERNAME (bawaan: OPERATORID riwayat terakhir). `-jalankan` ditolak bila `IS_PEGA_PROD=true`.

## Pemuat dokumen lama (tiket 22)

Memindah seluruh dokumen polis generasi NB (`PRODKE 0`) dari `POOLDATA.JSON_POLIS` ke 8 tabel diagram
grilling - setiap polis, tanpa penyaring (KEPUTUSAN-RONDE-12 butir 5). Dijalankan **manusia** dari akar
repo, sesudah lingkungan dimuat seperti `cmd/api` (`ORACLE_DSN`, `ORACLE_SCHEMA`, `IS_PEGA_PROD`); tidak
pernah berjalan saat aplikasi menyala. Urutan resmi (F7): bab *Migrasi* di atas.

```powershell
go run ./modul/nbtreatyin/backend/alat/pemuatlama -keluaran D:\laporan-pemuat            # uji-kering
go run ./modul/nbtreatyin/backend/alat/pemuatlama -keluaran D:\laporan-pemuat -jalankan  # tulis
```

| Flag | Arti |
| --- | --- |
| `-keluaran <folder>` | wajib; folder berkas laporan, dibuat bila belum ada |
| `-jalankan` | tulis ke tabel baru, satu transaksi per dokumen, aman diulang (kasus ber-IDPEGA sama dilewati). Tanpa flag ini: uji-kering - hanya `JSON_POLIS` yang dibaca, nol pernyataan ke tabel mana pun. Ditolak bila `IS_PEGA_PROD=true` |

Berkas per jalankan (`<stempel>` = `YYYYMMDD-HHMMSS`):

| Berkas | Kolom | Aturan |
| --- | --- | --- |
| `nbtreatyin-arsip-medan-<stempel>.csv` | `POLIS_ID`, `JALUR`, `NILAI`, `KEPUTUSAN` | **arsip audit pemuatan** (F3), bukan penampung: setiap medan dokumen yang tidak ditulis ke kolom maupun ke riwayat produksi, `dibuang: <kunci alasan>` (keputusan tertulis berbukti, `backend/models/medan_abaikan_lama.json`) atau `BELUM DIPUTUSKAN`. Baris `BELUM DIPUTUSKAN` **wajib 0** sebelum pekerjaan dinyatakan selesai (spec-penyimpanan AC 57, 59, RALAT F3) |
| `nbtreatyin-galat-<stempel>.csv` | `IDPEGA`, `NOPOLIS`, `JALUR`, `NILAI`, `SEBAB` | dokumen yang tidak dimuat beserta sebabnya (AC 58); tanggal ambigu tidak ditebak (K15) |

Ringkasan dicetak ke layar (cacah per jenis galat, per alasan medan dibuang, per pola medan belum
diputuskan, baris SuggestList yang disalin beserta cacah PIC kosong (AKSES_LOGIN selalu NULL - baris dokumen
tanpa anggota operator) dan dokumen yang salinannya
dilewati penjaga dobel, nomor kasus `NB-` terbesar yang dimuat - `SEQ_WORK_POLIS` wajib dimajukan
melewatinya). Kode keluar 0 hanya bila nol dokumen gagal dan nol medan belum diputuskan. Penulisan lewat
antarmuka yang sama dengan aplikasi (`SisipKasus`, `SimpanHalaman`, `SetelNomorPolis`, `TutupKasus`,
`CatatUsulan`, ditambah kolom datar json_polis); `JSON_POLIS` hanya dibaca. `SuggestList` dokumen lama
disalin ke `POOLDATA.HISTORYAKSEPTASIPRODUCTION` (F3, pemetaan `SaveViewSuggest` langkah 2) kecuali IDPEGA-nya
sudah punya baris di sana (penjaga dobel `repository.SalinUsulanLama`). Baris hasil pemuat dikenali dari
`IDPEGA` (`<kelas> <pyID>`) dan status `Resolved-Completed` - tanpa kolom penanda (F6). Generasi endorsemen
(`PRODKE > 0`) milik pemuat EDM (`modul/edmtreatyin`, tiket 10) dan hanya dihitung.

### Tombol Copy Old (perintah WO 07-10-2026)

*"nb ttreatyin tobol copy untuk data lama mana?"* -> *"langsung anda kerjakan!"* (dikerjakan sesi EDM TREATY IN dengan
izin WO) - tombol **Copy Old** di samping *Create* di portal (bukan layar Pega), pola sama dengan Copy Old EDM Treaty In
dan Product Name Life. Popup berisi dokumen polis NB lama yang belum ada di tabel flat (NB Number, Master ID, Policy
Number, Insured Name, Group Business, SOB, Ceding, Production Date, Notes); yang dicentang disalin lewat **Process
Copy** dengan `muat` pemuat di atas - satu transaksi per dokumen, hasil per dokumen (Copied / Already in the new tables /
Cannot be copied / Failed). Baris yang tidak dapat disalin tampil dengan alasan (tanpa nilai dokumen) dan tidak dapat
dicentang.

- **Sumber data lama** (perintah WO 07-10-2026): `SELECT * FROM DATAPEGA.PC_ASM_FW_GISFW_WORK a, json_polis b,
  treatyinproduction c WHERE a.pzinskey = b.idpega AND b.idpega = c.idpega` - ditulis EXISTS
  (`repository.sqlKunciJSONPolisCopyOld`) karena TREATYINPRODUCTION berbaris banyak per IDPEGA. DEV 07-10-2026: 27
  dokumen NB (dari 29; 2 belum ada di TREATYINPRODUCTION). Alat pemuat massal di atas tetap TANPA saringan.
- **Hanya superadmin** = pemegang menu Kelola User dengan menu NB Treaty In ber-hak penuh (pola Bordereaux; View only
  berlaku juga). Tombol menurut `GET /hak`; `GET /lama`, `POST /lama/salin` menolak 403.
- **ID kasus salinan = IDPEGA Pega UTUH** (perintah WO 07-10-2026 *"IDPEGA BAWAAN PEGA JANGAN DI POTONG, BERLAKU UNTUK
  SEMUA NB TREATY DAN EDM TREATY"*): `ASM-FW-GISFW-WORK NB-<n>`, bukan `NB-<n>` (`models.IDKasusDariIDPega`; kolom ID
  32 karakter, pzInsKey terpanjang DEV 27). Riwayat SuggestList / HISTORYAKSEPTASIPEGA berkas Pega berkunci pzInsKey yang
  sama, jadi tampil lewat `KunciInstans(ID)`. Berlaku juga untuk alat pemuat CLI.
  - **Layar menampilkan pyID saja** (perintah WO 07-10-2026 *"TAMPILAN NYA HANYA NB-XXX AJA, BERLAKU NB DAN EDM
    TREATY"*): `sajian.idTampil` (= `models.PyIDKasus`) di Offer No portal, Beranda, judul + popup nomor layar kasus,
    popup Copy Old. Kunci buka / kirim / centang tetap ID utuh; saringan portal tetap cocok (LIKE `%NB-<n>%`).
  - Perbaikan data DEV 07-10-2026 (perintah WO *"ID NYA KEPOTONG! CEK SEMUA TABLE YANG PAKE ID ITU! PERBAIKI !"*):
    satu-satunya salinan lama ber-ID `NB-<n>` diganti ke IDPEGA utuh di T_WORK_POLIS, T_GENERAL_POLIS_TREATY,
    T_POLIS_QUOTATION, T_POLIS_CEDING (QUOTATION_ID), T_POLIS_INSTALMENT, T_POLIS_SPREADING - satu transaksi (FK
    tidak deferrable: sisip induk baru dengan NOPOLIS NULL, pindahkan anak, hapus induk lama, pasang NOPOLIS),
    didahului mode cek ber-ROLLBACK. Pindai semua kolom teks tabel T_* sesudahnya: ID terpotong tinggal di
    `ID_NEW_BISNIS` - itu nilai asli dokumen Pega (`IDNewBisnis = pyWorkPage.pyID`), tidak diubah.
- **Cek sumber Copy Old** (WO 08-10-2026 *"CEK APAKAH COPY OLDNYA SUDAH BENERA MENGCOPY SUMBER DATANYA?"*; alat
  baca-saja di atas `PecahDokumenLama` + `BacaHalaman`): DEV 29 dokumen NB Treaty, semua Proportional; NB NonProp
  dan EDM Treaty 0 dokumen. Berkas tersalin cocok dengan sumbernya untuk setiap medan berkolom (spreading
  dibulatkan 8 desimal: DEV `T_POLIS_SPREADING/INSTALMENT(_DETAIL)/XOL(_LAYER)` masih NUMBER(38,8), migrasi
  NUMBER(38,10)). Medan tanpa kolom yang BELUM DIPUTUSKAN diputuskan WO 08-10-2026: `IsSOAUpload` ("IsSOAUpload
  ITU PERLU") jadi kolom T_GENERAL_POLIS_TREATY `IS_SOA_UPLOAD` (320 + `SCRIPT-TABEL-KOLOM-BARU.xlsx` NB TREATY
  bagian E untuk skema lama); `GuaranteeFund` dibuang ("GUARANTEE_FUND NUMBER(38,10), buang!", alasan
  `keputusan_wo` - walau dibaca InsetTreatyInProd_Act); sisanya dibuang dengan bukti korpus di
  `medan_abaikan_lama.json` (`f3_tanpa_pembaca`).
- **Pembuat berkas salinan** (WO 07-10-2026 *"PXCREATEOPERATOR,PXCREATEOPNAME"*): CREATE_OP / CREATE_OP_NAME =
  pembuat kasus Pega (`DATAPEGA.PC_ASM_FW_GISFW_WORK` menurut PZINSKEY = IDPEGA, `repository.PembuatPega`); tanpa
  baris Pega = NULL (tidak dikarang). Berkas tampil di portal akun yang LOGIN_ID-nya = PXCREATEOPERATOR (DEV: 23 dari
  27 dokumen NB kandidat punya akun). Dua berkas NB yang tersalin sebelum aturan ini diisi dari DATAPEGA di DEV
  07-10-2026 (izin WO; UPDATE 2 baris T_WORK_POLIS, mode cek ROLLBACK lebih dulu).
- **Proteksi dobel riwayat** (perintah WO 07-10-2026 *"TAMBAKAN PROTEKSI UNTUK 2 TABLE INI
  historyakseptasiproduction,historyakseptasiPEGA - SAAT COPY, JIKA UDAH ADA PADA 2 TABLE ITU JANGAN DI COPY, SUPAYA
  TIDAK DOUBLE"*): SuggestList dokumen lama ditulis ke HISTORYAKSEPTASIPRODUCTION HANYA bila IDPEGA-nya belum punya
  baris di HISTORYAKSEPTASIPRODUCTION maupun HISTORYAKSEPTASIPEGA (`services.salinUsulanLama` +
  `repository.AdaRiwayatIDPega`; uji seam `TestCopyOldNBTanpaDobelRiwayat`). HISTORYAKSEPTASIPEGA tidak pernah ditulis Copy Old
  / alat pemuat - riwayat Pega berkunci pzInsKey yang sama, dibaca lewat `KunciInstans`. DEV 07-10-2026: dari 27
  dokumen NB kandidat, 2 sudah punya SuggestList dan 1 hanya punya History (SuggestList-nya tidak disalin). Baris
  kembar yang sudah ada di kedua tabel berasal dari data Pega lama, bukan dari Copy Old.
- **Nomor polis dipegang berkas NB lain** (indeks unik `UQ_GP_TREATY_NOPOLIS`): `SetelNomorPolis` menjawab
  `ErrNomorPolisDipakai`, Process Copy menolak dokumen itu dengan *"this policy number is already used by another NB
  case in the new tables"* (dulu tampil "database error").
- **Commencement / Termination berkas salinan**: TREATY_IN_ID kosong (DATA_JSON hanya PolicyTreatyIn), master dibaca lewat
  `NoOffer` = TREATYID view, baris pertama (`DetailKontrakTreaty`, `muatMaster`). Kontrak yang tidak ada di view = tanpa
  tanggal, berkas tetap terbuka.
- Kode: `backend/services/copyold.go` (`GudangPemuat` di `pemuat.go`), `backend/models/copyold.go`,
  `frontend/components/DialogCopyOld.tsx`, `frontend/lama.ts`; uji `backend/handlers/copyold_test.go`,
  `backend/repository/lama_test.go`, `frontend/components/DialogCopyOld.test.tsx`, `frontend/lama.test.ts`.

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
