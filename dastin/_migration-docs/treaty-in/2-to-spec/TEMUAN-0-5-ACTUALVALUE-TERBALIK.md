# Temuan langkah 0.5 — `7-2-ACTUALVALUE.md` §1.2 dan §1.3 TERBALIK terhadap sumbernya

**Tanggal:** 24 September 2026 · **Lahir dari:** inventaris mekanis langkah 0.5
**Sifat:** koreksi atas berkas induk. **Dilaporkan, TIDAK disunting diam-diam.**

> Berkas `4-erd-dan-tabel-datar/7-2-ACTUALVALUE.md` disebut **satu kali** di seluruh korpus.
> Inventaris langkah 0.5 memunculkannya, dan pembacaan ulang atas sumbernya membalik
> kesimpulannya. Ini alasan langkah 0.5 ada.

---

## 1. Yang tertulis di berkas induk

§1.2 memberi tabel:

| `EDMState` | Dipotret? (menurut `7-2`) |
|---|---|
| `1` Internal Edit | **TIDAK** |
| `2` External Addendum | **TIDAK** |
| `3` Addendum Premium | **YA** |

dan §1.3 menyimpulkan: *"Untuk addendum jenis 1 dan 2, `TreatyInTemp` tidak pernah diisi pada
permintaan itu — dan isinya tetap disalin menimpa `ActualValue`."*

## 2. Yang terbaca dari sumbernya

`Activity/SaveTreatyIn_EDM_Act.xml`, langkah 2, `pyStepsPreCondParams` baris 1 — **dikutip utuh**:

```xml
<pyStepsPreCondParamsWhen>TreatyIn.EDMState=="3"</pyStepsPreCondParamsWhen>
<pyStepsPreCondParamsWhenTrue>1</pyStepsPreCondParamsWhenTrue>
<pyStepsPreCondParamsWhenTruePrms>jmp</pyStepsPreCondParamsWhenTruePrms>
<pyStepsPreCondParamsWhenFalse />
```

dan langkah 5 ber-`pyStepsBlockName = jmp`.

Langkah 3 dan 4 ber-`WhenTrue = 2`, `WhenFalse = 2`, `When` **kosong** — tanpa prasyarat.

### 2.1 Arti kode 1 — DIKALIBRASI, bukan diandaikan

Sapuan atas **kedua ekspor** (`Treaty In` dan `Treaty In Adjustment`), seluruh
`pyStepsPreCondParamsWhenTrue`/`WhenFalse`:

| Kode | Jumlah | Membawa nama blok (`…Prms`)? |
|---:|---:|---|
| **1** | **66** | **66 dari 66 — SELALU** |
| 2 | 18.412 | tidak pernah |
| 3 | 4.890 | tidak pernah |
| 4 | 412 | tidak pernah |
| 5 | 16 | tidak pernah |
| 6 | 55 | 2 dari 55 |

**Kode 1 satu-satunya yang selalu disertai nama blok**, dan pada **46** dari 66 nama blok itu
**ada di aktivitas yang sama**. Kode 2 adalah nilai yang muncul pada langkah **tanpa kondisi sama
sekali** — yaitu "lanjut".

> **Maka kode 1 = LOMPAT KE BLOK.** Ini terbaca dari sebaran, bukan dari kamus Pega mana pun.
>
> **Sisa 20 yang nama bloknya tidak ditemukan di aktivitas yang sama belum diperiksa sebabnya**
> dan **tidak** dihitung sebagai cacat di sini. Ia masuk daftar sisa.

## 3. Maka pembacaannya terbalik

| `EDMState` | Langkah 2-3-4 | `ActualValue` pada penyimpanan |
|---|---|---|
| `1` Internal · `2` External | **berjalan** | **ditimpa salinan utuh pohon utama, setiap kali simpan** |
| **`3` Addendum Premium** | **dilompati seluruhnya** — lompat ke blok `jmp` di langkah 5 | **tidak disentuh** — isinya tetap apa yang disunting pengguna |

Tiga akibat, dan ketiganya membalik kalimat di berkas induk:

1. **§1.2 terbalik.** Yang dipotret adalah jenis **1 dan 2**, bukan jenis 3.
2. **§1.3 terbalik, dan cacatnya berpindah alamat.** Bukan jenis 1 dan 2 yang menerima potret
   kosong-atau-basi; merekalah yang menerima potret **utuh dan segar**. Yang tidak pernah
   ditimpa adalah jenis **3** — dan untuk jenis 3 itu **bukan cacat, melainkan justru yang
   membuatnya bekerja**, sebab `ActualValue` adalah pohon yang disunting penggunanya.
3. **§1.4 berpindah alamat juga.** Sarang `ActualValue(ActualValue(…))` tumbuh pada addendum
   **jenis 1 dan 2** — merekalah yang menjalankan langkah 4 berulang kali — bukan pada addendum
   premi. **Uji AM-3** tetap sah; yang berubah **kontrak yang disasarnya**.

## 4. Apa yang TIDAK berubah

* **§2** (empat halaman potret dan bentuk pengisian masing-masing) tidak tersentuh.
* **§3** — `BESARAN_DAPAT_DISESUAIKAN` bersumber pada sembilan besaran `ValueDifference` — tidak
  tersentuh: ia dibaca dari sisi penghitung selisih, bukan dari sisi potret.
* **Uji AM** tetap perlu dijalankan; **pertanyaannya** yang berubah, dari *"apakah potret jenis 1
  dan 2 kosong"* menjadi *"apakah `ActualValue` jenis 3 pernah terisi sama sekali"*.

## 5. Dan ia menguatkan `TDA-16` di modul Adjustment

`TDA-16` menemukan `TreatyEDMDifferencePremium` beriterasi atas `TreatyIn.EGNPI`, bukan atas
`TreatyIn.ActualValue.EGNPI`. Temuan ini menambahkan sisi keduanya: pada addendum premi,
`ActualValue` **memang** tempat nilai barunya tinggal — dan mesin selisih premi **memang** tidak
pernah menengoknya. Keduanya sisi yang sama dari satu cacat.

> **Tidak ada yang disunting di `PENGETAHUAN.md` Adjustment maupun di `7-2-ACTUALVALUE.md` oleh
> berkas ini.** Usulan suntingannya ada di §6.

## 6. Usulan suntingan — menunggu persetujuan pemilik proses

| Berkas induk | Bagian | Usulan |
|---|---|---|
| `4-erd-dan-tabel-datar/7-2-ACTUALVALUE.md` | §1.2 tabel | balik kolom "Dipotret?": `1` **YA**, `2` **YA**, `3` **TIDAK** |
| idem | §1.3 | tulis ulang: yang menerima potret utuh adalah jenis 1 dan 2; jenis 3 dilompati dan itu disengaja |
| idem | §1.4 | sarang tumbuh pada jenis 1 dan 2 |
| idem | kepala berkas | *"Jawaban: SELURUHNYA — dan hanya untuk satu dari tiga jenis"* → **dua dari tiga** |
| `CONTEXT.md` §2.1 | — | tambahkan: **prasyarat yang hidup pun harus dibaca bersama KODE TINDAKANNYA.** `pyStepsPreCondParamsWhen` sendirian tidak memberi tahu apakah langkahnya berjalan atau dilompati |

## 7. Pemilik dan penagih

| | |
|---|---|
| **Siapa menutup** | pemilik proses — ia koreksi atas berkas induk, dan berkas induk tidak disunting tanpa persetujuan |
| **Yang menagih** | `C3` di grilling Adjustment berdiri di atas pembacaan ini. Selama §1.2 berbunyi terbalik, pembaca berikutnya yang membuka `7-2` akan menyimpulkan kebalikan dari yang dipakai `C3` |
