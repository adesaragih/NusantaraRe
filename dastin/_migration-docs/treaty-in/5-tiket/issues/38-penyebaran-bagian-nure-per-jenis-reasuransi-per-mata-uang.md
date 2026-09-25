---
status: tertahan
golongan: belum pasti — L atau B, lihat PENGHALANG
---

# 38: Bagian NuRe disebarkan ke susunan retro internal per jenis reasuransi per mata uang

*Asal: `DAFTAR-PEKERJAAN.md` `P-18` dan §2.3 · `SPEC-MODEL-DATA.md` §10.6, §10.7, §10.8, §12.3, §12.4 · `SPEC-INVARIAN.md` `INV-16`, `INV-47`, `INV-50`, `INV-51`, `INV-65`.*

**What to build:** **PK** menyebarkan bagian NuRe ke **susunan retro internal**, dirinci **per jenis
reasuransi** dan **per mata uang**, dan jumlah persen penyebarannya **per induk, dikelompokkan menurut
jenis reasuransi, sama dengan 100**.

Artefak: `PENYEBARAN` dengan **dua kolom induk bernama** (`KTV-B`), `RINCIAN_PENYEBARAN`,
`NILAI_PENYEBARAN`; `INV-16` sebagai **dua** `UNIQUE`; `INV-65`; dan **`INV-47`, `INV-50`, `INV-51`
beserta pemantau kebasiannya**.

**PEMBUAT PERTAMA** untuk `PENYEBARAN`, `RINCIAN_PENYEBARAN`, `NILAI_PENYEBARAN`.

**Kenapa begini:** Tiga invarian terpenting modul ini **melintasi entitas** dan ditegakkan lewat
*materialized view* — dan **tidak satu pun punya pemantau kebasian hari ini**.

> **MV yang gagal me-refresh berhenti menegakkan tanpa satu galat pun.** Irisan yang hanya memasang
> constraint-nya akan **dinyatakan selesai padahal tidak menegakkan apa pun**, dan tidak ada yang akan
> tahu sampai data yang melanggar sudah masuk.

Dan sumbunya **sempat ditulis salah**: kedua baris rincian semula berbunyi *"per pihak"*.
`SetSpreadName.xml` mengisi `BreakDownSprdList(…).ReinsID` dari `pxResults(…).ReinsTypeID` — **nama
ruasnya menyebut pihak; yang ditulis ke dalamnya jenis**. Seluruh 866 pasangan kelas–properti disapu:
identitas pihak **tidak ada sama sekali** di keluarga penyebaran.

**Persyaratan:** **`INV-16`** (jenis reasuransi unik di dalam satu induk) — **dua** `UNIQUE`, satu per
pelekatan · **`INV-65`** (jenis reasuransi unik di dalam satu `PENYEBARAN`) · **`INV-47`** (jumlah
baris turunan menurut `SUMBU_REKONSILIASI`) · **`INV-50`** (jumlah persen penyebaran per induk,
**dikelompokkan menurut jenis reasuransi**, sama dengan 100) · **`INV-51`** (retensi + penyerahan =
100) — bergolongan **CONSTRAINT BELUM DIBUKTIKAN** · `INV-17` · `KTV-B` · `INV-58` (penyebaran adalah
**fakta terbukukan yang membawa penunjuk asalnya**, dan `INV-58` menyediakan pengecualiannya)

**Tidak termasuk:** **Cabang retro keluar** — `RETRO_KELUAR` adalah **GEL-2**, tegas di luar
penyerahan pertama. Yang di sini **susunan retro internal**, dan pembedanya **arah**.
**Tetapan 15% ke `ORS`** — `FetchQSfromMasterXOL` menulis `Pct = "15.00"` dan `ReinsTypeID = "10007"`
sebagai tetapan di dalam kode. Di sistem baru ia **baris data di susunan retro**, bukan baris kode;
irisan ini **tidak menyalinnya sebagai tetapan**. Pertanyaan apakah 15% ke `ORS` ketentuan umum atau
tambalan **terdaftar sebagai butir wawancara** dan tidak menahan irisan ini.
**Pemantau kebasian untuk `INV-69`/`INV-70`** — tiket `11`, yang **menetapkan polanya**. Irisan ini
**mengikuti** pola itu, bukan sebaliknya.

**Jalur gagal:** Kedua kolom induk terisi atau kedua-duanya kosong -> **ditolak** `CHECK` · Dua
penyebaran berjenis reasuransi sama pada induk yang sama -> **ditolak** `INV-16` · Jumlah persen
penyebaran per kelompok jenis ≠ 100 -> **ditolak** `INV-50` · MV yang **basi** -> **kegagalan
penegakan bernama** yang menyebut **sejak kapan**, bukan diam.

**Uji:** **Negatif:** kedua induk terisi; kedua induk kosong; jenis reasuransi kembar di kedua
pelekatan; jumlah persen 99; jumlah persen 101; **dan uji negatif MV yang dijalankan, bukan
diargumentasikan**.
**Positif:** `BAGIAN` bernomor 7 dan `DETAIL_PROPORSIONAL` bernomor 7, masing-masing dengan penyebaran
berjenis sama -> **keduanya diterima**.
**Positif kedua — dan ia menguji pengelompokan `INV-50` yang dikoreksi:** satu induk dengan **dua
kelompok jenis reasuransi**, masing-masing berjumlah 100 -> **diterima**. Rumusan lama —
*"per susunan"* tanpa pengelompokan — menolaknya, sebab jumlah seluruhnya 200.
**Positif ketiga:** pemantau kebasian **berbunyi kepada seseorang** ketika MV sengaja dibuat basi.

**Menggantikan:** `P-18` melestarikan `SpreadingList` dan `BreakDownSprdList`. Yang bergeser: sumbunya
**per jenis reasuransi**, bukan per pihak — dan itu koreksi yang sudah masuk `INV-50` dan
`SPEC-INVARIAN.md` §4.4a.

**PENGHALANG:**

| | |
|---|---|
| **Apa yang ditunggu (1)** | **`Uji AD`** — apakah persentase rincian penyebaran **pernah menyimpang** dari tabel master. Bila **pernah**, `P-18` bergolongan **L**; bila **tidak pernah**, ia **B** dan irisan ini **menuntut medan `CARA MENYALAKANNYA`** yang belum ada |
| **Siapa dapat menjawabnya (1)** | **kantor** — izin kueri baca-saja ke produksi |
| **Apa yang berubah (1)** | golongan tiket, dan ada-tidaknya medan `CARA MENYALAKANNYA`. **Bentuk tabelnya tidak berubah** |
| **Apa yang ditunggu (2)** | **`L-3`** — tidak ada instans Oracle yang terjangkau, sehingga uji negatif MV untuk `INV-50` dan `INV-51` **belum pernah dijalankan di mana pun** |
| **Siapa dapat menjawabnya (2)** | **kantor** — sediakan instans yang terjangkau |
| **Apa yang berubah (2)** | dari **klaim** menjadi **penegakan**. Tanpa itu, ketiga MV berstatus klaim dan tiket ini **tidak dapat dinyatakan selesai** |

**Taksiran: kosong**, dan kekosongan itu **disengaja** — `G3` melarang taksiran pada tiket yang
penghalangnya belum terjawab.

**Blocked by:** 33 · 34 · 15 · **`Uji AD`** · **`L-3`**

**Dasar:**
```
EVIDENCED(SetSpreadName@ekspor-2026-09 - BreakDownSprdList.ReinsID diisi dari pxResults.ReinsTypeID)
        EVIDENCED(sapuan 866 pasangan kelas-properti - identitas pihak NOL di keluarga penyebaran)
        EVIDENCED(FetchQSfromMasterXOL@ekspor-2026-09 - Pct="15.00", ReinsTypeID="10007" sebagai tetapan)
        DECIDED(KTV-B, INV-16, INV-50, INV-58, INV-65, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(Uji AD)
        DIASUMSIKAN-CLEAR(L-3)
```

- [ ] `PENYEBARAN`, `RINCIAN_PENYEBARAN`, `NILAI_PENYEBARAN` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `PENYEBARAN` memakai **dua kolom induk bernama** + `CHECK`, sama bentuknya dengan `POTONGAN` di irisan `37`
- [ ] `INV-16` terpasang sebagai **dua** `UNIQUE`; uji positif membuktikan keduanya tidak mencampur ruang pengenal
- [ ] `INV-50` terpasang **dengan pengelompokan menurut jenis reasuransi**, dan uji positif kedua lulus
- [ ] **uji negatif MV dijalankan** untuk `INV-47`, `INV-50`, `INV-51` — bukan diargumentasikan
- [ ] **pemantau kebasian berdiri untuk ketiganya**: `REFRESH_MODE`, `STALENESS`, terjadwal, dan **berbunyi kepada seseorang yang bernama**
- [ ] pemantau memakai **pola yang tiket `11` tetapkan**, bukan pola kedua
- [ ] kegagalan pemantau tercatat sebagai **kegagalan penegakan bernama** yang menyebut **sejak kapan**
- [ ] tetapan 15% ke `ORS` **tidak disalin sebagai tetapan**; pertanyaan wawancaranya terdaftar
- [ ] golongan tiket ditetapkan sesudah `Uji AD` kembali, dan bila **B**, medan `CARA MENYALAKANNYA` ditambahkan dengan keempat butirnya
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
