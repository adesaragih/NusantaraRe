# Grilling — Claim — Life — Ronde 4 (penutup frontier)

Status: answered (work owner, 2026-09-14) — **frontier Claim — Life kosong**
Konteks: `claim-life` (Claim — Life)
Tanggal: 2026-09-14
Skill: `/mattpocock-skills:grill-with-docs` (grilling + domain-modeling)
Ronde sebelumnya: `grilling-ronde-1.md` (Q1–Q7), `grilling-ronde-2.md` (Q8–Q15),
`grilling-ronde-3.md` (Blok A + OQ-063)

> Work owner menjawab langsung empat butir sisa dari `kesiapan-to-spec.md`.
> Setiap klaim yang **tidak** terbukti di korpus ditandai `[keputusan work owner]`, bukan
> `[terverifikasi]`.

---

## 1. OQ-039 — akibat *reject* oleh Admin

**Jawaban work owner** `[keputusan work owner]`:

> Reject oleh Admin ("Reject Outstanding") MEMBATALKAN BARIS AdjustmentList ITU SAJA — klaim TIDAK
> tertutup; Admin lalu input baris adjustment baru. Memakai STS_REJECT = 2 (nilai SAMA dengan reject
> Komite). Konsekuensi: STS_REJECT = 2 SELALU berarti "baris ini ditolak", TIDAK PERNAH "klaim
> selesai" — konsisten dgn ADR-0011 (klaim tidak terminal, terminal hanya per baris). Dua sumber
> penolakan (Admin langsung / Komite lewat SPV) menulis nilai yang sama dengan makna operasional
> setara pada tingkat baris. Tutup OQ-039 untuk Claim Life.

**Pemilik:** Product+UW
**Dicatat ke:** **ADR-0011** §"Arti tunggal nilai `2`" + langkah 1b; `CONTEXT.md` (`STS_REJECT`,
`ReasLifeAdmin`, `AdjustmentList`); **OQ-039** ditutup untuk Claim — Life.

`[terverifikasi]` Konsisten dengan korpus: `RejectOSClaimLife_Act` dan `KomitePostAdjustment`
menulis `2` ke tingkat baris **dan** ke `PremiumListDetail(idx)` dalam bentuk yang sama persis —
korpus memang tidak membedakan keduanya, dan menurut keputusan ini memang tidak seharusnya berbeda.

**Konsekuensi yang saya tarik dan perlu Anda ketahui:** karena tidak ada nilai `STS_REJECT` yang
berarti "klaim selesai", **selesainya sebuah klaim menjadi keadaan turunan** dari kumpulan barisnya,
bukan status tersimpan. Aturan turunannya belum ditetapkan — itu keputusan desain, bukan fakta yang
hilang, jadi dicatat di `kesiapan-to-spec.md`, bukan sebagai OQ baru.

---

## 2. Arti `Type` (OQ-020 / OQ-063)

**Jawaban work owner** `[keputusan work owner]`:

> TP = Payable, TR = Receivable. (Definisi ini tidak ada di korpus — sumber work owner.)
> Perbarui CONTEXT.md dan register. Type menggerbangi DUA hal terverifikasi: (a) wewenang
> send-Komite (TP/TR bebas peran, non-TP/TR hanya SPV) dan (b) jendela validasi DOL di
> ValidasiDOL_Act (QP/QR -> GROSS_VALUATION_*; TP/TR -> RETROCESSION_VALUATION_*).
> Tutup butir "arti TP/TR" pada OQ-063.

**Pemilik:** Product+UW
**Dicatat ke:** **ADR-0012**; `CONTEXT.md` (entri baru `Type`, `TP`/`TR`, dan catatan `QP`/`QR`);
**OQ-063** ditutup; **OQ-020** diperbarui — **`QP` dan `QR` tetap terbuka**.

### V12 `[terverifikasi]` — kedua gerbang membaca **salinan berbeda** dari nilai yang sama

Ini nuansa yang perlu diketahui sebelum menulis spec:

| Yang digerbangi | Properti yang dibaca |
| --- | --- |
| Wewenang send-Komite (`AdjustmentDetail_Section.xml`) | `pyWorkPage.Type` |
| Jendela validasi DOL (`ValidasiDOL_Act.xml`) | `pyWorkPage.PolicyDataLife.Type` |

Penyalinannya: `Claim Life/Activity/LoadDataPeserta_Act.xml` →
`pyWorkPage.Type = pyWorkPage.PolicyDataLife.Type`. Sumber otoritatif adalah
**`PolicyDataLife.Type`** (tipe polis); `pyWorkPage.Type` hanya salinan kerja. Di sistem baru
keduanya menjadi satu field — tetapi perlu disadari bahwa di Pega mereka dua, sehingga secara teori
dapat berbeda bila `LoadDataPeserta_Act` tidak berjalan.

### V13 `[terverifikasi]` — cabang validasi DOL, kutipan penuh

| Cabang | Jendela | Pergeseran tanggal |
| --- | --- | --- |
| `Type=="QR" \|\| Type=="QP"` | `GROSS_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `@addCalendar(.DATE_OF_LOSS,0,0,0,0,0,0,0)` — **nol** |
| `Type=="TR" \|\| Type=="TP"` | `RETROCESSION_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `@addCalendar(.DATE_OF_LOSS,0,0,0,1,0,0,0)` — **+1 hari** |

Gagal → `local.errmsg = "Invalid DOL"`, dengan `local.Begin==false || local.Expired==true`.

**`[dugaan]` yang saya catat di OQ-020 — jangan dipakai sebagai fakta:** pola di atas
memperlihatkan huruf **pertama** memisahkan gross (`Q*`) dari retrosesi (`T*`). Bila huruf **kedua**
memang Payable/Receivable seperti pada `TP`/`TR`, maka `QP`/`QR` mengikuti pola yang sama — dan satu
konfirmasi Anda akan menutup sisa OQ-020 untuk kode ini. Belum saya tanyakan.

---

## 3. Kelonggaran `TP`/`TR` — keputusan desain

**Jawaban work owner** `[keputusan work owner]`:

> PARITAS + tandai risiko: Bawa perilaku apa adanya (untuk TP/TR, pengiriman ke Komite tidak
> dibatasi peran SPV). JANGAN diperketat di migrasi ini. Catat sebagai RISIKO RBAC di ADR terkait
> (wewenang bergantung Type), untuk ditinjau saat konteks Komite/IAM digarap. OQ-063 -> ditutup.

**Pemilik:** Product+UW + IAM
**Dicatat ke:** **ADR-0012** (baru, accepted) — keputusan paritas beserta tiga risiko yang diterima.

---

## 4. `SaveAdjustment_Act`

**Jawaban work owner** `[keputusan work owner]`:

> Claim Life/Activity/SaveAdjustment_Act.xml (set STS_REJECT=1 tanpa precondition Komite) SUDAH
> TIDAK DIPAKAI. Jangan dimigrasikan. Status dead = keputusan work owner, tidak dapat diverifikasi
> dari korpus -> tandai [keputusan work owner]. (Penulis STS_REJECT=1 yang berlaku = rule sisi
> Komite KomitePostAdjustment.)

**Pemilik:** Product+UW
**Dicatat ke:** **ADR-0011** §`SaveAdjustment_Act`; `CONTEXT.md`; **OQ-039** butir 3.

⚠️ Dicatat apa adanya, termasuk pertentangannya dengan korpus: `[terverifikasi]` rule itu
**terpasang di UI** — dirujuk 2× sebagai `<pyActivity>` dari
`Claim Life/Section/ClaimLifeDetailGCNM.xml` (824.562 byte). Pertentangan ini disimpan **sengaja**
sebagai titik mula bila keputusan "jangan dimigrasikan" perlu ditinjau ulang.

---

## Yang dicatat

| Artefak | Perubahan |
| --- | --- |
| `docs/adr/0012-wewenang-kirim-komite-bergantung-type.md` | **baru**, accepted — paritas + risiko RBAC |
| `docs/adr/0011-…` | langkah 1b diralat; §"Arti tunggal nilai `2`" ditambahkan; §`SaveAdjustment_Act` menggantikan "Catatan yang belum berstatus"; tabel OQ diperbarui |
| `CONTEXT.md` | entri baru `Type`, `TP`/`TR`, `QP`/`QR`; `STS_REJECT`, `ReasLifeAdmin`, `AdjustmentList`, `Penyerahan ke Komite` diralat; `SaveAdjustment_Act` ditandai dead |
| `discovery/open-questions.md` | **OQ-039** ditutup untuk Claim — Life (3 modul Claim lain tetap terbuka); **OQ-063** ditutup penuh; **OQ-020** diperbarui — `TP`/`TR` terjawab, `QP`/`QR` tetap terbuka |
| `.scratch/claim-life/kesiapan-to-spec.md` | dihitung ulang → **MATANG** |

`spec.md` dan `issues/` **tidak disentuh** — menunggu persetujuan eksplisit.
