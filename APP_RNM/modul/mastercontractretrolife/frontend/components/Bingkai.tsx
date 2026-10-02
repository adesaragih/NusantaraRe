// Kerangka panel bersama keempat harness popup Pega.
//
// Pega membuka `InboxRetroLimitReinsurers`, `InboxRetroLifeReinsurersList`,
// `InboxSecurityReinsurerLife`, dan `InboxBusinessLifeReinsurers` sebagai POPUP (`showHarness`
// `pyTarget=popup`) di atas layar pemanggilnya. Di sini: panel yang menggantikan layar pemanggil
// selama terbuka, dan tombol tutup mengembalikannya - keadaan layar pemanggil tetap (komponennya
// tidak dilepas).
//
// ⚠️ PENYIMPANGAN SADAR TAMPILAN (keputusan work owner 02-10-2026: "perbaiki tampilannya sama seperti
// treaty contract out", lingkup tampilan saja): kepala panel meniru Treaty Contract Out - tombol TEKS
// `Close` (Pega: ikon `pxIconCancel` tanpa teks) dan medan baca-saja sebagai kotak isian hanya-baca
// (`Field readOnly` di `form-grid`, Pega: `pyLabelFieldValue`). Label tetap VERBATIM; alur, pager,
// dan judul kolom (huruf besar Pega) tidak berubah.
//
// ⚠️ Kecuali `InboxSecurityReinsurerLife` sejak 02-10-2026 (keputusan work owner: akses Security Reinsurer
// SAMA dengan Treaty Contract Out): panel Security tampil di bawah daftar reinsurer, di dalam panel
// Reinsurer, tanpa menggantikannya - lihat `PanelReinsurer.tsx`.

import { Field, Halaman } from '../../../../inti/frontend/components/ui/dasar'
import { UMUM_MCRL } from '../labels'
import { UKURAN_HALAMAN_MCRL, sel } from '../tampilan'

/** Kepala panel: judul section + tombol `Close` + medan baca-saja (label VERBATIM) - pola Treaty Contract Out. */
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
        <button type="button" className="btn btn--ghost btn--sm" onClick={onTutup}>
          {UMUM_MCRL.tutup}
        </button>
      </header>
      <div className="form-grid">
        {medan.map(([label, nilai]) => (
          <Field key={label} label={label} value={sel(nilai)} onChange={() => undefined} readOnly />
        ))}
      </div>
    </>
  )
}

/** `pyGridPaginator` - 10 baris per halaman; tidak tampil untuk grid kosong. */
export function Penomoran({ halaman, total, onPindah }: { halaman: number; total: number; onPindah: (h: number) => void }) {
  if (total === 0) return null
  return <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN_MCRL} total={total} onPindah={onPindah} />
}
