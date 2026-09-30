# Grilling Ronde 6 — Claim Prop · SAPU BERSIH KORPUS

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Claim Prop\` (READ-ONLY)
**Sasaran:** 42 berkas yang belum pernah disebut (§B), gerbang induk `CopyOldataCurr_act` (§C),
sensus tuntas 303 catatan pengembang (§D).

> ⛔ **Nilai rahasia NOL.** Aturan §A1 ronde 4 tetap berlaku penuh.
> ⛔ **Modul lain NOL** — termasuk Komite Claim Prop.
> ⛔ Yang disunting: **`spec.md` (§E)** dan **lima tiket** yang disebut §F1 lebih dulu.

---

## §A — Satu jawaban yang menggantung: **BELUM DIJAWAB**

`[terbuka]` Pertanyaan A1 — *`CountPersen_act`: yang tercoret di layar Pega langkah 5 (seluruh
blok) atau 5.1 saja?* — **tidak diisi** pada blok GANTI KONTEKS ronde ini.

**Yang saya lakukan karenanya, persis seperti yang diperintahkan:**

| | |
| --- | --- |
| Menebak jawabannya | ⛔ **TIDAK** |
| Butir `[terbuka]` ke-6 di `spec.md` | ✅ **dibiarkan berdiri**, tidak disentuh |
| Menyinggungnya lagi di §F | ⛔ **TIDAK** — tidak muncul di §F5 |

⚠️ **Satu hal yang wajib saya catat:** §B ronde ini menemukan **kejadian kedua** dengan bentuk
**sama persis** (`CheckPeriodPolicy_Act` langkah 1). Itu temuan baru dan berdiri sendiri, bukan
pengulangan pertanyaan yang menggantung — ia masuk sebagai butir `[terbuka]` ke-**7**.

---

## §B — Sapu bersih 42 berkas ⭐ sasaran utama

### B1 — Hitung ulang: **42**, cocok

`[terverifikasi]` Dihitung ulang sendiri terhadap seluruh **329** berkas, disaring **tidak peka
huruf besar-kecil**, dengan pembanding **seluruh** keluaran (`grilling-ronde-1…5`,
`periksa-ulang-aturan-baru`, `spec.md`, dan 16 tiket):

| Jenis | Belum pernah disebut | Dari |
| --- | --- | --- |
| Activity | **10** | 112 |
| Section | **11** | 36 |
| Harness | **5** | 11 |
| When | **5** | 61 |
| DataTransform | **4** | 11 |
| ReportDefinition | **4** | 19 |
| FlowAction | **3** | 12 |
| ConnectREST · RDBList · DecisionTable · Flow · SystemSettings | **0** | 67 |
| **TOTAL** | ⭐ **42** | **329** |

⛔ **Angkanya tidak berbeda dari yang saya ajukan di draf brief.** Tidak ada yang perlu diralat.

⚠️ **Batas kejujuran yang tetap berlaku:** *"belum pernah disebut"* **bukan** *"belum dibaca"*.
Saringan ini mengukur nama yang tertulis di keluaran, bukan isi yang pernah dibuka. Saringan yang
kedua tidak dapat dibuat tanpa menebak.

### B2 — Keempat puluh dua, satu per satu

⛔ **Nol berkas yatim.** Keempat puluh dua punya **sekurang-kurangnya satu perujuk** di korpus.
Tidak satu pun dapat disebut mati karena tidak dipanggil.

#### Activity (10)

| Rule (`pxInsName`) | Perujuk | Yang dikerjakannya | Uang | Efek keluar | Langkah ber-remark |
| --- | --- | --- | --- | --- | --- |
| `…WORK-CLAIMTREATY!CHECKDATERECEIVED_ACT` | 2 Section | memeriksa Date Received terhadap hari ini, DOL, Report Date, periode polis | ⛔ | ⛔ | ⚠️ **2** — langkah 10 · 13 |
| `…WORK-CLAIMTREATY!CHECKPERIODPOLICY_ACT` | `SetEndDate_Act` + 2 Section | mencocokkan periode polis dengan periode treaty | ⛔ | ⛔ | ⚠️ **1** — langkah **1** |
| `…WORK-CLAIMTREATY!CHECKREPORTDATE_ACT` | 2 Section | memeriksa Report Date terhadap hari ini, DOL, Date Received, periode polis | ⛔ | ⛔ | ⚠️ **2** — langkah 11 · 14 |
| `…INT-ADJUSTERCONSULTANT!EDITMSTCONSULTANT_ACT` | `MstAdjusterConsultant` | menyalin baris konsultan ke halaman sunting | ⛔ | ⛔ | 0 |
| `…WORK-CLAIMTREATY!MAKELOWERCASE_ACT` | 2 Section | ⚠️ menimpa lokasi kerugian + uraian laporan dengan huruf kecil | ⛔ | ⛔ | 0 |
| `…WORK-CLAIMTREATY!REMOVELOSSALLOCTION_ACT` | 2 Section | menghapus baris loss allocation + menulis kronologi | ⛔ | ⛔ | 0 |
| `…DATA-ADJUSTMENT!SETDLACEDINGSOB` | `AdjustmentDetail` | menyalin `DLANoCeding` dan `DLANoSOB` ke halaman kerja | ⛔ | dokumen | 0 |
| `…WORK-CLAIMTREATY!SETENDDATE_ACT` | 2 Section | menghitung End Date polis dari Start Date | ⛔ | ⛔ | ⚠️ **1** — langkah 8 |
| `…DATA-ADJUSTMENT!SETKOMITENO_ACT` | `AdjustmentDetail` | ⚠️ mengiris nomor komite dari kunci internal Pega | ⛔ | ⛔ | 0 |
| `DATA-PORTAL!SETVIEWDATAMASTERTREATY_ACT` | `OutstandingClaim` | menyiapkan tampilan master treaty | ⛔ | ⛔ | 0 |

#### Section (11) · Harness (5)

| Rule | Perujuk | Isi | Uang |
| --- | --- | --- | --- |
| `…DATA-ADJUSTMENT!ADJUSTMENTTYPE` | `InputAcceptation_Adjs` | pilihan jenis adjustment; 1 penanda gerbang mati | ⛔ |
| `@BASECLASS!BROWSEDETAILCAUSEOFLOSS` | 2 | daftar rinci cause of loss, sumber `BrowseVDCauseOfLoss_RD` | ⛔ |
| `DATA-PORTAL!GRIDCAUSEOFLOSS` | `TambahMasterCauseOfLoss` | kisi cause of loss | ⛔ |
| `…INT-TREATYGROUP!ISDELETETREATYGROUP` | `MessageBeforeDeleteTreatyGroup` | konfirmasi hapus kelompok treaty | ⛔ |
| `…WORK-CLAIMTREATY!LISTPAYMENTCLAIM_SC` | `ListPaymentClaim_Harness` | daftar pembayaran klaim | ⚠️ `.ClaimPaid` · `.Currency` |
| `…WORK-CLAIMTREATY!LISTPOLICYNOTREATY` | `ListPolicyNoTreaty_Harness` | daftar nomor polis; 1 penanda gerbang mati (*"hide button choose"*) | ⛔ |
| `…WORK-CLAIMTREATY!PAYMENTPREMILIST_SECTION` | `ListPaymentPremi_Harness` | daftar pembayaran premi | ⚠️ `.AgingAmount` · `.TotalPayment` · `.CurrencyID` |
| `…WORK-CLAIMTREATY!PRINTFILEDLA` | `AdjustmentDetail` | pembungkus cetak DLA | ⛔ |
| `…WORK-CLAIMTREATY!SUMMARYOUTSCLAIM` | 3 | ringkasan outstanding; **3** penanda gerbang mati | ⛔ |
| `…WORK-CLAIMTREATY!VIEWATTACHMENT` | 2 | penampil lampiran | ⛔ |
| `…DATA-CURRENCY!VIEWDETAILPAYMENT_SC` | `ViewDetailPayment_FA` | rincian pembayaran | ⚠️ `.PaymentAmount` |
| `…!LISTPAYMENTCLAIM_HARNESS` · `…!LISTPAYMENTPREMI_HARNESS` · `…!LISTPOLICYNOTREATY_HARNESS` · `…!SUMMARYOUTSCLAIM` · `…!VIEWATTACHMENT` | 1–3 | pembungkus layar bagi Section di atasnya | mengikuti Section-nya |

#### DataTransform (4) · ReportDefinition (4) · FlowAction (3)

| Rule | Isi | Catatan pengembang |
| --- | --- | --- |
| `@BASECLASS!CNMREFRESHCAUSEOFLOSS_DT` · `…CNMSHOWINSERTCAUSEOFLOSS_DT` · `…CNMSHOWINSERTDETAILCAUSEOFLOSS_DT` | penyegar/penampil layar cause of loss | *"CHANGE CLASS"* ×2 · *"DONE"* |
| `…WORK-CLAIMTREATY!SETLABEL_DT` | menyetel label tampilan | *"set label\\"* |
| `…INT-BUSINESS!BROWSEBUSINESS_RD` | daftar lini bisnis | *"balikin sts aktif"* |
| `…INT-CURRENCY!BROWSECURRENCY_RD` | ⭐ **daftar mata uang — 12 perujuk**, yang terbanyak di antara 42 | *"Tambah Param untuk filter ID"* |
| `…INT-RW!BROWSEPROVINCE_RD` | daftar provinsi | *"OK"* |
| `…INT-V_M_CAUSE_OF_LOSS!SELECTVMCAUSEOFLOSS_RD` | pemilih cause of loss | *"done"* |
| `…INT-TREATYGROUP!MESSAGEBEFOREDELETETREATYGROUP` | konfirmasi sebelum hapus | *"Hide messageBeforeDeleteTreatyGroup"* |
| `…WORK-CLAIMTREATY!PRINTFILEDLA` | aksi cetak DLA | — |
| `…DATA-CURRENCY!VIEWDETAILPAYMENT_FA` | aksi lihat rincian pembayaran | — |

### B3 — Yang MENGUBAH pembacaan: **sebelas**

| # | Temuan | Bab `spec.md` | Tiket |
| --- | --- | --- | --- |
| 1 | ⭐ ⚠️ **Proteksi periode polis DIMATIKAN sengaja** — `CHECKREPORTDATE_ACT` langkah 11 · 14 ber-remark, catatan *"matikan protek dalam periode polis"* | Sapu bersih korpus | **01** |
| 2 | ⚠️ Pola sama untuk Date Received — `CHECKDATERECEIVED_ACT` langkah 10 · 13 ber-remark | Sapu bersih korpus | **01** |
| 3 | ⭐ ⚠️ **`CHECKPERIODPOLICY_ACT` — pengisi menggantung, kejadian KEDUA** | Sapu bersih korpus + `[terbuka]` 7 | **01** |
| 4 | ⭐ ⚠️ **Periode polis boleh "TBA"** — konsep belum pernah tercatat | Sapu bersih korpus + `[terbuka]` 8 | **01** |
| 5 | ⚠️ **Nomor komite diiris berdasarkan posisi karakter** — `SETKOMITENO_ACT` | Sapu bersih korpus | **11** |
| 6 | ⚠️ **Teks bebas pengguna diturunkan huruf secara merusak** — `MAKELOWERCASE_ACT` | Sapu bersih korpus | **14** · 01 |
| 7 | **Batas waktu 300 000 ms (5 menit)** — `SERVICEGOOGLE` | Sapu bersih korpus | **13** |
| 8 | ⚠️ **`SETPAYABLETREATY_ACT` tidak pernah menyimpan** — `Obj-Save` + `Obj-Refresh-And-Lock` ber-remark | Sapu bersih korpus | **13** |
| 9 | ⚠️ *"samain dengan prod"* pada `COUNTVALUEADJTREATY_ACT` — menyiratkan beda antar-lingkungan | `[terbuka]` 9 | — |
| 10 | `SETDLACEDINGSOB` — dua nomor dokumen yang belum tercatat | Sapu bersih korpus | 12 *(tidak ditambal)* |
| 11 | `REMOVELOSSALLOCTION_ACT` + `ADDINTEREST_ACT` — penulis kronologi ketiga | Sapu bersih korpus | **14** |

### B4 — Enam pasang nama kembar: **tiga sepasang, tiga kebetulan nama**

`[terverifikasi]` Dibuktikan lewat **`pxInsName`**, bukan nama berkas:

| Pasangan | `pxInsName` | Vonis |
| --- | --- | --- |
| `PrintFileDLA` Section ↔ FlowAction | `…WORK-CLAIMTREATY!PRINTFILEDLA` di keduanya | ⭐ **rule yang SAMA** |
| `ViewAttachment` Section ↔ Harness | `…WORK-CLAIMTREATY!VIEWATTACHMENT` di keduanya | ⭐ **rule yang SAMA** |
| `SummaryOutsClaim` Section ↔ `SummaryOutSClaim` Harness | `…WORK-CLAIMTREATY!SUMMARYOUTSCLAIM` di keduanya | ⭐ **rule yang SAMA** *(beda besar-kecil huruf pada nama berkas saja)* |
| `ListPaymentClaim_SC` ↔ `ListPaymentClaim_Harness` | `…!LISTPAYMENTCLAIM_SC` vs `…!LISTPAYMENTCLAIM_HARNESS` | ⛔ **KEBETULAN NAMA** |
| `ListPolicyNoTreaty` ↔ `ListPolicyNoTreaty_Harness` | `…!LISTPOLICYNOTREATY` vs `…!LISTPOLICYNOTREATY_HARNESS` | ⛔ **KEBETULAN NAMA** |
| `ViewDetailPayment_SC` ↔ `ViewDetailPayment_FA` | `…DATA-CURRENCY!VIEWDETAILPAYMENT_SC` vs `…_FA` | ⛔ **KEBETULAN NAMA** |

⭐ **Tiga berbanding tiga.** Kalau vonisnya diambil dari kemiripan nama berkas, **separuhnya
salah** — persis jebakan yang pernah mengenai saya pada tiga `IsCLM*.xml`.

### B5 — Lima rule `When`: **nol yang hidup**

| Rule (`pxInsName`) | Halaman yang diuji | Ada di objek kerja Claim Prop? |
| --- | --- | --- |
| `ASM-FW-GISFW-DATA!ISMBD` | `pyWorkPage.Quotation.BusinessType` | ⛔ tidak |
| `ASM-FW-GISFW-DATA!ISMBU` | `pyWorkPage.OfferFacIn.QuotationData.BusinessType` | ⛔ tidak |
| `ASM-FW-GISFW-DATA!ISMARINEHULL` | `pyWorkPage.Quotation.BusinessType` | ⛔ tidak |
| `@BASECLASS!ISBILLBOARDNEONSYARIAH` | `pyWorkPage.Quotation.BusinessType` + `.BusinessCode` | ⛔ tidak |
| `@BASECLASS!ISMAINTENANCE` | `pyWorkPage.Quotation.BusinessType` | ⛔ tidak |

⭐ **Kelimanya termasuk 52 rule sisa impor. AC 72 tidak berubah**, dan penggolongannya tetap
`[terbuka]` milik work owner.

> ### ⚠️ Aturan "ketiadaan" baru saja menyelamatkan satu kesimpulan salah
>
> `isMaintenance` tampak **tanpa kondisi** bila dibaca dari medan *viewer*
> (`pyConditionValue1StringLabel` berisi pola kosong `[first value][relation][second value]`).
> Kalau saya berhenti di situ, saya akan menulis *"rule tanpa kondisi"* — dan itu **salah**.
>
> Seluruh **41 medan kandidat** disisir. Kondisi sebenarnya ada di **`pyConditionString`**:
> `pyWorkPage.Quotation.BusinessType = "Maintenance"`.
>
> ⭐ **Ini gigitan keempat yang BERHASIL DICEGAH**, bukan yang terjadi. Aturannya bekerja.

---

## §C — Gerbang induk `CopyOldataCurr_act` 7 · 8 · 9 · 10

`[terverifikasi]` Catatan pengembang rule: *"set ganti currency other list kalau currency parent
di ganti"*.

### C1 — Apa adanya

**Keempat langkah induk berflag prakondisi `true`** — gerbangnya berlaku. Kode arah dibaca **tidak
terbalik** (benar→kode kiri, salah→kode kanan). Nol baris ber-`//`.

| Langkah | Halaman yang diulang | Baris syarat | benar→ | salah→ |
| --- | --- | --- | --- | --- |
| **7** | `ClaimData.SpreadingRisk` | 1 `CurrencyIDOld!=""` | 2 | 3 |
| | | 2 `panjang daftar > 0` | 2 | 3 |
| | | 3 `CARI1=="Interest"` | ⭐ **5** | 2 |
| | | 4 `CARI1=="ClaimAmount"` | 2 | 3 |
| **8** | `ClaimData.EstimationList` | bentuk **sama persis dengan langkah 7** | | |
| **9** | `ClaimData.SpreadingClaim` | 1 `CurrencyIDOld!=""` | 2 | 3 |
| | | 2 `panjang daftar > 0` | 2 | 3 |
| | | 3 `CARI1=="Interest"` | ⭐ **5** | 2 |
| | | 4 `CARI1=="ClaimAmount"` | ⭐ **5** | 2 |
| | | 5 `CARI1=="Estimation"` | 2 | 3 |
| **10** | `ClaimData.SpreadingBreakQS` | bentuk **sama persis dengan langkah 9** | | |

**Enam baris berkode 5** ada di keempat langkah ini (7→1 · 8→1 · 9→2 · 10→2). Baris **ketujuh**
ada di `Activity/SetMOClaimTreaty_Act.xml` langkah 1 — cocok dengan hitungan ronde 4.

### C2 — Rantai yang berlaku

Kode **5** = *lewati sisa syarat, jalankan langkah* → ia **memutus rantai AND menjadi OR**.

```
langkah 7   CurrencyIDOld != ""  AND  daftar tidak kosong
            AND ( Interest  OR  ClaimAmount )

langkah 8   sama bentuknya dengan langkah 7

langkah 9   CurrencyIDOld != ""  AND  daftar tidak kosong
            AND ( Interest  OR  ClaimAmount  OR  Estimation )

langkah 10  sama bentuknya dengan langkah 9
```

**Pembanding yang menentukan** — langkah **6** (`ListClaimAmount`) **tidak** punya baris berkode 5:

```
langkah 6   CurrencyIDOld != ""  AND  daftar tidak kosong  AND  Interest      (murni AND)
```

### C3 — Ketujuh baris berkode 5 membuat langkahnya **HIDUP**

| | |
| --- | --- |
| **Apa yang kurang** | Bentuk rantai keempat langkah induk **tidak pernah dituliskan** di ronde mana pun. Pembaca yang menurunkannya sendiri dengan aturan AND mendapat `Interest AND ClaimAmount` — **mustahil** — lalu menyimpulkan langkahnya mati dan **membuang rumusnya** |
| **Apa yang sudah cukup** | ✅ Gerbangnya kini **terbaca utuh**: flag prakondisi, syarat, kode arah, dan arahnya. Rantainya **dapat terpenuhi**. Keempat langkah **HIDUP** |
| **Apa yang masih perlu** | ⚠️ Nilai `FlagGantiCurrVal.CARI1` **tidak pernah diisi** di dalam rule ini — pengisinya di luar. Ketiga nilai yang diuji (`Interest`, `ClaimAmount`, `Estimation`) **habis mencakup** semua cabang, jadi rantainya praktis selalu terpenuhi bila dua syarat pertama lolos — tetapi **daftar nilai yang sah belum dapat dibuktikan lengkap** |

⛔ **Butir goyah ini TIDAK saya nyatakan tertutup.** Keputusannya milik work owner.

### C4 — Yang terpengaruh

| Sasaran | Bagaimana |
| --- | --- |
| **tiket 00** *(prefactor, jalur uang)* | rumus penyalin mata uang **tidak boleh dibuang** — langkahnya hidup |
| **tiket 06** *(loss allocation dan spreading)* | tiga dari empat langkah menulis daftar spreading |
| **`spec.md` bab 4** *(Uang — satu jalur perhitungan)* | ✅ ditambal §E3 |

---

## §D — 303 catatan pengembang, sensus tuntas

### D1 — Hitung ulang: **303 catatan di 295 berkas** — cocok

`[terverifikasi]` **303** catatan berisi · **295** berkas · **303** pasangan berkas+catatan
berbeda (nol catatan kembar dalam satu berkas). ⛔ **Tidak berbeda** dari angka ronde 3.

Sebarannya: Activity **108** · When **50** · RDBList **48** · Section **35** ·
ReportDefinition **19** · Harness **17** · DataTransform **9** · FlowAction **9** ·
ConnectREST **6** · DecisionTable **1** · Flow **1**.

### D2 — Keempat golongan

| Golongan | Jumlah | Isinya |
| --- | --- | --- |
| **(a) administratif** | **187** | *"ok"* · *"done"* · *"save"* · *"OK"* · *"."* · *"test"* · nama operator — tidak mengubah apa pun |
| **(b) menerangkan yang sudah tercatat** | **105** | *"Untuk History Claim"* · *"AMBIL KODE PRODUKSI"* · *"tambah auth"* · *"SET FAIL"* · *"edit kondisi"* pada rule `When` — menguatkan, tidak mengubah |
| **(c) ⭐ MENGUBAH pembacaan** | ⭐ **11** | daftar lengkap di bawah |
| **(d) ⚠️ melarang/memerintahkan sesuatu yang MASIH DIKERJAKAN** | ⭐ **0** | lihat D3 |

#### Golongan (c) — sebelas, dengan rule dan langkahnya

| Rule | Langkah | Catatan | Yang berubah |
| --- | --- | --- | --- |
| `CHECKREPORTDATE_ACT` | 11 · 14 | *"matikan protek dalam periode polis"* | proteksi periode polis **mati sengaja** |
| `CHECKDATERECEIVED_ACT` | 10 · 13 | *"change"* | pola sama untuk Date Received |
| `CHECKPERIODPOLICY_ACT` | 1 | *"perbaiki"* | pengisi menggantung, kejadian kedua |
| `ADDADJUSTMENT_ACT` | — | *"tambah…Policy period is TBA…"* | konsep periode polis TBA |
| `SETKOMITENO_ACT` | 1 | *"perbaiki substringnya"* | nomor komite diiris posisi karakter |
| `MAKELOWERCASE_ACT` | 1 | *"make location to temp"* | teks pengguna ditimpa huruf kecil |
| `SERVICEGOOGLE` | — | *"buat jadi 300rb"* | batas waktu 300 000 ms |
| `SETPAYABLETREATY_ACT` | 4 · 5 | *"BUKA PROTEKSI SALVAGE KASIR"* | rule tidak pernah menyimpan |
| `COUNTVALUEADJTREATY_ACT` | — | *"samain dengan prod"* | menyiratkan beda antar-lingkungan |
| `GETMOBYNOPOL_SQL` | — | *"UBAH BIAR GA SINGLE ROW"* | `[terbuka]` — penanda baris-tunggal masih terbaca |
| `INSERTGOOGLESTORAGE_ACT` · `GETURLGOOGLESTORAGE_ACT` | — | *"FIX ERROR HANDLING"* | goyah 1, sudah berdiri |

### D3 — Selisih terhadap ronde 3 dan ronde 4

| | Ronde 3 | Ronde 4 | ⭐ Ronde 6 (tuntas) | Selisih |
| --- | --- | --- | --- | --- |
| Catatan yang **dinilai** | 8 dari 303 | +4 larangan | ⭐ **303 dari 303** | — |
| Golongan **(c)** mengubah pembacaan | **2** | — | ⭐ **11** | **+9** |
| Golongan **(d)** masih dikerjakan | — | **2** | ⭐ **0** | **−2** |

**Penjelasan selisih (c): +9.** Ronde 3 hanya menguji 8 catatan, dan menandainya sendiri sebagai
*"paling rawan salah"* karena basisnya terlalu kecil. ⭐ **Kekhawatiran itu terbukti benar** —
angkanya naik lebih dari lima kali lipat. ⛔ Ronde 3 **tidak keliru**, hanya belum tuntas; ia sudah
menyatakan keterbatasannya sendiri.

**Penjelasan selisih (d): −2, dan sebabnya §A1 ronde 5.** Ronde 4 menandai dua larangan sebagai
*"masih dikerjakan"*: `CountPersen_act` *"jgn kali share ceding lagi"* dan `CheckNoPolicy`
*"jangan pake yg di query"*. ⭐ **Keduanya GUGUR** begitu `//` dipastikan mematikan langkah — kode
yang dilarang itu memang sudah dimatikan, hanya tidak dihapus.

**Keempat larangan langsung lain, diperiksa ulang:**

| Larangan | Masih dikerjakan? | Bukti |
| --- | --- | --- |
| *"hapus save"* — `SETCATASTROPE_ACT` · `SETDEFNONCATASTROPE_ACT` | ⛔ tidak | nol langkah simpan |
| *"hapus yg hapus attachment list"* — `INSERTJSONCLAIMTREATY_ACT` | ⛔ tidak | nol langkah penghapus |
| *"Hapus filter menggunkan text"* — `CAUSEOFLOSS_HARNESS` · `CAUSEOFLOSS_SECTION` | ⛔ tidak | **nol** filter berlabel, **nol** penanda penyaringan di kedua berkas |
| *"ganti IsBonding jadi IsBondingAndCustomBonds"* — `ISANEKA` | ⛔ tidak | rujukan yang berlaku sudah yang baru |

⭐ **Golongan (d) = NOL.** Tidak ada satu pun catatan pengembang yang melarang sesuatu yang masih
benar-benar dikerjakan.

---

## §E — `spec.md` ditambal di **lima titik**

| Butir | Di mana | Isi |
| --- | --- | --- |
| **E1** | *Further Notes → Aturan baca korpus* | RALAT angka ronde 2: **89 → 23** (kode 5 = 7, kode 6 = 16) dan **380 → 363**; *"44 TERBALIK"* **cocok** |
| **E2** | bab baru *Sapu bersih korpus — 2026-09-19 (ronde 6)* | sebelas temuan, enam pasang kembar, lima `When`, pelajaran parser |
| **E3** | bab **4** *(Uang)* | rantai syarat `CopyOldataCurr_act` 7·8·9·10 ditulis; keempatnya **HIDUP** |
| **E4** | *Pertanyaan terbuka → Status sekarang* | butir `[terbuka]` ringan **7 · 8 · 9** ditambahkan |
| **E5** | blok RALAT di *Status sekarang* | jumlah akhir **SEMBILAN**; **AC tidak berubah** |

### Daftar `[terbuka]` sesudah ronde 6

| | |
| --- | --- |
| **HILANG** | ⛔ **NOL** |
| **BERTAMBAH** | **tiga** — butir 7 (`CheckPeriodPolicy_Act` pengisi menggantung) · 8 (periode polis TBA) · 9 (*"samain dengan prod"*) |
| **JUMLAH AKHIR** | ⭐ **SEMBILAN** butir ringan · **yang memblokir: NIHIL** |
| **Jumlah AC** | ⛔ **TIDAK BERUBAH** — tetap **124**. Temuan ronde 6 dicatat sebagai keterangan dan butir terbuka, **bukan** AC baru, jadi sebaran penanda tidak bergeser dan `sensus.py` tidak perlu dijalankan ulang |

---

## §F — Tiket, lalu apa lagi

### F1 — Lima tiket, disebut SEBELUM disunting

| Tiket | Alasan satu baris |
| --- | --- |
| **01** registrasi klaim dan nomor polis | proteksi tanggal terhadap periode polis **dimatikan sengaja** di dua rule, ditambah konsep **periode polis TBA** yang belum pernah tercatat |
| **04** klasifikasi lini bisnis | lima rule `When` terakhir dibaca — **nol hidup**, sapu bersih jalur ini selesai |
| **11** penyerahan komite dan penutupan klaim | nomor komite **diiris dari kunci internal Pega** menurut posisi karakter |
| **13** efek keluar | batas waktu **300 000 ms** pada penyimpanan berkas, dan `SetPayableTreaty_Act` yang **tidak pernah menyimpan** |
| **14** jejak audit | teks bebas pengguna **ditimpa huruf kecil tanpa menyimpan aslinya**, ditambah penulis kronologi ketiga |

⛔ **Tiket 01 ditambal untuk alasan yang BERBEDA dari yang dilarang §A1.** Butir
*"`CheckNoPolicy` menembak pencarian polis tiga kali"* memang **GUGUR** dan **tidak disentuh** —
baris `CheckNoPolicy` di tabel tiket 01 dibiarkan apa adanya, dan catatan lingkup ditulis di dalam
tiketnya supaya tidak tertukar.

⛔ **Tiket 12** *(dokumen PLA/DLA)* **tidak ditambal** meskipun `SetDLACedingSOB` menyentuhnya —
temuannya baru berupa *"ada dua nomor dokumen"* tanpa bukti bagaimana keduanya dipakai. Dicatat di
`spec.md` saja.

### F2 — Bukti tiket lain tidak bergerak

Lihat §F-lampiran.

### F3 — Korpus: **HABIS**

⭐ **Sesudah ronde 6, korpus Claim Prop habis untuk tujuan audit ini.**

| Permukaan | Keadaan |
| --- | --- |
| 329 berkas | ⭐ **nol** yang belum pernah disebut |
| 303 catatan pengembang | ⭐ **303 dari 303** dinilai |
| 1366 langkah activity | seluruhnya tersisir untuk gerbang, remark, perulangan, Java, Obj-* |
| Empat keluarga gerbang · lima keluarga gerbang Section | tersisir |
| Goyah yang bisa dijawab korpus | ⭐ **nol tersisa** |

**Yang tersisa BUKAN pekerjaan korpus:**

| Sisa | Besar | Menunggu |
| --- | --- | --- |
| **9 butir `[terbuka]` ringan** | kecil per butir | **work owner** — nol memblokir |
| **3 butir goyah** | kecil | work owner / riwayat rule Pega — nol dapat dijawab korpus |
| **5 pertanyaan tergantung dari modul Komite Claim Prop** | sedang | work owner |
| Mulai membangun (`/implement` dari tiket **00**) | besar | ⭐ **tidak diblokir apa pun** |

### F4 — Kesimpulan yang sekarang PALING RAWAN SALAH

⛔ **Ditunjuk, tidak diperbaiki.**

> ⚠️ **Paling rawan: §D2 golongan (a) dan (b) — 187 + 105 = 292 catatan yang saya nyatakan
> "tidak mengubah pembacaan".**

Ini **pernyataan ketiadaan atas 292 butir sekaligus**, dan saya menilainya **dari teks catatannya
saja** — bukan dengan membuka tiap rule dan memeriksa apakah catatannya cocok dengan isinya.
Sebelas yang saya golongkan (c) ketahuan justru **karena** saya membuka rule-nya. Kalau 292 sisanya
dibuka satu per satu, sebagian hampir pasti pindah golongan.

⚠️ **Bentuknya sama dengan pola §D2 ronde 5** — hanya arahnya terbalik: bukan *"medan belum
disisir"*, melainkan *"bukti belum dibuka"*. ⭐ Aturan "ketiadaan" seharusnya juga berlaku di sini,
dan **saya belum memenuhinya untuk 292 butir itu**.

**Paling rawan kedua:** §C3 *"keempat langkah HIDUP"*. Ia benar untuk bentuk rantainya, tetapi
bersandar pada `FlagGantiCurrVal.CARI1` yang **pengisinya di luar rule ini dan tidak saya lacak**.

### F5 — Pertanyaan BARU untuk work owner — **2**

> ⛔ Tidak mengulang apa pun yang sudah dijawab, dan **tidak** menyinggung pertanyaan §A yang masih
> menggantung.

#### Pertanyaan 1 — Proteksi periode polis dimatikan sengaja: ditiru, atau dihidupkan lagi?

**Apa yang ditanyakan.** Dua rule pemeriksa tanggal — Report Date dan Date Received — punya
pemeriksaan terhadap **periode polis** yang **dicoret**, dan catatan pengembangnya terang-terangan
berbunyi *"matikan protek dalam periode polis"*. Jadi hari ini klaim bisa masuk dengan tanggal di
luar periode polis tanpa peringatan apa pun.

**Kenapa ini penting.** Aturan induk proyek ini adalah **tiru perilaku, bukan niat**. Menurut aturan
itu, aplikasi Go juga **tidak** memeriksa. Tetapi ini pemeriksaan yang melindungi uang, dan
mematikannya tampak sebagai keputusan sementara — bukan keputusan rancangan.

**Pilihan yang saya lihat.** (a) ditiru apa adanya — Go juga tidak memeriksa; (b) dihidupkan
kembali sebagai **penyimpangan sadar**, dengan alasan tertulis; (c) dihidupkan tetapi hanya sebagai
**peringatan**, bukan penolakan.

**Yang saya sarankan.** **(a) ditiru apa adanya.** Alasannya: ia dimatikan **dengan sengaja** dan
catatannya menyatakan itu — jadi ini perilaku yang dipilih, bukan kerusakan. ⛔ Tetapi ini
keputusan work owner, dan (c) layak dipertimbangkan kalau Finance menginginkan jaringnya.

#### Pertanyaan 2 — Periode polis "TBA": sampai mana klaimnya boleh jalan?

**Apa yang ditanyakan.** Ada properti `.ClaimData.PeriodPolicyTBA` dan pesan *"Policy period is
TBA. Please verify the dates."* di tiga berkas. Artinya klaim dapat didaftarkan ketika periode
polisnya **belum pasti**. Yang tidak terbaca dari korpus: apakah klaim seperti itu boleh berjalan
sampai **akseptasi**, sampai **komite**, atau sampai **pembayaran**.

**Kenapa ini penting.** Kalau boleh sampai pembayaran, maka uang keluar atas polis yang periodenya
belum dipastikan. Kalau ada batas, batas itu **tidak ada di korpus** — dan aplikasi Go akan
membangunnya tanpa batas kecuali diberi tahu.

**Pilihan yang saya lihat.** (a) tidak ada batas — persis seperti sekarang, pesan hanya peringatan;
(b) berhenti sebelum akseptasi; (c) berhenti sebelum pembayaran ke Kasir.

**Yang saya sarankan.** ⚠️ Saya **tidak menyarankan**. Korpus hanya memuat pesannya, **nol
gerbang** yang menghentikan apa pun karena TBA — jadi menebak batas di sini berarti mengarang
aturan bisnis yang tidak ada buktinya.

---

## §F-lampiran — bukti berkas yang tidak bergerak

Sidik jari MD5 atas **23 berkas** diambil **sebelum** penyuntingan dan dibandingkan sesudahnya.

### ✅ BERUBAH — tepat **6** berkas, semuanya disebut brief atau §F1

```
spec.md                                              <-- §E saja (lima sisipan)
issues/01-registrasi-klaim-dan-nomor-polis.md        <-- §F1
issues/04-klasifikasi-lini-bisnis.md                 <-- §F1
issues/11-penyerahan-komite-dan-penutupan-klaim.md   <-- §F1
issues/13-efek-keluar-kasir-arasapas-konversi-email.md <-- §F1
issues/14-jejak-audit.md                             <-- §F1
```

### ⛔ TIDAK BERGERAK — MD5 identik sebelum dan sesudah

```
grilling-ronde-1.md · grilling-ronde-2.md · grilling-ronde-3.md
grilling-ronde-4.md · grilling-ronde-5.md
periksa-ulang-aturan-baru.md
issues/00 · 02 · 03 · 05 · 06 · 07 · 08 · 09 · 10 · 12 · 15   (11 tiket)
```

**17 berkas terbukti tidak bergerak** — 5 berkas ronde, `periksa-ulang-aturan-baru.md`, dan
**11 dari 16 tiket**.

⭐ **Modul lain: NOL disentuh.** Seluruh pembacaan korpus ronde ini terbatas pada
`D:\XML\RNM_BRD\Claim Prop\`, dibuka **baca-saja**.

### Berkas BARU yang ditulis ronde ini — **satu**

```
grilling-ronde-6.md
```
