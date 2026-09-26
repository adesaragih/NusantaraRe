# 06: Validasi Date of Loss per `Type` + `ContentNote` dari `BusinessCode`

**Status:** claimed

**Blocked by:** 03 (baris `AdjustmentList` + Save ke Outstanding)

## Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, sistem menolak *Date of Loss* yang berada di luar jendela valuasi polis —
sehingga klaim yang tidak tertanggung tidak pernah masuk ke siklus — dan jenis klaim terisi otomatis
dari kode produk sehingga saya tidak salah memilih. *(User story 4 dan 5 di spec)*

## Area codebase

`internal/services` (aturan validasi + penurunan jenis klaim), `internal/models` (data acuan
`BusinessCode` → `ContentNote`), `internal/handlers` (pesan kesalahan yang dapat dibaca),
`frontend/` (penampilan pesan pada formulir).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/ValidasiDOL_Act.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `VALIDASIDOL_ACT` / `RULE-OBJ-ACTIVITY`, 59.747 byte | dua cabang jendela valuasi per `Type` |
| `Claim Life/Activity/SaveOutStandingLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | penurunan `ContentNote` dari `BusinessCode` |

`[terverifikasi]` Dua cabang validasi DOL:

| Cabang | Jendela | Argumen `@addCalendar` |
| --- | --- | --- |
| `Type=="QR" \|\| Type=="QP"` | `GROSS_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `(.DATE_OF_LOSS,0,0,0,0,0,0,0)` |
| `Type=="TR" \|\| Type=="TP"` | `RETROCESSION_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `(.DATE_OF_LOSS,0,0,0,1,0,0,0)` |

Gagal → `local.errmsg = "Invalid DOL"`, dengan `local.Begin==false || local.Expired==true`.

## ADR terkait

**ADR-0012** (`Type` menggerbangi dua hal: wewenang **dan** jendela validasi; `TP` = Payable,
`TR` = Receivable — `[keputusan work owner]`, tidak ada di korpus).

## Acceptance criteria

- [x] Klaim ber-`Type` `QP`/`QR` dengan *Date of Loss* di luar `GROSS_VALUATION_BEGIN_DATE` …
      `_EXPIRED_DATE` ditolak dengan pesan yang setara `"Invalid DOL"`. *(AC 13 spec)*
- [x] Klaim ber-`Type` `TP`/`TR` diuji terhadap `RETROCESSION_VALUATION_*`, dengan pergeseran
      tanggal yang sama seperti Pega. *(AC 14 spec)*
- [ ] `ContentNote` terisi sesuai tabel `BusinessCode` `L1`–`L21` di `CONTEXT.md`. *(AC 15 spec)*
- [x] Pemetaan `BusinessCode` → `ContentNote` diwujudkan sebagai **data acuan**, bukan rangkaian
      `if` bercabang.
- [x] Validasi dan penurunan jenis klaim membaca **satu** field `Type` yang sama — bukan dua salinan
      seperti di Pega.

## Catatan penutupan (2026-09-14)

**DOL dibawa sebagai paritas** `[keputusan work owner]` — validasi *Date of Loss* per `Type`
dipertahankan persis seperti perilaku Pega. Tercatat di `CONTEXT.md` (Lampiran, butir 2).

`[terbuka — catatan kecil, NON-PEMBLOKIR]` Satuan pergeseran tanggal pada cabang `TP`/`TR`
(`@addCalendar(.DATE_OF_LOSS,0,0,0,1,0,0,0)`) belum dikonfirmasi **Product+UW**; definisi
`@addCalendar` tidak ada di korpus. **Ini bukan blocker**: pergeseran direplikasi apa adanya dari
argumen yang terbaca. Bila kelak Product+UW menyatakan satuannya berbeda, itu perubahan satu baris
disertai testnya.

`[terbuka]` **OQ-020** — arti `QP`/`QR` belum dijawab. **Tidak memblokir**: perilaku kedua cabang
sudah terbaca penuh; hanya namanya yang belum.

`[data DBA]` Kolom terkait di `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`: `TYPE VARCHAR2(10)`.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Pembacaan ulang XML — 26 September 2026 malam (aturan brief modul §1.3)

Berkas dipecah dengan cara brief §2; nomor baris di bawah dari berkas pecahan.
`ValidasiDOL_Act.xml` 59.747 byte → 1.408 baris; `LoadDataPeserta_Act.xml` 80.424 byte → 1.916 baris.
`SaveOutStandingLife_Act.xml` dibaca ulang dari pecahan tiket 03.

| Rule | Tag yang dibaca | Yang diambil |
| --- | --- | --- |
| `Activity/ValidasiDOL_Act.xml` | langkah 2, 3, 4 beserta **aksi** precondition | dua cabang jendela, pergeseran, gerbang galat, properti sasaran pesan |
| `Activity/SaveOutStandingLife_Act.xml` | precondition baris 3261; jalur `ContentNote` 3446 · 3679 · 3823 · 4054 · 4199 | penguatan `L1`–`L11` → `DEATH` |
| `Section/EditDateClaimLife_Section.xml`, `Section/ClaimLifeDetailGCNM.xml` | `pyActivity` | **siapa** memanggil validasi ini |

### ⭐ Yang bertambah dari yang tiket tulis

**1. Validasi ini milik PESERTA, dan berjalan saat tanggal diubah di layar — bukan saat simpan.**
`[terverifikasi]` Kelas seluruh activity dan langkah 4-nya `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL`,
dan langkah 4 adalah `Property-Set-Messages` dengan `Field` = **`.DATE_OF_LOSS`**,
`Category` = `pyMessageLabel`, `Message` = `local.errmsg`. Pemanggilnya **dua Section**
(`EditDateClaimLife_Section` dan `ClaimLifeDetailGCNM`, masing-masing dua tempat) lewat
`<pyActivity>ValidasiDOL_Act</pyActivity>` — bukan Flow, bukan `SaveOutStandingLife_Act`.

⚠️ `Activity/ValidasiSTNC_Act.xml` **bukan** pemanggil: `pzOriginalInstanceKey`-nya berbunyi
`… VALIDASIDOL_ACT #20231213T100309.692 GMT`, jadi ia **salinan** rule ini yang kemudian disunting.
Mudah terbaca sebagai pemanggil oleh pencarian teks biasa.

**2. Kolomnya sudah ada, dan ia kolom PESERTA.** `DATE_OF_LOSS DATE` ada di migrasi `003`
(`T_CLAIMLF_PREMIUMLIST_DETAIL`) sejak tiket 14, tetapi **tidak pernah** masuk daftar kolom tulis/baca
maupun `models.Peserta`. Ditambahkan sebagai `TanggalKejadian`. Nol kolom baru, nol langkah migrasi.

**3. Nama properti Pega ≠ nama kolom kita.** `RETROCESSION_VALUATION_BEGIN_DATE` /
`_EXPIRED_DATE` di XML adalah `RETRO_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` di `003`
`[terverifikasi katalog]`. Dicatat supaya tidak lahir kolom ketiga.

**4. `ContentNote`: nol selisih antara XML dan `CONTEXT.md`, jadi nol ralat.** `[terverifikasi]`
Satu precondition (baris 3261) menguji `BusinessCode` `"L1"` sampai `"L11"` berderet, dan jalur yang
dijaganya adalah jalur `ContentNote = "DEATH"` (3446, 3679, 3823; kebalikannya `!= "DEATH"` di 4054
dan 4199). Kesebelasnya memang `DEATH` di tabel work owner. `L12`–`L21` **tidak diuji modul klaim
sama sekali** — keduanya tetap ditulis, sebab kode yang belum pernah lewat modul ini tidak boleh
jatuh ke galat hanya karena modul klaim belum pernah melihatnya.

⚠️ Di Pega `ContentNote` **tidak disimpan**: ia dibaca dari `Business.pxResults(1)` — hasil
pencarian saat itu juga. Tidak ada kolom `CONTENT_NOTE` di skema kita, dan menurut XML memang tidak
perlu ada.

### ⛔ Dua asumsi pustaka Pega — `[dugaan]`, bukan fakta korpus

Korpus **tidak** memuat definisi `@addCalendar` maupun `@CompareDates`. Keduanya menentukan seluruh
tabel batas, jadi keduanya ditulis sebagai **konstanta bernama** di `services/dol.go`:

| Konstanta | Asumsi | Alasannya |
| --- | --- | --- |
| `pergeseranTPTR = time.Hour` | argumen keempat `@addCalendar` adalah **jam** | cabang TP/TR memanggil `@addCalendar(.DATE_OF_LOSS,0,0,0,1,0,0,0)` — **tujuh** argumen angka sesudah tanggalnya, yang cocok dengan tanda tangan `(tanggal, tahun, bulan, hari, JAM, menit, detik, milidetik)`. Posisi keempat = jam |
| `bandingKetat = true` | `@CompareDates(a,b)` benar bila `a` **sesudah** `b`, ketat | tidak terbaca di korpus; ini bacaan lazim pustaka Pega |

Akibat keduanya, jendelanya **asimetris** — dan asimetri itu berasal dari argumen yang **berbeda**
di kedua cabang, bukan dari kami:

```
QP/QR : (BEGIN, EXPIRED]   tanggal mulai DITOLAK, tanggal berakhir diterima
TP/TR : [BEGIN, EXPIRED)   tanggal mulai diterima, tanggal berakhir DITOLAK
```

**Kasus mana yang berbalik bila Product+UW memutuskan lain** *(dituntut brief §4-06)*: bila satuannya
**hari**, `TP/TR` pada *satu jam sebelum berakhir* berbalik menjadi **tidak sah** (`DOL + 1 hari`
melewati `EXPIRED`), sedangkan `TP/TR` tepat pada tanggal mulai tetap sah. Bila perbandingannya
**tidak ketat**, `QP/QR` tepat pada tanggal mulai berbalik menjadi **sah**. Keduanya perubahan satu
baris disertai testnya. `[terbuka — Product+UW]`

### Pertanyaan yang XML tidak jawab

| # | Pertanyaan | Pemilik |
| ---: | --- | --- |
| 1 | Satuan pergeseran dan keketatan perbandingan *(di atas)* | Product + UW |
| 2 | Arti `QP`/`QR` — **OQ-020** | work owner |
| 3 | Siapa yang mengisi `DATE_OF_LOSS` dan kapan; Section-nya ada, layar kita belum | work owner |

---

## Implementasi — 26 September 2026 malam (tiket 06)

**Status: `claimed`** — **4 dari 5 AC tertutup.** Titik tetap `2a99841`.

Verifikasi penuh sesudah perbaikan review ada di bab hasil review di bawah, ditulis dari keluaran perintah.

### Keputusan yang dipakai dan yang ditunggu

Tiket ini **tidak memerlukan satu pun** keputusan yang masih `[USULAN]` — sesuai brief lanjutan §1.
Yang ditunggu hanya `[terbuka — Product+UW]` di atas, dan keduanya **tidak memblokir**: pergeseran
direplikasi apa adanya dari argumen yang terbaca, persis seperti bab penutupan tiket ini minta.

### Yang dibangun

| Berkas | Isi |
| --- | --- |
| `services/dol.go` *(baru)* | `ValidasiDOL(tipe, dol, peserta)` murni; `ContentNoteDari(kodeBisnis)`; konstanta `TypeQR/QP/TR/TP`, `pergeseranTPTR`, `bandingKetat`; empat galat yang dapat dikenali `errors.Is` |
| `models/businesscode.go` *(baru)* | tabel data acuan 21 baris + `ContentNoteUntuk` + `BusinessCodeContentNote()` yang mengembalikan **salinan** |
| `models/klaimlife.go`, `repository/kolompeserta.go` | `TanggalKejadian` ⇄ kolom `DATE_OF_LOSS` yang sudah ada di `003` |
| `models/satutype_test.go` *(baru)* | penjaga statik AC 5 |

### Penjaga baru, masing-masing DIBUKTIKAN dapat gagal

| Cacat yang dipasang | Yang menangkap |
| --- | --- |
| pergeseran TP/TR dihapus | 3 kasus batas TP/TR |
| `bandingKetat` dijadikan longgar | `TestBatasJendelaTiapType`, `TestPesanDOLPersisSepertiXML` |
| cabang QR/QP memakai jendela retro | 3 test |
| `Type` asing masuk cabang QR/QP | `TestTypeTidakDikenalDitolak` |
| valuasi kosong dianggap waktu nol | `TestValuasiKosongDitolakBukanDianggapNol` |
| medan `Type` kedua ditambahkan | `TestTypeKlaimHanyaSatuMedan` |

Pola mutasi ditulis **tanpa backslash** dan skripnya menggagalkan diri bila polanya tidak ditemukan —
pelajaran dari tiket 03, di mana satu pembuktian pernah hijau justru karena cacatnya tak pernah
terpasang.

### ⛔ Penyimpangan sadar dari rule Pega, satu buah

`Type` di luar `QR`/`QP`/`TR`/`TP` menjadi **galat**. Di Pega kedua precondition tidak terpenuhi,
activity-nya diam, dan **tanggal apa pun lolos**. Diam seperti itu adalah lubang, bukan aturan:
validasi yang tidak berjalan tidak dapat dibedakan dari validasi yang lulus. Dinyatakan, bukan
didiamkan.

### ⛔ Yang BELUM dikerjakan — 1 AC

| AC | Sebab |
| --- | --- |
| *"`ContentNote` **terisi** sesuai tabel `L1`–`L21`"* | nilainya **benar untuk seluruh 21 kode** dan teruji satu per satu, tetapi **belum ada satu pun tempat yang menampilkan atau menyimpannya**. Kata *terisi* menuntut konsumen, dan layar yang menampilkannya belum ada. Dicentang hanya bila jujur — pelajaran dari tiket 03, di mana empat AC bertuliskan *"tersimpan"* sempat tercentang tanpa satu pun pemanggil |

⚠️ **Ralat, sesudah `/code-review`:** ronde pertama bab ini memakai ukuran itu pada AC 15 saja,
sedangkan AC 13 dan 14 tercentang walau `ValidasiDOL` pun nol pemanggil produksi. Ukurannya kini
dipakai sama rata — dan keduanya ditutup dengan **memasang jalurnya**, bukan dengan membatalkan
centangnya.

⚠️ Tidak ada kolom `CONTENT_NOTE`, dan menurut XML memang tidak perlu ada — Pega pun menurunkannya
saat itu juga. Jadi yang kurang adalah **layarnya**, bukan kolomnya.

### Hasil `/code-review` atas titik tetap `2a99841`

Dua sub-agen paralel. **Sembilan temuan diterima**, nol ditolak. Yang terberat menyentuh
konsistensi saya sendiri.

| # | Sumbu | Temuan | Tindakan |
| ---: | --- | --- | --- |
| 1 | Spec | **AC 13 dan 14 tercentang padahal `ValidasiDOL` nol pemanggil produksi** — persis ukuran yang saya pakai sendiri untuk menolak mencentang AC 15 *("kata terisi menuntut konsumen")*, tidak saya pakai pada dua AC di atasnya | ✅ **Dikabelkan, bukan dibatalkan centangnya.** `KlaimLife.TypeKlaim` + `PerbaruiTanggalKejadian`, `services.TanggalKejadian().Set`, dan pintu `PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-kejadian` yang menjawab **422** berisi kalimat `"Invalid DOL"` beserta `pesertaId` di medannya sendiri |
| 2 | Standards | **ADR-U-0022 Akibat 1 dilanggar**: *"konversi adalah tanggung jawab pemuat di `repository`, bukan tersebar di lapisan layanan"*. `utils.ParseTanggal` di `services` adalah satu-satunya di luar `repository` | ✅ Pengurai pindah ke `models.Peserta.JendelaValuasi` — pada tipe yang memiliki teksnya, satu tempat, tidak tersebar. ⚠️ Perbaikan **sebenarnya** adalah pemuat mengurai sekali, dan itu menuntut kesembilan medan tanggal peserta berhenti bertipe teks: keputusan sekali untuk seluruh model, `[terbuka — tiket 14]` |
| 3 | Standards | **Nama kolom ketiga dikarang**: galat menyebut `VALUASI_GROSS_MULAI`, padahal kolomnya `GROSS_VALUATION_BEGIN_DATE` — justru yang bab pembacaan XML ini peringatkan sendiri | ✅ Nama kolom **sebenarnya** dari migrasi `003` dipakai, dan test mengunci namanya |
| 4 | Standards | **DOL kosong tidak dijaga**: waktu nol = tahun 1 Masehi, pasti di luar jendela, sehingga pengguna mendapat *"Invalid DOL"* untuk tanggal yang **tidak pernah ia isi** | ✅ `ErrDOLKosong` terpisah. Persis jebakan yang komentar ⛔ saya jaga untuk kolom valuasi tetapi lupa untuk DOL-nya sendiri |
| 5 | Standards | **`satutype_test.go` mendekati tautologi**: ia memindai `internal/models` saja sambil mengaku `Type` hidup *"HANYA di WorkClaim"*; `Type string` juga ada di `services/pendaftaran.go` dan `handlers/register.go` | ✅ Penjaga ditulis ulang: menelusuri **seluruh** `internal/`, dengan daftar tempat yang sah beserta **sebabnya** *(muatan permintaan, hidup satu permintaan lalu hilang)*. Pengakuan yang lebih luas daripada yang diperiksa adalah penjaga yang menenangkan tanpa menjaga |
| 6 | Spec | **Delapan dari enam belas sel batas hilang** — brief menuntut keempat batas untuk **tiap** `Type` | ✅ Tabelnya kini lengkap 16 sel, ditambah empat kasus silang jendela |
| 7 | Spec | **Pengenal peserta tidak ada di medan tersendiri** seperti brief minta | ✅ `GalatDOL{PesertaID}` dengan `Error()` mengembalikan kalimat XML apa adanya dan `Unwrap()` menjaga `errors.Is` |
| 8 | Spec | **Daftar "kasus mana yang berbalik" tidak lengkap** — `bandingKetat` yang longgar membalik LEBIH dari satu kasus | ✅ Daftar lengkap di bawah, ditulis di **satu** tempat; kode dan test menunjuk ke sini alih-alih menyalinnya |
| 9 | Standards | `pergeseranTPTR` diklaim lebih keras di kode daripada di tiket, dan tanpa path + baris; `ADR-0012` tanpa awalan seri | ✅ Komentar dilunakkan menjadi *"kecocokan cacah argumen, bukan bacaan langsung"*, diberi path dan nomor baris; `ADR-U-0012` |

### ⭐ Daftar LENGKAP kasus yang berbalik bila asumsi diputuskan lain

Ditulis **hanya di sini**; `services/dol.go` dan `dol_test.go` menunjuk ke bab ini.

**Bila `pergeseranTPTR` ternyata HARI, bukan jam** — hanya `TP`/`TR` terpengaruh:

| Kasus | Sekarang | Menjadi |
| --- | --- | --- |
| `TP`/`TR` tepat di tanggal mulai | sah | tetap sah |
| `TP`/`TR` satu jam sebelum berakhir | sah | **tidak sah** *(`DOL + 1 hari` melewati berakhir)* |
| `TP`/`TR` satu jam sebelum mulai | tidak sah | **sah** *(`DOL + 1 hari` sudah melewati mulai)* |

**Bila `bandingKetat` ternyata longgar (`>=`)** — **kedua** ujung dan **kedua** cabang berbalik:

| Kasus | Sekarang | Menjadi |
| --- | --- | --- |
| `QP`/`QR` tepat di tanggal mulai | tidak sah | **sah** |
| `QP`/`QR` tepat di tanggal berakhir | sah | **tidak sah** |
| `TP`/`TR` satu jam sebelum mulai | tidak sah | **sah** |
| `TP`/`TR` tepat di tanggal mulai | sah | tetap sah |

⚠️ Longgar membuat jendela `QP`/`QR` **bergeser**, bukan melebar: ia memperoleh tanggal mulai dan
kehilangan tanggal berakhir. Ronde pertama bab ini hanya menyebut yang pertama — kehilangan di ujung
lain adalah regresi yang tidak terlihat bila hanya satu ujung yang diperiksa.

### Verifikasi penuh sesudah perbaikan review

`go vet` bersih · `go vet -tags=db` bersih · `gofmt -l` nol · `go build` · **158 PASS · 0 FAIL**
*(dari 144)* · **24 SKIP** bertag `db` · `tsc --noEmit` bersih · **5** test JS · `vite build`
**88 modul**.

### ⛔ Yang masih tersisa

| Yang kurang | Sebab |
| --- | --- |
| AC *"`ContentNote` **terisi**"* | derivasinya benar untuk seluruh 21 kode dan teruji satu per satu, tetapi **belum ada layar** yang menampilkannya. Tidak ada kolom `CONTENT_NOTE`, dan menurut XML memang tidak perlu ada — Pega pun menurunkannya saat itu juga. Yang kurang layarnya, bukan kolomnya |
| Test `db` untuk pintu tanggal kejadian | jalurnya diuji di seam `services` *(pelaku anonim, pengenal kosong, tanpa Oracle)*; uji pulang-pergi terhadap Oracle menunggu G1 |
| Frontend formulir DOL | pintunya ada dan menjawab kalimat XML; layarnya belum |
