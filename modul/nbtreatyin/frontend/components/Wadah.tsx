// Wadah layout Section Pega TANPA judul: `pyContainerType=NOHEADER` /
// `pyIncludeHeader=false` (judul wadah tidak dirender walau `pyTitle` terisi).
// Bentuknya panel inti tanpa `panel__title` - judul buatan tidak ditambahkan
// (W6 audit silang P3; bab 0 butir 7 prompt putaran 3).

import type { ReactNode } from 'react'

export default function Wadah({ children }: { children: ReactNode }) {
  return <section className="panel">{children}</section>
}
