# 05: Reject Outstanding oleh Admin — tanpa Komite

**Status:** selesai — 28-09-2026, diverifikasi atas `5ac6571`

**Blocked by:** 04 (mesin status per baris)

## Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, saya dapat menolak baris adjustment yang saya input sendiri tanpa melalui
Komite — sehingga kesalahan input dapat saya batalkan cepat — dan penolakan itu **hanya membatalkan
baris itu**, klaimnya tetap hidup dan saya dapat menginput baris pengganti. *(User story 6 dan 7 di
spec)*

Ini jalur penolakan **kedua** di sistem ini. Jalur pertama (lewat Komite) datang di tiket 11.

## Area codebase

`internal/services` (aturan penolakan + gerbang status baris), `internal/handlers` (endpoint reject),
`frontend/` (kontrol "Reject Outstanding" pada baris adjustment).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `REJECTOSCLAIMLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` menulis `2` ke **dua tingkat berbarengan**: `.STS_REJECT` dan `…PremiumListDetail(idx).STS_REJECT` |
| `Claim Life/Section/RejectOSClaimLife_Sec.xml` | layar pemicu — `[terverifikasi]` merujuk activity di atas 2× sebagai `<pyActivity>` | |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | `[terverifikasi]` gerbang tombol "Reject Outstanding": `pyWorkPage.pyPosition =='ReasLifeAdmin' && pyWorkPage.ClaimData.PremiumListSummary.CLAIM_NO !='' && .STS_REJECT == 0` |

Gerbang itu menguji **`.STS_REJECT` tingkat baris** — bukti bahwa wewenang pun diukur per baris.

## ADR terkait

**ADR-0011** (dua sumber penolakan dengan makna setara pada tingkat baris; `RejectOSClaimLife_Act`
**bukan** dead rule), **ADR-0002** (peran).

## Acceptance criteria

- [x] Penolakan oleh `ReasLifeAdmin` membuat **hanya baris itu** berstatus Ditolak; baris lain pada
      klaim yang sama tidak berubah. *(AC 3 spec)* — bukti: `APP_RNM/internal/services/tolak.go:Status.Tolak`, `APP_RNM/internal/repository/klaimlife.go:KlaimLife.PerbaruiStatusBaris` (`WHERE ID`, wajib tepat satu baris); uji `TestUbahStatusMencerminkanTigaTingkat` (uji db)
- [x] Klaim **tidak** tertutup oleh penolakan itu, dan baris adjustment baru dapat diinput
      sesudahnya. — bukti: `APP_RNM/internal/services/hasilkomite.go:Putaran.Tambah`; uji `TestMenolakSatuBarisTidakMenutupKlaim`, `TestKlaimTidakTerminalSetelahPenolakan`
- [x] Penolakan hanya mungkin pada baris yang **masih Outstanding** — baris yang sudah Aksep atau
      Ditolak menolak upaya itu. — bukti: uji `TestTolakBarisFinalDitolak`, `TestTransisiHanyaDariOutstanding`
- [x] Penolakan hanya mungkin bila klaim sudah punya nomor (padanan `CLAIM_NO !=''`). — bukti: `APP_RNM/internal/services/tolak.go:PeriksaKlaimBernomor`; uji `TestKlaimBelumBernomorTidakDapatDitolak`
- [x] Nilai status hasil penolakan Admin **sama** dengan hasil penolakan Komite; tidak ada nilai
      khusus yang membedakan keduanya. — bukti: uji `TestNilaiTolakSamaDenganKomite`
- [x] Pencerminan ke tingkat `PremiumListDetail` terjadi bersamaan, bukan menyusul. *(AC 7 spec)* — bukti: `APP_RNM/internal/services/statusbaris.go:Status.ubah` (baris dan peserta dalam satu `DalamTransaksi`); uji `TestUbahStatusMencerminkanTigaTingkat` (uji db)
- [x] ⭐ **BARU menurut XML 26-09-2026** — penolakan juga **mencabut penanda dipilih** peserta:
      `PremiumListDetail(idx).IsCheck = "false"` ditulis dalam `Property-Set` yang **sama** dengan
      kedua `STS_REJECT`-nya, jadi ketiganya terjadi bersamaan. Peserta yang barisnya dibatalkan
      berhenti terhitung "dipilih untuk diklaim" dan dapat dipilih ulang dengan baris pengganti.
      *(`[terverifikasi]` `RejectOSClaimLife_Act.xml` pecahan 517–518)* — bukti: `APP_RNM/internal/services/tolak.go:Status.Tolak` → `Status.ubah` (`cabutPenanda`) → `KlaimLife.CabutPenandaDipilih` di transaksi yang sama; uji `TestTolakMencabutPenandaDipilihDiTransaksiYangSama` (uji db)
- [x] ⚠️ **Diselaraskan 2026-09-16:** penolakan langsung oleh `ReasLifeAdmin` menulis `STS_REJECT`
      = **`2`** sebagai **nilai sebenarnya menurut aksi** — bukan nilai yang di-hardcode seperti di
      Pega. Test yang menemukan nilai di-hardcode **gagal**. *(AC 42 spec; penyimpangan sadar 4)* — bukti: `APP_RNM/internal/services/statusbaris.go:Transisi` (`models.KodeDitolak`); uji `TestKodeStatusLiteralHanyaDiModels`, `TestNilaiTolakSamaDenganKomite`

## Catatan

`[keputusan work owner 2026-09-14]` Bahwa penolakan Admin membatalkan **baris saja** — bukan klaim —
berasal dari work owner. `[terverifikasi]` Korpus **tidak membedakan** kedua sumber penolakan: kedua
rule menulis `2` dalam bentuk yang sama persis. Perbedaan itu memang tidak seharusnya ada.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Pembacaan ulang XML — 26 September 2026 malam

XML yang sudah dicatat tiket 03 dan 04 tidak dibaca ulang. Yang **baru** dibaca giliran ini:
`RDBList/UpdateOsAkseptasiClaimLife_sql.xml` *(11.968 byte → 339 baris pecahan)* dan
`Section/RejectOSClaimLife_Sec.xml` *(127.445 byte → 4.145 baris)*.

| Fakta `[terverifikasi]` | Bukti | Akibat |
| --- | --- | --- |
| `RejectOSClaimLife_Act` langkah 2 menulis **tiga** hal sekaligus: `.STS_REJECT = 2`, `PremiumListDetail(idx).STS_REJECT = 2`, dan **`PremiumListDetail(idx).IsCheck = "false"`** | pecahan 443–518 *(dicatat tiket 04)* | ⭐ **AC baru** — pencabutan penanda dipilih; tiket ini semula tidak menyebutnya |
| Gerbang tombol: `pyPosition=='ReasLifeAdmin' && CLAIM_NO !='' && .STS_REJECT == 0` | `AdjustmentDetail_Section.xml` 15399 | tiga syarat, dan yang ketiga diukur **per baris** |
| Langkah 5 memanggil `UpdateOsAkseptasiClaimLife_sql`, yang **`INSERT`** ke `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` dengan **55 kolom** — termasuk `POLICY_HOLDER` dan `NAME_OF_INSURED` | pecahan 79–137 | ⛔ jalur datar warisan **tidak** dikerjakan di sini; lihat di bawah |

### ⛔ Jalur datar warisan — tidak dikerjakan, dan sebabnya dua

⚠️ **Diralat sesudah review — alasannya SATU, bukan dua.** Ronde pertama juga menyebut nama orang
*(`POLICY_HOLDER`, `NAME_OF_INSURED`)* sebagai penghalang; itu **bukan** penghalang, sebab brief sudah
menjawabnya: *"bila baris warisan menuntutnya, ambil dari sumber warisan saat menulis baris warisan
saja, dan catat"*. Yang benar-benar menahan hanya ini: **pemetaan** baris adjustment kita ke baris
datar warisan belum ada — tabel warisan berkunci `CASEID` + `ID`-nya sendiri, dan baris adjustment
tidak menyimpan pengenal itu, sedangkan langkah 5 rule itu **meng-INSERT baris baru**, bukan
mem-`UPDATE` yang ada. Membuat pemetaannya pekerjaan **tiket 13**. `[terbuka — tiket 13]`.
Pembacanya (`AmbilBarisLama`) sudah ada; penulisnya belum, dan tidak dikarang.

---

## Implementasi — 26 September 2026 malam (tiket 05)

**Status: `claimed`** — **8 dari 8 AC tertutup** *(7 asli + 1 baru menurut XML)*. Titik tetap
`d6499ff`. Angka verifikasi di bab hasil review.

⭐ Pintu tiket ini **juga menutup dua AC sisa tiket 04**: pencerminan status ke peserta dan ke header
baru benar-benar terjadi ketika ada yang memanggilnya.

| Berkas | Isi |
| --- | --- |
| `services/tolak.go` *(baru)* | `Status.Tolak`, `PeranRejectOutstanding`, `PeriksaKlaimBernomor` |
| `services/statusbaris.go` | `ubah` dengan satu titik sisip **di dalam** transaksinya |
| `repository/klaimlife.go` | `CabutPenandaDipilih` |
| `handlers/tolak.go` *(baru)* + rute | `POST …/adjustment/{adjId}/tolak` → 403 · 409 · 422 · 501 · 204 |
| `frontend/` | `tolakBarisAdjustment` + tombol pada baris, bergerbang dua syarat yang dapat diperiksa layar |

**Gerbang peran dipasang di sini, bukan ditunggu tiket 07.** AC tiket ini menuntut penolakan oleh
peran lain ditolak, dan gerbang yang tidak ada tidak dapat diuji. `WajibPeran` — yang sampai kini nol
pemanggil — akhirnya punya pemanggilnya. Tiket 07 kelak **menggeneralisasi**, bukan memperkenalkan.

**Penjaga aturan bisnis, dibuktikan dapat gagal:** gerbang peran dilepas → 1 test merah · syarat klaim
bernomor dilepas → 1 test merah.

⭐ **Satu utang yang tidak jadi dibawa.** Ronde pertama mencabut `IS_CHECK` di transaksi **terpisah**
sesudah `Ubah`, dan saya menandainya "utang" di komentar. Itu persis cacat yang tiket 04 tolak untuk
pencerminan statusnya sendiri: peserta dapat tertinggal masih "dipilih" padahal barisnya sudah batal.
`Ubah` karena itu diberi satu titik sisip di dalam transaksinya, dan ketiganya kini jadi bersama.

### ⚠️ Akibat butir **o** yang harus dinyatakan

Sampai penomoran diputuskan, `PenomorBelumDiputuskan` membuat **tidak ada** klaim bernomor. Jalur HTTP
nyata karena itu menjawab **422** untuk setiap permintaan tolak. Itu keadaan yang **benar** — gerbang
`CLAIM_NO != ''` memang menutupnya — bukan kerusakan. Uji `services` memakai pelaku dan nilai
langsung, sehingga aturannya tetap teruji tanpa menunggu **o**.

### ⛔ Yang terbuka

| Butir | Pemilik |
| --- | --- |
| Tulisan ke baris datar warisan *(55 kolom, memuat nama orang; pemetaan baris belum ada)* | tiket 13 |
| Jejak audit — `JejakBelumDiputuskan` masih menggagalkan jalur nyata dengan 501 | butir **am**, work owner |
| Test `db` ujung-ke-ujung untuk pintu tolak | Oracle *(G1)* |

### Hasil `/code-review` atas titik tetap `d6499ff`

⛔ **Temuan terberat: dua AC inti tiket ini TIDAK DAPAT BERJALAN.** Kaitan `sisip` diletakkan
**sesudah** perekaman jejak, dan jejak bawaan (`JejakBelumDiputuskan`) selalu gagal — sehingga di
jalur nyata mana pun pencabutan `IS_CHECK` tidak pernah tercapai, dan nol test menjangkaunya.
Bentuknya benar; pengirimannya tidak terbukti. Diperbaiki: bendera bernama, dijalankan **bersama**
tulisan yang lain sebelum jejak, dan dua test `db` baru membuktikannya.

| # | Sumbu | Temuan | Tindakan |
| ---: | --- | --- | --- |
| 1 | Spec | `sisip` sesudah `Rekam` yang selalu gagal → AC pencabutan `IS_CHECK` dan AC pencerminan tidak terbukti | ✅ urutan dibalik; `TestTolakMencabutPenandaDipilihDiTransaksiYangSama` + `TestTolakYangGagalTidakMencabutPenanda` |
| 2 | Spec | ⛔ **`STS_REJECT` peserta DITULIS tetapi tidak pernah DIBACA** — AC pencerminan tiket 04 tidak dapat dibuktikan, dan separuh pencerminan dapat mati tanpa satu pun test gagal | ✅ `Peserta.KodeStatus` + kolom di daftar baca/tulis. Cacat yang persis sama pernah terjadi pada `CLAIM_RETRO` di tiket 14 |
| 3 | Spec | test `db` tiket 04 membaca ulang peserta lalu **tidak memeriksa apa pun** — pembacaan yang menenangkan tanpa membuktikan | ✅ pernyataannya dipasang |
| 4 | Spec | alasan menunda jalur datar warisan **terlalu lebar**: brief sudah menjawab soal nama orang *("ambil dari sumber warisan saat menulis baris warisan saja")* | ✅ alasan dipersempit menjadi satu: **pemetaan** baris adjustment ↔ baris datar (`CASEID` + `ID`) belum ada, dan itu tiket 13 |
| 5 | Standards | satu sentinel `ErrTanpaWewenang` berarti **dua** hal, dan dua handler menerjemahkannya ke dua kode HTTP berbeda | ✅ dipisah menjadi `ErrTanpaIdentitas` (401) dan `ErrTanpaWewenang` (403); ketiga handler kini seragam |
| 6 | Standards | pesan 403 **berbohong**: pemanggil ber-`X-Peran` tanpa `X-Pelaku` diberi tahu "bukan ReasLifeAdmin" | ✅ identitas diperiksa **sebelum** peran, dan sebelum apa pun dibaca |
| 7 | Standards | **ADR-U-0004 sudah dicabut** *(`status: superseded`, digantikan ADR-0013)* tetapi dikutip di tiga berkas | ✅ ketiganya diralat |
| 8 | Standards | gerbang layar memakai `kodeStatus === '0'` — salinan ketiga kode mentah, pada medan yang kontraknya sendiri tandai "tidak ditampilkan" | ✅ `STATUS_OUTSTANDING`, kata yang datang dari `models.StatusBaris.String()` |
| 9 | Standards | `sisip` sebagai kaitan fungsi: namanya menyebut mekanisme bukan akibat, parameternya membayangi yang sudah ada, dan satu pemanggil | ✅ bendera `cabutPenanda bool` |

**Diterima tetapi belum dikerjakan** *(dicatat, bukan didiamkan)*: `pemilikBaris` dan `ubah`
membaca klaim yang sama dua kali; `PeriksaKlaimBernomor` diperiksa di luar transaksi tanpa penjaga
`WHERE`; blok `Exec`+`RowsAffected`+`"menyentuh %d baris"` kini salinan ketiga.

⚠️ Akibat yang harus dinyatakan: sampai butir **o**, `klaim.nomorKlaim` selalu kosong, sehingga
**tombolnya tidak pernah tampil** dan cabang 422 di layar tidak terjangkau. Itu benar — gerbang
`CLAIM_NO != ''` memang menutupnya — dan disebut di sini supaya tidak terbaca sebagai fitur mati.

**Verifikasi sesudah perbaikan:** vet · vet db · gofmt nol · build · **172 PASS · 0 FAIL** ·
**28 SKIP** · `tsc` · 5 JS · 88 modul.
