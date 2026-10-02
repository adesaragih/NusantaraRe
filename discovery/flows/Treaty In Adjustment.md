# Telusur — Treaty In Adjustment (modul **tanpa** rule `Flow`)

STEP D2, Tahap 5 konteks #17. Ditelusur 2026-09-13. Konvensi: **`_METHOD.md`** +
**`_METHOD-noflow.md`**. **379 file** (`find "Treaty In Adjustment" -name '*.xml' | wc -l`).

Artefak ini **mengutamakan yang BEDA dari `Treaty In`**. Untuk perilaku yang identik, rujuk
**`Treaty In.md`** — tidak diulang di sini.

**Titik masuk `[terverifikasi]`:** `Harness/InputTreatyInAdjustment.xml`
→ `RULE-HTML-HARNESS` / **`DATA-PORTAL`** / `INPUTTREATYINADJUSTMENT`, 943.918 byte.
Modul ini juga membawa salinan `Harness/InputTreatyInOffer.xml` (8.226.491 byte) — lihat §1.1.

---

## 1. Berapa banyak yang benar-benar berbeda dari `Treaty In`

### 1.1 Pengukuran `[terverifikasi]`

```
(cd "Treaty In" && find . -name '*.xml' | sed 's|^\./||' | sort) > ti.txt
(cd "Treaty In Adjustment" && find . -name '*.xml' | sed 's|^\./||' | sort) > tia.txt
comm -12 ti.txt tia.txt | wc -l        # 323 file bernama sama
comm -23 ti.txt tia.txt | wc -l        #   6 hanya di Treaty In
comm -13 ti.txt tia.txt | wc -l        #  56 hanya di Adjustment
# lalu untuk tiap file bersama: bandingkan hash ternormalisasi 18 tag
```

| Ukuran | Jumlah |
| --- | ---: |
| File bernama sama | **323** |
| — **identik** setelah normalisasi 18 tag | **277** (85,8 %) |
| — **berbeda** | **46** |
| Hanya di `Treaty In` | 6 |
| Hanya di `Treaty In Adjustment` | **56** |

**Uji lanjutan dengan normalisasi 21 tag** (menambah `pyDelete`, `pyVersionSecure`,
`pzIsPrivateCheckOut` — tiga tag metadata yang ditemukan di Tahap 4):

| Hasil | Jumlah |
| --- | ---: |
| Tetap berbeda | **43** |
| **Konflik semu** (runtuh jadi identik) | **3** — `RDBList/GetLinkStorage_SQL.xml`, `RDBList/GetTokenStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml` |

`[terverifikasi]` Ketiga konflik semu itu **keluarga yang sama** dengan #367/#371/#377 yang
ditemukan di Tahap 4 (`flows/_SUMMARY-facultative.md` §8) — rule penyimpanan berkas
`T_STORAGE_IMAGE`. Polanya konsisten lintas domain.

`[terverifikasi]` **Mesin status `Akseptasi_DT` IDENTIK** (hash `58b8e650`), demikian pula
`Activity/Akseptasi_Act.xml` (`027a42f4`), `TreatyInAkseptasi_Act.xml` (`226e2ca4`),
`TreatyInAkseptasiEDM_Act.xml` (`61e866e1`), `Section/TreatyInActionButtons.xml`,
`Activity/SaveTreatyIn_Act.xml`, `Activity/TreatyInDeclineConfirmation_postact.xml`.

**Konsekuensinya: tangga persetujuan kedua modul adalah rule yang sama.** Seluruh §2
`Treaty In.md` berlaku apa adanya di sini — **tidak ditelusur ulang**.

### 1.2 Yang berbeda dan bagaimana bedanya

`[terverifikasi]` Untuk rule yang tetap berbeda, **struktur langkah dan Section yang di-include
justru sama**:

```
diff <(grep -oE "<pyStepsActivityName>[^<]*" "Treaty In/Activity/TreatyInSubmit.xml" | sed 's/<[^>]*>//') \
     <(grep -oE "<pyStepsActivityName>[^<]*" "Treaty In Adjustment/Activity/TreatyInSubmit.xml" | sed 's/<[^>]*>//')
# -> tidak ada selisih (7 langkah sama, precondition sama: OutputParam.ERRMSG=="")

diff <(grep -oE "<pySection>[^<]*" "Treaty In/Harness/InputTreatyInOffer.xml" | sed 's/<[^>]*>//' | sort | uniq -c) \
     <(grep -oE "<pySection>[^<]*" "Treaty In Adjustment/Harness/InputTreatyInOffer.xml" | sed 's/<[^>]*>//' | sort | uniq -c)
# -> tidak ada selisih; ukuran file berbeda 1.396 byte (8.225.095 vs 8.226.491)
```

**Selisih isinya ada, tetapi tidak pada tulang punggung yang terbaca dari tag struktural.**
Letak persisnya **tidak dikarakterisasi** — itu memerlukan pembacaan penuh 43 file besar, di luar
cakupan telusur ini. Dicatat sebagai **batas cakupan**, dan sebagai bahan OQ-010/OQ-011.

Daftar 43 yang tetap berbeda mencakup: `Harness/InputTreatyInOffer.xml`,
`Activity/TreatyInSubmit.xml`, `Activity/TreatyInActualShare.xml`,
`Activity/TreatyInDifferenceFacShare.xml`, `Activity/TreatyInDifferenceSummaryFacShare.xml`,
`Activity/TreatyInXOLAddSpreadingDetail(Actual).xml`, `Activity/DetailCalculation(ROL).xml`,
`Activity/AddSpreadingXOL.xml`, `Activity/SetTreatyinRetro_Act.xml`,
`DataTransform/CountTotalPctSpead.xml`, `DataTransform/CountTotalPctSpreadXOL.xml`,
`DataTransform/ShareRetroFacultative_DT.xml`, `Section/Layers.xml`, `Section/MaxRetention.xml`,
`Section/TotalLimits.xml`, `Section/TreatyInShareProp.xml`, dan 27 lainnya.

---

## 2. Yang khas modul ini — mesin revisi & addendum

`[terverifikasi]` **56 file hanya ada di sini.** Dua kelompok besar:

### 2.1 Tampilan *data lama* — 23 rule

| Jenis | Jumlah | Contoh |
| --- | ---: | --- |
| `FlowAction/*OldData*` | **9** | `CoBListOldData`, `DetailEGNPIOldData`, `DetailLimitsOldData`, `DetailShareOldData`, `LayersOldData`, `LimitProportionalOldData`, `MaxRetentionOldData`, `TotalLimitsOldData`, `ShareOldData` |
| `Section/*OldData*` | **14** | ditambah `TreatyInNONProportionalOldData`, `TreatyInTabsProportionalOldData`, `TreatyInTabsNonProportionalOldData(Share)`, `TreatyInFacultativeShareCalculationOldData` |
| `Harness/*OldData*` | 1 | `TreatyInFacultativeShareCalculationOldData.xml` (762.247 B) |
| Lainnya | | `FlowAction/Installments_ReadOnly.xml`, `Section/Installments_ReadOnly.xml` |

`[terverifikasi]` Polanya: **setiap layar data treaty punya pasangan `OldData` sebagai pembanding**
— limit, layer, share, retensi, CoB, EGNPI, dan perhitungan share facultative. Inilah tampilan
"sebelum vs sesudah" yang tidak ada di `Treaty In`.

Perintah audit:
```
ls "Treaty In Adjustment/FlowAction/" | grep -ci olddata   # 9
ls "Treaty In Adjustment/Section/"    | grep -ci olddata   # 14
```

### 2.2 Mesin penomoran revisi — **terbaca penuh**

`[terverifikasi]` `Activity/TreatyInRevisi_post.xml`
(`RULE-OBJ-ACTIVITY` / `DATA-PORTAL` / `TREATYINREVISI_POST`, 52.947 byte, 5 langkah):

```
1  Call TreatyInSetAddendumToHistory
2  Property-Set    InputData.CARI1 := TreatyIn.ID + "%"
3  RDB-List        -> RDBList/GetTreatyRevisionID.xml
4  Property-Set    TreatyIn.ID    := TreatyIn.ID + "/R01"
                   TreatyIn.OLDID := TreatyIn.ID
5  Property-Set    TreatyIn.ID := @If(@substring(TreatyIn.ID,10,12) < 10,
                                      @substring(TreatyIn.ID,0,7) + "/R0" + (@toInt(@substring(TreatyIn.ID,10,12))+1),
                                      @substring(TreatyIn.ID,0,7) + "/R"  + (@toInt(@substring(TreatyIn.ID,10,12))+1))
```

`[terverifikasi]` `RDBList/GetTreatyRevisionID.xml` berisi:
```sql
select TO_NUMBER(SUBSTR(ID,10,2))+1 as HASIL1 …
```

**Nomor revisi disimpan di dalam string ID**, pada posisi karakter 10–12, dengan sufiks
`/R01`, `/R02`, … dan nilai sebelumnya disalin ke `TreatyIn.OLDID`. Posisi karakter
**ter-hardcode di dua tempat** (rule Pega dan SQL). → **OQ-055**.

`[terverifikasi]` `Activity/TreatyInSetAddendumToHistory.xml`
(`DATA-PORTAL!TREATYINSETADDENDUMTOHISTORY`, 32.214 byte) berisi 3 langkah:
`Page-Copy`, `Page-Copy`, `Page-Remove` — penyalinan halaman kerja ke riwayat.

### 2.3 Dua mode: revisi vs adjustment `[terverifikasi]`

`DataTransform/TreatyCreateEDM.xml` (`DATA-PORTAL!TREATYCREATEEDM`, 14.324 byte):

```
SET TreatyIn := ""
WHEN  Param.type == "revision"    ->  SET TreatyIn.EDMState := "1"
WHEN  Param.type == "adjustment"  ->  SET TreatyIn.EDMState := "3"
```

| Kode | Nilai literal | Arti |
| --- | --- | --- |
| `Param.type` | `"revision"`, `"adjustment"` | **belum terverifikasi** |
| `TreatyIn.EDMState` | `"1"` (revision), `"3"` (adjustment) | **belum terverifikasi** |

`[terverifikasi]` Nilai `EDMState` yang muncul dalam kondisi: `!= '3'` 28×, `=="3"` 9×, `!= 3` 8×,
`=3` 6×. Nilai `"2"` **tidak muncul** — **arti dan keberadaannya belum terverifikasi**.

Perintah audit:
```
grep -rhoE "EDMState[^<]{0,20}" "Treaty In Adjustment" --include="*.xml" \
 | sed 's/&amp;#61;/=/g;s/&amp;quot;/"/g;s/\]\[/ /g;s/[][]//g' | grep -E "[=!]" | sort | uniq -c | sort -rn
```

**Jadi satu modul menangani DUA proses berbeda** — revisi dan adjustment — yang dibedakan oleh
parameter `Param.type`, bukan oleh rule terpisah. Ini bahan langsung untuk OQ-010.

### 2.4 Rule khas lainnya

| Rule | Ukuran | Langkah | Keterangan |
| --- | ---: | ---: | --- |
| `Activity/TreatyRevisionCopyAttachment.xml` | 103.284 | 10 | salin lampiran ke revisi |
| `Activity/SetTreatyInEDM_Act.xml` | 79.163 | 6 | |
| `Activity/TreatyInEDMSetValue.xml` | 74.017 | 9 | |
| `Activity/TreatyLoadMasterJoinEdm(XOL).xml` | — | — | + RDBList pasangannya |
| `Activity/RefreshAchievement.xml` | — | — | |
| `FlowAction/PickerTreatyInMaster.xml`, `PickerTreatyInMasterRevisi.xml` | — | — | pemilih master yang direvisi |
| `RDBList/BrowseTreatyInEDM.xml`, `BrowseTreatyOutDetailEDM.xml`, `BrowseOffer(convert).xml` | — | — | nama file memuat tanda kurung |
| `Menu/MasterNavTreaty.xml` | — | — | **satu-satunya `RULE-OBJ-MENU`** di kelompok Tahap 5 |
| `Harness/ActivityStatusSuccess.xml` | 66.542 | — | class `@BASECLASS` (OQ-009) |

---

## 3. Status / state

`[terverifikasi]` **Identik dengan `Treaty In`** untuk `StatusAkseptasi`
(`"Accept"`, `"Reject"`, `"Decline"`, `"Resolve Complete"`), `ChooseStatusAkseptasi`
(3 nilai), `Position` (6 nilai), `ProportionType` (`"Proportional"` / `"NonProportional"`).
Lihat `Treaty In.md` §2.2 — tidak diulang.

**Tambahan khas modul ini:** `TreatyIn.EDMState` (§2.3), `TreatyIn.OLDID` (§2.2), dan pemakaian
`RevisionState` yang lebih padat:

```
grep -rhoE "RevisionState[^<]{0,24}" "Treaty In Adjustment" --include="*.xml" | … | sort | uniq -c
```
→ `RevisionState='1'` 8×, `==1 && …` 2×, `==1` 2×, `=1` 2×, `!= 1 && …` 2×.

`[terverifikasi]` `ProportionType` lebih sering diuji di sini: `='Proportional'` 12×,
`='NonProportional' &&…` 10×, `=="NonProportional"` 9×, `=="Proportional"` 5×,
`='NonProportional'` 5×, `='Proportional' && O…` 4×.

### 3.1 Dua `When` yang kondisinya tidak terbaca

`[terverifikasi]` `<pyLabel>` hanya berisi template kosong `[first value][relation][second value]`:

| Rule | Jumlah kondisi kosong |
| --- | ---: |
| `When/IsTreatyUser.xml` | **6** |
| `When/recordEvent.xml` (`<pyLabel>` = "Record User Waiting Time Events") | 1 |

Keduanya **hanya ada di modul ini**. `IsTreatyUser` namanya menyangkut otorisasi pengguna treaty —
tetapi **kondisinya tidak dapat dibaca**, jadi apa yang diujinya **tidak dinyatakan**. Ini keluarga
kelima pola When tak terbaca di D2 → OQ-029.

---

## 4. Objek Oracle yang disentuh

```
awk -F'\t' '$1 ~ "^Treaty In Adjustment/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
  | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

| Objek | Rule perujuk |
| --- | ---: |
| `M_ATTACHMENTTREATY_2` | 5 |
| `T_STORAGE_IMAGE`, `POOLDATA.TREATYINPRODUCTION`, `POOLDATA.ACHIEVEMENT` | 3 masing-masing |
| `TREATY_IN_EDM`, **`POOLDATA.TREATY_IN_EDM`**, **`POOLDATA.TREATY_IN`**, `POOLDATA.OS_AKSEPTASI_KLAIM`, `POOLDATA.M_TREATY_IN_EDM`, `POOLDATA.M_KATEGORIMASTERTREATY`, **`M_TREATY_IN_EDM`** | 2 masing-masing |
| `TREATYINDETAILEDM`, `TREATYINDETAIL`, `TREATYEXCHANGEYEARLY`, `PROPORTIONALARRG`, `POOLDATA.T_FOLDER_IMAGE` | 1 masing-masing |

`[terverifikasi]` **Selisih dari `Treaty In`:** modul ini merujuk `POOLDATA.TREATY_IN` dan
`POOLDATA.TREATY_IN_EDM` (masing-masing 2 rule) serta `M_TREATY_IN_EDM` — objek yang **tidak**
muncul di daftar `Treaty In`. Sebaliknya `Treaty In` merujuk `POOLDATA.TREATYINOFFER` dan
`CATEGORY_ATTACH_REAS` yang tidak menonjol di sini.

`[terverifikasi]` **Stored procedure: sama persis dengan `Treaty In`** — 5 procedure
(`PEGA_TREATY_IN`, `PEGA_M_TREATY_IN_EDM`, `PEGA_M_TREATY_IN_DETAIL`,
`PEGA_M_TREATY_IN_DETAIL_EDM`, `GET_TOKEN_STORAGE`). Isinya tidak ada di korpus (OQ-002).

---

## 5. Integrasi eksternal

`[terverifikasi]` Identik dengan `Treaty In`: satu `ConnectREST/ServiceGoogle.xml`
(`SETTING` / `LinkService!LinkService`, tanpa URL literal), penyimpanan berkas lewat
`T_STORAGE_IMAGE` + `POOLDATA.GET_TOKEN_STORAGE`, dan resolusi alamat lewat `GetLinkService`
(tabel `M_LINK_SERVICE`, OQ-047).

Tidak ada Arasapas / Kasir / Konversi / Gemini AI.

---

## 6. Batas pengetahuan

1. **Selisih isi 43 rule tidak dikarakterisasi** (§1.2) — batas **cakupan telusur**.
2. **Stored procedure** (5, §4) — batas korpus, OQ-002.
3. **Kolom JSON** `JSONDATA` pada `M_TREATY_IN` / `M_TREATY_IN_EDM` — OQ-012.
4. **Rumus limit/layer/spreading/ROL** di activity besar — sama dengan `Treaty In.md` §6.2;
   modul ini menambah `Activity/TreatyLoadMasterJoinEdmXOL.xml` yang juga **belum dibaca**.
5. **`When/IsTreatyUser.xml`** — 6 kondisi tak terbaca (§3.1), OQ-029.
6. **Harness `InputTreatyInOffer.xml` 8.226.491 byte** — hanya di-grep, tidak dibaca utuh.

---

## 7. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Treaty In Adjustment/`) | OQ-011? |
| --- | --- | --- |
| `DATA-PORTAL` / `INPUTTREATYINADJUSTMENT` / `RULE-HTML-HARNESS` | `Harness/InputTreatyInAdjustment.xml` | tidak — **khas modul ini** |
| `DATA-PORTAL` / `INPUTTREATYINOFFER` / `RULE-HTML-HARNESS` | `Harness/InputTreatyInOffer.xml` | **YA** — beda dgn `Treaty In` |
| `DATA-PORTAL` / `TREATYINFACULTATIVESHARECALCULATIONOLDDATA` / `RULE-HTML-HARNESS` | `Harness/TreatyInFacultativeShareCalculationOldData.xml` | tidak — khas |
| `@BASECLASS` / `ACTIVITYSTATUSSUCCESS` / `RULE-HTML-HARNESS` | `Harness/ActivityStatusSuccess.xml` | tidak (OQ-009) |
| **`DATA-PORTAL` / `AKSEPTASI_DT` / `RULE-OBJ-MODEL`** | `DataTransform/Akseptasi_DT.xml` | tidak — **identik dgn `Treaty In` (`58b8e650`)** |
| `DATA-PORTAL` / `TREATYCREATEEDM` / `RULE-OBJ-MODEL` | `DataTransform/TreatyCreateEDM.xml` | tidak — khas |
| `DATA-PORTAL` / `TREATYINSETEDITPRE` / `RULE-OBJ-MODEL` | `DataTransform/TreatyInSetEditPre.xml` | tidak — khas |
| `DATA-PORTAL` / `TREATYINREVISI_POST` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyInRevisi_post.xml` | tidak — khas |
| `DATA-PORTAL` / `TREATYINSETADDENDUMTOHISTORY` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyInSetAddendumToHistory.xml` | tidak — khas |
| `DATA-PORTAL` / `TREATYREVISIONCOPYATTACHMENT` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyRevisionCopyAttachment.xml` | tidak — khas |
| `DATA-PORTAL` / `SETTREATYINEDM_ACT`, `TREATYINEDMSETVALUE` / `RULE-OBJ-ACTIVITY` | `Activity/` | tidak — khas |
| `DATA-PORTAL` / `TREATYINSUBMIT` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyInSubmit.xml` | **YA** — beda dgn `Treaty In` |
| `ASM-FW-GISFW-INT-TREATY_IN` / `AKSEPTASI_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/Akseptasi_Act.xml` | tidak — identik (`027a42f4`) |
| `…` / `GETTREATYREVISIONID` / `RULE-CONNECT-SQL` | `RDBList/GetTreatyRevisionID.xml` | tidak — khas |
| `…` / `BROWSETREATYINEDM`, `BROWSETREATYOUTDETAILEDM`, `TREATYLOADMASTERJOINEDM(XOL)`, `COPYALLATTACHMENT2_SQL` / `RULE-CONNECT-SQL` | `RDBList/` | tidak — khas |
| `ASM-FW-GISFW-INT-TREATY_IN` / `ASM!SAVETREATYIN`, `ASM!BROWSETREATYIN` / `RULE-CONNECT-SQL` | `RDBList/` | tidak — identik 6 modul |
| `…` / `ISTREATYUSER`, `RECORDEVENT` / `RULE-OBJ-WHEN` | `When/` | tidak — **kondisi tak terbaca** |
| `…` / `MASTERNAVTREATY` / `RULE-OBJ-MENU` | `Menu/MasterNavTreaty.xml` | tidak — satu-satunya Menu |
| 9 `RULE-OBJ-FLOWACTION` `*OldData*` + `Installments_ReadOnly` + 2 `PickerTreatyInMaster*` | `FlowAction/` | tidak — khas |
| 14 `RULE-HTML-SECTION` `*OldData*` | `Section/` | tidak — khas |

**31 rule ditelusur** (4 Harness; 3 DataTransform; 7 Activity; 7 Connect-SQL; 2 When; 1 Menu;
kelompok `OldData` FlowAction/Section dihitung sebagai dua kelompok) — di luar 277 rule yang
**terbukti identik** dengan `Treaty In` dan karena itu tidak ditelusur ulang.

---

## 8. Pertanyaan terbuka

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-055** | Nomor revisi treaty disimpan di dalam string `ID` pada posisi karakter 10–12, ter-hardcode di rule Pega **dan** di SQL (`SUBSTR(ID,10,2)`) | DBA + Product+UW |

### 8.1 Kontribusi berbukti untuk OQ-010 (`Treaty In` vs `Treaty In Adjustment`)

`[terverifikasi]` Bukti yang dapat diaudit ulang:

| Pertanyaan | Bukti |
| --- | --- |
| Berbagi basis kode? | **ya** — 277 dari 323 file bernama sama **identik isinya** (85,8 %) |
| Berbagi mesin status? | **ya** — `Akseptasi_DT` identik (`58b8e650`); tangga persetujuannya satu rule yang sama |
| Ada proses yang khas? | **ya** — 56 file eksklusif: 23 layar *data lama*, mesin penomoran revisi, `TreatyCreateEDM` |
| Dibedakan oleh apa? | **`Param.type`** (`"revision"` / `"adjustment"`) dan **`TreatyIn.EDMState`** (`"1"` / `"3"`); ditambah `RevisionState` yang memendekkan tangga persetujuan |
| Satu modul atau dua? | **tidak diputuskan di sini** — itu D4 |

`[terverifikasi]` Pola ini **sejajar dengan OQ-015** (facultative NB/RNW/EDM): satu basis rule,
dibedakan saat runtime oleh nilai properti. Bedanya, di facultative pembedanya
`Quotation.StatusBusiness` (1/2/3) dan di sini `Param.type` + `EDMState`. **Kesamaan pola dicatat;
penetapan konteks bukan di sini.**

OQ dikuatkan: **OQ-005** (titik masuk), **OQ-010** (§8.1 — bukti utama), **OQ-011** (§1.1 — 43
tetap beda, 3 konflik semu), **OQ-002** (§4), **OQ-012**, **OQ-020** (§2.3), **OQ-023**,
**OQ-027**, **OQ-029** (§3.1 — keluarga kelima), **OQ-047**, **OQ-050**–**OQ-054**
(berlaku sama di sini karena rule-nya identik).
