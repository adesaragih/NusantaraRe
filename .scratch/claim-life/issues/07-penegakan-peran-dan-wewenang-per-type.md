# 07: Penegakan peran di lapisan layanan + wewenang kirim-Komite per `Type`

**Status:** selesai — 28-09-2026, diverifikasi atas `5ac6571`

**Blocked by:** 05 (reject Outstanding oleh Admin) — gerbang perlu tindakan nyata untuk dijaga

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin setiap tindakan pada klaim hanya dapat dilakukan peran yang
berhak — dan penegakannya berada di **lapisan layanan**, bukan sekadar di visibilitas layar —
sehingga wewenang tidak dapat dilewati dengan memanggil API langsung. *(User story 15 di spec)*

## Area codebase

`internal/services` (penegakan wewenang sebelum setiap transisi), `internal/handlers` (identitas
pemanggil), `frontend/` (kontrol yang tampil/tersembunyi mengikuti peran — **sebagai kenyamanan,
bukan sebagai penegakan**).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | tiga gerbang `<pyCondition>` di bawah |
| `Claim Life/Flow/Register_Flow.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW` | `[terverifikasi]` 15 kemunculan `pyPosition` — penegakan sesungguhnya di sistem lama bertumpu pada penugasan tahap di sini |

`[terverifikasi]` Tiga gerbang pada `AdjustmentDetail_Section`:

| Kontrol | `<pyCondition>` |
| --- | --- |
| "Reject Outstanding" | `pyWorkPage.pyPosition =='ReasLifeAdmin' && …CLAIM_NO !='' && .STS_REJECT == 0` |
| Jalur Komite (`GetListKomiteLife`) | `pyWorkPage.pyPosition =='ReasLifeSPV' \|\| pyWorkPage.Type = 'TP' \|\| pyWorkPage.Type = 'TR'` |
| "Save to Outstanding" | `pyWorkPage.pyPosition =='ReasLifeSPV'` |

`[terverifikasi]` `Type` dibaca dari **dua salinan** di Pega — `pyWorkPage.Type` (gerbang wewenang)
dan `pyWorkPage.PolicyDataLife.Type` (validasi DOL) — disalin di
`Claim Life/Activity/LoadDataPeserta_Act.xml` (`ASM-FW-GCNMFW-WORK` / `LOADDATAPESERTA_ACT` /
`RULE-OBJ-ACTIVITY`): `pyWorkPage.Type = pyWorkPage.PolicyDataLife.Type`. **Sumber otoritatif =
`PolicyDataLife.Type`.**

## ADR terkait

**ADR-0002** (tiga peran; rangkap peran tidak diperbolehkan), **ADR-0012** (wewenang kirim-Komite
bergantung `Type`; **dibawa apa adanya sebagai paritas**, risiko RBAC diterima dan ditinjau saat
konteks Komite/IAM digarap).

## Acceptance criteria

- [x] `ReasLifeMedicalAdvisor` **tidak dapat** mengubah status akseptasi baris mana pun. *(AC 9 spec)* — bukti: `APP_RNM/internal/services/wewenang.go:WajibPeranPengubahStatus`; uji `TestMedicalAdvisorTidakDapatMengubahStatus`
- [x] Untuk klaim ber-`Type` `QP` atau `QR`, **hanya `ReasLifeSPV`** yang dapat mengirim ke Komite;
      upaya oleh peran lain ditolak. *(AC 10 spec)* — bukti: `APP_RNM/internal/services/wewenang.go:WajibWewenangKomite`, dipanggil `APP_RNM/internal/services/komite.go:Penyerahan.Serahkan`; uji `TestWewenangKirimKomitePerType`
- [x] Untuk klaim ber-`Type` `TP` atau `TR`, `ReasLifeAdmin` **dapat** mengirim ke Komite.
      *(AC 11 spec)* — bukti: `APP_RNM/internal/services/wewenang.go:WajibWewenangKomite`, dipanggil `APP_RNM/internal/services/komite.go:Penyerahan.Serahkan`; uji `TestWewenangKirimKomitePerType`
- [x] Penolakan wewenang terjadi **di lapisan layanan**, dan tetap terjadi meskipun kontrol UI-nya
      ditampilkan. *(AC 12 spec)* — bukti: uji `TestWewenangDitegakkanDiLayanan`, `TestSetiapPenulisStatusBergerbangPeran`
- [x] Wewenang dan validasi membaca **satu** field `Type` yang sama — tidak ada dua salinan yang
      dapat berbeda. — bukti: `APP_RNM/internal/repository/klaimlife.go:KlaimLife.TypeKlaim` (dibaca `Penyerahan.Serahkan` dan `TanggalKejadian.Set`); uji `TestTypeKlaimHanyaSatuRumahTersimpan`
- [x] Tidak ada nama orang ter-hardcode di lapisan mana pun; wewenang diukur dari peran akun. — bukti: uji `TestNolNamaOrangDiKode`

## Catatan

⚠️ `[terverifikasi]` Penegakan di sistem lama **tidak seragam**: `pyPosition` tidak ada sama sekali
di `Section/InputOSClaimLife.xml`, `Section/RejectOSClaimLife_Sec.xml`, `Section/ClaimComite.xml`,
maupun `Harness/Committe_Life.xml`. **Ketidakseragaman itu tidak ditiru** — sistem baru menegakkan
di lapisan layanan untuk semua tindakan.

⚠️ Kelonggaran `TP`/`TR` **adalah risiko yang diterima sadar** (**ADR-0012**): wewenang bergantung
pada atribut data, bukan peran. Jangan "memperbaikinya" dalam tiket ini — itu keputusan terpisah.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Pembacaan ulang XML — 26 September 2026 malam

XML tiket ini **sudah dibaca dan dicatat** di tiket 03 *(gerbang `AdjustmentDetail_Section`)* dan
tiket 04 *(sensus penulis status)*; sesuai brief lanjutan 2 §1.4-c, keduanya dirujuk, bukan dibaca
ulang. Yang **baru** dibaca giliran ini: `Flow/Register_Flow.xml` — peran per tahap.

| Tahap | `pyPosition` `[terverifikasi]` |
| --- | --- |
| `Assignment2` *(Input Register)* | `ReasLifeAdmin` |
| `Assignment1` *(Outstanding Claim)* | `ReasLifeAdmin` |
| `Assignment3` *(Medical Check)* | `ReasLifeMedicalAdvisor` |
| `Decision3` | `ReasLifeMedicalAdvisor` |
| `Decision1` | `ReasLifeSPV` |

⭐ **`ReasLifeMedicalAdvisor` tidak muncul di satu pun gerbang status.** Ia memegang tahap telaah
medis dan keputusan medis, bukan keputusan akseptasi — dan itulah dasar AC 9.

⚠️ **Selisih tiket ↔ XML yang dicatat, bukan dipersempit.** Gerbang Komite berbunyi
`pyPosition=='ReasLifeSPV' || Type='TP' || Type='TR'`. Bacaan **harfiah**nya: untuk `TP`/`TR`
gerbangnya terbuka **tanpa memeriksa peran sama sekali**. Tiket menuliskannya *"`ReasLifeAdmin`
dapat mengirim"*. Keduanya sejalan selama Admin memang yang memegang tahapnya, tetapi mempersempit
kode menjadi *"hanya Admin"* akan **menolak SPV pada TP/TR** padahal XML menerimanya. Kode mengikuti
XML; identitas tetap wajib — gerbang yang terbuka bukan gerbang yang hilang.

## Implementasi — 26 September 2026 malam (tiket 07)

**Status: `claimed`** — **4 dari 6 AC tertutup** *(6 sebelum `/code-review`; dua AC Komite dibatalkan centangnya karena fungsinya nol pemanggil)*. Titik tetap `c7fcd27`.
Verifikasi: vet · vet db · gofmt nol · build · **182 PASS · 0 FAIL** *(dari 176)* · 28 SKIP ·
`tsc` · 7 JS · 88 modul.

| Berkas | Isi |
| --- | --- |
| `services/wewenang.go` *(baru)* | `PeranMedicalAdvisor`, `WajibPeranPengubahStatus`, `WajibWewenangKomite` |
| `services/statusbaris.go` | gerbang peran dipasang di `ubah` — **satu** jalan menuju tulisan status, dilewati `Ubah` dan `Tolak` |
| `models/satutype_test.go` | dua penjaga statik baru |

⛔ **Peran tidak diberi nama kedua.** `PeranSimpanOutstanding` *(SPV, tiket 03)* dan
`PeranRejectOutstanding` *(Admin, tiket 05)* dipakai apa adanya. Satu peran, satu nama, walau tempat
lahirnya berbeda.

### ⛔ Penjaga nama orang: percobaan pertama melanggar aturan yang dijaganya

Ronde pertama memuat **daftar nama operator nyata** dari korpus sebagai pola terlarang — dan itu
sendiri melanggar pagar keamanan brief induk: *"nilai nama orang tidak disalin ke artefak mana pun"*.
Penjaga yang melanggar aturannya sendiri salah **bentuk**, bukan perlu dikecualikan. Ia ditulis ulang
menjadi **struktural**: pemberian teks tetap ke medan yang namanya menandakan nama orang
(`OpName`, `PolicyHolder`, `NameOfInsured`, …), dengan nilai sintetis `UJI-*` dikecualikan — sebab
bentuk itu justru yang pagar keamanan **tuntut** untuk fixture. Dibuktikan menggigit: menyisipkan
`CreateOpName: "<nama>"` membuatnya merah.

**Penjaga aturan bisnis lain, dibuktikan dapat gagal:** gerbang peran dilepas dari `ubah` → penjaga
statik merah; Medical Advisor mengubah status → 3 test merah.

### Hasil `/code-review` atas titik tetap `c7fcd27`

⛔ **Review membuktikan penjaga saya HIJAU di bawah mutasi.** Sub-agen menambahkan satu fungsi
penulis status di berkas `services` yang **baru**, dan penjaga statik saya tetap hijau — sebab ia
membaca **satu** berkas dan mencari **satu** potongan teks, sambil mengaku memeriksa *"setiap fungsi
layanan yang mengubah status"*. Pengakuan yang lebih luas daripada yang diperiksanya; pelajaran yang
berkas itu sendiri catat, saya ulangi.

| # | Temuan | Tindakan |
| ---: | --- | --- |
| 1 | ⛔ **Daftar peran yang DATAR** — `{SPV, Admin}` untuk tujuan apa pun, sehingga **SPV dapat menolak baris** lewat `Ubah`, melewati gerbang Admin di `Tolak` beserta syarat klaim-bernomornya | ✅ gerbangnya kini bergantung **tujuan**: Ditolak → Admin *(gerbang 15399)*; Aksep → **ditolak bagi siapa pun di modul ini**, sebab `[terverifikasi]` hanya Komite yang menulis `1` |
| 2 | ⛔ **Dua penulis status melewati gerbang**: `Pendaftaran.Daftar` *(lewat `TandaiOutstanding`)* dan komentar `hapus.go` yang masih menunggu tiket ini | ✅ `Daftar` bergerbang `PeranInputRegister` *(`Assignment2` = `ReasLifeAdmin`)*; komentar `hapus.go` diralat |
| 3 | ⛔ **Penjaga statik tidak menjaga** *(terbukti lewat mutasi di atas)* | ✅ ditulis ulang: menelusuri **seluruh** `internal/`, berdaftar-izin per berkas beserta gerbang yang wajib dipanggil, dan menolak penulis status di berkas yang tidak terdaftar. Dibuktikan menangkap mutasi yang sama |
| 4 | **`WajibWewenangKomite` nol pemanggil produksi** — AC 10 dan 11 tercentang padahal tidak ada yang dapat dicoba, apalagi ditolak | ✅ **dibatalkan centangnya**; keduanya ditutup di **tiket 10** bersama pintunya |
| 5 | Teks `"ReasLifeAdmin"` tertulis di **tiga** berkas walau tiket ini menulis *"satu peran, satu nama"* | ✅ `PeranAdmin`/`PeranSPV`/`PeranMedicalAdvisor` satu tempat; izin tetap bernama sendiri sebab artinya berbeda |
| 6 | `switch` keempat `Type` di dua berkas, urutan cabang berbeda | ✅ `TypeDikenal` satu tempat |
| 7 | Penjaga nama orang: pola dapat dielakkan; letaknya di `models_test` yang membaca `../services` — membalik arah ketergantungan | ✅ dipindah ke `services`; pola menerima kutip-balik dan `UJI-` diperiksa di **awal** teks, bukan di mana saja; **batasnya dinyatakan** di komentar — nama lewat konstanta perantara atau perbandingan `AkunID == "..."` tidak tertangkap |

⚠️ Dua catatan yang tetap terbuka: `PeranSimpanOutstanding` masih nol pemanggil *(pintunya milik
sisa tiket 03)*, dan `WajibWewenangKomite` menerima `Type` dari pemanggil — saat tiket 10
mengabelnya, ia wajib membacanya lewat `KlaimLife.TypeKlaim`, bukan menerima teks.

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

## ⛔ PERUBAHAN AUTHZ — 28 September 2026 (GILIRAN-11 paket 3, butir **bj**)

`[DIPUTUSKAN; veto work owner]` — OQ-M8 **ditutup**. Rute DOL lama
`PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-kejadian` kini bergerbang **sama** dengan tiga
tanggal lain dialog Edit Date, menurut XML (*prinsip proyek: XML menang*):

| | Sebelum 28-09-2026 | Sesudah |
| --- | --- | --- |
| Peran | siapa pun yang beridentitas | `ReasLifeAdmin` (`WajibPeran(pelaku, PeranAdmin)`) |
| Tahap | tahap mana pun (kasus terbuka) | Outstanding Claim saja (`gerbangTahapDialogTanggal`) |

Buktinya: isian DOL `EditDateClaimLife_Section` berprasyarat baca-saja
`pyWorkPage.pyPosition!='ReasLifeAdmin' || …CLAIM_NO!=''` (b1000) — prasyarat yang **sama** dengan
ketiga isian lainnya (b1313, b1550, b1788). Separuh `CLAIM_NO` tetap tidak ditiru (OQ-M1). Satu
gerbang untuk keempat isian: `services/dol.go` `gerbangTahapDialogTanggal`, dipakai `Set` dan
`SetTanggalKlaim`; layar mematikan kotak DOL di luar tahap Outstanding (`bolehUbahTanggalKlaim`).

Akibat yang diterima: Medical Advisor dan SPV tidak lagi dapat mengubah DOL — jawabannya 403; di
tahap Medical Check dan Claim Analis — 409. Uji: `TestSetTanggalKejadianHanyaAdmin`.
