# 06: Jurnal balik — `Delete` selektif dan `Batal` menyeluruh

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 05 (penandaan `EDMStatus` harus sudah berjalan)

## Hasil & nilai pengguna

Sebagai **Finance**, saya ingin pembatalan polis dan penghapusan peserta menghasilkan **baris
bernilai negatif** — bukan baris nol dan bukan penghapusan — sehingga jejak transaksi asli tetap ada
dan nettonya dapat dihitung dengan menjumlahkan. *(User story 18–22 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Nilai uang sebagai desimal presisi arbitrer, **bertanda** |
| `internal/repository` | Penulisan baris negatif ke tabel yang sama |
| `internal/services` | **Mesin pembalikan tanda** — cakupan selektif versus menyeluruh |
| `internal/handlers` | Ringkasan endorsement menampilkan net |
| `frontend/` | Grid detail membedakan baris positif dan negatif secara kasatmata |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `SetPremi_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETPREMI_EDM` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/SetPremi_EDM.xml` (345.689 byte, **9 langkah, nol `<pyStepsBlockName>`**) |

`[terverifikasi]` **Peta langkah pembalikan:**

| Step | Deskripsi | Precondition | Akibat |
| --- | --- | --- | --- |
| **2.1** | **"Set 0 jika EDM Batal"** | `.EdmBatal=="True" \|\| pyWorkPage.EdmType==3` (baris 1273) | ⚠️ **mengalikan 32 kolom uang dengan `-1`** |
| 2.2 | "flag pengurangan" | `.EdmBatal=="True"` (1441) | `.EditInput = 1`, `.EDMStatus = "Delete"` (1384) |
| 2.3 | "flag batal" | `pyWorkPage.EdmType==3` (1583) / `.EdmBatal=="True"` (1636) | `.EDMStatus = "Batal"` (1506) |
| 5–7 | "Set COB", "Set property PremiumListSummary", "Insert to summary" | — | ringkasan ikut terbentuk |
| 9 | `Obj-Save` | — | |

### ⚠️ Deskripsi Pega **berbohong** terhadap kodenya

`[terverifikasi]` Deskripsi langkah berbunyi **"Set 0 jika EDM Batal"** — tetapi kodenya **tidak
menulis nol**. Ia menulis `<kolom> = <kolom> * -1` pada **32 kolom uang**:

```
SUM_INSURED, CEDING_RETENTION, SUM_REASURED, SHARE_NUSANTARA_RE, SHARE_NUSANTARA_RE_GROSS,
SUM_AT_RISK_GROSS, SUM_AT_RISK_RETRO, RETROCEDED_SHARE, SHARE_RETRO, RATE, FACTOR,
GROSS_PREMIUM, NET_PREMIUM, DEDUCTION, CLAIM_AMOUNT, RI_ADMIN_FEE, BROKERAGE_FEE,
beserta seluruh kelompok *_REFUND, *_RETRO, dan *_REFUND_RETRO
```

**Baca kodenya, jangan namanya.** `[keputusan work owner]` pembalikan tanda memang perilaku
akuntansi yang dikehendaki.

## ADR terkait

**ADR-0003** (uang non-float, desimal presisi arbitrer — **berlaku penuh pada nilai negatif**),
**ADR-0011** (status per baris), **ADR-0001** (kontrak hilir).

## Acceptance criteria

- [ ] ⚠️ **Penyimpangan sadar — jurnal balik.** Baris `Delete` dan `Batal` menghasilkan nilai
      **negatif**, **bukan nol** dan **bukan penghapusan**. *(AC 16 spec; `[keputusan work owner]`)*
- [ ] **Seluruh 32 kolom uang** ikut dibalik tandanya; tidak ada komponen yang tertinggal positif.
      Test memeriksa **kolom demi kolom**. *(AC 17 spec)*
- [ ] Baris positif asli **tetap ada** berdampingan dengan baris negatif di tabel yang sama.
      *(AC 18 spec)*
- [ ] **Jumlah** baris positif dan negatif untuk peserta yang dibatalkan = **nol**. *(AC 19 spec)*
- [ ] Cakupan minus benar: **`Delete` menegatifkan hanya peserta terpilih**; **`Batal` menegatifkan
      seluruh peserta**. *(AC 20 spec)*
- [ ] Nilai uang **tidak** melewati `float` di lapisan mana pun maupun di JSON API — **termasuk
      nilai negatif**. *(AC 21 spec; **ADR-0003**)*
- [ ] Nilai yang ditulis dan dibaca kembali **identik**; tidak ada pembulatan diam pada nilai negatif
      maupun pada nilai berpecahan panjang. *(AC 22 spec)*
- [ ] Ringkasan premium endorsement mencerminkan **net** (positif + negatif), bukan hanya salah satu
      sisi.
- [ ] Baris `Old` yang tidak ditandai `Delete` **tidak ikut dibalik** pada endorsement Perubahan

### Hapus = flag + minus ⚠️ BARU 2026-09-16 — spec §16

- [ ] ⚠️ Peserta yang dikeluarkan ditandai **`Delete`** dan nilainya menjadi **pengurang penuh** —
      **barisnya tetap tersimpan**. Test yang menemukan penghapusan fisik baris peserta **gagal**.
      *(AC 66 spec; penyimpangan sadar 11)*
- [ ] ⚠️ `EdmType` **`3`** (Batal) menjadikan **setiap** peserta pengurang penuh dan menandainya
      **`Batal`**, **tanpa menghapus satu baris pun**. *(AC 67 spec; penyimpangan sadar 11)*
- [ ] `EdmType` `1` dan `3` memakai **jalur simpan yang sama**; yang berbeda hanya nilai dan flag.
      *(AC 68 spec)*
- [ ] Kolom uang **menerima nilai negatif** tanpa kehilangan presisi — jurnal balik memang menulis
      minus. *(**ADR-0003**; spec §16)*
      Data.

## Catatan — apa yang dilihat siapa

Baris negatif ini **masuk ke tabel yang sama** yang dibaca Claim Life. Aturan hilirnya: **akuntansi
melihat seluruh baris; klaim hanya peserta hidup.** Penegakan penyaringan ada di konteks Claim Life
(**CL-02**, spec Claim Life §16 AC 25–30) dan diikat sebagai kontrak di **PL-08** serta tiket **11**
konteks ini.

## Blocker

**Tidak ada.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
