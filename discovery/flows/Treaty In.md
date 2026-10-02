# Telusur — Treaty In (modul **tanpa** rule `Flow`)

STEP D2, Tahap 5 konteks #16. Ditelusur 2026-09-13 dari korpus READ-ONLY `D:\XML\RNM_BRD\`.
Konvensi: **`_METHOD.md`** + **`_METHOD-noflow.md`**. **329 file**
(`find "Treaty In" -name '*.xml' | wc -l`).

**Titik masuk `[terverifikasi]`:** `Treaty In/Harness/InputTreatyInOffer.xml`
→ `RULE-HTML-HARNESS` / **`DATA-PORTAL`** / `INPUTTREATYINOFFER`, **8.225.095 byte** —
file terbesar di korpus. Dibaca **hanya lewat grep** (`_METHOD.md` §2.5).

Modul ini **tidak punya rule `Flow`** (OQ-005), sehingga tidak ada `<pyStartActivity>`,
`<pyFrom>`/`<pyTo>`, maupun `<pyShapeType>`. Mesin statusnya direkonstruksi sesuai
`_METHOD-noflow.md` §3.

**Catatan OQ-018:** setiap pernyataan menyebut file yang dibaca. Korpus memuat hostname DEV dan
dirakit dari >1 server Pega — **belum tentu cerminan production**.

---

## 1. Alur — direkonstruksi dari tombol + guard

### 1.1 Titik masuk dan layar

`[terverifikasi]` Tiga `Harness`:

| Harness | `pxInsName` | Ukuran |
| --- | --- | ---: |
| `InputTreatyInOffer.xml` | `DATA-PORTAL!INPUTTREATYINOFFER` | **8.225.095** |
| `TreatyInFacultativeShareCalculation.xml` | `DATA-PORTAL!TREATYINFACULTATIVESHARECALCULATION` | 761.861 |
| `ShowAttachmentTreaty.xml` | `ASM-FW-GISFW-INT-TREATY_IN!SHOWATTACHMENTTREATY` | 281.454 |

Section yang di-`<pySection>` oleh harness utama `[terverifikasi]`:
`WorkAttachments` (6×), `TreatyInFacultativeRetro` (4×), `InputTreatyInAdjustment` (4×),
`TreatyInTabsNonProportionalValueDifference` (2×), `TotalLimitsRetro` (2×),
`TreatyInNONProportional` (1×).

Perintah audit:
```
grep -oE "<pySection>[^<]*" "Treaty In/Harness/InputTreatyInOffer.xml" \
  | sed 's/<[^>]*>//' | sort | uniq -c | sort -rn
```

`[terverifikasi]` Harness modul **Treaty In** meng-*include* Section bernama
**`InputTreatyInAdjustment`** — layar yang nama utamanya milik modul `Treaty In Adjustment`.
Dicatat sebagai fakta; kaitannya dibahas di `Treaty In Adjustment.md` §1 dan OQ-010.

### 1.2 Panel tombol — daftar transisi yang mungkin

`[terverifikasi]` `Treaty In/Section/TreatyInActionButtons.xml`
(`RULE-HTML-SECTION` / `DATA-PORTAL` / `TREATYINACTIONBUTTONS`, 369.789 byte) memuat seluruh
tombol aksi. Label yang terbaca dari `<pyLabel>`:

| Tombol (label) | Rule terikat (`<pyName>`) |
| --- | --- |
| **`Submit`** | `Akseptasi_DT` (lewat `TreatyInSubmit`) |
| **`Submit Revision`** | `Revision` / `Akseptasi_DT` |
| `Save` | — |
| `Save ROL Profile` | — (`Activity/TreatyInSaveROL.xml`) |
| `Fix Soa Upload RIComm` | — |
| `Actions`, `Close` | — (navigasi) |
| `Force Edit (dev)` | `TreatyInForceEdit` |
| `Force Resolve Complete(dev)` | `TreatyInForceResolveComplete` |
| `ReturnToInputor(dev)` | `TreatyInReturntoInputor` |
| `RemoveLastComment (dev)` | `TreatyInRemoveLastComment` |
| `Populate treatyindetail (dev)`, `Save(dev)`, `Save EDM(dev)`, `Test Error (dev)`, `TestCopyDifference(dev)`, `fix Spl(dev)` | — |

Perintah audit:
```
grep -oE "<pyLabel>[^<]{1,45}" "Treaty In/Section/TreatyInActionButtons.xml" | sed 's/<pyLabel>//' | sort -u
grep -oE "<pyName>[^<]{1,45}"  "Treaty In/Section/TreatyInActionButtons.xml" | sed 's/<pyName>//'  | sort -u
```

`[terverifikasi]` **Sepuluh label memuat penanda `(dev)`**
(`grep -oE "<pyLabel>[^<]{1,45}" … | grep -ci "(dev)"` → 10), termasuk tombol yang memaksa status
akhir (`Force Resolve Complete`) dan memaksa mode edit (`Force Edit`). Jumlah yang sama di
`Treaty In Adjustment`. → **OQ-050**.

### 1.3 Guard tombol

`[terverifikasi]` Ekspresi visibilitas yang terbaca (apa adanya, `_METHOD-noflow.md` §3.2):

```
TreatyIn.ViewState != '1' && TreatyIn.StatusAkseptasi != 'Resolve Complete'
TreatyIn.IsEditData != '1' && TreatyIn.StatusAkseptasi != 'Resolve Complete'
((TreatyIn.ViewState == '1' && TreatyIn.ProportionType == 'Proportional')
   && TreatyIn.StatusAkseptasi = 'Resolve Complete') && OperatorID.pyWorkBasketList(2).pyWorkBasketName…
((TreatyIn.ViewState == '1' && TreatyIn.ProportionType == 'NonProportional')
   && TreatyIn.StatusAkseptasi = 'Resolve Complete') && OperatorID.pyWorkBasketList(2).pyWorkBasketName…
FALSE && OperatorID.pyWorkBasketList(2).pyWorkBasketName = 'ReasTreatyInAdmin'
   && TreatyIn.Position = '' && TreatyIn.StatusAkseptasi = 'Resolve Complete'
TreatyIn.ViewState != '1' && 1=2
```

`[terverifikasi]` **Dua guard tidak dapat bernilai benar menurut ekspresi yang terbaca**:
satu diawali `FALSE &&`, satu memuat `1=2`. Dicatat apa adanya — **tidak disimpulkan** bahwa
tombolnya mati (`_METHOD-noflow.md` §3.2). Paralel dengan OQ-023.

`[terverifikasi]` Otorisasi tombol diuji lewat **indeks tetap** pada daftar workbasket operator:
`OperatorID.pyWorkBasketList(2).pyWorkBasketName` (36 kemunculan) dan
`pyWorkBasketList(1).pyWorkBasketName` (4). Perintah audit:
```
grep -rhoE "pyWorkBasketList\([0-9]+\)\.pyWorkBasketName[^<]{0,34}" "Treaty In" "Treaty In Adjustment" \
  --include="*.xml" | sed 's/&amp;#61;/=/g;s/&amp;quot;/"/g' | sort | uniq -c | sort -rn
```
→ `(2)... = 'ReasTreatyInAdmin'` 36×, `(2)... = TreatyIn.Position` 8×,
`(1)... = 'ReasTreatyInSecHead'` 4×, `(2)... = 'ReasTreatyInSecHead'` 1×.

**Hak akses bergantung pada urutan workbasket dalam profil operator**, bukan pada pencarian
berdasarkan nama. → **OQ-051**.

### 1.4 Rantai aksi `Submit` `[terverifikasi]`

```
Section/TreatyInfoSubmit  ──tombol "Submit"──> Activity/TreatyInSubmit
                                                (DATA-PORTAL!TREATYINSUBMIT, 62.246 B, 7 langkah)
   1  Property-Set
   2  Call TreatyInCheckID                (DATA-PORTAL!TREATYINCHECKID, 109.830 B, 9 langkah)
   3  call TreatyInCheckError             (DATA-PORTAL!TREATYINCHECKERROR, 55.265 B, 6 langkah,
                                           2× Property-Set-Messages)
      precondition langkah berikutnya: OutputParam.ERRMSG == ""
   4  Apply-DataTransform  Akseptasi_DT   <-- MESIN STATUS, §2
   5  Call AddCommentList_Act
   6  Call SaveTreatyIn_Act               (DATA-PORTAL!SAVETREATYIN_ACT, 92.681 B, 11 langkah)
   7  Call TreatyInInputVis               (DATA-PORTAL!TREATYININPUTVIS, 40.124 B, 3 langkah)
```

`SaveTreatyIn_Act` langkahnya `[terverifikasi]`: `Property-Set`, `call CheckDuplicateOffer`,
`Property-Set`, `RDB-List`, `Property-Set`, `Page-Remove`, `Property-Set`,
`call Data-Portal.SaveTreatyInDetail_Act`, `call Data-Portal.SaveTreatyInOffer_Act`,
`Property-Set`, `call TreatyInInputVis`.

### 1.5 Rantai aksi `Decline` `[terverifikasi]`

```
Section/TreatyInDeclineConfirmation ──> Activity/TreatyInDeclineConfirmation_postact
                                        (DATA-PORTAL!TREATYINDECLINECONFIRMATION_POSTACT, 43.316 B)
   1  Property-Set   TreatyIn.CommentList(<LAST>).IsApproved := "Decline"
                     TreatyIn.StatusAkseptasi                := "Decline"
                     Param.SkipVisibility := "0"; Param.Info := TreatyIn.Information;
                     Param.Comment := TreatyIn.Comment
   2  Call AddCommentList_Act
   3  Property-Set
   4  Call SaveTreatyIn_Act
```

Ini **satu-satunya jalur yang men-set `StatusAkseptasi` secara langsung** di luar `Akseptasi_DT`.

### 1.6 Bentuk proses yang dapat dinyatakan

**Tahapan `[terverifikasi]`:** layar input master treaty → `Submit` (validasi ID + error) →
`Akseptasi_DT` menentukan posisi & status berikutnya → simpan master + detail + offer.
Jalur terpisah: `Decline` langsung men-set status.

**Yang TIDAK dapat dinyatakan** (`_METHOD-noflow.md` §3.5): urutan wajib antar layar, mekanisme
assignment Pega, SLA/timer, dan titik akhir proses — tidak ada sumber tag untuk keempatnya.

---

## 2. Mesin status — `Akseptasi_DT`

`[terverifikasi]` `Treaty In/DataTransform/Akseptasi_DT.xml`
(`RULE-OBJ-MODEL` / `DATA-PORTAL` / `AKSEPTASI_DT`) — **10 `WHEN`, 18 `OTHERWISE_WHEN`, 64 `SET`**
(`grep -oE "<pyActionName>[^<]*" … | sort | uniq -c`).

Hash ternormalisasi `58b8e650` — **identik dengan `Treaty In Adjustment`**.

Dibaca dengan awk berpenghitung kedalaman (`_METHOD-noflow.md` §3.3). **Nilai nama orang tidak
disalin**; ditulis `<orang>`.

### 2.1 Tangga lengkap

```
WHEN  TreatyIn.RevisionState != 1 && TreatyIn.StatusAkseptasi != "Resolve Complete"
  │
  ├─ WHEN  TreatyIn.Position == "" || TreatyIn.Position == "ReasTreatyInAdmin"
  │     ├─ WHEN            OperatorID.pyTelephone == "SPVTREATY1"
  │     │      SET Position := "ReasTreatyInSecHead" ; StatusAkseptasi := "Accept" ; PositionUsername := "<orang>"
  │     ├─ WHEN            OperatorID.pyTelephone == "SPVTREATY2"   → idem
  │     ├─ OTHERWISE_WHEN  OperatorID.pyTelephone == "TREATY1"      → idem
  │     └─ OTHERWISE_WHEN  OperatorID.pyTelephone == "TREATY2"      → idem
  │
  ├─ OTHERWISE_WHEN  TreatyIn.Position == "ReasTreatyInSecHead"
  │     ├─ WHEN            ChooseStatusAkseptasi == "Accept"
  │     │      SET Position := "ReasTreatyInDeptHead" ; StatusAkseptasi := "Accept" ; PositionUsername := "<orang>"
  │     ├─ OTHERWISE_WHEN  ChooseStatusAkseptasi == "Reject"
  │     │      SET Position := "ReasTreatyInAdmin"    ; StatusAkseptasi := "Reject"
  │     │          PositionUsername := TreatyIn.CommentList(1).OperatorName
  │     └─ OTHERWISE_WHEN  ChooseStatusAkseptasi == "Decline"
  │            SET Position := ""                    ; StatusAkseptasi := "Decline"  ; PositionUsername := ""
  │
  ├─ OTHERWISE_WHEN  TreatyIn.Position == "ReasTreatyInDeptHead"
  │     ├─ Accept  → Position := "ReasTreatyInDirector" ; Status := "Accept" ; PositionUsername := "<orang>"
  │     ├─ Reject  → Position := "ReasTreatyInAdmin"    ; Status := "Reject" ; PositionUsername := CommentList(1).OperatorName
  │     └─ Decline → Position := ""                     ; Status := "Decline"
  │
  ├─ OTHERWISE_WHEN  TreatyIn.Position == "ReasTreatyInGroupLeader"
  │     ├─ Accept  → Position := "ReasTreatyInDirector" ; Status := "Accept" ; PositionUsername := "<orang>"
  │     ├─ Reject  → Position := "ReasTreatyInAdmin"    ; Status := "Reject"
  │     └─ Decline → Position := ""                     ; Status := "Decline"
  │
  └─ OTHERWISE_WHEN  TreatyIn.Position == "ReasTreatyInDirector"
        ├─ Accept  → Position := ""                     ; Status := **"Resolve Complete"**
        ├─ Reject  → Position := "ReasTreatyInAdmin"    ; Status := "Reject"
        └─ Decline → Position := ""                     ; Status := "Decline"

OTHERWISE_WHEN  TreatyIn.RevisionState == 1 && TreatyIn.StatusAkseptasi != "Resolve Complete"
  ├─ WHEN  TreatyIn.Position == "" || TreatyIn.Position == "ReasTreatyInAdmin"
  │     → Position := "ReasTreatyInSecHead" ; Status := "Accept" ; PositionUsername := "<orang>"
  └─ OTHERWISE_WHEN  TreatyIn.Position == "ReasTreatyInSecHead"
        ├─ Accept  → Position := "" ; Status := "Resolve Complete" ; RevisionState := "" ; ViewState := ""
        ├─ Reject  → Position := "ReasTreatyInAdmin" ; Status := "Reject" ; RevisionState := "1"
        │            PositionUsername := TreatyIn.CommentList(<LAST>).OperatorName
        └─ Decline → Position := "" ; Status := "Decline" ; RevisionState := ""
```

`[terverifikasi]` **Jalur revisi jauh lebih pendek**: dari `ReasTreatyInSecHead`, `Accept` langsung
menjadi `Resolve Complete` — tanpa Dept Head maupun Direktur.

### 2.2 Nilai literal

| Properti | Nilai literal yang muncul | Arti |
| --- | --- | --- |
| `TreatyIn.ChooseStatusAkseptasi` | `"Accept"`, `"Reject"`, `"Decline"` — **3 nilai, dipilih pengguna** | **belum terverifikasi** |
| `TreatyIn.StatusAkseptasi` | `"Accept"`, `"Reject"`, `"Decline"`, **`"Resolve Complete"`** — **4 nilai, dihitung mesin** | **belum terverifikasi** |
| `TreatyIn.Position` | `""`, `"ReasTreatyInAdmin"`, `"ReasTreatyInSecHead"`, `"ReasTreatyInDeptHead"`, `"ReasTreatyInDirector"`, `"ReasTreatyInGroupLeader"` | nama workbasket sebagai **nilai data** |
| `TreatyIn.RevisionState` | `""`, `"1"` | **belum terverifikasi** |
| `TreatyIn.ViewState` | `"1"`, `""` | 506+288+233+52+40 kemunculan berkondisi |
| `TreatyIn.IsEditData` | `"1"` | **belum terverifikasi** |
| `TreatyIn.ProportionType` | **`"Proportional"`, `"NonProportional"`** | pembeda lini treaty — §3 |

Perintah audit:
```
grep -rhoE "(StatusAkseptasi|ChooseStatusAkseptasi|ProportionType|RevisionState|ViewState|IsEditData)[^<]{0,26}" \
  "Treaty In" --include="*.xml" | sed 's/&amp;#61;/=/g;s/&amp;quot;/"/g;s/\]\[/ /g;s/[][]//g' \
  | grep -E "[=!]" | sort | uniq -c | sort -rn
```

`[terverifikasi]` **`ChooseStatusAkseptasi` punya 3 nilai, `StatusAkseptasi` punya 4.**
`"Resolve Complete"` **tidak pernah dipilih pengguna** — ia hanya dihasilkan mesin di ujung tangga.

`[pertanyaan terbuka]` Mengapa ada **`Reject` dan `Decline` sekaligus**, dan apa bedanya? Keduanya
adalah pilihan pengguna yang berbeda, dengan akibat berbeda (`Reject` mengembalikan case ke Admin
dan menyimpan nama komentator pertama; `Decline` mengosongkan posisi). **Arti bisnisnya tidak ada
di korpus** → OQ-020, **OQ-052**.

### 2.3 Otorisasi: kode peran + identitas orang ter-hardcode

`[terverifikasi]` Tangga masuk dipilih oleh **`OperatorID.pyTelephone`** dengan **empat** nilai:
`"TREATY1"`, `"TREATY2"`, `"SPVTREATY1"`, `"SPVTREATY2"`.

D1/Tahap 1 hanya menemukan `TREATY1` dan `SPVTREATY1` (OQ-027). **Dua nilai baru ditemukan di
sini.** → OQ-027 diperluas.

`[terverifikasi]` `Akseptasi_DT` men-`SET` `TreatyIn.PositionUsername` dengan **5 identitas orang
berbeda yang ter-hardcode**. Perintah audit yang **menghitung tanpa menampilkan nilai**:
```
awk '/<pyPropertiesName>/{n=$0;gsub(/<[^>]*>/,"",n)}
     /<pyPropertiesValue>/{v=$0;gsub(/<[^>]*>/,"",v);
        if(n=="TreatyIn.PositionUsername" && v!="" && v !~ /\./){gsub(/"/,"",v); print v}; n=""}' \
  "Treaty In/DataTransform/Akseptasi_DT.xml" | sort -u | wc -l
```
→ **5**.

**Ini pola baru dan lebih berat dari OQ-021.** Di OQ-021 identitas orang dipakai sebagai *guard*
(menguji siapa yang sedang bekerja). Di sini identitas orang **ditetapkan sebagai pemilik tugas
berikutnya** — jadi tangga persetujuan treaty inward menunjuk **orang tertentu**, bukan peran.
Bila orang itu berganti jabatan, rule harus diubah. → **OQ-053**.

Nilai nama orang **tidak disalin** ke artefak ini (`_METHOD.md` §1.4).

---

## 3. Proportional vs Non-Proportional

`[terverifikasi]` `TreatyIn.ProportionType` bernilai literal **`"Proportional"`** dan
**`"NonProportional"`**; keduanya dipakai di guard visibilitas dan di perhitungan.

```
grep -rhoE "ProportionType[^<]{0,28}" "Treaty In" --include="*.xml" \
 | sed 's/&amp;#61;/=/g;s/&amp;quot;/"/g;s/\]\[/ /g;s/[][]//g' | grep -E "[=!]" | sort | uniq -c | sort -rn
```
→ `=="NonProportional"` 9×, `='Proportional'` 9×, `='NonProportional' &&…` 8×,
`=="Proportional"` 5×, `= "Proportional"` 4×, `= "NonProportional"` 4×.

`[terverifikasi]` Rule yang bercabang menurutnya, antara lain:
`Activity/TreatyInLimitsListValueProportional.xml`, `Activity/TreatyInNPSetTotalActualLimits.xml`,
`Section/TreatyInShareProp.xml`, `Section/TreatyInNONProportional*.xml`,
`FlowAction/LimitProportional.xml`, `FlowAction/Layers.xml`, `FlowAction/MaxRetention.xml`.

`[terverifikasi]` FlowAction yang membentuk struktur data treaty (32 buah), dikelompokkan:

| Kelompok | FlowAction |
| --- | --- |
| Struktur limit & layer | `DetailLimits`, `TotalLimits`, `TotalLimitsRetro`, `Layers`, `LayersEDM`, `MaxRetention`, `LimitProportional`, `LimitFacRetro` |
| Share / retrosesi | `Share`, `ShareOldData`, `ShareRetro`, `DetailShare`, `DetailShareRetro`, `TreatyRetroList`, `TreatyinChooseRetro` |
| Spreading | `SpreadingTPDtl`, `SpreadingTXOLDtl` |
| Premi & angsuran | `Installments`, `DetailEGNPI`, `AchievementCombine` |
| Aksi status | **`TreatyInAction`**, `TreatyInActionEDM`, **`TreatyInDeclineConfirmation`**, `TreatyInDeclineConfirmationEDM` |
| Master/pencarian | `CoBList`, `CoBListReadOnly`, `TreatyInSearchReinsured`, `TreatyInSearchSoB`, `TreatyInUploadCSV`, `TreatyAttachContent`, `ShowSummary`, `Action` |

`[terverifikasi]` `EGNPI` muncul **1.047 kali** di modul ini, dengan properti
`EGNPIAMOUNT`, `EGNPIPROPORTION`, `EGNPIAMOUNTNP`, dan rule `Activity/TreatyInEGNPIListValue.xml`,
`Activity/TotalEgnpi.xml`. **Kepanjangan `EGNPI` belum terverifikasi.**

---

## 4. Objek Oracle yang disentuh

```
awk -F'\t' '$1 ~ "^Treaty In/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
  | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

| Objek | Rule perujuk |
| --- | ---: |
| `M_ATTACHMENTTREATY_2` | 4 |
| `T_STORAGE_IMAGE`, `POOLDATA.TREATYINPRODUCTION`, `POOLDATA.ACHIEVEMENT` | 3 masing-masing |
| `TREATY_IN_EDM`, `POOLDATA.OS_AKSEPTASI_KLAIM`, `POOLDATA.M_KATEGORIMASTERTREATY`, `CATEGORY_ATTACH_REAS` | 2 masing-masing |
| `TREATYINDETAILEDM`, `TREATYINDETAIL`, `TREATYEXCHANGEYEARLY`, `PROPORTIONALARRG`, `POOLDATA.T_FOLDER_IMAGE`, `POOLDATA.TREATYINOFFER`, **`POOLDATA.M_TREATY_IN`**, `POOLDATA.M_TREATY_IN_EDM` | 1 masing-masing |

### 4.1 Kaitan ke realisasi treaty (NB Treaty In) `[terverifikasi]`

Modul ini **menulis master** `M_TREATY_IN`; modul `NB Treaty In` menjalankan **realisasi** lewat
`Flow/InputRealizationTreatyIn.xml` (`NB Treaty In.md`). Titik temunya terbaca:

| Rule di Treaty In | Menyentuh |
| --- | --- |
| `RDBList/SaveTreatyIn.xml` (`ASM-FW-GISFW-INT-TREATY_IN!ASM!SAVETREATYIN`) | procedure `POOLDATA.PEGA_TREATY_IN` |
| `RDBList/BrowseTreatyIn.xml` (`…!ASM!BROWSETREATYIN`) | `POOLDATA.M_TREATY_IN`, `POOLDATA.M_TREATY_IN_EDM` |
| `Activity/FetchTreatyExistingProduction.xml`, `RDBList/FetchTreatyInProductionUsingNooffer.xml` | `POOLDATA.TREATYINPRODUCTION` |

`[terverifikasi]` `ASM!SAVETREATYIN` (hash `e371c194`) dan `ASM!BROWSETREATYIN` (`196d49b5`)
**identik di 6 modul**: Treaty In, Treaty In Adjustment, NB Treaty In, EDM Treaty In, NB FacIn,
**Master Product Name Life** — dan **tidak terdaftar di register OQ-011**. Satu rule bersama,
bukan salinan yang menyimpang. Perintah audit:
```
find . -path ./OUTPUT_HASIL_RNM -prune -o -type f -name "SaveTreatyIn.xml" -print \
| while IFS= read -r f; do grep -o "<pxInsName>[^<]*" "$f" | head -1; sh nhash.sh "$f"; done
```

`[terverifikasi]` `POOLDATA.OS_AKSEPTASI_KLAIM` (2 rule) juga disentuh modul ini — tabel yang sama
yang dipakai domain Claim dan Komite (Tahap 2 & 3). Dicatat sebagai titik temu; **tidak
disimpulkan ulang**.

---

## 5. Integrasi eksternal

`[terverifikasi]` Satu rule `RULE-CONNECT-REST`: `ConnectREST/ServiceGoogle.xml`,
`pyBaseURLSelectionType = SETTING`, `pyBaseURLSetting = LinkService!LinkService` — **tanpa URL
literal endpoint**. Sama polanya dengan domain facultative (Tahap 4).

`[terverifikasi]` Penyimpanan berkas lewat Google Storage: `Activity/InsertGoogleStorage_Act.xml`,
`Activity/GetLinkService.xml`, dan 3 Connect-SQL `GetTokenStorage_SQL`, `GetLinkStorage_SQL`,
`Update_T_Storage_SQL` (tabel `T_STORAGE_IMAGE`, `POOLDATA.T_FOLDER_IMAGE`, procedure
`POOLDATA.GET_TOKEN_STORAGE`).

`[terverifikasi]` Alamat layanan diambil dari tabel Oracle `M_LINK_SERVICE` lewat `GetLinkService`
— pola yang sama dengan temuan Tahap 4 → **OQ-047**. Isi tabel tidak ada di korpus.

`[terverifikasi]` **Tidak ada integrasi Arasapas, Kasir, Konversi, maupun Gemini AI di modul ini**
(`ls "Treaty In/Activity/" | grep -icE "arasapas|kasir|konversi|gemini"` → 0).

---

## 6. Batas pengetahuan

### 6.1 Stored procedure — 5 dipanggil, isinya tidak diketahui

`[terverifikasi]`
```
POOLDATA.PEGA_TREATY_IN             POOLDATA.PEGA_M_TREATY_IN_EDM
POOLDATA.PEGA_M_TREATY_IN_DETAIL    POOLDATA.PEGA_M_TREATY_IN_DETAIL_EDM
POOLDATA.GET_TOKEN_STORAGE
```
**Isinya tidak ada di korpus** (OQ-002). Empat dari lima adalah **penulis master treaty inward** —
artinya cara master treaty disimpan (validasi, versi, kunci) **berada di sisi database**.

`[terverifikasi]` Tidak ada objek `@ASMD.SINARMAS.CO.ID` di modul ini (OQ-017 tidak berlaku).

### 6.2 Activity besar yang **belum habis dibaca**

| Activity | Ukuran |
| --- | ---: |
| `TreatyInNPSetTotal.xml` | **825.279** |
| `SaveTreatyInDetailEdm_Act.xml` | 764.703 |
| `SaveTreatyInDetail_Act.xml` | 655.709 |
| `GetAchievement.xml` | 518.157 |
| `SetSpreadingXOL.xml` | 509.259 |
| `FetchQSfromMasterXOL.xml` | 479.714 |

`[terverifikasi]` **31 activity** bernama limit/layer/spreading/retention/EGNPI/ROL
(`ls "Treaty In/Activity/" | grep -icE "layer|spread|limit|retention|egnpi|rol"` → 31).
**Rumus limit, layer, spreading, dan ROL tidak dinyatakan di artefak ini** — batas **cakupan
telusur** (volume), bukan batas korpus. `ROL`, `QS`, `XOL`, `TP`, `TXOL` **kepanjangan belum
terverifikasi**.

### 6.3 Kolom JSON

`[terverifikasi]` `BrowseTreatyIn` menyeleksi kolom `JSONDATA` dari `POOLDATA.M_TREATY_IN` dan
`POOLDATA.M_TREATY_IN_EDM` — ini justru **temuan asal OQ-012**. Strukturnya tidak ada di korpus.

### 6.4 Berkas non-rule di dalam korpus

`[terverifikasi]` D1 menginventarisasi hanya `*.xml`. Modul lain di korpus memuat folder
**`Claude outputs`** berisi berkas non-Pega. Di kelompok Tahap 5: `Treaty Contract Out/Claude
outputs/Struktur_InboxTreatyContract.xlsx` (306.979 B) dan `Master Contract Retro Life/Claude
outputs/` (`Struktur_GridRetrocessionLife.xlsx` 11.944 B, `perubahan_skill.diff` 10.278 B).
Modul `Treaty In Adjustment` juga memuat `Struktur_InputRenewalFacultativeIn.xlsx`-sejenis di
modul lain. **Isinya tidak dibaca** — bukan rule Pega, dan korpus READ-ONLY. → **OQ-054**.

---

## 7. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Treaty In/`) | OQ-011? |
| --- | --- | --- |
| `DATA-PORTAL` / `INPUTTREATYINOFFER` / `RULE-HTML-HARNESS` | `Harness/InputTreatyInOffer.xml` | **YA** (beda dgn Adjustment) |
| `DATA-PORTAL` / `TREATYINFACULTATIVESHARECALCULATION` / `RULE-HTML-HARNESS` | `Harness/TreatyInFacultativeShareCalculation.xml` | tidak (identik) |
| `ASM-FW-GISFW-INT-TREATY_IN` / `SHOWATTACHMENTTREATY` / `RULE-HTML-HARNESS` | `Harness/ShowAttachmentTreaty.xml` | tidak (identik) |
| `DATA-PORTAL` / `TREATYINACTIONBUTTONS` / `RULE-HTML-SECTION` | `Section/TreatyInActionButtons.xml` | tidak (identik) |
| **`DATA-PORTAL` / `AKSEPTASI_DT` / `RULE-OBJ-MODEL`** | `DataTransform/Akseptasi_DT.xml` | tidak — **identik (`58b8e650`)** |
| `DATA-PORTAL` / `TREATYINSUBMIT` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyInSubmit.xml` | **YA** (beda dgn Adjustment) |
| `DATA-PORTAL` / `TREATYINCHECKID` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyInCheckID.xml` | tidak (identik) |
| `DATA-PORTAL` / `TREATYINCHECKERROR` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyInCheckError.xml` | tidak (identik) |
| `DATA-PORTAL` / `SAVETREATYIN_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/SaveTreatyIn_Act.xml` | tidak (identik) |
| `DATA-PORTAL` / `TREATYININPUTVIS` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyInInputVis.xml` | tidak (identik) |
| `DATA-PORTAL` / `TREATYINDECLINECONFIRMATION_POSTACT` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyInDeclineConfirmation_postact.xml` | tidak (identik) |
| `DATA-PORTAL` / `TREATYINAKSEPTASI_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/TreatyInAkseptasi_Act.xml` | tidak — identik (`226e2ca4`) |
| **`ASM-FW-GISFW-INT-TREATY_IN` / `AKSEPTASI_ACT` / `RULE-OBJ-ACTIVITY`** | `Activity/Akseptasi_Act.xml` | tidak — **identik (`027a42f4`); lihat §8 koreksi** |
| `DATA-PORTAL` / `TREATYINFORCEEDIT`, `TREATYINFORCERESOLVECOMPLETE`, `TREATYINRETURNTOINPUTOR`, `TREATYINREMOVELASTCOMMENT` / `RULE-OBJ-MODEL` | `DataTransform/` | tidak |
| `ASM-FW-GISFW-INT-TREATY_IN` / `ASM!SAVETREATYIN` / `RULE-CONNECT-SQL` | `RDBList/SaveTreatyIn.xml` | tidak — identik 6 modul |
| `ASM-FW-GISFW-INT-TREATY_IN` / `ASM!BROWSETREATYIN` / `RULE-CONNECT-SQL` | `RDBList/BrowseTreatyIn.xml` | tidak — identik 6 modul |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETTOKENSTORAGE_SQL`, `RNM!GETLINKSTORAGE_SQL`, `RNM!UPDATE_T_STORAGE_SQL` / `RULE-CONNECT-SQL` | `RDBList/` | **YA — konflik semu, §8** |
| `DATA-PORTAL` / `SERVICEGOOGLE` / `RULE-CONNECT-REST` | `ConnectREST/ServiceGoogle.xml` | sebagian |
| 32 `RULE-OBJ-FLOWACTION` (§3) | `FlowAction/` | sebagian |

**24 rule ditelusur** (3 Harness; 2 Section; 1 DataTransform mesin status + 4 DataTransform aksi;
8 Activity; 5 Connect-SQL; 1 ConnectREST — FlowAction dihitung sebagai kelompok).

Rule **dirujuk tapi belum ditelusur**: 31 activity limit/layer/spreading (§6.2),
`SaveTreatyInDetail_Act`, `SaveTreatyInOffer_Act`, `AddCommentList_Act`, `CheckDuplicateOffer`,
`TreatyInSaveROL`, `GetAchievement`, seluruh 50 Section.

---

## 8. Pertanyaan terbuka & koreksi

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-050** | 10 tombol ber-label `(dev)` di panel aksi produksi, termasuk `Force Resolve Complete` dan `Force Edit` | Product+UW + pemilik export |
| **OQ-051** | Otorisasi tombol memakai **indeks tetap** `OperatorID.pyWorkBasketList(2)` — bergantung urutan workbasket di profil operator | IAM |
| **OQ-052** | `Reject` dan `Decline` hidup berdampingan sebagai dua pilihan pengguna dengan akibat berbeda; bedanya tidak dijelaskan korpus | Product+UW |
| **OQ-053** | `Akseptasi_DT` **menetapkan 5 identitas orang ter-hardcode sebagai pemilik tugas berikutnya** — pola lebih berat dari OQ-021 (yang hanya menguji identitas) | IAM + Product+UW |
| **OQ-054** | Folder `Claude outputs` (berkas `.xlsx`/`.diff` non-Pega) berada di dalam korpus ekspor 4 modul | pemilik export |

### 8.1 Koreksi terhadap teks OQ-011 `[terverifikasi]`

Narasi OQ-011 mencontohkan `ASM-FW-GISFW-INT-TREATY_IN / AKSEPTASI_ACT` sebagai rule yang
**berbeda** antara `Treaty In` dan `Treaty In Adjustment`. Pengukuran ulang menunjukkan keduanya
**identik** (hash `027a42f4`), dan identitas itu **tidak ada di register**
`../inventory/_oq011-konflik-isi.md`. Contoh tersebut berasal dari pengukuran batch 1 dengan
normalisasi 7 tag, sebelum perbaikan ke 18 tag (yang memangkas konflik 976 → 543). **Register
tetap berlaku; contohnya yang usang.** Diperbaiki di `../open-questions.md`.

OQ dikuatkan: **OQ-005** (§1 — titik masuk ditetapkan), **OQ-002** (§6.1), **OQ-007/OQ-024**
(§1.3, §2 — nama workbasket sebagai nilai data), **OQ-010** (lihat `Treaty In Adjustment.md`),
**OQ-011** (§7, §8.1), **OQ-012** (§6.3), **OQ-020** (§2.2), **OQ-021** (§2.3),
**OQ-027** (§2.3 — dua nilai `pyTelephone` baru), **OQ-023** (§1.3 — guard yang tak dapat benar),
**OQ-047** (§5).

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **3 file**,
`OperatorID.pyTelephone` **1 file**, `pxCreateOperator` sebagai guard **0 file**.
`IsPEGAPROD` **tidak ada** di modul ini.
