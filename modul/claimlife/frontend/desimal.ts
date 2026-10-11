// Desimal teks untuk isian Claim RNM (11-10-2026, keputusan work owner "hanya bisa angka, dan nilainya tidak boleh melebihi
// share nusantara re. diberikan pesan").
//
// ⛔ Nol float (ADR-U-0003): angka uang tidak pernah melewati Number. Perbandingan bekerja pada DIGIT teks yang sudah
// dinormalkan — inti tidak punya pembanding desimal di frontend (lib/format hanya memformat), jadi ini pembanding kecil
// berbasis digit, setara cara formatNumber membulatkan pada digit.

/**
 * Menyaring ketikan: hanya digit dan SATU pemisah desimal (titik atau koma). Karakter lain dibuang, pemisah kedua dibuang.
 * Tanpa pemisah ribuan saat diketik.
 */
export function saringKetikanAngka(teks: string): string {
  let pemisah = false
  let out = ''
  for (const c of teks) {
    if (c >= '0' && c <= '9') out += c
    else if ((c === '.' || c === ',') && !pemisah) {
      pemisah = true
      out += c
    }
  }
  return out
}

/** Bentuk baku "123.45" dari isian (koma = titik), atau null bila bukan angka tak bertanda yang sah. */
export function normalisasiDesimal(teks: string): string | null {
  const t = teks.trim().replace(',', '.')
  if (!/^\d+(\.\d*)?$|^\.\d+$/.test(t)) return null
  const [bulatMentah = '', pecahanMentah = ''] = t.split('.')
  const bulat = bulatMentah.replace(/^0+(?=\d)/, '') || '0'
  const pecahan = pecahanMentah.replace(/0+$/, '')
  return pecahan === '' ? bulat : `${bulat}.${pecahan}`
}

/** Membandingkan dua desimal baku tak bertanda: -1 / 0 / 1. Lewat panjang bagian bulat lalu digit, tanpa float. */
export function bandingDesimal(a: string, b: string): -1 | 0 | 1 {
  const [ab = '0', ap = ''] = a.split('.')
  const [bb = '0', bp = ''] = b.split('.')
  const ka = ab.replace(/^0+(?=\d)/, '')
  const kb = bb.replace(/^0+(?=\d)/, '')
  if (ka.length !== kb.length) return ka.length < kb.length ? -1 : 1
  if (ka !== kb) return ka < kb ? -1 : 1
  const n = Math.max(ap.length, bp.length)
  const pa = ap.padEnd(n, '0')
  const pb = bp.padEnd(n, '0')
  if (pa === pb) return 0
  return pa < pb ? -1 : 1
}

/** Desimal baku tak bertanda ("123.45") menjadi [bilangan bulat BigInt, skala]. */
function keSkala(baku: string): [bigint, number] {
  const [bulat = '0', pecahan = ''] = baku.split('.')
  return [BigInt(bulat + pecahan), pecahan.length]
}

/** Desimal desimal-tetap (BigInt berskala) menjadi teks baku tanpa nol ekor. */
function dariSkala(n: bigint, skala: number): string {
  const t = n.toString().padStart(skala + 1, '0')
  const bulat = t.slice(0, t.length - skala) || '0'
  const pecahan = skala > 0 ? t.slice(t.length - skala).replace(/0+$/, '') : ''
  return pecahan === '' ? bulat : `${bulat}.${pecahan}`
}

/**
 * Percent Claim (%) di layar — SAMA dengan backend `models.HitungPersenKlaim`: pecahan Claim RNM / Share dibulatkan
 * SETENGAH-KE-ATAS pada 5 desimal (`@divide(…,100,5)` b702; apd.BaseContext backend = half-up, diperiksa 11-10-2026), lalu
 * × 100. Masukan desimal baku tak bertanda (normalisasiDesimal). Share kosong / nol = "". Tanpa float (BigInt berskala).
 */
export function persenKlaimTeks(klaimRnm: string, share: string): string {
  const [a, sa] = keSkala(klaimRnm)
  const [b, sb] = keSkala(share)
  if (b === 0n) return ''
  // a/10^sa ÷ b/10^sb × 10^5 = a·10^(sb+5) / (b·10^sa)
  const pembilang = a * 10n ** BigInt(sb + 5)
  const penyebut = b * 10n ** BigInt(sa)
  let q = pembilang / penyebut
  if ((pembilang % penyebut) * 2n >= penyebut) q += 1n
  // q = pecahan × 10^5; persen = q × 100 / 10^5 = q / 10^3.
  return dariSkala(q, 3)
}
