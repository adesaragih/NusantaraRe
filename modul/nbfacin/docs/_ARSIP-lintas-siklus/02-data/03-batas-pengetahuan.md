# Batas Pengetahuan — apa yang **tidak** dapat dijawab dari korpus

> Sumber: konsolidasi seksi "Pertanyaan terbuka" dari seluruh dokumen discovery di
> `D:\migrasi\RNM\OUTPUT\`, yang dibangun dari korpus
> `D:\migrasi\RNM\{NB FacIn, RNW Fac In, Endorsment Fac In}\` (6.071 berkas `.xml`).
> Label mengikuti `CLAUDE.md` §3.

**Dokumen ini sengaja mendaftar kekurangan.** Mengosongkannya membuat kekurangan itu tidak terlihat,
dan aplikasi yang berjalan dengan tebakan diam-diam jauh lebih berbahaya daripada aplikasi yang
berhenti (`CLAUDE.md` §4.5).

**Cara membaca:** ⛔ = memblokir implementasi · ⚠️ = harus diputuskan bisnis, tidak memblokir ·
Pemilik = siapa yang dapat menjawab.

---

## 0. Ringkasan

| Kategori | ⛔ Memblokir | ⚠️ Perlu keputusan |
| --- | ---: | ---: |
| Arti enumerasi (§1) | 7 | 3 |
| Data yang tidak ada di korpus (§2) | 12 | 1 |
| Integritas ekspor & rule hilang (§3) | 5 | 2 |
| Aritmetika & representasi uang (§4) | 8 | 9 |
| Kontradiksi internal sistem lama (§5) | 9 | 12 |
| Authz, keamanan & guard identitas (§6) | 4 | 2 |
| **Total** | **45** | **29** |

Dua dokumen memuat pertanyaan tidak-memblokir tambahan yang tidak diulang di sini:
`02-data/02-skema-oracle.md` §13 (**13 butir** — antara lain arah tanda kolom `_SELISIH`, master
ganda: okupasi ×3, mata uang ×3, coverage ×3, klausula ×4) dan `01-flow/03-alur-endorsement.md` §10
(**17 butir**).

⚠️ **Butir yang paling mengubah rencana adalah R0 di §3** — selisih versi ekspor antar folder membuat
sebagian angka "percabangan antar siklus" di seluruh discovery ini tidak dapat dipisahkan dari
artefak ekspor. Baca itu sebelum memakai angka mana pun untuk keputusan.

Perintah audit untuk angka korpus yang dikutip dokumen ini ada di dokumen asalnya masing-masing;
setiap butir menyebut lokasinya.

---

## 1. Arti enumerasi — nilai terbaca, artinya tidak

`[terverifikasi]` Nilai literal di bawah **terbaca langsung** dari korpus. Artinya **tidak**.
Menebaknya melanggar `CLAUDE.md` §3 butir 4.

| # | Properti | Nilai literal yang terbaca | Blokir | Pemilik |
| --- | --- | --- | :-: | --- |
| **E1** | `ProposalAcceptStatus` | `1`,`2`,`3`,`4`,`7`,`9` — **`5`,`6`,`8` tidak pernah muncul** | ⛔ | Underwriting |
| **E2** | `QuotationData.Type` | `0`–`9`, `11`, `12` — **`10` tidak dipakai** | ⛔ | Product |
| **E3** | `QuotationData.EdmType` | `4` diuji 37 baris di rule `When`; nilai `1`,`2`,`3` muncul di Activity dengan **deskripsi yang saling bertentangan** — satu rule menyebut `EdmType=1` sebagai "batal sejak semula", rule lain memberi label itu pada `EdmType==2`. `EdmType==3` tidak dijelaskan tetapi **membalik tanda seluruh selisih persentase** | ⛔ | Product |
| **E4** | `BusinessCode` | **98 kode** tanpa tabel referensi | ⛔ | Product |
| **E5** | `BusinessOldId` | **87 kode** (`L1`–`L16`, `A1`–`A9`, …, numerik `22`–`99`) | ⛔ | Product |
| **E6** | `DeductibleType` | **dua ruang nilai tidak beririsan**: 7 kode 6-digit (`100413`…`100420`) vs 6 kode 1-digit (`1`…`8`), diuji terhadap properti yang sama | ⛔ | Product |
| E7 | `JN_SERVICE` | `"FACIN"` / `"FACOUT"` | ⚠️ | IT |
| E8 | `STS_KONVERSI`, `STS_KONVERSI_RETRO` | `1`, `8`, `9` | ⚠️ | IT |
| E9 | `TeamGroup` | `1`–`5` diuji; daftar lengkap tidak diketahui | ⚠️ | Underwriting |

### E1 memblokir paling luas

`[terverifikasi]` `DecisionTable/IsUWAccepted` mendeklarasikan **enam hasil** — `confirm` · `reject` ·
`ask` · `banding` · `revise` · `decline` — dan Data Transform menulis nilai `1`/`2`/`3`/`4`/`7`/`9`.
**Pemetaan antara keduanya tidak ada di ekspor**: baris hasil DecisionTable tidak ikut terekspor.

⚠️ **Korelasi nama rule di sini terbantah langsung:** nilai `4` ditulis `SetBandingProposal_DT`
(banding) tetapi **diuji** `When/IsFacout.xml` (fac out) — `RNW Fac In\When\IsFacout.xml:148` →
`pyWorkPage.ProposalAcceptStatus = 4`. Satu nilai, dua nama pemakai bertema berbeda.

Menebak akan menyalahkan arah **setiap** keputusan underwriting di seluruh node persetujuan.

⛔ **E10 — baris hasil 12 DecisionTable tidak terekspor**, termasuk `IsUWAccepted` (menentukan
penerimaan UW) dan `isApproved` (menentukan persetujuan treaty). Pemilik: Underwriting + IT (ekspor
ulang dari Pega dapat menyelesaikannya).

`[terverifikasi]` `IsUWAccepted` **berbeda isinya** antara siklus endorsement dan NB/RNW — aturan
persetujuan memang bercabang antar siklus, dan **apa** yang bercabang tidak terbaca.

### Celah enumerasi adalah temuan, bukan kebetulan

`[terverifikasi]` (`04-aturan/02-formula-dan-status.md` §3, T16) Nilai yang **tidak pernah muncul**:
`Type=10` · `StatusBusiness` 0 dan ≥4 · `EdmType` 0 dan ≥5 · `EmailType*` 0, 5, 6, 8 ·
`ProposalAcceptStatus` 0 dan 5–8 · `ProRateType` 0 · `metodeKalkulasi` 0 dan ≥3.

Korpus tidak memisahkan "nilai itu tidak ada" dari "nilai itu ada tetapi tidak pernah diperiksa".

---

## 2. Data yang tidak ada di korpus

| # | Hal | Konsekuensi | Blokir | Pemilik |
| --- | --- | --- | :-: | --- |
| **D1** | Isi tabel limit `M_LIMIT_*` | Nilai `LIMIT_BOTTOM`, `LIMIT_BOTTOM2`, `MAX_LIMIT_IDR/USD`, `BATAS_WAKTU`, ejaan pasti kolom `JABATAN` — **tangga akseptasi tidak dapat direkonsiliasi tanpa ini** | ⛔ | DBA |
| **D2** | Isi `M_LINK_SERVICE` (kunci `KATEGORI_1`+`KATEGORI_2`) | Seluruh alamat servis luar | ⛔ | IT |
| **D3** | Isi **24 stored procedure `POOLDATA`** (+ 7 fungsi) | **Seluruh jalur tulis melewatinya.** Nama terbaca, isi tidak — mustahil direplikasi tanpa `ALL_SOURCE` | ⛔ | DBA |
| **D4** | Query di balik `OutData.pxResults(1).CARI2` | 6 rule klasifikasi lini bisnis EDM (`IsAneka`, `IsFire`, `IsGolfInsurance`, `IsMarineCargo`, `IsMBU`, `IsPA`) membacanya; **klasifikasi lini bisnis endorsement tidak dapat diimplementasikan** | ⛔ | IT |
| **D5** | Asal `LimitAkseptasi.pxResults(n).CARI10` | Dibandingkan dengan `PositionNote` untuk menentukan approver berikutnya, tetapi **tidak pernah ditulis di korpus** dan tidak dialiaskan di SQL mana pun (SQL hanya CARI1–CARI5) | ⛔ | IT |
| **D6** | Sumber pengganti `DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK` | **Dua query membacanya untuk routing**, bukan sekadar audit. Tabel internal Pega — hilang bersama Pega. Data historis harus **diekstraksi sebelum dekomisioning**, dan itu **tidak dapat dibatalkan** | ⛔ | IT + DBA |
| **D7** | Ejaan nilai kolom `M_LIMIT_*.JABATAN` di Oracle | Ruang nama `Reas*` (seperti `PositionNote`) atau HURUF BESAR (seperti `LetterNo`)? Tidak ada satu baris data pun di korpus | ⛔ | DBA |
| **D8** | Nilai `getDataSystemSetting("PegaCRM-","SellingMode")` | 3 rule `isSellingMode*` tidak dapat dievaluasi | ⛔ | IT |
| **D10** | **Tipe kolom seluruh tabel** — korpus memuat **nol DDL** | Terutama `FACINOFFER`: bila `VARCHAR2`, seluruh perbandingan angka atasnya adalah perbandingan **string** | ⛔ | DBA |
| **D11** | Tipe dan skema `JSON_POLIS.DATA_JSON` | Hanya 11 jalur terbaca dari SQL; dokumen sesungguhnya tidak terdokumentasi — padahal agregat case diserialkan ke sini | ⛔ | DBA + IT |
| **D12** | Determinisme `TGL_TRANSFER` | `sysdate` + `ROWNUM = 1` membuat jalur banding **non-deterministik** bila dua baris masuk di detik yang sama | ⛔ | DBA |
| **D13** | `DEDUCTION2_MENJADI` / `_SELISIH` pada baris endorsement | Tidak pernah ditulis siklus EDM; perlu konfirmasi apakah agregasi hilir mengandaikan sebaliknya | ⛔ | Keuangan |
| D9 | Titik masuk portal (`Rule-Portal`, access group) | Layar pertama pengguna tidak dapat direkam | ⚠️ | Bisnis + IT |

⚠️ `[terverifikasi]` **Tidak ada rule `Property` di korpus.** Konsekuensinya menyeluruh: **seluruh
tipe data properti berstatus `belum terverifikasi`**. Struktur objek hanya dapat disimpulkan dari
rujukan di Section/Activity/DataTransform — itu `[dugaan]`, bukan `[terverifikasi]`.

---

## 3. Integritas ekspor — rule yang hilang dan versi yang tidak sepadan

### ⛔ R0 — 86 rule bernama sama berbeda **versi ruleset** antara NB dan EDM

`[terverifikasi]` (`01-flow/03-alur-endorsement.md`) Dari 1.571 rule bernama sama, **86 punya
`pyRuleSetVersion` berbeda** antara folder NB dan EDM. Bukan perbedaan kosmetik:

- `SetToInbox_ACT` di EDM bernomor versi lebih **lama** (2025-09) dibanding NB/RNW (2026-08), dan
  salinan EDM **kehilangan satu cabang routing**.
- Sebaliknya `SetBanding_ACT` di EDM lebih **baru** dan menambah klausa persetujuan yang tidak ada di
  NB/RNW.

**Arahnya tidak konsisten** — bukan "EDM tertinggal", melainkan tiga folder diekspor dari titik waktu
yang berbeda-beda.

⚠️ **Konsekuensi terhadap angka discovery ini:** sebagian dari 460 rule yang terhitung "bercabang
antar siklus" kemungkinan adalah **selisih versi ekspor, bukan percabangan bisnis**. Kedua hal itu
tidak dapat dipisahkan dari korpus ini. **Rekonsiliasi paralel run mustahil sampai ada satu ekspor
produksi tunggal yang konsisten.**

**Yang diminta:** satu ekspor ulang dari lingkungan **produksi** pada satu titik waktu, untuk
ketiga siklus sekaligus. Pemilik: IT.

### Rule yang dirujuk tetapi hilang dari ekspor

| # | Rule | Dampak | Blokir | Pemilik |
| --- | --- | --- | :-: | --- |
| **R1** | `GetLimitAkseptasi1SA_Act`, `GetLimitAkseptasi2SA_Act`, `GetLimitAkseptasi_Act2` | Jalur **Special Acceptance** dan "limit putaran 2" tidak dapat dimigrasikan | ⛔ | IT (ekspor ulang) |
| **R2** | `serviceInsertArasapas_act` kelas `ASM-FW-GISFW-Work` | Endpoint, payload, dan penanganan galat konversi ke produksi **tidak dapat direkam** | ⛔ | IT (ekspor ulang) |
| **R3** | `UpdateErrorNoteJsonPolis` | `UPDATE` apa, dan arti `sts_konversi = 9` | ⛔ | IT (ekspor ulang) |
| **R4** | 21 rujukan Section yang berkasnya tidak ada | Layar tidak lengkap | ⛔ | IT (ekspor ulang) |
| R5 | `pyIsIPad` | 1 dependensi predikat | ⚠️ | IT |
| R6 | 5 FlowAction tanpa `pySectionReference` | Aksi tanpa layar | ⚠️ | IT |

**Catatan penting:** R1–R4 kemungkinan besar dapat diselesaikan dengan **ekspor ulang dari Pega**
sebelum sistem lama dimatikan. Ini pekerjaan yang punya tenggat keras — setelah Pega padam, tidak
ada lagi sumbernya.

---

## 4. Aritmetika dan representasi uang

### ⛔ Memblokir

| # | Pertanyaan | Bukti |
| --- | --- | --- |
| **U1** | **Satuan `.Rate`.** `CountPremi_ACT` membagi dengan `1e9` (konsisten rate per mille); `CountProrateExtension_Act` step 7 membagi dengan `1e4` (konsisten persen). **Satu ordo besaran 10 berbeda.** | `04-aturan/02` §1.2 vs §1.5 |
| **U2** | **Perbandingan ambang uang sebagai string.** 8 bentuk, mis. `.TSILiability > "3000000000"`. Ambang **yang sama** (30 miliar) dibandingkan numerik di rule lain. Leksikografis: `"4000000000" > "30000000000"` bernilai **BENAR** padahal 4 miliar < 30 miliar. Menentukan tangga akseptasi dan proteksi spreading. | `04-aturan/02` §7.1 B2 |
| **U3** | **Arah konversi pemisah desimal berlawanan** di dua rule INSERT bersaudara untuk kolom yang sama: `InsertIntoFacinSPreadLife_Sql` titik→koma, `InsertIntoFacinSPreadLifeMonthly_Sql` koma→titik. Salah satunya menulis angka dengan pemisah yang salah ke Oracle. | `04-aturan/02` §7.1 B4 |
| **U4** | **Skala hasil `Tree_ShortPeriod`.** Cabang `GroupPanel="007"` mengembalikan 12,5…100; cabang `"002"` mengembalikan 0,125…1. Pemanggil memperlakukan keduanya sama. | `04-aturan/02` §1.6 |
| **U5** | **`Tree_ShortPeriod` rentang 8–44 tidak punya hasil** — baris memakai `= 45`, bukan `<= 45`. | `04-aturan/02` §1.6 |
| **U6** | **`Tree_ShortPeriod` tidak punya cabang default** untuk `GroupPanel` selain `"002"`/`"007"`. Nilai balik tidak diketahui. | `04-aturan/02` §1.6 |
| **U7** | **Tanda PPh/PPN dalam net payable** — empat konvensi berbeda hidup berdampingan: `+PPh +PPN`, `−PPh −PPN`, `−PPh +PPN`, dan varian dengan/tanpa `Deduction2`. | `04-aturan/02` §2.6 |
| **U8** | **`EDMPremiMenjadi` dihitung berbeda antar lini bisnis** — cabang FIRE memakai `EDMOldPayment + EDMNewPremi`, cabang ANEKA/GOLF memakai `EDMOldPremi + EDMNewPremi`. Dua basis berbeda untuk nilai yang sama. | `01-flow/03` |

### ⚠️ Perlu keputusan bisnis

`ProRateType != 3` menihilkan `PremiLifeNusantaraRe` — tiga dari empat nilai memicunya ·
presisi default `@Math.divide` dua argumen untuk nilai uang · pembulatan ulang 4 desimal hanya untuk
sebagian lini bisnis, sisanya 20 desimal · `Param.Spreaded` tidak direset per coverage di sebagian
cabang · prorate tertukar antara `.Premium` dan `.PremiRp` di cabang endorsement · `FirstLossScale`
dijaga di satu field tetapi tidak di field pasangannya · `@divide(.LostLimit,100,1)` membulatkan LOL
ke 1 desimal · `.PercentageAdjustment!="100" || !="75"` adalah tautologi · `Policy.Payment.Premium`
diset di luar loop mata uang sehingga menyimpan mata uang terakhir.

---

## 5. Kontradiksi internal sistem lama

`[terverifikasi]` Semuanya terbaca di korpus. Semuanya **kandidat perbaikan** — memutuskannya milik
bisnis (`CLAUDE.md` §1). Migrasi yang diam-diam memperbaikinya menghasilkan selisih angka yang tidak
dapat dijelaskan saat rekonsiliasi.

### ⛔ Memblokir

| # | Kontradiksi |
| --- | --- |
| **K1** | **Tautologi pembanding limit** (`TotalTSI >= Limit \|\| TotalTSI <= Limit`) di dua blok mesin akseptasi — selalu benar. Disengaja (penyaringan diserahkan ke SQL) atau bug yang menonaktifkan pemeriksaan `MAX_LIMIT_IDR`? |
| **K2** | **Loop tanpa transisi keluar**: bila beberapa baris limit cocok, yang menang adalah **baris terakhir**, bukan pertama. Bertentangan dengan pemahaman "baris pertama = approver berikutnya". |
| **K3** | **`IsEdmAdjCeding` dan `IsEdmPPNPPH` punya kondisi identik** — dua jenis endorsement tidak dapat dibedakan. |
| **K4** | **`IsT1T3` menguji TeamGroup 1 dan 2; `IsT2T4` menguji TeamGroup 2 dua kali (tidak pernah 4).** Nama rule bertentangan dengan isinya. Mana yang mencerminkan maksud bisnis? |
| **K5** | **`.Quotation.BusinessType` vs `.OfferFacIn.QuotationData.BusinessType`** — 30 rule membaca properti berbeda tergantung siklus. Bila keduanya tidak sinkron, klasifikasi produk berbeda antara NB dan EDM untuk kasus yang sama. Hal serupa berlaku untuk **dua properti `StatusBusiness`**. |
| **K6** | **Blok `EdmType==1` tidak pernah dapat dieksekusi** — blok "batal" dan blok "batal sejak semula" **keduanya** bergerbang `EdmType==2`. Satu jenis endorsement kehilangan penanganannya. |
| **K7** | **`CountDataEDMElse` langkah 2 berlabel MBU tetapi bergerbang `IsMarineCargo`** — akibatnya MBU **tidak terhitung** dan Marine Cargo **terhitung dua kali**. Contoh lain bahwa label bukan bukti. |
| **K8** | **`Decision39` memutus jalur simpan JSON** untuk case ber-fac-retro. Bila `IsRISlip=1`, case berakhir **tanpa baris `JSON_POLIS` maupun panggilan servis produksi**. Disengaja atau cacat? |
| **K9** | **Endorsement Life melewati seluruh tangga akseptasi** (`Decision19 --When IsLife--> Decision39`), dan `GetLimitAkseptasiLife_Act` tidak ada di korpus EDM. Apakah endorsement jiwa memang tanpa persetujuan berjenjang? |

### ⚠️ Perlu keputusan

`IsCedingConfirm` dibandingkan dengan `Offer`/`Policy` **tanpa tanda kutip** (referensi properti atau
literal?) · `" ElectronicEquipment "` berspasi di depan dan belakang · `IsCedingConfirm` memakai
`"Accepted"` dan `"accepted"` · tiga ejaan `.FlagDelete` (405 kemunculan, satu membandingkan sebagai
string) · `StepStatusFail` arah predikat `inString(...) = 0` · `@LengthOfPageList(...) = Local.Index`
memakai satu `=` di tengah kondisi · `IsNotActive` memakai hasil `CompareDates` langsung sebagai
boolean · `EmailTypeQuotation == "2"` berarti DECLINE sementara `2` di semua `EmailType*` lain
berarti REJECT · `GetAksepBanding_SQL` selalu memakai tabel limit Property apa pun lini bisnisnya ·
`GetLimitAkseptasiBanding_SQL` memfilter `LIMIT_BOTTOM2` tetapi `ORDER BY LIMIT_BOTTOM` ·
kualifikasi skema `POOLDATA.` tidak konsisten antar varian tabel yang sama · `Decision8` ("Which
team?") tidak punya cabang `Else` — case macet bila `LetterNo` kosong.

---

## 6. Authz dan guard berbasis identitas

⚠️ `CLAUDE.md` §6 mensyaratkan **persetujuan manusia** sebelum perubahan authn/authz.

| # | Hal | Blokir | Pemilik |
| --- | --- | :-: | --- |
| **A1** | **`IsAdmin` fail-open**: operator dengan workbasket kosong dihitung sebagai admin | ⛔ | Keamanan + Underwriting |
| **A2** | **`pyWorkBasketList(1)`** — hanya workbasket pertama yang diuji, dengan urutan tak ditentukan. Operator multi-workbasket berperilaku tidak deterministik | ⛔ | Keamanan + IT |
| A3 | **Enam rule memuat guard berbasis identitas orang.** Salah satunya melewati seluruh proses spreading untuk 3 identitas operator + 1 kode kontak marketing. Dipertahankan apa adanya (paralel run cocok) atau dipindah ke atribut peran (paralel run berbeda)? | ⚠️ | Underwriting |
| **A5** | **Prompt AI disimpan sebagai data tabel.** `M_PROMPT_AI` (kunci `KATEGORI_1`+`KATEGORI_2`, kolom `PROMPT_AI`) memasok prompt yang menggerakkan model AI di alur underwriting. Artinya perilaku model **dapat diubah tanpa deployment dan tanpa jejak version control**. Ini memperluas cakupan tinjauan keamanan `CLAUDE.md` §6 dari rule ke **tabel** | ⛔ | Keamanan + IT |
| **A6** | **Kolom PII di jalur produksi endorsement.** `FACINPRODUCTION` menulis 79 kolom di NB/RNW tetapi **82 di EDM** — dan bukan sekadar tambahan: EDM menulis `NO_KONTRAK`, `NO_NPWP`, `NO_KTP`, `STARTDATE_DEBITUR`, `ENDDATE_DEBITUR` sebagai **pengganti** `DEDUCTION2_MENJADI`/`_SELISIH`. Dua di antaranya data pribadi. `CLAUDE.md` §6 mensyaratkan persetujuan sebelum akses *regulated data* | ⛔ | Keamanan + Kepatuhan |
| A4 | **`IsUW` bercabang antar siklus** — satu workbasket diakui di EDM tetapi tidak di NB/RNW. Ini menentukan **siapa yang dianggap underwriter** | ⚠️ | Underwriting |

`[terverifikasi]` **Nilai nama orang tidak disalin ke dokumen ini** (`CLAUDE.md` §3 butir 5). Guard
dicatat sebagai jumlah dan mekanismenya saja.

⚠️ `[terverifikasi]` Selain di dalam logika, nama orang juga muncul di **teks status yang tampil ke
pengguna**: 669 kemunculan / 23 label status, dan 92 kemunculan pola *"IS IN … INBOX"* dengan 14 nama
unik di satu flow saja. Di sistem baru ini harus menjadi rujukan ke antrean/jabatan. **Keputusan
bisnis:** apakah teks status yang dilihat pengguna boleh berubah.

---

## 7. Cakupan yang belum didiscovery

Jangan dianggap "tidak ada".

| Area | Status |
| --- | --- |
| **Spreading · capacity · scoring** | Tersentuh lewat rumus dan proteksi spreading; **belum ada dokumen alurnya sendiri** |
| Isi 460 rule bercabang EDM | Terhitung, **belum dipetakan field-per-field**. Sebelum dipetakan, tidak boleh diklaim "hanya beda kosmetik" |
| `ReportDefinition\` (379 berkas) | Belum dibaca sama sekali |
| Kedalaman modal UI maksimum | Terbukti ≥2 tingkat; maksimumnya belum diukur |
| 141 Section yatim | Termasuk rumpun `*FacOut*` dan `*TreatyIn*` — apakah termasuk lingkup Fac **In**? |

⚠️ `[pertanyaan terbuka]` **Dua FlowAction bernuansa AI** (`AnalysLocationbyAI`, `AttachDoc_AI`) ada
di `NB FacIn\FlowAction\` dan `RNW Fac In\FlowAction\`, **tidak ada** di `Endorsment Fac In\`
`[terverifikasi]`. **Isinya belum dibaca** — nama bukan bukti perilaku. Bila benar memanggil model AI
pihak ketiga, `CLAUDE.md` §6 mengharuskan **tinjauan keamanan sebelum dipindahkan**.

```powershell
foreach ($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  foreach ($n in @('AnalysLocationbyAI','AttachDoc_AI')) {
    "{0,-18} {1,-22} {2}" -f $f, $n, (Test-Path "D:\migrasi\RNM\$f\FlowAction\$n.xml") } }
```

---

## 8. Yang **tidak** perlu ditanyakan

`[terverifikasi]` Berbeda dari enumerasi di §1, **`BusinessType` adalah string deskriptif** dengan
**39 nilai unik** yang terbaca langsung: `Fire`, `MarineCargo`, `MBUCar`, `MBUMotorCycle`, `PA`,
`Travel`, dan seterusnya. Tidak perlu kamus dari bisnis.

⚠️ Satu-satunya catatan: sebagian nilai muncul **dengan spasi di depan/belakang**. Normalisasi di
batas input — dan catat sebagai kandidat perbaikan, jangan diperbaiki diam-diam (§5).

---

## 9. Koreksi terhadap `CLAUDE.md` §4.5

`[terverifikasi]` §4.5 memerintahkan `panic()` untuk lima rule `When` yang "kondisinya kosong di
ekspor". **Pemeriksaan ulang atas korpus membantahnya**: jumlah rule `When` yang kondisinya tidak
terbaca sama sekali adalah **0** dari 601 berkas. Kelimanya terbaca:

| Rule | Kondisi terbaca | Letak |
| --- | --- | --- |
| `IsPKSASM` | `OfferFacIn.IsB2B = "ASM"` | hanya `<pyConditionValue1String>`; `<pyConditionString>` berisi placeholder |
| `ToUW` | `LetterNo = "UNDERWRITER"` | **kedua** tag terisi dan sepakat |
| `ToJUW_A` | `LetterNo = "JUW_A"` | **kedua** tag terisi dan sepakat |
| `LetterNoNull` | `LetterNo = ""` | hanya `<pyConditionValue1String>` |
| `IsEdmInternalRetro` | `QuotationData.EndorsementInternalRetro = 1` | hanya `<pyConditionValue1String>` |

```powershell
foreach ($n in @("IsPKSASM","ToUW","ToJUW_A","LetterNoNull","IsEdmInternalRetro")) {
  "=== $n ==="
  foreach ($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
    $p = "D:\migrasi\RNM\$f\When\$n.xml"
    if (Test-Path $p) {
      $cs = (Select-String -Path $p -Pattern '<pyConditionString>(.*?)</pyConditionString>' | ForEach-Object { $_.Matches[0].Groups[1].Value }) -join ' | '
      $v1 = (Select-String -Path $p -Pattern '<pyConditionValue1String>(.*?)</pyConditionValue1String>' | ForEach-Object { $_.Matches[0].Groups[1].Value }) -join ' | '
      "  [$f] Str='$cs'  V1='$v1'" } } }
```

**Penyebab kekeliruan sebelumnya:** discovery terdahulu hanya membaca `<pyConditionString>`.
**170 berkas (28,3 %) menyembunyikan kondisinya** dari tag itu — 13 berkas dengan
`<pyConditionString/>` kosong dan 157 dengan placeholder `[Double click to add condition]`.

⚠️ `[pertanyaan terbuka]` **`IsPKSASM` tetap perlu konfirmasi.** Baris kondisinya membawa
`<pyTempText>true</pyTempText>` dengan label belum ter-resolve
(`[first value][relation][second value]`), berbeda dari `IsNB` yang labelnya ter-resolve penuh.
Ekspresinya terbaca, tetapi apakah itu yang benar-benar dieksekusi belum dapat dipastikan. Catatan:
`pyTempText = true` muncul pada **294 berkas (48,9 %)**, jadi flag itu sendiri **bukan** penanda yang
membedakan.

**Usulan:** ubah `CLAUDE.md` §4.5 agar tidak lagi menyebut kelima rule ini, dan ganti mekanismenya
menjadi: *setiap predikat yang kondisinya tidak diketahui, atau enumerasi yang artinya belum dijawab
bisnis, wajib `panic()` — bukan `return false`.* Butir-butir di §1 dokumen ini adalah daftar
sesungguhnya. **Perubahan `CLAUDE.md` menunggu persetujuan work owner.**

---

*Setiap butir di dokumen ini menunggu jawaban manusia. Tidak satu pun boleh diisi dengan tebakan.*
