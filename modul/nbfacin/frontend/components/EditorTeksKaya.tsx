// Editor teks berformat untuk Message korespondensi (tiket 56) - padanan pxRichTextEditor Pega
// (`NB FacIn\Section\CorrespondenceContent.xml`, `.Message`). Tanpa pustaka baru: `contentEditable` + `execCommand`.
// Toolbar = bagian toolbar Pega yang hasilnya lolos pembersih (`htmlKorespondensi.ts`): Bold, Italic, Underline, daftar
// berpoin / bernomor, hapus format, Undo, Redo. Font, Size, warna teks / latar TIDAK disediakan - memerlukan atribut gaya
// yang dibuang pembersih (keputusan agent K56-3). Isi awal ditulis sekali saat dipasang (komponen tidak dikendalikan
// React sesudahnya - pemanggil memberi `key` per email); setiap ketikan dikirim ke `onChange` sesudah dibersihkan.

import { useEffect, useRef } from 'react'

import { TEKS_KORESPONDENSI as T } from '../labels'
import { bersihkanHtml } from './htmlKorespondensi'

const TOMBOL: readonly { perintah: string; label: string; tanda: string }[] = [
  { perintah: 'bold', label: T.editor.tebal, tanda: 'B' },
  { perintah: 'italic', label: T.editor.miring, tanda: 'I' },
  { perintah: 'underline', label: T.editor.garisBawah, tanda: 'U' },
  { perintah: 'insertUnorderedList', label: T.editor.daftarTitik, tanda: '•' },
  { perintah: 'insertOrderedList', label: T.editor.daftarNomor, tanda: '1.' },
  { perintah: 'removeFormat', label: T.editor.hapusFormat, tanda: 'Tx' },
  { perintah: 'undo', label: T.editor.urung, tanda: '↶' },
  { perintah: 'redo', label: T.editor.ulang, tanda: '↷' },
]

export default function EditorTeksKaya({
  label,
  awal,
  onChange,
  bacaSaja = false,
}: {
  label: string
  /** HTML awal (sudah dibersihkan server); dibersihkan lagi sebelum ditaruh. */
  awal: string
  onChange: (html: string) => void
  bacaSaja?: boolean
}) {
  const kotak = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (kotak.current) kotak.current.innerHTML = bersihkanHtml(awal)
    // Sengaja hanya saat dipasang: menulis ulang innerHTML tiap render memindahkan kursor.
  }, [])

  const jalankan = (perintah: string) => {
    kotak.current?.focus()
    document.execCommand(perintah)
    if (kotak.current) onChange(bersihkanHtml(kotak.current.innerHTML))
  }

  return (
    <div className="nbf-teks-kaya">
      <span className="nbf-teks-kaya__label">{label}</span>
      {!bacaSaja && (
        <div className="nbf-teks-kaya__alat" role="toolbar" aria-label={label}>
          {TOMBOL.map((t) => (
            <button
              key={t.perintah}
              type="button"
              className="btn nbf-tema__tombol btn--ghost btn--sm"
              title={t.label}
              aria-label={t.label}
              onMouseDown={(e) => e.preventDefault()}
              onClick={() => jalankan(t.perintah)}
            >
              {t.tanda}
            </button>
          ))}
        </div>
      )}
      <div
        ref={kotak}
        className="nbf-teks-kaya__isi"
        role="textbox"
        aria-multiline="true"
        aria-label={label}
        aria-readonly={bacaSaja}
        contentEditable={!bacaSaja}
        suppressContentEditableWarning
        onInput={(e) => onChange(bersihkanHtml(e.currentTarget.innerHTML))}
      />
    </div>
  )
}
