---
status: aktif
golongan: baru
---

# 42: Pelaksana migrasi menyimpan dokumen JSON sistem lama sebagai arsip, dan dapat menunjukkan arsipnya ada

*Asal: `DAFTAR-PEKERJAAN.md` `P-55` · `ADR-0034` · `SPEC-INVARIAN.md` `INV-61`.*

**What to build:** **PM** menyimpan dokumen `JSONDATA` sistem lama **sebagaimana adanya** pada saat
pemindahan, dan dapat **menunjukkan bahwa arsip itu ada** untuk kontrak mana pun — tanpa aplikasi
pernah membacanya sebagai sumber.

**Kenapa begini:** `JSONDATA` adalah **tempat kebenaran sistem lama** — seluruh halaman clipboard
sebagai satu kolom — dan model baru menggantikannya dengan tabel relasional. Selama beberapa tahun
pertama, pertanyaan *"apakah migrasi kehilangan sesuatu"* akan muncul, dan **satu-satunya cara
menjawabnya adalah membuka dokumen aslinya**. Tetapi arsip yang punya jalur baca akan **dipakai
sebagai sumber** oleh orang yang sedang buru-buru, dan sejak saat itu ada dua kebenaran. `ADR-0034`
menyelesaikannya dengan tegas: **arsipnya disimpan, jalur bacanya tidak ada.**

> Bentuk kemampuannya sempat ditulis sebagai *"arsip tidak punya jalur baca"* dan **gugur uji U-6** —
> itu **larangan**, bukan kemampuan, dan tidak ada pelaku yang "melakukannya". Yang nyata:
> **menyimpannya**, dipegang `PM`, dan dapat dinyatakan selesai. `INV-61` menjadi invarian yang dibawa,
> bukan tiket.

**Persyaratan:** `ADR-0034` · **`INV-61`** (dokumen JSON arsip **tidak punya jalur baca**; ia bukan
sumber kanonik) · `ADR-0042`

**Tidak termasuk:** **Jalur baca isi arsip** — **tegas tidak dibangun**, dan itu **pernyataan
keputusan**. Siapa pun yang menambahkannya membalikkan `ADR-0034` tanpa membukanya.
**Pembandingan isi arsip terhadap hasil migrasi** — itu **uji paritas migrasi**, milik irisan migrasi
batch 2.
**Kalimat untuk pembaca arsip** — **`D-7`**, diff modul Adjustment yang **masih diparkir**; dasarnya
sudah lewat dan ia menunggu satu kalimat izin pemilik proses. Irisan ini **tidak mendahuluinya**.
**Penyimpanan berkasnya di luar basis data** — bentuk fisik penyimpanan bukan keputusan irisan ini.

**Jalur gagal:** Sebuah kueri aplikasi yang **membaca isi** arsip -> **tidak ada jalurnya**, dan
ketiadaannya diperiksa dengan sapuan · Arsip tersimpan tetapi **tidak dapat ditunjukkan ada** untuk
sebuah kontrak -> kriteria selesai tidak terpenuhi · Arsip disimpan sesudah migrasi, bukan **pada
saat** migrasi -> ditolak; isinya sudah tidak asli.

**Uji:** **Negatif:** sapu seluruh basis kode untuk kueri yang menyentuh isi arsip; hasilnya **harus
nol**, dan angkanya dicetak.
**Positif — dan ia yang menangkap arsip yang tidak lengkap:** untuk **setiap** kontrak warisan yang
dipindahkan, keberadaan arsipnya dapat ditunjukkan. **Cacah arsip = cacah kontrak warisan**, dan
selisihnya **dilaporkan per kontrak**, bukan sebagai satu angka. Sebuah mekanisme yang menyimpan
arsip untuk sebagian kontrak lulus uji negatif di atas dan **kehilangan justru kontrak yang paling
aneh** — yang biasanya yang paling ingin diperiksa orang.

**Menggantikan:** tidak ada. **Golongannya BARU**: sistem lama **adalah** `JSONDATA`-nya; gagasan
"arsip" baru ada karena ada penggantinya.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **PM, satu kali pada tiap gelombang pemindahan**, atas selisih cacah arsip terhadap cacah kontrak warisan. Sesudah gelombang terakhir, **pemilik proses Treaty In, sekali** |
| 3 | **ambang berangka** — pemindahan sebuah gelombang dinyatakan lengkap ketika **selisih cacahnya nol**. Selisih bukan nol **menahan gelombang itu**, dan tiap kontrak yang hilang arsipnya disebut namanya |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(M_TREATY_IN.JSONDATA@Table/ - seluruh halaman clipboard sebagai satu kolom)
        DECIDED(ADR-0034, INV-61)
        DIASUMSIKAN-CLEAR(D-7)
```

- [ ] arsip tersimpan **pada saat** pemindahan, dan keberadaannya dapat ditunjukkan per kontrak
- [ ] sapuan basis kode untuk pembaca isi arsip kembali **nol**, dan angkanya **dicetak** — bukan disimpulkan
- [ ] uji positif lulus: cacah arsip = cacah kontrak warisan, selisihnya dilaporkan **per kontrak**
- [ ] ketiadaan jalur baca tertulis sebagai **pernyataan keputusan** yang menyebut `ADR-0034` dan `INV-61`
- [ ] `D-7` tercatat di `ASUMSI-CLEAR.md`; kalimat untuk pembaca arsip **bukan** bagian tiket ini
- [ ] keempat butir **CARA MENYALAKANNYA** terisi
