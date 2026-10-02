# Audit 11 — Daftar periksa **UI Pega**

> **21 September 2026.** Berisi **hanya** pertanyaan yang **tidak dapat dijawab dari korpus XML** dan
> menuntut membuka rule di layar Pega selagi sistem lama masih hidup.
>
> ⛔ **Jangan menebak jawabannya.** Tiap butir menyebut **berkas + alamat langkah** yang harus dibuka,
> **apa yang harus dicatat**, dan **keputusan/tiket mana yang tergantung** padanya. Sampai terjawab,
> statusnya tetap `[pertanyaan terbuka]`.
>
> Paket permintaan pihak luar (DBA/IT/Product) ada terpisah di
> `..\_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md`. Berkas ini **bukan** duplikatnya: yang di sini hanya dapat
> dijawab dengan **melihat layar**, bukan dengan query atau ekspor.

---

## Ringkasan

| # | Pertanyaan | Bergantung padanya |
| :-: | --- | --- |
| **U-1** | Arti kode transisi **5** | **F10** arah gerbang idempotensi · seluruh pemetaan alur activity |
| **U-2** | Arti kode transisi **1** dan **4** | `10-audit\02` §2.3 · setiap activity bercabang |
| **U-3** | Berapa rule bernama `IsFire`, `IsPA`, `IsUW`, dan di kelas apa saja | **K-050** premis "196 predikat identik" · **F01** · `10-audit\10` |
| **U-4** | Apakah rule `CountRateRetroCov` kelas `Data-Cargo` ada | **K-060** pertanyaan terbuka 1 · **F06** |

⚠️ **U-1 dan U-2 saling menopang.** Kode `2`, `3` dan `6` sudah tertambat pada tingkat `[dugaan kuat]`
(`02-peta-kode-transisi.md` §2.2). Bila U-1/U-2 terjawab, seluruh peta kode transisi naik ke
`[terverifikasi]` sekaligus.

---

## U-1 · Arti kode transisi **5**

### Pertanyaan

`<pyStepsPreCondParamsWhenTrue>` / `<…WhenFalse>` bernilai **`5`** — **apa yang dilakukan Pega**?

`[terverifikasi]` Kode 5 adalah yang **paling sering** di antara yang belum tertambat: **908
kemunculan** pada satu pasangan saja (T=5, F=2), seluruhnya bergerbang
(`02-peta-kode-transisi.md` §1).

⚠️ Hipotesis "kendali perulangan" **sudah diuji dan gagal**: langkah pembawa kode 5 beriterasi hanya
**9,8 %**, lebih rendah daripada kode 2 (**30,8 %**). Jadi **jangan** menawarkan jawaban itu kembali.

### Di mana membukanya — **empat contoh, seluruhnya T=5 / F=2**

| # | Berkas | Alamat langkah | Gerbang |
| ---: | --- | --- | --- |
| 1 | `NB FacIn\Activity\SetDataFacOut_Act` | `RH_1.pySteps(1)` | `IsGroup` |
| 2 | `NB FacIn\Activity\SetDataFacOut_Act` | `RH_1.pySteps(8).pySteps(1).pySteps(5)` | `IsAneka` |
| 3 | `NB FacIn\Activity\SetDataFacOut_Act` | `RH_1.pySteps(8).pySteps(1).pySteps(8)` | `IsPA` |
| 4 | `NB FacIn\Activity\CopyToAllSpreading_ACT` | `RH_1.pySteps(9)` | `IsPA` |

📌 Contoh 2 dan 3 **bersarang di dalam** langkah 8 → berguna untuk melihat apakah arti kode 5 berbeda
di dalam langkah bersarang dibanding di puncak.

### Apa yang harus dicatat

1. **Teks persis** yang tampil di kolom transisi form Activity untuk baris itu (mis. *"Jump to…"*,
   *"Exit iteration"*, *"Continue whens"*, atau apa pun bunyinya) — **salin apa adanya**, jangan
   diringkas.
2. Apakah teksnya **sama** untuk keempat contoh, atau berbeda antara langkah puncak dan bersarang.
3. Apakah ada **medan sasaran** yang terisi di sebelah baris itu (label/step name). Korpus mencatat
   tag `…Prms` **kosong** untuk kode 5 — konfirmasi apakah layar juga menampilkannya kosong.
4. Bila teksnya menyebut lompatan: **ke langkah mana** contoh 1 melompat.

### Tergantung padanya

- **F10** — arah gerbang idempotensi `InsertFacoutProduction` bersandar pada peta §2.2. Bila peta itu
  berubah, `10-audit\02` §3 ikut gugur dan F10 ditinjau ulang.
- Setiap pemetaan alur activity yang menyebut "langkah dilewati/dijalankan".

---

## U-2 · Arti kode transisi **1** dan **4**

### Pertanyaan

Sama seperti U-1, untuk kode **1** dan **4**. `02-peta-kode-transisi.md` §2.3 menandai keduanya
**tidak tertambat** dan **tidak memuat contoh beralamat** — contoh di bawah **dicari untuk berkas
ini**.

`[terverifikasi]` Pengukuran sendiri, seluruh korpus, parser anak-langsung, **hanya baris precondition
yang bergerbang** (`<pyStepsPreCondParamsWhen>` tidak kosong):

| Kode | Baris precondition bergerbang yang membawanya di T atau F |
| ---: | ---: |
| **1** | **347** |
| **4** | **120** |

⚠️ **Angka ini tidak sebanding langsung dengan §2.3 audit 02**, yang menghitung **pasangan (T,F)** dan
**menyertakan baris tanpa gerbang**. Keduanya benar untuk apa yang masing-masing diukur; definisinya
yang berbeda.

### Di mana membukanya — **kode 1**

| # | Berkas | Alamat langkah | Metode | Gerbang | T | F | `…Prms` |
| ---: | --- | --- | --- | --- | ---: | ---: | --- |
| 1 | `NB FacIn\Activity\CheckDeductible_Act` | `RH_2.pySteps(4)` | *(kosong)* | `IsFire` | **1** | **1** | T=`FIRE` · F=`NXTT` |
| 2 | `NB FacIn\Activity\ASMForceCaseClose` | `RH_1.pySteps(13)` | `Commit` | `param.Commit` | 2 | **1** | F=`Exit` |
| 3 | `NB FacIn\Activity\CekLimitSpreading_Act` | `RH_1.pySteps(6)` | *(kosong)* | `pyWorkPage.IsFlagUW!=""` | **1** | 2 | T=`save` |

⛔ **Contoh 1 adalah yang paling menjelaskan:** **kedua** cabang bernilai `1`, dengan **sasaran
berbeda** (`FIRE` vs `NXTT`). Nilai `…Prms` yang terisi hanya muncul pada kode 1
(`02-peta-kode-transisi.md` §2.2 catatan terakhir) — layar Pega akan memperlihatkan medan apa itu.

### Di mana membukanya — **kode 4**

| # | Berkas | Alamat langkah | Metode | Gerbang | T | F |
| ---: | --- | --- | --- | --- | ---: | ---: |
| 1 | `NB FacIn\Activity\CheckDeductible_Act` | `RH_2.pySteps(3).pySteps(2).pySteps(2).pySteps(2).pySteps(2)` | `Property-Set` | `@String.equals(local.LIndex,local.RIndex)` | **4** | 2 |
| 2 | `NB FacIn\Activity\CheckDeductible_Act` | `RH_2.pySteps(3)` | *(kosong)* | `@String.equals(local.IsSame,"0")` | *(kosong)* | **4** |
| 3 | `NB FacIn\Activity\AddCurencyList_ACT` | `RH_1.pySteps(7).pySteps(1).pySteps(2).pySteps(1).pySteps(1)` | `Property-Set` | `Local.MasterCurrency==.Currency.Name` | *(kosong)* | **4** |

📌 Contoh 1 dan 3 berada **dalam** langkah beriterasi yang bersarang dalam — layak diperiksa apakah
kode 4 berkaitan dengan **iterasi** (mis. berhenti dari loop), yang akan menjelaskan mengapa ia hampir
selalu muncul di kedalaman.

### Apa yang harus dicatat

1. **Teks persis** kolom transisi untuk tiap contoh — kode 1 dan kode 4 **terpisah**.
2. Untuk kode 1: **nama medan** tempat `FIRE` / `NXTT` / `Exit` / `save` tampil, dan apakah nilai itu
   **cocok dengan nama langkah lain** di activity yang sama (label lompatan) atau sesuatu yang lain.
3. Untuk kode 4: apakah artinya berubah bila langkahnya **beriterasi** dibanding tidak.
4. Apakah nilai T atau F yang **kosong** di XML tampil sebagai pilihan tertentu di layar, atau
   benar-benar kosong.

### Tergantung padanya

- `10-audit\02` §2.3 — tiga kode tak tertambat; kode 5 (U-1) dan 1 + 4 (U-2) adalah seluruh sisanya.
- Tiap activity bercabang yang alurnya dipetakan tanpa mengetahui arti kode ini.

---

## U-3 · Berapa rule bernama `IsFire`, `IsPA`, `IsUW` — dan di kelas apa saja

### Pertanyaan

Untuk masing-masing dari ketiga nama: **berapa rule `When` yang ada di Pega, dan di kelas apa saja?**

`[terverifikasi]` Di korpus, ketiganya berwujud **tiga berkas → dua identitas berbeda**, karena
salinan Endorsement berada di **kelas lain**:

| Nama | NB · RNW | Endorsment | Identitas berbeda |
| --- | --- | --- | :-: |
| `IsFire` | `ASM-FW-GISFW-Work` | **`ASM-SFAGIS-Work-Endorsement`** | **2** |
| `IsPA` | `Data-Party-Person` | **`ASM-SFAGIS-Work-Endorsement`** | **2** |
| `IsUW` | `ASM-FW-GISFW-Work` | **`ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance`** | **2** |

Basis `pzInsKey` masing-masing (stempel waktu dibuang):

```
RULE-OBJ-WHEN ASM-FW-GISFW-WORK ISFIRE      |  RULE-OBJ-WHEN ASM-SFAGIS-WORK-ENDORSEMENT ISFIRE
RULE-OBJ-WHEN DATA-PARTY-PERSON ISPA        |  RULE-OBJ-WHEN ASM-SFAGIS-WORK-ENDORSEMENT ISPA
RULE-OBJ-WHEN ASM-FW-GISFW-WORK ISUW        |  RULE-OBJ-WHEN ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE ISUW
```

⚠️ **Korpus tidak dapat menjawab pertanyaannya.** Ekspor hanya memperlihatkan rule yang **ikut
terekspor**; ia tidak dapat membuktikan bahwa **tepat dua** yang ada di Pega, dan tidak dapat
menunjukkan **rule mana yang menang** saat runtime (pewarisan kelas Pega). Ketiganya termasuk **109
nama rule kembar** di `10-audit\10-rule-kembar-pzinskey.md`, yang memuat **42 rule `When`**.

### Di mana membukanya

Di **Application Explorer / Records → When**, cari tiap nama dan buka **daftar seluruh instansinya**
(bukan satu yang terbuka pertama). Ketiga berkas korpus sebagai titik awal:

| Nama | Berkas korpus |
| --- | --- |
| `IsFire` | `NB FacIn\When\IsFire.xml` · `RNW Fac In\When\IsFire.xml` · `Endorsment Fac In\When\IsFire.xml` |
| `IsPA` | `NB FacIn\When\IsPA.xml` · `RNW Fac In\When\IsPA.xml` · `Endorsment Fac In\When\IsPA.xml` |
| `IsUW` | `NB FacIn\When\IsUW.xml` · `RNW Fac In\When\IsUW.xml` · `Endorsment Fac In\When\IsUW.xml` |

### Apa yang harus dicatat

1. **Jumlah instansi** tiap nama di Pega, dan **kelas** masing-masing — termasuk yang **tidak** ada
   di korpus.
2. **Ruleset + versi** tiap instansi.
3. **Kondisi** tiap instansi (apakah dua instansi menguji hal yang **sama** atau **berbeda**).
4. Ketika sebuah kasus Endorsement berjalan, **instansi mana** yang dipakai — dan apakah itu akibat
   pewarisan kelas atau ruleset.
5. ⛔ **Bila kondisinya berbeda:** apakah itu **disengaja per siklus**, atau salah satu tertinggal.

### Tergantung padanya

- **K-050** — premisnya "**196 predikat identik** dalam satu registry". `[terverifikasi]` **42 dari
  109** nama rule kembar bertipe `When`, hampir seluruhnya keluarga predikat COB. ⛔ **Apakah premis
  K-050 berubah adalah keputusan work owner, bukan kesimpulan berkas ini.**
- **F01** — registry predikat Fac Out; komentar ketertelusuran **wajib menyebut kelas** (`CLAUDE.md`
  §4.6), bukan hanya nama rule.
- `10-audit\10-rule-kembar-pzinskey.md` — ketiga nama ini adalah sampel yang mewakili pola tersebut.

---

## U-4 · Apakah rule `CountRateRetroCov` kelas `ASM-FW-GISFW-Data-Cargo` ada

### Pertanyaan

**Ada / pernah ada lalu dihapus / tidak pernah ada?**

`[terverifikasi]` Indeks `Embed-Reference-Rule` pada **7 berkas** mencatat `CountRateRetroCov` pada
**tiga** kelas, masing-masing **7 kali**: `ASM-FW-GISFW-Data-PropertyItem`,
`ASM-FW-GISFW-Data-Aneka`, dan **`ASM-FW-GISFW-Data-Cargo`**.

⛔ **Rule berkelas `Data-Cargo` tidak ada** di ketiga folder korpus maupun di `DDL\`.
⚠️ **Work owner menyatakan hanya ada DUA rule.** Dua pembacaan sama-sama konsisten dengan korpus dan
**korpus tidak dapat memutuskan** — rincian `09-pemanggil-countrateretrocov.md` §5.

### Di mana membukanya

**Rule yang dicari:**

| Kelas | Status di korpus |
| --- | --- |
| `ASM-FW-GISFW-Data-PropertyItem` | ada → `DDL\CountRateRetroCov.xml` |
| `ASM-FW-GISFW-Data-Aneka` | ada → `DDL\CountRateRetroCov(ANEKA).xml` |
| **`ASM-FW-GISFW-Data-Cargo`** | ⛔ **tidak ada** |

**Pemanggil yang harus dibuka** — tiap-tiap memuat **tiga** langkah `Call CountRateRetroCov` yang
seragam, seluruhnya bergerbang `.CoverageList(1).PremiumRetro==""`:

| Berkas | Alamat ketiga langkah `Call` |
| --- | --- |
| `NB FacIn\Activity\InsertFacoutProduction` | `RH_1.pySteps(8).pySteps(1)` · `RH_1.pySteps(9).pySteps(2).pySteps(2).pySteps(1)` · `RH_1.pySteps(10).pySteps(3).pySteps(1)` |
| `NB FacIn\Activity\InsertFacoutProductionEDM` | `RH_1.pySteps(10).pySteps(1)` · `RH_1.pySteps(11).pySteps(2).pySteps(2).pySteps(1)` · `RH_1.pySteps(12).pySteps(3).pySteps(1)` |

📌 Juga layak dibuka: `NB FacIn\Activity\CopyAllObjFacOutFireAneka_ACT`, alamat
`RH_1.pySteps(4).pySteps(2).pySteps(1).pySteps(9)` dan
`RH_1.pySteps(5).pySteps(2).pySteps(1).pySteps(1).pySteps(5)` — **dua panggilan tanpa gerbang sama
sekali**, dan **nol entri indeks** `Embed-Reference-Rule` meski memanggilnya dua kali.

### Apa yang harus dicatat

1. Apakah rule `CountRateRetroCov` berkelas **`Data-Cargo`** **ada sekarang**. Bila ada: ekspornya.
2. Bila **tidak ada**: apakah ia **pernah ada** (riwayat/audit trail rule), dan kapan dihapus.
3. Untuk tiap **tiga** langkah `Call` di kedua pemanggil produksi: **ke rule mana** langkah itu
   resolve saat dibuka — Pega menampilkan kelas yang di-resolve di samping nama activity.
4. Apakah langkah `Call` **ketiga** masih **dapat tercapai**, atau cabangnya sudah mati.
5. Untuk `CopyAllObjFacOutFireAneka_ACT`: **cabang mana memanggil varian mana** — namanya memuat
   "FireAneka", tetapi ⛔ **nama bukan bukti** (`PANDUAN-KERJA` §3) dan korpus tidak memisahkannya.

### Tergantung padanya

- **K-060** `[pertanyaan terbuka]` butir 1 — kemungkinan varian ketiga.
- **F06** — bagian *Kemungkinan varian KETIGA*; dan bila varian Cargo ada, **pembagi ketiga** harus
  diturunkan resolver **K-018** seperti dua varian lain.
- **T-7** di `..\_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md` — permintaan tertulis yang sepadan untuk IT.

---

## Yang **tidak** masuk berkas ini

Dicatat agar tidak dianggap terlewat:

1. **Status remark `//` pemanggil `GetLimitAkseptasi_Act`/`_Act2`** — sudah tercatat sebagai butir
   terbuka di `..\_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md` **P-12 catatan W-3**, lengkap dengan daftar
   pemanggil yang harus dicek. **Tidak diduplikasi di sini.**
2. **Apakah urutan `REPEATINGINDEX` = urutan eksekusi** sebagai **aturan umum** — tetap
   `[pertanyaan terbuka]` `[di luar korpus]`. ⚠️ Yang **sudah** dijawab work owner (21 September)
   hanyalah **nilai `Local.prorate` pada `CountRateRetroCov`**, bukan aturan umumnya
   (`01-pohon-langkah-countrateretrocov.md` §7D.4).
3. **Apakah penomoran `pySteps(n)` = urutan eksekusi** — `[di luar korpus]`; akan ikut terjawab bila
   U-1/U-2 terjawab, jadi tidak dijadikan butir tersendiri.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
