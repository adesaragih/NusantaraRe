# Grilling Ronde 6 — Komite Claim Prop

**Tanggal:** 2026-09-18 · **Korpus:** `D:\XML\RNM_BRD\Komite Claim Prop\` (READ-ONLY)
**Lingkup:** Komite Claim Prop saja `[keputusan work owner]` 2026-09-18.
**Metode:** parser XML bersarang. Bukti = path berkas + nama rule + nomor step Pega.
**Tidak ada nomor baris XML di berkas ini.**

> **Ralat dicatat DI SINI**, berkas ronde 1–5 tidak disunting.

---

## 0. Ringkasan

| # | Hasil | Status |
| --- | --- | --- |
| **§A1** | Urutan efek keluar | ⭐ **TIDAK DITIRU** — titik yang **sengaja diubah** |
| **§A2** | Hak akses per activity | ✅ diabaikan · **nama privilege terisi NOL** — cocok |
| **§B** | `M_LINK_SERVICE` | ⛔ **RALAT** — saringannya **ADA**. Nol cacat |
| **§B2** | Aturan parser baru | ⭐ **2 tertulis** — dan satu **koreksi letak** |
| **§B3** | `Obj-*` tanpa saringan | ⭐ **NOL cacat nyata** — 5 bersaring, 5 pakai *handle* |
| **§B4** | Kesimpulan lain dari "tidak ketemu" | ⚠️ **3** |
| **§C** | Definisi kolom layar | ✅ **(a) dipilih** — **komite 62 · klaim 31 = 93** |
| **§D** | Berkas `Flow` | **1 berkas · 4 bentuk · jalur balik ADA** |
| **§E2** | Wadah belum terparse | ⚠️ **ADA** — `pyXMLSignature` 18/35 · `pyStepsJavaSource` 20 langkah |

---

## §A — Tiga keputusan work owner `[keputusan work owner]` 2026-09-18

### A1 — Urutan efek keluar terhadap penyimpanan **TIDAK DITIRU**

> ## ⭐ TITIK YANG SENGAJA DIUBAH, BUKAN DISALIN
>
> Kedelapan efek keluar berjalan **sebelum `COMMIT`** di Pega (§B2 ronde 5). Itu dicatat sebagai
> **FAKTA**, tetapi **TIDAK MENGIKAT rancangan**. Urutannya **akan disesuaikan nanti**.
>
> ⛔ **Siapa pun yang menulis spec: jangan mengunci urutan Pega sebagai syarat.** Fakta "Kasir
> dikirim di langkah 34, `COMMIT` di langkah 41" adalah **catatan perilaku lama**, bukan
> **kebutuhan**.
>
> ⛔ **Urutan penggantinya belum diputuskan dan tidak dirancang di sini.**

Butir `[terbuka]` *"seluruh efek keluar mendahului `COMMIT`"* **berubah sifat**: bukan lagi
pertanyaan terbuka, melainkan **titik perubahan yang sudah disadari**.

### A2 — Hak akses per activity **DIABAIKAN**

`pyActivityPrivilegeList` **tidak dipindahkan dan tidak disisir**. Butir `[terbuka]` keluarga
keempat **TUTUP**. Butir pekerjaannya **dicoret** dari daftar §E1.

`[terverifikasi]` **Temuan pendukung, dihitung ulang sendiri:**

| Ukuran | Jumlah |
| --- | --- |
| Baris privilege di 35 Activity | **35** |
| `pyPrivilegeClass` **terisi** | **34** |
| **Nama privilege terisi** | ⭐ **NOL** |

✅ **Angka work owner COCOK PERSIS.** Satu-satunya tag yang pernah terisi adalah
`pyPrivilegeClass`; **tidak ada satu pun nama privilege**. Daftar itu **menyebut kelas tanpa
menyebut hak** — jadi memang tidak menegakkan apa pun.

### A3 — Kelas `ASM-FW-GCNMFW-Work-Komite` ikut **GUGUR**

Ia muncul hanya sebagai **kelas hak akses** (4× di `pyActivityPrivilegeList`), dan karena A2
mengabaikan daftar itu, kemunculannya **bukan temuan kelas kerja baru**. Butir `[terbuka]`-nya
**dicabut**.

⚠️ Catatan: ia juga muncul di `pyPrivilegeList` milik berkas `Flow` (§D1) — **dengan pola yang
sama**: kelas terisi, nama privilege kosong. Konsisten dengan A2.

---

## §B — ⛔ RALAT BESAR: `M_LINK_SERVICE` **ADA** saringannya

> ### ⛔ RONDE 5 §B3 DICABUT
>
> Ronde 5 menulis alamat layanan diambil **"baris pertama tanpa saringan"**, dan menaikkannya jadi
> **Pertanyaan 11** ke work owner. **Itu keliru.** Saringannya **ada**; ia tidak terbaca karena
> letaknya di wadah yang **tidak saya parse**.
>
> ✅ **`pxResults(1)` aman** karena hasilnya sudah disaring menjadi satu baris. **NOL cacat.**
> ✅ Butir `[terbuka]` `M_LINK_SERVICE` **dicabut**, dan **Pertanyaan 11 ditarik**.

### B1 — Saringan dan keempat pemanggil, dibuktikan sendiri `[terverifikasi]`

`Komite Claim Prop/Activity/GetLinkService.xml`, rule
**`ASM-FW-GISFW-INT-M_LINK_SERVICE!GETLINKSERVICE`**, langkah **`Obj-Browse` step 2**:

```
Select  .URL
where   .KATEGORI_1 = Param.Kategori_1
  and   .KATEGORI_2 = Param.Kategori_2
```

`[terverifikasi]` Parameter `Kategori_1` dan `Kategori_2` memang dideklarasikan di
`pyXMLSignature` milik rule-nya — **bukan** di daftar parameter biasa.

**Keempat pemanggil, dihitung ulang sendiri:**

| Activity | Step | `Kategori_1` | `Kategori_2` |
| --- | --- | --- | --- |
| `HitServiceToKasirKMT_Act` | **8** | `"Kasir"` | `"insertAllPaymentKasir"` |
| `KonversiKlaim_Act` | **2** | `"Klaim"` | `"insertClaimAccept"` |
| `InsertGoogleStorage_Act` | **10** | `"Google"` | `"upload"` |
| `GetUrlGoogleStorage_Act` | **6.4** | `"Google"` | `"geturl"` |

✅ **Keempat pasangan COCOK PERSIS** dengan yang disebut work owner.

⭐ Artinya `M_LINK_SERVICE` memang **berisi banyak baris**, satu per layanan, dan dipilih dengan
**dua kunci**. Kekhawatiran ronde 5 sepenuhnya tidak berdasar.

### B2 — ⭐ DUA ATURAN PARSER BARU

**ATURAN P1 — saringan `Obj-Browse` ada di `pyParamArray`, bukan tag tunggal.**
`[terverifikasi]` Satu baris saringan adalah **satu `<rowdata>`** yang membawa
`<Select>` · `<Field>` · `<Condition>` · `<Value>` bersama-sama.

⚠️ **Koreksi letak terhadap brief:** brief menyebut baris itu ada *"di dalam callparams
langkahnya"*. **Jalur sebenarnya berbeda** — ia di **`pyParamArray/rowdata`**, wadah yang **sama**
dengan tempat `Property-Set` menyimpan `PropertiesName`/`PropertiesValue`:

```
/pySteps/rowdata/pyParamArray/rowdata/Field
/pySteps/rowdata/pyParamArray/rowdata/Condition
/pySteps/rowdata/pyParamArray/rowdata/Value
/pySteps/rowdata/pyParamArray/rowdata/Select
```

`pyStepsCallParams` hanya memuat `<Logic>` yang **kosong**. ⭐ **`pyParamArray` bermuka dua:** untuk
`Property-Set` ia daftar penugasan, untuk `Obj-Browse` ia daftar kolom dan saringan. **Membacanya
sebagai satu jenis saja adalah cara paling mudah kehilangan separuh isinya** — dan itulah yang
terjadi di ronde 5.

**Membedakan kolom-pilih dari saringan:** baris ber-`Condition` **dan** `Value` terisi = **saringan**;
baris ber-`Select=true` tanpa keduanya = **kolom yang diambil**.

**ATURAN P2 — parameter activity ada di `pyXMLSignature`, sebagai XML TER-ESCAPE di dalam teks.**
`[terverifikasi]` Isinya bukan elemen XML biasa melainkan **dokumen XML yang di-escape** ke dalam
satu simpul teks:

```
<pyXMLSignature>&lt;?xml version="1.0"?&gt;&lt;methodSignature&gt;&lt;pyStepsCallParams&gt;
  &lt;Kategori_1 DESCRIPTION="" TYPE="STRING" REQUIRED="-1" INOUT="IN" … /&gt;
  &lt;Kategori_2 DESCRIPTION="" TYPE="STRING" INOUT="IN" … /&gt;
&lt;/pyStepsCallParams&gt;&lt;/methodSignature&gt;</pyXMLSignature>
```

⛔ **Parser XML biasa membacanya sebagai teks polos.** Ia harus **di-unescape lalu diurai ulang**.

### B3 — Sisir ulang seluruh `Obj-Browse` / `Obj-Open` / `Obj-List` `[terverifikasi]`

**10 langkah** di 35 Activity modul ini.

**Lima ber-`Obj-Browse`, SELURUHNYA punya saringan:**

| Activity · step | Kelas | Saringan |
| --- | --- | --- |
| `GetBase64Attachment` **3** | `Int-DOCUMENT_CLAIM` | `.IDPEGA = Local.Inskey` + `.KATEGORI_1` = `"Premium"` / `"DLA"` / `"AcceptanceNote"` / `"DLARetro"` |
| `GetLinkService` **2** | `Int-M_LINK_SERVICE` | `.KATEGORI_1` · `.KATEGORI_2` |
| `HitServiceToKasirKMT_Act` **12** | `Int-BANKACCOUNT` | `.NAMEOFBANK` · `.BRANCHOFBANK` · `.ACCOUNTNO` |
| `SendEmailKlaim_KMT` **13** | `Data-Admin-Operator-ID` | `.pyUserIdentifier = pyWorkPage.pxCreateOperator` |
| `SendErrorDirectKasir` **2.1** | `Data-Admin-Operator-ID` | `.pyUserIdentifier = pyWorkPage.pxCreateOperator` |

**Lima ber-`Obj-Open-By-Handle`, dan itu METODE LAIN** — ia membuka **satu instance lewat kunci**,
jadi **tidak perlu saringan**. Kelimanya memang mengirim kunci:

| Activity · step | `InstanceHandle` |
| --- | --- |
| `KomitePostAdjustment` **4** | `pyWorkPage.pxCoverInsKey` |
| `KomitePost_Close` **1** | `pyWorkPage.pxCoverInsKey` |
| `KomitePost_Reject` **1** | `pyWorkPage.pxCoverInsKey` |
| `PrintFileAcceptance_TKMT` **1** | `Param.IDPega` |
| `GetBase64Attachment` **5.6** | `Local.Key` |

⭐ **Jawaban B3: NOL yang benar-benar tanpa saringan. Nol calon cacat.**

⭐ Temuan sampingan: ketiga `KomitePost*` membuka **kasus klaim induk** lewat
`pyWorkPage.pxCoverInsKey` **dengan `Lock = true`** — mengunci kasus klaim selama komite memproses.

### B4 — Kesimpulan lain yang lahir dari "tidak ketemu" — ⚠️ **3**

⛔ **Didaftarkan, tidak diperbaiki.**

| # | Di mana | Kesimpulannya | Kenapa rawan |
| --- | --- | --- | --- |
| **1** | **ronde 5 §B4** | *"Empat efek keluar tidak punya penanganan gagal"* — PDF, Google Storage, email | Ditarik dari **dua keluarga** (`PreCondParams`, `TransParams`). **`pyStepsJavaSource` belum dibaca** (20 langkah), dan sebuah langkah `Java` bisa menangani gagal sendiri. **Pola yang sama persis dengan ralat §B ini.** |
| **2** | **ronde 2 §A2 Aturan 1** *(sudah diralat ronde 3)* | *"Sel tidak pernah bergerbang sendiri"* | Sudah terbukti salah — `pyUserData` tidak dicari. **Dicatat lagi di sini karena ia leluhur pola yang sama.** |
| **3** | **ronde 3 §D** | *"Pewarisan kelas tidak ada di jendela ini"* | Ditarik dari ketiadaan berkas `Rule-Obj-Class`. Tetapi `pyXMLSignature` membuktikan **isi bisa bersembunyi di dalam teks**. Sebelum F6 dinyatakan mustahil, wadah teks di berkas kelas perlu dicek. |

⭐ **Polanya satu dan sama: saya menyatakan "tidak ada" berdasarkan wadah yang saya tahu, bukan
wadah yang ada.** Dua dari tiga sudah terbukti salah.

---

## §C — Definisi "kolom layar", **dikunci**

### C1 — Definisi terpilih: **(a) FIELD YANG TAMPIL DI LAYAR**

Yaitu **`pyValue` di dalam sel `Embed-Display-Table-Cell`**, dan hanya itu.

**Alasannya.** Pertanyaan yang perlu dijawab adalah **"kolom apa yang dilihat dan diisi
pengguna"**. Hanya sel `Embed-Display-Table-Cell` yang **benar-benar menjadi kotak di layar**;
`pyValue`-nya adalah properti yang **isinya ditampilkan atau diketik**.

**Kenapa (b) tidak menjawab itu.** Definisi (b) menghitung **setiap penyebutan** properti, termasuk:
sebuah properti yang hanya muncul **di dalam teks syarat gerbang** (`.AcceptStatus = 1 &&
.TransferType =2`) — ia **menentukan** apakah kotak tampil, tetapi **bukan kotak**; sebuah properti
yang muncul di `pyLabelFor` atau `pyPrompt` — itu **label**, bukan medan; dan properti internal
Pega. Menghitungnya sebagai kolom akan **melebih-lebihkan layar** dan memasukkan hal yang tidak
pernah dilihat siapa pun.

⛔ **Definisi (b) tetap berguna** untuk menjawab pertanyaan lain — *"properti apa saja yang
disentuh berkas ini"* — dan angka 117 ronde 5 **tidak dicabut**; ia hanya **menjawab pertanyaan
yang berbeda**.

### C2 — Angka final `[terverifikasi]`

| Ember | Properti berbeda | Sel |
| --- | --- | --- |
| **KOMITE** — `.Xxx` + `pyWorkPage.Xxx` | **62** | 91 |
| **KLAIM induk** — `pyWorkCover.Xxx` | **31** | 32 |
| **TOTAL** | **93** | 123 |

*(dari 348 sel `Embed-Display-Table-Cell` seluruhnya)*

⚠️ **Angka saya 93, bukan 98.** Ember **KLAIM cocok persis: 31**. Selisihnya **seluruhnya di ember
KOMITE** — 62 lawan 67 — dan sebabnya **penyaringan yang saya sebutkan di C3**. **Saya tulis apa
adanya dan tidak saya kejar agar cocok.** Siapa pun yang menghitung ulang harus menyebut
saringannya.

**31 properti klaim yang benar-benar tampil** — data yang dibaca penyetuju:
`NoClaim` · `InsuredName` · `PolicyNo` · `PolicyData.PolicyNo` · `PolicyData.StartDateTime` ·
`PolicyData.EndDateTime` · `DateOfLoss` · `DateReceived` · `CauseOfLoss` · `Location` ·
`ReportDate` · `ReportType` · `ReportAddress` · `ReportDescription` · `ReporterName` ·
`ReporterStatus` · `ReporterTelp` · `InsuredRelationshipOthers` · `AppointedADJ` ·
`ConsultantName` · `StsKatastrofe` · `KatastrofeNote` · `NonKatastrofeType` · `NoPla` ·
`TotalEstimasi` · `TotalEstimasiIDR` · `TotalGrossEstimateIDR` · `TotalGrossEstimateTreaty` ·
`OfferFacIn.QuotationData.BusinessName` · `TreatyInMaster.Ceding` ·
`TreatyInMaster.LeadingReinsSource`.

**Delapan properti komite berawalan `pyWorkPage.` yang tampil:**
`CLMNO` · `pxCreateDateTime` · `pxCreateOpName` · `Komite.AdjusterFee` ·
`Komite.CircumtansesCouseOfLoss` · `Komite.Occupation` · `Komite.Remarks` · `Komite.Salvage`.

⭐ **Enam dari delapan itu adalah `pyWorkPage.Komite.*`** — yang di ronde 5 baru dicatat sebagai
temuan. Sekarang terbukti **benar-benar tampil di layar**, bukan sekadar disebut.

### C3 — Yang dibuang dari hitungan, supaya selisihnya terlacak `[terverifikasi]`

| Yang dibuang | Jumlah sel | Alasan |
| --- | --- | --- |
| **Teks label** di `pyValue` | **82** | isinya kalimat, bukan properti — *"Accepted Date"*, *"Claim Amount"*, *"Are you sure to accept this document?"* |
| **Properti internal Pega** | **3** | `.pyTemplateInputBox` · `.pyTemplateButton` ×2 — kerangka, bukan data |
| Penyebutan di **teks syarat gerbang** | — | menentukan tampil/tidak, bukan kotak |
| `pyLabelFor` **41** · `pyPrompt` **7** · `pyDisabledWhen` **4** | — | label dan penanda, bukan medan |

⛔ **Sampai keputusan ini dicabut, daftar properti TIDAK dipakai menetapkan kolom tabel.**

---

## §D — Berkas `Flow`: daur hidup kasus komite

### D1 — Berkasnya **satu** `[terverifikasi]`

| | |
| --- | --- |
| Berkas | `Flow/KomiteTreaty_Flow.xml` |
| `pxInsName` | **`ASM-FW-GCNMFW-WORK-KOMITETREATY!KOMITETREATY_FLOW`** |
| Kelas | `ASM-FW-GCNMFW-Work-KomiteTreaty` |
| `pyFlowType` | `KomiteTreaty_Flow` · `pyCategory` `FlowStandard` |

### D2 — Tahapnya **empat bentuk**, dalam bahasa biasa `[terverifikasi]`

| Bentuk | Jenis | Nama | Siapa yang mengerjakan |
| --- | --- | --- | --- |
| **`Start1`** | Start | — | sistem, saat kasus komite dibuat *(harness `NewSample`)* |
| **`ASSIGNMENT63`** | **Assignment** | **`KomiteRouter`** | ⭐ **penyetuju komite** — masuk ke *worklist*-nya |
| **`Decision1`** | **Gateway keputusan** | **`KomiteLoop`** | sistem — memilih jalur |
| **`END52`** | End | — | sistem · status akhir **`Resolved-Completed`** |

**Yang memindahkan antar tahap — empat penghubung:**

| Penghubung | Dari → ke | Pemicu |
| --- | --- | --- |
| `Transition1` | `Start1` → `ASSIGNMENT63` | **`[Always]`** — langsung |
| **`TRANSITION54`** | `ASSIGNMENT63` → `Decision1` | ⭐ **aksi pengguna: `ViewTransferDtl`** |
| **`Transition2`** | `Decision1` → **`ASSIGNMENT63`** | ⭐ **when `IsKomiteLoop`** — **JALUR BALIK** |
| `Transition3` | `Decision1` → `END52` | **`[Else]`**, bernama **`NoLoop`** |

Ceritanya, berurutan: **kasus komite dibuat → langsung masuk kotak kerja penyetuju
(`KomiteRouter` yang menentukan ke siapa) → penyetuju membuka layar dan menekan kirim
(`ViewTransferDtl`) → sistem memeriksa `IsKomiteLoop` → bila masih ada penyetuju berikutnya, kasus
KEMBALI ke kotak kerja; bila tidak, kasus SELESAI dengan status `Resolved-Completed`.**

⭐ **Penemuan baru:** ada **satu *ticket* bernama `komiteAccept_ticket`** (`Ticket1`) di dalam flow.
*Ticket* di Pega adalah **titik lompat darurat** yang bisa dipicu dari luar alur. ⛔ **Siapa yang
memicunya tidak ditelusuri.** `[terbuka]` **BARU**.

### D3 — Di mana `ViewTransferDtl` dan `KomitePost` `[terverifikasi]`

- **`ViewTransferDtl`** bukan tahap — ia **flow action pada penghubung `TRANSITION54`**, yaitu
  **aksi yang dikerjakan penyetuju saat berada di `ASSIGNMENT63`**. Menyelesaikan aksi itulah yang
  memindahkan kasus ke `Decision1`.
- **`KomitePost`** **tidak disebut di berkas `Flow` sama sekali**. Ia dipanggil sebagai
  **pasca-proses flow action `ViewTransferDtl`** (§B2 ronde 2). Jadi urutannya:
  **penyetuju menekan kirim → `KomitePost` berjalan → baru kasus berpindah.**

⚠️ Konsisten dengan §A1: **seluruh efek keluar `KomitePost*` terjadi SEBELUM kasus berpindah
tahap** — dan sebelum `COMMIT`.

`Flow` juga merujuk `pyTransferAssignment`, `pyCreateAdhocCase`, `EngageExternal` — **flow action
bawaan Pega**, dan `pyLocalActionsString = (3)` menandai **tiga aksi lokal** pada assignment itu.
⛔ Ketiganya **tidak ditelusuri**. `[terbuka]` **BARU**.

### D4 — Jalur balik: ⭐ **ADA**

**Dari `Decision1` kembali ke `ASSIGNMENT63`**, lewat penghubung **`Transition2`**, bersyarat
**`when IsKomiteLoop`** — rule `ASM-FW-GCNMFW-Work-KomiteTreaty!ISKOMITELOOP`.

✅ Cocok dengan **§C2 ronde 1**, yang sudah menutup pertanyaan penghubung dua baris `IsKomiteLoop`
dengan jawaban **`AND`**.

**Hanya satu jalur balik.** Tidak ada jalur balik lain — tidak ada yang kembali ke `Start1`, dan
tidak ada yang melompati `Decision1`.

### D5 — Tangga penyetuju: ⭐ **berputar DI DALAM satu tahap**

**Bukan satu tahap per penyetuju.** Buktinya menumpuk:

1. Flow hanya punya **satu assignment** — `ASSIGNMENT63`.
2. Jalur balik **menuju assignment yang sama**, bukan ke assignment berikutnya.
3. `KomiteRouter` adalah **router** assignment itu (`pyRouteToType = Custom`), dan ronde 1 §5 sudah
   menunjukkan ia menugaskan ke `param.AssignTo := .KomiteID` — **penyetuju yang sedang giliran**.
4. `KomiteCount` / `KomiteLoop` **tidak pernah tampil di layar** (§B4 ronde 2) — keduanya
   **keadaan dalam** yang menghitung putaran.

⭐ **Jadi: setiap penyetuju adalah KUNJUNGAN BARU ke assignment yang sama**, bukan tahap baru.
Yang berubah tiap putaran hanyalah **kepada siapa assignment itu dirutekan** dan **nilai
`KomiteCount`**.

⛔ **Tidak disimpulkan rancangan status untuk Go.** Yang di atas adalah perilaku Pega apa adanya.

---

## §E — Apa lagi yang seharusnya dikerjakan

### E1 — Urutan sesudah ronde 6

> ⛔ **Butir hak akses DICORET** — §A2.

| # | Yang dikerjakan | Kenapa perlu | Besar | Menunggu |
| --- | --- | --- | --- | --- |
| **1** | **20 langkah `Java`** — `pyStepsJavaSource` | ⭐ **naik ke puncak karena §B4 no.1**: kesimpulan *"empat efek keluar tak punya penanganan gagal"* bersandar pada wadah yang belum dibaca. Java bisa menangani gagal sendiri | sedang | korpus |
| **2** | **Keluarga ketiga `pyStepsRepeatDef`** — 85 langkah berulang | perulangan menentukan **berapa kali** langkah jalan | sedang | korpus |
| **3** | **`ViewDetailInterest`** — Section 115 KB + FlowAction | **layar kedua**, belum pernah dibuka | sedang | korpus |
| **4** | **Tiga aksi lokal + ticket `komiteAccept_ticket`** | keduanya **jalur masuk/keluar alur** yang belum terpetakan | sekali sisir | korpus |
| **5** | **12 activity penghitung + jalur konversi** | mengisi angka uang di layar | sedang | korpus |
| **6** | **`pyActionSets` tombol + `pyWorkPage.Edit`** | `.pyTemplateButton` 2× tak tertelusur | sekali sisir | korpus |
| **7** | **Sisa 14 berkas** — 8 RDBList · 3 ReportDefinition · 2 DataTransform · 1 DecisionTable | melengkapi peta | sekali sisir | korpus |

⭐ **Berkas `Flow` turun dari daftar — selesai ronde ini.**

### E2 — Wadah yang **belum pernah diparse**: ⚠️ **ADA**

⛔ **Dicari, tidak ditunggu.**

| Wadah | Sebaran | Keadaan |
| --- | --- | --- |
| **`pyXMLSignature`** | **35 activity**, **18** benar-benar mendeklarasikan parameter | ⭐ **baru terbaca ronde ini** |
| **`pyStepsJavaSource`** | **20 langkah**, 8 di antaranya bermarkup | ⚠️ **BELUM** — hanya 1 yang pernah dibaca (§C2 ronde 4) |
| `pyMemo` | 35 | ⚠️ BELUM — catatan penulis rule |
| `pyUsage` | 212 | ⚠️ BELUM |

**Parameter yang baru terbaca lewat `pyXMLSignature`** — 18 activity, contoh yang paling berbicara:

| Activity | Parameter |
| --- | --- |
| `InsertDocument_Act` | `IDPEGA` · `NAMAFILE` · `MIME` · `KATEGORI_1` · `KATEGORI_2` · **`NOAKSEP`** · **`NOPREKAS`** · **`PAYMENTDATE`** · `BASE64` |
| `InsertLogServiceClaim` | `IDPega` · `ParamInsert` · `JenisService` · `NoAkseptasi` · `NoDLA` · `StsMessage` · `ResponMessage` |
| `KonversiKlaim_Act` | `INSKEY` · `NOPOLIS` · `STSREJECT` |
| `GetUrlGoogleStorage_Act` | `Durasi` · `ImageID` · `Url` |
| `InsertGoogleStorage_Act` | `Folder` · `Namafile` · `Image` · `Durasi` · `Ext` · `ImageID` |
| `PrintFileAcceptance_TKMT` | `idxAdj` · `IDPega` |
| `getStatusKonversi_Act` | `STS_KONVERSI` |

⭐ **Ini melengkapi §B ronde 5**: isi tiap efek keluar sekarang punya **daftar masukan yang
tercatat**, yang sebelumnya tidak terbaca sama sekali.

### E3 — Kesimpulan yang **paling rawan salah** — **2**

1. ⛔ **Ronde 5 §B4 — "empat efek keluar tidak punya penanganan gagal"**. **Paling rawan**, dan
   §B4 no.1 di atas menjelaskan kenapa: ia kesimpulan **ketiadaan** dari wadah yang belum lengkap.
   **20 langkah `Java` belum dibaca.** Sudah dua kali pola ini membakar saya.
2. ⚠️ **Ronde 4 §A2 — "keluarga kedua diuji SESUDAH langkah jalan"**. Masih belum teruji silang;
   seluruh pembacaan penanganan gagal bersandar padanya.

⛔ **Ditunjuk, tidak diperbaiki.**

### E4 — Pertanyaan **BARU** untuk work owner — **2**

> ⚠️ **Nol pertanyaan lama yang menggantung.** Pertanyaan 11 ronde 5 **saya tarik sendiri** —
> §B membuktikan ia lahir dari kesalahan baca saya, bukan dari korpus.

#### Pertanyaan 13 — Kasus komite mengunci kasus klaim induknya selama diproses. Dipertahankan?

**Apa yang ditanyakan.** Ketiga jalur komite — adjustment, tolak, tutup — membuka **kasus klaim
induk** dan **menguncinya**, lalu baru memproses. Selama komite bekerja, klaim itu **tidak bisa
disentuh proses lain**. Apakah penguncian itu dipertahankan.

**Kenapa muncul.** Baru terbaca ronde ini. Selama enam ronde, hubungan komite–klaim hanya terlihat
sebagai pembacaan data; ternyata ada **penguncian**.

**Bedanya jawaban A atau B.** Bila **dipertahankan**, dua orang tidak bisa menggarap klaim yang
sama bersamaan — aman, tetapi bisa memblokir. Bila **tidak**, perlu ditetapkan apa yang terjadi
ketika klaim berubah di tengah komite menimbang.

**Yang tertahan.** Perilaku bersamaan antara modul klaim dan modul komite.

#### Pertanyaan 14 — Ada "titik lompat darurat" bernama `komiteAccept_ticket`. Masih dipakai?

**Apa yang ditanyakan.** Di dalam alur komite ada satu **titik lompat** yang bisa dipicu dari luar
alur, artinya kasus bisa **dipaksa berpindah** tanpa melewati jalur biasa. Apakah jalur itu masih
dipakai, dan siapa yang boleh memicunya.

**Kenapa muncul.** Baru terbaca ronde ini. Ia satu-satunya jalan masuk ke alur komite yang
**bukan** lewat persetujuan penyetuju.

**Bedanya jawaban A atau B.** Bila **masih dipakai**, ia jalur sah yang perlu ikut dipindahkan, dan
perlu diketahui siapa pemicunya. Bila **sudah tidak**, ia jalur mati yang dibuang seperti lima
lompatan menggantung di §A3 ronde 5.

**Yang tertahan.** Kelengkapan daur hidup kasus komite — selama belum dijawab, alur yang tercatat
belum tentu seluruh jalannya.

---

## §F — Daftar `[terbuka]` modul ini sesudah ronde 6 — **10**

| # | Butir | Menunggu |
| --- | --- | --- |
| 1 | `OPERATORID` `HISTORYAKSEPTASIPEGA` tak pernah diisi jalur komite | korpus / DBA |
| 2 | **F6** — pewarisan kelas kerja · ⚠️ **cek dulu wadah teks** (§B4 no.3) | work owner |
| 3 | beda `CONDITION` vs `ExpressionCondition` tidak terbaca | korpus |
| 4 | keluarga ketiga `pyStepsRepeatDef`, 85 langkah berulang | korpus |
| 5 | beda `REPEAT` vs `EMBEDDED`, dan halaman apa yang diulang | korpus |
| 6 | urutan evaluasi bila kedua keluarga terisi | korpus / work owner |
| 7 | ⭐ **BARU** — **20 langkah `Java`** belum dibaca; menggantung kesimpulan §B4 ronde 5 | korpus |
| 8 | ⭐ **BARU** — ticket **`komiteAccept_ticket`**: siapa pemicunya | work owner *(pertanyaan 14)* |
| 9 | ⭐ **BARU** — **tiga aksi lokal** pada assignment `KomiteRouter` | korpus |
| 10 | ⭐ **BARU** — penguncian kasus klaim induk oleh jalur komite | work owner *(pertanyaan 13)* |

**Keluar dari daftar ronde 5 — lima butir:** keluarga keempat *(A2)* · `Work-Komite` sebagai kelas
hak akses *(A3)* · `M_LINK_SERVICE` *(B)* · efek keluar sebelum `COMMIT` *(A1 — berubah sifat
menjadi titik perubahan)* · `pyWorkPage.Komite.*` *(terjawab §C2 — enam dari tujuh tampil di layar)*.

⛔ **Tidak ada butir lain dinyatakan tertutup.**

---

## §G — Ralat, dikumpulkan

| # | Di mana | Yang salah | Yang benar |
| --- | --- | --- | --- |
| 1 | **ronde 5 §B3** | alamat layanan *"baris pertama tanpa saringan"*, dinaikkan jadi Pertanyaan 11 | **saringannya ADA** — `.KATEGORI_1` dan `.KATEGORI_2`. **Nol cacat.** Pertanyaan 11 **ditarik** (§B1) |
| 2 | **ronde 5 §C1** | 117 disebut sebagai *"kolom layar"* | **93** adalah kolom layar (definisi a); **117 menjawab pertanyaan lain** dan tidak dicabut (§C) |
| 3 | brief ronde 6 | saringan `Obj-Browse` *"di dalam callparams langkahnya"* | jalurnya **`pyParamArray/rowdata`** (§B2) |
| 4 | ronde 5 §E2 | `pyActivityPrivilegeList` diangkat sebagai keluarga keempat yang **berperilaku** | **nama privilege NOL** — ia tidak menegakkan apa pun (§A2) |

⛔ **Tidak satu pun berkas ronde 1–5 disunting.**
