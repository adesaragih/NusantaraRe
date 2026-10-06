// Ketikan angka di isian layar NB Treaty In (permintaan work owner 06-10-2026): hanya angka yang dapat diketik,
// ribuan dipasang otomatis selama mengetik, titik ATAU koma yang diketik = pemisah desimal, teks berformat yang
// ditempel (`123.456,678`, `123,456.678`) dikenali. Tampilan mengikuti pola inti `sajian.ts` (titik ribuan, koma
// desimal); nilai yang dikirim ke backend selalu MENTAH (`123456.678`) - nilai tidak pernah dibulatkan di sini.

/** Pecahan desimal terpanjang yang diterima - skala kolom NUMBER(38,10). */
const DESIMAL_MAKS = 10

const MENTAH = /^-?(\d+)(\.(\d*))?$/

function kelompokRibuan(bulat: string): string {
  return bulat.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
}

/** Tampilan isian sebuah nilai mentah (`1234.5` -> `1.234,5`); bukan angka = apa adanya. */
export function tampilKetik(mentah: string): string {
  const t = mentah.trim()
  const m = MENTAH.exec(t)
  if (m === null) return t
  const negatif = t.startsWith('-')
  return (negatif ? '-' : '') + kelompokRibuan(m[1]!) + (m[2] !== undefined ? ',' + (m[3] ?? '') : '')
}

interface Angka {
  negatif: boolean
  bulat: string
  /** `null` = tanpa pemisah desimal; `''` = pemisah baru diketik. */
  pecahan: string | null
}

function rakit(a: Angka): { mentah: string; tampil: string } {
  let bulat = a.bulat.replace(/^0+(?=\d)/, '')
  if (bulat === '' && a.pecahan !== null) bulat = '0'
  const pecahan = a.pecahan === null ? null : a.pecahan.slice(0, DESIMAL_MAKS)
  if (bulat === '' && pecahan === null) return { mentah: a.negatif ? '-' : '', tampil: a.negatif ? '-' : '' }
  const tanda = a.negatif ? '-' : ''
  const mentah = tanda + bulat + (pecahan ? '.' + pecahan : '')
  const tampil = tanda + kelompokRibuan(bulat) + (pecahan !== null ? ',' + pecahan : '')
  return { mentah, tampil }
}

/** Teks BERFORMAT apa pun (ditempel) ke angka: kedua pemisah ada = yang terakhir desimal; satu jenis pemisah
 *  muncul sekali = desimal; muncul berkali-kali = ribuan. Selain angka, pemisah, dan minus di depan dibuang. */
function bacaBerformat(teks: string): Angka {
  const t = teks.trim()
  const negatif = /^-/.test(t.replace(/^[^\d.,-]+/, ''))
  const s = t.replace(/[^\d.,]/g, '')
  const koma = s.split(',').length - 1
  const titik = s.split('.').length - 1
  let desimal = -1
  if (koma > 0 && titik > 0) desimal = Math.max(s.lastIndexOf(','), s.lastIndexOf('.'))
  else if (koma === 1) desimal = s.indexOf(',')
  else if (titik === 1) desimal = s.indexOf('.')
  if (desimal < 0) return { negatif, bulat: s.replace(/[.,]/g, ''), pecahan: null }
  return { negatif, bulat: s.slice(0, desimal).replace(/[.,]/g, ''), pecahan: s.slice(desimal + 1).replace(/[.,]/g, '') }
}

/** Teks isian layar (titik = ribuan buatan sendiri, koma = desimal) ke angka. */
function bacaIsian(teks: string): Angka {
  const negatif = teks.trimStart().startsWith('-')
  const s = teks.replace(/[^\d,]/g, '')
  const i = s.indexOf(',')
  if (i < 0) return { negatif, bulat: s, pecahan: null }
  return { negatif, bulat: s.slice(0, i), pecahan: s.slice(i + 1).replace(/,/g, '') }
}

/** Nilai mentah dari teks isian layar (`1.234,5` -> `1234.5`). */
export function angkaMentah(tampil: string): string {
  return rakit(bacaIsian(tampil)).mentah
}

/** Jumlah "karakter bermakna" (digit, koma, minus) sebelum posisi `n` - titik ribuan tidak dihitung. */
function bermakna(teks: string, n: number): number {
  return teks.slice(0, n).replace(/[^\d,-]/g, '').length
}

/** Posisi di `teks` sesudah `k` karakter bermakna. */
function posisiBermakna(teks: string, k: number): number {
  if (k <= 0) return 0
  let hitung = 0
  for (let i = 0; i < teks.length; i++) {
    if (/[\d,-]/.test(teks[i]!)) hitung++
    if (hitung === k) return i + 1
  }
  return teks.length
}

/**
 * Satu perubahan isian: `lama` = tampilan sebelumnya, `baru` = isi kotak sesudah pengguna mengetik / menghapus /
 * menempel, `kursor` = posisi kursor di `baru`. Jawab tampilan rapi, nilai mentah, dan posisi kursor baru.
 */
export function ubahKetikan(lama: string, baru: string, kursor: number): { tampil: string; mentah: string; kursor: number } {
  // bagian yang disisipkan: awalan dan akhiran bersama dibuang
  let a = 0
  while (a < lama.length && a < baru.length && lama[a] === baru[a]) a++
  let b = 0
  while (b < lama.length - a && b < baru.length - a && lama[lama.length - 1 - b] === baru[baru.length - 1 - b]) b++
  const sisip = baru.slice(a, baru.length - b)

  let angka: Angka
  let teks = baru
  if (sisip.length > 1 || (lama === '' && /[.,]/.test(baru) && baru.length > 1)) {
    // menempel (atau mengganti seluruh isi) teks berformat
    angka = bacaBerformat(baru)
    const r = rakit(angka)
    return { ...r, kursor: r.tampil.length }
  }
  if (sisip === '.' || sisip === ',') {
    // titik atau koma yang DIKETIK = pemisah desimal; pemisah kedua diabaikan
    const sebelum = baru.slice(0, a) + baru.slice(a + 1)
    teks = /,/.test(sebelum) ? sebelum : baru.slice(0, a) + ',' + baru.slice(a + 1)
    if (/,/.test(sebelum)) kursor = a
  }
  if (sisip === '-' && a !== 0) {
    teks = baru.slice(0, a) + baru.slice(a + 1)
    kursor = a
  }
  angka = bacaIsian(teks)
  const r = rakit(angka)
  const k = bermakna(teks, kursor)
  // "0," otomatis di depan pecahan: kursor ikut bergeser bila nol disisipkan
  const tambahNol = angka.bulat === '' && angka.pecahan !== null ? 1 : 0
  return { ...r, kursor: posisiBermakna(r.tampil, k + tambahNol) }
}
