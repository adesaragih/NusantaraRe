# Grilling — Claim — Life — Ronde 3

Status: answered (work owner, 2026-09-14)
Konteks: `claim-life` (Claim — Life)
Tanggal: 2026-09-14
Skill: `/mattpocock-skills:grill-with-docs` (grilling + domain-modeling)
Ronde sebelumnya: `grilling-ronde-1.md` (Q1–Q7), `grilling-ronde-2.md` (Q8–Q15)

> Ronde ini **tidak berbentuk pertanyaan baru**. Work owner langsung menjawab **Blok A** dari
> `kesiapan-to-spec.md` — tiga pemblokir mesin status — plus persetujuan ADR-0010.
> Jawaban dicatat apa adanya; setiap klaim yang **tidak** terbukti di korpus ditandai
> `[keputusan work owner]`, bukan `[terverifikasi]`.

---

## Jawaban work owner — mesin status Claim Life

`[keputusan work owner 2026-09-14]` **Unit status = baris `AdjustmentList`.**

> 1. Admin (ReasLifeAdmin) input AdjustmentList pertama, lalu "Save ke OS" → STS_REJECT = 0
>    (via SaveOutStandingLife_Act - AKTIF). Admin TIDAK mengirim ke Komite.
>    Admin JUGA dapat me-reject adjustment yang ia input sendiri, TANPA lewat Komite.
> 2. Submit ke Medical Check (ReasLifeMedicalAdvisor).
> 3. Submit ke SPV (ReasLifeSPV).
> 4. "Send ke Komite" HANYA dilakukan SPV (untuk semua adjustment).
> 5. Komite memutus: aksep → STS_REJECT = 1; tolak → STS_REJECT = 2
>    (ditulis rule sisi Komite: KomitePostAdjustment, precondition AcceptStatus).
> 6. Jika Komite tolak (2): SPV add AdjustmentList BARU → Save ke OS (0) → send Komite. Berulang.

**Dua sumber penolakan (koreksi work owner):**

> - Admin reject langsung (adjustment yang ia input) — TANPA Komite.
> - Komite reject — lewat SPV.
> Karena itu Claim Life/Activity/RejectOSClaimLife_Act.xml BUKAN dead rule — ia jalur reject oleh
> Admin. (SaveAdjustment_Act yang set STS_REJECT=1 tanpa precondition Komite: status aktif/dead
> BELUM dipastikan work owner → biarkan sebagai catatan, jangan diklaim dead.)

**OQ-061 (tiga tingkat):**

> Unit keputusan = baris AdjustmentList. STS_REJECT di PremiumListDetail = CERMINAN hasil
> AdjustmentList TERAKHIR (menyalin status baris adjustment paling akhir), bukan unit keputusan
> sendiri. Header klaim mengikuti adjustment terakhir. Tutup OQ-061 untuk Claim Life.

**OQ-062 (final/buka-ulang):**

> Terminal PER BARIS AdjustmentList (sekali 1/2 tak berubah; tak ada rule reset 1/2 → 0 pada baris
> sama). KLAIM tidak terminal: setelah tolak, SPV buat baris baru. Revisi = baris baru. Tutup
> OQ-062 untuk Claim Life.

**OQ-039 (jalur aksep/tolak) — dijawab sebagian:**

> Ada DUA jalur reject: (a) Admin reject adjustment yang ia input, tanpa Komite; (b) Komite reject
> lewat SPV. Aksep final selalu lewat Komite. MASIH TERBUKA (jangan ditebak): apakah reject oleh
> Admin membatalkan adjustment itu saja atau menutup seluruh klaim, dan apakah memakai
> STS_REJECT=2 yang sama atau mekanisme lain. → OQ-039 dipersempit, tetap terbuka pada dua detail
> itu, pemilik Product+UW.

**ADR-0010:**

> SETUJUI, naikkan dari draft ke accepted. Ganti storage di luar cakupan; kontrak
> GET_TOKEN_STORAGE tetap OQ-002 (DBA).

**Pemilik:** Product+UW (mesin status, jalur reject), IT-infra + DBA (ADR-0010)

---

## Verifikasi korpus atas jawaban di atas

Dijalankan sebelum mencatat, sesuai aturan "klaim wajib bukti path+rule".

### V1 `[terverifikasi]` — `RejectOSClaimLife_Act` **aktif**, dan memang jalur Admin

Dirujuk 2× sebagai `<pyActivity>` dari `Claim Life/Section/RejectOSClaimLife_Sec.xml`
(127.445 byte). Kontrol yang memicunya di
`Claim Life/Section/AdjustmentDetail_Section.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE!ADJUSTMENTDETAIL_SECTION`, 557.563 byte) berlabel
**"Reject Outstanding"**, digerbangi:

```
pyWorkPage.pyPosition =='ReasLifeAdmin'
  && pyWorkPage.ClaimData.PremiumListSummary.CLAIM_NO !=''
  && .STS_REJECT == 0
```

**Koreksi work owner terbukti.** Dan gerbangnya menguji `.STS_REJECT` **tingkat baris** — bukti
kuat bahwa unit keputusan memang baris.

### V2 `[terverifikasi]` — jalur Komite digerbangi SPV, **tetapi ada jalan pintas**

Pada berkas yang sama, kontrol yang memicu `GetListKomiteLife` digerbangi:

```
pyWorkPage.pyPosition =='ReasLifeSPV' || pyWorkPage.Type = 'TP' || pyWorkPage.Type = 'TR'
```

⚠️ Disjungsi `Type = 'TP'` / `'TR'` **melewati pemeriksaan peran sama sekali**. Pernyataan
"Send ke Komite HANYA SPV" **tidak berlaku** untuk kedua tipe itu. → **OQ-063** (baru).

### V3 `[terverifikasi]` — `SaveOutStandingLife_Act` aktif, dipanggil dari **dua** tempat

Hanya ada **satu** berkas activity (`Claim Life/Activity/SaveOutStandingLife_Act.xml`), dirujuk
dengan dua ejaan huruf besar-kecil yang berbeda (Pega tidak membedakan):

| Pemanggil | Gerbang peran |
| --- | --- |
| `Claim Life/Section/InputOSClaimLife.xml` (799.901 byte) | **tidak ada** (`pyPosition` = 0 kemunculan) |
| `Claim Life/Section/AdjustmentDetail_Section.xml`, tombol "Save to Outstanding" | `pyWorkPage.pyPosition =='ReasLifeSPV'` |

Jalur SPV pada langkah 6 terbukti. Jalur Admin pada langkah 1 **ada** tetapi **tidak digerbangi
peran di UI** — lihat V6.

### V4 `[terverifikasi]` — pencerminan ditulis berbarengan, bukan disalin belakangan

| Rule | Tingkat baris | Tingkat `PremiumListDetail` |
| --- | --- | --- |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml` | `.STS_REJECT = 2` | `= 2` |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `1` (6×) / `2` (2×) | nilai sama, berpasangan |

`[keputusan work owner]` Bahwa yang tercermin adalah baris **terakhir** berasal dari work owner;
korpus hanya menunjukkan yang tercermin adalah **baris yang sedang diputus**.

### V5 `[terverifikasi]` — baris baru diwarisi dari baris pertama, tanpa status

`Claim Life/Activity/SetIndexAdjustmentList.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!SETINDEXADJUSTMENTLIST`, 47.748 byte) menyalin delapan kolom
dari `AdjustmentList(1)` ke `AdjustmentList(<LAST>)`: `SHARE_NUSANTARA_RE`, `CEDING_RETENTION`,
`SUM_REASURED`, `SUM_INSURED`, `SHARE_RETRO`, `RETROCEDED_SHARE`, `CURRENCYID`, `CURRENCY`.
`STS_REJECT` **tidak** ikut. Persis yang dibutuhkan langkah 6.

Ikutan untuk **OQ-060**: `CURRENCY` ada di **tingkat baris** dan disalin dari baris 1 — jadi dalam
praktiknya seragam, meskipun strukturnya membolehkan campur.

### V6 `[terverifikasi]` — RBAC tidak seragam di lapisan UI

`pyPosition` hanya muncul di 17 berkas `Claim Life`; yang memuat gerbang sungguhan:
`ClaimLifeDetailGCNM.xml` (21), `Register_Flow.xml` (15), `AdjustmentDetail_Section.xml` (4),
`EditDateClaimLife_Section.xml` (4), `MedicalCheckClaimLife.xml` (2),
`InputAkseptasiClaimLife.xml` (2).

**Tidak** bergerbang peran: `InputOSClaimLife.xml`, `RejectOSClaimLife_Sec.xml`,
`ClaimComite.xml`, `Harness/Committe_Life.xml`. Penegakan peran sesungguhnya bertumpu pada
penugasan tahap di `Register_Flow`, bukan pada tiap layar. **ADR-0002** menetapkan tiga peran;
sistem baru perlu menegakkannya **di lapisan layanan**, bukan meniru ketidakseragaman ini.

### V7 `[terverifikasi]` — tidak ada jalur buka-ulang pada baris yang sama

Sensus penulis lengkap (register **OQ-061**): tidak ada rule yang menulis `0` setelah `1`/`2`.
Mendukung "terminal per baris".

### V8 `[terverifikasi]` — `SaveAdjustment_Act` **terpasang di UI**

Dirujuk 2× sebagai `<pyActivity>` dari `Claim Life/Section/ClaimLifeDetailGCNM.xml` (824.562 byte).
Ia menulis `.AdjustmentList(<LAST>).STS_REJECT = 1` bersama `ACCEPTEDNO` dan `ACCEPTATION_DATE`,
tanpa precondition Komite.

`[pertanyaan terbuka]` Sesuai instruksi work owner, **tidak diklaim dead**. Dicatat sebagai catatan
di **ADR-0011** dan masuk **OQ-039**.

---

## Yang dicatat

| Artefak | Perubahan |
| --- | --- |
| `docs/adr/0011-unit-status-adalah-baris-adjustmentlist.md` | **baru**, accepted |
| `docs/adr/0010-penyimpanan-berkas-tetap-di-google-storage.md` | draft → **accepted** |
| `CONTEXT.md` | mesin status, dua sumber reject, pencerminan, peran per tahap |
| `discovery/open-questions.md` | **OQ-061** dan **OQ-062** ditutup untuk Claim — Life; **OQ-039** dipersempit; **OQ-060** ditambah bukti; **OQ-063** dibuat |
| `.scratch/claim-life/kesiapan-to-spec.md` | dihitung ulang |

`spec.md` dan `issues/` **tidak disentuh** — menunggu persetujuan eksplisit.

---

# Lanjutan Ronde 3 — OQ-063

Tanggal: 2026-09-14 (pesan susulan work owner)

## Jawaban work owner

> OQ-063 → DIJAWAB SEBAGIAN.
>
> [terverifikasi] Gerbang visibilitas "send ke Komite" di
> Claim Life/Section/AdjustmentDetail_Section.xml (ASM-FW-GISFW-DATA-ADJUSTMENTLIFE /
> ADJUSTMENTDETAIL_SECTION) = pyWorkPage.pyPosition=='ReasLifeSPV' || pyWorkPage.Type='TP'
> || pyWorkPage.Type='TR'.
>
> [keputusan work owner] Maknanya: untuk klaim Type 'TP' atau 'TR', pengiriman ke Komite
> TIDAK dibatasi peran SPV — Admin pun dapat langsung mengirim ke Komite. Untuk tipe lain,
> hanya SPV. Ini melengkapi/meralat aturan "send ke Komite hanya SPV" di ADR-0011: berlaku
> untuk non-TP/TR saja.
>
> MASIH TERBUKA (jangan ditebak, pemilik Product+UW):
> - Arti Type 'TP' dan 'TR' (kepanjangan belum ada di korpus).
> - Ekspresi memakai '=' (assignment) bukan '==' pada Type='TP'/'TR'. Disengaja, atau bug
>   kurang kurung/operator yang membuat gerbang selalu lolos? Perlu konfirmasi apakah perilaku
>   "TP/TR bebas peran" memang diinginkan.

**Pemilik:** Product+UW (arti kode, aturan bisnis) + IAM (wewenang)

## Verifikasi korpus — butir "`=` bukan `==`" dapat ditutup

### V9 `[terverifikasi]` — tidak ada ambiguitas presedensi

Ekspresinya memakai `||` saja, tanpa satu pun `&&`. Pembacaan alternatif `SPV && (TP || TR)`
menuntut sebuah `&&` yang **tidak ada** di ekspresi itu. Jadi bukan soal kurang tanda kurung.

### V10 `[terverifikasi]` — `=` di korpus ini adalah **pembanding**, bukan assignment

Sensus korpus-wide:

| Tag | memakai `==` | memakai `=` tanpa `==` |
| --- | ---: | ---: |
| `<pyCondition>` | 954 | **2.282** |
| `<pyStepsPreCondParamsWhen>` | 16.390 | 273 |

Pada `<pyCondition>` — jenis tag yang dipakai gerbang ini — bentuk `=` justru **mayoritas**.

Lebih tegas: banyak ekspresi **mencampur keduanya dalam satu baris**, di tempat yang akan rusak
secara kasatmata kalau `=` berarti assignment:

| Ekspresi | Akibat kalau `=` assignment |
| --- | --- |
| `.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"` | **setiap** baris terpilih |
| `pyWorkPage.AcceptStatus = 1 && pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop` | cabang aksep **dan** tolak sama-sama jalan |
| `(.Type== 1 && …) \|\| ((.Type = 2 \|\|.Type = 4) && …) \|\| (.Type = 3 && …)` | `==` dan `=` pada properti **yang sama**; cabang 2/3/4 selalu benar |
| `pyWorkPage.pxFlow(InputOfferFacultativeIn).pxRouteTo="ReasFacInAdmin"` | assignment ke properti sistem yang bersifat baca |

**Kesimpulan:** `=` dan `==` dipakai bergantian sebagai pembanding. Gerbang itu memang bebas-peran
untuk `TP`/`TR` — **bukan** karena bug operator. Butir "cacat ekspresi" pada OQ-063 **ditutup**.

Perintah audit:
```
grep -rhoE "<pyCondition>[^<]*" . --include="*.xml" | grep -c "=="
grep -rhoE "<pyCondition>[^<]*" . --include="*.xml" | grep -v "==" | grep -cE "[A-Za-z0-9_) ]=[ ']"
```

### V11 `[terverifikasi]` — arti `TP`/`TR` memang tidak ada di korpus; pencarian sudah tuntas

- Tidak ada satu pun label, caption, atau opsi yang memasangkan `TP`/`TR` dengan teks penjelas.
- Satu-satunya tempat nilainya **ditetapkan**: `PremiumList Life/Activity/SubmitPremiumList_Act.xml`
  (`ASM-FW-GISFW-WORK-LIFE!SUBMITPREMIUMLIST_ACT`, 214.154 byte) — empat langkah `Property-Set`
  mengisi `InputData.CARI20` dengan `"QR"`, `"QP"`, `"TP"`, `"TR"`, masing-masing bergerbang
  `.Type=="<kode itu sendiri>"`. **Melingkar** — meneruskan, bukan mendefinisikan.
- `GetKodeProdLife_SQL` bukan pemetaannya: ia membaca `POOLDATA.KODE_PRODUKSI WHERE TYPE ='LIFE'`,
  kolom `TYPE` yang berbeda.

Artinya definisi berada **di luar korpus** (kemungkinan master data `POOLDATA`). Dicatat agar
pencarian ini tidak diulang. Berkaitan **OQ-020**, **OQ-059**, **OQ-001**.

## Yang dicatat

| Artefak | Perubahan |
| --- | --- |
| `docs/adr/0011-…` | bagian **"Ralat langkah 4 — `Type` `TP`/`TR` bebas peran"** ditambahkan; langkah 4 dan kalimat "Admin tidak mengirim ke Komite" diralat |
| `CONTEXT.md` | entri `ReasLifeAdmin`, `ReasLifeSPV`, dan `Penyerahan ke Komite` diralat; pengecualian `TP`/`TR` dicatat dengan buktinya |
| `discovery/open-questions.md` | **OQ-063** dijawab sebagian (butir "cacat ekspresi" ditutup dengan bukti korpus); **OQ-020** ditambah catatan bahwa `.Type` belum terjawab dan kini menggerbangi wewenang |
| `.scratch/claim-life/kesiapan-to-spec.md` | dihitung ulang |
