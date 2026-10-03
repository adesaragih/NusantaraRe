# 14: Special Acceptance dan tangga akseptasi putaran kedua

**What to build:** Dua cabang inti tangga akseptasi yang sempat **ditangguhkan** kini berjalan penuh —
jalur **Special Acceptance** dan **tangga putaran kedua**.

Keduanya pernah ditangguhkan karena rule yang dipanggilnya tidak ada di ekspor yang kami terima.
Rule-nya kini tiba, rujukannya sudah diverifikasi ulang ke korpus, dan aturannya berlaku: **rule ada →
cabang dipulihkan**.

⚠️ **Dua rule Special Acceptance BUKAN Activity.** Keduanya dieksekusi sebagai **`RequestType` pada
langkah RDB-List**, berkelas integrasi, berdampingan dengan penanda akses dan halaman browse. Artinya
yang perlu ditulis adalah **query dan pemanggilan bacanya** — bukan sebuah fungsi layanan. Menyebutnya
"activity" akan menyesatkan porting.

Cabang putaran kedua berbeda: ia **benar-benar** panggilan activity, dipanggil dari post-activity
validasi tanggal.

⚠️ **Jalur ini hanya membaca.** Tidak ada bagian tiket ini yang menulis ke basis data; jalur tulis
produksi masih menunggu dan berada di luar lingkup.

**Blocked by:** 11, 13

**Status:** wontfix — kedua cabang di-remark (butir 43); penutupan dikonfirmasi work owner 01-10-2026 (butir 50)

- [ ] Jalur Special Acceptance berjalan lewat pemanggilan **baca** bergaya `RequestType`, bukan dipaksa menjadi panggilan fungsi layanan
- [ ] Kedua rule Special Acceptance diperlakukan sebagai **rule integrasi/SQL**, dengan komentar menyebut asal dan kelasnya
- [ ] Tangga putaran kedua berjalan sebagai panggilan activity dari post-activity validasi tanggal
- [ ] Keduanya tetap mematuhi aturan **satu keputusan = satu transisi** (tiket 11)
- [ ] Tidak ada operasi tulis di jalur ini
- [ ] Fixture menutupi keduanya; kasus uji membuktikan cabangnya **benar-benar tercapai**, bukan sekadar ada
- [ ] Status "ditangguhkan" dicabut di dokumentasi rancangan, dengan menyebut bukti rujukan yang memulihkannya

## Comments

### 2026-10-01 — ditahan (agent)

Bergantung pada NB-11 dan NB-13. Langkah 15–16 dan blok 21 `_ActFlow` (cabang Special Acceptance,
bergerbang `IsB2B == "ASM"` dan `IsSpecialAcceptance == "true"`) memakai kode transisi dan ejaan yang
sama dengan blok tangga (`../PERTANYAAN-AKSEPTASI.md` B1–B2).

### 2026-10-01 — kedua cabang di-remark: diusulkan ditutup (penerapan butir 43)

`[terverifikasi]` Setiap rujukan ke rule cabang ini berada di langkah berlabel `//` (= di-remark, tidak dipakai lagi,
butir 43). Dibaca per langkah dengan pengurai XML atas `pySteps` (label `pyStepsBlockName`), dicocokkan dengan
`grep -n` atas label dan pemanggilan:

| Rujukan | Pemanggil (`D:\migrasi\RNM\NB FacIn\Activity\`) | Langkah | Label |
| --- | --- | --- | --- |
| `GetLimitAkseptasi1SA_Act` (Rule-Connect-SQL, `ASM-FW-GISFW-INT-POLICYJSON!ASM!GETLIMITAKSEPTASI1SA_ACT`) | `GetLimitAkseptasi_ActFlow`, `GetLimitAkseptasi_Act` | 15 "Limit SA FOR SUW" | `//` |
| `GetLimitAkseptasi2SA_Act` (Rule-Connect-SQL) | idem | 16 "Limit SA FOR DEPHEAD UW" | `//` |
| `GetLimitAkseptasi_Act2` (`ASM-FW-GISFW-WORK!GETLIMITAKSEPTASI_ACT2`, hanya di `DDL\`) | `SetValidateDate_PostAct` | 5 "Untuk cek limit putaran 2" | `//` |
| idem | `InputOfferFacInEngineerUW_preACT` | 6 "Putaran 2" | `//` |

Blok 21 `_ActFlow` ("IsFire untuk spesial acceptance") juga `//`. Aturan "rule ada → cabang dipulihkan" (K-006)
tidak berlaku: rule-nya ada, tetapi pemanggilnya di-remark. Tidak ada kode.

