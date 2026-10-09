// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { Panel } from '../../../../inti/frontend/components/ui/dasar'
import type { TabTeksWarisan } from '../api'
import { FORM_KONTRAK } from '../labels'
import { bacaProperti, usePenampung, useProperti } from '../halaman'
import type { ModeForm } from '../mode'
import { SelKosongPega, TataPegaBlok } from './tataPega'

/**
 * DT `TreatyInCopyConditions` — perilaku `change` `ExclusionsP` dan
 * `SpecialConditionsP` (Section `TreatyInTabsProportional`):
 *   1  `TreatyIn.Exclusions = TreatyIn.ExclusionsP`
 *   2  `TreatyIn.SpecialConditions = TreatyIn.SpecialConditionsP`
 * KEDUANYA disalin, dari medan mana pun ia dipicu.
 */
const SALIN_SYARAT: Readonly<Record<string, string>> = {
  ExclusionsP: 'Exclusions',
  SpecialConditionsP: 'SpecialConditions',
}

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
  properti,
  tab,
  petunjukKosong,
  mode = 'lihat',
}: {
  judul: string
  /**
   * Properti halaman `TreatyIn` tab ini — `ExclusionsP` / `SpecialConditionsP`
   * (Prop) atau `Exclusions` / `SpecialConditions` (Non-Prop), ejaan ekspor.
   */
  properti: string
  tab: TabTeksWarisan | undefined
  petunjukKosong: string
  mode?: ModeForm
}) {
  const asli = tab?.isi ?? ''
  const lain = tab?.ejaanLain ?? []
  const bisaUbah = mode === 'ubah'
  // ⭐ Teks yang sedang disunting. Disalin SEKALI saat tab lahir; pemanggil
  // memberi tab lain lewat `judul` yang berbeda, dan React melahirkan
  // komponen baru — pola yang sama dengan `TabGridWarisan`.
  // ⭐ Penampung halaman — isian bertahan saat pindah tab (`halaman.tsx`).
  const [teks, setTeks] = useProperti(properti, asli)
  const isi = teks
  const penampung = usePenampung()
  /**
   * ⭐ `TreatyInCopyConditions` saat isian DITINGGALKAN (`change` Pega).
   * ⛔ Sumber yang belum pernah dibuka (belum ada di halaman) TIDAK disalin
   * — nilai tersimpannya tidak diketahui layar ini, dan menyalin kosong
   * akan menghapusnya.
   */
  const salinSyarat = () => {
    if (penampung === null || SALIN_SYARAT[properti] === undefined) return
    for (const [sumber, tujuan] of Object.entries(SALIN_SYARAT)) {
      const v = sumber === properti ? teks : bacaProperti(penampung.halaman, sumber)
      if (typeof v === 'string') penampung.ubah(tujuan, () => v)
    }
  }
  return (
    <Panel judul={judul}>
      {lain.length > 0 && (
        <span className="trin__teks-lain" role="note">
          {FORM_KONTRAK.ejaanLainBerisi} {lain.join(', ')}
        </span>
      )}
      {/* ⭐ MODE UBAH — textarea, seperti layar lama.
          ⛔ Sebelum 6 Oktober 2026 tab ini BACA-SAJA walau mode Edit, dan
          itu bertentangan dengan keputusan pemilik proses bahwa seluruh
          fungsi dapat dipakai di mode Edit. Kosong pun dapat diketik: tab
          yang menolak isian pertama tidak akan pernah terisi. */}
      {/* ⭐ TATA LETAK PEGA (8 Oktober 2026) — `Inline grid double` dengan
          SATU sel `pxTextArea`: kotak teks di separuh KIRI, separuh kanan
          kosong. `Section/TreatyInTabsProportional.xml` L49386 (ExclusionsP)
          dan L50813 (SpecialConditionsP); `TreatyInTabsNonProportional.xml`
          L134068 / L135274 sama. */}
      <TataPegaBlok tata="g2">
      <div>
      {bisaUbah ? (
        <textarea
          className="field__input trin__teks-isian"
          aria-label={judul}
          value={teks}
          rows={16}
          onChange={(e) => {
            setTeks(e.target.value)
          }}
          onBlur={salinSyarat}
        />
      ) : isi === '' ? (
        // ⭐ Bentuk Pega (8 Oktober 2026): `pxTextArea` baca-saja yang KOSONG —
        // kotak teks tanpa isi, bukan ikon "tanpa teks".
        <div className="trin__teks" tabIndex={0} aria-label={judul} title={petunjukKosong} />
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
      </div>
      <SelKosongPega />
      </TataPegaBlok>
    </Panel>
  )
}
