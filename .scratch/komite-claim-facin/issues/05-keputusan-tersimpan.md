# 05: Keputusan tersimpan — akun dirujuk, jabatan disalin

**Status:** ready-for-agent
**Blocked by:** 02 · 03
**Menutup:** AC 36 · 37 · 38 · 39 · 40 · 41 · 42 *(7 AC)* — US 12 · 24 · 25

## Hasil & nilai pengguna

Hari ini Keputusan tiap jenjang **belum tersimpan dengan jejak yang dapat dipertanggungjawabkan** — ⚠️ dan di sistem lama, jabatan yang tercatat **diambil dari daftar yang ditanam di kode**, dengan **nilai bawaan jabatan direksi** bagi siapa pun yang tak dikenal.

Sesudah tiket ini, Setiap keputusan tersimpan bersama **akun pemutus, salinan kode jabatan saat itu, waktu, dan keputusannya** — dan ⭐ **tidak dapat diubah** sesudah tersimpan.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Penyimpan keputusan | dua rule khas modul ini, keduanya dari **salinan produksi**; ⛔ **nol gerbang wewenang** di dalamnya |
| ⚠️ Nilai bawaan jabatan | ⛔ siapa pun di luar daftar tercatat sebagai **jabatan direksi** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K13** — ⭐⭐ **Identitas orang DIRUJUK · jabatan DISALIN · tulisan jabatan DIRUJUK**
- **K11** — ⭐ Jabatan dibaca dari **data pengguna**; ⛔ bila tidak diketahui, **tolak** — ⛔ **tidak ada jabatan bawaan**

## Yang harus diuji

- [ ] Keputusan tersimpan bersama **akun**, **salinan kode jabatan saat itu**, **waktu**, dan keputusannya
- [ ] ⭐ **Jabatan disimpan sebagai SALINAN kode** — ⛔ naik jabatan **tidak mengubah catatan lama**
- [ ] ⭐ **Akun disimpan sebagai RUJUKAN** — nama tampil yang berubah **memang** ikut berubah
- [ ] ⭐ **Tulisan jabatan dirujuk** dari daftar induk — ⛔ catatan lama **tidak boleh berpindah ke jabatan yang berbeda**
- [ ] ⛔ Jabatan pemutus **tidak diketahui** ⇒ penyimpanan **ditolak**; ⛔ **tidak ada bawaan**
- [ ] ⛔ Keputusan yang tersimpan **tidak dapat diubah**

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi daftar induk jabatan | ⚠️ menahan pencatatan **jabatan**; ⭐ tidak menahan pencatatan **akun** |

## Seam & verifikasi

**Seam:** lapisan layanan komite — penyimpanan keputusan.
1. Simpan satu keputusan ⇒ ⭐ akun, kode jabatan, waktu, dan keputusan tercatat.
2. Coba ubah keputusan yang sudah tersimpan ⇒ ⛔ **ditolak**.
3. Simpan keputusan oleh akun yang **jabatannya tidak diketahui** ⇒ ⛔ **ditolak**, tanpa bawaan.
