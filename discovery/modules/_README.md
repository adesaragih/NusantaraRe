# modules/ — STEP D3

Status: **SELESAI** — 20/20 modul, disintesis dari D1 (inventaris) + D2 (telusur 20 konteks).

Setiap `<modul>.md` memakai **7 bagian**: peran modul; proses/fitur utama (merujuk
`../flows/<konteks>.md`); entitas & tabel data; integrasi eksternal; ketergantungan ke modul lain
(dengan bukti); batasan & batas pengetahuan; OQ yang menyentuh modul.
Setiap pernyataan berlabel `[terverifikasi]` / `[dugaan]` / `[pertanyaan terbuka]`; setiap angka
disertai perintah audit; nomor OQ diambil dari register `../open-questions.md`.

## Daftar modul

| Modul | File `.xml` | Titik masuk D2 | Catatan |
| --- | ---: | --- | --- |
| `NB FacIn.md` | 2.083 | 6 `Flow`; utama `InputInwardFacultativeOffer` | modul terbesar; graf terbesar korpus |
| `Endorsment Fac In.md` | 2.061 | `InputAddendumFacultativeIn` → `Start2` | taksonomi 12 tipe endorsement |
| `RNW Fac In.md` | 1.927 | `InputRenewalFacultativeIn` → `Start2` | 99,0 % identitas dibagi dgn NB FacIn |
| `Claim Fac In.md` | 482 | `Register_Flow` → `Start1` (`…WORK-PNC`) | 8 `ConnectREST`, terbanyak di korpus |
| `Treaty In Adjustment.md` | 379 | **Harness** `InputTreatyInAdjustment` | 277/323 file identik dgn Treaty In |
| `Treaty In.md` | 329 | **Harness** `InputTreatyInOffer` (8,2 MB) | mesin status `Akseptasi_DT` |
| `Treaty Contract Out.md` | 303 | **Harness** `InboxTreatyContract` | **bukan outward** (OQ-022) |
| `Claim Non Prop.md` | 279 | `Flow_TreatyIn` → `Start1` (`…CLAIMTREATYNONPROP`) | limit ter-hardcode; baca treaty outward |
| `NB Treaty In.md` | 278 | `InputRealizationTreatyIn` → `Start1` | blocker OQ-025 |
| `Claim Prop.md` | 270 | `Flow_TreatyIn` → `Start1` (`…CLAIMTREATY`) | limit dari database |
| `EDM Treaty In.md` | 163 | `InputAddendumTreatyIn` → `Start1` | memuat implementasi Arasapas |
| `Claim Life.md` | 136 | `Register_Flow` → `Start2` (`…CLAIMLIFE`) | paling terpisah (12–14 % overlap) |
| `PremiumList Life.md` | 124 | `InputPolicyHolder` (root modul) → `Start1` | OQ-004 terjawab |
| `Komite Claim FacIn.md` | 114 | `Komite_Flow` → `Start1` | ambang nominal ter-hardcode |
| `Master Product Name Life.md` | 114 | **Harness** `InwardProductName` | membawa jalur tulis `M_TREATY_IN` |
| `Komite Claim Prop.md` | 80 | `KomiteTreaty_Flow` → `Start1` (`…KOMITETREATY`) | 3 guard nama orang |
| `Endorsement Life.md` | 75 | `InputEDMLife` → `Start1` | flow paling sederhana |
| `Master Contract Retro Life.md` | 66 | **Harness** `InboxRetroLifeReinsurersList` | modul terkecil; paling terisolasi |
| `Komite Claim Non Prop.md` | 59 | `KomiteTreaty_Flow` → `Start1` (`…KOMITETREATYNONPROP`) | otorisasi berbasis roster |
| `Komite Claim Life.md` | 47 | `KomiteLife_Flow` → `Start2` | routing berbasis data |

Total **9.369 file**. Audit:
`for m in */; do :; done` — lihat `../D2-CLOSING-REPORT.md` §1.1 untuk daftar titik masuk lengkap.

## Keluaran lain STEP D3

- **`../glossary.md`** — glossary final (kandidat seed `CONTEXT.md` untuk FASE B).
- `../glossary-draft.md` — draf berjalan D1–D2, **dipertahankan sebagai jejak**.
- `../open-questions.md` — 57 terbuka / 2 terjawab, per pemilik peran.

## Langkah berikutnya

STEP D4: `../context-map.md` + `../understanding-report.md`, lalu **GATE manusia** sebelum FASE B.
