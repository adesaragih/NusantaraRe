# Sintesis Domain Komite — D2 Tahap 2

Perbandingan berbukti atas **empat** tangga persetujuan Komite yang ditelusur 2026-09-13:
Komite Claim Life, Komite Claim FacIn, Komite Claim Prop, Komite Claim Non Prop.

**Ini bukan penetapan bounded context** — itu STEP D4. Dokumen ini menyajikan bukti.

---

## 1. Empat modul, empat rule — bukan salinan

`[terverifikasi]` Keempat flow **berbeda class**, dan hash ternormalisasi 18 tag membuktikan isinya
juga berbeda:

| Modul | Class flow | Nama file Flow | Hash |
| --- | --- | --- | --- |
| Komite Claim Life | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `KomiteLife_Flow.xml` | `b69929c407` |
| Komite Claim FacIn | `ASM-FW-GCNMFW-WORK-KOMITE` | `Komite_Flow.xml` | `f9366db22d` |
| Komite Claim Prop | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `KomiteTreaty_Flow.xml` | `5eacdb3783` |
| Komite Claim Non Prop | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `KomiteTreaty_Flow.xml` | `63b913a161` |

**Dua modul terakhir memakai nama file yang sama persis** — ditelusur sebagai dua rule terpisah
sesuai `_METHOD.md` §1.1.

Perintah audit:
```
grep -o "<pxInsName>[^<]*" "<modul>/Flow/<file>.xml" | head -1
grep -vE "$VOL18" "<modul>/Flow/<file>.xml" | sort | md5sum
```

---

## 2. Struktur tangga: **identik di keempat modul**

`[terverifikasi]` Keempatnya punya graf yang **persis sama bentuknya**: 1 Assignment, 1 Decision,
4 connector.

```
Start ─(Always)─> Assignment "KomiteRouter" [WorkList, route=Custom]
                      │ FlowAction: ViewTransferDtl
                      v
                  Decision "KomiteLoop"
                      ├─ When IsKomiteLoop ─> Assignment   ← LOOP
                      └─ Else ─────────────> End
```

**Tangga persetujuan tidak dimodelkan sebagai shape.** Ia adalah **loop satu assignment** yang
dikendalikan data. Ini temuan struktural utama Tahap 2.

Kondisi `IsKomiteLoop` **sama teksnya di keempat modul** `[terverifikasi]`:

```
.AcceptStatus = "1"   DAN   .KomiteCount <= .KomiteLoop
```

meski keempat rule `When`-nya berbeda (class berbeda, hash berbeda).

### 2.1 Jejak per tingkat adalah list data

`[terverifikasi]` Di Komite Claim Life, tiap tingkat menulis satu entri
`KomiteList(KomiteCount)` berisi `KomiteAproval` (= `AcceptStatus`), `KomiteComment`, dan
`DateApprove` (`@CurrentDateTime()`).

Jadi **berapa tingkat tangga dan siapa penyetujunya ditentukan data, bukan struktur proses.**

---

## 3. Routing: **dua pola berbeda**

`[terverifikasi]` Keempat modul memakai `<pyImplementation>WorkList` + `<pyRouteTo>Custom` dengan
router `KomiteRouter`. Isi router terbelah dua:

| Pola | Modul | Sasaran `param.AssignTo` |
| --- | --- | --- |
| **Berbasis data** | Komite Claim Life | **hanya** `.KomiteID`, digerbangi `TransferType == '2'` |
| **Ter-hardcode + fallback** | FacIn, Prop, Non Prop | `.KomiteCount == 1..4` → `"komitepnc"`, `"komitepnc2"`, `"komitepnc3"`, `"komitepnc4"`; lalu cabang `TransferType == '2'` / `!= '2'` → `.KomiteID` / `.Komite.KomiteID` |

`[terverifikasi]` Tiga modul memuat **empat nama sasaran ter-hardcode** yang sama persis. Hanya
Komite Claim Life yang sepenuhnya berbasis data. → **OQ-036**.

Perbedaan kecil: router Prop memuat nilai `.KomiteCount+1`; Non Prop tidak.

`[terverifikasi]` Keempat modul memakai **`WorkList`**, tidak satu pun `WorkBasket` — berbeda dari
treaty inward (D2 Tahap 1) yang seluruhnya `WorkBasket` → OQ-028.

---

## 4. Otorisasi: **dua pendekatan berbeda untuk proses setara**

Ini temuan paling material Tahap 2.

| Modul | Guard identitas orang ter-hardcode | Pola alternatif |
| --- | ---: | --- |
| Komite Claim FacIn | **4 file** (`KomitePost_Adjustment`, `KomitePost_CloseClaim`, `KomitePost_Reject`, `SetProteksiSubmiteKomite`) | — |
| Komite Claim Prop | **3 file** (`KomitePostAdjustment`, `KomitePost_Close`, `KomitePost_Reject`) | — |
| Komite Claim Life | **0** | sasaran routing dari `.KomiteID` |
| Komite Claim Non Prop | **0** | **mencocokkan pengguna dengan roster** (§4.1) |

Perintah audit:
```
grep -rlE "OperatorID\.pyUser(Identifier|Name)[ ]*[=!]+[ ]*[\"']" "<modul>" --include="*.xml"
```

`[terverifikasi]` Guard di FacIn dan Prop berada di `<pyExpression>` — **di dalam activity keputusan
komite**, bukan di lapisan UI. **Nilai nama orang tidak disalin** ke artefak D2 (OQ-021).

### 4.1 Pola berbasis roster di Non Prop `[terverifikasi]`

`Komite Claim Non Prop/Activity/KomitePostAdjustment.xml` step 2–3:

```
@contains( @toUpperCase(pyWorkPage.KomiteList(Local.IdxKomite).KomiteID),
           @toUpperCase(OperatorID.pyUserIdentifier) )
```

Pengguna yang login dicocokkan dengan **entri roster**, bukan dengan daftar nama ter-hardcode.

**Korpus memuat dua pendekatan otorisasi untuk proses yang setara**, dan salah satunya sudah
berbasis data. Ini titik awal paling konkret untuk RBAC di domain Komite.

### 4.2 Sumber roster `[terverifikasi]`

`Komite Claim FacIn/Activity/ApprovalKomite_Act.xml` membangun roster lewat `Obj-Browse` terhadap
class **`ASM-FW-GCNMFW-Int-EMAILKOMITE`** ke halaman `GetKomite`, lalu membuang duplikat
(`Property-Remove` bila `KomiteList(1).KomiteID == .OPERATOR_ID`).

Jadi **daftar anggota komite berasal dari tabel database**, sementara sasaran routing justru
ter-hardcode (§3). Kedua mekanisme itu hidup berdampingan.

### 4.3 Ambang nominal `[terverifikasi]`

Di `ApprovalKomite_Act` (FacIn):

```
Local.TotalAdj > 30000000.00 && Local.TotalAdj <= 57750000.00
```

Sebuah **pita nilai ter-hardcode** menentukan `Obj-Browse` mana yang dijalankan — artinya
**komposisi roster bergantung besaran nilai**. `Local.TotalAdj` diakumulasi dari `+ .ValueAdjustment`.

Mata uang **tidak disebut**. Activity ini 146 KB dan **belum habis dibaca**, jadi jumlah pita
seluruhnya **belum terukur** → **OQ-037**.

---

## 5. Efek samping: empat rule penomoran berbeda

`[terverifikasi]` Keempat modul menulis akseptasi lewat rule penomoran yang **semuanya berbeda**:

| Modul | Rule penomoran | Pencabangan |
| --- | --- | --- |
| Komite Claim Life | `Generate_NoAccept_KMT_Life` / `Generate_NoAccept_KMT_LifeRetro` | menurut kode `Type`: `QP`/`QR` vs `TP`/`TR` |
| Komite Claim FacIn | `GenerateNoAccept` / `GenerateNoAcceptNonFire` | menurut `When IsFire` |
| Komite Claim Prop | `GenerateNoAcceptTreaty` | — |
| Komite Claim Non Prop | `GenerateNOAccCTNP` | — |

**Dua modul memakai kriteria pencabangan yang berlainan untuk keputusan yang sama** (kode produk
vs lini bisnis).

Tiga rule dipanggil **oleh keempat modul**: `GETTanggalClosing_SQL`, `GetSequenceNumber_SQL`,
`InsertHistoryAkseptasiPega_Sql` (yang terakhir **berkonflik** di register OQ-011).

### 5.1 Tabel akseptasi dapat dibaca langsung

`[terverifikasi]` `Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` adalah blok PL/SQL
`INSERT INTO POOLDATA.OS_AKSEPTASI_KLAIM_LIFE (...)` + `COMMIT` — **bukan** stored procedure,
sehingga **55 nama kolom terbaca** (daftar lengkap di `Komite Claim Life.md` §3.1).

Ini fragmen skema nyata yang meringankan sebagian **OQ-001**. Tipe kolom tetap tidak diketahui.

---

## 6. Apakah Komite = wrapper approval murni?

D1 batch 3 mengukur: **59 dari 133 identitas modul Komite (44,4%) tidak muncul di modul Claim mana
pun** (`../inventory/_summary.md` §15.4).

Bukti D2 Tahap 2 **menguatkan** bahwa Komite bukan lapisan tipis:

| Bukti | Keterangan |
| --- | --- |
| Class kerja sendiri | `...WORK-KOMITELIFE`, `...WORK-KOMITE`, `...WORK-KOMITETREATY`, `...WORK-KOMITETREATYNONPROP` — empat class, bukan menumpang class Claim |
| Model data sendiri | `KomiteList` dengan `KomiteAproval`, `KomiteComment`, `DateApprove`, `KomiteID`, `KomiteCount`, `KomiteLoop` |
| Roster sendiri | class `ASM-FW-GCNMFW-Int-EMAILKOMITE` |
| Rule penomoran akseptasi sendiri | empat rule berbeda (§5) |
| Activity besar khas | Non Prop: `KomitePostAdjustmentCWP` 434 KB, `GenerateAccCNP_act` 325 KB; FacIn: `KomitePost_Adjustment` 626 KB |
| Tulis ke tabel akseptasi | `OS_AKSEPTASI_KLAIM_LIFE`, `SaveOSClaim_SQL`, `InsertOSKlaimCNP` |

`[dugaan]` Komite adalah **tahap proses dengan model data, roster, penomoran, dan penulisan
database sendiri** — bukan wrapper approval tanpa data. **Belum terverifikasi** sepenuhnya karena
sebagian besar activity belum ditelusur. **Penetapan konteks tetap D4.**

---

## 7. Ketidakkonsistenan yang terukur

`[terverifikasi]` Empat proses yang secara bentuk identik, tetapi berbeda dalam hal-hal berikut:

| Aspek | Life | FacIn | Prop | Non Prop |
| --- | --- | --- | --- | --- |
| Properti batas tangga | `.KomiteLoop` | `Local.TotalKomite` | `.KomiteLoop` | `.KomiteLoop` |
| Sasaran routing | data (`.KomiteID`) | hardcode + data | hardcode + data | hardcode + data |
| Guard identitas orang | tidak ada | 4 file | 3 file | tidak ada (roster) |
| Rule penomoran | 2 (per kode `Type`) | 2 (per `IsFire`) | 1 | 1 |
| Activity `KomitePost_Close`/`_Reject` | — | ada | ada | **tidak ada** |
| Ambang nominal di roster | belum terlihat | ada (1 pita) | belum terlihat | belum terlihat |

**Empat implementasi berbeda untuk satu konsep bisnis yang sama.** Ini fakta terukur, bukan
penilaian — konsekuensinya untuk migrasi adalah keputusan FASE B.

---

## 8. Kode yang menggerbangi logika (OQ-020)

`[terverifikasi]` Terkumpul di Tahap 2:

| Kode | Nilai | Modul |
| --- | --- | --- |
| `AcceptStatus` | `"1"`, `"2"` | keempat |
| `TransferType` | `'1'`, `'2'` | keempat |
| `PaymentType` | `1`–`6` | FacIn, Prop, Non Prop |
| `TypeComentAnalysis` | `"5"` | Prop |
| `.Type` (Life) | `QP`, `QR`, `TP`, `TR` | Life |
| `STS_REJECT` | `"0"` | Life |
| `.AcceptanceStatus` | `""`, `"0"` | FacIn |
| `RetroID` | `"1000013"` | Life |
| ambang tanggal | `25` | Life, Non Prop |
| tanggal cutover | `20250207T000000.000 GMT` | Life |

**Arti seluruhnya belum terverifikasi.** Dicatat literal apa adanya.

`[terverifikasi]` Yang perlu diperhatikan: di `Komite Claim Prop/Activity/KomitePost_Reject.xml`,
panggilan `SaveRejectTreatyIn_Act_KMT` digerbangi **`AcceptStatus=="1"`** — bukan `"2"`. Pemetaan
`AcceptStatus` → setuju/tolak **tidak dapat ditebak dari nama rule**.

---

## 9. Batas pengetahuan Tahap 2

- **`IsPEGAPROD`** ada di keempat modul, kondisinya **tidak terbaca** dari tag, dan ia
  **berkonflik** (OQ-011 #305). Di Komite Claim Life ia menggerbangi **tiga** pemanggilan keluar.
  **Cabang mana yang aktif tidak ditebak** → OQ-029.
- **32 identitas berkonflik** menyentuh modul Komite (`grep "Komite" _oq011-konflik-isi.md`).
  Yang tersentuh telusur ini dirujuk per-varian; yang di luar jalur inti tidak dibaca.
- **`UpdateWorkObject`** dan **`serviceInsertArasapasClaimLife_act`** dipanggil tetapi salinannya
  hanya ada di modul lain → OQ-035.
- **`ApprovalKomite_Act`** (146 KB) dan enam activity CNP besar **belum habis ditelusur**.
- Sumber nilai **`.KomiteLoop`** tidak ditemukan di modul Komite mana pun → OQ-032.
