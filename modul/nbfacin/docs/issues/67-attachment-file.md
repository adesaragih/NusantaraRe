# 67: Attachment File (lampiran) case NB FacIn

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 08-10-2026 (gambar layar grid
> Attachment File Pega) "dipega ada add attachment, kenapa disini blm ada ya, tambahin dong". Riset: scratchpad sesi 0f
> `riset-F-attachment.md`.

**Status:** frontend selesai 08-10-2026 (uji hijau); backend selesai sesi c3 08-10-2026 (lihat "Backend sesi c3"); pelaksana NYATA menunggu persetujuan penjaga klien HTTP (lihat "Terbuka").

## Keputusan work owner 08-10-2026

- **K67-1** Konsep penyimpanan sama dengan Treaty Contract Out (Google Storage, ADR-0010), tetapi NB FacIn menyimpan ke tabel Pega
  **`DOCUMENT_POLIS`** dan **`T_STORAGE_IMAGE`** ("konsepnya sama, tapi untuk nb facin disimpan ke tabel DOCUMENT_POLIS,
  T_STORAGE_IMAGE"); DDL `D:\migrasi\RNM\DDL\DOCUMENT_POLIS.txt` / `T_STORAGE_IMAGE.txt` dari work owner.
- **K67-2** ~~Upload di semua tahap oleh pemegang workbasket posisi case;~~ (bagian upload DIGANTI K67-6) Delete hanya tahap Offer
  (FlagOnGoingPolicy 0) atau IT.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\AttachmentGridReas.xml` (`ASM-FW-GISFW-DATA-FACOFFER!ATTACHMENTGRIDREAS`): "Attachment File", Category · Count · Upload
  File · View File; kategori `CategoryAttach_SQL` = `select * from CATEGORY_ATTACH_REAS order by note`; Count = dokumen
  `Int-DOCUMENT_POLIS` IDPEGA + KATEGORI_2. Di Pega grid ini di popup Preview R/I Slip (`PrintRISlip(s)`).
- Upload: `AttachContentGIS` → `AttachFileGIS` → `InsertDocument_Act` (KATEGORI_1 "Reas", KATEGORI_2 = kategori, ID cap waktu) →
  `InsertGoogleStorage_Act` (`GetMimeType`, Folder "Policy", token `GET_TOKEN_STORAGE` — sudah port Go di inti, ADR-0043).
- View: harness `ViewPictureList` (File · Note · Upload Date · Delete); unduh `DownloadDocumentPolis` → `GetUrlGoogleStorage_Act`.
  Delete Pega = hapus fisik → di sini JUGA hapus fisik: work owner 08-10-2026 "ikutin pega aja" (pengecualian sadar ADR-0031;
  rancangan awal "hapus LOGIS" diralat).
- `CheckSpreadingProtect_ACT` langkah 10: Count APPROVAL SPREADING ≥ 1 membebaskan "Can not proceed spreading ORS".

- ⚠️ **RALAT acuan 08-10-2026** (work owner: "AttachmentGridReas km bisa cek xml itu untuk referensi, itu dipakai dipega"):
  grid layar case Pega = `D:\migrasi\RNM\DDL\AttachmentGridReas.xml` **`ASM-FW-GISFW-WORK!ATTACHMENTGRIDREAS`** (pemuat
  `Activity\LoadAttachmentData.xml`), BUKAN varian `DATA-FACOFFER` (popup Preview R/I Slip, Position "FACOUT") yang dipakai riset-F.
  Kategori: `DDL\CATEGORY_ATTACH_REAS.csv` (20 baris, kolom NOTE · BISNIS · POSITION; FACIN 16, FACOUT 3, ALL 1) — gambar layar memuat
  ke-20, `[dugaan]` tanpa saringan POSITION (diverifikasi sesi c3 dari LoadAttachmentData).

## Keputusan agent

- **K67-3** Blok di layar case di bawah riwayat akseptasi (permintaan work owner), bukan popup Preview R/I Slip.
- **K67-4** Unggah banyak berkas = satu permintaan multipart per berkas (pola TCO); konfirmasi sebelum Delete (Pega tanpa konfirmasi).

## Frontend (sesi 0f)

- `components/BlokLampiran.tsx`, `api.ts` (`ambilLampiran`, `daftarLampiran`, `unggahLampiran`, `jalurIsiLampiran`,
  `hapusLampiran`), `labels.ts` `LAMPIRAN` / `TEKS_LAMPIRAN`, `pages/InwardFacultative.tsx`, `nbfacin.css` `.nbf-lampiran*`.
  Uji `BlokLampiran.test.tsx`.

## Backend (diminta ke sesi c3, 08-10-2026)

Lihat pesan kontrak tiket 67.

## Keputusan backend 08-10-2026 (dikumpulkan sesi 9d; work owner lewat AskUserQuestion kecuali ditandai)

- **K-A** `[work owner]` "ikutin pega aja, kalau km perlu xml nya infoin saja": Delete = HAPUS FISIK seperti Pega (objek +
  T_STORAGE_IMAGE + DOCUMENT_POLIS), tanpa migrasi 214 - PENGECUALIAN SADAR ADR-0031 (bagian "Pengecualian" ADR-0031).
- **K-B** `[work owner]` "Salin ke NB FacIn": mesin penyimpan di modul nbfacin dengan parameter Pega NB; inti dan modul
  Treaty Contract Out tidak disentuh; config inti yang sudah ada dibaca (PELAKSANA_STORAGE, UNGGAHAN_DIR,
  STORAGE_TOKEN_SALT).
- **K-C** `[keputusan agent, ADR-0008 mengikat]` asinkron pola Treaty Contract Out; `status` '' / 'menunggu' / 'gagal' +
  `POST …/ulangi`.
- **K-D** `[keputusan agent]` ID Pega verbatim `yyyyMMddhhmmssSSS` WIB (`hh` jam 12-an) + geser 1 ms bila terpakai / ORA-00001.
- **K-E** `[keputusan agent]` 400 "Jenis berkas tidak didukung: .<ext>"; MIME = ekstensi huruf kecil; objek Google sesuai Pega.
- **K-F** `[work owner]` "Persis Pega": bolehHapus = FlagOnGoingPolicy "0" ATAU akun IT, siapa pun yang membuka case.

## Acuan grid (9d, 08-10-2026): varian kelas WORK `[terverifikasi]`

`D:\migrasi\RNM\DDL\AttachmentGridReas.xml` = `ASM-FW-GISFW-WORK!ATTACHMENTGRIDREAS` (pemuat
`LoadAttachmentData`). Dibandingkan dengan riset-F (varian `DATA-FACOFFER`): **rantai sama** -
`NB FacIn\Activity\LoadAttachmentData.xml` (`WORK-!LOADATTACHMENTDATA`) hanya langkah 5 `Call
ASM-FW-GISFW-Data-OfferFacIn.SetCategoryAttach` tanpa parameter (langkah 1-4 `//`) -> kategori `CategoryAttach_SQL` `select *
from CATEGORY_ATTACH_REAS order by note` (BusinessFac F/T) **tanpa saringan POSITION / BISNIS** (kolom `POSITION` CSV tidak
dibaca jalur ini; 20 baris = gambar layar), Count `InputParamUploadReas_act`; sel Upload b2400-2520 `SetCategoryAttachment_Reas`
(category "Reas", inskey pzInsKey, GCNMCategory .NOTE) -> `AttachContentGIS`; sel View b3117-3213 `InputParamUploadReas_act`
-> harness `ViewPictureList` popup 1600×800; nol syarat tampil / baca-saja per sel (grid `readOnly` baris saja). Delete di
`ViewPictureList` b4605 `pyWorkPage.FlagOnGoingPolicy = '0' || OperatorID.pyOrgDivision = 'IT'`. Kontrak api.ts TIDAK berubah.
`[pertanyaan terbuka]` section induk yang menyertakan varian Work di layar case tidak ada di korpus / DDL (grep `AttachmentGridReas`
korpus: hanya PrintRISlip(s) / *_Endorsement varian DATA-FACOFFER) - syarat tampil blok di layar Pega belum terverifikasi.

## Backend sesi c3 (08-10-2026)

- **K67-B1 Rute** (`handlers/lampiran.go`): `GET …/lampiran` {kategori: [{kategori, jumlah}], bolehUnggah, bolehHapus};
  `GET …/lampiran/berkas?kategori=` {baris: [{id, namaFile, catatan, tanggal, status}]}; `POST …/lampiran/berkas` multipart
  (berkas, kategori) 201 {id}; `GET …/berkas/{id}/isi` (proksi, attachment + nosniff); `POST …/berkas/{id}/ulangi` 204;
  `DELETE …/berkas/{id}` 204. 401 tanpa sesi; 404 case / dokumen (dibatasi IDPEGA case); 403 unggah / ulangi bukan pemegang
  workbasket POSITION, hapus bukan (flag 0 / IT); 400 masukan; 413 > 25 MiB; 409 belum / sudah tersimpan / sumber hilang;
  422 kategori di luar master; 503 tanpa Oracle / UNGGAHAN_DIR.
- **K67-B2 Status** kabel: `""` tersimpan (baris T_STORAGE_IMAGE ada), `"menunggu"` (efek unggah terakhir antre / jalan),
  `"gagal"` (efek terakhir GAGAL_PERMANEN, atau tanpa objek dan tanpa efek - dokumen yang tidak pernah diantre sistem ini).
- **K67-B3 Repository** (`repository/lampiran.go`): SQL persis rule (`CategoryAttach_SQL` tanpa DISTINCT seperti Pega;
  Count = GROUP BY KATEGORI_2 per IDPEGA, KATEGORI_1 tidak disaring; GetAllDocument maks 500, ORDER BY ID `[keputusan agent]`
  karena RD tanpa urutan; INSERT DOCUMENT_POLIS TANGGAL SYSDATE, INSKEY_* NULL; Insert/Update/GetLink/Delete T_STORAGE_IMAGE
  tanpa blok PL/SQL / COMMIT; GetAppName_SQL - prasyarat IsPEGAPROD di InsertGoogleStorage_Act langkah 6 NONAKTIF, selalu
  dibaca); token GCP_IMAGE tanpa CALL (`inti/backend/layanan`).
- **K67-B4 Penyimpanan** (`services/lampiranstorage.go` stub + token, `services/lampiranpengirim.go` HTTP): Folder `Policy/Doc/
  YYYY/MM/`, Namafile `yyyyMMdd-hhmmss-S - <nama>` DIHITUNG DARI ID DOKUMEN (bukan jam unggah) supaya setiap percobaan ulang
  menulis ke objek yang sama; `YYYY` = tahun-minggu Java `[dugaan]` lokal en_US (Desember di minggu 1 Januari -> tahun
  berikutnya); Durasi 1800; URL tersimpan dipakai selama EXPDATE (GMT `[dugaan]`) > kini + 1 menit `[keputusan agent]`, selain
  itu geturl + Update_T_Storage_SQL; EXPDATE kosong = tanpa URL (GetUrlGoogleStorage_Act langkah 4); delete Namafile =
  APPFOLDER tanpa awalan bucket; https saja, pengalihan tidak diikuti; stub tidak menyentuh objek NYATA (APPFOLDER terisi:
  efek gagal permanen, catatan tetap).
- **K67-B5 Outbox** MODUL "NBFACIN", LINI "FAC" `[keputusan agent]`; efek dijalankan seketika sesudah aksi; pekerja latar
  setiap 1 menit `[keputusan agent]` (tanpa env baru - K-B melarang menyentuh inti; mati tanpa Oracle).
- **K67-B6 Unggah**: berkas < 7 byte = 400 "Berkas kosong atau terlalu kecil." (InsertGoogleStorage_Act langkah 5
  `@length(Image) < 10` base64); tanpa ekstensi = "Jenis berkas tidak didukung: (tanpa ekstensi)" `[keputusan agent]`;
  kategori dicocokkan tanpa beda huruf, teks MASTER yang disimpan.

**Penyimpangan sadar dari Pega**: baris DOCUMENT_POLIS lahir sebelum unggah sukses (Pega: sesudahnya) - selama "menunggu" Pega
melihat dokumen tanpa catatan objek (Count ikut menghitungnya); Pega batal diam untuk jenis tak dikenal / berkas kecil -> di
sini 400 berpesan; URL bertanda tangan diproksi backend; View Office Online (alamat layanan luar) tidak diport.

**Uji** - `models.TestIDDokumenPolis` / `TestLokasiUnggah` / `TestLokasiDariLink` / `TestStatusDokumen` / `TestBentukStorage`;
`repository.TestSQLLampiran`; `services.TestTampilanLampiran` / `TestUnggahLampiran` / `TestUnggahLampiranDitolak` /
`TestUnggahLampiranStorageGagal` / `TestUlangiSumberHilang` / `TestHapusLampiran` / `TestUnduhLampiranCaseLain` /
`TestLayakUlangLampiran`; transport httptest TLS `TestStorageJarakJauhSimpan` / `Buka` / `HapusAda` + `TestSumberTokenStorage`;
`handlers.TestLampiranHandler` (kunci api.ts TampilanLampiran / KategoriLampiran / BarisLampiran dua arah). **Mutasi 42 / 42
mati.** gofmt / vet bersih; `go test ./modul/nbfacin/... ./inti/backend/config/ ./cmd/...` hijau.

## Deploy

- Tanpa migrasi. Biner MENULIS tabel Pega `POOLDATA.DOCUMENT_POLIS` (INSERT / DELETE), `T_STORAGE_IMAGE` (INSERT / UPDATE /
  DELETE), `GCP_IMAGE` (INSERT token), outbox `T_LOG_SERVICE_RNM`; membaca `CATEGORY_ATTACH_REAS`, `T_FOLDER_IMAGE`,
  `M_LINK_SERVICE`.
- `UNGGAHAN_DIR` wajib (folder antrean + stub); `PELAKSANA_STORAGE=nyata` menuntut `STORAGE_TOKEN_SALT` (aturan config inti).

## Terbuka

- ✅ **Persetujuan klien HTTP (08-10-2026)**: work owner lewat AskUserQuestion "Setuju" - NB FacIn mengirim lampiran ke Google
  Storage seperti Treaty Contract Out; entri `berkasKlienHTTPDisetujui` untuk `lampiranpengirim.go` ditambahkan sesi 9d di
  `inti/backend/penjaga/lintasaplikasi_test.go`; `TestNolAlamatLayananDiKode` hijau. Catatan lama di bawah dipertahankan:
- ~~Penjaga MERAH~~ (selesai) untuk `modul/nbfacin/backend/services/lampiranpengirim.go` (klien HTTP keluar):
  menghubungkan layanan luar menuntut persetujuan manusia; daftar `berkasKlienHTTPDisetujui` tinggal di
  `inti/backend/penjaga/lintasaplikasi_test.go` (inti - tidak disentuh, K-B). Butuh: persetujuan work owner + satu baris entri
  (pola Treaty Contract Out / Master Product Name Life). Tanpa itu jangan menyetel PELAKSANA_STORAGE=nyata.
- DBA: izin tulis / baca tabel di atas; isi CATEGORY_ATTACH_REAS dan baris M_LINK_SERVICE Google upload / geturl / delete di DEV
  (OQ-047); `[saran]` indeks DOCUMENT_POLIS(IDPEGA, KATEGORI_2).
- `[belum terverifikasi]` lokal JVM Pega (tahun-minggu `YYYY`), zona EXPDATE, `.pyFileMimeType` drag-drop; geser ID / ORA-00001
  dan SQL belum diuji di Oracle DEV.
- Temuan untuk pemilik Treaty Contract Out (tidak diubah): `tco_lampiran.go` `sqlHapusLampiranTCO` / `sqlHapusObjekTCO` =
  hapus fisik tanpa pengecualian ADR-0031 tercatat.
- Tiket 67b: proteksi "Can not proceed spreading ORS" (CheckSpreadingProtect_ACT 10-13, Count APPROVAL SPREADING).


## Keputusan work owner 09-10-2026 — lampiran di tahap Policy / case selesai (lewat sesi 9d)

- Laporan sesi c3 (K69-8, tahap Policy baca-saja): `bolehUnggahLampiran` mengikuti K67-2 (pemegang workbasket posisi, semua
  tahap) sehingga Unggah / Ulangi tetap boleh di FlagOnGoingPolicy 1 dan di case Resolved-Completed (POSITION tidak berubah);
  Hapus tetap FlagOnGoingPolicy 0 atau IT.
- **K67-5** Work owner: *"Biarkan seperti sekarang"* — perilaku di atas dipertahankan, tanpa perubahan kode.
- **K67-6** (work owner 09-10-2026, lewat 9d; menggantikan bagian upload K67-2 dan K67-5): *"INI KENAPA TIDAK BISA UPLOAD
  ATTACHMENT, YANG TIDAK BISA UPLOAD ITU SAAT VIEW POLICY ATAU SUDAH JADI POLIS ATAU VIEW FOLOWING"*, diralat *"Resolved-Decline,
  Resolved-Reject JUGA DIKUNCI JUGA"*.
  - Unggah / Ulangi BOLEH oleh siapa pun yang membuka case, tanpa syarat workbasket.
  - DITOLAK (403, `bolehUnggah` false) bila `T_WORK_POLIS.STATUS_WORK` termasuk daftar EKSPLISIT (bukan awalan `Resolved-`):
    "Resolved-Completed" (polis jadi, K69-7), "Resolved-Decline", "Resolved-Rejected" (ralat ejaan di bawah).
  - Dasar korpus (sesi 0f): `NB FacIn\Section\AttachmentGridReas.xml` tombol Upload sel L2330 `pyVisible=ALWAYS`,
    `pyReadOnly=false`, tanpa syarat.
  - View Policy / View Following ditangani frontend (mode hanyaLihat). Kabel tidak berubah. Aturan Delete (K67-2 / K-F) tidak
    berubah.
  - RALAT EJAAN (work owner 10-10-2026, lewat 9d): "Resolved-Rejected" mengikuti ejaan Pega. Sensus dua sesi
    (`grep -rhoE "Resolved-[A-Za-z]+"` atas NB FacIn; c3 dan 9d): Resolved-Completed 11, Resolved-Rejected 6
    (`NB FacIn\ReportDefinition\GetListOpportunity*.xml`), tidak ada "Resolved-Reject". Daftar kunci = {"Resolved-Completed",
    "Resolved-Decline", "Resolved-Rejected"}; "Resolved-Reject" (bukan status Pega) TIDAK dikunci. Sistem baru belum menulis
    status reject apa pun (yang ditulis: "Resolved-Completed", "Resolved-Decline").
  - Uji: `services/lampiran_test.go` `TestUnggahLampiranK676`, `handlers/lampiran_test.go` `TestLampiranHandler`.
- **K67-7** (work owner 10-10-2026, lewat 9d; sesudah merge GitHub `7d627608` yang menghapus `PELAKSANA_STORAGE` /
  `bacaPelaksanaStorage` dari config inti):
  1. **Selalu nyata (Google Storage).** `modul.go` memasang `svc.DenganLampiran(p.Dasar(), true, cfg.StorageTokenSalt)`;
     alamat dari M_LINK_SERVICE, token dari STORAGE_TOKEN_SALT, sama dengan `inti/backend/penyimpanan`. Stub lokal hanya untuk
     uji (`DenganLampiranUji`). Konsekuensi: di DEV tanpa M_LINK_SERVICE "Google" / STORAGE_TOKEN_SALT, baris DOCUMENT_POLIS
     tercatat tetapi efek outbox gagal (status "gagal") dan dapat diulang (Ulangi).
  2. **Pilihan (A)**: nbfacin tetap memakai port tiket 67 sendiri (outbox asinkron, `services/lampiranpengirim.go`, aturan
     K67-6). Adaptor di atas `inti/backend/penyimpanan` (B) dan pindah ke `inti/backend/dokumenpolis` + `PanelLampiranReas` (C)
     TIDAK dikerjakan.
  - Teks persetujuan penjaga `inti/backend/penjaga/lintasaplikasi_test.go` L208 diperbarui sesi 9d / 0f.
- **K67-8** (work owner 10-10-2026, lewat 9d — cara 1): bug DEV 11:18 `nbfacin lampiran: … memungut efek keluar: ORA-02014:
  cannot select FOR UPDATE from view with DISTINCT, GROUP BY, etc.`, muncul sejak K67-7 menghidupkan pekerja.
  - Sebab: pungutan memakai `ORDER BY … FETCH FIRST 1 ROWS ONLY FOR UPDATE SKIP LOCKED`. Oracle melarang row_limiting_clause
    bersama `FOR UPDATE` (ORA-02014), walau T_LOG_SERVICE_RNM tabel biasa (migrasi claimlife 015 L25 `CREATE TABLE`).
    `[terverifikasi di kode]` batasan yang sama sudah tercatat di `modul/masterproductnamelife/backend/repository/
    mpnl_lampiran.go` L114; `[dugaan]` sumber dokumentasi Oracle tidak dibaca di sesi ini.
  - Perbaikan: kandidat dipilih lewat subkueri ROWNUM berurut JADWAL_BERIKUT, lalu `SELECT … FROM T WHERE ID = (…) FOR UPDATE
    SKIP LOCKED` atas tabel dasar. Diterapkan di dua tempat:
    - `inti/backend/outbox/penyimpan.go` `sqlPungutEfek` (inti lintas modul; saringan `AND MODUL = :3` tetap);
    - `repository/lampiran.go` `sqlPungutEfekRujukan` (pungutan per rujukan nbfacin, bentuk yang sama + `AND RUJUKAN = :4`).
    Penandaan `UPDATE … AND STATUS = :4` tetap pengaman pemungutan ganda. Harga: bila kandidat sedang dikunci worker lain,
    putaran itu kosong dan dicoba lagi putaran berikutnya.
  - Uji: `TestPungutEfekTanpaRowLimitingDenganForUpdate` (teks SQL utuh, bind :1-:3 berurutan sekali, kunci di tabel dasar);
    `TestSQLLampiran` (bentuk baru + penjaga: kueri FOR UPDATE tanpa FETCH FIRST). Tanpa query / tulisan ke DEV.

## Bug DEV 11-10-2026 — semua unggah lampiran 'gagal-permanen' "STORAGE_TOKEN_SALT kosong" (sesi c3, lewat sesi 9d)

**Gejala `[terverifikasi]`** (query work owner): T_LOG_SERVICE_RNM MODUL 'NBFACIN' JENIS_EFEK 'storage-unggah' ID 14 / 15 / 16 / 18
ber-GALAT_TERAKHIR "token penyimpanan gagal diambil: STORAGE_TOKEN_SALT kosong". `APP_RNM\.env` tidak memuat kunci
STORAGE_TOKEN_SALT (dicek ada / tidaknya kunci saja). Product Name Life di backend yang sama bisa unggah / view.

**Sebab `[terverifikasi]`:**

- `services/lampiranstorage.go` `sumberTokenStorage.TokenBaru` menolak garam kosong SEBELUM membaca token GCP_IMAGE.
- `masterproductnamelife/backend/services/mpnl_storage.go` `tokenStorage` (urutan procedure GET_TOKEN_STORAGE) memakai ulang token
  berlaku TANPA garam. Garam hanya dituntut saat merakit token baru.

**Perbaikan** (keputusan work owner "tiru menu Product Name Life"):

- APPNAME diperiksa lebih dulu.
- Token bersisa > margin 15 detik dipakai (seperti sebelumnya).
- Bila tidak ada token bersisa > margin dan garam kosong: token berlaku APA PUN (`INPUTDATE > saat`, sisa > 0) dipakai seperti
  MPNL. Baru bila itu pun tidak ada → `ErrGaramTokenKosong`.
- Dengan garam terisi, perilaku lama (margin, rakit baru) TIDAK berubah.
- Baris efek lama yang sudah gagal-permanen tidak disentuh; work owner klik Ulangi sesudah backend direstart.

**Uji** (`services/lampiranpengirim_test.go`):

- `TestTokenStorageTanpaGaram` (a): garam kosong + token sisa 10 detik → dipakai tanpa galat / simpan.
- (b): garam kosong + tanpa token → `ErrGaramTokenKosong`.
- (c): garam terisi + token sisa 10 detik → dirakit baru.
- Sisa > margin → dipakai dengan / tanpa garam.
- `TestSumberTokenStorage`: APPNAME kosong ditolak sebelum garam.
- Tiga mutasi semuanya tertangkap.

**Catatan:** tanpa garam, unggah tetap bergantung pada token GCP_IMAGE yang ditulis pihak lain (MPNL / Pega) dan masih berlaku.
Bila tidak ada, unggah gagal dengan galat yang sama. Pengisian STORAGE_TOKEN_SALT di `.env` tetap jalan yang lengkap.

## View / View Office Online di popup View File (permintaan work owner 11-10-2026, lewat sesi 9d; sesi c3)

- **View pdf / gambar:** frontend memakai `GET …/lampiran/berkas/{id}/isi` yang sudah ada. Tanpa perubahan backend.
- **`GET /api/nbfacin/kasus/{caseId}/lampiran/berkas/{id}/office` → 200 `{"url": "<https bertanda tangan>"}`**
  (`services/lampiranoffice.go` `TautanOfficeLampiran`, pola `inti/backend/dokumenpolis` TautanOffice):
  - otorisasi / tahap sama dengan `/isi` (tahapLampiran, dokumen milik case);
  - ekstensi hanya xls / xlsx / doc / docx / ppt / pptx (`ReasViewAttachment` b2951; `DOCUMENT_POLIS.MIME` huruf kecil); selain
    itu 400;
  - belum tersimpan (tanpa catatan T_STORAGE_IMAGE: menunggu / gagal) → 409 `ErrLampiranBelumTersimpan`;
  - stub lokal (klien penyimpanan tanpa tautan, atau URLPUBLIC kosong) → 409 `ErrLampiranOfficeTanpaStorage` "View Office Online
    butuh penyimpanan nyata (Google Storage) …";
  - objek hilang di penyimpanan → 409 `ErrLampiranTanpaBerkas`;
  - URL = GetUrlGoogleStorage_Act seperti `Buka`: URL tersimpan dipakai selama berlaku (sisa ≥ 1 menit), selain itu geturl
    (Durasi 1800) + Update_T_Storage_SQL; hanya https. `Buka` sendiri tidak diubah (`storageJarakJauh.Tautan` baru);
  - ⛔ URL tidak ditulis ke log / pesan galat; jawaban `Cache-Control: no-store`.
- **Uji:**
  - `services/lampiranoffice_test.go`: layanan (pdf 400, stub 409, case lain 404, nyata tanpa URLPUBLIC 409, nyata → URL, objek
    hilang 409, belum tersimpan 409) dan klien nyata (tersimpan tanpa geturl, geturl, http ditolak tanpa URL di pesan, tanpa
    EXPDATE);
  - `handlers/lampiranoffice_test.go`: 200 + no-store, 400 / 409 / 409 / 404 / 401.
  - Mutasi 5/5 tertangkap.

## Catatan silang — IDPEGA = case id polos (11-10-2026)

Tulisan IDPEGA tiket ini kini memakai case id polos; pembaca menerima dua bentuk (polos / berawalan). Rincian: `docs/issues/69-jalur-policy.md`
bagian "IDPEGA = case id polos".
