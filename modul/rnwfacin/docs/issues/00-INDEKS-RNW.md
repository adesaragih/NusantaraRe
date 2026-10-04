# Indeks Tiket — Siklus Renewal (RNW)

> **Delapan tiket, 14 berkas delta.** Sumber: `..\..\04-spec\04-spec-rnw.md` · register **K-030…K-041**
> · ADR-0001…0006.
>
> ⛔ **16 tiket New Business (`..\01`…`..\16`) tidak disentuh** dan tidak diulang. Renewal
> **memakai ulang seluruh modul NB** — perhitungan, registry predikat, tangga akseptasi, spreading —
> tanpa satu pun seam atau modul hitung baru.

**Dasar yang membuat lingkup ini kecil:** `[terverifikasi]` **1.907 dari 1.927** berkas RNW identik
**byte-per-byte** dengan NB (K-031). Delta sejati **14 berkas**, seluruhnya lapisan **masuk dan
tampilan**.

---

## Urutan dependency

| # | Slug | Blocked by | Berkas | Yang dihasilkan |
| ---: | --- | --- | ---: | --- |
| **R01** | `tracer-alur-masuk-renewal` | **NB-08** · **NB-11** | 3 | No. Polis + Renewal Date + Note → **OK** → kasus tercipta (`StatusBusiness=2`), data polis lama tersalin, **diterima tangga akseptasi tanpa penyesuaian** |
| **R02** | `layar-periode-renewal` | R01 | 2 | Periode polis baru + `OldPolicyNo` · `RNWDate` · pembeda siklus |
| **R03** | `layar-input-renewal-wadah` | R02 | 2 | Layar **menyisipkan** komponen NB, tidak menyalinnya |
| **R04** | `layar-detail-renewal` | R03 | 2 | 14 & 17 section disisipkan · `InputDtlObject` usang diport apa adanya |
| **R05** | `daftar-kandidat-renewal` | R01 | 2 | Daftar kandidat berkunci `OldPolicyNo` + portal. **Tanpa kelas kerja baru** |
| **R06** | `kelompok-bisnis` | R01 | 2 | Kode bisnis → tabel `business` → kelompok bisnis. `pySaveSQL` **dibuang** |
| **R07** | `konversi-produksi-renewal` | R01 · ⛔ **eksternal** | 1 | Konversi lewat jalur NB. **Belum dapat dieksekusi** |
| **R08** | `rekonsiliasi-kasus-renewal` | **NB-16** · R01 · R04 | — | Nol selisih memakai kerangka NB-16, **tanpa pembanding baru** |

**Cakupan:** 3 + 2 + 2 + 2 + 2 + 2 + 1 = **14 berkas** ✓

**Frontier — dapat dimulai segera setelah blocker NB-nya selesai:** **R01**. Setelah R01 hijau, empat
tiket terbuka sekaligus: R02, R05, R06, dan (secara formal) R07.

---

## Ketergantungan ke tiket New Business

| Tiket RNW | Tiket NB | Mengapa |
| --- | --- | --- |
| R01 | **NB-08** `registry-rules-eval` | `[terverifikasi]` flow renewal merujuk **19 rule `When`**, seluruhnya predikat routing yang diwarisi |
| R01 | **NB-11** `tangga-akseptasi-bentuk-a` | Kriteria penerimaan R01 adalah "kasus diterima tangga akseptasi" |
| R08 | **NB-16** `rekonsiliasi-eksak-tahap-1` | Renewal tidak punya pembanding sendiri; NB-16 sendiri diblok NB-15 (de-identifikasi) |

📌 Salah satu dari lima berkas kasus yang diurus **NB-15** adalah **kasus renewal nyata** — masukan
R08 sudah tersedia tanpa permintaan tambahan.

---

## ⛔ Satu blocker yang berada di luar kedua set tiket

**R07 tidak dapat dieksekusi** sampai jalur simpan produksi tersedia. Itu **Out of Scope butir 3 spec
NB**, menunggu **struktur tabel flat** (keputusan work owner) dan **`ALL_SOURCE`**. Tiketnya ditulis
penuh; hanya eksekusinya menunggu.

Celah **K-004** juga tetap terbuka: rule konversi asli kelas `Work` tidak ada di korpus mana pun.

---

## Keputusan yang tercermin di tiket

| Keputusan | Tercermin di |
| --- | --- |
| **K-031** renewal memakai ulang modul NB | seluruh tiket — tidak ada modul hitung baru |
| **K-034** gerbang masuk ditetapkan eksplisit (keputusan bisnis, bukan porting) | R01 |
| **K-035** konversi mengikuti pola NB | R07 |
| **K-037** `pySaveSQL` `PROSESCOPY` dibuang — perubahan perilaku yang disengaja | R06 |
| **K-040** empat section diwarisi dari NB, bukan menggantung | R03 |
| **K-041** `InputDtlObject` usang, diport apa adanya | R04 |
| **Keputusan work owner:** pasangan `_IsUW` = **dua komponen terpisah** | R02 · R03 · R04 |

⚠️ Keputusan `_IsUW` **sengaja berbeda** dari rekomendasi spec model-data NB §6 (satu komponen dua
mode). Alasannya reproduksi perilaku terekam (`CLAUDE.md` §1) — **bukan inkonsistensi**, dan dicatat
di ketiga tiket yang terkena.

---

## Di luar lingkup tiket RNW

| Butir | Alasan |
| --- | --- |
| Modul inti — perhitungan · registry · tangga akseptasi · spreading | Diwarisi NB (K-031); sudah punya tiket 01–16 |
| `EditMarketing` (Section + FlowAction) | Dibuang (K-036) |
| Fitur pilih-tertanggung (4 berkas) | Dibuang (K-038); tertanggung diambil dari polis lama |
| `SetOldData` / properti `TSIOld` | ⏸ Ditunda ke fase **Endorsement**. ⚠️ `TSIOld` (sufiks) ≠ `.OldTSI` |
| Celah K-004 — rule konversi asli kelas `Work` | Belum ada di korpus mana pun; arah = pola NB |

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
