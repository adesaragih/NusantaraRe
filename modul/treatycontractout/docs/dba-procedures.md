# Body Procedure — dari DBA (Treaty Contract Out)

Tanggal: 2026-09-15
Sumber: **DBA / work owner** (dikirim langsung; bukan korpus Pega). `[data DBA]`

---

## PEGA_PROPORTIONALARRG (induk, 35 param) & PEGA_M_PROPORTIONALARRG_CHILD (26 param)

### ⚠️ TEMUAN PASTI — menjawab "tersimpan di mana"

1. **KEDUANYA menulis ke SATU tabel: `POOLDATA.PROPORTIONALARRG`** (relasional, kolom bernama).
   **TIDAK ADA** yang menulis ke `M_PROPORTIONALARRG` (JSON). → sumber kebenaran tulis =
   **`PROPORTIONALARRG` relasional**. (Konfirmasi tebakan work owner: "kesimpan di proportionalarrg".)
2. **`m_PROPORTIONALARRG` (JSON) DIBACA sebagian kueri tapi TIDAK PERNAH DITULIS 2 procedure ini.**
   Dari mana isinya? tak terbukti korpus → OQ. Memperkuat Q1=(a): buang JSON, satu sumber relasional.
3. **`PEGA_M_PROPORTIONALARRG_CHILD` nama MENIPU (OQ-066):** menulis ke tabel yang **sama**
   (`PROPORTIONALARRG`), bukan tabel child terpisah. Bedanya hanya **jumlah kolom**:
   - Induk (`PEGA_PROPORTIONALARRG`): **35 kolom** — termasuk 9 kolom khusus induk:
     `ID_OCCUPATION, OCCUPATION, ID_CLAUSE, CLAUSE, TREATYLIMIT, COINS_MIN, COINS_MAX, MORERP, MOREUSD`.
   - "Child" (`PEGA_M_PROPORTIONALARRG_CHILD`): **26 kolom** — tanpa 9 kolom itu.
   → Di skema baru: **satu tabel `proportionalarrg`**, 9 kolom induk NULL untuk baris "child".
     Ini **menguatkan Q2=(a)** (satu tabel generik untuk semua klausul).

### Pola simpan (kedua procedure identik)

| Aspek | Perilaku |
| --- | --- |
| Mode | **UPSERT dikunci `ID`** (`SELECT COUNT(1) WHERE ID=P_ID` → UPDATE / INSERT) |
| ID baru | `'1' \|\| lpad(PROPORTIONALARRG_SEQ.nextval, 7, '0')` — **7 digit**, mis. `10000001` |
| `TGLUPDATE` | `SYSDATE` (param `P_TGLUPDATE` diabaikan) |
| Transaksi | **TIDAK commit sendiri** — hanya `ROLLBACK` on error; commit diserahkan pemanggil ✅ |
| Keluaran | `StsSimpan` **1=sukses / 0=gagal**; `ErrMsg` teks (bilang "JSON_KLAIM" = copy-paste template, **bukan** berarti JSON) |

⚠️ **Batas transaksi BEDA dari procedure produk Life** (yang commit sendiri): di sini **tidak
commit sendiri**, jadi Go BISA membungkus banyak klausul dalam satu transaksi atomik. Bagus untuk
simpan-banyak-klausul-atomik.

### Kolom `PROPORTIONALARRG` (dari daftar INSERT — 35 kolom, terverifikasi body)

`ID, TREATYYEAR, TREATYYEARID, TREATYGROUPID, TREATYGROUPNAME, TREATYDESCID, TREATYDESCNAME,
REINSTYPEID, REINSTYPENAME, LAYER, LAYERPART, LAYERPARTTYPE, LAYERTYPE, KURS, TGLUPDATE, USERID,
LINE, PCT, PCTME, YDCF, METHOD, TERRITORIALLIMIT, PARENTREINSTYPEID, SPREADINGORDER, RP, USD,
ID_OCCUPATION, OCCUPATION, ID_CLAUSE, CLAUSE, TREATYLIMIT, COINS_MIN, COINS_MAX, MORERP, MOREUSD`

- Jenis klausul dibedakan **`TREATYDESCID`/`TREATYDESCNAME`** (dari master `TREATYDESC`).
- Uang: `RP`, `USD`, `MORERP`, `MOREUSD`, `TREATYLIMIT`, `COINS_MIN`, `COINS_MAX` → decimal (ADR-0003).
- Persen: `PCT`, `PCTME`.
- Belum jelas artinya (fakta bisnis Q3): `YDCF`, `METHOD`, `SPREADINGORDER`, `LAYERPARTTYPE` vs
  `LAYERTYPE`, `PCTME` vs `PCT`, `LINE`.

## Dampak ke frontier

- **Q1** ✅ dikuatkan — tulis hanya ke relasional `PROPORTIONALARRG`; JSON tak ditulis → buang JSON.
- **Q2** ✅ dikuatkan (a) — satu tabel, "child" = subset kolom (9 kolom induk NULL). Bukan 2 tabel.
- **Batas transaksi** — procedure tak commit sendiri → Go bisa atomik lintas klausul (beda dari Life).
- ✅ **`M_PROPORTIONALARRG` (JSON) = MATI** `[keputusan work owner]`: "sudah lama tidak dipakai;
  jika ada yang baca dari situ, kemungkinan lupa dihapus." → **dead code, jangan dimigrasikan.**
  Semua kueri `FROM m_PROPORTIONALARRG` (`GetMasterDescriptionEPIParentList`, `GetMasterPortfolioListDetail`,
  `GetMasterPanggilID`, `GetMasterDescriptionPLAParentList`, `GetMasterDescriptionExGratiaList`,
  `GetMasterDescriptionFACINParentList`, dll) = dead read → di sistem baru **baca dari
  `PROPORTIONALARRG`**. Migrasi TIDAK ambil dari JSON (tak ada data hidup di sana).
- **OQ tersisa:** DDL `PROPORTIONALARRG` (tipe kolom, PK, nullability); body 5 procedure master lain
  + PROSESCOPY.


---

## DDL + 4 procedure master lain — dari DBA (menutup OQ-001 & OQ-002)

Tanggal: 2026-09-15. `[data DBA]`. Modul ini **satu-satunya penulis** tabel-tabel ini `[keputusan work owner]`; modul lain (Claim Prop/Fac In/Komite) **read-only**.

### Tipe kolom (temuan mengikat)

| Tabel | Catatan tipe |
| --- | --- |
| `PROPORTIONALARRG` | ⚠️ **CAMPUR**: `TREATYLIMIT, COINS_MIN, COINS_MAX, MORERP, MOREUSD` = **NUMBER**; tapi `RP, USD, PCT, PCTME` = **VARCHAR2(1000)** (teks). Skema baru: **semua uang/persen → decimal** (ADR-0003). Ada `TGLUPDATE DATE`, `OBJECT VARCHAR2(50)`, `PROPORTIONALLIST VARCHAR2(1000)` — 2 terakhir **tidak di-set procedure** → kemungkinan dead kolom, catat |
| `MTREATYSECURITY` | ⚠️ kotor: `THN_TREATY VARCHAR2(4) DEFAULT '1' NOT NULL`, `TOP_ID VARCHAR2(9)`, `TP_TREATY CHAR(2)`, `REAS_ID CHAR(7) NOT NULL`, `PCT_SHARE VARCHAR2(99)`, `USER_ID CHAR(99)`, `REAS_SECURITY CHAR(10) NOT NULL`. Tanpa PK. INSERT posisional mengosongkan `TOP_ID, TP_TREATY, USER_ID`. Skema baru: **PK surrogate, tipe wajar, PCT_SHARE→decimal, REAS_SECURITY atribut biasa** (Q10) |
| `TREATYCONTRACT` | `ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME, TREATYSTARTDATE DATE, TREATYENDDATE DATE, USERID, TGLUPDATE VARCHAR2(1000)`. ⚠️ tanggal start/end sudah DATE, tapi TGLUPDATE teks |
| `TREATYYEAR` | `ID, TREATYYEAR, UNDERWRITINGYEAR, TREATYGROUPID, TREATYGROUPNAME, USERID, TGLUPDATE, PROPORTION, STARTDATE, ENDDATE` — semua VARCHAR2 (tanggal teks) |
| `TREATYREINSURER` | `... RICOMM NUMBER, PCTSHARE NUMBER ...` (uang angka ✅), sisanya VARCHAR2; `STDRATING` (field dipakai-ulang, sama seperti Retro Life) |
| `TREATYBUSINESS` | `ID, ISACTIVE, TREATYYEAR, TREATYYEARID, TREATYGROUPID, TREATYGROUPNAME, REINSTYPEID, REINSTYPENAME, BIZCODE, BIZNAME, USERID, TGLUPDATE` — semua VARCHAR2 |
| `TREATYEXCHANGEYEARLY` | tabel kurs (nama sebenarnya, bukan `TREATYEXCHANGE`): `TOIDR, TOUSD, IDCURRENCY, CURRENCY, QUARTER, STARTDATE, ENDDATE, ...` semua VARCHAR2. Kurs USD→IDR: `IDCURRENCY='10001'` (=USD), `QUARTER='0'` |
| `TREATYDESC` | master jenis klausul: `ID, DESCNAME, ISXOL, STATUSAKTIF` — semua VARCHAR2 |

### 4 procedure (PEGA_TREATYCONTRACT/YEAR/REINSURER/BUSINESS) — pola identik

- **UPSERT dikunci `ID`**; ID baru `'1'+lpad(seq,6)` (**6 digit**), seq: `treatycontract_seq`,
  `TreatyYear_seq`, `M_TREATYREINSURER_SEQ`, `TREATY_BUSINESS_SEQ`.
- **TIDAK commit sendiri** (rollback on error) → Go bisa atomik lintas tabel. ✅
- `StsSimpan` **1=sukses/0=gagal**; `ErrMsg` teks (copy-paste "JSON_KLAIM" — kejanggalan, bukan JSON).
- `PEGA_TREATYCONTRACT`: `TREATYSTARTDATE`/`ENDDATE` di-`to_date(..,'DD/MM/YYYY')` → tanggal.
- ⚠️ `PEGA_TREATYBUSINESS` UPDATE hanya set `ISACTIVE,BIZCODE,BIZNAME,USERID,TGLUPDATE` (tidak set
  REINSTYPEID dll saat update — hanya saat insert). Catat.

### Status OQ setelah DDL+body
- **OQ-002 ✅ DITUTUP** — 6 procedure penulis diterima (PEGA_PROPORTIONALARRG + CHILD + 4 master).
  PROSESCOPY tidak perlu (fitur copy dibuang, Q7).
- **OQ-001 ✅ DITUTUP** — DDL 8 tabel diterima; tipe/nullability/index diketahui.
- **Perbaikan sadar tipe:** semua uang/persen → decimal; semua tanggal → DATE; `MTREATYSECURITY` PK
  surrogate + tipe wajar; kolom dead (`PROPORTIONALLIST`, `OBJECT`) tidak dibawa kecuali dipastikan.

## `PEGA_M_ATTACHMENT` — rekonsiliasi lampiran (OQ-TCO-24 — **ditutup** lanjutan 5, 29-09-2026)

### Pasti dari korpus `[terverifikasi]`

- `InsertAtatchment_Sql` b60: `POOLDATA.PEGA_M_ATTACHMENT(IDPEGA, DATAPEGA CLOB, ERRMSG out, STSSAVE out)` lalu `COMMIT`.
- **Satu pemanggil hidup**: `TreatyOutSaveAttachment` b1674 (`pyStepsBlockName` kosong). `IDPEGA` = `TreatyYear + TreatyYearID`
  (b1402), `DATAPEGA` = `@ASM.GetPageJSONString()` halaman `TempData` (b1510), `CARI19 = 1` (b1355).
- `Delete_act` b1022 memanggil prosedur yang sama, tetapi langkahnya **ter-remark** (`pyStepsBlockName` `//` b1034; juga
  b783, b900), jadi tidak pernah jalan. Jalur hapus yang hidup: `DeleteGoogleStorage_Act` b396 lalu `DeleteAttachment2_Sql` b565.
- Pembaca dan penghapus: `GetAllAttachment2_Sql`, `GetAttachment2_Sql`, `DeleteAttachment2_Sql` membaca `M_ATTACHMENTTREATY_2`.
- Sensus seluruh korpus `RNM_BRD` (tanpa folder keluaran): **nol** RDB/activity yang menyebut `M_ATTACHMENTTREATY` tanpa `_2`,
  dan nol yang menyebut `ID_COUNT`. Kedua nama itu hanya dikenal dari badan prosedur (agregat brief lanjutan 4 §1:
  35 baris `ALL_SOURCE`).

### Badan prosedur `[data DBA — ALL_SOURCE POOLDATA, dibaca asisten, baca-saja; 35 baris, baris komentar dibuang]`

Hanya strukturnya yang dicatat di sini. Kepala komentar prosedur tidak disalin.

| Baris | Isi |
| --- | --- |
| 1 | `PEGA_M_ATTACHMENT(IDPega IN VARCHAR2, DataPega IN CLOB, ErrMsg OUT VARCHAR2, StsSimpan OUT NUMBER)` |
| 7 | `SELECT count(id) INTO id_count FROM pooldata.M_ATTACHMENTTREATY WHERE ID = IDPega` |
| 14–16 | bila ada: `UPDATE pooldata.M_ATTACHMENTTREATY SET DATA_JSON = DataPega WHERE ID = IDPega; COMMIT` |
| 24–26 | bila tidak: `INSERT INTO pooldata.M_ATTACHMENTTREATY (ID, DATA_JSON, DATEINPUT) VALUES (IDPega, DataPega, TO_DATE(to_char(sysdate,'dd/MM/yyyy'),'dd/MM/yyyy'))`, tanpa `COMMIT` di cabang ini (`COMMIT` datang dari RDB `InsertAtatchment_Sql` b60) |

- **`id_count` adalah variabel lokal, bukan tabel.** Ini meralat dugaan lanjutan 4 bahwa `ID_COUNT` sebuah tabel penghitung
  ID. Karena itu pula namanya nol di korpus: ia hanya hidup di dalam badan prosedur.
- Katalog: `M_ATTACHMENTTREATY` **dan** `M_ATTACHMENTTREATY_2` keduanya ada di `POOLDATA`.
- `M_ATTACHMENTTREATY` adalah **simpanan JSON halaman**: `DATA_JSON` CLOB per `ID`, di-upsert dan ditimpa tiap simpan.
  `DATEINPUT` = tanggal hari ini tanpa jam. Isinya menurut pemanggil: `ID` = `TreatyYear + TreatyYearID` (b1402), `DATA_JSON`
  = JSON halaman `TempData`, salinan halaman berkas `Embed-DragDropFile` (b683; `pyCategory` "File" b558).
- `M_ATTACHMENTTREATY_2` adalah tabel **relasional** yang dibaca dan dihapus seluruh RDB modul.

### Kesimpulan dan keputusan

**Temuan warisan — dicatat, tidak disatukan.** Di Pega, langkah hidup `TreatyOutSaveAttachment` hanya menulis simpanan JSON
`M_ATTACHMENTTREATY`. Langkah `Java` b1016 hanya menyegarkan clipboard `.pyAttachments` / `D_AttachmentList`, dan tidak ada
langkah hidup yang menulis `M_ATTACHMENTTREATY_2`. Padahal daftar dan hapus modul (`GetAllAttachment2_Sql`,
`GetAttachment2_Sql`, `DeleteAttachment2_Sql`) membaca `_2`. Jadi lampiran yang disimpan dari modul ini di Pega tidak
tampil di daftarnya sendiri, kecuali `_2` diisi penulis lain. Hal ini sejalan dengan penyimpangan sadar 9 (jalur lampiran
Pega belum rampung) dan dengan tiket 12 sebagai fitur baru.

**Keputusan `[asisten dari bukti; veto work owner]`** (brief lanjutan 5 §1):

- aplikasi **tetap menulis dan membaca `M_ATTACHMENTTREATY_2`**, dengan kolom dari `Treaty In/InsertAttachment2_Sql` b84;
- aplikasi **tidak** menulis simpanan JSON `M_ATTACHMENTTREATY`.

Keputusan ini konsisten dengan keputusan proyek bahwa simpanan JSON halaman Pega tidak diteruskan (spec §15 "Dualitas JSON
dibuang"; PremiumList: revisi penyimpanan JSON dibuang). Tidak ada pembaca hilir `M_ATTACHMENTTREATY` di korpus: sensus di
atas menemukan nol sebutan. Kode memakai stempel `YYYYMMDDHH24MISSFF3` WIB sebagai `ID` baris `_2`.

## `Update_T_Storage_SQL` — ditiru (OQ-TCO-26, lanjutan 4, 29-09-2026)

- Urutan `GetUrlGoogleStorage_Act`:
  - langkah 6 "JIKA EXPDATE SUDAH EXPIRED": `geturl`, lalu `UpdateDoc` (b2145–b2295);
  - `Update_T_Storage_SQL` b2427 (`pyStepsBlockName` kosong) mengisi `URLPUBLIC`, `APPFOLDER`,
    `EXPDATE = To_date(exp,'DD/MM/YYYY HH24:MI:SS')`, dan `TANGGAL_UPLOAD = To_date(DateTime,'MM/DD/YYYY HH24:MI:SS')`,
    `WHERE imageid`.
- Kode: `repository.PerbaruiObjek` (`sqlPerbaruiObjekTCO`, teks SQL sama, tanpa `COMMIT`) dijalankan oleh
  `PenyimpananJarakJauhTCO.segarkan`. Pemicunya: tiap `geturl` yang berhasil, yaitu `Buka`, dan `Ada` bila objeknya ada.
  - Penyegaran berjalan dalam transaksi pendeknya sendiri, karena pemanggil tidak memegang kunci baris itu saat `geturl`.
  - Bila gagal, kegagalan dicatat di log (IMAGEID dan sebab, tanpa URL) dan unduhan tetap jalan.
- **Penyimpangan sadar**: Pega memakai ulang `URLPUBLIC` tersimpan bila belum kedaluwarsa (transisi
  `UploadDoc.exp==""` b704), jadi `geturl` hanya dipanggil bila perlu. Kode memanggil `geturl` setiap unduh/periksa (URL
  bertanda tangan segar), lalu menyegarkan baris. Isi tabel sesudahnya sama: URL dan `exp` terakhir.
- `exp` diubah seperti Pega (b2146/b2211): `-`/`:` dibuang, dibaca sebagai GMT, ditulis `dd/MM/yyyy HH:mm:ss`
  (`models.ExpStorageTCO`). Bentuk tak terbaca menjadi NULL.
- **RALAT jalur unggah**: `exp` kini diubah dengan cara yang sama (`InsertGoogleStorage_Act` b2366/b2431). Sebelumnya `exp`
  diteruskan mentah ke `TO_DATE(...,'DD/MM/YYYY HH24:MI:SS')`, sehingga jawaban berbentuk ISO akan gagal dengan ORA-01861.

## Kolom pelaku `USERID`/`TGLUPDATE` (OQ-TCO-25, lanjutan 4, 29-09-2026)

- `SaveMasterTreatyReinsurer_SQL` b79 dan `SaveMasterTreatyBusiness_SQL` b85 meneruskan `UserId`/`TglUpdate` apa adanya.
- Pega tidak mengisi kolom itu:
  - reinsurer: `NewTreatyReinsurerDetail_Act` b917 mengosongkan `UserId`; `SetUbahTreatyReinsurerList_Act` b1097 hanya
    menyalin nilai dari barisnya; tidak ada langkah yang mengisi `TglUpdate`;
  - business: tidak ada langkah yang mengisi keduanya.
- Data DEV cocok: 0/430 dan 2/4.621 baris terisi.
- Kode: layanan tidak mengisi kedua kolom. Pelaku tercatat di log aplikasi
  (`treaty contract out: reinsurer|business <ID> disimpan oleh akun <akun>`) dan di `OPERATORNAME` reinsurer
  (`OperatorID.pyUserName` b377).
