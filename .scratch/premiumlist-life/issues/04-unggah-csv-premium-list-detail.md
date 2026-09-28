# 04: Unggah CSV premium list detail — staging, validasi, tinjau, simpan

**Status:** sebagian — pemeriksaan keberadaan medan header langkah 3–8 belum ada (verifikasi "master agen" dibantah XML, sensus 28-09-2026); kolom uang kosong → 0 menunggu OQ-PL-12; tabel staging sengaja ditiadakan (ralat 28-09)

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**, 03 (premium list detail — unggahan mengisi struktur yang dibentuk di sana)

## Hasil & nilai pengguna

Sebagai **inputor Life**, saya ingin mengunggah berkas CSV berisi ratusan baris peserta dan
**melihat dulu** apa yang lolos dan apa yang ditolak beserta **nama kolomnya**, sebelum apa pun
tersimpan permanen — supaya satu baris rusak tidak mencemari premium list dan saya tidak perlu
menebak kolom mana yang salah. *(User story 17–24 di spec)*

## Area codebase

`internal/handlers` (endpoint unggah, endpoint tinjau, endpoint simpan permanen), `internal/services`
(mesin validasi per kolom; normalisasi uang; deteksi duplikat), `internal/repository` (tabel staging
+ pembersihannya), `frontend/` (form unggah, tabel hasil validasi berlabel kolom, tombol simpan).

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `UploadCSVLifePremium_Act` | `@BASECLASS` / `UPLOADCSVLIFEPREMIUM_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/UploadCSVLifePremium_Act.xml` | impor berkas — **rule base-class, dipakai bersama** |
| `ValidasiUploadPL_act` | `ASM-FW-GISFW-WORK-LIFE` / `VALIDASIUPLOADPL_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/ValidasiUploadPL_act.xml` | **43 langkah** validasi |
| `InsertDataUploadLife` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!INSERTDATAUPLOADLIFE` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/InsertDataUploadLife.xml` | `insert into POOLDATA.M_TEMPUPLOADLIFE (…)` — **staging** |
| `DeleteTempUploadDataLife` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!DELETETEMPUPLOADDATALIFE` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/DeleteTempUploadDataLife.xml` | `… M_TEMPUPLOADLIFE where idpega = {pyWorkPage.pyID}` |
| `CekDoubleInsured` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!CEKDOUBLEINSURED` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/CekDoubleInsured.xml` | duplikat peserta + jumlah retensi ceding |


⚠️ **Di luar cakupan tiket ini:** jalur unggah CSV **endorsement**
(`Endorsement Life/Activity/UploadCSVEDMLifePremium_Act.xml`, `@BASECLASS` /
`UPLOADCSVEDMLIFEPREMIUM_ACT`, dan `Endorsement Life/Activity/SaveCSVEDMLife.xml`,
`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE`). `[keputusan work owner]` **Endorsement Life
adalah konteks terpisah** — lihat `.scratch/endorsement-life/`.

⚠️ **RALAT 28-09-2026 — `M_TEMPUPLOADLIFE` BUKAN staging penuh.** `InsertDataUploadLife.xml`
menyisipkan **TUJUH** kolom saja: `INSURED, DOB, BEGINDATE, ENDDATE, POLIVYHOLDER, NO,
CEDING_RETENTION` *(ejaan `POLIVYHOLDER` apa adanya di korpus)*. Ia melayani **satu** keperluan —
pemeriksaan duplikat `CekDoubleInsured`. Baris lengkapnya di Pega tinggal di **halaman kerja**, bukan
di tabel itu; validasi berjalan atas halaman kerja.

Karena itu implementasi kami **tidak membuat tabel staging**: tinjauan mengurai dan memvalidasi tanpa
menulis apa pun, dan yang permanen baru tersentuh saat pemakai menekan simpan — saat mana berkasnya
diurai dan divalidasi **ulang**. Substansi AC 18 terpenuhi *(kegagalan validasi tidak menyentuh tabel
permanen mana pun)*; tabel staging penuh adalah **keputusan skema**, dan migrasi baru hanya dari
keputusan yang tercatat.

`[terverifikasi]` **`CekDoubleInsured`** mendeteksi duplikat dengan `upper(INSURED)` + `DOB`,
mengembalikan `min(NO)` baris pertama, dan menjumlahkan `to_number(CEDING_RETENTION)` per peserta.
⚠️ Variabel `v_Count2` dan `v_PLRetensiCeding` **dideklarasikan tetapi tidak pernah diisi** — selalu
`NULL`, dinetralkan `nvl(…,0)`. Sisa pemeriksaan kedua yang dicabut; **jangan** direplikasi.

⚠️ **RALAT 28-09-2026 — 32, bukan 33.** Sensus dihitung DUA cara dan keduanya berbeda; selisihnya
disebut, tidak didiamkan:

| Cara | Jendela | Hasil |
| --- | --- | ---: |
| A | precondition `@PropertyHasValue(.X)` di langkah 2 (`RH_1.pySteps(2)`..`(3)`) | **32** |
| B | `<pyStepsDescription>` sub-langkahnya yang berisi nama kolom huruf besar | 31 |

Selisihnya **satu dan bernama**: `FLEET_DISCOUNT` punya precondition tetapi deskripsi sub-langkahnya
**kosong**. Yang benar **32** — precondition-lah yang **berjalan**; deskripsi hanya nama yang dibaca
manusia. Sensus yang memakai cara B akan menghilangkan satu kolom uang dari validasi, dan kolom uang
yang tidak divalidasi adalah kolom uang yang menerima apa saja.

`[terverifikasi]` **32 kolom uang** divalidasi (`<pyStepsDescription>` per langkah): `NET_PREMIUM`,
`GROSS_PREMIUM`, `SHARE_NUSANTARA_RE`, `SUM_INSURED`, `CEDING_RETENTION`, `SUM_REASURED`, `CLAIM`,
`TAX`, `BROKERAGE_FEE`, `OVR_COMM`, `PROF_COMM`, `EM_PERCENT`, `COMM`, `FLEET_DISCOUNT`, seluruh
kelompok `*_REFUND`, seluruh kelompok `*_RETRO`, `*_REFUND_RETRO`, dan `CLAIM_AMOUNT`.

`[terverifikasi]` Validasi non-uang: `CERTIFICATE_NO` tidak boleh duplikat; `NAME_OF_INSURED` dan
`POLICY_HOLDER` **harus ada di `M AGENT`**; `DOB`, `BEGIN_DATE`, `EXPIRED_DATE` format `dd/mm/yyyy`;
`PLAN`, `CURRENCY` wajib; `MEDICAL_STATUS` salah satu dari `FCL` / `M` / `NM`; lampiran wajib
("BELUM ADA LAMPIRAN").

> ⛔ **Ralat 28-09-2026 (sensus remark GILIRAN-12).** Sertifikat sudah dipakai (9.2–9.4, 12),
> sertifikat ganda (9.5, 10, 13), `MEDICAL_STATUS` (9.11, 19), dan lampiran (11, 33) **ter-remark**
> (`//` b8025…b16535) — sistem lama tidak menolak karenanya. "`POLICY_HOLDER` harus di M AGENT"
> adalah deskripsi langkah 12 yang ter-remark (pesannya err1). Yang hidup: langkah 3–8 (b6874–b7767,
> keberadaan medan header: ceding, policy holder, MO, COB, retro, SOB — pesannya ada di kode, belum
> dipakai) dan langkah 2, yang **mengisi 0** kolom uang yang kosong (OQ-PL-12). `CekDoubleInsured`
> hidup, tetapi di `Calculate1_Act` 7.5 untuk akumulasi retensi, bukan penolakan unggah.

⚠️ `[terverifikasi]` **Normalisasi desimal Pega berbahaya.** Langkah "Rubah decimal dari koma jadi
titik" (baris 12480) memakai `@replaceAll(.KOLOM, ",", ".")` atas tiap kolom uang — **mengganti
setiap koma**, tanpa membedakan pemisah desimal dari pemisah ribuan. Nilai `1,234,567.89` menjadi
`1.234.567.89`, yang bukan angka. **Jangan** direplikasi apa adanya.

## ADR terkait

**ADR-0003** (uang non-float, desimal presisi arbitrer — inti tiket ini), **ADR-0010** (penyimpanan
berkas tetap Google Storage untuk lampiran), **ADR-0007** (jejak audit unggahan).

## Acceptance criteria

- [ ] Berkas yang diunggah masuk **staging** lebih dulu; kegagalan validasi **tidak** menyentuh tabel
      permanen mana pun. *(AC 18 spec)* — belum: sengaja tanpa tabel staging (ralat 28-09); kegagalan validasi memang tidak menyentuh tabel permanen (uji `TestTinjauTidakMenyentuhApaPun`, `TestSimpanMemvalidasiUlang`)
- [x] Pengguna dapat **meninjau** hasil unggahan — baris lolos dan baris ditolak — sebelum menyimpan
      permanen. *(AC 19 spec)* — bukti: `services/polis_unggah.go:Tinjau`; uji `TestTinjauTidakMenyentuhApaPun`
- [x] Setiap penolakan menyebut **nama kolom** yang salah dan **nomor baris** CSV-nya. *(AC 17 spec)* — bukti: uji `TestSetiapPenolakanMenyebutKolomDanBaris`
- [x] Nilai uang di-parse dengan **format yang dinyatakan eksplisit** (pemisah desimal dan pemisah
      ribuan ditentukan, bukan ditebak). `1,234,567.89` dan `1.234.567,89` **tidak** boleh keduanya
      diterima diam-diam sebagai angka yang sama. Test wajib memuat kasus pemisah ribuan. — bukti: `models/polis_unggah.go:UangCSV`; uji `TestUangCSVMenolakPemisahRibuan`
- [x] Nilai uang **tidak** melewati `float` pada tahap mana pun — parse langsung ke desimal presisi
      arbitrer. *(AC 14–15 spec; **ADR-0003**)* — bukti: `models/polis_unggah.go:UangCSV` (desimal apd); uji `TestUangCSVTidakMembulatkan`, `TestUangDikirimSebagaiTeks`
- [ ] Nilai uang yang lolos validasi, dibaca kembali dari staging, **identik** dengan yang diunggah —
      tidak ada pembulatan diam. *(AC 16 spec)* — belum: tidak ada staging, dan nol uji baca-kembali terhadap Oracle
- [x] ~~Duplikat peserta terdeteksi dengan pembandingan nama tanpa peduli huruf besar/kecil +
      tanggal lahir~~ — *disunting di tempat 28-09-2026 (sensus remark):* **tidak ada penolakan
      duplikat peserta saat unggah** — `CekDoubleInsured` menghitung retensi di `Calculate1_Act` 7.5,
      bukan menolak; kalimat err1 yang dipinjam milik langkah ter-remark. Penolakannya dibuang — bukti: uji `TestPenolakanLangkahTerRemarkTidakDitegakkan`
- [x] ~~`CERTIFICATE_NO` ganda di dalam satu berkas ditolak~~ — *disunting di tempat 28-09-2026:*
      **tidak** ditolak — langkah 9.5, 10, 13 ter-remark — bukti: uji `TestPenolakanLangkahTerRemarkTidakDitegakkan`
- [ ] ~~`NAME_OF_INSURED` dan `POLICY_HOLDER` diverifikasi terhadap master agen~~ — *disunting di tempat
      28-09-2026 (sensus remark):* tidak ada verifikasi master agen di sistem lama — 9.6 hanya menguji
      keberadaan (`@PropertyHasValue(.NAME_OF_INSURED)` b8915; sudah: `kolomTeksWajib`), dan langkah 12
      ("HARUS ADA DI M AGENT") ter-remark. Yang HIDUP: langkah 3–8 menguji **keberadaan** medan header
      polis — ceding (b6979), policy holder (b7132), MO (b7283), COB (b7434), retro bila bukan QR/QP
      (b7585/b7616), SOB (b7767) — dengan pesan `… TIDAK TERDAFTAR` (langkah 27–32). — belum: pesannya
      ada (`models/polis_unggah.go` `PesanCedingCo`…`PesanSOB`), pemeriksanya belum
- [x] Tanggal diverifikasi berformat `dd/mm/yyyy`; ~~`MEDICAL_STATUS` hanya menerima `FCL`, `M`, `NM`~~
      *(ter-remark 9.11/19 — dibuang, sensus 28-09-2026)*; `CURRENCY` wajib ada. — bukti: `models/polis_unggah.go:ValidasiUnggah` (`CURRENCY` di `kolomTeksWajib`); uji `TestTanggalCSVKetatDDMMYYYY`, `TestPenolakanLangkahTerRemarkTidakDitegakkan`
- [x] ~~Unggahan tanpa lampiran ditolak~~ — *disunting di tempat 28-09-2026:* langkah 11/33 ter-remark;
      permintaan tanpa berkas kini galat permintaan biasa (400 di handler, `ErrPermintaanTidakSah` di
      layanan), bukan penolakan korpus — bukti: `services/polis_unggah.go:periksaBerkas`
- [ ] Staging dibersihkan per case setelah simpan permanen **maupun** setelah pembatalan — tidak ada
      sisa baris menggantung milik case lain. — belum: tidak ada tabel staging (ralat 28-09), jadi pembersihan staging tidak dibangun
- [x] Mengunggah berkas kedua atas case yang sama **mengganti** isi staging, tidak menumpuk. — bukti: `services/polis_unggah.go:Simpan` (`HapusPesertaPolis` lalu `SisipPeserta` dalam satu transaksi; atas tabel peserta, tanpa staging)

## Blocker

⚠️ **OQ-069 terbuka** — pesan Pega `"NET PREMIUM HARUS ADA DAN LEBIH BESAR DARI GROSS PREMIUM"`
menjanjikan aturan yang **tidak pernah diperiksa**: satu-satunya precondition atas `NET_PREMIUM` di
seluruh berkas adalah `@PropertyHasValue(.NET_PREMIUM)`, dan **tidak ada** perbandingan terhadap
`GROSS_PREMIUM`. Arah perbandingannya pun janggal — lazimnya net lebih **kecil** dari gross.

**Selama OQ-069 terbuka:** tiket ini memeriksa **keberadaan** `NET_PREMIUM` saja, persis seperti
korpus. **JANGAN** mengarang aturan perbandingan. Pesan Pega dicatat apa adanya di kode sebagai
rujukan, dengan penunjuk ke OQ-069.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## Implementasi

**Dikerjakan 28-09-2026.** `main`, sesudah `2cbdf2d`.

### Yang dikirim

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `polis_unggah.go` | 32 pesan VERBATIM · `UangCSV` · `TanggalCSV` · `StatusMedisDiterima` · `KolomUangUnggah` (32) · `KolomWajibUnggah` · `ValidasiUnggah` |
| pkg/utils | `decimal.go` | `ParseDecimal` menolak `NaN`/`Infinity` — lihat di bawah |
| services | `polis_unggah.go` | `BacaCSVUnggah` (judul, BOM, batas) · `Tinjau` (nol tulisan) · `Simpan` (validasi ulang, satu transaksi) |
| repository | `polis_unggah.go` | `HapusPesertaPolis` · `SisipPeserta` · `PengenalPesertaUnggah` |
| handlers | `rute_unggahpolis.go` | `POST …/unggah/tinjau` · `POST …/unggah/simpan` |
| frontend | `UnggahCSVPeserta.tsx` | pilih → tinjau → simpan; tabel penolakan berkolom **Baris · Kolom · Pesan · Sebab** |

### Aturan uang, dan kenapa berbeda dari Pega

`[terverifikasi]` Korpus menyatakan aturannya **enam kali** di nama langkah
(*"SEPARATOR MENGGUNAKAN TITIK"*) dan **menegakkannya sekali** — satu-satunya pemeriksaan koma di
seluruh berkas adalah `@contains(.GROSS_PREMIUM,",")`. Pesan yang dilihat pemakai **tidak pernah**
menyebut aturan itu; ia hanya berbunyi `SUM INSURED HARUS ADA`.

⛔ **Normalisasi Pega tidak ditiru.** `@replaceAll(.KOLOM, ",", ".")` mengganti **setiap** koma:
`1,234,567.89` menjadi `1.234.567.89`, yang bukan angka sama sekali — lalu tersimpan sebagai uang.

**Yang kami kerjakan:** koma **ditolak**, untuk **seluruh** 32 kolom uang, dan sebabnya disebut di
pesan kami sendiri *(pesan verbatim tetap dibawa berdampingan)*. Aturannya juga **dinyatakan di
muka** di layar, bukan hanya saat menolak — orang yang baru tahu aturannya setelah berkasnya ditolak
sudah terlanjur menyiapkan berkas yang salah.

`1,234,567.89` dan `1.234.567,89` keduanya **ditolak**, tidak satu pun diterima diam-diam.

### Cacat lintas modul yang ditemukan dan diperbaiki

⛔ **`utils.ParseDecimal` menerima `NaN` dan `Infinity`.** `apd.NewFromString` mengikuti spesifikasi
desimal, dan di sana keduanya nilai yang sah. Fungsi itu **satu-satunya** jalan masuk teks-ke-desimal
(**ADR-U-0034**) dan punya **sepuluh** pemanggil — salah satunya pembaca uang dari **JSON**, yaitu
batas yang dilewati permintaan dari luar. Badan permintaan berisi `{"amount":"NaN"}` akan lolos
seluruh validasi.

Ditemukan oleh satu uji yang sengaja mencoba nilai aneh, bukan oleh tinjauan. Ditolak sekarang di
`ParseDecimal` — menutup kesepuluh pemanggil sekaligus. Penjaganya
(`TestParseDecimalMenolakYangTidakBerhingga`) **dibuktikan merah** lebih dahulu, dan pasangannya
(`…TetapMenerimaBilanganBiasa`) menjaga supaya penjepitan itu tidak menolak terlalu banyak.

### Pengenal baris peserta — nol sequence

Migrasi 050–056 tidak membuat satu pun sequence, dan tiket **00** sudah memutuskannya:
*"Nol sequence. Pengenalnya dirakit di `repository` mengikuti pola `PengenalWorkBerikut`."*
`PengenalPesertaUnggah(polisID, nomorBaris)` menghasilkan 32 heksa — **tepat** selebar kolom
`VARCHAR2(32)` — dan **deterministik**, sehingga unggah ulang menghasilkan baris yang sama persis dan
dapat dibandingkan. MD5 dipakai sebagai **pemadat**, bukan pengaman: ini pengenal baris, bukan kunci
penyimpanan *(bandingkan `models.ImageIDBaru`, yang justru harus tidak dapat ditebak)*.

### Penjaga yang menuduh hal yang BENAR, dan dipersempit

`TestNolNamaOrangDiKode` menyalakan dua konstanta pesan VERBATIM
(`PesanNamaTertanggung`, `PesanPolicyHolder`) — namanya memuat kata yang dijaga sebab itulah **kolom**
yang divalidasi, dan nilainya kalimat galat, bukan nama siapa pun. Dipersempit dengan **daftar
bernama** yang kuncinya diambil **dari `models`**, bukan diketik ulang: ronde pertama penyempitan
membuat daftarnya **menuduh dirinya sendiri**. Alasannya di komentar, bukan di nilai — sebab nilai
teks pada baris itu juga cocok dengan polanya. Dibuktikan masih menggigit dengan nama sungguhan.

### Butir `[terbuka]` yang LAHIR di tiket ini

- **OQ-PL-05 — `CekDoubleInsured` menyaring `idpega` di SATU dari TIGA pernyataan.**
  `count(1)` dan `min(NO)` **tidak** menyaring `idpega`; hanya `sum(CEDING_RETENTION)` yang
  menyaring. Artinya pemeriksaan duplikat di sistem lama berlaku **lintas case** sementara jumlah
  retensinya per case. Asimetri dalam satu procedure adalah tell bahwa salah satunya kelalaian —
  tetapi korpus tidak memberi tahu yang mana. **Kami menyaring per polis** *(sesuai AC: duplikat di
  dalam satu berkas)*, dan lintas-case dicatat di sini. ⚠️ *Sensus 28-09-2026:* pemanggil
  `CekDoubleInsured` adalah `Calculate1_Act` 7.5 (akumulasi retensi); penolakan duplikat unggah
  dibuang — OQ ini kini soal hitungan retensi, bukan validasi.
- **OQ-PL-06 — `upper(INSURED)` hanya di SATU sisi.** Procedure membandingkan `upper(INSURED)`
  dengan parameternya **apa adanya**; kecocokannya karena itu bergantung pemanggil yang mengirim
  nama dalam huruf besar. Kami menaikkan **kedua** sisi, dan ujinya mengunci itu — satu sisi yang
  lupa membuat duplikat lolos tanpa jejak.
- **OQ-PL-07 — pengecualian `02/01/1970`.** `@if(.DOB=="02/01/1970",true,@toDate(.DOB)!=0)`: Pega
  mengecualikannya sebab `@toDate` di sana mengembalikan **nol** untuk tanggal itu dan pemeriksanya
  membandingkan dengan nol. Pengurai kami tidak punya cacat itu, jadi tanggal lahir 2 Januari 1970
  diterima **tanpa** pengecualian. Dicatat supaya tidak ada yang "merapikannya" menjadi aturan yang
  menolak tanggal itu.

### OQ-069 tetap TERBUKA, dan kodenya mematuhinya

Pesan `NET PREMIUM HARUS ADA DAN LEBIH BESAR DARI GROSS PREMIUM` disalin **verbatim**, dan yang
diperiksa **hanya keberadaannya** — persis korpus. `TestNetPremiumHanyaDiperiksaKeberadaannya`
menegakkan itu dengan baris ber-`NET < GROSS` yang **harus lolos**. Mengarang perbandingannya akan
menolak berkas yang di sistem lama diterima, dan arah yang salah *(net lazimnya lebih kecil)* akan
menolak setiap berkas yang benar.

### Yang TIDAK dikerjakan, dan sebabnya

| Butir | Sebab |
| --- | --- |
| ~~`NAME_OF_INSURED` / `POLICY_HOLDER` diverifikasi terhadap **master agen**~~ | *dibantah XML (sensus 28-09-2026):* langkah 12 "M AGENT" ter-remark; yang hidup langkah 3–8, pemeriksaan **keberadaan** medan header — belum dibangun, lihat AC-nya |
| Spreading dan summary uang | tiket **05a** |
| Duplikat **lintas case** | OQ-PL-05 di atas |

### Verifikasi

```
gofmt -l .              bersih
go vet ./... · -tags db bersih
go test ./... -tags db  497 PASS · 0 FAIL · 38 SKIP   (dari 463)
npx tsc --noEmit        bersih
npx vitest run          334 PASS                       (dari 322)
npm run build           bersih
```

Nol migrasi baru. Nol procedure dipanggil. Nol `COMMIT` di teks SQL. Nol nama orang di fixture —
seluruh nilai uji berawalan `UJI-`.

## ⛔ Ralat bertanggal — 28 September 2026 (GILIRAN-12 paket 2: sensus remark)

`ValidasiUploadPL_act` dicetak beserta `pyStepsBlockName`: langkah 9.2–9.5, 9.11, 10, 11, 12, 13, 19,
33 ber-`//`. Lima penolakan yang ditiru dari sana — sertifikat terpakai, sertifikat ganda,
`MEDICAL_STATUS`, lampiran, dan duplikat nama+DOB (kalimat err1 dipinjam) — **dibuang** dari
`models.ValidasiUnggah` beserta ujinya; `MEDICAL_STATUS` tidak lagi kolom judul wajib. Teks AC yang
menuntutnya disunting di tempat. Celah hidup yang dicatat: langkah 2 mengisi 0 kolom uang kosong
(**OQ-PL-12**), langkah 3–8 memeriksa keberadaan medan header (pesannya ada, belum dipakai).
