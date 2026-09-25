# ADR-0036 — Beku saat disetujui, hitung saat dibaca

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In, dan diusulkan untuk seluruh modul

## Konteks

Bentuk masalah yang sama muncul tiga kali dalam satu sesi:

1. **Penyebaran** — baris penyebaran dihapus dan disusun ulang dari tabel susunan baku pada setiap
   perhitungan. Tabel itu disunting orang. Kontrak yang sudah melewati persetujuan empat tingkat
   berubah pembagian kapasitasnya begitu ia disentuh lagi.
2. **Angka yang dilihat approver** — tidak ada tempat yang mencatat berapa nilai kontrak pada saat
   disetujui. Yang tersimpan hanya hasil terakhir.
3. **Kurs** — konversi memakai kurs tahun berjalan saat perhitungan dijalankan, bukan kurs yang
   berlaku pada peristiwanya. Nilai rupiah kontrak lama bergeser setiap kali dihitung ulang di
   tahun berbeda.

Ditambah satu lagi yang ditemukan belakangan: **bagian NuRe** dipakai sebagai faktor konversi
tingkat, dan bagian yang dipakai adalah bagian hari ini, bukan bagian yang berlaku saat jumlah itu
terjadi.

## Keputusan

> Angka yang menjadi dasar suatu persetujuan **dibekukan pada saat persetujuan itu**, bersama
> seluruh masukan yang membentuknya.
>
> Angka yang menggambarkan **posisi terkini** dihitung saat dibaca dan tidak pernah disimpan
> sebagai atribut benda yang dilaporkannya.

Penerapan yang mengikat:

- Penyebaran adalah **turunan** selama kontrak belum disetujui, dan menjadi **fakta tercatat**
  begitu disetujui. Sesudah akseptasi, perubahan pada susunan baku **tidak menjalar**.
- Menyimpan **tidak pernah** menarik ulang dari data acuan. Tidak ada mode yang berpindah sendiri.
- Hasil yang dibekukan membawa **penunjuk** ke masukan yang dipakai: susunan mana, periode mana,
  kurs berapa dari sumber apa, bagian berapa, dan kapan dihitung — serta **asal-usulnya**
  (disemai dari acuan baku, diisi tangan, atau disemai lalu disunting).
- Bagian NuRe **boleh berubah** lewat addendum, tetapi perubahannya **berlaku ke depan**; ia tidak
  menulis ulang tingkat jumlah yang sudah terjadi.
- Keterlambatan jatuh tempo dihitung saat dibaca. **Penanda terlambat tidak pernah disimpan.**

## Kenapa ini bukan sekadar kebersihan model

Pemilik proses menetapkan bahwa Treaty In Adjustment akan **mengambil data lama lewat SELECT**,
bukan menyalinnya. Itu hanya aman bila data lama tidak bergerak. Kalau nilai kontrak masih bisa
berubah sesudah akseptasi, selisih yang dihitung hari ini dan bulan depan atas kontrak yang sama
akan berbeda tanpa ada yang mengubah apa pun.

Keputusan ini karena itu adalah **prasyarat teknis** bagi cara kerja Adjustment yang sudah
ditetapkan, bukan pilihan rancangan yang bisa ditunda.

## Konsekuensi

- Setiap versi kontrak yang disetujui harus punya identitas yang stabil dan bisa ditunjuk.
- Jejak perubahan (ADR-0045) menjadi bukti bahwa pembekuan benar-benar terjadi. Tanpa jejak,
  pembekuan adalah janji, bukan fakta.
- Sebagian besar angka pelaporan tidak punya kolom sama sekali.
