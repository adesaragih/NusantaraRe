---
status: accepted
tanggal: 2026-09-23
sumber: spec Endorsement Life dan spec penyimpanan EDM Treaty In, keputusan work owner
---

# Penghapusan adalah penanda dan nilai balik, bukan hapus fisik

`[keputusan work owner]` Baris yang "dihapus" oleh pengguna **tetap tersimpan**. Yang berubah
hanya **penandanya**, dan - bila baris itu membawa nilai - ditambahkan **nilai balik** yang
menolkan pengaruhnya.

**Tidak pernah hapus fisik.**

## Bentuknya

| Peristiwa | Penanda | Nilai |
| --- | --- | --- |
| nilai diubah | penanda ubah | selisih terhadap nilai lama |
| baris baru | penanda baru | penuh, tanpa pengurang |
| baris ditandai keluar | penanda hapus | **pengurang penuh** |
| seluruh objek dibatalkan | penanda batal | **pengurang penuh untuk setiap baris** |

## Kenapa

1. Baris tetap **dapat diperiksa**. Yang hilang secara fisik tidak dapat dijelaskan kemudian hari.
2. Jumlah akhir tetap benar tanpa jalur perhitungan khusus - nilai balik menolkannya dengan
   sendirinya.
3. Sistem lama pun tidak menghapus. Perilaku ini **ditiru**, bukan diciptakan.

## Hubungannya dengan pembatalan generasi

ADR-0024 menetapkan pembatalan sebagai generasi bernilai nol. Catatan ini adalah bentuk yang sama
pada tingkat **baris**, bukan tingkat objek. Keduanya sejalan: tidak ada yang hilang, yang berubah
hanya penanda dan nilainya.

## Akibat

1. Nol perintah hapus fisik pada jalur pengguna di lapisan mana pun.
2. Pembaca hilir **wajib** menyaring pada penanda. Pembacaan tanpa penyaring akan menghitung baris
   yang sudah ditandai keluar.
3. Test yang berhasil menghapus baris secara fisik lewat jalur pengguna **gagal**.
