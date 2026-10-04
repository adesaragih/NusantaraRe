# 26: Proyeksi selisih yang dapat dibaca langsung

> ## ⛔⛔ DIGANTIKAN — 23 September 2026
>
> **Tiket ini tidak dikerjakan.** Ia digantikan oleh
> **`.scratch\edm-treaty-in\issues\07`**.
>
> ### Kenapa
>
> Brief ronde tiket sebelumnya berlingkup **dua spec sekaligus**, sehingga pekerjaan **endorsemen**
> jatuh ke folder tiket **polis baru**. Ronde berikutnya membaginya ulang ke folder yang benar,
> lebih halus, dan — yang menentukan — **setiap tiket menyebut tiket NB mana yang membuat
> tabelnya**. Tanpa gate itu, tiket endorsemen bisa dikerjakan sebelum tabelnya ada.
>
> ⚠️ **Kekeliruan lingkup ini milik penyusun brief, bukan pelaksana ronde ini.**
>
> ⭐ Berkas ini **tidak dihapus** supaya jejaknya tidak hilang — tanpa ini, nomor tiket NB akan
> tampak melompat dari 23 ke 29 tanpa sebab.
>
> ### Yang berlaku
>
> | | |
> | --- | --- |
> | Tiket penyimpanan **NB** | `nb-treaty-in\issues\` **16–23** |
> | Tiket penyimpanan **EDM** | `edm-treaty-in\issues\` **01–11** |
>
> ⛔ Isi di bawah dibiarkan apa adanya sebagai catatan, **bukan sebagai pekerjaan.**

---


**Status:** ~~blocked~~ — **DIGANTIKAN**
**Blocked by:** **20** · **25** · ⛔ `[work owner]` **anak proyeksi dibutuhkan pembaca SQL atau tidak** — belum dijelaskan
**Menutup:** EDM AC **23–34** *(12 AC)*
**Sumber:** `edm-treaty-in\spec-penyimpanan-relasional.md` ID-6 · ID-7 · ID-23..ID-27 · ID-20..ID-22

## Hasil & nilai pengguna

Angka selisih disimpan sebagai **proyeksi** yang dapat dibaca langsung dengan SQL, di luar layar
aplikasi — `[keputusan work owner]` pembaca seperti itu memang ada.

⛔⛔ **Rumusnya tetap hanya di satu tempat.** View yang menghitung selisih **dibatalkan**: rumus di
dua tempat — di aplikasi dan di basis data — cepat atau lambat bercabang.

## Yang dibangun

Tabel proyeksi selisih beserta **tiga aturan yang membuatnya bukan tabel sumber**:

1. ia **hanya ditulis** lapisan aplikasi, dalam transaksi yang sama dengan generasinya;
2. ia **boleh dihapus total dan dibangun ulang** — ⛔ **hanya baris yang dihitung sistem baru**;
3. bila isinya berbeda dari hasil hitung ulang, ⭐ **tabelnya yang salah**, bukan operannya.

⭐ Kolom asal membedakan baris **hasil migrasi** (beku) dari baris **hasil hitung** (boleh dibangun
ulang). ⛔ **Tanpa penjaga ini, satu perintah bangun ulang menimpa seluruh angka historis dan tidak
dapat dikembalikan.**

Kunci penyaring disertakan supaya pembaca SQL tidak perlu join balik.

⭐ **Dua proyeksi sengaja TIDAK dibuat:** untuk induk lapisan *(hanya salinan penanda yang sudah ada
di anaknya)* dan untuk rincian angsuran *(selisihnya hanya satu tingkat)*.

Ditambah tiga pernyataan kolom khas: kolom endorsemen kosong pada polis baru, kolom polis baru
kosong pada endorsemen, dan keadaan layar **tidak disimpan**.

## Batas — yang TIDAK termasuk

⛔ Rumus selisihnya sendiri — tiket **25**.
⛔ Penanda migrasi pada baris proyeksi — tiket **28**.

## Cara mengujinya

Lewat seam `repository`. ⭐ **Uji yang paling penting:** bangun ulang proyeksi dijalankan, lalu
baris hasil migrasi diperiksa **satu per satu** — satu baris yang berubah berarti gagal.

## Acceptance criteria

- [ ] **AC 23** — rumus selisih **tidak pernah** ditulis di dalam basis data
- [ ] **AC 24** — tabel selisih **hanya ditulis** lapisan aplikasi
- [ ] **AC 25** — tabel selisih ditulis **dalam transaksi yang sama** dengan generasinya
- [ ] **AC 26** — bangun ulang **hanya menyentuh baris hasil hitung**
- [ ] **AC 27** — bila isinya beda dari hitung ulang, **tabelnya** yang dinyatakan salah
- [ ] **AC 28** — tabel selisih memuat kunci penyaring, sehingga pembaca tidak perlu join balik
- [ ] **AC 29** — **nol** proyeksi untuk induk lapisan
- [ ] **AC 30** — **nol** proyeksi untuk rincian angsuran
- [ ] **AC 31** — lapisan penyimpanan mengambil **dua baris**, bukan menjalankan agregasi
- [ ] **AC 32** — kolom khas endorsemen **kosong** pada baris polis baru
- [ ] **AC 33** — kolom khas polis baru **kosong** pada endorsemen, dan **tidak** dihapus dari skema
- [ ] **AC 34** — keadaan layar **tidak tersimpan** di tabel mana pun

## ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Belum dijelaskan apakah anak tabel proyeksi — untuk penyebaran dan untuk angsuran —
memang dibutuhkan pembaca SQL, atau cukup induknya.** Membuat keduanya berarti memelihara tabel yang
mungkin tidak pernah dibaca; tidak membuatnya berarti pembaca harus join balik ke tabel sumber.

⭐ **Induknya dapat dikerjakan lebih dulu** — AC 23–28 dan 31–34 tidak bergantung pada keputusan itu.
⛔ AC 29 dan 30 menyatakan proyeksi mana yang **tidak** dibuat, dan itu sudah pasti.