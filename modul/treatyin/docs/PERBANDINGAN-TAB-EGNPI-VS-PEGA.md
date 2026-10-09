# Tab EGNPI (non-prop) — aplikasi lawan ekspor Pega

Dibaca 6 Oktober 2026 dari `D:\XML_NURE\Treaty In`. Berkasnya:

```
Section/TreatyInTabsNonProportional.xml   wadah TABBED ke-3 `EGNPI`
Section/DetailEGNPI.xml                   rincian satu baris
Activity/SetAmountConversion.xml          .AmountIDR = .Amount x kurs
Activity/TreatyInEGNPIListValue.xml       tombol `Update EGNPI Value`
Activity/TreatyInNPSetTotal.xml           tombol `Update Total` (type=egnpi)
Activity/TreatyInNonAddItem.xml           tombol `Add` (Type=egnpi)
Activity/TotalEgnpi.xml                   dibaca tab Limits, bukan tab ini
```

---

## 1 · Hasil perbandingan

| | Pega | Aplikasi SEBELUM | Sekarang |
|---|---|---|---|
| Grid 6 kolom | ✅ | grid 8 kolom hanya-baca | ✅ kartu lipat + rincian |
| Tombol `Add` | ✅ | ⛔ **tidak ada** | ✅ |
| Tombol `Delete` | ✅ | ⛔ **tidak ada** | ✅ |
| Panel `Total EGNPI Amount` per mata uang | ✅ | ⛔ **tidak ada** | ✅ |
| `Total Amount in IDR` | ✅ | ⛔ **tidak ada** | ✅ |
| `Total Proportion %` | ✅ | ⛔ **tidak ada** | ✅ |
| Tombol `Update Total` | ✅ | ⛔ **tidak ada** | ✅ |
| Tombol `Update EGNPI Value` | ✅ | ⛔ **tidak ada** | ✅ |
| Rincian baris (`DetailEGNPI`) | ✅ | ⛔ **tidak ada** | ✅ |

⛔ **Akibat yang paling mahal dari keadaan sebelumnya**: tanpa kedua tombol,
`Amount in IDR` dan `Proportion %` nol pernah terisi — padahal tab **Limits**
MEMBACA baris EGNPI lewat `TotalEgnpi`. Jadi tab yang terlihat pasif itu
sebenarnya memasok angka ke tab lain.

---

## 2 · Rumus yang dipindahkan

Seluruhnya di `backend/services/hitung_egnpi.go`, diuji di
`hitung_egnpi_test.go`, dipanggil lewat `POST /api/treaty-in/hitung/egnpi`.

### 2.1 `SetAmountConversion`

```
.AmountIDR = .Amount x kurs(.Currency)
```

⚠️ Kursnya dari kecocokan **TERAKHIR** di `CurrencyList`, bukan pertama —
ekspor menapaki seluruh daftar tanpa `break`, dan setiap kecocokan menimpa
`local.ConvValue`.

### 2.2 `TreatyInNPSetTotal` (`param.type == "egnpi"`)

```
[4] kosongkan TotalEgnpiAmountNP · TotalEgnpiAmount · TotalEgnpiProportion
[5] TotalEgnpiAmount    = Σ .AmountIDR
[6] bila total == 0 → BERHENTI (pagar bagi-nol)
[7] .Proportion         = @divide(.AmountIDR, TotalEgnpiAmount, 20) × 100
    TotalEgnpiProportion = Σ .Proportion
    TotalEgnpiAmountNP   = Σ .Amount PER MATA UANG
```

⛔ Langkah 5 menjumlah `.AmountIDR`, langkah 7 menjumlah `.Amount`. Dua kolom
berbeda di baris yang sama — menyeragamkannya membuat grid per mata uang
memperlihatkan rupiah.

⚠️ Skala **20** bukan hiasan: itulah yang membuat jumlah proporsi jatuh tepat
di 100 untuk nilai yang habis dibagi. Untuk yang tidak habis dibagi ia
berbunyi `99.999999999999999999` — di Pega juga.

### 2.3 `TreatyInNonAddItem` (`Type == "egnpi"`)

Baris baru mewarisi mata uang **baris pertama tab Maximum Retention**
(`TreatyIn.Retention(1)`). Nol retensi berarti baris tanpa mata uang; ia
tidak dikarang menjadi `IDR`.

---

## 3 · ⛔ Temuan yang harus dinyatakan

### 3.1 Kurs cadangan memakai TAHUN BERJALAN, bukan tahun treaty

`SetAmountConversion` langkah 3, bila mata uangnya nol di grid:

```
param.pyReportName   "BrowseTreatyExchangeYearly_RD"
param.pyReportClass  "ASM-FW-GISFW-Int-TREATYEXCHANGEYEARLY"
param.TreatyYear     @substring(@CurrentDateTime(),0,4)
```

⚠️ `@CurrentDateTime()` — **tahun kalender hari ini**, bukan
`TreatyIn.TreatyYear`. Dua kontrak bertahun treaty berbeda akan memakai kurs
yang sama. Itu bunyi ekspornya.

⛔ Cadangan ini **tidak dibangun**: rute `/hitung/` murni menghitung, nol baca
basis data. Bila kursnya nol di grid, `.AmountIDR` **dibiarkan** dan satu
pesan menyebut mata uangnya — mengisinya nol akan terbaca sebagai "kursnya
nol" padahal yang benar "kursnya tidak diketahui".

### 3.2 Cabang `ClassOfBusiness` di `TotalEgnpi` MATI

Langkah `1.4.2` berketerangan *"Copy matching CoB to EgnpiMatch"* dan membawa
`When` `.ClassOfBusiness==Local.cob`, tetapi:

- `pyStepsActivityName` **kosong** — nol metode, jadi nol yang dijalankan;
- `pyStepsPreCondition` = `false` — pra-syaratnya pun dimatikan.

Yang hidup `1.4.2.1`, dan ia mencocokkan **`.TreatyGroup`**, bukan
`ClassOfBusiness`. ⭐ Itu menguatkan pengukuran yang sudah tercatat di
`hitung_limit_np.go`: pencocokan lewat `.TreatyGroup` menjelaskan 5.020 dari
5.069 baris (99,0 %).

### 3.3 Ekspor berselisih dengan dirinya sendiri — `Treaty Group`

| | `pyReadOnly` |
|---|---|
| grid tab EGNPI | *(tidak ada)* → dapat diubah |
| `Section/DetailEGNPI.xml` | `true` |

Yang dipakai bentuk **grid**: `Add` melahirkan baris kosong, dan baris yang
kelompok treaty-nya tidak dapat diisi nol pernah cocok dengan layer mana pun
di `TotalEgnpi`.

### 3.4 ~~`Note` nol pernah dapat diisi~~ — ⛔ RALAT 7 Oktober 2026

**Bacaan itu KELIRU dan ditarik.** `Note` DAN `As At` memang ber-`pyReadOnly
= true`, tetapi sel yang SAMA juga memuat
`pyReadOnlyCondition = TreatyIn.ViewState = 1` — dan **syarat itulah yang
berlaku**. `ViewState 1` = mode lihat, jadi keduanya **aktif di mode Edit**.

⛔ Akibatnya nyata: kedua medan sempat **dikunci mati** di tab EGNPI
aplikasi. Sudah diperbaiki — keduanya kini `readOnly={!bisaUbah}`.

⚠️ Kekeliruannya satu jenis dengan membaca `pyStepsPreCondition` tanpa
`pyStepsBlockName`: membaca SATU tag lalu berhenti, padahal tag kedua di sel
yang sama membatalkannya.

⭐ Aturannya kini terkodekan SEKALI, berikut ujinya atas XML sungguhan:
`D:\XML_NURE\_migration-docslat-baca-ekspor\` (`hanya_baca`, Aturan 3).

### 3.5 Pesan galat total-nol HAMPA di Pega

Langkah 6 `Property-Set-Messages` dengan parameter **kosong**, jadi Pega
memunculkan pesan tanpa teks. Teks `"Error Divide by Zero"` dideklarasikan di
langkah 1 dan tidak pernah dipasang. ⭐ Aplikasi memakai teks yang
dideklarasikan itu — pesan hampa tidak memberi tahu pemakai apa pun.

### 3.6 Pemicu konversi: setiap perubahan lawan kehilangan fokus

Ekspor memanggil `SetAmountConversion` pada **setiap** perubahan `.Currency`
dan `.Amount`. Layar ini memanggilnya saat kotak `Amount` kehilangan fokus
(mata uang tetap seketika, sebab ia pilihan). Satu panggilan jaringan per
ketukan papan ketik bukan kesetiaan, melainkan layar yang tersendat.

---

## 3.7 ⛔ BENTUK RINCIAN DIPERBAIKI 7 Oktober 2026

Atas tangkapan layar Pega pemilik proses, tiga hal yang sebelumnya salah:

| | Sebelumnya | Ekspor / tangkapan layar |
|---|---|---|
| `As At` | kotak teks biasa | **`pxDateTime`** — inputan tanggal berikon kalender |
| `Treaty Group`, `Currency` | autocomplete ketik-bebas | **dropdown**, disamakan dengan `Ceding` dan `Source of Business` atas permintaan pemilik proses |
| `Amount`, `Amount in IDR` | dua medan tunggal | **SEPASANG**: satuan di kiri, nilai di kanan tanpa label sendiri |

⭐ Dropdownnya memakai komponen yang **SAMA** (`DropdownWarisan`), bukan
tiruannya — jadi perilaku ketik-saring, papan tik, dan penandaan nama kembar
ikut apa adanya.

⛔ Baris `Amount in IDR` satuannya **TETAP `IDR`** dan hanya-baca: sel
kirinya di ekspor `.pyTemplateRichTextEditor`, pemegang tempat yang hanya
memperlihatkan satuan — bukan pilihan.

⚠️ `ambil` adalah tanggungan efek `DropdownWarisan`; fungsi baru tiap render
akan mengambil ulang daftarnya tanpa henti. Keduanya dibekukan `useCallback`,
dan daftarnya diambil SEKALI untuk seluruh baris.

---

## 4 · Yang BELUM tersentuh

- **Menyimpan.** Aturan B tetap berlaku: suntingan hidup di keadaan layar
  sampai Save/Submit. Rute yang tab ini panggil ber-awalan `/hitung/`, yang
  menyatakan dirinya nol tulis basis data.
- **`ClassOfBusiness`** ada di tipe barisnya tetapi nol medan layar — ekspor
  tidak menampilkannya di grid maupun rincian.
- **`AltValue`** pada `TotalEgnpiAmountNP(<LAST>)` diisi `local.Altvalue` yang
  nol pernah disetel di Activity mana pun. Tidak dibangun.
