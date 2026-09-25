# ADR-0047 — Data uji: identitas boleh disamarkan, angka tidak

**Status:** diterima, 23 September 2026
**Berlaku untuk:** migrasi dan pengujian Treaty In

## Keputusan pemisah

> **Identitas boleh disamarkan. Angka tidak boleh disentuh sama sekali.**

| Boleh disamarkan | Tidak boleh diubah walau satu digit |
|---|---|
| nama cedant, leading reinsurer, broker, retrosesioner | seluruh nilai uang |
| nama operator pada catatan komentar | seluruh persentase |
| nama berkas lampiran | seluruh tanggal |
| | mata uang, jumlah lines, jumlah reinstatement, struktur layer |

Alasannya dari ADR-0043: kriteria migrasi adalah **identik**, nol toleransi, dan uji paritas
membandingkan hasil kalkulator baru terhadap angka lama. Penyamaran yang menyentuh angka
menghancurkan kedua uji itu sekaligus, dan tidak ada cara mengetahuinya sampai terlambat.

**Teks bebas pada catatan komentar dibuang seluruhnya, bukan disamarkan.** Ia berisiko tertinggi —
orang menulis apa saja di kolom komentar — dan paling tidak diperlukan, karena tidak satu pun uji
paritas membacanya.

## Kejujuran yang harus dicatat

**Menyamarkan nama tidak menganonimkan data treaty.** Sebuah treaty dengan limit tertentu, periode
tertentu, susunan layer tertentu, dan leading reinsurer tertentu dapat dikenali kembali oleh siapa
pun yang bekerja di pasar reasuransi Indonesia. Data ini kecil dan pesertanya saling kenal.

Jangan membangun rasa aman dari transformasi. Perlindungan yang nyata ada pada **kendali akses dan
lingkungan**:

- data uji tinggal di lingkungan terkendali dengan akses bernama, tidak disalin ke laptop;
- siapa yang mengaksesnya tercatat;
- salinan penuh dipakai untuk run paritas lalu dibersihkan, bukan dibiarkan menganggur.

## Dua kumpulan, bukan satu

| | Kumpulan kurasi | Salinan paritas |
|---|---|---|
| Ukuran | kecil, satu kontrak per kelas cacat | lengkap |
| Angka | boleh disamarkan berat, sebagian boleh dibuat | **asli**, tidak disentuh |
| Teks bebas | boleh dibuat | dibuang |
| Dipakai untuk | membangun dan menguji sehari-hari | run paritas saja |
| Letaknya | boleh di laptop orang | lingkungan terkendali |

Dengan pemisahan itu, pekerjaan harian tidak pernah menyentuh salinan penuh, dan salinan penuh
tidak pernah kehilangan presisinya.

## Kumpulan kurasi adalah KELUARAN uji, bukan disusun tangan

Setiap uji A sampai V **sudah menghasilkan daftar ID kontrak** yang memenuhi kriterianya — itu
justru inti kerjanya. Kumpulan kurasi adalah keluaran uji-uji itu: ambil satu atau dua ID dari
setiap daftar.

Kelas yang harus terwakili: layer bermata uang ganda; kontrak berbagian per baris; addendum
berprorata; addendum yang mengubah bagian NuRe; penyebaran manual; penyebaran otomatis yang
masternya sudah berubah; kontrak bertanda kegagalan Rate on Line; kontrak FAC-OUT dengan sisa
daftar penyebaran; kontrak berpenanda konversi generasi lama; addendum ber-`EDMSTATE` 3; baris
reinstatement yang persennya bukan 100; dan kontrak dengan rasio potongan terhadap bruto yang
menyimpang.

**Konsekuensi jadwal, direvisi 23 September 2026.** Kumpulan kurasi tetap merupakan **keluaran**,
bukan susunan tangan — tetapi bukan keluaran uji A–V, melainkan keluaran **run paritas**. Run
paritas sendiri **tidak menunggu apa pun**: ia berjalan atas seluruh data yang dimigrasi, dan
justru run itulah yang **menamai** kontrak mana yang masuk kumpulan kurasi.

Urutannya: **migrasi → run paritas → kumpulan kurasi.** Bukan sebaliknya, dan tidak ada yang
tertahan di dalamnya.

Yang tetap harus ditulis **sebelum** run dijalankan adalah **bentuk** selisih yang diharapkan per
kelas cacat — arah, rumus penyimpangan, tanda pengenal. Hanya **populasinya** yang datang dari run.
Lihat `KEPUTUSAN-TANPA-VERIFIKASI.md` §5.
