---
status: accepted
tanggal: 2026-09-23
sumber: spec Claim Prop, keputusan work owner
---

# Skema relasional dirancang baru bila ia tidak pernah ada

`[terverifikasi]` Pada sebagian modul, skema relasional **tidak dapat direkayasa-balik dari
korpus - sebab ia belum pernah ada**. Persistensi terjadi lewat penyimpanan objek kerja sebagai
dokumen, bukan lewat tabel.

`[keputusan work owner]` Untuk modul semacam itu, skema **dirancang baru**, dan perancangannya
dinyatakan sebagai pekerjaan tersendiri - bukan disamarkan sebagai hasil penelusuran.

## Bahan perancangan yang sah

| Sumber | Yang diambil darinya |
| --- | --- |
| Bentuk dokumen nyata `[data DBA]` | nama dan tipe medan |
| Hierarki daftar bersarang pada layar | bentuk induk-anak |
| Sensus kolom tabel proyeksi yang sudah ada | medan yang terbukti dipakai hilir |

## Kenapa ini perlu dinyatakan

Tanpa pernyataan ini, pembaca spec akan mengira bentuk tabel **ditemukan** di sistem lama, dan
memperlakukannya sebagai fakta yang tidak boleh diganggu.

Kenyataannya ia **rancangan**, dan rancangan boleh diperdebatkan dengan alasan. Menyamarkan
rancangan sebagai temuan menghilangkan hak orang berikutnya untuk memperbaikinya.

## Akibat

1. Modul yang skemanya dirancang baru menandainya di spec, dengan bukti bahwa skema lama memang
   tidak ada.
2. Tabel proyeksi yang sudah dipakai sistem hilir **diikuti apa adanya** - itu kontrak luar, bukan
   rancangan bebas.
3. Setiap tabel hasil rancangan baru menyebut **dari mana** daftar kolomnya berasal (ADR-0017,
   ADR-0037).
