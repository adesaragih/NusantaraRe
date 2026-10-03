# Bahan spec — Modul 4: Predikat `When` Endorsement (EDM)

> **Ini BAHAN untuk `/to-spec`, bukan spec.** `/to-spec` ber-`disable-model-invocation: true` dan
> **belum dijalankan**. Work owner yang menjalankannya manual.
>
> **Sumber:** `07-edm\02-e1-pola-perbedaan-when.md` §5–§7. Prior art bentuk:
> `04-spec\03-spec-modul-terverifikasi.md` (registry predikat NB = **Seam 1**) dan
> `04-spec\05/06/07-spec-edm`.
>
> **Keputusan mengikat:** K-005 · K-018 (COB→skala) · K-029 · K-044 · K-046 · K-047 · K-048 · K-049.
>
> ⛔ **Modul ini dependensi tiga modul sebelumnya** — before-image, selisih, dan alur masuk
> semuanya bercabang per lini bisnis (COB).
>
> ✅ **E-Q5 TUNTAS** di putaran ini — §3. ⚠️ **Satu pertanyaan registry/seam** diajukan — §5.

---

## 1. Klasifikasi 202 rule `When` EDM

`[terverifikasi]`

| Ukuran | Jumlah |
| --- | ---: |
| Total `When` di `Endorsment Fac In\` | **202** |
| Bernama sama dengan NB (peka huruf) | **175** |
| — identik isinya *(kontrak K-043)* | **119** |
| — **berbeda isinya** | **56** |
| **EDM-only** | **27** |

⚠️ Angka identik/berbeda memakai kontrak volatile **K-043 yang sudah diamandemen** (23 tag + isi
blok `pzIndexes`). Di bawah kontrak lama angkanya 112/63; perbedaannya artefak penomoran indeks,
bukan isi.

### 1.1 Sumber data — hanya 6 yang membaca query

`[terverifikasi]` Dari **202** rule `When` EDM:

| Sumber kondisi | Jumlah |
| --- | ---: |
| **`OutData.pxResults(...)` — hasil query** | **6** |
| Properti kasus biasa | **196** |

Keenam yang membaca query, seluruhnya predikat lini bisnis:

```
IsAneka · IsFire · IsGolfInsurance · IsMarineCargo · IsMBU · IsPA
```

⛔ **JANGAN menyimpulkan "predikat EDM umumnya membaca query".** Itu **terbatas pada 6 dari 202**
(3,0 %). Sembilan puluh tujuh persen predikat EDM membaca properti kasus, persis seperti NB.

```powershell
$edW="D:\migrasi\RNM\Endorsment Fac In\When"
$q=0;$p=0
Get-ChildItem $edW -Filter *.xml -File | % {
  $v=[regex]::Matches([IO.File]::ReadAllText($_.FullName),
       '<pyConditionValue1String>([^<]*)</pyConditionValue1String>') | % { $_.Groups[1].Value.Trim() }
  if (@($v | ? { $_ -match 'OutData\.pxResults' }).Count) { $q++ } else { $p++ } }
"query=$q  properti=$p"        # => query=6  properti=196
```

---

## 2. Enam predikat COB — beda **sumber data**, bukan beda kelas deklarasi

### 2.1 Penegasan yang wajib masuk spec

`[terverifikasi]` **`pxObjClass` keenamnya = `Rule-Obj-When` di NB maupun EDM.** Tipe dan
deklarasinya **sama**.

⛔ Yang berbeda adalah **kelas halaman yang DIBACA kondisinya** — `Code-Pega-List` /
`ASM-FW-GISFW-Data-Search` di EDM, karena `OutData` adalah halaman **hasil query**.

📌 Ini jebakan keluarga "nama bukan bukti": **nama sama, tipe sama, sumber data berbeda.**
Versi awal discovery sempat menyebutnya "beda kelas deklarasi" — itu keliru dan sudah dikoreksi.

### 2.2 Tabel enam predikat

`[terverifikasi]` Kondisi dibaca dari **`pyConditionValue1String`** (yang dieksekusi), bukan
`pyConditionString` (label tampilan) — `CLAUDE.md` §4.5.

| Rule | Kondisi di **EDM** | Cacah kondisi EDM | Kondisi di **NB** | Cacah NB |
| --- | --- | ---: | --- | ---: |
| `IsFire` | `OutData.pxResults(1).CARI2 = "KPR"` · `"OilGas"` · `"FireStyle1"` · `"FireStyle2"` · `"Fire"` | **5** | **delegasi** ke `IsKPR`, `IsOilGas`, `IsFireStyle1`, `IsFireStyle2` + 1 properti | 5 |
| `IsAneka` | 25 literal `OutData.pxResults(1).CARI2 = "<COB>"` — `"Liability"`, `"MarineHull"`, `"GrowingTrees"`, … | **25** | **delegasi** ke 25 sub-rule + 1 properti langsung | 26 |
| `IsPA` | `OutData.pxResults(1).CARI2 = "PA"` | **1** | `pyWorkPage.Quotation.BusinessType = "PA"` | 1 |
| `IsMBU` | `= "MBUCar"` · `"MBUMotorCycle"` · **keduanya berulang** | **4** | `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "MBUCar"` / `"MBUMotorCycle"` | 2 |
| `IsMarineCargo` | `= "MarineCargo"` | **1** | `.OfferFacIn.QuotationData.BusinessType` **dan** `pyWorkPage.Quotation.BusinessType` | 2 |
| `IsGolfInsurance` | `= "GolfInsurance"` | **1** | `pyWorkPage.Quotation.BusinessType = "GolfInsurance"` | 1 |

**Total cabang yang harus diimplementasikan di EDM: 37 kondisi.**

### 2.3 ⚠️ Dua cara menghitung — jangan tertukar

`[terverifikasi]` Angka **28 · 128 · 7 · 23 · 8 · 9** yang beredar adalah **cacah token mentah**
`OutData.pxResults` di seluruh tag berkas (termasuk `pyConditionString`, `pyPagesAndClasses`,
`pxRuleReferences`). Yang menentukan implementasi adalah **cacah kondisi
`pyConditionValue1String`**: **5 · 25 · 1 · 4 · 1 · 1**.

| Rule | token mentah | **kondisi dieksekusi** |
| --- | ---: | ---: |
| `IsFire` | 28 | **5** |
| `IsAneka` | 128 | **25** |
| `IsPA` | 7 | **1** |
| `IsMBU` | 23 | **4** |
| `IsMarineCargo` | 8 | **1** |
| `IsGolfInsurance` | 9 | **1** |

⛔ Spec memakai **kolom kanan**. Mengimplementasikan 203 cabang alih-alih 37 adalah kekeliruan yang
mudah terjadi bila angka mentah dipakai.

### 2.4 Delegasi **tidak seragam** di NB

`[terverifikasi]` Bukan "NB mendelegasikan, EDM literal". Hitungan penuh:

| Rule | NB: delegasi ke sub-rule | NB: properti langsung | EDM: delegasi | EDM: literal |
| --- | ---: | ---: | ---: | ---: |
| `IsAneka` | **25** | 1 | 0 | 25 |
| `IsFire` | **4** | 1 | 0 | 5 |
| `IsPA` | 0 | 1 | 0 | 1 |
| `IsGolfInsurance` | 0 | 1 | 0 | 1 |
| `IsMarineCargo` | 0 | 2 | 0 | 1 |
| `IsMBU` | 0 | 2 | 0 | 4 |

⛔ **NB mendelegasikan HANYA pada `IsFire` dan `IsAneka`.** Empat lainnya membaca properti langsung.
EDM memakai literal pada keenamnya.

**Akibat:** di endorsement, klasifikasi COB **tidak melewati sub-rule sama sekali**. Sub-rule seperti
`IsKPR`, `IsOilGas`, `IsLiability` **tidak dipanggil** di jalur EDM — literalnya dibandingkan
langsung terhadap kolom query.

---

## 3. ✅ E-Q5 TUNTAS — pengisi `OutData` dan arti `CARI2`

### 3.1 Siapa mengisi `OutData`

`[terverifikasi]` **Satu berkas saja** di seluruh korpus EDM memakai `<BrowsePage>OutData</BrowsePage>`:

```
Activity\SetErrorBatalEndorsement_Act  langkah 7
  RDB-List  « get business type »
    RequestType = GetBusinessType_Sql
    ClassName   = ASM-FW-GISFW-Int-OFFERJSON
    BrowsePage  = OutData
```

📌 Langkah **7** — yaitu **di dalam modul 3 (alur masuk)**, sebelum kasus EDM lahir dan sebelum
before-image berjalan.

### 3.2 Apa isi `CARI2`

`[terverifikasi]` `RDBList\GetBusinessType_Sql.xml` · `pxObjClass` = **`Rule-Connect-SQL`** ⚠️
(bukan tipe RDBList meski di folder `RDBList\`) · `pyClassName` = `ASM-FW-GISFW-Int-OFFERJSON`:

```sql
select a.Data_json.QuotationData.BusinessType as CARI2
  from json_polis a
 where nopolis = {InputData.CARI17}
   and PRODKE  = (SELECT COUNT(NOPOLIS)-1 FROM JSON_POLIS b WHERE NOPOLIS={InputData.CARI17})
```

> ### ⛔ **`CARI2` = `QuotationData.BusinessType` milik POLIS LAMA (versi terakhir).**

**Inilah jawaban E-Q5, dan maknanya besar:** keenam predikat COB di endorsement mengklasifikasikan
menurut **lini bisnis polis yang sedang di-endors**, bukan menurut lini bisnis kasus yang sedang
berjalan. Secara bisnis itu masuk akal — endorsement tidak mengubah lini bisnis polis.

### 3.3 Tiga konsekuensi yang mengikat

`[terverifikasi]`

1. **`GetBusinessType_Sql` EDM-only** — nihil di `NB FacIn\`, dan dipakai **satu activity saja**.
2. **Membaca dokumen JSON**, bukan kolom relasional: `a.Data_json.QuotationData.BusinessType` —
   dot-notation JSON Oracle. Di sistem baru wajib diterjemahkan ke `JSON_VALUE(...)` atau kolom
   tersendiri.
3. **Memakai pola `PRODKE = COUNT(NOPOLIS)-1`** — pola yang sama dengan lapis A before-image, dan
   **rentan terhadap lubang `PRODKE`** (keluarga pertanyaan kerapatan `PRODKE`). Bila `PRODKE`
   berlubang, predikat COB endorsement mengklasifikasikan menurut **versi polis yang salah**.

⚠️ Butir 3 adalah **risiko baru yang belum tercatat**: sebelumnya kerapatan `PRODKE` hanya diketahui
memengaruhi nilai lama. Kini terlihat ia juga memengaruhi **klasifikasi lini bisnis** — yang
menentukan cabang perhitungan premi, skala rasio (K-018), dan jalur produksi.

### 3.4 Urutan eksekusi menjadi kontrak

⛔ `OutData` diisi di **modul 3 langkah 7**. Keenam predikat COB dibaca oleh modul 1, 2, dan 3.

**Urutan wajib:** `GetBusinessType_Sql` → `OutData` terisi → predikat COB dapat dievaluasi. Bila
predikat dievaluasi sebelum query berjalan, keenamnya **bernilai salah** dan seluruh percabangan COB
runtuh diam-diam.

📌 Sejalan dengan kontrak urutan di K-048 §8.1 (porsi periode). **Endorsement punya urutan eksekusi
yang mengikat, bukan modul-modul bebas urutan.**

---

## 4. Tiga perbedaan perilaku

`[terverifikasi]` Ketiganya `pxObjClass` = `Rule-Obj-When` di kedua korpus.

### 4.1 `IsUW` — endorsement mengakui satu workbasket tambahan

| | `pyLogic` | Kondisi (`pyConditionValue1String`) |
| --- | --- | --- |
| `NB FacIn\When\IsUW.xml` | `A OR B OR C` | `ReasFacInDirector` · `ReasFacInGroupLeader` · `ReasFacInUnderwriting` |
| `Endorsment Fac In\When\IsUW.xml` | **`A OR B OR C OR D`** | `ReasFacInDirector` · **`ReasFacInMarketing`** · `ReasFacInGroupLeader` · `ReasFacInUnderwriting` |

`pyDesignatedClass` = `Code-Pega-Requestor` **di kedua sisi** — sumber data sama, murni perbedaan
himpunan kondisi.

⛔ **Di endorsement, operator ber-workbasket marketing dihitung sebagai underwriter.** Ini menentukan
**siapa yang boleh menyetujui**.

📌 Menyatu dengan modul 3: kasus EDM **lahir** di workbasket `ReasFacInMarketing`. Jadi orang yang
membuat endorsement termasuk yang dapat bertindak sebagai UW atasnya.

⚠️ **E-Q6 `[pertanyaan terbuka]` — milik work owner + Underwriting.** Disengaja atau tidak, **bukan
keputusan teknis**. Sampai dijawab, perilakunya **diport apa adanya**.

### 4.2 `IsClaim` — uji berbeda sifat, kelas sumber berbeda

| | `pyDesignatedClass` | Kondisi |
| --- | --- | --- |
| `NB FacIn\When\IsClaim.xml` | `ASM-FW-GCNMFW-Work-PNC` | `pyWorkPage.pyWorkIDPrefix = "CLM-"` — menguji **prefiks ID kasus** |
| `Endorsment Fac In\When\IsClaim.xml` | **`ASM-FW-GISFW-Work-Claim`** | `pyWorkPage.ClaimData has a value` — menguji **keberadaan halaman data** |

⛔ Bukan jalur berbeda, bukan nilai berbeda: **dua pengujian berbeda sifatnya**. Keduanya dapat
memberi jawaban berlainan pada kasus yang sama — misalnya kasus berprefiks `CLM-` yang halaman
`ClaimData`-nya kosong.

### 4.3 `IsTravel` — dua sumber di EDM, satu di NB

| | `pyLogic` | `pyConditionValue1String` | `pyConditionString` (label) |
| --- | --- | --- | --- |
| `NB FacIn\When\IsTravel.xml` | `A` | `pyWorkPage.Quotation.BusinessType = "Travel"` | `BusinessCode = "02"/"58"/"SB"/"SG"` |
| `Endorsment Fac In\When\IsTravel.xml` | **`A OR B`** | `.Quotation.BusinessType = "Travel"` **OR** `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "Travel"` | `Kode Bisnis = "77"` |

⛔ **Label tidak sejalan dengan ekspresi tersimpan di kedua korpus.** Yang dieksekusi
`pyConditionValue1String` (`CLAUDE.md` §4.5). Label `pyConditionString` **tidak boleh** dijadikan
dasar perilaku.

📌 `IsTravel` **tidak** memakai pola query — ia membaca properti. Pola query terbatas pada 6 COB.

---

## 5. ⚠️ Kontrak registry — pertanyaan untuk work owner (E-Q12)

### 5.1 Siapa memakai predikat ini

| Modul | Memakai predikat untuk |
| --- | --- |
| **Modul 1** — before-image | `SetOldData` blok 4: tujuh cabang lini (`IsFire`, `IsGolfInsurance`, `IsAneka`, `IsPA`, `IsMarineCargo`, `IsMBU`, `IsTravel`) + gerbang keluar `IsLife`, `IsEDM`; lapis C tujuh varian |
| **Modul 2** — selisih | tabel rumus per (lini × `EdmType` × jalur); gerbang lini di `CountEndorsementData` / `CountDataEDMElse` |
| **Modul 3** — alur masuk | 7 gerbang lini bisnis di `SetErrorBatalEndorsement_Act`; gerbang siklus `IsEDM`/`IsNotEDM` |
| **Jalur produksi** *(belum dispec)* | 7 cabang `SaveFacinProdAllEDM_Act` |

⛔ **Predikat COB adalah simpul pusat seluruh siklus endorsement.** Salah di sini menjalar ke semua.

### 5.2 Masalahnya

`[terverifikasi]` Registry NB (**Seam 1**, `internal/rules.Eval`) memuat predikat yang membaca
**properti kasus**. Keenam predikat COB EDM membaca **kolom hasil query dari polis lama**.

⛔ **Registry NB tidak dapat dipakai ulang apa adanya untuk keenam COB EDM.** Bukan karena namanya
berbeda — namanya sama persis — melainkan karena **sumber datanya berbeda**.

Tambahan yang memperumit: **196 dari 202** predikat EDM membaca properti kasus **persis seperti NB**.
Jadi bukan "registry EDM berbeda"; yang berbeda hanya **6 anggota**.

### 5.3 Dua pilihan — keputusan work owner

| | **Pilihan A — Seam 1 diperluas** | **Pilihan B — registry kedua untuk 6 COB** |
| --- | --- | --- |
| Bentuk | satu registry; tiap predikat mendeklarasikan **sumber datanya** (properti kasus / kolom query), `rules.Eval` menerima konteks berisi keduanya | registry NB tetap apa adanya untuk 196; registry EDM terpisah khusus 6 COB |
| Untung | satu tempat untuk seluruh 202; 196 predikat yang sama persis **benar-benar dipakai ulang**; tidak menambah seam | pemisahan tegas; registry NB tidak tersentuh sama sekali |
| Rugi | mengubah kontrak Seam 1 yang sudah disetujui untuk NB | 196 predikat identik terduplikasi atau harus di-*delegate*; dua tempat mencari satu nama rule |
| Seam | **tetap 5** — Seam 1 diperluas, bukan ditambah | **tetap 5** — registry kedua diuji lewat Seam 1 yang sama |

⚠️ **Keduanya tidak menambah seam.** Total tetap **5** (K-049). Yang diputuskan adalah **bentuk
registry**, bukan jumlah seam.

**Usulan saya: Pilihan A.** Alasannya: **196 dari 202 predikat identik dengan NB** — memisahkan
seluruh registry demi 6 anggota akan menduplikasi 97 % yang sama. Mendeklarasikan sumber data
per-predikat juga lebih jujur terhadap korpus: sumber data **memang** atribut per-rule, bukan atribut
per-siklus (terbukti dari P1 yang terbantah — NB sendiri campuran jalur properti).

⛔ **Tidak saya putuskan.** Menunggu keputusan work owner, seperti Seam 4 dan Seam 5 dulu.

### 5.4 Yang berlaku apa pun pilihannya

- **Resolver per-rule, bukan per-siklus.** Pola P1 (*"EDM baca agregat tersimpan, NB baca halaman
  aktif"*) **terbantah** — NB sendiri memakai kedua jalur, kadang keduanya dalam satu rule
  (`IsMarineCargo` NB: 2 kondisi, dua jalur berbeda).
- **`IsNotEDM` bukan negasi `IsEDM`** — jalur propertinya berbeda (agregat tersimpan vs halaman
  aktif). Wajib implementasi sendiri.
- **Nama rule tak dikenal → `panic`** (kontrak Seam 1 NB, berlaku penuh).
- **Kedua tag kondisi wajib dibaca**, dan **yang mengikat `pyConditionValue1String`**
  (`CLAUDE.md` §4.5).
- **Predikat COB menyalakan resolver skala K-018** — salah COB berarti salah skala (‰ vs %), yang
  menggeser premi faktor seribu.

---

## 6. Kejanggalan K-046 — kasus uji bernama, diport apa adanya

⛔ Nama test **wajib menyebut K-046**.

| # | Kejanggalan | Bukti | Nama test yang disarankan |
| ---: | --- | --- | --- |
| 1 | **`IsMBU` EDM cabang ganda** — `"MBUCar"` dan `"MBUMotorCycle"` masing-masing **dua kali** (4 kondisi untuk 2 nilai). Redundan, hasil tidak berubah | §2.2 | `K046_IsMBU_CabangGanda_HasilTidakBerubah` |
| 2 | **`IsAneka` label tidak sinkron** — 6 `pyConditionString` berlabel `Kode Bisnis = "24"/"18"/…` yang tidak berkaitan dengan 25 `pyConditionValue1String`-nya | §2.2 | `K046_IsAneka_LabelKodeBisnis_Diabaikan_IkutiValue1` |
| 3 | **`IsTravel` label vs ekspresi** berbeda di kedua korpus | §4.3 | `K046_IsTravel_LabelTidakDieksekusi` |
| 4 | **`IsUW` mengakui workbasket marketing** di EDM | §4.1 | `K046_IsUW_MarketingSebagaiUnderwriter_EDM` |
| 5 | **`IsClaim` dua uji berbeda sifat** antara NB dan EDM | §4.2 | `K046_IsClaim_UjiPrefiks_vs_UjiHalaman` |
| 6 | **Predikat COB EDM tidak melewati sub-rule** — sub-rule NB seperti `IsKPR`, `IsOilGas` tidak dipanggil di jalur EDM | §2.4 | `K046_COB_EDM_TanpaDelegasiSubRule` |

---

## 7. Out of Scope

1. **Tangga akseptasi endorsement** — termasuk klaim arsip §4.3 *"nilai dasar akseptasi = selisih
   TSI"*, tetap **`[belum diuji]`**.
2. **27 rule `When` EDM-only** — belum dibedah isinya; masuk pekerjaan 354 EDM-only.
3. **56 rule `When` bernama sama yang berbeda isinya** di luar 6 COB dan 3 perbedaan perilaku —
   belum ditelusuri satu per satu.
4. **Isi tabel `json_polis`** dan kerapatan `PRODKE` — milik DBA.
5. **Registry predikat NB** — sudah dispec (`03-spec-modul-terverifikasi.md`); modul ini hanya
   membahas perluasan atau pendampingnya.
6. **Sub-rule COB NB** (`IsKPR`, `IsOilGas`, `IsLiability`, …) — tetap dipakai jalur NB, tidak
   disentuh modul ini.

---

## 8. Pertanyaan terbuka

| # | Pertanyaan | Status | Pemilik |
| --- | --- | --- | --- |
| **E-Q5** | Bagaimana `OutData` diisi dan apa arti `CARI2` | ✅ **TUNTAS** — `GetBusinessType_Sql`, `CARI2` = `BusinessType` polis lama (§3) | — |
| **E-Q12** | Bentuk registry: Seam 1 diperluas atau registry kedua | ⚠️ **diajukan** (§5.3), usulan Pilihan A | **work owner** |
| **E-Q6** | `IsUW` — marketing dihitung underwriter di EDM. Disengaja? | `[pertanyaan terbuka]` | **work owner + Underwriting** |
| **baru** | **Kerapatan `PRODKE` juga memengaruhi klasifikasi COB**, bukan hanya nilai lama (§3.3 butir 3) | `[pertanyaan terbuka]` | **DBA** |
| **E-Q7** | `IsClaim` dua uji berbeda sifat — disengaja? | `[pertanyaan terbuka]`, tidak memblokir | work owner |
| **E-Q8** | `IsTravel` dua sumber di EDM — disengaja? | `[pertanyaan terbuka]`, tidak memblokir | work owner |

---

## 9. Kesiapan

✅ **Bahan lengkap untuk `/to-spec`:**

| Butir | Status |
| --- | --- |
| Klasifikasi 202 `When`: 175 bernama sama (119 identik / 56 beda), 27 EDM-only | §1 |
| **6 baca query vs 196 baca properti** + perintah audit | §1.1 |
| Tabel 6 predikat COB — EDM vs NB, dengan penegasan **beda sumber data bukan kelas deklarasi** | §2.1–§2.2 |
| **Dua cara menghitung dibedakan** — 37 kondisi dieksekusi, bukan 203 token | §2.3 |
| Delegasi tidak seragam — NB hanya `IsFire` & `IsAneka` | §2.4 |
| ✅ **E-Q5 tuntas** — `GetBusinessType_Sql`, `CARI2` = BusinessType polis lama, + 3 konsekuensi | §3 |
| **Kontrak urutan eksekusi** — `OutData` diisi modul 3 langkah 7 | §3.4 |
| 3 perbedaan perilaku dengan kontrak masing-masing | §4 |
| Kontrak registry + **E-Q12 diajukan** (Pilihan A/B, usulan A, **seam tetap 5**) | §5 |
| 6 kejanggalan K-046 dengan **nama test** | §6 |

⛔ **`/to-spec` dan `/to-tickets` TIDAK dijalankan.**
⚠️ **E-Q12 (§5.3) sebaiknya diputuskan sebelum spec ditulis** — ia menentukan bentuk modul.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
