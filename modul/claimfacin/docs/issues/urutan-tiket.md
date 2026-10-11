# Urutan tiket — Claim Fac In

⛔ **Ini bukan tiket.** Berkas ini hanya menjelaskan **tiket mana harus selesai sebelum tiket mana**.
Isinya **tidak menambah keputusan apa pun**.

Sumber: `claim-facin\spec.md` · `STRUKTUR-TABEL-CLAIM-FACIN.md` · `RELASI-TABEL-CLAIM-FACIN.md`.

> ⛔⛔ **RALAT — 2026-09-20.** Kalimat di bawah ini **SALAH**, dan dikutip utuh di sini alih-alih
> dihapus:
>
> *"⚠️⚠️ **Nomor di folder ini BERSAMBUNG ke putaran sisi komite.** `[keputusan work owner]`
> **K6** — komite lini FAC dan lini PROP disimpan di tabel yang sama, jadi **tiketnya satu
> folder**. ⛔ **Tiket sisi komite belum ditulis**; ia akan mulai dari **15**. ⭐ Jangan memakai
> nomor 15 ke atas untuk apa pun yang lain."*
>
> ⭐ **Yang benar:** tiket sisi komite ada di **folder modulnya sendiri** —
> `komite-claim-facin\issues\` — bernomor **00–12**.
> ⭐ **Nomor 15 ke atas di folder ini BEBAS dipakai.**
>
> ⚠️ **Kenapa kalimat lama keliru:** ia mencampur dua hal. ⭐ **Tabelnya** memang dipakai bersama
> *(K6)*, ⛔ **tetapi tiket memerikan perilaku sebuah MODUL**, dan modul komite punya **spec
> sendiri dengan 100 AC sendiri**. ⭐ **Tiket tinggal bersama spec yang ia tutup** — sehingga
> `Menutup: AC …` tidak pernah ambigu.
>
> ⛔ **Rujukan lintas-modul WAJIB berawalan foldernya** — contoh: `komite-claim-facin\issues\00`.
> ⚠️ Tiket `00` ada di **kedua** folder dan **bukan tiket yang sama**.

---

## Lima belas tiket sisi klaim

| # | Judul | Blocked by |
| --- | --- | --- |
| ⭐ **00** | **PREFACTOR** — dua tabel baru dan kolom induk kedua, **expand–contract** | ⭐ **None** |
| **01** | Registrasi klaim, pengikatan polis, dan daur hidup kasus | 00 |
| **02** | Klasifikasi lini produk — 13 penggolong hidup, 36 tidak dialihkan | 00 · 01 |
| **03** | Estimasi nilai kerugian dan validasinya | 00 · 01 |
| **04** | Pembagian klaim per treaty, dan aturan keseragaman berbagi | 00 · 03 |
| **05** | Penyesuaian nilai klaim, pembagian dan quota share di atasnya | 00 · 03 · 04 |
| **06** | Uang, ketelitiannya, dan batas transaksi | 00 · 05 |
| **07** | Layar klaim | 01 · 03 · 05 |
| ⛔ **08** | **Wewenang** — ditegakkan di lapisan layanan | 01 · 05 |
| **09** | Dokumen akseptasi dan pencetakannya | 05 · 06 |
| **10** | Penyimpanan berkas lampiran | 01 |
| **11** | Efek keluar — kasir, surel, konversi | 06 · 09 |
| **12** | Jejak audit dan kronologi klaim | 01 · 06 |
| **13** | Penyerahan ke komite — **kontrak muatan** | 05 · 06 · 08 |
| **14** | Migrasi data lama, paritas, dan penamaan ulang | 00 · seluruh 01–13 |

> ⛔ **RALAT 10-10-2026.** Sel lamanya dikutip utuh, tidak dihapus: *"**PREFACTOR** — dua tabel baru dan kolom induk
> kedua, **expand–contract**"* dan *"Klasifikasi lini produk — 13 penggolong hidup, 36 tidak dialihkan"* →
>
> - **00**: expand–contract **ditolak** (OQ-CFI-01, 09-10-2026) — nol `MODIFY`, nol `CHECK`, hanya `ADD` (`560`–`567`);
>   baris FAC mengisi `CLAIM_ID` **dan** `OBJECT_ITEM_ID`. Lihat RALAT di tiket 00.
> - **02**: korpus memuat **60** rule `When`, dibangun **satu fungsi per `When`** di `backend/models/lini.go`; angka
>   13 / 36 / 49 tidak berdasar, dan premis "MBU / Travel tak pernah terbit" bertentangan dengan AC 18 / 109. Lihat
>   RALAT di tiket 02.
>
> Tahap 1 modul ini **dibangun 10-10-2026** (`MODUL.md`, status per tombol di `docs/PARITAS.md`); urutan di bawah adalah
> rencana 20-09-2026.

---

## ⭐ Rantai terdalam

```
00 ──▶ 03 ──▶ 04 ──▶ 05 ──▶ 06 ──▶ 09 ──▶ 11 ──▶ 14
```

⭐ **Delapan tingkat.** ⛔ Tiket **14** menunggu seluruhnya, sebab ia menguji **paritas** atas
perilaku yang dibangun tiket-tiket sebelumnya.

## ⭐ Yang dapat berjalan bersamaan

| Sesudah selesai | Dapat mulai bersamaan |
| --- | --- |
| **00** | 01 · **10** *(berkas lampiran hanya butuh registrasi)* |
| **01** | 02 · 03 · 10 |
| **05** | 06 · 07 · 08 |
| **06** | 09 · 12 |

---

## ⛔ Satu tiket yang TERTAHAN, dan sebabnya

| Tiket | Penahan | Sifat penahan |
| --- | --- | --- |
| ⛔ **08 · Wewenang** | **butir 6** — ⛔ **isi daftar jabatan dan susunan jenjang belum ada** | ⛔ **MENAHAN PEMBANGUNAN.** ⭐ Tiketnya **tetap ditulis** dan jahitannya disebut; ⛔ ia **tidak dapat selesai** sebelum work owner memberikan daftarnya |

⚠️ **Tiket 13 ikut tertahan sebagian** — ⭐ kontrak muatannya dapat ditulis, ⛔ tetapi **jumlah
jenjang** menunggu butir 6 yang sama.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⛔ **MENAHAN PEMBANGUNAN.** ⭐ Tiketnya **tetap
> ditulis** dan jahitannya disebut; ⛔ ia **tidak dapat selesai** sebelum work owner memberikan daftarnya"* → penahan
> **terjawab**: wewenang = **ADR-0030 + keanggotaan workbasket** — Input Register / Input Estimasi = pembuat kasus,
> Choose Surveyor = anggota workbasket `ReasKlaimTeknik` (`ReasPNCTeknik` XML tidak ada di DEV); server menolak tulisan
> ke kasus tertutup dan ke adjustment yang sedang di komite (`backend/services/layanan.go`). ⚠️ Tiket 08 sendiri masih
> berstatus *"ready-for-agent"* padahal tabel ini menahannya — status itu keliru sampai 10-10-2026 (RALAT di tiket 08).
> Tiket 13: tangga awal dari roster `EMAILKOMITE` FACIN; pemutus menurut jabatan → workbasket di tahap 2
> (`komiteclaimfacin`, OQ-CFI-04). Tidak ada lagi tiket yang tertahan oleh butir 6.

---

## ⭐⭐ Liputan Acceptance Criteria

`claim-facin\spec.md` memuat **114 AC**, bernomor **1–114**, ⛔ tanpa nomor hilang.

| | Jumlah |
| --- | ---: |
| ⭐ tertutup **tepat satu kali** | ⭐ **109** |
| tertutup **lebih dari satu kali** | ⛔ **0** |
| ⛔ **tidak tertutup** | ⚠️ **5** |
| **TOTAL** | **114** |

⭐ **Dihitung DUA CARA** — *(a)* dari baris `Menutup:` tiap tiket; *(b)* sisiran pola nomor AC pada
seluruh berkas tiket. ✅ **Keduanya sepakat.**

### ⚠️ Lima AC yang TIDAK tertutup — dan kenapa

⛔ **Bukan kelalaian.** ⭐ Kelimanya berasal dari bab spec berjudul **"Butir yang belum punya sasaran
uji"** — ⭐ spec **sendiri** menyatakan mereka tidak punya sasaran uji.

| AC | Isinya | Kenapa tidak dapat menjadi tiket |
| --- | --- | --- |
| **87** | Dua penunjuk bernama nyaris sama, diisi dari dua sumber berbeda; mana yang mana **tidak terbaca** | ⛔ tidak ada perilaku yang dapat diuji sebelum artinya diketahui |
| **88** | Apakah rule di luar ekspor dapat melempar ke mesin pembangkit tiket **tidak terbukti** | ⛔ buktinya ada **di luar korpus** |
| **89** | Arti sebuah nilai pada tanda tangan parameter **belum diketahui** — ⛔ **dan tidak dipakai sebagai dasar apa pun** | ⭐ tidak menyentuh perilaku |
| **90** | Isi penampung generik bernomor **tidak terbaca** dari rule yang memakainya | ⛔ artinya ditentukan halaman pembawanya, per pemakaian |
| **91** | **327 catatan pengembang belum diuji** dengan membuka rule-nya | ⛔ pekerjaan **pembacaan korpus**, bukan pembangunan |

⭐ **Kelimanya dibawa ke daftar `[terbuka]`, bukan ke tiket.** ⚠️ Bila salah satunya kelak terjawab
dan **ternyata menyentuh perilaku**, ⛔ ia **wajib melahirkan tiket baru** — dan nomornya diambil
dari deret sesudah tiket sisi komite.

> ⛔ **RALAT 10-10-2026.** Sel lamanya dikutip utuh, tidak dihapus: *"⭐ tertutup **tepat satu kali** | ⭐ **109**"* →
> dua kekeliruan di liputan ini dan di tabel *Butir `[terbuka]`* tiket-tiket:
>
> - **AC 104** ber-`[terbuka]` (butir 21 spec) tetapi "ditutup" tiket 00 — AC tanpa sasaran uji tidak dapat ditutup
>   tiket, sama seperti kelima AC di atas (RALAT di tiket 00).
> - ⚠️ **Ruang nomor butir tercampur.** Tiket menyebut "butir n" dari **dua register berbeda tanpa menyebut yang mana**:
>
>   | Register | Butir yang disebut tiket |
>   | --- | --- |
>   | **spec** (*Butir `[terbuka]` — daftar penuh*, 1–25) | **1** (tiket 00 · 01 · 06 · 09 · 10 · 14) · **6** (tiket 02 · 08 · 12 · 13) |
>   | **`STRUKTUR-TABEL-CLAIM-FACIN.md` §6** (1–13) | **9** (tiket 04) · **11** (tiket 00 · 07) · **12** (tiket 03) · **13** (tiket 05) |
>   | ⛔ **tidak ada di register mana pun** | **26 · 27** (tiket 11) · **29** (tiket 14) |
>
>   Nomor yang sama berarti hal berbeda di kedua register (mis. butir 11 spec = mesin tiket, butir 11 STRUKTUR = 29 medan
>   jenis objek). Setiap tiket kini memuat RALAT yang menyebut registernya. Pertanyaan terbuka yang lahir saat
>   membangun dicatat di `docs/OQ.md` (`OQ-CFI-nn`), bukan sebagai butir baru.

---

## Bukti berkas lain tidak disentuh

| Berkas / folder | Keadaan |
| --- | --- |
| ⛔ `claim-life\` · `komite-claim-life\` · `premiumlist-life\` · `endorsement-life\` | ✅ **NOL disentuh** |
| `claim-facin\spec.md` · `STRUKTUR` · `RELASI` | ✅ md5 tidak berubah — hanya dibaca |
| `komite-claim-facin\` · `claim-prop\` · `komite-claim-prop\` · `docs\adr\` · `CLAUDE.md` | ✅ **NOL disunting** |
| korpus | ✅ md5 tidak berubah |

> ⛔ **RALAT 10-10-2026.** Sel lamanya dikutip utuh, tidak dihapus: *"`komite-claim-facin\` · `claim-prop\` ·
> `komite-claim-prop\` · `docs\adr\` · `CLAUDE.md`"* → letak folder `.scratch` lama. Kini:
> `APP_RNM/modul/komiteclaimfacin/docs/` (tiket sisi komite di `issues/` di sana — juga untuk rujukan
> `komite-claim-facin\issues\` di RALAT 2026-09-20 di atas),
> `APP_RNM/modul/claimprop/docs/`, `APP_RNM/modul/komiteclaimprop/docs/`, dan ADR di `OUTPUT_HASIL_RNM/docs/bersama/adr/`
> — **44 berkas**: `00-INDEKS.md` + **43 ADR** (0001–0043).

⛔ Kode **NOL** · DDL **NOL** · `CREATE TABLE` **NOL** · jalur berkas Go **NOL** · nomor baris XML
**NOL** · tiket sisi komite **NOL**.
