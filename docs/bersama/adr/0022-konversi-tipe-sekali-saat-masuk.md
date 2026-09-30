---
status: accepted
tanggal: 2026-09-23
sumber: spec penyimpanan NB Treaty In - `.scratch/nb-treaty-in/`, keputusan work owner
---

# Konversi tipe terjadi sekali saat masuk, bukan tiap kali dibaca

Dokumen JSON menyimpan hampir seluruh nilai sebagai **teks**. Konversi ke tipe yang benar terjadi
**satu kali saat dimuat**, dan kolom disimpan dengan tipe aslinya - bukan disimpan sebagai teks
lalu dikonversi tiap kali dibaca.

| Golongan | Tipe kolom |
| --- | --- |
| Uang, persen | desimal berskala tetap *(ADR-0003, ADR-0016)* |
| Tanggal | tipe tanggal |
| Bilangan | tipe bilangan |
| Kode dan penanda | **tetap teks** |

## Kenapa kode dan penanda tetap teks

`[terverifikasi]` **Nol di depan membawa makna.** Sebuah kode bernilai `"006"` yang disimpan
sebagai bilangan lalu dibaca kembali menjadi `"6"` akan **lolos uji pulang-pergi** tetapi
memecahkan penggolongnya.

Penanda juga tetap teks, sebab di sistem lama ia **dibandingkan sebagai teks**.

Karena itu **uji pulang-pergi saja tidak cukup.** Sebagian test wajib memeriksa **nilai kolom
langsung**.

## Dua hal yang mudah terlewat

1. **Dua format tanggal dalam satu dokumen** - bentuk ringkas dan cap waktu lengkap. Pemuat wajib
   mengenali keduanya.
2. **Teks kosong pada kolom angka atau tanggal menjadi kosong**, bukan nol dan bukan tanggal nol.

## Akibat

1. Konversi adalah tanggung jawab pemuat di `repository`, bukan tersebar di lapisan layanan.
2. Perbandingan nilai uang **tidak boleh sama-persis biner** - ia dibandingkan sebagai desimal.
3. Test memeriksa nilai kolom langsung untuk kode berawalan nol.
