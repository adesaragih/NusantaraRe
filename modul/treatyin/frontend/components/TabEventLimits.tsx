// Tab **Event Limits** — EMPAT BARIS BERLABEL, bukan grid datar.
//
// ⛔ Bentuknya dari gambar `28`, bukan dari bentuk datanya. Pega menyusunnya
// sebagai empat baris berlabel, masing-masing dengan mata uangnya sendiri
// di sebelah kiri dan nilainya di kanan:
//
//	RSMD Limit                  [ Currency ]  [ 0,00 ]
//	Earthquake Limit            [ Currency ]  [ 0,00 ]
//	Flood Limit (Jabodetabek)   [ Currency ]  [ 0,00 ]
//	Flood Limit (Nationwide)    [ Currency ]  [ 0,00 ]
//
// ⚠️ Bentuk sebelumnya grid sembilan kolom — benar isinya, salah bentuknya.
// Sembilan kolom sempit berisi empat pasang nilai memaksa pembacanya
// mencocokkan kepala kolom dengan sel, sementara layar lama menaruh labelnya
// tepat di samping nilainya.
//
// ⛔ SATU BEDA dari gambar, dan ia disengaja: gambar memperlihatkan SATU
// himpunan empat baris, sementara dokumen menyimpannya per
// LAYER × TREATY GROUP. Menampilkan satu himpunan berarti memilih satu baris
// dan menyembunyikan sisanya — persis cacat `nilaiPertama` yang baru ditutup
// 6 Oktober 2026. Jadi tiap baris layer memperoleh himpunannya sendiri, dan
// kepalanya menyebut layer mana.
//
// ⚠️ Pada kontrak contoh gambar 28 keempat nilainya `0,00` dan nol layer
// terlihat — jadi gambar itu TIDAK memperlihatkan apa yang Pega lakukan
// ketika layernya lebih dari satu. Pertanyaannya di
// `PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §15.

import { Kosong, Panel } from '../../../../inti/frontend/components/ui/dasar'
import type { BarisLayerWarisan } from '../api'
import { EVENT_LIMITS, FORM_KONTRAK } from '../labels'
import { selAngka } from './angka'

export default function TabEventLimits({
  layer,
}: {
  layer: readonly BarisLayerWarisan[]
}) {
  return (
    <Panel judul={EVENT_LIMITS.judul}>
      {layer.length === 0 && (
        <Kosong pesan={FORM_KONTRAK.tanpaBaris} petunjuk={FORM_KONTRAK.petunjukEventLimits} />
      )}
      {layer.map((b, i) => (
        <section key={b.layer + '|' + b.kelompokTreaty + '|' + String(i)} className="trin__ev">
          {/* ⚠️ Kepala menyebut layer DAN treaty group: dua baris dapat
              bernomor layer sama tetapi bertreaty group berbeda, dan tanpa
              keduanya himpunan di bawahnya tidak dapat ditelusuri. */}
          <h5 className="trin__ev-kepala">
            {EVENT_LIMITS.layer} {b.layer === '' ? '—' : b.layer}
            {b.kelompokTreaty !== '' && (
              <span className="trin__redup"> · {b.kelompokTreaty}</span>
            )}
          </h5>
          <dl className="trin__ev-medan">
            {(
              [
                [EVENT_LIMITS.rsmd, b.mataUangRSMD, b.batasRSMD],
                [EVENT_LIMITS.gempa, b.mataUangGempa, b.gempa],
                [EVENT_LIMITS.banjirJab, b.mataUangBanjirJab, b.batasBanjirJab],
                [EVENT_LIMITS.banjirNas, b.mataUangBanjirNas, b.batasBanjirNas],
              ] as const
            ).map(([label, mataUang, nilai]) => (
              <div key={label} className="trin__ev-baris">
                <dt>{label}</dt>
                {/* ⛔ Mata uang TEKS, tidak pernah diformat. */}
                <dd className="trin__ev-mu">
                  {mataUang === '' ? <span className="trin__redup">—</span> : mataUang}
                </dd>
                {/* ⚠️ Desimalnya `null` — gambar 28 memperlihatkan `0,00`
                    pada kontrak yang keempat nilainya nol, dan nol tidak
                    memberitahu presisi kolomnya. Jadi aturan lama berlaku,
                    dan kolom ini masuk daftar §13 pertanyaan terbuka. */}
                <dd className="trin__ev-nilai">
                  {nilai === '' ? <span className="trin__redup">—</span> : selAngka('uang', nilai)}
                </dd>
              </div>
            ))}
          </dl>
        </section>
      ))}
    </Panel>
  )
}
