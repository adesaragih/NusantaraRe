# 08: Tahap Medical Check & Claim Analis + jalur balik

**Status:** claimed

**Blocked by:** 07 (penegakan peran) — tiap tahap milik peran tertentu

## Hasil & nilai pengguna

Sebagai **ReasLifeMedicalAdvisor**, saya menerima klaim yang menunggu telaah medis dan dapat
mencatat hasilnya; sebagai **ReasLifeSPV**, saya menerima klaim yang siap dianalisis. Keduanya dapat
**mengembalikan** kasus ke admin bila data kurang, dan SPV dapat mengembalikan ke medis bila telaah
perlu diulang — sehingga tidak ada keputusan yang diambil di atas data tidak lengkap.
*(User story 8, 11–14, 23, 24 di spec)*

## Area codebase

`internal/models` (posisi tahap pada klaim), `internal/services` (transisi antar tahap + jalur
balik), `internal/handlers` (endpoint submit dan kembalikan), `frontend/` (layar Medical Check dan
Claim Analis, kontrol kembalikan).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Flow/Register_Flow.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW` | `[terverifikasi]` empat tahap: `Assignment2` "Input Register", `Assignment1` "Outstanding Claim", `Assignment3` "Medical Check", `Assignment4` "Claim Analis" |
| `Claim Life/When/IsSendtoAdmin.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` / `RULE-OBJ-WHEN` | `[terverifikasi]` `pyWorkPage.SendtoAdmin = 1`; menggerbangi pengembalian dari **tiga titik** di `Register_Flow` |
| `Claim Life/When/IsSendtoMedical.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOMEDICAL` / `RULE-OBJ-WHEN` | pengembalian SPV → medis; `[keputusan work owner]` artinya — kondisinya TERBACA - lihat blok ralat |
| `Claim Life/Section/MedicalCheckClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `MEDICALCHECKCLAIMLIFE` / `RULE-OBJ-HTML-SECTION` | layar telaah medis; memuat gerbang `pyPosition` |
| FlowAction `MEDICALCHECK`, `AKSEPTASICLAIMLIFE` | — | `[terverifikasi]` tindakan pada kedua tahap |

## ADR terkait

**ADR-0002** (peran per tahap), **ADR-0011** (tahap 2 dan 3 **tidak mengubah** status baris —
baris tetap Outstanding sepanjang Medical Check dan sampai keputusan Komite).

## Acceptance criteria

- [ ] Klaim dapat berpindah Register → Outstanding → Medical Check → Claim Analis, masing-masing
      hanya oleh peran yang berhak.
- [x] Status baris adjustment **tetap Outstanding** sepanjang perpindahan tahap — perpindahan tahap
      bukan keputusan akseptasi.
- [x] `ReasLifeMedicalAdvisor` dapat mengembalikan kasus ke `ReasLifeAdmin`.
- [ ] `ReasLifeSPV` dapat mengembalikan kasus ke `ReasLifeAdmin`.
- [ ] `ReasLifeSPV` dapat mengembalikan kasus ke `ReasLifeMedicalAdvisor`.
- [ ] Kasus yang dikembalikan muncul kembali di antrean peran tujuan.
- [x] Pengembalian **tidak** mengubah status baris adjustment mana pun.

## Catatan

`[terverifikasi]` Pengembalian bukan jalur langka: `SendtoAdmin` menggerbangi pengembalian dari
**tiga titik** di `Register_Flow`. Perekaman pelakunya ditangani tiket **09** — di sistem lama kedua
penanda hanya menyimpan nilai `1`, tanpa pelaku dan tanpa waktu.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Pembacaan ulang XML — 26 September 2026 malam

| Rule | Byte | Yang diambil |
| --- | ---: | --- |
| `Flow/Register_Flow.xml` | 264.216 | `pyPosition` per shape — peran tiap tahap |
| `When/IsSendtoAdmin.xml` | 15.776 | `pyConditionString` = `pyWorkPage.SendtoAdmin = 1` |
| `When/IsSendtoMedical.xml` | 18.349 | `pyConditionString` = `pyWorkPage.SendtoMedical = 1` |
| `Section/Diagnose_Section.xml` | 128.697 | bukti butir **al** |

### ⭐ `IsSendtoMedical` `[ditutup oleh XML — 26-09-2026]`

Tiket ini menandainya *"tidak terbaca dari tag"* dan `[keputusan work owner]`. **Ia terbaca.**
Kondisinya memang tidak ada di label rule — label-nya masih cetakan kosong
*"[first value][relation][second value]"* — tetapi `pyConditionString` berbunyi
**`pyWorkPage.SendtoMedical = 1`** dan `pyConditionValue1` berbunyi
`compareTwoValues(pyWorkPage.SendtoMedical, "=", 1)`. Simetris persis dengan `IsSendtoAdmin`.
Pelajarannya: label rule bukan tempat kondisi tinggal.

### ⚠️ `Assignment4` (Claim Analis) — pemegangnya TIDAK ada di Flow

Keempat `Assignment` ada sebagai shape, tetapi **nol** di antara empat belas `pyPosition` di
`Register_Flow.xml` yang menyebut `Assignment4`. Peran pemegangnya `[terbuka — work owner]` dan
**tidak ditebak**: `models.PeranTahap` sengaja tidak memuat entrinya, dan memindahkan kasus dari
tahap itu menghasilkan galat terang.

### ⭐ Bukti butir **al** — ditulis apa pun keadaan PAKET

Kelas `ASM-FW-GISFW-Data-DiagnoseLife` dipakai **enam** rule: `SearchDiagnose_act`, `SetDisease`,
`SetSTS_Reject`, `Diagnose_Harness`, `Diagnose_Section`, `ClaimLifeDetailGCNM`. Layar diagnosa
**ada**, dan `Diagnose_Section` mengisi **dua** kolom: **`.ICD_Code`** dan **`.Disease`**. Itulah
bentuk tabel `T_CLAIMLF_DIAGNOSE` bila work owner menyetujui **al** — bukti dulu, keputusan kemudian.

## Implementasi — 26 September 2026 malam (tiket 08)

**Status: `claimed`** — **3 dari 7 AC tertutup** *(6 sebelum `/code-review`)*. Titik tetap `9b48e32`.
Verifikasi: vet · vet db · gofmt nol · build · **186 PASS · 0 FAIL** *(dari 182)* · 28 SKIP ·
`tsc` · 7 JS · 88 modul.

| Berkas | Isi |
| --- | --- |
| `models/tahap.go` *(baru)* | `Tahap` tertutup, `Kode()`/`TahapDariKode` pulang-pergi, `PeranTahap` |
| `services/tahap.go` *(baru)* | `JalurBalik`, `WajibPeranTahap`, `TahapLayanan.Pindah` |
| `repository/klaimlife.go` | `PerbaruiTahap`, `TahapKlaim` — kolomnya sudah ada di `001`, nol migrasi |

⛔ **Satu cacat desain yang hampir lolos.** Ronde pertama memeriksa peran tahap **tujuan**. Itu
membalik seluruh jalur balik: pengembalian ke Admin oleh Medical Advisor akan menuntut pelakunya
berperan **Admin** — yang justru bukan dia. Orang memindahkan pekerjaan yang **sedang ia pegang**,
jadi yang diperiksa adalah peran tahap **asal**, dibaca dari `PY_POSITION`.

⛔ **Nol tulisan status.** `tahap.go` menyebut `KodeStatus` hanya di dua komentar; penjaga statik
tiket 07 memeriksanya, dan berkas ini sengaja tidak masuk daftar penulis status.

### ⛔ AC yang belum tertutup — 1

*"Klaim dapat berpindah Register → Outstanding → Medical Check → Claim Analis"* — tiga tahap pertama
bekerja; **Claim Analis** tidak dapat dimasuki maupun ditinggalkan, sebab pemegangnya tidak ada di
Flow. `[terbuka — work owner]`

### Ralat menurut XML — 27 September 2026

Bukti: berkas pecahan `Flow/Register_Flow.xml` **114.144 byte** *(bukan 264.216 seperti tertulis di
bab Pembacaan — angka itu keliru dan diralat di sini)*.

| # | Teks lama | Teks baru | Bukti |
| ---: | --- | --- | --- |
| 1 | `PY_POSITION` berisi pengenal shape `"Assignment<n>"` | ⛔ ia berisi **NAMA PERAN** | pecahan 582, 605, 628, 668, 731 — seluruhnya `pyWorkPage.pyPosition` ← `"ReasLifeAdmin"` / `"ReasLifeMedicalAdvisor"` / `"ReasLifeSPV"`; nol baris menyetelnya ke `"Assignment<n>"`; 11 gerbang di korpus membandingkannya dengan nama peran |
| 2 | *"`IsSendtoMedical` kondisinya tidak terbaca dari tag"* `[keputusan work owner]` | `[ditutup oleh XML]` — `pyWorkPage.SendtoMedical = 1` | `When/IsSendtoMedical.xml` pecahan 165 (`pyConditionString`) dan 318 (`compareTwoValues`); label 317 memang cetakan kosong |
| 3 | Claim Analis `[terbuka]` *("nol `pyPosition` menyebut `Assignment4`")* | **`ReasLifeSPV`** | ADR-U-0002 tabel pemetaan tahap; ADR itu sendiri menjelaskan diamnya korpus — *"model peran DIRANCANG, bukan dimigrasikan"* |
| 4 | bukti butir **al** dari `MedicalCheckClaimLife.xml` *(penunjuk brief)* | dari `Section/Diagnose_Section.xml` | `MedicalCheckClaimLife.xml` (745.295 byte) memuat **nol** `DiagnoseLife`, `ICD_Code`, maupun `.Disease`; ia menampilkan peserta (`.ClientName`, `.Note`, `.SourceOfBusiness`). Daftar diagnosanya ada di `Diagnose_Section` pecahan 1657–1658 dan 2973–2974 |
| 5 | bentuk tabel **al** = dua kolom | **empat**: `.ICD_Code`, `.Disease`, `STS_REJECT`, dan kunci peserta | `SetSTS_Reject` pecahan 241/250/257–258 menulis `.STS_REJECT` pada `.DiagnoseList` kelas `Data-DiagnoseLife`, induknya `Int-LIFE_PREMIUM_DETAIL` — jadi per peserta |

### Hasil `/code-review` atas titik tetap `9b48e32`

⛔ **Cacat terbesar tiket ini milik saya, dan ia tidak akan terlihat tanpa Oracle.** Saya membalik
arti `PY_POSITION`. Akibatnya setiap pembacaan baris nyata berakhir *"tidak dikenal"* dan
perpindahan tahap **selalu gagal** — sedangkan seluruh test hijau, sebab tak satu pun menyentuh
kolomnya. Model ditulis ulang: yang tersimpan adalah **peran pemegang**, dan perpindahan adalah
**serah terima peran**.

| # | Temuan | Tindakan |
| ---: | --- | --- |
| 1 | `PY_POSITION` dibalik artinya | ✅ model ditulis ulang; `TestKolomPyPositionMenyimpanNamaPeran` mengunci arahnya, termasuk bahwa `"Assignment1"` **harus** tetap tak dikenal |
| 2 | Claim Analis ditandai `[terbuka]` padahal ADR-U-0002 memutuskannya | ✅ `ReasLifeSPV` dipasang; saya memperlakukan diamnya korpus sebagai pertanyaan, padahal ia celah yang ADR itu tutup |
| 3 | **Nol aturan kesahan perpindahan** — Input Register langsung ke Claim Analis lolos | ✅ `SerahTerimaSah` dari penyambung Flow + ADR-U-0002 |
| 4 | `JalurBalik` sebagai bendera BEBAS, dan test saya bahkan menegaskan keduanya menyala — keadaan yang mustahil | ✅ diturunkan dari pasangan perannya; keadaan mustahil itu kini tidak dapat dibentuk |
| 5 | TOCTOU: tahap asal dibaca di luar transaksi tanpa syarat `WHERE` | ✅ `AND PY_POSITION = :asal` |
| 6 | ADR-U-0007 — jalur balik disebut namanya di ADR, dan `Pindah` merekam nol | ✅ `Jejak` + `JejakBelumDiputuskan`, direkam di dalam transaksi |
| 7 | `PeranTahap` peta ekspor yang dapat ditulis pemanggil | ✅ `PeranPemegangTahap` |
| 8 | Blok `Exec`+`RowsAffected` lima salinan | ✅ `pastikanSatuBaris` |
| 9 | Tiga AC tercentang tanpa terkirim *(dua jalur balik SPV, satu antrean)* | ✅ **dibatalkan centangnya** — 6/7 menjadi **3/7** |

⛔ Tersisa dan dinyatakan: antrean per peran *(nol query, nol pintu, nol layar)*, dua layar tahap,
dan test `db` untuk perpindahan tahap. Ditambah `[terbuka — work owner]` baru: `PY_POSITION` tidak
dapat membedakan Input Register dari Outstanding, sebab keduanya dipegang peran yang sama.

**Verifikasi sesudah perbaikan:** vet · vet db · gofmt nol · build · **187 PASS · 0 FAIL** ·
28 SKIP · `tsc` · 7 JS · 88 modul.

---

### Ralat menurut XML — 27 September 2026 (audit A0, brief lanjutan 4 bab 7)

⛔ **Akseptasi punya DUA jalur; tiket ini hanya mengenal satu.**

| Butir | Teks lama | Teks baru | Bukti |
| --- | --- | --- | --- |
| siapa menulis status Aksep | *"Aksep ditulis modul Komite, bukan Claim Life"* — `ErrAksepBukanDariModulIni` menolak tujuan Aksep bagi siapa pun | **Claim Life mengaksep sendiri** lewat `SaveAdjustment_Act`; sentinelnya **DIHAPUS** | `Claim Life/Activity/SaveAdjustment_Act.xml` pecahan baris **1833** (`ACCEPTEDNO`), **1879** (`STS_REJECT = 1`), **1899** (`ACCEPTATION_DATE = @CurrentDateTime()`) |
| gerbang peran jalur itu | — *(tidak dikenal)* | **pemegang TAHAP**, bukan daftar peran datar | `[terverifikasi — pohon XML]` tombol "Save Adjustment" (`Section/ClaimLifeDetailGCNM.xml` **22641 → 22665**) tidak dibungkus gerbang peran mana pun; seluruh leluhurnya ALWAYS. `pyCondition 1=2` yang tampak bertetangga adalah `pyContainerVisibleWhen` **layout lain** |
| prasyarat dagangnya | — | peserta `.IsCheck=true`, `.ACCEPTEDNO==""`, baris `.STS_REJECT=="0"`, dan `Type` | pecahan **854** (QP/QR) dan **1048** (TP/TR) |

⚠️ **DUA cacat rule Pega yang sengaja TIDAK ditiru — dilaporkan ke work owner:**

1. **Prasyarat tanpa kurung.** Baris **837** dan **1031** berbunyi
   `.IsCheck=true && Type=="QP" || Type=="QR"`. Karena `&&` mengikat lebih erat daripada `||`,
   bacaan harfiahnya meloloskan `QR` dan `TR` **tanpa** memeriksa `IsCheck` sama sekali. Go memakai
   bacaan yang **dimaksud** — `IsCheck && (QP||QR)`.
2. **Tahun tidak bergeser.** Baris **605** berbunyi
   `@if(MM=="12" && NextMonth=="01", @toDecimal(@CurrentDate("YY")), @toDecimal(@CurrentDate("YY")))`
   — **kedua cabangnya identik**, sehingga nomor Januari memakai tahun Desember. Go menggeser
   tahunnya.

⭐ **Temuan yang menghentikan jalur ini di satu tempat:** nomor akseptasi memuat **kode bisnis**
(`'RNML-A'||{pyWorkPage.BusinessCode}||…`), tetapi model relasional kita **tidak menyimpannya** —
`T_GENERAL_CLAIM` hanya punya `BUSINESS_NAME`, dan `BUSINESSID` bukan salah satu dari 18 kolom datar
warisan yang `Simpan` tulis. Jalurnya **gagal terang** dengan `ErrKodeBisnisBelumTersimpan` (HTTP 501)
sampai kolomnya lahir di **A1**. Tidak dikarang.

## Ralat menurut XML — 27 September 2026 (butir at)

**Yang diralat:** `[terbuka — work owner]` "pembedaan Input Register vs Outstanding Claim tidak
tersimpan di kolom mana pun". **Ditutup oleh XML**, bukan oleh keputusan manusia.

**Bukti, path + baris** *(korpus READ-ONLY)*:

| Bukti | Isi |
| --- | --- |
| `Flow/Register_Flow.xml` 358 · 343 · 268 · 313 | `<pyTaskName>` keempat assignment: `Input Register`, `Outstanding Claim`, `Medical Check`, `Claim Analis`. Keempatnya **keadaan** kasus |
| `Section/InputOSClaimLife.xml` **21404** | `<pyLabel>Send Back to Register</pyLabel>` |
| `Section/InputOSClaimLife.xml` **21433** | `<pyLocalAction>SendtoAdmin</pyLocalAction>` |
| `Section/InputOSClaimLife.xml` 21349 · 21839 · 21863 | `Send to Medical Check` → `SendtoAdmin_Act1` |

`Input Register` **dapat dituju kembali**. Ia bukan tahap yang hanya dilewati sekali, melainkan
keadaan yang tombol di layar Outstanding kembalikan kasus kepadanya. Tanpa kolom yang menyatakannya,
kotak masuk Admin menyatukan dua antrian yang di Pega terpisah.

**Keputusan at** `[DIPUTUSKAN 27-09-2026; work owner dapat memveto sebelum migrasi dijalankan di
Oracle mana pun]`: `T_WORK_CLAIM.TAHAP VARCHAR2(32)` menyimpan nama assignment VERBATIM `pyTaskName`.
`PY_POSITION` **tetap** — ia peran pemegangnya (ADR-U-0002); kedua kolom menjawab dua pertanyaan
berbeda.

**Akibatnya pada kode:**

| Sebelum | Sesudah |
| --- | --- |
| `serahTerimaSah` peta **PERAN** | peta **TAHAP** — peta peran tidak dapat menyatakan Admin→Admin |
| `JalurBalikPeran(dari, ke string)` | `JalurBalikTahap(dari, ke Tahap)`; Register ⇄ Outstanding **bukan** jalur balik (tidak ada yang dikembalikan kepada siapa pun) |
| `TahapKlaim` membaca `PY_POSITION` | `TahapDanPeran` membaca **keduanya dalam satu kueri** |
| penjaga optimis `UPDATE … AND PY_POSITION = :n` | `AND NVL(TAHAP, :n) = :n` — pada Register ⇄ Outstanding `PY_POSITION` **tidak berubah**, sehingga penjaga lama tidak dapat mendeteksi kasus yang sudah dipindah orang lain |
| jejak mencatat peran asal → peran tujuan | mencatat **tahap** asal → tujuan; jejak berperan akan berbunyi "dari Admin ke Admin" |

**AC yang bergeser:** tidak ada AC tiket ini yang berubah centangnya — yang ditutup adalah `[terbuka]`
di badan tiket, bukan sebuah AC.

⚠️ **Masih terbuka, dan disebut namanya:** prasyarat `pyWorkPage.pyPosition=="ReasLifeMedicalAdvisor"`
pada `SendtoAdmin_Act` **dan** `SendtoAdmin_Act1` belum dijelaskan — dengan posisi Admin langkahnya
dilewati (`WhenFalse=3`), lalu bagaimana `Send Back to Register` bekerja dari layar Outstanding yang
dipegang Admin? Pertanyaan **pohon**, dijawab di A3 kelompok Outstanding dengan membaca
`Register_Flow.xml` 583–772 beserta urutan shape-nya. Bukan alasan menunda at.

## Ralat menurut XML — 27 September 2026 (butir aw)

**Cacat rule warisan, ditiru MAKSUDnya bukan hurufnya.**

Dua tombol di layar Outstanding memanggil perpindahan tahap:

| Tombol | Baris | Jalur |
| --- | ---: | --- |
| `Send Back to Register` | `InputOSClaimLife.xml` **21404** | → `<pyLocalAction>SendtoAdmin` **21433** → `SendtoAdmin_Act` |
| `Send to Medical Check` | **21349** / **21839** | → `SendtoAdmin_Act1` **21863** |

**Yang XML sebenarnya lakukan:** kedua activity berprasyarat
`pyWorkPage.pyPosition=="ReasLifeMedicalAdvisor"` dengan `WhenTrue=2` dan `WhenFalse=3` *(lewati)* —
`SendtoAdmin_Act` baris **338**. Layar Outstanding dipegang **`ReasLifeAdmin`**. Jadi **pada posisi
Admin kedua tombol itu tidak menulis apa pun**: `SendtoAdmin` tidak pernah menjadi `"1"`, `Decision3`
tidak pernah mengirim kasus kembali ke `Assignment2`, dan kasus **tidak pernah dapat kembali ke Input
Register**.

Penulis `SendtoAdmin`/`SendtoMedical` di seluruh modul hanya **tiga** activity — `SendtoAdmin_Act`,
`SendtoAdmin_Act1`, `SendtoMedical_Act` — dan tidak satu pun berjalan pada posisi Admin.

**Mengapa ini dinilai CACAT, bukan maksud bisnis:** tiga hal menyebut jalur balik ini sebagai fitur —
label tombolnya sendiri, penyambung `Decision3 → Assignment2` yang ada di alurnya, dan ADR-U-0002.
Prasyaratnya tampak terbalik atau salah tempel.

**Keputusan aw** `[DIPUTUSKAN — 27 September 2026, dari maksud XML yang terang; veto work owner
terbuka sampai kode ini dijalankan di Oracle mana pun]`:

| Tombol | Perbuatan kita |
| --- | --- |
| `Send Back to Register` | `Pindah` ke **Input Register**, jejak audit direkam |
| `Send to Medical Check` | `Pindah` ke **Medical Check** *(maksud `SendtoAdmin_Act1` = `"0"`)* |

Keduanya bergerbang **pemegang tahap Outstanding** *(services; 403 bila bukan)*, dan perpindahan yang
tidak ada di tangga dijawab **409** — bukan 400: permintaannya berbentuk benar, keadaan kasusnya yang
tidak mengizinkan.

Rute: `POST /api/klaim-life/{id}/tahap/{tujuan}`, tujuan berupa **kata** *(`input-register`,
`outstanding`, `medical-check`, `claim-analis`)* — bukan angka: jalur `/tahap/2` tidak terbaca siapa
pun, dan angka yang bergeser bila urutan `models.Tahap` berubah memindahkan kasus ke tempat yang
salah tanpa satu pun galat.

**Cacatnya dilaporkan** ke `OQ-untuk-tim.md` **OQ-C** dengan barisnya, untuk dikonfirmasi pengembang
Pega.

**AC:** tidak ada AC tiket ini yang berubah centangnya — jalur baliknya kini punya rute dan kontrol,
tetapi pembuktian perilakunya menuntut Oracle.
