// Pembersih HTML pesan korespondensi (tiket 56) - `.Message` Pega adalah pxRichTextEditor; sistem baru menyimpan dan
// mengirim HTML-nya, jadi isinya WAJIB dibersihkan sebelum ditaruh ke `innerHTML` / dikirim. Aturan sama dengan backend
// (kontrak tiket 56): hanya tag di `TAG_HTML_KORESPONDENSI`, TANPA atribut apa pun (gaya, tautan, gambar dibuang); isi
// `script` / `style` dibuang utuh; tag lain dibuang tetapi teksnya dipertahankan; teks di-escape (`&` kecuali entitas
// yang sah, `<`, `>`). Hasil hanya tersusun dari tag polos yang diizinkan + teks ter-escape.

export const TAG_HTML_KORESPONDENSI = ['b', 'strong', 'i', 'em', 'u', 'p', 'br', 'div', 'ul', 'ol', 'li'] as const

const DIIZINKAN = new Set<string>(TAG_HTML_KORESPONDENSI)
const ISI_DIBUANG = new Set(['script', 'style'])
const ENTITAS = /^&(?:[a-zA-Z][a-zA-Z0-9]{1,31}|#\d{1,7}|#[xX][0-9a-fA-F]{1,6});/

function escapeTeks(teks: string): string {
  let hasil = ''
  for (let i = 0; i < teks.length; i++) {
    const c = teks[i]
    if (c === '&') {
      const m = ENTITAS.exec(teks.slice(i))
      if (m) {
        hasil += m[0]
        i += m[0].length - 1
      } else hasil += '&amp;'
    } else if (c === '<') hasil += '&lt;'
    else if (c === '>') hasil += '&gt;'
    else hasil += c
  }
  return hasil
}

export function bersihkanHtml(html: string): string {
  let hasil = ''
  let buang: string | null = null
  const pola = /<(\/?)([a-zA-Z][a-zA-Z0-9]*)\b[^>]*>|<!--[\s\S]*?-->/g
  let akhir = 0
  for (let m = pola.exec(html); m !== null; m = pola.exec(html)) {
    const teks = html.slice(akhir, m.index)
    akhir = pola.lastIndex
    if (buang === null) hasil += escapeTeks(teks)
    if (m[2] === undefined) continue // komentar
    const tutup = m[1] === '/'
    const nama = m[2].toLowerCase()
    if (buang !== null) {
      if (tutup && nama === buang) buang = null
      continue
    }
    if (ISI_DIBUANG.has(nama)) {
      if (!tutup) buang = nama
      continue
    }
    if (!DIIZINKAN.has(nama)) continue
    if (nama === 'br') hasil += tutup ? '' : '<br>'
    else hasil += tutup ? `</${nama}>` : `<${nama}>`
  }
  if (buang === null) hasil += escapeTeks(html.slice(akhir))
  return hasil
}

/** Batas MESSAGE di backend (`T_CORRESPONDENCELIST.MESSAGE` VARCHAR2(4000); penjaga inti melarang CLOB di migrasi modul). */
export const BATAS_BITA_PESAN = 4000

/** Ukuran HTML bersih dalam bita UTF-8 - yang dihitung backend. */
export function bitaPesan(html: string): number {
  return new TextEncoder().encode(bersihkanHtml(html)).length
}

/** Pesan kosong secara isi (hanya tag / spasi / &nbsp;). */
export function pesanKosong(html: string): boolean {
  return bersihkanHtml(html).replace(/<[^>]*>/g, '').replace(/&nbsp;/g, ' ').trim() === ''
}
