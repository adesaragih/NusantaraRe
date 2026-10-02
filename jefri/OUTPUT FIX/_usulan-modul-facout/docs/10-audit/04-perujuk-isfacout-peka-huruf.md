# Audit — perujuk `IsFacout` dihitung ulang secara PEKA HURUF

> **Tanggal:** 21 September 2026 · **Memicu:** amandemen dasar pengukuran **K-019**
> **Tidak mengubah keputusan K-019.** Fitur fac out tetap usang secara bisnis dan **kodenya tetap
> diport apa adanya**. Yang dikoreksi di sini **hanya dasar pengukuran jumlah perujuknya**.

---

## 1. Sebab audit ini

K-019 mencatat **10 perujuk** `IsFacout`, diukur dengan:

```powershell
Select-String -Path "D:\migrasi\RNM\NB FacIn\*\*.xml" -Pattern '\bIsFacout\b' -List
# -> 11 berkas; K-019 mengurangi 1 (berkas definisi When\IsFacout) -> 10 perujuk
```

⛔ **`Select-String` dan operator `-match` di PowerShell TIDAK peka huruf secara bawaan.** Korpus
memuat **tiga ejaan berbeda** yang seluruhnya tertangkap perintah di atas:

| Ejaan | Apa sebenarnya |
| --- | --- |
| `IsFacout` | **nama rule** `Rule-Obj-When` |
| `isFacOut` · `isFacOutLoc` | **variabel / parameter lokal** di dalam Activity |
| `ISFACOUT` | ejaan huruf besar di Section/Harness tampilan |

Aritmetika K-019 benar; **dasarnya** yang keliru.

---

## 2. Tiga angka, tiga definisi — semuanya dari korpus yang sama

`[terverifikasi]` Folder `NB FacIn`:

| Definisi pencarian | Berkas | Perintah |
| --- | ---: | --- |
| `\bIsFacout\b` **tak peka huruf** (sebagaimana K-019) | **11** | baris di §1 |
| `\bIsFacout\b` **PEKA HURUF** | **3** | `Select-String … -List -CaseSensitive` |
| substring `IsFacout` tak peka huruf, tanpa `\b` | **17** | `[regex]::new('IsFacout',IgnoreCase)` |

```powershell
# menghasilkan 11 lalu 3
(Select-String -Path "D:\migrasi\RNM\NB FacIn\*\*.xml" -Pattern '\bIsFacout\b' -List).Count
(Select-String -Path "D:\migrasi\RNM\NB FacIn\*\*.xml" -Pattern '\bIsFacout\b' -List -CaseSensitive).Count
```

**Mengapa 17 > 11:** tanpa `\b`, substring `isFacOut` ikut tertangkap di dalam kata majemuk seperti
`ViewGeneralPol**isFacOut**`, `ViewTabGroupPol**isFacOut**Savior`, `ConvertPropViewPol**isFacOut**_Act`
— kebetulan perangkaian huruf, **bukan** rujukan apa pun. Angka 17 dicantumkan hanya agar terlihat
bahwa pencarian substring polos lebih buruk lagi, **bukan** untuk dipakai.

---

## 3. Hasil peka huruf — **3 berkas per folder**, bukan 11

`[terverifikasi]` Peka huruf, per folder, dengan tag pembawa tiap kecocokan:

| Folder | Berkas | Tag pembawa | Perujuk sejati? |
| --- | --- | --- | :-: |
| NB FacIn | `Section\InputCoverageFire` | `pyContainerVisibleWhen`, `pyRuleName` | ya |
| NB FacIn | `Section\InputCoverageFire_IsUW` | `pyDisplayWhen`, `pyRuleName` | ya |
| NB FacIn | `When\IsFacout` | `pxTabLabel`, `pyBlockName`, `pyJavaClassName`, `pyLabel`, `pyRuleName` | **tidak — ini definisinya** |
| RNW Fac In | *identik NB* (ketiga berkas sama) | idem | idem |
| Endorsment | `Section\InputCoverageAneka_FacIn` | `pyContainerVisibleWhen`, `pyDisplayWhen`, `pyRuleName` | ya |
| Endorsment | `Section\InputCoverageFire` | `pyContainerVisibleWhen`, `pyRuleName` | ya |
| Endorsment | `When\IsFacout` | idem NB | **tidak — definisinya** |

⚠️ **Himpunannya tidak sama antar folder.** NB/RNW memakai `InputCoverageFire_IsUW`; Endorsment
memakai `InputCoverageAneka_FacIn`. Jumlahnya kebetulan sama, isinya tidak.

### Isi tag pemanggilnya — **kutipan verbatim ketiga folder**

`[terverifikasi]` Diperbarui 21 September 2026 dengan kutipan lengkap dari **ketiga** folder.

**NB FacIn** — 3 berkas peka huruf:

```
Section\InputCoverageFire            peka huruf=2  tak peka huruf=3
  <pyContainerVisibleWhen> = (.FlagDelete != 1 && pyPortal.IsShowOpenPolicyMarine!= true
                              && pyWorkPage.IsInFacRetro != 1 && !IsClaim && !IsRenewal
                              && !IsUW && !IsFacout)&&1==2
  <pyRuleName>             = IsFacout

Section\InputCoverageFire_IsUW       peka huruf=2  tak peka huruf=3
  <pyDisplayWhen>          = !IsUW && !IsFacout
  <pyRuleName>             = IsFacout

When\IsFacout                        peka huruf=5  tak peka huruf=16
  <pyRuleName> <pxTabLabel> <pyBlockName> <pyLabel> = IsFacout
  <pyJavaClassName> = Rule_Obj_When_ASM_FW_GISFW_Data_IsFacout_When_<TS>
```

**RNW Fac In** — 3 berkas, **identik NB** (berkas yang sama, isi kutipan sama persis):

```
Section\InputCoverageFire            peka huruf=2  tak peka huruf=3
  <pyContainerVisibleWhen> = (.FlagDelete != 1 && pyPortal.IsShowOpenPolicyMarine!= true
                              && pyWorkPage.IsInFacRetro != 1 && !IsClaim && !IsRenewal
                              && !IsUW && !IsFacout)&&1==2
Section\InputCoverageFire_IsUW       peka huruf=2  tak peka huruf=3
  <pyDisplayWhen>          = !IsUW && !IsFacout
When\IsFacout                        peka huruf=5  tak peka huruf=16
```

**Endorsment Fac In** — 3 berkas, **himpunannya BERBEDA**:

```
Section\InputCoverageAneka_FacIn     peka huruf=3  tak peka huruf=4
  <pyContainerVisibleWhen> = pyPortal.IsShowOpenPolicyMarine!= true
                              && pyWorkPage.IsInFacRetro != 1 && !IsClaim && !IsRenewal
                              && !IsUW && !IsFacout
  <pyDisplayWhen>          = !IsUW && !IsFacout && .FlagDelete!=1
  <pyRuleName>             = IsFacout

Section\InputCoverageFire             peka huruf=2  tak peka huruf=3
  <pyContainerVisibleWhen> = (.FlagDelete != 1 && pyPortal.IsShowOpenPolicyMarine!= true
                              && pyWorkPage.IsInFacRetro != 1 && !IsClaim && !IsRenewal
                              && !IsUW && !IsFacout)&&1==2
  <pyRuleName>             = IsFacout

When\IsFacout                         peka huruf=5  tak peka huruf=16
```

### Mana yang aktif, mana yang dinetralkan

`[terverifikasi]`

| Folder | Perujuk | Ekspresi | Status |
| --- | --- | --- | :-: |
| NB · RNW | `InputCoverageFire` | ditutup **`&&1==2`** | ⛔ **DINETRALKAN** — selalu salah |
| NB · RNW | `InputCoverageFire_IsUW` | `!IsUW && !IsFacout` | ✅ **AKTIF** |
| Endorsment | `InputCoverageFire` | ditutup **`&&1==2`** | ⛔ **DINETRALKAN** |
| Endorsment | `InputCoverageAneka_FacIn` | dua tag, **tanpa** `&&1==2` | ✅ **AKTIF** (dua kondisi) |

📌 **Setiap folder punya tepat 1 perujuk aktif dan 1 perujuk mati.** Yang mati selalu
`InputCoverageFire`, karena konjungsi `&&1==2` membuat seluruh kondisinya salah apa pun nilai
`IsFacout`.

### Kelas resolusinya tidak seragam

`[terverifikasi]` Resolusi dibawa tag **`<pxRuleClassName>`** di dalam blok `<pxObjClass>` =
`Embed-Reference-Rule`, dengan `<pxRuleObjClass>` = `Rule-Obj-When`. ⚠️ `<pyClassName>` **kosong** —
mencarinya di tag itu akan gagal.

| Berkas | `<pxRuleClassName>` |
| --- | --- |
| NB `Section\InputCoverageFire` | `ASM-FW-GISFW-Data-PropertyItem` |
| NB `Section\InputCoverageFire_IsUW` | `ASM-FW-GISFW-Data-PropertyItem` |
| Endorsment `Section\InputCoverageAneka_FacIn` | **`ASM-FW-GISFW-Data-Aneka`** |

**Seluruh pemakaian berbentuk `!IsFacout` pada kondisi tampilan UI** (`pyContainerVisibleWhen` /
`pyDisplayWhen`) — menyembunyikan komponen layar, **bukan** menggerakkan alur. Ini sejalan dengan
`GLOSARIUM.md`: isi rule-nya menguji **banding**, bukan fac out.

### ⛔ Temuan tambahan — satu perujuk dinetralkan tautologi

`[terverifikasi]` Ekspresi `NB\Section\InputCoverageFire` ditutup **`&& 1==2`**. Konjungsi itu selalu
salah, sehingga **seluruh kondisi selalu salah** dan komponen itu **tidak pernah tampil**, apa pun
nilai `IsFacout`. Perujuknya ada di teks, **pengaruhnya nol**.

Ini **kandidat K-046** (diport apa adanya). Polanya sama dengan tautologi pembanding limit yang sudah
tercatat di `PANDUAN-KERJA` §6. ⚠️ **Jangan "dirapikan"** — menghapus `&& 1==2` akan **memunculkan**
komponen yang selama ini tersembunyi, dan itu perubahan perilaku.

---

## 4. Delapan berkas yang K-019 hitung sebagai perujuk, ternyata **bukan**

`[terverifikasi]` Kedelapan Activity ini ber-`IsFacout` **peka huruf = 0**. Yang ada di dalamnya
adalah **variabel lokal** `Local.isFacOut` / `Local.isFacOutLoc` dan **deklarasi parameter**
`<pyParametersParamName>` = `isFacOut` / `isFacOutLoc`:

| Berkas (folder NB) | Tag pembawa `isFacOut` | Ejaan persis `IsFacout` |
| --- | --- | ---: |
| `SetDataFacOut_Act` | `pyParametersParamName` | **0** |
| `SetDataFacOutFire_Act` | `pyParametersParamName` · `PropertiesName` (`Local.isFacOut`, `Local.isFacOutLoc`) · `pyStepsPreCondParamsWhen` (`Local.isFacOut==1`, `Local.isFacOutLoc==1`) | **0** |
| `SetDataFacOutAnekaGolf_Act` | idem | **0** |
| `SetDataFacOutCargoMBU_Act` | idem | **0** |
| `SetDataFacOutPATravel_Act` | idem | **0** |
| `SetValidateDateUW_PostAct` | `pyParametersParamName` **saja** | **0** |
| `SetValidateDate_PostAct` | `pyParametersParamName` **saja** | **0** |
| `SumFacOutPA_Act` | `pyParametersParamName` · `PropertiesName` · `pyStepsPreCondParamsWhen` | **0** |

📌 `Local.isFacOut==1` di `<pyStepsPreCondParamsWhen>` **terlihat persis seperti** pemanggilan rule
`When` di tag yang sama. Pembedanya hanya awalan `Local.` dan kapitalisasi huruf pertama. Inilah
jebakan yang menghasilkan angka 10.

---

## 5. Angka perujuk yang sebenarnya

| | K-019 (lama) | Audit ini |
| --- | ---: | ---: |
| Berkas cocok | 11 | **3** per folder |
| Dikurangi berkas definisi | −1 | −1 |
| **Perujuk** | **10** | **2** per folder |
| Di antaranya **berpengaruh** | — | **1** di NB/RNW (satu lagi dinetralkan `&& 1==2`) · **2** di Endorsment |

⛔ **Kedelapan Activity `SetDataFacOut*` / `SetValidateDate*` / `SumFacOutPA_Act` BUKAN perujuk
`IsFacout`.** Mereka tetap **inti jalur Fac Out** — tetapi lewat flag `.IsFacRetro`, bukan lewat rule
ini. Lihat `09-facout\01-temuan-dan-rancangan-facout.md` §0.1.

---

## 6. Bahan amandemen K-019 — **untuk work owner, belum diputuskan**

**Yang TIDAK berubah:**

- Keputusan intinya: fitur fac out usang secara bisnis, **kode tetap diport apa adanya**.
- Konsekuensi 1–4 di K-019 (status ⏸ dicabut · tidak ada klaim "kode mati" · nama dipertahankan demi
  ketertelusuran · penghapusan butuh keputusan tersendiri).
- Catatan rekonsiliasi: volume kasus lewat jalur ini mungkin nol; **uji dengan fixture**.

**Yang perlu diamandemen:**

1. Angka **10 perujuk → 2 perujuk per folder**, beserta tabelnya.
2. Perintah audit diberi **`-CaseSensitive`**, dengan catatan mengapa.
3. Tabel kelompok perujuk diganti: kelompok "Activity penyetel data fac out per lini bisnis" (5),
   "Post-activity validasi tanggal" (2) dan "Activity penjumlah fac out PA" (1) **dihapus dari daftar
   perujuk** — kedelapannya memakai variabel lokal `isFacOut`, bukan rule.
4. Ditambahkan: himpunan perujuk **berbeda antar folder** (NB/RNW vs Endorsment).
5. Ditambahkan: satu perujuk NB/RNW **dinetralkan `&& 1==2`**, kandidat K-046.

⚠️ **Konsekuensi yang perlu ditimbang work owner:** kalimat K-019 "menghapus predikat yang masih
dirujuk akan mematikan gerbang lima jalur penyetelan data fac out sekaligus" **tidak lagi didukung
bukti** — kelima jalur itu tidak pernah memanggil `IsFacout`. Alasan mempertahankannya kini bertumpu
pada dua kondisi tampilan UI, bukan lima gerbang alur. **Keputusannya tetap milik work owner**; audit
ini hanya meluruskan dasarnya.

---

## 7. Catatan samping — keluarga nama yang mirip

`[terverifikasi]` Empat properti berbeda hidup berdampingan dan **mudah tertukar**:

| Nama | Catatan |
| --- | --- |
| `.IsFacRetro` | flag pemicu Fac Out, nilai telanjang `0`/`1` |
| `.IsInputFacRetro` | gerbang re-input saat UW |
| `.IsFacRetroOffer` | disetel **`"0"` berkutip** di `OfferFacOut_PreAct` — gaya berbeda dari `.IsFacRetro` |
| `pyWorkPage.IsInFacRetro` | dipakai di kondisi tampilan Section (§3) |

Jumlah berkas folder NB yang memuat tiap nama (peka huruf): `IsFacRetro` **37**, `IsInputFacRetro`
**15**, `IsFacRetroOffer` **14**, `IsInFacRetro` **4**.

⚠️ Angka `IsFacRetro` **37 bersifat substring** — ia ikut menghitung berkas yang sebenarnya memuat
`IsFacRetroOffer`. Keduanya **tidak** boleh dikurangkan begitu saja; pemisahan tepatnya belum diukur.
`[dugaan]`

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
