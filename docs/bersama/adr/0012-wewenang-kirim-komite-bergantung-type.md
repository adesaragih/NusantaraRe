---
status: accepted
tanggal: 2026-09-14
sumber: grilling Ronde 4 (`.scratch/claim-life/grilling-ronde-4.md`), keputusan work owner
---

# Wewenang kirim ke Komite bergantung `Type` — dibawa apa adanya, dengan risiko RBAC tercatat

Di Claim — Life, **siapa yang boleh mengirim kasus ke Komite ditentukan oleh `Type` klaim**, bukan
oleh peran saja:

| `Type` | Siapa yang boleh mengirim |
| --- | --- |
| **`TP`** (Payable) atau **`TR`** (Receivable) | **siapa pun** — `ReasLifeAdmin` termasuk |
| lainnya (`QP`, `QR`) | **hanya `ReasLifeSPV`** |

**Perilaku ini dibawa apa adanya ke sistem baru (paritas). Tidak diperketat pada migrasi ini.**
Risikonya dicatat di bawah untuk ditinjau saat konteks **Komite Life** / **IAM** digarap.

## Arti kode `[keputusan work owner 2026-09-14]`

**`TP` = Payable. `TR` = Receivable.**

⚠️ Definisi ini **tidak ada di korpus** — sumbernya work owner. `[terverifikasi]` Pencarian korpus
sudah tuntas dan nihil: tidak ada label, caption, atau opsi mana pun yang memasangkan kode itu
dengan teks; satu-satunya tempat nilainya **ditetapkan** adalah
`PremiumList Life/Activity/SubmitPremiumList_Act.xml` (`ASM-FW-GISFW-WORK-LIFE!SUBMITPREMIUMLIST_ACT`,
214.154 byte), dan di sana keempat kode hanya diteruskan ke slot generik `InputData.CARI20` dengan
precondition `.Type=="<kode itu sendiri>"` — melingkar, tidak mendefinisikan.

`QP` dan `QR` **belum** dijawab dan tetap terbuka di **OQ-020**.

## Bukti: `Type` menggerbangi dua hal `[terverifikasi]`

### 1. Wewenang kirim ke Komite

`Claim Life/Section/AdjustmentDetail_Section.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE!ADJUSTMENTDETAIL_SECTION`, 557.563 byte), kontrol jalur Komite
(`GetListKomiteLife`):

```
pyWorkPage.pyPosition =='ReasLifeSPV' || pyWorkPage.Type = 'TP' || pyWorkPage.Type = 'TR'
```

Seluruhnya `||`, tanpa `&&` — tidak ada ambiguitas presedensi. Dan `=` di korpus ini adalah
**pembanding**, bukan assignment: pada `<pyCondition>` bentuk `=` justru mayoritas (2.282 vs 954
`==`), dan banyak ekspresi mencampur keduanya dalam satu baris di tempat yang akan rusak kasatmata
kalau `=` berarti assignment — mis. `.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"`.

`[terverifikasi]` Gerbang seperti ini **hanya ada di `Claim Life`** — tidak ada di 19 modul lain.

### 2. Jendela validasi Date of Loss

`Claim Life/Activity/ValidasiDOL_Act.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!VALIDASIDOL_ACT`, 59.747 byte):

| Cabang | Jendela yang dipakai | Pergeseran tanggal |
| --- | --- | --- |
| `Type=="QR" \|\| Type=="QP"` | `GROSS_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `@addCalendar(.DATE_OF_LOSS,0,0,0,0,0,0,0)` — **nol** |
| `Type=="TR" \|\| Type=="TP"` | `RETROCESSION_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `@addCalendar(.DATE_OF_LOSS,0,0,0,1,0,0,0)` — **+1 hari** |

Gagal → `local.errmsg = "Invalid DOL"`, dengan `local.Begin==false || local.Expired==true`.

### 3. Kedua gerbang membaca **dua salinan** properti yang sama

`[terverifikasi]` Gerbang Komite membaca `pyWorkPage.Type`; `ValidasiDOL_Act` membaca
`pyWorkPage.PolicyDataLife.Type`. Keduanya berisi nilai yang sama karena disalin:

```
Claim Life/Activity/LoadDataPeserta_Act.xml
  pyWorkPage.Type = pyWorkPage.PolicyDataLife.Type
```

Sumber otoritatifnya adalah **`PolicyDataLife.Type`** — tipe polis; `pyWorkPage.Type` hanya salinan
kerja. Di sistem baru keduanya menjadi **satu** field; yang perlu disadari adalah di Pega mereka
dua, sehingga secara teori dapat berbeda bila `LoadDataPeserta_Act` tidak berjalan.

## Considered Options

- **Paritas — bawa apa adanya, catat risikonya** — dipilih
- Perketat menjadi `SPV && (TP || TR)` — ditolak pada migrasi ini: itu **mengubah perilaku**, dan
  wewenang lintas konteks (Claim Life ↔ Komite Life ↔ IAM) belum ditetapkan menyeluruh.
  **OQ-021** menunjukkan RBAC adalah blok terbesar konteks Komite

Membiarkan pilihan ini tidak tercatat bukan opsi: seorang pembaca di kemudian hari akan menyangka
ini kekeliruan porting, lalu "memperbaikinya" tanpa tahu bahwa ia disengaja.

## Risiko yang diterima

⚠️ **Wewenang bergantung pada atribut data, bukan pada peran.** Siapa pun yang dapat membuat klaim
ber-`Type` `TP`/`TR` dengan sendirinya memperoleh wewenang mengirim ke Komite. Peran bukan lagi
satu-satunya penentu akses.

⚠️ **`Type` adalah arah akuntansi (Payable/Receivable), bukan tingkat kewenangan.** Tidak terbaca —
di korpus maupun dari definisinya — mengapa arah akuntansi menentukan siapa boleh mengeskalasi.
Kemungkinannya: kelonggaran operasional lama, atau `Type` membawa makna lain yang belum dinyatakan.
**Tidak ditebak.**

⚠️ **Ditinjau ulang saat Komite Life / IAM digarap**, bukan sekarang. Bila peninjauan itu
memutuskan memperketat, ADR ini yang diganti.

## Consequences

- Lapisan layanan harus menegakkan aturan bergantung-`Type` ini **secara eksplisit** — bukan
  mewarisinya sebagai efek samping visibilitas UI, sebagaimana di Pega. Penegakan di lapisan layanan
  mengikuti **ADR-0002**.
- `Type` menjadi **field yang menyentuh keamanan**, bukan sekadar data. Perubahan nilainya harus
  masuk jejak audit (**ADR-0007**).
- Mesin status (**ADR-0011**) tidak berubah karenanya — yang berbeda hanya **siapa** yang boleh
  melakukan langkah 4.
- Validasi DOL mewarisi percabangan yang sama; keduanya harus memakai **sumber `Type` yang sama**
  agar tidak dapat berbeda seperti di Pega.

## Tangga komite sepenuhnya data-driven (2026-09-14) `[terverifikasi]`

Penutupan **OQ-032** melengkapi ADR ini dari sisi lain: `Type` menentukan **siapa** yang boleh
menyerahkan; **roster** menentukan **berapa tingkat** tangga persetujuannya.

> **`KomiteLoop` = COUNT baris roster `EMAILKOMITE` yang aktif (`STS_AKTIF = "1"`) dan
> ber-`LIMIT_BOTTOM <= CLAIM_AMOUNT`.** Tidak ada konstanta di mana pun.

Tiga bukti berantai:

1. `Claim Life/Activity/GetListKomiteLife.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
   `GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`) — `Local.IsADj = .CLAIM_AMOUNT`, lalu memanggil report
   `FilterEmailKomiteWithLimit` (`Param.pyReportClass = "ASM-FW-GCNMFW-Int-EMAILKOMITE"`), hasilnya
   ke `Primary.KomiteList(<APPEND>)`.
2. `Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml` (`ASM-FW-GCNMFW-INT-EMAILKOMITE` /
   `FILTEREMAILKOMITEWITHLIMIT` / `RULE-OBJ-REPORT-DEFINITION`) — filter
   `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM AND .STS_KLAIM = Param.STS_KLAIM AND .STS_AKTIF = "1"`.
3. `Claim Life/Activity/CreateKMTLife_Act.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
   `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`) —
   `childPageKomite.KomiteLoop = @Utilities.SizeOfPropertyList(childPageKomite.KomiteList)`.

⚠️ **Detail yang wajib direplikasi:** ambang yang dikirim ke roster adalah **nilai mutlak** klaim —
`Param.LIMIT_BOTTOM = @if(Local.IsADj<0, Local.IsADj * -1, Local.IsADj)`. Klaim bernilai negatif
dicari dengan tandanya dihilangkan.

**Akibat:** menambah atau menonaktifkan satu baris roster **mengubah jumlah tingkat persetujuan**
klaim yang sedang berjalan. Roster adalah **data operasional yang menyentuh alur keputusan**, bukan
sekadar daftar alamat email — perubahannya layak masuk jejak audit (**ADR-0007**).

⚠️ `[data DBA]` Roster **tidak punya kolom mata uang**; pita nilainya berlaku atas satu mata uang
implisit. Sahih selama invariant **ADR-0003** (satu klaim satu mata uang) dipegang.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-020** | Arti `QP` dan `QR` belum dijawab. `[dugaan]` Bukti §2 memperlihatkan huruf **pertama** memisahkan gross (`Q*`) dari retrocession (`T*`); bila huruf **kedua** memang Payable/Receivable seperti pada `TP`/`TR`, maka `QP`/`QR` mengikuti pola yang sama — **belum dikonfirmasi, jangan dipakai sebagai fakta** |
| **OQ-021** | RBAC lintas konteks belum ditetapkan; peninjauan ulang ADR ini bergantung padanya |
| **OQ-001** | Tidak ada DDL — nilai `Type` yang sah menurut basis data tidak diketahui |
