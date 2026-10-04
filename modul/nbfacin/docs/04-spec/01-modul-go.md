# Spesifikasi Modul Go — Siklus New Business

> Sumber: dokumen discovery NB di `D:\migrasi\RNM\OUTPUT\` (`01-activity\`, `02-layar\`, `03-celah\`),
> dibangun dari korpus `D:\migrasi\RNM\NB FacIn\` (2.083 berkas `.xml`).
> Rujukan lintas-siklus di `_ARSIP-lintas-siklus\`.
> Label mengikuti `CLAUDE.md` §3. Rancangan ditandai **[usulan]**; temuan korpus yang memaksanya
> disebut di setiap tempat.

---

## 0. Ruang lingkup

Dokumen ini mencakup **siklus New Business dari folder `NB FacIn\`**.

**Treaty Inward termasuk lingkup proyek** (K-004), tetapi dicakup **hanya sejauh ia muncul di dalam
`NB FacIn\`** — batas kerja tidak diperluas (K-005). `[terverifikasi]` 1.149 berkas korpus Treaty
Inward tersendiri tidak dibaca; 106 berkas bernuansa Treaty di dalam folder Fac In yang dipakai.
Ini batas cakupan yang disadari, bukan blocker.

Konsekuensi baiknya: tidak ada rule yang perlu dikeluarkan dari NB. `Activity\serviceInsertArasapas_act.xml`
berkelas `ASM-FW-GISFW-Data-PolicyTreatyIn` — panggilan servis konversi ke produksi Fac In yang versi
kelas `Work`-nya tidak ada di korpus — **tetap di dalam lingkup**, dan paket `production` di §2 tidak
terpotong.

Dua alasan NB aman dilakukan lebih dulu:

`[terverifikasi]` **NB dan RNW identik.** Nol rule bernama sama yang berbeda versi antara NB dan RNW
(dari 1.719 rule bersama NB↔EDM, 95 berbeda versi — seluruhnya melibatkan EDM). Implementasi NB
karena itu mencakup RNW hampir seluruhnya.

`[terverifikasi]` **Blocker versi ekspor tidak mengenai NB.** Ketidaksepadanan versi rule hanya muncul
pada salinan EDM. NB dapat direkonsiliasi tanpa menunggu pertanyaan itu dijawab.

---

## 1. Empat kendala arsitektur yang **dipaksa korpus**

Keempatnya mengubah rancangan lapisan repository dari bentuk yang biasa. Mengabaikannya menghasilkan
kode yang tampak benar tetapi berperilaku lain.

### 1.1 ⛔ Tidak ada `RDB-Save` — seluruh tulisan Oracle lewat metode **baca**

`[terverifikasi]` Hitungan langsung atas 609 Activity NB:

```powershell
foreach ($m in @('RDB-Save','RDB-List','RDB-Delete','RDB-Open','Obj-Save','Commit')) {
  $n = (Select-String -Path "D:\migrasi\RNM\NB FacIn\Activity\*.xml" `
        -Pattern "<pyStepsActivityName>$m</pyStepsActivityName>" -AllMatches | Measure-Object).Count
  "{0,-12} : {1}" -f $m, $n }
```

| Metode | Kemunculan |
| --- | ---: |
| `RDB-Save` | **0** |
| `RDB-List` | 614 |
| `RDB-Delete` | 2 |
| `RDB-Open` | 0 |
| `Obj-Save` | 45 |
| `Commit` | 5 |

Dan **55 dari 57** pernyataan SQL yang menulis tersimpan di tag `pyBrowseSQL` — tag untuk *browse*.

**Konsekuensi [usulan]:** antarmuka repository **tidak boleh** diturunkan dari nama metode Pega.
Arah operasi (baca vs tulis) hanya boleh ditentukan dengan membaca **verba SQL di dalam** rule
`RDBList` yang dipanggil. Setiap fungsi repository Go menyebut rule asalnya **dan** verba SQL-nya:

```go
// Asal: NB FacIn/RDBList/InsertOfferProduction_Sql.xml
// PERINGATAN: rule ini dipanggil lewat langkah RDB-List (metode baca) dan SQL-nya
// tersimpan di tag pyBrowseSQL, tetapi isinya INSERT. Arah operasi = tulis.
func (r *OfferRepo) InsertOfferProduction(ctx context.Context, ...) error
```

### 1.2 ⛔ `COMMIT` berada di dalam SQL, bukan di Pega

`[terverifikasi]` 36 pernyataan SQL memuat `COMMIT`, sementara hanya **5** langkah `Commit` ada di
seluruh 609 Activity. Batas transaksi **milik database**, bukan aplikasi.

⛔ **Dan 21 SQL penulis tidak memuat `COMMIT` sama sekali** — termasuk jalur produksi
(`InsertTreatyProduction_Sql`, `InsertOfferProduction_Sql`, `UpdatePolisEndorsement_SQL`), yang juga
tidak punya langkah `Commit` di jalurnya. **[pertanyaan terbuka] MEMBLOKIR:** siapa yang meng-commit
jalur produksi? Autocommit driver, atau memang tidak pernah ter-commit?

**Konsekuensi [usulan]:** jangan membungkus panggilan repository dalam `tx, _ := db.Begin()` dengan
asumsi biasa. Sebelum pertanyaan di atas dijawab, lapisan repository **mereproduksi** perilaku lama:
tiap pernyataan dieksekusi apa adanya, dan `COMMIT` yang tertanam di SQL dibiarkan. Setiap
penyimpangan dicatat, bukan diam-diam "dirapikan".

### 1.3 ⛔ Seluruh tulisan produksi melewati stored procedure yang isinya tidak diketahui

`[terverifikasi]` 35 stored procedure `POOLDATA.*` dipanggil; **tidak satu pun badan prosedurnya ada
di korpus** — korpus juga memuat nol DDL.

**Konsekuensi:** lapisan `repository/oracle` untuk jalur produksi **tidak dapat dispesifikasikan**
sekarang. Yang dapat dispesifikasikan adalah *pemanggilnya*: nama prosedur, parameter, dan urutan
panggilan. Isi prosedur menunggu `ALL_SOURCE` dari DBA.

### 1.4 ⚠️ Langkah bersarang sampai kedalaman 9

`[terverifikasi]` 11.090 langkah = 4.069 tingkat atas + **7.021 bersarang**. Membaca hanya `pySteps`
tingkat atas kehilangan **63 %** logika.

**Konsekuensi:** ini jebakan bagi siapa pun yang mengaudit ulang, dan alasan mengapa porting
per-langkah harus memakai parser rekursif, bukan grep tingkat atas.

---

## 2. Pemetaan paket

`[usulan]`, diturunkan dari inventaris 609 Activity dengan dua aturan penimpa berbasis fakta
(penulis Oracle → `repository`; pemanggil Connect-REST → `integration`).

| Paket | Jumlah activity | Isi |
| --- | ---: | --- |
| `services/faccase` | 275 | siklus hidup case, data penawaran, lokasi & objek risiko |
| `services/premium` | 97 | rumus premi, komisi, brokerage |
| `services/production` | 81 | konversi ke produksi, pasangan nilai-sesudah/delta |
| `services/spreading` | 79 | spreading, kapasitas, proteksi, scoring |
| `services/integration` | 61 | panggilan servis luar |
| `services/acceptance` | 16 | tangga persetujuan |

⚠️ **Jumlah activity bukan ukuran kompleksitas.** `acceptance` hanya 16 activity tetapi merupakan
inti aplikasi; `spreading` 79 activity tetapi sebagian besarnya varian per lini bisnis.

---

## 3. `services/acceptance` — mesin tangga

`[terverifikasi]` Mesin keadaan **satu langkah per keputusan manusia**, bukan loop. Rinciannya di
`_ARSIP-lintas-siklus\01-flow\04-mesin-akseptasi.md`.

```go
// Satu transisi = satu keputusan manusia. JANGAN merancang ini sebagai loop yang
// menghitung seluruh rantai approver sekaligus — bentuk itu tidak ada di sistem lama
// dan berperilaku berbeda saat rantai terputus di tengah.
type Decision struct {
    CaseID   string
    Decision int    // ProposalAcceptStatus — arti nilainya BELUM DIKETAHUI, lihat §6
    Note     string
}

type LadderState struct {
    CurrentQueue         string // PositionNote — ruang nama ANTREAN
    NextApproverPosition string // LetterNo     — ruang nama KODE JABATAN. Bukan nomor surat.
    LastDecision         int
    Finished             bool   // tidak ada To* yang cocok → tangga selesai, BUKAN galat
}
```

⚠️ `[terverifikasi]` **Dua ruang nama yang tidak boleh disatukan**: token antrean
(`ReasFacInUnderwriting`) dan kode jabatan (`SENIORUW`). Tipe Go yang berbeda untuk keduanya mencegah
tertukar saat kompilasi — ini salah satu tempat di mana sistem tipe dapat menangkap cacat warisan.

⛔ `[pertanyaan terbuka]` `GetLimitAkseptasi_ActFlow` (eksklusif NB) memuat **15 tautologi**
pembanding limit dan 4 nomor polis literal sebagai gerbang alur. Sebelum bisnis memutuskan, port
**mereproduksi** tautologi itu apa adanya dan menandainya dengan komentar.

---

## 4. `services/spreading` — tiga mesin, bukan satu

`[terverifikasi]` (`03-celah\01-spreading-capacity-scoring-nb.md`) NB memuat **tiga mesin spreading
yang berbeda**, dengan presisi berbeda:

| Mesin | Cakupan | Presisi |
| --- | --- | --- |
| A | per-coverage Fac In | 20 desimal |
| B | kapasitas treaty QS/SPL (`GetKapasitasTreaty`) | 20 desimal |
| C | `PolicyTreatyIn.SpreadingRiskList` — satu-satunya dengan pembagian rata `100/n` | 10 dan 8 desimal |

**Empat sifat yang mengubah rancangan:**

1. `[terverifikasi]` **Spreading tidak mengubah `LetterNo` maupun `PositionNote`** — nol `Property-Set`
   di 74 Activity terkait. Keduanya hanya dibaca sebagai gerbang. Jadi `spreading` **tidak** memanggil
   `acceptance`; arah dependensinya satu arah.
2. `[terverifikasi]` **Akibat pelanggaran ambang adalah tombol kirim mati, bukan eskalasi approver.**
   Kondisinya hanya muncul di `pyDisabledWhen` Section email — nol rujukan di `Flow`/`FlowAction`/`When`.
   Merancangnya sebagai aturan eskalasi akan salah.
3. `[terverifikasi]` **Scoring tidak menggerakkan alur sama sekali.** Keluarannya hanya dirujuk
   produsennya dan Section tampilan; satu-satunya pemakaian non-tampilan memeriksa *kelengkapan*
   (`FinalScore == ""`), bukan nilainya. Di target, scoring adalah layanan tampilan.
4. `[terverifikasi]` **Validasi sisa tanpa toleransi** — menuntut `Σ Share == 100` dan
   `Σ PremiumSpreaded == PremiNusantaraRe` **persis** (keduanya dibulatkan 4 desimal lebih dulu).
   Tidak ada distribusi galat pembulatan. Port wajib mereproduksi pembulatan-lalu-banding ini,
   termasuk `@divide(...,1,4)` yang berada **di dalam** loop akumulasi sehingga galatnya menumpuk.

⛔ `[pertanyaan terbuka]` Mesin B — Activity terbesar di NB (1,3 MB) — bergerbang parameter
`spread == "Spreading"`, sementara mayoritas pemanggil mengirim parameter bernama `InSpreading`.
Audit terverifikasi: parameter `spread` hanya dikirim **2** berkas
(`Harness\ViewPolis.xml`, `Section\InputInwardFacultativeDtl_IsUW.xml`), sedangkan `InSpreading`
dikirim 7 berkas. Apakah mesin B memang hanya hidup di dua layar itu, atau ini cacat? Jawabannya
menentukan apakah 1,3 MB logika perlu diport sama sekali.

### 4.1 `IsSpreadingDepan` — bukan rule hilang

`[terverifikasi]` Rule ini dirujuk `NB FacIn\Activity\SpreadingAdditionalProtection.xml` sebagai
pra-kondisi langkah, dan **berkasnya tidak ada di `NB FacIn\When\`**. Tetapi ia **ada di folder
Endorsment**, dan kondisinya terbaca penuh:

```
pyLogic = A OR B OR C OR D
  A: Rule IsPA     evaluates to true
  B: Rule IsTravel evaluates to true
  C: Rule IsMBU    evaluates to true
  D: Rule IsFire   evaluates to true
```

⚠️ **Caveat yang wajib dibawa:** salinan yang terbaca berversi `01-01-52` (commit 2025-05-13) dari
folder EDM. Karena salinan NB tidak ada, **tidak dapat dipastikan** NB memakai versi yang sama — dan
95 rule di korpus ini terbukti berbeda versi antar folder. Implementasikan kondisinya, tetapi catat
asalnya dari folder lain.

**Tidak perlu `panic()`** — dikukuhkan keputusan work owner K-003. Ini rule yang hilang dari satu
folder, bukan rule yang kondisinya tidak diketahui. Predikatnya **tetap diimplementasikan** sesuai
kondisi di atas; menghapusnya akan mengubah perilaku pra-kondisi langkah di
`SpreadingAdditionalProtection.xml`.

### 4.2 `IsOfferFacIn` — gerbang masuk NB

`[terverifikasi]` + keputusan work owner **K-002**:

```go
// Asal: NB FacIn/When/IsOfferFacIn.xml (pyConditionValue1String)
// K-002 (2026-09-15): ekspresi tersimpan yang berlaku.
// Teks tampilan rule menyebut kondisi lain (Kode Bisnis "02"/"58"/"SB"/"SG") —
// kandidat perbaikan, TIDAK diimplementasikan. Tersangka pertama bila paralel run
// memperlihatkan selisih pada gerbang masuk NB.
func IsOfferFacIn(c *Case) bool { return c.Quotation.BusinessFac == "F" }
```

---

## 5. `repository/oracle`

`[terverifikasi]` 62 activity menulis Oracle lewat 225 langkah RDB penulis; 57 dari 221 pernyataan
SQL menulis.

**Aturan [usulan] yang mengikat:**

1. Satu berkas Go per tabel; setiap fungsi menyebut rule `RDBList` asalnya (`CLAUDE.md` §4.6).
2. Arah operasi ditentukan dari **verba SQL**, bukan nama rule atau nama tag (§1.1).
3. Nilai uang masuk dan keluar sebagai `decimal.Decimal` + mata uang; kolom Oracle tetap `NUMBER`.
4. ⚠️ Parser masukan menerima **string berkoma desimal** — 1.057 titik konversi di korpus.
5. ⛔ Alias koneksi `ASM` (529 langkah) dan `RNM` (87 langkah) menunjuk instance/skema yang
   **belum diketahui**. Menentukan jumlah pool koneksi dan apakah transaksi lintas-alias mungkin.

---

## 6. Yang **tidak boleh** diimplementasikan sekarang

| Bagian | Penghalang |
| --- | --- |
| Enum `ProposalAcceptStatus` | Arti nilai tidak ada di korpus; baris DecisionTable tidak terekspor |
| Jalur tulis produksi | Isi 35 stored procedure `POOLDATA.*` tidak diketahui |
| Endpoint lookup | Nama tabel Oracle untuk 64 kelas `Int-*` tidak diketahui |
| **Enam cabang DITANGGUHKAN** (K-006) — Special Acceptance · tangga akseptasi putaran kedua · limit tambahan treaty type · kapasitas treaty · simpan produksi endorsement bonding · penanganan galat konversi produksi | ⏸ Menunggu verifikasi terhadap **ekspor produksi tunggal**. **Bukan kode mati** — "hilang dari ekspor" ≠ "usang", karena ekspor diambil dari titik waktu berbeda antar folder (R0). Bila activity-nya **ada** di ekspor produksi → **wajib diport**; bila benar tidak ada → baru dikeluarkan, dengan bukti. **Tidak boleh dihapus dari rancangan sebelum itu.** |
| Rumus premi | Satuan `.Rate` tidak konsisten — `1e9` vs `1e4`, selisih satu ordo besaran 10 |
| Mesin spreading B | Gerbang parameter `spread`/`InSpreading` belum dijelaskan |

**Yang dapat dimulai sekarang tanpa menunggu jawaban bisnis:** `pkg/money`, `internal/rules`
(registry predikat `When`), dan struktur repository untuk jalur **baca**.

---

## 7. Aturan penulisan kode

Berlaku penuh, dari `CLAUDE.md`:

- Uang tidak pernah `float` — `decimal.Decimal` + mata uang wajib (§4.1).
- `handlers → services → repository`, searah, tanpa memotong lapisan (§4.2).
- Setiap query, predikat, dan transisi menyebut rule Pega asalnya dalam komentar (§4.6).
- Endpoint & secret dari konfigurasi, tidak pernah literal (§4.4).
- Presisi pembulatan adalah **parameter per-rule**, bukan konstanta global. Jangan menyeragamkan.
- Kandidat perbaikan direproduksi apa adanya + diberi komentar; memperbaikinya adalah keputusan
  bisnis terpisah (§1).
