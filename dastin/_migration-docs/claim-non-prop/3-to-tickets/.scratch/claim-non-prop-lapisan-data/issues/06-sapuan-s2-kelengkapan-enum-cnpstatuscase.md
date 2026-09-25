---
status: selesai
---

# 06: Keadaan kasus: nilai mana yang ada, dan siapa yang membacanya

> **SELESAI 19 September 2026.** Dijalankan di atas indeks utuh 45 tag; hasilnya `SAPUAN-CELAH-DAN-JAVA.md` §5.2.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

> **DIBUKA KEMBALI 2026-09-18 — judulnya diganti, karena pertanyaannya bukan yang semula ditulis.**
> Judul lama, *“Sapuan S2: kelengkapan enum `CNPStatusCase`”*, mengandaikan ada enum yang tinggal dilengkapi. Sapu ulang dengan pola tiga-bentuk membatalkan andaian itu dari dua arah:
>
> 1. **`"CLAIM ACCEPTED"` dan `"CLAIM REJECTED"` muncul nol kali** di 279 berkas — bukan nol sebagai nilai `CNPStatusCase`, melainkan nol **di mana pun**. Dua dari empat nilai yang didaftar memori §7.3 tidak pernah ada. Butir **C2** ditutup sebagai premis gugur atas dasar ini.
> 2. **`CNPStatusCase` tidak pernah dibaca.** Tiga penulisan, nol pembacaan pada empat lapisan. **Enum yang tidak pernah dibaca bukan enum yang lengkap** — ia bukan enum.
>
> Yang tersisa karena itu bukan *melengkapi daftar* melainkan: nilai mana yang benar-benar ada, siapa yang membacanya (bila ada), dan apakah keadaan kasus disimpan di tempat lain. Dicatat sebagai **D35**; penanda yatim lainnya di **D25**.
> Ejaan tersimpan `COMITEE` satu T, sementara `COMMITTEE` dua T muncul 10 kali sebagai teks tampilan.

*Asal: `T-33` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Empat lapisan disapu; **2 nilai ditulis**, 3 berkas, **nol** rule menguji, dan `"CLAIM ACCEPTED"`/`"CLAIM REJECTED"` **nihil di 279 berkas**.

**Blocked by:**

- None (can start immediately)


**Dasar:** EVIDENCED. Menguatkan `BLUEPRINT.md` §3.1 — semula berlaku untuk lapisan Activity saja.

- [ ] Terpenuhi. Akibatnya: `STATUS_KLAIM` **tanpa `CHECK`** (T-12).

**Ketidakpastian:** Kelengkapan enum hanya dapat dibuktikan dengan membuka folder Komite — **belum diizinkan**.

## Batas jawaban — ditetapkan 19 September 2026

Folder `Komite Claim Non Prop` **tertutup untuk batch ini**. Jawaban tiket ini karena itu berbunyi:

> *"Tidak ada pembaca `CNPStatusCase` **di folder `Claim Non Prop`**"*

bukan *"tidak ada pembaca"*. **Batas itu ditulis di badan jawabannya**, bukan disimpulkan pembaca nanti.

Batas kedua, dari pemeriksaan cakupan: pernyataan itu juga berbatas pada **tag yang benar-benar dibaca indeks**. Tiket ini dikerjakan **sesudah** celah tag ditutup, supaya batas keduanya tinggal satu — folder, bukan folder *dan* tag.

## Hasil

- **Dua nilai** yang benar-benar ditulis: `"COMITEE ACCEPTANCE (DEPT. HEAD)"` dan `"INPUT ACCEPTATION CLAIM"`, di tiga tempat.
- **Nol pembacaan, nol perbandingan** — diuji terhadap 45 tag pada 279 berkas, 12 jenis rule, ditambah 480 baris Java.
- **Batas jawaban, dan ia bagian dari jawabannya**: pernyataan ini berlaku **di folder `Claim Non Prop`**. Folder Komite tertutup untuk batch ini, dan ketiga penulisnya berurusan dengan Komite — jadi pembacanya, bila ada, kemungkinan besar di sana.
- `"CLAIM ACCEPTED"` dan `"CLAIM REJECTED"` **nol kemunculan di mana pun** (D35). Dua dari empat nilai memori §7.3 tidak pernah ada.
- **Ia bukan enum yang perlu dilengkapi.** Kolom teks, diisi tiga tempat, tidak mengatur apa pun di modul ini. Diperlakukan menurut kebijakan **E17**: dimigrasi apa adanya, ditandai tak berpemilik, tanpa domain tertutup.
- Ejaan tersimpan **`COMITEE` satu T**; `COMMITTEE` dua T hanya teks tampilan, 10 kali di 7 berkas.
