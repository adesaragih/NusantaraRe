---
status: accepted
tanggal: 2026-09-23
sumber: spec penyimpanan NB dan EDM Treaty In, keputusan work owner
---

# Generasi polis adalah rantai baris, bukan penyuntingan di tempat

Endorsemen **tidak menyunting** baris polis yang ada. Ia menambah **baris baru** yang menunjuk
baris sebelumnya. Riwayat polis adalah rantai baris, dan baris lampau **tidak boleh disunting**.

## Bentuknya

| Unsur | Ketetapan |
| --- | --- |
| Penunjuk generasi sebelumnya | kolom penunjuk ke baris induk, **nullable** |
| Polis baru | penunjuk **kosong**, nomor generasi `0` |
| Endorsemen | penunjuk **terisi**, nomor generasi `>= 1` |
| Kunci alami | `(nomor polis, nomor generasi)` - **unik** |
| Larangan percabangan | penunjuk generasi sebelumnya juga **unik** |

Keunikan pada **penunjuk** itulah yang melarang percabangan: satu generasi hanya boleh punya satu
penerus. Tanpanya, dua endorsemen dapat menunjuk induk yang sama dan riwayat bercabang tanpa ada
yang tahu mana yang sah.

## Kenapa bukan menyunting di tempat

1. Angka selisih endorsemen dihitung sebagai `generasi baru - generasi yang ditunjuk`. Menyunting
   di tempat **menghapus operand pengurangnya**.
2. Nilai lama wajib dapat dibaca kembali apa adanya untuk pemeriksaan. Penyuntingan membuatnya
   hilang tanpa jejak.
3. Nomor generasi di sistem lama disimpan pada medan **dua digit**, batas 99. Sistem baru **tidak
   memakai batas itu** - nomor generasi adalah **bilangan bulat**, bukan teks dua digit. Urutannya
   harus benar secara angka: sebagai teks, `"100"` jatuh sebelum `"99"`.

## Akibat

1. Lapisan `repository` tidak menyediakan jalur pembaruan untuk baris generasi lampau.
2. Pembekuan itu ditegakkan **di lapisan layanan**, bukan hanya diandalkan pada disiplin pemakai.
3. Test yang berhasil menyunting baris generasi lampau, atau berhasil membuat dua generasi dengan
   induk yang sama, **gagal**.
4. Berlaku untuk setiap modul yang punya endorsemen - bukan hanya tempat keputusan ini lahir.

## Yang catatan ini TIDAK putuskan

Berapa kali satu polis boleh di-endorse. Tidak ada batas dagang; batas teknis lama tidak
dipertahankan.

Bagaimana baris antar generasi dipasangkan - itu ADR-0019.
