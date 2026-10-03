# Grilling — Master Product Name Life — Ronde 1

Tanggal: 2026-09-15
Konteks: `master-product-name-life` — **konteks/menu sendiri** `[keputusan work owner]`
Modul: **Master Product Name Life** (114 berkas)
Skill: `/mattpocock-skills:grilling`
Sumber: korpus `D:\XML\RNM_BRD\` (READ-ONLY), `discovery/modules|flows/Master Product Name Life.md`,
`discovery/context-map.md` §2.9, `CONTEXT.md`, `docs/adr/`,
`.scratch/master-contract-retro-life/` (pola konteks Life terbaru)

> **Konvensi.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[fakta bisnis — work owner]`; `[data DBA]`; `[dugaan]`;
> `[terbuka]` = OQ.

> `[keputusan work owner]` **Master Contract Retro Life adalah konteks terpisah** dan sudah selesai
> (spec + 12 tiket). Modul ini digrill sendiri.

---

## Bagian A — Jalur simpan: **JSON CLOB**, bukan kolom per kolom

Ini perbedaan paling tajam dari Master Contract Retro Life, dan ia mengubah bentuk seluruh spec.

`[terverifikasi]` **Dua procedure penulis produk menerima SATU CLOB**, bukan daftar kolom:

| Rule Connect-SQL | Class / Nama | Isi |
| --- | --- | --- |
| `SaveProductNameLIfe` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!SAVEPRODUCTNAMELIFE` / `RULE-CONNECT-SQL` | ⬇️ |
| `SaveProductNameInwardLIfe` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!SAVEPRODUCTNAMEINWARDLIFE` / `RULE-CONNECT-SQL` | ⬇️ |

```sql
DECLARE
  DATAPEGA CLOB;
  IDPEGA   VARCHAR2(32767);
BEGIN
  DBMS_LOB.CREATETEMPORARY(DATAPEGA, true);
  DATAPEGA := {InputParam.DATAPEGA};
  IDPEGA   := {InputParam.IDPEGA};
  POOLDATA.PEGA_M_PRODUCT_LIFE(DATAPEGA, IDPEGA,
      {OutputParam.ERRMSG out}, {OutputParam.IDPEGAOUT out}, {OutputParam.STSSAVE out});
  COMMIT;
END;
```

(`PEGA_M_PRODUCT_INWARD_LIFE` identik bentuknya.)

**Empat hal yang langsung terbaca:**

1. **Seluruh rekam produk menyeberang sebagai satu CLOB** bernama `DATAPEGA`. Langkah pemanggilnya
   `[terverifikasi]` berdeskripsi **"JSON"** (`SaveProductName_Act` step 9, baris 1848). Ini
   **menguatkan OQ-012**: master produk berbentuk JSON.
2. **Tiga keluaran**, bukan satu: `ERRMSG`, **`IDPEGAOUT`**, `STSSAVE`. Bandingkan Master Contract
   Retro Life yang hanya `HASIL1`. `IDPEGAOUT` tampak mengembalikan identitas yang terbentuk.
3. **`COMMIT` di dalam pembungkus** — pola sama dengan konteks Life lain.
4. `DBMS_LOB.CREATETEMPORARY` adalah **built-in Oracle**, bukan procedure bisnis — D3 memasukkannya
   ke daftar "5 stored procedure"; sebenarnya yang bisnis hanya **empat**.

### ⚠️ A1. Jalur ketiga yang mem-bypass procedure

`[terverifikasi]` `SaveProductNameLIfeFlat` (`ASM-FW-GISFW-INT-PRODUCT_LIFE` /
`ASM!SAVEPRODUCTNAMELIFEFLAT` / `RULE-CONNECT-SQL`) **tidak memanggil procedure apa pun**:

```sql
UPDATE POOLDATA.M_PRODUCT_LIFE SET
  RIRISKID = {InputData.CARI2},
  RIRISK   = {InputData.CARI3}
WHERE ID = {InputData.CARI1}
```

**`UPDATE` langsung atas dua kolom saja.** Dua akibat:

- **`M_PRODUCT_LIFE` punya kolom nyata** (`ID`, `RIRISKID`, `RIRISK`) — jadi ia **bukan** murni
  penyimpan JSON. CLOB `DATAPEGA` agaknya **diurai procedure menjadi kolom**. → **Q7**
- Ada **jalur tulis yang melewati procedure**. Bila procedure memegang aturan (validasi, versi,
  audit), jalur ini **melompatinya**. → **Q2**

`[terverifikasi]` Rule ini yang **paling baru** di antara ketiganya (cap `20240821`).

---

## Bagian B — Dua jalur simpan produk: kapan yang mana

`[terverifikasi]` **`SaveProductName_Act`** (`ASM-FW-GISFW-INT-PRODUCT_LIFE` /
`SAVEPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY`, 165.682 byte, 16 langkah):

| Step | Isi | Precondition |
| ---: | --- | --- |
| 1 | `Property-Set` | `ProductName.TYPE=="" \|\| ProductName.GRUP==""` |
| 2 | `Property-Set-Messages` **"error Product Name"** | `ProductName.PRODUCTNAME==""` |
| 3 | `Property-Set-Messages` **"error Ceding"** | `ProductName.CEDING==""` |
| 4 | `Property-Set-Messages` **"error Policy Holder"** | `ProductNameInward.POLICYHODERNAME==""` |
| 5 | `Property-Set-Messages` **"error SOB"** | `ProductName.SOBNAME==""` |
| 6 | `Property-Set` "set operator create" | `@PropertyHasValue(ProductName.CREATEOP)` |
| 7 | `call AddCommentList_Act` | |
| 8 | `Property-Set` | ⚠️ **`PoductName.TYPE==""`** |
| **9** | `RDB-List` **"JSON"** → **`SaveProductNameLIfe`** | ⚠️ **`PoductName.TYPE==""`** |

`[terverifikasi]` **`SaveInwardProductName_Act`** (`ASM-FW-GISFW-INT-PRODUCT_LIFE` /
`SAVEINWARDPRODUCTNAME_ACT`, 64.796 byte, 6 langkah):

| Step | Isi | Catatan |
| ---: | --- | --- |
| 1 | `Property-Set` | |
| ~~2~~ | `Property-Set-Messages` "error Type" | **REMARK** (`//` baris 414) |
| ~~3~~ | `Property-Set-Messages` "error Grup" | **REMARK** (`//` baris 601) |
| 4 | `Property-Set` | |
| **5** | `RDB-List` → **`SaveProductNameInwardLIfe`** | ⚠️ `PoductName.TYPE==""` |
| 6 | `Property-Set` "Set ID to the page" | ⚠️ `PoductName.TYPE==""` |

**Jalur inward lebih pendek dan validasinya dimatikan.** → **Q3**

### ⚠️ B1. Salah ketik yang mengenai **langkah simpan itu sendiri**

`[terverifikasi]` Sensus seluruh modul:

| Ejaan | Kemunculan | Berkas |
| --- | ---: | --- |
| `ProductName` | **714** | tersebar |
| **`PoductName`** (tanpa `r`) | **20** | **hanya 2** — kedua activity simpan |

`[terverifikasi]` **`PoductName` tidak pernah dideklarasikan sebagai halaman** — nol
`<pyPagesAndClassesPage>` dan nol `<pxPageName>` dengan nama itu.

`[terverifikasi]` Kedua ejaan memakai kode transisi yang **sama** (`WhenTrue = 3`), jadi maksudnya
memang guard yang sama. Tetapi yang bersalah ketik justru **menjaga langkah 8–9 di jalur utama dan
langkah 5–6 di jalur inward** — yakni **langkah yang memanggil procedure penyimpan**.

⚠️ Akibat sebenarnya **tidak dapat dipastikan dari korpus** (bergantung semantik Pega atas rujukan
halaman tak dikenal). Yang pasti: **guard pada langkah simpan bukanlah guard yang sama dengan
langkah 1**, meski ditulis seolah identik. → **Q4**

---

## Bagian C — `OR` akhirnya terbukti, dan ia menautkan dua konteks Master

`[terverifikasi]` `BrowseReinstypeOR_SQL` (`ASM-FW-GISFW-INT-PRODUCT_LIFE` /
`ASM!BROWSEREINSTYPEOR_SQL` *(class terbaca dari berkas)* / `RULE-CONNECT-SQL`):

```sql
SELECT tc.*, ty.*
FROM pooldata.treatycontract_life tc
JOIN pooldata.treatyyear_life ty ON ty.id = tc.idtreatyyear
WHERE TO_DATE({Temp.CARI1},'DD/MM/YYYY') >= tc.treatystartdate
  AND TO_DATE({Temp.CARI2},'DD/MM/YYYY') <= tc.treatyenddate
  AND REINSTYPEID = '10200'
```

✅ **Ini menutup sisa OQ-057.** D3 menyatakan kode `OR` "tidak dapat dibuktikan dari korpus" karena
string `"OR"` tidak pernah muncul sebagai literal — **benar untuk stringnya, keliru untuk kodenya**.
`REINSTYPEID = '10200'` **adalah** literal, dan `[fakta bisnis — work owner]` dari konteks Master
Contract Retro Life: **`10200` = `OR`**.

⚠️ **Tautan lintas konteks yang nyata:** modul ini **membaca** `TREATYCONTRACT_LIFE` dan
`TREATYYEAR_LIFE` — tabel yang **ditulis** Master Contract Retro Life — disaring jenis reasuransi
`OR` **dan rentang tanggal berlaku**. → **Q8**

### ⚠️ C1. Bendera `.IsORS` dibandingkan dengan dua cara berbeda

`[terverifikasi]` `GetReinsTypeOR_Life` (`ASM-FW-GISFW-INT-PRODUCT_LIFE` / `GETREINSTYPEOR_LIFE` /
`RULE-OBJ-ACTIVITY`, 66.519 byte, 4 langkah):

| Baris | Precondition |
| ---: | --- |
| 343 | **`.IsORS=="'true'"`** ← dengan kutip di dalam kutip |
| 678 | `.IsORS=="true"` |
| 993 | **`.IsORS=="'true'"`** |

`[terverifikasi]` Sensus modul: **3× `=="true"`**, **2× `=="'true'"`**, dan **nol** tempat yang
**men-set** `.IsORS` — nilainya datang dari luar modul.

Dua dari lima perbandingan **tidak mungkin cocok** dengan tiga lainnya. → **Q5**

---

## Bagian D — OQ-056: jalur tulis `M_TREATY_IN` **nyata, dan bergerbang**

`[terverifikasi]` `SetTreatyIn_Act` (**`DATA-PORTAL`** / `SETTREATYIN_ACT` / `RULE-OBJ-ACTIVITY`,
159.737 byte, 15 langkah) — peta langkahnya:

| Step | Isi | Precondition |
| ---: | --- | --- |
| 2 | `call TreatyInInputVis` | |
| 4 | `RDB-List` → **`BrowseTreatyIn`** (baca `M_TREATY_IN`, `M_TREATY_IN_EDM`) | |
| 5 | `Java` "Map oracle column to clipboard" | |
| 7 | `Property-Set` "RevisionState" | |
| 8 | `RDB-List` → `GetCurrentDate` | |
| **9** | `Property-Set` **"add coment revisi"** | **`param.revisionstate==1`** |
| **10** | `Property-Set` **"Insert to Json"** | **`param.revisionstate==1`** |
| **11** | `RDB-List` → **`SaveTreatyIn`** → `POOLDATA.PEGA_TREATY_IN` (baris 2027) | **`param.revisionstate==1`** |
| 12 | `Call ConvertHistoryDate` | |
| 13 | `Call TreatySetReinstatement` | `TreatyIn.ProportionType=="NonProportional"` |
| 14 | `Call CheckDuplicateOffer` | |

**Jawaban parsial OQ-056, dari korpus:**

- ✅ **Jalur tulis `M_TREATY_IN` BENAR-BENAR ADA** di modul ini — `SetTreatyIn_Act` step 11.
- ✅ **Ia bergerbang `param.revisionstate == 1`** — bukan jalur biasa, melainkan **jalur revisi**.
- ⚠️ `[terverifikasi]` **`revisionstate` adalah parameter MASUK** activity ini (terdaftar di
  `pyXMLSignature` dan `pyParametersParamName` baris 249), dan **tidak ada satu pun tempat di modul
  ini yang menetapkannya**. Pemanggilnya berada **di luar** modul.

**Yang tersisa untuk work owner:** dalam keadaan bisnis apa `revisionstate=1` dikirim, dan mengapa
**master produk** yang merevisi **master treaty inward**. → **Q9**

`[terverifikasi]` Tangga persetujuan treaty **tidak ikut**: `Akseptasi_DT`, `TreatyInSubmit`,
`TreatyInActionButtons` **tidak ada** di modul ini.

---

## Bagian E — Tujuh pemilih master

`[terverifikasi]` Enam dari tujuh ber-class sama dengan entitas kerja; satu berbeda:

| FlowAction | Class / Nama |
| --- | --- |
| `ChooseCauseOfLoss` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECAUSEOFLOSS` |
| `ChooseCeding` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECEDING` |
| `ChooseCurrency` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECURRENCY` |
| `ChoosePolicyHolder` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSEPOLICYHOLDER` |
| **`ChooseRIRate`** | **`ASM-FW-GISFW-DATA-PLAN`** / `CHOOSERIRATE` ← **class berbeda** |
| `ChooseRIRisk` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSERIRISK` |
| `ChooseSOB` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSESOB` |

→ **Q6**

---

## Bagian F — Google Storage, token, dan lampiran

`[terverifikasi]` Seluruh rantai berkas ber-class **`ASM-FW-GISFW-INT-T_STORAGE_IMAGE`**:

| Rule | Class / Nama / Tipe |
| --- | --- |
| `ServiceGoogle` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `SERVICEGOOGLE` / `RULE-CONNECT-REST` |
| `GetTokenStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETTOKENSTORAGE_SQL` / `RULE-CONNECT-SQL` |
| `GetLinkStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETLINKSTORAGE_SQL` / `RULE-CONNECT-SQL` |
| `Insert_T_Storage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!INSERT_T_STORAGE_SQL` / `RULE-CONNECT-SQL` |
| `GetMimeType` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `GETMIMETYPE` / `RULE-OBJ-DECISIONTABLE` |

`[terverifikasi]` Satu-satunya `When` di modul: `@BASECLASS` / `RECORDEVENT` — bukan aturan bisnis.

⚠️ Berbeda dari Master Contract Retro Life yang **tanpa integrasi luar sama sekali**, modul ini
punya **efek keluar nyata** → **ADR-0013** (resolusi endpoint) dan **ADR-0015** (efek keluar) kini
berlaku. → **Q10**

---

# Frontier Ronde 1 — 10 pertanyaan

**Empat fakta bisnis (OQ tim — TIDAK saya tebak): Q1, Q6, Q8, Q9.**
**Enam keputusan desain / konfirmasi (saya beri rekomendasi): Q2, Q3, Q4, Q5, Q7, Q10.**

❓ **Q1** — **`PRODUCT_LIFE` versus `PRODUCTINWARD_LIFE`: dua entitas atau satu entitas dua sudut
pandang?** *(fakta bisnis)* `[terverifikasi]` Ada **dua procedure penulis terpisah**
(`PEGA_M_PRODUCT_LIFE`, `PEGA_M_PRODUCT_INWARD_LIFE`), **dua activity simpan terpisah**, dan **dua
class** (`ASM-FW-GISFW-INT-PRODUCT_LIFE`, `ASM-FW-GISFW-Int-PRODUCTINWARD_LIFE`). Atribut inward
bercorak underwriting: `MINAGE`, `MAXAGE`, `BROKERAGE`, `EXTRAPREMI`, `INSURED`, `RNMLIMITNUM`,
`MAXSUMREASURED`, `MAXEXPIREDCLAIM`, `MAXDATARECEIVE`, `SUBJECTTO`, `ADDENDUMWORD`.

Apakah: (a) **dua entitas berbeda** yang kebetulan diedit di satu layar; (b) **satu produk** dengan
dua kelompok atribut (umum + syarat underwriting inward); atau (c) hubungan induk-anak?

Dan apa arti atribut limitnya — khususnya `RNMLIMITNUM`, `MAXSUMREASURED`, `MAXEXPIREDCLAIM`,
`MAXDATARECEIVE`?

➡️ **Tidak ada rekomendasi — ini fakta bisnis.** Jawabannya menentukan bentuk skema dan apakah satu
layar menyimpan satu rekam atau dua.

---

❓ **Q2** — **Jalur `UPDATE` langsung yang mem-bypass procedure: dipertahankan?**
`[terverifikasi]` `SaveProductNameLIfeFlat` meng-`UPDATE` **hanya `RIRISKID` dan `RIRISK`** pada
`M_PRODUCT_LIFE`, **tanpa** melewati procedure — dan ia yang **paling baru** (cap `20240821`).

➡️ **Rekomendasi: satukan ke jalur procedure.** Bila procedure memegang validasi, versi, atau audit
— dan itu baru akan diketahui dari bodinya (**OQ-002**) — maka jalur pintas ini **melompatinya**,
sehingga perubahan RI Risk tidak tercatat seperti perubahan lain. Menyatukannya membuat satu pintu
tulis. ⚠️ Bila ternyata jalur pintas ini **sengaja** dibuat untuk alasan kinerja pada pembaruan
massal, itu keputusan berbeda — saya perlu tahu.

---

❓ **Q3** — **Validasi jalur inward yang dimatikan: dihidupkan atau memang tidak berlaku?**
`[terverifikasi]` Di `SaveInwardProductName_Act`, langkah **2 ("error Type")** dan **3 ("error
Grup")** keduanya **REMARK** (`//` baris 414 dan 601). Jalur utama `SaveProductName_Act` **tidak**
me-remark padanannya dan malah punya **lima** pemeriksaan (Product Name, Ceding, Policy Holder,
SOB, ditambah Type/Grup).

➡️ **Rekomendasi: hidupkan, samakan dengan jalur utama.** Dua jalur simpan atas entitas yang sama
sebaiknya tidak berbeda ketatnya — kecuali memang ada alasan bisnis. ⚠️ Per **OQ-066**, saya
**tidak** menyimpulkan sendiri bahwa remark itu kelalaian; butuh keputusan.

---

❓ **Q4** — **Salah ketik `PoductName` pada langkah simpan: diperbaiki?** `[terverifikasi]` Ejaan
tanpa `r` muncul **20 kali di dua berkas saja** — keduanya activity simpan — dan **tidak pernah
dideklarasikan sebagai halaman**, sementara `ProductName` muncul **714 kali**. Yang bersalah ketik
justru menjaga **langkah yang memanggil procedure penyimpan**.

➡️ **Rekomendasi: perlakukan sebagai salah ketik, dan tegakkan guard yang benar pada langkah
simpan.** Maksudnya jelas — kode transisinya identik (`WhenTrue = 3`) dengan guard bereja benar di
langkah 1. ⚠️ Konsekuensinya **tidak netral**: bila selama ini guard itu tak pernah menyala,
memperbaikinya akan **mulai menolak** penyimpanan yang dahulu lolos. Itu perubahan perilaku yang
harus disadari, bukan diselundupkan.

---

❓ **Q5** — **`.IsORS` dibandingkan dua cara: mana yang benar?** `[terverifikasi]` Tiga kali
`=="true"`, dua kali **`=="'true'"`** (kutip di dalam kutip), dan **nol** tempat yang men-set nilai
itu di modul ini.

➡️ **Rekomendasi: satu bentuk — perlakukan `IsORS` sebagai boolean sejati, bukan teks.** Dua
perbandingan berbeda atas satu bendera berarti **sebagian cabang tidak pernah menyala**; mana yang
mati bergantung pada nilai yang dikirim pemanggil. Di sistem baru: satu tipe, satu perbandingan.
⚠️ Saya perlu tahu **siapa yang menyalakan `IsORS`** — itu di luar modul ini.

---

❓ **Q6** — **Tujuh pemilih master: dari mana masing-masing datanya?** *(fakta bisnis)*
`[terverifikasi]` Enam ber-class `ASM-FW-GISFW-INT-PRODUCT_LIFE`; **`ChooseRIRate` ber-class
`ASM-FW-GISFW-DATA-PLAN`** — berbeda sendiri.

Yang saya butuhkan per pemilih: **tabel/master sumbernya**, dan apakah ada penyaring lini (seperti
`Flag = 1` untuk jenis reasuransi di Master Contract Retro Life). Khususnya: mengapa **RI Rate**
berada di class `DATA-PLAN`, bukan di master produk?

➡️ **Tidak ada rekomendasi — ini fakta bisnis.** Saya dapat menelusur kelas dan RD-nya di ronde
berikutnya, tetapi **arti dan cakupan** tiap master adalah keputusan tim.

---

❓ **Q7** — **`M_PRODUCT_LIFE`: tabel berkolom, penyimpan JSON, atau keduanya?** `[terverifikasi]`
Dua bukti yang tampak bertentangan: procedure menerima **satu CLOB `DATAPEGA`** (langkahnya
berdeskripsi **"JSON"**), tetapi `SaveProductNameLIfeFlat` meng-`UPDATE` **kolom nyata**
(`RIRISKID`, `RIRISK`) pada tabel yang sama.

➡️ **Rekomendasi: perlakukan sebagai tabel berkolom yang procedure-nya menerima payload JSON lalu
menguraikannya.** Itu satu-satunya pembacaan yang menjelaskan kedua bukti. ⚠️ **Perlu dipastikan
DDL** (**OQ-001** untuk modul ini) — apakah ada kolom `JSONDATA`/`DATA_JSON` di samping kolom
biasa, seperti pola `RATE_LIFE`/`CURRENCY` yang merupakan **VIEW atas `JSONDATA`**. **Butuh DBA.**

---

❓ **Q8** — **Mengapa master produk membaca kontrak treaty retro?** *(fakta bisnis)*
`[terverifikasi]` `BrowseReinstypeOR_SQL` menggabungkan `treatycontract_life` × `treatyyear_life`,
disaring **`REINSTYPEID = '10200'`** (= **`OR`**) **dan** rentang tanggal berlaku. Ia dipanggil
`GetReinsTypeOR_Life` ketika `.IsORS` menyala.

Apa yang dilakukan produk dengan kontrak `OR` yang ditemukannya — apakah produk **mewarisi** sesuatu
dari kontrak itu (limit? rate? masa berlaku?), atau sekadar menandai keterkaitan?

➡️ **Tidak ada rekomendasi — ini fakta bisnis.** ✅ Sekaligus mencatat: temuan ini **menutup sisa
OQ-057** — kode `OR` kini terbukti lewat ID `10200`, bukan lewat nama rule.

---

❓ **Q9** — **`revisionstate = 1`: kapan master produk merevisi master treaty inward?**
*(fakta bisnis — inti OQ-056)* `[terverifikasi]` Jalur tulis `M_TREATY_IN` **ada dan aktif**:
`SetTreatyIn_Act` step 11 → `SaveTreatyIn` → `POOLDATA.PEGA_TREATY_IN`, bersama step 9 ("add coment
revisi") dan step 10 ("Insert to Json") — **ketiganya bergerbang `param.revisionstate==1`**.
`revisionstate` adalah **parameter masuk** yang **tidak ditetapkan di modul ini**.

Dalam keadaan bisnis apa pemanggil mengirim `revisionstate = 1`, dan **mengapa penyuntingan master
produk life berujung menulis master treaty inward umum**?

➡️ **Tidak ada rekomendasi — ini fakta bisnis, dan ia menentukan cakupan konteks.** Bila jawabannya
"jalur itu tidak dipakai lagi", modul ini menyusut banyak. Bila dipakai, konteks ini punya **efek
tulis lintas domain** yang harus masuk kontrak batas.

---

❓ **Q10** — **Efek keluar Google Storage: wajib berhasil atau boleh tertunda?** `[terverifikasi]`
Modul ini punya rantai berkas lengkap ber-class `ASM-FW-GISFW-INT-T_STORAGE_IMAGE`: `ServiceGoogle`
(REST), token (`GetTokenStorage_SQL`), tautan (`GetLinkStorage_SQL`), rekam
(`Insert_T_Storage_SQL`), jenis berkas (`GetMimeType`). Berbeda dari Master Contract Retro Life yang
**tanpa integrasi luar sama sekali**.

➡️ **Rekomendasi: perlakukan unggah lampiran sebagai efek keluar yang BOLEH tertunda, bukan wajib
berhasil — tetapi kegagalannya harus terlihat.** Alasan: lampiran adalah pelengkap master produk,
bukan syarat sahnya; memblokir penyimpanan produk karena Google Storage sedang tak dapat dihubungi
akan menghentikan pekerjaan tanpa perlu. Yang wajib: **status lampiran terlihat**, dan **dapat
diulang**. Alamat tetap di-resolve runtime dari `M_LINK_SERVICE` (**ADR-0013**, OQ-047).
⚠️ Bandingkan **ADR-0015** (efek keluar Komite **wajib** berhasil) — saya sarankan **tidak**
menyamakan keduanya, karena taruhannya berbeda.

---

## Rekap frontier Ronde 1

| # | Jenis | Pemblokir spec? |
| --- | --- | --- |
| **Q1** dua entitas produk | **fakta bisnis** | **ya** — bentuk skema |
| Q2 jalur `UPDATE` pintas | desain | **ya** |
| Q3 validasi inward dimatikan | desain | **ya** |
| Q4 salah ketik `PoductName` | desain | **ya** — mengubah perilaku |
| Q5 `.IsORS` dua bentuk | desain | tidak |
| **Q6** sumber tujuh pemilih | **fakta bisnis** | **ya** |
| Q7 `M_PRODUCT_LIFE` JSON atau kolom | desain (+ **butuh DBA**) | **ya** |
| **Q8** guna kontrak `OR` | **fakta bisnis** | **ya** |
| **Q9** `revisionstate` & `M_TREATY_IN` | **fakta bisnis** | **ya** — menentukan cakupan |
| Q10 efek keluar Storage | desain | tidak |

**Frontier belum kosong.** Ronde 2 terbuka setelah Q1, Q7, dan Q9 dijawab — terutama untuk menelusur
Harness `InwardProductName` (540.636 byte, belum dibuka), `SetProductNameInward` (125.938 byte),
dan ketujuh RD pemilih master.

## Status OQ setelah Ronde 1

| OQ | Status | Pemilik |
| --- | --- | --- |
| **OQ-002** | `[terbuka]` — **kemungkinan pemblokir**. Body **empat** procedure bisnis: `PEGA_M_PRODUCT_LIFE`, `PEGA_M_PRODUCT_INWARD_LIFE`, `PEGA_TREATY_IN`, `GET_TOKEN_STORAGE`. (`DBMS_LOB.CREATETEMPORARY` = built-in Oracle, bukan procedure bisnis) | **DBA** |
| **OQ-001** (modul ini) | `[terbuka]` — DDL `M_PRODUCT_LIFE`, `M_PRODUCTINWARD_LIFE`: adakah kolom JSON di samping kolom biasa (Q7) | **DBA** |
| **OQ-056** | **menyempit tajam** — jalur tulis terbukti ada dan bergerbang `revisionstate==1`; yang tersisa **keadaan bisnisnya** (Q9) | Product+UW |
| **OQ-057** | ✅ **DITUTUP** — `OR` terbukti lewat `REINSTYPEID = '10200'` di `BrowseReinstypeOR_SQL`, berpadu enumerasi Master Contract Retro Life | — |
| **OQ-058** | `[terbuka]` — `DATAPEGA.PC_ASM_FW_GCNMFW_WORK`: peran akses SQL langsung ke tabel internal Pega, dan nasibnya pasca-migrasi | **DBA** + work owner |
| **OQ-012** | **menguat** — payload simpan memang CLOB berdeskripsi "JSON"; bentuk tabelnya menunggu DDL (Q7) | DBA |
| **OQ-047** | `[terbuka]` — alamat Google Storage dari `M_LINK_SERVICE`; **ADR-0013** berlaku di sini | — |
| **OQ-066** | `[terbuka]` — diterapkan: remark di `SaveInwardProductName_Act` **tidak** saya simpulkan sendiri sebagai kelalaian (Q3) | — |
