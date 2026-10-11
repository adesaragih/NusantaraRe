// Klik baris = buka / tutup detail (work owner 10-10-2026: "untuk semua nb facin, kalau ada tabel gini, tabelnya bisa juga
// dibuka saat barisnya di klik, bukan harus di tanda panah itu. di klik di judulnya atau namanya selama msh baris dia, harus bisa
// terbuka").
//
// Satu penangan di tingkat dokumen (dipasang `RuteNbFacIn`, aktif selama modul terpasang): klik di dalam baris tabel yang sel
// LANGSUNG-nya memuat tombol buka `.nbf-buka` (penanda master-detail khusus modul ini) diteruskan ke tombol itu. Dikecualikan:
// klik pada tombol / tautan / isian / label (aksinya sendiri tetap jalan), klik di baris detail bersarang (bukan baris itu
// sendiri), dan klik yang sedang menyeleksi teks.

const INTERAKTIF = 'button, a, input, select, textarea, label, summary, [role="button"], [role="tab"], [contenteditable="true"]'

/** Bentuk minimum elemen yang dibaca (tanpa `instanceof` - diuji di lingkungan node tanpa DOM). */
interface ElemenMin {
  tagName: string
  children: ArrayLike<ElemenMin>
  classList?: { contains: (c: string) => boolean }
  closest?: (s: string) => ElemenMin | null
  disabled?: boolean
  click?: () => void
}

/** Tombol buka baris yang sel langsungnya memuat `.nbf-buka`, atau null. */
export function tombolBukaBaris(baris: ElemenMin): ElemenMin | null {
  for (const sel of Array.from(baris.children)) {
    if (sel.tagName !== 'TD' && sel.tagName !== 'TH') continue
    const b = Array.from(sel.children).find((x) => x.tagName === 'BUTTON' && x.classList?.contains('nbf-buka') === true)
    if (b) return b
  }
  return null
}

/** Tangani satu klik; true bila diteruskan ke tombol buka baris. */
export function tanganiKlikBaris(
  e: Pick<MouseEvent, 'defaultPrevented' | 'button'> & { target: unknown },
  seleksi: () => string = () => window.getSelection()?.toString() ?? '',
): boolean {
  if (e.defaultPrevented || e.button !== 0) return false
  const t = e.target as ElemenMin | null
  if (!t || typeof t.closest !== 'function') return false
  if (t.closest(INTERAKTIF)) return false
  const baris = t.closest('tr')
  if (!baris) return false
  const tombol = tombolBukaBaris(baris)
  if (!tombol || tombol.disabled === true || typeof tombol.click !== 'function') return false
  if (seleksi().trim() !== '') return false
  tombol.click()
  return true
}

/** Pasang penangan di dokumen; kembalikan pelepasnya. */
export function pasangKlikBaris(doc: Document = document): () => void {
  const h = (e: MouseEvent) => {
    tanganiKlikBaris(e)
  }
  doc.addEventListener('click', h)
  return () => doc.removeEventListener('click', h)
}
