---
status: aktif
golongan: baru
---

# 17: Kontrak ditemukan kembali lewat nomor warisan sistem lama, berdampingan dengan pengenal baru

*Asal: `DAFTAR-PEKERJAAN.md` `P-03` · `SPEC-MODEL-DATA.md` §10.1 · `ADR-0042`.*

**What to build:** **PK** yang memegang **nomor kontrak sistem lama** — angka yang tertulis di kertas,
di surel, dan di sistem tetangga — dapat menemukan kontraknya di sistem baru dengan mengetik nomor
itu, **tanpa perlu tahu pengenal barunya**.

**Kenapa begini:** Pada hari peralihan, **setiap rujukan yang beredar di luar sistem memakai nomor
lama**. Bordereau yang sudah dikirim, surat yang sudah ditandatangani, dan sistem hilir yang belum
dipindahkan seluruhnya menyebut nomor itu. Tanpa jalan masuk lewat nomor lama, orang yang memegang
selembar kertas **tidak punya cara menemukan barisnya** — dan yang terjadi adalah ia membuat kontrak
baru.

**Persyaratan:** `ADR-0042` (sejarah pindah apa adanya) · `SPEC-MODEL-DATA.md` §10.1
`NOMOR_KONTRAK_WARISAN` *("hanya terisi pada baris hasil migrasi")* · `INV-02` (pengenal baru tetap
dari sequence; nomor lama **tidak pernah** menjadi pengenal)

**Tidak termasuk:** **Pengisian `NOMOR_KONTRAK_WARISAN`** — itu jalur migrasi, dan migrasinya batch 2.
Irisan ini membangun **jalan masuknya** dan mengujinya dengan baris uji.
**`NOMOR_PENAWARAN_WARISAN`** — kolom kedua yang `KEPUTUSAN-TANPA-VERIFIKASI.md` §6 putuskan,
bergantung pada **Uji AL**. Ia dicari lewat jalur yang sama begitu terisi, tetapi keputusan
mencabutnya **belum jatuh**.
**Pencarian menurut cedant, periode, atau keadaan** — irisan batch 2 (`P-56`), sebab ia mencari
**menurut keadaan siklus hidup**.

**Jalur gagal:** Mencari nomor warisan yang tidak ada -> hasil **kosong yang dinyatakan**, bukan galat
dan bukan daftar penuh · Baris non-warisan bernomor warisan kosong -> **tidak muncul** pada pencarian
nomor warisan mana pun, dan **tidak menghalangi** pencarian yang lain · Nomor warisan dipakai sebagai
kunci gabung ke tabel lain -> **tidak ada jalurnya**; ia kolom pencarian, bukan pengenal.

**Uji:** **Negatif:** cari nomor yang tidak ada; pastikan baris berkolom kosong tidak ikut terambil
oleh pencarian bernilai kosong.
**Positif — dan ia yang menangkap jalur yang terlalu sempit:** dua kontrak warisan bernomor lama
**berbeda** ditemukan masing-masing, dan sebuah kontrak yang **lahir di sistem baru** — nomor
warisannya kosong — tetap ditemukan lewat pengenal barunya. Jalur pencarian yang diam-diam menuntut
nomor warisan terisi lulus setiap uji negatif yang pernah ditulis untuknya.

**Menggantikan:** tidak ada. Golongannya **BARU**: sistem lama tidak punya "nomor warisan" sebab ia
**adalah** sistem lamanya. `ID`-nya bertipe teks dan menjadi pengenal sekaligus pembawa nomor revisi
lewat pola `@substring(ID,10,12)` — dua arti di satu kolom, dan itu persis yang `ADR-0040` larang
dibawa.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **PK yang menangani pertanyaan masuk dari cedant**, pada setiap pertanyaan yang menyebut nomor lama. Bukan pemantau berkala: kegagalannya muncul saat seseorang mencari dan tidak menemukan |
| 3 | **ambang berangka** — nyala penuh ketika **seluruh** kontrak warisan sudah punya `NOMOR_KONTRAK_WARISAN` terisi, dan pencarian atas contoh acak **50 nomor lama** menemukan **50**. Di bawah itu jalurnya tetap hidup tetapi kekosongan dilaporkan, bukan didiamkan |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(SaveTreatyInDetail_Act@ekspor-2026-09 - SaveData.ID; @substring(ID,10,12) sebagai nomor revisi)
        DECIDED(ADR-0040, ADR-0042)
```

- [ ] pencarian lewat `NOMOR_KONTRAK_WARISAN` menemukan barisnya, dan hasil kosong **dinyatakan kosong**
- [ ] baris bernomor warisan kosong tidak ikut terambil oleh pencarian bernilai kosong
- [ ] nomor warisan **tidak** dipakai sebagai kunci gabung di mana pun — diperiksa dengan sapuan
- [ ] uji positif lulus: kontrak yang lahir di sistem baru tetap ditemukan lewat pengenal barunya
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk **ambang berangka**
- [ ] `NOMOR_PENAWARAN_WARISAN` dinyatakan **di luar irisan ini**, dengan Uji AL sebagai penagihnya
