# Telusur — Master Contract Retro Life (modul **tanpa** rule `Flow`)

STEP D2, Tahap 5 konteks #20 — **konteks terakhir D2**. Ditelusur 2026-09-13.
Konvensi: **`_METHOD.md`** + **`_METHOD-noflow.md`**.
**66 file** — **modul terkecil di korpus** (`find "Master Contract Retro Life" -name '*.xml' | wc -l`).

**Titik masuk `[terverifikasi]`:** `Harness/InboxRetroLifeReinsurersList.xml`
→ `RULE-HTML-HARNESS` / **`DATA-PORTAL`** / `INBOXRETROLIFEREINSURERSLIST`, 530.052 byte.

---

## 1. Alur — empat layar master, tanpa proses

### 1.1 Empat Harness = empat grid master `[terverifikasi]`

| Harness | `pxInsName` | Ukuran |
| --- | --- | ---: |
| `InboxRetroLifeReinsurersList.xml` | `DATA-PORTAL!INBOXRETROLIFEREINSURERSLIST` | 530.052 |
| `InboxRetroLimitReinsurers.xml` | `DATA-PORTAL!INBOXRETROLIMITREINSURERS` | 575.179 |
| `InboxBusinessLifeReinsurers.xml` | `DATA-PORTAL!INBOXBUSINESSLIFEREINSURERS` | 486.911 |
| `InboxSecurityReinsurerLife.xml` | `DATA-PORTAL!INBOXSECURITYREINSURERLIFE` | 444.874 |

`[terverifikasi]` Empat Harness untuk 66 file — **rasio Harness tertinggi di korpus**. Modul ini
pada dasarnya **empat layar grid** di atas satu keluarga tabel `*_LIFE`.

### 1.2 Hanya dua FlowAction

`[terverifikasi]` `FlowAction/ViewRate.xml` dan `FlowAction/ViewRateTable.xml` — keduanya
**tampilan**, bukan aksi proses. Sama sedikitnya dengan `Treaty Contract Out` (2 FlowAction).

### 1.3 Pola CRUD simetris — sama dengan `Treaty Contract Out`

`[terverifikasi]` 27 Activity, seluruhnya berpasangan per entitas:

| Pola | Activity |
| --- | --- |
| **`New…`** (4) | `NewInputBusinessLife_Act`, `NewInputSecurityLife_Act`, `NewInputTreatyLimit_Life`, `NewInputTreatyYear_Life_Act`, `NewTreatyReinsurerDetail_Act` |
| **`Set…`** (7) | `SetBusinessListLife_Act`, `SetRetroListLife_Act`, `SetSecurityLife_Act`, `SetSecurityReinsurerLife_Act`, `SetTreatyYearLife_Act`, `SetValueRetroLimit_TreatyYearLife`, `SetParamRate`, `SetParamRateTable` |
| **`Save…`** (6) | `SaveBusinessLife_Act`, `SaveBusinessToAllLife_Act`, `SaveSecurityLife_Act`, `SaveSecurityReinsurerLife_Act`, `SaveTreatyLimit_Act`, `SaveTreatyYearLife_Act` |
| **`Delete…`** (4) | `DeleteRowBusiness`, `DeleteSecurityLife_Act`, `DeleteSecurityReinsurerLife_Act`, `DeleteTreatyLimit_Act` |
| **`CancelActivity…`** (1) | `CancelActivityTreatyContract` |
| Perhitungan | `CountingPercentShare_Act` (59.429 B), `TreatyLimit_TypeProtect` (54.893 B) |
| Pesan | `SetErrorMessageReinsurer` |

`[terverifikasi]` **Tidak ada `StatusAkseptasi`, `Position`, `Akseptasi_DT`, maupun tombol
Submit/Decline.** Perintah audit:
```
grep -rl "StatusAkseptasi" "Master Contract Retro Life" --include="*.xml" | wc -l   # 0
ls "Master Contract Retro Life/DataTransform/"                                       # 1 file
```

**Tahapan yang dapat dinyatakan `[terverifikasi]`:** buka salah satu dari empat grid → tambah/ubah/
hapus baris (`New…` → `Set…` → `Save…` / `Delete…`) → simpan ke master `*_LIFE` lewat Connect-SQL.
**Tidak ada tangga persetujuan dan tidak ada status proses.**

`[terverifikasi]` Modul ini **satu-satunya di korpus tanpa rule `ConnectREST`, tanpa `When`, dan
tanpa `SystemSettings`** (`ls "Master Contract Retro Life"` → hanya Activity, Claude outputs,
DataTransform, FlowAction, Harness, RDBList, ReportDefinition, Section).

---

## 2. Kontrak retro Life — apa adanya

### 2.1 Objek Oracle: **seluruhnya keluarga `_LIFE`** `[terverifikasi]`

```
awk -F'\t' '$1 ~ "^Master Contract Retro Life/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
  | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

| Objek | Rule perujuk |
| --- | ---: |
| `POOLDATA.TREATYREINSURER_LIFE` | 2 |
| `POOLDATA.TREATYCONTRACT_LIFE` | 2 |
| `POOLDATA.TREATYBUSINESS_LIFE` | 2 |
| `POOLDATA.TREATYSECURITYREINSURER_LIFE` | 1 |

`[terverifikasi]` **Hanya 4 objek, seluruhnya bersufiks `_LIFE`, seluruhnya di skema `POOLDATA`.**
Tidak ada `M_TREATY_IN`, tidak ada objek facultative, tidak ada objek outward, **tidak ada db-link**.
Modul paling terisolasi di korpus.

### 2.2 Dua belas Connect-SQL, lima procedure penulis `[terverifikasi]`

| Rule | Procedure |
| --- | --- |
| `RDBList/SaveMasterTreatyContract_Life_SQL.xml` | `POOLDATA.INSERTTREATYCONTRACT_LIFE` |
| `RDBList/SaveMasterTreatyYear_Life_SQL.xml` | `POOLDATA.INSERTTREATYYEAR_LIFE` |
| `RDBList/SaveMasterTreatyReinsurer_Life_SQL.xml` | `POOLDATA.INSERTREINSURER_LIFE` |
| `RDBList/SaveMasterTreatyBusiness_Life_SQL.xml` | `POOLDATA.INSERTBUSINESS_LIFE` |
| `RDBList/SaveMasterTreatySecurityReinsurer_Life_SQL.xml` | `POOLDATA.INSERTSECURITYREINSURER_LIFE` |
| `RDBList/SaveTreatyBusinessAll_Life_SQL.xml`, `DeleteRowBusinessList.xml`, `DeleteSecurityReinsurer(Life)_SQL.xml`, `DeleteTreatyLimit_SQL.xml`, `GetMasterReinsurerLifeList_SQl.xml`, `GetTreatyContract_life.xml` | — |

`[terverifikasi]` Kelima procedure **hanya muncul di modul ini** di seluruh korpus
(`awk -F'\t' '$8 ~ /INSERTTREATYCONTRACT_LIFE|INSERTREINSURER_LIFE|INSERTBUSINESS_LIFE|INSERTSECURITYREINSURER_LIFE|INSERTTREATYYEAR_LIFE/{print $1}' all-rules.tsv`).

**Seluruh jalur tulis melewati stored procedure** — sama dengan `Treaty Contract Out`. Aturan
penyimpanan kontrak retro life **berada di sisi database, tidak ada di korpus** (OQ-002).

### 2.3 `REINSTYPEID` dan parameter procedure

`[terverifikasi]` `REINSTYPEID` muncul **129 kali** di modul ini
(`grep -rho "REINSTYPEID" "Master Contract Retro Life" --include="*.xml" | wc -l`). Bentuk
pemakaiannya, apa adanya dari blok PL/SQL:

```
DECLARE
  … REINSTYPEID VARCHAR2(32767);
  REINSTYPEID := {TempInputData.CARI3};
  REINSTYPEID := {TempInputData.CARI4};
  REINSTYPEID := {TempInputData.CARI5};
  … ( p_REINSTYPEID, p_REINSTYPENAME, p_TREATYSTARTDATE, p_TREATY… )
  … ( p_REINSTYPEID, p_REINSTYPENAME, p_REINSURERID, p_REINSURERN… )
  … ( p_REINSTYPEID, p_REINSTYPENAME, p_BIZCODE, p_BIZNAME, p_RIR… )
```

`[terverifikasi]` `REINSTYPEID` diisi dari properti generik `TempInputData.CARI3/4/5` — **bukan
dari konstanta**, sehingga **nilai literalnya tidak muncul di korpus sama sekali**. Ia diteruskan
sebagai parameter `p_REINSTYPEID` ke tiga procedure berbeda, bersama `p_REINSTYPENAME`.

**Arti dan daftar nilai `REINSTYPEID` tidak dapat dinyatakan** — tidak ada satu pun nilai literal
untuk dikutip. → OQ-020, dan **OQ-057** (masalah yang sama dengan kode `OR` di
`Master Product Name Life`).

`[terverifikasi]` Properti bernama `CARI1`…`CARI5` dipakai sebagai **slot parameter generik**
menuju SQL — pola yang sama terlihat di `NB Treaty In` (`InputData.CARI20`, `CARI21`, `CARI3`)
dan `Treaty In Adjustment` (`InputData.CARI1`). → **OQ-059**.

### 2.4 Perhitungan share dan tipe proteksi

`[terverifikasi]` Dua activity perhitungan, **keduanya di class `@BASECLASS`** dan **belum dibaca**:

| Activity | `pxInsName` | Ukuran | Langkah |
| --- | --- | ---: | ---: |
| `CountingPercentShare_Act.xml` | `@BASECLASS!COUNTINGPERCENTSHARE_ACT` | 59.429 | 4 |
| `TreatyLimit_TypeProtect.xml` | `@BASECLASS!TREATYLIMIT_TYPEPROTECT` | 54.893 | 4 |
| `SetValueRetroLimit_TreatyYearLife.xml` | `@BASECLASS!SETVALUERETROLIMIT_TREATYYEARLIFE` | 63.678 | 2 |

**Rumus pembagian share retro dan jenis proteksi limit tidak dinyatakan** — batas cakupan telusur.

---

## 3. Status / state dan kode

`[terverifikasi]` **Tidak ada properti status proses.** Yang ada:

| Properti | Nilai literal | Arti |
| --- | --- | --- |
| `REINSTYPEID` | **tidak ada nilai literal di korpus** — diisi dari `TempInputData.CARI3/4/5` | **tidak dapat dinyatakan** (§2.3) |
| `REINSTYPENAME` | diteruskan sebagai `p_REINSTYPENAME` | **belum terverifikasi** |
| `BIZCODE`, `BIZNAME`, `REINSURERID`, `REINSURERNAME`, `TREATYSTARTDATE` | parameter procedure | **belum terverifikasi** |

Perintah audit:
```
grep -rhoE "REINSTYPEID[^<]{0,46}" "Master Contract Retro Life" --include="*.xml" \
  | sed 's/&amp;#61;/=/g;s/&amp;quot;/"/g' | sort -u
```

---

## 4. Integrasi eksternal

`[terverifikasi]` **Tidak ada.** Modul ini tidak punya rule `RULE-CONNECT-REST` sama sekali —
satu-satunya modul di korpus yang demikian
(`ls "Master Contract Retro Life/ConnectREST" 2>&1` → direktori tidak ada).

Tidak ada Google Storage, tidak ada `M_LINK_SERVICE`, tidak ada Arasapas/Kasir/Konversi/Gemini AI,
tidak ada db-link. **Seluruh interaksinya dengan dunia luar adalah 12 Connect-SQL ke `POOLDATA`.**

---

## 5. Batas pengetahuan

| # | Batas | Sifat |
| ---: | --- | --- |
| 1 | Isi 5 stored procedure `INSERT*_LIFE` — **seluruh aturan penyimpanan** | korpus (OQ-002) |
| 2 | Nilai dan arti `REINSTYPEID` — tidak ada literal untuk dikutip | korpus (OQ-020, OQ-057) |
| 3 | Rumus `CountingPercentShare_Act`, `TreatyLimit_TypeProtect`, `SetValueRetroLimit_TreatyYearLife` | **cakupan telusur** |
| 4 | Isi 4 Harness (444–575 KB) — hanya di-grep | **cakupan telusur** |
| 5 | Arti slot `TempInputData.CARI1`…`CARI5` | korpus (OQ-059) |

`[terverifikasi]` **Tidak ada kolom JSON** di modul ini
(`grep -rho "JSONDATA\|DATA_JSON" "Master Contract Retro Life" --include="*.xml" | wc -l` → **0**) —
satu-satunya modul Tahap 5 yang bebas dari OQ-012.

`[terverifikasi]` Folder **`Claude outputs`** berisi `Struktur_GridRetrocessionLife.xlsx`
(11.944 B) dan **`perubahan_skill.diff`** (10.278 B) — berkas non-Pega di dalam korpus READ-ONLY.
**Tidak dibaca.** → OQ-054.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Master Contract Retro Life/`) | OQ-011? |
| --- | --- | --- |
| `DATA-PORTAL` / `INBOXRETROLIFEREINSURERSLIST` / `RULE-HTML-HARNESS` | `Harness/InboxRetroLifeReinsurersList.xml` | tidak |
| `DATA-PORTAL` / `INBOXRETROLIMITREINSURERS` / `RULE-HTML-HARNESS` | `Harness/InboxRetroLimitReinsurers.xml` | tidak |
| `DATA-PORTAL` / `INBOXBUSINESSLIFEREINSURERS` / `RULE-HTML-HARNESS` | `Harness/InboxBusinessLifeReinsurers.xml` | tidak |
| `DATA-PORTAL` / `INBOXSECURITYREINSURERLIFE` / `RULE-HTML-HARNESS` | `Harness/InboxSecurityReinsurerLife.xml` | tidak |
| `…` / `SAVEMASTERTREATYCONTRACT_LIFE_SQL` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterTreatyContract_Life_SQL.xml` | tidak |
| `…` / `SAVEMASTERTREATYYEAR_LIFE_SQL` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterTreatyYear_Life_SQL.xml` | tidak |
| `…` / `SAVEMASTERTREATYREINSURER_LIFE_SQL` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterTreatyReinsurer_Life_SQL.xml` | tidak |
| `…` / `SAVEMASTERTREATYBUSINESS_LIFE_SQL` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterTreatyBusiness_Life_SQL.xml` | tidak |
| `…` / `SAVEMASTERTREATYSECURITYREINSURER_LIFE_SQL` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterTreatySecurityReinsurer_Life_SQL.xml` | tidak |
| `…` / `SAVETREATYBUSINESSALL_LIFE_SQL`, `GETMASTERREINSURERLIFELIST_SQL`, `GETTREATYCONTRACT_LIFE`, `DELETEROWBUSINESSLIST`, `DELETESECURITYREINSURERLIFE_SQL`, `DELETESECURITYREINSURER_SQL`, `DELETETREATYLIMIT_SQL` / `RULE-CONNECT-SQL` | `RDBList/` | sebagian |
| `@BASECLASS` / `COUNTINGPERCENTSHARE_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/CountingPercentShare_Act.xml` | tidak |
| `@BASECLASS` / `TREATYLIMIT_TYPEPROTECT` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyLimit_TypeProtect.xml` | tidak |
| `@BASECLASS` / `SETVALUERETROLIMIT_TREATYYEARLIFE` / `RULE-OBJ-ACTIVITY` | `Activity/SetValueRetroLimit_TreatyYearLife.xml` | tidak |
| `…` / 4 `NEWINPUT*`, 7 `SET*`, 6 `SAVE*`, 4 `DELETE*`, `CANCELACTIVITYTREATYCONTRACT`, `SETERRORMESSAGEREINSURER` / `RULE-OBJ-ACTIVITY` | `Activity/` | sebagian |
| `…` / `VIEWRATE`, `VIEWRATETABLE` / `RULE-OBJ-FLOWACTION` | `FlowAction/` | tidak |

**23 rule ditelusur** (4 Harness; 12 Connect-SQL; 3 Activity perhitungan; 2 FlowAction; kelompok
CRUD dihitung sebagai satu kelompok). Dari 66 file, **5 entri terdaftar di register OQ-011**
(`grep -cF "Master Contract Retro Life/" ../inventory/_oq011-konflik-isi.md` → 5) — jumlah paling
sedikit di korpus.

Rule **dirujuk tapi belum ditelusur**: 12 ReportDefinition, 8 Section, sisa 24 Activity.

---

## 7. Pertanyaan terbuka

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-059** | Properti generik `TempInputData.CARI1`…`CARI5` / `InputData.CARI*` dipakai sebagai slot parameter menuju SQL di beberapa modul — artinya per pemanggilan tidak terbaca | DBA + Product+UW |

OQ dikuatkan: **OQ-002** (§2.2 — 5 procedure eksklusif, seluruh tulis lewat SP),
**OQ-005** (§1 — titik masuk ditetapkan), **OQ-009** (§2.4 — activity perhitungan di `@BASECLASS`),
**OQ-016** (`POOLDATA`), **OQ-020** (§3), **OQ-054** (§5), **OQ-057** (§2.3).

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **0 file**,
`OperatorID.pyTelephone` **0**, `pxCreateOperator` sebagai guard **0**. `IsPEGAPROD` **tidak ada**.
**Modul tanpa jejak guard identitas sama sekali** — satu-satunya di Tahap 5.
