# Telusur — Master Product Name Life (modul **tanpa** rule `Flow`)

STEP D2, Tahap 5 konteks #19. Ditelusur 2026-09-13. Konvensi: **`_METHOD.md`** +
**`_METHOD-noflow.md`**. **114 file** — modul terkecil kedua di korpus
(`find "Master Product Name Life" -name '*.xml' | wc -l`).

**Titik masuk `[terverifikasi]`:** `Harness/InwardProductName.xml`
→ `RULE-HTML-HARNESS` / **`ASM-FW-GISFW-INT-PRODUCT_LIFE`** / `INWARDPRODUCTNAME`, 540.636 byte.
**Satu-satunya Harness** di modul ini, dan satu-satunya titik masuk Tahap 5 yang **tidak** berclass
`DATA-PORTAL` — ia menempel langsung pada class entitas data produk life.

---

## 1. Alur — editor master produk life

### 1.1 Titik masuk dan aksi

`[terverifikasi]` Sebelas `FlowAction`, seluruhnya bercorak **pemilih (picker) dan konfirmasi**,
bukan langkah proses:

| FlowAction | Corak |
| --- | --- |
| `ChooseCauseOfLoss`, `ChooseCeding`, `ChooseCurrency`, `ChoosePolicyHolder`, `ChooseRIRate`, `ChooseRIRisk`, `ChooseSOB` | **7 pemilih master** |
| `SaveProductName_Confirm`, `EditProductName_Confirm` | 2 konfirmasi simpan/ubah |
| `ViewRate` | tampilan rate |
| `ProductNameAttachContent` | lampiran |

`[terverifikasi]` **Tidak ada tombol Submit/Akseptasi/Decline, tidak ada `StatusAkseptasi`,
tidak ada `Position`, tidak ada `Akseptasi_DT`.** Perintah audit:
```
grep -rl "StatusAkseptasi" "Master Product Name Life" --include="*.xml" | wc -l   # 0
ls "Master Product Name Life/DataTransform/"                                      # 10, tanpa Akseptasi_DT
```

**Bentuk modul `[terverifikasi]`:** editor master **satu entitas** (nama produk life) dengan
7 pemilih master pendukung dan 2 konfirmasi — **bukan proses berjenjang**. Tahapan yang dapat
dinyatakan: buka layar produk → pilih nilai master (ceding, currency, policy holder, SOB, RI rate,
RI risk, cause of loss) → konfirmasi → simpan.

### 1.2 Rantai simpan `[terverifikasi]`

| Activity | `pxInsName` | Ukuran | Langkah |
| --- | --- | ---: | ---: |
| `SaveProductName_Act.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE!SAVEPRODUCTNAME_ACT` | 165.682 | 16 |
| `SaveInwardProductName_Act.xml` | `…!SAVEINWARDPRODUCTNAME_ACT` | 64.796 | 6 |
| `SetProductNameInward.xml` | `…!SETPRODUCTNAMEINWARD` | 125.938 | 3 |
| `NewProductLife.xml` | `…!NEWPRODUCTLIFE` | 40.598 | 3 |
| `GetReinsTypeOR_Life.xml` | `…!GETREINSTYPEOR_LIFE` | 66.519 | 4 |

`[terverifikasi]` Dua jalur simpan yang **berbeda nama**: `SaveProductName_Act` (produk) dan
`SaveInwardProductName_Act` (produk **inward**), ditambah `SetProductNameInward`. Pembeda
"inward" muncul di nama rule dan di class `ASM-FW-GISFW-Int-PRODUCTINWARD_LIFE` yang dirujuk
`GetReinsTypeOR_Life`.

---

## 2. **UJI: apakah modul ini membawa salinan subtree Treaty In?**

### 2.1 Pengukuran `[terverifikasi]`

```
for rel in $(cd "Master Product Name Life" && find . -name '*.xml' | sed 's|^\./||'); do
  [ -f "Treaty In/$rel" ] && { bandingkan hash ternormalisasi 18 tag; }
done
```

| Ukuran | Jumlah |
| --- | ---: |
| File `Master Product Name Life` yang **bernama sama** dengan file di `Treaty In` | **35** dari 114 (30,7 %) |
| — **identik** isinya (18 tag) | **16** |
| — berbeda (18 tag) | **19** |
| — berbeda setelah normalisasi **21 tag** | **12** |
| — **konflik semu** (runtuh oleh 3 tag metadata) | **7** |

`[terverifikasi]` Ketujuh konflik semu seluruhnya keluarga penyimpanan berkas:
`RDBList/DeleteStorage_SQL.xml`, `GetAppName_SQL.xml`, `GetLinkStorage_SQL.xml`,
`GetTokenStorage_SQL.xml`, `Insert_T_Storage_SQL.xml`, `Update_T_Storage_SQL.xml`,
**`SystemSettings/LinkService.xml`** — pola yang sama persis dengan temuan Tahap 4
(`flows/_SUMMARY-facultative.md` §8, entri #324/#367/#369/#371/#376/#377).

### 2.2 Jawaban: **ya, sebagian — tetapi bukan subtree proses**

`[terverifikasi]` Yang dibawa adalah **infrastruktur bersama**, bukan mesin proses treaty:
lampiran/penyimpanan berkas (`LoadAttachment`, `DownloadAll_Act`, `*Storage*`, `GetLinkService`,
`ServiceGoogle`, `LinkService`, `GetMimeType`), komentar (`AddCommentList_Act`), dan
**tiga rule treaty inward** yang nyata:

| Rule | Identitas | Sifat |
| --- | --- | --- |
| `Activity/SetTreatyIn_Act.xml` | `DATA-PORTAL!SETTREATYIN_ACT` | 159.737 B, **15 langkah — urutan langkah IDENTIK dengan `Treaty In`**, isi berbeda 47 byte |
| `RDBList/SaveTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN!ASM!SAVETREATYIN` | **identik lintas 6 modul** (`e371c194`) |
| `RDBList/BrowseTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN!ASM!BROWSETREATYIN` | **identik lintas 6 modul** (`196d49b5`) |

Ditambah `Activity/TreatyInInputVis.xml`, `Activity/TreatyInitAttach.xml`,
`Activity/TreatyInDownloadAll.xml`, `Activity/TreatySetReinstatement.xml`,
`Activity/SetReinstatementPct.xml`, `Activity/CheckDuplicateOffer.xml`,
`Activity/ConvertHistoryDate.xml`, `DataTransform/TreatyInIDSetPyPortal.xml`,
`Activity/SetCategoryAttachTreatyin.xml`.

`[terverifikasi]` **`Akseptasi_DT`, `TreatyInSubmit`, `TreatyInActionButtons`,
`TreatyInDeclineConfirmation_postact` — tidak ada** di modul ini. Jadi **tangga persetujuan treaty
tidak ikut tersalin**; yang tersalin adalah rule pembentuk dan penyimpan data treaty.

### 2.3 **Penulis kedua `M_TREATY_IN`** — terbukti `[terverifikasi]`

```
awk -F'\t' '$1 ~ "^Master Product Name Life/" && $8 ~ /M_TREATY_IN|PEGA_TREATY_IN/{print $5"  "$1"  "$8}' all-rules.tsv
```
→
```
BROWSETREATYIN  …/RDBList/BrowseTreatyIn.xml  sqlKind=QUERY;sqlOps=SELECT;
                tables=POOLDATA.M_TREATY_IN,POOLDATA.M_TREATY_IN_EDM
SAVETREATYIN    …/RDBList/SaveTreatyIn.xml    sqlKind=PLSQL;procs=POOLDATA.PEGA_TREATY_IN
```

`[terverifikasi]` **Ya — modul master produk life membawa jalur tulis ke master treaty inward**,
lewat procedure `POOLDATA.PEGA_TREATY_IN` yang sama dengan yang dipakai modul `Treaty In`.
Pemanggilnya di modul ini adalah `Activity/SetTreatyIn_Act.xml`
(`grep -rl "SaveTreatyIn" "Master Product Name Life" --include="*.xml"` → `Activity/SetTreatyIn_Act.xml`,
`RDBList/SaveTreatyIn.xml`).

**Nuansa yang penting dan tidak boleh dilewat:** rule tulisnya **bukan salinan yang menyimpang** —
`ASM!SAVETREATYIN` identik di 6 modul dan **tidak terdaftar di register OQ-011**. Jadi ini bukan
"penulis kedua dengan logika berbeda", melainkan **rule tulis bersama yang ikut terekspor ke modul
master life**. Apakah jalur itu benar-benar dieksekusi dari modul ini **tidak dapat dipastikan dari
korpus** — tidak ada graf yang menunjukkannya. → **OQ-056**.

---

## 3. Status / state dan kode

`[terverifikasi]` **Tidak ada properti status proses.** Kode yang ada bersifat master/identitas:

| Properti | Nilai literal / sumber | Arti |
| --- | --- | --- |
| `ReinsTypeID` | dari `BrowseReinstypeOR_SQL` (`SELECT tc.*, ty.* …`) | **belum terverifikasi** |
| kode **`OR`** (pada nama `GetReinsTypeOR_Life`, `BrowseReinstypeOR_SQL`) | muncul hanya sebagai bagian nama rule; **tidak ditemukan sebagai nilai literal** | **kepanjangan & arti belum terverifikasi** |

Perintah audit:
```
grep -rhoE "\"OR[A-Z0-9]*\"|'OR[A-Z0-9]*'" "Master Product Name Life" --include="*.xml" | sort | uniq -c
```
→ kosong. **Kode `OR` tidak dapat dikutip sebagai nilai** — hanya nama rule yang memuatnya,
dan `_METHOD.md` §2.2 melarang memakai nama sebagai bukti. → **OQ-057**.

`[terverifikasi]` `GetReinsTypeOR_Life` merujuk **tiga class berbeda**:
`ASM-FW-GISFW-Int-PRODUCT_LIFE`, `ASM-FW-GISFW-Int-PRODUCTINWARD_LIFE`,
`ASM-FW-GISFW-Int-TREATYBUSINESS_LIFE` — jadi ia menjembatani produk life, produk inward life, dan
business treaty life.

---

## 4. Objek Oracle yang disentuh

```
awk -F'\t' '$1 ~ "^Master Product Name Life/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
  | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

| Objek | Rule perujuk |
| --- | ---: |
| `T_STORAGE_IMAGE` | 3 |
| `M_ATTACHMENTTREATY_2`, `M_ATTACHMENTPRODUCTNAME` | 2 masing-masing |
| `POOLDATA.T_FOLDER_IMAGE`, `POOLDATA.TREATYYEAR_LIFE`, `POOLDATA.TREATYCONTRACT_LIFE`, **`POOLDATA.M_TREATY_IN`**, **`POOLDATA.M_TREATY_IN_EDM`**, `POOLDATA.M_PRODUCT_LIFE`, `M_PRODUCT_LIFE`, **`DATAPEGA.PC_ASM_FW_GCNMFW_WORK`**, `CATEGORY_ATTACH_REAS` | 1 masing-masing |

`[terverifikasi]` **`DATAPEGA.PC_ASM_FW_GCNMFW_WORK`** — tabel **internal Pega** (skema `DATAPEGA`,
tabel work object class `ASM-FW-GCNMFW-WORK`) diakses langsung lewat SQL, bukan lewat API Pega.
Skema `DATAPEGA` adalah salah satu dari 11 skema OQ-016. Perintah audit:
```
awk -F'\t' '$8 ~ /DATAPEGA/{print $1"  "$5}' all-rules.tsv
```
→ **OQ-058**.

`[terverifikasi]` Modul ini menyentuh **tiga keluarga master sekaligus**: produk life
(`M_PRODUCT_LIFE`), kontrak treaty life (`TREATYCONTRACT_LIFE`, `TREATYYEAR_LIFE`), dan master
treaty inward umum (`M_TREATY_IN`, `M_TREATY_IN_EDM`).

---

## 5. Integrasi eksternal

`[terverifikasi]` Satu `RULE-CONNECT-REST`: `ConnectREST/ServiceGoogle.xml`
(`SETTING` / `LinkService!LinkService`, **tanpa URL literal endpoint**).
Empat activity Google Storage: `InsertGoogleStorage_Act`, `GetUrlGoogleStorage_Act`,
`DeleteGoogleStorage_Act`, `GetLinkService` — alamat dari tabel `M_LINK_SERVICE` (OQ-047).

**Tidak ada `GeminiAIGoogle_Act`**, tidak ada Arasapas/Kasir/Konversi.

---

## 6. Batas pengetahuan

### 6.1 Stored procedure — 5 dipanggil `[terverifikasi]`

```
POOLDATA.PEGA_M_PRODUCT_LIFE        POOLDATA.PEGA_M_PRODUCT_INWARD_LIFE
POOLDATA.PEGA_TREATY_IN             POOLDATA.GET_TOKEN_STORAGE
DBMS_LOB.CREATETEMPORARY
```
Isinya tidak ada di korpus (OQ-002). `PEGA_M_PRODUCT_INWARD_LIFE` **hanya muncul di modul ini**
`[terverifikasi]`.

### 6.2 Yang belum dibaca

`[terverifikasi]` Activity terbesar yang **belum habis dibaca**:
`SaveProductName_Act.xml` (165.682 B, 16 langkah), `SetTreatyIn_Act.xml` (159.737 B, 15 langkah),
`SetProductNameInward.xml` (125.938 B), `GetReinsTypeOR_Life.xml` (66.519 B),
`SaveInwardProductName_Act.xml` (64.796 B), serta Harness `InwardProductName.xml` (540.636 B).

Selisih isi 12 rule terhadap `Treaty In` (§2.1) **tidak dikarakterisasi** — batas cakupan telusur.

### 6.3 Kolom JSON

`[terverifikasi]` `BrowseTreatyIn` menyeleksi `JSONDATA` dari `POOLDATA.M_TREATY_IN` dan
`POOLDATA.M_TREATY_IN_EDM` → OQ-012. Strukturnya tidak ada di korpus.

---

## 7. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Master Product Name Life/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `INWARDPRODUCTNAME` / `RULE-HTML-HARNESS` | `Harness/InwardProductName.xml` | tidak |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/SaveProductName_Act.xml` | tidak |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEINWARDPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/SaveInwardProductName_Act.xml` | tidak |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SETPRODUCTNAMEINWARD` / `RULE-OBJ-ACTIVITY` | `Activity/SetProductNameInward.xml` | tidak |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `NEWPRODUCTLIFE` / `RULE-OBJ-ACTIVITY` | `Activity/NewProductLife.xml` | tidak |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `GETREINSTYPEOR_LIFE` / `RULE-OBJ-ACTIVITY` | `Activity/GetReinsTypeOR_Life.xml` | tidak |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!BROWSEREINSTYPEOR_SQL` / `RULE-CONNECT-SQL` | `RDBList/BrowseReinstypeOR_SQL.xml` | tidak |
| **`DATA-PORTAL` / `SETTREATYIN_ACT` / `RULE-OBJ-ACTIVITY`** | `Activity/SetTreatyIn_Act.xml` | **YA** — beda dgn `Treaty In` (langkah identik) |
| **`ASM-FW-GISFW-INT-TREATY_IN` / `ASM!SAVETREATYIN` / `RULE-CONNECT-SQL`** | `RDBList/SaveTreatyIn.xml` | tidak — **identik 6 modul** (§2.3) |
| **`ASM-FW-GISFW-INT-TREATY_IN` / `ASM!BROWSETREATYIN` / `RULE-CONNECT-SQL`** | `RDBList/BrowseTreatyIn.xml` | tidak — identik 6 modul |
| `…` / `TREATYININPUTVIS`, `TREATYINITATTACH`, `TREATYINDOWNLOADALL`, `TREATYSETREINSTATEMENT`, `SETREINSTATEMENTPCT`, `CHECKDUPLICATEOFFER`, `CONVERTHISTORYDATE` / `RULE-OBJ-ACTIVITY` | `Activity/` | sebagian |
| `…` / `TREATYINIDSETPYPORTAL` / `RULE-OBJ-MODEL` | `DataTransform/TreatyInIDSetPyPortal.xml` | **YA** |
| 11 `RULE-OBJ-FLOWACTION` (§1.1) | `FlowAction/` | sebagian |
| `DATA-PORTAL` / `SERVICEGOOGLE` / `RULE-CONNECT-REST` | `ConnectREST/ServiceGoogle.xml` | **YA** |
| `LINKSERVICE` / `LINKSERVICE` / `RULE-ADMIN-SYSTEM-SETTINGS` | `SystemSettings/LinkService.xml` | **YA — konflik semu** (§2.1) |

**22 rule ditelusur** (1 Harness; 8 Activity; 4 Connect-SQL; 1 DataTransform; 1 ConnectREST;
1 SystemSettings; 11 FlowAction dihitung sebagai satu kelompok).

Rule **dirujuk tapi belum ditelusur**: 22 RDBList selain yang disebut, 17 ReportDefinition,
12 Section, sisa 29 Activity.

---

## 8. Pertanyaan terbuka

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-056** | Modul master produk life membawa jalur tulis ke master **treaty inward** (`ASM!SAVETREATYIN` → `POOLDATA.PEGA_TREATY_IN`); apakah jalur itu dieksekusi dari sini, dan mengapa? | Product+UW + DBA |
| **OQ-057** | Kode **`OR`** (`GetReinsTypeOR_Life`, `BrowseReinstypeOR_SQL`) hanya ada sebagai nama rule, **tidak pernah muncul sebagai nilai literal** — artinya tidak dapat dibuktikan dari korpus | Product+UW |
| **OQ-058** | Tabel internal Pega `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` diakses langsung lewat SQL, bukan lewat API Pega | DBA + Arsitektur Pega |

OQ dikuatkan: **OQ-002** (§6.1), **OQ-005** (§1 — titik masuk), **OQ-009** (Harness di class
entitas, bukan `DATA-PORTAL`), **OQ-011** (§2.1 — 12 tetap beda, 7 konflik semu termasuk
`LINKSERVICE`), **OQ-012** (§6.3), **OQ-016** (§4 — skema `DATAPEGA`), **OQ-018/OQ-047** (§5),
**OQ-020** (§3).

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **3 file**,
`OperatorID.pyTelephone` **0**, `pxCreateOperator` sebagai guard **0**. `IsPEGAPROD` **tidak ada**.
