# Arsip — discovery lintas-siklus (Prompt 2)

Dokumen di folder ini **tidak dibatalkan**. Ia diarsipkan pada 15 September 2026 pukul 18:01 WIB
atas keputusan work owner untuk mengalihkan fokus ke **satu siklus dahulu: NB FacIn**.

Isinya dibangun dari ketiga folder korpus sekaligus (`NB FacIn`, `RNW Fac In`, `Endorsment Fac In`).
Temuan **lintas-siklus** di sini tidak diproduksi ulang oleh pass NB, dan tetap menjadi rujukan:

| Temuan | Tetap berlaku |
| --- | --- |
| Tiga siklus berjalan di atas **satu basis rule** — 1.220 dari 1.680 identitas rule bersama berlogika identik | ya |
| **NB dan RNW tidak pernah berbeda satu sama lain** (1.680/1.680; UI 575/575) | ya |
| ⛔ **Salinan EDM memuat versi rule yang tidak sepadan** — lihat koreksi di bawah | ya, **dan memblokir — tetapi hanya untuk EDM** |
| ⛔ Rule `When` yang kondisinya tidak terbaca = **0 dari 601** — mengoreksi `CLAUDE.md` §4.5 | ya |
| 45 butir memblokir + 29 perlu keputusan bisnis | ya |
| Tiga jebakan metodologis pembandingan korpus | ya |

---

## Koreksi R0 — diukur ulang 15 September 2026, 18:2x WIB

Rumusan asli R0 ("ketiga folder diekspor dari titik waktu berbeda") **terlalu longgar** dan sebagian
salah sasaran. Pengukuran ulang langsung atas `<pyRuleSetVersion>` dan `<pxCommitDateTime>` **di
dalam** berkas — bukan atas tanggal berkas — memberi gambaran yang lebih tajam:

| Ukuran | Nilai |
| --- | ---: |
| Rule bernama sama hadir di NB **dan** EDM | 1.719 |
| — berbeda versi ruleset | **95** |
| — salinan EDM lebih **tua** | 62 |
| — salinan EDM lebih **baru** | 33 |
| Rule bernama sama NB vs RNW berbeda versi | **0** |

```powershell
# perbandingan versi: ekstrak <pyRuleSetVersion> dan <pxCommitDateTime> per berkas,
# kelompokkan menurut (tipe rule + nama berkas), bandingkan antar folder.
# Contoh terverifikasi — Activity\SetToInbox_ACT.xml:
#   NB  = 01-01-95, commit 2026-08-07
#   RNW = 01-01-95, commit 2026-08-07
#   EDM = 01-01-87, commit 2025-09-18   <- selisih 11 BULAN pada rule-nya sendiri
```

**Apa yang berubah dari rumusan lama:**

1. Ini **bukan** artefak tanggal ekspor. Ketiga folder diekspor dalam rentang 8 hari. Yang tidak
   sepadan adalah **versi rule di dalamnya** — pada contoh di atas, selisihnya 11 bulan.
2. **Arahnya tetap tidak konsisten** (62 lebih tua, 33 lebih baru), jadi ini bukan sekadar "EDM
   tertinggal". Sebagian salinan EDM justru mendahului NB.
3. ⚠️ **Ini masalah EDM, bukan masalah NB.** `NB vs RNW = 0` — tidak satu pun rule bernama sama
   berbeda versi antara NB dan RNW. **Konsekuensinya penting: pekerjaan siklus NB dapat dilanjutkan
   dan direkonsiliasi tanpa menunggu pertanyaan ekspor ini dijawab.** Yang tertahan hanyalah EDM.

Angka lintas-folder lain di arsip ini (mis. 460 rule "bercabang") tetap merupakan **batas atas**:
sebagian selisihnya berasal dari 95 rule di atas, bukan dari percabangan bisnis.

---

## Koreksi R1 — 21 September 2026 · ⛔ **nama berkas yang sama TIDAK menjamin rule yang sama**

> **Rumusan lama di atas tidak dihapus.** Ia tetap berlaku untuk apa yang diukurnya. Yang ditambahkan
> di sini adalah **syarat yang mendahuluinya**.

⛔ **Pembandingan lintas folder WAJIB memeriksa `pzInsKey` dan `pyClassName` LEBIH DULU.**

`[terverifikasi]` Blok kode di atas mengelompokkan menurut **"(tipe rule + nama berkas)"**. Itu
**tidak cukup**: dua berkas bernama sama di folder berbeda dapat merupakan **DUA RULE BERBEDA** di
kelas berbeda — bukan dua versi dari satu rule.

**Contoh terverifikasi** — `Activity\CountRateRetroCov.xml`:

| Folder | `pyClassName` | Basis `pzInsKey` |
| --- | --- | --- |
| NB · RNW | `ASM-FW-GISFW-Data-PropertyItem` | `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-PROPERTYITEM COUNTRATERETROCOV` |
| Endorsment | **`ASM-FW-GISFW-Data-Aneka`** | **`RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-ANEKA COUNTRATERETROCOV`** |

Keduanya **rule berbeda**, keduanya **berlaku**, dipilih menurut COB (**K-060 amandemen**). Pembagi
preminya pun berbeda — **100.000** vs **10.000**.

### Skala masalahnya

`[terverifikasi]` Sapuan seluruh korpus (`10-audit\10-rule-kembar-pzinskey.md`), 5.288 pasangan
berkas bernama sama:

| Golongan | Jumlah |
| --- | ---: |
| **G1** kunci sama + isi sama | 4.590 |
| **G2** kunci sama + isi **beda** — *drift sejati* | **483** |
| **G3** kunci/kelas **beda** — ⛔ *rule kembar* | **215** (109 nama unik) |

**Dari 526 pasangan "versi sama tetapi isi beda", 93 ternyata G3** — bukan drift. Drift sejati pada
populasi itu **433**, bukan 526.

### Cara membaca `pzInsKey`

⚠️ `pzInsKey` memuat **stempel waktu per-penyimpanan**; membandingkan kunci utuh menandai hampir
semuanya berbeda. Identitas dibaca dari **basisnya**:

```
RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-PROPERTYITEM COUNTRATERETROCOV #20260921T095254.245 GMT
└──────────────────────── basis (identitas) ─────────────────────┘ └──── stempel (dibuang) ────┘
```

```powershell
$rxBasis=[regex]'\s*#\s*\d{8}T\d{6}.*$'
$basis = $rxBasis.Replace($pzInsKey,'')
# perbandingan basis dan kelas WAJIB peka huruf: -ceq / -cne
```

### Konsekuensi

1. **Angka drift mana pun yang diukur tanpa memeriksa `pzInsKey` menggabungkan dua hal berbeda** —
   berlaku bagi 95, 354, dan 526.
2. **Rule kembar tidak diekspor ulang dan tidak dipilih salah satu — keduanya diport.**
3. Komentar ketertelusuran (`CLAUDE.md` §4.6) **wajib menyebut kelas**, bukan hanya nama rule.

⚠️ **Metode deteksi drift itu sendiri tetap `[pertanyaan terbuka]` di K-058** — apakah diganti dari
banding `pyRuleSetVersion` menjadi banding hash 23 tag **belum diputuskan**, dan **tidak diputuskan di
sini**. Koreksi ini hanya menambahkan syarat: apa pun metodenya, ia wajib memeriksa identitas rule
lebih dulu.

---

Discovery NB yang baru berada satu tingkat di atas folder ini, di `D:\migrasi\RNM\OUTPUT\`.
