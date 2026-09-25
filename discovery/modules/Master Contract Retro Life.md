# Modul — Master Contract Retro Life

STEP D3. Sintesis dari `../inventory/Master Contract Retro Life.md` (D1) +
`../flows/Master Contract Retro Life.md` (D2 Tahap 5) + `../flows/_SUMMARY-treaty-master.md`.
**66 file — modul terkecil di korpus**. Modul **tanpa rule `Flow`** (OQ-005).
Audit: `find "Master Contract Retro Life" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` **Master kontrak retrosesi lini life** — empat layar grid di atas satu keluarga
tabel `POOLDATA.*_LIFE`. Class terbanyak **`@BASECLASS`** (38 rule), lalu
`ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` (6), `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` (5),
`DATA-PORTAL` (4).

`[terverifikasi]` **Bukan proses berjenjang**: nol rujukan `StatusAkseptasi`, tanpa `Akseptasi_DT`,
tanpa tombol Submit/Decline. Hanya **1 DataTransform** di seluruh modul.
Audit: `grep -rl "StatusAkseptasi" "Master Contract Retro Life" --include="*.xml" | wc -l` → **0**.

`[terverifikasi]` **Modul paling terisolasi di korpus**: tanpa `ConnectREST`, tanpa `When`, tanpa
`SystemSettings`, tanpa db-link, tanpa kolom JSON, tanpa jejak guard identitas.

## 2. Proses / fitur utama

Rincian: **`../flows/Master Contract Retro Life.md`** (23 rule ditelusur).

`[terverifikasi]` **Titik masuk**: `Harness/InboxRetroLifeReinsurersList.xml` → `RULE-HTML-HARNESS`
/ `DATA-PORTAL` / `INBOXRETROLIFEREINSURERSLIST`, 530.052 byte.

`[terverifikasi]` **Empat Harness = empat grid master** — rasio Harness tertinggi di korpus:

| Harness | `pxInsName` | Ukuran |
| --- | --- | ---: |
| `InboxRetroLifeReinsurersList.xml` | `DATA-PORTAL!INBOXRETROLIFEREINSURERSLIST` | 530.052 |
| `InboxRetroLimitReinsurers.xml` | `DATA-PORTAL!INBOXRETROLIMITREINSURERS` | 575.179 |
| `InboxBusinessLifeReinsurers.xml` | `DATA-PORTAL!INBOXBUSINESSLIFEREINSURERS` | 486.911 |
| `InboxSecurityReinsurerLife.xml` | `DATA-PORTAL!INBOXSECURITYREINSURERLIFE` | 444.874 |

`[terverifikasi]` Hanya **dua FlowAction**, keduanya tampilan: `ViewRate.xml`, `ViewRateTable.xml`.

`[terverifikasi]` Pola CRUD simetris (27 Activity): 5 `New…`, 8 `Set…`, 6 `Save…`, 4 `Delete…`,
1 `CancelActivityTreatyContract`, ditambah `CountingPercentShare_Act`, `TreatyLimit_TypeProtect`,
`SetErrorMessageReinsurer`. Tahapan: buka grid → tambah/ubah/hapus baris → simpan ke master
`*_LIFE` lewat Connect-SQL.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 27 Activity, 12 ReportDefinition, 12 Connect-SQL,
8 Section, 4 Harness, 2 FlowAction, 1 DataTransform. **Nol `Flow`, nol `When`, nol `ConnectREST`.**

`[terverifikasi]` **Hanya 4 objek Oracle, seluruhnya bersufiks `_LIFE` di skema `POOLDATA`**:

| Objek | Rule perujuk |
| --- | ---: |
| `POOLDATA.TREATYREINSURER_LIFE` | 2 |
| `POOLDATA.TREATYCONTRACT_LIFE` | 2 |
| `POOLDATA.TREATYBUSINESS_LIFE` | 2 |
| `POOLDATA.TREATYSECURITYREINSURER_LIFE` | 1 |

Audit:
```
awk -F'\t' '$1 ~ "^Master Contract Retro Life/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
 | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

`[terverifikasi]` **Tidak ada** `M_TREATY_IN`, objek facultative, maupun objek outward.

## 4. Integrasi eksternal

`[terverifikasi]` **Tidak ada** — modul ini **tidak punya rule `RULE-CONNECT-REST` sama sekali**,
satu-satunya di korpus yang demikian. Tidak ada Google Storage, tidak ada `M_LINK_SERVICE`, tidak
ada Arasapas/Kasir/Konversi/Gemini AI, tidak ada db-link.

**Seluruh interaksinya dengan dunia luar adalah 12 Connect-SQL ke `POOLDATA`.**

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Claim Life** | modul itu membaca `RETROCESSIONLIFE`, `TREATYYEAR_LIFE` | `[terverifikasi]` |
| **Master Product Name Life** | `POOLDATA.TREATYCONTRACT_LIFE`, `TREATYYEAR_LIFE` disentuh keduanya | `[terverifikasi]` |
| **Komite Claim Life** | class `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` dirujuk modul itu | `[terverifikasi]` |

`[terverifikasi]` **Entri register OQ-011 paling sedikit di korpus: 5**
(`grep -cF "Master Contract Retro Life/" ../inventory/_oq011-konflik-isi.md`).

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Lima stored procedure, seluruhnya penulis dan seluruhnya eksklusif modul ini**,
isinya tidak ada di korpus (OQ-002):

| Rule Connect-SQL | Procedure |
| --- | --- |
| `RDBList/SaveMasterTreatyContract_Life_SQL.xml` | `POOLDATA.INSERTTREATYCONTRACT_LIFE` |
| `RDBList/SaveMasterTreatyYear_Life_SQL.xml` | `POOLDATA.INSERTTREATYYEAR_LIFE` |
| `RDBList/SaveMasterTreatyReinsurer_Life_SQL.xml` | `POOLDATA.INSERTREINSURER_LIFE` |
| `RDBList/SaveMasterTreatyBusiness_Life_SQL.xml` | `POOLDATA.INSERTBUSINESS_LIFE` |
| `RDBList/SaveMasterTreatySecurityReinsurer_Life_SQL.xml` | `POOLDATA.INSERTSECURITYREINSURER_LIFE` |

**Seluruh jalur tulis melewati stored procedure** — aturan penyimpanan kontrak retro life
(validasi, versi, kunci, kaskade) **berada di sisi database dan tidak ada di korpus**.

### 6.1 `REINSTYPEID` — tidak ada satu pun nilai literal

`[terverifikasi]` `REINSTYPEID` muncul **129 kali** di modul ini
(`grep -rho "REINSTYPEID" "Master Contract Retro Life" --include="*.xml" | wc -l`), tetapi
**selalu diisi dari slot generik** dan **tidak pernah sebagai konstanta**:

```
DECLARE
  … REINSTYPEID VARCHAR2(32767);
  REINSTYPEID := {TempInputData.CARI3};
  REINSTYPEID := {TempInputData.CARI4};
  REINSTYPEID := {TempInputData.CARI5};
  … ( p_REINSTYPEID, p_REINSTYPENAME, p_TREATYSTARTDATE, … )
  … ( p_REINSTYPEID, p_REINSTYPENAME, p_REINSURERID, p_REINSURERNAME, … )
  … ( p_REINSTYPEID, p_REINSTYPENAME, p_BIZCODE, p_BIZNAME, p_RIR… )
```

**Daftar nilai `REINSTYPEID` tidak dapat dinyatakan** → OQ-020, **OQ-057**.

`[terverifikasi]` Pola slot generik `CARI<n>` juga dipakai di `NB Treaty In`
(`InputData.CARI20`, `CARI21`, `CARI3`), `Treaty In Adjustment` (`InputData.CARI1`),
`Treaty In` (`StatusDoc.CARI30`), dan `NB FacIn` (`FlagBanding.CARI11`) → **OQ-059**.

### 6.2 Rumus yang belum dibaca

`[terverifikasi]` Tiga activity perhitungan, **seluruhnya di class `@BASECLASS`** dan **belum
dibaca** — batas **cakupan telusur**:

| Activity | `pxInsName` | Ukuran | Langkah |
| --- | --- | ---: | ---: |
| `SetValueRetroLimit_TreatyYearLife.xml` | `@BASECLASS!SETVALUERETROLIMIT_TREATYYEARLIFE` | 63.678 | 2 |
| `CountingPercentShare_Act.xml` | `@BASECLASS!COUNTINGPERCENTSHARE_ACT` | 59.429 | 4 |
| `TreatyLimit_TypeProtect.xml` | `@BASECLASS!TREATYLIMIT_TYPEPROTECT` | 54.893 | 4 |

**Rumus pembagian share retro dan jenis proteksi limit tidak dinyatakan.**
Isi 4 Harness (444–575 KB) juga hanya di-grep, tidak dibaca.

### 6.3 Berkas non-rule di dalam ekspor

`[terverifikasi]` Folder **`Claude outputs`** berisi `Struktur_GridRetrocessionLife.xlsx`
(11.944 B) dan **`perubahan_skill.diff`** (10.278 B) — berkas non-Pega di dalam korpus READ-ONLY.
**Tidak dibaca** → **OQ-054**.

`[terverifikasi]` **Tidak ada kolom JSON** di modul ini
(`grep -rho "JSONDATA\|DATA_JSON" … | wc -l` → **0**) — satu-satunya modul Tahap 5 yang bebas dari
OQ-012. `IsPEGAPROD` **tidak ada** — OQ-029 tidak berlaku.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **0 file**,
`OperatorID.pyTelephone` **0**, `pxCreateOperator` sebagai guard **0** — **modul tanpa jejak guard
identitas sama sekali**, satu-satunya di Tahap 5.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-005**, **OQ-009**, **OQ-020**, **OQ-054**, **OQ-057**, **OQ-059**.
