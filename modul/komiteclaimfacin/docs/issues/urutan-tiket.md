# Urutan tiket — Komite Claim Fac In

⛔ **Ini bukan tiket.** Berkas ini hanya menjelaskan **tiket mana harus selesai sebelum tiket mana**.
Isinya **tidak menambah keputusan apa pun**.

Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5 ·
`claim-facin\RELASI-TABEL-CLAIM-FACIN.md` relasi 12–15.

---

## ⛔⛔ Nomor di dua folder BERDIRI SENDIRI

⚠️ **Modul ini dan modul Claim Fac In punya folder tiket masing-masing**, dan **keduanya bernomor
mulai dari `00`**.

> ⛔ **Tiket `00` di sini BUKAN tiket `00` di `claim-facin\issues\`.**
> ⭐ **Rujukan lintas-modul WAJIB berawalan foldernya** — contoh: `claim-facin\issues\13`.

⭐ **Sebabnya folder dipisah:** tiket memerikan **perilaku sebuah modul**, dan modul ini punya
**spec sendiri dengan 100 AC sendiri**. ⭐ **Tiket tinggal bersama spec yang ia tutup** — sehingga
`Menutup: AC …` tidak pernah ambigu.

⚠️ **Tabel basis datanya tetap satu** — `[keputusan work owner]` **K4 · K6**, komite lini FAC dan
lini PROP memakai tabel yang sama. ⭐ **Tabel bersama tidak menuntut tiket bersama.**

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⚠️ **Tabel basis datanya tetap satu** —
> `[keputusan work owner]` **K4 · K6**, komite lini FAC dan lini PROP memakai tabel yang sama."* → **"K4 · K6" di sini
> memakai nomor `STRUKTUR-TABEL-CLAIM-FACIN.md`** (`modul/claimfacin/docs/`, §0 K4 "Modul Komite tidak punya tabel baru"
> dan §5d K6 "Komite lini PROP dan lini FAC disimpan di TABEL YANG SAMA"), **bukan** K4 / K6 spec modul ini (ID-13 dua
> kotak centang, ID-9 surel galat kasir). Isinya tetap: nol tabel baru. Tetapi `T_GENERAL_KOMITE` kini **diubah** oleh
> KCF-03 (migrasi komiteclaimfacin 640–679: `MODIFY ADJUSTMENT_ID` boleh kosong + kolom `TRANSFER_TYPE` 2 / 3 / 4, izin
> work owner). Jalur folder lama di berkas ini (`komite-claim-facin\…`, `claim-facin\issues\…`) kini
> `modul/komiteclaimfacin/docs/…` dan `modul/claimfacin/docs/issues/…`.

---

## Tiga belas tiket

| # | Judul | Blocked by |
| --- | --- | --- |
| ⭐ **00** | Kasus komite lahir dari penyesuaian — jenjang dibekukan | `claim-facin\issues\13` |
| ⚠️ **01** | Pengeluaran pengaju — jenjang pertama saja | 00 |
| ⭐ **02** | Tangga berjalan — giliran, maju, selesai, tutup seketika | 00 · 01 |
| ⛔ **03** | Wewenang komite — ditegakkan di lapisan layanan | 00 · `claim-facin\issues\08` |
| **04** | Layar komite — 92 medan, dua kotak centang terkunci | 00 · 02 |
| ⭐ **05** | Keputusan tersimpan — akun dirujuk, jabatan disalin | 02 · 03 |
| ⭐ **06** | Efek akhir — terbit sekali, oleh jenjang terakhir | 05 |
| **07** | Nomor akseptasi — satu rangkaian, dikunci kelas kasus induk | 06 |
| **08** | Penggolongan lini usaha di sisi komite | 00 · `claim-facin\issues\02` |
| ⚠️ **09** | Jalur kasir — penjaga ganda-bayar eksplisit | 06 · `claim-facin\issues\11` |
| ⚠️ **10** | Jalur Fac Retro — jalur yang melompati penyiapan wewenang | 00 · 03 |
| ⭐ **11** | Jabatan, tabel login, dan jejak audit | 05 · `claim-facin\issues\12` |
| **12** | Kronologi komite dan penandaan data lama | 05 · 11 |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"| ⚠️ **01** | Pengeluaran pengaju — jenjang
> pertama saja | 00 |"*, *"| ⭐ **02** | Tangga berjalan — giliran, maju, selesai, tutup seketika | 00 · 01 |"*, *"| ⭐
> **11** | Jabatan, tabel login, dan jejak audit | 05 · `claim-facin\issues\12` |"* dan *"| **12** | Kronologi komite dan
> penandaan data lama | 05 · 11 |"* → **tiket 01 GUGUR** (K9 salah baca: `ApprovalKomite_Act` L6.2 hanya membuang calon
> yang sama dengan anggota tingkat 1) dan **tiket 11 GUGUR** (K13 digantikan workbasket, KCF-01). Tersisa **sebelas tiket
> hidup**. Gantungan baru: 02 ← 00; 12 ← 05. Tiket 00 kini membaca kasus yang sudah dilahirkan claimfacin tahap 1
> (`8c3b2e71`) dan membangun perluasan KCF-02 + TT3 / TT4 KCF-03.

---

## ⭐ Rantai terdalam

```
claim-facin\issues\00 ──▶ … ──▶ claim-facin\issues\13
                                     └─▶ 00 ──▶ 01 ──▶ 02 ──▶ 05 ──▶ 06 ──▶ 09
```

⭐ **Enam tingkat di dalam modul ini**, ⚠️ **menumpang rantai delapan tingkat** di modul Claim Fac
In. ⛔ Jadi tiket **09** adalah yang **paling jauh dari titik mulai** pada seluruh pekerjaan Fac In.

## ⭐ Yang dapat berjalan bersamaan

| Sesudah selesai | Dapat mulai bersamaan |
| --- | --- |
| **00** | 01 · **08** *(penggolongan hanya butuh kasus lahir)* · **10** |
| **02** | 04 · *(03 bila wewenang klaim sudah selesai)* |
| **05** | 06 · **11** |
| **06** | 07 · 09 |
| **11** | 12 |

> ⛔ **RALAT 10-10-2026.** Bagan dan kalimat lamanya dikutip utuh, tidak dihapus: *"└─▶ 00 ──▶ 01 ──▶ 02 ──▶ 05 ──▶ 06
> ──▶ 09"*, *"⭐ **Enam tingkat di dalam modul ini**, ⚠️ **menumpang rantai delapan tingkat** di modul Claim Fac In."*,
> *"| **00** | 01 · **08** *(penggolongan hanya butuh kasus lahir)* · **10** |"*, *"| **05** | 06 · **11** |"* dan
> *"| **11** | 12 |"* → rantai terdalam kini **claimfacin tahap 1 (`8c3b2e71`, kasus `KMT-` lahir) ──▶ 00 ──▶ 02 ──▶ 05
> ──▶ 06 ──▶ 09**: lima tingkat di modul ini, dan gantungan Claim Fac In sudah terpenuhi. Sesudah 00: 02 · 08 · 10.
> Sesudah 05: 06 · 12. Baris "11 → 12" gugur bersama tiket 11.

---

## ⚠️ Lima tiket bergantung pada modul Claim Fac In

| Tiket di sini | Menunggu | Kenapa |
| --- | --- | --- |
| **00** | `claim-facin\issues\13` | ⭐ **kontrak muatan** — kasus komite lahir dari penyerahan |
| **03** | `claim-facin\issues\08` | ⭐ **wewenang dibangun sekali**, dipakai kedua modul |
| **08** | `claim-facin\issues\02` | ⭐ **penggolongan lini usaha dibangun sekali** |
| **09** | `claim-facin\issues\11` | ⭐ **jalur kasir dibangun sekali** |
| **11** | `claim-facin\issues\12` | ⭐ **jejak audit dibangun sekali** |

⭐ **Kelimanya benar** — ⛔ keduanya **bukan** salinan yang berdiri sendiri.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Kelimanya benar** — ⛔ keduanya **bukan**
> salinan yang berdiri sendiri."* → keadaan 10-10-2026:
>
> - **00** ← `claim-facin\issues\13`: **terpenuhi** (claimfacin tahap 1 `8c3b2e71`, `CreateKMTNo_Act`).
> - **03** ← `claim-facin\issues\08`: **terpenuhi** (ADR-0030 + keanggotaan workbasket, dibangun tahap 1). Wewenang
>   komite memakai pola yang sama (`Pemegang`), dibangun di modul ini, tidak dipinjam dari claimfacin.
> - **08** ← `claim-facin\issues\02`: komite **tidak** memakai penggolong claimfacin lewat impor; halaman klaim dibaca
>   lewat kontrak `kontrak.KlaimFacInKomite` (`inti/backend/kontrak/klaimfacin.go`).
> - **09** ← `claim-facin\issues\11`: muatan kasir claimfacin masuk outbox hanya di produksi; dedupe = OQ-CFI-26.
> - **11** ← `claim-facin\issues\12`: **gugur** bersama tiket 11.

---

## ⛔ Dua tiket yang TERTAHAN

| Tiket | Penahan | Sifat |
| --- | --- | --- |
| ⛔ **03 · Wewenang** | **butir 6** — ⛔ isi daftar jabatan dan susunan jenjang **belum ada** | ⛔⛔ **MENAHAN PEMBANGUNAN.** ⭐ Tiketnya **tetap ditulis**, jahitannya disebut |
| ⛔ **11 · Jabatan & jejak audit** | **butir 6** yang sama | ⛔⛔ **MENAHAN PEMBANGUNAN** |

⚠️ **Tiket 00 tertahan sebagian** — ⭐ kasus dapat lahir, ⛔ tetapi **jumlah jenjang** menunggu butir
6, dan **butir 39** *(satu jabatan, beberapa pemegang)* menahan bila keadaan itu terjadi.

> ⛔ **RALAT 10-10-2026.** Baris dan kalimat lamanya dikutip utuh, tidak dihapus: *"| ⛔ **03 · Wewenang** | **butir 6** —
> ⛔ isi daftar jabatan dan susunan jenjang **belum ada** | ⛔⛔ **MENAHAN PEMBANGUNAN.** ⭐ Tiketnya **tetap ditulis**,
> jahitannya disebut |"*, *"| ⛔ **11 · Jabatan & jejak audit** | **butir 6** yang sama | ⛔⛔ **MENAHAN PEMBANGUNAN** |"* dan
> *"⚠️ **Tiket 00 tertahan sebagian** — ⭐ kasus dapat lahir, ⛔ tetapi **jumlah jenjang** menunggu butir 6, dan **butir
> 39** *(satu jabatan, beberapa pemegang)* menahan bila keadaan itu terjadi."* → **nol tiket tertahan.** Butir 6 tertutup
> ADR-0030 (`OUTPUT_HASIL_RNM/docs/bersama/adr/0030-aturan-peran-ditetapkan-sekali-lintas-modul.md`) + KCF-01. Butir 39
> tertutup pola Komite Prop 09-10-2026 (setiap anggota workbasket tingkat berjalan boleh memutus). Tiket 11 gugur.

---

## ⭐⭐ Liputan Acceptance Criteria

`komite-claim-facin\spec.md` memuat **100 AC**, bernomor **1–100**, ⛔ tanpa nomor hilang.

| | Jumlah |
| --- | ---: |
| ⭐ tertutup **tepat satu kali** | ⭐ **100** |
| tertutup **lebih dari satu kali** | ⛔ **0** |
| ⛔ **tidak tertutup** | ⭐ **0** |

⭐ **Dihitung DUA CARA** — *(a)* dari baris `Menutup:` tiap tiket; *(b)* sisiran pola nomor pada
seluruh berkas tiket. ✅ **Keduanya sepakat.**

⭐ **Liputan penuh, dan itu berbeda dari sisi klaim.** ⚠️ Di `claim-facin\issues\`, **5 dari 114 AC**
tidak tertutup — sebab spec di sana punya bab **"Butir yang belum punya sasaran uji"**.
⭐ `[terverifikasi]` **Spec modul ini TIDAK punya bab serupa** — 15 bab AC-nya seluruhnya perilaku
yang dapat diuji.

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"| ⭐ tertutup **tepat satu kali** | ⭐ **100** |"* →
> angka itu berlaku 20-09-2026. Sesudah RALAT 10-10-2026 di spec dan tiket, AC berikut **gugur** dan tidak lagi punya
> tiket penutup: **9–13** (tiket 01), **20–22** (eskalasi), **25** (percobaan ditolak terekam), **77 · 78 · 82 · 83**
> (tiket 11), **92–95** (K8 / KCF-04). AC **79 · 80 · 81 · 84 · 85** dari tiket 11 yang gugur tetap berlaku dengan sumber
> roster + KCF-01; ditutup migrasi roster dan tiket 03 / 05. Cacah ulang liputan dibuat sesudah PARITAS komite selesai,
> bukan di blok ini.

---

## ⭐ Dua uji yang MENGUNCI keputusan, bukan menguji fungsi

⛔ **Keduanya wajib ada, dan sebab keberadaannya wajib tertulis di tiketnya.**

| Uji | Tiket | Kenapa |
| --- | --- | --- |
| ⭐ **Pengeluaran pengaju — DUA ARAH** | **01** | ⛔ Arah kedua membuktikan pengaju di jenjang **kedua TETAP di daftar** — ⭐ **mengunci K9** supaya tidak "terperbaiki" diam-diam |
| ⭐ **Salinan jabatan** | **11** | ⛔ Naik jabatan ⇒ catatan lama **tidak berubah** — ⭐ inilah yang **membedakan salinan dari rujukan** |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"| ⭐ **Pengeluaran pengaju — DUA ARAH** | **01** |
> ⛔ Arah kedua membuktikan pengaju di jenjang **kedua TETAP di daftar** — ⭐ **mengunci K9** supaya tidak "terperbaiki"
> diam-diam |"* dan *"| ⭐ **Salinan jabatan** | **11** | ⛔ Naik jabatan ⇒ catatan lama **tidak berubah** — ⭐ inilah yang
> **membedakan salinan dari rujukan** |"* → uji pengeluaran pengaju **gugur** bersama K9; penggantinya uji L6.2 (calon
> perluasan yang sama dengan anggota tingkat 1 tidak masuk dua kali) di tiket 00. Uji salinan jabatan **pindah ke tiket
> 05**: ubah JABATAN roster sesudah keputusan ⇒ `KOMITE_JABATAN` dan teks kronologi lama tidak berubah.

---

## Bukti berkas lain tidak disentuh

| Berkas / folder | Keadaan |
| --- | --- |
| ⛔ `claim-life\` · `komite-claim-life\` · `premiumlist-life\` · `endorsement-life\` | ✅ **NOL disentuh** |
| `komite-claim-facin\spec.md` dan 5 berkas grilling/keputusan | ✅ **NOL disunting** — hanya dibaca |
| `claim-facin\issues\00-14` | ✅ **NOL disunting** |
| `claim-facin\issues\urutan-tiket.md` | ⭐ disunting **hanya pada blok RALAT** |
| `claim-prop\` · `komite-claim-prop\` · `docs\adr\` · `CLAUDE.md` · struktur mana pun | ✅ **NOL disunting** |
| korpus | ✅ md5 tidak berubah |

⛔ Kode **NOL** · DDL **NOL** · `CREATE TABLE` **NOL** · jalur berkas Go **NOL** · nomor baris XML
**NOL**.
