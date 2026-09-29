# 12: Lampiran di tahun treaty — **FITUR BARU**

**Status:** selesai (29-09-2026)

**Blocked by:** 03 (tahun treaty sebagai induk lampiran)

⚠️ **Ini fitur baru, bukan migrasi.** `[keputusan work owner]` Jalur lampiran di Pega
(`M_ATTACHMENTTREATY_2`) **belum rampung di-develop**. Di sistem baru fungsinya **dilengkapi** dan
melekat pada **tahun treaty**.

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin **melampirkan berkas pada sebuah tahun treaty** beserta
kategorinya, supaya dokumen kontraknya tersimpan bersama datanya — **tanpa** itu menjadi syarat
tersimpannya data tahun treaty; dan sebagai **tim operasi**, saya ingin **kegagalan unggah terlihat
dan dapat diulang**, serta alamat penyimpanan **dibaca dari konfigurasi saat dijalankan**.
*(User story 31–36 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Rekam lampiran; resolusi alamat penyimpanan; **cache token** |
| `internal/clients` | Klien penyimpanan berkas **di balik interface** |
| `internal/services` | Orkestrasi efek keluar; penandaan kegagalan; pengulangan |
| `internal/handlers` | Endpoint unggah / unduh / hapus / status |
| `frontend/` | Panel lampiran di layar tahun treaty; penanda "terkirim" / "tertunda" |

## Rule Pega sumber

`[terverifikasi]` Rantai berkas yang ada sekarang — pola yang sama dengan Master Product Name Life:

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `ServiceGoogle` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `SERVICEGOOGLE` / `RULE-CONNECT-REST` | `Treaty Contract Out/ConnectREST/ServiceGoogle.xml` | **satu-satunya** ConnectREST modul ini |
| `LinkService` | `LINKSERVICE` / `LINKSERVICE` / `RULE-ADMIN-SYSTEM-SETTINGS` | `Treaty Contract Out/SystemSettings/LinkService.xml` | sumber alamat |
| `GetTokenStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetTokenStorage_SQL.xml` | → `POOLDATA.GET_TOKEN_STORAGE` |
| `GetLinkStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetLinkStorage_SQL.xml` | ambil tautan berkas |
| `Update_T_Storage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/Update_T_Storage_SQL.xml` | perbarui rekam |
| `DeleteStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteStorage_SQL.xml` | hapus rekam |
| `InsertAtatchment_Sql` | `ASM-FW-GISFW-INT-TREATY_IN` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/InsertAtatchment_Sql.xml` | → `POOLDATA.PEGA_M_ATTACHMENT` (CLOB) |
| `GetAllAttachment2_Sql` | `ASM-FW-GISFW-INT-TREATY_IN` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetAllAttachment2_Sql.xml` | ⚠️ berkunci `TreatyIn.ID` |
| `CategoryAttach_SQL` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/CategoryAttach_SQL.xml` | → `CATEGORY_ATTACH_REAS` |
| `GetLinkService` | `@BASECLASS` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/GetLinkService.xml` | resolusi alamat |
| `TreatyOutSaveAttachment`, `TreatyOutDownloadAll_Act`, `TreatyOutDownloadOne`, `LoadAttachmentTreatyOut` | `Treaty Contract Out/Activity/` | jalur berkas yang belum rampung |

`[terverifikasi]` `ServiceGoogle` **tanpa URL literal** — alamat datang dari `LinkService`.
⚠️ `pzOriginalInstanceKey`-nya menunjuk `…GOOGLESTORAGE_UPLOAD` — nama berkas ≠ nama asal
(**OQ-066** lagi).

⚠️ **Kebocoran batas yang tidak dibawa.** `[terverifikasi]` `GetAllAttachment2_Sql` berclass
**`ASM-FW-GISFW-INT-TREATY_IN`** dan membaca `M_ATTACHMENTTREATY_2 WHERE treatyid = {TreatyIn.ID}` —
lampiran existing dikunci ke **ID treaty inward**, bukan ke tahun treaty. `[keputusan work owner]`
Di sistem baru lampiran melekat pada **tahun treaty**; kepemilikan `M_ATTACHMENTTREATY_2` tetap di
konteks treaty inward.

⚠️ **Penyimpangan sadar 9 — fitur baru.** `[keputusan work owner]`

## ADR terkait

**ADR-0013** (alamat di-resolve **runtime**; **dilarang** sebagai literal, konstanta, **maupun env
var** — menggantikan ADR-0004), **ADR-0015** (efek keluar — kegagalan tidak boleh diam-diam),
**ADR-0010** (penyimpanan berkas tetap Google Storage), **ADR-0007** (jejak audit).

## Acceptance criteria

- [ ] ⚠️ Berkas dapat **dilampirkan pada sebuah tahun treaty** — fitur yang **belum ada** di Pega.
      *(AC 54 spec; User story 31; penyimpangan sadar 9)*
- [ ] ⚠️ Lampiran **opsional**: kegagalan unggah **tidak** membatalkan tersimpannya tahun treaty —
      hanya status lampiran yang berubah. *(AC 55 spec; User story 32)*
- [ ] Tiap lampiran dapat diberi **kategori** dari master kategori. *(AC 56 spec; User story 33)*
- [ ] Lampiran dapat **diunduh** dan **dihapus**. *(AC 57 spec; User story 34)*
- [ ] ⚠️ Kegagalan unggah **tercatat dan dapat diulang**; pengulangan **tidak** menggandakan berkas.
      *(AC 58 spec; User story 35; **ADR-0015**)*
- [ ] Alamat penyimpanan di-resolve **runtime** dari konfigurasi. Test yang memindai kode untuk URL
      sebagai **literal, konstanta, atau pembacaan env var** **gagal** bila menemukannya.
      *(AC 59 spec; User story 36; **ADR-0013**)*
- [ ] Token **di-cache** dan **diperbarui sebelum kedaluwarsa**; token yang **gagal diambil**
      menghasilkan kegagalan yang **terlihat**, bukan unggahan yang diam saja tidak terjadi.
      *(AC 60 spec)*
- [ ] Rekam berkas di basis data dan berkas di penyimpanan **tetap sejalan**: rekam tanpa berkas
      terdeteksi dan dapat diperbaiki. *(AC 61 spec)*
- [ ] Klien penyimpanan berkas berada **di balik interface** dan **di-fake** di test; yang diperiksa
      adalah **efeknya**. *(AC 62 spec)*
- [ ] Penghapusan berkas yang **sudah tidak ada** di penyimpanan **tidak** menggagalkan penghapusan
      rekamnya.
- [ ] Lampiran melekat pada **tahun treaty**, bukan pada ID treaty inward. Test yang menemukan
      rujukan ke ID treaty inward **gagal**.
- [ ] Setiap unggah dan penghapusan mencatat **jejak audit**. *(**ADR-0007**)*

## Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **OQ-047 — tidak memblokir:** alamat fisik Google Storage tersimpan di
`M_LINK_SERVICE`, dan isinya belum dibaca. **ADR-0013** sudah mengatur **cara** meresolusinya, jadi
tiket ini dapat selesai tanpa mengetahui alamatnya — yang mengikat adalah **alamat tidak ditanam**.

## Catatan

⚠️ **Taruhannya berbeda dari Komite Claim Life.** Di sana efek keluar **wajib berhasil**
(transactional outbox, **ADR-0015** versi Komite). Di sini `[keputusan work owner]` lampiran
**opsional** — yang wajib adalah **kegagalannya terlihat dan dapat diulang**, bukan berhasil.
Sama dengan perlakuan di Master Product Name Life.

⚠️ **Fitur ini tidak lengkap di Pega** `[keputusan work owner]` — jadi korpus bukan sumber
kebenaran perilaku di sini, hanya sumber **rantai teknis** (token, alamat, rekam berkas, kategori).
Perilaku yang tidak ada di korpus **tidak ditebak**; ia ditetapkan oleh AC di atas.

⚠️ **Satu-satunya efek keluar konteks ini.** `[terverifikasi]` Modul ini punya **satu**
`RULE-CONNECT-REST`, dan tidak ada integrasi Arasapas / Kasir / Konversi / Gemini AI. Berbeda dari
Master Contract Retro Life yang **tanpa `ConnectREST` sama sekali** — di sana ADR-0013 dan ADR-0015
tidak berlaku; di sini berlaku keduanya.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**; **penyimpanan berkas di-fake** di balik
interface — ia satu-satunya yang dipalsukan di seluruh konteks ini.

```
go test ./internal/...
cd frontend && npm test
make check
```


---

## Pembacaan ulang XML — 29-09-2026 (sesi modul, lanjutan 1)

Nomor baris = baris mentah berkas korpus (satu tag per baris), diverifikasi dengan `awk 'NR==n'`. Korpus di sini adalah
sumber RANTAI TEKNIS dan LABEL; perilakunya ditetapkan AC tiket ini (penyimpangan sadar 9).

| Unsur | Bukti | Dibawa sebagai |
| --- | --- | --- |
| letak panel | `Section/InputTreatyContract.xml` b11721 `Attachment for`, b13074 `<pyInclude>GridTreatyArrangementAttachment</pyInclude>` | panel di form tahun treaty yang sudah ber-ID (`PanelLampiranTahun.tsx`) |
| tombol dan kolom | `Section/GridTreatyArrangementAttachment.xml`: b578 `Add attachment` (→ `SetCategory_act` b596 → flow action `TreatyOutAttachContent` b642), b1023 `Refresh` (→ `LoadAttachmentTreatyOut` b1041), b1785 `For Treaty Contract Out`, b2659 `Download All` (→ `TreatyOutDownloadAll_Act` b2677), b3032 `File Name` (sel `.pyFileName` b3428 → `TreatyOutDownloadOne` b3488), b3170 `Type` (sel `.pyCategory` b3705), b3897 `Delete` (→ `DeleteAttachmentTreaty` b3915) | `LAMPIRAN_TCO`, 10 baris diuji terhadap korpus |
| tombol `Download` | b2391 → `DownloadAll_Act` generik | ➖ tidak dibawa sebagai tombol kedua; `Download All` adalah aksi yang sama untuk modul ini |
| pesan tanpa berkas | `Activity/TreatyOutSaveAttachment.xml` b376 `"Tidak ada file yg diattach"` | 400 dengan teks VERBATIM |
| kunci lampiran Pega | `TreatyOutSaveAttachment.xml` b1402 dan `DeleteAttachmentTreaty.xml` b252: `TreatyYear + TreatyYearID`; `RDBList/GetAllAttachment2_Sql.xml` `M_ATTACHMENTTREATY_2 where treatyid = {TreatyIn.ID}` | ➖ tidak dibawa; kunci = `IDTREATYYEAR`; penjaga statik Go dan JS menolak rujukan treaty inward |
| penulisan Pega | `RDBList/InsertAtatchment_Sql.xml` → `POOLDATA.PEGA_M_ATTACHMENT` (CLOB JSON + `COMMIT`) | ➖ tidak dibawa (prosedur tidak dipanggil, nol COMMIT, nol JSON/CLOB) |
| master kategori | `RDBList/CategoryAttach_SQL.xml` `select * from CATEGORY_ATTACH_REAS order by note`; `Activity/SetCategoryAttachTreatyin.xml` b500 `.NOTE` | dibaca saja: `SELECT DISTINCT NOTE … ORDER BY NOTE`; kolom `CATEGORY` menyimpan teks `NOTE` master apa adanya |
| rantai penyimpanan | `RDBList/GetTokenStorage_SQL.xml` (`GET_TOKEN_STORAGE` + `COMMIT`), `GetLinkStorage_SQL.xml`, `Update_T_Storage_SQL.xml`, `DeleteStorage_SQL.xml` atas `T_STORAGE_IMAGE`; `Activity/GetLinkService.xml` b705–b706 | bentuknya ditiru di balik antarmuka (`PenyimpananJarakJauhTCO`: resolver `M_LINK_SERVICE` saat jalan + cache token); `T_STORAGE_IMAGE` dan prosedur tidak disentuh |
| `ConnectREST/ServiceGoogle.xml`, `SystemSettings/LinkService.xml` | masing-masing tepat satu `://`, dan keduanya di `<pyHelpURI>` (tautan bantuan platform), bukan alamat layanan | klaim tiket "tanpa URL literal" terkonfirmasi; nilai tautan tidak disalin |

### Ralat bertanggal 29-09-2026

1. **Area codebase `internal/clients`** — paket itu tidak ada di codebase (`internal/` = config, handlers, models,
   repository, services). Klien penyimpanan (`KlienPenyimpananTCO`) diletakkan di `services`, sejajar dengan
   `ResolverEndpoint` dan `PelaksanaEfek` yang sudah ada di sana.
2. **"Cache token"** — jalur token yang ada (`services.TokenStorage` atas `GCP_IMAGE`) menulis tabel warisan dan menuntut
   garam `STORAGE_TOKEN_SALT`. Karena endpoint penyimpanan nyata tidak dipanggil, yang dibangun adalah `CacheTokenTCO`
   (pakai ulang, perbarui 15 detik sebelum kedaluwarsa, gagal terang) di atas antarmuka `SumberTokenTCO`. Sumber
   token Oracle belum dipasang; ia bagian dari penyambungan nyata yang menuntut persetujuan manusia.
3. **Pekerja outbox** — tidak ada pekerja yang berjalan di `cmd/api` untuk modul mana pun (nol pemanggil `SatuPutaran`
   di kode produksi). Antrean modul ini dijalankan sesudah unggah, hapus, dan ulangi; `SatuPutaran` diekspor untuk
   pekerja latar kelak. Pekerja bawaan `PekerjaEfek` tidak dipakai karena jalur menyerahnya menulis jejak Claim Life.
4. **Kolom master kategori** — hanya `NOTE` yang terbukti dipakai; `CATEGORY_ID` di `M_ATTACHMENTTREATY_2` tidak dibawa
   karena kolom ID master tidak terbukti di korpus (`select *`).

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| skema | `migrations/307_t_treatyyear_lampiran.sql` (+`_down`) | `T_TREATYYEAR_LAMPIRAN`, FK ke `T_TREATYYEAR` tanpa kaskade, `IMAGEID` unik, `SEQ_T_TREATYYEAR_LAMPIRAN`, index tahun |
| models | `tco_lampiran.go` | `LampiranTCO`, `StatusLampiranTCO` (terkirim menang), `KategoriLampiranSah`, `NamaBerkasAntreLampiranTCO` (nama unggahan tidak menentukan jalur), `NamaEntriZipLampiranTCO` |
| repository | `tco_lampiran.go` | `KategoriLampiran` (dibaca saja), `MasterLampiranTCO` (daftar + efek unggah terakhir dari outbox bersama, ambil berbatas tahun, kunci `FOR UPDATE`, sisip, hapus, tandai) |
| services | `tco_lampiran.go`, `tco_penyimpanan.go` | `LampiranTahunTCO` (unggah, daftar, unduh, unduh semua, hapus, ulangi, periksa selaras, pekerja modul sendiri), stub lokal, rangkaian jarak jauh, cache token |
| handlers | `tco_lampiran.go` (+ uji, + uji `db`) | 8 rute; pemetaan galat 400/404/409/413/422/503 |
| frontend | `PanelLampiranTahun.tsx` (+ uji), `LAMPIRAN_TCO`, `api.ts` (+9) | unggah multipart, unduh lewat `fetch` berheader identitas, status + galat terlihat, `Ulangi`, `Periksa keselarasan` |

**Status:** selesai 29-09-2026 — commit `treaty-contract-out: tiket 12 — lampiran di tahun treaty`.

## Keputusan work owner 29-09-2026

- **OQ-TCO-08 — ditutup.** Jawaban: *"sekarang"*. Pelaksana **nyata** efek penyimpanan lampiran terpasang di balik
  `PELAKSANA_STORAGE=nyata`; bawaan tetap `stub` (`PenyimpananLokalTCO`, folder `UNGGAHAN_DIR`). Rangkaiannya:
  `PenyimpananLampiranTCO` (pemilih) → `PenyimpananJarakJauhTCO` (alamat per operasi dari `M_LINK_SERVICE` **saat jalan**:
  `Google/upload`, `Google/geturl`, `Google/delete` — tidak disalin ke mana pun) → `CacheTokenTCO` → sumber token
  `GET_TOKEN_STORAGE` ditiru di Go (`GCP_IMAGE`: token berlaku dipakai ulang beserta kedaluwarsanya, atau
  `RakitToken(garam, saat)` disimpan dengan umur 1 menit, satu transaksi; `App` = `T_FOLDER_IMAGE.APPNAME`, dibaca saja)
  → transport `NewPengirimBerkasHTTPTCO` (`services/tco_pengirim_storage.go`: POST JSON halaman `UploadDoc`, batas waktu
  300 s dari `ServiceGoogle.xml` b29; unduh lewat URL bertanda tangan yang diberikan layanan). Garam dari env
  `STORAGE_TOKEN_SALT`; `nyata` tanpa garam **ditolak saat menyala**. Galat tidak pernah memuat alamat, token, atau garam.
  Uji memakai **server tiruan lokal** (`httptest`); tidak ada layanan sungguhan yang dipanggil dari uji maupun sesi ini.
- **OQ-TCO-09 — ditutup.** Jawaban: *"perlu"*. `LampiranTahunTCO.JalankanPekerja(ctx, interval, catat)` menjalankan
  `JalankanAntrean` (paling banyak 50 efek per ketukan) setiap `TCO_PEKERJA_LAMPIRAN_INTERVAL`; kosong/0 = **mati**
  (bawaan, juga di uji). Dihidupkan dari `cmd/api` (`jalankanPekerjaLampiranTCO`) hanya bila Oracle terpasang, dan berhenti
  bersama sinyal proses. Coba ulang dan anti-dobel memakai yang sudah ada: `Backoff` + `percobaanMaksimum` outbox,
  `FOR UPDATE SKIP LOCKED` pemungutan, nama objek = `IMAGEID` (unggah ulang menimpa objek yang sama, AC 58).
- **Galat permanen baru** (tidak diputar ulang antrean): layanan menolak bentuk permintaan (400/422), garam kosong, `APPNAME`
  kosong. Jaringan, batas waktu, 401/403/429/5xx tetap dicoba ulang.
- **Pembacaan ulang XML — `pyStepsBlockName` dicetak:** Claim Life `InsertGoogleStorage_Act` — seluruh langkah kosong
  kecuali `EXIT` (b2953, `Page-Remove`); Treaty Contract Out `GetUrlGoogleStorage_Act` — kosong kecuali `EXIT` (b2790);
  `DeleteGoogleStorage_Act` — kosong kecuali `EXIT` (b1765). Nol langkah ter-remark `//`: setiap langkah yang dikutip
  transport memang berjalan di Pega.
- **Dibuka: OQ-TCO-22** — `[keputusan kami]` untuk yang korpus modul ini tidak memuat (unggah lampiran fitur baru):
  `Folder = "TreatyContractOut/"` (Claim Life menyusun `Param.Folder + "/Doc/" + tahun/bulan + "/"`), `Durasi = 60`,
  dan `Namafile = IMAGEID`.
- **Kode bersama yang disentuh (aditif):** medan `Service.penyimpananTCO` (`services/services.go`), dua kunci
  `internal/config` + `.env.example`, `cmd/api/main.go`, dan peta `berkasKlienHTTPDisetujui` di
  `services/efekkeluar_statik_test.go` (satu baris, jumlahnya dikunci; `://` dan env tetap diperiksa untuk berkas itu).

### Ralat bertanggal 29-09-2026 — temuan /code-review kelompok 6

- **Jalur hapus salah (diperbaiki).** Delete kini mengirim `Namafile` = jalur objek PENUH `folder + IMAGEID` tanpa `Folder`
  (`DeleteGoogleStorage_Act` b1091); folder berakhiran `/` seperti Pega (Insert b1407–b1408, GetUrl b1260–b1326). Dulu
  delete menunjuk akar bucket, dan objek sebenarnya tertinggal.
- **404 titik layanan ≠ berkas hilang (diperbaiki).** 404 dari upload/geturl/delete = galat layanan (dicoba ulang,
  terlihat); "berkas tidak ada" hanya dari URL bertanda tangan. `Hapus` memeriksa keberadaan lebih dulu (`geturl` + GET
  `Range: bytes=0-0`), jadi berkas yang sudah tidak ada tetap tidak menggagalkan (AC lama) tanpa menelan jalur
  `M_LINK_SERVICE` yang salah.
- **Token (diperbaiki).** Token dipakai ulang hanya bila sisa umurnya > `MarginTokenTCO` (AC 60; penyimpangan sadar kecil
  dari `GET_TOKEN_STORAGE` yang memakai ulang token apa pun yang belum lewat). Sisa umur dihitung DI ORACLE terhadap
  waktu yang di-bind — `INPUTDATE` tidak dibaca ke jam aplikasi (DATE tanpa zona). 401/403 mengosongkan cache token.
- **`ext` (diperbaiki).** Dari nama berkas asli, huruf kecil tanpa titik (`@toLowerCase(Param.Ext)` b587), bukan
  ditebak dari tabel MIME OS. `KlienPenyimpananTCO.Simpan` dan `PengirimBerkasTCO.Kirim` menerima `ekstensi`.
- **URL bertanda tangan (diperbaiki).** Hanya https; pengalihan tidak diikuti.
- **Pekerja (diperbaiki).** `TCO_PEKERJA_LAMPIRAN_INTERVAL` minimal `1s` (atau 0); `cmd/api` menunggu pekerja berhenti
  sebelum db ditutup. Pemilih kini `DenganPenyimpananLampiranTCO(nyata bool, garam)` — nilai env diputuskan satu kali di
  `internal/config`.
- **Dibiarkan, dengan alasan:** (a) I/O jaringan di dalam transaksi outbox — pola bersama `PekerjaEfek` (pungut,
  jalankan, tuntaskan satu transaksi; `SKIP LOCKED` hidup selama transaksi); memisahkannya menuntut desain sewa baris
  lintas modul; risikonya: baris lampiran terkunci selama unggahan berjalan (paling lama batas waktu 300 s). (b) Memori
  unggahan ~3× ukuran berkas (base64 + JSON) — batas 25 MiB, Pega pun menaruh base64 utuh di halaman. (c) Kueri token
  modul tidak memakai `PohonKlaim.TokenBerlaku` bersama — ia butuh sisa umur + margin, dan kode token bersama dipakai modul
  klaim. (d) `APPNAME` dibaca per operasi — seperti resolver `M_LINK_SERVICE`, perubahan DBA berlaku tanpa restart.

## Ralat bertanggal 29-09-2026 — tco4 (nol tabel baru) `[keputusan work owner]`

- **RALAT "kunci treaty inward"**: `TreatyIn.ID` di modul ini diisi `TreatyYear + TreatyYearID`
  (`TreatyOutSaveAttachment` b1402, `DeleteAttachmentTreaty` b252). Lampiran kini di tabel warisan
  `M_ATTACHMENTTREATY_2` (kolom = `Treaty In/RDBList/InsertAttachment2_Sql.xml` b84; ID `YYYYMMDDHH24MISSFF3`;
  `CATEGORY` "File", kategori pilihan di `CATEGORY_ID`) dan objeknya di `T_STORAGE_IMAGE` (`Insert_T_Storage_SQL`
  Claim Fac In b85, `DeleteStorage_SQL` b85). `T_TREATYYEAR_LAMPIRAN` dibuang.
- Badan `PEGA_M_ATTACHMENT` (penulis Treaty Contract Out, `InsertAtatchment_Sql` b60) **[terbuka — DBA]** → OQ-TCO-24.
- Terkirim = objek tercatat di `T_STORAGE_IMAGE`. Ukuran berkas tidak disimpan (kolom tidak ada; layar Pega tidak
  menampilkannya) — dibuang dari API dan layar. **Jejak lampiran gugur**; menyerah terlihat dari status outbox.
- OQ-TCO-22 (`Folder`/`Durasi`/`Namafile`) tetap untuk work owner; nilai sekarang berlabel `[terbuka — OQ-TCO-22]`.

## Ralat bertanggal 29-09-2026 — lanjutan 4, OQ-TCO-24/26 `[asisten dari data; veto work owner]`

- **OQ-TCO-26 ditiru**: `Update_T_Storage_SQL` dijalankan sesudah tiap `geturl` yang berhasil
  (`PenyimpananJarakJauhTCO.segarkan`), dengan `exp` dan `DateTime` dalam bentuk To_date-nya. `exp` pada jalur unggah
  kini juga diubah seperti `InsertGoogleStorage_Act` b2366/b2431. Rinciannya di `dba-procedures.md`.
- **OQ-TCO-24**: rekonsiliasi korpus sudah dicatat. Satu pemanggil hidup `PEGA_M_ATTACHMENT` (`TreatyOutSaveAttachment`
  b1674); pembaca/penghapus membaca `M_ATTACHMENTTREATY_2`; nol sebutan `M_ATTACHMENTTREATY`/`ID_COUNT` di korpus.
  Badan prosedurnya belum terbaca (kueri siap di `dba-procedures.md`). Kode tetap menulis `_2` dan tidak disatukan.
- Uji: `TestExpStorageTCOSepertiPega`, `TestSQLPerbaruiObjekTCOSepertiUpdateTStorage`,
  `TestPenyimpananMenyegarkanObjekSesudahGetURL`, `TestPenyimpananNyataMenyegarkanObjekSesudahGetURL`, dan tag `db`
  `TestLampiranTahunTreatyLingkaranPenuh`.
