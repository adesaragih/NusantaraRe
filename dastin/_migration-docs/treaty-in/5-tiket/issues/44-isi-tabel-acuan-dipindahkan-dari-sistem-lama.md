---
status: tertahan
golongan: pelestarian
---

# 44: Isi keenam tabel acuan dipindahkan dari sistem lama

*Asal: `DAFTAR-PEKERJAAN.md` `P-50` (**pecahan tabel acuan** — lihat `Tidak termasuk`) · `ADR-0042` · `ADR-0043` · `SPEC-MODEL-DATA.md` §10.22.*

**What to build:** **PM** memindahkan isi keenam himpunan acuan — **mata uang, jenis potongan, jenis
reasuransi, bahaya, kelompok treaty, kelas bisnis** — dari sumber lamanya **apa adanya, tanpa
menghitung ulang apa pun**. Baris yang tidak dapat dipetakan menjadi **pengecualian migrasi
bernomor**, bukan baris yang didiamkan.

**Kenapa begini:** `P-50` adalah **payung**, dan `DAFTAR-PEKERJAAN.md` memerintahkan ia **dipecah per
kelompok entitas sebelum penaksiran** — kemampuan payung yang ditaksir sebagai satu **selalu**
ditaksir terlalu rendah, sebab yang ditaksir **kalimatnya**, bukan isinya. Pecahan tabel acuan
**sendirian masuk batch 1**, dan sebabnya tajam: ia **tidak menulis satu pun nilai keadaan**. Kelima
pecahan lain — kepala kontrak, daftar anak versi, cabang non-proporsional, cabang proporsional,
potongan dan penyebaran — menggantung pada pemindahan kepala, yang **menulis `KEADAAN_SIKLUS_HIDUP`**
termasuk `WARISAN_TAK_TERPETAKAN`, dan karena itu **menunggu `REV-3`**.

Dan ia **didahulukan** bukan karena kecil: **setiap kunci asing di seluruh skema menunggunya**.
Memindahkan kontrak sebelum tabel acuannya terisi berarti setiap baris ditolak kunci asing.

**Persyaratan:** `ADR-0042` (sejarah pindah apa adanya) · **`ADR-0043`** (menghitung ulang adalah
**peristiwa bisnis tersendiri**, bukan bagian migrasi) · `INV-68` (kode unik) · `INV-44` · `INV-62`

**Tidak termasuk:** **Kelima pecahan `P-50` yang lain** — **batch 2**, dan sebabnya tertulis di
`USULAN-IRISAN-TREATY-IN-BATCH-1.md` §0.2.
**Perbaikan data lama dan hitung ulang** — `G2` menyatakannya **tegas tidak masuk** penyerahan
pertama: `ADR-0043` menjadikannya peristiwa bisnis tersendiri. Baris acuan yang **salah eja** di
sistem lama dipindahkan **salah eja**, dan pembetulannya pekerjaan lain yang punya pelakunya sendiri.
**Penambahan baris acuan baru** — irisan `15` sudah membuktikan ia mungkin; yang di sini
**pemindahan isi yang sudah ada**.

**Jalur gagal:** Baris acuan lama berkode ganda -> **ditolak** `INV-68`, dan **kedua** baris dilaporkan
sebagai pengecualian bernomor — bukan salah satunya dibuang diam-diam · Baris acuan lama tanpa kode ->
pengecualian bernomor · Migrasi membetulkan ejaan -> **dilarang** `ADR-0042`; ia hitung ulang yang
menyamar sebagai kerapian · Sakelar trigger menyala saat pemindahan berjalan -> pemindahan **tidak
dimulai**.

**Uji:** **Negatif:** sisipkan sumber berkode ganda; sumber tanpa kode; jalankan pemindahan dengan
sakelar menyala.
**Positif — dan ia pokok irisan ini:** sesudah pemindahan, **cacah baris acuan baru = cacah baris
sumber dikurangi cacah pengecualian bernomor**, dan persamaan itu **menutup untuk keenam tabel**.
Sebuah pemindahan yang diam-diam membuang baris yang tidak dikenalinya lulus setiap uji negatif di
atas — sebab yang dibuangnya **tidak pernah muncul di mana pun**.
**Positif kedua:** sebuah baris acuan yang **salah eja di sumber** mendarat **salah eja** — dan itu
**keberhasilan**, bukan kegagalan. Uji paritas yang menandainya sebagai cacat **salah membaca
`ADR-0042`**.

**Menggantikan:** `P-50` bunyi lama: *"**PM** dapat memindahkan kontrak lama apa adanya, **tanpa
menghitung ulang apa pun**"* — utuh, sebagai payung. Yang bergeser: **lingkupnya dipecah**, dan
pecahan ini membawa **hanya tabel acuan**. Bunyi lengkapnya tetap berlaku untuk kelima pecahan lain.

**PENGHALANG:**

| | |
|---|---|
| **Apa yang ditunggu** | **`KTV-A`**, dan ia **bertenggat**. Presisi kolom boleh dipersempit **hanya sebelum data dimuat** — dan irisan inilah **yang memuat data pertama**. Sesudah ia berjalan, mempersempit kolom menuntut memeriksa setiap baris, dan baris yang tidak muat **tidak punya tempat pergi** |
| **Siapa dapat menjawabnya** | **gerbang pembuka sesi DDL**, bersama pemilik proses |
| **Apa yang berubah** | tidak ada isi tiket yang berubah. Yang berubah **kapan ia boleh dijalankan**: sesudah sesi DDL menutup presisinya, atau sesudah pemilik proses menyatakan presisi `KTV-A` diterima apa adanya |

**Taksiran: kosong**, dan kekosongan itu **disengaja**.

**Blocked by:** 15 · 43 · **`KTV-A`**

**Dasar:**
```
EVIDENCED(SPEC-MODEL-DATA.md §10.22 - enam tabel acuan, bentuk seragam)
        DECIDED(ADR-0042, ADR-0043, INV-68)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] keenam tabel acuan terisi dari sumber lama, dan **persamaan cacah menutup** untuk keenamnya
- [ ] pengecualian migrasi **bernomor** dan dilaporkan **per baris**, bukan sebagai satu angka
- [ ] baris berkode ganda melaporkan **keduanya**, bukan membuang salah satunya
- [ ] uji positif kedua lulus: salah eja di sumber mendarat salah eja, dan **tidak** ditandai cacat
- [ ] pemindahan berjalan **hanya** di dalam jendela sakelar irisan `43`, dan jendelanya tercatat
- [ ] **sebelum dijalankan**, `ASUMSI-CLEAR.md` diperiksa — `KTV-A` bertenggat pada momen ini, dan `KTV-2` menuntut hal yang sama untuk tiket `08`
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md` sebagai asumsi **bertenggat** yang tiket ini menutup tenggatnya
