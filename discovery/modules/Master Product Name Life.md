# Modul — Master Product Name Life

STEP D3. Sintesis dari `../inventory/Master Product Name Life.md` (D1) +
`../flows/Master Product Name Life.md` (D2 Tahap 5) + `../flows/_SUMMARY-treaty-master.md`.
**114 file**. Modul **tanpa rule `Flow`** (OQ-005).
Audit: `find "Master Product Name Life" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` **Master nama produk lini life** — editor satu entitas dengan pemilih master
pendukung. Class kerja `ASM-FW-GISFW-INT-PRODUCT_LIFE` (55 rule).

`[terverifikasi]` **Titik masuk**: `Harness/InwardProductName.xml` → `RULE-HTML-HARNESS` /
**`ASM-FW-GISFW-INT-PRODUCT_LIFE`** / `INWARDPRODUCTNAME`, 540.636 byte.
**Satu-satunya Harness** di modul ini, dan **satu-satunya titik masuk Tahap 5 yang tidak berclass
`DATA-PORTAL`** — ia menempel langsung pada class entitas data produk life.

`[terverifikasi]` **Bukan proses berjenjang**: nol rujukan `StatusAkseptasi`, tanpa `Akseptasi_DT`,
tanpa tombol Submit/Akseptasi/Decline.
Audit: `grep -rl "StatusAkseptasi" "Master Product Name Life" --include="*.xml" | wc -l` → **0**.

## 2. Proses / fitur utama

Rincian: **`../flows/Master Product Name Life.md`** (22 rule ditelusur).

`[terverifikasi]` Sebelas `FlowAction`, seluruhnya bercorak **pemilih dan konfirmasi**, bukan
langkah proses:

| Kelompok | FlowAction |
| --- | --- |
| **7 pemilih master** | `ChooseCauseOfLoss`, `ChooseCeding`, `ChooseCurrency`, `ChoosePolicyHolder`, `ChooseRIRate`, `ChooseRIRisk`, `ChooseSOB` |
| 2 konfirmasi | `SaveProductName_Confirm`, `EditProductName_Confirm` |
| Lainnya | `ViewRate`, `ProductNameAttachContent` |

`[terverifikasi]` Tahapan: buka layar produk → pilih nilai master → konfirmasi → simpan.

`[terverifikasi]` **Dua jalur simpan berbeda nama**: `SaveProductName_Act`
(`ASM-FW-GISFW-INT-PRODUCT_LIFE!SAVEPRODUCTNAME_ACT`, 165.682 B, 16 langkah) dan
`SaveInwardProductName_Act` (64.796 B, 6 langkah), ditambah `SetProductNameInward` (125.938 B).
Pembeda "inward" muncul di nama rule dan di class `ASM-FW-GISFW-Int-PRODUCTINWARD_LIFE`.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 37 Activity, 22 Connect-SQL, 17 ReportDefinition,
12 Section, 11 FlowAction, 10 DataTransform, 1 When, 1 Harness, 1 DecisionTable, 1 ConnectREST,
1 SystemSettings. **Nol rule `Flow`.**

| Objek | Rule perujuk |
| --- | ---: |
| `T_STORAGE_IMAGE` | 3 |
| `M_ATTACHMENTTREATY_2`, `M_ATTACHMENTPRODUCTNAME` | 2 masing-masing |
| `POOLDATA.T_FOLDER_IMAGE`, `POOLDATA.TREATYYEAR_LIFE`, `POOLDATA.TREATYCONTRACT_LIFE`, **`POOLDATA.M_TREATY_IN`**, **`POOLDATA.M_TREATY_IN_EDM`**, `POOLDATA.M_PRODUCT_LIFE`, `M_PRODUCT_LIFE`, **`DATAPEGA.PC_ASM_FW_GCNMFW_WORK`**, `CATEGORY_ATTACH_REAS` | 1 masing-masing |

`[terverifikasi]` **Menyentuh tiga keluarga master sekaligus**: produk life (`M_PRODUCT_LIFE`),
kontrak treaty life (`TREATYCONTRACT_LIFE`, `TREATYYEAR_LIFE`), dan **master treaty inward umum**
(`M_TREATY_IN`, `M_TREATY_IN_EDM`).

`[terverifikasi]` **`DATAPEGA.PC_ASM_FW_GCNMFW_WORK`** — tabel **internal Pega** (skema `DATAPEGA`,
tabel work object class `ASM-FW-GCNMFW-WORK`) diakses **langsung lewat SQL**, bukan lewat API Pega
→ **OQ-058**.
Audit: `awk -F'\t' '$8 ~ /DATAPEGA/{print $1"  "$5}' all-rules.tsv`

## 4. Integrasi eksternal

`[terverifikasi]` Satu `RULE-CONNECT-REST`: `ServiceGoogle` (`SETTING` → `LinkService!LinkService`,
tanpa URL literal). Empat activity Google Storage: `InsertGoogleStorage_Act`,
`GetUrlGoogleStorage_Act`, `DeleteGoogleStorage_Act`, `GetLinkService` — alamat dari tabel
`M_LINK_SERVICE` → OQ-047.

**Tidak ada `GeminiAIGoogle_Act`, Arasapas, Kasir, maupun Konversi.**

## 5. Ketergantungan ke modul lain — **penulis master treaty inward kedua**

`[terverifikasi]` **35 dari 114 file (30,7 %) bernama sama** dengan file di `Treaty In`:
16 identik, 19 berbeda (12 setelah normalisasi 21 tag; **7 konflik semu**, termasuk
`SystemSettings/LinkService.xml`).

Yang tersalin adalah **infrastruktur bersama** (lampiran, penyimpanan berkas, komentar)
**ditambah tiga rule treaty inward nyata**:

| Rule | Sifat |
| --- | --- |
| `RDBList/SaveTreatyIn.xml` (`ASM-FW-GISFW-INT-TREATY_IN!ASM!SAVETREATYIN`) | memanggil **`POOLDATA.PEGA_TREATY_IN`**; **identik di 6 modul** (hash `e371c194`) |
| `RDBList/BrowseTreatyIn.xml` (`…!ASM!BROWSETREATYIN`) | membaca `M_TREATY_IN`, `M_TREATY_IN_EDM`; identik di 6 modul (`196d49b5`) |
| `Activity/SetTreatyIn_Act.xml` (`DATA-PORTAL!SETTREATYIN_ACT`) | 159.737 B, **15 langkah — urutan identik** dengan versi `Treaty In`, isi berbeda 47 byte |

`[terverifikasi]` **Nuansa yang tidak boleh dilewat:** rule tulisnya **bukan salinan yang
menyimpang** — identik di 6 modul dan **tidak terdaftar di register OQ-011**. Jadi ini **rule tulis
bersama yang ikut terekspor**, bukan penulis kedua dengan logika berbeda.

`[terverifikasi]` **`Akseptasi_DT`, `TreatyInSubmit`, `TreatyInActionButtons` TIDAK ada** di modul
ini — **tangga persetujuan treaty tidak ikut tersalin**.

`[pertanyaan terbuka]` Apakah jalur tulis `M_TREATY_IN` benar-benar dieksekusi dari modul ini, dan
dalam keadaan apa? **Tidak dapat dipastikan dari korpus** — tidak ada graf yang menunjukkannya
→ **OQ-056**.

Ketergantungan lain:

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Claim Life / PremiumList Life / Endorsement Life** | `M_PRODUCT_LIFE`, `PRODUCT_LIFE` dibaca modul-modul itu | `[terverifikasi]` |
| **Master Contract Retro Life** | `POOLDATA.TREATYCONTRACT_LIFE`, `TREATYYEAR_LIFE` disentuh keduanya | `[terverifikasi]` |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Lima stored procedure**, isinya tidak ada di korpus (OQ-002):
`POOLDATA.PEGA_M_PRODUCT_LIFE`, **`POOLDATA.PEGA_M_PRODUCT_INWARD_LIFE`** (hanya di modul ini),
**`POOLDATA.PEGA_TREATY_IN`**, `POOLDATA.GET_TOKEN_STORAGE`, `DBMS_LOB.CREATETEMPORARY`.

`[terverifikasi]` **Kode `OR` tidak dapat dibuktikan dari korpus.** Ia muncul hanya sebagai bagian
nama rule (`Activity/GetReinsTypeOR_Life.xml`, `RDBList/BrowseReinstypeOR_SQL.xml`) dan **tidak
pernah sebagai nilai literal**:
```
grep -rhoE "\"OR[A-Z0-9]*\"|'OR[A-Z0-9]*'" "Master Product Name Life" --include="*.xml" | sort | uniq -c
```
→ **kosong**. `../flows/_METHOD.md` §2.2 melarang memakai nama rule sebagai bukti → **OQ-057**.

`[terverifikasi]` `GetReinsTypeOR_Life` merujuk **tiga class**: `ASM-FW-GISFW-Int-PRODUCT_LIFE`,
`ASM-FW-GISFW-Int-PRODUCTINWARD_LIFE`, `ASM-FW-GISFW-Int-TREATYBUSINESS_LIFE` — menjembatani produk
life, produk inward life, dan business treaty life.

`[terverifikasi]` Activity besar **belum habis dibaca**: `SaveProductName_Act.xml` (165.682 B),
`SetTreatyIn_Act.xml` (159.737 B), `SetProductNameInward.xml` (125.938 B),
`GetReinsTypeOR_Life.xml` (66.519 B), Harness `InwardProductName.xml` (540.636 B) — batas
**cakupan telusur**. Selisih isi 12 rule terhadap `Treaty In` **tidak dikarakterisasi**.

`[terverifikasi]` `BrowseTreatyIn` menyeleksi kolom `JSONDATA` → OQ-012.

`[terverifikasi]` `IsPEGAPROD` **tidak ada** di modul ini — OQ-029 tidak berlaku.
**Tidak ada objek db-link** — OQ-017 tidak berlaku.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **3 file**,
`OperatorID.pyTelephone` **0**.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-005**, **OQ-009**, **OQ-011**, **OQ-012**, **OQ-016**, **OQ-018**, **OQ-021**,
**OQ-047**, **OQ-056**, **OQ-057**, **OQ-058**.
