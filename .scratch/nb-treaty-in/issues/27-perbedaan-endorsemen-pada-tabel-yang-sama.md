# 27: Perbedaan endorsemen pada tabel yang sama

> ## ⛔⛔ DIGANTIKAN — 23 September 2026
>
> **Tiket ini tidak dikerjakan.** Ia digantikan oleh
> **`.scratch\edm-treaty-in\issues\08`**.
>
> ### Kenapa
>
> Brief ronde tiket sebelumnya berlingkup **dua spec sekaligus**, sehingga pekerjaan **endorsemen**
> jatuh ke folder tiket **polis baru**. Ronde berikutnya membaginya ulang ke folder yang benar,
> lebih halus, dan — yang menentukan — **setiap tiket menyebut tiket NB mana yang membuat
> tabelnya**. Tanpa gate itu, tiket endorsemen bisa dikerjakan sebelum tabelnya ada.
>
> ⚠️ **Kekeliruan lingkup ini milik penyusun brief, bukan pelaksana ronde ini.**
>
> ⭐ Berkas ini **tidak dihapus** supaya jejaknya tidak hilang — tanpa ini, nomor tiket NB akan
> tampak melompat dari 23 ke 29 tanpa sebab.
>
> ### Yang berlaku
>
> | | |
> | --- | --- |
> | Tiket penyimpanan **NB** | `nb-treaty-in\issues\` **16–23** |
> | Tiket penyimpanan **EDM** | `edm-treaty-in\issues\` **01–11** |
>
> ⛔ Isi di bawah dibiarkan apa adanya sebagai catatan, **bukan sebagai pekerjaan.**

---


**Status:** ~~blocked~~ — **DIGANTIKAN**
**Blocked by:** **19** · **24** · ⛔ `[work owner]` **beda dagang tabel pecahan penyebaran** · ⛔ `[data DBA]` **presisi fisik**
**Menutup:** EDM AC **35–38** · EDM AC **54–58** *(9 AC)*
**Sumber:** `edm-treaty-in\spec-penyimpanan-relasional.md` ID-5..ID-7 · ID-39..ID-42 · ID-46

## Hasil & nilai pengguna

Tabel dasarnya sama dengan polis baru; yang berbeda **isinya**. Tiket ini menyatakan perbedaan itu
dengan tepat, dan membuktikan bahwa bentuk tabel dasar **tidak berubah**.

⭐⭐ **Nol tabel dasar dirancang ulang.** Bila pekerjaan ini menuntutnya, itu **temuan** yang wajib
dilaporkan, bukan perubahan yang dijalankan diam-diam.

## Yang dibangun

| Perbedaan | Endorsemen | Polis baru |
| --- | --- | --- |
| presisi pembagian penyebaran | **20** | **10** |
| baris pertama penyebaran | nilai bawaan **100** bila persentasenya kosong | — |
| rincian angsuran bertingkat | **hidup** pada bentuk non-proporsional | tidak ada |
| asal rincian angsuran | ⭐ dapat berasal dari **salinan master kontrak** | diketik |

⚠️ **Perbedaan presisi ditiru apa adanya**, dan **wajib punya test tersendiri** — ia perilaku lama
yang disalin sadar, bukan cacat.

Ditambah pernyataan bentuk: cacah tabel pada kedua bentuk endorsemen, panjang minimum kolom
keterangan, dan skala kolom uang.

## Batas — yang TIDAK termasuk

⛔ Rantai generasi — tiket **24**.
⛔ Proyeksi selisih — tiket **26**.
⛔ `CREATE TABLE` — presisi fisik dicocokkan DBA **di dalam** tiket ini.

## Cara mengujinya

Lewat seam `repository`. ⭐ Uji presisi dijalankan **berdampingan** — kasus yang sama disimpan
sebagai polis baru dan sebagai endorsemen, dan hasilnya **harus berbeda** sesuai ketetapan.

## Acceptance criteria

- [ ] **AC 35** — penyebaran endorsemen memakai presisi **20**; polis baru **10**
- [ ] **AC 36** — baris pertama penyebaran menerima nilai bawaan **100** ketika persentasenya kosong
- [ ] **AC 37** — rincian angsuran bertingkat **tersimpan** pada endorsemen non-proporsional
- [ ] **AC 38** — rincian dari salinan master kontrak tersimpan **sama** seperti yang diketik
- [ ] **AC 54** — bentuk tabel dasar **tidak berubah**
- [ ] **AC 55** — cacah tabel pada kedua bentuk endorsemen sesuai ketetapan
- [ ] **AC 56** — tabel pecahan penyebaran menyimpan empat medan yang ditetapkan
- [ ] **AC 57** — kolom keterangan berpanjang sekurangnya batas yang ditetapkan
- [ ] **AC 58** — nilai berdesimal lebih dari skala kolom **dibulatkan, bukan ditolak**

## ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Beda dagang tabel pecahan penyebaran dari tabel penyebaran biasa belum
dijelaskan.** AC 56 menyatakan isinya, tetapi **kapan baris masuk ke yang satu dan bukan ke yang
lain** belum dapat diuji.

`[data DBA]` **Presisi fisik kolom uang** — sama dengan penahan tiket **18**.

⚠️ **Dan satu pertentangan di dalam spec, dicatat karena spec tidak boleh disunting:**
`[terverifikasi]` AC **49** pada spec yang sama masih menuntut **sembilan** desimal, sementara
AC **58** menetapkan **delapan**. ⛔ **Tidak diputuskan di tiket ini.**