// Penyusun rupa layar kasus Claim Prop - pola Kelola User (keputusan work owner 08-10-2026): kartu `panel` bertitel,
// isinya `form-grid`, setiap medan berlabel di atas kotak isian. Isi pohon tata server TIDAK diubah; di sini hanya
// dikelompokkan:
//   - bagian tanpa label diratakan ke induknya (pengelompokan tata letak Pega, bukan blok bertitel);
//   - label tepat sebelum medan tanpa label menjadi label medan itu (`Q`, `/`, `U/Y`);
//   - label satuan ("%") tepat sesudah medan menjadi satuan medan itu;
//   - tombol sesudah medan menempel di kanan kotak isiannya (Choose Cause of Loss, ikon), kecuali tombol berlabel
//     di akhir isi - itu baris aksi;
//   - tingkat layar: label judul membuka kartu, bagian berlabel menjadi kartu sendiri (bagian kecil tanpa grid di
//     tengah kartu yang terbuka tetap di kartu itu), tombol tanpa kartu = baris aksi.

import type { Tata } from '../api'

/** Unsur tata sesudah label pendamping digabung. */
export interface Butir extends Tata {
  satuan?: string
}

export type Unsur =
  | { jenis: 'medan'; t: Butir; tombol: Tata[] }
  | { jenis: 'tombol'; tombol: Tata[]; akhir: boolean }
  | { jenis: 'judul'; label: string }
  | { jenis: 'sub'; t: Tata }
  | { jenis: 'grid'; t: Tata }

export type Blok = { jenis: 'panel'; judul: string; isi: Butir[] } | { jenis: 'aksi'; tombol: Tata[] }

const SATUAN = new Set(['%'])

/** Ratakan bagian tanpa label, lalu gabungkan label pendamping ke medannya. */
export function ratakan(tata: readonly Tata[]): Butir[] {
  const datar: Butir[] = []
  const jalan = (ts: readonly Tata[]) => {
    for (const t of ts) {
      if (t.jenis === 'bagian' && !t.label) jalan(t.anak ?? [])
      else datar.push(t)
    }
  }
  jalan(tata)
  const out: Butir[] = []
  for (let i = 0; i < datar.length; i++) {
    const t = datar[i]!
    const berikut = datar[i + 1]
    const sebelum = out[out.length - 1]
    if (t.jenis === 'label' && berikut?.jenis === 'medan' && !berikut.label) {
      out.push({ ...berikut, label: t.label })
      i++
    } else if (t.jenis === 'label' && SATUAN.has(t.label ?? '') && sebelum?.jenis === 'medan') {
      out[out.length - 1] = { ...sebelum, satuan: t.label }
    } else {
      out.push(t)
    }
  }
  return out
}

/** Isi satu kartu / sub-bagian / modal menjadi sel form-grid. */
export function susunIsi(tata: readonly Tata[]): Unsur[] {
  const b = ratakan(tata)
  const out: Unsur[] = []
  b.forEach((t, i) => {
    const akhir = out[out.length - 1]
    switch (t.jenis) {
      case 'medan':
        out.push({ jenis: 'medan', t, tombol: [] })
        return
      case 'tombol': {
        const sisaTombol = b.slice(i).every((x) => x.jenis === 'tombol')
        if (akhir?.jenis === 'medan' && (!t.label || !sisaTombol)) akhir.tombol.push(t)
        else if (akhir?.jenis === 'tombol') akhir.tombol.push(t)
        else out.push({ jenis: 'tombol', tombol: [t], akhir: false })
        return
      }
      case 'label':
        out.push({ jenis: 'judul', label: t.label ?? '' })
        return
      case 'grid':
        out.push({ jenis: 'grid', t })
        return
      default:
        out.push({ jenis: 'sub', t })
    }
  })
  const terakhir = out[out.length - 1]
  if (terakhir?.jenis === 'tombol') terakhir.akhir = true
  return out
}

function adaGrid(t: Tata): boolean {
  return (t.anak ?? []).some((a) => a.jenis === 'grid' || (a.jenis === 'bagian' && (!!a.label || adaGrid(a))))
}

/** Tingkat layar: kartu bertitel dan baris aksi. */
export function susunLayar(tata: readonly Tata[]): Blok[] {
  const b = ratakan(tata)
  const out: Blok[] = []
  let kini: { judul: string; isi: Butir[] } | null = null
  const tutup = () => {
    if (kini === null) return
    const { judul, isi } = kini
    if (judul === '' && isi.length > 0 && isi.every((x) => x.jenis === 'tombol'))
      out.push({ jenis: 'aksi', tombol: isi })
    else if (judul !== '' || isi.length > 0) out.push({ jenis: 'panel', judul, isi })
    kini = null
  }
  b.forEach((t, i) => {
    if (t.jenis === 'bagian') {
      const lepas = b[i + 1]?.jenis === 'medan' || b[i + 1]?.jenis === 'tombol'
      if (kini !== null && kini.judul !== '' && lepas && !adaGrid(t)) {
        kini.isi.push(t)
        return
      }
      tutup()
      out.push({ jenis: 'panel', judul: t.label ?? '', isi: t.anak ?? [] })
      return
    }
    if (t.jenis === 'label') {
      if (kini === null) kini = { judul: t.label ?? '', isi: [] }
      else if (kini.judul === '' && kini.isi.length === 0) kini.judul = t.label ?? ''
      else if (kini.isi.length === 0) kini.isi.push(t)
      else {
        tutup()
        kini = { judul: t.label ?? '', isi: [] }
      }
      return
    }
    kini ??= { judul: '', isi: [] }
    kini.isi.push(t)
  })
  tutup()
  return out
}
