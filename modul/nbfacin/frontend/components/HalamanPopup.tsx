// Pager untuk popup `Modal` ber-`onKirim` (work owner 08-10-2026: popup Change SOB "kenapa tidak bisa ya, di klik
// selanjutnya"). `Modal` ber-`onKirim` merender isinya sebagai `<form>`, sedangkan tombol `Halaman` inti
// (`inti/frontend/components/ui/dasar.tsx`) tidak ber-`type="button"` - jadi klik Berikutnya / Sebelumnya juga men-submit form:
// `onPindah(h)` jalan, lalu `onKirim` memuat ulang halaman 1 dan jawaban halaman h dibuang. Pembungkus ini membatalkan
// submit (preventDefault sesudah onClick tombol berjalan) tanpa mengubah komponen inti milik bersama.

import type { ComponentProps, MouseEvent } from 'react'

import { Halaman } from '../../../../inti/frontend/components/ui/dasar'

/** Klik tombol di dalam pager tidak boleh men-submit form pembungkus. */
export function cegahSubmit(e: Pick<MouseEvent, 'target' | 'preventDefault'>): void {
  if (e.target instanceof Element && e.target.closest('button') !== null) e.preventDefault()
}

export default function HalamanPopup(props: ComponentProps<typeof Halaman>) {
  return (
    <div onClick={cegahSubmit}>
      <Halaman {...props} />
    </div>
  )
}
