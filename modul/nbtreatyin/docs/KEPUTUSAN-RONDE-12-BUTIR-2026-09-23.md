# Keputusan ronde dua belas butir — 23 September 2026

> Butir tertahan dari tiket penyimpanan NB dan EDM Treaty In, dibahas satu per satu.
> Setiap keputusan di sini milik **work owner**, bukan agen. Yang belum dijawab tetap `[terbuka]`.
> Penerapan ke spec dan tiket dilakukan **sesudah** seluruh ronde selesai, sekali jalan.

| Butir | Menahan | Keadaan |
| ---: | --- | --- |
| **1** | EDM 02 — batas endorsemen | ✅ **DIPUTUSKAN** |
| **2** | EDM 07 — anak tabel proyeksi | ✅ **DIPUTUSKAN** |
| **3** | EDM 08 — beda dagang sebaran tambahan | ✅ **DIPUTUSKAN — TABEL DIBATALKAN** |
| **4** | EDM 05 — endorse sesudah batal | ✅ **DIPUTUSKAN** |
| **5** | NB 22 · EDM 09 · 10 — lingkup pemindahan dokumen lama | ✅ **DIPUTUSKAN** |
| **6** | NB 13 — P30 berkas aturan penjumlah | ✅ **DICABUT — DI LUAR LINGKUP** |
| **7** | NB 13 — P8 muatan layanan luar | ✅ **DICABUT — IKUTI YANG TERBACA** |
| 8 | NB 18 · EDM 11 — presisi digit kiri koma | `[data DBA]` surat terkirim |
| 9 | NB 19 — panduan bentuk dokumen basi | `[data DBA]` surat terkirim |
| 10 | `ID-27b` medan `REMARK` | dikerjakan agen, tanpa keputusan |

---

## Butir 1 · Batas endorsemen — **TANPA BATAS**

**Menahan:** EDM 02, dan enam tiket di belakangnya — EDM 03 · 04 · 06 · 07 · 09 · 10.

### Keputusan

`[keputusan work owner]` **23-09-2026** — Satu polis boleh di-endorse **tanpa batas**. Batas
teknis warisan **99** tidak dipertahankan. Kolom nomor generasi dilebarkan supaya nilai di atas
99 muat.

Bunyi lama yang digantikan: ~~*"batas teknis 99, batas dagang belum ditetapkan"*~~.

### Pemeriksaan korpus sesudah keputusan

Keputusan ini diuji terhadap korpus 21 modul sebelum dicatat, karena "lebarkan kolom" mudah
memecahkan hal lain. Hasilnya **aman**, dengan tiga temuan:

| | Temuan | Akibat |
| --- | --- | --- |
| 1 | Nomor generasi dinaikkan **secara aritmetika** — `Prodke+1`, bukan penyambungan teks | ✅ pelebaran kolom tidak mengubah cara menaikkannya |
| 2 | Nol `LPAD` · nol `RPAD` · nol `TO_CHAR` bertopeng lebar pada nomor generasi | ✅ tidak ada tampilan berpadding nol yang ikut berubah |
| 3 | Satu-satunya pemotongan tetap dalam kueri generasi jatuh pada **nomor polis**, bukan nomor generasi | ⚠️ lihat di bawah |

**Cacah:** 117 berkas korpus menyebut kolom nomor generasi. Diperiksa dengan dua pola berbeda
*(pola lebar-tetap, dan pola kenaikan nilai)*.

### ⚠️ Satu hal yang **tidak** dibuka keputusan ini

Kueri pencari generasi terakhir memotong nomor polis pada **24 karakter tetap**:

```
select PRODKE from json_polis where substr(nopolis,1,24) = ?
```

Angka 24 itu **asumsi keras** dan tidak berhubungan dengan nomor generasi — jadi keputusan
"tanpa batas" tidak memperburuknya. Tetapi ia tetap rapuh sendiri: nomor polis yang panjangnya
bukan 24 akan tercocokkan salah. Ini butir lama, tetap `[terbuka]`, **bukan** bagian butir 1.

### Yang berubah di rancangan

1. Kolom nomor generasi disimpan sebagai **bilangan bulat**, bukan dua digit teks.
2. Go **tidak** memasang penolakan pada angka tertentu. Tidak ada galat "batas endorsemen
   tercapai" yang perlu ditulis.
3. Migrasi data lama berjalan apa adanya — nilai 1..99 muat tanpa pemetaan ulang.
4. Tiket EDM 02 kehilangan penahannya. Label berubah `blocked` → `ready-for-agent`, dan enam
   tiket di belakangnya ikut mengalir.

⛔ Belum diterapkan ke berkas spec dan tiket. Penerapan sekali jalan di akhir ronde.

---

## Butir 2 · Anak tabel proyeksi — **KEDUANYA DIBUAT**

**Menahan:** EDM 07, dan dua tiket di belakangnya — EDM 09 · 10.

### Keputusan

`[keputusan work owner]` **23-09-2026** — *"Ikuti saja yang sudah didiskusikan, sesuai Excel."*

`T_POLIS_DIFFERENCE_SPREADING` dan `T_POLIS_DIFFERENCE_INSTALMENT` **dibuat keduanya**. Pembaca
SQL langsung membutuhkan selisih sampai rincian per baris, bukan hanya total tingkat polis.

### Dasar

Bukan keputusan baru — pembacaan dari artefak rancangan yang sudah disepakati,
`Diagram-Skema-Tabel-NusantaraRe.xlsx`, sheet **EDM Treaty In Prop**:

| Sel | Tabel | Warna |
| --- | --- | --- |
| `J85` | `T_POLIS_DIFFERENCE` — 1:1 · TABEL PROYEKSI | ungu |
| `R103` | `T_POLIS_DIFFERENCE_SPREADING` — 1:N | ungu |
| `R116` | `T_POLIS_DIFFERENCE_INSTALMENT` — 1:N | ungu |

Ketiganya juga muncul di daftar tabel resmi sheet itu, urutan **9 · 10 · 11** dari sebelas.

### ⚠️ Satu cacah yang perlu diluruskan

Sapuan kasar pertama saya membaca sheet EDM Prop sebagai **14** tabel dan NonProp **16**, lalu
melaporkannya sebagai bertentangan dengan spec *(11 dan 15)*. **Laporan itu salah.** Sapuan itu
ikut menghitung dua hal yang bukan tabel penyimpanan sheet bersangkutan:

| Terhitung keliru | Sebenarnya |
| --- | --- |
| `T_POLIS_XOL` · `T_POLIS_XOL_LAYER` di sheet Prop | catatan abu-abu *"NOL BARIS pada proporsional"* — disebut justru untuk menyatakan **tidak dipakai** |
| `POOLDATA.TREATYINPRODUCTION` | sasaran konversi produksi, bukan tabel penyimpanan |

✅ Daftar resmi di sheet dan cacah di spec **cocok**: EDM Prop **11**, EDM NonProp **15**.

### Yang berubah di rancangan

1. Kedua tabel anak ikut ditulis Go **dalam transaksi yang sama** dengan induknya *(ID-26 ①)*.
2. Keduanya ikut aturan bangun ulang — ⛔ **hanya baris ber-`SUMBER='GO'`** yang boleh dihapus.
3. Tiket EDM 07 kehilangan penahannya. `blocked` → `ready-for-agent`, menunggu EDM 06 saja.
4. Sisi NonProp menambah anak ketiga `T_POLIS_XOL_LAYER_DIFFERENCE`, sudah ada di sheet NonProp
   — **di luar ronde ini**, karena prop diselesaikan dulu.

⛔ Belum diterapkan ke berkas spec dan tiket. Penerapan sekali jalan di akhir ronde.

---

## Butir 3 · Sebaran tambahan — ⛔ **`T_POLIS_BREAKDOWN_SPREAD` DIBATALKAN**

**Menahan:** EDM 08. Nol tiket lain di belakangnya.

### Keputusan

`[keputusan work owner]` **23-09-2026** — *"Itu tidak ada."*

⛔ **`T_POLIS_BREAKDOWN_SPREAD` tidak dibuat.** Bukan tabel sumber, bukan tabel apa pun.

Bunyi lama yang digantikan — `ID-6b`, ditulis oleh agen 23-09-2026 pagi:
~~*"`T_POLIS_BREAKDOWN_SPREAD` — tabel dasar **baru**, ditemukan dari dokumen produksi... Empat
medan: `TREATY_TYPE` · `SHARE_PERCENTAGE` · `PREMIUM_SPREADED` · `CLAIM_SPREADED`"*~~

### Bukti korpus — seluruh isinya turunan

`[terverifikasi]` 16 berkas korpus. Dua aturan menentukan:

**1 · Dua medan pertama diambil dari tabel master, bukan data polis**

`NB Treaty In/RDBList/GetBreakDownSpread_SQL.xml`:

```sql
select DISTINCT Reinstypeid AS "TreatyType", PCT AS "SharePercentage"
from pooldata.proportionalarrg a
where a.treatygroupid = ? and a.PARENTREINSTYPEID = ? and TreatyDescID = '10001'
```

**2 · Dua medan uangnya dihitung saat itu juga**

`NB Treaty In/Activity/BreakDownSpreading_Act.xml`:

```
.PremiumSpreaded = pyWorkPage.PolicyTreatyIn.TotalPremium * @divide(.SharePercentage,100,8)
.ClaimSpreaded   = @toDecimal(pyWorkPage.PolicyTreatyIn.TotalClaim) * @divide(.SharePercentage,100,8)
```

⭐ **Nol dari empat medan adalah data asli.** Dua dari master `POOLDATA.PROPORTIONALARRG`, dua
dihitung dari `TotalPremium` dan `TotalClaim` yang **sudah** tersimpan di polis. Menyimpannya
berarti menyimpan hasil hitungan — pelanggaran aturan yang sama yang membatalkan
`V_POLIS_DIFFERENCE`.

### ⚠️ Koreksi laporan agen

Bunyi lama: ~~*"Sensus korpus melewatkannya — tidak satu pun aturan merujuk anggotanya."*~~
**Keliru.** Ada 16 berkas: satu aturan pengisi, satu kueri, empat layar NB Treaty In, dan satu
aturan di EDM Treaty In.

Sebabnya **titik buta sapuan berjangkar** — sapuan diikat pada `PolicyTreatyIn.…` sehingga
rujukan relatif di dalam loop *(`.PremiumSpreaded`, `.SharePercentage`)* tidak terlihat.

⭐ Pola yang sama sudah meleset **empat kali** dalam proyek ini: `SuggestList` · `CedingCoList` ·
lima properti EDM · dan sekarang `BreakDownSpreadList`. Tiga pertama **menghilangkan** yang ada;
yang keempat **menambah** yang seharusnya tidak ada.

### Yang berubah di rancangan

| | Semula | Menjadi |
| --- | ---: | ---: |
| Tabel EDM Prop | 11 | **10** |
| Tabel EDM NonProp | 15 | **14** |
| Tabel NB Prop | — | ikut berkurang satu |

1. `ID-6b` **dicabut** di spec EDM. Bunyi lamanya dikutip, tidak dihapus.
2. Excel: sel `J62` *(Prop)* dan `J67` *(NonProp)* ditandai dibatalkan; baris `K151` dan `K177`
   dikeluarkan dari daftar tabel resmi.
3. Tiket EDM 08 kehilangan separuh lingkupnya. Sisa lingkupnya — selisih rincian angsuran —
   tetap berlaku. Penahannya hilang: `blocked` → `ready-for-agent`.
4. Go tetap **menghitung** angka sebaran tambahan saat layar memintanya, memakai persentase dari
   tabel master dan dua total yang sudah tersimpan. Tidak ada yang hilang dari layar.

### ⛔ Butir baru yang lahir dari sini

`BreakDownSprdListXOL` — varian non-proporsional, 26 kemunculan di korpus. **Belum diperiksa.**
Dugaan kuat ia turunan juga, tetapi itu dugaan. Diperiksa saat ronde NonProp, **bukan** sekarang.

⛔ Belum diterapkan ke berkas spec, tiket, dan Excel. Penerapan sekali jalan di akhir ronde.

---

## Butir 4 · Endorse sesudah pembatalan — **BOLEH, DENGAN PERINGATAN**

**Menahan:** EDM 05. Nol tiket lain di belakangnya.

### Keputusan

`[keputusan work owner]` **23-09-2026** — *"Kasih warning aja dulu."*

⭐ Polis yang sudah dibatalkan **masih boleh di-endorse lagi**. Generasi baru boleh ditumpuk di
atas generasi pembatalan. **Tidak ada gerbang peran**, tidak ada alur persetujuan tambahan.

Layar memberi **peringatan** sebelum pengguna melanjutkan. Peringatan, bukan penghalang —
pengguna tetap dapat meneruskan.

### Yang ini TIDAK berarti

⛔ Ini **bukan** keputusan bahwa persetujuan tidak akan pernah diperlukan. Kata *"dulu"* dicatat
apa adanya: bila kemudian ternyata perlu gerbang peran, keputusan itu diambil terpisah.
Rancangannya dibuat supaya penambahan itu murah — lihat butir berikutnya.

⛔ Peran `DeptHead` · `UW` · `Manager` **tidak** dilibatkan ronde ini. Cacah korpusnya dicatat
sebagai bekal kalau keputusan keras menyusul: `DeptHead` 242 kemunculan di Treaty In,
`UW` 26, `Manager` 35, `KADIV` dan `Director` masing-masing **1** — praktis tidak ada tingkat
di atas `DeptHead` di modul ini.

### Yang berubah di rancangan

1. `services` **tidak** menolak generasi baru di atas generasi pembatalan. Nol aturan larangan.
2. `repository` cukup mengembalikan **jenis generasi sebelumnya**, supaya layar tahu kapan harus
   memperingatkan. Kolom penandanya sudah ada — kolom jenis berkas yang dipakai tiket EDM 05.
3. Peringatan ditulis di **React**, bukan di Go dan bukan di Oracle. Backend hanya menyediakan
   faktanya.
4. Tiket EDM 05 kehilangan penahannya. `blocked` → `ready-for-agent`, tetap menunggu NB 19
   *(pemecah dokumen — kolom jenis berkas)*.
5. ⭐ Test wajib: rantai **generasi 1 → 2 → pembatalan → 4** tersimpan utuh dan terbaca utuh.
   Sebelumnya rantai semacam ini belum pernah diuji.

⛔ Belum diterapkan ke berkas spec dan tiket. Penerapan sekali jalan di akhir ronde.

---

## Butir 5 · Lingkup pemindahan dokumen lama — **SELURUHNYA, SEMUA GENERASI**

**Menahan:** NB 22 · EDM 09 · EDM 10 — **tiga tiket**, jangkauan terluas dari seluruh butir.

### Keputusan

`[keputusan work owner]` **23-09-2026** — Seluruh dokumen polis di `POOLDATA.JSON_POLIS`
dipindahkan ke tabel baru. **Setiap polis, setiap generasinya.** Tidak ada penyaringan menurut
status, tahun buku, atau lini usaha.

Bunyi lama yang digantikan: ~~*"dokumen lama dipindahkan seluruhnya atau sebagian — belum
diputuskan"*~~.

### Kenapa ini keputusan yang benar untuk endorsemen

Setiap polis lama membawa **rantai generasi**. Angka selisih endorsemen adalah
`generasi baru − generasi yang ditunjuk OLD_POLIS_ID`. Bila hanya generasi terakhir dipindah,
**pengurangnya hilang** — dan selisih lama tidak dapat dihitung ulang maupun diperiksa.

⭐ Keputusan ini juga menghapus kebutuhan dua sumber baca selama masa peralihan. Sistem baru
membaca satu tempat saja.

### Yang berubah di rancangan

1. Pemuat migrasi tidak memerlukan penyaring apa pun. Lebih sederhana daripada yang diandaikan.
2. Baris hasil migrasi tetap ditandai `SUMBER='PEGA'` dan **beku** — ⛔ perintah bangun ulang
   proyeksi tidak boleh menyentuhnya. Aturan ini sekarang berlaku atas **seluruh** data historis,
   bukan sebagian, sehingga taruhannya naik.

   > ⛔ **RALAT 04-10-2026 — F6** `[keputusan work owner]` (PROMPT-NB-TREATY-IN-PUTARAN-3.md bab 2 F6;
   > `PERMINTAAN-TIM-INTI.md` F6). Bunyi lama, dikutip apa adanya: *"Baris hasil migrasi tetap ditandai
   > `SUMBER='PEGA'` dan **beku**"*. Bunyi baru: penanda `SUMBER='PEGA'` **gugur** — tidak ada kolom
   > `SUMBER` di diagram grilling (`Diagram-Skema-Tabel-NusantaraRe.xlsx` sheet *NB Treaty In Prop* /
   > *NonProp*, F9–F33) dan tidak dibuat (bab 0 butir 11–12). Baris hasil pemuat dokumen lama (NB
   > tiket 22) **dikenali dari `IDPEGA` dan status**: `T_GENERAL_POLIS.IDPEGA` berisi
   > `JSON_POLIS.IDPEGA` apa adanya, berbentuk `<kelas> <pyID>` = `pyWorkPage.pzInsKey` kasus Pega
   > (`RDBList\SavePolisTreatyIn_SQL.xml`: `PEGA_JSON_POLIS_TREATYIN({pyWorkPage.pzInsKey}, …)`;
   > `repository.SetelKolomDatarLama`), sedangkan jalur biasa menulis ID kasus `NB-<n>`
   > (`repository/kasus.go` `sqlSisipGenerasi`); kasusnya `Resolved-Completed` (`TutupKasus` — dokumen
   > JSON_POLIS hanya lahir di jalur `Decision8` *Nopolis not empty* → `Utility1` → `End3`, Flow
   > `InputRealizationTreatyIn`). Larangan bangun ulang proyeksi menyentuh baris itu (milik modul EDM)
   > memakai pengenal yang sama — `IDPEGA` dan status — bukan kolom penanda. Bukti uji: `repository/lama_db_test.go`
   > `TestPemuatLamaMenulisLewatAntarmukaSama` (IDPEGA kelas + `Resolved-Completed`; bertag `db`, belum
   > dijalankan — K11).
3. Angka selisih lama **disimpan apa adanya**, tidak dihitung ulang — keputusan lama yang tetap
   berlaku dan sekarang berlaku menyeluruh.
4. Tiga tiket kehilangan penahannya: NB 22 · EDM 09 · EDM 10, seluruhnya
   `blocked` → `ready-for-agent`, tinggal menunggu ketergantungan antar tiket.

### ⛔ Butir baru yang lahir dari sini — untuk DBA

Migrasi penuh perlu direncanakan, dan perencanaannya butuh angka yang tidak kami miliki:

| | Yang diminta |
| --- | --- |
| 1 | Cacah baris `POOLDATA.JSON_POLIS` |
| 2 | Tahun paling awal yang masih tersimpan |
| 3 | Ukuran total kolom dokumen JSON |
| 4 | Rantai generasi terpanjang yang pernah ada — nilai `PRODKE` tertinggi |

⭐ Butir keempat sekaligus menguji butir 1 ronde ini: bila `PRODKE` tertinggi ternyata mendekati
99, keputusan *"tanpa batas"* terbukti mendesak, bukan sekadar kehati-hatian.

⛔ Belum diterapkan ke berkas spec dan tiket. Penerapan sekali jalan di akhir ronde.

---

## Butir 6 · P30 — ⛔ **DICABUT, ITU MILIK FAC IN**

**Menahan:** NB 13 — separuhnya. Sisanya P8, masih berlaku.

### Keputusan

`[keputusan work owner]` **23-09-2026** — *"Cabut, itu punya Fac In."*

⛔ **P30 keluar dari NB 13.** Aturan `SumTSIPremiSpreadRNMMultiCob_Act` memang hilang dari
korpus, tetapi ia **bukan milik Treaty In**, jadi ketiadaannya tidak menahan pekerjaan Treaty In.
Tidak ada berkas yang perlu diminta ke pemilik export Pega untuk lingkup ini.

Bunyi lama yang digantikan:
~~*"P30 — satu aturan penjumlah yang dipanggil tetapi tidak ada di seluruh korpus 21 modul —
`[pemilik export Pega]`"*~~

### Bukti — seluruh rantai pemanggilnya milik Fac In

`[terverifikasi]` Empat aturan ditelusuri dari yang hilang ke atas:

| Aturan | Kelas | `OfferFacIn` | `PolicyTreatyIn` |
| --- | --- | ---: | ---: |
| `CheckSpreadingProtect_ACT` | `ASM-FW-GISFW-Work` | **41** | **0** |
| └ `cekSpreadingFactIn` | `ASM-FW-GISFW-Data-Coverage` | 6 | **0** |
|  └ `SumTSIPremiSpreadedRNM_Act` | `ASM-FW-GISFW-Data-Coverage` | **47** | **0** |
|   └ `SumTSIPremiSpreadRNMMultiCob_Act` | ⛔ **tidak ada di korpus** | — | — |

Tiga hal yang saling menguatkan:

1. ⭐ **Nol** rujukan `PolicyTreatyIn` di seluruh rantai. Seluruhnya `OfferFacIn`.
2. ⭐ **Nol flow** NB Treaty In menyebut rantai ini.
3. ⭐ Kepala rantai `CheckSpreadingProtect_ACT` **tidak dipanggil siapa pun** di dalam modul
   NB Treaty In.

### Kenapa ia ada di folder Treaty In

Pega mengekspor **menurut kelas**, bukan menurut modul. `ASM-FW-GISFW-Data-Coverage` dan
`ASM-FW-GISFW-Work` dipakai bersama lintas modul, sehingga aturan Fac In ikut terbawa ke folder
ekspor Treaty In.

⭐ **Pola ini sudah pernah menipu sekali** — pohon `LocationList` lima jenjang, yang juga
dikira milik Treaty In dan dikoreksi work owner: *"itu punya facin, jangan masukkan ke treaty."*

⚠️ **Ukuran gejalanya: 67 berkas di folder NB Treaty In menyebut `OfferFacIn`.** Setiap
kesimpulan yang diambil dari folder ini tanpa memeriksa kelas dan halaman kerjanya berisiko
mengulang kekeliruan yang sama.

### Yang berubah di rancangan

1. NB 13 tinggal tertahan **satu** butir: **P8** — muatan efek keluar, `[Product+Underwriting]`.
2. Baris P30 di tabel Blocker tiket 13 dicoret, bunyinya dikutip, **tidak dihapus**.
3. ⭐ Aturan hilang itu dicatat sebagai **butir Fac In**, untuk diangkat kembali bila dan ketika
   modul Fac In digarap. Ia bukan butir Treaty In.

⛔ Belum diterapkan ke berkas tiket. Penerapan sekali jalan di akhir ronde.

---

## Butir 7 · P8 — ⛔ **DICABUT, IKUTI BENTUK YANG TERBACA**

**Menahan:** NB 13 — penahan terakhirnya. ⭐ **Tiket 13 sekarang terbuka penuh.**

### Keputusan

`[keputusan work owner]` **23-09-2026** — *"Ikuti yang terbaca, cabut P8."*

Panggilan ke layanan luar dibangun mengikuti **bentuk yang terbaca dari saudara kelas induk**.
NB 13 tidak menunggu berkas apa pun.

`[penyimpangan sadar]` Bentuk yang ditiru berasal dari **kelas yang berbeda** dari yang benar-benar
jalan di Treaty In. Keputusan diambil sadar atas risiko itu.

### Temuan — dua aturan berbeda dengan nama sama

`[terverifikasi]` Nama `serviceInsertArasapas_act` menunjuk **dua aturan**, bukan satu:

| Kelas | Versi | Ukuran | Badan |
| --- | ---: | ---: | --- |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN` · 2020-12-02 | **1** | 18 KB | ⛔ **stub, nol penugasan** |
| `ASM-FW-GISFW-WORK` · terakhir 2026-06-23 | **19** | 232 KB | ✅ **terbaca penuh, 19 penugasan** |

⭐ Pega memilih menurut hierarki kelas. `Data-PolicyTreatyIn` **lebih khusus** daripada `Work`,
sehingga yang benar-benar jalan di Treaty In justru yang badannya tidak ikut ekspor. Itulah sebab
P8 selama ini tampak buntu.

⚠️ Berkas 18 KB di folder `NB Treaty In` dan di folder `NB FacIn` **berukuran sama persis** —
keduanya stub rujukan flow, parameternya generik *(`flowName`, `flowType`, `ReferencePageName`)*.
Bukan badan aturan.

### Bentuk yang ditiru

Muatan keluarnya **empat medan**, jauh lebih kecil daripada yang diduga:

```
InputData.CARI1   = pyWorkPage.pzInsKey                                  kunci work object
InputData.CARI2   = PolicyData.PolicyNo                                  nomor polis
InputData.CARI3   = PolicyData.PolicyNo                                  nomor polis
InputData.CARI21  = ProdDateTime, format "dd/MM/yyyy HH:mm:ss"           saat produksi
```

Ditambah tiga perilaku:

1. **Dua panggilan konektor keluar** — tanda tangan `ServiceName` · `EndPointURL` · `MethodName` ·
   `ExecutionMode` muncul dua kali.
2. **Status balik dibaca** — `pyWorkPage.StatusService.ResponseMessage`.
3. **Email dikirim bila gagal**, dengan penanda layar
   `FlagErrorKonversi = "Gagal Konversi, Silahkan Coba Lagi atau Hub IT !"`.

⭐ Kata **"Konversi"** menerangkan sifat layanan itu: Arasapas adalah sistem **tujuan konversi**
polis, dan panggilan ini melaporkan polis yang sudah jadi. Ia dipakai **52 berkas lintas modul** —
Claim, Life, Fac In, RNW — jadi ia milik perusahaan, bukan milik Treaty In.

### ⛔ Kerahasiaan

Alamat surat tujuan dan satu nilai tetap `CARI3` **disunting** dari keluaran pemeriksaan dan
**tidak dicatat** di berkas mana pun. Keduanya konfigurasi, bukan rancangan. Di sistem baru
keduanya menjadi **variabel lingkungan** *(ADR-0005)*, tidak pernah ditulis harfiah.

### Yang berubah di rancangan

1. ⭐ **NB 13 tidak lagi tertahan.** Kedua penahannya gugur hari ini: P30 di luar lingkup
   *(butir 6)*, P8 dicabut *(butir ini)*. `blocked` → `ready-for-agent`.
2. Panggilan keluar dibangun di lapisan `services`, dengan alamat dan nama layanan sebagai
   konfigurasi.
3. Kegagalan **tidak** membatalkan penyimpanan polis — di sistem lama ia hanya menyalakan penanda
   layar dan mengirim surat. Perilaku itu ditiru.
4. ⚠️ **Butir pencocokan tersisa:** bila berkas aturan kelas `Data-PolicyTreatyIn` kelak
   diperoleh, bentuk yang dibangun **wajib dicocokkan**. Dicatat sebagai risiko yang diterima,
   bukan sebagai penahan.

⛔ Belum diterapkan ke berkas tiket. Penerapan sekali jalan di akhir ronde.

---

## Butir 3b · ⛔ **SEBARAN TAMBAHAN TIDAK DIMIGRASI SAMA SEKALI**

`[keputusan work owner]` **23-09-2026 sore** — *"Itu tidak usah di migrasi perhitungannya."*

⭐ Ini **memperluas** butir 3, bukan mengulanginya. Bukan hanya tabelnya yang tidak dibuat —
**perhitungannya pun tidak dibawa** ke sistem baru. Nol tabel, nol rumus di `services`, nol layar.

### Tiga bunyi berturut-turut, dua ditarik

| | Bunyi | Keadaan |
| ---: | --- | --- |
| 1 | *(pagi)* `T_POLIS_BREAKDOWN_SPREAD` tabel dasar baru, empat medan | ⛔ ditarik |
| 2 | *(sore)* tidak ada tabel; angkanya **dihitung** dari master × total tersimpan | ⛔ ditarik |
| 3 | *(sore, berlaku)* **tidak dimigrasi sama sekali** | ✅ berlaku |

⚠️ Bunyi kedua adalah **ralat agen yang masih terlalu jauh**. Ia menghapus tabelnya tetapi
mempertahankan rumusnya di lapisan `services`, padahal yang diminta adalah menghapus keduanya.

### Akibatnya

1. **AC 56 ditarik dari lingkup.** Cacah AC berlaku **58 → 57**. Nomor 56 tidak dipakai ulang
   supaya rujukan lama tidak salah arah.
2. Tiket EDM 08 kembali **murni soal penyimpanan baris** — AC 35 · 36 *(penyebaran risiko)* dan
   AC 37 · 38 *(rincian angsuran)*, seluruhnya tidak terdampak.
3. `ID-6b` tetap dicabut. Tidak ada `ID` pengganti — tidak ada yang perlu diputuskan lagi.

---

## Butir 8 · Tipe kolom nomor generasi — **BILANGAN BULAT**

`[keputusan work owner]` **23-09-2026** — Nomor generasi disimpan sebagai **bilangan bulat
berlebar wajar**, bukan teks.

| Alasan | |
| --- | --- |
| ⭐ urutan benar | teks menaruh `"100"` **sebelum** `"99"`; bilangan tidak |
| ⭐ kenaikan utuh | `Prodke+1` di sistem lama memang aritmetika — sudah diperiksa, 117 berkas |
| ⭐ migrasi bersih | nilai lama 1..99 muat tanpa pemetaan ulang |
| ⭐ nol tampilan patah | nol `LPAD` · nol `RPAD` · nol `TO_CHAR` bertopeng lebar di korpus |

⛔ Tidak ada padding nol yang perlu dipertahankan. Bila kelak layar menginginkannya, itu urusan
penyajian di React, bukan tipe kolom.

---

## Butir 10 · Medan catatan — **IKUT DIMIGRASI**

`[keputusan work owner]` **23-09-2026** — Medan catatan bebas tingkat atas *(panjang 128)*
**ikut dipindahkan**.

⭐ Alasannya menentukan: isinya **tidak tersimpan di tempat lain**. Bila tidak ikut, tidak dapat
dipulihkan.

⛔ Dua medan tampilan yang bertetangga dengannya tetap **tidak** dimigrasi — keduanya keadaan
layar, bukan data dagang. Keputusan lama, tetap berlaku.

⇒ Tiket **EDM 12** berlaku persis seperti ditulis. Tidak ada perubahan.

---

## ⭐ RONDE SELESAI — 23 September 2026

| | |
| --- | ---: |
| Butir dibahas | **12** |
| Diputuskan work owner | **9** |
| Diserahkan ke DBA | **3** |
| Tiket dibuka | **8** |
| Tiket baru | **1** |
| Penahan `[work owner]` tersisa | ⭐ **0** |

⛔ Yang tersisa seluruhnya `[data DBA]`: NB 18 · NB 19 · EDM 11 · EDM 12 — empat tiket, satu surat.
