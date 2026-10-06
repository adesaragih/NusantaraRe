# 06: Penggolongan jenis usaha — 36 baris, berhenti di yang pertama cocok, bawaan UNKNOWN

**Status:** selesai *(putaran 2, 03-10-2026: tujuh AC dibangun, AC 66 tidak dibangun dengan bukti (a) —
RALAT di bawah; implementasi 2026-10-03 semula: sebagian; semula: ready-for-agent)*
**Blocked by:** 01
**Menutup:** AC 19 · 20 · 21 · 22 · 66 · 67 · 75 · 76 *(8 AC)* — US 25 · 26 · 32

## Hasil & nilai pengguna

Hari ini jenis usaha digolongkan otomatis dari kode bisnisnya lewat sebuah tabel keputusan.
⚠️ `[terverifikasi]` Keterangan yang ditulis manusia pada sebagian aturan penggolong **berbeda dari
syarat yang benar-benar dijalankan** — sembilan di antaranya bertentangan.

Sesudah tiket ini, penggolongan berjalan **persis seperti sistem lama**, ⭐ termasuk berhenti di
baris pertama yang cocok dan memberi nilai bawaan yang jelas ketika tidak ada yang cocok — sehingga
berkas tetap dapat diproses.

## Area codebase

- Lapisan service: penggolongan jenis usaha
- Fungsi murni: penilaian baris penggolong

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Tabel penggolong | `DecisionTable\BusinessType_DeT.xml` — **36 baris**, dua kolom penguji, **128** nilai kode seluruhnya unik |
| Cara evaluasi | `pyEvaluateAllRows` = `no` ⇒ ⭐ **berhenti di baris pertama yang cocok** |
| Nilai bawaan | ⭐ **`"UNKNOWN"`** |
| Penggolong lini jiwa | `When\IsLife.xml` — **16** nilai kode, dirujuk dua aktivitas penjaga |
| Pembeda jalur simpan | penanda bernilai `EDM` *(endorsemen)* / `POLICY` *(polis baru)* |
| Penanda penempatan keluar | `DataTransform\TestTreatyToFacStatus.xml` — dua kode angka |

## ADR terkait

- **ADR-0005** — aturan lingkungan tidak ditiru sebagai rule

## Acceptance criteria

- [x] **AC 19** — penggolongan **berhenti di baris pertama yang cocok**
- [x] **AC 20** — kode yang tidak cocok baris mana pun menghasilkan **`"UNKNOWN"`**
- [x] **AC 21** — ke-**128** kode menghasilkan penggolongan yang sama seperti sistem lama
- [x] **AC 22** — syarat diambil dari **yang dijalankan**, ⛔ bukan dari keterangannya — termasuk
      pada **9** aturan yang bertentangan
- [ ] ⛔ **AC 66** — penanda `EDM` ⇒ jalur endorsemen; `POLICY` ⇒ polis baru; ⛔ keduanya **tidak**
      menempuh cara penyimpanan yang sama — ⛔ **tidak dibangun, terbukti nol efek di XML** (RALAT
      putaran 2 di bawah, alasan (a))
- [x] **AC 67** — penggolong lini jiwa dimigrasi apa adanya, **16** nilai kode tanpa perubahan
- [x] **AC 75** — berkas non-proporsional ditandai sesuai jenis proporsinya
- [x] **AC 76** — penanda penempatan keluar dipasang **hanya** pada dua kode yang dikenal

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **10** | arti kode lini bisnis, termasuk satu kode berskala penomoran berbeda | tidak menahan |
| **11** | arti dua kode penanda penempatan keluar | tidak menahan |

## Perintah verifikasi

1. Jalankan ke-**128** kode melalui penggolong — ⭐ hasilnya **sama** dengan sistem lama, satu per satu.
2. Kirim kode yang tidak dikenal — ⭐ hasilnya **`"UNKNOWN"`**, bukan kosong dan bukan galat.
3. Kirim kode yang cocok **dua** baris — ⭐ yang dipakai baris **pertama**.

## Catatan

⚠️ `[keputusan work owner]` P23 menetapkan **keterangan diabaikan seluruhnya** — tidak dipakai
bahkan sebagai petunjuk. ⛔ Membangun dari keterangan akan salah pada **9 dari 50** penggolong yang
punya keduanya.

## ⛔ RALAT putaran 2 (P4, 03-10-2026) — AC 66 `isFOR`

Bunyi lama (spec AC 66): *"`isFOR` bernilai `EDM` ⇒ jalur **endorsemen**; `POLICY` ⇒ **polis baru**. Test
yang menemukan keduanya menempuh cara penyimpanan yang sama **gagal**."*

Bukti XML (graf keterjangkauan, `docs/alat/pemakai.py "\bisFOR\b"` di atas `graf.Graf(...).terjangkau()`:
278 rule, 176 terjangkau) — nama `isFOR` muncul di **tepat dua** rule, keduanya terjangkau:

| Rule | Peran | Isi |
| --- | --- | --- |
| `Flow\InputRealizationTreatyIn.xml` shape `Utility1` | satu-satunya pemanggil | `<pyCallParams><isFOR/></pyCallParams>` — **kosong** |
| `Activity\SaveJsonPolisTreatyIn_Act.xml` | satu-satunya pembaca | langkah 3 syarat `@equals(param.isFOR,"EDM")` ber-`pyStepsBlockName` **`//`** (mati); langkah 5 syarat `@equals(param.isFOR,"POLICY")` dengan kotak *When* **tidak dicentang** (`pyStepsPreCondition=false`) ⇒ syaratnya tidak pernah dievaluasi, langkahnya selalu berjalan |

Jadi `isFOR` **tidak memilih jalur apa pun**, bahkan di Pega: nilainya selalu kosong dan kedua syarat
yang membacanya tidak aktif; NB Treaty In selalu menempuh satu cara simpan (polis baru — generasi
baru `T_GENERAL_POLIS`, `OLD_POLIS_ID` kosong di NB, diagram R82). Endorsemen milik modul EDM.
`SaveJsonPolisTreatyIn_Act` sendiri diganti penyimpanan relasional (AC 16).

Bunyi baru: **"`isFOR` tidak dibangun — nol efek di XML (pemanggil tunggal mengirim kosong; syarat
EDM berlabel `//`, syarat POLICY tidak aktif). NB Treaty In hanya punya jalur polis baru."** Status:
⛔ tidak dibangun, alasan (a) bab 4 — cabang EDM/POLICY tidak terjangkau dari titik masuk nyata.
