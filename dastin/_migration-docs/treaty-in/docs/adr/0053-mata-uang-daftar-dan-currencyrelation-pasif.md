# ADR-0053 — Mata uang sebagai daftar; `CurrencyRelation` disimpan tetapi tidak menggerakkan angka

**Status:** diterima **tanpa verifikasi**, 23 September 2026
**Berlaku untuk:** Treaty In

## Konteks

Sistem lama menyimpan mata uang layer dalam **dua slot** (`Currency`/`Currency2`,
`Limit`/`Limit2`), dan sebuah atribut `CurrencyRelation` bernilai `AND` atau `OR` yang membuat
rumus limit bercabang. Tidak ada satu pun aturan, memo, maupun label yang menjelaskan apa arti
bisnis kedua nilai itu.

Uji I dan J, yang akan menunjukkan apakah atribut itu masih berarti dan apakah asumsi slot keliru,
**tidak dijalankan**.

## Keputusan

### a. Bentuk: daftar, bukan dua slot

Mata uang per layer **tidak dibatasi**. Bila batas dua kelak dikonfirmasi bisnis, ia menjadi
**aturan validasi**, bukan bentuk tabel.

### b. `CurrencyRelation`: disimpan, ditampilkan, **tidak dipakai menghitung**

Ia bertahan sebagai atribut bernilai bernama dan tetap terlihat pengguna, tetapi **tidak dipakai
dalam perhitungan apa pun** sampai bisnis menamai artinya.

> Atribut yang artinya tidak diketahui **tidak boleh menggerakkan angka.** Menyimpannya aman;
> membiarkannya bercabang di rumus adalah mewarisi perilaku yang tidak bisa dijelaskan siapa pun.

### c. Penyebut Rate on Line dihitung sekali

Setiap mata uang dikonversi sekali lalu dijumlahkan. Tidak ada percabangan atas `CurrencyRelation`,
dan tidak ada penjumlahan limit per mata uang EGNPI.

## Dasar

**Bagian a — penalaran.** Asimetri di cabang lama (satu sisi menambahkan limit tanpa konversi, sisi
lain dengan konversi) lebih mungkin merupakan jejak **asumsi bahwa slot pertama adalah IDR**
daripada semantik AND/OR. Bentuk daftar menampung kedua kemungkinan tanpa kehilangan apa pun, dan
tidak memaksa kita memilih sebelum tahu.

**Bagian b — penalaran, dan ini yang terpenting.** Membawa percabangan yang tidak bisa dijelaskan
berarti memindahkan perilaku yang tidak ada pemiliknya ke sistem baru, di mana ia akan tampak
disengaja.

**Bagian c — artefak.** Penjumlahan limit di sistem lama berada **di dalam** perulangan daftar
EGNPI, sehingga limit terjumlah sekali untuk setiap mata uang EGNPI; dan pada relasi `AND`, dua
langkah berurutan menambahkan limit yang sama dua kali dalam satu iterasi. Keduanya cacat, bukan
kebijakan.

## Syarat pembalikan

Bisnis menamai arti `AND` dan `OR`. Saat itu `CurrencyRelation` **boleh naik menjadi masukan
perhitungan**, dan rumus limit bercabang atasnya.

**Bila terbalik:** ini salah satu dari enam keputusan yang mengubah bentuk model — atribut pasif
menjadi masukan aktif.
