# Spec — Treaty Contract Out (migrasi Pega → Go + React + Oracle)

Status: **ready-for-agent** — ✅ **OQ-001, OQ-002, OQ-020, OQ-042 DITUTUP**
Konteks: `treaty-contract-out` — **konteks non-Life pertama** `[keputusan work owner]`
Modul: **Treaty Contract Out** (303 berkas)
Tanggal: 2026-09-16
Sumber: `.scratch/treaty-contract-out/grilling-ronde-1-jawaban.md` (**12 verdict final**),
`.scratch/treaty-contract-out/dba-procedures.md` `[data DBA]`,
`.scratch/treaty-contract-out/grilling-ronde-1.md`, `discovery/modules/Treaty Contract Out.md`,
`discovery/flows/Treaty Contract Out.md`, `discovery/context-map.md` §2.4 & §3, `CONTEXT.md`,
`docs/adr/ADR-0001`–`ADR-0015`
Skill: `/mattpocock-skills:to-spec`

> ⛔ **RALAT BERTANGGAL 29-09-2026 — tco4 dan tco5 `[DIPUTUSKAN work owner]`** (brief lanjutan 3).
> **tco4** menggantikan tco1: modul ini **nol tabel baru** — menulis dan membaca tabel warisan `POOLDATA`
> (`TREATYYEAR`, `TREATYCONTRACT`, `TREATYREINSURER`, `MTREATYSECURITY`, `TREATYBUSINESS`, `PROPORTIONALARRG`,
> `M_ATTACHMENTTREATY_2`, `T_STORAGE_IMAGE`) dengan nama tabel/kolom VERBATIM, persis RDB XML; procedure tetap tidak
> dipanggil (keputusan o), logikanya ditiru, sequence warisan dipakai. Penyimpangan sadar "tipe dirapikan", "PK
> surrogate security", "tabel jejak", dan migrasi data gugur; peta tabel: `STRUKTUR-TABEL-TREATY-CONTRACT-OUT.md`.
> **tco5**: menu kelompok Treaty Contract Out SATU butir `Treaty Contract Out`; `InboxTreatyContractReinsType` dan
> `InboxTreatyContractDescription` adalah popup form kontrak (`InputTreatyContract.xml` b20778/b22196), bukan menu.
> Bagian spec di bawah yang menyebut tabel `T_*`, skema relasional baru, atau migrasi data tunduk pada ralat ini.

> ⛔ **RALAT BERTANGGAL 29-09-2026 — lanjutan 4 `[asisten dari data DEV; veto work owner]`.** Bentuk nilai mengikuti
> DATA warisan: `TREATYYEAR.STARTDATE/ENDDATE` ditulis dan dibaca `YYYYMMDD` (bentuk lain ditolak, OQ-TCO-01);
> `TREATYREINSURER.STARTDATE/ENDDATE` tidak ditulis; `USERID`/`TGLUPDATE` reinsurer dan business tidak diisi layanan
> (OQ-TCO-25; pelaku di log aplikasi); desimal teks bertitik (OQ-TCO-23); `Update_T_Storage_SQL` ditiru sesudah tiap
> geturl, `exp` diubah seperti Pega (OQ-TCO-26). Rekonsiliasi `PEGA_M_ATTACHMENT`: `dba-procedures.md` (OQ-TCO-24).

> **Konvensi penandaan.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[fakta bisnis — work owner]`; `[data DBA]`; `[terbuka]` = OQ.
> **Identitas rule wajib menyertakan class** — nama sama di class berbeda = rule berbeda.

> **Sumber tunggal.** Korpus Pega `D:\XML\RNM_BRD\` (READ-ONLY) dan artefak di
> `OUTPUT_HASIL_RNM\`, sekaligus repo target tunggal. ⛔ `D:\XML\nusantara-re\` di-blacklist.

> `[keputusan work owner]` **Nama modul menyesatkan.** `Treaty Contract Out` **bukan** modul treaty
> outward. Ia **editor master term/arrangement kontrak treaty non-life**. Lima uji berbukti ada di
> `discovery/flows/Treaty Contract Out.md` §2. Penetapan nama konteks di sistem baru = keputusan
> GATE (**OQ-022**), tidak memblokir spec ini.

---

## Problem Statement

Seluruh syarat sebuah kontrak treaty non-life — siapa reinsurernya, berapa sharenya, bisnis apa
yang ditanggung, dan dua puluh lima jenis **klausul** yang mengikat (limit, EPI, PLA, Ricomm,
profit commission, ex-gratia, cash loss limit, territorial limit, dan seterusnya) — hari ini dikelola
di satu layar Pega dan disimpan ke delapan tabel Oracle lewat enam stored procedure.

Tujuh hal membuat modul ini berisiko dipindahkan secara salah:

1. **Ada dua tabel kembar untuk hal yang sama, dan yang satu sudah mati.**
   `[terverifikasi]` Sebagian kueri membaca `PROPORTIONALARRG` (kolom datar), sebagian lagi membaca
   `m_PROPORTIONALARRG` dengan `a.JSONDATA.…` dan `JSON_TABLE`. Bahkan satu rule memegang keduanya:
   `RDBList/DeleteRowBusinessList.xml` (`ASM-FW-GISFW-INT-TREATYBUSINESS` /
   `ASM!DELETEROWBUSINESSLIST` / `RULE-CONNECT-SQL`) menghapus `treatybusiness` di `pyBrowseSQL`
   dan `m_treatybusiness` di `pyDeleteSQL`. `[data DBA]` Body procedure membuktikan **tidak satu pun
   procedure menulis ke tabel JSON** — jadi yang JSON **mati**, dan setiap kueri yang masih
   membacanya mengembalikan data basi.

2. **Nama-nama berbohong, berulang kali.** `[terverifikasi]` Rule bernama
   `BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull` justru yang **terbaru** (ruleset 01-01-91, commit
   2026-02-05) dan **dipakai 11 grid**; yang tanpa "Old" dipakai satu layar.
   `[data DBA]` Procedure bernama `PEGA_M_PROPORTIONALARRG_CHILD` **tidak** menulis ke tabel child —
   ia menulis ke tabel yang **sama**. Dan activity bernama `testingKurs` (`@BASECLASS!TESTINGKURS`)
   adalah jalur kurs yang **benar-benar dipakai**, sementara `SetTreatyArrangementDesc_Act`
   (`ASM-FW-GISFW-INT-PROPORTIONALARRG`) **lima dari enam langkahnya di-remark**.

3. **Uang dan tanggal disimpan sebagai teks.** `[data DBA]` Di `PROPORTIONALARRG`, `TREATYLIMIT`,
   `COINS_MIN`, `COINS_MAX`, `MORERP`, `MOREUSD` bertipe `NUMBER` — tetapi `RP`, `USD`, `PCT`,
   `PCTME` bertipe `VARCHAR2(1000)`. Di `TREATYYEAR` dan `TREATYBUSINESS`, **semua** kolom
   `VARCHAR2`, tanggal termasuk. Nilai uang dalam satu tabel diperlakukan dua cara berbeda.

4. **Satu tabel master ditulis dengan SQL mentah dan struktur yang kotor.**
   `[terverifikasi]` `RDBList/InsertToMTreatySecurity.xml` (`ASM-FW-GISFW-INT-MTREATYSECURITY` /
   `ASM!INSERTTOMTREATYSECURITY`) melakukan `insert into mtreatysecurity values (…)` **tanpa daftar
   kolom**, tujuh nilai berposisi, tiga di antaranya string kosong. Pembaruan dan penghapusannya
   berkunci `REAS_ID` + `trim(REAS_SECURITY)` — nama dipakai sebagai bagian kunci, dan `trim()`
   mengakui datanya bertabur spasi. `[data DBA]` DDL-nya **tanpa primary key**.

5. **Fitur salin tahun membawa angka tahun lalu ke tahun ini.**
   `[fakta bisnis — work owner]` `PROSESCOPY` menyalin isi satu tahun treaty ke tahun lain. Dalam
   praktiknya ia membawa nilai lama — misalnya batas QS tahun 2025 — ke tahun yang semestinya
   berbeda, dan itu menjadi sumber kesalahan, bukan penghemat waktu.

6. **Lampiran belum selesai dibangun.** `[keputusan work owner]` Jalur `M_ATTACHMENTTREATY_2` di
   Pega **belum rampung di-develop**; di sistem baru fungsinya perlu **dilengkapi**, bukan sekadar
   dipindahkan.

7. **Klausul tidak punya foreign key ke kontrak.** `[fakta bisnis — work owner]`
   `PROPORTIONALARRG` menggantung pada **(TreatyYear, TreatyGroupID, ReinsTypeID)**, bukan pada
   `ID` kontrak. Ini bukan kelalaian: klausul memang **milik level tahun/grup/jenis reasuransi**
   dan dipakai bersama. Migrasi yang "merapikan" dengan menambahkan FK ke kontrak akan **mengubah
   arti data**.

## Solution

Membangun ulang Treaty Contract Out sebagai **editor master arrangement kontrak treaty non-life**
di atas Go + React + Oracle — **editor master murni**: tidak ada tangga persetujuan, tidak ada
Submit/Decline, tidak ada status akseptasi. `[terverifikasi]` Nol rule `Flow`, nol rule `When`, nol
`StatusAkseptasi` di seluruh 303 berkas.

### Sembilan penyimpangan sadar

| # | Penyimpangan | Alasan | Sumber |
| --- | --- | --- | --- |
| 1 | **Dualitas JSON dibuang** — satu sumber kebenaran `PROPORTIONALARRG` relasional | `[data DBA]` tidak ada procedure yang menulis ke tabel JSON; ia sudah lama tidak dipakai | `[keputusan work owner]` |
| 2 | **SATU tabel generik untuk 25 jenis klausul**; baris "child" = subset kolom (9 kolom induk NULL) | `[data DBA]` kedua procedure menulis ke tabel yang sama; bedanya hanya jumlah kolom | `[keputusan work owner]` |
| 3 | **Fitur salin tahun treaty DIBUANG** | menyalin membawa nilai tahun lalu ke tahun yang semestinya berbeda | `[fakta bisnis — work owner]` |
| 4 | **Kaskade hapus 3 anak + popup konfirmasi**; **klausul sengaja dikecualikan** | klausul milik level tahun/grup/jenis, bukan milik kontrak | `[keputusan work owner]` |
| 5 | **`MTREATYSECURITY` dibersihkan** — PK surrogate, kolom bernama, `REAS_SECURITY` jadi atribut biasa | INSERT posisional + kunci `trim()` adalah kelas bug, bukan aturan bisnis | `[keputusan work owner]` |
| 6 | **Semua uang & persen → desimal; semua tanggal → `DATE`** | existing menyimpan sebagian sebagai `VARCHAR2` | `[keputusan work owner]`, **ADR-0003** |
| 7 | **`IDCURRENCY='10001'` (USD) jadi rujukan master**, bukan literal di kode | konversi USD→IDR memang aturan bisnis, tetapi identitas mata uang bukan konstanta program | `[keputusan work owner]` |
| 8 | **Nama jujur di sistem baru** — buang `"Old"`, `"testing"`, dan teks galat `"JSON_KLAIM"` | nama yang berbohong merambat ke sistem baru bila ditiru | `[keputusan work owner]` |
| 9 | **FITUR BARU: unggah lampiran di level tahun treaty** | di Pega belum rampung di-develop | `[keputusan work owner]` |

### Kontrak batas

| Arah | Isi |
| --- | --- |
| **Keluar** | Modul ini adalah **satu-satunya penulis** `TREATYREINSURER`, `TREATYBUSINESS`, `PROPORTIONALARRG`, `TREATYCONTRACT`, `TREATYYEAR`, `MTREATYSECURITY`. `Claim Prop`, `Komite Claim Prop`, `Claim Fac In` **hanya membaca** (**OQ-042 terkonfirmasi**) |
| **Masuk** | Master yang **dibaca saja**: `REINSURANCETYPE` (jenis reasuransi), `TREATYDESC` (jenis klausul), `TREATYGROUP`, mata uang, `TREATYEXCHANGEYEARLY` (kurs), `CATEGORY_ATTACH_REAS` (kategori lampiran); berkas lampiran ke **Google Storage** |
| **Tidak ada** | ⚠️ Modul ini **tidak** menyentuh objek treaty **outward** sama sekali `[terverifikasi]` — nol berkas. Pekerjaan outward yang sesungguhnya ada di `NB Treaty In` dan `Claim Non Prop` |

---

## User Stories

### Tahun treaty dan kontrak

1. Sebagai **admin master treaty**, saya ingin membuat **tahun treaty** beserta grup treaty,
   underwriting year, proporsi, dan masa berlakunya, supaya seluruh kontrak tahun itu punya wadah.
2. Sebagai **admin master**, saya ingin **tidak perlu menentukan nomor identitas** tahun treaty
   maupun kontrak, supaya penomoran tidak pernah bentrok antar petugas.
3. Sebagai **admin master**, saya ingin identitas tetap berbentuk seperti yang sudah dikenal
   perusahaan, supaya rujukan lama tetap terbaca.
4. Sebagai **admin master**, saya ingin membuat **kontrak treaty** di dalam sebuah tahun treaty,
   dengan jenis reasuransinya dan masa berlakunya sendiri.
5. Sebagai **admin master**, saya ingin mengubah tahun atau kontrak yang sudah ada **tanpa membuat
   duplikat**.
6. Sebagai **admin master**, saya ingin **diberi tahu bila penyimpanan gagal**, supaya saya tidak
   mengira data tersimpan padahal tidak.
7. Sebagai **underwriter**, saya ingin masa berlaku kontrak tidak boleh berakhir sebelum ia mulai.

### Jenis reasuransi

8. Sebagai **admin master**, saya ingin memilih **jenis reasuransi** dari daftar master, supaya
   tidak ada yang diketik bebas.
9. Sebagai **admin master**, saya ingin daftar jenis reasuransi yang muncul **hanya berisi jenis
   yang berlaku untuk non-life**, supaya saya tidak salah pilih jenis milik lini lain.
10. Sebagai **organisasi**, saya ingin master jenis reasuransi **dikelola di tempatnya sendiri**,
    supaya satu daftar melayani seluruh lini.

### Reinsurer dan security

11. Sebagai **admin master**, saya ingin mencatat **para reinsurer** pada sebuah kombinasi tahun,
    grup, dan jenis reasuransi, beserta **share** dan **komisi reasuransi** masing-masing.
12. Sebagai **underwriter**, saya ingin melihat **total share** seluruh reinsurer, supaya saya tahu
    apakah penempatan sudah penuh.
13. Sebagai **admin master**, saya ingin mencatat **security** di bawah seorang reinsurer beserta
    porsinya, supaya eksposur berjenjang terlihat.
14. Sebagai **organisasi**, saya ingin baris security punya **identitas sendiri**, supaya mengubah
    nama security tidak memutus rujukannya.

### Business

15. Sebagai **admin master**, saya ingin mencatat **jenis bisnis** yang ditanggung sebuah kontrak,
    dengan kode dan namanya.
16. Sebagai **admin master**, saya ingin **menonaktifkan** satu baris bisnis tanpa menghapusnya.

### Klausul kontrak

17. Sebagai **underwriter**, saya ingin mengelola **dua puluh lima jenis klausul** kontrak, supaya
    seluruh syarat treaty tercatat di satu tempat.
18. Sebagai **underwriter**, saya ingin daftar jenis klausul datang dari **masternya sendiri**,
    supaya jenis baru dapat ditambahkan tanpa mengubah aplikasi.
19. Sebagai **underwriter**, saya ingin sebagian jenis klausul punya **baris rincian di bawahnya**,
    supaya klausul berlapis dapat dinyatakan.
20. Sebagai **underwriter**, saya ingin **setiap jenis klausul memeriksa field yang memang relevan
    baginya**, bukan satu daftar wajib yang seragam.
21. Sebagai **underwriter**, saya ingin **membatalkan** pengeditan satu klausul tanpa kehilangan
    klausul lain yang sedang saya kerjakan.
22. Sebagai **underwriter**, saya ingin sistem **menolak baris klausul yang sudah pernah saya input**,
    supaya tidak ada duplikat diam-diam.
23. Sebagai **underwriter**, saya ingin klausul **bertahan ketika kontrak dihapus**, karena klausul
    itu milik tahun/grup/jenis reasuransi, bukan milik satu kontrak.

### Menyimpan

24. Sebagai **organisasi**, saya ingin seluruh perubahan pada satu kontrak — reinsurer, security,
    business, dan klausul — tersimpan **bersama atau tidak sama sekali**.
25. Sebagai **admin master**, saya ingin kegagalan penyimpanan **terlihat terang-terangan** dengan
    pesan yang menyebut apa yang gagal.

### Menghapus

26. Sebagai **admin master**, saya ingin menghapus sebuah kontrak **beserta** reinsurer, security,
    dan business-nya, supaya tidak ada sisa yang menggantung.
27. Sebagai **admin master**, saya ingin **diberi peringatan berisi jumlah baris yang akan ikut
    terhapus** sebelum penghapusan dijalankan, supaya saya dapat membatalkan.

### Kurs

28. Sebagai **underwriter**, saya ingin nilai dalam **USD** dapat dilihat padanannya dalam **IDR**
    menurut kurs yang berlaku pada periode kontrak.
29. Sebagai **Finance**, saya ingin **tidak ada nilai uang yang berubah** saat menyeberang batas
    penyimpanan.
30. Sebagai **Finance**, saya ingin nilai **Rp** dan **USD** tetap tercatat sebagai dua nilai
    terpisah, karena keduanya memang dua angka yang berbeda.

### Lampiran (fitur baru)

31. Sebagai **admin master**, saya ingin **melampirkan berkas pada sebuah tahun treaty**, supaya
    dokumen kontraknya tersimpan bersama datanya.
32. Sebagai **admin master**, saya ingin melampirkan berkas **tanpa itu menjadi syarat** tersimpannya
    data tahun treaty.
33. Sebagai **admin master**, saya ingin memberi **kategori** pada tiap lampiran.
34. Sebagai **admin master**, saya ingin **mengunduh** dan **menghapus** lampiran.
35. Sebagai **tim operasi**, saya ingin **kegagalan unggah terlihat dan dapat diulang**, supaya tidak
    ada berkas yang hilang diam-diam.
36. Sebagai **tim operasi**, saya ingin alamat penyimpanan berkas **dibaca dari konfigurasi saat
    dijalankan**, supaya perpindahan lingkungan tidak menuntut penempelan ulang.

### Migrasi

37. Sebagai **tim migrasi**, saya ingin seluruh data master arrangement pindah **tanpa kehilangan
    satu nilai pun**.
38. Sebagai **tim migrasi**, saya ingin nilai uang dan tanggal yang hari ini tersimpan sebagai teks
    menjadi **tipe yang semestinya**.
39. Sebagai **tim migrasi**, saya ingin **tidak mengambil apa pun dari tabel JSON yang sudah mati**.
40. Sebagai **konteks hilir** (`Claim Prop`, `Komite Claim Prop`, `Claim Fac In`), saya ingin tetap
    dapat **membaca** master arrangement dalam bentuk yang saya kenal sampai saya ikut bermigrasi.

---

## Implementation Decisions

### 1. Batas konteks dan kepemilikan tulis

`[keputusan work owner]` Modul ini **satu-satunya penulis** enam tabel master arrangement. Konteks
hilir **read-only**. `[terverifikasi]` `discovery/context-map.md` §3 mencatat `Claim Prop`,
`Komite Claim Prop`, dan `Claim Fac In` membaca `TREATYREINSURER`, `TREATYBUSINESS`, dan
`PROPORTIONALARRG` — **OQ-042 terkonfirmasi**.

Konsekuensi mengikat: sampai ketiga konteks itu ikut bermigrasi, skema baru **wajib menyediakan
bentuk relasional yang mereka baca hari ini**. Perubahan nama kolom yang mereka pakai adalah
perubahan kontrak lintas konteks, bukan perapian internal.

### 2. Entitas dan hierarki `[data DBA]`

```
TREATYYEAR   (ID, TREATYYEAR, UNDERWRITINGYEAR, TREATYGROUPID, TREATYGROUPNAME,
              PROPORTION, STARTDATE, ENDDATE, USERID, TGLUPDATE)
   └─ TREATYCONTRACT (ID, IDTREATYYEAR → TREATYYEAR.ID, REINSTYPEID, REINSTYPENAME,
                      TREATYSTARTDATE, TREATYENDDATE, USERID, TGLUPDATE)

Digantung pada kunci gabungan (TREATYYEAR, TREATYGROUPID, REINSTYPEID) — BUKAN pada ID kontrak:
   ├─ TREATYREINSURER    (RICOMM, PCTSHARE, STDRATING, STARTDATE, ENDDATE, STATUSON, …)
   │     └─ MTREATYSECURITY (REAS_ID → TREATYREINSURER.ID, PCT_SHARE, REAS_SECURITY)
   ├─ TREATYBUSINESS     (BIZCODE, BIZNAME, ISACTIVE)  [+ TREATYYEARID]
   └─ PROPORTIONALARRG   (klausul; + TREATYDESCID, TREATYYEARID, PARENTREINSTYPEID)
```

⚠️ **Bedakan tegas dari Master Contract Retro Life.** Di konteks Life, anak-anak menggantung pada
**ID kontrak**. Di sini pada **kunci gabungan**. Itulah sebabnya kaskade hapus existing
(`RDBList/DeleteFromTREATYCONTRACT_SQL.xml`) harus mengulang seluruh kunci gabungan, bukan satu `ID`.
`[keputusan work owner]` **Bentuk ini dipertahankan** — menambahkan FK ke kontrak akan mengubah arti
data.

Master yang **dibaca saja**: `REINSURANCETYPE`, `TREATYDESC` (`ID`, `DESCNAME`, `ISXOL`,
`STATUSAKTIF`), `TREATYGROUP`, mata uang, `TREATYEXCHANGEYEARLY`, `CATEGORY_ATTACH_REAS`.

### 3. ⚠️ Satu tabel untuk dua puluh lima jenis klausul (penyimpangan sadar 2)

`[terverifikasi]` Sensus penuh atas 25 activity `SaveTreatyArr*` di
`Treaty Contract Out/Activity/`: **18** memanggil `RDBList/SaveMasterProportionalArrg.xml`
(→ `POOLDATA.PEGA_PROPORTIONALARRG`), **7** memanggil
`RDBList/SaveMasterProportionalArrgChild.xml` (→ `POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD`).
Mekanismenya: tiap klausul punya halaman masukannya sendiri, lalu langkah `Page-Copy` menyalinnya
ke halaman bersama `InputTreatyArrTreatyLimit` / `InputTreatyArrTreatyLimitChild` sebelum langkah
`RDB-List` memanggil procedure.

| | Jenis klausul |
| --- | --- |
| **Induk (18)** | `BordereAux`, `CashLossLimit`, `ClaimCoorp`, `CoinsPanel`, `EPI`, `ExGratia`, `ExclutionTreaty`, `FacIn`, `LimitMB`, `MaxCoinsPanel`, `MinLOL`, `MinLOLMB`, `PLA`, `Portfolio`, `ProfitComm`, `Ricomm`, `TerrLimit`, `TreatyLimit` |
| **Anak (7)** | `CashLossLimitList`, `ClaimCoorpChild`, `EpiList`, `ExGratiaChildList`, `FacInList`, `PLAList`, `TreatyLimitChild` |

`[data DBA]` **Kedua procedure menulis ke satu tabel yang sama: `POOLDATA.PROPORTIONALARRG`.**
Induk mengisi **35 kolom**; "anak" mengisi **26** — tanpa sembilan kolom khusus induk:
`ID_OCCUPATION`, `OCCUPATION`, `ID_CLAUSE`, `CLAUSE`, `TREATYLIMIT`, `COINS_MIN`, `COINS_MAX`,
`MORERP`, `MOREUSD`.

**Keputusan** `[keputusan work owner]`: satu tabel `proportionalarrg` di skema baru. Baris "anak"
= sembilan kolom itu **NULL**. Jenis klausul dibedakan **`TREATYDESCID`** (+ `TREATYDESCNAME`).
**Bukan** dua puluh lima tabel, **bukan** dua tabel induk/anak.

Kolom `PROPORTIONALARRG` `[data DBA]` — daftar lengkap dari body procedure:

```
ID, TREATYYEAR, TREATYYEARID, TREATYGROUPID, TREATYGROUPNAME, TREATYDESCID, TREATYDESCNAME,
REINSTYPEID, REINSTYPENAME, LAYER, LAYERPART, LAYERPARTTYPE, LAYERTYPE, KURS, TGLUPDATE, USERID,
LINE, PCT, PCTME, YDCF, METHOD, TERRITORIALLIMIT, PARENTREINSTYPEID, SPREADINGORDER, RP, USD,
ID_OCCUPATION, OCCUPATION, ID_CLAUSE, CLAUSE, TREATYLIMIT, COINS_MIN, COINS_MAX, MORERP, MOREUSD
```

`[terbuka]` Kolom `PROPORTIONALLIST` dan `OBJECT` ada di DDL tetapi **tidak di-set procedure mana
pun** — **jangan dibawa** ke skema baru kecuali migrasi membuktikan ada data hidup di sana.

### 4. `[keputusan work owner]` Istilah diikuti apa adanya

Nama klausul (`EPI`, `PLA`, `Ricomm`, `MB`, `LOL`, `BordereAux`, `ExGratia`, `CashLossLimit`,
`ClaimCoorperation`, `ProfitCommision`, `pTerrLimit`, `Portfolio`, `CoinsPanel`, `ExclutionTreaty`)
dan nama field (`Ydcf`, `PctMe`, `LayerPartType`, `LayerType`, `SpreadingOrder`, `Method`, `Line`)
**dipakai apa adanya** — diambil dari teks tampilan Pega, **tidak diterjemahkan dan tidak ditebak**.

`[terbuka]` Arti bisnis tiap istilah = **OQ terbuka untuk Product + UW**. Tidak memblokir: layar dan
skema dapat dibangun dengan istilah asli.

### 5. Validasi wajib-isi — berbeda per jenis klausul, ditegakkan di Go

`[terverifikasi]` Sensus prasyarat langkah (`<pyStepsPreCondParamsWhen> … ==""`) atas 25 activity:

| Jenis klausul | Field wajib |
| --- | --- |
| `CashLossLimit`, `ClaimCoorp`, `EPI`, `ExGratia`, `FacIn`, `PLA`, `TreatyLimit` | `ReinsTypeID`, `Rp`, `Usd` |
| `CashLossLimitList`, `ClaimCoorpChild`, `EpiList`, `ExGratiaChildList`, `FacInList`, `PLAList`, `TreatyLimitChild` | `Pct`, `ReinsTypeID`, `Rp`, `Usd` |
| `ProfitComm` | `Pct`, `PctMe`, `ReinsTypeID` |
| `Ricomm` | `Method`, `Pct`, `ReinsTypeID` |
| `CoinsPanel` | `TerritorialLimit`, `TreatyLimit` |
| `MaxCoinsPanel` | `CoIns_Max`, `TerritorialLimit` |
| `MinLOL`, `MinLOLMB` | `Pct`, `TerritorialLimit` |
| `TerrLimit` | `TerritorialLimit` |
| `ExclutionTreaty` | `ID_Occupation`, `TerritorialLimit` |
| `BordereAux` | `Method` |
| **`LimitMB`** | `[terbuka]` — hari ini **nol validasi** |
| **`Portfolio`** | `[terbuka]` — hari ini **nol validasi** |

`[keputusan work owner]` **Perbedaan dipertahankan** — tiap klausul memang meminta hal berbeda.
Pemetaan jenis → aturan wajib-isi berada di **kode Go**, bukan di data. Dua baris terakhir tetap
**OQ terbuka**, bukan "tanpa validasi": aturannya ditetapkan Product + UW sebelum kedua klausul itu
dinyatakan selesai.

`[terverifikasi]` Penanda baris baru di Pega: `InputTreatyArr….ID = "UnknownId"`. Anti-dobel memakai
pesan `OutputParam.ERRMSG4 = "Data sudah pernah di Input"`
(`Activity/SaveTreatyArrEPI_Act.xml`, `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `SAVETREATYARREPI_ACT`).

### 6. ⚠️ Simpan atomik — procedure **tidak** commit sendiri (penyimpangan yang menguntungkan)

`[data DBA]` Keenam procedure penulis (`PEGA_PROPORTIONALARRG`, `PEGA_M_PROPORTIONALARRG_CHILD`,
`PEGA_TREATYCONTRACT`, `PEGA_TREATYYEAR`, `PEGA_TREATYREINSURER`, `PEGA_TREATYBUSINESS`) berpola
identik:

| Aspek | Perilaku |
| --- | --- |
| Mode | **UPSERT dikunci `ID`** (`SELECT COUNT(1) WHERE ID = P_ID` → `UPDATE` / `INSERT`) |
| Identitas baru | dibuat di basis data: `'1' \|\| lpad(<seq>.nextval, N, '0')` |
| `TGLUPDATE` | diisi `SYSDATE` |
| **Transaksi** | **TIDAK commit sendiri** — hanya `ROLLBACK` bila galat; commit diserahkan pemanggil |
| Keluaran | `StsSimpan` **1 = sukses / 0 = gagal**; `ErrMsg` teks |

⚠️ Ini **berbeda dari lima konteks Life sebelumnya**, yang procedure-nya `COMMIT` sendiri dan
karenanya menjadi "titik potong" transaksi. Di sini **tidak ada titik potong**: Go **dapat dan
harus** membungkus seluruh perubahan satu kontrak — kontrak, reinsurer, security, business, dan
seluruh klausul — dalam **satu transaksi**, lalu commit sekali.

⚠️ `[data DBA]` Teks galat sebagian procedure menyebut `"JSON_KLAIM"`. Itu **sisa salin-tempel
template**, bukan petunjuk bahwa data berbentuk JSON. Pesan galat di sistem baru **jangan meniru
teksnya**.

### 7. Penomoran identitas via sequence `[data DBA]` — **ADR-0006**

| Tabel | Sequence | Format |
| --- | --- | --- |
| `PROPORTIONALARRG` | `PROPORTIONALARRG_SEQ` | `'1' + lpad(seq, 7, '0')` — **7 digit** |
| `TREATYCONTRACT` | `treatycontract_seq` | `'1' + lpad(seq, 6, '0')` — **6 digit** |
| `TREATYYEAR` | `TreatyYear_seq` | `'1' + lpad(seq, 6, '0')` |
| `TREATYREINSURER` | `M_TREATYREINSURER_SEQ` | `'1' + lpad(seq, 6, '0')` |
| `TREATYBUSINESS` | `TREATY_BUSINESS_SEQ` | `'1' + lpad(seq, 6, '0')` |

Pengguna **tidak pernah** mengetik identitas. Bentuknya dipertahankan supaya rujukan lama terbaca.

### 8. ⚠️ Kaskade hapus + popup konfirmasi — klausul **sengaja dikecualikan** (penyimpangan sadar 4)

`[terverifikasi]` Existing (`RDBList/DeleteFromTREATYCONTRACT_SQL.xml`) menghapus empat tabel
berurut dalam SQL mentah lalu `COMMIT`: `treatycontract` → `treatybusiness` → `MTREATYSECURITY`
(lewat subselect atas `TREATYREINSURER`) → `TREATYREINSURER`. **`PROPORTIONALARRG` tidak ikut.**

`[fakta bisnis — work owner]` Ketidakikutan itu **desain, bukan bug**: klausul menggantung pada
`(TreatyYear, TreatyGroupID, ReinsTypeID)` dan **milik level tahun/grup/jenis**, dipakai bersama
lintas kontrak.

**Keputusan** `[keputusan work owner]`:
- Menghapus kontrak **mengkaskade** ke `TREATYBUSINESS`, `MTREATYSECURITY`, dan `TREATYREINSURER`.
- Didahului **popup konfirmasi Ya/Batal** yang **menyebut jumlah baris tiap jenis** yang akan ikut
  terhapus (pola yang sudah ditetapkan di Master Contract Retro Life).
- **`PROPORTIONALARRG` tetap hidup.**
- Seluruh kaskade berjalan dalam **satu transaksi**.

⚠️ Existing juga **tidak konsisten** dalam hal tabel kembar: kaskade ini menghapus `treatycontract`
tetapi bukan `M_TREATYCONTRACT`, sedangkan `DeleteRowBusinessList` menghapus **keduanya**.
Inkonsistensi itu **lenyap dengan sendirinya** begitu penyimpangan 1 diterapkan.

### 9. ⚠️ `MTREATYSECURITY` dibersihkan (penyimpangan sadar 5)

`[data DBA]` DDL sekarang:
`THN_TREATY VARCHAR2(4) DEFAULT '1' NOT NULL`, `TOP_ID VARCHAR2(9)`, `TP_TREATY CHAR(2)`,
`REAS_ID CHAR(7) NOT NULL`, `PCT_SHARE VARCHAR2(99)`, `USER_ID CHAR(99)`,
`REAS_SECURITY CHAR(10) NOT NULL` — **tanpa primary key**.
`[terverifikasi]` INSERT posisional mengosongkan `TOP_ID`, `TP_TREATY`, `USER_ID`; update dan delete
berkunci `REAS_ID` + `trim(REAS_SECURITY)`.

**Skema baru** `[keputusan work owner]`: tabel anak `TREATYREINSURER` yang wajar — **primary key
surrogate**, seluruh kolom bernama dan diisi eksplisit, `PCT_SHARE` **desimal**, `REAS_SECURITY`
**atribut biasa** (bukan bagian kunci). `trim()` di kunci **tidak dibawa**.

### 10. ⚠️ Kurs USD → IDR (penyimpangan sadar 7)

`[data DBA]` Tabel kurs bernama **`TREATYEXCHANGEYEARLY`** (bukan `TREATYEXCHANGE`), kolom `TOIDR`,
`TOUSD`, `IDCURRENCY`, `CURRENCY`, `QUARTER`, `STARTDATE`, `ENDDATE` — **seluruhnya `VARCHAR2`**.

`[terverifikasi]` Jalur yang **benar-benar dipakai** adalah activity `testingKurs`
(`@BASECLASS!TESTINGKURS`, `Treaty Contract Out/Activity/testingKurs.xml`, nol langkah di-remark) —
**meskipun namanya "testing"**. Sebaliknya `SetTreatyArrangementDesc_Act`
(`ASM-FW-GISFW-INT-PROPORTIONALARRG`) **lima dari enam langkahnya di-remark** → mati.
⚠️ **OQ-066 terbalik**: yang bernama "testing" hidup, yang bernama wajar mati.

Aturan `[keputusan work owner]`:
- Kurs dicari menurut **periode** (tanggal kontrak berada di antara `STARTDATE` dan `ENDDATE`).
- **`IDCURRENCY = '10001'` berarti USD**, dan konversi **USD → IDR** memang aturan bisnis.
  Identitasnya menjadi **rujukan ke master mata uang**, **bukan literal di kode**.
- **`QUARTER = '0'` diikuti apa adanya**; artinya `[terbuka]` — OQ kecil, tidak memblokir.
- **`Rp` dan `Usd` tetap dua kolom terpisah** — keduanya memang dua angka berbeda, bukan satu nilai
  dengan kode mata uang.

### 11. Jenis reasuransi — master bersama, filter diikuti apa adanya

`[terverifikasi]` Class `ASM-FW-GISFW-INT-REINSURANCETYPE` adalah **class yang sama** dengan master
jenis reasuransi Life (`Master Contract Retro Life/ReportDefinition/BrowseReinsuranceTypeLimit_RD.xml`).
Jadi jenis reasuransi adalah **master bersama life & non-life**.

⚠️ Dua rule di class itu, dan **yang bernama "Old" justru yang terbaru dan dominan**:

| Rule (`pxInsName`) | RuleSet | Commit | Dipakai |
| --- | --- | --- | ---: |
| `…!BROWSEREINSURANCETYPE_RD_OLD_LJT_ID_ISNOTNULL` | **01-01-91** | **2026-02-05** | **11** grid |
| `…!BROWSEREINSURANCETYPE_RD` | 01-01-83 | 2025-06-02 | 1 layar |

`[keputusan work owner]` Yang **dimigrasikan** adalah perilaku yang dipakai 11 grid, **dan diberi
nama jujur** (kata "Old" dibuang):

```
.ID NOT IN ("10004","10011","10012","10021","10022","10025","10026","10028",
            "10248","10249","10018","10217")
AND .Flag = "active"
AND .Type IN ("1","2","3")
```

`[keputusan work owner]` **Blacklist dua belas ID itu ditiru persis — masih benar, jangan
digeneralkan.** Rule satunya adalah perilaku **layar master jenis reasuransi**, yang berada di luar
konteks ini (§Out of Scope).

⚠️ Catat beda ejaan antar konteks: `.Flag` bernilai **`1`** di Life, **`"active"`** di sini.

### 12. Jenis klausul — master `TREATYDESC`, dibaca saja

`[terverifikasi]` `ReportDefinition/BrowseTreatyDesc_RD.xml` (`ASM-FW-GISFW-INT-TREATYDESC` /
`BROWSETREATYDESC_RD`): field `.ID`, `.DescName`, `.IsXOL`, `.StatusAktif`; satu filter `.IsXOL`.
`TreatyDescID` **tidak pernah literal** di modul ini — selalu datang dari parameter.

`[data DBA]` DDL `TREATYDESC(ID, DESCNAME, ISXOL, STATUSAKTIF)`.
`[keputusan work owner]` Master **terpisah**, di-input di layarnya sendiri; konteks ini **baca saja**.

### 13. ⚠️ FITUR BARU — lampiran di level tahun treaty (penyimpangan sadar 9)

`[keputusan work owner]` Jalur `M_ATTACHMENTTREATY_2` di Pega **belum rampung di-develop**. Di
sistem baru fitur ini **dibangun sampai selesai**, melekat pada **tahun treaty**.

`[terverifikasi]` Rantai berkas yang ada sekarang, seluruhnya berclass
`ASM-FW-GISFW-INT-T_STORAGE_IMAGE` kecuali yang disebut lain:
`ConnectREST/ServiceGoogle.xml` (tanpa URL literal; alamat dari `SystemSettings/LinkService.xml`
→ `LINKSERVICE!LINKSERVICE`), `RDBList/GetTokenStorage_SQL.xml` (→ `POOLDATA.GET_TOKEN_STORAGE`),
`RDBList/GetLinkStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml`,
`RDBList/DeleteStorage_SQL.xml`, `RDBList/InsertAtatchment_Sql.xml` (→ `POOLDATA.PEGA_M_ATTACHMENT`,
CLOB), `RDBList/CategoryAttach_SQL.xml` (→ `CATEGORY_ATTACH_REAS`).

⚠️ `pzOriginalInstanceKey` `ServiceGoogle` menunjuk `…GOOGLESTORAGE_UPLOAD` — sekali lagi nama
berkas ≠ nama asal.

Aturan: **lampiran opsional** — kegagalannya **tidak** membatalkan tersimpannya tahun treaty; alamat
penyimpanan **di-resolve runtime** (**ADR-0013**, `[terbuka]` **OQ-047** tidak memblokir); kegagalan
efek keluar **terlihat dan dapat diulang** (**ADR-0015**); berkas tetap di **Google Storage**
(**ADR-0010**).

⚠️ **Taruhan berbeda dari Komite Claim Life.** Di sana efek keluar **wajib berhasil** (transactional
outbox). Di sini lampiran **opsional** — yang wajib adalah kegagalannya **terlihat**, bukan berhasil.

### 14. ⚠️ Tipe data dirapikan (penyimpangan sadar 6)

`[data DBA]` Keadaan sekarang — **tidak konsisten di dalam satu tabel**:

| Tabel | Temuan tipe |
| --- | --- |
| `PROPORTIONALARRG` | `TREATYLIMIT`, `COINS_MIN`, `COINS_MAX`, `MORERP`, `MOREUSD` = **NUMBER**; tetapi `RP`, `USD`, `PCT`, `PCTME` = **`VARCHAR2(1000)`** |
| `TREATYCONTRACT` | `TREATYSTARTDATE`/`TREATYENDDATE` sudah `DATE`; tetapi `TGLUPDATE` = `VARCHAR2(1000)` |
| `TREATYYEAR` | **seluruh** kolom `VARCHAR2` — termasuk `STARTDATE`, `ENDDATE` |
| `TREATYBUSINESS` | **seluruh** kolom `VARCHAR2` |
| `TREATYREINSURER` | `RICOMM`, `PCTSHARE` = `NUMBER`; sisanya `VARCHAR2` |
| `TREATYEXCHANGEYEARLY` | **seluruh** kolom `VARCHAR2` |

**Skema baru** `[keputusan work owner]`: **seluruh uang dan persen → desimal presisi arbitrer**
(**ADR-0003**, tidak pernah melewati `float`); **seluruh tanggal → `DATE`**.

⚠️ `[data DBA]` `PEGA_TREATYBUSINESS` saat **UPDATE** hanya mengisi `ISACTIVE`, `BIZCODE`,
`BIZNAME`, `USERID`, `TGLUPDATE` — `REINSTYPEID` dan kawan-kawan hanya diisi saat **INSERT**.
Perilaku ini dicatat; di sistem baru pembaruan **memperbarui seluruh field yang dikirim**.

### 15. ⚠️ Dualitas JSON dibuang (penyimpangan sadar 1)

`[data DBA]` **Tidak satu pun dari enam procedure penulis menyentuh tabel JSON.** `[keputusan work
owner]` `M_PROPORTIONALARRG` **sudah lama tidak dipakai**; kueri yang masih membacanya adalah sisa
yang lupa dihapus.

**Keputusan:** sumber kebenaran tunggal = **`PROPORTIONALARRG` relasional**. Setiap kueri yang hari
ini berbunyi `FROM m_PROPORTIONALARRG` **dibaca ulang dari `PROPORTIONALARRG`** di sistem baru.
**Migrasi tidak mengambil apa pun dari tabel JSON.**

### 16. Migrasi data — **ready**

Bentuk skema berubah (tipe dirapikan, `MTREATYSECURITY` diberi PK, JSON dibuang), tetapi **bentuk
relasionalnya sudah ada** — berbeda dari Master Product Name Life yang harus dibangun dari nol.
Karena hilir **membaca** tabel-tabel ini, migrasi wajib **menjaga bentuk yang mereka baca**
(**ADR-0009** migrasi penuh; koeksistensi ditolak untuk data, bukan untuk kontrak baca hilir).

Sequence `PROPORTIONALARRG_SEQ`, `treatycontract_seq`, `TreatyYear_seq`, `M_TREATYREINSURER_SEQ`,
`TREATY_BUSINESS_SEQ` pindah dengan **nilai berjalan yang benar**.

---

## Acceptance Criteria

### Batas konteks & kepemilikan

1. [ ] Modul ini **satu-satunya penulis** `TREATYYEAR`, `TREATYCONTRACT`, `TREATYREINSURER`,
   `MTREATYSECURITY`, `TREATYBUSINESS`, `PROPORTIONALARRG`. *(§1; OQ-042)*
2. [ ] Konteks hilir (`Claim Prop`, `Komite Claim Prop`, `Claim Fac In`) dapat **membaca** bentuk
   relasional yang mereka pakai hari ini. *(§1, §16)*
3. [ ] Modul ini **tidak** menulis satu pun objek treaty **outward**. Test yang menemukan tulisan ke
   objek outward **gagal**. *(§Kontrak batas)*

### Tahun treaty & kontrak

4. [ ] Tahun treaty dapat dibuat dengan grup treaty, underwriting year, proporsi, dan masa
   berlakunya. *(User story 1)*
5. [ ] Identitas tahun treaty **tidak pernah diketik pengguna** — dibuat dari sequence. *(US 2; §7)*
6. [ ] Identitas berbentuk `'1'` diikuti nomor urut ber-*padding* nol: **6 digit** untuk tahun,
   kontrak, reinsurer, dan business; **7 digit** untuk klausul. *(US 3; §7; **ADR-0006**)*
7. [ ] Kontrak treaty dibuat **di dalam** sebuah tahun treaty dan menyimpan rujukan ke tahun itu.
   *(US 4; §2)*
8. [ ] Menyimpan tahun atau kontrak yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
   *(US 5; §6 — upsert dikunci `ID`)*
9. [ ] Masa berlaku yang **berakhir sebelum dimulai** — tahun treaty maupun kontrak — **ditolak**,
   dengan pesan yang menyebut field-nya. *(US 7)*

### Jenis reasuransi

10. [ ] Jenis reasuransi dipilih dari **master**, tidak pernah diketik bebas. *(US 8; §11)*
11. [ ] ⚠️ Daftar jenis reasuransi disaring **persis** seperti existing: **dua belas ID
    di-blacklist** (`10004`, `10011`, `10012`, `10021`, `10022`, `10025`, `10026`, `10028`, `10248`,
    `10249`, `10018`, `10217`) **dan** `Flag = "active"` **dan** `Type` termasuk `1`, `2`, `3`.
    *(US 9; §11; `[keputusan work owner]` — ditiru apa adanya, jangan digeneralkan)*
12. [ ] Master jenis reasuransi **tidak ditulis** oleh konteks ini. *(US 10; §11)*
13. [ ] ⚠️ Nama komponen di sistem baru **tidak mengandung kata "Old"** meskipun rule sumbernya
    bernama demikian. *(§11; penyimpangan sadar 8)*

### Reinsurer & security

14. [ ] Reinsurer dicatat pada kombinasi **(tahun, grup, jenis reasuransi)**, masing-masing dengan
    **share** dan **komisi reasuransi**. *(US 11; §2)*
15. [ ] **Total share** seluruh reinsurer pada satu kombinasi **terlihat** bagi pengguna. *(US 12)*
16. [ ] Share dan komisi diperlakukan sebagai **desimal presisi arbitrer**; **tidak** melewati
    `float`. *(US 29; §14; **ADR-0003**)*
17. [ ] Security dicatat **di bawah seorang reinsurer** beserta porsinya. *(US 13; §2)*
18. [ ] ⚠️ Baris security punya **primary key surrogate** sendiri; `REAS_SECURITY` adalah **atribut
    biasa**, **bukan** bagian kunci. Test yang menemukan kunci berbasis nama security **gagal**.
    *(US 14; §9; penyimpangan sadar 5)*
19. [ ] ⚠️ Seluruh kolom baris security **diisi eksplisit dan bernama**; **tidak ada** penulisan
    berposisi. *(§9; penyimpangan sadar 5)*
20. [ ] ⚠️ Pencocokan baris security **tidak** memakai `trim()` atas nilai kunci. *(§9)*

### Business

21. [ ] Jenis bisnis dicatat dengan **kode** dan **nama**-nya pada sebuah kontrak. *(US 15; §2)*
22. [ ] Satu baris bisnis dapat **dinonaktifkan** tanpa dihapus. *(US 16)*
23. [ ] ⚠️ Memperbarui baris bisnis **memperbarui seluruh field yang dikirim** — bukan hanya
    sebagian seperti procedure existing. *(§14)*

### Klausul

24. [ ] ⚠️ Seluruh **dua puluh lima** jenis klausul tersimpan dalam **satu tabel**
    `proportionalarrg`, dibedakan oleh **`TREATYDESCID`**. Test yang menemukan tabel terpisah per
    jenis klausul **gagal**. *(US 17; §3; penyimpangan sadar 2)*
25. [ ] ⚠️ Baris klausul jenis "anak" menyimpan **sembilan kolom khusus induk sebagai NULL**
    (`ID_OCCUPATION`, `OCCUPATION`, `ID_CLAUSE`, `CLAUSE`, `TREATYLIMIT`, `COINS_MIN`, `COINS_MAX`,
    `MORERP`, `MOREUSD`). *(§3; penyimpangan sadar 2)*
26. [ ] Daftar jenis klausul dibaca dari master **`TREATYDESC`**; jenis baru dapat ditambahkan di
    master **tanpa mengubah skema**. *(US 18; §12)*
27. [ ] Master jenis klausul **tidak ditulis** oleh konteks ini. *(§12)*
28. [ ] Klausul berlapis dapat dinyatakan lewat baris "anak" pada kombinasi yang sama
    (+ `PARENTREINSTYPEID`). *(US 19; §2)*
29. [ ] Membatalkan pengeditan **satu** klausul **tidak** membuang klausul lain yang sedang
    dikerjakan. *(US 21)*
30. [ ] Baris klausul yang **sudah pernah diinput** **ditolak** dengan pesan yang jelas. *(US 22;
    §5)*
31. [ ] Urutan baris dalam daftar klausul **dipertahankan** saat dibaca kembali.
32. [ ] ⚠️ Nama klausul dan nama field **dipakai apa adanya** dari Pega — **tidak diterjemahkan**.
    *(§4; `[keputusan work owner]`)*

### Validasi per jenis klausul

33. [ ] Setiap jenis klausul memeriksa **field yang relevan baginya sendiri**, sesuai tabel §5.
    *(US 20; §5)*
34. [ ] Aturan wajib-isi berada di **kode**, bukan di data master. *(§5)*
35. [ ] Pesan penolakan **menyebut field** yang kurang. *(§5)*
36. [ ] `[terbuka]` Jenis **`LimitMB`** dan **`Portfolio`** **tidak** dinyatakan selesai sebelum
    aturan wajib-isinya ditetapkan Product + UW. Keduanya **tidak** dilepas sebagai "tanpa
    validasi". *(§5; OQ terbuka)*

### Menyimpan

37. [ ] ⚠️ Seluruh perubahan satu kontrak — kontrak, reinsurer, security, business, **dan seluruh
    klausul** — ditulis dalam **satu transaksi**; kegagalan di mana pun **membatalkan seluruhnya**.
    Dibuktikan dengan menyuntikkan kegagalan pada satu baris anak lalu memastikan **tidak ada**
    perubahan tersimpan. *(US 24; §6)*
38. [ ] Penyimpanan yang **gagal** menghasilkan kegagalan **terang-terangan** dengan pesan yang
    menyebut apa yang gagal — bukan diam-diam dianggap sukses. *(US 6, 25; §6; **ADR-0015**)*
39. [ ] Status simpan **1 = sukses / 0 = gagal** ditegakkan; nilai selain `1` **selalu** dibaca
    sebagai kegagalan. *(§6)*
40. [ ] ⚠️ Pesan galat di sistem baru **tidak memuat teks `"JSON_KLAIM"`**. Test yang menemukannya
    **gagal**. *(§6; penyimpangan sadar 8)*
41. [ ] Setiap penyimpanan mencatat **jejak audit** — siapa dan kapan. *(**ADR-0007**)*

### Menghapus

42. [ ] ⚠️ Menghapus kontrak **mengkaskade** ke `TREATYBUSINESS`, `MTREATYSECURITY`, dan
    `TREATYREINSURER`. *(US 26; §8; penyimpangan sadar 4)*
43. [ ] ⚠️ Penghapusan didahului **popup konfirmasi Ya/Batal** yang **menyebut jumlah baris tiap
    jenis** yang akan ikut terhapus. **Batal** membatalkan seluruhnya. *(US 27; §8; penyimpangan
    sadar 4)*
44. [ ] ⚠️ **`PROPORTIONALARRG` (klausul) TIDAK ikut terhapus** dan tetap dapat dibaca setelah
    kontrak dihapus. Test yang menemukan klausul ikut terhapus **gagal**. *(US 23; §8;
    `[fakta bisnis — work owner]` — desain, bukan bug)*
45. [ ] Seluruh kaskade berjalan dalam **satu transaksi**. *(§8)*

### Kurs

46. [ ] Padanan **IDR** untuk nilai **USD** dihitung dari kurs yang **berlaku pada periode**
    kontrak. *(US 28; §10)*
47. [ ] ⚠️ Identitas mata uang USD (`IDCURRENCY = '10001'`) adalah **rujukan ke master mata uang**,
    **bukan literal di kode**. Test yang menemukan `'10001'` sebagai konstanta program **gagal**.
    *(§10; penyimpangan sadar 7)*
48. [ ] `QUARTER = '0'` **diikuti apa adanya**; artinya dicatat sebagai `[terbuka]` dan **tidak
    ditebak**. *(§10)*
49. [ ] ⚠️ Nilai **Rp** dan **Usd** tetap **dua nilai terpisah**. *(US 30; §10)*
50. [ ] ⚠️ Jalur kurs yang dimigrasikan adalah yang **benar-benar dipakai** (`testingKurs`), bukan
    yang namanya lebih wajar tetapi **di-remark**. Nama barunya **tidak mengandung kata "testing"**.
    *(§10; penyimpangan sadar 8; OQ-066)*

### Tipe data

51. [ ] ⚠️ **Seluruh** nilai uang — `RP`, `USD`, `MORERP`, `MOREUSD`, `TREATYLIMIT`, `COINS_MIN`,
    `COINS_MAX`, `PCT_SHARE`, `RICOMM` — diperlakukan sebagai **desimal presisi arbitrer**; **tidak**
    melewati `float`. Test yang menemukan kolom uang bertipe teks **gagal**. *(US 29; §14;
    penyimpangan sadar 6; **ADR-0003**)*
52. [ ] ⚠️ **Seluruh** persentase — `PCT`, `PCTME` — desimal, **tidak** dibulatkan ke bilangan bulat
    dan **tidak** disimpan sebagai teks. *(§14; penyimpangan sadar 6)*
53. [ ] ⚠️ **Seluruh** tanggal — `STARTDATE`, `ENDDATE`, `TREATYSTARTDATE`, `TREATYENDDATE`,
    `TGLUPDATE` — bertipe **`DATE`**, bukan teks. *(§14; penyimpangan sadar 6)*

### Lampiran (fitur baru)

54. [ ] ⚠️ Berkas dapat **dilampirkan pada sebuah tahun treaty** — fitur yang **belum ada** di Pega.
    *(US 31; §13; penyimpangan sadar 9)*
55. [ ] ⚠️ Lampiran **opsional**: kegagalan unggah **tidak** membatalkan tersimpannya tahun treaty —
    hanya status lampiran yang berubah. *(US 32; §13)*
56. [ ] Tiap lampiran dapat diberi **kategori** dari master kategori. *(US 33; §13)*
57. [ ] Lampiran dapat **diunduh** dan **dihapus**. *(US 34)*
58. [ ] ⚠️ Kegagalan unggah **tercatat dan dapat diulang**; pengulangan **tidak** menggandakan
    berkas. *(US 35; §13; **ADR-0015**)*
59. [ ] Alamat penyimpanan berkas di-resolve **runtime** dari konfigurasi. Test yang memindai kode
    untuk URL sebagai **literal, konstanta, atau pembacaan env var** **gagal** bila menemukannya.
    *(US 36; §13; **ADR-0013**)*
60. [ ] Token penyimpanan **di-cache** dan **diperbarui sebelum kedaluwarsa**; kegagalan
    mengambilnya menghasilkan kegagalan yang **terlihat**. *(§13)*
61. [ ] Rekam berkas di basis data dan berkas di penyimpanan **tetap sejalan**: rekam tanpa berkas
    terdeteksi dan dapat diperbaiki. *(§13)*
62. [ ] Klien penyimpanan berkas berada **di balik interface** dan **di-fake** di test; yang diperiksa
    adalah **efeknya**. *(§Testing)*

### Migrasi

63. [ ] ⚠️ Skema target **relasional penuh**. Test yang menemukan kolom JSON sebagai penyimpan
    atribut arrangement **gagal**. *(US 39; §15; penyimpangan sadar 1)*
64. [ ] ⚠️ Migrasi **tidak mengambil apa pun** dari `M_PROPORTIONALARRG`, `M_TREATYCONTRACT`,
    `M_TREATYBUSINESS`, atau tabel `M_*` lain. Sumbernya **hanya** tabel relasional.
    *(US 39; §15; penyimpangan sadar 1)*
65. [ ] Seluruh baris `PROPORTIONALARRG`, `TREATYCONTRACT`, `TREATYYEAR`, `TREATYREINSURER`,
    `TREATYBUSINESS`, `MTREATYSECURITY` pindah **tanpa kehilangan satu nilai pun**. *(US 37;
    **ADR-0009**)*
66. [ ] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara
    tepat**, bukan dengan toleransi. *(US 38; **ADR-0003**)*
67. [ ] Nilai tanggal yang hari ini berupa **teks** menjadi `DATE` **tanpa pergeseran zona waktu**;
    teks yang **tidak dapat diurai** dilaporkan, **tidak** didiamkan. *(US 38; §14)*
68. [ ] Baris `MTREATYSECURITY` mendapat **primary key surrogate** saat migrasi, dan rujukannya ke
    reinsurer tetap utuh. *(§9)*
69. [ ] Sequence `PROPORTIONALARRG_SEQ`, `treatycontract_seq`, `TreatyYear_seq`,
    `M_TREATYREINSURER_SEQ`, `TREATY_BUSINESS_SEQ` pindah dengan **nilai berjalan yang benar**,
    sehingga identitas baru **tidak bertabrakan** dengan yang lama. *(§16)*
70. [ ] `[terbuka]` Kolom `PROPORTIONALLIST` dan `OBJECT` **tidak dibawa** ke skema baru kecuali
    migrasi membuktikan ada data hidup di sana; temuannya **dilaporkan**. *(§3; OQ terbuka)*
71. [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.

### Tambahan validasi `[keputusan work owner]`

73. [ ] **Tahun treaty yang sama tidak boleh dibuat dua kali.** Kombinasi **(`STARTDATE`,
    `ENDDATE`, `TREATYGROUPID`)** yang **sudah ada** ditolak, dengan pesan yang menyebut tahun
    treaty mana yang sudah memakainya. `[keputusan work owner]` — aturan **baru**; existing tidak
    memeriksanya (`Activity/CheckYear.xml`, `@BASECLASS!CHECKYEAR`, hanya memeriksa
    `isNumber(TreatyYear)`).

### Fitur yang sengaja dibuang

72. [ ] ⚠️ **Fitur salin tahun treaty tidak dibangun.** Tidak ada jalur — layar, endpoint, maupun
    pekerjaan latar — yang menyalin isi satu tahun treaty ke tahun lain. Test yang menemukan jalur
    semacam itu **gagal**. *(§Out of Scope; `[fakta bisnis — work owner]` — penyimpangan sadar 3:
    menyalin membawa nilai tahun lalu ke tahun yang semestinya berbeda)*

---

## Testing Decisions

**Apa yang membuat test bagus di sini.** Test menguji **perilaku yang terlihat dari luar** — apa yang
tersimpan, apa yang ditolak, apa yang terlihat pengguna — **bukan** bentuk internal. Test yang
menyebut nama fungsi, urutan pemanggilan repository, atau bentuk struct adalah test yang akan pecah
pada perapian pertama tanpa menemukan satu pun bug.

**Seam: API HTTP** — seam yang **sama** dengan enam konteks sebelumnya (CL-01). **Tidak ada seam
baru.** Seluruh AC dapat diuji lewat seam ini.

**Oracle tidak pernah dipalsukan.** Yang diuji di sini — atomisitas lintas enam tabel, kaskade hapus
yang **mengecualikan** satu tabel, upsert dikunci `ID`, identitas dari sequence, presisi desimal,
dan konversi tanggal — **hanya berperilaku benar pada basis data sungguhan**. Memalsukannya berarti
tidak menguji apa pun yang penting. Test berjalan terhadap **skema uji Oracle nyata**.

**Satu-satunya yang di-fake**: klien **penyimpanan berkas** (Google Storage), di balik interface.
Yang diperiksa adalah **efeknya** — rekam lampiran, status, dan kemampuan mengulang — bukan
panggilannya.

**Yang wajib punya test tersendiri** karena di sinilah bug akan muncul:
- **Atomisitas** — suntikkan kegagalan pada baris klausul ke-N, pastikan **tidak ada** perubahan
  tersimpan pada kelima tabel lain. *(AC 37)*
- **Kaskade hapus yang mengecualikan klausul** — hapus kontrak, pastikan tiga anak hilang dan
  `PROPORTIONALARRG` **masih ada**. *(AC 42, 44)*
- **Filter jenis reasuransi** — dua belas ID ter-blacklist **tidak** muncul; `Flag`/`Type` ditegakkan.
  *(AC 11)*
- **Validasi per jenis klausul** — satu kasus per jenis, memakai tabel §5 sebagai tabel kebenaran.
  *(AC 33)*
- **Presisi uang** — nilai dengan banyak angka di belakang koma melewati simpan-baca **tanpa berubah
  satu digit pun**. *(AC 51)*
- **Rekonsiliasi migrasi** — perbandingan **tepat**, bukan bertoleransi. *(AC 66)*

**Prior art**: pola seam, pemakaian Oracle nyata, dan bentuk test kaskade+popup mengikuti
`.scratch/master-contract-retro-life/` (tiket 09 dan 10); pola efek keluar lampiran mengikuti
`.scratch/master-product-name-life/` (tiket 08 dan 09).

---

## Out of Scope

**Kode mati — tidak dimigrasikan** `[keputusan work owner]`:

| Yang dibuang | Bukti |
| --- | --- |
| **`M_PROPORTIONALARRG` (JSON)** dan **seluruh** kueri `FROM m_PROPORTIONALARRG` — `GetMasterDescriptionEPIParentList`, `GetMasterDescriptionPLAParentList`, `GetMasterDescriptionExGratiaList`, `GetMasterDescriptionFACINParentList`, `GetMasterPortfolioListDetail`, `GetMasterPanggilID`, dan lainnya | `[data DBA]` tidak ada procedure yang menulis ke sana |
| **Fitur salin tahun treaty** — `POOLDATA.PROSESCOPY`, `RDBList/SaveMasterCopyData_SQL.xml`, `Activity/BrowseCopyData.xml`, `Activity/SaveTreatyYearMultiple_Act.xml`, bagian salin `Activity/BrowseDeleteRowTreatyInContract.xml` | `[fakta bisnis — work owner]` menyalin membawa nilai tahun lalu → sumber kesalahan |
| **`SetTreatyArrangementDesc_Act`** (`ASM-FW-GISFW-INT-PROPORTIONALARRG`) | `[terverifikasi]` **5 dari 6 langkah di-remark** |
| **Rule jenis reasuransi yang tanpa "Old"** — kecuali sebagai perilaku layar master jenis reasuransi | `[terverifikasi]` dipakai 1 layar; yang dominan 11 grid |
| **Kolom `PROPORTIONALLIST` dan `OBJECT`** | `[data DBA]` tidak di-set procedure mana pun → konfirmasi saat migrasi (AC 70) |

`[terverifikasi]` Sisa kode mati yang **wajib disensus sebelum ditiketkan**, karena berkas ada ≠
dipakai (**OQ-066**): **33** berkas memuat `1=2`, **1** memuat `1==2`, **41** memuat visibilitas
`never`, **46** langkah di-remark (`<pyStepsBlockName>//`), **2** langkah `EXIT`, **1** memo
"not used".

**Di luar konteks ini:**

- **Master jenis reasuransi** (`REINSURANCETYPE`) dan **master jenis klausul** (`TREATYDESC`) —
  konteks/menu tersendiri; di sini **baca saja**. `[keputusan work owner]`
- **Master grup treaty**, **master mata uang**, **master kategori lampiran** — dibaca saja.
- **Treaty outward** — ada di `NB Treaty In` dan `Claim Non Prop`, **nol** di modul ini.
- **Konteks hilir** `Claim Prop`, `Komite Claim Prop`, `Claim Fac In` — mereka **membaca**; migrasi
  mereka adalah pekerjaan tersendiri.
- **Penetapan nama bounded context** untuk 303 rule ini — **OQ-022**, keputusan GATE.
- **`Claude outputs/Struktur_InboxTreatyContract.xlsx`** — berkas non-Pega di dalam korpus
  READ-ONLY; **tidak dibaca** (**OQ-054**).
- **RBAC / siapa boleh mengedit klausul** — `[terverifikasi]` `discovery/flows/_METHOD-noflow.md`
  §3.5: modul ini **tidak memuat** ekspresi visibilitas ber-workbasket; tidak dapat direkonstruksi
  dari korpus.

---

## Further Notes

### OQ terbuka yang **tidak memblokir**

| OQ | Isi |
| --- | --- |
| **OQ-047** | Alamat fisik Google Storage ada di `M_LINK_SERVICE`, isinya belum dibaca. **ADR-0013** sudah mengatur **cara** meresolusinya — yang mengikat adalah **alamat tidak ditanam**. |
| **OQ-022** | *Mengapa* modul dinamai "Out" dan ke bounded context mana ia ditempatkan — keputusan GATE. |
| **OQ kecil** | Arti istilah klausul dan field (§4) — Product + UW · field wajib `LimitMB` & `Portfolio` (AC 36) · arti `QUARTER = '0'` (AC 48) · kolom `PROPORTIONALLIST`/`OBJECT` (AC 70). |

### Urutan irisan yang disarankan untuk `/to-tickets`

1. **Skema + migrasi** — tipe dirapikan, `MTREATYSECURITY` diberi PK, JSON dibuang. **PREFACTOR**:
   tidak ada irisan lain yang berdiri sebelum bentuk barunya ada.
2. **Tahun treaty + kontrak** — identitas dari sequence, gagal terang-terangan, gerbang masa berlaku.
3. **Jenis reasuransi** — master dibaca, filter blacklist 12 ID, nama jujur.
4. **Reinsurer + total share**.
5. **Security reinsurer** — PK surrogate, struktur bersih.
6. **Business**.
7. **Klausul — satu tabel, 25 jenis** — inti konteks ini.
8. **Validasi per jenis klausul** — tabel §5 sebagai tabel kebenaran.
9. **Simpan atomik lintas enam tabel**.
10. **Kaskade hapus + popup; klausul dikecualikan**.
11. **Kurs USD → IDR**.
12. **Lampiran di tahun treaty** — fitur baru.

### Motif yang terus berulang di proyek ini

**"Baca kodenya, jangan namanya."** Konteks ini menyumbang tiga contoh baru sekaligus: rule bernama
**`_Old`** yang justru terbaru dan dipakai 11 grid; procedure bernama **`_CHILD`** yang menulis ke
tabel yang **sama**; dan activity bernama **`testingKurs`** yang justru jalur produksi, sementara
`SetTreatyArrangementDesc_Act` yang bernama wajar **di-remark**. Tambahan tahun ini:
**"baca juga body procedure-nya, jangan tanda tangannya"** — `PEGA_M_PROPORTIONALARRG_CHILD` hanya
terbuka artinya setelah bodinya dibaca.

Dan satu motif baru dari konteks ini: **"periksa apakah procedure commit sendiri."** Lima konteks
Life sebelumnya menemukan `COMMIT` internal dan karenanya harus menerima titik potong transaksi.
Di sini `[data DBA]` membuktikan **sebaliknya** — dan itulah yang membuat simpan atomik lintas enam
tabel menjadi mungkin tanpa membuang satu pun procedure.
