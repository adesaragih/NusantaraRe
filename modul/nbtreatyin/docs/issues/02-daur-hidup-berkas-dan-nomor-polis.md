# 02: Daur hidup berkas realisasi dan nomor polis — satu nomor, sekali, tanpa bentrok

**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
**Blocked by:** 01
**Menutup:** AC 31 · 59 · 73 · 74 *(4 AC)* — US 1 · 4 · 5

## Hasil & nilai pengguna

Hari ini admin treaty membuka penawaran yang sudah disetujui dan melengkapinya menjadi realisasi.
`[terverifikasi]` Nomor polis dibentuk dari sebuah deret, dan ⚠️ salah satu bagiannya diambil dari
slot parameter generik yang **tidak terlihat diisi** di aktivitas mana pun.

Sesudah tiket ini, sebuah realisasi treaty **lahir dari penawaran yang disetujui**, mendapat
**satu nomor polis yang tidak bentrok**, dan ⭐ pengguna **diperingatkan** ketika ia hendak membuat
penawaran yang sudah pernah ada.

## Area codebase

- Lapisan service: daur hidup berkas realisasi
- Lapisan service: pembentukan nomor polis
- Lapisan handler: peringatan penawaran ganda

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Titik masuk alur | `Flow\InputRealizationTreatyIn.xml` — `pyStartActivity` = `Start1` |
| Pembentukan nomor polis | `RDBList\GenerateNoPolicy.xml` |
| Peringatan penawaran ganda | `Activity\CheckDuplicateOffer.xml` langkah 6 dan 8, lewat `RDBList\GetCountClaim.xml` |
| Rantai pemanggilnya | `TreatyRealizationCheckXOLList` → `SetTreatyIn_Act` → `CheckDuplicateOffer` |

## ADR terkait

- **ADR-0006** — penomoran lewat deret basis data

## Acceptance criteria

- [x] **AC 73** — nomor polis memuat awalan tetap, penanda treaty, bulan-tahun, dan nomor urut
      berdigit tetap
- [x] **AC 74** — nomor dibentuk **sekali** per berkas; tidak berubah pada penyimpanan berikutnya
- [ ] 🟡 **AC 31** — dua berkas **tidak pernah** bernomor sama
- [ ] ⛔ **AC 59** — peringatan muncul ketika jumlah berkas klaim terhubung **lebih dari nol**

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **P7** | asal salah satu bagian nomor polis **tidak terlihat diisi** | ⚠️ tidak menahan — polanya diketahui |

## Perintah verifikasi

1. Realisasikan dua penawaran berbarengan — ⭐ nomor polisnya **berbeda**.
2. Simpan ulang salah satunya — ⭐ nomornya **tidak berubah**.
3. Buat penawaran yang sudah pernah ada — ⭐ **peringatan muncul**.

## ⛔ RALAT implementasi 2026-10-03

1. **Rule pembentuk nomor polis.** Bunyi lama: *"Pembentukan nomor polis | `RDBList\GenerateNoPolicy.xml`"*.
   ⛔ Keliru: `GenerateNoPolicy` dipanggil `SaveJsonPolisTreatyIn_Act` langkah 4 dan **hasilnya tidak
   dipakai**. Nomor polis dibentuk `Activity\GeneratePolicyNoTreaty_Act` (langkah efektif 3-13, 25-30;
   14-24 berlabel `//`): `KODE_PRODUKSI(NONLIFE) + QR/QP/TP + ".T" + OJKBusinessID + "." + MM.YYYY +
   "." + urut5`, urut dari `PROC_GENERATE_SEQUENCE_NUMBER` (kini `inti/backend/penomor`). Langkah 28
   bersyarat `PolicyNo == ""` — nomor sekali (AC 74).
2. **P7 terjawab.** `InputData.CARI20` diisi langkah 7-9 (`DueTo` 1 → QR, 0 → QP, `ClaimType`
   "XOL Retro" → TP).
3. **AC 59 tidak dapat dipenuhi seperti tertulis.** `CheckDuplicateOffer` langkah 1-4 berlabel `//`
   dan dipanggil `SetTreatyIn_Act` (pembongkar JSON, tidak dimigrasi — AC 62) TANPA parameter, sehingga
   `GetCountClaim` selalu menghitung `masterid = NULL` — peringatan **tidak pernah menyala** di sistem
   lama. Yang dibangun: `TreatyRealizationCheckDuplicate` (pasca-submit admin, IsApproved 1).
