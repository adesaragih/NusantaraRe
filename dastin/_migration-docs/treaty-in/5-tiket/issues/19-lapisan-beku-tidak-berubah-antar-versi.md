---
status: aktif
golongan: baru
---

# 19: Lapisan beku tidak berubah antar versi

*Asal: `DAFTAR-PEKERJAAN.md` `P-43` · `SPEC-MODEL-DATA.md` §10.1 · `ADR-0040`.*

**What to build:** Versi baru pada kontrak yang sama **tidak dapat membawa** cedant, asal bisnis,
sifat proporsi, atau periode yang berbeda dari kontraknya. Percobaannya ditolak, dan pesannya
menyebut ruas yang menyimpang beserta nilai yang berlaku.

**Kenapa begini:** Pemisahan `KONTRAK` / `VERSI_KONTRAK` hanya berarti bila **lapisan bekunya
benar-benar beku**. Bila sebuah versi dapat menyatakan cedant yang berbeda, maka kontraknya bukan
lagi satu benda, dan pertanyaan *"siapa cedant kontrak ini"* punya jawaban yang berbeda-beda
tergantung versi mana yang dibaca. Di sistem lama pertanyaan itu memang punya jawaban berbeda-beda —
setiap addendum menyimpan salinan penuh kepala kontraknya di `JSONDATA`-nya sendiri, dan tidak ada
yang mencocokkannya.

**Persyaratan:** `ADR-0040` · `SPEC-MODEL-DATA.md` §10.1 · `INV-29` (sifat proporsi dua nilai) ·
`INV-53` (tanggal mulai ≤ tanggal berakhir) · `INV-59` (tidak ada nilai yang disalin dari entitas
lain; hilir diberi rujukan)

**Tidak termasuk:** **Pembekuan satu baris `KONTRAK` terhadap `UPDATE`** — irisan `18`.
**Mesin pro rata.** `GRL-15` memutuskan **ia sengaja tidak dibangun**: atribut *"berlaku sejak"*
dibawa, mesinnya tidak. Ini **pernyataan keputusan**, bukan `TODO` — siapa pun yang membangunnya
membalikkan `GRL-15` tanpa membukanya. Akibat yang disengaja: perubahan yang berlaku di tengah
periode **dicatat tanggalnya** dan **tidak dihitung prorata oleh sistem**.
**Keadaan versi baru saat lahir** — batch 2.

**Jalur gagal:** Menyimpan versi kedua yang membawa `ID_CEDANT` berbeda -> **ditolak**, pesannya
menyebut ruas dan nilai yang berlaku · Menyalin lapisan beku ke kolom pada `VERSI_KONTRAK` ->
**tidak ada kolomnya**; `INV-59` melarangnya dan ketiadaannya diperiksa dengan sapuan atas
`KAMUS-KOLOM.md`.

**Uji:** **Negatif:** untuk masing-masing dari kelima ruas beku, buat versi kedua yang menyimpang.
**Lima uji.**
**Positif — dan ia yang menangkap pemeriksaan yang terlalu lebar:** versi kedua yang mengubah
`NAMA_KONTRAK`, `LINGKUP_WILAYAH`, dan seluruh kepala yang **memang boleh berbeda** -> **diterima**.
Sebuah pemeriksaan yang membandingkan seluruh kepala versi terhadap versi sebelumnya lulus kelima uji
negatif dan **membekukan hal yang tidak pernah dimaksudkan beku**.

**Menggantikan:** tidak ada. **Golongannya BARU.** Sistem lama menyimpan salinan penuh kepala di
setiap addendum, sehingga *"tidak berubah antar versi"* bukan aturan melainkan **kebetulan** — dan
kebetulan yang tidak pernah diperiksa. Uji apa pun terhadap data lama yang menemukan penyimpangan
**bukan bukti sistem baru salah**; ia bukti bahwa aturan ini memang belum pernah ada.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemilik proses Treaty In, mingguan**. Pelanggaran di sini menunjuk **data warisan yang memang menyimpang**, bukan orang yang salah mengisi |
| 3 | **ambang berangka** — kuasa memblokir dinyalakan ketika penyimpangan pada baris **yang lahir di sistem baru** nol selama 4 minggu. Baris warisan **dihitung terpisah** dan tidak ikut menahan penyalaan |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(M_TREATY_IN_EDM@Table/ - tiap addendum menyimpan salinan penuh kepala di JSONDATA-nya sendiri)
        DECIDED(ADR-0040, GRL-15, INV-59)
```

- [ ] kelima ruas beku diuji **satu per satu**; lima uji negatif lulus
- [ ] pesan penolakan menyebut ruas **dan** nilai yang berlaku pada kontraknya
- [ ] uji positif lulus: kepala versi yang memang boleh berbeda tetap dapat berbeda
- [ ] sapuan atas `KAMUS-KOLOM.md` membuktikan **tidak ada** kolom salinan lapisan beku di `VERSI_KONTRAK`
- [ ] pemisahan hitungan baris warisan versus baris baru **tertulis**, bukan digabung
- [ ] ketiadaan mesin pro rata tertulis sebagai **pernyataan keputusan** yang menyebut `GRL-15`
- [ ] keempat butir **CARA MENYALAKANNYA** terisi
