---
status: aktif
golongan: baru
---

# 18: Kunci alami kontrak tidak dapat diubah sesudah kontraknya lahir, dan penolakannya menyebut apa yang diubah

*Asal: `DAFTAR-PEKERJAAN.md` `P-05` · `SPEC-INVARIAN.md` `INV-19` · `ADR-0040`.*

**What to build:** Sesudah sebuah kontrak tersimpan, **cedant, asal bisnis, sifat proporsi, tanggal
mulai, dan tanggal berakhirnya tidak dapat diubah oleh siapa pun**. Percobaannya ditolak, dan
pesannya **menyebut ruas mana** yang ia coba ubah.

**Kenapa begini:** Kelima ruas itu **adalah identitas kontraknya**. Mengubah salah satunya bukan
menyunting kontrak melainkan **menjadikannya kontrak yang lain**, sementara seluruh versi, layer,
penyebaran, dan catatan persetujuan di bawahnya tetap menggantung padanya. Sistem lama membiarkannya:
`SaveTreatyInDetail_Act` menulis ulang seluruh kepala pada setiap penyimpanan, tanpa satu pun
pemeriksaan bahwa yang beku tetap beku. Dan penolakan yang tidak menyebut ruasnya memaksa orang
menebak — pada borang berkolom puluhan, menebak berarti mencoba satu per satu.

**Persyaratan:** **`INV-19`** — *"kunci alami kontrak tidak berubah sepanjang hidup kontrak"*,
bergolongan **TRIGGER** · `ADR-0040` · `SPEC-MODEL-DATA.md` §10.1 · `INV-29` · `INV-53`

**Tidak termasuk:** **Peringatan kunci alami ganda saat lahir** — irisan `16`. Yang di sini soal
**mengubahnya**, yang di sana soal **kembarannya**.
**Lapisan beku yang tidak berubah antar versi** — irisan `19`. Bedanya nyata: irisan ini menjaga
**satu baris `KONTRAK`** terhadap `UPDATE`; irisan `19` menjaga **versi baru** agar tidak membawa
nilai yang berbeda.
**Sakelar penegakan trigger selama migrasi** — irisan `43`. Trigger ini **harus dapat dimatikan** saat
pemindahan, dan mekanismenya milik irisan itu.

**Jalur gagal:** `UPDATE` atas `ID_CEDANT` pada kontrak yang sudah ada -> **ditolak** `INV-19`,
pesannya menyebut `ID_CEDANT` · `UPDATE` atas dua ruas beku sekaligus -> ditolak, pesannya menyebut
**keduanya**, bukan yang pertama saja · `UPDATE` atas ruas yang **bukan** kunci alami -> **diterima**.

**Uji:** **Negatif:** ubah masing-masing dari kelima ruas, satu per satu — **lima uji, bukan satu**.
Satu uji atas satu ruas tidak membuktikan keempat lainnya dijaga, dan trigger yang hanya memeriksa
ruas pertama lulus uji tunggal mana pun.
**Positif — dan ia yang menangkap trigger yang terlalu lebar:** ubah `NAMA_KONTRAK` dan
`LINGKUP_WILAYAH` pada versi yang sama -> **diterima**. Trigger yang menolak setiap `UPDATE` atas
`KONTRAK` lulus kelima uji negatif di atas dan **tetap salah**.

**Menggantikan:** tidak ada aturan lama yang digantikan — **golongannya BARU**, dan itu yang membuat
medan `CARA MENYALAKANNYA` wajib. Sistem lama **tidak pernah memeriksanya**; kontrak yang cedantnya
diganti diam-diam akan lolos, dan **tidak ada cara mengetahui berapa kali itu terjadi**, sebab
penyimpanan tidak meninggalkan jejak siapa maupun kapan (`ADR-0045` konteks).

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemilik proses Treaty In, mingguan**, atas catatan pelanggaran mode peringatan |
| 3 | **ambang berangka** — kuasa memblokir dinyalakan ketika **nol pelanggaran tercatat selama 4 minggu berturut-turut**. Bila masih ada, tiap pelanggaran diperiksa satu per satu lebih dulu: ia dapat berarti data lama yang memang perlu dibetulkan, bukan orang yang salah |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(SaveTreatyInDetail_Act@ekspor-2026-09 - kepala ditulis ulang tiap simpan, nol pemeriksaan)
        DECIDED(ADR-0040, INV-19)
```

- [ ] `INV-19` terpasang sebagai **trigger**, bukan sebagai pemeriksaan aplikasi
- [ ] kelima ruas diuji **satu per satu**; lima uji negatif lulus
- [ ] pesan penolakan menyebut **seluruh** ruas yang dicoba, bukan yang pertama saja
- [ ] uji positif lulus: ruas non-kunci tetap dapat diubah
- [ ] trigger ini **terdaftar** di daftar yang irisan `43` matikan saat pemindahan — bukan ditemukan belakangan saat migrasi gagal
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk ambang berangka
