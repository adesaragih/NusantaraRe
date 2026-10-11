# 05: Keputusan tersimpan — akun dirujuk, jabatan disalin

**Status:** ready-for-agent
**Blocked by:** 02 · 03
**Menutup:** AC 36 · 37 · 38 · 39 · 40 · 41 · 42 *(7 AC)* — US 12 · 24 · 25

## Hasil & nilai pengguna

Hari ini Keputusan tiap jenjang **belum tersimpan dengan jejak yang dapat dipertanggungjawabkan** — ⚠️ dan di sistem lama, jabatan yang tercatat **diambil dari daftar yang ditanam di kode**, dengan **nilai bawaan jabatan direksi** bagi siapa pun yang tak dikenal.

Sesudah tiket ini, Setiap keputusan tersimpan bersama **akun pemutus, salinan kode jabatan saat itu, waktu, dan keputusannya** — dan ⭐ **tidak dapat diubah** sesudah tersimpan.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Hari ini Keputusan tiap jenjang **belum
> tersimpan dengan jejak yang dapat dipertanggungjawabkan** — ⚠️ dan di sistem lama, jabatan yang tercatat **diambil dari
> daftar yang ditanam di kode**, dengan **nilai bawaan jabatan direksi** bagi siapa pun yang tak dikenal."* dan
> *"Sesudah tiket ini, Setiap keputusan tersimpan bersama **akun pemutus, salinan kode jabatan saat itu, waktu, dan
> keputusannya** — dan ⭐ **tidak dapat diubah** sesudah tersimpan."* →
>
> - **Jabatan dari roster, bukan tabel login.** Jalur hidup Fac In menulis jabatan dari baris tangga:
>   `KomitePost_Adjustment` S3 / S4 `KomiteList(KomiteCount).IDKomite`, dan `IDKomite` diisi `.JABATAN` roster
>   (`ApprovalKomite_Act` L7.1, `SetListKomite_act` sisi klaim). Rantai `@if` jabatan + nilai bawaan direksi adalah sisa
>   editor (spec RALAT lintas-ronde #2), bukan rule berjalan.
> - Yang tersimpan per tingkat (`T_KOMITE_KOMITELIST`): `KOMITE_APPROVAL`, `KOMITE_COMMENT`, `DATE_APPROVE`
>   (`KomitePost_Adjustment` S7.2.1.4–S7.2.1.8), `KOMITE_OPERATORID` ditimpa **akun pemutus** (pola Komite Prop), dan
>   `KOMITE_JABATAN` = salinan teks JABATAN roster saat tangga dibentuk. Tidak ada "kode jabatan" / daftar induk
>   (K13 gugur). Tanpa kolom baru.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Penyimpan keputusan | dua rule khas modul ini, keduanya dari **salinan produksi**; ⛔ **nol gerbang wewenang** di dalamnya |
| ⚠️ Nilai bawaan jabatan | ⛔ siapa pun di luar daftar tercatat sebagai **jabatan direksi** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K13** — ⭐⭐ **Identitas orang DIRUJUK · jabatan DISALIN · tulisan jabatan DIRUJUK**
- **K11** — ⭐ Jabatan dibaca dari **data pengguna**; ⛔ bila tidak diketahui, **tolak** — ⛔ **tidak ada jabatan bawaan**

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"**K13** — ⭐⭐ **Identitas orang DIRUJUK · jabatan
> DISALIN · tulisan jabatan DIRUJUK**"* dan *"**K11** — ⭐ Jabatan dibaca dari **data pengguna**; ⛔ bila tidak
> diketahui, **tolak** — ⛔ **tidak ada jabatan bawaan**"* → K13 (kolom jabatan di tabel login) **digantikan workbasket**
> (KCF-01). Bentuk "akun dirujuk · jabatan disalin" tetap berlaku pada baris tangga (`KOMITE_OPERATORID` akun pemutus,
> `KOMITE_JABATAN` salinan), tetapi "tulisan jabatan dirujuk dari daftar induk" gugur: jabatan berupa teks roster
> `EMAILKOMITE.JABATAN`, disalin apa adanya. K11 dibaca begini: jabatan dibaca dari **baris tangga** (bukan data
> pengguna) dan **tanpa jabatan bawaan**.

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

> ⛔ **RALAT 10-10-2026.** Butir, sel, dan langkah lamanya dikutip utuh, tidak dihapus: *"⭐ **Jabatan disimpan sebagai
> SALINAN kode** — ⛔ naik jabatan **tidak mengubah catatan lama**"*, *"⭐ **Tulisan jabatan dirujuk** dari daftar induk —
> ⛔ catatan lama **tidak boleh berpindah ke jabatan yang berbeda**"*, *"⚠️ menahan pencatatan **jabatan**; ⭐ tidak
> menahan pencatatan **akun**"* dan *"3. Simpan keputusan oleh akun yang **jabatannya tidak diketahui** ⇒ ⛔
> **ditolak**, tanpa bawaan."* →
>
> - Uji salinan jabatan (dulu tiket 11, kini di sini): putuskan satu tingkat, **ubah JABATAN baris roster** sesudahnya ⇒
>   `KOMITE_JABATAN` baris tangga lama dan teks kronologinya **tidak berubah**. "Naik jabatan" orang tidak lagi
>   bermakna, sebab jabatan melekat pada workbasket, bukan pada akun.
> - Butir daftar induk **gugur** (tidak ada daftar induk).
> - Butir 6 **tertutup** (ADR-0030 + KCF-01), tidak menahan.
> - Langkah 3 diganti: Submit oleh akun yang bukan anggota workbasket tingkat berjalan ⇒ **403** (tiket 03). Jabatan
>   tidak pernah "tidak diketahui", sebab ia dibaca dari baris tangga.
