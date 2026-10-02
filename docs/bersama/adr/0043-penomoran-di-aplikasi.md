---
status: accepted
tanggal: 2026-09-27
sumber: keputusan work owner (butir o1, o2, o3), brief lanjutan 4 §1
supersedes: ADR-U-0006
---

# Penomoran dijalankan aplikasi, bukan stored procedure

`[keputusan work owner]` *"jangan ada lagi pemanggilan procedure, segala procedure hardcode dalam
skrip"*. Seluruh logika penomoran ditulis ulang di Go; nol `CALL` ke procedure mana pun.

Ini **meng-supersede ADR-U-0006** pada bagian yang menyerahkan penerbitan nomor bisnis kepada
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`. Bagian ADR-U-0006 yang menetapkan bahwa **identitas
tabel** `T_CLAIMLF_*` berasal dari **sequence** tetap berlaku — keduanya hal berbeda: yang satu
nomor bisnis yang dibaca manusia, yang lain pengenal baris.

## Yang tetap di basis data, dan alasannya

| Tetap di Oracle | Alasan |
| --- | --- |
| `SELECT … FOR UPDATE` atas `(CLASS, JENIS, TAHUN)` | kunci baris memang milik basis data. Menirunya di aplikasi berarti menulis ulang penguncian, dan dua pemanggil serentak akan membaca urut yang sama lalu menerbitkan nomor kembar |
| `sequence.NEXTVAL` untuk nomor akseptasi | urut yang aman dari dua pemanggil serentak adalah persis yang sequence jamin |

## Yang pindah ke aplikasi

| Pindah ke Go | Bentuk |
| --- | --- |
| penentuan periode | hari tutup buku dari `TANGGAL_CLOSING`; bila hari berjalan melewatinya, periode digeser satu bulan — **beserta tahunnya** |
| cabang cutover | `TRUNC(now) <= 02/01/2026` → periode `12.2025`, tahun `2025`. Sudah lewat hari ini, **dipertahankan** sebab pengurai nomor lama (tiket 13) memerlukannya |
| penaikan dan penyimpanan penghitung | baca terkunci, naikkan, perbarui; bila baris belum ada, sisipkan dengan urut `1` |
| perakitan nomor | `RNML-K<kode bisnis>.<MM.YYYY>.<5 digit>` untuk klaim; `RNML-A<kode>.<MM>.<YY>.<5 digit>` dan `RNML-AR…` untuk akseptasi |

## Akibat 1 — dua penomoran, dua perilaku, keduanya ditiru apa adanya

Nomor **klaim** menggeser tahun ketika bulan berguling Desember → Januari *(`ADD_MONTHS` memang
begitu)*. Nomor **akseptasi** di Pega **tidak**, sebab kedua cabang `@if`-nya identik
*(`SaveAdjustment_Act` baris 605)* — dan di sana kami **menyimpang sadar** dengan menggeser
tahunnya, sebab nomor Januari bertahun Desember dapat bertabrakan.

Perbedaan itu bukan kelalaian: ia dua rule yang berbeda, dan masing-masing dibaca sendiri.

## Akibat 2 — tulisan ke tabel warisan yang disengaja

`GENERATE_SEQUENCE_NUMBER` di-`UPDATE`/`INSERT` oleh aplikasi. Itu inheren pada keputusan ini:
penghitung yang tidak disimpan bukan penghitung. Dicatat di sini supaya ia tidak terbaca sebagai
kebocoran.

## Akibat 3 — label sumber tidak dinaikkan

Bentuk procedure-nya dibaca dari katalog instance **pengembangan** dan berlabel
`[data DBA — belum dikonfirmasi DBA]` di sumbernya. Kode yang menirunya memakai label **yang sama**;
menaikkannya menjadi `[terverifikasi]` akan mengaku lebih banyak daripada yang diketahui.

## Yang tetap terbuka

`[terbuka — DBA]` **OQ-013**: procedure aslinya tidak memuat `COMMIT`, dan `COMMIT` yang terlihat di
rule Pega berada di blok pemanggil. Batas transaksi di sisi basis data belum ditetapkan. Jalur kita
sendiri tidak pernah menulis `COMMIT` (ADR-U-0029).
