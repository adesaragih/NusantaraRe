// Halaman depan NB FacIn - portal Opportunity (tiket 25).
//
// Port harness `SFAPortalOpportunities` (`D:\migrasi\RNM\NB FacIn\Harness\SFAPortalOpportunities.xml`,
// PEGACRM-PORTAL!SFAPORTALOPPORTUNITIES) beserta section `SFAPortalOpportunitiesHeader` dan
// `SFAPortal_OpportunitiesList`. Hanya unsur yang TAMPIL di Pega yang diport; yang tersembunyi
// permanen (`1=2` / `NEVER`: Stage view, List view, All/Individual/Corporate, Phase, Export, Refresh)
// tidak dirender sama sekali.
//
// Grid = `GetListOpportunityF` (keputusan agent B-1). Sumber datanya - RD atas tabel kerja Pega -
// belum ada padanannya (B-2), jadi daftar tampil sebagai `BelumTersedia`, BUKAN tabel kosong: tabel
// kosong berbunyi "tidak ada opportunity", dan itu tidak benar. Kotak saring dan Filter nonaktif
// selama daftar tidak dapat dimuat - menyaring daftar yang tidak ada hanya memberi kesan bekerja.
//
// `Create opportunity` membuka form Opportunity (B-3, tiket 26). Nol catatan pengembang di layar (B-4).

import { BelumTersedia } from '../../../../inti/frontend/components/ui/dasar'
import { KEPALA_PORTAL, KOLOM_PORTAL, SARING_PORTAL, TEKS_PORTAL } from '../labels'

export default function PortalOpportunity({ onBuat }: { onBuat: () => void }) {
  return (
    <div className="nbfacin">
      <section className="panel">
        <div className="nbf-kepala">
          <h4 className="panel__title">{KEPALA_PORTAL.judul.label}</h4>
          <button type="button" className="btn btn--primary" onClick={onBuat}>
            {KEPALA_PORTAL.buat.label}
          </button>
        </div>
        <div className="nbf-saring">
          <input
            className="field__input nbf-saring__kotak"
            type="text"
            aria-label={SARING_PORTAL.label.label}
            placeholder={SARING_PORTAL.placeholder.label}
            value=""
            disabled
            readOnly
          />
          <button type="button" className="btn btn--ghost btn--sm" aria-label={TEKS_PORTAL.hapusIsian} disabled>
            ✕
          </button>
          <button type="button" className="btn btn--sm" disabled>
            {SARING_PORTAL.tombol.label}
          </button>
        </div>
        <div className="table-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                {KOLOM_PORTAL.map((k) => (
                  <th key={k.sel} scope="col">
                    {k.label}
                  </th>
                ))}
              </tr>
            </thead>
          </table>
        </div>
        <BelumTersedia apa={TEKS_PORTAL.daftar} />
      </section>
    </div>
  )
}
