// Isi tab polis NB Fac In hanya-baca Loss Record (FIRE) dan Correspondence - HANYA akun IT (`PolisNB.tsx`; keputusan
// work owner 11-10-2026 "Ikut NB: hanya akun IT"). Susunan disalin dari `modul/nbfacin/frontend/components/
// {TabLossRecord,TabCorrespondence}.tsx` baca-saja (pola, bukan impor): premi berlabel kiri, grid kerugian, kotak Total
// Loss per Class Of Business, kotak Loss Ratio 1 / 3 / 5 Year; grid email ▸ membuka pengirim, penerima, subjek, pesan.
// Data `polisNBIT.ts`.

import { useEffect, useRef } from 'react'

import { angkaNB, tanggalNB } from '../polisNB'
import {
  KORESPONDENSI_NB as K,
  LOSS_RECORD_NB as L,
  TEKS_KORESPONDENSI_NB as TK,
  TEKS_LOSS_RECORD_NB as TL,
  bersihkanHtmlNB,
  daftarNamaPenerimaNB,
  tanggalKorespondensiNB,
  totalCobNB,
  type KorespondensiNB,
  type LossRecordNB,
} from '../polisNBIT'
import { GridNB, TeksNB } from './PolisNBDasar'

const uang = (k: string) => ({ label: k, angka: true })

/** `persenRasio` nbfacin: angka 4 desimal + "%", kosong tetap kosong. */
const persenRasio = (v: string) => (v.trim() === '' ? '' : `${angkaNB(v)}%`)

export function TabLossRecordNB({ lr }: { lr: LossRecordNB }) {
  const rasio: [string, string][] = [
    [lr.gross1, lr.net1],
    [lr.gross3, lr.net3],
    [lr.gross5, lr.net5],
  ]
  return (
    <div className="claimfacin__nb-lr">
      <section className="claimfacin__nb-lr-premi">
        <TeksNB kiri angka label={L.premiPar} nilai={angkaNB(lr.premiPar)} />
        <TeksNB kiri angka label={L.premiMb} nilai={angkaNB(lr.premiMb)} />
        <TeksNB kiri angka label={L.premiPl} nilai={angkaNB(lr.premiPl)} />
        <TeksNB kiri angka label={L.premiOther} nilai={angkaNB(lr.premiOther)} />
      </section>
      <GridNB
        kosong={TL.kosong}
        kolom={[
          L.uwYear,
          L.classOfBusiness,
          L.dateOfLoss,
          L.causeOfLoss,
          L.currency,
          uang(L.lossAmount),
          uang(L.lossAmountIdr),
        ]}
        baris={lr.baris.map((b) => [
          b.uwYear,
          b.classOfBusiness,
          tanggalNB(b.dateOfLoss),
          b.causeOfLoss,
          b.currency,
          angkaNB(b.lossAmount),
          angkaNB(b.lossAmountIdr),
        ])}
      />
      <section className="claimfacin__nb-lr-total">
        <GridNB
          kosong={TL.kosong}
          kolom={[L.classOfBusiness, L.currency, uang(L.totalLoss), uang(L.totalLossIdr)]}
          baris={totalCobNB(lr.baris).map((t) => [
            t.classOfBusiness,
            t.currency,
            angkaNB(t.lossAmount),
            angkaNB(t.lossAmountIdr),
          ])}
        />
      </section>
      <section className="claimfacin__nb-lr-rasio">
        {L.rasio.map((judul, k) => (
          <div key={judul} className="claimfacin__nb-lr-rasio-baris">
            <span className="claimfacin__nb-lr-rasio-judul">{judul}</span>
            <dl>
              <dt>{L.gross}</dt>
              <dd>{persenRasio(rasio[k]?.[0] ?? '')}</dd>
            </dl>
            <dl>
              <dt>{L.net}</dt>
              <dd>{persenRasio(rasio[k]?.[1] ?? '')}</dd>
            </dl>
          </div>
        ))}
      </section>
    </div>
  )
}

/** Pesan email: HANYA keluaran `bersihkanHtmlNB` yang ditaruh ke innerHTML (pola `EditorTeksKaya` nbfacin baca-saja). */
function PesanNB({ html }: { html: string }) {
  const kotak = useRef<HTMLDivElement>(null)
  useEffect(() => {
    if (kotak.current) kotak.current.innerHTML = bersihkanHtmlNB(html)
  }, [html])
  return <div ref={kotak} className="claimfacin__nb-kor-pesan" role="document" aria-label={K.message} />
}

function EmailNB({ k }: { k: KorespondensiNB }) {
  return (
    <div className="claimfacin__nb-kor">
      <div className="claimfacin__nb-kolom">
        <TeksNB kiri label={K.sender} nilai={k.sender || '—'} />
        <TeksNB kiri label={K.senderMail} nilai={k.senderMail || '—'} />
      </div>
      <section>
        <h6 className="claimfacin__nb-kor-judul">{K.penerima}</h6>
        {k.penerima.length === 0 && <p className="claimfacin__nb-mentah">{TK.kosong}</p>}
        {k.penerima.map((p, i) => (
          <div key={i} className="claimfacin__nb-kor-penerima">
            <span>{p.recipient || '—'}</span>
            <span>{p.recipientMail}</span>
          </div>
        ))}
      </section>
      <TeksNB kiri label={K.subject} nilai={k.subject} />
      <div>
        <h6 className="claimfacin__nb-kor-judul">{K.message}</h6>
        <PesanNB html={k.message} />
      </div>
    </div>
  )
}

export function TabCorrespondenceNB({ daftar }: { daftar: readonly KorespondensiNB[] }) {
  return (
    <GridNB
      kosong={TK.kosong}
      kolom={K.kolom}
      baris={daftar.map((k) => [
        tanggalKorespondensiNB(k.tanggal),
        k.sender,
        k.subject,
        daftarNamaPenerimaNB(k.penerima),
        <span
          key="status"
          className={'claimfacin__nb-kor-status' + (k.terkirim ? ' claimfacin__nb-kor-status--terkirim' : '')}
        >
          {k.terkirim ? K.terkirim : K.belumTerkirim}
        </span>,
      ])}
      rinci={(i) => (daftar[i] ? <EmailNB k={daftar[i]} /> : null)}
    />
  )
}
