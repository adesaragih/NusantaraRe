# 01: Terima kasus dari Claim Life + inbox komite per posisi

**Status:** ready-for-agent

**Blocked by:** **00 (skema penyimpanan komite — PREFACTOR)** · **CL-01** (kerangka aplikasi + seam API) · **kontrak muatan penyerahan CL-10**
— bukan penyelesaian CL-10. Lihat §"Kontrak muatan penyerahan".

## Hasil & nilai pengguna

Sebagai **anggota komite**, kasus yang diserahkan Claim — Life muncul di **inbox saya** — dan hanya
kasus pada tingkat yang saya miliki — lengkap dengan rincian klaim dan baris `AdjustmentList` yang
diputuskan. Saya tahu apa yang harus saya kerjakan tanpa melihat pekerjaan orang lain.
*(User story 1, 2, 7 di spec)*

Ini **tracer bullet** konteks Komite: menembus React → HTTP → handler → service → repository →
Oracle, memakai ulang seam yang sudah ditetapkan **CL-01**.

## Area codebase

`internal/models` (kasus komite + koleksi `KomiteList` per tingkat), `internal/repository` (pembacaan
kasus + roster), `internal/services` (penyaringan inbox per posisi), `internal/handlers` (endpoint
inbox + detail kasus), `frontend/` (halaman inbox komite + layar detail).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/Flow/KomiteLife_Flow.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITELIFE_FLOW` / `RULE-OBJ-FLOW` | `[terverifikasi]` Assignment "KomiteRouter", `pyImplementation = WorkList`, `pyRouteTo = Custom` — antrean **per-pengguna** |
| `Komite Claim Life/FlowAction/ViewTransferDtl.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `VIEWTRANSFERDTL` / `RULE-OBJ-FLOWACTION`, 31.963 byte | layar detail kasus |
| `Komite Claim Life/Section/ShowTransfer.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SHOWTRANSFER` / `RULE-HTML-SECTION`, 1.125.234 byte | bentuk tampilan kasus — **berkas terbesar modul** |

`[terverifikasi]` Kelas kasus: `ASM-FW-GCNMFW-Work-KomiteLife`. Muatan penyerahan membawa `CLMNO`,
`KomiteCount`, `KomiteLoop`, `IndexAdjustment`, `IndexPremiumList`, dan `KomiteList` berisi
`KomiteID` / `IDKomite` / `KomiteAproval` / `KomiteEmail` — ditambah nilai klaim, `CURRENCY`, dan
`STS_REJECT` saat penyerahan (**ADR-0001**, keputusan Ronde 2 Q15 Claim Life).

## ADR terkait

**ADR-0001** (kontrak batas — penyerahan masuk), **ADR-0014** (inbox hanya menampilkan kasus sesuai
posisi roster), **ADR-0003** (uang non-float pada tampilan nilai klaim), **ADR-0011** (yang
diputuskan adalah **baris `AdjustmentList`**, bukan klaim).

## Acceptance criteria

- [ ] Kasus yang diserahkan Claim — Life dapat dibuka lewat API dan menampilkan klaim beserta
      **baris `AdjustmentList`** yang diputuskan.
- [ ] Inbox seorang anggota komite **hanya** memuat kasus pada tingkat yang ia miliki — bukan
      seluruh antrean komite. *(AC 10 spec)*
- [ ] Muatan penyerahan yang diterima memuat nilai klaim dan `CURRENCY`; nilai uang ditampilkan
      tanpa melewati *binary floating point*. *(AC 17 spec)*
- [ ] `KomiteLoop` yang diterima dari muatan dipakai apa adanya sebagai batas tangga — **tidak**
      dihitung ulang di konteks ini.
- [ ] ⚠️ Muatan ber-`KomiteLoop < 1` atau ber-`KomiteList` kosong **ditolak di lapisan layanan**;
      `T_GENERAL_KOMITE` dan `T_KOMITE_KOMITELIST` **tidak dibuat**, dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID`
      tidak terisi. **Menolak bukan menghitung ulang** — checkbox di atas tetap berlaku.
      *(AC 4 spec; `[keputusan work owner]` — penjaga berlapis; `[terverifikasi]` Pega tidak
      memuat gerbang ini di sisi mana pun)*
- [ ] Status ditampilkan sebagai kata (Outstanding / Aksep / Ditolak), **bukan** nama field
      `STS_REJECT` dan bukan angka. *(AC 29 spec)*
- [ ] Ada satu test ujung-ke-ujung yang menggerakkan sistem lewat HTTP dan memeriksa hasilnya lewat

### Penyimpanan kasus komite ⚠️ BARU 2026-09-16 — spec §9

- [ ] ⚠️ Menerima penyerahan **membuat `T_GENERAL_KOMITE`** (header) **beserta satu baris
      `T_KOMITE_KOMITELIST` per anggota roster**, masing-masing ber-`KOMITE_APROVAL = 0` dan
      ber-`KOMITE_URUT` sesuai jenjangnya. *(AC 30 spec; penyimpangan sadar 1)*
- [ ] ⚠️ Menerima penyerahan **melahirkan baris baru di `T_WORK_CLAIM`** — `ID` = identitas kasus
      komite (**teks berformat `KMT-xxxxxx`**), `COVER_KEY` = `ID` baris klaim — dan header
      `T_GENERAL_KOMITE` **memakai `ID` yang sama persis** (**shared primary key**).
      **REVISI 2026-09-18:** ⛔ **tidak ada kolom `WORK_CLAIM_ID`** — dibuang; hubungannya dijamin
      oleh ID identik, bukan kolom penyambung. Test yang menemukan kolom `WORK_CLAIM_ID` **gagal**.
      *(§9 spec; REVISI 2026-09-17 — menggantikan `T_WORK_CLAIM.KMT_NO` yang dibuang;
      REVISI 2026-09-18 — shared PK, `[keputusan work owner]`)*
- [ ] ⚠️ **Penunjuk dua arah konsisten dalam SATU transaksi**: `T_GENERAL_KOMITE.ADJUSTMENT_ID` dan
      `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` terisi bersama. Test yang menemukan salah satunya kosong
      sementara yang lain terisi **gagal**. *(AC 32 spec)*
- [ ] `KOMITE_LOOP` diisi **jumlah tingkat** hasil hitung roster aktif; `KOMITE_COUNT` dimulai pada
      tingkat pertama. *(AC 30 spec; tiket 02)*
- [ ] ⚠️ Roster dibaca dari master `EMAILKOMITE`; master itu **tidak ditulis**. *(`[data DBA]`)*
      HTTP, terhadap skema uji Oracle.

## Kontrak muatan penyerahan — disepakati di muka

Tiket ini bergantung pada **bentuk data** yang diserahkan Claim — Life, **bukan** pada selesainya
CL-10. Bentuk itu ditetapkan di sini agar kedua konteks dapat dikembangkan **paralel**, dengan
muatan penyerahan **di-fake di seam** sampai CL-10 nyata.

| Bagian muatan | Isi |
| --- | --- |
| Penunjuk baris | `IndexAdjustment`, `IndexPremiumList` — baris `AdjustmentList` yang diputuskan |
| Nilai klaim | jumlah klaim yang menjadi dasar pita roster |
| `CURRENCY` | mata uang nilai klaim (**ADR-0003**) |
| `KomiteList` | roster tingkat — tiap entri berisi `KomiteID`, `IDKomite`, `KomiteAproval`, `KomiteEmail` |
| `KomiteLoop` | jumlah tingkat tangga — **sudah dihitung** Claim Life |
| Status baris saat serah | `STS_REJECT` pada saat penyerahan |

`[terverifikasi]` Bentuk ini terbaca dari `Claim Life/Activity/CreateKMTLife_Act.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`, 121.652 byte) —
properti `childPageKomite.*` sebelum `Call pxAddChildWork`. Tiga tambahan (nilai klaim, `CURRENCY`,
`STS_REJECT`) adalah `[keputusan work owner]` dari Ronde 2 Q15 Claim Life (**ADR-0001**).

**Perubahan pada bentuk ini adalah perubahan kontrak lintas konteks**, bukan perubahan internal.

## Catatan — batas kepemilikan dengan Claim Life

`[terverifikasi]` **Perhitungan roster dan `KomiteLoop` adalah kode Claim Life, bukan Komite.**
`Claim Life/Activity/GetListKomiteLife.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
`GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`) dan `Claim Life/Activity/CreateKMTLife_Act.xml` keduanya
berada di modul `Claim Life`.

Konsekuensi: AC "COUNT roster", "nilai mutlak klaim negatif", dan "gagal bila roster kosong" adalah
AC **CL-10**. Di konteks ini ia **ekspektasi kontrak** — Komite **menerima** `KomiteList` dan
`KomiteLoop` yang sudah dihitung, dan **tidak menghitung ulang**.

`[terbuka]` **OQ-007 / OQ-021** — pemetaan `KomiteID` → identitas akun bergantung pada model RBAC
yang belum ditetapkan. **Asumsi tiket ini: satu `KomiteID` memetakan ke satu identitas akun.**

`[terbuka]` **OQ-035** — `UpdateWorkObject` dan `serviceInsertArasapasClaimLife_act` dipanggil dari
Komite tetapi salinannya ada di modul lain. **Bukan pemblokir isi**; perilakunya sudah terbaca.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
