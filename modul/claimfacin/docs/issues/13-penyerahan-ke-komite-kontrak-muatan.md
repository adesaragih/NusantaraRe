# 13: Penyerahan ke komite — kontrak muatan

**Status:** ready-for-agent
**Blocked by:** 05 (penyesuaian) · 06 (uang) · 08 (wewenang)
**Menutup:** AC 43 · 44 · 45 · 46 · 47 · 48 · 49 · 50 *(8 AC)* — US 13 · 14

## Hasil & nilai pengguna

Hari ini Penyesuaian di atas kewenangan penilai **belum dapat diserahkan ke komite**, dan tidak ada kesepakatan tentang **apa yang diserahkan**.

Sesudah tiket ini, Penilai dapat **menyerahkan penyesuaian ke komite**, dan ⭐ **muatan yang diserahkan lengkap** — termasuk data kutipan **utuh**, sehingga penggolongan lini usaha di sisi komite tidak pincang.

## Area codebase

- Lapisan layanan klaim: penyerahan ke komite
- ⭐ Kontrak muatan — apa yang disalin dan seberapa lengkap

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pembentukan kasus komite | rule pembuat nomor komite; menyemai jumlah jenjang dan giliran mulai |
| Jalur satu jenjang | rule kirim tutup klaim dan kirim tolak klaim — ⭐ keduanya **berkomite satu** |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0014** — wewenang ditegakkan di lapisan layanan

## Acceptance criteria

- [ ] **AC 43–50** — penyerahan ke komite
- [ ] ⭐ **Data kutipan disalin UTUH**, ⛔ bukan daftar medan bernama — ⚠️ daftar bernama itulah yang melahirkan cacat MBU dan Travel
- [ ] ⭐ Jalur **tutup klaim** dan **tolak klaim** membentuk komite **satu jenjang**
- [ ] ⛔ Tiket ini **berhenti di kontrak muatan** — ⭐ perilaku komite ada di putaran sisi komite

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"⭐ **Data kutipan disalin UTUH**, ⛔ bukan daftar
> medan bernama — ⚠️ daftar bernama itulah yang melahirkan cacat MBU dan Travel"* dan *"⭐ Jalur **tutup klaim** dan
> **tolak klaim** membentuk komite **satu jenjang**"* →
>
> - Premis *"cacat MBU dan Travel"* **bertentangan dengan AC 109** (properti klasifikasi benar-benar diisi). Tahap 1
>   **tidak menyalin data kutipan** ke kasus `KMT-`: kasus komite TT2 lahir di `T_WORK_CLAIM` (`LINI = 'FACIN'`,
>   `TAHAP Komite_Flow`) + `T_GENERAL_KOMITE` (`ADJUSTMENT_ID` = ID stabil adjustment) + `T_KOMITE_KOMITELIST`; muatan
>   outbox surel hanya pengenal (`PARITAS.md` §5). Apa yang dibaca sisi komite = **tahap 2** (`komiteclaimfacin`);
>   *"disalin UTUH"* **perlu dicek terhadap XML** Komite Claim FacIn sebelum dijadikan kontrak.
> - Jalur **tolak klaim (TT3)** dan **tutup tanpa bayar (TT4)** melahirkan kasus komite **tanpa adjustment**, padahal
>   `T_GENERAL_KOMITE.ADJUSTMENT_ID` NOT NULL — tombol **Yes** keduanya **nonaktif** (**OQ-CFI-27**); Close langsung
>   (tanpa CWP) dibangun. Bundel adjustment PT 2 ke satu `KMT-` tidak dibangun (**OQ-CFI-28**).
> - Komite mengikuti **pola Komite Claim Prop, tanpa menu** (OQ-CFI-04); nomor `KMT-` dari `SEQ_WORK_CLAIM`
>   (OQ-CFI-02).

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi daftar jabatan dan susunan jenjang | ⛔⛔ **MENAHAN** penentuan jumlah jenjang |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"Isi daftar jabatan dan susunan jenjang | ⛔⛔
> **MENAHAN** penentuan jumlah jenjang"* → **tidak menahan lagi**: tahap 1 menulis tangga awal dari roster `EMAILKOMITE`
> FACIN (`SetListKomite_act`: anggota aktif ber-`LIMIT_BOTTOM` ≤ batas; aturan khusus objek retro memakai `LIMIT_TOP`
> baris *Technic Div. Head*), dengan `T_WORK_CLAIM.POSITION` = `OPERATOR_ID` tingkat 1 (`backend/models/komite.go`,
> `repository/komite.go`). ⭐ **Pemutus menurut jabatan → workbasket**, pola Komite Claim Prop, di **tahap 2**
> (`komiteclaimfacin`, OQ-CFI-04). Wewenang membuka / menyerahkan = ADR-0030 + workbasket (RALAT tiket 08).

## Perintah verifikasi

1. Serahkan satu penyesuaian ke komite — ⭐ kasus komite lahir dengan muatan lengkap.
2. Periksa data kutipan di sisi komite — ⭐ **utuh**, bukan dua medan.
3. Serahkan lewat jalur tutup klaim — ⭐ komitenya **satu jenjang**.
