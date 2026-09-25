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

---

## ⛔ Dua tiket yang TERTAHAN

| Tiket | Penahan | Sifat |
| --- | --- | --- |
| ⛔ **03 · Wewenang** | **butir 6** — ⛔ isi daftar jabatan dan susunan jenjang **belum ada** | ⛔⛔ **MENAHAN PEMBANGUNAN.** ⭐ Tiketnya **tetap ditulis**, jahitannya disebut |
| ⛔ **11 · Jabatan & jejak audit** | **butir 6** yang sama | ⛔⛔ **MENAHAN PEMBANGUNAN** |

⚠️ **Tiket 00 tertahan sebagian** — ⭐ kasus dapat lahir, ⛔ tetapi **jumlah jenjang** menunggu butir
6, dan **butir 39** *(satu jabatan, beberapa pemegang)* menahan bila keadaan itu terjadi.

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

---

## ⭐ Dua uji yang MENGUNCI keputusan, bukan menguji fungsi

⛔ **Keduanya wajib ada, dan sebab keberadaannya wajib tertulis di tiketnya.**

| Uji | Tiket | Kenapa |
| --- | --- | --- |
| ⭐ **Pengeluaran pengaju — DUA ARAH** | **01** | ⛔ Arah kedua membuktikan pengaju di jenjang **kedua TETAP di daftar** — ⭐ **mengunci K9** supaya tidak "terperbaiki" diam-diam |
| ⭐ **Salinan jabatan** | **11** | ⛔ Naik jabatan ⇒ catatan lama **tidak berubah** — ⭐ inilah yang **membedakan salinan dari rujukan** |

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
