// Penjumlahan DESIMAL EKSAK untuk rumus yang dijalankan di layar.
//
// ⛔ Bukan `Number`: `0.1 + 0.2` dan nilai uang sebelas digit berpecahan
// tidak boleh melewati float. Nilai Pega teks desimal; dijumlahkan sebagai
// bilangan bulat ber-skala (BigInt), seperti BigDecimal Pega.

/** Teks desimal → (mantissa, skala); kosong/tak terbaca = 0, seperti properti Decimal Pega yang kosong. */
function urai(v: string): { m: bigint; s: number } {
  const t = v.trim()
  const c = /^([+-]?)(\d*)(?:\.(\d*))?$/.exec(t)
  if (c === null || (c[2] === '' && (c[3] ?? '') === '')) return { m: 0n, s: 0 }
  const tanda = c[1] === '-' ? -1n : 1n
  const pecahan = c[3] ?? ''
  return { m: tanda * BigInt(`${c[2] === '' ? '0' : c[2]}${pecahan}`), s: pecahan.length }
}

function tulis(m: bigint, s: number): string {
  const negatif = m < 0n
  const mutlak = (negatif ? -m : m).toString().padStart(s + 1, '0')
  const bulat = s === 0 ? mutlak : mutlak.slice(0, mutlak.length - s)
  const pecahan = s === 0 ? '' : `.${mutlak.slice(mutlak.length - s)}`
  return `${negatif ? '-' : ''}${bulat}${pecahan}`
}

/** a + b, skala hasil = skala terbesar (perilaku BigDecimal.add). */
export function tambah(a: string, b: string): string {
  const x = urai(a)
  const y = urai(b)
  const s = Math.max(x.s, y.s)
  const m = x.m * 10n ** BigInt(s - x.s) + y.m * 10n ** BigInt(s - y.s)
  return tulis(m, s)
}
