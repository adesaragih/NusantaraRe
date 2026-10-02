---
status: accepted
tanggal: 2026-09-23
sumber: spec penyimpanan EDM Treaty In - `.scratch/edm-treaty-in/`, keputusan work owner
---

# Pembatalan adalah generasi bernilai nol, bukan penghapusan

Membatalkan polis **tidak menghapus** apa pun. Ia menambah **generasi baru** yang seluruh kolom
nilainya **nol**.

Tiket atau kode yang membuat jalur penghapusan untuk pembatalan **salah**.

## Kenapa

Pembatalan adalah peristiwa dagang yang harus terbaca di riwayat, bukan ketiadaan data. Sebagai
generasi bernilai nol ia ikut aturan rantai *(ADR-0018)* dan menghasilkan selisih yang benar
dengan sendirinya: `0 - nilai lama = -nilai lama`, tanpa jalur perhitungan khusus.

Sistem lama pun tidak menghapus. Perilaku ini **ditiru**, bukan diciptakan.

## Sesudah dibatalkan

`[keputusan work owner]` Polis yang sudah dibatalkan **masih boleh di-endorse lagi**. Generasi baru
boleh ditumpuk di atas generasi pembatalan.

**Nol gerbang peran**, nol alur persetujuan tambahan. Layar memberi **peringatan** sebelum pengguna
melanjutkan - peringatan, bukan penghalang.

| Lapisan | Tugasnya |
| --- | --- |
| layanan | **tidak** menolak generasi di atas pembatalan |
| `repository` | mengembalikan **jenis generasi sebelumnya** supaya layar tahu kapan memperingatkan |
| antarmuka | menulis peringatannya |

Keputusan ini diambil dengan kata *"dulu"*. Bila kelak gerbang peran diperlukan, keputusan itu
diambil terpisah - dan rancangan ini dibuat supaya penambahannya murah.

## Akibat

1. Nol jalur `DELETE` untuk pembatalan di lapisan mana pun.
2. Test wajib: rantai **generasi 1 - 2 - pembatalan - 4** tersimpan utuh dan terbaca utuh.
