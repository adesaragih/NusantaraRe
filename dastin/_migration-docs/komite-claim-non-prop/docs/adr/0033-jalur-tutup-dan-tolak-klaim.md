> Modul  : Komite Claim Non Prop · Tahap 3 · 2026-09-21
> Peran  : juru catat
> Masukan: KETETAPAN.md K6-4 · KETETAPAN.md §8 F-20 · REGISTER-DEVIASI.md bagian 6 (deviasi 23, 24, 25) · SPEC-KOMITE-01.md S-001 dan S-008
> Status : TERBUKA
> Sifat  : HIDUP

# ADR-0033 · Jalur menutup dan menolak klaim

**Status:** Diterima · 2026-09-21 · sumber normatif `KETETAPAN.md` `K6-4`

## Konteks

Modul komite punya dua jalur yang mekanismenya identik sampai ke tingkat langkah tetapi
akibatnya berbeda: mengedarkan **usulan pembayaran**, dan mengedarkan **usulan menutup atau
menolak klaim**. Jalur kedua lahir dari rule tersendiri, `CreateChildKomiteCloseNP_Act`, dan
cara ia menentukan siapa memutus tidak punya kemiripan dengan jalur pertama.

Jalur pembayaran menyusun jenjangnya dari daftar anggota di basis data, disaring oleh ambang
kewenangan. Jalur menutup dan menolak tidak melakukan keduanya. Ia **tidak membaca daftar
anggota sama sekali** dan **tidak mengenal ambang**. Sebagai gantinya, langkah 8 dan 9 menulis
pemegangnya langsung di dalam rule: **nama orang, alamat surel, inisial, dan dua sebutan
jabatan yang berbeda untuk satu kedudukan yang sama** (`F-20`). Empat bentuk pengenal untuk
satu orang, seluruhnya tetapan di dalam kode.

Akibatnya berlipat. Ketika orang itu pindah jabatan, cuti, atau berhenti, tidak ada tempat
untuk mengubahnya selain kode — dan karena satu kedudukan ditulis dengan dua sebutan jabatan
yang berbeda, mengubahnya di satu tempat meninggalkan yang lain. Keadaan "tidak ada pemegang"
juga **tidak dapat terjadi** di jalur ini: pemegangnya tetapan, sehingga sirkulasi selalu
lahir, bahkan ketika orangnya sudah tidak berwenang.

Jenis sirkulasinya pun tidak dinyatakan. Langkah 12 membedakan menutup dari menolak dengan
membaca **nilai sebuah kolom analisis** — kolom yang diisi untuk keperluan lain, dan yang
dapat berselisih dengan tindakan yang benar-benar dipilih pengaju tanpa ada satu pun jalur
yang memeriksanya. Dua maksud berbeda karena itu dapat tertukar diam-diam.

Ketetapan ronde 4 yang mengatur jalur ini, `H-2`, **bunyinya hilang** bersama berkas yang
memuatnya. Yang tersisa hanyalah catatan bahwa keputusan itu pernah diambil. Itu sendiri
alasan untuk menuliskannya sebagai ADR: keputusan ini sudah sekali kehilangan bunyinya.

## Keputusan

Sirkulasi menutup klaim dan menolak klaim **tetap satu jenjang**. Kedalamannya tidak
disamakan dengan jalur pembayaran.

Pemegang jenjang itu **diselesaikan dari daftar anggota berdasarkan peran**, bukan dari nama
yang ditulis di dalam aturan. Peran adalah atribut baris daftar anggota; nama orang tidak
muncul di jalur mana pun.

Bila **tidak ada baris daftar anggota yang aktif** memegang peran itu, sirkulasi **ditolak
saat dibuat**, dengan galat yang menyebut peran yang kosong.

**Jenis sirkulasi dinyatakan oleh pemanggil.** Ia tidak diturunkan dari nilai kolom analisis
mana pun.

## Konsekuensi

- Mutasi jabatan menjadi penyuntingan data. Tidak ada permintaan perubahan kode ketika orang
  yang menandatangani penutupan klaim berganti.
- Satu kedudukan punya **satu** sebutan. Dua sebutan berbeda untuk satu kedudukan tidak dapat
  lahir kembali, karena sebutannya kini kolom, bukan dua tetapan di dua langkah.
- Keadaan "tidak ada pemegang" menjadi **dapat terjadi dan terlihat**. Sistem lama tidak dapat
  menghasilkannya; sistem baru menolak pembentukan alih-alih melahirkan sirkulasi yang tidak
  ada yang berwenang memutuskannya.
- Dua maksud berbeda **tidak dapat tertukar**. Menutup dan menolak dibedakan oleh tindakan
  yang dipilih pengaju, bukan oleh kolom yang diisi untuk keperluan lain.
- Tiga perbedaan perilaku terhadap sistem lama, dan ketiganya **disengaja**: deviasi 23
  (jenis dinyatakan pemanggil), 24 (pemegang dari peran), dan 25 (penolakan saat peran
  kosong). Masing-masing membawa uji yang **dirancang gagal** terhadap sistem lama —
  `PS-01`, `PS-02`, `PS-03`. Uji yang lulus paritas adalah tanda deviasinya tidak terpasang.
- Jalur ini karena itu **tidak dapat dibangun sebelum daftar anggota berdiri sebagai data**.
  Urutan itu mengikat.
- `H-2` digantikan. Ia tetap tercatat `BUNYI HILANG`; ADR ini tidak memulihkannya, melainkan
  menggantikannya dengan bunyi yang ditulis dari langkah yang terbaca.

## Alternatif yang ditolak

**Menyamakan kedalamannya dengan jalur pembayaran — banyak jenjang, disaring ambang.**
Ditolak: menutup klaim adalah keputusan yang **membatalkan**, bukan yang mengeluarkan uang.
Kedalaman berjenjang ada untuk menjaga pengeluaran, dan memakainya di tempat yang tidak
mengeluarkan apa pun menambah persetujuan tanpa menambah penjagaan. Sistem lama pun
memperlakukannya satu jenjang; yang kami ubah adalah **siapa** jenjang itu, bukan berapa.

**Menurunkan jenis sirkulasi dari kolom analisis, sebagaimana sistem lama — paritas.**
Ditolak: paritas terhadap perilaku yang dapat berselisih dengan maksud pengaju adalah paritas
terhadap cacat. Kolom itu diisi untuk keperluan lain, dan tidak ada satu pun jalur yang
memeriksa apakah isinya sejalan dengan tindakan yang dipilih. Mempertahankannya berarti
mewarisi kemungkinan tertukar tanpa mewarisi satu pun manfaat.

**Memindahkan nama orang dari kode ke berkas konfigurasi.** Ditolak: ia memindahkan masalahnya
tanpa menyelesaikannya. Yang salah bukan **tempat** nama itu ditulis, melainkan bahwa
kewenangan melekat pada **orang** alih-alih pada **kedudukan**. Konfigurasi berisi nama orang
tetap memerlukan penyuntingan setiap mutasi, tetap dapat berselisih dengan daftar anggota
yang sebenarnya, dan tetap melanggar larangan bercabang berdasarkan identitas orang.

**Membiarkan sirkulasi lahir tanpa pemegang, lalu menugaskannya belakangan.** Ditolak: ia
mengulang persis cacat yang sudah ada di jalur pembayaran — berkas yang lahir lalu diam,
tanpa ada yang tahu ia menunggu. Sirkulasi tanpa jenjang ditolak saat dibuat, dan jalur ini
tidak dikecualikan.
