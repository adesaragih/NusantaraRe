# Tiket `02` dan `03` — apa yang tersisa, dan satu pertentangan yang diadili

Dibuat 1 Oktober 2026, bersamaan migrasi `440`–`441`.

---

## 1 · Apa yang tersisa dari tiket `02` dan `03`

> ⚠️ **RALAT 1 Oktober 2026.** Bab ini pernah berjudul *"Tiket `02` dan `03` tidak punya sisa
> pekerjaan SKEMA"*. **Itu keliru**, dan kekeliruannya ketahuan saat daftar periksa kedua tiket dibaca
> baris demi baris alih-alih disimpulkan dari kolomnya. Keduanya **masih** punya sisa pekerjaan skema.
> Bunyi lamanya tidak dihapus — ia dicoret di tabel di bawah — sebab kekeliruan yang disamarkan akan
> diulang.

### 1.1 Yang memang sudah berdiri

| Tiket | Artefak | Keadaan |
| --- | --- | --- |
| `02` | kolom `SIFAT_MATERIAL_ADDENDUM` — tiket menyebut dirinya **PEMBUAT PERTAMA** | **sudah berdiri** sejak migrasi `401` modul `treatyin` |
| `03` | atribut tanggal berlaku pada `VERSI_KONTRAK` | **sudah berdiri** sejak migrasi `401` (`TANGGAL_BERLAKU_ADDENDUM`) |

**Sebabnya satu kalimat:** `KAMUS-KOLOM.md` §10.2 adalah **satu model untuk kedua modul** — pemisahan
25 September 2026 memindahkan **papan tiketnya**, bukan model datanya, dan §10.2 sudah memuat seluruh
kolom penyesuaian. Tiket `14`, yang menyatakan kamus itu MENGIKAT, karena itu membangun keduanya.

### 1.2 Yang MASIH tersisa, dan dua di antaranya SKEMA

| Tiket | Daftar periksa | Keadaan |
| --- | --- | --- |
| `02` | *"`SIFAT_MATERIAL_ADDENDUM` berdiri **dengan dua nilainya**, boleh kosong pada versi pertama"* | **SKEMA, belum ada.** Kolomnya ada; kedua nilainya (`MATERIAL`, `TIDAK_MATERIAL`) tidak dinyatakan di mana pun — nol `CHECK`, nol tabel acuan. **INV-62 tidak menutupinya**: ADR-0038 melarang `CHECK` untuk himpunan yang DAPAT BERTAMBAH, dan dua nilai ini **tertutup** |
| `03` | *"atribut tanggal berlaku berdiri pada `VERSI_KONTRAK` **dengan bawaan tanggal mulai kontrak**"* | **SKEMA, belum ada.** Tidak ada `DEFAULT`, dan bawaannya pun tidak dapat berupa `DEFAULT` Oracle: nilainya datang dari baris `KONTRAK` lain, bukan tetapan. Ia bawaan **jalur simpan** |
| `02` | *"daftar ruas yang dikunci tiap nilai **tertulis**"* · *"simpan yang melanggar ditolak"* · uji positif | penegakan — menunggu jalur simpan |
| `03` | *"trigger `INV-54` menolak tanggal di luar periode"* · uji positif batas inklusif | penegakan — lihat §2 |
| `03` | *"pernyataan keputusan **mesin pro rata tidak dibangun** tertulis di tiket ini dan di §10.2"* | ada di tiket; **di §10.2 belum diperiksa** — §10.2 di luar repo |

~~"Tiket `02` dan `03` tidak punya sisa pekerjaan skema."~~ — bunyi lama, **dicabut**.

**Yang benar:** pokok kedua tiket adalah penegakan di sisi simpan, dan jalur simpan menuntut
spesifikasi layar yang `L-4` nyatakan belum ada. **Tetapi dua butir skema di atas tidak menunggu
layar** — keduanya dapat dikerjakan begitu pemilik proses menjawab bentuknya (§3).

> ⛔ **Jangan menandai `02` dan `03` selesai karena kolomnya ada.** Kolom tanpa penegakan adalah
> persis bentuk yang `TDA-10` catat sebagai cacat sistem lama.

## 2 · Tiket `03` menuntut trigger yang ADR-0056 larang — diadili

| | |
| --- | --- |
| **Tiket `03`** | *"Artefak: atribut tanggal berlaku pada `VERSI_KONTRAK`, bawaannya, dan **trigger `INV-54`**."* Jalur gagalnya: *"ditolak **trigger**, pesannya menyebut periode kontraknya."* |
| **ADR-0056 (K-4)** | *"TIDAK ADA `CREATE PROCEDURE`, `CREATE FUNCTION`, maupun trigger pembawa aturan bisnis di seluruh folder ini."* Seluruh 36 berkas `ddl-usulan/` berdiri tanpa satu pun. |
| **ADR-U-0029 (repo)** | aturan bisnis ditahan di lapisan services; nol `COMMIT`, nol trigger di teks SQL. |

**DIPUTUSKAN: `INV-54` ditegakkan di lapisan services, bukan sebagai trigger basis data.**

**Kenapa:** tiga alasan yang menunjuk arah sama.

1. **ADR-0056 lebih luas daripada tiket `03`.** Ia berlaku atas seluruh folder DDL; tiket `03` satu
   irisan. Melanggarnya di sini berarti modul ini sendirian memuat satu-satunya trigger di aplikasi,
   dan pembaca berikutnya tidak akan tahu mana yang kebijakan dan mana yang kekecualian.
2. **Preseden sudah berdiri satu tiket sebelumnya.** `INV-53` — `TANGGAL_MULAI` ≤ `TANGGAL_BERAKHIR`
   pada `KONTRAK` — berbentuk **persis sama** (perbandingan dua tanggal pada satu baris), dan migrasi
   `401` sudah menyatakannya ditegakkan di services. Dua invarian sebentuk yang ditegakkan di dua
   tempat berbeda adalah cacat yang ditanam dengan tangan sendiri.
3. **`INV-54` membandingkan antar-tabel.** Tanggal berlaku ada di `VERSI_KONTRAK`, periode kontrak di
   `KONTRAK`. `CHECK` Oracle tidak dapat menyatakannya sama sekali — satu-satunya bentuk basis data
   yang bisa justru trigger, yaitu hal yang dilarang.

**Akibat:** sampai jalur simpan berdiri, `TANGGAL_BERLAKU_ADDENDUM` menerima tanggal di luar periode
kontraknya. **Itu keadaan yang diketahui**, bukan yang tersembunyi.

**Ditagih:** tiket lapisan aplikasi yang menulis jalur simpan. Daftar periksa tiket `03` menuntut
empat uji — sebelum mulai, sesudah berakhir, tepat di batas luar (ditolak), dan **tepat pada** kedua
batas (diterima, sebab batasnya inklusif, ADR-0022) — dan keempatnya dijalankan di seam services.

**Yang dapat membalikkan ini:** pemilik proses menyatakan ADR-0056 tidak berlaku bagi `INV-54`. Maka
trigger-nya ditulis, dan `INV-53` ikut dipindahkan ke sana — keduanya, bukan salah satunya.

---

## 3 · Dua pertanyaan yang menutup sisa skema `02` dan `03`

Keduanya satu kalimat, dan keduanya milik pemilik proses — bukan keputusan pengembang modul.

| # | Pertanyaan | Yang berubah menurut jawabannya |
| --- | --- | --- |
| 1 | Bolehkah `CHECK` dipasang untuk himpunan **TERTUTUP**? | ADR-0038 melarangnya untuk himpunan yang dapat bertambah; ADR-0056 melarang aturan bisnis di basis data. Dua nilai `SIFAT_MATERIAL_ADDENDUM` tidak masuk larangan pertama, dan masuk-tidaknya ke larangan kedua belum pernah diputuskan. **Jawabannya juga menentukan `INV-29`** (`SIFAT_PROPORSI`, dua nilai) dan daftar keadaan `KEADAAN_SIKLUS_HIDUP` tiket `45` — tiga tempat, satu putusan |
| 2 | Di mana bawaan tanggal berlaku berdiri? | Bukan `DEFAULT` Oracle — nilainya dibaca dari `KONTRAK.TANGGAL_MULAI` baris lain. Pilihannya: jalur simpan mengisinya, atau layar menawarkannya. Tiket `03` menulis *"dengan bawaan tanggal mulai kontrak"* tanpa menyebut pembawanya |
