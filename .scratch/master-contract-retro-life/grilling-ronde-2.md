# Grilling Ronde 2 — Master Contract Retro Life

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Master Contract Retro Life\` (READ-ONLY)

> ⭐ **ATURAN INDUK: MENAMBAH, BUKAN MENGUBAH** `[keputusan work owner 2026-09-19]`
> ⛔ `grilling-ronde-1.md` · `grilling-ronde-1-jawaban.md` · `spec.md` · `issues/01…12` ·
> `ddl-tables-from-dba.md` · `procedure-bodies-from-dba.md` — **NOL disentuh, satu byte pun.**
> ⛔ Angka lama **dikutip**, tidak dihapus dan tidak ditimpa.
>
> ⛔ Modul lain NOL · kode NOL · `CREATE TABLE` NOL · DDL NOL · nol nomor baris XML ·
> ⛔ nol butir `[terbuka]`/OQ dinyatakan tertutup · ⛔ nol revisi ADR diusulkan.
>
> ⚠️ **Tabrakan dilaporkan:** blok ini datang saat `/grill-with-docs` baru dipanggil.
> Blok GANTI KONTEKS **menang**; skill itu dibatalkan.

---

## §A — Hitung ulang dengan aturan baca yang baru

⚠️ **Ronde 1 dikerjakan 15 September; `ATURAN-BACA-KORPUS-PEGA.md` lahir 19 September.** Seluruh
ronde 1 karena itu disusun **tanpa** aturan bacanya. Inilah yang paling menahan.

### A1 — Sensus berkas

**Jendela (S1):** seluruh isi folder modul, **rekursif**, tanpa pengecualian.

| | Jumlah |
| --- | --- |
| **TOTAL isi folder** | **69** |
| ⭐ berkas `.xml` rule Pega | **66** ← ⛔ *ronde 1 menulis* **"66 berkas"** — **COCOK** |
| berkas **bukan** rule | **3** |
| di antaranya di `Claude outputs\` | **2** |

**Ketiga berkas bukan rule:** `Struktur_GridRetrocessionLife.xlsx` *(akar modul)* ·
`Claude outputs\perubahan_skill.diff` · `Claude outputs\Struktur_GridRetrocessionLife.xlsx`.

⚠️ `Claude outputs\` adalah **OQ-054** dan **tidak dihitung sebagai korpus** — dikecualikan sesuai
brief. ⭐ Tetapi **satu `.xlsx` berada di AKAR modul, di luar folder itu** — ⛔ belum pernah
tercatat di mana pun. `[terbuka]`.

**Sebaran 66 rule per jenis:**

| Jenis | Jumlah |
| --- | --- |
| Activity | **27** |
| RDBList | **12** |
| ReportDefinition | **12** |
| Section | **8** |
| Harness | **4** |
| FlowAction | **2** |
| DataTransform | **1** |

### A4 — ⭐ DUA JANGKAR: **angkanya BERBEDA** *(aturan P4 terpicu)*

| Cara hitung | Angka |
| --- | --- |
| Jangkar **A** — `<pxObjClass>Embed-ActivitySteps</pxObjClass>` | ⭐ **94** |
| Jangkar **B** — `<pyStepsActivityNameUC>` | ⭐ **91** |
| Pengurai XML — `rowdata` bersarang di dalam `pySteps` | **94** |
| Baris syarat — `Embed-ActivityPreConditions` | **108** |

> ## ⚠️ ATURAN P4 BERLAKU — **BELUM PUNYA DATA, BUKAN TINGGAL PILIH**
>
> **Jangkar A = 94 · Jangkar B = 91. Selisih 3.** Ketiganya dilaporkan apa adanya.
>
> ⚠️ Pengurai XML **setuju dengan jangkar A** *(94)*, yang menjadikan **B** kandidat yang meleset —
> kemungkinan besar karena `pyStepsActivityNameUC` **tidak terbit** pada langkah yang metodenya
> kosong *(langkah blok murni, mis. pengulangan tanpa metode)*. ⛔ **Itu dugaan, bukan bacaan.**
> `[terbuka]`.

### A2 — ⭐ Langkah ber-remark: ronde 1 **tidak pernah memeriksanya**

**Jendela (S1):** 27 Activity, **94** langkah *(jangkar A)*, medan `pyStepsBlockName`.

⭐ **Langkah ber-`//` di seluruh modul: SATU.**

| Rule | Langkah | Metode |
| --- | --- | --- |
| `NewInputBusinessLife_Act` | **1** | `Page-New` |

#### Ketiga "rumus", diperiksa satu per satu

| Activity | ⛔ Ronde 1 menulis | ⭐ Langkah **sebenarnya** | Ber-`//` | Kesimpulan ronde 1 |
| --- | --- | --- | --- | --- |
| `CountingPercentShare_Act` | **4 langkah** | ⭐ **5** | **0** | ✅ **BERTAHAN** |
| `TreatyLimit_TypeProtect` | **4 langkah** | ⭐ **5** | **0** | ✅ **BERTAHAN** |
| `SetValueRetroLimit_TreatyYearLife` | **2 langkah** | **2** | **0** | ✅ **BERTAHAN** |

> ### ⭐ KESIMPULAN RONDE 1 **BERTAHAN** — dan satu parafrase brief perlu diluruskan
>
> Brief menyebut kesimpulan induk ronde 1 sebagai *"ketiga rumus itu **BUKAN** rumus"*.
> ⚠️ **Ronde 1 tidak menulis begitu.** Kalimatnya: *"…dan **dua** di antaranya ternyata bukan
> rumus"*, dan ia **secara eksplisit** menyebut `CountingPercentShare_Act` sebagai
> **"penjumlah share"**. `spec.md` §15 mencatatnya sama: *"**penjumlah** `PCTSHARE`"*.
>
> ⭐ **Jadi ronde 1 sudah benar**, dan bahaya yang brief khawatirkan **tidak terjadi**: nol dari
> ketiganya ber-remark, sehingga kesimpulan dari isi langkah **tidak dirusak oleh langkah mati**.
>
> ⚠️ **Yang memang keliru adalah jumlah langkahnya** — dua dari tiga, masing-masing **4 → 5**,
> karena ronde 1 tampaknya tidak menghitung langkah **bersarang** *(3.1 dan 4.1)*.

**Isi `CountingPercentShare_Act` apa adanya** *(5 langkah, nol remark)*:

```
lgk 1    Property-Set   TempSearch.CARI1 := Param.TREATYYEARID
                        TempSearch.CARI2 := Param.TREATYCONTRACTID
                        InputSecurityLife.STDRATING := ""
                        Local.TotalShare := ""
lgk 2    RDB-List
lgk 3    (blok)         Local.TotalShare := 0
lgk 3.1  Property-Set   Local.TotalShare := Local.TotalShare + @toDecimal(.PCTSHARE)
lgk 4    Property-Set   InputSecurityLife.STDRATING := Local.TotalShare
```

⚠️ ⭐ **Total share ditulis ke medan bernama `STDRATING`** — *standard rating*. Nama medannya
**tidak menyebut share sama sekali**. Ini pola "nama berbohong" yang sama dengan yang tercatat di
modul lain. ⛔ Nol kesimpulan; dicatat sebagai temuan.

### A3 — ⭐ Kode arah gerbang: ronde 1 **tidak pernah membacanya**

**Jendela (S1):** 27 Activity · **94** langkah · **108** baris keluarga-1, di antaranya
**42** bersyarat *(medan `when` terisi)*. Nol dikecualikan.

**Flag prakondisi** *(penyebut 94 langkah)*: `(kosong)` **66** · `true` **20** · `false` **8**.

**Sebaran kode arah** *(penyebut 108 baris)*: `2/2` **67** · `3/2` **23** · `_/3` **9** · `2/3` **9**.

**Kode khusus:** `1` = **0** · `3` = **41** · `4` = **0** · `5` = **0** · `6` = **0**.

> ⭐ **Nol kode 5 berarti NOL rantai OR.** Setiap rantai syarat di modul ini adalah **AND** murni —
> *Continue Whens* + *Skip Step*. Nol lompatan, nol keluar-iterasi, nol keluar-activity.

#### ⭐ TEMUAN TERBESAR §A: **21 dari 42 baris bersyarat bergerbang MATI**

| Status | Baris |
| --- | --- |
| **BERLAKU** *(flag `true`)* | **21** |
| ⭐ **MATI** *(flag `false` — langkah jalan **tanpa saringan**)* | ⭐ **21** |

**Yang mati, seluruhnya — dan seluruhnya validasi simpan:**

| Rule | Langkah | Syarat | Kode |
| --- | --- | --- | --- |
| `SaveSecurityLife_Act` | **4 · 5 · 6** *(9 baris)* | `InputRetrocessionLife.REINSURERNAME==""` · `TREATYSTARTDATE==""` · `@substring(...)` | `3/2` |
| `SaveSecurityReinsurerLife_Act` | **4 · 5 · 6** *(9 baris)* | syarat yang **sama persis** | `3/2` |
| `SaveBusinessLife_Act` | **6** *(3 baris)* | syarat yang **sama persis** | `3/2` |
| `SaveBusinessToAllLife_Act` | **3** *(1 baris)* | `.REINSTYPEID==Param.REINSTYPEID` | `2/2` |

#### Keempat gerbang yang brief minta divonis

| Gerbang | Vonis |
| --- | --- |
| `TreatyLimit_TypeProtect` langkah **4.1** — `TempInputData.CARI1 == .ID` | ✅ **BERLAKU** — flag `true`, `2/3` |
| `SetErrorMessageReinsurer` — dua baris batas 0–100 | ✅ **BERLAKU** — keduanya flag `true`, `2/3` *(langkah 2 `PctShare>…`, langkah 3 `Ricomm>100`)* |
| ⭐ **gerbang konsistensi tahun** — `@substring(...) <> TREATYYEAR_LIFE` | ⛔ **MATI** — muncul **3 kali** *(`SaveSecurityLife_Act` 4·5·6 dan kembarannya)*, **seluruhnya flag `false`** |
| `SaveBusinessToAllLife_Act` — `.REINSTYPEID == Param.REINSTYPEID` | ⚠️ **TERBELAH** — langkah **3** flag `false` **MATI**; langkah **3.1** dan **3.2** flag `true` **BERLAKU** |

> ## ⭐ AKIBATNYA UNTUK SPEC — dan ini menyentuh bab yang sudah ditulis
>
> `spec.md` **§8 "Gerbang konsistensi tahun"** dan **`issues/03-gerbang-tahun-lewat-nilai-tanggal.md`**
> dibangun di atas gerbang yang, menurut korpus, ⭐ **TIDAK PERNAH BERJALAN**.
>
> ⚠️ Itu **tidak** berarti keputusannya salah — `grilling-ronde-1-jawaban.md` **Q8** sudah
> memutuskannya sebagai `[keputusan work owner]`, yaitu **dibangun**, bukan ditiru. ⭐ Tetapi
> **dasar faktualnya berubah**: dari *"Pega memeriksanya dengan cara yang keliru (substring)"*
> menjadi *"Pega **tidak memeriksanya sama sekali**"*.
>
> ⛔ Ralatnya ditulis di `TAMBAHAN-SPEC-ronde-2.md`. ⛔ `spec.md` **tidak disunting**.

### A5 — Identitas rule empat bagian

**Jendela (S1):** 66 berkas; medan `pxInsName` · `pyRuleSet` · `pyRuleSetVersion`.

| | |
| --- | --- |
| ⭐ **Ruleset berbeda** | ⭐ **SATU — `GISFW`** |
| Kombinasi ruleset+versi | **12** |
| `pxInsName` muncul di >1 berkas | **1** |

**Sebaran versi:** `01-01-92` **23** · `01-01-93` **21** · `01-01-94` **7** · `01-01-53` **6** ·
`01-01-58` **2** · `01-01-54` · `01-01-56` · `01-01-67` · `01-01-73` · `01-01-83` · `01-01-91` ·
`01-01-95` **1 masing-masing**.

⭐ **Modul ini bebas dari jebakan yang menggigit modul lain** — di Claim Fac In dua rule bernama
sama hidup di ruleset **GISFW** dan **GCNMFW**; di sini **hanya ada satu ruleset**, jadi
`pxInsName` **cukup** sebagai identitas **untuk modul ini saja**.

⚠️ ⛔ Itu **tidak** berlaku saat membandingkan dengan modul lain — dan modul ini memang memuat
rule ber-`@BASECLASS` yang dipakai bersama.

### A6 — ⭐ Sensus lingkungan: **belum pernah dikerjakan di modul ini**

**Jendela (S1):** 66 berkas, tiga medan, **kosong 0** di ketiganya.

| Medan | Nilai berbeda | Sebaran |
| --- | --- | --- |
| **`pxUpdateSystemID`** | **3** | ⭐ `pegadevnusare2` **53** · `pega` **12** · `pegaprdnusare` **1** |
| **`pxHostId`** | ⭐ **5** | `jboss117` **44** · `jboss1073` **9** · `a11207acce0fa39e49c3b08f8b22d60f` **6** · `pega-nusre` **6** · ⭐ `jboss122117` **1** |
| `pxCreateSystemID` | 3 | `pegadevnusare2` **51** · `pega` **14** · `pegaprdnusare` **1** |

> ⭐ **Nilai KELIMA yang belum pernah dinyatakan: `jboss122117`** — satu berkas. Ia menyerupai
> gabungan `jboss122` dan `117`, tetapi ⛔ **itu dugaan**; korpus tidak menerangkannya.
>
> ⚠️ **Konteks yang sudah diketahui** `[data work owner]`: `jboss1073` = production ·
> `jboss117` = dev · sistem memakai *mirroring*. ⛔ `pega-nusre`, id hash, dan `jboss122117`
> **belum dinyatakan**.
>
> ⛔ **Saya TIDAK menyimpulkan modul ini production atau dev.** Angkanya disajikan; artinya
> keputusan work owner.

⚠️ Sebagai pembanding *(bukan kesimpulan)*: Claim Fac In `pega` **66,4%**, modul ini
`pegadevnusare2` **80,3%** — **sebaran yang berlawanan**.

---

## §B — Tiga hal yang ronde 1 sendiri nyatakan belum dikerjakan

### B1 — Empat Harness: **isinya dibaca**

**Jendela (S1):** 4 berkas Harness; lima keluarga gerbang layar dari aturan **H**; penanda mati
`NEVER` dan `1=2`.

| Harness | Byte | Section yang disambung | Gerbang **HIDUP** | Gerbang **MATI** |
| --- | --- | --- | --- | --- |
| `InboxBusinessLifeReinsurers` | 486 911 | `InputBusinessLifeReinsurers` | per-sel **9** · per-layout **4** | per-sel **2** |
| `InboxRetroLifeReinsurersList` | 530 052 | `InputSecurityLifeReinsurers` | per-sel **6** · per-layout **3** | per-sel **3** |
| `InboxRetroLimitReinsurers` | 575 179 | `InputRetroLimitReinsurers` | per-sel **6** · per-layout **3** | ⭐ per-sel **6** |
| `InboxSecurityReinsurerLife` | 444 874 | `InputSecurityReinsurerLife` | per-sel **6** · per-layout **3** | per-sel **4** |

⭐ **Total: HIDUP 40 · MATI 15.** ⭐ Setiap Harness menyambung **tepat satu** Section — jadi keempat
"Inbox" adalah **pembungkus layar**, bukan layar gabungan.

⚠️ **Yang TIDAK berhasil saya baca:** kolom yang ditampilkan. Disisir medan `pyPropertyName` —
**tidak ketemu**. Sesuai aturan **I2**, kalimat yang sah: *"tidak ketemu di medan `pyPropertyName`"*,
**bukan** *"nol kolom"*. `[terbuka]`, dan ini **kekurangan ronde ini**.

⚠️ **Gerbang MATI di sini bermakna berbeda dari gerbang mati langkah.** Penanda `NEVER` / `1=2`
adalah **niat menyembunyikan** — elemen sengaja dibuat tak pernah tampil. ⭐ `InboxRetroLimitReinsurers`
punya **6**, terbanyak.

### B4 — ⭐ Section: **4 dari 8 belum pernah disebut ronde 1**

| Section | Byte | Status |
| --- | --- | --- |
| `GridRetrocessionLife` | 97 817 | disebut ronde 1 |
| ⭐ `InputBusinessLifeReinsurers` | **449 805** | ⛔ **BELUM DISEBUT** |
| ⭐ `InputDtlRetrocessionLife` | **203 911** | ⛔ **BELUM DISEBUT** |
| ⭐ `InputRetroLimitReinsurers` | **541 845** | ⛔ **BELUM DISEBUT** |
| `InputRetrocessionLife` | 468 908 | disebut ronde 1 |
| ⭐ `InputSecurityLifeReinsurers` | **496 437** | ⛔ **BELUM DISEBUT** |
| `InputSecurityReinsurerLife` | 408 684 | disebut ronde 1 |
| `ViewRate` | 131 052 | disebut ronde 1 |

⭐ **Keempat yang belum disebut berjumlah 1,69 MB** — lebih besar dari keempat yang sudah disebut
*(1,10 MB)*. ⚠️ Ketiganya adalah Section yang disambung oleh Harness di B1, jadi **layar utama
modul ini justru yang belum dibaca**.

### B5 — ReportDefinition: ⭐ **11 dari 12 belum pernah disebut**

| ReportDefinition | Status |
| --- | --- |
| `BrowseReinsuranceTypeLimit_RD` | disebut ronde 1 |
| `BrowseBusinessLife_RD` · `BrowseCedingCoLife_RD` · `BrowseDetailTreatyReisurerLife_RD` · `BrowseRateLifeSummary` · `BrowseRateLife_RD` · `BrowseRetrocessionLife_RD` · `BrowseSecurityReinsurer_Life_RD` · `BrowseTreatyBusiness_Life_RD` · `BrowseTreatyContract_Life_RD` · `BrowseTreatyYear_Life_RD` · `BrowseTreatyYear_RD` | ⛔ **BELUM DISEBUT — 11** |

⚠️ **Kelas sumber, kolom, dan saringan TIDAK berhasil saya baca.** Disisir `pyObjClass` ·
`pyColumnName` · `pyCriteriaValue` — ketiganya **nol** di kedua belas berkas. Sesuai **I2**:
*"tidak ketemu di tiga medan yang disisir"*. ⭐ **Ini kekurangan ronde ini yang paling besar di §B.**

⚠️ ⭐ Dua di antaranya bernama **`BrowseRateLife_RD`** dan **`BrowseRateLifeSummary`** — sumber
data yang B2 cari.

### B2 — `ViewRate` / `ViewRateTable`

⚠️ **Belum dibaca isinya di ronde ini.** Yang terbaca hanyalah keberadaannya:
`Section/ViewRate.xml` *(131 052 byte, disebut ronde 1)* · dua FlowAction · dua activity
`SetParamRate` / `SetParamRateTable`.

⭐ **Sumber datanya kemungkinan besar `BrowseRateLife_RD` dan `BrowseRateLifeSummary`** — ⛔ tetapi
itu **dugaan dari nama**, bukan bacaan. `[terbuka]`. **Kekurangan ronde ini.**

### B3 — ⭐ `RIRATEID` dan `RIRATE`: **jawabannya dari DDL, dan mengejutkan**

**Jendela (S1):** 66 berkas; penulis dari `PropertiesName`, pembaca dari `PropertiesValue` ·
`pyValue` · `pyBrowseSQL` · `pyColumnName`.

| | Penulis | Pembaca |
| --- | --- | --- |
| `RIRATEID` | **4** — `NewInputBusinessLife_Act` · `SetBusinessListLife_Act` · `SetParamRate` | **9** |
| `RIRATE` | **4** — sama | **8** |

> ## ⭐ `RIRATE` = **`VARCHAR2(1000)`** — TEKS, bukan angka
>
> `[data DBA]` `ddl-tables-from-dba.md` mencatatnya sebagai kejanggalan: *"Rate reasuransi disimpan
> sebagai string… mungkin format `0.5%` / rate bertingkat; **jangan** paksa jadi angka tanpa
> konfirmasi."*
>
> ⭐ **Jawaban B3: ia BUKAN angka uang dan BUKAN angka persen — ia TEKS.** Karena itu ⛔ **aturan
> ketelitian angka TIDAK berlaku padanya**, dan pertanyaan *"uang atau persen"* **tidak punya
> jawaban** selama isinya masih teks bebas.
>
> ⚠️ `RIRATEID` berdampingan dengannya — ⛔ apakah `RIRATE` adalah **nilai** dan `RIRATEID`
> **penunjuk ke tabel rate**, belum terbaca. `[terbuka]`.

---

## §C — Korpus diadu dengan jawaban DBA

⚠️ Kedua berkas DBA lahir **15 September**, **setelah** ronde 1 ditulis. Korpus belum pernah dibaca
ulang dengan keduanya di tangan.

### C1 — Lima procedure penulis: tanda tangan vs body

`[data DBA]` Kelimanya **berpola IDENTIK**. Yang terbaca dari body dan **tidak ada di korpus Pega**:

| # | Perilaku procedure | Ada di Pega? |
| --- | --- | --- |
| 1 | ⭐ **UPSERT dikunci `ID`** — `SELECT COUNT(1) WHERE ID=p_ID` → ada = UPDATE, tidak = INSERT | ⛔ **TIDAK** — Pega memanggil satu rule "Save", tidak membedakan |
| 2 | ⭐ **ID diciptakan DB** — `'1' \|\| lpad(seq.nextval,6,'0')`; `p_ID` **diabaikan** saat INSERT | ⛔ **TIDAK** |
| 3 | ⭐ **`TGLUPDATE` selalu `SYSDATE`** — parameter `p_TGLUPDATE` **diabaikan** | ⛔ **TIDAK** — Pega mengirimnya, dan nilainya dibuang |
| 4 | ⭐ **`COMMIT` di dalam tiap procedure**, `ROLLBACK` pada exception terluar | ⛔ **TIDAK** |
| 5 | ⭐ **`o_message` berisi HTML** — `<span style="color:red">…</span>` + `SQLERRM` | ⛔ **TIDAK** |
| 6 | **Validasi bisnis NOL** — nol cek 100%, nol cek anak, nol unique selain PK | ✅ sejalan — Pega juga nol |

**Validasi yang ada di Pega tetapi procedure TIDAK mengeceknya:** ⭐ **nol yang hidup.** Empat baris
validasi Pega yang menyentuh wilayah ini *(`SaveSecurityLife_Act` 4·5·6 dan kembarannya)* justru
**bergerbang MATI** — §A3. ⭐ **Jadi tidak ada penjaga di kedua sisi.**

**Kapan `HASIL1` diisi nilai gagal:** `[data DBA]` **kosong/NULL = sukses; berisi teks = gagal.**

⭐ **Parameter cocok semua?** ⚠️ **Tidak dapat dipastikan dari korpus** — kelima rule SQL Pega
mengirim nilai lewat penampung `InputData.CARI1 … CARI14`, **bukan** lewat nama kolom. Pemetaan
posisi ke parameter **tidak terbaca**. `[terbuka]`, dan ini **kekurangan ronde ini**.

### C2 — ⭐ Klaim ketiadaan `HASIL1`: **sebagian GUGUR**

**Jendela (S1):** **66** berkas `.xml`; medan `PropertiesName` · `PropertiesValue` · `pyBrowseSQL` ·
`pyValue` · `pyCriteriaValue`. Nol dikecualikan.

⛔ **Ronde 1 §C2 menulis:** *"hanya `DeleteRowBusiness.xml` yang menyebut `HASIL1`; tidak satu pun
activity `Save*` membacanya."*

| Yang diperiksa | Hasil |
| --- | --- |
| Berkas menyebut `HASIL1` | ⭐ **7 dari 66**, bukan 1 |
| Di antaranya **kelima rule SQL `Save*`** | ⭐ **ADA** — kelimanya **mendeklarasikan** `o_message` di dalam blok anonimnya |
| **Activity** yang **membaca** `HASIL1` | ⭐ **NOL** — satu-satunya activity yang cocok, `DeleteRowBusiness.xml`, memakai **`HASIL12`**, bukan `HASIL1` |

> ### ⭐ Klaimnya GUGUR sebagian, dan kesimpulannya justru MENGUAT
>
> **Salah:** *"hanya satu berkas menyebut `HASIL1`"* — sebenarnya **tujuh**, dan kelima jalur simpan
> **memang mendeklarasikannya**.
>
> **Benar, dan kini lebih tegas:** ⛔ **NOL activity membaca `HASIL1`.** Dengan body procedure di
> tangan, kini dapat dinyatakan **apa yang jatuh diam-diam**: procedure menaruh pesan galat
> *(berikut `SQLERRM`)* di `o_message`, melakukan `ROLLBACK`, dan ⭐ **Pega tidak pernah
> melihatnya** — pengguna melihat simpan yang seolah berhasil.
>
> ⚠️ Dan `DeleteRowBusiness.xml` memakai **`HASIL12`** — penampung **berbeda**. ⛔ Apakah itu salah
> ketik atau memang penampung lain **belum terbaca**. `[terbuka]`.

### C3 — Kolom DDL vs kolom yang dipakai Pega

⚠️ **Metode saya lemah, dan saya laporkan itu.** Kelima rule SQL memakai penampung
`InputData.CARI1 … CARI14`, sehingga **nama kolom sebagian besar tidak muncul** di teks SQL Pega.

| | Hasil |
| --- | --- |
| Kolom DDL unik *(5 tabel)* | **30** |
| Tidak muncul sebagai token di SQL Pega | **19** |
| Token huruf-besar di SQL yang bukan kolom DDL | **29** — ⭐ **seluruhnya `CARI1…CARI14` dan sejenisnya** |

> ⭐ **Kolom yang DITULIS Pega tetapi TIDAK ADA di DDL: NOL ditemukan.**
>
> ⚠️ Sesuai aturan **I2**, kalimat yang sah: ***"tidak ketemu dengan metode pencocokan nama
> kolom"*** — **bukan** *"tidak ada"*. Metodenya tidak dapat membuktikannya, karena Pega tidak
> memakai nama kolom. `[terbuka]`.

### C4 — Tipe dan presisi kolom uang dan share

`[data DBA]`:

| Kolom | Tipe DDL | Terhadap aturan ketelitian |
| --- | --- | --- |
| `IDR` · `USD` · `B_IDR` · `B_USD` · `IDR_SELISIH` · `USD_SELISIH` | **`NUMBER` tanpa presisi** | ✅ **selaras** — `NUMBER` tanpa `(p,s)` berarti presisi penuh Oracle, dan Go memakai tipe desimal presisi arbitrer |
| `PCTSHARE` · `COMMISION` · `OVR_COMM` | **`NUMBER` tanpa presisi** | ✅ selaras |
| ⭐ **`RIRATE`** | ⛔ **`VARCHAR2(1000)`** | ⛔ **TIDAK selaras** — ia teks |

⭐ **Berapa digit?** `NUMBER` tanpa presisi **tidak menetapkan jumlah digit** — jadi ⛔ **DDL tidak
membatasi apa pun**, dan batas 20 digit dengan 8 di belakang koma adalah **keputusan aplikasi**,
bukan batas basis data. ⚠️ Itu berarti nilai yang lebih panjang **akan tersimpan di Oracle** dan
**ditolak di Go** — arah kegagalan yang berlawanan dari modul lain.

### C5 — ⭐ Kaskade hapus: **berubah artinya**

⛔ **Ronde 1 §C menulis:** *"nol kaskade"* — dari sisi **Pega**.

`[data DBA]` **Keadaan FINAL basis data**, sesudah pembaruan work owner:

| | |
| --- | --- |
| **PK `ID`** | ✅ **ada di kelima tabel** |
| ⭐ **FK antar tabel** | ✅ **4 FK, mode `ON DELETE CASCADE`** — contract→year · reinsurer→contract · security→reinsurer · business→contract |
| **NOT NULL** | ⛔ **tetap nol** → wajib-isi ditegakkan di Go |

> ⭐ **"Risiko yatim" ronde 1 BERUBAH ARTINYA, bukan menguat.**
>
> Dari sisi **Pega** memang nol kaskade — itu **tetap benar**. Tetapi basis data **kini
> mengaskade sendiri**. ⚠️ Artinya risikonya **terbalik**: bukan lagi *"anak menjadi yatim"*,
> melainkan ⭐ ***"menghapus induk menghapus seluruh anaknya — dan Pega tidak tahu itu terjadi"***.
>
> ⚠️ `grilling-ronde-1-jawaban.md` **Q5** sudah memutuskan **kaskade + popup konfirmasi Ya/Batal**,
> dan `[data DBA]` menyatakan FK CASCADE **selaras** dengan keputusan itu. ⭐ Jadi keputusannya
> **tidak berubah**; yang berubah adalah **siapa yang melakukan kaskade** — basis data, bukan
> aplikasi. ⛔ Butir ini **tidak saya nyatakan tertutup**.

---

## §D — Korpus diadu dengan 15 ADR

⚠️ **Sama seperti modul lain:** kelima belas ADR bersumber dari **Claim Life** dan **Komite Claim
Life**. Nol menyebut Master Contract Retro Life.

| ADR | Vonis | Bukti satu baris |
| --- | --- | --- |
| **0001** batas konteks | ⚠️ **menguatkan sebagian** | modul ini **memang konteks/menu sendiri** — 4 Harness masing-masing menyambung satu Section tersendiri *(§B1)*; nol persinggungan dengan alur klaim |
| **0002** RBAC tiga peran | ⛔ **tidak menyentuh** | korpus nol memuat model peran |
| **0003** uang non-float | ⚠️ **MENENTANG sebagian** | kolom uang `NUMBER` ✅, tetapi ⭐ **`RIRATE` = `VARCHAR2(1000)`** — angka rate disimpan sebagai teks |
| **0004** endpoint env-var | ⛔ **tidak menyentuh** | `superseded`; dan modul ini **nol ConnectREST** *(jendela: 66 berkas, folder `ConnectREST` tidak ada)* |
| **0005** flag lingkungan | ⚠️ **tidak dapat diputuskan** | `IsPEGAPROD` **tidak ketemu** di 66 berkas — tetapi lihat silang di bawah |
| **0006** penomoran via procedure | ✅ **MENGUATKAN** | ⭐ ID diciptakan DB: `'1'+lpad(seq.nextval,6,'0')`, sequence per tabel · `TREATYYEAR_LIFE_SEQ START WITH 44` |
| **0007** jejak audit tiap transisi | ⚠️ **MENENTANG** | procedure hanya menyimpan `USERID`; ⭐ `p_TGLUPDATE` **diabaikan**, diganti `SYSDATE` — lihat D3 |
| **0008** efek keluar asinkron | ⛔ **tidak menyentuh** | nol efek keluar |
| **0009** migrasi penuh Life | ⛔ **tidak menyentuh** | lingkupnya data klaim Life |
| **0010** berkas di Google Storage | ⛔ **tidak menyentuh** | nol penyimpanan berkas |
| **0011** unit status = baris `AdjustmentList` | ⛔ **tidak menyentuh** | modul master, nol `AdjustmentList` |
| **0012** wewenang kirim komite | ⛔ **tidak menyentuh** | nol komite |
| **0013** endpoint via `M_LINK_SERVICE` | ⛔ **tidak menyentuh** | nol ConnectREST |
| **0014** pemutus komite per `KomiteID` | ⛔ **tidak menyentuh** | nol komite |
| **0015** outbox wajib berhasil | ⚠️ **MENENTANG** | ⭐ `o_message` berisi galat, **nol activity membacanya** *(§C2)* — kegagalan simpan hilang tanpa jejak |

⭐ **Menguatkan 2 · MENENTANG 3 · tidak menyentuh 9 · tidak dapat diputuskan 1.**
**Nomor yang ditentang: 0003 · 0007 · 0015.**

### D3 — ADR-0007 diperiksa lebih dalam

Kelima procedure menerima `p_USERID` **dan** `p_TGLUPDATE`. ⭐ **`p_TGLUPDATE` diabaikan** —
`TGLUPDATE` selalu `SYSDATE`.

**Apa yang KURANG terhadap ADR-0007** *(merekam **siapa** dan **kapan** untuk **setiap** transisi
status **dan setiap jalur balik**)*:

| Tuntutan ADR-0007 | Keadaan modul ini |
| --- | --- |
| **siapa** | ✅ ada — `USERID` |
| **kapan** | ⚠️ ada, tetapi **waktu basis data**, bukan waktu aksi pengguna |
| **setiap transisi status** | ⛔ **tidak ada** — modul master tidak punya status; yang tersimpan hanya **keadaan terakhir**, satu baris ditimpa |
| **setiap jalur balik** | ⛔ **tidak ada** — upsert **menimpa**, riwayat sebelumnya **hilang** |

⭐ **Kesimpulan: bukan jejak audit, melainkan stempel penyunting terakhir.**

### D4 — ADR-0005 diperiksa ulang sebagai klaim ketiadaan

**Jendela (S1):** 66 berkas; penyaring **tidak peka huruf besar-kecil** atas `IsPEGAPROD` /
`pzProductionLevel`. ⭐ **Tidak ketemu.**

⚠️ **Silang dengan A6 — dan keduanya tidak bertemu.** `IsPEGAPROD` menguji **tingkat produksi
proses berjalan**; `pxUpdateSystemID` dan `pxHostId` mencatat **tempat rule terakhir disimpan**.
⛔ Ketiadaan `IsPEGAPROD` **tidak menerangkan** kenapa modul ini 80% `pegadevnusare2`.

### D5 — ADR-0004 / 0013 diperiksa ulang

**Jendela (S1):** 66 berkas, 7 jenis folder. ⭐ **Folder `ConnectREST` tidak ada** di modul ini, dan
penyaringan `Connect-REST` atas seluruh medan metode **tidak ketemu**. ✅ **Klaim ronde 1 bertahan.**

### D6 — ADR-0001 sesudah §B1

⭐ **Batas "konteks/menu sendiri" BERTAHAN, dan menguat.** Keempat Harness masing-masing menyambung
**satu** Section milik modul ini sendiri; nol menyambung Section modul lain.

---

## §E — Sebelas pertanyaan ronde 1

⭐ **Fakta yang menentukan §E:** `grilling-ronde-1-jawaban.md` **sudah memuat jawaban untuk
Q1–Q11 seluruhnya** — Q1 · Q4 · Q5 *(DIREVISI)* · Q8 · Q9 · Q10 · Q11 sebagai
`[keputusan work owner]`, dan Q2 · Q3 · Q6 · Q7 sebagai `[fakta bisnis — work owner]`.

⛔ **Karena itu §E berubah bentuk:** bukan *"mana yang kini terjawab"*, melainkan ***"apakah bahan
baru mengubah atau memperkuat jawaban yang sudah ada"***. ⛔ **Nol Q dinyatakan tertutup** — itu
wewenang work owner.

| Q | Pertanyaannya | Status ronde 1 | Bahan baru mengubahnya? |
| --- | --- | --- | --- |
| **Q1** | Nama jenis/tahun disimpan berulang di beberapa tabel — disengaja? | dijawab, keputusan | ⛔ tidak |
| **Q2** | Enam kolom batas pada kontrak — artinya apa masing-masing | dijawab, fakta bisnis | ⭐ **menguat** — DDL memastikan keenamnya `NUMBER`, jadi keenamnya **angka**, bukan teks |
| **Q3** | Dua tingkat share — reinsurer lalu security reinsurer | dijawab, fakta bisnis | ⛔ tidak |
| **Q4** | Total share wajib 100% atau hanya ditampilkan | dijawab, keputusan | ⭐ **menguat tegas** — §A3 membuktikan **nol penjaga** di Pega *(21 baris mati)*, dan `[data DBA]` membuktikan **nol penjaga** di procedure. **Tidak ada penjaga di kedua sisi** |
| **Q5** | Menghapus induk yang punya anak | dijawab, **DIREVISI** | ⭐ **BERUBAH ARTINYA** — §C5: basis data **kini mengaskade sendiri** |
| **Q6** | Tahun treaty tidak dapat dihapus — benar? | dijawab, fakta bisnis | ⛔ tidak |
| **Q7** | Enumerasi jenis reasuransi | dijawab, fakta bisnis | ⛔ tidak |
| **Q8** | Gerbang konsistensi tahun | dijawab, keputusan | ⭐ **DASAR FAKTUALNYA BERUBAH** — §A3: gerbangnya **MATI**, bukan sekadar keliru cara |
| **Q9** | Lima procedure dipakai apa adanya | dijawab, keputusan | ⭐ **menguat** — body kini di tangan; upsert, ID dari DB, `COMMIT` internal |
| **Q10** | Fitur "terapkan ke semua" | dijawab, keputusan | ⚠️ **satu baris gerbangnya MATI** *(`SaveBusinessToAllLife_Act` langkah 3)*, dua lainnya hidup |
| **Q11** | Galat procedure diabaikan | dijawab, keputusan | ⭐ **menguat tegas** — §C2: **nol activity membaca `HASIL1`**, dan kini terbaca **apa** yang hilang |

⭐ **Bahan baru mengubah atau menguatkan 6 dari 11** — Q2 · Q4 · Q5 · Q8 · Q9 · Q11.
⛔ **Nol dinyatakan tertutup.**

---

## §F — Dua berkas tambahan

✅ Ditulis sebagai berkas **BARU**, keduanya menunjuk balik ke bagian yang ditambalnya:

- `TAMBAHAN-SPEC-ronde-2.md`
- `issues/TAMBAHAN-TIKET-ronde-2.md`

⛔ `spec.md` dan `issues/01…12` **tidak disunting**.

---

## §G — Pertanyaan, lengkap dan dibedah

### G5 — Yang dibuang karena sudah terjawab: ⭐ **SEBELAS**

⛔ Q1 sampai Q11 **seluruhnya** sudah dijawab di `grilling-ronde-1-jawaban.md` dan **tidak diulang**.

### G1 · G2 · G3 — Pertanyaan yang masih berdiri: **4**

#### ⭐ Pertanyaan A — `RIRATE` disimpan sebagai teks. Apa sebenarnya isinya?

| | |
| --- | --- |
| **1 · Apa yang ditanyakan** | Kolom `RIRATE` pada tabel business menyimpan **teks sepanjang 1000 karakter**, bukan angka. Apa yang sebenarnya diketik pengguna ke situ? |
| **2 · Kenapa muncul** | DDL dari DBA mencatatnya sebagai kejanggalan, dan ronde 1 menandai `RIRATE` *"belum saya telusur"*. Ronde ini menemukan **4 penulis dan 8 pembaca** — jadi ia dipakai, bukan sisa |
| **3 · Bedanya kalau A atau B** | **A — ia satu angka** *(mis. `0.5`)*: kolomnya menjadi desimal, tunduk aturan ketelitian, dan layar memvalidasinya. **B — ia teks berstruktur** *(mis. rate bertingkat per usia, atau `"0.5%"` berikut satuannya)*: kolomnya **tetap teks**, tidak boleh dipaksa jadi angka, dan layar memerlukan bentuk masukan tersendiri. ⚠️ Data lama **tidak dapat dimigrasikan** tanpa tahu bentuknya |
| **4 · Apa yang MACET** | ⭐ **Tiket 07** *(business dan tampilan rate)* macet: bentuk kolom, validasi, dan migrasi nilainya ketiganya bergantung jawaban ini |

**Golongan:** ⭐ **FAKTA BISNIS** — ⛔ jangan ditebak, nol rekomendasi.
**Pemilik:** **Product + UW**, dengan **DBA** untuk memastikan isi nyata kolomnya.

#### ⭐ Pertanyaan B — Sebelas validasi simpan bergerbang mati. Dibangun atau tidak?

| | |
| --- | --- |
| **1 · Apa yang ditanyakan** | **21 dari 42** baris syarat di jalur simpan bergerbang **mati** — termasuk seluruh pemeriksaan nama reinsurer kosong, tanggal mulai kosong, dan konsistensi tahun. Apakah sistem baru **membangun** pemeriksaan itu, atau **meniru keadaan sekarang** yang tanpa pemeriksaan? |
| **2 · Kenapa muncul** | Ronde 1 tidak pernah membaca flag prakondisi. Ronde ini membacanya, dan menemukan pemeriksaannya **tertulis tetapi tidak berjalan** |
| **3 · Bedanya kalau A atau B** | **A — dibangun:** simpan ditolak bila nama reinsurer atau tanggal kosong; ⚠️ **data lama yang kosong menjadi tidak dapat disimpan ulang**. **B — ditiru apa adanya:** simpan selalu diterima, sama seperti sekarang; ⚠️ dan baris tanpa nama reinsurer **tetap boleh lahir** |
| **4 · Apa yang MACET** | ⛔ **Tidak ada yang macet.** Keduanya dapat dibangun; yang berubah hanya **perilaku yang dijanjikan**. Tetapi ia menyentuh **tiket 05, 06, 07, dan 10** |

**Golongan:** **KEPUTUSAN DESAIN** — boleh diberi rekomendasi.
⭐ **REKOMENDASI:** **A — dibangun**, dengan dua syarat: *(i)* dinyatakan sebagai **penyimpangan
sadar** di spec, dan *(ii)* migrasi data lama **tidak menolak** baris yang sudah terlanjur kosong.
Alasannya: keduanya pemeriksaan yang **penulisnya sendiri sudah menuliskannya** — mematikannya
tampak seperti keputusan sementara, bukan keputusan rancangan. ⛔ **Ini rekomendasi, bukan keputusan.**
**Pemilik:** **work owner**.

#### Pertanyaan C — Modul ini 80% tercap `pegadevnusare2`. Ekspornya dari lingkungan mana?

| | |
| --- | --- |
| **1 · Apa yang ditanyakan** | Dari 66 berkas: **53** tercap `pegadevnusare2`, **12** `pega`, **1** `pegaprdnusare`. Dan `pxHostId` punya **lima** nilai, dua di antaranya *(`pega-nusre`, `jboss122117`)* belum pernah dinyatakan. Ekspor ini diambil dari lingkungan mana? |
| **2 · Kenapa muncul** | Sensus lingkungan **belum pernah dikerjakan** di modul ini. Di modul lain sebarannya **berlawanan** *(Claim Fac In 66% `pega`)* |
| **3 · Bedanya kalau A atau B** | **A — ekspornya dari dev:** yang dibaca ronde 1 dan 2 adalah **rule yang belum tentu berjalan di produksi**, dan seluruh kesimpulan modul ini perlu diberi catatan itu. **B — campuran sah karena mirroring:** sebarannya tidak berarti apa-apa, dan tidak ada yang perlu diubah |
| **4 · Apa yang MACET** | ⛔ **Tidak ada yang macet sekarang.** Tetapi bila jawabannya A, **seluruh spec modul ini** perlu diberi pernyataan lingkup |

**Golongan:** **FAKTA** — ⛔ jangan ditebak.
**Pemilik:** **IT-infra**, dengan **work owner** untuk memutuskan akibatnya.

#### Pertanyaan D — `DeleteRowBusiness` memakai `HASIL12`, bukan `HASIL1`. Salah ketik atau penampung lain?

| | |
| --- | --- |
| **1 · Apa yang ditanyakan** | Satu-satunya activity yang menyentuh penampung galat memakai **`HASIL12`**; kelima rule SQL simpan memakai **`HASIL1`**. Apakah itu salah ketik, atau `HASIL12` memang penampung yang berbeda? |
| **2 · Kenapa muncul** | Ronde 1 menyebut `DeleteRowBusiness.xml` sebagai satu-satunya pembaca `HASIL1`. Ronde ini menemukan yang dipakainya **bukan `HASIL1`** |
| **3 · Bedanya kalau A atau B** | **A — salah ketik:** jalur hapus **memang** membaca galat, dan hanya penamaannya keliru. **B — penampung berbeda:** jalur hapus membaca **sesuatu yang lain**, dan ⭐ **nol jalur di seluruh modul membaca galat simpan** |
| **4 · Apa yang MACET** | ⛔ **Tidak macet**, tetapi ia menentukan seberapa tegas **tiket 10** *(penegakan `HASIL1`)* harus berbunyi |

**Golongan:** **FAKTA** — terbaca dari korpus bila ditelusuri lebih dalam; ⛔ belum saya telusuri.
**Pemilik:** **work owner** — atau **korpus**, bila ronde berikutnya membacanya.

### G4 — Urutan menahan

| # | Pertanyaan | Yang tertahan |
| --- | --- | --- |
| **1** | ⭐ **A — `RIRATE` teks** | **tiket 07** macet: bentuk kolom, validasi, dan migrasi |
| **2** | **B — 21 validasi mati** | menyentuh **tiket 05 · 06 · 07 · 10**; tidak macet, tetapi mengubah janji perilaku |
| **3** | **C — lingkungan ekspor** | tidak macet; menentukan apakah spec perlu pernyataan lingkup |
| **4** | **D — `HASIL12`** | tidak macet; menentukan ketegasan **tiket 10** |

---

## §H — Apa lagi

### H1 — Sesudah ronde ini

| # | Yang dikerjakan | Kenapa menahan |
| --- | --- | --- |
| **1** | ⭐ **Baca 4 Section yang belum pernah disebut** *(1,69 MB)* | layar utama modul ini justru yang belum dibaca |
| **2** | ⭐ **Baca 11 ReportDefinition** dengan medan yang benar | ronde ini **gagal** membaca kelas, kolom, dan saringannya |
| **3** | **Baca `ViewRate` / `ViewRateTable`** dan sumber datanya | B2 hanya sampai dugaan dari nama |
| **4** | **Temukan medan kolom Harness** *(B1)* | ronde ini gagal di `pyPropertyName` |
| **5** | **Petakan `CARI1…CARI14` ke parameter procedure** | tanpa itu, C1 *"parameter cocok semua?"* tidak terjawab |
| **6** | **Telusuri `HASIL12`** *(Pertanyaan D)* | terbaca dari korpus, belum ditelusuri |
| **7** | **Periksa `.xlsx` di akar modul** | berkas non-rule di luar `Claude outputs\`, belum pernah tercatat |

### H2 — Kesimpulan SAYA SENDIRI yang paling rawan salah

⛔ **Ditunjuk, tidak diperbaiki.**

> ⚠️ **Paling rawan: §A3 — "21 dari 42 baris bersyarat bergerbang MATI".**

Ia bersandar penuh pada aturan **C1**, yaitu bahwa `pyStepsPreCondition = false` berarti gerbangnya
**tersimpan tetapi dimatikan**. ⭐ Aturan itu ditulis dari pengamatan di **modul lain**, dan
⛔ **belum pernah diuji di modul ini** — misalnya dengan melihat layar Pega seperti yang dilakukan
work owner untuk penanda `//`.

⚠️ Bila `false` di sini ternyata berarti sesuatu yang lain, **21 baris berbalik menjadi hidup**,
dan temuan terbesar §A **runtuh** — berikut ralat §A3 terhadap spec §8.

**Paling rawan kedua:** §C3 *"nol kolom ditulis Pega yang tidak ada di DDL"*. Saya sendiri
menyatakan metodenya **tidak dapat membuktikannya**, karena Pega tidak memakai nama kolom.
⚠️ Menuliskannya sebagai "nol ditemukan" tetap berisiko dibaca sebagai "tidak ada".

**Paling rawan ketiga:** §A4 memilih **94** sebagai angka yang dipakai seluruh §A, dengan alasan
pengurai XML setuju dengan jangkar A. ⚠️ Aturan **P4** menyatakan dua angka berbeda berarti
**belum punya data** — dan saya tetap memakai salah satunya untuk menghitung §A2 dan §A3.

### H3 — RALAT terhadap ronde 1

⛔ **Angka lama dikutip semua.** ⛔ `grilling-ronde-1.md` **tidak disentuh.**

| # | ⛔ Pernyataan LAMA *(ronde 1)* | ⭐ Pernyataan BARU *(ronde 2)* | Sebabnya |
| --- | --- | --- | --- |
| **1** | `CountingPercentShare_Act` **4 langkah** | ⭐ **5 langkah** | langkah bersarang **3.1** tidak terhitung |
| **2** | `TreatyLimit_TypeProtect` **4 langkah** | ⭐ **5 langkah** | langkah bersarang **4.1** tidak terhitung |
| **3** | *"hanya `DeleteRowBusiness.xml` yang menyebut `HASIL1`"* | ⭐ **7 berkas dari 66**, termasuk kelima rule SQL simpan | jendela ronde 1 terlalu sempit |
| **4** | *(tidak pernah dinyatakan)* — gerbang tidak pernah dibaca | ⭐ **21 dari 42 baris bersyarat bergerbang MATI** | aturan C1 belum ada saat ronde 1 |
| **5** | *(tidak pernah dinyatakan)* — langkah ber-remark tidak pernah diperiksa | ⭐ **1 langkah ber-`//`** di seluruh modul | aturan F2 belum ada saat ronde 1 |
| **6** | `"66 berkas"` | **66 rule** ✅ **COCOK**, tetapi folder berisi **69** berkas | tiga berkas bukan rule tidak terhitung |
| **7** | §C *"nol kaskade"* *(dari sisi Pega)* | ✅ **tetap benar untuk Pega**, ⭐ tetapi basis data **kini punya 4 FK `ON DELETE CASCADE`** | DDL datang sesudah ronde 1 |
| **8** | *(tidak pernah dinyatakan)* — sensus lingkungan belum ada | ⭐ `pxUpdateSystemID` **3 nilai** · `pxHostId` **5 nilai** | belum pernah dikerjakan |

⭐ **Delapan butir ralat. Nol angka lama dihapus.**

⚠️ **Satu pelurusan terhadap BRIEF, bukan terhadap ronde 1:** brief menyebut kesimpulan ronde 1
sebagai *"ketiga rumus itu BUKAN rumus"*. ⛔ **Ronde 1 menulis "dua di antaranya"**, dan
menyebut `CountingPercentShare_Act` sebagai **penjumlah**. Parafrase brief lebih keras dari
teks aslinya.

### H4 — ⭐ Yang seharusnya dikerjakan tetapi TIDAK diperintahkan blok ini

⛔ **Disebutkan, tidak dikerjakan.**

| # | Butir |
| --- | --- |
| **1** | ⭐ **Menguji arti `pyStepsPreCondition = false` di layar Pega** — seperti yang sudah dilakukan untuk `//`. Temuan terbesar ronde ini bergantung padanya *(H2)* |
| **2** | **Menyisir `pxUpdateSystemID` dan `pxHostId` di 20 modul korpus lain** — tiga modul kini menunjukkan sebaran yang sangat berbeda |
| **3** | **Memeriksa apakah `STDRATING` dipakai untuk hal lain** — total share ditulis ke medan bernama *standard rating* |
| **4** | **Mengadu Claim Prop dan Komite Claim Prop dengan 15 ADR** — keduanya selesai tanpa pernah diadu dengan dokumen |
| **5** | **Memeriksa keempat rule `@BASECLASS` modul ini di modul lain** — identitas empat bagian menunjukkan modul ini hanya punya satu ruleset, tetapi `@BASECLASS` dipakai bersama |
| **6** | **Menanyakan ke DBA apakah `NUMBER` tanpa presisi memang disengaja** — DDL tidak membatasi apa pun, sehingga batas 20/8 murni keputusan aplikasi |
| **7** | **Memeriksa apakah 4 FK `ON DELETE CASCADE` sudah ada saat data lama dibuat** — menentukan apakah data yatim lama masih ada |

---

## Lampiran — bukti berkas lama tidak disentuh

MD5 atas **34 berkas** diambil di LANGKAH 0, sebelum pekerjaan apa pun:
kedua berkas grilling lama · `spec.md` · kedua berkas DBA · **12 tiket** · **15 ADR** ·
`ATURAN-BACA-KORPUS-PEGA.md` · `CONTEXT.md`.

**Tiga berkas baru:** `grilling-ronde-2.md` · `TAMBAHAN-SPEC-ronde-2.md` ·
`issues/TAMBAHAN-TIKET-ronde-2.md`.
