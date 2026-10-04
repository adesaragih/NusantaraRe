// Layar rantai versi kontrak — tiket 01 dan 05 Treaty In Adjustment.
//
// ⛔ BACA-SAJA, dan itu disengaja. Papan tiket menyatakan `L-4` ("tidak ada
// spesifikasi layar di mana pun") dan melarang mengarang layar untuk memenuhi
// bentuk irisan tegak. Halaman ini TIDAK mengarang alur penyesuaian - ia
// memperlihatkan kolom yang tiket 01 dan 05 buat, tidak lebih.
//
// Rupa (02-10-2026): isinya duduk di dalam `Panel`, bukan menempel di `--bg`;
// keadaan kosong memakai `Kosong` (pesan + petunjuk), bukan paragraf telanjang.
// Cabang `Gagal` dan `Memuat` TIDAK disentuh - ketiganya sengaja terpisah,
// lihat `dasar.tsx` baris 344.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Panel, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { ambilKontrak, ambilRantaiVersi, type Kontrak, type Versi } from '../api'
import { MENU_TREATYINADJUSTMENT, RANTAI_VERSI } from '../labels'

/** Label pilihan kontrak: nomor warisan bila ada, selain itu pengenalnya. */
export function labelKontrak(k: Kontrak): string {
  const nomor = k.nomorKontrakWarisan !== '' ? k.nomorKontrakWarisan : `#${k.id}`
  return `${nomor} — ${k.sifatProporsi} (${k.tanggalMulai} s.d. ${k.tanggalBerakhir})`
}

/**
 * Sel yang KOSONG karena memang tidak punya nilai - bukan karena belum diisi.
 *
 * Teksnya redup supaya terbaca sebagai keterangan, bukan sebagai data: "versi
 * pertama" di kolom Dasar adalah jawaban, dan sel yang benar-benar kosong di
 * sebelahnya akan terbaca sebagai data yang hilang.
 */
function Redup({ teks }: { teks: string }) {
  return <span className="tria__redup">{teks}</span>
}

export default function RantaiVersi() {
  const [kontrak, setKontrak] = useState<Kontrak[] | null>(null)
  const [pilih, setPilih] = useState('')
  const [versi, setVersi] = useState<Versi[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    ambilKontrak()
      .then(setKontrak)
      .catch((e: unknown) => {
        setGalat(e)
      })
  }, [])

  useEffect(() => {
    if (pilih === '') {
      setVersi(null)
      return
    }
    let dibuang = false
    setVersi(null)
    setGalat(null)
    ambilRantaiVersi(Number(pilih))
      .then((v) => {
        if (!dibuang) setVersi(v)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [pilih])

  const terpilih = kontrak?.find((k) => String(k.id) === pilih)

  return (
    <div className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{MENU_TREATYINADJUSTMENT.rantaiVersi}</h2>
      </header>
      {/* Judul panel menyebut kontrak yang sedang dibuka begitu ada yang
          dipilih - rantai versi tanpa nama kontraknya menuntut pembacanya
          mengingat apa yang baru saja ia pilih. */}
      <Panel judul={terpilih ? labelKontrak(terpilih) : RANTAI_VERSI.panelPilih}>
        {galat !== null && <Gagal galat={galat} />}
        {kontrak === null && galat === null && <Memuat />}
        {kontrak !== null && kontrak.length === 0 && (
          <Kosong pesan={RANTAI_VERSI.kosong} petunjuk={RANTAI_VERSI.kosongPetunjuk} />
        )}
        {kontrak !== null && kontrak.length > 0 && (
          <div className="form-grid">
            <Pilih
              label={RANTAI_VERSI.pilihKontrak}
              value={pilih}
              onChange={setPilih}
              opsi={kontrak.map((k) => ({ value: String(k.id), label: labelKontrak(k) }))}
            />
          </div>
        )}
        {pilih !== '' && versi === null && galat === null && <Memuat />}
        {versi !== null && (
          // `.table-wrap` menggulir tabel tujuh kolom ini pada layar sempit.
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th scope="col">{RANTAI_VERSI.kolomNomor}</th>
                  <th scope="col">{RANTAI_VERSI.kolomNama}</th>
                  <th scope="col">{RANTAI_VERSI.kolomKeadaan}</th>
                  <th scope="col">{RANTAI_VERSI.kolomJenis}</th>
                  <th scope="col">{RANTAI_VERSI.kolomMaterial}</th>
                  <th scope="col">{RANTAI_VERSI.kolomBerlaku}</th>
                  <th scope="col">{RANTAI_VERSI.kolomDasar}</th>
                </tr>
              </thead>
              <tbody>
                {versi.map((v) => (
                  <tr key={v.id}>
                    <td>
                      {v.nomorUrutVersi ?? <Redup teks={RANTAI_VERSI.belumDinomori} />}
                    </td>
                    <td>{v.namaKontrak}</td>
                    <td>{v.keadaanSiklusHidup}</td>
                    <td>{v.jenisAddendum}</td>
                    <td>{v.sifatMaterialAddendum}</td>
                    <td>{v.tanggalBerlakuAddendum}</td>
                    <td>{v.idVersiDasar ?? <Redup teks={RANTAI_VERSI.versiPertama} />}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Panel>
    </div>
  )
}
