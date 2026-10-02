# Modul — Treaty In Adjustment

STEP D3. Sintesis dari `../inventory/Treaty In Adjustment.md` (D1) +
`../flows/Treaty In Adjustment.md` (D2 Tahap 5) + `../flows/_SUMMARY-treaty-master.md`.
**379 file**. Modul **tanpa rule `Flow`** (OQ-005).
Audit: `find "Treaty In Adjustment" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` **Revisi dan adjustment atas master treaty inward**. Class terbanyak
**`DATA-PORTAL`** (171 rule), lalu `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` (36),
`ASM-FW-GISFW-INT-TREATY_IN` (29), `ASM-FW-GISFW-DATA-TREATYINLIMITS` (29).

`[terverifikasi]` **Satu modul menangani DUA proses berbeda** — revisi dan adjustment — dibedakan
oleh parameter `Param.type`, bukan oleh rule terpisah (§2.3).

## 2. Proses / fitur utama

Rincian: **`../flows/Treaty In Adjustment.md`** (31 rule ditelusur).

`[terverifikasi]` **Titik masuk**: `Harness/InputTreatyInAdjustment.xml` → `RULE-HTML-HARNESS` /
`DATA-PORTAL` / `INPUTTREATYINADJUSTMENT`, 943.918 byte. Modul ini juga membawa salinan
`Harness/InputTreatyInOffer.xml` (8.226.491 B).

### 2.1 Berapa banyak yang benar-benar berbeda dari `Treaty In`

`[terverifikasi]`

| Ukuran | Jumlah |
| --- | ---: |
| File bernama sama | **323** |
| — **identik** (normalisasi 18 tag) | **277** (85,8 %) |
| — berbeda | 46 |
| — berbeda setelah normalisasi **21 tag** | **43** |
| — **konflik semu** (3 tag metadata) | **3** — keluarga `T_STORAGE_IMAGE` |
| Hanya di `Treaty In` | 6 |
| Hanya di `Treaty In Adjustment` | **56** |

Audit:
```
(cd "Treaty In" && find . -name '*.xml' | sed 's|^\./||' | sort) > ti.txt
(cd "Treaty In Adjustment" && find . -name '*.xml' | sed 's|^\./||' | sort) > tia.txt
comm -12 ti.txt tia.txt | wc -l   # 323
comm -13 ti.txt tia.txt | wc -l   # 56
```

`[terverifikasi]` **Mesin status identik**: `DataTransform/Akseptasi_DT.xml` (`58b8e650`),
`Activity/Akseptasi_Act.xml` (`027a42f4`), `TreatyInAkseptasi_Act.xml` (`226e2ca4`),
`Section/TreatyInActionButtons.xml`, `SaveTreatyIn_Act.xml`,
`TreatyInDeclineConfirmation_postact.xml`. **Tangga persetujuan kedua modul adalah rule yang sama**
— lihat `Treaty In.md` §2.2.

`[terverifikasi]` Untuk 43 rule yang tetap berbeda, **struktur langkah dan Section yang di-include
justru identik** (diperiksa pada `Activity/TreatyInSubmit.xml` dan `Harness/InputTreatyInOffer.xml`;
selisih ukuran harness 1.396 byte). Letak persis selisihnya **tidak dikarakterisasi** — batas
**cakupan telusur**, bahan OQ-010/OQ-011.

### 2.2 Yang khas modul ini — 56 file eksklusif

`[terverifikasi]` **Tampilan data lama (23 rule)**: 9 `FlowAction/*OldData*`
(`CoBListOldData`, `DetailEGNPIOldData`, `DetailLimitsOldData`, `DetailShareOldData`,
`LayersOldData`, `LimitProportionalOldData`, `MaxRetentionOldData`, `TotalLimitsOldData`,
`ShareOldData`) + 14 `Section/*OldData*` + `Harness/TreatyInFacultativeShareCalculationOldData.xml`
(762.247 B), ditambah `Installments_ReadOnly`.

**Setiap layar data treaty punya pasangan `OldData` sebagai pembanding** — tampilan "sebelum vs
sesudah" yang tidak ada di `Treaty In`.

Audit: `ls "Treaty In Adjustment/FlowAction/" | grep -ci olddata` → 9;
`ls "Treaty In Adjustment/Section/" | grep -ci olddata` → 14.

### 2.3 Mesin penomoran revisi — terbaca penuh

`[terverifikasi]` `Activity/TreatyInRevisi_post.xml` (`DATA-PORTAL!TREATYINREVISI_POST`,
52.947 B, 5 langkah):

```
1  Call TreatyInSetAddendumToHistory
2  Property-Set   InputData.CARI1 := TreatyIn.ID + "%"
3  RDB-List       -> RDBList/GetTreatyRevisionID.xml
4  Property-Set   TreatyIn.ID := TreatyIn.ID + "/R01" ; TreatyIn.OLDID := TreatyIn.ID
5  Property-Set   TreatyIn.ID := @If(@substring(TreatyIn.ID,10,12) < 10,
                                     @substring(TreatyIn.ID,0,7)+"/R0"+(@toInt(@substring(TreatyIn.ID,10,12))+1),
                                     @substring(TreatyIn.ID,0,7)+"/R"+(@toInt(@substring(TreatyIn.ID,10,12))+1))
```

dan di sisi database `RDBList/GetTreatyRevisionID.xml`:
`select TO_NUMBER(SUBSTR(ID,10,2))+1 as HASIL1 …`

`[terverifikasi]` **Nomor revisi disimpan di dalam string ID**, posisi karakter 10–12,
**ter-hardcode di dua tempat** (rule Pega dan SQL) → **OQ-055**.

`[terverifikasi]` `Activity/TreatyInSetAddendumToHistory.xml` (32.214 B): `Page-Copy`, `Page-Copy`,
`Page-Remove` — penyalinan halaman kerja ke riwayat.

### 2.4 Dua mode: revisi vs adjustment

`[terverifikasi]` `DataTransform/TreatyCreateEDM.xml` (`DATA-PORTAL!TREATYCREATEEDM`, 14.324 B):

```
SET TreatyIn := ""
WHEN  Param.type == "revision"    ->  SET TreatyIn.EDMState := "1"
WHEN  Param.type == "adjustment"  ->  SET TreatyIn.EDMState := "3"
```

`[terverifikasi]` Nilai `EDMState` dalam kondisi: `!= '3'` 28×, `=="3"` 9×, `!= 3` 8×, `=3` 6×.
Nilai `"2"` **tidak muncul** — arti dan keberadaannya **belum terverifikasi** → OQ-020.

`[terverifikasi]` `TreatyIn.RevisionState = "1"` **memendekkan tangga persetujuan**: dari Sec Head
langsung `Resolve Complete`, melewati Dept Head dan Direktur.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 155 Activity, 68 Section, 44 Connect-SQL, 43 FlowAction,
38 DataTransform, 18 ReportDefinition, 6 Harness, 3 When, **1 Menu** (`Menu/MasterNavTreaty.xml` —
satu-satunya `RULE-NAVIGATION` di kelompok Tahap 5), 1 DecisionTable, 1 ConnectREST,
1 SystemSettings. **Nol rule `Flow`.**

| Objek | Rule perujuk |
| --- | ---: |
| `M_ATTACHMENTTREATY_2` | 5 |
| `T_STORAGE_IMAGE`, `POOLDATA.TREATYINPRODUCTION`, `POOLDATA.ACHIEVEMENT` | 3 masing-masing |
| `TREATY_IN_EDM`, **`POOLDATA.TREATY_IN_EDM`**, **`POOLDATA.TREATY_IN`**, `POOLDATA.OS_AKSEPTASI_KLAIM`, `POOLDATA.M_TREATY_IN_EDM`, `POOLDATA.M_KATEGORIMASTERTREATY`, **`M_TREATY_IN_EDM`** | 2 masing-masing |

`[terverifikasi]` **Selisih dari `Treaty In`**: modul ini merujuk `POOLDATA.TREATY_IN`,
`POOLDATA.TREATY_IN_EDM`, `M_TREATY_IN_EDM` yang **tidak** muncul di daftar `Treaty In`; sebaliknya
`POOLDATA.TREATYINOFFER` dan `CATEGORY_ATTACH_REAS` tidak menonjol di sini.

`[terverifikasi]` Menyentuh objek treaty **outward** di **2 file**
(`RDBList/BrowseTreatyOutDetailEDM.xml`, class `ASM-FW-GISFW-INT-TREATYOUTDETAIL`) — salah satu dari
lima modul yang menyentuhnya → OQ-022.

## 4. Integrasi eksternal

`[terverifikasi]` Identik dengan `Treaty In`: satu `ConnectREST/ServiceGoogle.xml`
(`SETTING` → `LinkService!LinkService`, tanpa URL literal), penyimpanan berkas lewat
`T_STORAGE_IMAGE` + `POOLDATA.GET_TOKEN_STORAGE`, resolusi alamat lewat `GetLinkService`
(tabel `M_LINK_SERVICE`, OQ-047).

**Tidak ada Arasapas / Kasir / Konversi / Gemini AI.**

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Treaty In** | 277/323 file identik; mesin status `Akseptasi_DT` satu rule yang sama | `[terverifikasi]` — OQ-010 |
| **NB Treaty In / EDM Treaty In** | berbagi `ASM!SAVETREATYIN` (`e371c194`), `ASM!BROWSETREATYIN` (`196d49b5`); `BrowseTreatyOutDetailEDM` juga ada di EDM Treaty In | `[terverifikasi]` |
| **Domain Claim/Komite** | `POOLDATA.OS_AKSEPTASI_KLAIM` (2 rule) | `[terverifikasi]` titik temu |

`[terverifikasi]` Pola `Treaty In` ↔ `Adjustment` **sejajar dengan OQ-015** (facultative
NB/RNW/EDM): satu basis rule, dibedakan **saat runtime** oleh nilai properti. Di facultative
pembedanya `Quotation.StatusBusiness` (1/2/3); di sini `Param.type` + `EDMState`.
**Penetapan bounded context tetap D4.**

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Lima stored procedure — sama persis dengan `Treaty In`**:
`POOLDATA.PEGA_TREATY_IN`, `PEGA_M_TREATY_IN_EDM`, `PEGA_M_TREATY_IN_DETAIL`,
`PEGA_M_TREATY_IN_DETAIL_EDM`, `GET_TOKEN_STORAGE`. Isinya tidak ada di korpus → OQ-002.

`[terverifikasi]` **Tidak ada objek db-link** — OQ-017 tidak berlaku.

`[terverifikasi]` **Dua `When` yang kondisinya tidak terbaca** (`<pyLabel>` hanya template kosong),
keduanya **hanya ada di modul ini**:

| Rule | Kondisi kosong | Catatan |
| --- | ---: | --- |
| `When/IsTreatyUser.xml` | **6** | namanya menyangkut otorisasi pengguna treaty — **apa yang diujinya tidak dapat dinyatakan** |
| `When/recordEvent.xml` | 1 | `<pyLabel>` = "Record User Waiting Time Events" |

→ OQ-029 (keluarga kelima).

`[terverifikasi]` Rumus limit/layer/spreading/ROL di activity besar **belum dibaca**; modul ini
menambah `Activity/TreatyLoadMasterJoinEdmXOL.xml` yang juga belum dibaca — batas **cakupan
telusur**.

`[terverifikasi]` Kolom `JSONDATA` pada `M_TREATY_IN` / `M_TREATY_IN_EDM` → OQ-012.

`[terverifikasi]` 10 tombol `(dev)` (section identik dengan `Treaty In`) → OQ-050; otorisasi lewat
indeks tetap `pyWorkBasketList(2)` → OQ-051; 5 identitas orang di `Akseptasi_DT` → OQ-053.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **3 file**,
`OperatorID.pyTelephone` **1 file**.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-005**, **OQ-009**, **OQ-010**, **OQ-011**, **OQ-012**, **OQ-018**, **OQ-020**,
**OQ-021**, **OQ-022**, **OQ-027**, **OQ-029**, **OQ-047**, **OQ-050**, **OQ-051**, **OQ-053**,
**OQ-055**, **OQ-059**.
