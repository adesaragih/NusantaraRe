# Telusur — Treaty Contract Out (modul **tanpa** rule `Flow`)

STEP D2, Tahap 5 konteks #18. Ditelusur 2026-09-13. Konvensi: **`_METHOD.md`** +
**`_METHOD-noflow.md`**. **303 file** (`find "Treaty Contract Out" -name '*.xml' | wc -l`).

**Titik masuk `[terverifikasi]`:** `Harness/InboxTreatyContract.xml`
→ `RULE-HTML-HARNESS` / **`DATA-PORTAL`** / `INBOXTREATYCONTRACT`, 128.846 byte.

Modul ini adalah subjek **OQ-022** — apakah namanya menyesatkan. §2 menjawabnya dengan bukti.

---

## 1. Alur — direkonstruksi dari Harness + panel

### 1.1 Tiga Harness `[terverifikasi]`

| Harness | `pxInsName` | Ukuran | Section yang di-`<pySection>` |
| --- | --- | ---: | --- |
| `InboxTreatyContract.xml` | `DATA-PORTAL!INBOXTREATYCONTRACT` | 128.846 | — |
| `InboxTreatyContractDescription.xml` | `DATA-PORTAL!INBOXTREATYCONTRACTDESCRIPTION` | 502.598 | `GridTreatyArrangementPortfolioList`, `GridTreatyArrangementProfitCommision`, `NitipKurs`, `ViewDetailDescription` |
| `InboxTreatyContractReinsType.xml` | `DATA-PORTAL!INBOXTREATYCONTRACTREINSTYPE` | 85.622 | — |

`[terverifikasi]` Modul ini punya **hanya 2 `FlowAction`** — `DetailTreatyExclustion.xml` dan
`TreatyOutAttachContent.xml` — jumlah paling sedikit di korpus untuk modul sebesar ini
(bandingkan `Treaty In Adjustment`: 43). Perilakunya **tidak digerakkan tombol FlowAction**,
melainkan oleh 168 `Activity` yang dipanggil dari Section grid.

### 1.2 Bentuk modul: **editor master, bukan proses berjenjang**

`[terverifikasi]` **Tidak ada `Akseptasi_DT`, tidak ada `StatusAkseptasi`, tidak ada `Position`,
tidak ada tombol Submit/Akseptasi/Decline.** Perintah audit:
```
ls "Treaty Contract Out/DataTransform/"                 # 10 file, tidak ada Akseptasi_DT
grep -rl "StatusAkseptasi" "Treaty Contract Out" --include="*.xml" | wc -l    # 0
```

Yang ada sebagai gantinya adalah **pasangan simetris** per jenis klausul kontrak:

| Pola nama | Jumlah | Guna |
| --- | ---: | --- |
| `SaveTreatyArr*_Act` | 27 | simpan satu jenis klausul |
| `NewTreatyArr*`, `SetTreatyArr*` | — | buat / isi baris baru |
| `BrowseTreatyArr*` | — | tampilkan daftar |
| **`CancelActivity*`** | **16** | batalkan pengeditan satu jenis klausul |
| Total activity bernama `*TreatyArr*` | **81** | — |

`ls "Treaty Contract Out/Activity" | grep -ci treatyarr` → **81**.

`[terverifikasi]` Enam belas `CancelActivity*` menamai jenis klausul yang dapat diedit:
`BordereAux`, `CashLossLimit`, `ClaimCoorperation`, `Epi`, `ExGratia`, `ExGratiaLimitChild`,
`FacIn`, `FacInList`, `PLA`, `ProfitCommision`, `Ricomm`, `TreatyContract`, `TreatyLimit`,
`TreatyLimitChild`, `pPortfolio`, `pTerrLimit`.

Kepanjangan `EPI`, `PLA`, `Ricomm`, `LOL`, `MB` **belum terverifikasi**.

**Tahapan yang dapat dinyatakan `[terverifikasi]`:** buka inbox kontrak → pilih kontrak → edit per
klausul lewat grid (`New…` / `Set…` / `Save…` / `CancelActivity…`) → simpan ke master lewat
Connect-SQL. **Tidak ada tangga persetujuan.** Tidak ada assignment, SLA, atau titik akhir proses
(`_METHOD-noflow.md` §3.5).

---

## 2. **UJI NAMA MODUL** — outward atau inward?

Task ini menguji apakah `Treaty Contract Out` benar tentang treaty **outward**, atau justru master
*term/arrangement* treaty. Lima uji terpisah, seluruhnya dapat diaudit ulang.

### 2.1 Uji objek database — **nol objek outward** `[terverifikasi]`

```
awk -F'\t' '$1 ~ "^Treaty Contract Out/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
  | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

| Objek | Rule perujuk |
| --- | ---: |
| **`M_PROPORTIONALARRG`** | **10** |
| **`MTREATYSECURITY`** | 5 |
| `TREATYREINSURER` | 3 |
| `TREATYBUSINESS` | 3 |
| `T_STORAGE_IMAGE`, `M_ATTACHMENTTREATY_2` | 2 masing-masing |
| `TREATYEXCHANGE`, `TREATYCONTRACT`, `PROPORTIONALARRG`, `M_TREATYYEAR`, `CATEGORY_ATTACH_REAS` | 1 masing-masing |

**Tidak satu pun bernama `*_OUT*`.**

```
for pat in M_TREATY_OUT TREATY_OUT TREATYOUTDETAIL FACOUTPRODUCTION; do
  echo "$pat: $(grep -rli "$pat" "Treaty Contract Out" --include='*.xml' | wc -l)"
done
# -> M_TREATY_OUT 0 ; TREATY_OUT 0 ; TREATYOUTDETAIL 0 ; FACOUTPRODUCTION 0
```

### 2.2 Uji class rule — **tidak ada class outward** `[terverifikasi]`

```
awk -F'\t' '$1 ~ "^Treaty Contract Out/"{print $3}' all-rules.tsv | sort | uniq -c | sort -rn
```

| Class | Rule |
| --- | ---: |
| `@BASECLASS` | **202** |
| **`ASM-FW-GISFW-INT-PROPORTIONALARRG`** | **40** |
| `DATA-PORTAL` | 8 |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 7 |
| **`ASM-FW-GISFW-INT-TREATY_IN`** | **5** |
| `ASM-FW-GISFW-INT-TREATYREINSURER` | 5 |
| `ASM-FW-GISFW-INT-TREATYYEAR` | 4 |
| `ASM-FW-GISFW-INT-TREATYBUSINESS` | 4 |
| `ASM-FW-GISFW-INT-MTREATYSECURITY` | 4 |
| `ASM-FW-GISFW-INT` | 4 |

`awk -F'\t' '$1 ~ "^Treaty Contract Out/"{print $3}' all-rules.tsv | grep -i out` → **kosong**.

Sebaliknya, **5 rule di modul ini berclass `ASM-FW-GISFW-INT-TREATY_IN`** — class treaty
**inward**: `DELETE_ACT`, `DELETEATTACHMENT2_SQL`, `GETALLATTACHMENT2_SQL`, `GETATTACHMENT2_SQL`,
`INSERTATATCHMENT_SQL`. Dan **11 file** menyebut class itu
(`grep -rl "ASM-FW-GISFW-INT-TREATY_IN" "Treaty Contract Out" --include="*.xml" | wc -l`).

### 2.3 Uji penamaan rule `[terverifikasi]`

```
ls "Treaty Contract Out/Activity" | grep -ci treatyarr        # 81
ls "Treaty Contract Out/Activity" | grep -ci out              #  4
ls "Treaty Contract Out/Activity" | grep -ciE "facin|treatyin" # 12
```

Empat rule ber-nama `Out` seluruhnya soal **lampiran**: `LoadAttachmentTreatyOut`,
`TreatyOutDownloadAll_Act`, `TreatyOutDownloadOne`, `TreatyOutSaveAttachment` (ditambah
`FlowAction/TreatyOutAttachContent.xml`). Isinya pun hanya soal berkas:

```
grep -rl "TREATYOUT" "Treaty Contract Out" --include="*.xml"
# -> 5 file, seluruhnya Activity/FlowAction lampiran
```

Dua belas rule ber-nama inward: `BrowseDeleteRowTreatyInContract`, `BrowseTreatyArrFacInParentList`,
`CancelActivityFacIn(List)`, `NewTreatyArrFacIn(List)`, `SaveTreatyArrFacIn(List)_Act`,
`SetCategoryAttachTreatyin`, `SetTreatyArrFacIn(List)_Act`, `TreatyInitAttach`.

### 2.4 Uji arah tulis — modul ini **MENULIS** master, bukan membacanya `[terverifikasi]`

Sembilan rule `RULE-CONNECT-SQL` memanggil stored procedure penulis:

| Rule | Procedure |
| --- | --- |
| `RDBList/SaveMasterTreatyContract_SQL.xml` | `POOLDATA.PEGA_TREATYCONTRACT` |
| `RDBList/SaveMasterTreatyYear_SQL.xml` | `POOLDATA.PEGA_TREATYYEAR` |
| `RDBList/SaveMasterTreatyReinsurer_SQL.xml` | `POOLDATA.PEGA_TREATYREINSURER` |
| `RDBList/SaveMasterTreatyBusiness_SQL.xml` | `POOLDATA.PEGA_TREATYBUSINESS` |
| `RDBList/SaveMasterProportionalArrg.xml` | `POOLDATA.PEGA_PROPORTIONALARRG` |
| `RDBList/SaveMasterProportionalArrgChild.xml` | `POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD` |
| `RDBList/SaveMasterCopyData_SQL.xml` | `POOLDATA.PROSESCOPY` |
| `RDBList/InsertAtatchment_Sql.xml` | `POOLDATA.PEGA_M_ATTACHMENT` |
| `RDBList/GetTokenStorage_SQL.xml` | `POOLDATA.GET_TOKEN_STORAGE` |

Contoh isi: `SaveMasterProportionalArrg` = `BEGIN POOLDATA.PEGA_PROPORTIONALARRG ({InputTreatyArrTreatyLimit.ID}, …`.

**Modul ini adalah pemilik tulis master arrangement treaty.**

### 2.5 Uji pembanding — **siapa yang sebenarnya menyentuh treaty outward** `[terverifikasi]`

```
for m in <20 modul>; do grep -rli "M_TREATY_OUT\|TREATY_OUT\|TREATYOUTDETAIL" "$m" --include="*.xml" | wc -l; done
```

| Modul | File yang menyentuh objek outward |
| --- | ---: |
| **NB Treaty In** | **8** |
| **Claim Non Prop** | **4** |
| **EDM Treaty In** | **4** |
| **Treaty In Adjustment** | 2 |
| **Claim Prop** | 1 |
| **`Treaty Contract Out`** | **0** |

Satu-satunya class treaty-outward di **seluruh korpus** adalah
`ASM-FW-GISFW-INT-TREATYOUTDETAIL` (5 rule), dan letaknya:
`NB Treaty In/Activity/SetValueRetro_Act.xml`, `NB Treaty In/RDBList/BrowseTreatyOutDetail.xml`,
`NB Treaty In/ReportDefinition/BrowseTreatyOutDetail.xml`,
`Treaty In Adjustment/RDBList/BrowseTreatyOutDetailEDM.xml`,
`EDM Treaty In/RDBList/BrowseTreatyOutDetailEDM.xml` — **tidak satu pun di modul ini**.

### 2.6 Jawaban

`[terverifikasi]` **`Treaty Contract Out` BUKAN modul treaty outward.** Ia adalah **editor master
*term / arrangement* kontrak treaty**: `TREATYCONTRACT`, `M_TREATYYEAR`, `TREATYREINSURER`,
`TREATYBUSINESS`, `M_PROPORTIONALARRG` (+ child), `MTREATYSECURITY`, `TREATYEXCHANGE` — dengan
16 jenis klausul yang dapat diedit dan 9 jalur tulis ke stored procedure.

Kata "Out" pada nama modul **hanya bertahan di 5 rule lampiran**. Sebaliknya modul ini membawa
5 rule berclass treaty **inward** dan 12 rule bernama FacIn/TreatyIn.

`[terverifikasi]` Pekerjaan treaty **outward** yang sesungguhnya ada di **modul lain** — terbanyak
di `NB Treaty In` (8 file) dan `Claim Non Prop` (4 file). Temuan Tahap 3 (OQ-042: klaim
non-proporsional membaca master treaty outward) **tetap berdiri**, dan kini diketahui bahwa ia
tidak membaca dari modul ini.

**Yang masih terbuka:** *mengapa* modul ini dinamai "Out", dan ke bounded context mana 303 rule ini
seharusnya ditempatkan. Penetapan konteks adalah **D4** (`_METHOD.md` §0). → OQ-022 diperbarui,
**belum ditutup**.

---

## 3. Status / state dan kode

`[terverifikasi]` **Tidak ada properti status proses** di modul ini (§1.2). Kode yang ada bersifat
identifikasi data:

| Properti | Nilai literal | Arti |
| --- | --- | --- |
| `ReinsTypeID` | diuji `== ""` **114×**; diisi dari `InputData.CARIREINS…` | **belum terverifikasi** |
| `ReinsTypeName` | — (properti deskripsi) | **belum terverifikasi** |
| `ReinsType.TreatyYear` | dipakai dalam ekspresi `(… % 4) ==` | **belum terverifikasi** |

Perintah audit:
```
grep -rhoE "ReinsType[A-Za-z]*[^<]{0,20}" "Treaty Contract Out" --include="*.xml" \
 | sed 's/&amp;#61;/=/g;s/&amp;quot;/"/g;s/\]\[/ /g;s/[][]//g' | grep -E "[=!]" | sort | uniq -c | sort -rn
```

`[terverifikasi]` `Harness/InboxTreatyContractReinsType.xml` adalah layar tersendiri untuk
`ReinsType` — jadi "jenis reasuransi" adalah **dimensi master tersendiri**, bukan sekadar atribut.

`[terverifikasi]` Dua `RULE-DECLARE-PAGES`: `DataPage/D_TreatyContract.xml` dan
`DataPage/D_EnumerationList.xml` — satu-satunya modul Tahap 5 yang punya DataPage.

---

## 4. Integrasi eksternal

`[terverifikasi]` Satu `RULE-CONNECT-REST`: `ConnectREST/ServiceGoogle.xml`,
`pyBaseURLSelectionType = SETTING`, `pyBaseURLSetting = LinkService!LinkService` — **tanpa URL
literal endpoint**.

Penyimpanan berkas: `T_STORAGE_IMAGE` (2 rule), `POOLDATA.GET_TOKEN_STORAGE`,
`POOLDATA.PEGA_M_ATTACHMENT`, dan 7 rule berclass `ASM-FW-GISFW-INT-T_STORAGE_IMAGE`.

Tidak ada Arasapas / Kasir / Konversi / Gemini AI.

---

## 5. Batas pengetahuan

### 5.1 Stored procedure — 11 dipanggil, isinya tidak diketahui `[terverifikasi]`

```
POOLDATA.PEGA_TREATYCONTRACT        POOLDATA.PEGA_TREATYYEAR
POOLDATA.PEGA_TREATYREINSURER       POOLDATA.PEGA_TREATYBUSINESS
POOLDATA.PEGA_PROPORTIONALARRG      POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD
POOLDATA.PEGA_M_ATTACHMENT          POOLDATA.PROSESCOPY
POOLDATA.GET_TOKEN_STORAGE          DBMS_LOB.CREATETEMPORARY
DBMS_OUTPUT.PUT_LINE
```

**Jumlah procedure terbanyak dari kelima modul Tahap 5** (OQ-002). Karena **seluruh jalur tulis
master melewati procedure**, aturan penyimpanan master arrangement treaty — validasi, versi,
kunci, kaskade ke child — **seluruhnya berada di sisi database dan tidak ada di korpus**.

`[terverifikasi]` `DBMS_OUTPUT.PUT_LINE` muncul — **satu-satunya kemunculannya di korpus**
(`awk -F'\t' '$8 ~ /DBMS_OUTPUT/{print $1}' all-rules.tsv`). Ini pemanggilan keluaran debug
di dalam SQL produksi. `POOLDATA.PROSESCOPY` juga hanya muncul di modul ini.

### 5.2 Rule di class bawaan Pega

`[terverifikasi]` **202 dari 303 rule (66,7 %) berada di `@BASECLASS`** — proporsi tertinggi di
korpus. Ini memperkuat **OQ-009** secara material: dua pertiga modul tidak tinggal di class
aplikasi, sehingga pemetaan rule → entitas domain tidak dapat diturunkan dari class.

### 5.3 Berkas non-rule di dalam ekspor

`[terverifikasi]` `Treaty Contract Out/Claude outputs/Struktur_InboxTreatyContract.xlsx`
(306.979 byte) — berkas **non-Pega** di dalam korpus READ-ONLY. **Tidak dibaca.** → OQ-054.

### 5.4 Yang tidak dapat direkonstruksi

Sesuai `_METHOD-noflow.md` §3.5: urutan wajib antar layar, mekanisme assignment, SLA, titik akhir
proses, dan siapa yang boleh mengedit klausul — **tidak ada sumber tag** untuk kelimanya di modul
ini (bahkan ekspresi visibilitas ber-workbasket pun tidak ditemukan, berbeda dari `Treaty In`).

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Treaty Contract Out/`) | OQ-011? |
| --- | --- | --- |
| `DATA-PORTAL` / `INBOXTREATYCONTRACT` / `RULE-HTML-HARNESS` | `Harness/InboxTreatyContract.xml` | tidak |
| `DATA-PORTAL` / `INBOXTREATYCONTRACTDESCRIPTION` / `RULE-HTML-HARNESS` | `Harness/InboxTreatyContractDescription.xml` | tidak |
| `DATA-PORTAL` / `INBOXTREATYCONTRACTREINSTYPE` / `RULE-HTML-HARNESS` | `Harness/InboxTreatyContractReinsType.xml` | tidak |
| `ASM-FW-GISFW-INT-…` / `SAVEMASTERTREATYCONTRACT_SQL` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterTreatyContract_SQL.xml` | tidak |
| `…` / `SAVEMASTERTREATYYEAR_SQL` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterTreatyYear_SQL.xml` | tidak |
| `…` / `SAVEMASTERTREATYREINSURER_SQL` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterTreatyReinsurer_SQL.xml` | tidak |
| `…` / `SAVEMASTERTREATYBUSINESS_SQL` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterTreatyBusiness_SQL.xml` | tidak |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `SAVEMASTERPROPORTIONALARRG` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterProportionalArrg.xml` | tidak |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `SAVEMASTERPROPORTIONALARRGCHILD` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterProportionalArrgChild.xml` | tidak |
| `…` / `SAVEMASTERCOPYDATA_SQL` / `RULE-CONNECT-SQL` | `RDBList/SaveMasterCopyData_SQL.xml` | tidak |
| **`ASM-FW-GISFW-INT-TREATY_IN`** / `DELETE_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/Delete_act.xml` | **§2.2 — class inward** |
| **`ASM-FW-GISFW-INT-TREATY_IN`** / `INSERTATATCHMENT_SQL`, `GETATTACHMENT2_SQL`, `GETALLATTACHMENT2_SQL`, `DELETEATTACHMENT2_SQL` / `RULE-CONNECT-SQL` | `RDBList/` | **§2.2 — class inward** |
| `@BASECLASS` / 16 `CANCELACTIVITY*` / `RULE-OBJ-ACTIVITY` | `Activity/` | sebagian |
| `@BASECLASS` / 27 `SAVETREATYARR*_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/` | sebagian |
| `@BASECLASS` / `LOADATTACHMENTTREATYOUT`, `TREATYOUTDOWNLOADALL_ACT`, `TREATYOUTSAVEATTACHMENT` / `RULE-OBJ-ACTIVITY` | `Activity/` | sebagian — **satu-satunya jejak "Out"** |
| `DATA-PORTAL` / `TREATYOUTDOWNLOADONE` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyOutDownloadOne.xml` | sebagian |
| `…` / `D_TREATYCONTRACT`, `D_ENUMERATIONLIST` / `RULE-DECLARE-PAGES` | `DataPage/` | sebagian |
| `…` / `DETAILTREATYEXCLUSTION`, `TREATYOUTATTACHCONTENT` / `RULE-OBJ-FLOWACTION` | `FlowAction/` | tidak |
| `DATA-PORTAL` / `SERVICEGOOGLE` / `RULE-CONNECT-REST` | `ConnectREST/ServiceGoogle.xml` | sebagian |

**26 rule ditelusur** (3 Harness; 9 Connect-SQL penulis master; 5 rule berclass TREATY_IN;
2 DataPage; 2 FlowAction; 1 ConnectREST; kelompok `CancelActivity*` dan `SaveTreatyArr*` dihitung
sebagai dua kelompok; 4 rule lampiran "Out").

Rule **dirujuk tapi belum ditelusur**: 168 Activity (mayoritas `New*`/`Set*`/`Browse*` per klausul),
49 Section, 30 ReportDefinition.

---

## 7. Pertanyaan terbuka

Tidak ada OQ baru eksklusif modul ini.

**OQ-022 diperbarui secara material** — pertanyaan "outward atau bukan" **terjawab
`[terverifikasi]`: bukan outward**; yang tersisa adalah *mengapa* dinamai demikian dan ke konteks
mana modul ini ditempatkan (D4). Bukti lengkap: §2.

OQ dikuatkan: **OQ-002** (§5.1 — 11 procedure, seluruh jalur tulis lewat SP),
**OQ-005** (§1 — titik masuk ditetapkan), **OQ-009** (§5.2 — 202/303 rule di `@BASECLASS`),
**OQ-016** (skema `POOLDATA`; `DBMS_OUTPUT` muncul sekali di korpus), **OQ-018/OQ-047** (§4),
**OQ-020** (§3 — `ReinsTypeID`), **OQ-042** (§2.5 — arah pembacaan outward kini terpetakan),
**OQ-054** (§5.3).

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **2 file**,
`OperatorID.pyTelephone` **0**, `pxCreateOperator` sebagai guard **0**. `IsPEGAPROD` **tidak ada**.
