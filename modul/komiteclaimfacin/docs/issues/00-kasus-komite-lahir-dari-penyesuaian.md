# 00: Kasus komite lahir dari penyesuaian — jenjang dibekukan

**Status:** ready-for-agent
**Blocked by:** `claim-facin\issues\13` *(kontrak muatan)*
**Menutup:** AC 1 · 2 · 3 · 4 · 5 · 6 · 7 · 8 · 96 · 97 · 98 · 99 · 100 *(13 AC)* — US 1–5

## Hasil & nilai pengguna

Hari ini Penyesuaian di atas kewenangan penilai **belum dapat berubah menjadi kasus komite**, dan tidak ada yang menetapkan **berapa jenjang** harus menyetujui.

Sesudah tiket ini, Penilai menyerahkan penyesuaian, dan **kasus komite lahir** dengan **jumlah jenjang yang ditetapkan susunan jenjang komite**. ⭐ Tiap jenjang **dibekukan bersama jabatan dan akun pemegangnya** — pergantian pemegang jabatan sesudah itu **tidak memindahkan kasus**.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Pembentukan kasus | rule pembuat nomor komite menyemai **jumlah jenjang** dan **giliran mulai**, lalu membuat kasus anak |
| Jalur satu jenjang | rule kirim tutup klaim dan kirim tolak klaim — ⭐ keduanya menyemai **satu jenjang** |
| Penguncian induk | modul komite **mengunci kasus induk** sebelum mengerjakan apa pun |
| Empat jalur masuk | penyesuaian · tolak · tutup klaim · survey — ⛔ **jalur survey tidak dialihkan**, rule-nya tidak ada di korpus mana pun |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K13** — ⭐ Jenjang dibekukan bersama **akun pemegangnya**, bukan hanya jabatannya
- **K6** — ⭐ Komite lini FAC dan lini PROP memakai **tabel yang sama** — ⛔ nol tabel baru

## Yang harus diuji

- [ ] Kasus komite lahir dengan **jumlah jenjang benar** dari susunan jenjang
- [ ] Giliran mulai pada **jenjang pertama**
- [ ] ⭐ Jenjang **dibekukan**: pergantian pemegang jabatan sesudah pembentukan **tidak memindahkan kasus**
- [ ] ⛔ Jabatan **tanpa pemegang aktif** ⇒ pembentukan **ditolak dengan galat yang menyebut jabatannya**
- [ ] ⭐ Jalur **tutup klaim** dan **tolak klaim** membentuk komite **satu jenjang**
- [ ] ⭐ Seluruh data kutipan yang dibutuhkan penggolongan **tersalin dari kasus induk**
- [ ] ⛔ Jalur **survey** **tidak dibangun**

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **39** | ⚠️ **Satu jabatan dipegang LEBIH DARI SATU orang — siapa menerima giliran?** ⭐ Usul asisten: tolak dengan galat yang jelas, ⛔ jangan memilih diam-diam | ⚠️ **MENAHAN pembentukan kasus** bila keadaan itu terjadi |
| **6** | Isi daftar jabatan dan susunan jenjang | ⛔⛔ **MENAHAN** penentuan jumlah jenjang |

## Seam & verifikasi

**Seam:** lapisan layanan komite — ⭐ **satu pintu masuk**, bukan layar.
2. Serahkan satu penyesuaian — ⭐ kasus komite lahir, jumlah jenjang benar.
3. Ganti pemegang jabatan **sesudah** pembentukan — ⭐ kasus **tidak berpindah**.
4. Kosongkan pemegang salah satu jabatan — ⛔ pembentukan **ditolak**, galatnya **menyebut jabatannya**.
5. Serahkan lewat jalur tutup klaim — ⭐ komitenya **satu jenjang**.
