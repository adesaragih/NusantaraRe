# 01: Skema relasional + migrasi + tipe dirapikan — **PREFACTOR**

**Status:** ready-for-agent

**Blocked by:** **CL-01** (kerangka aplikasi + seam API — scaffolding lintas konteks, tidak dibuat
di sini)

⚠️ **Ini PREFACTOR, dan ia tiket PERTAMA.** Bentuk skema berubah di tiga sumbu sekaligus —
tabel JSON dibuang, `MTREATYSECURITY` dibersihkan, dan tipe kolom diperbaiki — sehingga **tidak ada
irisan lain yang dapat berdiri** sebelum bentuk barunya ada. *"Make the change easy, then make the
easy change."*

## Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin seluruh data master arrangement treaty non-life pindah ke bentuk
yang dapat dipercaya — uang sebagai angka, tanggal sebagai tanggal, setiap baris punya identitas —
tanpa kehilangan satu nilai pun, dan **tanpa mengambil apa pun dari tabel JSON yang sudah mati**.
*(User story 37–40 di spec)*

Dan sebagai **konteks hilir** (`Claim Prop`, `Komite Claim Prop`, `Claim Fac In`), saya ingin tetap
dapat **membaca** master arrangement dalam bentuk yang saya kenal sampai saya ikut bermigrasi.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL delapan tabel + sequence; skrip migrasi & rekonsiliasi |
| `internal/models` | Bentuk tahun treaty, kontrak, reinsurer, security, business, klausul |
| — | Skrip rekonsiliasi nilai uang & tanggal |

## Keadaan lama `[data DBA]`

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

## Bentuk baru

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

### Kolom `proportionalarrg` — 35 kolom `[data DBA]`

```
ID, TREATYYEAR, TREATYYEARID, TREATYGROUPID, TREATYGROUPNAME, TREATYDESCID, TREATYDESCNAME,
REINSTYPEID, REINSTYPENAME, LAYER, LAYERPART, LAYERPARTTYPE, LAYERTYPE, KURS, TGLUPDATE, USERID,
LINE, PCT, PCTME, YDCF, METHOD, TERRITORIALLIMIT, PARENTREINSTYPEID, SPREADINGORDER, RP, USD,
ID_OCCUPATION, OCCUPATION, ID_CLAUSE, CLAUSE, TREATYLIMIT, COINS_MIN, COINS_MAX, MORERP, MOREUSD
```

Sembilan kolom terakhir yang bercetak — `ID_OCCUPATION`, `OCCUPATION`, `ID_CLAUSE`, `CLAUSE`,
`TREATYLIMIT`, `COINS_MIN`, `COINS_MAX`, `MORERP`, `MOREUSD` — hanya diisi baris **induk**; baris
"anak" mengisinya **NULL**.

### Sequence `[data DBA]`

| Tabel | Sequence | Format identitas |
| --- | --- | --- |
| `PROPORTIONALARRG` | `PROPORTIONALARRG_SEQ` | `'1' + lpad(seq, 7, '0')` — **7 digit** |
| `TREATYCONTRACT` | `treatycontract_seq` | `'1' + lpad(seq, 6, '0')` |
| `TREATYYEAR` | `TreatyYear_seq` | `'1' + lpad(seq, 6, '0')` |
| `TREATYREINSURER` | `M_TREATYREINSURER_SEQ` | `'1' + lpad(seq, 6, '0')` |
| `TREATYBUSINESS` | `TREATY_BUSINESS_SEQ` | `'1' + lpad(seq, 6, '0')` |

## Rule Pega sumber

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

## ADR terkait

**ADR-0003** (uang non-float), **ADR-0006** (identitas lewat sequence basis data),
**ADR-0009** (migrasi penuh; koeksistensi ditolak untuk data — **bukan** untuk kontrak baca hilir).

## Acceptance criteria

- [ ] Modul ini adalah **satu-satunya penulis** `TREATYYEAR`, `TREATYCONTRACT`, `TREATYREINSURER`,
      `MTREATYSECURITY`, `TREATYBUSINESS`, `PROPORTIONALARRG`. *(AC 1 spec; OQ-042)*
- [ ] Konteks hilir (`Claim Prop`, `Komite Claim Prop`, `Claim Fac In`) tetap dapat **membaca**
      bentuk relasional yang mereka pakai hari ini. *(AC 2 spec)*
- [ ] Skema **tidak memuat** satu pun objek treaty **outward**. *(AC 3 spec)*
- [ ] ⚠️ **Seluruh** nilai uang — `RP`, `USD`, `MORERP`, `MOREUSD`, `TREATYLIMIT`, `COINS_MIN`,
      `COINS_MAX`, `PCT_SHARE`, `RICOMM` — bertipe **desimal presisi arbitrer**; **tidak** melewati
      `float`. Test yang menemukan kolom uang bertipe teks **gagal**. *(AC 51 spec; **ADR-0003**;
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
      **ADR-0009**)*
- [ ] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara
      tepat**, bukan dengan toleransi. *(AC 66 spec; **ADR-0003**)*
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
      bertabrakan** dengan yang lama. *(AC 69 spec; **ADR-0006**)*
- [ ] `[terbuka]` Kolom `PROPORTIONALLIST` dan `OBJECT` **tidak dibawa** ke skema baru kecuali
      migrasi membuktikan ada data hidup di sana; temuannya **dilaporkan**. *(AC 70 spec)*
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.
      *(AC 71 spec)*

## Blocker

**Tidak ada pemblokir.** `[data DBA]` DDL delapan tabel dan body enam procedure sudah diterima —
**OQ-001 dan OQ-002 ditutup**.

## Catatan

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

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — presisi desimal, konversi tanggal, dan
keutuhan rujukan setelah pemberian PK **hanya berperilaku benar pada basis data sungguhan**;
memalsukannya berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```
