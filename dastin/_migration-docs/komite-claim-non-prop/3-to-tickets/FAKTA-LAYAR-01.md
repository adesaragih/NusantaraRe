> Modul  : Komite Claim Non Prop · Tahap 3 · 2026-09-21
> Peran  : pembaca berkas
> Masukan: `Komite Claim Non Prop/Section/ShowTransfer.xml` (1.939.867 byte, ekspor 2026-09-09)
> Status : DITUTUP 2026-09-21
> Sifat  : TAMBAH-SAJA

# FAKTA-LAYAR-01 — enam kendali ber-`pyDisabledWhen`

Berkas ini **bukan** berkas normatif dan tidak menetapkan apa pun. Ia memuat satu pembacaan
yang selama ini tertunda, agar tiket Layar dipecah dengan cacah yang benar.

Sebab ia ada: `F-2` mencacah enam medan dapat-sunting dan menyatakan empat di antaranya
hanya untuk jenjang pertama, tetapi **tidak menamai** keempatnya. `PENGETAHUAN.md` §10.1
mendaftarkan medan rekening sebagai dapat-sunting, dan `F-1` serta `F-6` membatalkan
pendaftaran itu. Dua nama karena itu hanya hidup sebagai angka.

## Cara membacanya — dan mengapa ini bukan kedekatan posisi

`F-12` menolak pembacaan sebelumnya dengan alasan yang benar: **kedekatan posisi bukan
bukti**. Pembacaan ini tidak memakai kedekatan. Tiap `pyDisabledWhen` berada di dalam blok
`<pyModes><rowdata>` milik satu `Embed-Display-Table-Cell`, dan sel itu punya `<pyValue>`
yang menamai properti terikatnya. Yang dibaca adalah **isi sel**, bukan apa yang kebetulan
berada di dekatnya.

Cacah: `pyDisabledWhen` muncul **enam kali** di seluruh berkas. Keenamnya ada di bawah ini.
Tidak ada yang ketujuh.

## Keenam kendali

| `pyValue` | `pyFormat` | `pyReadOnly` | `pyRequiredNew` | `pyDisabledWhen` |
|---|---|---|---|---|
| `.IsSubjectivity` | `pxCheckbox` | `false` | — | `.KomiteCount!='1'` |
| `.SubjectivityNote` | **`pxDropdown`** | `false` | `true` | `.KomiteCount != 1` |
| `.Adjustment.IsProposeClose` | `pxCheckbox` | `false` | — | `.KomiteCount!='1'` |
| `.Adjustment.IsPropReserved` | `pxCheckbox` | `false` | — | `.KomiteCount!='1'` |
| `pyWorkCover.ClaimData.ReporterStatus` | `pxDropdown` | **`true`** | `false` | `pyWorkCover.IsOutstanding==1` |
| `.TotalClaim` | `pxNumber` | **`true`** | `false` | `.CNPFlagOuts==1` |

## Tiga fakta yang keluar dari tabel itu

**1 · Keempat medan berkunci jenjang pertama bernama.** `.IsSubjectivity`,
`.SubjectivityNote`, `.Adjustment.IsProposeClose`, `.Adjustment.IsPropReserved`. Keempatnya
memakai pembanding `.KomiteCount` terhadap satu — dua ejaan berbeda untuk perbandingan yang
sama, yang `K5-1` gabungkan menjadi satu turunan. Daftar tujuh medan pada `SPEC-KOMITE-01.md`
bagian *Medan layar* **cocok**; yang berubah hanyalah bahwa keempatnya kini `EVIDENCED`,
bukan diterima atas wewenang instruksi.

**2 · Dua `pyDisabledWhen` sisanya menunjuk kendali yang sudah `pyReadOnly=true`.**
`pyWorkCover.ClaimData.ReporterStatus` dan `.TotalClaim` keduanya baca-saja, sehingga syarat
menonaktifkannya **tidak pernah berpengaruh** — kendali yang tidak dapat disunting tidak
dapat dinonaktifkan lebih jauh. Ini menutup `F-12` sebagai fakta: sasaran `IsOutstanding` dan
`CNPFlagOuts` **terbaca**, dan keduanya mati.

Akibatnya pada spesifikasi: kesimpulan bagian *Medan layar* — kedua kendali itu **dibuang** —
**tetap berlaku**, dan alasannya menguat. Bunyi lamanya bersandar pada "sasarannya tidak
dapat ditentukan"; dasar yang sekarang adalah "sasarannya terbaca, dan keduanya baca-saja".
Kesimpulan tidak berubah, sehingga tidak ada ketetapan yang perlu dibuka.

**3 · `.SubjectivityNote` adalah `pxDropdown`, bukan kotak teks.** `SPEC-KOMITE-01.md`
menuliskannya sebagai "teks panjang", dan kolom `CATATAN_BERSYARAT` berbentuk
`VARCHAR2(2000 CHAR)`. Bentuk kendalinya di sistem lama adalah **daftar pilihan**. Ini
selisih yang nyata dan **belum diselesaikan**; ia masuk daftar hal terbuka pada tiket Layar,
bukan diputuskan di sini. Isi daftar pilihannya tidak terbaca dari berkas ini — `Rule-Obj-FieldValue`
tidak ikut diekspor (`KETETAPAN.md` §10.1, tambahan ronde 5).

## Apa yang tidak dibaca

Berkas ini hanya menjawab pertanyaan `pyDisabledWhen`. Cacah 296 kendali, aturan tampil per
blok, dan daftar medan baca-saja **tidak** dibaca ulang — ketiganya sudah `EVIDENCED` lewat
`F-1`, `F-3`, `F-6`, `F-7`, `F-9`, dan tidak ada yang meragukannya.
