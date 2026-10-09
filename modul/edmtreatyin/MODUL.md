# Modul `edmtreatyin` — EDM Treaty In

Endorsemen (addendum) polis NB Treaty In yang sudah terbit: setiap endorsemen melahirkan **generasi** polis baru
(PRODKE n+1) yang menunjuk generasi tepat sebelumnya, melalui tangga **Admin → Sec Head → Dept Head**
(`docs/spec.md`, `docs/spec-penyimpanan-relasional.md`). Padanan Pega: portal
`Harness/SFAPortal_Endorsement_Treaty` (tombol *Create* → `Section/TreatyCreateEdm` → `Activity/CreateEDMT`) dan
alur realisasi endorsemen di korpus `EDM Treaty In` (163 rule). Nasib setiap rule (dibangun / tidak + alasan + sumber)
di `docs/INVENTARIS-XML.md`; hasil per acceptance criteria di `docs/HASIL-IMPLEMENTASI.md`; koreksi dokumen lawan XML
di `docs/KOREKSI-DOKUMEN-2026-10-06.md` dan `docs/REKONSILIASI-AC.md`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `edmtreatyin` |
| Folder korpus | `EDM Treaty In` |
| GROUPMENU | `TREATY` |
| Pemilik | `@PEMILIK-EDMTREATYIN` |
| Status | dimigrasi |
| Rentang migrasi | `360-399` |
| Slot menu | `970-971` |
| Prefix rute API | `/api/edm-treaty-in` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | spec, spec-penyimpanan, tiket (`issues/` 00-12), grilling (tidak disunting), `KEADAAN-EDM-TREATY-IN.md`, `STRUKTUR-TABEL-EDM-TREATY-IN.md`, `PERBANDINGAN-KOLOM-DIAGRAM.md`, `INVENTARIS-XML.md`, `HASIL-IMPLEMENTASI.md`, `KOREKSI-DOKUMEN-2026-10-06.md`, `REKONSILIASI-AC.md` |
| `backend/` | `models/` (fungsi murni: halaman kerja, katalog, rantai uang, selisih, pilih bisnis EDM + XOL, tangga, layar) · `repository/` (seluruh SQL) · `services/` (aturan, transaksi) · `handlers/` (HTTP) · `tiruan/` (gudang dalam memori untuk uji seam) · `migrations/` · `alat/pemuatlama/` (pemuat dokumen lama endorsemen — dijalankan manusia) |
| `frontend/` | `pages/` (portal, layar kasus) · `components/` · `api.ts` `labels.ts` `medan.ts` `edmtreatyin.css` dan `*.test.ts` |

Kode yang disalin dari `modul/nbtreatyin` (modul tidak saling mengimpor) membawa komentar asal
`salinan modul/nbtreatyin/... (06-10-2026)`.

## Aturan (ringkas)

- **Generasi, bukan salinan dokumen.** Endorsemen = baris baru `T_GENERAL_POLIS_TREATY` ber-`PRODKE` = generasi
  sebelumnya + 1, `OLD_POLIS_ID` = ID generasi tepat sebelumnya (UNIQUE → nol percabangan), `NOENDORS` = EDMNo,
  `EDM_TYPE`. `NOPOLIS` KOSONG selama berjalan dan diisi nomor polis induk saat Dept Head menyetujui
  (`SetelNomorPolisSelesai`; indeks fungsi `UQ_GP_TREATY_NOPOLIS` hanya mencakup baris bernomor). Daftar NB
  menyaring `PRODKE = 0`, daftar EDM `PRODKE >= 1` dan `ID LIKE 'EDMT-%'`.
- **ID kasus** `EDMT-<SEQ_WORK_POLIS>` (`CreateEDMT`); **EDMNo** = nomor polis + `/E` + sekurangnya dua digit
  (`SetEDMTNoPolis` langkah 3). Satu endorsemen berjalan per polis (409 bila ada).
- **OldData** dibaca dari generasi yang ditunjuk `OLD_POLIS_ID` (bukan `JSON_POLIS.DATA_JSON`); tab *Old Data 2*
  memperlihatkan selisih generasi sebelumnya.
- **Selisih** (`TreatyDifference`) dihitung Go di transaksi yang sama dan disimpan ke tabel proyeksi 360-363
  (`SUMBER = 'GO'`; baris `'PEGA'` hasil migrasi BEKU). Proporsional: satu rumus *baru − lama* untuk semua generasi
  (varian berlapis Pega TIDAK ditiru — keputusan WO 23-09-2026, ID-28/30). NonProp: `CalculateDifferenceEDM_act`
  apa adanya (batas bawah 0, prorata, pajak, pembatalan). `Total*` diturunkan, tidak disimpan.
- **Tangga** selalu tiga tingkat (prompt eksekusi WO): `CekLimitTreatyAcc_Act`/LetterNo tidak dibangun; operator
  hardcode Pega diganti posisi/workbasket. Admin menolak → kasus ditolak dan generasinya dilepas (`OLD_POLIS_ID`
  NULL); Sec Head / Dept Head menolak → kembali ke Admin (Position `4` / `5`).
- **NBStatus**: `EDM IS IN <akun>'S INBOX` (create), `EDMT IS IN <NAMA>'S INBOX`, `EDM WAS DECLINED BY  <NAMA>`,
  kembali ke Admin `EDM IS IN <PEMBUAT>'S INBOX`.
- **Master treaty** dibaca JSON BACA-SAJA (pengecualian K8) lewat SATU metode `repository.pembacaMasterEDM.baca`,
  keluarannya selalu lewat penyaring daftar medan tertutup (`uraiMasterXOLPerBaris`) — dijaga
  `repository/masterxol_penjaga_test.go`. Nol penulisan JSON.
- **Catatan usulan** (`SuggestList`) ke `HISTORYAKSEPTASIPRODUCTION` `TYPE_POLIS 'EDMT'` (pola NB K4 — penyimpangan
  sadar: korpus EDM menyimpannya di blob kasus).
- **Konversi Arasapas** (`serviceInsertArasapas_act`) TIDAK disambung (sama dengan NB; keputusan WO 07-10-2026); muatan
  `noPolis / caseId / tglInput` disiapkan `models.RakitMuatanKonversi`.
- **Keputusan WO 07-10-2026** (jawaban laporan 06-10, rinci di `docs/HASIL-IMPLEMENTASI.md` bab 5):
  - **Kotak masuk Beranda**: EDM ikut (`frontend/menu.ts` `antreanBeranda` + `daftarBeranda`) - Sec Head / Dept Head
    membuka berkas dari Beranda.
  - **Tombol "Calculate Value Difference" DIBUANG** (WO 08-10-2026 *"COBA CEK TOMBOLITU, APAKAH MASIH
    DIPERLUKAN? KALAU SUDAH TIDAK DIHAPUS AJA!"*; penyimpangan sadar dari `PropNewData2` S24): sel uang tab
    New Data sudah menghitung ulang tab Value Difference; % spreading dan `.Installment` kini ikut
    (`medan.ts` `AKSI_SPREADING`, aksi Installment); server tetap menghitung ulang saat Save / Submit / produksi
    (`models.HitungSelisihGenerasi`).
  - **Portal = aturan portal NB**: switch *In Progress* (buatan akun, masih proses) / *Resolved* (`?status=selesai`,
    berkas selesai BUATAN akun ini, dibuka hanya-baca - RALAT WO 07-10-2026 *"TAMBAHKAN KAN UNTUK PEMBUAT. MENU ITU
    HANYA UNTUK SI PEMBUAT, NB DAN EDM TREATY"*; dulu semua berkas selesai).
    Tab Resolved menambah kolom **Production Date** DD-MM-YYYY (`T_GENERAL_POLIS_TREATY.TGL_PROD`) di ujung; Policy Number
    sudah kolom XML ke-3.
    Kolom Status berkas selesai = `T_WORK_POLIS.STATUS_WORK` (`Resolved-Completed` / `Resolved-Rejected`), bukan NBStatus
    terakhir (*"KALO DAH RESOLVE STATUS NYA PAKE STATUS RESOLVE"*, pola portal NB).
  - **Header terkunci bagi atasan**: With Tax, Type Tax, Overiding Commision, Marketing Officer hanya diterima dari layar
    Admin (`models.GabungMasukanLayar`, `services/aksiposisi.go`) - XML tidak menguncinya per posisi.
  - **Label EDMType** = DT `TreatyEDMListType` (screenshot Pega): 1 Internal · 2 External · 3 Adjustment Premium · 4 Cancel
    Input (`models.LabelJenisEDM`, `frontend/labels.ts` `LABEL_JENIS_EDM`).
  - **Pembatalan** ikut XML (16 medan `SetEDMTCancel`); kode mati salinan NB dibuang.
  - **Grid spreading tab New Data hanya-baca** (sama dengan NB): tanpa Add / Delete, Type Treaty teks; server hanya
    menerima % Share dan % Share klaim. Choose Business (`FillSpreading`) membawa SEMUA baris generasi lama (XML hanya
    baris 1 - sisanya dulu ditambah lewat Add). Tombol Save di dalam tab dibuang; Save kaki halaman tetap.
  - **With Tax / Type Tax diwarisi** dari generasi lama bila belum diisi (saat Create dan Choose Business,
    `models.WarisPajakLama`) - grid XOL Current Premium NonProp ikut berpajak; "false" yang dipilih admin dihormati.
  - **With Tax / Type Tax diubah SESUDAH Choose Business** (NonProp baru): server menjalankan lagi rantai Choose Business
    dengan master yang sama (`models.HitungUlangPajakNonPropEDM`, pola NB `pajakNonProp`) - pajak grid XOL dan
    selisihnya terhitung; Installment, StatementDate, Marketing Officer dikembalikan, angsuran disusun ulang.
  - **Class Of Business diganti Source Of Business** (*"SAMAIN DENGAN NB NYA!"*): sel kiri `.SOBName` tanpa syarat,
    sel `.SOBName` kolom kanan dipindah (tidak dobel); `.BizName` tidak tampil lagi.
  - **Grid XOL NonProp** (Previous / Current Premium / Total Difference + rincian layer): uang 4 desimal, kosong / nol
    tampil "0" rata kanan (*"JANGAN NULL TAPI 0"*, *"BUAT 4 ANGKA BELAKANG KOMA"*; `frontend/xol.ts` `UANG4`) - hanya
    tampilan, nilai tersimpan tidak diubah.
  - **Baris spreading tambahan tanpa % Share** ditolak dengan pesan di baris itu (`PesanShareSpreadingKosong`), bukan
    dibagi `PolicyTreatyIn.RNMShare` (tak diisi rule mana pun).
  - **Kotak saring portal diperluas** (*"pencarian ... buat bisa mencari nomor nb/edm, insured name dll"*): setiap kata
    (paling banyak 5) wajib cocok dengan salah satu kolom - nomor kasus EDMT, Offer/Master ID, nomor polis (generasi dan
    polis NB), EDM No, insured (generasi dan quotation), group business, SOB, ceding, marketing, treaty group, class of
    business, nama pembuat; tanpa beda huruf, `%` `_` harfiah (`repository.kolomCariPortal`, `models.KataCari`). XML
    hanya `OldPolicyNo`. Seragam dengan portal NB (sesi NB TREATY). Dicoba DEV 07-10-2026 (COUNT saja).
  - Bila konversi Arasapas kelak disambung: P50 `IsFacRetro` dan gerbang `IsPEGAPROD` / `IsTreatyIn` mengikuti XML.

## Migrasi

Rentang `360-399`. 360-363 — TEPAT empat tabel **proyeksi selisih** (diagram grilling sheet *EDM Treaty In Prop* /
*NonProp*, golongan UNGU *proyeksi*); nol tabel dasar dibuat, nol kolom ditambah ke tabel NB. Perbandingan kolom
lawan diagram dan XML: `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`; bentuk kolom: `docs/STRUKTUR-TABEL-EDM-TREATY-IN.md`.
Slot menu `970`: satu `UPDATE DIMIGRASI = '1'` baris `edmtreatyin` (sudah ada sejak inti 900), nol `INSERT`.
Seluruhnya dikunci `backend/migrasi_test.go`.

| Migrasi | Isi | Diagram |
| --- | --- | --- |
| 360 | `T_POLIS_DIFFERENCE` (1:1 generasi, `SUMBER` PEGA/GO) | Prop J85–J101 · NonProp J100–J116 |
| 361 | `T_POLIS_DIFFERENCE_SPREADING` | Prop R103–R114 · NonProp R118–R129 |
| 362 | `T_POLIS_DIFFERENCE_INSTALMENT` (satu tingkat) | Prop R116–R122 · NonProp R131–R137 |
| 363 | `T_POLIS_XOL_LAYER_DIFFERENCE` (NonProp saja) | NonProp R139–R147 |
| 970 | `UPDATE M_NAV_MENU SET DIMIGRASI = '1' WHERE KODE = 'edmtreatyin'` | — |

**Urutan:** sesudah NB 320-328 (FK `T_POLIS_DIFFERENCE.POLIS_ID` → `T_GENERAL_POLIS_TREATY`). Agen **tidak**
menjalankan migrasi: `-migrate` oleh work owner (menolak `IS_PEGA_PROD=true`), lalu mulai ulang backend. Skrip SQL
Developer setara (bagian A, pembatalan bagian B): `OUTPUT_HASIL_RNM/SCRIPT-TABEL-KOLOM-BARU.xlsx` sheet
**EDM TREATY**. Status DEV 07-10-2026: `-migrate` WO sempat berhenti di pra-terbang 880 masterprovince (efek 880-883
sudah ada di DEV tetapi tak tercatat); atas izin WO keempatnya dicatat di `T_MIGRASI` (4 baris, tabel tidak diubah) - `-migrate`
berikutnya membuat 360-363, menjalankan 970, dan mencatat NB 328. Jalur manual bagian A tetap aman: tabel berbentuk sama
hanya dicatat (ORA-00955 ditelan, `inti/backend/migrasi`).

## Tabel yang dibaca dan ditulis, tidak dibuat

| Tabel | Akses | Dasar |
| --- | --- | --- |
| `T_GENERAL_POLIS_TREATY` + `T_POLIS_QUOTATION`, `T_POLIS_CEDING`, `T_POLIS_INSTALMENT`, `T_POLIS_INSTALMENT_DETAIL`, `T_POLIS_SPREADING`, `T_POLIS_XOL`, `T_POLIS_XOL_LAYER`, `T_POLIS_SURVEY` | tulis + baca (generasi `PRODKE >= 1` saja) | milik `nbtreatyin` (320-328); spec-penyimpanan ID-1..ID-5 |
| `T_WORK_POLIS`, `SEQ_WORK_POLIS` | tulis + baca | tabel kasus lintas-lini (premiumlistlife); ID `EDMT-n` |
| `HISTORYAKSEPTASIPEGA` | tulis + baca | `InsertHistoryAkseptasiPega` pasca-submit |
| `HISTORYAKSEPTASIPRODUCTION` | tulis + baca | `SuggestList`, `TYPE_POLIS 'EDMT'` (pola NB K4) |
| `JSON_POLIS` | tulis + baca | baca: `CheckNopolisAvailability` / `FetchNopolisCount` (ada-tidaknya nomor polis, generasi terakhir); tulis: Utility1 `SaveJsonPolisTreatyInEDM_Act` **tanpa `DATA_JSON`** (keputusan WO 06-10-2026) |
| `TREATYINPRODUCTION` | tulis + baca | tulis: `InsertTreatyInProdEDMT_SQL` saat selesai, dilewati bila IDPEGA sudah ada; baca: `FetchNoOfferFromNoPolis` (*No Master Treaty* layar Create, `TrtEdmCheckPolicyError` 6-7) |
| `ACHIEVEMENT` | tulis | `SaveAchievementSQL` (prosedur `InsertUpdateAchievment` = satu INSERT, nol prosedur) |
| `M_TREATY_IN`, `M_TREATY_IN_EDM`, `M_TREATY_OUT` | baca (JSON, satu metode) | `BrowseTreatyIn*`, `BrowseTreatyOutDetailEDM`; K8 |
| `TREATY_IN`, `TREATY_IN_EDM`, `TREATY_OUT2` | baca | popup *Choose Business* EDM (tiga varian RDB) |
| `CURRENCY`, `MARKETINGOFFICER`, `REINSURANCETYPE`, `BUSINESS`, `AGENT`, `TANGGAL_CLOSING` | baca | RDB terjangkau |
| `M_LOGIN_GO`, `M_LOGIN_GO_WORKBASKET`, `M_WORKBASKET` | baca | nama tampilan dan pemegang kotak masuk NBStatus (divisi IT tidak dihitung) |
| `DOCUMENT_POLIS`, `CATEGORY_ATTACH_REAS`, `T_STORAGE_IMAGE`, `T_FOLDER_IMAGE` | tulis + baca (`DOCUMENT_POLIS`, `T_STORAGE_IMAGE`); baca (`CATEGORY_ATTACH_REAS`, `T_FOLDER_IMAGE`) | panel lampiran **Attachment File** di bawah layar kasus (grid `AttachmentGridReas` korpus NB FacIn, keputusan WO 08-10-2026) lewat `inti/backend/dokumenpolis` + `inti/backend/penyimpanan` - SQL di inti, bukan di modul ini. Lampiran baru berkunci `KunciInstans` (ID `T_WORK_POLIS` polos), lampiran Pega lama dibaca lewat pzInsKey `ASM-FW-GISFW-WORK <pyID>`. Upload / Delete selama kasus belum Resolve (WO 08-10-2026 *"SEMUA BISA ASAL BELUM RESOLVE"*); rute `/api/edm-treaty-in/kasus/{id}/lampiran` |

## API (`/api/edm-treaty-in`)

| Rute | Guna |
| --- | --- |
| `GET /kasus` | portal endorsemen (saringan XML; 500 baris, urut `TGL_UPDATE`); `?status=selesai` = switch Resolved |
| `GET /kotak-masuk`, `GET /kotak-masuk/kasus` | kartu kotak masuk Beranda per workbasket dan daftarnya |
| `GET /periksa-polis` | `CheckNopolisAvailability` + `TrtEdmCheckPolicyError` sebelum create |
| `POST /kasus` | `CreateEDMT` — melahirkan generasi berikut (satu transaksi) |
| `GET /kasus/{id}`, `PUT /kasus/{id}` | buka / simpan draf |
| `POST /kasus/{id}/hitung` | aksi hitung bernama (selisih, angsuran EDMT, pajak, spreading, ...) |
| `POST /kasus/{id}/bisnis`, `POST /kasus/{id}/pilih-bisnis` | popup *Choose Business* EDM dan `EDMChooseBusiness_Act` |
| `POST /kasus/{id}/kirim` | submit per posisi (Utility1 saat Dept Head setuju) |
| `GET /acuan` | daftar pilihan (mata uang, MO, jenis EDM, ...) |
| `GET /hak` | hak layar portal akun (`{"copyOld": bool}`) |
| `GET /kasus/{id}/lampiran`, `GET /kasus/{id}/lampiran/dokumen?kategori=`, `POST` (sama), `GET /kasus/{id}/lampiran/dokumen/{did}/isi`, `GET .../office`, `POST .../hapus` | lampiran **Attachment File** (`inti/backend/dokumenpolis.Pasang`): grid kategori + `bolehUbah`, daftar / unggah, unduh / View, View Office Online, Delete |
| `GET /lama`, `POST /lama/salin` | popup **Copy Old** dan *Process Copy* `{"ids": [...]}` — superadmin saja |

## Pemuat dokumen lama (tiket 09-10)

Memindah **seluruh generasi endorsemen** lama (`POOLDATA.JSON_POLIS` ber-`PRODKE > 0`, `pxObjClass` PolicyTreatyIn,
IDPEGA kelas Work EndorsementTreaty `EDMT-n` — bukti `SaveJsonPolisTreatyInEDM_Act` / `SavePolisTreatyInEDM_SQL`)
ke tabel generasi + proyeksi selisih **`SUMBER = 'PEGA'`**, lewat antarmuka penyimpanan yang SAMA dengan aplikasi
(AC 44): `SisipKasus` (UNIQUE `OLD_POLIS_ID` menolak percabangan) → `SimpanHalaman` → `SetelNomorPolisSelesai` →
`SimpanSelisih` → `CatatUsulan` → `TutupKasus`, satu transaksi per dokumen. Nilai selisih disalin **apa adanya**
(AC 39 — tidak dihitung ulang, digit galat utuh); penanda migrasi `PASANGAN_BERGESER` / `RUMUS_BERLAPIS`
(ID-35..ID-37, AC 40-43) dihitung tanpa menyentuh satu angka dan hanya pada baris `SUMBER = 'PEGA'`.

```powershell
go run ./modul/edmtreatyin/backend/alat/pemuatlama -keluaran <folder>             # uji-kering (bawaan, baca saja)
go run ./modul/edmtreatyin/backend/alat/pemuatlama -keluaran <folder> -jalankan   # tulis
```

- Dijalankan **manusia** (WO/DBA), tidak pernah agen; `-jalankan` ditolak bila `IS_PEGA_PROD=true`.
- **Urutan wajib:** migrasi 360-363 → pemuat NB (`modul/nbtreatyin/backend/alat/pemuatlama -jalankan`, generasi
  `PRODKE 0`) → pemuat EDM uji-kering → keputusan medan `BELUM DIPUTUSKAN` (`backend/models/medan_abaikan_lama.json`,
  pola F3 NB) → pemuat EDM `-jalankan`. Generasi diurut NOPOLIS lalu PRODKE **sebagai bilangan** (di Go; kolom
  `JSON_POLIS.PRODKE` bertipe teks). Generasi sebelumnya yang belum ada = galat dokumen, bukan tebakan.
- Aman diulang: generasi yang sudah dimuat dilewati; proyeksi `PEGA` beku (`ErrSelisihBeku`, juga atas `PEGA`).
- Keluaran: `edmtreatyin-arsip-medan-<stempel>.csv` dan `edmtreatyin-galat-<stempel>.csv`; kode keluar 0 hanya bila
  nol dokumen gagal dan nol medan `BELUM DIPUTUSKAN`.
- Status 06-10-2026: dibangun dan diuji di atas gudang tiruan (rantai tiga generasi, percabangan ditolak, nilai AC 39
  identik digit demi digit); **belum pernah dijalankan** di DEV.

### Tombol Copy Old (perintah WO 07-10-2026)

*"BUATKAN TOMBOL COPY OLD SAMA SEPERTI MASTER PRODUCTNAME LIFE, KHUSUS BUAT SUPERUSER"* — tombol **Copy Old** di samping
*Create New Addendum Treaty* di portal (bukan layar Pega). Popup berisi dokumen endorsemen lama di `JSON_POLIS` yang
belum ada di tabel flat (EDM Number, Policy Number, EDM No, Generation, EDM Type, SOB, Ceding, Production Date, Notes);
yang dicentang disalin lewat **Process Copy** dengan `muat` pemuat di atas — satu transaksi per dokumen, urutan
generasi (bukan urutan centang), dokumen yang gagal tidak membatalkan yang lain, hasil per dokumen (Copied / Already in
the new tables / Cannot be copied / Failed). Baris yang tidak dapat disalin (generasi sebelumnya belum ada,
percabangan, ID bentrok, galat dokumen) tampil dengan alasan dan tidak dapat dicentang; pesan tanpa nomor polis / nilai
dokumen (`models.AlasanSalinLama`).

- **Sumber data lama** (perintah WO 07-10-2026): `SELECT * FROM DATAPEGA.PC_ASM_FW_GISFW_WORK a, json_polis b,
  treatyinproduction c WHERE a.pzinskey = b.idpega AND b.idpega = c.idpega` - ditulis EXISTS
  (`repository.sqlPmKunciJSONPolisEDMCopyOld`), sama dengan Copy Old NB. Alat pemuat massal tetap tanpa saringan.
- **ID kasus salinan = IDPEGA Pega UTUH** (perintah WO 07-10-2026 *"IDPEGA BAWAAN PEGA JANGAN DI POTONG"*):
  `ASM-FW-GISFW-WORK EDMT-<n>`. RALAT: pzInsKey Pega memakai kelas GRUP `ASM-FW-GISFW-WORK` (DEV: 4 kasus EDMT, 264 NB),
  bukan `KelasKerjaEDM` - dulu setiap dokumen EDM asli akan ditolak `ErrIDPega`. Daftar portal / kotak masuk / cek EDM
  terbuka mengenali kedua bentuk ID (`repository.sqlIDKasusEDM`). Berlaku juga untuk alat pemuat CLI.
  - **Layar menampilkan pyID saja** (perintah WO 07-10-2026 *"TAMPILAN NYA HANYA NB-XXX AJA, BERLAKU NB DAN EDM
    TREATY"*): `sajian.idTampil` (= `models.PyIDKasus`) di EDM Number portal, Beranda, judul + popup nomor layar kasus,
    popup Copy Old. Kunci buka / kirim / centang tetap ID utuh; saringan portal tetap cocok (LIKE `%EDMT-<n>%`).
- **Cek sumber Copy Old** (WO 08-10-2026 *"CEK APAKAH COPY OLDNYA SUDAH BENERA MENGCOPY SUMBER DATANYA?"*; alat
  baca-saja di atas `PecahDokumenLama` + `BacaHalaman`): DEV 29 dokumen NB Treaty, semua Proportional; NB NonProp
  dan EDM Treaty 0 dokumen. Berkas tersalin cocok dengan sumbernya untuk setiap medan berkolom (spreading
  dibulatkan 8 desimal: DEV `T_POLIS_SPREADING/INSTALMENT(_DETAIL)/XOL(_LAYER)` masih NUMBER(38,8), migrasi
  NUMBER(38,10)). Medan tanpa kolom yang BELUM DIPUTUSKAN diputuskan WO 08-10-2026: `IsSOAUpload` ("IsSOAUpload
  ITU PERLU") jadi kolom T_GENERAL_POLIS_TREATY `IS_SOA_UPLOAD` (320 + `SCRIPT-TABEL-KOLOM-BARU.xlsx` NB TREATY
  bagian E untuk skema lama); `GuaranteeFund` dibuang ("GUARANTEE_FUND NUMBER(38,10), buang!", alasan
  `keputusan_wo`; produksi EDM membaca TreatyDifference.GuaranteeFund); sisanya dibuang dengan bukti korpus di
  `medan_abaikan_lama.json` (`f3_tanpa_pembaca`). EDM: `IS_SOA_UPLOAD` ikut katalog EDM (tabel
  NB yang sama); `RNMShare` TETAP belum diputuskan di EDM - dibaca `CountSpreading_Act` korpus EDM.
- **Pembuat berkas salinan** (WO 07-10-2026 *"PXCREATEOPERATOR,PXCREATEOPNAME"*): CREATE_OP / CREATE_OP_NAME =
  pembuat kasus Pega (`DATAPEGA.PC_ASM_FW_GISFW_WORK` menurut PZINSKEY = IDPEGA, `repository.PembuatPega`); tanpa
  baris Pega = NULL (tidak dikarang). Berkas tampil di portal akun yang LOGIN_ID-nya = PXCREATEOPERATOR.
- **Proteksi dobel riwayat** (perintah WO 07-10-2026 *"TAMBAKAN PROTEKSI UNTUK 2 TABLE INI
  historyakseptasiproduction,historyakseptasiPEGA - SAAT COPY, JIKA UDAH ADA PADA 2 TABLE ITU JANGAN DI COPY, SUPAYA
  TIDAK DOUBLE"*): SuggestList dokumen lama ditulis ke HISTORYAKSEPTASIPRODUCTION HANYA bila IDPEGA-nya belum punya
  baris di HISTORYAKSEPTASIPRODUCTION maupun HISTORYAKSEPTASIPEGA (`services.salinUsulanLama` +
  `repository.AdaRiwayatIDPega`; uji seam `TestCopyOldTanpaDobelRiwayat`). HISTORYAKSEPTASIPEGA tidak pernah ditulis Copy Old
  / alat pemuat - riwayat Pega berkunci pzInsKey yang sama, dibaca lewat `KunciInstans`. DEV 07-10-2026: dari 27
  dokumen NB kandidat, 2 sudah punya SuggestList dan 1 hanya punya History (SuggestList-nya tidak disalin). Baris
  kembar yang sudah ada di kedua tabel berasal dari data Pega lama, bukan dari Copy Old.
- **Hanya superadmin** = pemegang menu Kelola User dengan menu EDM Treaty In ber-hak penuh (pola Copy Old Data
  Bordereaux; View only berlaku juga bagi superadmin). Tombol tampil menurut `GET /hak`; rute `lama` menolak 403.
- Generasi 1 butuh polis NB-nya sudah di tabel flat (pemuat NB). DEV 07-10-2026: nol dokumen endorsemen Treaty In lama
  di `JSON_POLIS`, jadi popup kosong di DEV.
- Baris `JSON_POLIS` tulisan Utility1 aplikasi baru (IDPEGA = ID kasus polos tanpa kelas Pega, `DATA_JSON` kosong) BUKAN
  dokumen Pega: dilewati popup dan pemuat (`models.ErrBarisAplikasiBaru`, dihitung terpisah di ringkasan). Ditemukan
  WO 07-10-2026: popup menampilkan baris EDMT-22449 sendiri sebagai "the old JSON cannot be read".
- Kode: `backend/services/copyold.go`, `backend/models/copyold.go`, `frontend/components/DialogCopyOld.tsx`,
  `frontend/lama.ts`; uji `backend/services/copyold_test.go` (daftar, urutan generasi, superadmin, HTTP, View only).

## Menjalankan uji modul ini saja

Dari akar repo:

```powershell
go test ./modul/edmtreatyin/...
npx vitest run modul/edmtreatyin
```

Uji seam HTTP (`backend/handlers/alur_test.go`) berjalan di atas `backend/tiruan` (gudang dalam memori, transaksi
bersnapshot, penjaga UNIQUE `OLD_POLIS_ID`) — tanpa Oracle. SQL baru dicoba baca-saja di DEV (hitungan saja)
sebelum dinyatakan selesai; SQL tulis ke tabel 360-363 baru dapat dicoba sesudah migrasinya dijalankan WO.
