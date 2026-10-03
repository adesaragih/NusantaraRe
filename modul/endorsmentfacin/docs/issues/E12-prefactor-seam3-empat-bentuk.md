# E12: Prefactor — Seam 3 menjadi satu pintu dengan empat bentuk

**What to build:** Pintu perhitungan premi yang sudah ada diperluas agar dapat memilih di antara
**empat bentuk perhitungan**, tanpa menambah seam. *"Make the change easy, then make the easy
change."*

`[terverifikasi]` Endorsement **memakai ulang** mesin premi New Business — puluhan activity
perhitungan ada di kedua korpus. Yang ditambahkan endorsement adalah **bentuk**, bukan mesin baru.

⛔ **SATU pintu, empat bentuk — bukan seam per bentuk.**

| Bentuk | Isi |
| --- | --- |
| **1 — dasar** | `TSI × Rate × faktor ÷ pembagi-komposit`; pembagi diturunkan **resolver K-018** dari faktor yang ikut |
| **2 — dua-bagian endorsement** | porsi baru untuk sisa periode **+** porsi lama untuk periode berjalan |
| **3 — tabel tarif** | premi dari lookup **+** tambahan per minggu; masukan dari **repository** |
| **4 — dekomposisi delta coverage** | `(Δrate × TSI_lama) + (ΔTSI × rate)`, dengan cabang periode pendek |

⛔ **Dua keluarga rasio prorata tetap terpisah** — yang untuk premi berbeda penyebut, satuan, dan
skala dari yang untuk delta spreading. Namanya nyaris sama; menyatukannya menggeser angka.

⛔ **Porsi periode dipakai sebagai pecahan, langsung** — tidak dibagi seratus lagi.

**Asal (Pega).** `Activity\CalculatePremiFire` · `HitungPremiOnChange` · `calculatePremiPA` ·
`CalculatePremiumTravel` · `SetLocalNonMbuProrate` · `CalculatePremiPA_FacIn`

**Keputusan.** **K-051** · K-018 · K-010/K-012 · K-027 · ADR-0004/0005

**Blocked by:** `..\01-tracer-money-ratio-premi-pa.md` · `..\03-resolver-cob-skala.md` ·
`..\04-rumus-premi-pa.md` · `..\05-rumus-premi-mbu.md` · `..\06-rumus-premi-layering.md`
⛔ Kelimanya **belum dikerjakan** — lihat §0 indeks.

**Status:** sebagian — 01-10-2026. Work owner memutus: `CalculatePremiFire` diport DI MODUL INI
(perintah 1). `backend/services/premifire.go` (`HitungPremiFire`): langkah 2/3/6/7/8/10 +
`SetLocalNonMbuProrate` langkah 1 & 10 (FIRE) — bentuk 1 dan 2 untuk FIRE EDM; HALF_UP (A37); tiga
asimetri K-046 diport apa adanya. ⛔ Menyimpang dari tiket, mengikuti korpus: (1) rumus FIRE EDM
berbeda dari mesin NB (`CountPremi_ACT`), jadi bentuk 1 BUKAN satu implementasi NB/RNW/EDM; (2) porsi
FIRE dipakai sebagai PERSEN (`@Math.divide(hari,365,6)*100`, pembagi 1e9), bukan pecahan. Belum:
satu pintu empat bentuk, resolver K-018, bentuk 3/4, dua keluarga rasio bertipe terpisah. Polis
berjalan hampir selalu → `ErrSatuanSelisihWaktuBelumTerverifikasi` (+12 jam → n,5 hari).

- [ ] Satu pintu memilih di antara **empat bentuk** berdasarkan lini bisnis dan konteks endorsement
- [ ] Bentuk 1 tetap **satu implementasi**, dipakai New Business, Renewal, dan Endorsement
- [ ] Pembagi komposit diturunkan resolver (K-018) — **bukan** angka hafalan per lini
- [ ] Bentuk 2 adalah **komposisi** dua panggilan bentuk 1, bukan cabang di dalam rumus
- [ ] Bentuk 3 menerima hasil lookup sebagai **masukan**, tidak menghitungnya sendiri
- [ ] **Dua keluarga rasio prorata terpisah** dan tidak dapat tertukar
- [ ] Porsi periode **pecahan, dipakai langsung**
- [ ] Pembulatan literal per langkah (K-011); nilai berkoma desimal (K-027)
- [ ] `Money + Ratio` **gagal saat kompilasi**; jembatan tunggal `Money × Ratio → Money`

**Catatan sebelum implementasi** — bukan blocker spec:
rumus PA & travel penuh (tarif dari repository, sumbernya nihil di korpus) ·
perhitungan PA bersama sudah terbukti **identik** New Business↔Endorsement ·
**keputusan desain:** empat bentuk sebagai satu fungsi bercabang **atau** empat implementasi di balik
satu pintu · **lima dari enam** salinan mesin premi berbeda isinya dan perlu dibedah sebelum dipakai
ulang.
