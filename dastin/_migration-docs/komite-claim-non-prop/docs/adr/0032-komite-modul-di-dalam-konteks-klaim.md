> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : juru catat
> Masukan: KETETAPAN.md D-2, D-3 · GRILL-05/06-PUTUSAN.md T5-04 · ADR-0016 (sisi Claim)
> Status : TERBUKA
> Sifat  : HIDUP

# ADR-0032 · Komite adalah modul di dalam konteks Klaim, bukan konteks berjaringan

**Status:** Diterima · 2026-09-20 · sumber normatif `KETETAPAN.md` `D-2`

## Konteks

Di sistem lama, komite adalah **case anak** dari klaim. Keterhubungan keduanya bukan
panggilan antar-sistem melainkan penyuntingan halaman bersama: activity keputusan membuka
klaim induk dengan kunci (`Obj-Open-By-Handle` atas `pxCoverInsKey`), memutakhirkan
`TempMainWork.ClaimData.AdjustmentList(...)`, lalu menyimpan dan melakukan commit.

Pola itu tidak berdiri di satu berkas. Kedua activity pembuat sirkulasi menjalankan
`Obj-Refresh-And-Lock` atas klaim induk **dua kali** — sekali di awal dan sekali menjelang
akhir — dengan penyimpanan di antaranya. Batas antara "apa yang milik komite" dan "apa yang
milik klaim" tidak pernah digambar; yang ada hanyalah halaman yang keduanya sunting.

Sementara itu langkah commit yang sama menjalankan enam efek ke luar berturut-turut tanpa
kompensasi, dan `ADR-0016` sisi Claim sudah melarang tautan basis data lintas sistem.

Menjadikan komite konteks berjaringan akan berarti menggambar batas di tempat yang sistem
lama tidak pernah gambar, lalu menanggung konsistensi akhir pada perkara yang selama ini
bertransaksi bersama.

## Keputusan

Komite adalah **modul di dalam konteks Klaim**, dengan agregat sendiri.

Sirkulasi, jenjang, dan keputusan membentuk satu agregat yang berubah bersama atau tidak
sama sekali. Perubahan yang menyeberang ke klaim — penanda maksud, penanda akibat, keadaan
klaim — terjadi di dalam batas transaksi yang sama, bukan lewat pesan yang mungkin sampai.

Batas transaksi itu menutupi **kedua** activity pembuat sirkulasi dan activity keputusan,
bukan hanya yang terakhir.

Efek ke luar yang tidak dapat dibatalkan (keputusan beku no. 6) berada di luar batas
transaksi. Jenis dan penanganannya per aliran ditetapkan setelah `PG-04`, `PG-05`, dan
`PG-06` dibuka.

## Konsekuensi

- Tidak ada tautan basis data, tidak ada panggilan sinkron antar-konteks, dan tidak ada
  konsistensi akhir antara sirkulasi dan klaim. `ADR-0016` tetap berlaku penuh.
- Skema `KLAIMNP` menampung objek komite; tidak ada skema terpisah dan tidak ada basis data
  terpisah.
- Penyuntingan halaman bersama digantikan oleh kepemilikan yang tegas: sirkulasi memiliki
  jenjang dan keputusan; klaim memiliki usulan dan penandanya. Tidak ada satu pun medan
  yang dimiliki keduanya.
- Pola "buka dengan kunci, sunting, buka dengan kunci lagi" tidak dibawa. Perkara `K-11`
  tidak dibuka kembali; yang berubah adalah bahwa batas transaksinya kini diputuskan,
  bukan diwarisi.
- Enam efek luar keluar dari langkah commit — itu memang sudah diperintahkan keputusan beku
  no. 6; ADR ini menyatakan bahwa mereka berada di luar batas transaksi, bukan ke mana
  mereka pindah.

## Alternatif yang ditolak

**Komite sebagai konteks berjaringan dengan pesan antar-konteks.** Ditolak: batasnya harus
digambar di tempat yang sistem lama tidak pernah gambar, dan perkara yang selama ini
bertransaksi bersama akan menanggung konsistensi akhir tanpa satu pun kebutuhan yang
memintanya.

**Komite sebagai bagian dari agregat klaim, tanpa agregat sendiri.** Ditolak: sebuah klaim
dapat memiliki banyak sirkulasi yang hidup bersamaan, dan menjadikan seluruhnya satu agregat
membuat dua sirkulasi atas dua usulan berbeda saling mengunci tanpa sebab.

**Menunda keputusan sampai beban nyata terukur.** Ditolak: batas agregat menentukan bentuk
tabel dan bentuk transaksi, dan keduanya ditulis sebelum beban apa pun terukur. Menunda
berarti memilih tanpa menyatakannya.
