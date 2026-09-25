> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : juru catat
> Masukan: KETETAPAN.md K5-1 · GRILL-05/01-TEMUAN.md N-01, N-02 · GRILL-05/06-PUTUSAN.md T5-01, T5-02
> Status : TERBUKA
> Sifat  : HIDUP

# ADR-0030 · Satu penentu jenjang aktif

**Status:** Diterima · 2026-09-20 · sumber normatif `KETETAPAN.md` `K5-1`

## Konteks

Sistem lama memakai **dua** penentu jenjang aktif di dalam satu rule yang sama,
`KomiteRouter`:

- **Hitungan.** Langkah 1–4 memilih sasaran penugasan dari `.KomiteCount`, sebuah bilangan
  yang lahir bernilai 1 dan bertambah satu di `KomitePostAdjustment` langkah 29.
- **Keadaan baris.** Langkah 6 beriterasi atas daftar jenjang, berhenti pada baris pertama
  ber-`KomiteAproval==0` (transisi `6/6`, keluar iterasi), dan menugaskan ke pemegangnya.

Keduanya tidak pernah saling memeriksa, dan langkah 6 menang karena berjalan belakangan.

Hitungan itu juga dipakai untuk hal ketiga: penutupan case. `KomitePostAdjustment` langkah
20 menetapkan `KomiteCount = KomiteLoop` ketika hasilnya penolakan — bukan untuk mencatat
apa pun, melainkan agar pra-syarat langkah 28 (`KomiteCount >= KomiteLoop`) terpenuhi dan
`ASMForceCaseClose` terpanggil. Penutupan sirkulasi dengan demikian bergantung pada
perbandingan dua bilangan, dan salah satunya sengaja dipalsukan.

Bilangan pembandingnya, `KomiteLoop`, adalah jumlah baris **hasil laporan** roster, bukan
jumlah jenjang yang benar-benar terbentuk. Pada jalur akseptasi bersyarat daftar jenjang
tidak dibersihkan sebelum diisi, sehingga jumlah anggota tumbuh sementara `KomiteLoop`
tidak. Sejak itu kedua bilangan berbeda, dan tiga hal ikut salah sekaligus: kapan giliran
berhenti, siapa yang tidak pernah mendapat giliran, dan apakah case tertutup.

Keadaan keputusan sendiri bersifat sesaat: `KomiteRouter` langkah 5 mengosongkan
`.AcceptStatus` tanpa pra-syarat, sedangkan gerbang lingkar menuntutnya bernilai `"1"`.

## Keputusan

Jenjang aktif ditentukan oleh **satu hal saja**: jenjang berderajat terendah yang belum
memutuskan.

Hitungan jenjang adalah **turunan** yang dihitung dari daftar jenjang. Ia tidak pernah
menjadi syarat penugasan, tidak pernah menjadi syarat penutupan, dan tidak pernah disimpan
sebagai penyimpan kedua yang dapat berselisih dengan daftar jenjangnya.

Sirkulasi ditutup dari **keadaan keputusan** — tidak ada jenjang tersisa yang belum
memutuskan, atau sebuah keputusan mengakhiri sirkulasi menurut jenisnya — bukan dari
perbandingan dua bilangan.

Keputusan tiap jenjang adalah **catatan tetap** pada jenjang itu, bukan satu medan yang
dipakai bergantian oleh seluruh sirkulasi.

## Konsekuensi

- Urutan giliran mengikuti derajat naik, dan itu **paritas** — sistem lama sudah berperilaku
  demikian. Uji `P5-01` dan `P5-02` dirancang lulus.
- Keadaan "jenjang aktif menurut hitungan berbeda dari jenjang aktif menurut daftar" menjadi
  **tidak dapat direpresentasikan**. Uji `P5-02b` karena itu dirancang gagal, dan itu
  dinyatakan, bukan disembunyikan.
- Empat nama keranjang yang ditulis di dalam rule hilang bersama penentu keduanya; sasaran
  penugasan menjadi data per jenjang.
- Riwayat keputusan menjadi dapat dibaca ulang setelah sirkulasi selesai — kemampuan yang
  sistem lama tidak punya.
- Migrasi data harus menghitung ulang jenjang aktif dari daftar jenjang, dan **tidak boleh**
  mempercayai `KomiteCount` yang tersimpan.

## Alternatif yang ditolak

**Membawa kedua penentu dan menambahkan pemeriksaan konsistensi.** Ditolak: pemeriksaan
konsistensi hanya memberi tahu bahwa keduanya sudah berselisih, tanpa memberi tahu yang mana
yang benar. Sistem lama sudah menjalankan percobaan itu selama bertahun-tahun.

**Membawa hitungan sebagai penyimpan dan menurunkan daftar jenjang darinya.** Ditolak: arah
turunannya terbalik terhadap sumber kebenarannya. Siapa yang berhak memutus adalah fakta
tentang orang dan derajat, bukan tentang bilangan.

**Menutup sirkulasi dengan penanda "sudah selesai" yang ditulis pemutus terakhir.** Ditolak:
ia mengulang pola yang sama — satu fakta disimpan di dua tempat — dan menyerahkan kebenaran
penutupan kepada langkah yang dapat gagal.
