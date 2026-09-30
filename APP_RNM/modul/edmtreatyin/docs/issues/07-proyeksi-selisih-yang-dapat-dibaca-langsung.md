# 07: Proyeksi selisih yang dapat dibaca langsung

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Kedua anak tabel proyeksi **dibuat** — pembaca SQL langsung membutuhkan selisih sampai rincian per baris. Sesuai artefak rancangan Excel, sheet EDM Treaty In Prop.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 2.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** **06** · ⛔ `[work owner]` **anak proyeksi dibutuhkan pembaca SQL atau tidak**~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **20** *(transaksi tunggal dan skema eksplisit)*
**Menutup:** AC **23–31** *(9 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-6 · ID-7 · ID-23..ID-27 · ID-38 · ID-42

## Hasil & nilai pengguna

Angka selisih disimpan sebagai **proyeksi** yang dapat dibaca langsung dengan SQL, di luar layar
aplikasi — `[keputusan work owner]` pembaca seperti itu memang ada.

⛔⛔ **Tabel proyeksi BUKAN tabel sumber**, dan ketiga aturan yang membuatnya begitu wajib dibangun
utuh — bukan diringkas.

## Yang dibangun

### Tiga aturan yang membuatnya bukan tabel sumber

**①** Ia **hanya ditulis lapisan aplikasi**, dalam **transaksi yang sama** dengan generasinya.

**②** Ia **boleh dihapus total dan dibangun ulang** — ⛔⛔ **tetapi hanya baris yang dihitung sistem
baru**. Baris hasil migrasi **beku**.
⛔ **Tanpa aturan ini, satu perintah bangun ulang menimpa seluruh angka historis dan tidak dapat
dikembalikan.**

**③** Bila isinya berbeda dari hasil hitung ulang, ⭐ **tabelnya yang salah**, bukan operannya.

### Sisanya

Kolom asal membedakan baris **hasil migrasi** dari baris **hasil hitung**. Kunci penyaring
disertakan supaya pembaca SQL tidak perlu join balik.

⭐ **Dua proyeksi sengaja TIDAK dibuat:** untuk induk lapisan *(hanya salinan penanda yang sudah ada
di anaknya)* dan untuk rincian angsuran *(selisihnya hanya satu tingkat)*.

## Batas — yang TIDAK termasuk

⛔ Rumus selisihnya sendiri — tiket **06**.
⛔ Penanda migrasi pada baris proyeksi — tiket **09**.
⛔ **View yang menghitung selisih** — dibatalkan; jangan dihidupkan.

## Cara mengujinya

Lewat seam `repository`. ⭐ **Uji yang paling penting:** bangun ulang proyeksi dijalankan, lalu
**setiap baris hasil migrasi diperiksa satu per satu** — satu baris yang berubah berarti gagal.

⚠️ Uji itu wajib berdiri sendiri, tidak digabung dengan uji bangun ulang yang lain.

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

## ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Belum dijelaskan apakah anak tabel proyeksi — untuk penyebaran dan untuk angsuran —
memang dibutuhkan pembaca SQL, atau cukup induknya.** Membuat keduanya berarti memelihara tabel yang
mungkin tidak pernah dibaca; tidak membuatnya berarti pembaca harus join balik.

⭐ **Induknya dapat dikerjakan lebih dulu** — AC 23–28 dan 31 tidak bergantung pada keputusan itu,
dan AC 29–30 justru menyatakan proyeksi mana yang **tidak** dibuat, yang sudah pasti.