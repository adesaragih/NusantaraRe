> Modul  : Komite Claim Non Prop · Tahap 3 · 2026-09-21
> Peran  : juru catat
> Untuk  : pemilik proses komite klaim
> Status : MENUNGGU JAWABAN
> Sifat  : TAMBAH-SAJA

# Tiga keputusan yang kami perlukan

Pembangunan modul komite sudah berjalan dan **tidak berhenti menunggu memo ini**. Ketiga hal
di bawah punya jalan yang berjalan hari ini. Yang kami minta adalah persetujuan atau koreksi,
bukan izin untuk mulai.

Tiap butir menyebut apa yang berjalan bila Anda tidak menjawab, dan apa yang berubah bila
Anda menjawab lain. Ketiganya milik orang yang sama: yang memutuskan siapa menandatangani apa.

---

## 1 · Siapa yang berhak mengubah daftar anggota komite?

Daftar anggota komite — siapa menduduki jenjang apa, siapa menggantikan siapa saat cuti —
kini menjadi **data yang dapat disunting**, bukan lagi sesuatu yang tertanam di dalam
program. Itu perbaikan besar: mutasi jabatan tidak lagi memerlukan permintaan perubahan ke
tim pengembang.

Tetapi data yang dapat disunting perlu penjaga. Pertanyaannya: **jabatan apa yang berhak
menyuntingnya?**

**Yang berjalan hari ini.** Penjaganya sudah dibangun. Kami memakai satu jabatan
administratif tersendiri — seseorang yang mengurus daftar tetapi **bukan anggota komite**.
Pemisahan itu disengaja: orang yang memutuskan klaim sebaiknya bukan orang yang dapat
mengubah siapa berhak memutuskan.

**Bila Anda menjawab lain.** Yang berubah hanya **nama jabatan yang dicocokkan** sistem.
Penjaganya tetap ada, dan pekerjaannya tidak diulang.

**Yang tertahan.** Tidak ada. Pembangunan berjalan; namanya dapat dipasang belakangan.

**Rekomendasi kami:** satu jabatan administratif tersendiri di luar keanggotaan komite,
sebagaimana yang berjalan sekarang. Sebutkan jabatannya, dan kami pasang.<sup>1</sup>

---

## 2 · Apakah batas nilai rupiah dikembalikan?

Dahulu sistem dimaksudkan memilih tingkat komite berdasarkan **besarnya nilai** yang
diusulkan — usulan besar naik ke jenjang lebih tinggi. Pembanding nilainya **dimatikan** di
dalam program, dengan catatan "untuk sementara", dan sejak itu tidak pernah dinyalakan
kembali. Hari ini yang berlaku hanya dua tingkat, dan nilai rupiah tidak lagi ikut menentukan.

Kami **tidak menghidupkannya kembali sendiri**, karena kami tidak punya sumber tertulis untuk
angka ambangnya. Angka yang pernah kami temukan hanya tertulis di dalam komentar program, dan
komentar bukan kebijakan.

**Yang berjalan hari ini.** Dua tingkat, persis seperti perilaku yang berlaku sekarang.
Tidak ada perubahan yang dirasakan siapa pun.

**Bila Anda menjawab lain.** Aturan pemilihan tingkat kini berbentuk **tabel yang disunting**,
bukan program. Menghidupkan batas nilai berarti menambah baris ke tabel itu — **tidak perlu
menunggu rilis, tidak perlu tim pengembang.**

**Yang tertahan.** Tidak ada, sekarang maupun nanti.

**Rekomendasi kami:** biarkan dua tingkat untuk sementara, dan putuskan batas nilainya setelah
sistem berjalan beberapa bulan — saat itu Anda punya angka nyata untuk mendasarinya, bukan
tebakan.<sup>2</sup>

---

## 3 · Dua kiriman keluar yang tidak ada dalam rancangan baru

Sesudah komite memutuskan, sistem lama mengirim dua hal ke luar: satu memberi tahu **sistem
konversi klaim**, satu menuliskan **salinan data klaim** ke tempat lain. Rancangan baru
mendaftarkan lima kiriman keluar; **kedua ini tidak ada di antaranya.**

Kami tidak menambahkannya sendiri, karena menambah kiriman keluar berarti memutuskan bahwa
sistem lain masih memerlukannya — dan itu bukan pertanyaan teknis, melainkan pertanyaan
tentang siapa memakai apa.

Perlu ditegaskan: ini **bukan** perkara kekurangan bukti. Berkas kedua program itu ada pada
kami dan sudah kami baca. Yang belum ada adalah keputusan apakah keduanya masih dipakai.

**Yang berjalan hari ini.** Tidak ada yang dibangun untuk keduanya. Bila jawabannya "masih
dipakai", keduanya ditambahkan sebagai kiriman keluar seperti lima yang lain.

**Bila Anda menjawab lain.** Bila keduanya tidak lagi dipakai, kami catat sebagai dibuang,
dan tidak ada yang perlu dibangun.

**Yang tertahan.** Tidak ada pekerjaan tertahan, tetapi **pertanyaan ini perlu tertutup
sebelum sistem baru menggantikan yang lama** — bila ada sistem lain yang diam-diam menunggu
kiriman itu, ketiadaannya baru terasa saat peralihan.

**Rekomendasi kami:** tanyakan kepada pemilik sistem konversi klaim apakah mereka masih
menerima kiriman dari komite. Satu percakapan menutup keduanya.

---

## Ringkas

| # | Yang diputuskan | Menahan pekerjaan? | Bila diam |
|---|---|---|---|
| 1 | Jabatan yang berhak menyunting daftar anggota komite | tidak | jabatan administratif tersendiri di luar keanggotaan komite |
| 2 | Batas nilai rupiah penentu tingkat komite | tidak | dua tingkat, seperti yang berlaku sekarang |
| 3 | Dua kiriman keluar sesudah keputusan komite | tidak, tetapi perlu tertutup sebelum peralihan | tidak dibangun |

---

<sup>1</sup> Ketetapan `H-4`; gerbangnya dibangun pada tiket `01`.
<sup>2</sup> Ketetapan `K6-1`; aturan pemilihan tingkat sebagai data ditetapkan `K5-5`, dibangun pada tiket `01`.
<sup>3</sup> `KonversiKlaim_Act` dan `InsertJsonClaimTreatyNonProp_act`; berkasnya terdaftar pada `INVENTARIS-BUKTI.md` §2.4. Keadaannya `BELUM DIPUTUSKAN` pada bagian *Ketertelusuran* `SPEC-KOMITE-01.md`.
