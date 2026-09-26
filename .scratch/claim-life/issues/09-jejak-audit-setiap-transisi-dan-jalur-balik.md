# 09: Jejak audit — siapa + kapan untuk setiap transisi dan setiap jalur balik

**Status:** claimed

**Blocked by:** 08 (tahap + jalur balik) — seluruh transisi harus ada dulu untuk dapat direkam

## Hasil & nilai pengguna

Sebagai **auditor**, saya dapat mengetahui **siapa** dan **kapan** untuk setiap transisi status dan
setiap pengembalian kasus — sehingga setiap keputusan dapat dipertanggungjawabkan, dan pengembalian
kasus dapat ditelusuri. Sebagai **ReasLifeAdmin**, saya melihat siapa yang mengembalikan kasus
kepada saya dan kapan, sehingga saya tahu apa yang diminta.
*(User story 10, 29, 30 di spec)*

**Ini penyimpangan sadar dari sistem lama — sebuah perbaikan, bukan paritas.**

## Area codebase

`internal/models` (entri jejak audit), `internal/repository` (penyimpanan jejak),
`internal/services` (perekaman pada setiap transisi — satu tempat, bukan tersebar),
`internal/handlers` (identitas pelaku), `frontend/` (tampilan riwayat pada klaim).

Jejak direkam **per baris `AdjustmentList`**, bukan per klaim — karena unit statusnya baris
(**ADR-0011**).

## Rule Pega sumber

| Rule | Identitas | Keadaan sekarang |
| --- | --- | --- |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` **INSERT** (bukan update, walau namanya begitu) ke `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`; dari seluruh kolomnya yang **berciri audit** hanya `CREATEOPNAME` + empat tanggal: `ACCEPTATION_DATE`, `CONFIRMATION_DATE`, `CLAIM_RECEIVED_DATE`, `COMPLETE_DATE`. Memuat `COMMIT;` - OQ-013 |
| `Claim Life/When/IsSendtoAdmin.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` / `RULE-OBJ-WHEN` | `[terverifikasi]` hanya menyimpan nilai `1` — **tanpa pelaku, tanpa waktu** |
| `Claim Life/When/IsSendtoMedical.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOMEDICAL` / `RULE-OBJ-WHEN` | idem |
| `Claim Life/RDBList/InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` / `RNM!INSERTLOGSERVICECLAIM` / `RULE-CONNECT-SQL` | `[terverifikasi]` `INSERT INTO pooldata.monitoring_klaim_log` — **log layanan**, bukan jejak keputusan |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml`, `SendEmailKlaimLF.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `…` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` memakai `OperatorID.pyUserIdentifier` / `pyUserName` sebagai data |

## ADR terkait

**ADR-0007** (jejak audit setiap transisi dan setiap jalur balik — penyimpangan sadar),
**ADR-0002** (pelaku diidentifikasi lewat akun berperan, bukan nama ter-hardcode),
**ADR-0011** (jejak per baris).

## Acceptance criteria

- [ ] Setiap transisi status baris menghasilkan catatan berisi **pelaku dan waktu**. *(AC 16 spec)*
- [ ] Setiap pengembalian (`SendtoAdmin`, `SendtoMedical`) menghasilkan catatan berisi **pelaku dan
      waktu**. *(AC 17 spec)*
- [ ] Perubahan nilai `Type` menghasilkan catatan berisi pelaku dan waktu. *(AC 18 spec)* — `Type`
      menyentuh keamanan, bukan sekadar data (**ADR-0012**).
- [ ] Jejak melekat pada **baris** yang bersangkutan, dan riwayat satu klaim dapat dibaca utuh
      lintas seluruh barisnya.
- [ ] Rekam akseptasi lama tetap ditulis sebagaimana adanya — kontrak dengan Komite tidak berubah
      karena tiket ini.
- [x] Pelaku dicatat sebagai identitas akun; **tidak ada nama orang ter-hardcode**.

## Catatan

`[keputusan work owner 2026-09-14]` Jejak lama **tidak dapat direkonstruksi ke belakang** — data
sebelum cutover hanya punya `CREATEOPNAME` + empat tanggal. Riwayat transisi lengkap hanya ada untuk
kejadian **setelah** cutover. Konsekuensi ini disadari dan diterima.

`[terbuka]` **OQ-013** (pemilik **DBA**) — `COMMIT` berada di dalam blok PL/SQL, sehingga
atomisitas "tulis akseptasi + tulis jejak audit" belum dapat dipastikan. **Tidak memblokir** tiket
ini; memblokir jaminan atomisitasnya.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Pembacaan ulang XML — 27 September 2026

XML tiket ini **sudah dibaca dan dicatat** di tiket 04 *(sensus penulis `STS_REJECT`)* dan tiket 08
*(`IsSendtoAdmin`, `IsSendtoMedical`, `Register_Flow`)*; sesuai brief lanjutan 2 §1.4-c keduanya
dirujuk, bukan dibaca ulang.

**Sensus transisi yang WAJIB direkam** — diturunkan dari kedua bacaan itu:

| Transisi | Penulisnya | Sudah merekam? |
| --- | --- | --- |
| baris → Outstanding `0` | `SaveOutStandingLife_Act` ⇄ `TandaiOutstanding` | ⛔ belum — jalurnya lewat `Pendaftaran.Daftar`, yang belum merekam |
| baris → Ditolak `2` | `RejectOSClaimLife_Act` ⇄ `Status.Tolak` | ✅ lewat `Status.ubah` |
| baris → Aksep `1` | `SaveAdjustment_Act`, `KomitePostAdjustment` | ⛔ modul Komite — tiket 11 |
| peserta + header dicerminkan | `RejectOSClaimLife_Act`, `serviceInsertArasapasClaimLife_act` | ✅ satu transaksi dengan barisnya |
| jalur balik `SendtoAdmin = 1` | `Activity/SendtoAdmin_Act.xml` baris **259** ⇄ `TahapLayanan.Pindah` | ⛔ belum — `Pindah` menulis penandanya, tetapi `Jejak` produksi belum terpasang |
| jalur balik `SendtoAdmin = 0` *(pembersihan)* | `Activity/SendtoAdmin_Act1.xml` baris **282** | ⛔ belum — dan **belum ada padanannya** di sistem ini |
| jalur balik `SendtoMedical = 1` | `Activity/SendtoMedical_Act.xml` baris **260** ⇄ `TahapLayanan.Pindah` | ⛔ belum — sebab yang sama |
| perubahan `Type` | **dua** rule menulisnya: `Activity/LoadDataPeserta_Act.xml` **726**, `Activity/InsertJsonClaimLife_Act.xml` **745**; dan ia **dapat disunting manusia** — `pxDropdown` di `Section/InputRegisterClaimLife.xml` **9110**, `InputOSClaimLife.xml` **6905**, `InputAkseptasiClaimLife.xml` **6721**, `MedicalCheckClaimLife.xml` **6751** | ⛔ belum — dan ini transisi yang **nyata**, bukan yang mustahil |

### Ralat menurut XML — 27 September 2026

**Sensus dihitung dua cara** *(CLAUDE.md §4a)*, jendela `D:\XML\RNM_BRD\Claim Life\**\*.xml`:
**cara 1** `<PropertiesName>` memuat `Type` → **2 berkas**; **cara 2** `.PolicyDataLife.Type`
dirender kontrol → **4 section**. Keduanya berbeda jenis — penulis vs penyunting — dan justru
selisihnya yang membongkar kekeliruan di bawah.

| Butir | Teks lama | Teks baru | Bukti |
| --- | --- | --- | --- |
| penulis `Type` | *"nol rule menulisnya"*, lalu diralat jadi *"satu rule"* | **dua** rule menulisnya | `LoadDataPeserta_Act.xml` 726; `InsertJsonClaimLife_Act.xml` 745 |
| `Type` diubah manusia? | *"nol jalur mengubahnya"* | **dapat disunting di empat layar** | `Section/InputRegisterClaimLife.xml` 9110 `pyFormat=pxDropdown`, `pyDisabledNew=false`, tanpa `pyReadOnly` — padahal tabel yang sama memakai `<pyReadOnly>true</pyReadOnly>` di 9058 dan 9257 |
| penulis `SendtoAdmin` | `Register_Flow` | `SendtoAdmin_Act.xml` 259 (`"1"`) dan `SendtoAdmin_Act1.xml` 282 (**`"0"`**) | Flow hanya memanggilnya |
| penulis `SendtoMedical` | `Register_Flow` | `SendtoMedical_Act.xml` 260 (`"1"`) | idem |
| `UpdateOsAkseptasiClaimLife_sql` | *"hanya CREATEOPNAME + empat tanggal"* | **INSERT** besar; lima kolom itu yang berciri audit | namanya `Update…`, isinya `INSERT INTO` |

⛔ **Kekeliruan ini milik saya, dua ronde berturut-turut.** Ronde pertama: *"nol rule menulis
`Type`"*. Ronde kedua saya ralat jadi *"satu rule, penyalinan saat muat, jadi tetap tak ada yang
diubah manusia"* — dan **ralat itu sendiri salah**. `Type` dipasang sebagai dropdown yang dapat
disunting di layar Register, OS, Akseptasi, dan Medical Check. Akibatnya AC 18 bukan AC yang
subjeknya mustahil; subjeknya **ada**, dan yang belum ada hanya tabelnya. Saya menutup pertanyaan
dengan mencari pembuktian untuk kesimpulan yang sudah saya pegang, bukan dengan mencari
sanggahannya — dan sensus satu arah membiarkannya lolos dua kali.

⚠️ `RDBList/InsertLogServiceClaim.xml` menulis `pooldata.monitoring_klaim_log` dengan kolom
`IDPEGA, PARAMETER, JN_SERVICE, NO_AKSEPTASI, NO_DLA, STS_MESSAGE, RESPON_MESSAGE` `[terverifikasi]`
— nama layanan, pesan, dan responsnya. Itu **log panggilan layanan**, bukan jejak keputusan, dan
**tidak** dipakai sebagai tabel audit *(brief modul §4-09)*. ⛔ SQL-nya bahkan memuat `COMMIT;`,
yang `PeriksaSQL` tolak di jalur kita (ADR-U-0029).

## Implementasi — 27 September 2026 (tiket 09)

**Status: `claimed`** — **1 dari 6 AC tertutup.** Titik tetap `7c17bbf`.

⛔ **Semula saya tulis 3 dari 6.** AC 16 dan AC 17 saya centang padahal tidak terkirim: satu-satunya
`Jejak` produksi adalah `JejakBelumDiputuskan`, yang selalu galat, sehingga **nol** transisi
menghasilkan catatan. Yang tertutup hanya *"pelaku dicatat sebagai identitas akun"*. Mencentang AC
yang bentuknya sudah ada tetapi isinya belum adalah cara paling halus membuat papan hijau berbohong.

| Angka | Nilai | Perintah audit | Label |
| --- | --- | --- | --- |
| test Go | **189 PASS · 0 FAIL · 28 SKIP** *(dari 187)* | `go test -tags=db ./internal/... -v` lalu cacah awalan `--- PASS` / `--- FAIL` / `--- SKIP` | `[terverifikasi]` |
| vet | bersih, termasuk `-tags=db` | `go vet ./... && go vet -tags=db ./...` | `[terverifikasi]` |
| gofmt | nol berkas | `gofmt -l internal/ cmd/ pkg/` | `[terverifikasi]` |
| test JS | **7 PASS** | `cd frontend && npm test` | `[terverifikasi]` |
| modul frontend | **88** | `cd frontend && npx tsc --noEmit && npm run build` | `[terverifikasi]` |

⛔ **Tabel jejaknya tidak dibuat.** Butir **am** masih `[USULAN]`; langkah migrasi baru hanya lahir
dari paket keputusan yang disahkan. Yang dibangun adalah **bentuknya**: antarmuka `Jejak`
*(tiket 04)*, `CatatanJejak`, dan `JejakBelumDiputuskan` yang **gagal terang** — sehingga transisi
tidak dapat terjadi tanpa terekam, dan yang belum diputuskan terlihat sebagai satu galat yang
menyebut apa yang ditunggu.

| Berkas | Isi |
| --- | --- |
| `services/jejak_statik_test.go` *(baru)* | dua penjaga statik |

**Penjaga, dibuktikan dapat gagal — tiga ronde:**

| Ronde | Celah | Bukti gagalnya |
| --- | --- | --- |
| 1 | per **berkas**: dua fungsi bebas di satu berkas lolos | perekaman dilepas dari `tahap.go` → merah |
| 2 | per **fungsi**, tetapi komentar hanya dibuang di sisi penulis | `tulisTanpaJejakBungkam` + komentar prosa yang menyebut `jejak.Rekam(ctx, tx,` → **hijau padahal salah** |
| 3 | komentar dibuang di **kedua** sisi | kasus yang sama → merah: *"memanggil penulis transisi repository tetapi tidak merekam jejaknya di fungsi yang sama"* |

⛔ Cakupan penjaga kini **dinyatakan, bukan dilebihkan**. Ia memeriksa pemanggil
`PerbaruiStatusBaris` / `CerminkanHeader` / `PerbaruiTahap`; ia **tidak** memeriksa `KodeStatus =`,
sehingga jalur `Pendaftaran.Daftar` di luar jangkauannya. Ronde sebelumnya pesannya berbunyi
*"SETIAP transisi"* sementara tiketnya mengaku ada celah — pengakuan yang hanya ada di prosa dan
tidak ada di penjaganya sama saja dengan tidak mengaku.

### ⛔ AC yang belum tertutup — 3

| AC | Sebab |
| --- | --- |
| *"Jejak melekat pada baris, dan riwayat satu klaim dapat dibaca utuh"* | menuntut **tabelnya** — butir **am**, `[USULAN]` |
| *"Setiap transisi status baris menghasilkan catatan"* **(AC 16)** | `Jejak` produksi **belum terpasang**: satu-satunya yang ada `JejakBelumDiputuskan`, yang selalu galat. Nol transisi menghasilkan catatan hari ini — yang dihasilkan HTTP 501. Menunggu **am** |
| *"Setiap pengembalian menghasilkan catatan"* **(AC 17)** | sebab yang sama, ditambah: **nol handler** memanggil `TahapLayanan.Pindah`, jadi jalur baliknya belum punya pintu masuk |
| *"Perubahan nilai `Type` menghasilkan catatan"* **(AC 18)** | `Type` **dapat** diubah — dropdown tersunting di empat layar, dan `services/pendaftaran.go:177` menulisnya. Yang belum ada tabelnya, bukan subjeknya. Menunggu **am** |
| *"Rekam akseptasi lama tetap ditulis sebagaimana adanya"* | jalur tulis datar warisan milik **tiket 13** *(lihat tiket 05)* |

⚠️ Satu celah yang sensus di atas ungkap dan **tidak** saya tutup diam-diam: `Pendaftaran.Daftar`
menulis status Outstanding lewat `TandaiOutstanding` **tanpa** merekam jejak. Penjaga statik tidak
menangkapnya sebab ia menulis lewat `PohonKlaim.Simpan`, bukan lewat `PerbaruiStatusBaris`.
Dinyatakan di sini; penutupannya menunggu **am** bersama tabelnya.

### Hasil /code-review — titik tetap `7c17bbf`

**Standards — 5 temuan, 4 diperbaiki.**

| Temuan | Tindakan |
| --- | --- |
| ⛔ lubang sisa: `buangKomentar` hanya di sisi penulis, tidak di sisi perekam; `pecahPerFungsi` memotong di `^func` sehingga komentar kepala fungsi berikutnya jatuh ke tubuh sebelumnya — satu baris prosa membungkam penjaga | **diperbaiki**, dan dibuktikan merah lagi dengan kasus bungkam |
| cakupan tak jujur: pola tanpa `KodeStatus =`, tetapi pesannya berbunyi *"SETIAP transisi"* | **diperbaiki** — batasnya dinyatakan di komentar pola **dan** pesan galatnya diperkecil |
| §4 rule 9: blok angka tanpa perintah audit dan tanpa label | **diperbaiki** — tabel angka + perintah + label |
| §4a: sensus satu arah | **diperbaiki** — dihitung dua cara; selisihnya yang membongkar kekeliruan `Type` |
| §4 rule 1: sensus menyebut rule tanpa path + baris | **diperbaiki** |
| duplikasi jalan-telusur dengan `wewenang_statik_test.go`; komentar basi *"kelak oleh tiket 09"* | **diterima, tidak diperbaiki** — penyatuannya menunggu `Jejak` nyata (butir **am**); menyatukan dua penjaga sekarang menukar duplikasi dengan abstraksi yang belum tahu bentuknya |

**Spec — 6 temuan, semuanya sahih, semuanya diperbaiki.** Review ini menemukan yang paling mahal:

| Temuan | Tindakan |
| --- | --- |
| ⛔ AC 16 dicentang padahal nol transisi menghasilkan catatan | **dicabut** |
| ⛔ AC 17 dicentang padahal nol handler memanggil `Pindah` | **dicabut** |
| ⛔ sensus `Type` salah — dan **ralat saya atasnya juga salah** | **diralat**, dengan blok ralat bertanggal |
| penulis `Sendto*` dikreditkan ke `Register_Flow`, padahal Activity | **diralat**, dengan nomor baris |
| pembersihan `SendtoAdmin = "0"` hilang dari sensus | **ditambahkan** |
| `UpdateOsAkseptasiClaimLife_sql` terbaca seolah tulisan lima kolom | **diberi kualifikasi**: INSERT besar, lima kolom berciri audit |

⭐ Nilai review di sini bukan menemukan kode yang salah — kodenya benar. Ia menemukan **tiket yang
mengaku lebih banyak daripada yang dikerjakan**, dua kali di butir yang sama.

**lanjut dari sini:** tiket 09 selesai — AC 16/17/18 dan dua AC lain menunggu butir **am**
disahkan work owner. Berikutnya tiket 10.
