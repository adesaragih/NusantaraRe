# 15: Hapus klaim — popup konfirmasi, kaskade tiga tingkat, dan baris work

**Status:** claimed

**Blocked by:** 14 (skema relasional klaim — PREFACTOR), 03 (peserta, adjustment, dan dokumen harus
ada agar dapat dihitung dan dihapus)

> Tiket ini menutup **AC 47** — yang sebelumnya tidak dirujuk tiket mana pun — dan **sisi perilaku**
> AC 48. Tiket 14 hanya membuat *constraint* `ON DELETE CASCADE`-nya; yang dilihat pengguna —
> peringatan berisi jumlah, dan pembatalan yang benar-benar tidak mengubah apa pun — belum ada
> pemiliknya.

## Hasil & nilai pengguna

Sebagai **`ReasLifeAdmin`**, saya dapat menghapus sebuah klaim **beserta seluruh isinya** —
peserta, seluruh baris adjustment, seluruh spreading dan spreading retro, dan seluruh dokumen —
tetapi **tidak sebelum diberi peringatan berisi jumlah baris yang akan ikut terhapus**, supaya saya
dapat membatalkan ketika angkanya tidak sesuai dugaan. *(User story 54 di spec)*

⚠️ **RALAT 2026-09-18** — kata **"polis, marketing"** **DICABUT** dari daftar di atas: kedua
tabelnya **dihapus**. Data polis dan marketing **dibaca hidup** dari tabel polis dan **bukan milik
klaim**, jadi menghapus klaim **tidak menyentuhnya**.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Hapus berkaskade tiga tingkat + baris `T_WORK_CLAIM`, **dalam satu transaksi** |
| `internal/services` | Hitung jumlah baris terdampak sebelum menghapus; orkestrasi |
| `internal/handlers` | Endpoint pratinjau dampak + endpoint hapus |
| `frontend/` | Popup konfirmasi Ya/Batal dengan rincian jumlah per jenis |

## Bentuk yang dihapus — spec §2b

⚠️ **RALAT 2026-09-18** `[keputusan work owner]` — **`T_CLAIM_POLICY` dan `T_CLAIM_MARKETING`
DIHAPUS**, bukan diganti nama, jadi keduanya **tidak ada lagi untuk ikut terhapus**. ⛔ Jangan
membuat `T_CLAIMLF_POLICY`/`T_CLAIMLF_MARKETING`. Data polis dan marketing **dibaca hidup** dari
tabel polis — menghapus klaim **tidak boleh menyentuhnya sama sekali**.

```
T_GENERAL_CLAIM                       ← yang dihapus pengguna
  └─ T_CLAIMLF_PREMIUMLIST_DETAIL  1:N   ikut
        ├─ T_CLAIMLF_ADJUSTMENT     1:N   ikut  ⬅ CUCU
        │     └─ T_CLAIMLF_ADJUSTMENT_SPREADING        1:N  ikut  ⬅ CICIT
        │            └─ T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO  1:N  ikut  ⬅ CICIT-CUCU
        └─ DOCUMENT_CLAIM         1:N   ikut  ⬅ CUCU

T_WORK_CLAIM                      ← baris work klaim life itu ikut terhapus
```

⚠️ **RALAT 2026-09-18 — klaim "tiga tingkat" di tiket ini SUDAH SALAH SEJAK 2026-09-16, bukan
karena ralat hari ini.** Dua tabel spreading ditemukan pada verifikasi 2026-09-16 dan **tidak pernah
masuk** ke pohon di tiket ini. Kaskade sebenarnya menyentuh **lima tingkat** — klaim → peserta →
adjustment → spreading → spreading retro — persis seperti **AC 60 spec** dan tiket **14**. Uji wajib
memeriksa **sampai cicit-cucu** (`T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`), bukan hanya dua tingkat
cucu. Pohon yang berlaku: `spec.md` §2b RALAT D.

## ADR terkait

**ADR-0007** (jejak audit), **ADR-0015** (kegagalan ditangani eksplisit, tidak ditelan).

## Acceptance criteria

- [ ] ⚠️ Menghapus klaim **mengkaskade** ke peserta, **seluruh baris adjustment**, **seluruh
      spreading**, **seluruh spreading retro**, dan **seluruh dokumen**. **REVISI 2026-09-18:** kata
      **"polis, marketing"** **DICABUT** — kedua tabelnya dihapus. Test wajib memeriksa **setiap
      tingkat** secara terpisah, **sampai cicit-cucu**. *(AC 48 spec; penyimpangan sadar 8)*
- [ ] ⚠️ Menghapus klaim **tidak menyentuh tabel polis maupun marketing** — data itu **dibaca
      hidup** dari `T_PREMIUM_LIST` dkk dan **bukan milik klaim**. Test yang menemukan penghapusan
      menyentuh tabel polis **gagal**. *(REVISI 2026-09-18; `[keputusan work owner]`)*
- [x] ⚠️ Penghapusan **didahului popup konfirmasi Ya/Batal** yang menyebut **jumlah baris tiap
      jenis** yang akan ikut terhapus. *(AC 48 spec)*
- [x] ⚠️ Memilih **Batal** **tidak mengubah apa pun** — tidak ada baris terhapus, tidak ada status
      berubah, tidak ada jejak audit penghapusan. *(AC 48 spec)*
- [x] Jumlah yang ditampilkan popup **sama persis** dengan jumlah yang benar-benar terhapus —
      dihitung dari data, bukan dari perkiraan.
- [ ] ⚠️ Menghapus klaim life **menghapus juga baris `T_WORK_CLAIM`**-nya; tidak ada keadaan tangga
      yang tertinggal tanpa klaim. *(AC 47 spec; penyimpangan sadar 6)*
- [ ] Seluruh penghapusan berjalan dalam **satu transaksi**: kegagalan di tingkat mana pun
      **membatalkan seluruhnya**, dan klaim tetap utuh. *(AC 49 spec)*
- [ ] Penghapusan mencatat **jejak audit** — siapa, kapan, dan berapa baris tiap jenis.
      *(**ADR-0007**)*
- [ ] Penghapusan yang gagal menghasilkan kegagalan **terang-terangan**, bukan sebagian terhapus
      diam-diam. *(**ADR-0015**)*
- [x] Menghapus klaim **tidak menyentuh** `M_LIFE_PREMIUM_DETAIL` — ia hanya dibaca sebagai sumber
      snapshot peserta.
- [ ] ⚠️ **Perlakuan `OS_AKSEPTASI_KLAIM_LIFE` saat klaim dihapus** ditetapkan eksplisit. Ia **tetap
      ditulis** saat klaim disimpan (koreksi 2026-09-16, AC 32 spec), sehingga menghapus klaim
      menimbulkan pertanyaan: barisnya ikut dihapus, atau ditinggal karena hilir sudah membacanya?
      `[terbuka]` — **keputusan work owner**, **jangan tebak**. Tiket ini **tidak dinyatakan selesai**
      sebelum jawabannya ada.
- [x] Menghapus klaim **tidak menghapus** peserta di premium list sumbernya — yang terhapus hanya
      **snapshot** milik klaim itu.

## Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka — Non-Life]` **tidak memblokir:** relasi & cascade formal `T_WORK_CLAIM` ditetapkan saat
konteks Non-Life digarap. Untuk Claim Life yang mengikat hanya AC di atas: klaim life dihapus →
baris work-nya ikut.

## Catatan

⚠️ **Mengapa kaskade tidak cukup diuji lewat constraint basis data.** Tiket 14 memasang
`ON DELETE CASCADE`; itu menjamin *baris* hilang, bukan bahwa **angka di popup benar** dan bahwa
**Batal benar-benar membatalkan**. Keduanya perilaku layanan dan layar — dan keduanya tempat bug
biasanya muncul.

⚠️ **Pola yang diikuti.** Kaskade + popup konfirmasi sudah ditetapkan di Master Contract Retro Life
(tiket 09) dan Treaty Contract Out (tiket 10). Bedanya di sini: **tiga tingkat**, dan ada **satu
tabel di luar pohon** (`T_WORK_CLAIM`) yang ikut.

## Seam & perintah verifikasi

**Seam: API HTTP Claim — Life** terhadap **skema uji Oracle nyata** — kaskade tiga tingkat dan
keutuhan setelah Batal **hanya berperilaku benar pada basis data sungguhan**; memalsukannya berarti
tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Pembacaan ulang XML — 26 September 2026 malam

Tiket ini penyimpangan sadar berbentuk relasional; yang dicari di korpus adalah rule yang menghapus
**kasus**. Hasilnya ditulis walau nol, seperti brief minta.

| Yang dicari | Cara | Hasil |
| --- | --- | --- |
| rule penghapus kasus klaim | `grep -rli 'Obj-Delete\|pxDelete\|OBJ-DELETE'` atas **136** berkas `Claim Life` | **NOL**. Yang ditemukan hanya penghapusan **dokumen** (`DeleteDocument_Act`, `DeleteGoogleStorage_Act`, `ConfirmDeleteAttachment`) dan `DeletePesertaClaimLife` |
| `DeletePesertaClaimLife` | dipecah, 44.082 byte → 995 baris; tiga langkah dibaca utuh | ⭐ ia **bukan** kaskade: ia mengulang peserta, menyetel `.IndexPremiumList = Local.IndexPremium` pada tiap `.AdjustmentList`, lalu `Obj-Save`. Artinya **penomoran ulang subscript** sesudah satu peserta dikeluarkan dari daftar tertanam |

⭐ **Akibatnya bagi kita: nol kode.** Subscript Pega digantikan `PREMIUM_LIST_DETAIL_ID`, yang tidak
berubah ketika saudaranya hilang — jadi langkah penomoran ulang itu **tidak punya padanan**, dan
kaskadenya milik `ON DELETE CASCADE` (tiket 14), bukan milik kode.

---

## Implementasi — 26 September 2026 malam (tiket 15)

**Status: `claimed`** — **5 dari 12 AC tertutup.** Titik tetap `f610178`.
Angka verifikasi di bab hasil review.

### ⛔ Temuan yang mengubah seluruh bentuk tiket ini

**ADR-U-0031 melarang hapus fisik di jalur pengguna.** *"Penghapusan adalah penanda dan nilai balik,
bukan hapus fisik … Nol perintah hapus fisik pada jalur pengguna di lapisan mana pun … Test yang
berhasil menghapus baris secara fisik lewat jalur pengguna **gagal**."* `[keputusan work owner,
23-09-2026]`

Teks tiket ini ditulis **18-09-2026** dan berbicara tentang kaskade **fisik** lima tingkat. ADR itu
**lima hari lebih muda** dan memutuskan sebaliknya. Penjaga statik di `batasanpemakaian_test.go` pun
sudah menyebut tiket ini sebagai pelaksananya — jadi yang diminta memang bentuk **logis**.

Yang menahan pelaksanaannya: **kolom penandanya belum ada**, dan menambahkannya adalah langkah
migrasi baru yang hanya lahir dari paket keputusan yang masih `[USULAN]`. Menebak nama kolom dan
artinya berarti mengarang bentuk penyimpanan.

⛔ Karena itu `Penghapusan.Hapus` **gagal terang** dengan `ErrHapusFisikDilarang` → **501**, sedangkan
`Dampak` bekerja **penuh**. Bagian yang benar-benar dilihat pengguna — peringatan berisi angka, dan
Batal yang benar-benar membatalkan — dapat dibangun dan diuji sekarang.

| Berkas | Isi |
| --- | --- |
| `models/dampak.go` *(baru)* | `DampakHapus` bermedan per jenis, `Total`, `Kosong`, `String` |
| `repository/pohonklaim.go` | `Dampak` — tujuh `COUNT` terpisah, dihitung dari data |
| `services/hapus.go` *(baru)* | `Penghapusan.Dampak` *(hanya membaca)*, `Hapus` *(gagal terang)*, gerbang Komite |
| `handlers/hapus.go` *(baru)* + rute | `GET …/dampak-hapus`, `DELETE /api/klaim-life/{id}` |
| `frontend/` | popup Ya/Batal berisi rincian per jenis |
| `repository/batasanpemakaian_test.go` | penjaga hapus fisik **dipertajam** |

⭐ **Penjaga dipertajam, bukan dilonggarkan.** Pola lama mencocokkan `.Hapus(` apa pun, sehingga
`services.Penghapusan.Hapus` — yang justru **menolak** menghapus — ikut tertuduh. Pola baru menyasar
tanda tangan repository-nya (`.Hapus(ctx, tx`), dan dibuktikan masih menangkap pemanggilan
sungguhan yang disisipkan sengaja.

### ⚠️ Keputusan tiket 14 yang ternyata belum disahkan

`repository.PohonKlaim.Hapus` **menghapus baris datar warisan** (`DELETE … WHERE CASEID`), padahal
tiket ini menandai nasib `OS_AKSEPTASI_KLAIM_LIFE` sebagai `[terbuka — jangan tebak]`. Kode itu tidak
diubah — ia hanya dipanggil test kaskade — tetapi cacahnya kini **ditampilkan di popup** sebagai
barisnya sendiri, supaya keputusan yang belum diambil terlihat oleh yang menekan tombol.

### ⛔ AC yang BELUM tertutup — 7, dan semuanya satu sebab

Kaskade, penghapusan baris work, satu transaksi, jejak audit penghapusan, kegagalan terang-terangan,
dan tidak-menyentuh tabel polis: seluruhnya menunggu **bentuk logis** ditetapkan — kolom penanda,
dan apakah "terhapus" berarti baris hilang atau baris bertanda. `[terbuka — work owner, ADR-U-0031]`
Ditambah `[terbuka]` nasib `OS_AKSEPTASI_KLAIM_LIFE` yang tiket ini sendiri syaratkan.

### Lanjut dari sini

Bila konteks sesi harus diringkas, titik sambungnya: **tiket 15, sesudah `/code-review` dijalankan
atas titik tetap `f610178`, sebelum perbaikan temuan dan commit.** Yang sudah hijau: 175 PASS · 0
FAIL · 28 SKIP · `tsc` · 5 JS · 88 modul. Sesudah commit tiket 15, urutan berlanjut ke **07 → 08 →
09 → 10 → 12 → 11 → 13**.

Bekal yang sudah dibaca dan tidak perlu diulang: peran per tahap dari `Flow/Register_Flow.xml`
`[terverifikasi]` — `Assignment1` dan `Assignment2` → `ReasLifeAdmin`, `Assignment3` →
`ReasLifeMedicalAdvisor`, `Decision1` → `ReasLifeSPV`, `Decision3` → `ReasLifeMedicalAdvisor`
*(tiket 07 dan 08)*. Bukti butir **al**: kelas `Data-DiagnoseLife` dipakai enam rule, termasuk
`Diagnose_Harness.xml` dan `Diagnose_Section.xml` — layar diagnosa memang ada *(tiket 08)*.

### Hasil `/code-review` atas titik tetap `f610178`

**Kedua sumbu mengesahkan penolakan 501 itu benar, bukan penghindaran**: brief modul §1.2-c memang
mewajibkan executor menulis kedua sisi dan menandai `[terbuka — work owner]` ketika penyimpangan
sadar bertabrakan dengan keputusan lain. Tetapi keduanya juga menemukan enam cacat nyata.

| # | Temuan | Tindakan |
| ---: | --- | --- |
| 1 | ⛔ **Popup MENDAHULUI keputusan yang terbuka**: "Baris datar warisan" berdiri di bawah judul *"yang akan ikut terhapus"* dan ikut dijumlahkan — padahal AC 11 menandai nasibnya `[terbuka — jangan tebak]` | ✅ dikeluarkan dari daftar **dan** dari `Total()`; `TotalTermasukWarisan()` terpisah; di layar ia paragraf tersendiri bertuliskan *nasibnya belum diputuskan* |
| 2 | ⛔ **`Dampak` mencacah dengan kunci yang salah** — saya mengirim `klaimID` sebagai `caseID` dengan alasan "berbagi kunci utama". Butir **ae1** mengisinya begitu untuk klaim yang sistem ini buat; itu keputusan **pengisian**, bukan jaminan bentuk, dan klaim yang dimigrasikan membawa `CASE_ID` warisannya sendiri | ✅ `CaseIDKlaim` membacanya dari baris work |
| 3 | ⛔ **`Dampak` melewatkan baris header** — `Total()` kurang satu dari yang sebenarnya | ✅ dicacah |
| 4 | ⛔ **Penjaga hapus fisik lolos** oleh nama variabel lain (`trx`, `h.tx`) atau baris yang terpotong, padahal komentarnya mengaku menangkap *"setiap pemanggilan"* | ✅ yang diperbaiki **namanya**, bukan polanya: `PohonKlaim.Hapus` → **`HapusFisik`**, sehingga penjaga mencocokkan nama dan tidak dapat dielakkan bentuk pemanggilan |
| 5 | Test kaskade belum diperluas ke lima tingkat, dan AC *"jumlah popup sama persis"* tercentang tanpa satu pun perbandingan | ✅ `TestHapusMengkaskadeSampaiCicit` kini mencacah **sebelum** menghapus, menolak fixture yang kosong, lalu memeriksa **tiap tingkat** nol sesudahnya |
| 6 | Perkabelan yang menjanjikan yang tidak dilakukannya: `Hapus(…, saat)` tidak pernah membaca `saat`; `jejak`/`DenganJejak` tidak pernah dibaca; `CatatanJejak.Catatan` tanpa penulis; cabang 501 jejak tak terjangkau | ✅ keempatnya dibuang. Jejak penghapusan lahir bersama bentuk logisnya, bukan sebelum itu |

Ditambah dua penjaga baru: `TestJalurHapusTidakMenyentuhTabelSumber` *(AC `M_LIFE_PREMIUM_DETAIL`,
yang sebelumnya benar secara kebetulan)* dan dua test JS untuk bentuk pratinjau.

⚠️ Satu ketidaktepatan angka di bab pembacaan XML: `grep` `Obj-Delete|pxDelete` sebenarnya mengenai
**delapan** berkas, bukan empat — enam sisanya pesan peringatan bawaan Pega, bukan langkah. Nol rule
penghapus kasus tetap benar.

**Verifikasi sesudah perbaikan:** vet · vet db · gofmt nol · build · **176 PASS · 0 FAIL** ·
**28 SKIP** · `tsc` · **7** test JS *(dari 5)* · 88 modul.
