# Tiket - Penempatan Treaty Keluar

> Dokumen ini memuat **badan tiket lengkap**, disusun per modul lalu per nomor.
> Disusun 25 September 2026 dari berkas tiket proyek migrasi Nusantara Re.

## Matriks status

| Modul | Tiket | Siap | Tertahan | needs-info | wontfix | Lain |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Treaty Contract Out | **12** | 12 | 0 | 0 | 0 | 0 |
| **Jumlah** | **12** | **12** | **0** | **0** | **0** | **0** |

---

# Treaty Contract Out

Jumlah tiket: **12**

## Treaty Contract Out - 01 - Skema relasional + migrasi + tipe dirapikan — **PREFACTOR**

**Status:** ready-for-agent

**Blocked by:** **CL-01** (kerangka aplikasi + seam API — scaffolding lintas konteks, tidak dibuat
di sini)

⚠️ **Ini PREFACTOR, dan ia tiket PERTAMA.** Bentuk skema berubah di tiga sumbu sekaligus —
tabel JSON dibuang, `MTREATYSECURITY` dibersihkan, dan tipe kolom diperbaiki — sehingga **tidak ada
irisan lain yang dapat berdiri** sebelum bentuk barunya ada. *"Make the change easy, then make the
easy change."*

#### Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin seluruh data master arrangement treaty non-life pindah ke bentuk
yang dapat dipercaya — uang sebagai angka, tanggal sebagai tanggal, setiap baris punya identitas —
tanpa kehilangan satu nilai pun, dan **tanpa mengambil apa pun dari tabel JSON yang sudah mati**.
*(User story 37–40 di spec)*

Dan sebagai **konteks hilir** (`Claim Prop`, `Komite Claim Prop`, `Claim Fac In`), saya ingin tetap
dapat **membaca** master arrangement dalam bentuk yang saya kenal sampai saya ikut bermigrasi.

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL delapan tabel + sequence; skrip migrasi & rekonsiliasi |
| `internal/models` | Bentuk tahun treaty, kontrak, reinsurer, security, business, klausul |
| — | Skrip rekonsiliasi nilai uang & tanggal |

#### Keadaan lama `[data DBA]`

| Tabel | Temuan tipe |
| --- | --- |
| `PROPORTIONALARRG` | ⚠️ **CAMPUR**: `TREATYLIMIT`, `COINS_MIN`, `COINS_MAX`, `MORERP`, `MOREUSD` = **`NUMBER`**; tetapi `RP`, `USD`, `PCT`, `PCTME` = **`VARCHAR2(1000)`**. Ada `TGLUPDATE DATE`, `OBJECT VARCHAR2(50)`, `PROPORTIONALLIST VARCHAR2(1000)` |
| `MTREATYSECURITY` | ⚠️ `THN_TREATY VARCHAR2(4) DEFAULT '1' NOT NULL`, `TOP_ID VARCHAR2(9)`, `TP_TREATY CHAR(2)`, `REAS_ID CHAR(7) NOT NULL`, `PCT_SHARE VARCHAR2(99)`, `USER_ID CHAR(99)`, `REAS_SECURITY CHAR(10) NOT NULL` — **tanpa primary key** |
| `TREATYCONTRACT` | `TREATYSTARTDATE`/`TREATYENDDATE` sudah `DATE`; tetapi `TGLUPDATE VARCHAR2(1000)` |
| `TREATYYEAR` | **seluruh** kolom `VARCHAR2` — termasuk `STARTDATE`, `ENDDATE` |
| `TREATYBUSINESS` | **seluruh** kolom `VARCHAR2` |
| `TREATYREINSURER` | `RICOMM`, `PCTSHARE` = `NUMBER`; sisanya `VARCHAR2` |
| `TREATYEXCHANGEYEARLY` | **seluruh** kolom `VARCHAR2` |
| `TREATYDESC` | `ID, DESCNAME, ISXOL, STATUSAKTIF` — seluruhnya `VARCHAR2` |

#### Bentuk baru

```
treaty_year   (ID + grup + underwriting year + proporsi + STARTDATE/ENDDATE sebagai DATE)
  └─ treaty_contract  (IDTREATYYEAR → treaty_year.ID, REINSTYPEID, tanggal sebagai DATE)

Menggantung pada kunci gabungan (TREATYYEAR, TREATYGROUPID, REINSTYPEID) — BUKAN FK ke kontrak:
  ├─ treaty_reinsurer      (PCTSHARE, RICOMM sebagai desimal)
  │    └─ mtreaty_security (PK surrogate, REAS_ID → treaty_reinsurer.ID, PCT_SHARE desimal)
  ├─ treaty_business       (BIZCODE, BIZNAME, ISACTIVE)
  └─ proportionalarrg      (SATU tabel untuk 25 jenis klausul, dibedakan TREATYDESCID)
```

⚠️ **Kunci gabungan dipertahankan** `[fakta bisnis — work owner]`. Klausul dan anak-anak lain
menggantung pada **(TreatyYear, TreatyGroupID, ReinsTypeID)**, **bukan** pada `ID` kontrak.
Menambahkan foreign key ke kontrak akan **mengubah arti data** — jangan lakukan.

##### Kolom `proportionalarrg` — 35 kolom `[data DBA]`

```
ID, TREATYYEAR, TREATYYEARID, TREATYGROUPID, TREATYGROUPNAME, TREATYDESCID, TREATYDESCNAME,
REINSTYPEID, REINSTYPENAME, LAYER, LAYERPART, LAYERPARTTYPE, LAYERTYPE, KURS, TGLUPDATE, USERID,
LINE, PCT, PCTME, YDCF, METHOD, TERRITORIALLIMIT, PARENTREINSTYPEID, SPREADINGORDER, RP, USD,
ID_OCCUPATION, OCCUPATION, ID_CLAUSE, CLAUSE, TREATYLIMIT, COINS_MIN, COINS_MAX, MORERP, MOREUSD
```

Sembilan kolom terakhir yang bercetak — `ID_OCCUPATION`, `OCCUPATION`, `ID_CLAUSE`, `CLAUSE`,
`TREATYLIMIT`, `COINS_MIN`, `COINS_MAX`, `MORERP`, `MOREUSD` — hanya diisi baris **induk**; baris
"anak" mengisinya **NULL**.

##### Sequence `[data DBA]`

| Tabel | Sequence | Format identitas |
| --- | --- | --- |
| `PROPORTIONALARRG` | `PROPORTIONALARRG_SEQ` | `'1' + lpad(seq, 7, '0')` — **7 digit** |
| `TREATYCONTRACT` | `treatycontract_seq` | `'1' + lpad(seq, 6, '0')` |
| `TREATYYEAR` | `TreatyYear_seq` | `'1' + lpad(seq, 6, '0')` |
| `TREATYREINSURER` | `M_TREATYREINSURER_SEQ` | `'1' + lpad(seq, 6, '0')` |
| `TREATYBUSINESS` | `TREATY_BUSINESS_SEQ` | `'1' + lpad(seq, 6, '0')` |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterProportionalArrg` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `ASM!SAVEMASTERPROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterProportionalArrg.xml` | 35 kolom induk |
| `SaveMasterProportionalArrgChild` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterProportionalArrgChild.xml` | 26 kolom anak |
| `GetMasterDescriptionLimitParentList` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterDescriptionLimitParentList.xml` | baca **relasional** |
| `GetMasterDescriptionEPIParentList` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterDescriptionEPIParentList.xml` | ⚠️ baca **JSON** — mati |
| `DeleteRowBusinessList` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `ASM!DELETEROWBUSINESSLIST` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteRowBusinessList.xml` | ⚠️ memegang **kedua** tabel kembar |
| `InsertToMTreatySecurity` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `ASM!INSERTTOMTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/InsertToMTreatySecurity.xml` | ⚠️ INSERT posisional |

⚠️ **Penyimpangan sadar 1 — dualitas JSON dibuang.** `[keputusan work owner]`
`[data DBA]` **Tidak satu pun dari enam procedure penulis menyentuh tabel JSON.**
`M_PROPORTIONALARRG` sudah lama tidak dipakai; kueri yang masih membacanya adalah sisa yang lupa
dihapus. Sumber kebenaran tunggal = **`PROPORTIONALARRG` relasional**.

⚠️ **Penyimpangan sadar 5 — `MTREATYSECURITY` dibersihkan.** `[keputusan work owner]`
PK surrogate, kolom bernama dan diisi eksplisit, `PCT_SHARE` desimal, `REAS_SECURITY` atribut biasa.

⚠️ **Penyimpangan sadar 6 — seluruh uang/persen jadi desimal, seluruh tanggal jadi `DATE`.**
`[keputusan work owner]`

#### ADR terkait

**ADR-U-0003** (uang non-float), **ADR-U-0006** (identitas lewat sequence basis data),
**ADR-U-0009** (migrasi penuh; koeksistensi ditolak untuk data — **bukan** untuk kontrak baca hilir).

#### Acceptance criteria

- [ ] Modul ini adalah **satu-satunya penulis** `TREATYYEAR`, `TREATYCONTRACT`, `TREATYREINSURER`,
      `MTREATYSECURITY`, `TREATYBUSINESS`, `PROPORTIONALARRG`. *(AC 1 spec; OQ-042)*
- [ ] Konteks hilir (`Claim Prop`, `Komite Claim Prop`, `Claim Fac In`) tetap dapat **membaca**
      bentuk relasional yang mereka pakai hari ini. *(AC 2 spec)*
- [ ] Skema **tidak memuat** satu pun objek treaty **outward**. *(AC 3 spec)*
- [ ] ⚠️ **Seluruh** nilai uang — `RP`, `USD`, `MORERP`, `MOREUSD`, `TREATYLIMIT`, `COINS_MIN`,
      `COINS_MAX`, `PCT_SHARE`, `RICOMM` — bertipe **desimal presisi arbitrer**; **tidak** melewati
      `float`. Test yang menemukan kolom uang bertipe teks **gagal**. *(AC 51 spec; **ADR-U-0003**;
      penyimpangan sadar 6)*
- [ ] ⚠️ **Seluruh** persentase — `PCT`, `PCTME` — desimal, tidak dibulatkan ke bilangan bulat dan
      tidak disimpan sebagai teks. *(AC 52 spec; penyimpangan sadar 6)*
- [ ] ⚠️ **Seluruh** tanggal — `STARTDATE`, `ENDDATE`, `TREATYSTARTDATE`, `TREATYENDDATE`,
      `TGLUPDATE` — bertipe **`DATE`**. *(AC 53 spec; penyimpangan sadar 6)*
- [ ] ⚠️ Skema target **relasional penuh**. Test yang menemukan kolom JSON sebagai penyimpan
      atribut arrangement **gagal**. *(AC 63 spec; penyimpangan sadar 1)*
- [ ] ⚠️ Migrasi **tidak mengambil apa pun** dari `M_PROPORTIONALARRG`, `M_TREATYCONTRACT`,
      `M_TREATYBUSINESS`, atau tabel `M_*` lain. Sumbernya **hanya** tabel relasional.
      *(AC 64 spec; penyimpangan sadar 1)*
- [ ] Seluruh baris keenam tabel pindah **tanpa kehilangan satu nilai pun**. *(AC 65 spec;
      **ADR-U-0009**)*
- [ ] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara
      tepat**, bukan dengan toleransi. *(AC 66 spec; **ADR-U-0003**)*
- [ ] ⚠️ Nilai uang & persen yang hari ini berupa **teks berkoma desimal** (mis. `"12,5"`) terurai
      **benar** menjadi desimal; teks yang **tidak dapat diurai** dilaporkan, **tidak** didiamkan
      dan **tidak** diam-diam jadi nol. *(AC 66 spec; `[terverifikasi]` — existing memanggil
      `@replaceAll(.Pct,",",".")` di `Activity/TreatyTestChildTotal_Act.xml`,
      `@BASECLASS!TREATYTESTCHILDTOTAL_ACT`, dan `Activity/SetErrorMessageReinsurer.xml`)*
- [ ] Tanggal yang hari ini berupa **teks** menjadi `DATE` **tanpa pergeseran zona waktu**; teks
      yang tidak dapat diurai **dilaporkan**. *(AC 67 spec)*
- [ ] ⚠️ Baris `MTREATYSECURITY` mendapat **primary key surrogate** saat migrasi, dan rujukannya ke
      reinsurer tetap utuh. Test yang menemukan kunci berbasis nama security **gagal**.
      *(AC 68 spec; penyimpangan sadar 5)*
- [ ] Kelima sequence pindah dengan **nilai berjalan yang benar**, sehingga identitas baru **tidak
      bertabrakan** dengan yang lama. *(AC 69 spec; **ADR-U-0006**)*
- [ ] `[terbuka]` Kolom `PROPORTIONALLIST` dan `OBJECT` **tidak dibawa** ke skema baru kecuali
      migrasi membuktikan ada data hidup di sana; temuannya **dilaporkan**. *(AC 70 spec)*
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.
      *(AC 71 spec)*

#### Blocker

**Tidak ada pemblokir.** `[data DBA]` DDL delapan tabel dan body enam procedure sudah diterima —
**OQ-001 dan OQ-002 ditutup**.

#### Catatan

⚠️ **Kejanggalan yang wajib dicatat, bukan ditiru.**

- `RDBList/DeleteRowBusinessList.xml` menghapus `treatybusiness` di `pyBrowseSQL` **dan**
  `m_treatybusiness` di `pyDeleteSQL` — satu rule, dua tabel kembar. Sebaliknya
  `RDBList/DeleteFromTREATYCONTRACT_SQL.xml` menghapus `treatycontract` tetapi **bukan**
  `M_TREATYCONTRACT`. Inkonsistensi ini **lenyap dengan sendirinya** begitu JSON dibuang.
- `[data DBA]` Procedure bernama **`PEGA_M_PROPORTIONALARRG_CHILD`** **tidak** menulis ke tabel
  child — ia menulis ke tabel yang **sama**. Nama menipu (**OQ-066**).
- `[data DBA]` Teks galat sebagian procedure menyebut **`"JSON_KLAIM"`** — sisa salin-tempel
  template, **bukan** petunjuk bahwa data berbentuk JSON.
- `[terverifikasi]` `Activity/GetPeriode.xml` (`@BASECLASS!GETPERIODE`) memakai
  `ParentReinsTypeID = "00"` sebagai penanda baris tanpa induk. Nilai sentinel ini perlu dibawa
  utuh saat migrasi.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — presisi desimal, konversi tanggal, dan
keutuhan rujukan setelah pemberian PK **hanya berperilaku benar pada basis data sungguhan**;
memalsukannya berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Treaty Contract Out - 02 - Jenis reasuransi — master dibaca + saringan non-life

**Status:** ready-for-agent

**Blocked by:** 01 (skema harus ada)

#### Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin memilih **jenis reasuransi** dari daftar master — dan
daftar itu hanya memuat jenis yang **berlaku untuk non-life** — sehingga saya tidak pernah mengetik
bebas dan tidak pernah salah memilih jenis milik lini lain. *(User story 8–10 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Baca master jenis reasuransi (read-only) |
| `internal/services` | Saringan non-life |
| `internal/handlers` | Endpoint daftar jenis reasuransi |
| `frontend/` | Pemilih jenis reasuransi yang dipakai layar kontrak dan seluruh grid klausul |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| **RD dominan** | `ASM-FW-GISFW-INT-REINSURANCETYPE` / `BROWSEREINSURANCETYPE_RD_OLD_LJT_ID_ISNOTNULL` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml` | ⚠️ dipakai **11 grid** |
| RD kedua | `ASM-FW-GISFW-INT-REINSURANCETYPE` / `BROWSEREINSURANCETYPE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseReinsuranceType_RD.xml` | dipakai 1 layar — **layar master**, di luar konteks ini |
| `GetMasterReinsTypeContract` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterReinsTypeContract.xml` | ⚠️ membaca tabel **JSON** — mati |

⚠️ **Nama berbohong — dan yang "Old" justru yang benar** `[terverifikasi]`:

| Rule | RuleSet | Commit | Dipakai |
| --- | --- | --- | ---: |
| `…!BROWSEREINSURANCETYPE_RD_OLD_LJT_ID_ISNOTNULL` | **01-01-91** | **2026-02-05** | **11** grid |
| `…!BROWSEREINSURANCETYPE_RD` | 01-01-83 | 2025-06-02 | 1 layar |

`[keputusan work owner]` Yang **dimigrasikan** adalah perilaku yang dipakai 11 grid:

```
.ID NOT IN ("10004","10011","10012","10021","10022","10025","10026","10028",
            "10248","10249","10018","10217")
AND .Flag = "active"
AND .Type IN ("1","2","3")
```

⚠️ **Penyimpangan sadar 8 — nama jujur.** `[keputusan work owner]` Komponen di sistem baru
**tidak** mengandung kata `"Old"`.

`[terverifikasi]` Class `ASM-FW-GISFW-INT-REINSURANCETYPE` adalah **class yang sama** dengan master
jenis reasuransi Life (`Master Contract Retro Life/ReportDefinition/BrowseReinsuranceTypeLimit_RD.xml`).
Jadi master ini **melayani life dan non-life sekaligus**; saringanlah yang memisahkan.

#### ADR terkait

**ADR-U-0015** (kegagalan ditangani eksplisit, tidak ditelan).

#### Acceptance criteria

- [ ] Jenis reasuransi dipilih dari **master**, tidak pernah diketik bebas. *(AC 10 spec)*
- [ ] ⚠️ Daftar disaring **persis** seperti existing: **dua belas ID di-blacklist** (`10004`,
      `10011`, `10012`, `10021`, `10022`, `10025`, `10026`, `10028`, `10248`, `10249`, `10018`,
      `10217`) **dan** `Flag = "active"` **dan** `Type` termasuk `1`, `2`, `3`.
      *(AC 11 spec; `[keputusan work owner]` — ditiru apa adanya, **jangan digeneralkan**)*
- [ ] Master jenis reasuransi **tidak ditulis** oleh konteks ini. Test yang menemukan tulisan ke
      master itu **gagal**. *(AC 12 spec)*
- [ ] ⚠️ Nama komponen, endpoint, dan fungsi di sistem baru **tidak mengandung kata "Old"**
      meskipun rule sumbernya bernama demikian. *(AC 13 spec; penyimpangan sadar 8)*
- [ ] Daftar dibaca dari **satu tempat**, dipakai layar kontrak maupun seluruh grid klausul —
      bukan disalin per layar.
- [ ] Master yang **kosong atau tidak terbaca** menghasilkan kegagalan yang **terlihat**, bukan
      daftar kosong yang diam. *(**ADR-U-0015**)*

#### Blocker

**Tidak ada.** **OQ-020 ditutup** — arti `ReinsTypeID` sudah jelas: master bersama life & non-life,
dipisahkan oleh saringan di atas.

#### Catatan

⚠️ **Beda ejaan antar konteks.** `.Flag` bernilai **`1`** di Master Contract Retro Life, tetapi
**`"active"`** di sini. Jangan menyalin nilai dari konteks Life.

⚠️ **Jangan ikut memigrasikan `GetMasterReinsTypeContract`** — kueri itu membaca `M_TREATYCONTRACT`
dan `M_TREATYYEAR` lewat `a.JSONDATA.…`, dua tabel JSON yang sudah **mati** (tiket 01).

⚠️ **Layar master jenis reasuransi ada di luar konteks ini.** `[keputusan work owner]`
`Harness/InboxTreatyContractReinsType.xml` (`DATA-PORTAL!INBOXTREATYCONTRACTREINSTYPE`) dan RD kedua
adalah perilaku **layar master**, konteks/menu tersendiri. Di sini: **baca saja**.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — saringan blacklist hanya terbukti benar
bila diuji terhadap master yang benar-benar memuat ID yang dikecualikan.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Treaty Contract Out - 03 - Tahun treaty — CRUD, gerbang periode, dan anti-dobel

**Status:** ready-for-agent

**Blocked by:** 01 (skema + sequence harus ada)

#### Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin membuat dan mengubah **tahun treaty** beserta grup
treaty, underwriting year, proporsi, dan masa berlakunya — **tanpa pernah mengetik nomor
identitas**, **tanpa** dapat membuat periode yang berakhir sebelum dimulai, dan **tanpa** dapat
membuat tahun yang sama dua kali. *(User story 1–3, 5–6 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas tahun treaty |
| `internal/repository` | Upsert dikunci `ID`; identitas dari sequence |
| `internal/services` | Gerbang periode; gerbang anti-dobel |
| `internal/handlers` | Endpoint tahun treaty |
| `frontend/` | Layar daftar + form tahun treaty |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveTreatyYear_Act` | `@BASECLASS` / `SAVETREATYYEAR_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyYear_Act.xml` | orkestrator simpan |
| `NewInputTreatyYear_Act` | `@BASECLASS` / `NEWINPUTTREATYYEAR_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/NewInputTreatyYear_Act.xml` | baris baru |
| `SetTreatyYear_Act` | `@BASECLASS` / `SETTREATYYEAR_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetTreatyYear_Act.xml` | isi form dari baris terpilih |
| `CheckYear` | `@BASECLASS` / `CHECKYEAR` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/CheckYear.xml` | ⚠️ **satu-satunya** validasi: `isNumber(TreatyYear)` |
| `SaveMasterTreatyYear_SQL` | `ASM-FW-GISFW-INT-TREATYYEAR` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyYear_SQL.xml` | → `POOLDATA.PEGA_TREATYYEAR` |
| `BrowseTreatyYear_RD` | `ASM-FW-GISFW-INT-TREATYYEAR` / `BROWSETREATYYEAR_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseTreatyYear_RD.xml` | daftar tahun |

`[data DBA]` `POOLDATA.PEGA_TREATYYEAR` — **upsert dikunci `ID`**; identitas baru dibuat di basis
data sebagai `'1' || lpad(TreatyYear_seq.nextval, 6, '0')`; `TGLUPDATE` diisi `SYSDATE`;
**tidak `COMMIT` sendiri**; keluaran `StsSimpan` **1 = sukses / 0 = gagal**.

#### ADR terkait

**ADR-U-0006** (identitas lewat sequence basis data), **ADR-U-0007** (jejak audit),
**ADR-U-0015** (kegagalan ditangani eksplisit, tidak ditelan).

#### Acceptance criteria

- [ ] Tahun treaty dapat dibuat dengan **grup treaty, underwriting year, proporsi, dan masa
      berlakunya**. *(AC 4 spec; User story 1)*
- [ ] Identitas tahun treaty **tidak pernah diketik pengguna** — dibuat dari sequence.
      *(AC 5 spec; **ADR-U-0006**)*
- [ ] Identitas berbentuk `'1'` diikuti nomor urut ber-*padding* nol **6 digit**.
      *(AC 6 spec; `[data DBA]`)*
- [ ] Menyimpan tahun treaty yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
      *(AC 8 spec — sisi tahun; `[data DBA]` upsert dikunci `ID`)*
- [ ] Masa berlaku yang **berakhir sebelum dimulai** **ditolak**, dengan pesan yang menyebut
      field-nya. *(AC 9 spec — sisi tahun)*
- [ ] **Tahun treaty yang sama tidak dapat dibuat dua kali**: kombinasi **(`STARTDATE`, `ENDDATE`,
      `TREATYGROUPID`)** yang **sudah ada** ditolak, dengan pesan yang menyebut tahun treaty mana
      yang sudah memakainya. *(AC 73 spec; `[keputusan work owner]` — **aturan baru**)*
- [ ] Gerbang anti-dobel **tidak** menghalangi pembaruan baris itu sendiri — memperbarui tahun
      treaty yang sudah ada dengan periode & grup yang tidak berubah tetap **berhasil**.
- [ ] Tanggal mulai dan tanggal akhir tersimpan sebagai **tanggal**, bukan teks. *(AC 53 spec)*
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**; nilai status selain `1`
      **selalu** dibaca sebagai kegagalan. *(AC 38, 39 spec; **ADR-U-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit** — siapa dan kapan. *(AC 41 spec; **ADR-U-0007**)*
- [ ] ⚠️ **Fitur salin tahun treaty tidak dibangun.** Tidak ada jalur — layar, endpoint, maupun
      pekerjaan latar — yang menyalin isi satu tahun treaty ke tahun lain. Test yang menemukan
      jalur semacam itu **gagal**. *(AC 72 spec; `[fakta bisnis — work owner]` — penyimpangan
      sadar 3)*

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Mengapa fitur salin dibuang.** `[fakta bisnis — work owner]` `POOLDATA.PROSESCOPY` menyalin
isi satu tahun treaty ke tahun lain, dan dalam praktiknya membawa **nilai tahun lalu** — misalnya
batas QS tahun 2025 — ke tahun yang semestinya berbeda. Itu sumber kesalahan, bukan penghemat
waktu. Yang **tidak** dimigrasikan: `RDBList/SaveMasterCopyData_SQL.xml`,
`Activity/BrowseCopyData.xml`, `Activity/SaveTreatyYearMultiple_Act.xml`,
`Activity/NewInputTreatyYearMultiple_Act.xml`, dan bagian salin
`Activity/BrowseDeleteRowTreatyInContract.xml`.

⚠️ **Validasi existing sangat tipis.** `[terverifikasi]` `Activity/CheckYear.xml` hanya memeriksa
`@Default.isNumber(InputTreatyYear.TreatyYear)`. Tidak ada gerbang periode, tidak ada anti-dobel.
Kedua gerbang di tiket ini adalah **tambahan sadar** `[keputusan work owner]`, bukan tiruan.

⚠️ `[terverifikasi]` `RDBList/SaveMasterCopyData_SQL.xml` menulis
`dbms_output.put_line(errmsg)` dengan variabel lokal `errmsg` yang **dideklarasikan tetapi tidak
pernah diisi** — selalu mencetak NULL. Dicatat sebagai jejak; tidak dimigrasikan.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — gerbang anti-dobel adalah pertanyaan
keunikan di basis data; memalsukannya berarti tidak mengujinya.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Treaty Contract Out - 04 - Kontrak treaty di dalam tahun treaty

**Status:** ready-for-agent

**Blocked by:** 02 (pemilih jenis reasuransi), 03 (tahun treaty sebagai induk)

#### Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin membuat **kontrak treaty** di dalam sebuah tahun
treaty — dengan **jenis reasuransinya** dan **masa berlakunya sendiri** — dan mengubahnya kemudian
tanpa membuat duplikat. *(User story 4–5, 7 di spec)*

Kontrak inilah yang membuka kombinasi **(tahun, grup, jenis reasuransi)** yang dipakai reinsurer,
business, dan seluruh klausul di tiket-tiket berikutnya.

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas kontrak treaty |
| `internal/repository` | Upsert dikunci `ID`; rujukan ke tahun treaty |
| `internal/services` | Gerbang masa berlaku |
| `internal/handlers` | Endpoint kontrak |
| `frontend/` | Grid kontrak di dalam layar tahun treaty |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveTreatyContract_Act` | `@BASECLASS` / `SAVETREATYCONTRACT_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyContract_Act.xml` | orkestrator simpan |
| `NewInputTreatyContract_Act` | `@BASECLASS` / `NEWINPUTTREATYCONTRACT_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/NewInputTreatyContract_Act.xml` | baris baru |
| `SetUbahTreatyContract` | `@BASECLASS` / `SETUBAHTREATYCONTRACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetUbahTreatyContract.xml` | muat untuk diubah |
| `SetTanggalTreatyContract` | `@BASECLASS` / `SETTANGGALTREATYCONTRACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetTanggalTreatyContract.xml` | isi tanggal |
| `SaveMasterTreatyContract_SQL` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyContract_SQL.xml` | → `POOLDATA.PEGA_TREATYCONTRACT` |
| `InboxTreatyContract` | `DATA-PORTAL` / `INBOXTREATYCONTRACT` / `RULE-HTML-HARNESS` | `Treaty Contract Out/Harness/InboxTreatyContract.xml` | titik masuk modul |

`[terverifikasi]` Parameter `PEGA_TREATYCONTRACT`: `ID`, `IDTreatyYear`, `ReinsTypeID`,
`ReinsTypeName`, `TreatyStartDate`, `TreatyEndDate`, pengguna, `TglUpdate`, + 2 keluaran.

`[data DBA]` Upsert dikunci `ID`; identitas `'1' || lpad(treatycontract_seq.nextval, 6, '0')`;
`TREATYSTARTDATE`/`TREATYENDDATE` di-`to_date(…,'DD/MM/YYYY')`; **tidak `COMMIT` sendiri**;
`StsSimpan` **1 = sukses / 0 = gagal**.

#### ADR terkait

**ADR-U-0006** (identitas lewat sequence), **ADR-U-0007** (jejak audit), **ADR-U-0015** (kegagalan
eksplisit).

#### Acceptance criteria

- [ ] Kontrak treaty dibuat **di dalam** sebuah tahun treaty dan menyimpan **rujukan ke tahun itu**.
      *(AC 7 spec; User story 4)*
- [ ] Menyimpan kontrak yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
      *(AC 8 spec — sisi kontrak)*
- [ ] Masa berlaku kontrak yang **berakhir sebelum dimulai** **ditolak**, dengan pesan yang
      menyebut field-nya. *(AC 9 spec — sisi kontrak)*
- [ ] Jenis reasuransi kontrak dipilih dari daftar tersaring tiket **02** — tidak diketik bebas.
      *(AC 10, 11 spec)*
- [ ] Identitas kontrak **tidak pernah diketik pengguna**; berbentuk `'1'` + 6 digit.
      *(AC 5, 6 spec; **ADR-U-0006**)*
- [ ] Tanggal mulai dan akhir kontrak tersimpan sebagai **tanggal**, bukan teks. *(AC 53 spec)*
- [ ] Kontrak **tidak dapat dibuat** tanpa tahun treaty induk yang ada.
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38, 39 spec;
      **ADR-U-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit**. *(AC 41 spec; **ADR-U-0007**)*

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Kombinasi, bukan foreign key.** `[fakta bisnis — work owner]` Reinsurer, business, dan klausul
**tidak** menyimpan `ID` kontrak. Mereka menggantung pada **(TreatyYear, TreatyGroupID,
ReinsTypeID)**. Kontrak "membuka" kombinasi itu, tetapi tidak memilikinya. Ini **berbeda** dari
Master Contract Retro Life, di mana anak-anak menggantung pada `ID` kontrak — jangan menyalin pola
dari sana.

⚠️ **Editor master, bukan proses berjenjang.** `[terverifikasi]` Nol rule `Flow`, nol rule `When`,
nol `StatusAkseptasi` di seluruh 303 berkas modul. Tidak ada Submit/Decline, tidak ada assignment,
tidak ada SLA.

`[terverifikasi]` **RBAC tidak dapat direkonstruksi** dari korpus — modul ini tidak memuat ekspresi
visibilitas ber-workbasket (`discovery/flows/_METHOD-noflow.md` §3.5). Siapa yang boleh mengedit
kontrak adalah keputusan terpisah, di luar tiket ini.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Treaty Contract Out - 05 - Reinsurer + total share

**Status:** ready-for-agent

**Blocked by:** 04 (kombinasi tahun/grup/jenis dibuka oleh kontrak)

#### Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin mencatat **para reinsurer** pada sebuah kombinasi tahun,
grup, dan jenis reasuransi, masing-masing dengan **share** dan **komisi reasuransi**-nya; dan
sebagai **underwriter**, saya ingin melihat **total share** seluruh reinsurer, supaya saya tahu
apakah penempatan sudah penuh. *(User story 11–12 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas reinsurer pada kombinasi |
| `internal/repository` | Upsert dikunci `ID`; baca daftar per kombinasi |
| `internal/services` | Penjumlahan total share |
| `internal/handlers` | Endpoint daftar + simpan reinsurer |
| `frontend/` | Grid reinsurer dengan baris total |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveTreatyReinsurerDetail1_Act` | `@BASECLASS` / `SAVETREATYREINSURERDETAIL1_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyReinsurerDetail1_Act.xml` | orkestrator simpan |
| `NewTreatyReinsurerDetail_Act` | `@BASECLASS` / `NEWTREATYREINSURERDETAIL_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/NewTreatyReinsurerDetail_Act.xml` | baris baru |
| `SetUbahTreatyReinsurerList_Act` | `@BASECLASS` / `SETUBAHTREATYREINSURERLIST_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetUbahTreatyReinsurerList_Act.xml` | muat untuk diubah |
| `SetErrorMessageReinsurer` | `@BASECLASS` / `SETERRORMESSAGEREINSURER` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetErrorMessageReinsurer.xml` | ⚠️ normalisasi koma→titik |
| `SaveMasterTreatyReinsurer_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyReinsurer_SQL.xml` | → `POOLDATA.PEGA_TREATYREINSURER` |
| `GetMasterReinsurerList` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterReinsurerList.xml` | baca daftar per kombinasi |
| `BrowseDetailTreatyReisurer_RD` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseDetailTreatyReisurer_RD.xml` | daftar |

`[terverifikasi]` `GetMasterReinsurerList` menyaring **`TreatyYear` + `TreatyGroupID` +
`ReinsTypeID`** — kombinasi, bukan `ID` kontrak.
`[terverifikasi]` Parameter `PEGA_TREATYREINSURER` (19 in + 2 out): `ID`, `TreatyYear`,
`TreatyGroupID`, `TreatyGroupName`, `ReinsTypeID`, `ReinsTypeName`, `ReinsurerID`, `CLIENTID`,
`NAME`, `Ricomm`, `PctShare`, `IUDate`, `UserId`, `StartDate`, `EndDate`, `StatusOn`, `StdRating`,
`OperatorName`, `TglUpdate`.

`[data DBA]` `RICOMM` dan `PCTSHARE` sudah **`NUMBER`** di existing; sisanya `VARCHAR2`.
Upsert dikunci `ID`; identitas `'1' || lpad(M_TREATYREINSURER_SEQ.nextval, 6, '0')`; **tidak
`COMMIT` sendiri**; `StsSimpan` **1 = sukses / 0 = gagal**.

#### ADR terkait

**ADR-U-0003** (uang & persen non-float), **ADR-U-0006** (identitas lewat sequence),
**ADR-U-0007** (jejak audit), **ADR-U-0015** (kegagalan eksplisit).

#### Acceptance criteria

- [ ] Reinsurer dicatat pada kombinasi **(tahun, grup, jenis reasuransi)**, masing-masing dengan
      **share** dan **komisi reasuransi**. *(AC 14 spec; User story 11)*
- [ ] **Total share** seluruh reinsurer pada satu kombinasi **terlihat** bagi pengguna.
      *(AC 15 spec; User story 12)*
- [ ] ⚠️ Share dan komisi diperlakukan sebagai **desimal presisi arbitrer**; **tidak** melewati
      `float` dan **tidak** dibulatkan ke bilangan bulat. *(AC 16 spec; **ADR-U-0003**; penyimpangan
      sadar 6)*
- [ ] Identitas reinsurer **tidak pernah diketik pengguna**; berbentuk `'1'` + 6 digit.
      *(AC 5, 6 spec; **ADR-U-0006**)*
- [ ] Menyimpan reinsurer yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
      *(AC 8 spec)*
- [ ] Daftar reinsurer disaring **per kombinasi**, bukan per `ID` kontrak.
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38, 39 spec;
      **ADR-U-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit**. *(AC 41 spec; **ADR-U-0007**)*

#### Blocker

**Tidak ada.**

⚠️ `[terbuka]` **Total share = 100% bukan gerbang di tiket ini.** Existing **tidak** menolak
kombinasi yang totalnya ≠ 100 pada tingkat reinsurer; yang ada hanyalah penegakan 100% pada baris
**anak klausul** (tiket 08). Apakah total share reinsurer wajib 100% adalah **fakta bisnis yang
belum ditetapkan** — AC di atas hanya mewajibkan totalnya **terlihat**, bukan ditegakkan. Bila
Product + UW menetapkan sebaliknya, itu AC tambahan, bukan perubahan tiket ini.

#### Catatan

⚠️ **Angka hari ini disimpan sebagai teks ber-koma desimal.** `[terverifikasi]`
`Activity/SetErrorMessageReinsurer.xml` (`@BASECLASS!SETERRORMESSAGEREINSURER`) menjalankan
`@replaceAll(InputTreatyReinsurer.PctShare, ",", ".")` dan hal yang sama untuk `Ricomm` — bukti
bahwa nilai masuk dari layar dalam format berkoma. Di sistem baru **normalisasi dilakukan di batas
masukan**, sekali, bukan ditempel di tengah alur simpan.

⚠️ `[terverifikasi]` Kolom `STDRATING` adalah **field yang dipakai-ulang** — pola yang sama sudah
ditemukan di Master Contract Retro Life. Isi sebenarnya belum terverifikasi; bawa apa adanya dan
jangan menafsirkan.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — presisi desimal share dan komisi hanya
terbukti benar pada basis data sungguhan.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Treaty Contract Out - 06 - Security reinsurer — struktur bersih

**Status:** ready-for-agent

**Blocked by:** 05 (security menggantung pada reinsurer)

#### Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin mencatat **security** di bawah seorang reinsurer beserta
porsinya, supaya eksposur berjenjang terlihat; dan sebagai **organisasi**, saya ingin baris security
punya **identitas sendiri**, supaya mengubah nama security tidak memutus rujukannya.
*(User story 13–14 di spec)*

⚠️ Inilah tiket dengan pembersihan struktur paling dalam di konteks ini.

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas security reinsurer |
| `internal/repository` | Tulis/baca dengan **kolom bernama** dan **PK surrogate** |
| `internal/services` | Aturan per baris |
| `internal/handlers` | Endpoint security |
| `frontend/` | Grid security di bawah baris reinsurer |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveSecurityReinsurer_Act` | `@BASECLASS` / `SAVESECURITYREINSURER_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveSecurityReinsurer_Act.xml` | orkestrator simpan |
| `ShowEditSecurityReinsurer` | `@BASECLASS` / `SHOWEDITSECURITYREINSURER` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/ShowEditSecurityReinsurer.xml` | muat untuk diubah |
| `InsertToMTreatySecurity` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `ASM!INSERTTOMTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/InsertToMTreatySecurity.xml` | ⚠️ INSERT **posisional** |
| `UpdateMTreatySecurity` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/UpdateMTreatySecurity.xml` | ⚠️ kunci `trim()` |
| `DeleteSecurityReinsurer` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteSecurityReinsurer.xml` | ⚠️ kunci `trim()` |
| `SelectSecurityReinsurer` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/SelectSecurityReinsurer.xml` | daftar |

`[terverifikasi]` Keadaan existing — **satu-satunya master di modul ini yang ditulis SQL mentah**:

```
insert into mtreatysecurity values ({…}, '', '', {…}, {…PCT_SHARE}, '', {…REAS_SECURITY})
update mtreatysecurity set THN_TREATY=…, PCT_SHARE=…, REAS_SECURITY=…
 where REAS_ID = … and trim(REAS_SECURITY) = trim(…)
delete from mtreatysecurity where REAS_ID = … and trim(REAS_SECURITY) = trim(…)
```

⚠️ INSERT **tanpa daftar kolom** — tujuh nilai berposisi, **tiga di antaranya string kosong**.
⚠️ Ketiganya **tanpa `COMMIT`**.
⚠️ Kunci pembaruan dan penghapusan memakai **`trim(REAS_SECURITY)`** — nama dipakai sebagai bagian
kunci, dan `trim()` mengakui datanya bertabur spasi.

`[data DBA]` DDL existing: `THN_TREATY VARCHAR2(4) DEFAULT '1' NOT NULL`, `TOP_ID VARCHAR2(9)`,
`TP_TREATY CHAR(2)`, `REAS_ID CHAR(7) NOT NULL`, `PCT_SHARE VARCHAR2(99)`, `USER_ID CHAR(99)`,
`REAS_SECURITY CHAR(10) NOT NULL` — **tanpa primary key**. INSERT posisional mengosongkan
`TOP_ID`, `TP_TREATY`, `USER_ID`.

⚠️ **Penyimpangan sadar 5 — struktur dibersihkan.** `[keputusan work owner]` Tabel anak
`TREATYREINSURER` yang wajar: **PK surrogate**, seluruh kolom **bernama dan diisi eksplisit**,
`PCT_SHARE` **desimal**, `REAS_SECURITY` **atribut biasa** — bukan bagian kunci.

#### ADR terkait

**ADR-U-0003** (persen non-float), **ADR-U-0007** (jejak audit), **ADR-U-0015** (kegagalan eksplisit).

#### Acceptance criteria

- [ ] Security dicatat **di bawah seorang reinsurer** beserta porsinya. *(AC 17 spec; User story 13)*
- [ ] ⚠️ Baris security punya **primary key surrogate** sendiri; `REAS_SECURITY` adalah **atribut
      biasa**, **bukan** bagian kunci. Test yang menemukan kunci berbasis nama security **gagal**.
      *(AC 18 spec; User story 14; penyimpangan sadar 5)*
- [ ] ⚠️ Seluruh kolom baris security **diisi eksplisit dan bernama**; **tidak ada** penulisan
      berposisi. Test yang menemukan penulisan tanpa daftar kolom **gagal**. *(AC 19 spec;
      penyimpangan sadar 5)*
- [ ] ⚠️ Pencocokan baris security **tidak** memakai `trim()` atas nilai kunci. *(AC 20 spec;
      penyimpangan sadar 5)*
- [ ] ⚠️ `PCT_SHARE` bertipe **desimal**, bukan teks. *(AC 51 spec; **ADR-U-0003**; penyimpangan
      sadar 6)*
- [ ] Mengubah **nama security** **tidak** memutus rujukan barisnya ke reinsurer. *(User story 14)*
- [ ] Baris security dapat **ditambah dan dihapus** tanpa menyentuh baris reinsurer induknya —
      kecuali ketika penyimpanan dilakukan sebagai satu kesatuan (tiket 09).
- [ ] Menghapus seorang reinsurer **menghapus juga** baris security di bawahnya. *(lihat tiket 10
      untuk kaskade dari kontrak)*
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38, 39 spec;
      **ADR-U-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit**. *(AC 41 spec; **ADR-U-0007**)*

#### Blocker

**Tidak ada.**

⚠️ `[terbuka]` **Tiga kolom yang dikosongkan INSERT posisional** — `TOP_ID`, `TP_TREATY`,
`USER_ID`. `[data DBA]` DDL-nya sudah diketahui, tetapi **apakah ketiganya masih punya arti** tidak
terbaca dari korpus. Tiket ini **membawanya sebagai kolom bernama**; bila migrasi (tiket 01)
membuktikan seluruh barisnya kosong, pembuangannya adalah keputusan terpisah — **jangan dibuang di
sini atas inisiatif sendiri**.

#### Catatan

⚠️ **Mengapa `trim()` tidak dibawa.** Pemakaian `trim()` di kunci adalah **pengakuan bahwa datanya
kotor**, bukan aturan bisnis. Membawanya ke sistem baru berarti membawa kekotorannya, dan menutup
kemungkinan dua security yang namanya hanya berbeda spasi tetap dianggap berbeda. Migrasi (tiket 01)
memberi setiap baris identitas sendiri; sesudah itu nama tidak perlu lagi memikul beban kunci.

⚠️ `RDBList/DeleteFromTreatyReinsurer_Act.xml` menghapus `MTREATYSECURITY` lalu `TREATYREINSURER`
**tanpa `COMMIT`** — konsisten dengan temuan bahwa modul ini menyerahkan batas transaksi kepada
pemanggil (tiket 09).

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — keutuhan rujukan setelah pemberian PK dan
hilangnya kunci berbasis nama **hanya berperilaku benar pada basis data sungguhan**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Treaty Contract Out - 07 - Business + penonaktifan

**Status:** ready-for-agent

**Blocked by:** 04 (kombinasi tahun/grup/jenis dibuka oleh kontrak)

#### Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin mencatat **jenis bisnis** yang ditanggung sebuah
kontrak dengan kode dan namanya, dan **menonaktifkan** satu baris tanpa menghapusnya — supaya
riwayatnya tetap ada. *(User story 15–16 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas baris bisnis |
| `internal/repository` | Upsert dikunci `ID`; baca daftar per kombinasi |
| `internal/services` | Penonaktifan; pembaruan seluruh field |
| `internal/handlers` | Endpoint daftar + simpan bisnis |
| `frontend/` | Grid bisnis di layar kontrak |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveTreatyBusinessDetail_Act` | `@BASECLASS` / `SAVETREATYBUSINESSDETAIL_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyBusinessDetail_Act.xml` | orkestrator simpan |
| `NewTreatyBusinessDetail_Act` | `@BASECLASS` / `NEWTREATYBUSINESSDETAIL_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/NewTreatyBusinessDetail_Act.xml` | baris baru |
| `SetUbahTreatyBusinessList_Act` | `@BASECLASS` / `SETUBAHTREATYBUSINESSLIST_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetUbahTreatyBusinessList_Act.xml` | muat untuk diubah |
| `SaveMasterTreatyBusiness_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyBusiness_SQL.xml` | → `POOLDATA.PEGA_TREATYBUSINESS` |
| `GetMasterBusinessList` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterBusinessList.xml` | baca daftar per kombinasi |
| `DeleteRowBusinessList` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `ASM!DELETEROWBUSINESSLIST` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteRowBusinessList.xml` | ⚠️ memegang **kedua** tabel kembar |
| `BrowseTreatyBusiness_RD`, `BrowseFilterBusiness_RD` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/` | daftar & saringan |

`[terverifikasi]` Parameter `PEGA_TREATYBUSINESS` (12 in + 2 out): `ID`, `IsActive`, `TreatyYear`,
`TreatyYearID`, `TreatyGroupID`, `TreatyGroupName`, `ReinsTypeID`, `ReinsTypeName`, `BizCode`,
`BIZNAME`, `UserID`, `TglUpdate`.

`[data DBA]` Upsert dikunci `ID`; identitas `'1' || lpad(TREATY_BUSINESS_SEQ.nextval, 6, '0')`;
**tidak `COMMIT` sendiri**; `StsSimpan` **1 = sukses / 0 = gagal**.
⚠️ **Saat UPDATE, procedure existing hanya mengisi `ISACTIVE`, `BIZCODE`, `BIZNAME`, `USERID`,
`TGLUPDATE`** — `REINSTYPEID` dan kawan-kawan hanya diisi saat INSERT.

#### ADR terkait

**ADR-U-0006** (identitas lewat sequence), **ADR-U-0007** (jejak audit), **ADR-U-0015** (kegagalan
eksplisit).

#### Acceptance criteria

- [ ] Jenis bisnis dicatat dengan **kode** dan **nama**-nya pada sebuah kontrak. *(AC 21 spec;
      User story 15)*
- [ ] Satu baris bisnis dapat **dinonaktifkan** tanpa dihapus, dan tetap terbaca sebagai baris
      nonaktif. *(AC 22 spec; User story 16)*
- [ ] ⚠️ Memperbarui baris bisnis **memperbarui seluruh field yang dikirim** — **bukan** hanya lima
      kolom seperti procedure existing. Test yang menemukan field terkirim yang tidak tersimpan
      **gagal**. *(AC 23 spec; `[keputusan work owner]` — perbaikan sadar atas perilaku existing)*
- [ ] Identitas baris bisnis **tidak pernah diketik pengguna**; berbentuk `'1'` + 6 digit.
      *(AC 5, 6 spec; **ADR-U-0006**)*
- [ ] Menyimpan baris yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
      *(AC 8 spec)*
- [ ] Daftar bisnis disaring **per kombinasi** (tahun, grup, jenis reasuransi), bukan per `ID`
      kontrak.
- [ ] ⚠️ Penghapusan baris bisnis menyentuh **satu tabel saja** — tidak ada tabel kembar JSON yang
      ikut dihapus. *(AC 63, 64 spec; penyimpangan sadar 1)*
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38, 39 spec;
      **ADR-U-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit**. *(AC 41 spec; **ADR-U-0007**)*

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Rule ini adalah bukti paling telanjang dari dualitas tabel kembar.** `[terverifikasi]`
`RDBList/DeleteRowBusinessList.xml` memegang dua SQL sekaligus:

```
<pyBrowseSQL>  delete from treatybusiness   where id = {…}
<pyDeleteSQL>  delete from m_treatybusiness where id = {…}; commit;
```

Satu rule, dua tabel, dua ejaan. Setelah penyimpangan sadar 1 diterapkan (tiket 01), hanya yang
pertama tersisa.

⚠️ **`TREATYYEARID` boleh NULL di data lama.** `[terverifikasi]` Kaskade hapus existing
(`RDBList/DeleteFromTREATYCONTRACT_SQL.xml`) menulis
`(TREATYYEARID = {…} OR TREATYYEARID IS NULL)` — pengakuan bahwa sebagian baris bisnis lama tidak
punya nilai itu. Migrasi (tiket 01) dan pembacaan di sini harus **tahan terhadap NULL**, bukan
mengandaikannya selalu terisi.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Treaty Contract Out - 08 - Klausul — satu tabel, 25 jenis, validasi per jenis

**Status:** ready-for-agent

**Blocked by:** 04 (kombinasi tahun/grup/jenis dibuka oleh kontrak)

⚠️ **Inti konteks ini.** Tiket paling berisi: 13 acceptance criteria, dan satu-satunya tempat di
mana penyimpangan sadar 2 diterapkan.

#### Hasil & nilai pengguna

Sebagai **underwriter**, saya ingin mengelola **dua puluh lima jenis klausul** kontrak treaty di
satu tempat — sebagian berlapis dengan baris rincian di bawahnya — dan saya ingin **setiap jenis
memeriksa field yang memang relevan baginya**, bukan satu daftar wajib yang seragam.
*(User story 17–22 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Satu bentuk baris klausul; jenis sebagai atribut |
| `internal/repository` | Tulis/baca **satu tabel** `proportionalarrg` |
| `internal/services` | Pemetaan jenis klausul → aturan wajib-isi; anti-dobel; hitungan turunan |
| `internal/handlers` | Endpoint daftar + simpan per jenis klausul |
| `frontend/` | Grid per jenis klausul; grid rincian untuk jenis berlapis |

#### Rule Pega sumber

`[terverifikasi]` Sensus penuh atas **25** activity `SaveTreatyArr*` di
`Treaty Contract Out/Activity/` — seluruhnya berujung pada **dua** rule Connect-SQL:

| Menulis lewat | Jml | Jenis klausul |
| --- | ---: | --- |
| `SaveMasterProportionalArrg` → `POOLDATA.PEGA_PROPORTIONALARRG` (**35 kolom**) | **18** | `BordereAux`, `CashLossLimit`, `ClaimCoorp`, `CoinsPanel`, `EPI`, `ExGratia`, `ExclutionTreaty`, `FacIn`, `LimitMB`, `MaxCoinsPanel`, `MinLOL`, `MinLOLMB`, `PLA`, `Portfolio`, `ProfitComm`, `Ricomm`, `TerrLimit`, `TreatyLimit` |
| `SaveMasterProportionalArrgChild` → `POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD` (**26 kolom**) | **7** | `CashLossLimitList`, `ClaimCoorpChild`, `EpiList`, `ExGratiaChildList`, `FacInList`, `PLAList`, `TreatyLimitChild` |

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `SaveMasterProportionalArrg` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterProportionalArrg.xml` |
| `SaveMasterProportionalArrgChild` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterProportionalArrgChild.xml` |
| `SaveTreatyArrEPI_Act` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `SAVETREATYARREPI_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyArrEPI_Act.xml` (111.418 byte, 10 langkah) |
| `HitungRpUsd` | `@BASECLASS` / `HITUNGRPUSD` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/HitungRpUsd.xml` |
| `TreatyTestChildTotal_Act` | `@BASECLASS` / `TREATYTESTCHILDTOTAL_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/TreatyTestChildTotal_Act.xml` |
| `BrowseTreatyDesc_RD` | `ASM-FW-GISFW-INT-TREATYDESC` / `BROWSETREATYDESC_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseTreatyDesc_RD.xml` |
| 16 `CancelActivity*` | `Treaty Contract Out/Activity/` | pembatalan per jenis klausul |

`[terverifikasi]` **Mekanisme existing**: tiap jenis punya halaman masukannya sendiri
(`InputTreatyArrEpi`, `InputTreatyArrPLA`, `InputTreatyArrProfit`, `InputTreatyArrBord`, …), lalu
langkah **`Page-Copy`** menyalinnya ke halaman bersama `InputTreatyArrTreatyLimit` (induk) atau
`InputTreatyArrTreatyLimitChild` (anak) sebelum langkah `RDB-List` memanggil procedure.

⚠️ **Penyimpangan sadar 2 — SATU tabel untuk semua jenis klausul.** `[keputusan work owner]`
`[data DBA]` **Kedua procedure menulis ke tabel yang sama: `POOLDATA.PROPORTIONALARRG`.** Bedanya
hanya jumlah kolom. Baris "anak" mengisi **NULL** pada sembilan kolom khusus induk:
`ID_OCCUPATION`, `OCCUPATION`, `ID_CLAUSE`, `CLAUSE`, `TREATYLIMIT`, `COINS_MIN`, `COINS_MAX`,
`MORERP`, `MOREUSD`. Jenis dibedakan **`TREATYDESCID`**.

`[data DBA]` Identitas `'1' || lpad(PROPORTIONALARRG_SEQ.nextval, 7, '0')` — **7 digit**, berbeda
dari lima tabel lain yang 6 digit. **Tidak `COMMIT` sendiri**; `StsSimpan` **1 = sukses / 0 = gagal**.

#### Validasi wajib-isi per jenis `[terverifikasi]`

Sensus prasyarat langkah (`<pyStepsPreCondParamsWhen> … ==""`) atas 25 activity — **ini tabel
kebenaran untuk AC validasi**:

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

#### ADR terkait

**ADR-U-0003** (uang & persen non-float), **ADR-U-0006** (identitas lewat sequence),
**ADR-U-0007** (jejak audit), **ADR-U-0015** (kegagalan eksplisit).

#### Acceptance criteria

- [ ] ⚠️ Seluruh **dua puluh lima** jenis klausul tersimpan dalam **satu tabel**
      `proportionalarrg`, dibedakan oleh **`TREATYDESCID`**. Test yang menemukan tabel terpisah per
      jenis klausul **gagal**. *(AC 24 spec; penyimpangan sadar 2)*
- [ ] ⚠️ Baris jenis "anak" menyimpan **sembilan kolom khusus induk sebagai NULL**
      (`ID_OCCUPATION`, `OCCUPATION`, `ID_CLAUSE`, `CLAUSE`, `TREATYLIMIT`, `COINS_MIN`,
      `COINS_MAX`, `MORERP`, `MOREUSD`). *(AC 25 spec; penyimpangan sadar 2)*
- [ ] Daftar jenis klausul dibaca dari master **`TREATYDESC`**; jenis baru dapat ditambahkan di
      master **tanpa mengubah skema**. *(AC 26 spec; User story 18)*
- [ ] Master jenis klausul **tidak ditulis** oleh konteks ini. Test yang menemukan tulisan ke
      `TREATYDESC` **gagal**. *(AC 27 spec)*
- [ ] Klausul berlapis dapat dinyatakan lewat baris "anak" pada kombinasi yang sama, dibedakan
      **`PARENTREINSTYPEID`**. *(AC 28 spec; User story 19)*
- [ ] Membatalkan pengeditan **satu** jenis klausul **tidak** membuang jenis klausul lain yang
      sedang dikerjakan. *(AC 29 spec; User story 21)*
- [ ] Baris klausul yang **sudah pernah diinput** **ditolak** dengan pesan yang jelas.
      *(AC 30 spec; User story 22)*
- [ ] Urutan baris dalam daftar klausul **dipertahankan** saat dibaca kembali. *(AC 31 spec)*
- [ ] ⚠️ Nama jenis klausul dan nama field **dipakai apa adanya** dari Pega — **tidak
      diterjemahkan, tidak ditebak**. *(AC 32 spec; `[keputusan work owner]`)*
- [ ] Setiap jenis klausul memeriksa **field yang relevan baginya sendiri**, sesuai tabel di atas.
      *(AC 33 spec; User story 20)*
- [ ] Aturan wajib-isi berada di **kode**, bukan di data master. *(AC 34 spec)*
- [ ] Pesan penolakan **menyebut field** yang kurang. *(AC 35 spec)*
- [ ] `[terbuka]` Jenis **`LimitMB`** dan **`Portfolio`** **tidak** dinyatakan selesai sebelum
      aturan wajib-isinya ditetapkan Product + UW. Keduanya **tidak** dilepas sebagai "tanpa
      validasi". *(AC 36 spec)*
- [ ] ⚠️ Nilai uang (`Rp`, `Usd`, `MoreRp`, `MoreUsd`, `TreatyLimit`, `CoIns_Min`, `CoIns_Max`) dan
      persen (`Pct`, `PctMe`) bertipe **desimal**, bukan teks. *(AC 51, 52 spec; **ADR-U-0003**;
      penyimpangan sadar 6)*
- [ ] Identitas baris klausul **tidak pernah diketik pengguna**; berbentuk `'1'` + **7 digit** —
      berbeda dari lima tabel lain yang 6 digit. *(AC 6 spec; `[data DBA]`; **ADR-U-0006**)*

#### Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **Dua OQ yang tidak memblokir tiket ini, tetapi memblokir dua jenis klausul:**
field wajib `LimitMB` dan `Portfolio` (AC di atas), dan **arti bisnis** setiap istilah klausul
(`EPI`, `PLA`, `Ricomm`, `MB`, `LOL`, `BordereAux`, `ExGratia`, `CashLossLimit`,
`ClaimCoorperation`, `ProfitCommision`, `pTerrLimit`, `CoinsPanel`, `ExclutionTreaty`) dan field
(`Ydcf`, `PctMe`, `LayerPartType`, `LayerType`, `SpreadingOrder`, `Method`, `Line`). Keduanya
**Product + UW**. Layar dan skema dapat dibangun penuh dengan istilah asli.

#### Catatan

⚠️ **Dua perilaku existing yang mudah terlewat — keduanya nyata dan harus ikut.**

1. **`Rp` dan `Usd` pada baris anak adalah nilai TURUNAN.** `[terverifikasi]`
   `Activity/HitungRpUsd.xml` (`@BASECLASS!HITUNGRPUSD`) menghitung, untuk ketujuh jenis anak,
   `Usd = Pct × Usd_induk / 100` dan `Rp = Pct × Rp_induk / 100`. Jadi pengguna mengetik **`Pct`**,
   dan kedua nilai uang mengikuti. Perilaku ini **bukan** kosmetik layar — angka turunannya yang
   tersimpan.
2. **Total `Pct` baris anak wajib 100%.** `[terverifikasi]`
   `Activity/TreatyTestChildTotal_Act.xml` (`@BASECLASS!TREATYTESTCHILDTOTAL_ACT`) menjumlahkan
   `Pct` seluruh baris anak dan menolak bila `!= 100.00`, dengan pesan
   `"Please make sure spreading is 100%"`. Ini gerbang yang **memang ada** di existing — berbeda
   dari total share reinsurer (tiket 05) yang **tidak** ditegakkan.

⚠️ **Nama berbohong — dua kali dalam satu tiket.** `[data DBA]` Procedure
`PEGA_M_PROPORTIONALARRG_CHILD` **tidak** menulis ke tabel child; ia menulis ke tabel yang sama.
Dan `[terverifikasi]` `Activity/SetTreatyArrangementDesc_Act.xml`
(`ASM-FW-GISFW-INT-PROPORTIONALARRG`) — yang namanya paling wajar untuk "mengisi deskripsi
arrangement" — **lima dari enam langkahnya di-remark**; ia mati (**OQ-066**).

⚠️ **`ParentReinsTypeID = "00"`** adalah penanda baris tanpa induk `[terverifikasi]`
(`Activity/GetPeriode.xml`, `@BASECLASS!GETPERIODE`). Sentinel ini dibawa apa adanya.

⚠️ **Jangan memigrasikan kueri yang membaca JSON.** `GetMasterDescriptionEPIParentList`,
`GetMasterDescriptionPLAParentList`, `GetMasterDescriptionExGratiaList`,
`GetMasterDescriptionFACINParentList`, `GetMasterPortfolioListDetail`, `GetMasterPanggilID` —
seluruhnya membaca `m_PROPORTIONALARRG`. Di sistem baru jawabannya datang dari
**`proportionalarrg` relasional** (tiket 01).

⚠️ **`Portfolio` di existing berupa daftar bersarang di dalam JSON** `[terverifikasi]`
(`RDBList/GetMasterPortfolioListDetail.xml` memakai `JSON_TABLE` atas `'$.ProportionalList[*]'`
dengan kolom `PortfolioType`, `PremiLost`, `PortfolioValue`, `PortfolioDesc`). Karena JSON dibuang,
**bentuk penyimpanan Portfolio di skema baru ditetapkan bersama field wajibnya** — bagian dari OQ
`LimitMB`/`Portfolio` di atas. **Jangan tebak.**

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — satu tabel yang memuat 25 jenis dengan
pola NULL yang berbeda per jenis hanya terbukti benar pada basis data sungguhan.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Treaty Contract Out - 09 - Simpan atomik lintas enam tabel

**Status:** ready-for-agent

**Blocked by:** 05, 06, 07, 08 (seluruh penulis harus ada sebelum dapat dibungkus jadi satu)

#### Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin seluruh perubahan pada satu kontrak — kontrak, reinsurer,
security, business, **dan seluruh klausul** — tersimpan **bersama atau tidak sama sekali**, sehingga
tidak pernah ada kontrak yang tersimpan separuh; dan sebagai **admin master**, saya ingin
**diberi tahu bila penyimpanan gagal**, supaya saya tidak mengira data tersimpan padahal tidak.
*(User story 24–25 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | **Satu transaksi** melintasi enam tabel; commit sekali |
| `internal/services` | Orkestrasi simpan menyeluruh; rollback total |
| `internal/handlers` | Endpoint simpan kontrak utuh |
| `frontend/` | Tombol simpan tunggal; tampilan galat |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatyContract_SQL` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyContract_SQL.xml` | ⚠️ `COMMIT` di dalam SQL Pega |
| `SaveMasterTreatyYear_SQL` | `ASM-FW-GISFW-INT-TREATYYEAR` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyYear_SQL.xml` | ⚠️ idem |
| `SaveMasterTreatyReinsurer_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyReinsurer_SQL.xml` | ⚠️ idem |
| `SaveMasterTreatyBusiness_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyBusiness_SQL.xml` | ⚠️ idem |
| `SaveMasterProportionalArrg` / `…Child` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/` | ⚠️ idem |
| `RefreshErrorProportionalarrg` | `@BASECLASS` / `REFRESHERRORPROPORTIONALARRG` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/RefreshErrorProportionalarrg.xml` | tampilan galat |
| `SetErrorMessage` | `@BASECLASS` / `SETERRORMESSAGE` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetErrorMessage.xml` | tampilan galat |

⚠️ **Temuan yang membuat tiket ini mungkin** `[data DBA]`: keenam procedure penulis
(`PEGA_PROPORTIONALARRG`, `PEGA_M_PROPORTIONALARRG_CHILD`, `PEGA_TREATYCONTRACT`, `PEGA_TREATYYEAR`,
`PEGA_TREATYREINSURER`, `PEGA_TREATYBUSINESS`) **TIDAK `COMMIT` sendiri** — hanya `ROLLBACK` bila
galat; commit diserahkan kepada pemanggil.

`COMMIT` yang terlihat di rule Connect-SQL Pega berada **di luar** procedure, ditulis Pega sendiri.
Jadi **tidak ada titik potong transaksi di sisi basis data** — Go dapat dan harus membungkus
seluruhnya.

⚠️ Ini **berbeda dari lima konteks Life sebelumnya** (Claim Life, Komite, PremiumList, Endorsement,
Master Contract Retro Life), yang procedure-nya `COMMIT` sendiri dan karenanya memaksa titik potong.

#### ADR terkait

**ADR-U-0006** (identitas lewat sequence), **ADR-U-0007** (jejak audit setiap transaksi),
**ADR-U-0015** (kegagalan ditangani eksplisit, tidak ditelan).

#### Acceptance criteria

- [ ] ⚠️ Seluruh perubahan satu kontrak — kontrak, reinsurer, security, business, **dan seluruh
      klausul** — ditulis dalam **satu transaksi**; kegagalan di mana pun **membatalkan seluruhnya**.
      Dibuktikan dengan menyuntikkan kegagalan pada **satu baris klausul ke-N** lalu memastikan
      **tidak ada** perubahan tersimpan pada kelima tabel lain. *(AC 37 spec; User story 24;
      `[data DBA]` — procedure tidak commit sendiri)*
- [ ] Penyimpanan yang **gagal** menghasilkan kegagalan **terang-terangan** dengan pesan yang
      menyebut **apa** yang gagal — bukan diam-diam dianggap sukses. *(AC 38 spec; User story 6,
      25; **ADR-U-0015**)*
- [ ] Status simpan **1 = sukses / 0 = gagal** ditegakkan; nilai **selain `1`** — termasuk kosong
      dan NULL — **selalu** dibaca sebagai kegagalan. *(AC 39 spec; `[data DBA]`)*
- [ ] ⚠️ Pesan galat di sistem baru **tidak memuat teks `"JSON_KLAIM"`**. Test yang menemukannya
      **gagal**. *(AC 40 spec; penyimpangan sadar 8)*
- [ ] Setiap penyimpanan mencatat **jejak audit** — siapa dan kapan. *(AC 41 spec; **ADR-U-0007**)*
- [ ] Transaksi **di-commit sekali**, di akhir; **tidak ada** commit per tabel.
- [ ] Kegagalan mengembalikan basis data ke keadaan **persis sebelum** permintaan — termasuk
      **tidak menyisakan identitas terpakai** yang membuat nomor melompat tanpa alasan.

#### Blocker

**Tidak ada.** **OQ-002 ditutup** — body keenam procedure sudah diterima; perilaku transaksinya
terbukti.

#### Catatan

⚠️ **Inilah kebalikan dari Master Product Name Life.** Di sana `[data DBA]` membuktikan procedure
**`COMMIT` sendiri**, sehingga atomik hanya mungkin dengan **membuang** procedure-nya
(penyimpangan sadar di konteks itu). Di sini `[data DBA]` membuktikan **sebaliknya** — dan itulah
yang membuat simpan atomik lintas enam tabel menjadi mungkin **tanpa membuang satu pun procedure**.
Motif proyek: **"periksa apakah procedure commit sendiri."**

⚠️ **Nomor urut yang melompat bukan sekadar kosmetik.** Karena identitas dibuat basis data lewat
sequence (`ADR-U-0006`), rollback tidak mengembalikan `nextval` yang sudah diambil. AC terakhir
menuntut agar kegagalan **tidak mengambil identitas lebih dulu daripada perlu** — bukan agar
sequence dimundurkan.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — atomisitas lintas enam tabel adalah
**satu-satunya** hal yang diuji tiket ini, dan ia **hanya berperilaku benar pada basis data
sungguhan**. Memalsukan Oracle di sini berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Treaty Contract Out - 10 - Kaskade hapus + popup konfirmasi — klausul **tetap hidup**

**Status:** ready-for-agent

**Blocked by:** 06 (security), 07 (business), 09 (pembungkus transaksi)

#### Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin menghapus sebuah kontrak **beserta** reinsurer,
security, dan business-nya supaya tidak ada sisa yang menggantung — tetapi saya ingin **diberi
peringatan berisi jumlah baris yang akan ikut terhapus** lebih dulu, supaya saya dapat membatalkan.
*(User story 26–27 di spec)*

Dan sebagai **underwriter**, saya ingin **klausul bertahan** ketika kontrak dihapus, karena klausul
itu milik tahun/grup/jenis reasuransi, bukan milik satu kontrak. *(User story 23)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Kaskade hapus tiga anak + kontrak, **dalam satu transaksi** |
| `internal/services` | Hitung jumlah baris terdampak sebelum menghapus |
| `internal/handlers` | Endpoint pratinjau dampak + endpoint hapus |
| `frontend/` | Popup konfirmasi Ya/Batal dengan rincian jumlah |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `DeleteFromTREATYCONTRACT_SQL` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteFromTREATYCONTRACT_SQL.xml` | kaskade empat tabel |
| `DeleteFromTreatyReinsurer_Act` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteFromTreatyReinsurer_Act.xml` | hapus security + reinsurer |
| `DeleteTreatyReins_Act` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `DELETETREATYREINS_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/DeleteTreatyReins_Act.xml` | orkestrator |
| `DeleteRowBusiness` | `@BASECLASS` / `DELETEROWBUSINESS` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/DeleteRowBusiness.xml` | hapus baris bisnis |
| `DeleteSecurityReinsurer` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteSecurityReinsurer.xml` | hapus security |
| `BrowseDeleteRowTreatyInContract` | `@BASECLASS` / `BROWSEDELETEROWTREATYINCONTRACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/BrowseDeleteRowTreatyInContract.xml` | ⚠️ bagian **salin**-nya tidak dimigrasikan |

`[terverifikasi]` Kaskade existing — empat `DELETE` berurut dalam SQL mentah, lalu `COMMIT`:

1. `DELETE FROM treatycontract WHERE id = … AND IDTREATYYEAR = …`
2. `DELETE FROM treatybusiness WHERE TREATYYEAR = … AND (TREATYYEARID = … OR TREATYYEARID IS NULL)
   AND TREATYGROUPID = … AND REINSTYPEID = …`
3. `DELETE FROM MTREATYSECURITY WHERE REAS_ID IN (SELECT id FROM TREATYREINSURER a WHERE …)`
4. `DELETE FROM TREATYREINSURER WHERE TreatyYear = … AND TreatyGroupID = … AND ReinsTypeID = …`

⚠️ **`PROPORTIONALARRG` tidak ikut** — dan itu **benar**.

⚠️ **Penyimpangan sadar 4 — kaskade + popup; klausul sengaja dikecualikan.**
`[fakta bisnis — work owner]` Klausul menggantung pada **(TreatyYear, TreatyGroupID, ReinsTypeID)**
dan **milik level tahun/grup/jenis**, dipakai bersama lintas kontrak. Ketidakikutannya adalah
**desain, bukan bug**.

#### ADR terkait

**ADR-U-0007** (jejak audit setiap transisi), **ADR-U-0015** (kegagalan ditangani eksplisit).

#### Acceptance criteria

- [ ] ⚠️ Menghapus kontrak **mengkaskade** ke `TREATYBUSINESS`, `MTREATYSECURITY`, dan
      `TREATYREINSURER`. *(AC 42 spec; User story 26; penyimpangan sadar 4)*
- [ ] ⚠️ Penghapusan didahului **popup konfirmasi Ya/Batal** yang **menyebut jumlah baris tiap
      jenis** yang akan ikut terhapus. Memilih **Batal** membatalkan seluruhnya dan **tidak
      mengubah apa pun**. *(AC 43 spec; User story 27; penyimpangan sadar 4)*
- [ ] ⚠️ **`PROPORTIONALARRG` (klausul) TIDAK ikut terhapus** dan tetap dapat dibaca setelah
      kontrak dihapus. Test yang menemukan klausul ikut terhapus **gagal**. *(AC 44 spec;
      User story 23; `[fakta bisnis — work owner]` — desain, bukan bug)*
- [ ] Popup **menyebut secara eksplisit** bahwa klausul **tidak** akan terhapus, supaya pengguna
      tidak menyangka sebaliknya. *(turunan AC 43, 44)*
- [ ] Seluruh kaskade berjalan dalam **satu transaksi**; kegagalan di langkah mana pun
      **membatalkan seluruhnya**. *(AC 45 spec; tiket 09)*
- [ ] Jumlah yang ditampilkan popup **sama** dengan jumlah yang benar-benar terhapus — dihitung dari
      kombinasi yang sama, bukan dari perkiraan.
- [ ] ⚠️ Kaskade menyentuh **satu keluarga tabel saja** — tidak ada tabel kembar JSON yang ikut
      atau tertinggal. *(AC 63, 64 spec; penyimpangan sadar 1)*
- [ ] Penghapusan mencatat **jejak audit** — siapa, kapan, dan berapa baris tiap jenis.
      *(**ADR-U-0007**)*
- [ ] Penghapusan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38 spec;
      **ADR-U-0015**)*

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Existing tidak konsisten terhadap tabel kembar.** `[terverifikasi]` Kaskade ini menghapus
`treatycontract` tetapi **bukan** `M_TREATYCONTRACT`, sedangkan
`RDBList/DeleteRowBusinessList.xml` menghapus **keduanya**. Inkonsistensi itu **lenyap dengan
sendirinya** begitu penyimpangan sadar 1 diterapkan (tiket 01) — jangan mencoba "memperbaikinya"
dengan menambahkan penghapusan tabel JSON.

⚠️ **`DeleteFromTreatyReinsurer_Act` tanpa `COMMIT`** `[terverifikasi]` — konsisten dengan temuan
bahwa modul ini menyerahkan batas transaksi kepada pemanggil. Di sistem baru, pemanggil itu adalah
pembungkus transaksi tiket 09.

⚠️ **`InputData.CARI19` dipakai untuk dua arti berbeda di existing** `[terverifikasi]`:
`ReinsTypeID` pada kaskade hapus ini, tetapi `userId` pada jalur salin
(`RDBList/SaveMasterCopyData_SQL.xml`). Jebakan bagi migrasi yang meniru nama slot — di sistem baru
setiap nilai punya namanya sendiri. Jalur salin sendiri **tidak dimigrasikan** (tiket 03).

⚠️ **`TREATYYEARID` boleh NULL** `[terverifikasi]` — kaskade existing menulis
`(TREATYYEARID = … OR TREATYYEARID IS NULL)`. Kaskade baru harus **tahan terhadap NULL** pada
kolom itu, bukan mengandaikannya selalu terisi; bila tidak, baris bisnis lama akan tertinggal.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — kaskade yang **mengecualikan satu tabel**
adalah pernyataan tentang basis data; memalsukannya berarti tidak mengujinya. Uji dua arah: tiga
anak **hilang**, dan klausul **masih ada**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Treaty Contract Out - 11 - Kurs USD → IDR

**Status:** ready-for-agent

**Blocked by:** 08 (kurs melekat pada baris klausul)

#### Hasil & nilai pengguna

Sebagai **underwriter**, saya ingin nilai dalam **USD** dapat dilihat padanannya dalam **IDR**
menurut kurs yang berlaku **pada periode kontrak**; dan sebagai **Finance**, saya ingin nilai
**Rp** dan **USD** tetap tercatat sebagai **dua nilai terpisah**, karena keduanya memang dua angka
yang berbeda. *(User story 28–30 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Baca kurs per periode dari master (read-only) |
| `internal/services` | Konversi USD → IDR; rujukan mata uang |
| `internal/handlers` | Endpoint kurs berlaku |
| `frontend/` | Tampilan padanan IDR pada baris klausul |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `testingKurs` | `@BASECLASS` / `TESTINGKURS` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/testingKurs.xml` | ⚠️ **jalur yang benar-benar dipakai** — nol langkah di-remark |
| `RefreshKurs` | `@BASECLASS` / `REFRESHKURS` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/RefreshKurs.xml` | muat ulang kurs |
| `GetMasterKursList` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterKursList.xml` | baca kurs per periode |
| `NitipKurs` | `RULE-HTML-SECTION` | `Treaty Contract Out/Section/NitipKurs.xml` | tampilan kurs |
| `SetTreatyArrangementDesc_Act` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `SETTREATYARRANGEMENTDESC_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetTreatyArrangementDesc_Act.xml` | ⚠️ **MATI** — 5 dari 6 langkah di-remark |

`[terverifikasi]` Kueri kurs existing:

```
select TOIDR from treatyexchange
 where Quarter = '0'
   and to_date({tanggal},'YYYYMMDD') BETWEEN trunc(TO_TIMESTAMP_TZ(STARTDATE, 'YYYYMMDD"T"HH24MISS.FF3 TZR'))
                                         AND trunc(TO_TIMESTAMP_TZ(ENDDATE,   'YYYYMMDD"T"HH24MISS.FF3 TZR'))
   and IDCURRENCY = '10001'
```

`[data DBA]` **Nama tabel sebenarnya `TREATYEXCHANGEYEARLY`** (bukan `TREATYEXCHANGE`); kolom
`TOIDR`, `TOUSD`, `IDCURRENCY`, `CURRENCY`, `QUARTER`, `STARTDATE`, `ENDDATE` — **seluruhnya
`VARCHAR2`**, termasuk kedua tanggal.

⚠️ **Penyimpangan sadar 7 — `IDCURRENCY='10001'` jadi rujukan master, bukan literal.**
`[keputusan work owner]` `IDCURRENCY = '10001'` **berarti USD**, dan konversi **USD → IDR** memang
aturan bisnis. Yang tidak boleh adalah identitas mata uangnya ditanam sebagai konstanta program.

⚠️ **Penyimpangan sadar 8 — nama jujur.** `[keputusan work owner]` Komponen di sistem baru
**tidak** mengandung kata `"testing"`.

#### ADR terkait

**ADR-U-0003** (uang non-float), **ADR-U-0015** (kegagalan ditangani eksplisit).

#### Acceptance criteria

- [ ] Padanan **IDR** untuk nilai **USD** dihitung dari kurs yang **berlaku pada periode** kontrak —
      tanggal berada di antara tanggal mulai dan tanggal akhir baris kurs. *(AC 46 spec;
      User story 28)*
- [ ] ⚠️ Identitas mata uang USD (`IDCURRENCY = '10001'`) adalah **rujukan ke master mata uang**,
      **bukan literal di kode**. Test yang menemukan `'10001'` sebagai konstanta program **gagal**.
      *(AC 47 spec; penyimpangan sadar 7)*
- [ ] `QUARTER = '0'` **diikuti apa adanya**; artinya dicatat sebagai `[terbuka]` dan **tidak
      ditebak**. *(AC 48 spec)*
- [ ] ⚠️ Nilai **Rp** dan **Usd** tetap **dua nilai terpisah** pada baris klausul — bukan satu nilai
      dengan kode mata uang. *(AC 49 spec; User story 30)*
- [ ] ⚠️ Jalur kurs yang dimigrasikan adalah yang **benar-benar dipakai** (`testingKurs`), bukan
      yang namanya lebih wajar tetapi **di-remark**. Nama komponen barunya **tidak mengandung kata
      "testing"**. *(AC 50 spec; penyimpangan sadar 8; **OQ-066**)*
- [ ] Tanggal mulai/akhir baris kurs diperlakukan sebagai **tanggal**, bukan teks yang di-parse
      setiap kueri. *(AC 53 spec; penyimpangan sadar 6)*
- [ ] Nilai kurs diperlakukan sebagai **desimal presisi arbitrer**; **tidak** melewati `float`.
      *(AC 51 spec; **ADR-U-0003**)*
- [ ] Master kurs **tidak ditulis** oleh konteks ini. Test yang menemukan tulisan ke
      `TREATYEXCHANGEYEARLY` **gagal**.
- [ ] Periode yang **tidak punya baris kurs** menghasilkan kegagalan yang **terlihat**, bukan nilai
      nol atau kosong yang diam. *(**ADR-U-0015**)*

#### Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **Arti `QUARTER = '0'`** — OQ kecil, **tidak memblokir**: AC di atas mewajibkan
nilainya diikuti apa adanya, bukan dipahami.

#### Catatan

⚠️ **OQ-066 terbalik di sini.** `[terverifikasi]` Activity bernama **`testingKurs`**
(`@BASECLASS!TESTINGKURS`) **nol langkah di-remark** — ia hidup dan dipakai. Sebaliknya
`SetTreatyArrangementDesc_Act` (`ASM-FW-GISFW-INT-PROPORTIONALARRG`), yang namanya paling wajar,
**lima dari enam langkahnya di-remark** — ia mati. Biasanya nama "testing" menandai yang mati; di
sini justru sebaliknya. **Baca kodenya, jangan namanya.**

⚠️ **Nama tabel di korpus dan di basis data berbeda.** Kueri Pega menyebut `treatyexchange`;
`[data DBA]` menyebut tabel sebenarnya **`TREATYEXCHANGEYEARLY`**. Migrasi (tiket 01) memakai nama
DBA.

⚠️ **Dua kolom uang dipertahankan dengan sengaja.** `[keputusan work owner]` `Rp` dan `Usd` bukan
satu nilai yang ditampilkan dua cara — keduanya **dua angka berbeda** yang disimpan berdampingan
(dan untuk baris anak, keduanya **turunan** dari `Pct` × nilai induk — lihat tiket 08). Jangan
"merapikannya" jadi satu kolom + kode mata uang.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — pemilihan baris kurs menurut periode
adalah pertanyaan rentang tanggal di basis data.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Treaty Contract Out - 12 - Lampiran di tahun treaty — **FITUR BARU**

**Status:** ready-for-agent

**Blocked by:** 03 (tahun treaty sebagai induk lampiran)

⚠️ **Ini fitur baru, bukan migrasi.** `[keputusan work owner]` Jalur lampiran di Pega
(`M_ATTACHMENTTREATY_2`) **belum rampung di-develop**. Di sistem baru fungsinya **dilengkapi** dan
melekat pada **tahun treaty**.

#### Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin **melampirkan berkas pada sebuah tahun treaty** beserta
kategorinya, supaya dokumen kontraknya tersimpan bersama datanya — **tanpa** itu menjadi syarat
tersimpannya data tahun treaty; dan sebagai **tim operasi**, saya ingin **kegagalan unggah terlihat
dan dapat diulang**, serta alamat penyimpanan **dibaca dari konfigurasi saat dijalankan**.
*(User story 31–36 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Rekam lampiran; resolusi alamat penyimpanan; **cache token** |
| `internal/clients` | Klien penyimpanan berkas **di balik interface** |
| `internal/services` | Orkestrasi efek keluar; penandaan kegagalan; pengulangan |
| `internal/handlers` | Endpoint unggah / unduh / hapus / status |
| `frontend/` | Panel lampiran di layar tahun treaty; penanda "terkirim" / "tertunda" |

#### Rule Pega sumber

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

#### ADR terkait

**ADR-U-0013** (alamat di-resolve **runtime**; **dilarang** sebagai literal, konstanta, **maupun env
var** — menggantikan ADR-U-0004), **ADR-U-0015** (efek keluar — kegagalan tidak boleh diam-diam),
**ADR-U-0010** (penyimpanan berkas tetap Google Storage), **ADR-U-0007** (jejak audit).

#### Acceptance criteria

- [ ] ⚠️ Berkas dapat **dilampirkan pada sebuah tahun treaty** — fitur yang **belum ada** di Pega.
      *(AC 54 spec; User story 31; penyimpangan sadar 9)*
- [ ] ⚠️ Lampiran **opsional**: kegagalan unggah **tidak** membatalkan tersimpannya tahun treaty —
      hanya status lampiran yang berubah. *(AC 55 spec; User story 32)*
- [ ] Tiap lampiran dapat diberi **kategori** dari master kategori. *(AC 56 spec; User story 33)*
- [ ] Lampiran dapat **diunduh** dan **dihapus**. *(AC 57 spec; User story 34)*
- [ ] ⚠️ Kegagalan unggah **tercatat dan dapat diulang**; pengulangan **tidak** menggandakan berkas.
      *(AC 58 spec; User story 35; **ADR-U-0015**)*
- [ ] Alamat penyimpanan di-resolve **runtime** dari konfigurasi. Test yang memindai kode untuk URL
      sebagai **literal, konstanta, atau pembacaan env var** **gagal** bila menemukannya.
      *(AC 59 spec; User story 36; **ADR-U-0013**)*
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
- [ ] Setiap unggah dan penghapusan mencatat **jejak audit**. *(**ADR-U-0007**)*

#### Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **OQ-047 — tidak memblokir:** alamat fisik Google Storage tersimpan di
`M_LINK_SERVICE`, dan isinya belum dibaca. **ADR-U-0013** sudah mengatur **cara** meresolusinya, jadi
tiket ini dapat selesai tanpa mengetahui alamatnya — yang mengikat adalah **alamat tidak ditanam**.

#### Catatan

⚠️ **Taruhannya berbeda dari Komite Claim Life.** Di sana efek keluar **wajib berhasil**
(transactional outbox, **ADR-U-0015** versi Komite). Di sini `[keputusan work owner]` lampiran
**opsional** — yang wajib adalah **kegagalannya terlihat dan dapat diulang**, bukan berhasil.
Sama dengan perlakuan di Master Product Name Life.

⚠️ **Fitur ini tidak lengkap di Pega** `[keputusan work owner]` — jadi korpus bukan sumber
kebenaran perilaku di sini, hanya sumber **rantai teknis** (token, alamat, rekam berkas, kategori).
Perilaku yang tidak ada di korpus **tidak ditebak**; ia ditetapkan oleh AC di atas.

⚠️ **Satu-satunya efek keluar konteks ini.** `[terverifikasi]` Modul ini punya **satu**
`RULE-CONNECT-REST`, dan tidak ada integrasi Arasapas / Kasir / Konversi / Gemini AI. Berbeda dari
Master Contract Retro Life yang **tanpa `ConnectREST` sama sekali** — di sana ADR-U-0013 dan ADR-U-0015
tidak berlaku; di sini berlaku keduanya.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**; **penyimpanan berkas di-fake** di balik
interface — ia satu-satunya yang dipalsukan di seluruh konteks ini.

```
go test ./internal/...
cd frontend && npm test
make check
```
