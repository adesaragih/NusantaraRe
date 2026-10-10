# Tab Maximum Retention (non-prop) — aplikasi lawan ekspor Pega

Dibaca 6 Oktober 2026 dari `D:\XML_NURE\Treaty In`:

```
Section/TreatyInTabsNonProportional.xml   wadah TABBED ke-1
Section/MaxRetention.xml                  rincian satu baris
FlowAction/MaxRetention.xml               pembungkus rincian
Activity/TreatyInNPSetTotal.xml           param.type == "retention"
Activity/TreatyInNonAddItem.xml           param.Type == "retention"
```

---

## 1 · Hasil perbandingan

| | Pega | Aplikasi SEBELUM | Sekarang |
|---|---|---|---|
| Grid `Treaty Group` · `Currency` · `Amount` | ✅ | grid 5 kolom hanya-baca | ✅ kartu lipat + rincian |
| `Add` | ✅ | ⛔ tidak ada | ✅ |
| `Delete` | ✅ | ⛔ tidak ada | ✅ |
| Rincian baris (`MaxRetention`) | ✅ | ⛔ tidak ada | ✅ |
| Panel `Total Retention Amount` \| `Value` | ✅ | ✅ (nilai tersimpan) | ✅ + dihitung ulang |
| Tombol `Update Total` | ✅ | ⚠️ **ADA tetapi MATI** | ✅ hidup |

⭐ Tombolnya **mati dengan jujur** sebelum ini — rumusnya memang belum ada,
dan tombol hidup yang tidak menghitung apa pun lebih buruk daripada tombol
mati. Yang berubah: rumusnya kini ada.

---

## 2 · Rumus yang dipindahkan

`backend/services/hitung_retensi.go`, diuji di `hitung_retensi_test.go`,
dipanggil lewat `POST /api/treaty-in/hitung/retensi`.

### 2.1 `TreatyInNPSetTotal` (`param.type == "retention"`)

```
[2] kosongkan TreatyIn.TotalRetentionAmountNP
[3] per baris Retention: Σ .Amount PER MATA UANG
    (pola `Appendflag`: gabung ke baris bermata uang sama, atau tambah)
```

⛔ **Nol pagar bagi-nol** di cabang ini, dan itu benar: retensi tidak pernah
membagi. Menyalin pagar dari cabang EGNPI akan menolak data yang sah.

⭐ Urutan baris total mengikuti **kemunculan pertama** tiap mata uang — pola
`Appendflag`, bukan urutan abjad.

### 2.2 `TreatyInNonAddItem` (`Type == "retention"`)

Satu baris **kosong**, `ID = ""`. ⚠️ Nol nilai awal diwarisi dari mana pun —
beda dengan cabang `egnpi` yang mewarisi mata uang `Retention(1)`.
Ketergantungannya satu arah: **EGNPI membaca retensi**, bukan sebaliknya.

---

## 3 · ⛔ Temuan yang perlu dinyatakan

### 3.1 Satu tombol, bukan dua

Tab EGNPI punya `Update Total` **dan** `Update EGNPI Value`; tab ini hanya
yang pertama. Sebabnya: retensi **nol punya kolom `Amount in IDR`**, jadi nol
konversi kurs yang perlu dijalankan. Menambahkan tombol kedua berarti membuat
tombol yang di Pega tidak ada dan tidak punya pekerjaan.

### 3.2 `Update Total` MATI bila `EDMMaterialType = 2`

`pyDisabledWhen` @202558 — terukur di keempat sel tombol sejenis, dan sudah
tercatat di `labels.ts` sebelum ronde ini. Syaratnya kini dipakai, bukan
sekadar dicatat: `totalTerkunci()` di `labelsRetensi.ts`.

### 3.3 Tombol KEDUA di ekspor — MATI, tidak dibangun

Tepat di atas `Update Total` ada tombol tak berlabel (`.pyTemplateButton`)
ber-`pyCondition 1=2` @198452. Kemungkinan pendahulunya. **Yang dibangun yang
hidup**, dan uji menjaga jumlah tombol tetap satu.

### 3.4 Ekspor berselisih dengan dirinya sendiri — `Treaty Group`

| | `pyReadOnly` |
|---|---|
| grid tab | *(tidak ada)* → dapat diubah |
| `Section/MaxRetention.xml` | `true` |

Persis selisih yang sama dengan tab EGNPI. Yang dipakai bentuk **grid**:
`Add` melahirkan baris kosong, dan baris yang kelompok treaty-nya tidak dapat
diisi tidak berguna bagi siapa pun.

### 3.5 ~~`Note` nol pernah dapat diisi~~ — ⛔ RALAT 7 Oktober 2026

**Bacaan itu KELIRU dan ditarik.** Sel `Note` memuat `pyReadOnly = true`
**dan** `pyReadOnlyCondition = TreatyIn.IsEditData = 1`; syaratnya yang
berlaku, jadi medan itu **aktif di mode Edit**. Sudah diperbaiki.

⚠️ **Syaratnya BUKAN `ViewState`** — beda dengan tab EGNPI yang medan
serupanya memakai `TreatyIn.ViewState = 1`. Nol tambalan menyeluruh untuk
keduanya; tiap sel harus dibaca syaratnya sendiri.

⛔ `IsEditData` nol dibawa aplikasi. Ia DIDEKATI dengan mode Edit,
mengikuti preseden `TabPortofolio.tsx` (sel 112·113·114, syarat yang sama)
yang sudah melakukan hal itu dan diuji begitu. Pendekatan ini salah HANYA
ketika `IsEditData = 1` — keadaan yang belum pernah terukur di aplikasi.

⭐ Tetap `textarea` ber-`readOnly` sungguhan di mode lihat, bukan kotak yang
terlihat dapat diisi lalu membuang ketikannya.

⭐ Aturannya kini terkodekan SEKALI di
`D:\XML_NURE\_migration-docslat-baca-ekspor\` (`hanya_baca`, Aturan 3).

### 3.6 Baris `Amount` memuat DUA medan

`pyLabelFieldValue = Amount` menempel pada **`.Currency`**, dan `.Amount` di
sebelahnya tanpa label sendiri. Jadi satu baris berlabel `Amount` memuat
pilihan mata uang lalu kotak nilai — persis yang terlihat di tangkapan layar
pemilik proses.

### 3.7 Dua presisi untuk nilai yang sama

`Amount` di grid berbunyi `3.500.000.000` (0 desimal); panel
`Total Retention Amount` tepat di bawahnya `3.500.000.000,00` (2). Keduanya
dipertahankan — itu bukan kekeliruan tangkapan layar.

---

## 4 · Yang berubah di sekitarnya

- **`PanelTotalRetensi.tsx` DIHAPUS.** Seluruh isinya — grid total dan tombol
  matinya — kini hidup di `TabRetensi.tsx`. Berkas mati yang ditinggalkan
  akan terbaca sebagai pilihan kedua yang masih sah.
- **Total tersimpan tetap dipakai sebagai nilai AWAL** (`totalAwal`). Pega
  memperlihatkan total yang sudah ada di clipboard kontrak; panel kosong
  sampai tombol ditekan akan terbaca sebagai "kontrak ini nol retensi", dan
  itu berbeda dari "belum dihitung ulang".
  ⚠️ Tab **EGNPI** tidak punya padanannya: dokumen nol menyimpan
  `TotalEgnpiAmount`, jadi di sana panelnya memang kosong sampai ditekan.

## 5 · Yang BELUM tersentuh

- **Menyimpan.** Aturan B tetap: suntingan hidup di keadaan layar sampai
  Save/Submit. Rute yang dipanggil ber-awalan `/hitung/` — nol tulis.
- **`ClassOfBusiness`** dibawa di tipe barisnya tetapi nol medan layar —
  ekspor tidak menampilkannya di grid maupun rincian.
