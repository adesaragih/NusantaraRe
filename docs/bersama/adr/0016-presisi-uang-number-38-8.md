---
status: accepted
tanggal: 2026-09-23
sumber: ronde keputusan dua belas butir — `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, keputusan work owner
memperluas: ADR-0003 (representasi uang non-float)
---

# Kolom uang memakai `NUMBER(38,8)`

Seluruh kolom uang di sistem baru memakai **`NUMBER(38,8)`** — tiga puluh digit di depan koma,
delapan di belakang. Berlaku **seluruh sistem**, bukan satu modul.

ADR-0003 menetapkan bahwa uang **tidak** direpresentasikan sebagai `float`. Ia tidak menetapkan
presisinya. Catatan ini menutup lubang itu.

## Kenapa delapan desimal di belakang koma

`[terverifikasi]` Nilai berdesimal terpanjang yang ditemukan pada dokumen produksi membawa **24
angka** di belakang koma. Seluruh selisih yang bermakna jatuh pada desimal **ketujuh atau
sebelumnya**; sisanya ekor galat penjumlahan.

Nilai berdesimal lebih dari delapan **dibulatkan saat dimuat, bukan ditolak**.

## Kenapa tiga puluh digit di depan koma — bukan dua belas

Angka semula yang diusulkan adalah `NUMBER(20,8)`, yaitu dua belas digit di depan koma. Sapuan
korpus membuktikan itu **tidak cukup**:

| Nilai di korpus | Digit di depan koma | Sifat |
| --- | ---: | --- |
| `99999999999999999.99` | **17** | sentinel "tanpa batas" |
| `181500000000.00` | **12** | ambang kapasitas treaty |
| `150000000000.00` | **12** | ambang kapasitas treaty |
| `200000000.00` | 9 | batas premi neto |

⭐ Ambang dagang yang nyata sudah **tepat di batas** dua belas digit — nol ruang lega. Penjumlahan
lintas mata uang dapat melewatinya, dan menaikkan presisi **sesudah** data terisi jauh lebih mahal
daripada menaikkannya sekarang.

## Kenapa pelebaran ini tidak berbiaya

`NUMBER` di Oracle berpanjang **berubah-ubah**: hanya digit bermakna yang tersimpan. Kolom berisi
`1500` memakan ruang yang sama entah dideklarasi `(20,8)` atau `(38,8)`. Deklarasi menetapkan
**batas**, bukan ukuran tersimpan.

⭐ Karena itu tidak ada alasan menahan presisi di angka yang pas-pasan. **38 adalah presisi
tertinggi yang diterima Oracle.**

## Akibat

1. Seluruh kolom uang pada modul mana pun dideklarasi `NUMBER(38,8)`. Tidak ada modul yang
   menetapkan presisi sendiri.
2. Kolom **persen** tidak termasuk — ia bukan uang, dan tetap mengikuti ketetapan modulnya.
3. Lapisan `services` membulatkan pada desimal kedelapan **saat memuat**, bukan tiap kali membaca.
4. Test yang menemukan kolom uang bertipe `float`, atau berpresisi selain `(38,8)`, **gagal**.

## Yang catatan ini TIDAK putuskan

⛔ Bukan tentang pembulatan ke presisi mata uang untuk **tampilan**. Itu urusan penyajian.

⛔ Bukan tentang satuan minor. Sistem ini memakai desimal berskala tetap, bukan bilangan bulat
dalam satuan minor.
