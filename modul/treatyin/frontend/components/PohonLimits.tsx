// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { Fragment } from 'react'

import { Kosong, Panel } from '../../../../inti/frontend/components/ui/dasar'
import type { BarisLayerWarisan } from '../api'
import { FORM_KONTRAK, LIMITS_POHON } from '../labels'
import { selAngka } from './angka'

/**
 * Panel `Total Retention Amount` beserta tombol `Update Total` — keduanya
 * di bawah grid Maximum Retention, persis urutan layar lama.
 *
 * ⛔ Nilainya TIDAK dihitung di sini. Ia datang sudah jadi dari services,
 * yang menyalin rumusnya dari `Activity/TreatyInNPSetTotal.xml` cabang
 * `param.type=retention`. Menghitung ulang di layar berarti rumus yang sama
 * hidup di dua tempat, dan yang kedua selalu yang basi.
 */
/**
 * Panel `Existing Policy for Master ID` — kanan atas layar, KEDUA cabang.
 *
 * ⛔ Panel ini NOL di kode kita sampai 5 Oktober 2026, dan sumbernya bukan
 * dikarang: `Section/InputTreatyInOffer.xml` @89.949 ->
 * `Activity/FetchTreatyExistingProduction.xml` ->
 * `RDBList/FetchTreatyInProductionUsingNooffer.xml`, yang membaca
 * `POOLDATA.TREATYINPRODUCTION`.
 *
 * ⚠️ Nol baris adalah keadaan yang LAZIM — gambar 01 (kontrak 1001846)
 * memperlihatkannya berbunyi `No items`, dan pengukuran membenarkannya.
 */
/** Grid `RNM Share` — satu bentuk, dipakai sebagai tab DAN sebagai sub-tab. */
/**
 * Tab **Limits** — POHON TIGA TINGKAT, bukan grid datar.
 *
 * ⛔ Delapan belas gambar untuk satu tab (`02`–`15` prop, `30`–`33` non-prop),
 * dan sebabnya struktural. Susunan tingkatnya di `LIMITS_POHON`, beserta
 * bukti mengapa kedua cabang berbeda di puncak.
 *
 * ⚠️ Yang dibangun di sini KERANGKANYA — tiga tingkat yang dapat dibuka,
 * dengan medan yang jalur baca kita sudah punya. Yang BELUM:
 *
 *   - `Treaty Type` di tingkat puncak prop (gambar 03) — calonnya
 *     `Detail[].SpreadingType`, tetapi ia tingkat `Detail[]`; §14.
 *   - 11 sub-tab tingkat terdalam (gambar 05–15).
 *   - Panel `Summary of Limit` dan `Total All Layers` (gambar 30) — rumusnya
 *     menunggu jawaban `PremiumEarnedList`/`MDPList`.
 *
 * ⛔ Nol di antaranya dibangun atas tebakan.
 */
export default function PohonLimits({
  layer,
  nonProp,
}: {
  layer: readonly BarisLayerWarisan[]
  nonProp: boolean
}) {
  // ⛔ Pengelompokan puncak BERBEDA per cabang, dan medannya pun berbeda:
  // `jenisTreaty` terisi pada prop, `layer` pada non-prop. Memakai satu
  // medan untuk keduanya memberi satu kelompok bernama kosong di cabang
  // yang medannya nil.
  const puncak: { kunci: string; baris: BarisLayerWarisan[] }[] = []
  for (const b of layer) {
    const k = nonProp ? b.layer : b.jenisTreaty
    const ada = puncak.find((p) => p.kunci === k)
    // ⚠️ Urutan kemunculan, bukan abjad — urutan dokumen adalah urutan yang
    // layar lama tampilkan.
    if (ada) ada.baris.push(b)
    else puncak.push({ kunci: k, baris: [b] })
  }

  return (
    <Panel judul="Limits">
      <div className="trin__pohon" role="group" aria-label="Limits">
        <div className="trin__pohon-kepala">
          <span>{nonProp ? LIMITS_POHON.puncakNonProp : LIMITS_POHON.puncakProp}</span>
          {/* ⛔ MATI — nol jalur tulis. */}
          <button type="button" className="btn btn--ghost btn--sm" disabled>
            {nonProp ? LIMITS_POHON.tambahLayer : LIMITS_POHON.tambah}
          </button>
        </div>

        {puncak.length === 0 && (
          <Kosong pesan={LIMITS_POHON.tanpaBaris} petunjuk={FORM_KONTRAK.petunjukLayer} />
        )}

        {puncak.map((p) => (
          <details key={p.kunci} className="trin__pohon-tingkat">
            <summary>
              {/* Kelompok tanpa nama tetap terlihat — ia ada di dokumen. */}
              {p.kunci === '' ? LIMITS_POHON.tanpaBaris : p.kunci}
              {nonProp && p.baris[0] !== undefined && p.baris[0].jenisLayer !== '' && (
                <span className="trin__redup"> · {p.baris[0].jenisLayer}</span>
              )}
            </summary>

            {/* Medan tingkat LAYER — satu nilai untuk seluruh kelompok. */}
            <PerincianLayer baris={p.baris[0]} />

            <div className="trin__pohon-kepala">
              <span>{LIMITS_POHON.kelompokTreaty}</span>
              <button type="button" className="btn btn--ghost btn--sm" disabled>
                {nonProp ? LIMITS_POHON.tambahKelompok : LIMITS_POHON.tambah}
              </button>
            </div>

            {p.baris.map((b, i) => (
              <details key={b.kelompokTreaty + String(i)} className="trin__pohon-tingkat">
                <summary>
                  {b.kelompokTreaty === '' ? LIMITS_POHON.tanpaBaris : b.kelompokTreaty}
                </summary>
                <div className="trin__pohon-kepala">
                  <span>{LIMITS_POHON.kelasBisnis}</span>
                </div>
                {/* ⚠️ Pada kontrak contoh gambar 05 dan 31 grid ini berbunyi
                    `No items`; terukur 15.742 elemen di seluruh dokumen,
                    jadi ia berisi pada kontrak lain. */}
                {b.kelasBisnis.length === 0 ? (
                  <p className="trin__redup">{LIMITS_POHON.tanpaBaris}</p>
                ) : (
                  <ul className="trin__pohon-daun">
                    {b.kelasBisnis.map((k, j) => (
                      <li key={k + String(j)}>{k}</li>
                    ))}
                  </ul>
                )}
              </details>
            ))}
          </details>
        ))}
      </div>
    </Panel>
  )
}

/**
 * Medan tingkat LAYER — nilai yang berlaku untuk seluruh kelompok.
 *
 * ⚠️ `Deductible` ditandai REDUP beserta angka keyakinannya. Padanan
 * `CEDANT_RETENTION` ↔ `.Deductible` cocok **70,6%** atas 1.340 kontrak —
 * jauh di atas dugaan 51% dan cukup untuk memetakan, TIDAK cukup untuk
 * menampilkannya seolah pasti.
 */
function PerincianLayer({ baris }: { baris: BarisLayerWarisan | undefined }) {
  if (baris === undefined) return null
  const medan: { nama: string; nilai: string; belumPasti?: boolean }[] = [
    { nama: 'Cover', nilai: baris.dasarCover },
    { nama: 'Currency Relation', nilai: baris.relasiMataUang },
    { nama: 'Currency', nilai: baris.mataUang },
    { nama: '100% Limit', nilai: selAngka('uang', baris.limit100) },
    // ⛔ Satu-satunya medan yang ditandai belum pasti.
    { nama: 'Deductible', nilai: selAngka('uang', baris.retensiCedant), belumPasti: true },
    { nama: 'MDP', nilai: selAngka('uang', baris.mdp) },
    { nama: 'MDP %', nilai: selAngka('persen', baris.rasioMDP) },
    { nama: 'ROL %', nilai: selAngka('persen', baris.rol) },
    { nama: 'Adjustment Rate %', nilai: selAngka('persen', baris.adjRate) },
    { nama: 'Premium Earned', nilai: selAngka('uang', baris.premiEarned) },
  ]
  return (
    <dl className="trin__pohon-medan">
      {medan.map((m) => (
        <Fragment key={m.nama}>
          <dt>{m.nama}</dt>
          <dd>
            {m.nilai === '' ? <span className="trin__redup">&mdash;</span> : m.nilai}
            {m.belumPasti === true && (
              <span className="trin__redup" title={LIMITS_POHON.deductiblePetunjuk}>
                {' '}
                {LIMITS_POHON.belumPasti}
              </span>
            )}
          </dd>
        </Fragment>
      ))}
    </dl>
  )
}
