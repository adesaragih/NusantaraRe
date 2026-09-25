# ⛔ DIBATALKAN — surat ini tidak jadi dikirim

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22
>
> **Surat ini dibatalkan sebelum dikirim. Tidak ada yang perlu Anda kerjakan.**
>
> Surat ini disusun atas dasar pernyataan tim migrasi bahwa isi langkah penetapan nilai tidak
> ikut terekspor. **Pernyataan itu keliru, dan kekeliruannya milik tim migrasi sepenuhnya.**
>
> Isi langkah **ada di dalam ekspor yang sudah kami terima**, di tag `PropertiesName` dan
> `PropertiesValue` — **tanpa awalan `py`** — di dalam `pyParamArray` milik tiap langkah.
> Tim migrasi memeriksa `pyPropRef`, `pyPropertiesName`, dan `pyPropertiesValue` — **dengan**
> awalan `py`. Ketiganya memang kosong, dan kekosongan itu terlanjur dipercaya.
>
> ⭐ **Sebabnya dua huruf.**

---

## Yang terukur sesudah diperiksa ulang

| | NB Treaty In | EDM Treaty In |
| --- | ---: | ---: |
| langkah `Property-Set` | 939 | 380 |
| ⭐ **membawa isinya sendiri** | **938** *(99,9 %)* | **379** *(99,7 %)* |
| pasangan nama=nilai terisi | **2.481** | **1.309** |
| ⭐ pada 51 aturan / 269 langkah yang diminta surat ini | ⭐ **707** | — |
| tanda nilai terpotong, tujuh uji | **0** | **0** |

Rinciannya di `.scratch\nb-treaty-in\VERIFIKASI-P18.md`.

---

## Bunyi permintaan yang ditarik — dikutip utuh, tidak dihapus

> **Perihal:** isi langkah penetapan nilai tidak ikut terkirim
> **Lampiran:** `LAMPIRAN-P18-ATURAN-DIMINTA.md` — 51 nama aturan
>
> ### Yang kami minta
>
> *Ekspor ulang **51 aturan** pada lampiran, **dengan isi langkah disertakan**.*
>
> *Ketika ekspor itu dibuka, kami dapat melihat **bahwa** sebuah langkah menetapkan nilai, dan
> **urutannya** — tetapi tidak satu pun **rumus atau nilai** yang ditetapkannya ikut terkirim.*
>
> *Yang tersisa hanya potongan salinan sementara dari layar penyunting, dan potongan itu terputus
> di tengah kalimat sehingga tidak dapat dipakai.*
>
> ### Kenapa ini penting
>
> *Langkah-langkah itu adalah tempat perhitungan uang berada: premi, bagian, pengurangan, dan
> pajak.*
>
> *Tanpa isinya, perhitungan sistem lama **tidak dapat ditiru — hanya ditebak**. Dan tebakan pada
> perhitungan premi tidak menimbulkan pesan galat; selisihnya baru muncul di laporan keuangan
> berbulan-bulan kemudian, saat sudah sulit ditelusuri asalnya.*
>
> *Ada satu akibat lagi yang baru kami ketahui. Sembilan medan pada layar Kepala Departemen
> **wajib diisi tetapi terkunci** — pengguna tidak dapat mengetiknya sendiri. Yang mengisinya
> adalah langkah-langkah ini. Selama isinya belum ada, layar itu **tidak akan dapat disimpan sama
> sekali**.*
>
> ### Lingkupnya sudah kami persempit lebih dulu
>
> *Kami tidak ingin Anda mengerjakan yang tidak kami butuhkan. Folder yang kami terima memuat 92
> aturan dengan 938 langkah, tetapi sebagian besar bukan milik modul ini:*
>
> | Tahap penyaringan | Aturan | Langkah |
> | --- | ---: | ---: |
> | Seluruh isi folder yang kami terima | 92 | 938 |
> | Dibuang — tidak terjangkau dari titik masuk mana pun | −28 | −561 |
> | Dibuang — tidak dibangun ulang karena penyimpanan JSON ditinggalkan | −5 | −109 |
> | **Yang kami minta** | **51** | **268** |
>
> ***Turun 72 %** dari angka yang mungkin pernah Anda dengar dari kami sebelumnya.*
>
> ### Kalau tidak bisa seluruhnya
>
> ***Kerjakan dari urutan teratas lampiran.** Kami sudah mengurutkannya dari yang paling berat:*
> *10 aturan teratas mencakup **56 %**, 20 teratas **76 %**, 25 teratas **82 %**.*
>
> ### Kalau memang tidak pernah bisa
>
> *Bila ekspor jenis ini **tidak pernah** menyertakan isi langkah, mohon katakan begitu — itu
> jawaban yang sah dan tetap berguna bagi kami. Yang kami perlukan berikutnya adalah cara lain
> yang tersedia: cetakan layar aturan satu per satu, dokumentasi perhitungan bila pernah dibuat,
> atau akses baca ke sistem lama walau terbatas waktu.*

---

## ⛔ Tidak ada yang perlu dikirim

**Nol aturan. Nol ekspor ulang. Nol cetakan layar. Nol akses sistem lama.**

Seluruh isi yang diminta surat ini sudah ada di tangan kami sejak awal. Yang gagal bukan ekspornya,
melainkan cara tim migrasi membacanya.

⚠️ **Yang masih kami minta kepada pihak lain, dan tidak berubah oleh pembatalan ini:**
**P1** kepada DBA — badan **empat** stored procedure *(`POOLDATA.PEGA_TREATY_IN`,
`POOLDATA.PEGA_JSON_POLIS_TREATYIN`, `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`,
`POOLDATA.PEGA_DELETE_ERROR_KONVERSI`)* ditambah contoh isi kolom JSON polis. Suratnya
*(`SURAT-P1-KE-DBA.md`)* **tetap berlaku dan tetap dikirim.**

---

*Dibatalkan 2026-09-22, sesudah `VERIFIKASI-P18.md` mengukur ulang korpus dan menemukan bahwa
penahan terberat proyek ini tidak pernah ada.*
