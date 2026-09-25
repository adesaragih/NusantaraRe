# Perkakas pembedah ekspor Pega — modul Treaty In Adjustment

Disimpan sebagai berkas, bukan perintah sekali pakai, supaya seluruh angka dan daftar di
`PENGETAHUAN.md` dapat **diverifikasi ulang** atas ekspor aslinya, bukan diingat.

Jalankan dengan `PYTHONIOENCODING=utf-8` — memo penulis aturan memuat karakter seperti `→`.

| Skrip | Kegunaan | Contoh |
|---|---|---|
| `pre.py` | pohon langkah Activity **beserta precondition yang benar** | `python pre.py SaveTreatyIn_EDM_Act` |
| `act.py` | pohon langkah + parameter panggilan + pasangan `Property-Set` | `python act.py TreatyEDMDifferenceShare` |
| `dt.py` | aksi Data Transform (`pyProperties`) | `python dt.py Akseptasi_DT` |
| `sect.py` | ringkasan Section/Harness: kondisi tampil | `python sect.py "…/Section/PickerTreatyInMaster.xml"` |
| `cmp2.py` | bandingkan ekspor Adjustment dengan ekspor Treaty In | `python cmp2.py` |
| `cat56.py` | katalog 56 aturan khas Adjustment | `python cat56.py` |

## Dua penanda "mati" yang berbeda — keduanya wajib dicetak

`pre.py` ada karena dua kesalahan baca yang terjadi di sesi 23 September 2026:

1. **`pyStepsBlockName == "//"`** — **langkahnya** mati, isinya tidak berjalan.
2. **`pyStepsPreCondition == "false"`** — **syaratnya** mati sementara langkahnya **tetap
   berjalan**. Teks syaratnya masih tersimpan dan terbaca meyakinkan.

Dan `pyStepsPreCondParamsWhen` **tersarang** di bawah `pyStepsPreCondParams/rowdata`, bukan anak
langsung baris langkah. Mencarinya dengan `findtext` pada baris langkah selalu mengembalikan kosong,
sehingga setiap langkah tampak tanpa syarat. `pre.py` mencarinya sebagai keturunan **tanpa
melintasi `pySteps` bersarang**, supaya syarat anak tidak tertukar dengan syarat induk.

Kode tindakan pada `pyStepsPreCondParamsWhenTrue` / `…False`:
`1` = lompat ke blok bernama (`…WhenTruePrms` menyebut namanya), `2` = lanjut, `3` = lewati langkah,
`5` = keluar aktivitas, `6` = keluar iterasi.
