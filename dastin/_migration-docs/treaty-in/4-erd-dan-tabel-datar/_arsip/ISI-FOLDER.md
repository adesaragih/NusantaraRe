# `_arsip/` — berkas yang tergantikan, disimpan bukan dihapus

**Dibuat:** 25 September 2026

> **Kenapa disimpan, bukan dihapus.** Berkas yang pernah ada adalah bukti bahwa bentuknya pernah
> dicoba dan kenapa ia ditinggalkan. Menghapusnya membuat orang berikutnya membangunnya lagi dengan
> bentuk yang sama, dan menemukan kelemahan yang sama dari awal.

## Isi yang dituju

| Berkas | Digantikan oleh | Kenapa ditinggalkan |
|---|---|---|
| `ERD-TREATY-IN-DAN-EDM.xlsx` | `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx` | bentuk tabel biasa, **bukan** bentuk ERD `contooh.xlsx` yang ditetapkan pemilik proses |
| `Diagram-Skema-Tabel-TreatyIn-dan-EDM.xlsx` | `…-v2.xlsx` | versi pertama bentuk ERD: garis relasi belum tergambar, akar bertanda `1:1`, penanda "bersama" terpotong |

**Keduanya potret SISTEM LAMA**, sama seperti penggantinya — bukan rancangan.

## Yang TIDAK diarsipkan, dan sebabnya disebut

| Berkas | Kenapa tetap di tempatnya |
|---|---|
| `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx` | **satu-satunya gambar struktur sistem lama pada tingkat tabel.** Berbanner potret; masih dipakai sebagai AS-IS |
| `ERD-TREATY-MASUK.html` · `ERD-STRUKTUR-TREATYIN.html` · `Diagram-Skema-Tabel-TreatyMasuk.xlsx` | **basi, tetapi dirujuk 5–7 berkas.** Sudah berbanner basi, dan banner adalah perlakuan yang **lebih murah dan lebih aman** daripada memindahkan berkas yang dirujuk tujuh tempat |

> **Basi bukan berarti tidak dipakai.** Berkas basi yang **berbanner** tetap menjawab pertanyaan
> *"seperti apa bentuknya sebelum ini"*. Yang diarsipkan hanya yang **tergantikan utuh** oleh berkas
> lain yang menjawab pertanyaan yang sama dengan lebih baik.

## Berkas kunci Excel

Tiga berkas `~$…xlsx` di folder induk adalah **berkas kunci Excel**, bukan artefak. Ia muncul selama
berkasnya dibuka dan hilang sendiri saat ditutup dengan benar. Bila ia tetap ada padahal Excel sudah
ditutup, ia sisa sesi yang pernah berhenti tidak wajar dan **aman dihapus**.

## Cara memindahkan sisanya

Kedua berkas di tabel pertama **belum berpindah** — keduanya terbuka di Excel saat folder ini dibuat.
Sesudah Excel ditutup:

```
cd D:\XML_NURE\_migration-docs\treaty-in\4-erd-dan-tabel-datar
move "ERD-TREATY-IN-DAN-EDM.xlsx" _arsip\
move "Diagram-Skema-Tabel-TreatyIn-dan-EDM.xlsx" _arsip\
```

Lalu perbarui **dua** perujuknya, dan hanya dua — sudah diperiksa:

| Perujuk | Yang diubah |
|---|---|
| `../ISI-FOLDER.md` | jalurnya menjadi `_arsip/` |
| `../../alat/banner-potret-sistem-lama.py` | jalur sasarannya |

**Jangan memindahkan berkas tanpa memeriksa perujuknya lebih dulu.** Rujukan yang putus di berkas
markdown **tidak menimbulkan galat**, dan tidak ada yang menjalankannya untuk mengetahui.
