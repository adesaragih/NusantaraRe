# PROMPT — Keadaan EDM Treaty In, sekaligus koreksi ronde 1

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
> Jalankan **sebelum** `to-spec` EDM Treaty In.
>
> ⚠️ Ronde ini **menyunting** satu berkas dan **menulis** satu berkas baru.

---

## 0. LINGKUP — DIKUNCI

Modul: **`D:\XML\RNM_BRD\EDM Treaty In`**. `NB Treaty In` hanya sebagai pembanding.

`D:\XML\RNM_BRD\` **READ-ONLY**. Menulis hanya ke `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang**.

⛔ `grilling-ronde-1.md` EDM **tersegel** — jangan disunting. Koreksinya hidup di berkas keadaan.

---

## 1. DUA PEKERJAAN

| | Keluaran |
| --- | --- |
| **A** — tulis berkas keadaan | `.scratch\edm-treaty-in\KEADAAN-EDM-TREATY-IN.md` |
| **B** — koreksi lembar jawaban | `.scratch\edm-treaty-in\PERTANYAAN-RONDE-1.md` |

Keduanya wajib selesai sebelum `to-spec` dijalankan.

---

## 2. PEKERJAAN A — BERKAS KEADAAN

Bentuknya mengikuti `.scratch\nb-treaty-in\KEADAAN-NB-TREATY-IN.md`. Bab wajib:
urutan kewenangan · angka yang berlaku · lingkup · sumber data · keputusan yang mengikat ·
koreksi atas ronde 1 · cara menguji ulang.

### Angka yang harus diuji ulang, dua cara

| Pernyataan | Angka |
| --- | ---: |
| berkas dalam modul | 163 |
| ukuran modul | 17.639.259 B |
| `Activity` | 66 |
| langkah `Property-Set` | 378 |
| pasangan `PropertiesName`=`PropertiesValue` terisi | 1.062 |
| **`Activity` terjangkau / yatim** | **66 / 0** |
| langkah terjangkau / yatim | 378 / 0 |

⚠️ **Jangan percaya angka ini.** Sebagian dihitung sesi asisten dengan satu pola saja, dan sesi itu
sudah terbukti salah sekali hari ini. Ronde verifikasi NB mengukur 380 langkah dan 1.309 pasangan
untuk modul yang sama — **selisih itu harus dijelaskan**, bukan dipilih salah satunya.

### Yang membedakan EDM dari NB — wajib dinyatakan

⭐ **Nol aturan yatim.** NB punya 28 yatim dari 92 `Activity`; EDM nol dari 66. ⇒ **lingkup EDM
tidak akan menyusut** seperti NB menyusut 60 %. Nyatakan ini terang-terangan supaya tidak ada yang
mengharapkan penyusutan serupa.

⭐ **Dua tipe rule yang tidak dimiliki NB**: `ConnectREST` dan `SystemSettings`.

⭐ **Alur berkelas kerja sendiri** — `…WORK-ENDORSEMENTTREATY`. Endorsemen jenis kasus tersendiri,
dan alurnya menyebut **tiga antrean** — sejalan keputusan tangga tiga jenjang pada **P13**.

### Keputusan yang sudah mengikat dan berlaku juga di sini

Dua belas ketetapan pada `KEADAAN-NB-TREATY-IN.md` Bab 4 berlaku di EDM **kecuali terbukti
sebaliknya**. Ditambah keputusan EDM sendiri: **P50** sampai **P60**, seluruhnya terjawab.

Bila EDM berperilaku berbeda pada salah satu ketetapan NB, **itu temuan** dan wajib dinyatakan.

---

## 3. PEKERJAAN B — KOREKSI RONDE 1

Empat pernyataan ronde 1 kini terbukti keliru. Catat koreksinya di **berkas keadaan Bab koreksi**,
dan perbarui jawaban terkait di lembar pertanyaan. ⛔ Berkas grilling tetap tidak disentuh.

| Bunyi ronde 1 | Yang benar |
| --- | --- |
| *"isi langkah penetapan nilai tidak terekspor"* | **ada** — di `PropertiesName`/`PropertiesValue` **tanpa awalan `py`**. 379 dari 380 langkah membawa isinya *(99,7 %)* |
| *"pemeriksaan keberhasilan hanya memeriksa FacIn"* | `When\IsSuccessHitService.xml` memeriksa **keduanya**, tiga cabang: `FacIn = 1` **dan** (`FacOut = ""` **atau** `FacOut = 1`). Ronde 1 membaca `pyConditionString`, bukan bentuk yang dijalankan — jebakan **P23** |
| *"beda hanya metadata ekspor"* untuk 69 aturan | **enam berkas** punya catatan pengembang berbeda, empat di antaranya di dalam kelompok 69 itu |
| *"`BusinessType_DeT` tinggal 9 baris dari 34"* | **36 baris utuh**; `pyRowNum` pada `pyOrConditions` berbasis nol. Terverifikasi ulang dari **data produksi**: `GroupPanel=006` + `BusinessOldId=01` menghasilkan `FireStyle2`, dan dokumennya memang berisi itu |

### Dan satu yang perlu diperiksa, bukan diasumsikan

`P60` dijawab dari korpus. **Periksa bahwa jawabannya sudah tercatat** dan angkanya benar:
8 langkah EDM lawan 7 NB, langkah tambahan `.SpreadingRiskList(1).SplitRNMSharePct = 100`,
presisi 20 lawan 10.

---

## 4. TIGA TEMUAN DARI DATA PRODUKSI YANG MENYENTUH EDM

Tiga contoh `DATA_JSON` sungguhan diterima work owner 22 September 2026. Dua di antaranya
menyentuh EDM langsung — catat di berkas keadaan.

**Bentuk dokumen mengikuti jenis treaty.** `QuotationData.ProportionalType` bernilai
`"Proportional"` lawan `"NonProportional"`; `IsNewPolicyNonProp` `"0"` lawan `"1"`. Daftar
bersarangnya berbeda-beda, dan `ListInstallment` **bersarang** pada satu bentuk tetapi **datar**
pada bentuk lain. Pengurai seragam akan patah.

**`OldData` muncul di dokumen polis baru**, bukan hanya di endorsemen — berisi `TreatyXOLList`.
Ini menyentuh **P57** dan **P58**: periksa apakah keputusan di sana masih berdiri.

**Galat uang tersimpan di data produksi** — ekor sebesar 2,76 × 10⁻⁷ pada nilai premi, lahir di
rantai perhitungan, bukan di penyimpanan. Tercatat sebagai butir terbuka pada **P29** NB dengan
tiga pilihan; **belum dipilih work owner**. Nyatakan bahwa butir itu berlaku juga untuk EDM.

⛔ **Jangan menyimpan contoh JSON apa adanya** — ketiganya memuat nama orang.

---

## 5. DISIPLIN

Setiap angka membawa **perintah yang menghasilkannya**. Sensus dihitung **dua cara berbeda**.

Penanda: `[terverifikasi]` · `[dugaan]` · `[terbuka]` · `[data DBA]` · `[keputusan work owner]` ·
`[penyimpangan sadar]`.

⛔ Jangan menutup butir terbuka. Jangan membuat ADR baru. Jangan menulis kode atau DDL.
Jangan menyalin nilai berupa nama orang.

### Enam jebakan yang sudah menjerat proyek ini

1. **Menebak nama tag.** `pyPropRef`, `pyPropertiesValue`, `pyStepsPage`, `pyCriteriaValue`,
   `pyResult`, `pyShapeName` — semuanya pernah dicari dan tidak ada. **Periksa daftar tag lebih
   dulu.** Dan ingat: `Activity` memakai tag **tanpa** awalan `py`, `Flow` memakai **dengan** awalan.
2. `<rowdata REPEATINGINDEX="n"/>` yang menutup sendiri tidak tertangkap pola
   `<rowdata...>(.*?)</rowdata>`.
3. `pyRowNum` pada `pyOrConditions` berbasis **nol**; baris tabel berbasis satu.
4. Memasangkan dua daftar menurut **urutan**, bukan menurut blok yang sama.
5. Membaca `pyConditionString` — keterangan manusia — alih-alih bentuk yang dijalankan.
6. Menyatakan sesuatu nihil tanpa menyebut lingkup penelusuran.

Jalur korpus memuat spasi — pakai Python, bukan loop shell.

---

## 6. BAB WAJIB — TELEMETRI EKSEKUSI

Di akhir berkas keadaan. Ukur dari luar bila bisa:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Bila tidak, katakan begitu. Bila hasil pengurangan baseline, katakan itu juga.

---

## 7. YANG MENANDAKAN RONDE INI BERHASIL

- Ketujuh angka diuji ulang dua cara, dan **selisih 378/380 serta 1.062/1.309 dijelaskan**, bukan
  dipilih salah satunya.
- Keempat koreksi atas ronde 1 tercatat di berkas keadaan, dengan bunyi lamanya **dikutip**.
- Perbedaan EDM dari NB dinyatakan terang-terangan — terutama **nol yatim**.
- Tiga temuan dari data produksi tercatat, dan dampaknya ke **P57** dan **P58** diperiksa.
- Berkas grilling **tidak tersentuh**.
- Berkas keadaan memuat bab **cara menguji ulang** dengan kode siap jalan.
- Nol butir ditutup sendiri, nol ADR baru, nol nama orang tersalin.

---

*Disusun 22 September 2026. Ronde berikutnya: `to-spec` EDM Treaty In, yang membaca berkas keadaan
ini sebagai sumber keadaan terukur.*
