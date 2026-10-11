// Menu View Policy (tiket 71) - port section Pega `DATA-PORTAL!VIEWOPENPOLICYMARINEINBOX`
// (`D:\migrasi\RNM\DDL\ViewOpenPolicyMarineInbox.xml`), tiga tab (work owner 11-10-2026 "buatin untuk life dan treaty juga"):
// - Facultative: List Polis (SEARCHON POLIFAC, grid `GridOpenPolicyMarine`, View Option FilterForViewRetro, detail versi /
//   ringkasan / layar case - K71-1 / K71-2) dan List Offer (OFFERFAC, `BrowseDataViewOfferF` atas FACINOFFER, grid
//   `GridOpenoffer`; View Option = polis Facultative tanpa QQ Name - keputusan work owner, FACINOFFER tanpa kolom QQ);
// - Treaty: List Polis (TREATYPOL, `BrowsePolisTreaty` atas TREATYINPRODUCTION, grid `GridOpenPolicyMarineTreaty`, View Option
//   `FilterForViewPolicy`). List Offer Treaty TIDAK dibangun: grid `GridOpenOfferTreaty` tidak dikirim dan
//   `GetListViewPolicy_act` tidak punya SEARCHON-nya;
// - Life: List Polis (POLISLIFE, `BrowsePolisLife` atas LIFEINPRODUCTION, grid `GridOpenPolicyLife`) dan List Offer
//   (OFFERLIFE, `BrowseOfferLife` atas JSON_OFFER_LIFE, grid `GridOpenOfferLife`); tanpa View Option (work owner 11-10-2026).
// View Details Treaty (work owner 11-10-2026) = berkas NB / EDM Treaty In polis itu (pola View polis Claim Prop / Claim Non
// Prop OQ-CNP-13), dibuka DI TEMPAT seperti detail NB FacIn - `PropsRute.onBukaModul`; Back berkas itu kembali ke View Policy
// (work owner: "jangan pop up doang, buatin kaya nb facin"). Tanpa `onBukaModul`: jendela `onLihatBerkas`. Life belum diport.
// Semua tab tetap terpasang (`hidden`) supaya isian cari tiap tab bertahan saat berpindah tab.

import { useState } from 'react'

import { ApiFailure } from '../../../../inti/frontend/klien'

import { Gagal, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import {
  berkasPolisTreaty,
  cariOfferFacultative,
  cariOfferLife,
  cariPolis,
  cariPolisLife,
  cariPolisTreaty,
  type BarisFollowing,
  type BarisOfferFacultative,
  type BarisOfferLife,
  type BarisPolisLife,
  type BarisPolisTreaty,
  type OpsiFollowing,
} from '../api'
import DaftarCariPolis from '../components/DaftarCariPolis'
import { DetailPolis } from '../components/DetailPolis'
import { tanggalFollowing } from '../components/PopupFollowing'
import { FOLLOWING as F, VIEW_POLICY as V } from '../labels'
import InwardFacultative from './InwardFacultative'

// Detail polis (versi, ringkasan polis lama) dipindah ke components/DetailPolis.tsx - dipakai juga popup Following.
export { labelPolicyType, RingkasanPolisLama } from '../components/DetailPolis'

type TabVP = (typeof V.tab)[number]
type Detail = { jenis: 'polis'; baris: BarisFollowing } | { jenis: 'kasus'; caseId: string } | null

/** `onBukaModul` / `onLihatBerkas` - `PropsRute` (View Details Treaty: di tempat, cadangan jendela); tanpa keduanya tombolnya tidak ada. */
export default function ViewPolicy({
  onLihatBerkas,
  onBukaModul,
}: { onLihatBerkas?: (modul: string, id: string) => boolean; onBukaModul?: (modul: string, id: string) => boolean } = {}) {
  const bukaBerkas = onBukaModul ?? onLihatBerkas
  const [tab, setTab] = useState<TabVP>(V.tab[0])
  const [detail, setDetail] = useState<Detail>(null)
  const [pesanTreaty, setPesanTreaty] = useState('')
  const [galatTreaty, setGalatTreaty] = useState<unknown>(null)

  const lihatTreaty = (b: BarisPolisTreaty) => {
    setPesanTreaty('')
    setGalatTreaty(null)
    berkasPolisTreaty(b.policyNo).then(
      (x) => {
        if (bukaBerkas?.(x.modul, x.kasus) !== true) setPesanTreaty(V.teks.berkasTakTerpasang)
      },
      (g: unknown) => {
        // Hanya 404 JSON dari backend (`DITOLAK_BACKEND`) = polis tanpa berkas. 404 teks biasa = rute belum terpasang (backend
        // lama belum direstart) -> galat apa adanya, bukan pesan "polis lama" (bug DEV 11-10-2026).
        if (g instanceof ApiFailure && g.status === 404 && g.detail.code === 'DITOLAK_BACKEND') setPesanTreaty(V.teks.polisTanpaBerkas(b.policyNo))
        else setGalatTreaty(g)
      },
    )
  }

  // Templat layar = menu Kelola User (work owner 09-10-2026; pola portal NB FacIn): kepala `inbox__kepala` + judul besar,
  // switch tab, lalu daftar-daftar cari tab itu.
  return (
    <div className="nbfacin">
      <section className="inbox nbf-tema">
        {detail?.jenis === 'polis' && (
          <DetailPolis
            baris={detail.baris}
            judul={`${V.judul} · ${detail.baris.policyNo}`}
            onKembali={() => setDetail(null)}
            tampilKasus={(caseId) => <InwardFacultative kasus={{ caseId }} onBatal={() => setDetail(null)} hanyaLihat />}
          />
        )}
        {/* View Details List Offer Facultative: case NB (FACINOFFER.IDPEGA = nomor case) dibuka baca-saja. */}
        {detail?.jenis === 'kasus' && <InwardFacultative kasus={{ caseId: detail.caseId }} onBatal={() => setDetail(null)} hanyaLihat />}
        <div hidden={detail !== null}>
          <header className="inbox__kepala">
            <h2 className="inbox__judul nbf-tema__judul">{V.judul}</h2>
          </header>
          <StripTab tab={V.tab} aktif={tab} onPilih={setTab} />

          <div hidden={tab !== 'Facultative'}>
            <DaftarCariPolis<BarisFollowing>
              judul={V.daftar}
              opsi={F.opsi}
              placeholder={F.placeholder}
              kolom={F.kolom}
              cari={(o, f, h) => cariPolis(o as OpsiFollowing, f, h)}
              kunci={(b) => `${b.caseId}-${b.policyNo}`}
              sel={(b) => [
                b.caseId,
                b.policyNo,
                b.reffNumber,
                b.insuredName,
                b.sourceOfBusiness,
                b.qqName,
                tanggalFollowing(b.startPeriod),
                tanggalFollowing(b.endPeriod),
              ]}
              onLihat={(b) => setDetail({ jenis: 'polis', baris: b })}
              bisaLihat={(b) => b.policyNo !== ''}
            />
            <DaftarCariPolis<BarisOfferFacultative>
              judul={V.daftarOffer}
              opsi={V.opsiOfferFacultative}
              placeholder={V.placeholder.offer}
              kolom={V.kolom.offerFacultative}
              cari={cariOfferFacultative}
              kunci={(b) => `${b.caseId}-${b.noOffer}`}
              sel={(b) => [b.caseId, b.noOffer, b.insuredName, b.sourceOfBusiness, b.qqName, tanggalFollowing(b.startPeriod), tanggalFollowing(b.endPeriod)]}
              onLihat={(b) => setDetail({ jenis: 'kasus', caseId: b.caseId })}
              bisaLihat={(b) => b.caseId !== ''}
            />
          </div>

          <div hidden={tab !== 'Treaty'}>
            {pesanTreaty !== '' && <div className="alert alert--info">{pesanTreaty}</div>}
            <Gagal galat={galatTreaty} />
            <DaftarCariPolis<BarisPolisTreaty>
              judul={V.daftar}
              opsi={V.opsiTreaty}
              placeholder={V.placeholder.polis}
              kolom={V.kolom.polisTreaty}
              cari={cariPolisTreaty}
              kunci={(b) => `${b.caseId}-${b.policyNo}`}
              sel={(b) => [b.caseId, b.policyNo, b.insuredName, b.sourceOfBusiness, tanggalFollowing(b.startPeriod), tanggalFollowing(b.endPeriod)]}
              onLihat={bukaBerkas ? lihatTreaty : undefined}
              bisaLihat={(b) => b.policyNo !== ''}
            />
          </div>

          <div hidden={tab !== 'Life'}>
            <DaftarCariPolis<BarisPolisLife>
              judul={V.daftar}
              placeholder={V.placeholder.polis}
              kolom={V.kolom.polisLife}
              cari={(_o, f, h) => cariPolisLife(f, h)}
              kunci={(b, i) => `${b.caseId}-${b.policyNo}-${b.endorsNo}-${i}`}
              sel={(b) => [b.caseId, b.policyNo, b.typePolicy, b.sourceOfBusiness, b.cedingCoName, b.qq]}
            />
            <DaftarCariPolis<BarisOfferLife>
              judul={V.daftarOffer}
              placeholder={V.placeholder.offer}
              kolom={V.kolom.offerLife}
              cari={(_o, f, h) => cariOfferLife(f, h)}
              kunci={(b, i) => `${b.caseId}-${b.noOffer}-${i}`}
              sel={(b) => [b.caseId, b.noOffer, b.typePolicy, b.businessName, b.cedingCoName, b.qq, b.lastStatus]}
            />
          </div>
        </div>
      </section>
    </div>
  )
}
