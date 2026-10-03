# 07: Nomor akseptasi — satu rangkaian, dikunci kelas kasus induk

**Status:** ready-for-agent
**Blocked by:** 06
**Menutup:** AC 52 · 53 · 54 · 55 *(4 AC)* — US 28

## Hasil & nilai pengguna

Hari ini Nomor akseptasi **belum terbit**, dan di korpus ada **tujuh pembangkit per lini usaha** yang seluruhnya **ber-remark** — sehingga tidak jelas mana yang berlaku.

Sesudah tiket ini, ⭐ Nomor akseptasi terbit dari **satu rangkaian**, dikunci pada **jenis kasus induk** — ⛔ bukan per lini usaha — dan **tidak dibangkitkan ulang** untuk akseptasi yang sudah bernomor.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Pembangkit yang hidup | satu prosedur tersimpan, dikunci **kelas kasus induk** |
| ⛔ Tujuh pembangkit per lini | ⛔ **ber-remark seluruhnya**; sisiran menyeluruh **tidak menemukan penggantinya** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K1** — ⭐ Nomor terbit sesudah ringkasan tersimpan, di **jenjang terakhir**

## Yang harus diuji

- [ ] ⭐ Nomor berasal dari **satu rangkaian**, dikunci **jenis kasus induk**
- [ ] ⛔ Tujuh pembangkit per lini usaha **tidak dialihkan**
- [ ] ⛔ Nomor yang sudah terbit **tidak dibangkitkan ulang**
- [ ] ⚠️ Perilaku prosedur tersimpannya **belum terbaca dari korpus** — ⭐ sistem baru **memerikan sendiri** aturan pembangkitannya, ⛔ tidak mewarisi yang tak terbaca

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **31** | Isi prosedur tersimpan pembangkit nomor — `[data DBA]` | ⚠️ menahan **peniruan persis**; ⭐ tidak menahan pembuatan aturan sendiri |

## Seam & verifikasi

**Seam:** lapisan layanan komite — penerbitan nomor.
1. Selesaikan dua kasus komite ⇒ ⭐ nomornya **berurut dalam satu rangkaian**.
2. Selesaikan kasus dua lini usaha berbeda ⇒ ⭐ keduanya dari **rangkaian yang sama**.
3. Coba terbitkan ulang untuk akseptasi yang sudah bernomor ⇒ ⛔ **ditolak**.
