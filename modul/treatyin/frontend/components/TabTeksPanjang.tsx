// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { Kosong, Panel } from '../../../../inti/frontend/components/ui/dasar'
import type { TabTeksWarisan } from '../api'
import { FORM_KONTRAK } from '../labels'

/**
 * Tab yang isinya SATU medan teks panjang — Exclusions, Special Conditions.
 *
 * ⛔ Jalan B, keputusan §15: teksnya dibaca dari `JSONDATA`, nol tabel
 * pendaratan. Yang dikerjakan layar hanya menampilkannya.
 *
 * ⛔ TIGA keadaan, dan ketiganya terlihat berbeda:
 *
 *   ada isinya            teksnya, dapat digulir, plus ejaan asalnya.
 *   ejaan cabang kosong   `Kosong` — "belum ada DATA". TIDAK diisi dari
 *                         ejaan lain: isinya BERBEDA, bukan salinan basi.
 *   ejaan lain juga ada   peringatan yang MENYEBUT ejaannya, supaya
 *                         pembacanya tahu ada teks yang ia tidak lihat.
 */
export default function TabTeksPanjang({
  judul,
  tab,
  petunjukKosong,
}: {
  judul: string
  tab: TabTeksWarisan | undefined
  petunjukKosong: string
}) {
  const isi = tab?.isi ?? ''
  const lain = tab?.ejaanLain ?? []
  return (
    <Panel judul={judul}>
      {lain.length > 0 && (
        <span className="trin__teks-lain" role="note">
          {FORM_KONTRAK.ejaanLainBerisi} {lain.join(', ')}
        </span>
      )}
      {isi === '' ? (
        <Kosong pesan={FORM_KONTRAK.tanpaTeks} petunjuk={petunjukKosong} />
      ) : (
        <>
          <span className="trin__teks-asal">
            {FORM_KONTRAK.ejaanDipakai} {tab?.ejaan}
          </span>
          <div className="trin__teks" tabIndex={0}>
            {isi}
          </div>
        </>
      )}
    </Panel>
  )
}
