# PROMPT — lanjutan 3 modul Claim Life: sesudah tiket 04 (`65e937e`) — **rantai tiket tanpa mengakhiri giliran**

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief modul **`PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE.md`** §1–§11, brief lanjutan 1 §3–§4, dan brief
> lanjutan 2 §1–§7 berlaku **seluruhnya**. Berkas ini menambah **satu mekanisme** *(§1)* yang sudah
> dua kali gagal ditegakkan oleh kalimat saja, mengulang **paket keputusan** *(§2)*, dan memberi fakta
> yang sudah diverifikasi untuk tiket 05 *(§4)*.
>
> **GILIRAN INI:** Langkah 0 → **05 → 15 → 07 → 08 → 09 → 10 → 12 → 11 → 13**. Pesan ke manusia
> **hanya satu**, di akhir *(§5)*.

---

## 0. KEADAAN AWAL — 26 September 2026 malam, sesudah `65e937e`

| | Keadaan |
| --- | --- |
| `HEAD` | `65e937e` *claim-life: tiket 04 — mesin status per baris, status klaim turunan*; sebelumnya `e1bd6e3` *(docs lanjutan 2)*, `e2190bb` *(tiket 06)*. Working tree bersih |
| Uji yang lulus | vet · vet db · gofmt nol · build · **167 PASS · 0 FAIL · 26 SKIP** · `tsc` · **5** JS · **88** modul |
| Verifikasi independen `65e937e` | 12 berkas +1.160/−41 dan seluruh angka tereproduksi. **Sensus penulis `STS_REJECT` diulang sendiri** dengan pola `.*STS_REJECT` atas seluruh `Activity\` Claim Life + Komite Claim Life: persis **enam `Property-Set` di lima rule** Claim Life *(`RejectOSClaimLife_Act` ×2, `SaveAdjustment_Act` `.AdjustmentList(<LAST>).STS_REJECT = 1`, `SaveOutStandingLife_Act` `0`, `serviceInsertArasapasClaimLife_act` header, `SetSTS_Reject` diagnosa)* + `KomitePostAdjustment` *(`1` ×3 dan `2` ×1 pada baris **dan** pada peserta)*; `ProtectCloseClaim_act` hanya **membaca**. Penulis `ACCEPTATION_DATE`: `SaveAdjustment_Act` *(`<LAST>`, `@CurrentDateTime()`)* dan `KomitePostAdjustment` ×3. Ralat executor **benar**, dan ia jujur menandai sensus pertamanya salah |
| Tiket | 04 `claimed` **6/8** *(dua AC sisa menunggu pintu HTTP tiket 05)*; 06 `claimed` 4/5; 03 `claimed` 11/28; 02 7/26; 14 40/53; 01 4/14; **05, 15, 07, 08, 09, 10, 12, 11, 13 `ready-for-agent`** |
| Penjaga | `CREATE` = 20 · `Buka()` = 8 · `kolomSalin` = 24 · penyebut `EDMSTATUS` = 2 · literal kode status hanya di `models` · `Type` satu rumah |

---

## 1. MEKANISME RANTAI — mengapa giliran selalu berakhir sesudah satu tiket, dan aturannya kini

**Sebabnya, dari teks skill:** `implement` berbunyi *"Implement the work described by the user in the
spec or tickets … Once done, use /code-review … Commit"*. Ia dipanggil **manusia satu kali** untuk
seluruh pekerjaan yang dijelaskan *(seluruh modul)*, dan **tidak dapat dipanggil ulang oleh model**
*(`disable-model-invocation: true`)* — juga **tidak perlu**: langkah-langkahnya berulang per tiket.
Giliran berakhir bukan karena skill, melainkan karena executor **menyajikan laporan review dan
laporan giliran sebagai pesan terakhirnya** *("Berikutnya: tiket NN … Saya lanjut dari sana.")*.
Pesan tanpa panggilan alat sesudahnya = giliran selesai = manusia harus menempel lagi.

**Aturan `[DIPUTUSKAN — work owner]`:**

1. **Kedua laporan `/code-review`** *(Standards, Spec)* ditulis ke bab `### Hasil /code-review` **di
   berkas tiket** — seperti yang sudah dilakukan — dan **tidak** disajikan sebagai pesan. Skill
   code-review meminta "present"; yang dimaksud di proyek ini adalah **ke tiket**.
2. **Antara dua tiket tidak ada pesan penutup.** Sesudah commit: satu baris kemajuan
   *(`tiket NN — <SHA> — AC x/y — test a PASS · b SKIP`)* **diikuti, di pesan yang sama, panggilan alat
   pertama tiket berikutnya** *(mengubah `Status:` tiket berikutnya menjadi `claimed`)*. Baris kemajuan
   yang berdiri sendiri tanpa panggilan alat **melanggar** aturan ini.
3. **Putarannya, harfiah:**

   ```
   untuk tiket dalam [05, 15, 07, 08, 09, 10, 12, 11, 13]:
       Status: claimed
       → baca XML yang belum dicatat tiket lain (bab Pembacaan ulang, ≤ 40 baris)
       → test dulu di seam → kode → verifikasi penuh
       → bab Implementasi (≤ 50 baris)
       → /code-review atas titik tetap tiket ini → perbaiki → verifikasi penuh lagi
       → commit `claim-life: tiket NN — <judul>` + baris `Tiket:`
       → satu baris kemajuan + LANGSUNG tiket berikutnya
   laporan akhir (§5) HANYA sesudah tiket 13, atau berhenti sah (brief lanjutan 2 §1-2)
   ```

4. **Konteks yang menipis bukan alasan berhenti.** Bila konteks harus diringkas, tulis dulu catatan
   *"lanjut dari sini: <tiket, langkah, SHA>"* di bab Implementasi tiket yang sedang berjalan, lalu
   teruskan. Menyerahkan giliran ke manusia untuk "melanjutkan di sesi baru" **bukan** pilihan selama
   masih ada tiket yang dapat dikerjakan.
5. **Berhenti sah** persis brief lanjutan 2 §1-2. Kalimat *"saya lanjut dari sana"* hanya boleh muncul
   di laporan akhir bila giliran berhenti sah, dan harus menyebut **gerbang** yang menghentikannya.

---

## 2. PAKET KEPUTUSAN — masih `[USULAN]` saat berkas ini ditulis

Akibatnya sudah terlihat di tiket 04: jejak audit duduk di balik `JejakBelumDiputuskan`, diagnosa tidak
dicerminkan. Tanpa paket ini tiket **09, 10, 11** berhenti di antarmuka.

**PAKET MODUL: `[USULAN]`** ← ganti menjadi `[DIPUTUSKAN]` untuk menyetujui **aj + am + af + al**
sekaligus *(rincian tiap butir: brief lanjutan 2 §2)*. Work owner juga dapat menyatakannya kepada
asisten *("setuju paket")*, yang lalu mengubah kata ini sebelum sesi mulai. Butir per butir dapat
ditimpa di sini:

| | Keputusan | Keadaan |
| ---: | --- | --- |
| aj | migrasi `011` kolom bank | mengikuti PAKET |
| am | tabel jejak audit `T_CLAIMLF_JEJAK` — tiket 09; tiket 04 lalu mengganti `JejakBelumDiputuskan` | mengikuti PAKET |
| af | tabel Komite + `COVER_KEY` dibuat di modul ini — tiket 10, 11 | mengikuti PAKET |
| al | tabel diagnosa — bukti dari XML tiket 08 dulu | mengikuti PAKET, **bersyarat bukti** |
| ag, ah, ak, o1–o3 | tetap | `[USULAN]` / `[terbuka]` |

Bila PAKET `[DIPUTUSKAN]`: tiap tabel baru = langkah migrasi bernomor baru + STRUKTUR bertanggal + satu
baris di tiket 14 + penjaga cacah `CREATE` diperbarui *(bukan dilonggarkan)*; **aj** dikerjakan sebagai
commit tambahan tiket 03 **sesudah tiket 13**.

---

## 3. LANGKAH 0

`git add PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-3.md` → commit `docs: brief lanjutan 3 — mekanisme
rantai tiket` → `git status --porcelain` kosong → uji tanpa Oracle hijau *(167 · 26 SKIP · 5 JS · 88
modul)* → SHA = titik tetap **tiket 05**.

---

## 4. FAKTA YANG SUDAH DIVERIFIKASI UNTUK TIKET BERIKUTNYA

### Tiket 05 — Reject Outstanding oleh Admin

Bekal `[terverifikasi]` *(tiket 03 bab pembacaan, tiket 04 bab pembacaan, sensus §0)*:
`RejectOSClaimLife_Act` langkah 2 menulis `2` pada baris **dan** peserta serta `IsCheck = "false"` pada
peserta; langkah 3–7 merakit `TempInputDetail.CARI1…13` untuk **`UpdateOsAkseptasiClaimLife_sql`**
*(jalur tulis datar warisan — baca `pyBrowseSQL`/`pySaveSQL`-nya utuh dan cocokkan kolomnya dengan
`BarisLamaDari`)*. Mesin `services.Status.Ubah` *(tiket 04)* sudah mencerminkan baris + peserta +
header dan menstempel `ACCEPTATION_DATE` hanya pada Aksep.

Yang dibangun: `services.Tolak(ctx, pelaku, klaimID, adjID)` → `WajibPeran(pelaku, "ReasLifeAdmin")`
**di sini, bukan menunggu tiket 07** *(AC 05 menuntut 403; tiket 07 menggeneralisasi)* → klaim
bernomor → baris Outstanding → `Status.Ubah(…, Ditolak)` → `IS_CHECK = 0` peserta *(AC baru menurut
XML)* → baris datar warisan diperbarui lewat jalur yang ada, **tanpa** menyalin nama orang. Pintu
`POST /api/klaim-life/{id}/adjustment/{adjId}/tolak` → 200 / 403 / 409 / 422; pintu ini sekaligus
menutup **dua AC sisa tiket 04**. ⚠️ Sampai **o** diputuskan tidak ada klaim bernomor: jalur HTTP nyata
menjawab 422 — nyatakan di tiket; test `services` memakai `penomorUji`. Frontend: kontrol "Reject
Outstanding" pada baris, tampil hanya bila gerbang XML terpenuhi.

### Tiket 15, 07, 08 — brief lanjutan 1 §4 dan lanjutan 2 §5 berlaku apa adanya

Tambahan untuk 08: bukti **al** *(`MedicalCheckClaimLife.xml`, kelas `Data-DiagnoseLife`)* ditulis di
tiket 08 **apa pun** keadaan PAKET, supaya keputusan work owner berdiri di atas bukti.

### Tiket 09, 10, 12, 11, 13 — brief modul §4; gerbangnya §2

Tanpa PAKET: 09 = `Jejak` tetap antarmuka + bab pembacaan XML + daftar transisi yang wajib direkam;
10 = muatan penyerahan + roster + antarmuka `PencatatKomite`; 11 = pembaca di balik antarmuka;
12 = antarmuka efek keluar + antrean + stub; 13 = skrip + rekonsiliasi di skema uji. Semua tetap
**di-commit** sebagai tiket `claimed` dengan AC terbuka bertanda pemilik — bukan dilewati.

---

## 5. LAPORAN AKHIR GILIRAN — satu-satunya pesan ke manusia

Persis brief lanjutan 2 §6: tabel per tiket *(titik tetap, commit, AC ditutup/diralat/total, uji,
terkunci oleh)*; keputusan §2 yang dipakai/ditunggu; terbuka menurut pemilik; telemetri brief modul §10
+ cacah tiket ter-commit giliran ini, XML dibaca *(berkas, byte)*, ralat menurut XML, token sub-agen
per tiket, token sesi utama *(⛔ tidak terukur — nyatakan)*.

---

## 6. PERSETUJUAN MANUSIA

Brief modul §9 dan brief lanjutan 2 §7. Tidak ada yang baru.

---

*Disusun 26 September 2026 malam sesudah verifikasi independen `65e937e`: uji dijalankan ulang, diff
12 berkas dibaca, sensus penulis `STS_REJECT` dan `ACCEPTATION_DATE` diulang sendiri atas kedua korpus,
teks skill `implement` dan `code-review` dibaca untuk menemukan sebab giliran berakhir.*
