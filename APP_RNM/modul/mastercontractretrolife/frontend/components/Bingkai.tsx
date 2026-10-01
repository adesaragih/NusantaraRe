// Kerangka panel bersama keempat harness popup Pega.
//
// Pega membuka `InboxRetroLimitReinsurers`, `InboxRetroLifeReinsurersList`,
// `InboxSecurityReinsurerLife`, dan `InboxBusinessLifeReinsurers` sebagai POPUP (`showHarness`
// `pyTarget=popup`) di atas layar pemanggilnya. Di sini: panel yang menggantikan layar pemanggil
// selama terbuka, dan ikon tutup kerangka harness (`pxIconCancel`) mengembalikannya - keadaan layar
// pemanggil tetap (komponennya tidak dilepas).

import { Halaman, IkonTutup } from '../../../../inti/frontend/components/ui/dasar'
import { UMUM_MCRL } from '../labels'
import { UKURAN_HALAMAN_MCRL, sel } from '../tampilan'

/** Kepala panel: judul section + medan baca-saja (label VERBATIM) + ikon tutup harness. */
export function KepalaPanel({
  judul,
  medan,
  onTutup,
}: {
  judul: string
  medan: ReadonlyArray<readonly [label: string, nilai: string]>
  onTutup: () => void
}) {
  return (
    <>
      <header className="inbox__kepala">
        <h3 className="panel__title">{judul}</h3>
        <button type="button" className="btn btn--ghost btn--sm mcrl-tutup" aria-label={UMUM_MCRL.tutup} title={UMUM_MCRL.tutup} onClick={onTutup}>
          <IkonTutup />
        </button>
      </header>
      <dl className="mcrl-kepala">
        {medan.map(([label, nilai]) => (
          <div key={label} className="mcrl-kepala__medan">
            <dt>{label}</dt>
            <dd>{sel(nilai)}</dd>
          </div>
        ))}
      </dl>
    </>
  )
}

/** `pyGridPaginator` - 10 baris per halaman; tidak tampil untuk grid kosong. */
export function Penomoran({ halaman, total, onPindah }: { halaman: number; total: number; onPindah: (h: number) => void }) {
  if (total === 0) return null
  return <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN_MCRL} total={total} onPindah={onPindah} />
}
