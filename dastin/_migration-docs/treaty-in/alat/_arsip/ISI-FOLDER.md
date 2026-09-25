# `alat/_arsip/` — perkakas yang tergantikan, disimpan bukan dihapus

**Dibuat:** 25 September 2026 · **Isi:** 4 berkas · ~55 KB

> **Kenapa disimpan, bukan dihapus.** Perkakas yang pernah ada adalah bukti bahwa pendekatannya
> pernah dicoba dan kenapa ia ditinggalkan. Menghapusnya membuat orang berikutnya membangunnya lagi
> dengan bentuk yang sama, dan menemukan kelemahan yang sama dari awal.

---

## Isinya

| Berkas | Digantikan oleh | Kenapa ditinggalkan |
|---|---|---|
| `cocok-tiga-sumber-entitas.py` | [`../cocok-enam-sumber-struktur.py`](../cocok-enam-sumber-struktur.py) | ia mencetak **`lulus = True`** karena **mengecualikan** `NILAI_SELISIH` dan `BESARAN_DAPAT_DISESUAIKAN` ke daftar `LUAR` atas dasar embargo `SPEC-MODEL-DATA.md` §11.3 — embargo yang sudah lewat |
| `himpunan-entitas.json` | [`../himpunan-struktur.json`](../himpunan-struktur.json) | keluaran perkakas di atas; tidak ada lagi yang membacanya |
| `urai-sepuluh.py` | [`../buat-kamus-dan-ddl.py`](../buat-kamus-dan-ddl.py) | pembaca §10 versi pertama; **nol perujuk di seluruh korpus** |
| `definisi-sepuluh.json` | [`../definisi-skema-treaty-masuk.json`](../definisi-skema-treaty-masuk.json) | keluaran `urai-sepuluh.py`; satu-satunya yang menyebutnya adalah pembuatnya sendiri |

**Dua pasang alat-dan-keluarannya**, keduanya lingkaran tertutup: pembuatnya tidak dipakai lagi, dan
keluarannya tidak dibaca siapa pun selain pembuatnya.

---

## Kenapa `cocok-tiga-sumber-entitas.py` layak dibaca meski sudah diarsipkan

Ia contoh bersih dari cacat yang berulang di proyek ini: **pemeriksa yang mengecualikan justru apa
yang seharusnya diperiksanya**. Daftar pengecualiannya ditulis jujur di dalam kodenya sendiri —

```
'NILAI_SELISIH':             'modul Adjustment -- SPEC-MODEL-DATA sec 11.3',
'BESARAN_DAPAT_DISESUAIKAN': 'modul Adjustment -- tabel acuan, sec 11.3',
```

— dan justru karena jujur, ia lolos berkali-kali: pembacanya melihat sebab yang masuk akal dan
tidak memeriksa apakah sebab itu masih berlaku.

> **Uji yang dijawab sama oleh kedua kemungkinan bukan uji.** Perkakas itu akan mencetak `LULUS`
> baik ketika kedua entitas memang di luar lingkup maupun ketika keduanya hilang karena kelalaian.

Uraiannya di [`../../4-erd-dan-tabel-datar/COCOK-ENAM-SUMBER.md`](../../4-erd-dan-tabel-datar/COCOK-ENAM-SUMBER.md)
§0 dan temuan `S-1`.

---

## Perujuk yang sudah diperbaiki — tiga, dan hanya tiga

Diperiksa lebih dulu dengan [`../hitung-perujuk-erd.py`](../hitung-perujuk-erd.py), **sebelum** satu
berkas pun dipindahkan.

| Perujuk | Yang diubah |
|---|---|
| `4-erd-dan-tabel-datar/COCOK-ENAM-SUMBER.md` | jalurnya menjadi `../alat/_arsip/` |
| `4-erd-dan-tabel-datar/ISI-FOLDER.md` | idem |
| `alat/cocok-enam-sumber-struktur.py` | docstring-nya menyebut *"kini di alat/_arsip/"* |

`himpunan-entitas.json` dan `definisi-sepuluh.json` **tidak punya perujuk di luar pembuatnya
sendiri**, dan pembuatnya ikut pindah — sehingga tidak ada tautan yang putus.

---

## Yang TIDAK diarsipkan meski nol perujuk, dan sebabnya disebut

| Berkas | Kenapa tetap di `alat/` |
|---|---|
| `banner-potret-sistem-lama.py` · `banner-html-basi.py` · `inventaris-berkas.py` | **perkakas perawatan sekali-jalan.** Nol perujuk karena memang tidak dirujuk berkas lain — ia dijalankan, bukan dibaca. Mengarsipkannya berarti menulisnya ulang lain kali |
| `hitung-perujuk-erd.py` | dibuat hari ini; ia yang **mengukur** pemindahan ini |

> **Nol penyebutan bukan bukti tidak dipakai.** Perkakas dijalankan dari baris perintah tanpa satu
> berkas pun menyebutnya. Yang diarsipkan bukan yang "tidak disebut", melainkan yang
> **tergantikan utuh** oleh perkakas lain yang mengerjakan hal yang sama dengan lebih baik.
